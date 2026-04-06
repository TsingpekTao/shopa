package order

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/order-svc/api/v1"
	"github.com/TsingpekTao/shopa/order-svc/internal/dao"
	"github.com/TsingpekTao/shopa/order-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/order-svc/internal/model/entity"
	"github.com/TsingpekTao/shopa/order-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	idempotencyStatusProcessing = 0
	idempotencyStatusSuccess    = 1
	idempotencyStatusFailed     = 2

	idempotencyActionCreateFromCart = "CREATE_ORDER_FROM_CART"
	idempotencyActionCreateBuyNow   = "CREATE_ORDER_BUY_NOW"
	idempotencyActionRequestPay     = "REQUEST_PAY"
	idempotencyActionCancelOrder    = "CANCEL_MY_ORDER"
	idempotencyActionPayCallback    = "PAY_CALLBACK"

	defaultPageSize = 20
	maxPageSize     = 100
)

type sOrder struct {
}

// New 创建订单领域逻辑实例。
func New() *sOrder { return &sOrder{} }

// init 在包初始化阶段注册订单服务实现。
func init() {
	service.RegisterOrder(New())
}

// CreateOrderFromCart 基于购物车结算快照创建订单。
func (s *sOrder) CreateOrderFromCart(ctx context.Context, req *v1.CreateOrderFromCartReq) (*v1.CreateOrderFromCartRes, error) {
	if req == nil || strings.TrimSpace(req.GetCheckoutToken()) == "" || strings.TrimSpace(req.GetIdempotencyKey()) == "" || req.GetAddressId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "checkout_token/address_id/idempotency_key are required")
	}
	// 读取用户身份，校验必需的上下文信息。
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// 幂等记录：确保相同 key 的下单请求只会执行一次，避免重复扣库存或积分。
	hit, row, err := s.getOrCreateIdempotency(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateFromCart)
	if err != nil {
		return nil, err
	}
	if hit {
		agg, replay, err := s.replayCreateOrderWithFlag(ctx, row)
		if err != nil {
			return nil, err
		}
		return &v1.CreateOrderFromCartRes{Order: agg, IdempotentReplay: replay}, nil
	}
	snap, err := s.consumeCheckoutSnapshot(ctx, req.GetCheckoutToken())
	if err != nil {
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateFromCart, "TOKEN_CONSUME_FAILED")
		return nil, err
	}
	if snap.UserID != 0 && snap.UserID != userID {
		// 购物车快照指定了用户但与当前访问用户不一致，直接拒绝。
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateFromCart, "TOKEN_USER_MISMATCH")
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "checkout token user mismatch")
	}
	if req.GetExpectedSnapshotDigest() != "" && snap.SnapshotDigest != "" && req.GetExpectedSnapshotDigest() != snap.SnapshotDigest {
		// 前端传入的快照摘要与 Redis 中的实际快照不一致，防止脏数据提交。
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateFromCart, "SNAPSHOT_DIGEST_MISMATCH")
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "snapshot digest mismatch")
	}
	lines := snapshotToOrderLines(snap)
	if len(lines) == 0 {
		// 无行项目说明快照已经被消费或为空，直接标记失败。
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateFromCart, "EMPTY_SNAPSHOT")
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "snapshot has no items")
	}
	// 组装创建订单所需的关键字段。
	orderNo := generateBizNo("ORD")
	createInput := &createOrderInput{
		OrderNo:        orderNo,
		UserID:         userID,
		AddressID:      req.GetAddressId(),
		BuyerRemark:    req.GetBuyerRemark(),
		GoodsAmount:    snap.GoodsAmount,
		FreightAmount:  snap.FreightAmount,
		DiscountAmount: 0,
		PayableAmount:  snap.PayableAmount,
		Lines:          lines,
	}
	previewResp, err := s.previewOrderPoints(ctx, createInput, req.GetUsePoints(), req.GetIntentPoints(), req.GetExpectedPointsCashAmount(), req.GetPointsRuleSnapshotDigest(), req.GetSubmitSourceCode())
	if err != nil {
		// 预览失败直接回滚幂等状态，避免重复调用锁分接口。
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateFromCart, "POINTS_PREVIEW_FAILED")
		return nil, err
	}
	lockResp, err := s.lockPointsForOrder(ctx, createInput, req.GetUsePoints(), req.GetIntentPoints(), req.GetExpectedPointsCashAmount(), req.GetPointsRuleSnapshotDigest(), req.GetIdempotencyKey(), req.GetSubmitSourceCode())
	if err != nil {
		// 锁分失败也记录幂等失败，防止客户端重试导致重复锁定。
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateFromCart, "POINTS_LOCK_FAILED")
		return nil, err
	}
	// 将积分锁定的结果回填到订单创建输入，后续金额计算依赖。
	createInput.PointsReservationNo = lockResp.ReservationNo
	createInput.PointsUsed = lockResp.PointsUsed
	createInput.PointsDiscountAmount = lockResp.PointsDiscountAmount
	createInput.PointsRuleSnapshotJSON = firstNonEmpty(lockResp.PointsRuleSnapshotJSON, previewResp.PointsRuleSnapshotJSON)
	createInput.PointsRuleSnapshotDigest = firstNonEmpty(lockResp.PointsRuleSnapshotDigest, previewResp.PointsRuleDigest)
	createInput.PointsSubAllocations = pointsAllocationsToMap(lockResp.SubAllocations)
	if createInput.PointsDiscountAmount > createInput.PayableAmount {
		// 锁定的积分折现金额超过应付金额，属于异常数据，直接拒绝。
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateFromCart, "POINTS_AMOUNT_INVALID")
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "points discount amount exceeds payable amount")
	}
	createInput.PayableAmount -= createInput.PointsDiscountAmount

	// 预留库存，失败时会释放积分锁并标记幂等失败。
	reservationNo, err := s.reserveInventory(ctx, orderNo, userID, lines)
	if err != nil {
		if isNotImplementedErr(err) {
			g.Log().Warningf(ctx, "inventory reserve not implemented, skip reservation, order_no=%s", orderNo)
			reservationNo = ""
		} else {
			_ = s.compensateLockedPointsAfterCreateFailure(ctx, userID, orderNo, createInput.PointsReservationNo, err)
			_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateFromCart, "INVENTORY_RESERVE_FAILED")
			return nil, err
		}
	}
	createInput.ReservationNo = reservationNo
	// 真正落库订单主子单。失败时补偿库存与积分。
	agg, err := s.createOrderAggregate(ctx, createInput)
	if err != nil {
		s.cancelInventoryReservationBestEffort(ctx, reservationNo, orderNo, "ORDER_CREATE_FAILED")
		_ = s.compensateLockedPointsAfterCreateFailure(ctx, userID, orderNo, createInput.PointsReservationNo, err)
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateFromCart, "ORDER_CREATE_FAILED")
		return nil, err
	}
	if err = s.markIdempotencySuccess(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateFromCart, orderNo, nil); err != nil {
		return nil, err
	}
	return &v1.CreateOrderFromCartRes{Order: agg, IdempotentReplay: false}, nil
}

// CreateOrderBuyNow 基于立即购买参数直接创建订单。
func (s *sOrder) CreateOrderBuyNow(ctx context.Context, req *v1.CreateOrderBuyNowReq) (*v1.CreateOrderBuyNowRes, error) {
	if req == nil || len(req.GetItems()) == 0 || strings.TrimSpace(req.GetIdempotencyKey()) == "" || req.GetAddressId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "items/address_id/idempotency_key are required")
	}
	// 从上下文拿用户 ID，确保请求已经鉴权。
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// 幂等处理：同一个下单请求的重复调用直接重放历史结果。
	hit, row, err := s.getOrCreateIdempotency(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateBuyNow)
	if err != nil {
		return nil, err
	}
	if hit {
		agg, replay, err := s.replayCreateOrderWithFlag(ctx, row)
		if err != nil {
			return nil, err
		}
		return &v1.CreateOrderBuyNowRes{Order: agg, IdempotentReplay: replay}, nil
	}
	lines := make([]orderLine, 0, len(req.GetItems()))
	for _, item := range req.GetItems() {
		if strings.TrimSpace(item.GetSkuNo()) == "" || strings.TrimSpace(item.GetSpuNo()) == "" || strings.TrimSpace(item.GetShopNo()) == "" || item.GetQty() == 0 {
			_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateBuyNow, "INVALID_ITEM")
			return nil, gerror.NewCode(gcode.CodeInvalidParameter, "buy-now item invalid")
		}
		lines = append(lines, orderLine{SkuNo: item.GetSkuNo(), SpuNo: item.GetSpuNo(), ShopNo: item.GetShopNo(), Qty: item.GetQty()})
	}
	orderNo := generateBizNo("ORD")
	createInput := &createOrderInput{OrderNo: orderNo, UserID: userID, AddressID: req.GetAddressId(), BuyerRemark: req.GetBuyerRemark(), Lines: lines}
	previewResp, err := s.previewOrderPoints(ctx, createInput, req.GetUsePoints(), req.GetIntentPoints(), req.GetExpectedPointsCashAmount(), req.GetPointsRuleSnapshotDigest(), req.GetSubmitSourceCode())
	if err != nil {
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateBuyNow, "POINTS_PREVIEW_FAILED")
		return nil, err
	}
	lockResp, err := s.lockPointsForOrder(ctx, createInput, req.GetUsePoints(), req.GetIntentPoints(), req.GetExpectedPointsCashAmount(), req.GetPointsRuleSnapshotDigest(), req.GetIdempotencyKey(), req.GetSubmitSourceCode())
	if err != nil {
		// 锁分失败直接结束，避免后续库存已扣而积分未锁。
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateBuyNow, "POINTS_LOCK_FAILED")
		return nil, err
	}
	createInput.PointsReservationNo = lockResp.ReservationNo
	createInput.PointsUsed = lockResp.PointsUsed
	createInput.PointsDiscountAmount = lockResp.PointsDiscountAmount
	createInput.PointsRuleSnapshotJSON = firstNonEmpty(lockResp.PointsRuleSnapshotJSON, previewResp.PointsRuleSnapshotJSON)
	createInput.PointsRuleSnapshotDigest = firstNonEmpty(lockResp.PointsRuleSnapshotDigest, previewResp.PointsRuleDigest)
	createInput.PointsSubAllocations = pointsAllocationsToMap(lockResp.SubAllocations)

	reservationNo, err := s.reserveInventory(ctx, orderNo, userID, lines)
	if err != nil {
		if isNotImplementedErr(err) {
			g.Log().Warningf(ctx, "inventory reserve not implemented, skip reservation, order_no=%s", orderNo)
			reservationNo = ""
		} else {
			_ = s.compensateLockedPointsAfterCreateFailure(ctx, userID, orderNo, createInput.PointsReservationNo, err)
			_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateBuyNow, "INVENTORY_RESERVE_FAILED")
			return nil, err
		}
	}
	createInput.ReservationNo = reservationNo
	agg, err := s.createOrderAggregate(ctx, createInput)
	if err != nil {
		s.cancelInventoryReservationBestEffort(ctx, reservationNo, orderNo, "ORDER_CREATE_FAILED")
		_ = s.compensateLockedPointsAfterCreateFailure(ctx, userID, orderNo, createInput.PointsReservationNo, err)
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateBuyNow, "ORDER_CREATE_FAILED")
		return nil, err
	}
	if err = s.markIdempotencySuccess(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateBuyNow, orderNo, nil); err != nil {
		return nil, err
	}
	return &v1.CreateOrderBuyNowRes{Order: agg, IdempotentReplay: false}, nil
}

// RequestPay 为指定订单发起支付请求。
func (s *sOrder) RequestPay(ctx context.Context, req *v1.RequestPayReq) (*v1.RequestPayRes, error) {
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" || strings.TrimSpace(req.GetIdempotencyKey()) == "" || req.GetPayChannel() == v1.PayChannel_PAY_CHANNEL_UNSPECIFIED {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no/pay_channel/idempotency_key are required")
	}
	// 只有本人才能拉起支付，先解析 userID。
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// 幂等处理：重复支付请求直接重放历史结果，避免重复写入 payment 表。
	hit, row, err := s.getOrCreateIdempotency(ctx, userID, req.GetIdempotencyKey(), idempotencyActionRequestPay)
	if err != nil {
		return nil, err
	}
	if hit {
		return s.replayRequestPay(ctx, req.GetOrderNo(), row)
	}
	// 校验订单归属和状态，只有待支付订单才能创建支付单。
	mainRow, err := s.getOrderMainByNo(ctx, req.GetOrderNo())
	if err != nil {
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionRequestPay, "ORDER_NOT_FOUND")
		return nil, err
	}
	if mainRow.UserId != userID || v1.OrderStatus(mainRow.OrderStatus) != v1.OrderStatus_ORDER_STATUS_PENDING_PAY {
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionRequestPay, "ORDER_NOT_PAYABLE")
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order not payable")
	}
	payNo := generateBizNo("PAY")
	expireAt := mainRow.PayDeadlineAt
	if expireAt == nil {
		expireAt = gtime.NewFromTime(time.Now().Add(15 * time.Minute))
	}
	// 事务内同时插入支付记录与更新主单支付状态，保证状态一致。
	if err = dao.OrderMain.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err = tx.Model(dao.OrderPayment.Table()).Data(do.OrderPayment{OrderNo: req.GetOrderNo(), PayNo: payNo, PaymentEventId: generateBizNo("PEV"), PayChannel: uint(req.GetPayChannel()), PayStatusCode: "PAYING"}).Insert()
		if err != nil {
			return gerror.Wrap(err, "insert order_payment failed")
		}
		cols := dao.OrderMain.Columns()
		_, err = tx.Model(dao.OrderMain.Table()).Where(cols.OrderNo, req.GetOrderNo()).Where(cols.OrderStatus, uint(v1.OrderStatus_ORDER_STATUS_PENDING_PAY)).Data(do.OrderMain{PaymentStatus: uint(v1.PaymentStatus_PAYMENT_STATUS_PAYING), Version: gdb.Raw(cols.Version + " + 1")}).Update()
		return gerror.Wrap(err, "update order_main payment status failed")
	}); err != nil {
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionRequestPay, "REQUEST_PAY_FAILED")
		return nil, err
	}
	if err = s.markIdempotencySuccess(ctx, userID, req.GetIdempotencyKey(), idempotencyActionRequestPay, req.GetOrderNo(), map[string]any{"pay_no": payNo}); err != nil {
		return nil, err
	}
	return &v1.RequestPayRes{OrderNo: req.GetOrderNo(), PayNo: payNo, PaymentStatus: v1.PaymentStatus_PAYMENT_STATUS_PAYING, PayUrl: fmt.Sprintf("https://mock-pay.shopa.local/pay?pay_no=%s", payNo), PayPayloadJson: "{}", ExpireAt: toProtoTs(expireAt)}, nil
}

// CancelMyOrder 取消当前买家的订单。
func (s *sOrder) CancelMyOrder(ctx context.Context, req *v1.CancelMyOrderReq) (*v1.CancelMyOrderRes, error) {
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" || strings.TrimSpace(req.GetIdempotencyKey()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no/idempotency_key are required")
	}
	// 只有下单用户才能主动取消，先获取 userID。
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// 幂等：避免重复取消导致状态抖动。
	hit, row, err := s.getOrCreateIdempotency(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCancelOrder)
	if err != nil {
		return nil, err
	}
	if hit {
		return s.replayCancel(ctx, req.GetOrderNo(), row)
	}
	reason := req.GetReasonCode()
	if reason == v1.CancelReasonCode_CANCEL_REASON_CODE_UNSPECIFIED {
		// 默认设置为买家主动取消，方便运营统计。
		reason = v1.CancelReasonCode_CANCEL_REASON_CODE_BUYER_CANCEL
	}
	var (
		status              = v1.OrderStatus_ORDER_STATUS_UNSPECIFIED
		reservationNo       string
		pointsReservationNo string
	)
	// 事务内同时更新主单、子单状态，防止部分更新。
	err = dao.OrderMain.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cols := dao.OrderMain.Columns()
		var mainRow entity.OrderMain
		if err := tx.Model(dao.OrderMain.Table()).Where(cols.OrderNo, req.GetOrderNo()).Where(cols.UserId, userID).Scan(&mainRow); err != nil {
			return gerror.Wrap(err, "query order_main failed")
		}
		if mainRow.Id == 0 {
			return gerror.NewCode(gcode.CodeNotFound, "order not found")
		}
		reservationNo = mainRow.ReservationNo
		pointsReservationNo = mainRow.PointsReservationNo
		current := v1.OrderStatus(mainRow.OrderStatus)
		if current == v1.OrderStatus_ORDER_STATUS_CANCELED || current == v1.OrderStatus_ORDER_STATUS_CLOSED {
			status = current
			return nil
		}
		if current != v1.OrderStatus_ORDER_STATUS_PENDING_PAY {
			// 非待支付订单不允许前台取消，防止已付款的资金纠纷。
			return gerror.NewCode(gcode.CodeInvalidParameter, "only pending-pay order can be canceled")
		}
		result, err := tx.Model(dao.OrderMain.Table()).Where(cols.OrderNo, req.GetOrderNo()).Where(cols.UserId, userID).Where(cols.OrderStatus, uint(v1.OrderStatus_ORDER_STATUS_PENDING_PAY)).Data(do.OrderMain{OrderStatus: uint(v1.OrderStatus_ORDER_STATUS_CANCELED), CancelReasonCode: uint(reason), ClosedAt: gtime.Now(), Version: gdb.Raw(cols.Version + " + 1")}).Update()
		if err != nil {
			return gerror.Wrap(err, "cancel order_main failed")
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "order status changed")
		}
		subCols := dao.OrderSub.Columns()
		_, err = tx.Model(dao.OrderSub.Table()).Where(subCols.OrderNo, req.GetOrderNo()).Where(subCols.SubStatus, uint(v1.SubOrderStatus_SUB_ORDER_STATUS_PENDING_PAY)).Data(do.OrderSub{SubStatus: uint(v1.SubOrderStatus_SUB_ORDER_STATUS_CANCELED)}).Update()
		if err != nil {
			return gerror.Wrap(err, "cancel order_sub failed")
		}
		status = v1.OrderStatus_ORDER_STATUS_CANCELED
		return nil
	})
	if err != nil {
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCancelOrder, "CANCEL_ORDER_FAILED")
		return nil, err
	}
	if reservationNo != "" {
		s.cancelInventoryReservationBestEffort(ctx, reservationNo, req.GetOrderNo(), reason.String())
	}
	if pointsReservationNo != "" {
		s.compensateLockedPointsAfterOrderClose(ctx, userID, req.GetOrderNo(), pointsReservationNo, reason.String())
	}
	if err = s.markIdempotencySuccess(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCancelOrder, req.GetOrderNo(), nil); err != nil {
		return nil, err
	}
	return &v1.CancelMyOrderRes{OrderNo: req.GetOrderNo(), OrderStatus: status}, nil
}

// GetMyOrderDetail 查询当前买家的订单详情。
func (s *sOrder) GetMyOrderDetail(ctx context.Context, req *v1.GetMyOrderDetailReq) (*v1.GetMyOrderDetailRes, error) {
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no is required")
	}
	// 用户级别查询，必须保证订单属于本人。
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	mainRow, err := s.getOrderMainByNo(ctx, req.GetOrderNo())
	if err != nil {
		return nil, err
	}
	if mainRow.UserId != userID {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "order does not belong to current user")
	}
	agg, err := s.loadOrderAggregate(ctx, req.GetOrderNo())
	if err != nil {
		return nil, err
	}
	return &v1.GetMyOrderDetailRes{Order: agg}, nil
}

// ListMyOrders 分页查询当前买家的订单列表。
func (s *sOrder) ListMyOrders(ctx context.Context, req *v1.ListMyOrdersReq) (*v1.ListMyOrdersRes, error) {
	// 个人订单列表：按照创建时间倒序分页，使用游标避免翻页偏移。
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	pageSize := normalizePageSize(req.GetPageSize())
	cursorID, err := parseCursorID(req.GetNextCursor())
	if err != nil {
		return nil, err
	}
	cols := dao.OrderMain.Columns()
	model := dao.OrderMain.Ctx(ctx).Where(cols.UserId, userID).WhereNull(cols.DeletedAt)
	if cursorID > 0 {
		model = model.WhereLT(cols.Id, cursorID)
	}
	if len(req.GetStatuses()) > 0 {
		model = model.WhereIn(cols.OrderStatus, orderStatusToUintSlice(req.GetStatuses()))
	}
	if kw := strings.TrimSpace(req.GetKeyword()); kw != "" {
		model = model.WhereLike(cols.OrderNo, "%"+kw+"%")
	}
	var rows []*entity.OrderMain
	if err = model.OrderDesc(cols.Id).Limit(pageSize + 1).Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "list order_main failed")
	}
	hasMore := len(rows) > pageSize
	if hasMore {
		rows = rows[:pageSize]
	}
	orders := make([]*v1.OrderMain, 0, len(rows))
	for _, row := range rows {
		agg, loadErr := s.loadOrderAggregate(ctx, row.OrderNo)
		if loadErr != nil {
			return nil, loadErr
		}
		orders = append(orders, agg)
	}
	nextCursor := ""
	if hasMore && len(rows) > 0 {
		nextCursor = strconv.FormatUint(rows[len(rows)-1].Id, 10)
	}
	return &v1.ListMyOrdersRes{Orders: orders, NextCursor: nextCursor, HasMore: hasMore}, nil
}

// ListShopOrders 分页查询店铺维度的订单列表。
func (s *sOrder) ListShopOrders(ctx context.Context, req *v1.ListShopOrdersReq) (*v1.ListShopOrdersRes, error) {
	if req == nil || strings.TrimSpace(req.GetShopNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop_no is required")
	}
	// 店铺侧列表：以子单为粒度，便于商家查看本店订单。
	pageSize := normalizePageSize(req.GetPageSize())
	cursorID, err := parseCursorID(req.GetNextCursor())
	if err != nil {
		return nil, err
	}
	cols := dao.OrderSub.Columns()
	model := dao.OrderSub.Ctx(ctx).Where(cols.ShopNo, req.GetShopNo()).WhereNull(cols.DeletedAt)
	if cursorID > 0 {
		model = model.WhereLT(cols.Id, cursorID)
	}
	if len(req.GetStatuses()) > 0 {
		model = model.WhereIn(cols.SubStatus, subStatusToUintSlice(req.GetStatuses()))
	}
	if kw := strings.TrimSpace(req.GetKeyword()); kw != "" {
		model = model.WhereLike(cols.OrderNo, "%"+kw+"%")
	}
	var rows []*entity.OrderSub
	if err = model.OrderDesc(cols.Id).Limit(pageSize + 1).Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "list order_sub failed")
	}
	hasMore := len(rows) > pageSize
	if hasMore {
		rows = rows[:pageSize]
	}
	orders := make([]*v1.OrderSub, 0, len(rows))
	for _, row := range rows {
		items, itemErr := s.listOrderItemsBySubNo(ctx, row.SubOrderNo)
		if itemErr != nil {
			return nil, itemErr
		}
		orders = append(orders, toProtoSub(row, items))
	}
	nextCursor := ""
	if hasMore && len(rows) > 0 {
		nextCursor = strconv.FormatUint(rows[len(rows)-1].Id, 10)
	}
	return &v1.ListShopOrdersRes{Orders: orders, NextCursor: nextCursor, HasMore: hasMore}, nil
}

// GetShopOrderDetail 查询店铺侧的子单详情。
func (s *sOrder) GetShopOrderDetail(ctx context.Context, req *v1.GetShopOrderDetailReq) (*v1.GetShopOrderDetailRes, error) {
	if req == nil || strings.TrimSpace(req.GetSubOrderNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "sub_order_no is required")
	}
	// 先查子单确认存在，随后加载整单聚合数据，保证子单信息与主单保持一致。
	var sub entity.OrderSub
	err := dao.OrderSub.Ctx(ctx).Where(dao.OrderSub.Columns().SubOrderNo, req.GetSubOrderNo()).WhereNull(dao.OrderSub.Columns().DeletedAt).Scan(&sub)
	if err != nil {
		return nil, gerror.Wrap(err, "query sub order failed")
	}
	if sub.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "sub order not found")
	}
	agg, err := s.loadOrderAggregate(ctx, sub.OrderNo)
	if err != nil {
		return nil, err
	}
	var target *v1.OrderSub
	for _, item := range agg.GetSubOrders() {
		if item.GetSubOrderNo() == sub.SubOrderNo {
			target = item
			break
		}
	}
	if target == nil {
		items, itemErr := s.listOrderItemsBySubNo(ctx, sub.SubOrderNo)
		if itemErr != nil {
			return nil, itemErr
		}
		target = toProtoSub(&sub, items)
	}
	return &v1.GetShopOrderDetailRes{Order: agg, SubOrder: target}, nil
}

// MarkSubOrderShipped 将子单标记为已发货。
func (s *sOrder) MarkSubOrderShipped(ctx context.Context, req *v1.MarkSubOrderShippedReq) (*v1.MarkSubOrderShippedRes, error) {
	if req == nil || strings.TrimSpace(req.GetSubOrderNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "sub_order_no is required")
	}
	status := v1.SubOrderStatus_SUB_ORDER_STATUS_UNSPECIFIED
	// 事务内校验当前状态并更新为已发货，防止并发写覆盖。
	err := dao.OrderSub.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cols := dao.OrderSub.Columns()
		var sub entity.OrderSub
		if err := tx.Model(dao.OrderSub.Table()).Where(cols.SubOrderNo, req.GetSubOrderNo()).WhereNull(cols.DeletedAt).Scan(&sub); err != nil {
			return gerror.Wrap(err, "query sub order failed")
		}
		if sub.Id == 0 {
			return gerror.NewCode(gcode.CodeNotFound, "sub order not found")
		}
		current := v1.SubOrderStatus(sub.SubStatus)
		if current == v1.SubOrderStatus_SUB_ORDER_STATUS_SHIPPED {
			status = current
			return nil
		}
		if current != v1.SubOrderStatus_SUB_ORDER_STATUS_PAID && current != v1.SubOrderStatus_SUB_ORDER_STATUS_WAIT_SHIP {
			// 子单必须已付款或等待发货才能标记发货。
			return gerror.NewCode(gcode.CodeInvalidParameter, "sub order not shippable")
		}
		result, err := tx.Model(dao.OrderSub.Table()).Where(cols.SubOrderNo, req.GetSubOrderNo()).WhereIn(cols.SubStatus, []uint{uint(v1.SubOrderStatus_SUB_ORDER_STATUS_PAID), uint(v1.SubOrderStatus_SUB_ORDER_STATUS_WAIT_SHIP)}).Data(do.OrderSub{SubStatus: uint(v1.SubOrderStatus_SUB_ORDER_STATUS_SHIPPED), SellerRemark: req.GetShippedNote()}).Update()
		if err != nil {
			return gerror.Wrap(err, "update sub order failed")
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "sub order status changed")
		}
		status = v1.SubOrderStatus_SUB_ORDER_STATUS_SHIPPED
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &v1.MarkSubOrderShippedRes{SubOrderNo: req.GetSubOrderNo(), SubStatus: status}, nil
}

// CloseOrderIfUnpaid 关闭超时未支付的订单。
func (s *sOrder) CloseOrderIfUnpaid(ctx context.Context, req *v1.CloseOrderIfUnpaidReq) (*v1.CloseOrderIfUnpaidRes, error) {
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no is required")
	}
	reason := req.GetReasonCode()
	if reason == v1.CancelReasonCode_CANCEL_REASON_CODE_UNSPECIFIED {
		reason = v1.CancelReasonCode_CANCEL_REASON_CODE_TIMEOUT_CLOSE
	}
	var (
		closed              bool
		status              = v1.OrderStatus_ORDER_STATUS_UNSPECIFIED
		reservationNo       string
		pointsReservationNo string
		userID              uint64
	)
	// 定时任务或支付回调失败时会调用该流程，事务内关闭主单与子单。
	err := dao.OrderMain.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cols := dao.OrderMain.Columns()
		var mainRow entity.OrderMain
		if err := tx.Model(dao.OrderMain.Table()).Where(cols.OrderNo, req.GetOrderNo()).Scan(&mainRow); err != nil {
			return gerror.Wrap(err, "query order_main failed")
		}
		if mainRow.Id == 0 {
			return gerror.NewCode(gcode.CodeNotFound, "order not found")
		}
		userID = mainRow.UserId
		reservationNo = mainRow.ReservationNo
		pointsReservationNo = mainRow.PointsReservationNo
		status = v1.OrderStatus(mainRow.OrderStatus)
		if status != v1.OrderStatus_ORDER_STATUS_PENDING_PAY {
			closed = false
			return nil
		}
		// 仅当仍在待支付时才更新为已关闭，避免误关已支付订单。
		result, err := tx.Model(dao.OrderMain.Table()).Where(cols.OrderNo, req.GetOrderNo()).Where(cols.OrderStatus, uint(v1.OrderStatus_ORDER_STATUS_PENDING_PAY)).Data(do.OrderMain{OrderStatus: uint(v1.OrderStatus_ORDER_STATUS_CLOSED), CancelReasonCode: uint(reason), ClosedAt: gtime.Now(), Version: gdb.Raw(cols.Version + " + 1")}).Update()
		if err != nil {
			return gerror.Wrap(err, "close order_main failed")
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			closed = false
			return nil
		}
		subCols := dao.OrderSub.Columns()
		_, err = tx.Model(dao.OrderSub.Table()).Where(subCols.OrderNo, req.GetOrderNo()).Where(subCols.SubStatus, uint(v1.SubOrderStatus_SUB_ORDER_STATUS_PENDING_PAY)).Data(do.OrderSub{SubStatus: uint(v1.SubOrderStatus_SUB_ORDER_STATUS_CLOSED)}).Update()
		if err != nil {
			return gerror.Wrap(err, "close order_sub failed")
		}
		closed = true
		status = v1.OrderStatus_ORDER_STATUS_CLOSED
		return nil
	})
	if err != nil {
		return nil, err
	}
	if closed && reservationNo != "" {
		s.cancelInventoryReservationBestEffort(ctx, reservationNo, req.GetOrderNo(), reason.String())
	}
	if closed && pointsReservationNo != "" {
		s.compensateLockedPointsAfterOrderClose(ctx, userID, req.GetOrderNo(), pointsReservationNo, reason.String())
	}
	return &v1.CloseOrderIfUnpaidRes{OrderNo: req.GetOrderNo(), OrderStatus: status, Closed: closed}, nil
}

// HandlePayCallback 处理支付回调并推进订单状态。
func (s *sOrder) HandlePayCallback(ctx context.Context, req *v1.HandlePayCallbackReq) (*v1.HandlePayCallbackRes, error) {
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" || strings.TrimSpace(req.GetPayNo()) == "" || strings.TrimSpace(req.GetPaymentEventId()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no/pay_no/payment_event_id are required")
	}
	idem := strings.TrimSpace(req.GetIdempotencyKey())
	if idem == "" {
		idem = req.GetPaymentEventId()
	}
	// 支付回调幂等：支付网关可能多次推送相同事件，必须避免重复扣积分或重复变更订单。
	hit, row, err := s.getOrCreateIdempotency(ctx, 0, idem, idempotencyActionPayCallback)
	if err != nil {
		return nil, err
	}
	if hit {
		return s.replayPayCallback(ctx, req.GetOrderNo(), row)
	}
	paidAt := protoTsToGTime(req.GetPaidAt())
	if paidAt == nil {
		paidAt = gtime.Now()
	}
	payCode := strings.ToUpper(strings.TrimSpace(req.GetPayStatusCode()))
	isPaySuccess := payCode == "SUCCESS" || payCode == "PAID" || payCode == "TRADE_SUCCESS"
	var (
		orderStatus         = v1.OrderStatus_ORDER_STATUS_UNSPECIFIED
		paymentStatus       = v1.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
		callbackHit         bool
		refundRequired      bool
		refundReasonCode    string
		reservationNo       string
		pointsReservationNo string
		userID              uint64
	)
	// 事务确保支付记录和主单状态同步更新。
	err = dao.OrderMain.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		mainCols := dao.OrderMain.Columns()
		payCols := dao.OrderPayment.Columns()
		var byEvent entity.OrderPayment
		if err := tx.Model(dao.OrderPayment.Table()).Where(payCols.PaymentEventId, req.GetPaymentEventId()).Scan(&byEvent); err != nil {
			return gerror.Wrap(err, "query payment by event id failed")
		}
		if byEvent.Id > 0 {
			// 已处理过该事件，直接读取当前状态并返回。
			callbackHit = true
			mainRow, loadErr := s.getOrderMainByNoTx(ctx, tx, req.GetOrderNo())
			if loadErr != nil {
				return loadErr
			}
			orderStatus = v1.OrderStatus(mainRow.OrderStatus)
			paymentStatus = v1.PaymentStatus(mainRow.PaymentStatus)
			reservationNo = mainRow.ReservationNo
			pointsReservationNo = mainRow.PointsReservationNo
			userID = mainRow.UserId
			return nil
		}
		var byPayNo entity.OrderPayment
		if err := tx.Model(dao.OrderPayment.Table()).Where(payCols.PayNo, req.GetPayNo()).Scan(&byPayNo); err != nil {
			return gerror.Wrap(err, "query payment by pay_no failed")
		}
		if byPayNo.Id == 0 {
			_, err := tx.Model(dao.OrderPayment.Table()).Data(do.OrderPayment{OrderNo: req.GetOrderNo(), PayNo: req.GetPayNo(), PaymentEventId: req.GetPaymentEventId(), PayChannel: uint(req.GetPayChannel()), PayStatusCode: payCode, ChannelTradeNo: req.GetChannelTradeNo(), PaidAmount: req.GetPaidAmount(), PaidAt: paidAt, RawPayload: req.GetRawPayload()}).Insert()
			if err != nil {
				return gerror.Wrap(err, "insert order_payment failed")
			}
		} else {
			_, err := tx.Model(dao.OrderPayment.Table()).Where(payCols.PayNo, req.GetPayNo()).Data(do.OrderPayment{PaymentEventId: req.GetPaymentEventId(), PayStatusCode: payCode, ChannelTradeNo: req.GetChannelTradeNo(), PaidAmount: req.GetPaidAmount(), PaidAt: paidAt, RawPayload: req.GetRawPayload()}).Update()
			if err != nil {
				return gerror.Wrap(err, "update order_payment failed")
			}
		}
		mainRow, loadErr := s.getOrderMainByNoTx(ctx, tx, req.GetOrderNo())
		if loadErr != nil {
			return loadErr
		}
		reservationNo = mainRow.ReservationNo
		pointsReservationNo = mainRow.PointsReservationNo
		userID = mainRow.UserId
		if !isPaySuccess {
			// 支付失败仅更新支付状态，不动订单状态。
			orderStatus = v1.OrderStatus(mainRow.OrderStatus)
			paymentStatus = v1.PaymentStatus_PAYMENT_STATUS_PAY_FAILED
			_, err := tx.Model(dao.OrderMain.Table()).Where(mainCols.OrderNo, req.GetOrderNo()).Data(do.OrderMain{PaymentStatus: uint(v1.PaymentStatus_PAYMENT_STATUS_PAY_FAILED)}).Update()
			return gerror.Wrap(err, "update payment status failed")
		}
		result, err := tx.Model(dao.OrderMain.Table()).Where(mainCols.OrderNo, req.GetOrderNo()).Where(mainCols.OrderStatus, uint(v1.OrderStatus_ORDER_STATUS_PENDING_PAY)).Data(do.OrderMain{OrderStatus: uint(v1.OrderStatus_ORDER_STATUS_PAID), PaymentStatus: uint(v1.PaymentStatus_PAYMENT_STATUS_PAID), PaidAmount: req.GetPaidAmount(), PaidAt: paidAt, Version: gdb.Raw(mainCols.Version + " + 1")}).Update()
		if err != nil {
			return gerror.Wrap(err, "update order_main paid failed")
		}
		affected, _ := result.RowsAffected()
		if affected == 0 {
			latest, e := s.getOrderMainByNoTx(ctx, tx, req.GetOrderNo())
			if e != nil {
				return e
			}
			orderStatus = v1.OrderStatus(latest.OrderStatus)
			paymentStatus = v1.PaymentStatus(latest.PaymentStatus)
			if orderStatus == v1.OrderStatus_ORDER_STATUS_CLOSED || orderStatus == v1.OrderStatus_ORDER_STATUS_CANCELED {
				// 如果订单已关闭但仍收到成功回调，需要标记退款流程。
				refundRequired = true
				refundReasonCode = "ORDER_ALREADY_CLOSED"
			}
			callbackHit = orderStatus == v1.OrderStatus_ORDER_STATUS_PAID
			return nil
		}
		_, err = tx.Model(dao.OrderSub.Table()).Where(dao.OrderSub.Columns().OrderNo, req.GetOrderNo()).Where(dao.OrderSub.Columns().SubStatus, uint(v1.SubOrderStatus_SUB_ORDER_STATUS_PENDING_PAY)).Data(do.OrderSub{SubStatus: uint(v1.SubOrderStatus_SUB_ORDER_STATUS_PAID)}).Update()
		if err != nil {
			return gerror.Wrap(err, "update order_sub paid failed")
		}
		orderStatus = v1.OrderStatus_ORDER_STATUS_PAID
		paymentStatus = v1.PaymentStatus_PAYMENT_STATUS_PAID
		return nil
	})
	if err != nil {
		_ = s.markIdempotencyFailed(ctx, 0, idem, idempotencyActionPayCallback, "PAY_CALLBACK_FAILED")
		return nil, err
	}
	// 支付成功后尝试确认库存与积分，失败会回写幂等错误码。
	if !refundRequired && orderStatus == v1.OrderStatus_ORDER_STATUS_PAID && reservationNo != "" {
		if err = s.confirmInventoryReservation(ctx, reservationNo, req.GetOrderNo()); err != nil {
			if !isNotImplementedErr(err) {
				_ = s.markIdempotencyFailed(ctx, 0, idem, idempotencyActionPayCallback, "INVENTORY_CONFIRM_FAILED")
				return nil, err
			}
			g.Log().Warningf(ctx, "inventory confirm not implemented, skip confirm, order_no=%s reservation_no=%s", req.GetOrderNo(), reservationNo)
		}
	}
	if !refundRequired && orderStatus == v1.OrderStatus_ORDER_STATUS_PAID && pointsReservationNo != "" {
		if err = s.confirmLockedPoints(ctx, userID, req.GetOrderNo(), pointsReservationNo, idem); err != nil {
			_ = s.markIdempotencyFailed(ctx, 0, idem, idempotencyActionPayCallback, "POINTS_CONFIRM_FAILED")
			return nil, err
		}
	}
	if err = s.markIdempotencySuccess(ctx, 0, idem, idempotencyActionPayCallback, req.GetOrderNo(), nil); err != nil {
		return nil, err
	}
	return &v1.HandlePayCallbackRes{OrderNo: req.GetOrderNo(), PaymentStatus: paymentStatus, OrderStatus: orderStatus, CallbackIdempotentHit: callbackHit, RefundRequired: refundRequired, RefundReasonCode: refundReasonCode}, nil
}

// GetOrderSnapshotByNo 按订单号读取订单快照。
func (s *sOrder) GetOrderSnapshotByNo(ctx context.Context, req *v1.GetOrderSnapshotByNoReq) (*v1.GetOrderSnapshotByNoRes, error) {
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no is required")
	}
	agg, err := s.loadOrderAggregate(ctx, req.GetOrderNo())
	if err != nil {
		return nil, err
	}
	return &v1.GetOrderSnapshotByNoRes{Order: agg}, nil
}

type createOrderInput struct {
	OrderNo                  string
	UserID                   uint64
	AddressID                uint64
	BuyerRemark              string
	ReservationNo            string
	PointsReservationNo      string
	GoodsAmount              uint64
	FreightAmount            uint64
	DiscountAmount           uint64
	PayableAmount            uint64
	PointsUsed               uint64
	PointsDiscountAmount     uint64
	PointsRuleSnapshotJSON   string
	PointsRuleSnapshotDigest string
	PointsSubAllocations     map[string]pointsSubAllocation
	Lines                    []orderLine
}

type orderLine struct {
	ShopNo          string
	SpuNo           string
	SkuNo           string
	Qty             uint32
	SpuTitle        string
	SkuName         string
	SkuImageAssetID uint64
	SalePrice       uint64
	MarketPrice     uint64
	SaleAttrsJSON   string
}

type checkoutSnapshot struct {
	CheckoutToken  string                 `json:"checkout_token"`
	UserID         uint64                 `json:"user_id"`
	Items          []checkoutSnapshotItem `json:"items"`
	GoodsAmount    uint64                 `json:"goods_amount"`
	FreightAmount  uint64                 `json:"freight_amount"`
	PayableAmount  uint64                 `json:"payable_amount"`
	SnapshotDigest string                 `json:"snapshot_digest"`
}

type checkoutSnapshotItem struct {
	SkuNo           string `json:"sku_no"`
	SpuNo           string `json:"spu_no"`
	ShopNo          string `json:"shop_no"`
	Qty             uint32 `json:"qty"`
	SettlePrice     uint64 `json:"settle_price"`
	MarketPrice     uint64 `json:"market_price"`
	SpuTitle        string `json:"spu_title"`
	SkuName         string `json:"sku_name"`
	SkuImageAssetID uint64 `json:"sku_image_asset_id"`
	SaleAttrsJSON   string `json:"sale_attrs_json"`
}

// snapshotToOrderLines 将结算快照转换为订单行模型。
func snapshotToOrderLines(snapshot *checkoutSnapshot) []orderLine {
	if snapshot == nil || len(snapshot.Items) == 0 {
		return nil
	}
	lines := make([]orderLine, 0, len(snapshot.Items))
	for _, item := range snapshot.Items {
		lines = append(lines, orderLine{ShopNo: item.ShopNo, SpuNo: item.SpuNo, SkuNo: item.SkuNo, Qty: item.Qty, SpuTitle: item.SpuTitle, SkuName: item.SkuName, SkuImageAssetID: item.SkuImageAssetID, SalePrice: item.SettlePrice, MarketPrice: item.MarketPrice, SaleAttrsJSON: item.SaleAttrsJSON})
	}
	return lines
}

// createOrderAggregate 在事务内落库订单主表、子表与快照。
func (s *sOrder) createOrderAggregate(ctx context.Context, input *createOrderInput) (*v1.OrderMain, error) {
	if input == nil || strings.TrimSpace(input.OrderNo) == "" || len(input.Lines) == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "invalid create order input")
	}
	// subBuild 用于按店铺拆单，聚合金额后生成多个子单。
	type subBuild struct {
		SubNo                string
		ShopNo               string
		Items                []orderLine
		GoodsAmount          uint64
		FreightAmount        uint64
		PayableAmount        uint64
		PointsUsed           uint64
		PointsDiscountAmount uint64
	}
	mainGoods := input.GoodsAmount
	mainFreight := input.FreightAmount
	mainDiscount := input.DiscountAmount
	mainPayable := input.PayableAmount
	if mainGoods == 0 {
		for _, line := range input.Lines {
			mainGoods += line.SalePrice * uint64(line.Qty)
		}
	}
	if mainPayable == 0 {
		totalDiscount := mainDiscount + input.PointsDiscountAmount
		if mainGoods+mainFreight >= totalDiscount {
			mainPayable = mainGoods + mainFreight - totalDiscount
		}
	}
	// 根据店铺拆分子单，避免跨店铺物流结算互相污染。
	subMap := make(map[string]*subBuild)
	for _, line := range input.Lines {
		if strings.TrimSpace(line.SkuNo) == "" || strings.TrimSpace(line.SpuNo) == "" || strings.TrimSpace(line.ShopNo) == "" || line.Qty == 0 {
			return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order line invalid")
		}
		bucket, ok := subMap[line.ShopNo]
		if !ok {
			bucket = &subBuild{SubNo: generateBizNo("SUB"), ShopNo: line.ShopNo}
			subMap[line.ShopNo] = bucket
		}
		bucket.Items = append(bucket.Items, line)
		lineAmount := line.SalePrice * uint64(line.Qty)
		bucket.GoodsAmount += lineAmount
		bucket.PayableAmount += lineAmount
	}
	// 将积分抵扣分摊到对应店铺，保证金额一致。
	for shopNo, allocation := range input.PointsSubAllocations {
		if bucket, ok := subMap[shopNo]; ok {
			bucket.PointsUsed = allocation.PointsUsed
			bucket.PointsDiscountAmount = allocation.PointsDiscountAmount
			if bucket.PointsDiscountAmount <= bucket.PayableAmount {
				bucket.PayableAmount -= bucket.PointsDiscountAmount
			}
		}
	}
	payDeadline := gtime.NewFromTime(time.Now().Add(15 * time.Minute))
	// 事务写入主单、地址快照、子单、明细、库存关联，任何一步失败都会整体回滚。
	err := dao.OrderMain.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Model(dao.OrderMain.Table()).Data(do.OrderMain{OrderNo: input.OrderNo, UserId: input.UserID, OrderStatus: uint(v1.OrderStatus_ORDER_STATUS_PENDING_PAY), PaymentStatus: uint(v1.PaymentStatus_PAYMENT_STATUS_UNPAID), ReservationNo: input.ReservationNo, PointsReservationNo: input.PointsReservationNo, GoodsAmount: mainGoods, FreightAmount: mainFreight, DiscountAmount: mainDiscount, PayableAmount: mainPayable, PaidAmount: uint64(0), PointsUsed: input.PointsUsed, PointsDiscountAmount: input.PointsDiscountAmount, PointsRuleSnapshotJson: input.PointsRuleSnapshotJSON, PointsRuleSnapshotDigest: input.PointsRuleSnapshotDigest, BuyerRemark: input.BuyerRemark, CancelReasonCode: uint(v1.CancelReasonCode_CANCEL_REASON_CODE_UNSPECIFIED), PayDeadlineAt: payDeadline, Version: uint64(1)}).Insert()
		if err != nil {
			return gerror.Wrap(err, "insert order_main failed")
		}
		_, err = tx.Model(dao.OrderAddressSnapshot.Table()).Data(do.OrderAddressSnapshot{OrderNo: input.OrderNo, SourceAddressId: input.AddressID, SourceAddressVersion: uint64(0), ReceiverName: "", ReceiverPhone: "", CountryCode: "", ProvinceCode: "", ProvinceName: "", CityCode: "", CityName: "", DistrictCode: "", DistrictName: "", Street: "", Detail: "", PostalCode: "", Latitude: float64(0), Longitude: float64(0)}).Insert()
		if err != nil {
			return gerror.Wrap(err, "insert order_address_snapshot failed")
		}
		for _, sub := range subMap {
			_, err = tx.Model(dao.OrderSub.Table()).Data(do.OrderSub{SubOrderNo: sub.SubNo, OrderNo: input.OrderNo, ShopNo: sub.ShopNo, SubStatus: uint(v1.SubOrderStatus_SUB_ORDER_STATUS_PENDING_PAY), GoodsAmount: sub.GoodsAmount, FreightAmount: sub.FreightAmount, DiscountAmount: uint64(0), PayableAmount: sub.PayableAmount, PaidAmount: uint64(0), PointsUsed: sub.PointsUsed, PointsDiscountAmount: sub.PointsDiscountAmount, SellerRemark: "", BuyerRemark: input.BuyerRemark}).Insert()
			if err != nil {
				return gerror.Wrap(err, "insert order_sub failed")
			}
			for _, line := range sub.Items {
				_, err = tx.Model(dao.OrderItem.Table()).Data(do.OrderItem{ItemNo: generateBizNo("ITEM"), OrderNo: input.OrderNo, SubOrderNo: sub.SubNo, ShopNo: line.ShopNo, SpuNo: line.SpuNo, SkuNo: line.SkuNo, SpuTitle: line.SpuTitle, SkuName: line.SkuName, SkuImageAssetId: line.SkuImageAssetID, Qty: line.Qty, SalePrice: line.SalePrice, MarketPrice: line.MarketPrice, SaleAttrsJson: line.SaleAttrsJSON}).Insert()
				if err != nil {
					return gerror.Wrap(err, "insert order_item failed")
				}
			}
		}
		if strings.TrimSpace(input.ReservationNo) != "" {
			_, err = tx.Model(dao.OrderInventoryLink.Table()).Data(do.OrderInventoryLink{OrderNo: input.OrderNo, ReservationNo: input.ReservationNo}).Insert()
			if err != nil {
				return gerror.Wrap(err, "insert order_inventory_link failed")
			}
		}
		// 写入操作日志，记录创建动作与积分信息。
		appendOperateLogTx(ctx, tx, input.OrderNo, "", "ORDER_CREATE", "", v1.OrderStatus_ORDER_STATUS_PENDING_PAY.String(), map[string]any{
			"points_reservation_no":       input.PointsReservationNo,
			"points_used":                 input.PointsUsed,
			"points_discount_amount":      input.PointsDiscountAmount,
			"points_rule_snapshot_digest": input.PointsRuleSnapshotDigest,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.loadOrderAggregate(ctx, input.OrderNo)
}

// loadOrderAggregate 组装订单聚合视图返回给接口层。
func (s *sOrder) loadOrderAggregate(ctx context.Context, orderNo string) (*v1.OrderMain, error) {
	mainRow, err := s.getOrderMainByNo(ctx, orderNo)
	if err != nil {
		return nil, err
	}
	var (
		address entity.OrderAddressSnapshot
		subs    []*entity.OrderSub
		items   []*entity.OrderItem
	)
	// 逐表加载：地址快照、子单、明细，保持读取顺序简单可控。
	if err = dao.OrderAddressSnapshot.Ctx(ctx).Where(dao.OrderAddressSnapshot.Columns().OrderNo, orderNo).Scan(&address); err != nil {
		return nil, gerror.Wrap(err, "query address snapshot failed")
	}
	if err = dao.OrderSub.Ctx(ctx).Where(dao.OrderSub.Columns().OrderNo, orderNo).WhereNull(dao.OrderSub.Columns().DeletedAt).OrderAsc(dao.OrderSub.Columns().Id).Scan(&subs); err != nil {
		return nil, gerror.Wrap(err, "query order_sub failed")
	}
	if err = dao.OrderItem.Ctx(ctx).Where(dao.OrderItem.Columns().OrderNo, orderNo).OrderAsc(dao.OrderItem.Columns().Id).Scan(&items); err != nil {
		return nil, gerror.Wrap(err, "query order_item failed")
	}
	itemMap := make(map[string][]*v1.OrderItemSnapshot)
	for _, row := range items {
		itemMap[row.SubOrderNo] = append(itemMap[row.SubOrderNo], &v1.OrderItemSnapshot{ItemNo: row.ItemNo, OrderNo: row.OrderNo, SubOrderNo: row.SubOrderNo, ShopNo: row.ShopNo, SpuNo: row.SpuNo, SkuNo: row.SkuNo, SpuTitle: row.SpuTitle, SkuName: row.SkuName, SkuImageAssetId: row.SkuImageAssetId, Qty: uint32(row.Qty), SalePrice: row.SalePrice, MarketPrice: row.MarketPrice, SaleAttrsJson: row.SaleAttrsJson})
	}
	outSubs := make([]*v1.OrderSub, 0, len(subs))
	for _, row := range subs {
		outSubs = append(outSubs, toProtoSub(row, itemMap[row.SubOrderNo]))
	}
	var addr *v1.OrderAddressSnapshot
	if address.Id > 0 {
		addr = &v1.OrderAddressSnapshot{SourceAddressId: address.SourceAddressId, SourceAddressVersion: address.SourceAddressVersion, ReceiverName: address.ReceiverName, ReceiverPhone: address.ReceiverPhone, CountryCode: address.CountryCode, ProvinceCode: address.ProvinceCode, ProvinceName: address.ProvinceName, CityCode: address.CityCode, CityName: address.CityName, DistrictCode: address.DistrictCode, DistrictName: address.DistrictName, Street: address.Street, Detail: address.Detail, PostalCode: address.PostalCode, Latitude: address.Latitude, Longitude: address.Longitude}
	}
	return &v1.OrderMain{OrderNo: mainRow.OrderNo, UserId: mainRow.UserId, OrderStatus: v1.OrderStatus(mainRow.OrderStatus), PaymentStatus: v1.PaymentStatus(mainRow.PaymentStatus), Amount: buildAmount(mainRow.GoodsAmount, mainRow.FreightAmount, mainRow.DiscountAmount, mainRow.PayableAmount, mainRow.PaidAmount, mainRow.PointsDiscountAmount), Address: addr, ReservationNo: mainRow.ReservationNo, PayDeadlineAt: toProtoTs(mainRow.PayDeadlineAt), PaidAt: toProtoTs(mainRow.PaidAt), ClosedAt: toProtoTs(mainRow.ClosedAt), CreatedAt: toProtoTs(mainRow.CreatedAt), UpdatedAt: toProtoTs(mainRow.UpdatedAt), Version: mainRow.Version, SubOrders: outSubs, BuyerRemark: mainRow.BuyerRemark, CancelReasonCode: v1.CancelReasonCode(mainRow.CancelReasonCode), PointsReservationNo: mainRow.PointsReservationNo, PointsUsed: mainRow.PointsUsed, PointsDiscountAmount: mainRow.PointsDiscountAmount, PointsRuleSnapshotJson: mainRow.PointsRuleSnapshotJson, PointsRuleSnapshotDigest: mainRow.PointsRuleSnapshotDigest}, nil
}

// buildAmount 组装订单金额快照结构。
func buildAmount(goods, freight, discount, payable, paid, pointsDiscount uint64) *v1.OrderAmount {
	return &v1.OrderAmount{GoodsAmount: goods, FreightAmount: freight, DiscountAmount: discount, PayableAmount: payable, PaidAmount: paid, PointsDiscountAmount: pointsDiscount}
}

// listOrderItemsBySubNo 查询指定子单下的商品明细。
func (s *sOrder) listOrderItemsBySubNo(ctx context.Context, subOrderNo string) ([]*v1.OrderItemSnapshot, error) {
	var rows []*entity.OrderItem
	if err := dao.OrderItem.Ctx(ctx).Where(dao.OrderItem.Columns().SubOrderNo, subOrderNo).OrderAsc(dao.OrderItem.Columns().Id).Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "query order items failed")
	}
	out := make([]*v1.OrderItemSnapshot, 0, len(rows))
	for _, row := range rows {
		out = append(out, &v1.OrderItemSnapshot{ItemNo: row.ItemNo, OrderNo: row.OrderNo, SubOrderNo: row.SubOrderNo, ShopNo: row.ShopNo, SpuNo: row.SpuNo, SkuNo: row.SkuNo, SpuTitle: row.SpuTitle, SkuName: row.SkuName, SkuImageAssetId: row.SkuImageAssetId, Qty: uint32(row.Qty), SalePrice: row.SalePrice, MarketPrice: row.MarketPrice, SaleAttrsJson: row.SaleAttrsJson})
	}
	return out, nil
}

// toProtoSub 将子单实体转换为 protobuf 结构。
func toProtoSub(row *entity.OrderSub, items []*v1.OrderItemSnapshot) *v1.OrderSub {
	return &v1.OrderSub{SubOrderNo: row.SubOrderNo, OrderNo: row.OrderNo, ShopNo: row.ShopNo, SubStatus: v1.SubOrderStatus(row.SubStatus), Amount: buildAmount(row.GoodsAmount, row.FreightAmount, row.DiscountAmount, row.PayableAmount, row.PaidAmount, row.PointsDiscountAmount), SellerRemark: row.SellerRemark, BuyerRemark: row.BuyerRemark, CreatedAt: toProtoTs(row.CreatedAt), UpdatedAt: toProtoTs(row.UpdatedAt), Items: items, PointsUsed: row.PointsUsed, PointsDiscountAmount: row.PointsDiscountAmount}
}

// getOrderMainByNo 按订单号查询订单主记录。
func (s *sOrder) getOrderMainByNo(ctx context.Context, orderNo string) (*entity.OrderMain, error) {
	var row entity.OrderMain
	if err := dao.OrderMain.Ctx(ctx).Where(dao.OrderMain.Columns().OrderNo, orderNo).WhereNull(dao.OrderMain.Columns().DeletedAt).Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "query order_main failed")
	}
	if row.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "order not found")
	}
	return &row, nil
}

// getOrderMainByNoTx 在事务内按订单号查询订单主记录。
func (s *sOrder) getOrderMainByNoTx(ctx context.Context, tx gdb.TX, orderNo string) (*entity.OrderMain, error) {
	var row entity.OrderMain
	if err := tx.Model(dao.OrderMain.Table()).Where(dao.OrderMain.Columns().OrderNo, orderNo).Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "query order_main failed")
	}
	if row.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "order not found")
	}
	return &row, nil
}

// getOrCreateIdempotency 获取或创建幂等记录。
func (s *sOrder) getOrCreateIdempotency(ctx context.Context, userID uint64, key, action string) (bool, *entity.OrderIdempotency, error) {
	key = strings.TrimSpace(key)
	action = strings.TrimSpace(action)
	if key == "" || action == "" {
		return false, nil, gerror.NewCode(gcode.CodeInvalidParameter, "idempotency key/action are required")
	}
	// 先查是否存在，存在则表示命中幂等；不存在再插入 processing 记录。
	cols := dao.OrderIdempotency.Columns()
	var row entity.OrderIdempotency
	if err := dao.OrderIdempotency.Ctx(ctx).Where(cols.UserId, userID).Where(cols.IdempotencyKey, key).Where(cols.ActionCode, action).Scan(&row); err != nil {
		return false, nil, gerror.Wrap(err, "query idempotency failed")
	}
	if row.Id > 0 {
		return true, &row, nil
	}
	_, err := dao.OrderIdempotency.Ctx(ctx).Data(do.OrderIdempotency{UserId: userID, IdempotencyKey: key, ActionCode: action, Status: idempotencyStatusProcessing, ExpireAt: gtime.NewFromTime(time.Now().Add(24 * time.Hour))}).Insert()
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			var existed entity.OrderIdempotency
			if e := dao.OrderIdempotency.Ctx(ctx).Where(cols.UserId, userID).Where(cols.IdempotencyKey, key).Where(cols.ActionCode, action).Scan(&existed); e != nil {
				return false, nil, gerror.Wrap(e, "query duplicated idempotency failed")
			}
			// 并发插入导致唯一键冲突，返回已存在记录视作命中。
			return true, &existed, nil
		}
		return false, nil, gerror.Wrap(err, "insert idempotency failed")
	}
	return false, nil, nil
}

// markIdempotencySuccess 将幂等记录标记为成功。
func (s *sOrder) markIdempotencySuccess(ctx context.Context, userID uint64, key, action, orderNo string, response any) error {
	resp := ""
	if response != nil {
		b, err := json.Marshal(response)
		if err != nil {
			return gerror.Wrap(err, "marshal idempotency response failed")
		}
		resp = string(b)
	}
	_, err := dao.OrderIdempotency.Ctx(ctx).Where(dao.OrderIdempotency.Columns().UserId, userID).Where(dao.OrderIdempotency.Columns().IdempotencyKey, key).Where(dao.OrderIdempotency.Columns().ActionCode, action).Data(do.OrderIdempotency{Status: idempotencyStatusSuccess, OrderNo: orderNo, ResponseJson: resp, ErrorCode: ""}).Update()
	return gerror.Wrap(err, "mark idempotency success failed")
}

// markIdempotencyFailed 将幂等记录标记为失败。
func (s *sOrder) markIdempotencyFailed(ctx context.Context, userID uint64, key, action, errorCode string) error {
	_, err := dao.OrderIdempotency.Ctx(ctx).Where(dao.OrderIdempotency.Columns().UserId, userID).Where(dao.OrderIdempotency.Columns().IdempotencyKey, key).Where(dao.OrderIdempotency.Columns().ActionCode, action).Data(do.OrderIdempotency{Status: idempotencyStatusFailed, ErrorCode: errorCode}).Update()
	return gerror.Wrap(err, "mark idempotency failed failed")
}

// replayCreateOrderWithFlag 回放创建订单请求的历史结果。
func (s *sOrder) replayCreateOrderWithFlag(ctx context.Context, row *entity.OrderIdempotency) (*v1.OrderMain, bool, error) {
	if row == nil {
		return nil, false, gerror.NewCode(gcode.CodeInvalidParameter, "idempotency record nil")
	}
	switch row.Status {
	case idempotencyStatusSuccess:
		if strings.TrimSpace(row.OrderNo) == "" {
			return nil, false, gerror.NewCode(gcode.CodeInvalidParameter, "idempotency record missing order_no")
		}
		agg, err := s.loadOrderAggregate(ctx, row.OrderNo)
		if err != nil {
			return nil, false, err
		}
		return agg, true, nil
	case idempotencyStatusProcessing:
		return nil, false, gerror.NewCode(gcode.CodeInvalidParameter, "same idempotency request is processing")
	default:
		return nil, false, gerror.NewCodef(gcode.CodeInvalidParameter, "previous request failed, error_code=%s", row.ErrorCode)
	}
}

// replayRequestPay 回放请求支付的历史结果。
func (s *sOrder) replayRequestPay(ctx context.Context, orderNo string, row *entity.OrderIdempotency) (*v1.RequestPayRes, error) {
	if row == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "idempotency record nil")
	}
	switch row.Status {
	case idempotencyStatusSuccess:
		if row.OrderNo != "" {
			orderNo = row.OrderNo
		}
		return s.buildRequestPayReplay(ctx, orderNo)
	case idempotencyStatusProcessing:
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "same idempotency request is processing")
	default:
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, "previous request failed, error_code=%s", row.ErrorCode)
	}
}

// buildRequestPayReplay 构造支付请求的重放响应。
func (s *sOrder) buildRequestPayReplay(ctx context.Context, orderNo string) (*v1.RequestPayRes, error) {
	var pay entity.OrderPayment
	if err := dao.OrderPayment.Ctx(ctx).Where(dao.OrderPayment.Columns().OrderNo, orderNo).OrderDesc(dao.OrderPayment.Columns().Id).Scan(&pay); err != nil {
		return nil, gerror.Wrap(err, "query latest payment failed")
	}
	if pay.Id == 0 {
		mainRow, err := s.getOrderMainByNo(ctx, orderNo)
		if err != nil {
			return nil, err
		}
		return &v1.RequestPayRes{OrderNo: orderNo, PaymentStatus: v1.PaymentStatus(mainRow.PaymentStatus), PayPayloadJson: "{}", ExpireAt: toProtoTs(mainRow.PayDeadlineAt)}, nil
	}
	return &v1.RequestPayRes{OrderNo: orderNo, PayNo: pay.PayNo, PaymentStatus: paymentStatusFromPayCode(pay.PayStatusCode), PayUrl: fmt.Sprintf("https://mock-pay.shopa.local/pay?pay_no=%s", pay.PayNo), PayPayloadJson: "{}", ExpireAt: toProtoTs(pay.PaidAt)}, nil
}

// replayCancel 回放取消订单请求的历史结果。
func (s *sOrder) replayCancel(ctx context.Context, orderNo string, row *entity.OrderIdempotency) (*v1.CancelMyOrderRes, error) {
	if row == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "idempotency record nil")
	}
	switch row.Status {
	case idempotencyStatusSuccess:
		if row.OrderNo != "" {
			orderNo = row.OrderNo
		}
		mainRow, err := s.getOrderMainByNo(ctx, orderNo)
		if err != nil {
			return nil, err
		}
		return &v1.CancelMyOrderRes{OrderNo: orderNo, OrderStatus: v1.OrderStatus(mainRow.OrderStatus)}, nil
	case idempotencyStatusProcessing:
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "same idempotency request is processing")
	default:
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, "previous request failed, error_code=%s", row.ErrorCode)
	}
}

// replayPayCallback 回放支付回调处理结果。
func (s *sOrder) replayPayCallback(ctx context.Context, orderNo string, row *entity.OrderIdempotency) (*v1.HandlePayCallbackRes, error) {
	if row == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "idempotency record nil")
	}
	switch row.Status {
	case idempotencyStatusSuccess:
		if row.OrderNo != "" {
			orderNo = row.OrderNo
		}
		mainRow, err := s.getOrderMainByNo(ctx, orderNo)
		if err != nil {
			return nil, err
		}
		return &v1.HandlePayCallbackRes{OrderNo: orderNo, PaymentStatus: v1.PaymentStatus(mainRow.PaymentStatus), OrderStatus: v1.OrderStatus(mainRow.OrderStatus), CallbackIdempotentHit: true}, nil
	case idempotencyStatusProcessing:
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "same idempotency request is processing")
	default:
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, "previous request failed, error_code=%s", row.ErrorCode)
	}
}

// consumeCheckoutSnapshot 消费并校验结算快照。
func (s *sOrder) consumeCheckoutSnapshot(ctx context.Context, checkoutToken string) (*checkoutSnapshot, error) {
	checkoutToken = strings.TrimSpace(checkoutToken)
	if checkoutToken == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "checkout_token is required")
	}
	key := fmt.Sprintf("checkout:token:%s", checkoutToken)
	const script = `local v = redis.call("GET", KEYS[1]); if not v then return nil end; redis.call("DEL", KEYS[1]); return v`
	v, err := g.Redis().Do(ctx, "EVAL", script, 1, key)
	if err != nil {
		return nil, gerror.Wrap(err, "consume checkout token failed")
	}
	if v == nil || v.IsNil() {
		return nil, gerror.NewCode(gcode.CodeNotFound, "checkout token used or expired")
	}
	var snap checkoutSnapshot
	if err = json.Unmarshal([]byte(v.String()), &snap); err != nil {
		return nil, gerror.Wrap(err, "unmarshal checkout snapshot failed")
	}
	return &snap, nil
}

// reserveInventory 调用库存服务预留库存。
func (s *sOrder) reserveInventory(ctx context.Context, orderNo string, userID uint64, lines []orderLine) (string, error) {
	_ = ctx
	_ = orderNo
	_ = userID
	_ = lines
	// TODO(order-svc): 接入 inventory-svc 的 ReserveStock RPC。
	return "", gerror.NewCode(gcode.CodeNotImplemented, "TODO: inventory ReserveStock RPC not integrated yet")
}

// confirmInventoryReservation 确认库存预留结果。
func (s *sOrder) confirmInventoryReservation(ctx context.Context, reservationNo, orderNo string) error {
	if strings.TrimSpace(reservationNo) == "" {
		return nil
	}
	_ = ctx
	_ = orderNo
	// TODO(order-svc): 接入 inventory-svc 的 ConfirmReservation RPC。
	return gerror.NewCode(gcode.CodeNotImplemented, "TODO: inventory ConfirmReservation RPC not integrated yet")
}

// cancelInventoryReservationBestEffort 以最大努力方式取消库存预留。
func (s *sOrder) cancelInventoryReservationBestEffort(ctx context.Context, reservationNo, orderNo, reason string) {
	if strings.TrimSpace(reservationNo) == "" {
		return
	}
	_ = orderNo
	_ = reason
	g.Log().Warningf(ctx, "TODO: inventory CancelReservation RPC not integrated yet, reservation_no=%s", reservationNo)
}

// userIDFromContext 从上下文中提取当前登录用户 ID。
func userIDFromContext(ctx context.Context) (uint64, error) {
	if r := g.RequestFromCtx(ctx); r != nil {
		for _, key := range []string{"x-user-id", "X-User-Id", "user_id", "uid"} {
			if raw := strings.TrimSpace(r.Header.Get(key)); raw != "" {
				uid, err := strconv.ParseUint(raw, 10, 64)
				if err == nil && uid > 0 {
					return uid, nil
				}
			}
		}
	}
	md := grpcx.Ctx.IncomingMap(ctx)
	for _, key := range []string{"x-user-id", "user_id", "uid", "userid"} {
		if val := md.Get(key); val != nil {
			uid := gconv.Uint64(val)
			if uid > 0 {
				return uid, nil
			}
		}
	}
	return 0, gerror.NewCode(gcode.CodeNotAuthorized, "missing x-user-id")
}

// normalizePageSize 将分页大小收敛到允许范围内。
func normalizePageSize(reqSize int32) int {
	size := int(reqSize)
	if size <= 0 {
		size = defaultPageSize
	}
	if size > maxPageSize {
		size = maxPageSize
	}
	return size
}

// parseCursorID 解析游标中的主键 ID。
func parseCursorID(cursor string) (uint64, error) {
	cursor = strings.TrimSpace(cursor)
	if cursor == "" {
		return 0, nil
	}
	id, err := strconv.ParseUint(cursor, 10, 64)
	if err != nil {
		return 0, gerror.WrapCode(gcode.CodeInvalidParameter, err, "next_cursor must be uint64")
	}
	return id, nil
}

// orderStatusToUintSlice 将订单状态枚举切片转换为数据库值切片。
func orderStatusToUintSlice(in []v1.OrderStatus) []uint {
	out := make([]uint, 0, len(in))
	for _, item := range in {
		out = append(out, uint(item))
	}
	return out
}

// subStatusToUintSlice 将子单状态枚举切片转换为数据库值切片。
func subStatusToUintSlice(in []v1.SubOrderStatus) []uint {
	out := make([]uint, 0, len(in))
	for _, item := range in {
		out = append(out, uint(item))
	}
	return out
}

// paymentStatusFromPayCode 根据支付编码推导支付状态。
func paymentStatusFromPayCode(code string) v1.PaymentStatus {
	code = strings.ToUpper(strings.TrimSpace(code))
	switch code {
	case "PAYING":
		return v1.PaymentStatus_PAYMENT_STATUS_PAYING
	case "SUCCESS", "PAID", "TRADE_SUCCESS":
		return v1.PaymentStatus_PAYMENT_STATUS_PAID
	case "REFUNDED":
		return v1.PaymentStatus_PAYMENT_STATUS_REFUNDED
	case "PAY_FAILED", "FAILED":
		return v1.PaymentStatus_PAYMENT_STATUS_PAY_FAILED
	default:
		return v1.PaymentStatus_PAYMENT_STATUS_UNPAID
	}
}

// toProtoTs 将 GoFrame 时间转换为 protobuf 时间戳。
func toProtoTs(t *gtime.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(t.Time)
}

// protoTsToGTime 将 protobuf 时间戳转换为 GoFrame 时间。
func protoTsToGTime(ts *timestamppb.Timestamp) *gtime.Time {
	if ts == nil {
		return nil
	}
	return gtime.NewFromTime(ts.AsTime())
}

// generateBizNo 生成业务单号。
func generateBizNo(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%s%s%06d", prefix, now.Format("20060102150405"), now.UnixNano()%1000000)
}

// isNotImplementedErr 判断错误是否表示下游接口尚未实现。
func isNotImplementedErr(err error) bool {
	return err != nil && gerror.Code(err) == gcode.CodeNotImplemented
}

// firstNonEmpty 返回首个非空字符串。
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
