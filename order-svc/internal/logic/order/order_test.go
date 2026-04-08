package order

import (
	"context"
	"database/sql"
	"net/url"
	"strings"
	"testing"
	"time"

	cartv1 "github.com/TsingpekTao/shopa/cart-svc/api/v1"
	catalogv1 "github.com/TsingpekTao/shopa/catalog-svc/api/v1"
	inventoryv1 "github.com/TsingpekTao/shopa/inventory-svc/api/v1"
	v1 "github.com/TsingpekTao/shopa/order-svc/api/v1"
	"github.com/TsingpekTao/shopa/order-svc/internal/model/entity"
	paymentv1 "github.com/TsingpekTao/shopa/payment-svc/api/v1"
	userprofilev1 "github.com/TsingpekTao/shopa/user-profile-svc/api/v1"
	"github.com/gogf/gf/v2/os/gtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestIsIdempotencyMissErrorTreatsSQLNoRowsAsMiss(t *testing.T) {
	if !isIdempotencyMissError(sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows to be treated as idempotency miss")
	}
}

func TestDecodeCheckoutSnapshotPayloadSupportsProtoJSONNumbers(t *testing.T) {
	payload := []byte(`{
		"checkout_token":"chk_test",
		"user_id":"2040403725247320064",
		"items":[
			{
				"sku_no":"SKU1",
				"spu_no":"SPU1",
				"shop_no":"SHOP1",
				"qty":1,
				"settle_price":"1999",
				"market_price":"2999",
				"sale_attrs_json":"[]"
			}
		],
		"goods_amount":"1999",
		"freight_amount":"0",
		"payable_amount":"1999",
		"snapshot_digest":"digest"
	}`)

	snapshot, err := decodeCheckoutSnapshotPayload(payload)
	if err != nil {
		t.Fatalf("expected protojson checkout snapshot to decode, got error: %v", err)
	}
	if snapshot.UserID != 2040403725247320064 {
		t.Fatalf("expected user_id to decode, got %d", snapshot.UserID)
	}
	if snapshot.GoodsAmount != 1999 {
		t.Fatalf("expected goods_amount to decode, got %d", snapshot.GoodsAmount)
	}
	if len(snapshot.Items) != 1 {
		t.Fatalf("expected 1 snapshot item, got %d", len(snapshot.Items))
	}
}

func TestNormalizeOptionalJSONColumnValueBlankReturnsNil(t *testing.T) {
	if got := normalizeOptionalJSONColumnValue("   "); got != nil {
		t.Fatalf("expected blank json column value to become nil, got %#v", got)
	}
	if got := normalizeOptionalJSONColumnValue(`{"rule":"ok"}`); got != `{"rule":"ok"}` {
		t.Fatalf("expected non-empty json column value to stay unchanged, got %#v", got)
	}
}

func TestHydrateOrderLinesFromSnapshotsFillsMissingBuyNowFields(t *testing.T) {
	lines := []orderLine{
		{SkuNo: "SKU1", SpuNo: "SPU1", ShopNo: "SHOP1", Qty: 2},
	}
	snapshots := map[string]*catalogv1.SkuOrderSnapshot{
		"SKU1": {
			SkuNo:           "SKU1",
			SpuNo:           "SPU1",
			ShopNo:          "SHOP1",
			SpuTitle:        "Leather Jacket",
			SkuName:         "Black-S",
			SkuImageAssetId: 8899,
			SalePrice:       29900,
			MarketPrice:     39900,
			SaleAttrs: []*catalogv1.SkuSaleAttr{
				{AttrCode: "color", AttrName: "Color", Value: "Black"},
				{AttrCode: "size", AttrName: "Size", Value: "S"},
			},
		},
	}

	if err := hydrateOrderLinesFromSnapshots(lines, snapshots, true); err != nil {
		t.Fatalf("expected hydrateOrderLinesFromSnapshots to succeed, got error: %v", err)
	}
	if lines[0].SpuTitle != "Leather Jacket" {
		t.Fatalf("expected spu title to be filled, got %q", lines[0].SpuTitle)
	}
	if lines[0].SkuName != "Black-S" {
		t.Fatalf("expected sku name to be filled, got %q", lines[0].SkuName)
	}
	if lines[0].SkuImageAssetID != 8899 {
		t.Fatalf("expected sku image asset id to be filled, got %d", lines[0].SkuImageAssetID)
	}
	if lines[0].SalePrice != 29900 {
		t.Fatalf("expected sale price to be filled, got %d", lines[0].SalePrice)
	}
	if lines[0].MarketPrice != 39900 {
		t.Fatalf("expected market price to be filled, got %d", lines[0].MarketPrice)
	}
	if lines[0].SaleAttrsJSON == "" || lines[0].SaleAttrsJSON == "[]" {
		t.Fatalf("expected sale attrs json to be filled, got %q", lines[0].SaleAttrsJSON)
	}
}

func TestHydrateOrderLinesFromSnapshotsBackfillsImageWithoutOverwritingPrice(t *testing.T) {
	lines := []orderLine{
		{SkuNo: "SKU2", SpuNo: "SPU2", ShopNo: "SHOP2", Qty: 1, SalePrice: 18800, MarketPrice: 28800},
	}
	snapshots := map[string]*catalogv1.SkuOrderSnapshot{
		"SKU2": {
			SkuNo:           "SKU2",
			SpuNo:           "SPU2",
			ShopNo:          "SHOP2",
			SpuTitle:        "Jeans",
			SkuName:         "Blue-M",
			SkuImageAssetId: 5566,
			SalePrice:       19900,
			MarketPrice:     29900,
		},
	}

	if err := hydrateOrderLinesFromSnapshots(lines, snapshots, false); err != nil {
		t.Fatalf("expected hydrateOrderLinesFromSnapshots to succeed, got error: %v", err)
	}
	if lines[0].SpuTitle != "Jeans" {
		t.Fatalf("expected missing spu title to be backfilled, got %q", lines[0].SpuTitle)
	}
	if lines[0].SkuName != "Blue-M" {
		t.Fatalf("expected missing sku name to be backfilled, got %q", lines[0].SkuName)
	}
	if lines[0].SkuImageAssetID != 5566 {
		t.Fatalf("expected missing image asset id to be backfilled, got %d", lines[0].SkuImageAssetID)
	}
	if lines[0].SalePrice != 18800 {
		t.Fatalf("expected existing sale price to stay unchanged, got %d", lines[0].SalePrice)
	}
	if lines[0].MarketPrice != 28800 {
		t.Fatalf("expected existing market price to stay unchanged, got %d", lines[0].MarketPrice)
	}
}

func TestCreatePaymentIntentUsesPaymentServiceLaunch(t *testing.T) {
	oldFactory := newBuyerPaymentClient
	defer func() {
		newBuyerPaymentClient = oldFactory
	}()

	var captured *paymentv1.CreatePaymentIntentReq
	newBuyerPaymentClient = func(context.Context) (buyerPaymentIntentClient, error) {
		return buyerPaymentIntentClientFunc(func(ctx context.Context, req *paymentv1.CreatePaymentIntentReq, _ ...grpc.CallOption) (*paymentv1.CreatePaymentIntentRes, error) {
			captured = req
			return &paymentv1.CreatePaymentIntentRes{
				Intent: &paymentv1.PaymentIntent{
					PaymentNo:       "PAY202604060001",
					GatewayExpireAt: timestamppb.New(time.Date(2026, 4, 6, 13, 30, 0, 0, time.FixedZone("CST", 8*3600))),
				},
				PayUrl:         "https://openapi-sandbox.dl.alipaydev.com/gateway.do?foo=bar",
				PayPayloadJson: `{"gateway":"alipay"}`,
			}, nil
		}), nil
	}

	svc := &sOrder{}
	mainRow := &entity.OrderMain{
		OrderNo:       "ORD202604060001",
		PayableAmount: 1999,
		PayDeadlineAt: gtime.NewFromTime(time.Date(2026, 4, 6, 13, 45, 0, 0, time.FixedZone("CST", 8*3600))),
	}

	res, err := svc.createPaymentIntent(context.Background(), mainRow, 2040403725247320064, v1.PayChannel_PAY_CHANNEL_ALIPAY, "idem-request-pay-1")
	if err != nil {
		t.Fatalf("expected createPaymentIntent to succeed, got error: %v", err)
	}
	if res.GetPayUrl() != "https://openapi-sandbox.dl.alipaydev.com/gateway.do?foo=bar" {
		t.Fatalf("expected payment-svc pay_url to be returned, got %q", res.GetPayUrl())
	}
	if captured == nil {
		t.Fatalf("expected payment-svc request to be captured")
	}
	if captured.GetOrderNo() != "ORD202604060001" {
		t.Fatalf("expected order_no to be forwarded, got %q", captured.GetOrderNo())
	}
	if captured.GetUserId() != 2040403725247320064 {
		t.Fatalf("expected user_id to be forwarded, got %d", captured.GetUserId())
	}
	if captured.GetPayChannel() != paymentv1.PayChannel_PAY_CHANNEL_ALIPAY {
		t.Fatalf("expected pay_channel to be mapped to payment proto, got %v", captured.GetPayChannel())
	}
	if captured.GetPayableAmount() != 1999 {
		t.Fatalf("expected payable_amount to be forwarded, got %d", captured.GetPayableAmount())
	}
	if captured.GetSubject() == "" {
		t.Fatalf("expected a non-empty payment subject")
	}
	expectedReturnURL := "http://localhost:3100/me/orders?order_no=" + url.QueryEscape("ORD202604060001")
	if captured.GetReturnUrl() != expectedReturnURL {
		t.Fatalf("expected return_url to point at buyer order page, got %q", captured.GetReturnUrl())
	}
}

func TestReserveInventoryUsesInventoryServiceReservation(t *testing.T) {
	oldFactory := newOrderInventoryClient
	defer func() {
		newOrderInventoryClient = oldFactory
	}()

	var captured *inventoryv1.ReserveStockReq
	newOrderInventoryClient = func(context.Context) (orderInventoryClient, error) {
		return orderInventoryClientStub{
			reserve: func(ctx context.Context, req *inventoryv1.ReserveStockReq, _ ...grpc.CallOption) (*inventoryv1.ReserveStockRes, error) {
				captured = req
				return &inventoryv1.ReserveStockRes{
					Success: true,
					Reservation: &inventoryv1.ReservationRecord{
						ReservationNo: "RSV202604060001",
					},
				}, nil
			},
		}, nil
	}

	svc := &sOrder{}
	lines := []orderLine{
		{SkuNo: "SKU1", SpuNo: "SPU1", ShopNo: "SHOP1", Qty: 2},
		{SkuNo: "SKU2", SpuNo: "SPU2", ShopNo: "SHOP1", Qty: 1},
	}

	reservationNo, err := svc.reserveInventory(context.Background(), "ORD202604060001", 2040403725247320064, lines)
	if err != nil {
		t.Fatalf("expected reserveInventory to succeed, got error: %v", err)
	}
	if reservationNo != "RSV202604060001" {
		t.Fatalf("expected reservation number from inventory-svc, got %q", reservationNo)
	}
	if captured == nil {
		t.Fatalf("expected inventory reserve request to be captured")
	}
	if captured.GetOrderNo() != "ORD202604060001" {
		t.Fatalf("expected order_no to be forwarded, got %q", captured.GetOrderNo())
	}
	if captured.GetUserId() != 2040403725247320064 {
		t.Fatalf("expected user_id to be forwarded, got %d", captured.GetUserId())
	}
	if captured.GetMode() != inventoryv1.ReserveMode_RESERVE_MODE_ALL_OR_NOTHING {
		t.Fatalf("expected all-or-nothing reserve mode, got %v", captured.GetMode())
	}
	if strings.TrimSpace(captured.GetReservationNo()) == "" {
		t.Fatalf("expected generated reservation_no to be forwarded")
	}
	if len(captured.GetItems()) != 2 {
		t.Fatalf("expected 2 reserve items, got %d", len(captured.GetItems()))
	}
	if captured.GetItems()[0].GetSkuNo() != "SKU1" || captured.GetItems()[0].GetQty() != 2 {
		t.Fatalf("expected first reserve item to map sku/qty, got %#v", captured.GetItems()[0])
	}
	if captured.GetItems()[1].GetSkuNo() != "SKU2" || captured.GetItems()[1].GetQty() != 1 {
		t.Fatalf("expected second reserve item to map sku/qty, got %#v", captured.GetItems()[1])
	}
}

func TestConfirmInventoryReservationCallsInventoryService(t *testing.T) {
	oldFactory := newOrderInventoryClient
	defer func() {
		newOrderInventoryClient = oldFactory
	}()

	var captured *inventoryv1.ConfirmReservationReq
	newOrderInventoryClient = func(context.Context) (orderInventoryClient, error) {
		return orderInventoryClientStub{
			confirm: func(ctx context.Context, req *inventoryv1.ConfirmReservationReq, _ ...grpc.CallOption) (*inventoryv1.ConfirmReservationRes, error) {
				captured = req
				return &inventoryv1.ConfirmReservationRes{
					Reservation: &inventoryv1.ReservationRecord{
						ReservationNo: req.GetReservationNo(),
					},
				}, nil
			},
		}, nil
	}

	svc := &sOrder{}
	if err := svc.confirmInventoryReservation(context.Background(), "RSV202604060001", "ORD202604060001"); err != nil {
		t.Fatalf("expected confirmInventoryReservation to succeed, got error: %v", err)
	}
	if captured == nil {
		t.Fatalf("expected inventory confirm request to be captured")
	}
	if captured.GetReservationNo() != "RSV202604060001" {
		t.Fatalf("expected reservation_no to be forwarded, got %q", captured.GetReservationNo())
	}
	if captured.GetOrderNo() != "ORD202604060001" {
		t.Fatalf("expected order_no to be forwarded, got %q", captured.GetOrderNo())
	}
}

func TestCancelInventoryReservationBestEffortCallsInventoryService(t *testing.T) {
	oldFactory := newOrderInventoryClient
	defer func() {
		newOrderInventoryClient = oldFactory
	}()

	var captured *inventoryv1.CancelReservationReq
	newOrderInventoryClient = func(context.Context) (orderInventoryClient, error) {
		return orderInventoryClientStub{
			cancel: func(ctx context.Context, req *inventoryv1.CancelReservationReq, _ ...grpc.CallOption) (*inventoryv1.CancelReservationRes, error) {
				captured = req
				return &inventoryv1.CancelReservationRes{}, nil
			},
		}, nil
	}

	svc := &sOrder{}
	svc.cancelInventoryReservationBestEffort(context.Background(), "RSV202604060001", "ORD202604060001", "ORDER_TIMEOUT")

	if captured == nil {
		t.Fatalf("expected inventory cancel request to be captured")
	}
	if captured.GetReservationNo() != "RSV202604060001" {
		t.Fatalf("expected reservation_no to be forwarded, got %q", captured.GetReservationNo())
	}
	if captured.GetOrderNo() != "ORD202604060001" {
		t.Fatalf("expected order_no to be forwarded, got %q", captured.GetOrderNo())
	}
	if captured.GetReasonCode() != "ORDER_TIMEOUT" {
		t.Fatalf("expected reason_code to be forwarded, got %q", captured.GetReasonCode())
	}
}

func TestRequestPayPayloadJSONRoundTrip(t *testing.T) {
	expireAt := timestamppb.New(time.Date(2026, 4, 6, 13, 30, 0, 0, time.UTC))

	raw := buildRequestPayPayloadJSON("https://sandbox.pay/redirect", `{"gateway":"alipay"}`, expireAt)
	payload := parseRequestPayPayloadJSON(raw)

	if payload.PayURL != "https://sandbox.pay/redirect" {
		t.Fatalf("expected pay_url to round-trip, got %q", payload.PayURL)
	}
	if payload.PayPayloadJSON != `{"gateway":"alipay"}` {
		t.Fatalf("expected pay_payload_json to round-trip, got %q", payload.PayPayloadJSON)
	}
	if payload.ExpireAt == nil || !payload.ExpireAt.AsTime().Equal(expireAt.AsTime()) {
		t.Fatalf("expected expire_at to round-trip, got %#v", payload.ExpireAt)
	}
}

func TestBuildOrderAddressSnapshotDOMapsUserProfileSnapshot(t *testing.T) {
	addr := buildOrderAddressSnapshotDO(
		"ORD202604060001",
		&userprofilev1.AddressSnapshot{
			SourceAddressId:      7,
			SourceAddressVersion: 13,
			ReceiverName:         "Alice",
			ReceiverPhone:        "13800001111",
			CountryCode:          "CN",
			ProvinceCode:         "110000",
			ProvinceName:         "Beijing",
			CityCode:             "110100",
			CityName:             "Beijing",
			DistrictCode:         "110101",
			DistrictName:         "Dongcheng",
			Street:               "Chang'an Avenue",
			Detail:               "No.1",
			PostalCode:           "100000",
			Latitude:             39.9042,
			Longitude:            116.4074,
		},
	)

	if addr.OrderNo != "ORD202604060001" {
		t.Fatalf("expected order_no to be mapped, got %v", addr.OrderNo)
	}
	if addr.SourceAddressId != uint64(7) {
		t.Fatalf("expected source_address_id to be mapped, got %#v", addr.SourceAddressId)
	}
	if addr.ReceiverName != "Alice" {
		t.Fatalf("expected receiver_name to be mapped, got %#v", addr.ReceiverName)
	}
	if addr.Detail != "No.1" {
		t.Fatalf("expected detail to be mapped, got %#v", addr.Detail)
	}
}

func TestCanUpdateOrderAddressAllowsPendingOrdersBeforePaymentSuccess(t *testing.T) {
	allowed := &entity.OrderMain{
		OrderStatus:   uint(v1.OrderStatus_ORDER_STATUS_PENDING_PAY),
		PaymentStatus: uint(v1.PaymentStatus_PAYMENT_STATUS_UNPAID),
	}
	if !canUpdateOrderAddress(allowed) {
		t.Fatalf("expected unpaid pending order to allow address update")
	}

	paying := &entity.OrderMain{
		OrderStatus:   uint(v1.OrderStatus_ORDER_STATUS_PENDING_PAY),
		PaymentStatus: uint(v1.PaymentStatus_PAYMENT_STATUS_PAYING),
	}
	if !canUpdateOrderAddress(paying) {
		t.Fatalf("expected paying pending order to allow address update")
	}

	payFailed := &entity.OrderMain{
		OrderStatus:   uint(v1.OrderStatus_ORDER_STATUS_PENDING_PAY),
		PaymentStatus: uint(v1.PaymentStatus_PAYMENT_STATUS_PAY_FAILED),
	}
	if !canUpdateOrderAddress(payFailed) {
		t.Fatalf("expected failed pending order to allow address update")
	}

	paid := &entity.OrderMain{
		OrderStatus:   uint(v1.OrderStatus_ORDER_STATUS_PAID),
		PaymentStatus: uint(v1.PaymentStatus_PAYMENT_STATUS_PAID),
	}
	if canUpdateOrderAddress(paid) {
		t.Fatalf("expected paid order to reject address update")
	}
}

func TestShouldInsertOrderPaymentOnScanMiss(t *testing.T) {
	insert, err := shouldInsertOrderPayment(nil, sql.ErrNoRows)
	if err != nil {
		t.Fatalf("expected sql.ErrNoRows to be treated as insert path, got error: %v", err)
	}
	if !insert {
		t.Fatalf("expected missing pay_no row to require insert")
	}
}

func TestShouldInsertOrderPaymentOnExistingRow(t *testing.T) {
	insert, err := shouldInsertOrderPayment(&entity.OrderPayment{Id: 12}, nil)
	if err != nil {
		t.Fatalf("expected existing row to avoid error, got: %v", err)
	}
	if insert {
		t.Fatalf("expected existing pay_no row to update instead of insert")
	}
}

func TestShouldTreatMissingPaymentEventAsInsertPath(t *testing.T) {
	insert, err := shouldInsertOrderPayment(nil, sql.ErrNoRows)
	if err != nil {
		t.Fatalf("expected sql.ErrNoRows to be treated as insert path for payment_event_id lookup, got error: %v", err)
	}
	if !insert {
		t.Fatalf("expected missing payment_event_id row to require insert")
	}
}

type buyerPaymentIntentClientFunc func(ctx context.Context, req *paymentv1.CreatePaymentIntentReq, opts ...grpc.CallOption) (*paymentv1.CreatePaymentIntentRes, error)

func (f buyerPaymentIntentClientFunc) CreatePaymentIntent(ctx context.Context, req *paymentv1.CreatePaymentIntentReq, opts ...grpc.CallOption) (*paymentv1.CreatePaymentIntentRes, error) {
	return f(ctx, req, opts...)
}

type orderInventoryClientStub struct {
	reserve func(ctx context.Context, req *inventoryv1.ReserveStockReq, opts ...grpc.CallOption) (*inventoryv1.ReserveStockRes, error)
	confirm func(ctx context.Context, req *inventoryv1.ConfirmReservationReq, opts ...grpc.CallOption) (*inventoryv1.ConfirmReservationRes, error)
	cancel  func(ctx context.Context, req *inventoryv1.CancelReservationReq, opts ...grpc.CallOption) (*inventoryv1.CancelReservationRes, error)
}

func TestCollectSnapshotSkuNosKeepsStableUniqueOrder(t *testing.T) {
	snapshot := &checkoutSnapshot{
		Items: []checkoutSnapshotItem{
			{SkuNo: "SKU-1"},
			{SkuNo: " SKU-2 "},
			{SkuNo: "SKU-1"},
			{SkuNo: ""},
		},
	}

	skuNos := collectSnapshotSkuNos(snapshot)
	if len(skuNos) != 2 {
		t.Fatalf("expected 2 unique sku numbers, got %d", len(skuNos))
	}
	if skuNos[0] != "SKU-1" || skuNos[1] != "SKU-2" {
		t.Fatalf("expected stable trimmed sku order, got %#v", skuNos)
	}
}

func TestCollectOrderItemSkuNosKeepsStableUniqueOrder(t *testing.T) {
	order := &v1.OrderMain{
		SubOrders: []*v1.OrderSub{
			{
				Items: []*v1.OrderItemSnapshot{
					{SkuNo: "SKU-1"},
					{SkuNo: " SKU-2 "},
				},
			},
			{
				Items: []*v1.OrderItemSnapshot{
					{SkuNo: "SKU-1"},
					{SkuNo: ""},
					{SkuNo: "SKU-3"},
				},
			},
		},
	}

	skuNos := collectOrderItemSkuNos(order)
	if len(skuNos) != 3 {
		t.Fatalf("expected 3 unique sku numbers, got %d", len(skuNos))
	}
	if skuNos[0] != "SKU-1" || skuNos[1] != "SKU-2" || skuNos[2] != "SKU-3" {
		t.Fatalf("expected stable trimmed sku order from paid order aggregate, got %#v", skuNos)
	}
}

func TestRemoveOrderedCartItemsBestEffortCallsClient(t *testing.T) {
	oldFactory := buyerCartClientFactory
	defer func() {
		buyerCartClientFactory = oldFactory
	}()

	var captured []string
	buyerCartClientFactory = func(ctx context.Context) (cartv1.BuyerCartServiceClient, error) {
		return buyerCartClientStub{
			remove: func(ctx context.Context, req *cartv1.RemoveItemsReq, _ ...grpc.CallOption) (*cartv1.RemoveItemsRes, error) {
				captured = append([]string(nil), req.GetSkuNos()...)
				return &cartv1.RemoveItemsRes{}, nil
			},
		}, nil
	}

	s := &sOrder{}
	s.removeOrderedCartItemsBestEffort(context.Background(), 2040403725247320064, []string{"SKU-A", "SKU-B"})
	if len(captured) != 2 || captured[0] != "SKU-A" || captured[1] != "SKU-B" {
		t.Fatalf("expected cart removal to send requested sku nos, got %#v", captured)
	}
}

func TestRemoveOrderedCartItemsBestEffortSkipsInvalidInput(t *testing.T) {
	oldFactory := buyerCartClientFactory
	defer func() {
		buyerCartClientFactory = oldFactory
	}()

	callCount := 0
	buyerCartClientFactory = func(ctx context.Context) (cartv1.BuyerCartServiceClient, error) {
		callCount++
		return buyerCartClientStub{}, nil
	}

	s := &sOrder{}
	s.removeOrderedCartItemsBestEffort(context.Background(), 0, []string{"SKU-A"})
	if callCount != 0 {
		t.Fatalf("expected client not to be built when user_id is missing")
	}

	s.removeOrderedCartItemsBestEffort(context.Background(), 2040, nil)
	if callCount != 0 {
		t.Fatalf("expected client not to be built when no sku nos are provided")
	}
}

func (s orderInventoryClientStub) ReserveStock(ctx context.Context, req *inventoryv1.ReserveStockReq, opts ...grpc.CallOption) (*inventoryv1.ReserveStockRes, error) {
	if s.reserve == nil {
		return nil, nil
	}
	return s.reserve(ctx, req, opts...)
}

func (s orderInventoryClientStub) ConfirmReservation(ctx context.Context, req *inventoryv1.ConfirmReservationReq, opts ...grpc.CallOption) (*inventoryv1.ConfirmReservationRes, error) {
	if s.confirm == nil {
		return nil, nil
	}
	return s.confirm(ctx, req, opts...)
}

func (s orderInventoryClientStub) CancelReservation(ctx context.Context, req *inventoryv1.CancelReservationReq, opts ...grpc.CallOption) (*inventoryv1.CancelReservationRes, error) {
	if s.cancel == nil {
		return nil, nil
	}
	return s.cancel(ctx, req, opts...)
}

type buyerCartClientStub struct {
	remove func(ctx context.Context, req *cartv1.RemoveItemsReq, opts ...grpc.CallOption) (*cartv1.RemoveItemsRes, error)
}

func (buyerCartClientStub) AddItem(ctx context.Context, req *cartv1.AddItemReq, opts ...grpc.CallOption) (*cartv1.AddItemRes, error) {
	return nil, nil
}

func (buyerCartClientStub) UpdateItemQty(ctx context.Context, req *cartv1.UpdateItemQtyReq, opts ...grpc.CallOption) (*cartv1.UpdateItemQtyRes, error) {
	return nil, nil
}

func (buyerCartClientStub) ToggleItemChecked(ctx context.Context, req *cartv1.ToggleItemCheckedReq, opts ...grpc.CallOption) (*cartv1.ToggleItemCheckedRes, error) {
	return nil, nil
}

func (buyerCartClientStub) BatchToggleItems(ctx context.Context, req *cartv1.BatchToggleItemsReq, opts ...grpc.CallOption) (*cartv1.BatchToggleItemsRes, error) {
	return nil, nil
}

func (s buyerCartClientStub) RemoveItems(ctx context.Context, req *cartv1.RemoveItemsReq, opts ...grpc.CallOption) (*cartv1.RemoveItemsRes, error) {
	if s.remove != nil {
		return s.remove(ctx, req, opts...)
	}
	return &cartv1.RemoveItemsRes{}, nil
}

func (buyerCartClientStub) ClearInvalidItems(ctx context.Context, req *cartv1.ClearInvalidItemsReq, opts ...grpc.CallOption) (*cartv1.ClearInvalidItemsRes, error) {
	return nil, nil
}

func (buyerCartClientStub) GetMyCart(ctx context.Context, req *cartv1.GetMyCartReq, opts ...grpc.CallOption) (*cartv1.GetMyCartRes, error) {
	return nil, nil
}

func (buyerCartClientStub) PrepareCheckout(ctx context.Context, req *cartv1.PrepareCheckoutReq, opts ...grpc.CallOption) (*cartv1.PrepareCheckoutRes, error) {
	return nil, nil
}
