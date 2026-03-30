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

func New() *sOrder { return &sOrder{} }

func init() {
	service.RegisterOrder(New())
}

func (s *sOrder) CreateOrderFromCart(ctx context.Context, req *v1.CreateOrderFromCartReq) (*v1.CreateOrderFromCartRes, error) {
	if req == nil || strings.TrimSpace(req.GetCheckoutToken()) == "" || strings.TrimSpace(req.GetIdempotencyKey()) == "" || req.GetAddressId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "checkout_token/address_id/idempotency_key are required")
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
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
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateFromCart, "TOKEN_USER_MISMATCH")
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "checkout token user mismatch")
	}
	if req.GetExpectedSnapshotDigest() != "" && snap.SnapshotDigest != "" && req.GetExpectedSnapshotDigest() != snap.SnapshotDigest {
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateFromCart, "SNAPSHOT_DIGEST_MISMATCH")
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "snapshot digest mismatch")
	}
	lines := snapshotToOrderLines(snap)
	if len(lines) == 0 {
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateFromCart, "EMPTY_SNAPSHOT")
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "snapshot has no items")
	}
	orderNo := generateBizNo("ORD")
	reservationNo, err := s.reserveInventory(ctx, orderNo, userID, lines)
	if err != nil {
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateFromCart, "INVENTORY_RESERVE_FAILED")
		return nil, err
	}
	agg, err := s.createOrderAggregate(ctx, &createOrderInput{
		OrderNo:        orderNo,
		UserID:         userID,
		AddressID:      req.GetAddressId(),
		BuyerRemark:    req.GetBuyerRemark(),
		ReservationNo:  reservationNo,
		GoodsAmount:    snap.GoodsAmount,
		FreightAmount:  snap.FreightAmount,
		DiscountAmount: 0,
		PayableAmount:  snap.PayableAmount,
		Lines:          lines,
	})
	if err != nil {
		s.cancelInventoryReservationBestEffort(ctx, reservationNo, orderNo, "ORDER_CREATE_FAILED")
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateFromCart, "ORDER_CREATE_FAILED")
		return nil, err
	}
	if err = s.markIdempotencySuccess(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateFromCart, orderNo, nil); err != nil {
		return nil, err
	}
	return &v1.CreateOrderFromCartRes{Order: agg, IdempotentReplay: false}, nil
}

func (s *sOrder) CreateOrderBuyNow(ctx context.Context, req *v1.CreateOrderBuyNowReq) (*v1.CreateOrderBuyNowRes, error) {
	if req == nil || len(req.GetItems()) == 0 || strings.TrimSpace(req.GetIdempotencyKey()) == "" || req.GetAddressId() == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "items/address_id/idempotency_key are required")
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
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
	reservationNo, err := s.reserveInventory(ctx, orderNo, userID, lines)
	if err != nil {
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateBuyNow, "INVENTORY_RESERVE_FAILED")
		return nil, err
	}
	agg, err := s.createOrderAggregate(ctx, &createOrderInput{OrderNo: orderNo, UserID: userID, AddressID: req.GetAddressId(), BuyerRemark: req.GetBuyerRemark(), ReservationNo: reservationNo, Lines: lines})
	if err != nil {
		s.cancelInventoryReservationBestEffort(ctx, reservationNo, orderNo, "ORDER_CREATE_FAILED")
		_ = s.markIdempotencyFailed(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateBuyNow, "ORDER_CREATE_FAILED")
		return nil, err
	}
	if err = s.markIdempotencySuccess(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCreateBuyNow, orderNo, nil); err != nil {
		return nil, err
	}
	return &v1.CreateOrderBuyNowRes{Order: agg, IdempotentReplay: false}, nil
}

func (s *sOrder) RequestPay(ctx context.Context, req *v1.RequestPayReq) (*v1.RequestPayRes, error) {
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" || strings.TrimSpace(req.GetIdempotencyKey()) == "" || req.GetPayChannel() == v1.PayChannel_PAY_CHANNEL_UNSPECIFIED {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no/pay_channel/idempotency_key are required")
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	hit, row, err := s.getOrCreateIdempotency(ctx, userID, req.GetIdempotencyKey(), idempotencyActionRequestPay)
	if err != nil {
		return nil, err
	}
	if hit {
		return s.replayRequestPay(ctx, req.GetOrderNo(), row)
	}
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
func (s *sOrder) CancelMyOrder(ctx context.Context, req *v1.CancelMyOrderReq) (*v1.CancelMyOrderRes, error) {
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" || strings.TrimSpace(req.GetIdempotencyKey()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no/idempotency_key are required")
	}
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	hit, row, err := s.getOrCreateIdempotency(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCancelOrder)
	if err != nil {
		return nil, err
	}
	if hit {
		return s.replayCancel(ctx, req.GetOrderNo(), row)
	}
	reason := req.GetReasonCode()
	if reason == v1.CancelReasonCode_CANCEL_REASON_CODE_UNSPECIFIED {
		reason = v1.CancelReasonCode_CANCEL_REASON_CODE_BUYER_CANCEL
	}
	var (
		status        = v1.OrderStatus_ORDER_STATUS_UNSPECIFIED
		reservationNo string
	)
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
		current := v1.OrderStatus(mainRow.OrderStatus)
		if current == v1.OrderStatus_ORDER_STATUS_CANCELED || current == v1.OrderStatus_ORDER_STATUS_CLOSED {
			status = current
			return nil
		}
		if current != v1.OrderStatus_ORDER_STATUS_PENDING_PAY {
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
	if err = s.markIdempotencySuccess(ctx, userID, req.GetIdempotencyKey(), idempotencyActionCancelOrder, req.GetOrderNo(), nil); err != nil {
		return nil, err
	}
	return &v1.CancelMyOrderRes{OrderNo: req.GetOrderNo(), OrderStatus: status}, nil
}

func (s *sOrder) GetMyOrderDetail(ctx context.Context, req *v1.GetMyOrderDetailReq) (*v1.GetMyOrderDetailRes, error) {
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no is required")
	}
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

func (s *sOrder) ListMyOrders(ctx context.Context, req *v1.ListMyOrdersReq) (*v1.ListMyOrdersRes, error) {
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

func (s *sOrder) ListShopOrders(ctx context.Context, req *v1.ListShopOrdersReq) (*v1.ListShopOrdersRes, error) {
	if req == nil || strings.TrimSpace(req.GetShopNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "shop_no is required")
	}
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

func (s *sOrder) GetShopOrderDetail(ctx context.Context, req *v1.GetShopOrderDetailReq) (*v1.GetShopOrderDetailRes, error) {
	if req == nil || strings.TrimSpace(req.GetSubOrderNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "sub_order_no is required")
	}
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

func (s *sOrder) MarkSubOrderShipped(ctx context.Context, req *v1.MarkSubOrderShippedReq) (*v1.MarkSubOrderShippedRes, error) {
	if req == nil || strings.TrimSpace(req.GetSubOrderNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "sub_order_no is required")
	}
	status := v1.SubOrderStatus_SUB_ORDER_STATUS_UNSPECIFIED
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
func (s *sOrder) CloseOrderIfUnpaid(ctx context.Context, req *v1.CloseOrderIfUnpaidReq) (*v1.CloseOrderIfUnpaidRes, error) {
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no is required")
	}
	reason := req.GetReasonCode()
	if reason == v1.CancelReasonCode_CANCEL_REASON_CODE_UNSPECIFIED {
		reason = v1.CancelReasonCode_CANCEL_REASON_CODE_TIMEOUT_CLOSE
	}
	var (
		closed        bool
		status        = v1.OrderStatus_ORDER_STATUS_UNSPECIFIED
		reservationNo string
	)
	err := dao.OrderMain.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		cols := dao.OrderMain.Columns()
		var mainRow entity.OrderMain
		if err := tx.Model(dao.OrderMain.Table()).Where(cols.OrderNo, req.GetOrderNo()).Scan(&mainRow); err != nil {
			return gerror.Wrap(err, "query order_main failed")
		}
		if mainRow.Id == 0 {
			return gerror.NewCode(gcode.CodeNotFound, "order not found")
		}
		reservationNo = mainRow.ReservationNo
		status = v1.OrderStatus(mainRow.OrderStatus)
		if status != v1.OrderStatus_ORDER_STATUS_PENDING_PAY {
			closed = false
			return nil
		}
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
	return &v1.CloseOrderIfUnpaidRes{OrderNo: req.GetOrderNo(), OrderStatus: status, Closed: closed}, nil
}

func (s *sOrder) HandlePayCallback(ctx context.Context, req *v1.HandlePayCallbackReq) (*v1.HandlePayCallbackRes, error) {
	if req == nil || strings.TrimSpace(req.GetOrderNo()) == "" || strings.TrimSpace(req.GetPayNo()) == "" || strings.TrimSpace(req.GetPaymentEventId()) == "" {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "order_no/pay_no/payment_event_id are required")
	}
	idem := strings.TrimSpace(req.GetIdempotencyKey())
	if idem == "" {
		idem = req.GetPaymentEventId()
	}
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
		orderStatus      = v1.OrderStatus_ORDER_STATUS_UNSPECIFIED
		paymentStatus    = v1.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
		callbackHit      bool
		refundRequired   bool
		refundReasonCode string
		reservationNo    string
	)
	err = dao.OrderMain.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		mainCols := dao.OrderMain.Columns()
		payCols := dao.OrderPayment.Columns()
		var byEvent entity.OrderPayment
		if err := tx.Model(dao.OrderPayment.Table()).Where(payCols.PaymentEventId, req.GetPaymentEventId()).Scan(&byEvent); err != nil {
			return gerror.Wrap(err, "query payment by event id failed")
		}
		if byEvent.Id > 0 {
			callbackHit = true
			mainRow, loadErr := s.getOrderMainByNoTx(ctx, tx, req.GetOrderNo())
			if loadErr != nil {
				return loadErr
			}
			orderStatus = v1.OrderStatus(mainRow.OrderStatus)
			paymentStatus = v1.PaymentStatus(mainRow.PaymentStatus)
			reservationNo = mainRow.ReservationNo
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
		if !isPaySuccess {
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
	if !refundRequired && orderStatus == v1.OrderStatus_ORDER_STATUS_PAID && reservationNo != "" {
		if err = s.confirmInventoryReservation(ctx, reservationNo, req.GetOrderNo()); err != nil {
			_ = s.markIdempotencyFailed(ctx, 0, idem, idempotencyActionPayCallback, "INVENTORY_CONFIRM_FAILED")
			return nil, err
		}
	}
	if err = s.markIdempotencySuccess(ctx, 0, idem, idempotencyActionPayCallback, req.GetOrderNo(), nil); err != nil {
		return nil, err
	}
	return &v1.HandlePayCallbackRes{OrderNo: req.GetOrderNo(), PaymentStatus: paymentStatus, OrderStatus: orderStatus, CallbackIdempotentHit: callbackHit, RefundRequired: refundRequired, RefundReasonCode: refundReasonCode}, nil
}

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
	OrderNo        string
	UserID         uint64
	AddressID      uint64
	BuyerRemark    string
	ReservationNo  string
	GoodsAmount    uint64
	FreightAmount  uint64
	DiscountAmount uint64
	PayableAmount  uint64
	Lines          []orderLine
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
func (s *sOrder) createOrderAggregate(ctx context.Context, input *createOrderInput) (*v1.OrderMain, error) {
	if input == nil || strings.TrimSpace(input.OrderNo) == "" || len(input.Lines) == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "invalid create order input")
	}
	type subBuild struct {
		SubNo         string
		ShopNo        string
		Items         []orderLine
		GoodsAmount   uint64
		FreightAmount uint64
		PayableAmount uint64
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
	if mainPayable == 0 && mainGoods+mainFreight >= mainDiscount {
		mainPayable = mainGoods + mainFreight - mainDiscount
	}
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
	payDeadline := gtime.NewFromTime(time.Now().Add(15 * time.Minute))
	err := dao.OrderMain.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Model(dao.OrderMain.Table()).Data(do.OrderMain{OrderNo: input.OrderNo, UserId: input.UserID, OrderStatus: uint(v1.OrderStatus_ORDER_STATUS_PENDING_PAY), PaymentStatus: uint(v1.PaymentStatus_PAYMENT_STATUS_UNPAID), ReservationNo: input.ReservationNo, GoodsAmount: mainGoods, FreightAmount: mainFreight, DiscountAmount: mainDiscount, PayableAmount: mainPayable, PaidAmount: uint64(0), BuyerRemark: input.BuyerRemark, CancelReasonCode: uint(v1.CancelReasonCode_CANCEL_REASON_CODE_UNSPECIFIED), PayDeadlineAt: payDeadline, Version: uint64(1)}).Insert()
		if err != nil {
			return gerror.Wrap(err, "insert order_main failed")
		}
		_, err = tx.Model(dao.OrderAddressSnapshot.Table()).Data(do.OrderAddressSnapshot{OrderNo: input.OrderNo, SourceAddressId: input.AddressID, SourceAddressVersion: uint64(0), ReceiverName: "", ReceiverPhone: "", CountryCode: "", ProvinceCode: "", ProvinceName: "", CityCode: "", CityName: "", DistrictCode: "", DistrictName: "", Street: "", Detail: "", PostalCode: "", Latitude: float64(0), Longitude: float64(0)}).Insert()
		if err != nil {
			return gerror.Wrap(err, "insert order_address_snapshot failed")
		}
		for _, sub := range subMap {
			_, err = tx.Model(dao.OrderSub.Table()).Data(do.OrderSub{SubOrderNo: sub.SubNo, OrderNo: input.OrderNo, ShopNo: sub.ShopNo, SubStatus: uint(v1.SubOrderStatus_SUB_ORDER_STATUS_PENDING_PAY), GoodsAmount: sub.GoodsAmount, FreightAmount: sub.FreightAmount, DiscountAmount: uint64(0), PayableAmount: sub.PayableAmount, PaidAmount: uint64(0), SellerRemark: "", BuyerRemark: input.BuyerRemark}).Insert()
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
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.loadOrderAggregate(ctx, input.OrderNo)
}

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
	return &v1.OrderMain{OrderNo: mainRow.OrderNo, UserId: mainRow.UserId, OrderStatus: v1.OrderStatus(mainRow.OrderStatus), PaymentStatus: v1.PaymentStatus(mainRow.PaymentStatus), Amount: buildAmount(mainRow.GoodsAmount, mainRow.FreightAmount, mainRow.DiscountAmount, mainRow.PayableAmount, mainRow.PaidAmount), Address: addr, ReservationNo: mainRow.ReservationNo, PayDeadlineAt: toProtoTs(mainRow.PayDeadlineAt), PaidAt: toProtoTs(mainRow.PaidAt), ClosedAt: toProtoTs(mainRow.ClosedAt), CreatedAt: toProtoTs(mainRow.CreatedAt), UpdatedAt: toProtoTs(mainRow.UpdatedAt), Version: mainRow.Version, SubOrders: outSubs, BuyerRemark: mainRow.BuyerRemark, CancelReasonCode: v1.CancelReasonCode(mainRow.CancelReasonCode)}, nil
}

func buildAmount(goods, freight, discount, payable, paid uint64) *v1.OrderAmount {
	return &v1.OrderAmount{GoodsAmount: goods, FreightAmount: freight, DiscountAmount: discount, PayableAmount: payable, PaidAmount: paid}
}

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

func toProtoSub(row *entity.OrderSub, items []*v1.OrderItemSnapshot) *v1.OrderSub {
	return &v1.OrderSub{SubOrderNo: row.SubOrderNo, OrderNo: row.OrderNo, ShopNo: row.ShopNo, SubStatus: v1.SubOrderStatus(row.SubStatus), Amount: buildAmount(row.GoodsAmount, row.FreightAmount, row.DiscountAmount, row.PayableAmount, row.PaidAmount), SellerRemark: row.SellerRemark, BuyerRemark: row.BuyerRemark, CreatedAt: toProtoTs(row.CreatedAt), UpdatedAt: toProtoTs(row.UpdatedAt), Items: items}
}

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
func (s *sOrder) getOrCreateIdempotency(ctx context.Context, userID uint64, key, action string) (bool, *entity.OrderIdempotency, error) {
	key = strings.TrimSpace(key)
	action = strings.TrimSpace(action)
	if key == "" || action == "" {
		return false, nil, gerror.NewCode(gcode.CodeInvalidParameter, "idempotency key/action are required")
	}
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
			return true, &existed, nil
		}
		return false, nil, gerror.Wrap(err, "insert idempotency failed")
	}
	return false, nil, nil
}

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

func (s *sOrder) markIdempotencyFailed(ctx context.Context, userID uint64, key, action, errorCode string) error {
	_, err := dao.OrderIdempotency.Ctx(ctx).Where(dao.OrderIdempotency.Columns().UserId, userID).Where(dao.OrderIdempotency.Columns().IdempotencyKey, key).Where(dao.OrderIdempotency.Columns().ActionCode, action).Data(do.OrderIdempotency{Status: idempotencyStatusFailed, ErrorCode: errorCode}).Update()
	return gerror.Wrap(err, "mark idempotency failed failed")
}

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

func (s *sOrder) reserveInventory(ctx context.Context, orderNo string, userID uint64, lines []orderLine) (string, error) {
	_ = ctx
	_ = orderNo
	_ = userID
	_ = lines
	// TODO(order-svc): 接入 inventory-svc 的 ReserveStock RPC。
	return "", gerror.NewCode(gcode.CodeNotImplemented, "TODO: inventory ReserveStock RPC not integrated yet")
}

func (s *sOrder) confirmInventoryReservation(ctx context.Context, reservationNo, orderNo string) error {
	if strings.TrimSpace(reservationNo) == "" {
		return nil
	}
	_ = ctx
	_ = orderNo
	// TODO(order-svc): 接入 inventory-svc 的 ConfirmReservation RPC。
	return gerror.NewCode(gcode.CodeNotImplemented, "TODO: inventory ConfirmReservation RPC not integrated yet")
}

func (s *sOrder) cancelInventoryReservationBestEffort(ctx context.Context, reservationNo, orderNo, reason string) {
	if strings.TrimSpace(reservationNo) == "" {
		return
	}
	_ = orderNo
	_ = reason
	g.Log().Warningf(ctx, "TODO: inventory CancelReservation RPC not integrated yet, reservation_no=%s", reservationNo)
}

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

func orderStatusToUintSlice(in []v1.OrderStatus) []uint {
	out := make([]uint, 0, len(in))
	for _, item := range in {
		out = append(out, uint(item))
	}
	return out
}

func subStatusToUintSlice(in []v1.SubOrderStatus) []uint {
	out := make([]uint, 0, len(in))
	for _, item := range in {
		out = append(out, uint(item))
	}
	return out
}

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

func toProtoTs(t *gtime.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(t.Time)
}

func protoTsToGTime(ts *timestamppb.Timestamp) *gtime.Time {
	if ts == nil {
		return nil
	}
	return gtime.NewFromTime(ts.AsTime())
}

func generateBizNo(prefix string) string {
	now := time.Now()
	return fmt.Sprintf("%s%s%06d", prefix, now.Format("20060102150405"), now.UnixNano()%1000000)
}
