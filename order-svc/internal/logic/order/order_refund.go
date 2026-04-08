package order

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	inventoryv1 "github.com/TsingpekTao/shopa/inventory-svc/api/v1"
	v1 "github.com/TsingpekTao/shopa/order-svc/api/v1"
	"github.com/TsingpekTao/shopa/order-svc/internal/dao"
	"github.com/TsingpekTao/shopa/order-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/order-svc/internal/model/entity"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type orderInventoryAdminClient interface {
	BatchAdjustStockByAdmin(ctx context.Context, req *inventoryv1.BatchAdjustStockByAdminReq, opts ...grpc.CallOption) (*inventoryv1.BatchAdjustStockByAdminRes, error)
}

type refundRuleSnapshot struct {
	GrantPointsPerCent uint64 `json:"grant_points_per_cent"`
}

type pointsRefundReturnRequest struct {
	RefundNo       string `json:"refund_no"`
	OrderNo        string `json:"order_no"`
	SubOrderNo     string `json:"sub_order_no"`
	ShopNo         string `json:"shop_no"`
	UserID         uint64 `json:"user_id"`
	PointsToReturn uint64 `json:"points_to_return"`
	CashAmountCent uint64 `json:"cash_amount_cent"`
	IdempotencyKey string `json:"idempotency_key"`
	RequestSource  string `json:"request_source"`
}

type pointsRefundReturnResponse struct {
	Success            bool   `json:"success"`
	PointsReturnAmount uint64 `json:"points_return_amount"`
}

type pointsRefundReverseRequest struct {
	RefundNo             string `json:"refund_no"`
	OrderNo              string `json:"order_no"`
	SubOrderNo           string `json:"sub_order_no"`
	ShopNo               string `json:"shop_no"`
	UserID               uint64 `json:"user_id"`
	PointsToReverse      uint64 `json:"points_to_reverse"`
	ApprovedRefundAmount uint64 `json:"approved_refund_amount"`
	IdempotencyKey       string `json:"idempotency_key"`
	RequestSource        string `json:"request_source"`
}

type pointsRefundReverseResponse struct {
	Success                bool   `json:"success"`
	PointsReverseAmount    uint64 `json:"points_reverse_amount"`
	PointsCashOffsetAmount uint64 `json:"points_cash_offset_amount"`
	AccountDebtAfter       uint64 `json:"account_debt_after"`
}

type refundTarget struct {
	main           *entity.OrderMain
	sub            *entity.OrderSub
	items          []*entity.OrderItem
	paymentNo      string
	selectedItemNo string
}

var (
	newOrderInventoryAdminClient = defaultOrderInventoryAdminClient

	orderInventoryAdminClientOnce sync.Once
	orderInventoryAdminConn       *grpc.ClientConn
	orderInventoryAdminInst       orderInventoryAdminClient
	orderInventoryAdminErr        error
)

// PreviewSubOrderRefund 预览子单退款资格、现金退款金额和积分联动结果。
func (s *sOrder) PreviewSubOrderRefund(ctx context.Context, req *v1.PreviewSubOrderRefundReq) (*v1.PreviewSubOrderRefundRes, error) {
	// 预览退款时必须明确订单号和子单号，否则无法定位要退款的对象。
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" || strings.TrimSpace(req.GetSubOrderNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no/sub_order_no are required")
	}
	// 先装载当前退款目标快照，后续资格校验和金额计算都基于这份只读快照展开。
	target, err := s.loadRefundTarget(ctx, req.GetOrderNo(), req.GetSubOrderNo(), req.GetItemNo())
	if err != nil {
		return nil, err
	}
	// 只有已付款且未发货的子单才能进入预发货退款流程。
	if err = ensureRefundableSubStatus(v1.SubOrderStatus(target.sub.SubStatus)); err != nil {
		return nil, err
	}
	// 组装对外预览快照，供 aftersale-svc 和网关统一复用。
	return &v1.PreviewSubOrderRefundRes{
		Snapshot: buildRefundSubOrderSnapshot(target, v1.SubOrderStatus(target.sub.SubStatus), v1.OrderStatus(target.main.OrderStatus), v1.PaymentStatus(target.main.PaymentStatus), target.selectedItemNo, target.sub.PayableAmount),
	}, nil
}

// MarkSubOrderRefunding 在买家发起退款后立即冻结子单发货状态。
func (s *sOrder) MarkSubOrderRefunding(ctx context.Context, req *v1.MarkSubOrderRefundingReq) (*v1.MarkSubOrderRefundingRes, error) {
	// 冻结退款状态时必须知道目标订单、目标子单以及本次退款业务号。
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" || strings.TrimSpace(req.GetSubOrderNo()) == "" || strings.TrimSpace(req.GetRefundNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no/sub_order_no/refund_no are required")
	}
	// 事务里读取并推进子单状态，避免与发货或重复退款并发写冲突。
	err := dao.OrderMain.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 先读取订单主记录，后续主订单状态回写要基于它完成。
		mainRow, innerErr := s.getOrderMainByNoTx(ctx, tx, req.GetOrderNo())
		if innerErr != nil {
			return innerErr
		}
		// 再读取目标子单，确认它确实属于本订单。
		subRow, innerErr := s.getOrderSubByNoTx(ctx, tx, req.GetSubOrderNo())
		if innerErr != nil {
			return innerErr
		}
		// 子单和订单号不匹配时直接拒绝，防止跨单误操作。
		if subRow.OrderNo != mainRow.OrderNo {
			return gerror.NewCode(gcode.CodeInvalidParameter, "sub order does not belong to order")
		}
		// 已经处于退款中时直接幂等返回，不再重复改状态。
		if v1.SubOrderStatus(subRow.SubStatus) == v1.SubOrderStatus_SUB_ORDER_STATUS_REFUNDING {
			return nil
		}
		// 只有已付款或待发货子单才能被冻结到退款中。
		if innerErr = ensureRefundableSubStatus(v1.SubOrderStatus(subRow.SubStatus)); innerErr != nil {
			return innerErr
		}
		// 用 CAS 把目标子单推进到退款中，避免并发时覆盖别的状态迁移。
		result, innerErr := tx.Model(dao.OrderSub.Table()).
			Where(dao.OrderSub.Columns().SubOrderNo, subRow.SubOrderNo).
			WhereIn(dao.OrderSub.Columns().SubStatus, []uint{uint(v1.SubOrderStatus_SUB_ORDER_STATUS_PAID), uint(v1.SubOrderStatus_SUB_ORDER_STATUS_WAIT_SHIP)}).
			Data(do.OrderSub{SubStatus: uint(v1.SubOrderStatus_SUB_ORDER_STATUS_REFUNDING)}).
			Update()
		if innerErr != nil {
			return gerror.Wrap(innerErr, "mark sub order refunding failed")
		}
		// 没有命中说明子单状态已经被别的流程改变，需要上层重新读取。
		affected, _ := result.RowsAffected()
		if affected == 0 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "sub order status changed")
		}
		// 把主订单提升为退款中，确保卖家发货中心和买家订单页都能避开发货动作。
		_, innerErr = tx.Model(dao.OrderMain.Table()).
			Where(dao.OrderMain.Columns().OrderNo, mainRow.OrderNo).
			Data(do.OrderMain{
				OrderStatus: uint(v1.OrderStatus_ORDER_STATUS_REFUNDING),
				Version:     gdb.Raw(dao.OrderMain.Columns().Version + " + 1"),
			}).
			Update()
		if innerErr != nil {
			return gerror.Wrap(innerErr, "mark order refunding failed")
		}
		// 事务成功即可提交，之后再统一回读最新快照返回给上层。
		return nil
	})
	if err != nil {
		return nil, err
	}
	// 冻结成功后重新装载目标快照，保证响应返回的是数据库里的最新状态。
	target, err := s.loadRefundTarget(ctx, req.GetOrderNo(), req.GetSubOrderNo(), req.GetItemNo())
	if err != nil {
		return nil, err
	}
	// 返回退款中的最新快照，供 aftersale-svc 记录审批窗口和退款任务上下文。
	return &v1.MarkSubOrderRefundingRes{
		Snapshot: buildRefundSubOrderSnapshot(target, v1.SubOrderStatus_SUB_ORDER_STATUS_REFUNDING, v1.OrderStatus_ORDER_STATUS_REFUNDING, v1.PaymentStatus(target.main.PaymentStatus), target.selectedItemNo, target.sub.PayableAmount),
	}, nil
}

// RejectSubOrderRefund 在卖家驳回退款后恢复子单可发货状态。
func (s *sOrder) RejectSubOrderRefund(ctx context.Context, req *v1.RejectSubOrderRefundReq) (*v1.RejectSubOrderRefundRes, error) {
	// 驳回退款时必须明确订单、子单和退款业务号。
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" || strings.TrimSpace(req.GetSubOrderNo()) == "" || strings.TrimSpace(req.GetRefundNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no/sub_order_no/refund_no are required")
	}
	// 在事务里恢复目标子单，并根据所有子单状态重新汇总主订单状态。
	err := dao.OrderMain.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 先读取主订单，后续主订单状态回写需要它的业务号和版本字段。
		mainRow, innerErr := s.getOrderMainByNoTx(ctx, tx, req.GetOrderNo())
		if innerErr != nil {
			return innerErr
		}
		// 读取目标子单，确认状态是否合法。
		subRow, innerErr := s.getOrderSubByNoTx(ctx, tx, req.GetSubOrderNo())
		if innerErr != nil {
			return innerErr
		}
		// 子单与订单不一致时直接拦截，避免串单恢复。
		if subRow.OrderNo != mainRow.OrderNo {
			return gerror.NewCode(gcode.CodeInvalidParameter, "sub order does not belong to order")
		}
		// 已经恢复过待发货时直接幂等返回，不重复推进状态。
		if v1.SubOrderStatus(subRow.SubStatus) == v1.SubOrderStatus_SUB_ORDER_STATUS_WAIT_SHIP {
			return nil
		}
		// 只有退款中的子单才允许驳回恢复，其他状态说明流程已变化。
		if v1.SubOrderStatus(subRow.SubStatus) != v1.SubOrderStatus_SUB_ORDER_STATUS_REFUNDING {
			return gerror.NewCode(gcode.CodeInvalidParameter, "sub order is not refunding")
		}
		// 先把目标子单恢复为待发货，卖家后续才可以正常发货。
		result, innerErr := tx.Model(dao.OrderSub.Table()).
			Where(dao.OrderSub.Columns().SubOrderNo, subRow.SubOrderNo).
			Where(dao.OrderSub.Columns().SubStatus, uint(v1.SubOrderStatus_SUB_ORDER_STATUS_REFUNDING)).
			Data(do.OrderSub{SubStatus: uint(v1.SubOrderStatus_SUB_ORDER_STATUS_WAIT_SHIP)}).
			Update()
		if innerErr != nil {
			return gerror.Wrap(innerErr, "restore sub order wait ship failed")
		}
		// CAS 未命中说明子单状态在事务期间被别的流程先行修改。
		affected, _ := result.RowsAffected()
		if affected == 0 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "sub order status changed")
		}
		// 读取订单下所有子单状态，用于重新汇总主订单状态。
		subRows, innerErr := s.listOrderSubsByOrderNoTx(ctx, tx, mainRow.OrderNo)
		if innerErr != nil {
			return innerErr
		}
		// 用恢复后的状态覆盖内存切片，保证汇总口径和本次事务写入一致。
		statuses := make([]v1.SubOrderStatus, 0, len(subRows))
		for _, row := range subRows {
			// 当前目标子单已经恢复到待发货，内存里也按恢复后的值参与主订单汇总。
			if row.SubOrderNo == subRow.SubOrderNo {
				statuses = append(statuses, v1.SubOrderStatus_SUB_ORDER_STATUS_WAIT_SHIP)
				continue
			}
			// 其他子单沿用数据库当前状态参与汇总。
			statuses = append(statuses, v1.SubOrderStatus(row.SubStatus))
		}
		// 汇总主订单新状态，确保退款驳回后主订单从退款中回到正常履约态。
		orderStatus, paymentStatus := summarizeRefundOrderState(statuses)
		_, innerErr = tx.Model(dao.OrderMain.Table()).
			Where(dao.OrderMain.Columns().OrderNo, mainRow.OrderNo).
			Data(do.OrderMain{
				OrderStatus:   uint(orderStatus),
				PaymentStatus: uint(paymentStatus),
				Version:       gdb.Raw(dao.OrderMain.Columns().Version + " + 1"),
			}).
			Update()
		if innerErr != nil {
			return gerror.Wrap(innerErr, "update order state after refund reject failed")
		}
		// 汇总和恢复都成功时提交事务。
		return nil
	})
	if err != nil {
		return nil, err
	}
	// 事务提交后重新回读最新快照，供 aftersale-svc 和前端展示。
	target, err := s.loadRefundTarget(ctx, req.GetOrderNo(), req.GetSubOrderNo(), "")
	if err != nil {
		return nil, err
	}
	// 返回恢复后的子单快照，明确告诉上层该单重新回到待发货。
	return &v1.RejectSubOrderRefundRes{
		Snapshot: buildRefundSubOrderSnapshot(target, v1.SubOrderStatus(target.sub.SubStatus), v1.OrderStatus(target.main.OrderStatus), v1.PaymentStatus(target.main.PaymentStatus), target.selectedItemNo, target.sub.PayableAmount),
	}, nil
}

// ensureRefundableSubStatus 校验子单是否还处于预发货退款允许的状态。
func ensureRefundableSubStatus(status v1.SubOrderStatus) error {
	// 只有已付款或待发货状态才允许进入预发货退款流程。
	if status != v1.SubOrderStatus_SUB_ORDER_STATUS_PAID && status != v1.SubOrderStatus_SUB_ORDER_STATUS_WAIT_SHIP {
		return gerror.NewCode(gcode.CodeInvalidParameter, "sub order is not refundable before shipment")
	}
	return nil
}

// summarizeRefundOrderState 根据子单状态重算主订单和支付状态。
func summarizeRefundOrderState(subStatuses []v1.SubOrderStatus) (v1.OrderStatus, v1.PaymentStatus) {
	// 没有子单时回退为未指定，避免给出错误业务含义。
	if len(subStatuses) == 0 {
		return v1.OrderStatus_ORDER_STATUS_UNSPECIFIED, v1.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
	}
	// 先统计几类关键子单状态，为主订单汇总提供依据。
	var (
		hasRefunding bool
		hasShipped   bool
		hasPaid      bool
		allRefunded  = true
		allClosed    = true
	)
	for _, status := range subStatuses {
		// 只要存在退款中子单，主订单就应该停留在退款中。
		if status == v1.SubOrderStatus_SUB_ORDER_STATUS_REFUNDING {
			hasRefunding = true
		}
		// 只要出现待发货或已付款子单，就说明主订单仍有正常履约链路未完成。
		if status == v1.SubOrderStatus_SUB_ORDER_STATUS_PAID || status == v1.SubOrderStatus_SUB_ORDER_STATUS_WAIT_SHIP {
			hasPaid = true
		}
		// 只要出现已发货或已完成子单，主订单应表现为履约中。
		if status == v1.SubOrderStatus_SUB_ORDER_STATUS_SHIPPED || status == v1.SubOrderStatus_SUB_ORDER_STATUS_COMPLETED {
			hasShipped = true
		}
		// 只要存在非已退款子单，就不能把主订单视为整单已退款。
		if status != v1.SubOrderStatus_SUB_ORDER_STATUS_REFUNDED {
			allRefunded = false
		}
		// 只要出现非关闭/取消子单，就说明主订单不应该汇总成关闭态。
		if status != v1.SubOrderStatus_SUB_ORDER_STATUS_CLOSED && status != v1.SubOrderStatus_SUB_ORDER_STATUS_CANCELED {
			allClosed = false
		}
	}
	// 先处理最强语义：所有子单都已退款时，主订单和支付状态都改成已退款。
	if allRefunded {
		return v1.OrderStatus_ORDER_STATUS_REFUNDED, v1.PaymentStatus_PAYMENT_STATUS_REFUNDED
	}
	// 只要仍有子单在退款中，就继续把主订单维持在退款中。
	if hasRefunding {
		return v1.OrderStatus_ORDER_STATUS_REFUNDING, v1.PaymentStatus_PAYMENT_STATUS_PAID
	}
	// 没有活跃退款但仍有发货或完结子单时，主订单按履约中处理。
	if hasShipped {
		return v1.OrderStatus_ORDER_STATUS_FULFILLING, v1.PaymentStatus_PAYMENT_STATUS_PAID
	}
	// 剩余场景里只要还有待发货子单，主订单就按已付款待履约处理。
	if hasPaid {
		return v1.OrderStatus_ORDER_STATUS_PAID, v1.PaymentStatus_PAYMENT_STATUS_PAID
	}
	// 全部子单都已关闭或取消时，主订单同步回到关闭态。
	if allClosed {
		return v1.OrderStatus_ORDER_STATUS_CLOSED, v1.PaymentStatus_PAYMENT_STATUS_UNPAID
	}
	// 兜底返回当前系统已有的“已付款”语义，避免落成未指定态影响前端分类。
	return v1.OrderStatus_ORDER_STATUS_PAID, v1.PaymentStatus_PAYMENT_STATUS_PAID
}

// calculateRefundReversePoints 根据订单保存的积分规则快照计算本次退款需要冲回的赠分。
func calculateRefundReversePoints(ruleSnapshotJSON string, refundAmount uint64) (uint64, error) {
	// 没有规则快照时说明订单未开启赠分能力，本次直接返回 0。
	if strings.TrimSpace(ruleSnapshotJSON) == "" || refundAmount == 0 {
		return 0, nil
	}
	// 反序列化订单保存的规则快照，确保退款按下单时的规则回退赠分。
	var snapshot refundRuleSnapshot
	if err := json.Unmarshal([]byte(ruleSnapshotJSON), &snapshot); err != nil {
		return 0, gerror.Wrap(err, "unmarshal points rule snapshot failed")
	}
	// 赠分规则按“每 1 元多少积分”计算，金额单位是分，所以要除以 100。
	if snapshot.GrantPointsPerCent == 0 {
		return 0, nil
	}
	return refundAmount * snapshot.GrantPointsPerCent / 100, nil
}

// calculateRefundReturnPoints 按退款金额比例返还已使用积分。
func calculateRefundReturnPoints(pointsUsed, refundAmount, subPayableAmount uint64) uint64 {
	// 没有已使用积分或没有退款金额时，返还积分自然为 0。
	if pointsUsed == 0 || refundAmount == 0 {
		return 0
	}
	// 金额大于等于整笔子单可退现金时，直接全额返还该子单锁定的积分。
	if subPayableAmount == 0 || refundAmount >= subPayableAmount {
		return pointsUsed
	}
	// 否则按现金退款比例折算要返还的积分数量。
	return pointsUsed * refundAmount / subPayableAmount
}

// FinalizeSubOrderRefund 在支付退款成功后回写订单状态，并处理库存和积分联动。
func (s *sOrder) FinalizeSubOrderRefund(ctx context.Context, req *v1.FinalizeSubOrderRefundReq) (*v1.FinalizeSubOrderRefundRes, error) {
	// 收口退款时必须知道目标订单、目标子单和退款业务号。
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" || strings.TrimSpace(req.GetSubOrderNo()) == "" || strings.TrimSpace(req.GetRefundNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no/sub_order_no/refund_no are required")
	}
	// 先装载退款目标快照，为外部副作用和数据库回写提供统一输入。
	target, err := s.loadRefundTarget(ctx, req.GetOrderNo(), req.GetSubOrderNo(), "")
	if err != nil {
		return nil, err
	}
	// 只有退款中的子单才能进入退款成功收口；已退款则直接走幂等返回。
	currentSubStatus := v1.SubOrderStatus(target.sub.SubStatus)
	if currentSubStatus == v1.SubOrderStatus_SUB_ORDER_STATUS_REFUNDED {
		return &v1.FinalizeSubOrderRefundRes{
			Snapshot: buildRefundSubOrderSnapshot(target, currentSubStatus, v1.OrderStatus(target.main.OrderStatus), v1.PaymentStatus(target.main.PaymentStatus), target.selectedItemNo, target.sub.PayableAmount),
		}, nil
	}
	// 不是退款中的子单说明状态机已经偏离预期，不能继续收口。
	if currentSubStatus != v1.SubOrderStatus_SUB_ORDER_STATUS_REFUNDING {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "sub order is not refunding")
	}
	// 审批金额默认按整笔子单现金实付处理，并且不能超过子单可退现金金额。
	approvedRefundAmount := req.GetApprovedRefundAmount()
	if approvedRefundAmount == 0 {
		approvedRefundAmount = target.sub.PayableAmount
	}
	if approvedRefundAmount > target.sub.PayableAmount {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "approved_refund_amount exceeds refundable amount")
	}
	// 先按审批金额计算积分返还与赠分冲回口径，后续 points-svc 调用就按这份快照执行。
	pointsReturnAmount := calculateRefundReturnPoints(target.sub.PointsUsed, approvedRefundAmount, target.sub.PayableAmount)
	pointsReverseAmount, err := calculateRefundReversePoints(target.main.PointsRuleSnapshotJson, approvedRefundAmount)
	if err != nil {
		return nil, err
	}
	// 先恢复库存，确保退款成功后商品库存只释放一次。
	if err = s.restoreRefundInventory(ctx, req.GetRefundNo(), target.items); err != nil {
		return nil, err
	}
	// 再返还本次子单消费掉的积分，积分接口本身带幂等键，可以安全重试。
	if pointsReturnAmount > 0 {
		if _, err = s.returnRefundPoints(ctx, target.main, target.sub, req.GetRefundNo(), approvedRefundAmount, pointsReturnAmount); err != nil {
			return nil, err
		}
	}
	// 冲回支付成功时发放的赠分，余额不足时 points-svc 会返回现金抵扣金额和欠账结果。
	var reverseRes *pointsRefundReverseResponse
	if pointsReverseAmount > 0 {
		reverseRes, err = s.reverseRefundGrantedPoints(ctx, target.main, target.sub, req.GetRefundNo(), approvedRefundAmount, pointsReverseAmount)
		if err != nil {
			return nil, err
		}
	}
	// 外部副作用成功后，再在事务里把子单和主订单状态正式推进为已退款。
	err = dao.OrderMain.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// 重新读取目标子单，确保事务内仍然处于退款中。
		subRow, innerErr := s.getOrderSubByNoTx(ctx, tx, target.sub.SubOrderNo)
		if innerErr != nil {
			return innerErr
		}
		// 如果已经被别的重试推进为已退款，说明本次调用命中了幂等，直接返回即可。
		if v1.SubOrderStatus(subRow.SubStatus) == v1.SubOrderStatus_SUB_ORDER_STATUS_REFUNDED {
			return nil
		}
		// 只有退款中的子单才允许推进到已退款。
		if v1.SubOrderStatus(subRow.SubStatus) != v1.SubOrderStatus_SUB_ORDER_STATUS_REFUNDING {
			return gerror.NewCode(gcode.CodeInvalidParameter, "sub order status changed")
		}
		// 先把目标子单推进为已退款。
		result, innerErr := tx.Model(dao.OrderSub.Table()).
			Where(dao.OrderSub.Columns().SubOrderNo, subRow.SubOrderNo).
			Where(dao.OrderSub.Columns().SubStatus, uint(v1.SubOrderStatus_SUB_ORDER_STATUS_REFUNDING)).
			Data(do.OrderSub{SubStatus: uint(v1.SubOrderStatus_SUB_ORDER_STATUS_REFUNDED)}).
			Update()
		if innerErr != nil {
			return gerror.Wrap(innerErr, "mark sub order refunded failed")
		}
		// CAS 未命中时说明并发状态已经变化，本次请求需要重新读取。
		affected, _ := result.RowsAffected()
		if affected == 0 {
			return gerror.NewCode(gcode.CodeInvalidParameter, "sub order status changed")
		}
		// 再读取订单下全部子单状态，用于重算主订单状态。
		subRows, innerErr := s.listOrderSubsByOrderNoTx(ctx, tx, target.main.OrderNo)
		if innerErr != nil {
			return innerErr
		}
		// 用最新子单状态构造汇总所需的状态切片。
		statuses := make([]v1.SubOrderStatus, 0, len(subRows))
		for _, row := range subRows {
			// 目标子单在本次事务里已经推进到已退款，汇总时按新状态参与计算。
			if row.SubOrderNo == subRow.SubOrderNo {
				statuses = append(statuses, v1.SubOrderStatus_SUB_ORDER_STATUS_REFUNDED)
				continue
			}
			// 其他子单沿用数据库当前状态参与汇总。
			statuses = append(statuses, v1.SubOrderStatus(row.SubStatus))
		}
		// 根据所有子单状态重新汇总主订单和支付状态。
		orderStatus, paymentStatus := summarizeRefundOrderState(statuses)
		_, innerErr = tx.Model(dao.OrderMain.Table()).
			Where(dao.OrderMain.Columns().OrderNo, target.main.OrderNo).
			Data(do.OrderMain{
				OrderStatus:   uint(orderStatus),
				PaymentStatus: uint(paymentStatus),
				Version:       gdb.Raw(dao.OrderMain.Columns().Version + " + 1"),
			}).
			Update()
		if innerErr != nil {
			return gerror.Wrap(innerErr, "update order state after refund success failed")
		}
		// 主子订单状态都落库成功后提交事务。
		return nil
	})
	if err != nil {
		return nil, err
	}
	// 最后重新回读最新快照，确保返回给 aftersale-svc 的状态和数据库完全一致。
	target, err = s.loadRefundTarget(ctx, req.GetOrderNo(), req.GetSubOrderNo(), "")
	if err != nil {
		return nil, err
	}
	// 组装退款成功响应，把库存和积分联动结果一并回传给上层持久化到 refund_task。
	res := &v1.FinalizeSubOrderRefundRes{
		Snapshot:            buildRefundSubOrderSnapshot(target, v1.SubOrderStatus(target.sub.SubStatus), v1.OrderStatus(target.main.OrderStatus), v1.PaymentStatus(target.main.PaymentStatus), target.selectedItemNo, approvedRefundAmount),
		PointsReturnAmount:  pointsReturnAmount,
		PointsReverseAmount: pointsReverseAmount,
	}
	// 只有真的走了赠分冲回时才填充现金抵扣和欠账结果。
	if reverseRes != nil {
		res.PointsCashOffsetAmount = reverseRes.PointsCashOffsetAmount
		res.AccountDebtAfter = reverseRes.AccountDebtAfter
	}
	return res, nil
}

// restoreRefundInventory 在退款成功后把已扣减库存回补回去。
func (s *sOrder) restoreRefundInventory(ctx context.Context, refundNo string, items []*entity.OrderItem) error {
	// 没有商品明细时无需回补库存。
	if len(items) == 0 {
		return nil
	}
	// 建立库存管理端 gRPC 客户端，用于执行真实库存回补。
	client, err := newOrderInventoryAdminClient(ctx)
	if err != nil {
		return err
	}
	// 逐个商品构造库存回补项，并为每个 SKU 绑定稳定幂等 biz_no。
	adjustItems := make([]*inventoryv1.AdminAdjustItem, 0, len(items))
	for _, row := range items {
		// 子单退款成功后库存要按原购买数量回补到总库存。
		adjustItems = append(adjustItems, &inventoryv1.AdminAdjustItem{
			SkuNo:         row.SkuNo,
			SpuNo:         row.SpuNo,
			ShopNo:        row.ShopNo,
			DeltaTotalQty: int64(row.Qty),
			BizNo:         fmt.Sprintf("%s:%s", strings.TrimSpace(refundNo), strings.TrimSpace(row.SkuNo)),
			ReasonCode:    inventoryv1.AdjustReasonCode_ADJUST_REASON_CODE_RETURN_IN,
			Remark:        "refund restore inventory",
		})
	}
	// 执行库存回补，底层库存账本用 biz_no 保证重复调用不重复加库存。
	_, err = client.BatchAdjustStockByAdmin(ctx, &inventoryv1.BatchAdjustStockByAdminReq{Items: adjustItems})
	if err != nil {
		return gerror.Wrap(err, "restore refund inventory failed")
	}
	return nil
}

// returnRefundPoints 在退款成功后返还买家本次消费掉的积分。
func (s *sOrder) returnRefundPoints(ctx context.Context, mainRow *entity.OrderMain, subRow *entity.OrderSub, refundNo string, refundAmount uint64, pointsToReturn uint64) (*pointsRefundReturnResponse, error) {
	// 没有主订单或子单快照时无法发起积分返还。
	if mainRow == nil || subRow == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "refund points target is nil")
	}
	// 复用 order 已有的 points 配置，保证本地联调和部署环境都走同一套 points 地址。
	conf := loadPointsIntegrationConf(ctx)
	if !conf.Enabled || strings.TrimSpace(conf.BaseURL) == "" || pointsToReturn == 0 {
		return &pointsRefundReturnResponse{Success: true}, nil
	}
	// 调用 points 内部退款返还接口，幂等键稳定绑定到 refundNo + subOrderNo。
	resp := new(pointsRefundReturnResponse)
	err := doPointsJSONRequest(ctx, conf, strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.points.returnRefundPath", "/v1/points/internal/refund/return").String()), &pointsRefundReturnRequest{
		RefundNo:       refundNo,
		OrderNo:        mainRow.OrderNo,
		SubOrderNo:     subRow.SubOrderNo,
		ShopNo:         subRow.ShopNo,
		UserID:         mainRow.UserId,
		PointsToReturn: pointsToReturn,
		CashAmountCent: refundAmount,
		IdempotencyKey: buildRefundPointsReturnIdempotencyKey(refundNo, subRow.SubOrderNo),
		RequestSource:  "ORDER_REFUND_FINALIZE",
	}, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// reverseRefundGrantedPoints 在退款成功后冲回支付成功时已赠送给买家的积分。
func (s *sOrder) reverseRefundGrantedPoints(ctx context.Context, mainRow *entity.OrderMain, subRow *entity.OrderSub, refundNo string, refundAmount uint64, pointsToReverse uint64) (*pointsRefundReverseResponse, error) {
	// 主订单或子单快照缺失时无法发起赠分冲回。
	if mainRow == nil || subRow == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "refund reverse points target is nil")
	}
	// 如果没有启用积分服务或本次无需冲回赠分，直接返回空成功结果。
	conf := loadPointsIntegrationConf(ctx)
	if !conf.Enabled || strings.TrimSpace(conf.BaseURL) == "" || pointsToReverse == 0 {
		return &pointsRefundReverseResponse{Success: true}, nil
	}
	// 调用 points 内部赠分冲回接口，points-svc 会返回现金抵扣金额和账户欠账结果。
	resp := new(pointsRefundReverseResponse)
	err := doPointsJSONRequest(ctx, conf, strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.points.reverseRefundPath", "/v1/points/internal/refund/reverse").String()), &pointsRefundReverseRequest{
		RefundNo:             refundNo,
		OrderNo:              mainRow.OrderNo,
		SubOrderNo:           subRow.SubOrderNo,
		ShopNo:               subRow.ShopNo,
		UserID:               mainRow.UserId,
		PointsToReverse:      pointsToReverse,
		ApprovedRefundAmount: refundAmount,
		IdempotencyKey:       buildRefundPointsReverseIdempotencyKey(refundNo, subRow.SubOrderNo),
		RequestSource:        "ORDER_REFUND_FINALIZE",
	}, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// loadRefundTarget 装载退款所需的订单主记录、子单、明细和支付单号。
func (s *sOrder) loadRefundTarget(ctx context.Context, orderNo, subOrderNo, itemNo string) (*refundTarget, error) {
	// 先读取订单主记录，后续状态和积分规则都从这里拿。
	mainRow, err := s.getOrderMainByNo(ctx, orderNo)
	if err != nil {
		return nil, err
	}
	// 再读取目标子单，确认它真实存在。
	subRow, err := s.getOrderSubByNo(ctx, subOrderNo)
	if err != nil {
		return nil, err
	}
	// 子单不属于订单时直接拒绝，避免跨订单串写。
	if subRow.OrderNo != mainRow.OrderNo {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "sub order does not belong to order")
	}
	// 加载目标子单下的全部商品快照，退款入口虽然可以按商品点入，但 V1 实际按整笔子单处理。
	itemRows, err := s.listOrderItemRowsBySubNo(ctx, subRow.SubOrderNo)
	if err != nil {
		return nil, err
	}
	// 子单下没有商品快照说明订单数据已经损坏，不能继续退款。
	if len(itemRows) == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "sub order items not found")
	}
	// 记录本次前端点入的商品号，后续返回给 aftersale-svc 做“商品入口提升为整笔子单”的提示。
	selectedItemNo := strings.TrimSpace(itemNo)
	if selectedItemNo == "" {
		selectedItemNo = itemRows[0].ItemNo
	}
	// 如果前端带了 item_no，就必须确认该商品确实属于目标子单。
	if strings.TrimSpace(itemNo) != "" {
		found := false
		for _, row := range itemRows {
			// 命中目标商品时说明本次点入来源合法。
			if row.ItemNo == strings.TrimSpace(itemNo) {
				found = true
				break
			}
		}
		// 未命中时直接拒绝，避免跨子单商品错误提升。
		if !found {
			return nil, gerror.NewCode(gcode.CodeInvalidParameter, "item does not belong to sub order")
		}
	}
	// 读取最近一次支付单号，退款创建任务时需要把现金退款挂到对应支付流水上。
	paymentNo, err := s.getLatestPaymentNoByOrder(ctx, mainRow.OrderNo)
	if err != nil {
		return nil, err
	}
	// 返回退款目标统一快照，供预览、冻结、驳回、成功收口共用。
	return &refundTarget{
		main:           mainRow,
		sub:            subRow,
		items:          itemRows,
		paymentNo:      paymentNo,
		selectedItemNo: selectedItemNo,
	}, nil
}

// buildRefundSubOrderSnapshot 组装退款编排需要的子单快照。
func buildRefundSubOrderSnapshot(target *refundTarget, subStatus v1.SubOrderStatus, orderStatus v1.OrderStatus, paymentStatus v1.PaymentStatus, itemNo string, refundAmount uint64) *v1.RefundSubOrderSnapshot {
	// 空目标时没有可返回的业务快照。
	if target == nil || target.main == nil || target.sub == nil {
		return nil
	}
	// 按当前退款金额重新计算积分返还和赠分冲回口径，保证返回值和本次操作一致。
	pointsReturnAmount := calculateRefundReturnPoints(target.sub.PointsUsed, refundAmount, target.sub.PayableAmount)
	pointsReverseAmount, _ := calculateRefundReversePoints(target.main.PointsRuleSnapshotJson, refundAmount)
	// 把数据库里的订单项实体转换成对外统一的订单商品快照。
	items := make([]*v1.OrderItemSnapshot, 0, len(target.items))
	for _, row := range target.items {
		// 逐条转换订单项，确保退款预览页可以直接展示商品信息。
		items = append(items, &v1.OrderItemSnapshot{
			ItemNo:          row.ItemNo,
			OrderNo:         row.OrderNo,
			SubOrderNo:      row.SubOrderNo,
			ShopNo:          row.ShopNo,
			SpuNo:           row.SpuNo,
			SkuNo:           row.SkuNo,
			SpuTitle:        row.SpuTitle,
			SkuName:         row.SkuName,
			SkuImageAssetId: row.SkuImageAssetId,
			Qty:             uint32(row.Qty),
			SalePrice:       row.SalePrice,
			MarketPrice:     row.MarketPrice,
			SaleAttrsJson:   row.SaleAttrsJson,
		})
	}
	// 返回统一退款快照，后续 aftersale-svc 会把它展开成退款批次和子单售后单。
	return &v1.RefundSubOrderSnapshot{
		OrderNo:             target.main.OrderNo,
		SubOrderNo:          target.sub.SubOrderNo,
		ItemNo:              itemNo,
		ShopNo:              target.sub.ShopNo,
		UserId:              target.main.UserId,
		SubStatus:           subStatus,
		RefundableAmount:    refundAmount,
		PointsReturnAmount:  pointsReturnAmount,
		PointsReverseAmount: pointsReverseAmount,
		PaymentNo:           target.paymentNo,
		Items:               items,
		OrderStatus:         orderStatus,
		PaymentStatus:       paymentStatus,
	}
}

// getOrderSubByNo 按子单号查询子单主记录。
func (s *sOrder) getOrderSubByNo(ctx context.Context, subOrderNo string) (*entity.OrderSub, error) {
	// 预留实体承接数据库扫描结果。
	var row entity.OrderSub
	// 按业务号精确读取子单，退款流程的所有状态切换都以它为锚点。
	if err := dao.OrderSub.Ctx(ctx).Where(dao.OrderSub.Columns().SubOrderNo, subOrderNo).WhereNull(dao.OrderSub.Columns().DeletedAt).Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "query order_sub failed")
	}
	// 未命中时明确返回 not found，避免上层把空对象当成可退款子单。
	if row.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "sub order not found")
	}
	return &row, nil
}

// getOrderSubByNoTx 在事务里按子单号读取子单记录。
func (s *sOrder) getOrderSubByNoTx(ctx context.Context, tx gdb.TX, subOrderNo string) (*entity.OrderSub, error) {
	// 预留实体承接事务内读取结果。
	var row entity.OrderSub
	// 事务内直接读子单表，确保状态校验和后续更新看到的是同一视图。
	if err := tx.Model(dao.OrderSub.Table()).Where(dao.OrderSub.Columns().SubOrderNo, subOrderNo).Scan(&row); err != nil {
		return nil, gerror.Wrap(err, "query order_sub failed")
	}
	// 未命中时明确返回 not found。
	if row.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "sub order not found")
	}
	return &row, nil
}

// listOrderSubsByOrderNoTx 在事务里加载订单下全部子单。
func (s *sOrder) listOrderSubsByOrderNoTx(ctx context.Context, tx gdb.TX, orderNo string) ([]*entity.OrderSub, error) {
	// 预留切片承接订单下全部子单。
	var rows []*entity.OrderSub
	// 按订单号读取全部子单，用于主订单状态汇总。
	if err := tx.Model(dao.OrderSub.Table()).Where(dao.OrderSub.Columns().OrderNo, orderNo).WhereNull(dao.OrderSub.Columns().DeletedAt).OrderAsc(dao.OrderSub.Columns().Id).Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "query order sub list failed")
	}
	return rows, nil
}

// listOrderItemRowsBySubNo 查询指定子单下的订单项实体。
func (s *sOrder) listOrderItemRowsBySubNo(ctx context.Context, subOrderNo string) ([]*entity.OrderItem, error) {
	// 预留实体切片承接数据库扫描结果。
	var rows []*entity.OrderItem
	// 按子单号加载全部订单项，为退款预览和库存回补提供商品明细。
	if err := dao.OrderItem.Ctx(ctx).Where(dao.OrderItem.Columns().SubOrderNo, subOrderNo).OrderAsc(dao.OrderItem.Columns().Id).Scan(&rows); err != nil {
		return nil, gerror.Wrap(err, "query order items failed")
	}
	return rows, nil
}

// getLatestPaymentNoByOrder 读取订单最近一次支付单号。
func (s *sOrder) getLatestPaymentNoByOrder(ctx context.Context, orderNo string) (string, error) {
	// 预留支付记录实体承接最近一条支付流水。
	var row entity.OrderPayment
	// 取最新支付流水是为了把退款挂回真实支付单号。
	if err := dao.OrderPayment.Ctx(ctx).Where(dao.OrderPayment.Columns().OrderNo, orderNo).OrderDesc(dao.OrderPayment.Columns().Id).Scan(&row); err != nil {
		return "", gerror.Wrap(err, "query latest payment failed")
	}
	// 未支付流水时不能继续创建现金退款任务。
	if row.Id == 0 || strings.TrimSpace(row.PayNo) == "" {
		return "", gerror.NewCode(gcode.CodeNotFound, "payment no not found")
	}
	return row.PayNo, nil
}

// buildRefundPointsReturnIdempotencyKey 构造退款返还已用积分的幂等键。
func buildRefundPointsReturnIdempotencyKey(refundNo, subOrderNo string) string {
	// 用退款号和子单号构成稳定键，保证重试不会重复返还积分。
	return fmt.Sprintf("REFUND_RETURN_POINTS:%s:%s", strings.TrimSpace(refundNo), strings.TrimSpace(subOrderNo))
}

// buildRefundPointsReverseIdempotencyKey 构造退款冲回赠分的幂等键。
func buildRefundPointsReverseIdempotencyKey(refundNo, subOrderNo string) string {
	// 用退款号和子单号构成稳定键，保证赠分冲回幂等。
	return fmt.Sprintf("REFUND_REVERSE_POINTS:%s:%s", strings.TrimSpace(refundNo), strings.TrimSpace(subOrderNo))
}

// defaultOrderInventoryAdminClient 懒加载库存管理端 gRPC 客户端。
func defaultOrderInventoryAdminClient(ctx context.Context) (orderInventoryAdminClient, error) {
	// 先解析库存服务 gRPC 地址，没有配置时直接报错，避免假成功。
	addr := strings.TrimSpace(g.Cfg().MustGet(ctx, "upstream.inventory.grpcTarget", "127.0.0.1:9005").String())
	if addr == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "inventory grpc target not configured")
	}
	// 连接只初始化一次，避免每笔退款都重复拨号。
	orderInventoryAdminClientOnce.Do(func() {
		// 为首次建连增加超时，避免配置错误时请求长时间挂起。
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		// 本地开发环境统一使用明文 gRPC 连接。
		conn, err := grpc.DialContext(timeoutCtx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			orderInventoryAdminErr = gerror.Wrap(err, "dial inventory admin grpc failed")
			return
		}
		// 保存连接和 client，后续请求直接复用。
		orderInventoryAdminConn = conn
		orderInventoryAdminInst = inventoryv1.NewAdminInventoryServiceClient(conn)
	})
	// 首次建连失败时直接把错误回传给上层。
	if orderInventoryAdminErr != nil {
		return nil, orderInventoryAdminErr
	}
	return orderInventoryAdminInst, nil
}
