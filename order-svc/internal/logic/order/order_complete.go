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

const (
	idempotencyActionCompleteOrder        = "COMPLETE_ORDER"
	idempotencyActionConfirmOrderReceived = "CONFIRM_ORDER_RECEIVED"
)

// ConfirmMyOrderReceived 允许买家在待收货阶段手动确认收货。
func (s *sOrder) ConfirmMyOrderReceived(ctx context.Context, req *v1.CompleteOrderReq) (*v1.CompleteOrderRes, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "request is required")
	}
	nextReq := *req
	if strings.TrimSpace(nextReq.CompletionSourceCode) == "" {
		nextReq.CompletionSourceCode = "BUYER_CONFIRM_RECEIPT"
	}
	if strings.TrimSpace(nextReq.CompletionNote) == "" {
		nextReq.CompletionNote = "buyer confirmed receipt"
	}
	return s.completeOrder(ctx, &nextReq, idempotencyActionConfirmOrderReceived, userID, true)
}

// CompleteOrder 将订单推进到已完成状态并触发赠分。
func (s *sOrder) CompleteOrder(ctx context.Context, req *v1.CompleteOrderReq) (*v1.CompleteOrderRes, error) {
	return s.completeOrder(ctx, req, idempotencyActionCompleteOrder, 0, false)
}

func (s *sOrder) completeOrder(ctx context.Context, req *v1.CompleteOrderReq, action string, expectedUserID uint64, enforceOwner bool) (*v1.CompleteOrderRes, error) {
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" || strings.TrimSpace(req.GetIdempotencyKey()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no/idempotency_key are required")
	}
	hit, row, err := s.getOrCreateIdempotency(ctx, expectedUserID, req.GetIdempotencyKey(), action)
	if err != nil {
		return nil, err
	}
	if hit {
		switch row.Status {
		case idempotencyStatusSuccess:
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
		pointsGranted    bool
	)

	err = dao.OrderMain.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		mainRow, loadErr := s.getOrderMainByNoTx(ctx, tx, req.GetOrderNo())
		if loadErr != nil {
			return loadErr
		}
		if enforceOwner && mainRow.UserId != expectedUserID {
			return gerror.NewCode(gcode.CodeNotFound, "order not found")
		}
		userID = mainRow.UserId
		paidAmount = mainRow.PaidAmount
		orderStatus = v1.OrderStatus(mainRow.OrderStatus)
		if orderStatus == v1.OrderStatus_ORDER_STATUS_COMPLETED {
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
				return gerror.NewCode(gcode.CodeInvalidParameter, "sub order not completable")
			}
		}

		mainCols := dao.OrderMain.Columns()
		if _, err := tx.Model(dao.OrderMain.Table()).
			Where(mainCols.OrderNo, req.GetOrderNo()).
			Data(do.OrderMain{
				OrderStatus: uint(v1.OrderStatus_ORDER_STATUS_COMPLETED),
				Version:     gdb.Raw(mainCols.Version + " + 1"),
			}).Update(); err != nil {
			return gerror.Wrap(err, "update order_main completed failed")
		}

		if _, err := tx.Model(dao.OrderSub.Table()).
			Where(dao.OrderSub.Columns().OrderNo, req.GetOrderNo()).
			Where(dao.OrderSub.Columns().SubStatus, uint(v1.SubOrderStatus_SUB_ORDER_STATUS_SHIPPED)).
			Data(do.OrderSub{SubStatus: uint(v1.SubOrderStatus_SUB_ORDER_STATUS_COMPLETED)}).Update(); err != nil {
			return gerror.Wrap(err, "update order_sub completed failed")
		}

		appendOperateLogTx(ctx, tx, req.GetOrderNo(), "", "ORDER_COMPLETE", orderStatus.String(), v1.OrderStatus_ORDER_STATUS_COMPLETED.String(), map[string]any{
			"completion_note":        req.GetCompletionNote(),
			"completion_source_code": req.GetCompletionSourceCode(),
		})
		orderStatus = v1.OrderStatus_ORDER_STATUS_COMPLETED
		pointsGrantReady = true
		return nil
	})
	if err != nil {
		_ = s.markIdempotencyFailed(ctx, expectedUserID, req.GetIdempotencyKey(), action, "COMPLETE_ORDER_FAILED")
		return nil, err
	}
	if pointsGrantReady {
		if err = s.grantPointsByOrderCompleted(ctx, userID, req.GetOrderNo(), paidAmount, req.GetIdempotencyKey()); err != nil {
			s.compensateGrantPointsAfterOrderCompleted(ctx, userID, req.GetOrderNo(), err)
		} else {
			pointsGranted = true
		}
	}
	if err = s.markIdempotencySuccess(ctx, expectedUserID, req.GetIdempotencyKey(), action, req.GetOrderNo(), nil); err != nil {
		return nil, err
	}
	return &v1.CompleteOrderRes{
		OrderNo:              req.GetOrderNo(),
		OrderStatus:          orderStatus,
		PointsGrantTriggered: pointsGranted,
	}, nil
}
