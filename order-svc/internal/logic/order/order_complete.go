package order

import (
	"context"
	"strings"

	v1 "github.com/TsingpekTao/shopa/order-svc/api/v1"
	"github.com/TsingpekTao/shopa/order-svc/internal/dao"
	"github.com/TsingpekTao/shopa/order-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/order-svc/internal/model/entity"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

const idempotencyActionCompleteOrder = "COMPLETE_ORDER"

// CompleteOrder 将订单推进到已完成状态并触发赠分。
func (s *sOrder) CompleteOrder(ctx context.Context, req *v1.CompleteOrderReq) (*v1.CompleteOrderRes, error) {
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" || strings.TrimSpace(req.GetIdempotencyKey()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no/idempotency_key are required")
	}
	// 幂等校验：用 key 限定「谁」调用过「完成订单」动作，保证重复回调不会重复落库或重复赠分。
	hit, row, err := s.getOrCreateIdempotency(ctx, 0, req.GetIdempotencyKey(), idempotencyActionCompleteOrder)
	if err != nil {
		return nil, err
	}
	if hit {
		switch row.Status {
		case idempotencyStatusSuccess:
			// 成功重放时直接读取当前订单状态，避免再次修改数据。
			mainRow, loadErr := s.getOrderMainByNo(ctx, req.GetOrderNo())
			if loadErr != nil {
				return nil, loadErr
			}
			return &v1.CompleteOrderRes{
				OrderNo:              req.GetOrderNo(),
				OrderStatus:          v1.OrderStatus(mainRow.OrderStatus),
				PointsGrantTriggered: v1.OrderStatus(mainRow.OrderStatus) == v1.OrderStatus_ORDER_STATUS_COMPLETED,
			}, nil
		case idempotencyStatusProcessing:
			return nil, gerror.NewCode(gcode.CodeInvalidParameter, "same idempotency request is processing")
		default:
			return nil, gerror.NewCodef(gcode.CodeInvalidParameter, "previous request failed, error_code=%s", row.ErrorCode)
		}
	}

	var (
		userID           uint64
		paidAmount       uint64
		orderStatus      = v1.OrderStatus_ORDER_STATUS_UNSPECIFIED
		pointsGrantReady bool
	)
	// 事务内同步推进主单和子单状态，只在全部子单满足条件时才允许完成订单。
	err = dao.OrderMain.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		mainRow, loadErr := s.getOrderMainByNoTx(ctx, tx, req.GetOrderNo())
		if loadErr != nil {
			return loadErr
		}
		userID = mainRow.UserId
		paidAmount = mainRow.PaidAmount
		orderStatus = v1.OrderStatus(mainRow.OrderStatus)
		if orderStatus == v1.OrderStatus_ORDER_STATUS_COMPLETED {
			// 已经完成时直接短路，允许重复调用。
			pointsGrantReady = true
			return nil
		}
		if orderStatus != v1.OrderStatus_ORDER_STATUS_PAID && orderStatus != v1.OrderStatus_ORDER_STATUS_FULFILLING {
			return gerror.NewCode(gcode.CodeInvalidParameter, "order not completable")
		}
		var subs []*entity.OrderSub
		if err := tx.Model(dao.OrderSub.Table()).
			Where(dao.OrderSub.Columns().OrderNo, req.GetOrderNo()).
			WhereNull(dao.OrderSub.Columns().DeletedAt).
			Scan(&subs); err != nil {
			return gerror.Wrap(err, "query order_sub for completion failed")
		}
		for _, sub := range subs {
			if sub == nil {
				continue
			}
			subStatus := v1.SubOrderStatus(sub.SubStatus)
			if subStatus != v1.SubOrderStatus_SUB_ORDER_STATUS_SHIPPED && subStatus != v1.SubOrderStatus_SUB_ORDER_STATUS_COMPLETED {
				// 任意子单未发货完成则不允许主单完成，保持状态一致性。
				return gerror.NewCode(gcode.CodeInvalidParameter, "sub order not completable")
			}
		}
		// 主单改为已完成状态，并递增版本号用于并发控制。
		mainCols := dao.OrderMain.Columns()
		if _, err := tx.Model(dao.OrderMain.Table()).
			Where(mainCols.OrderNo, req.GetOrderNo()).
			Data(do.OrderMain{
				OrderStatus: uint(v1.OrderStatus_ORDER_STATUS_COMPLETED),
				Version:     gdb.Raw(mainCols.Version + " + 1"),
			}).Update(); err != nil {
			return gerror.Wrap(err, "update order_main completed failed")
		}
		// 将仍处于已发货状态的子单同步更新为已完成。
		if _, err := tx.Model(dao.OrderSub.Table()).
			Where(dao.OrderSub.Columns().OrderNo, req.GetOrderNo()).
			Where(dao.OrderSub.Columns().SubStatus, uint(v1.SubOrderStatus_SUB_ORDER_STATUS_SHIPPED)).
			Data(do.OrderSub{SubStatus: uint(v1.SubOrderStatus_SUB_ORDER_STATUS_COMPLETED)}).Update(); err != nil {
			return gerror.Wrap(err, "update order_sub completed failed")
		}
		// 记录操作日志，便于运营排障和审计追踪。
		appendOperateLogTx(ctx, tx, req.GetOrderNo(), "", "ORDER_COMPLETE", orderStatus.String(), v1.OrderStatus_ORDER_STATUS_COMPLETED.String(), map[string]any{
			"completion_note":        req.GetCompletionNote(),
			"completion_source_code": req.GetCompletionSourceCode(),
		})
		orderStatus = v1.OrderStatus_ORDER_STATUS_COMPLETED
		pointsGrantReady = true
		return nil
	})
	if err != nil {
		_ = s.markIdempotencyFailed(ctx, 0, req.GetIdempotencyKey(), idempotencyActionCompleteOrder, "COMPLETE_ORDER_FAILED")
		return nil, err
	}
	if pointsGrantReady {
		if err = s.grantPointsByOrderCompleted(ctx, userID, req.GetOrderNo(), paidAmount, req.GetIdempotencyKey()); err != nil {
			_ = s.markIdempotencyFailed(ctx, 0, req.GetIdempotencyKey(), idempotencyActionCompleteOrder, "POINTS_GRANT_FAILED")
			return nil, err
		}
	}
	if err = s.markIdempotencySuccess(ctx, 0, req.GetIdempotencyKey(), idempotencyActionCompleteOrder, req.GetOrderNo(), nil); err != nil {
		return nil, err
	}
	return &v1.CompleteOrderRes{
		OrderNo:              req.GetOrderNo(),
		OrderStatus:          orderStatus,
		PointsGrantTriggered: pointsGrantReady,
	}, nil
}
