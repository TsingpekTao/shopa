package points

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	httpv1 "github.com/TsingpekTao/shopa/points-svc/api/points/v1"
	"github.com/TsingpekTao/shopa/points-svc/internal/dao"
	"github.com/TsingpekTao/shopa/points-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/points-svc/internal/model/entity"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

const (
	// consumptionDetail* 描述积分消费明细在锁定、确认、取消之间的流转状态。
	consumptionDetailLocked    = "LOCKED"
	consumptionDetailConfirmed = "CONFIRMED"
	consumptionDetailCanceled  = "CANCELED"

	// refundAction* 用于标识退款返分与赠分冲回两类幂等动作。
	refundActionReturn  = "RETURN"
	refundActionReverse = "REVERSE"
)

// reservationPayload 保存锁分时的规则快照和子单分摊，用于后续确认、取消按快照执行。
type reservationPayload struct {
	RuleSnapshotJSON string                       `json:"rule_snapshot_json"`
	SubAllocations   []httpv1.PointsSubAllocation `json:"sub_allocations"`
}

// bucketLockPiece 表示本次锁分从某个 bucket 中预占了多少积分。
type bucketLockPiece struct {
	BucketNo string
	Points   uint64
}

// detailAllocation 把子单分摊和 bucket 锁定结果进一步展开成消费明细行。
type detailAllocation struct {
	SubOrderNo     string
	BucketNo       string
	Points         uint64
	CashAmountCent int64
}

// refundActionPayload 是退款动作的幂等回放载荷，用于重复请求直接返回历史结果。
type refundActionPayload struct {
	EffectivePoints       uint64 `json:"effective_points"`
	CashOffsetAmount      uint64 `json:"cash_offset_amount"`
	AccountAvailableAfter int64  `json:"account_available_after"`
	AccountDebtAfter      uint64 `json:"account_debt_after"`
	GraceBucketGranted    uint64 `json:"grace_bucket_granted"`
	GraceBucketNo         string `json:"grace_bucket_no,omitempty"`
}

// previewOrderHTTP 负责预览订单可抵扣积分、抵现金额和子单分摊结果。
func (s *sPoints) previewOrderHTTP(ctx context.Context, req *httpv1.PreviewOrderReq) (*httpv1.PreviewOrderRes, error) {
	if req == nil || req.OrderDraft == nil || req.OrderDraft.UserID == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_draft.user_id is required")
	}
	if !req.UsePoints {
		// 调用方明确不使用积分时，直接返回空预览结果。
		return &httpv1.PreviewOrderRes{UsePoints: false}, nil
	}

	// 先准备账户和规则快照，保证后续抵扣计算基于固定规则。
	account, err := s.ensureAccount(ctx, req.OrderDraft.UserID)
	if err != nil {
		return nil, err
	}
	_, snapshot, snapshotJSON, digest, err := s.loadRule(ctx)
	if err != nil {
		return nil, err
	}

	payableAmount := int64(req.OrderDraft.PayableAmount)
	if payableAmount < snapshot.MinOrderAmountCent {
		// 未达到起用门槛时仍返回规则快照，方便调用方展示原因。
		return &httpv1.PreviewOrderRes{
			UsePoints:              false,
			PointsRuleSnapshotJSON: snapshotJSON,
			PointsRuleDigest:       digest,
		}, nil
	}

	available := uint64(0)
	if account.AvailableBalance > 0 {
		available = uint64(account.AvailableBalance)
	}
	// 请求积分为 0 或超出可用余额时，统一收敛到真实可用积分。
	requested := req.IntentPoints
	if requested == 0 || requested > available {
		requested = available
	}

	maxCashByRule := uint64(payableAmount) * uint64(snapshot.MaxDeductionRateBps) / 10000
	if maxCashByRule == 0 || maxCashByRule > uint64(payableAmount) {
		maxCashByRule = uint64(payableAmount)
	}

	// 规则里未配置换算比例时，退化为 1 分 = 1 分现金单位，避免除 0。
	pointsPerCent := snapshot.DeductPointsPerCent
	if pointsPerCent == 0 {
		pointsPerCent = 1
	}

	maxPointsByRule := maxCashByRule * pointsPerCent
	pointsUsed := requested
	if pointsUsed > maxPointsByRule {
		pointsUsed = maxPointsByRule
	}

	cashDiscount := pointsUsed / pointsPerCent
	// 子单分摊遵循前 N-1 比例、最后一项兜底，保证总分摊和原值一致。
	allocations := allocateByShop(req.OrderDraft.SubOrders, pointsUsed, cashDiscount)

	return &httpv1.PreviewOrderRes{
		UsePoints:              pointsUsed > 0,
		PointsUsed:             pointsUsed,
		PointsDiscountAmount:   cashDiscount,
		PointsRuleSnapshotJSON: snapshotJSON,
		PointsRuleDigest:       digest,
		SubAllocations:         allocations,
	}, nil
}

// lockOrderHTTP 将预览结果落成正式 reservation，并同步锁定 bucket 资产。
func (s *sPoints) lockOrderHTTP(ctx context.Context, req *httpv1.LockOrderReq) (*httpv1.LockOrderRes, error) {
	if req == nil || req.OrderDraft == nil || req.OrderDraft.UserID == 0 || strings.TrimSpace(req.IdempotencyKey) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_draft.user_id and idempotency_key are required")
	}

	// 锁分前先走一次统一预览，确保 preview 与 lock 使用同一套规则和算法。
	preview, err := s.previewOrderHTTP(ctx, &httpv1.PreviewOrderReq{
		OrderDraft:    req.OrderDraft,
		UsePoints:     req.UsePoints,
		IntentPoints:  req.IntentPoints,
		RequestSource: req.RequestSource,
	})
	if err != nil {
		return nil, err
	}
	if req.ExpectedPointsCashAmount > 0 && preview.PointsDiscountAmount != req.ExpectedPointsCashAmount {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_points_cash_amount mismatch")
	}
	if strings.TrimSpace(req.ExpectedRuleSnapshotDigest) != "" && strings.TrimSpace(preview.PointsRuleDigest) != strings.TrimSpace(req.ExpectedRuleSnapshotDigest) {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "expected_rule_snapshot_digest mismatch")
	}
	if !preview.UsePoints || preview.PointsUsed == 0 {
		// 最终没有可锁积分时，直接回传规则快照，不创建 reservation。
		return &httpv1.LockOrderRes{
			PointsRuleSnapshotJSON:   preview.PointsRuleSnapshotJSON,
			PointsRuleSnapshotDigest: preview.PointsRuleDigest,
		}, nil
	}

	// 先查是否已经针对同一订单和幂等键成功创建 reservation。
	var existing entity.PointsReservation
	_ = dao.PointsReservation.Ctx(ctx).
		Where(dao.PointsReservation.Columns().OrderNo, req.OrderDraft.OrderNo).
		Where(dao.PointsReservation.Columns().IdempotencyKey, strings.TrimSpace(req.IdempotencyKey)).
		Scan(&existing)
	if existing.ReservationNo != "" {
		// 命中幂等时直接回放历史 reservation，避免重复锁 bucket。
		return &httpv1.LockOrderRes{
			ReservationNo:            existing.ReservationNo,
			PointsUsed:               existing.LockedPoints,
			PointsDiscountAmount:     uint64(existing.LockedCashAmountCent),
			PointsRuleSnapshotJSON:   preview.PointsRuleSnapshotJSON,
			PointsRuleSnapshotDigest: existing.DeductionDigest,
			SubAllocations:           preview.SubAllocations,
		}, nil
	}

	payloadBytes, err := json.Marshal(reservationPayload{
		RuleSnapshotJSON: preview.PointsRuleSnapshotJSON,
		SubAllocations:   preview.SubAllocations,
	})
	if err != nil {
		return nil, gerror.Wrap(err, "marshal reservation payload failed")
	}

	reservationNo := generateBizNo("PR")
	err = dao.PointsAccount.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 锁分事务中先校验账户状态，再去 bucket 侧做真实预占。
		account, innerErr := s.getAccountTx(ctx, tx, req.OrderDraft.UserID)
		if innerErr != nil {
			return innerErr
		}
		if account.StatusCode != accountStatusActive {
			return gerror.NewCode(gcode.CodeNotAuthorized, "points account is not active")
		}
		if account.AvailableBalance < int64(preview.PointsUsed) {
			return gerror.NewCode(gcode.CodeInvalidParameter, "insufficient points balance")
		}

		// 先按 FIFO 从 bucket 上真实锁定，避免确认阶段出现“空头支票”。
		lockPieces, innerErr := s.lockBucketsTx(ctx, tx, req.OrderDraft.UserID, preview.PointsUsed)
		if innerErr != nil {
			return innerErr
		}
		// 再把子单分摊拆到 bucket 维度，形成可回溯的消费明细。
		detailRows, innerErr := buildDetailAllocations(preview.SubAllocations, lockPieces)
		if innerErr != nil {
			return innerErr
		}

		availableAfter := account.AvailableBalance - int64(preview.PointsUsed)
		frozenAfter := account.FrozenBalance + int64(preview.PointsUsed)
		if _, innerErr = tx.Model(dao.PointsAccount.Table()).
			Where(dao.PointsAccount.Columns().UserId, req.OrderDraft.UserID).
			Data(do.PointsAccount{
				AvailableBalance: availableAfter,
				FrozenBalance:    frozenAfter,
			}).Update(); innerErr != nil {
			return gerror.Wrap(innerErr, "update account for lock failed")
		}

		if _, innerErr = tx.Model(dao.PointsReservation.Table()).Data(do.PointsReservation{
			ReservationNo:         reservationNo,
			UserId:                req.OrderDraft.UserID,
			OrderNo:               req.OrderDraft.OrderNo,
			ReservationStatusCode: reservationStatusLocked,
			RequestedPoints:       req.IntentPoints,
			LockedPoints:          preview.PointsUsed,
			LockedCashAmountCent:  int64(preview.PointsDiscountAmount),
			DeductionDigest:       preview.PointsRuleDigest,
			RuleSnapshotJson:      string(payloadBytes),
			IdempotencyKey:        strings.TrimSpace(req.IdempotencyKey),
			ExpireAt:              gtime.NewFromTime(time.Now().Add(30 * time.Minute)),
		}).Insert(); innerErr != nil {
			return gerror.Wrap(innerErr, "insert points_reservation failed")
		}

		// 每一行消费明细都带上子单和 bucket，用于后续退款返分与撤销赠分。
		for _, row := range detailRows {
			if _, innerErr = tx.Model(dao.PointsConsumptionDetail.Table()).Data(do.PointsConsumptionDetail{
				DetailNo:         generateBizNo("PCD"),
				UserId:           req.OrderDraft.UserID,
				ReservationNo:    reservationNo,
				OrderNo:          req.OrderDraft.OrderNo,
				RefundNo:         "",
				SubOrderNo:       row.SubOrderNo,
				BucketNo:         row.BucketNo,
				DetailStatusCode: consumptionDetailLocked,
				ConsumedPoints:   row.Points,
				ReturnedPoints:   0,
				CashAmountCent:   row.CashAmountCent,
			}).Insert(); innerErr != nil {
				return gerror.Wrap(innerErr, "insert locked points_consumption_detail failed")
			}
		}

		_, innerErr = s.insertLedgerTx(ctx, tx, &ledgerInput{
			UserID:         req.OrderDraft.UserID,
			EntryTypeCode:  "LOCK",
			BizType:        "ORDER_LOCK",
			BizNo:          req.OrderDraft.OrderNo,
			ReservationNo:  reservationNo,
			PointsDelta:    0,
			AvailableAfter: availableAfter,
			FrozenAfter:    frozenAfter,
			DebtAfter:      debtFromAvailable(availableAfter),
			CashAmountCent: int64(preview.PointsDiscountAmount),
			Remark:         "lock points for order with bucket reservation",
		})
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	return &httpv1.LockOrderRes{
		ReservationNo:            reservationNo,
		PointsUsed:               preview.PointsUsed,
		PointsDiscountAmount:     preview.PointsDiscountAmount,
		PointsRuleSnapshotJSON:   preview.PointsRuleSnapshotJSON,
		PointsRuleSnapshotDigest: preview.PointsRuleDigest,
		SubAllocations:           preview.SubAllocations,
	}, nil
}

// cancelOrderHTTP 释放尚未确认的 reservation，把冻结积分和 bucket 资产原路退回。
func (s *sPoints) cancelOrderHTTP(ctx context.Context, req *httpv1.CancelOrderReq) (*httpv1.SimpleAckRes, error) {
	reservation, err := s.findReservation(ctx, req.ReservationNo, req.OrderNo)
	if err != nil {
		return nil, err
	}
	if reservation == nil {
		// 不存在 reservation 时按幂等成功处理，方便补偿重试。
		return &httpv1.SimpleAckRes{Success: true}, nil
	}
	if reservation.ReservationStatusCode == reservationStatusCanceled {
		return &httpv1.SimpleAckRes{Success: true}, nil
	}
	if reservation.ReservationStatusCode == reservationStatusConfirmed {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "reservation already confirmed")
	}

	err = dao.PointsAccount.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		account, innerErr := s.getAccountTx(ctx, tx, reservation.UserId)
		if innerErr != nil {
			return innerErr
		}

		// 找出仍处于 LOCKED 状态的消费明细，只释放尚未确认的那部分资产。
		details, innerErr := s.listReservationDetailsTx(ctx, tx, reservation.ReservationNo, consumptionDetailLocked)
		if innerErr != nil {
			return innerErr
		}
		if innerErr = s.releaseLockedBucketsTx(ctx, tx, details); innerErr != nil {
			return innerErr
		}

		availableAfter := account.AvailableBalance + int64(reservation.LockedPoints)
		frozenAfter := account.FrozenBalance - int64(reservation.LockedPoints)
		if frozenAfter < 0 {
			frozenAfter = 0
		}
		if _, innerErr = tx.Model(dao.PointsAccount.Table()).
			Where(dao.PointsAccount.Columns().UserId, reservation.UserId).
			Data(do.PointsAccount{
				AvailableBalance: availableAfter,
				FrozenBalance:    frozenAfter,
			}).Update(); innerErr != nil {
			return gerror.Wrap(innerErr, "restore account by cancel failed")
		}

		if _, innerErr = tx.Model(dao.PointsConsumptionDetail.Table()).
			Where(dao.PointsConsumptionDetail.Columns().ReservationNo, reservation.ReservationNo).
			Where(dao.PointsConsumptionDetail.Columns().DetailStatusCode, consumptionDetailLocked).
			Data(do.PointsConsumptionDetail{DetailStatusCode: consumptionDetailCanceled}).
			Update(); innerErr != nil {
			return gerror.Wrap(innerErr, "mark locked details canceled failed")
		}

		if _, innerErr = tx.Model(dao.PointsReservation.Table()).
			Where(dao.PointsReservation.Columns().ReservationNo, reservation.ReservationNo).
			Data(do.PointsReservation{
				ReservationStatusCode: reservationStatusCanceled,
				CancelReasonCode:      strings.TrimSpace(req.ReasonCode),
				CanceledAt:            gtime.Now(),
			}).Update(); innerErr != nil {
			return gerror.Wrap(innerErr, "update reservation canceled failed")
		}

		_, innerErr = s.insertLedgerTx(ctx, tx, &ledgerInput{
			UserID:         reservation.UserId,
			EntryTypeCode:  "CANCEL",
			BizType:        "ORDER_CANCEL",
			BizNo:          reservation.OrderNo,
			ReservationNo:  reservation.ReservationNo,
			PointsDelta:    0,
			AvailableAfter: availableAfter,
			FrozenAfter:    frozenAfter,
			DebtAfter:      debtFromAvailable(availableAfter),
			CashAmountCent: reservation.LockedCashAmountCent,
			Remark:         "cancel order points reservation",
		})
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	return &httpv1.SimpleAckRes{Success: true}, nil
}

// confirmOrderHTTP 把 reservation 从 LOCKED 推进到 CONFIRMED，并累计已消费积分。
func (s *sPoints) confirmOrderHTTP(ctx context.Context, req *httpv1.ConfirmOrderReq) (*httpv1.SimpleAckRes, error) {
	reservation, err := s.findReservation(ctx, req.ReservationNo, req.OrderNo)
	if err != nil {
		return nil, err
	}
	if reservation == nil {
		// reservation 缺失时视作已经处理过，避免支付回调重放报错。
		return &httpv1.SimpleAckRes{Success: true}, nil
	}
	if reservation.ReservationStatusCode == reservationStatusConfirmed {
		return &httpv1.SimpleAckRes{Success: true}, nil
	}
	if reservation.ReservationStatusCode == reservationStatusCanceled {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "reservation already canceled")
	}

	err = dao.PointsAccount.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		account, innerErr := s.getAccountTx(ctx, tx, reservation.UserId)
		if innerErr != nil {
			return innerErr
		}

		// 确认时不再重新扫可用 bucket，而是消费锁单阶段已经预占的资产。
		details, innerErr := s.listReservationDetailsTx(ctx, tx, reservation.ReservationNo, consumptionDetailLocked)
		if innerErr != nil {
			return innerErr
		}
		if innerErr = s.confirmLockedBucketsTx(ctx, tx, details); innerErr != nil {
			return innerErr
		}

		frozenAfter := account.FrozenBalance - int64(reservation.LockedPoints)
		if frozenAfter < 0 {
			frozenAfter = 0
		}
		if _, innerErr = tx.Model(dao.PointsAccount.Table()).
			Where(dao.PointsAccount.Columns().UserId, reservation.UserId).
			Data(do.PointsAccount{
				FrozenBalance:   frozenAfter,
				TotalUsedPoints: gdb.Raw(fmt.Sprintf("%s + %d", dao.PointsAccount.Columns().TotalUsedPoints, reservation.LockedPoints)),
			}).Update(); innerErr != nil {
			return gerror.Wrap(innerErr, "confirm account points failed")
		}

		if _, innerErr = tx.Model(dao.PointsConsumptionDetail.Table()).
			Where(dao.PointsConsumptionDetail.Columns().ReservationNo, reservation.ReservationNo).
			Where(dao.PointsConsumptionDetail.Columns().DetailStatusCode, consumptionDetailLocked).
			Data(do.PointsConsumptionDetail{DetailStatusCode: consumptionDetailConfirmed}).
			Update(); innerErr != nil {
			return gerror.Wrap(innerErr, "mark locked details confirmed failed")
		}

		if _, innerErr = tx.Model(dao.PointsReservation.Table()).
			Where(dao.PointsReservation.Columns().ReservationNo, reservation.ReservationNo).
			Data(do.PointsReservation{
				ReservationStatusCode: reservationStatusConfirmed,
				ConfirmedAt:           gtime.Now(),
			}).Update(); innerErr != nil {
			return gerror.Wrap(innerErr, "update reservation confirmed failed")
		}

		_, innerErr = s.insertLedgerTx(ctx, tx, &ledgerInput{
			UserID:         reservation.UserId,
			EntryTypeCode:  "CONFIRM",
			BizType:        "ORDER_CONFIRM",
			BizNo:          reservation.OrderNo,
			ReservationNo:  reservation.ReservationNo,
			PointsDelta:    0,
			AvailableAfter: account.AvailableBalance,
			FrozenAfter:    frozenAfter,
			DebtAfter:      debtFromAvailable(account.AvailableBalance),
			CashAmountCent: reservation.LockedCashAmountCent,
			Remark:         "confirm locked points without re-consuming buckets",
		})
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	return &httpv1.SimpleAckRes{Success: true}, nil
}

// grantOrderHTTP 在订单完成后发放奖励积分，并写入 grant detail 与月度 bucket。
func (s *sPoints) grantOrderHTTP(ctx context.Context, req *httpv1.GrantOrderReq) (*httpv1.SimpleAckRes, error) {
	if req == nil || req.UserID == 0 || strings.TrimSpace(req.OrderNo) == "" || strings.TrimSpace(req.IdempotencyKey) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no user_id idempotency_key are required")
	}
	if req.PaidAmount == 0 {
		// 实付为 0 时不发奖励积分，但按幂等成功返回。
		return &httpv1.SimpleAckRes{Success: true}, nil
	}

	// 同一订单只允许生成一笔订单级赠分明细。
	var existing entity.PointsGrantDetail
	_ = dao.PointsGrantDetail.Ctx(ctx).
		Where(dao.PointsGrantDetail.Columns().OrderNo, req.OrderNo).
		Where(dao.PointsGrantDetail.Columns().SubOrderNo, "__ORDER__").
		Scan(&existing)
	if existing.GrantDetailNo != "" {
		return &httpv1.SimpleAckRes{Success: true}, nil
	}

	_, snapshot, snapshotJSON, _, err := s.loadRule(ctx)
	if err != nil {
		return nil, err
	}

	grantedPoints := req.PaidAmount * snapshot.GrantPointsPerCent
	err = dao.PointsAccount.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		account, innerErr := s.getAccountTx(ctx, tx, req.UserID)
		if innerErr != nil {
			return innerErr
		}

		availableAfter := account.AvailableBalance + int64(grantedPoints)
		if _, innerErr = tx.Model(dao.PointsAccount.Table()).
			Where(dao.PointsAccount.Columns().UserId, req.UserID).
			Data(do.PointsAccount{
				AvailableBalance:  availableAfter,
				TotalEarnedPoints: gdb.Raw(fmt.Sprintf("%s + %d", dao.PointsAccount.Columns().TotalEarnedPoints, grantedPoints)),
			}).Update(); innerErr != nil {
			return gerror.Wrap(innerErr, "grant account points failed")
		}

		// 月聚合 bucket 负责降低碎片数量，grant detail 负责售后冲回进度追踪。
		bucketNo, innerErr := s.upsertMonthlyGrantBucketTx(ctx, tx, req.UserID, grantedPoints)
		if innerErr != nil {
			return innerErr
		}

		if _, innerErr = tx.Model(dao.PointsGrantDetail.Table()).Data(do.PointsGrantDetail{
			GrantDetailNo:    generateBizNo("PGD"),
			UserId:           req.UserID,
			OrderNo:          req.OrderNo,
			SubOrderNo:       "__ORDER__",
			ShopNo:           "",
			GrantedPoints:    grantedPoints,
			ReversedPoints:   0,
			RuleCode:         snapshot.RuleCode,
			RuleSnapshotJson: snapshotJSON,
			GrantStatusCode:  "GRANTED",
		}).Insert(); innerErr != nil {
			return gerror.Wrap(innerErr, "insert points_grant_detail failed")
		}

		_, innerErr = s.insertLedgerTx(ctx, tx, &ledgerInput{
			UserID:          req.UserID,
			EntryTypeCode:   "GRANT",
			BizType:         "ORDER_GRANT",
			BizNo:           req.OrderNo,
			RelatedBucketNo: bucketNo,
			PointsDelta:     int64(grantedPoints),
			AvailableAfter:  availableAfter,
			FrozenAfter:     account.FrozenBalance,
			DebtAfter:       debtFromAvailable(availableAfter),
			CashAmountCent:  int64(req.PaidAmount),
			Remark:          "grant completed order points",
		})
		return innerErr
	})
	if err != nil {
		return nil, err
	}

	return &httpv1.SimpleAckRes{Success: true}, nil
}

// returnRefundHTTP 在售后退款时返还用户曾经消耗的积分，并优先用于抵债。
func (s *sPoints) returnRefundHTTP(ctx context.Context, req *httpv1.ReturnRefundReq) (*httpv1.ReturnRefundRes, error) {
	if req == nil || req.UserID == 0 || strings.TrimSpace(req.RefundNo) == "" || strings.TrimSpace(req.OrderNo) == "" || strings.TrimSpace(req.IdempotencyKey) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "refund_no order_no user_id idempotency_key are required")
	}

	res := &httpv1.ReturnRefundRes{Success: true}
	lookupKey := strings.TrimSpace(req.SubOrderNo)
	if lookupKey == "" {
		// 子单号缺失时退化到店铺维度匹配，兼容按店铺拆单的场景。
		lookupKey = strings.TrimSpace(req.ShopNo)
	}

	err := dao.PointsAccount.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		replayed, innerErr := s.findRefundActionReplayTx(ctx, tx, refundActionReturn, req.RefundNo, lookupKey, req.ShopNo, req.IdempotencyKey, req.UserID)
		if innerErr != nil {
			return innerErr
		}
		if replayed != nil {
			// 命中退款动作幂等回放时，直接返回上次的返分结果。
			res.PointsReturnAmount = replayed.EffectivePoints
			return nil
		}

		account, innerErr := s.getAccountTx(ctx, tx, req.UserID)
		if innerErr != nil {
			return innerErr
		}

		// 退款返分只允许返还尚未返过的确认消费积分。
		details, innerErr := s.listRefundableDetailsTx(ctx, tx, req.OrderNo, lookupKey)
		if innerErr != nil {
			return innerErr
		}
		actualReturn := capReturnPoints(details, req.PointsToReturn)
		if actualReturn == 0 {
			// 没有可返积分时也要记录 refund action，避免重复请求继续打数据库。
			return s.insertRefundActionTx(ctx, tx, refundActionReturn, req, 0, 0, account.AvailableBalance, debtFromAvailable(account.AvailableBalance), 0, "")
		}

		if innerErr = s.applyReturnToDetailsTx(ctx, tx, req.RefundNo, details, actualReturn); innerErr != nil {
			return innerErr
		}

		availableAfter := account.AvailableBalance + int64(actualReturn)
		gracePoints := actualReturn
		existingDebt := debtFromAvailable(account.AvailableBalance)
		if existingDebt > 0 {
			// 账户存在债务时，返还积分必须先填平债务，剩余部分才能入退款宽限桶。
			if existingDebt >= actualReturn {
				gracePoints = 0
			} else {
				gracePoints = actualReturn - existingDebt
			}
		}

		totalUsedExpr := fmt.Sprintf("CASE WHEN %s >= %d THEN %s - %d ELSE 0 END", dao.PointsAccount.Columns().TotalUsedPoints, actualReturn, dao.PointsAccount.Columns().TotalUsedPoints, actualReturn)
		if _, innerErr = tx.Model(dao.PointsAccount.Table()).
			Where(dao.PointsAccount.Columns().UserId, req.UserID).
			Data(do.PointsAccount{
				AvailableBalance: availableAfter,
				TotalUsedPoints:  gdb.Raw(totalUsedExpr),
			}).Update(); innerErr != nil {
			return gerror.Wrap(innerErr, "update account return points failed")
		}

		graceBucketNo := ""
		if gracePoints > 0 {
			graceBucketNo, innerErr = s.upsertRefundGraceBucketTx(ctx, tx, req.UserID, req.RefundNo, gracePoints)
			if innerErr != nil {
				return innerErr
			}
		}

		_, innerErr = s.insertLedgerTx(ctx, tx, &ledgerInput{
			UserID:          req.UserID,
			EntryTypeCode:   "RETURN",
			BizType:         "REFUND_RETURN",
			BizNo:           req.RefundNo,
			RelatedBucketNo: graceBucketNo,
			PointsDelta:     int64(actualReturn),
			AvailableAfter:  availableAfter,
			FrozenAfter:     account.FrozenBalance,
			DebtAfter:       debtFromAvailable(availableAfter),
			CashAmountCent:  int64(req.CashAmountCent),
			Remark:          "return used points by refund with debt-first policy",
		})
		if innerErr != nil {
			return innerErr
		}

		if innerErr = s.insertRefundActionTx(ctx, tx, refundActionReturn, req, actualReturn, 0, availableAfter, debtFromAvailable(availableAfter), gracePoints, graceBucketNo); innerErr != nil {
			return innerErr
		}

		res.PointsReturnAmount = actualReturn
		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// reverseRefundHTTP 在售后退款时撤销已发奖励积分，不足部分转为债务并计算现金抵扣。
func (s *sPoints) reverseRefundHTTP(ctx context.Context, req *httpv1.ReverseRefundReq) (*httpv1.ReverseRefundRes, error) {
	if req == nil || req.UserID == 0 || strings.TrimSpace(req.RefundNo) == "" || strings.TrimSpace(req.OrderNo) == "" || strings.TrimSpace(req.IdempotencyKey) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "refund_no order_no user_id idempotency_key are required")
	}

	res := &httpv1.ReverseRefundRes{Success: true}
	lookupKey := strings.TrimSpace(req.SubOrderNo)
	if lookupKey == "" {
		// 未传子单号时用店铺维度兜底匹配赠分明细。
		lookupKey = strings.TrimSpace(req.ShopNo)
	}

	err := dao.PointsAccount.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		replayed, innerErr := s.findRefundActionReplayTx(ctx, tx, refundActionReverse, req.RefundNo, lookupKey, req.ShopNo, req.IdempotencyKey, req.UserID)
		if innerErr != nil {
			return innerErr
		}
		if replayed != nil {
			// 命中幂等回放时直接使用历史撤销结果，不重复扣账。
			res.PointsReverseAmount = replayed.EffectivePoints
			res.PointsCashOffsetAmount = replayed.CashOffsetAmount
			res.AccountDebtAfter = replayed.AccountDebtAfter
			return nil
		}

		account, innerErr := s.getAccountTx(ctx, tx, req.UserID)
		if innerErr != nil {
			return innerErr
		}

		var grant entity.PointsGrantDetail
		model := tx.Model(dao.PointsGrantDetail.Table()).
			Where(dao.PointsGrantDetail.Columns().OrderNo, req.OrderNo)
		if lookupKey != "" {
			model = model.Where(dao.PointsGrantDetail.Columns().SubOrderNo, lookupKey)
		} else {
			model = model.Where(dao.PointsGrantDetail.Columns().SubOrderNo, "__ORDER__")
		}
		if err := model.Scan(&grant); err != nil {
			return gerror.Wrap(err, "query grant detail failed")
		}
		if grant.Id == 0 && lookupKey != "" {
			// 子单级赠分不存在时，兼容回落到整单级赠分明细。
			if err := tx.Model(dao.PointsGrantDetail.Table()).
				Where(dao.PointsGrantDetail.Columns().OrderNo, req.OrderNo).
				Where(dao.PointsGrantDetail.Columns().SubOrderNo, "__ORDER__").
				Scan(&grant); err != nil {
				return gerror.Wrap(err, "query order level grant detail failed")
			}
		}

		remainingReverse := req.PointsToReverse
		if grant.Id > 0 {
			grantRemain := uint64(0)
			if grant.GrantedPoints > grant.ReversedPoints {
				grantRemain = grant.GrantedPoints - grant.ReversedPoints
			}
			if remainingReverse > grantRemain {
				remainingReverse = grantRemain
			}
		}
		if remainingReverse == 0 {
			// 当前赠分已经全部冲回过时，记录空动作后直接返回。
			return s.insertRefundActionTx(ctx, tx, refundActionReverse, req, 0, 0, account.AvailableBalance, debtFromAvailable(account.AvailableBalance), 0, "")
		}

		bucketDeduct := remainingReverse
		if account.AvailableBalance <= 0 {
			// 当前已无可用资产时，不再从 bucket 扣减，剩余部分直接形成债务。
			bucketDeduct = 0
		} else if uint64(account.AvailableBalance) < bucketDeduct {
			bucketDeduct = uint64(account.AvailableBalance)
		}
		if bucketDeduct > 0 {
			if innerErr = s.consumeAvailableBucketsForReverseTx(ctx, tx, req.UserID, bucketDeduct); innerErr != nil {
				return innerErr
			}
		}

		beforeDebt := debtFromAvailable(account.AvailableBalance)
		availableAfter := account.AvailableBalance - int64(remainingReverse)
		afterDebtRaw := debtFromAvailable(availableAfter)
		debtIncrease := uint64(0)
		if afterDebtRaw > beforeDebt {
			debtIncrease = afterDebtRaw - beforeDebt
		}

		cashOffset := debtIncrease
		if cashOffset > req.ApprovedRefundAmount {
			// 现金抵扣不能超过本次批准退款金额，超出部分继续保留为债务。
			cashOffset = req.ApprovedRefundAmount
		}
		availableAfter += int64(cashOffset)
		debtAfter := debtFromAvailable(availableAfter)

		totalEarnedExpr := fmt.Sprintf("CASE WHEN %s >= %d THEN %s - %d ELSE 0 END", dao.PointsAccount.Columns().TotalEarnedPoints, remainingReverse, dao.PointsAccount.Columns().TotalEarnedPoints, remainingReverse)
		if _, innerErr = tx.Model(dao.PointsAccount.Table()).
			Where(dao.PointsAccount.Columns().UserId, req.UserID).
			Data(do.PointsAccount{
				AvailableBalance:  availableAfter,
				TotalEarnedPoints: gdb.Raw(totalEarnedExpr),
			}).Update(); innerErr != nil {
			return gerror.Wrap(innerErr, "reverse grant update account failed")
		}

		if grant.Id > 0 {
			statusCode := "PARTIAL_REVERSED"
			if grant.ReversedPoints+remainingReverse >= grant.GrantedPoints {
				statusCode = "REVERSED"
			}
			if _, innerErr = tx.Model(dao.PointsGrantDetail.Table()).
				Where(dao.PointsGrantDetail.Columns().Id, grant.Id).
				Data(do.PointsGrantDetail{
					ReversedPoints:  grant.ReversedPoints + remainingReverse,
					GrantStatusCode: statusCode,
				}).Update(); innerErr != nil {
				return gerror.Wrap(innerErr, "update grant reverse progress failed")
			}
		}

		_, innerErr = s.insertLedgerTx(ctx, tx, &ledgerInput{
			UserID:         req.UserID,
			EntryTypeCode:  "REVERSE",
			BizType:        "REFUND_REVERSE",
			BizNo:          req.RefundNo,
			PointsDelta:    -int64(remainingReverse) + int64(cashOffset),
			AvailableAfter: availableAfter,
			FrozenAfter:    account.FrozenBalance,
			DebtAfter:      debtAfter,
			CashAmountCent: int64(cashOffset),
			Remark:         "reverse granted points by refund and keep debt if cash refund is insufficient",
		})
		if innerErr != nil {
			return innerErr
		}

		if innerErr = s.insertRefundActionTx(ctx, tx, refundActionReverse, req, remainingReverse, cashOffset, availableAfter, debtAfter, 0, ""); innerErr != nil {
			return innerErr
		}

		res.PointsReverseAmount = remainingReverse
		res.PointsCashOffsetAmount = cashOffset
		res.AccountDebtAfter = debtAfter
		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// findReservation 按 reservation_no 或 order_no 查询最近一次锁分记录。
func (s *sPoints) findReservation(ctx context.Context, reservationNo, orderNo string) (*entity.PointsReservation, error) {
	var row entity.PointsReservation
	model := dao.PointsReservation.Ctx(ctx)
	if strings.TrimSpace(reservationNo) != "" {
		model = model.Where(dao.PointsReservation.Columns().ReservationNo, strings.TrimSpace(reservationNo))
	} else {
		model = model.Where(dao.PointsReservation.Columns().OrderNo, strings.TrimSpace(orderNo)).OrderDesc(dao.PointsReservation.Columns().CreatedAt)
	}
	if err := model.Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "query points_reservation failed")
	}
	if row.ReservationNo == "" {
		return nil, nil
	}
	return &row, nil
}

// allocateByShop 按支付金额比例分摊积分和抵现金额，最后一项负责兜底余数。
func allocateByShop(subs []httpv1.PointsOrderDraftSub, totalPoints, totalCash uint64) []httpv1.PointsSubAllocation {
	if len(subs) == 0 || totalPoints == 0 {
		return nil
	}

	totalPayable := uint64(0)
	for _, sub := range subs {
		totalPayable += sub.PayableAmount
	}
	if totalPayable == 0 {
		return nil
	}

	allocations := make([]httpv1.PointsSubAllocation, 0, len(subs))
	allocatedPoints := uint64(0)
	allocatedCash := uint64(0)
	for i, sub := range subs {
		points := uint64(0)
		cash := uint64(0)
		if i == len(subs)-1 {
			points = totalPoints - allocatedPoints
			cash = totalCash - allocatedCash
		} else {
			points = totalPoints * sub.PayableAmount / totalPayable
			cash = totalCash * sub.PayableAmount / totalPayable
			allocatedPoints += points
			allocatedCash += cash
		}
		allocations = append(allocations, httpv1.PointsSubAllocation{
			ShopNo:               sub.ShopNo,
			PointsUsed:           points,
			PointsDiscountAmount: cash,
		})
	}
	return allocations
}

// buildDetailAllocations 把子单分摊结果继续切到 bucket 维度，生成消费明细草稿。
func buildDetailAllocations(allocations []httpv1.PointsSubAllocation, pieces []bucketLockPiece) ([]detailAllocation, error) {
	if len(allocations) == 0 || len(pieces) == 0 {
		return nil, nil
	}

	result := make([]detailAllocation, 0, len(pieces))
	pieceIndex := 0
	pieceRemain := pieces[0].Points
	for _, allocation := range allocations {
		if allocation.PointsUsed == 0 {
			continue
		}

		pointsRemain := allocation.PointsUsed
		cashRemain := int64(allocation.PointsDiscountAmount)
		for pointsRemain > 0 {
			if pieceIndex >= len(pieces) {
				return nil, gerror.NewCode(gcode.CodeInternalError, "bucket lock pieces not enough for allocations")
			}

			take := pieceRemain
			if take > pointsRemain {
				take = pointsRemain
			}

			cashTake := int64(0)
			if take == pointsRemain {
				cashTake = cashRemain
			} else if allocation.PointsUsed > 0 {
				cashTake = int64(uint64(allocation.PointsDiscountAmount) * take / allocation.PointsUsed)
				if cashTake > cashRemain {
					cashTake = cashRemain
				}
			}

			result = append(result, detailAllocation{
				SubOrderNo:     allocation.ShopNo,
				BucketNo:       pieces[pieceIndex].BucketNo,
				Points:         take,
				CashAmountCent: cashTake,
			})

			pointsRemain -= take
			cashRemain -= cashTake
			pieceRemain -= take
			if pieceRemain == 0 {
				pieceIndex++
				if pieceIndex < len(pieces) {
					pieceRemain = pieces[pieceIndex].Points
				}
			}
		}
	}
	return result, nil
}

// lockBucketsTx 按 FIFO 锁定 bucket 中的可用积分，并返回每个 bucket 的锁定片段。
func (s *sPoints) lockBucketsTx(ctx context.Context, tx gdb.TX, userID uint64, points uint64) ([]bucketLockPiece, error) {
	if points == 0 {
		return nil, nil
	}

	var buckets []*entity.PointsExpireBucket
	if err := tx.Model(dao.PointsExpireBucket.Table()).
		Where(dao.PointsExpireBucket.Columns().UserId, userID).
		Where(dao.PointsExpireBucket.Columns().BucketStatus, bucketStatusActive).
		WhereGT(dao.PointsExpireBucket.Columns().RemainingPoints, 0).
		OrderAsc(dao.PointsExpireBucket.Columns().ExpireAt).
		OrderAsc(dao.PointsExpireBucket.Columns().CreatedAt).
		Scan(&buckets); err != nil {
		return nil, gerror.Wrap(err, "query buckets for lock failed")
	}

	remaining := points
	pieces := make([]bucketLockPiece, 0)
	for _, bucket := range buckets {
		if remaining == 0 {
			break
		}
		lockPoints := bucket.RemainingPoints
		if lockPoints > remaining {
			lockPoints = remaining
		}
		if lockPoints == 0 {
			continue
		}

		remainingAfter := bucket.RemainingPoints - lockPoints
		lockedAfter := bucket.LockedPoints + lockPoints
		if _, err := tx.Model(dao.PointsExpireBucket.Table()).
			Where(dao.PointsExpireBucket.Columns().BucketNo, bucket.BucketNo).
			Data(do.PointsExpireBucket{
				RemainingPoints: remainingAfter,
				LockedPoints:    lockedAfter,
				BucketStatus:    deriveBucketStatus(remainingAfter, lockedAfter),
			}).Update(); err != nil {
			return nil, gerror.Wrap(err, "lock bucket points failed")
		}

		pieces = append(pieces, bucketLockPiece{BucketNo: bucket.BucketNo, Points: lockPoints})
		remaining -= lockPoints
	}
	if remaining > 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "insufficient bucket assets for lock")
	}
	_ = ctx
	return pieces, nil
}

// listReservationDetailsTx 查询某个 reservation 下的消费明细，可选按状态过滤。
func (s *sPoints) listReservationDetailsTx(ctx context.Context, tx gdb.TX, reservationNo, statusCode string) ([]*entity.PointsConsumptionDetail, error) {
	if strings.TrimSpace(reservationNo) == "" {
		return nil, nil
	}

	model := tx.Model(dao.PointsConsumptionDetail.Table()).
		Where(dao.PointsConsumptionDetail.Columns().ReservationNo, reservationNo)
	if strings.TrimSpace(statusCode) != "" {
		model = model.Where(dao.PointsConsumptionDetail.Columns().DetailStatusCode, statusCode)
	}

	var details []*entity.PointsConsumptionDetail
	if err := model.OrderAsc(dao.PointsConsumptionDetail.Columns().Id).Scan(&details); err != nil {
		return nil, gerror.Wrap(err, "query reservation details failed")
	}
	_ = ctx
	return details, nil
}

// releaseLockedBucketsTx 把取消订单释放出来的锁定积分回滚回原 bucket。
func (s *sPoints) releaseLockedBucketsTx(ctx context.Context, tx gdb.TX, details []*entity.PointsConsumptionDetail) error {
	if len(details) == 0 {
		return nil
	}

	aggregated := aggregateDetailPointsByBucket(details)
	for bucketNo, points := range aggregated {
		var bucket entity.PointsExpireBucket
		if err := tx.Model(dao.PointsExpireBucket.Table()).
			Where(dao.PointsExpireBucket.Columns().BucketNo, bucketNo).
			Scan(&bucket); err != nil {
			return gerror.Wrap(err, "query bucket for cancel failed")
		}

		lockedAfter := bucket.LockedPoints
		if lockedAfter >= points {
			lockedAfter -= points
		} else {
			lockedAfter = 0
		}
		remainingAfter := bucket.RemainingPoints + points
		if _, err := tx.Model(dao.PointsExpireBucket.Table()).
			Where(dao.PointsExpireBucket.Columns().BucketNo, bucketNo).
			Data(do.PointsExpireBucket{
				RemainingPoints: remainingAfter,
				LockedPoints:    lockedAfter,
				BucketStatus:    deriveBucketStatus(remainingAfter, lockedAfter),
			}).Update(); err != nil {
			return gerror.Wrap(err, "release locked bucket points failed")
		}
	}
	_ = ctx
	return nil
}

// confirmLockedBucketsTx 将锁定 bucket 中的积分转为已使用积分，完成真实消费确认。
func (s *sPoints) confirmLockedBucketsTx(ctx context.Context, tx gdb.TX, details []*entity.PointsConsumptionDetail) error {
	if len(details) == 0 {
		return gerror.NewCode(gcode.CodeInternalError, "reservation locked details missing")
	}

	aggregated := aggregateDetailPointsByBucket(details)
	for bucketNo, points := range aggregated {
		var bucket entity.PointsExpireBucket
		if err := tx.Model(dao.PointsExpireBucket.Table()).
			Where(dao.PointsExpireBucket.Columns().BucketNo, bucketNo).
			Scan(&bucket); err != nil {
			return gerror.Wrap(err, "query bucket for confirm failed")
		}

		lockedAfter := bucket.LockedPoints
		if lockedAfter >= points {
			lockedAfter -= points
		} else {
			lockedAfter = 0
		}
		usedAfter := bucket.UsedPoints + points
		if _, err := tx.Model(dao.PointsExpireBucket.Table()).
			Where(dao.PointsExpireBucket.Columns().BucketNo, bucketNo).
			Data(do.PointsExpireBucket{
				LockedPoints: lockedAfter,
				UsedPoints:   usedAfter,
				BucketStatus: deriveBucketStatus(bucket.RemainingPoints, lockedAfter),
			}).Update(); err != nil {
			return gerror.Wrap(err, "confirm bucket points failed")
		}
	}
	_ = ctx
	return nil
}

// aggregateDetailPointsByBucket 把明细按 bucket 聚合，便于批量回滚或确认。
func aggregateDetailPointsByBucket(details []*entity.PointsConsumptionDetail) map[string]uint64 {
	result := make(map[string]uint64, len(details))
	for _, detail := range details {
		if detail == nil || detail.BucketNo == "" || detail.ConsumedPoints == 0 {
			continue
		}
		result[detail.BucketNo] += detail.ConsumedPoints
	}
	return result
}

// deriveBucketStatus 根据剩余和锁定积分判断 bucket 是否仍处于活跃状态。
func deriveBucketStatus(remainingPoints, lockedPoints uint64) string {
	if remainingPoints == 0 && lockedPoints == 0 {
		return bucketStatusClosed
	}
	return bucketStatusActive
}

// listRefundableDetailsTx 查询指定订单下可用于退款返分的已确认消费明细。
func (s *sPoints) listRefundableDetailsTx(ctx context.Context, tx gdb.TX, orderNo, lookupKey string) ([]*entity.PointsConsumptionDetail, error) {
	model := tx.Model(dao.PointsConsumptionDetail.Table()).
		Where(dao.PointsConsumptionDetail.Columns().OrderNo, orderNo).
		Where(dao.PointsConsumptionDetail.Columns().DetailStatusCode, consumptionDetailConfirmed)
	if lookupKey != "" {
		model = model.Where(dao.PointsConsumptionDetail.Columns().SubOrderNo, lookupKey)
	}

	var details []*entity.PointsConsumptionDetail
	if err := model.OrderAsc(dao.PointsConsumptionDetail.Columns().Id).Scan(&details); err != nil {
		return nil, gerror.Wrap(err, "query refundable details failed")
	}
	_ = ctx
	return details, nil
}

// capReturnPoints 计算本次退款最多还能返还多少积分，防止多次部分退款超返。
func capReturnPoints(details []*entity.PointsConsumptionDetail, requested uint64) uint64 {
	remaining := uint64(0)
	for _, detail := range details {
		if detail == nil {
			continue
		}
		if detail.ConsumedPoints > detail.ReturnedPoints {
			remaining += detail.ConsumedPoints - detail.ReturnedPoints
		}
	}
	if requested == 0 || requested > remaining {
		return remaining
	}
	return requested
}

// applyReturnToDetailsTx 把本次返还积分分摊回消费明细，累计 returned_points 进度。
func (s *sPoints) applyReturnToDetailsTx(ctx context.Context, tx gdb.TX, refundNo string, details []*entity.PointsConsumptionDetail, points uint64) error {
	remaining := points
	for _, detail := range details {
		if detail == nil || remaining == 0 {
			continue
		}
		refundable := uint64(0)
		if detail.ConsumedPoints > detail.ReturnedPoints {
			refundable = detail.ConsumedPoints - detail.ReturnedPoints
		}
		if refundable == 0 {
			continue
		}

		returnPoints := refundable
		if returnPoints > remaining {
			returnPoints = remaining
		}
		if _, err := tx.Model(dao.PointsConsumptionDetail.Table()).
			Where(dao.PointsConsumptionDetail.Columns().Id, detail.Id).
			Data(do.PointsConsumptionDetail{
				RefundNo:       refundNo,
				ReturnedPoints: detail.ReturnedPoints + returnPoints,
			}).Update(); err != nil {
			return gerror.Wrap(err, "update returned points failed")
		}
		remaining -= returnPoints
	}
	_ = ctx
	return nil
}

// consumeAvailableBucketsForReverseTx 从当前可用 bucket 中扣除需要冲回的已发奖励积分。
func (s *sPoints) consumeAvailableBucketsForReverseTx(ctx context.Context, tx gdb.TX, userID, points uint64) error {
	if points == 0 {
		return nil
	}

	var buckets []*entity.PointsExpireBucket
	if err := tx.Model(dao.PointsExpireBucket.Table()).
		Where(dao.PointsExpireBucket.Columns().UserId, userID).
		Where(dao.PointsExpireBucket.Columns().BucketStatus, bucketStatusActive).
		WhereGT(dao.PointsExpireBucket.Columns().RemainingPoints, 0).
		OrderAsc(dao.PointsExpireBucket.Columns().ExpireAt).
		OrderAsc(dao.PointsExpireBucket.Columns().CreatedAt).
		Scan(&buckets); err != nil {
		return gerror.Wrap(err, "query buckets for reverse failed")
	}

	remaining := points
	for _, bucket := range buckets {
		if remaining == 0 {
			break
		}
		deduct := bucket.RemainingPoints
		if deduct > remaining {
			deduct = remaining
		}
		if deduct == 0 {
			continue
		}

		remainingAfter := bucket.RemainingPoints - deduct
		totalAfter := bucket.TotalPoints
		if totalAfter >= deduct {
			totalAfter -= deduct
		} else {
			totalAfter = 0
		}
		if _, err := tx.Model(dao.PointsExpireBucket.Table()).
			Where(dao.PointsExpireBucket.Columns().BucketNo, bucket.BucketNo).
			Data(do.PointsExpireBucket{
				TotalPoints:     totalAfter,
				RemainingPoints: remainingAfter,
				BucketStatus:    deriveBucketStatus(remainingAfter, bucket.LockedPoints),
			}).Update(); err != nil {
			return gerror.Wrap(err, "deduct bucket points for reverse failed")
		}
		remaining -= deduct
	}
	if remaining > 0 {
		return gerror.NewCode(gcode.CodeInvalidParameter, "insufficient available bucket assets for reverse")
	}
	_ = ctx
	return nil
}

// findRefundActionReplayTx 按退款号和幂等键查询历史退款动作，用于重试回放。
func (s *sPoints) findRefundActionReplayTx(ctx context.Context, tx gdb.TX, actionType, refundNo, subOrderNo, shopNo, idempotencyKey string, userID uint64) (*refundActionPayload, error) {
	var action entity.PointsRefundAction
	model := tx.Model(dao.PointsRefundAction.Table()).
		Where(dao.PointsRefundAction.Columns().ActionType, actionType).
		Where(dao.PointsRefundAction.Columns().RefundNo, strings.TrimSpace(refundNo)).
		Where(dao.PointsRefundAction.Columns().IdempotencyKey, strings.TrimSpace(idempotencyKey)).
		Where(dao.PointsRefundAction.Columns().UserId, userID)
	if strings.TrimSpace(subOrderNo) != "" {
		model = model.Where(dao.PointsRefundAction.Columns().SubOrderNo, strings.TrimSpace(subOrderNo))
	}
	if strings.TrimSpace(shopNo) != "" {
		model = model.Where(dao.PointsRefundAction.Columns().ShopNo, strings.TrimSpace(shopNo))
	}
	if err := model.Scan(&action); err != nil {
		return nil, gerror.Wrap(err, "query refund action replay failed")
	}
	if action.Id == 0 {
		return nil, nil
	}

	payload := &refundActionPayload{}
	if strings.TrimSpace(action.ResultPayloadJson) != "" {
		if err := json.Unmarshal([]byte(action.ResultPayloadJson), payload); err != nil {
			return nil, gerror.Wrap(err, "unmarshal refund action replay failed")
		}
	}
	_ = ctx
	return payload, nil
}

// insertRefundActionTx 落退款动作幂等记录，保存本次返分或冲回的最终结果快照。
func (s *sPoints) insertRefundActionTx(ctx context.Context, tx gdb.TX, actionType string, req any, effectivePoints, cashOffset uint64, accountAvailableAfter int64, accountDebtAfter uint64, gracePoints uint64, graceBucketNo string) error {
	refundNo := ""
	orderNo := ""
	subOrderNo := ""
	shopNo := ""
	userID := uint64(0)
	idempotencyKey := ""
	cashAmountCent := int64(0)
	requestedPoints := effectivePoints

	switch typed := req.(type) {
	case *httpv1.ReturnRefundReq:
		refundNo = typed.RefundNo
		orderNo = typed.OrderNo
		subOrderNo = typed.SubOrderNo
		shopNo = typed.ShopNo
		userID = typed.UserID
		idempotencyKey = typed.IdempotencyKey
		cashAmountCent = int64(typed.CashAmountCent)
		requestedPoints = typed.PointsToReturn
	case *httpv1.ReverseRefundReq:
		refundNo = typed.RefundNo
		orderNo = typed.OrderNo
		subOrderNo = typed.SubOrderNo
		shopNo = typed.ShopNo
		userID = typed.UserID
		idempotencyKey = typed.IdempotencyKey
		cashAmountCent = int64(typed.ApprovedRefundAmount)
		requestedPoints = typed.PointsToReverse
	default:
		return gerror.NewCode(gcode.CodeInternalError, "unsupported refund action request type")
	}

	payloadJSON, err := json.Marshal(refundActionPayload{
		EffectivePoints:       effectivePoints,
		CashOffsetAmount:      cashOffset,
		AccountAvailableAfter: accountAvailableAfter,
		AccountDebtAfter:      accountDebtAfter,
		GraceBucketGranted:    gracePoints,
		GraceBucketNo:         graceBucketNo,
	})
	if err != nil {
		return gerror.Wrap(err, "marshal refund action payload failed")
	}

	_, err = tx.Model(dao.PointsRefundAction.Table()).Data(do.PointsRefundAction{
		ActionNo:          generateBizNo("PRA"),
		ActionType:        actionType,
		RefundNo:          refundNo,
		OrderNo:           orderNo,
		SubOrderNo:        subOrderNo,
		ShopNo:            shopNo,
		UserId:            userID,
		IdempotencyKey:    idempotencyKey,
		RequestedPoints:   requestedPoints,
		EffectivePoints:   effectivePoints,
		CashAmountCent:    cashAmountCent,
		ActionStatus:      "SUCCESS",
		ResultPayloadJson: string(payloadJSON),
	}).Insert()
	if err != nil && !strings.Contains(strings.ToLower(err.Error()), "duplicate") {
		return gerror.Wrap(err, "insert refund action failed")
	}
	_ = ctx
	return nil
}

// upsertMonthlyGrantBucketTx 以月为粒度聚合奖励积分 bucket，减少 bucket 碎片数量。
func (s *sPoints) upsertMonthlyGrantBucketTx(ctx context.Context, tx gdb.TX, userID, points uint64) (string, error) {
	bucketPeriod := time.Now().Format("200601")
	bucketNo := fmt.Sprintf("PB%s%06d", bucketPeriod, userID%1000000)

	var existing entity.PointsExpireBucket
	_ = tx.Model(dao.PointsExpireBucket.Table()).
		Where(dao.PointsExpireBucket.Columns().BucketNo, bucketNo).
		Scan(&existing)
	if existing.BucketNo == "" {
		expireAt := gtime.NewFromTime(time.Date(time.Now().Year()+1, 12, 31, 23, 59, 59, 0, time.Local))
		if _, err := tx.Model(dao.PointsExpireBucket.Table()).Data(do.PointsExpireBucket{
			BucketNo:        bucketNo,
			UserId:          userID,
			SourceType:      bucketSourceGrant,
			SourceNo:        bucketPeriod,
			SourceVersion:   1,
			BucketPeriod:    bucketPeriod,
			TotalPoints:     points,
			RemainingPoints: points,
			LockedPoints:    0,
			UsedPoints:      0,
			ExpiredPoints:   0,
			ReturnedPoints:  0,
			ExpireAt:        expireAt,
			BucketStatus:    bucketStatusActive,
		}).Insert(); err != nil {
			return "", gerror.Wrap(err, "insert monthly grant bucket failed")
		}
		return bucketNo, nil
	}

	if _, err := tx.Model(dao.PointsExpireBucket.Table()).
		Where(dao.PointsExpireBucket.Columns().BucketNo, bucketNo).
		Data(do.PointsExpireBucket{
			TotalPoints:     existing.TotalPoints + points,
			RemainingPoints: existing.RemainingPoints + points,
			BucketStatus:    bucketStatusActive,
		}).Update(); err != nil {
		return "", gerror.Wrap(err, "update monthly grant bucket failed")
	}
	_ = ctx
	return bucketNo, nil
}

// upsertRefundGraceBucketTx 为退款返分创建或累计宽限 bucket，给用户保留再次使用窗口。
func (s *sPoints) upsertRefundGraceBucketTx(ctx context.Context, tx gdb.TX, userID uint64, refundNo string, points uint64) (string, error) {
	if points == 0 {
		return "", nil
	}

	bucketNo := fmt.Sprintf("PBG%s", strings.TrimSpace(refundNo))
	var existing entity.PointsExpireBucket
	_ = tx.Model(dao.PointsExpireBucket.Table()).
		Where(dao.PointsExpireBucket.Columns().BucketNo, bucketNo).
		Scan(&existing)
	if existing.BucketNo == "" {
		expireAt := gtime.NewFromTime(time.Now().Add(30 * 24 * time.Hour))
		if _, err := tx.Model(dao.PointsExpireBucket.Table()).Data(do.PointsExpireBucket{
			BucketNo:        bucketNo,
			UserId:          userID,
			SourceType:      bucketSourceRefundGrace,
			SourceNo:        refundNo,
			SourceVersion:   1,
			BucketPeriod:    time.Now().Format("200601"),
			TotalPoints:     points,
			RemainingPoints: points,
			LockedPoints:    0,
			UsedPoints:      0,
			ExpiredPoints:   0,
			ReturnedPoints:  points,
			ExpireAt:        expireAt,
			BucketStatus:    bucketStatusActive,
		}).Insert(); err != nil {
			return "", gerror.Wrap(err, "insert refund grace bucket failed")
		}
		return bucketNo, nil
	}

	if _, err := tx.Model(dao.PointsExpireBucket.Table()).
		Where(dao.PointsExpireBucket.Columns().BucketNo, bucketNo).
		Data(do.PointsExpireBucket{
			TotalPoints:     existing.TotalPoints + points,
			RemainingPoints: existing.RemainingPoints + points,
			ReturnedPoints:  existing.ReturnedPoints + points,
			BucketStatus:    bucketStatusActive,
		}).Update(); err != nil {
		return "", gerror.Wrap(err, "update refund grace bucket failed")
	}
	_ = ctx
	return bucketNo, nil
}
