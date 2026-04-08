package order

import (
	"context"
	"strings"
	"testing"
	"time"

	inventoryv1 "github.com/TsingpekTao/shopa/inventory-svc/api/v1"
	"github.com/TsingpekTao/shopa/order-svc/internal/model/entity"
	v1 "github.com/TsingpekTao/shopa/order-svc/api/v1"
	"google.golang.org/grpc"
)

func TestSummarizeRefundOrderStateAllRefunded(t *testing.T) {
	subStatuses := []v1.SubOrderStatus{
		v1.SubOrderStatus_SUB_ORDER_STATUS_REFUNDED,
		v1.SubOrderStatus_SUB_ORDER_STATUS_REFUNDED,
	}

	orderStatus, paymentStatus := summarizeRefundOrderState(subStatuses)
	if orderStatus != v1.OrderStatus_ORDER_STATUS_REFUNDED {
		t.Fatalf("expected refunded order status, got %v", orderStatus)
	}
	if paymentStatus != v1.PaymentStatus_PAYMENT_STATUS_REFUNDED {
		t.Fatalf("expected refunded payment status, got %v", paymentStatus)
	}
}

func TestSummarizeRefundOrderStateWithRefundingSubOrder(t *testing.T) {
	subStatuses := []v1.SubOrderStatus{
		v1.SubOrderStatus_SUB_ORDER_STATUS_REFUNDING,
		v1.SubOrderStatus_SUB_ORDER_STATUS_WAIT_SHIP,
	}

	orderStatus, paymentStatus := summarizeRefundOrderState(subStatuses)
	if orderStatus != v1.OrderStatus_ORDER_STATUS_REFUNDING {
		t.Fatalf("expected refunding order status, got %v", orderStatus)
	}
	if paymentStatus != v1.PaymentStatus_PAYMENT_STATUS_PAID {
		t.Fatalf("expected paid payment status, got %v", paymentStatus)
	}
}

func TestCalculateRefundReversePointsUsesRuleSnapshot(t *testing.T) {
	points, err := calculateRefundReversePoints(`{"grant_points_per_cent":1}`, 12900)
	if err != nil {
		t.Fatalf("expected calculateRefundReversePoints to succeed, got error: %v", err)
	}
	if points != 129 {
		t.Fatalf("expected 129 reverse points, got %d", points)
	}
}

func TestCalculateRefundReversePointsBlankSnapshotReturnsZero(t *testing.T) {
	points, err := calculateRefundReversePoints("", 12900)
	if err != nil {
		t.Fatalf("expected blank snapshot to be accepted, got error: %v", err)
	}
	if points != 0 {
		t.Fatalf("expected zero reverse points for blank snapshot, got %d", points)
	}
}

func TestRestoreRefundInventoryAdjustsStock(t *testing.T) {
	ctx := context.Background()
	s := &sOrder{}
	stub := &mockInventoryClient{t: t}
	oldClient := newOrderInventoryAdminClient
	newOrderInventoryAdminClient = func(ctx context.Context) (orderInventoryAdminClient, error) {
		return stub, nil
	}
	defer func() { newOrderInventoryAdminClient = oldClient }()

	items := []*entity.OrderItem{
		{
			SkuNo:  "SKU123",
			SpuNo:  "SPU123",
			ShopNo: "SHOP1",
			Qty:    3,
		},
	}
	if err := s.restoreRefundInventory(ctx, "REFUND-1", items); err != nil {
		t.Fatalf("restoreRefundInventory failed: %v", err)
	}
	if len(stub.req.Items) != 1 {
		t.Fatalf("expected 1 adjust item, got %d", len(stub.req.Items))
	}
	item := stub.req.Items[0]
	if item.SkuNo != "SKU123" || item.DeltaTotalQty != 3 {
		t.Fatalf("unexpected adjust item: %+v", item)
	}
	if !strings.Contains(item.BizNo, "REFUND-1") {
		t.Fatalf("expected biz no to include refund identifier, got %s", item.BizNo)
	}
}

func TestReturnRefundPointsInvokesPointsService(t *testing.T) {
	ctx := context.Background()
	restoreConf := configurePointsIntegrationForTest(pointsIntegrationConf{
		Enabled:     true,
		BaseURL:     "http://points",
		Timeout:     3 * time.Second,
	})
	defer restoreConf()
	called := false
	var captured pointsRefundReturnRequest
	doPointsJSONRequest = func(ctx context.Context, conf pointsIntegrationConf, path string, reqBody any, respBody any) error {
		called = true
		captured = *(reqBody.(*pointsRefundReturnRequest))
		if resp, ok := respBody.(*pointsRefundReturnResponse); ok {
			resp.Success = true
			resp.PointsReturnAmount = captured.PointsToReturn
		}
		return nil
	}
	defer func() { doPointsJSONRequest = defaultDoPointsJSONRequest }()

	s := &sOrder{}
	mainRow := &entity.OrderMain{OrderNo: "ORD-1", UserId: 123}
	subRow := &entity.OrderSub{SubOrderNo: "SUB-1", ShopNo: "SHOP-1"}
	if _, err := s.returnRefundPoints(ctx, mainRow, subRow, "REFUND-A", 29900, 50); err != nil {
		t.Fatalf("returnRefundPoints failed: %v", err)
	}
	if !called {
		t.Fatalf("expected points return request to execute")
	}
	if captured.PointsToReturn != 50 {
		t.Fatalf("expected points to return 50, got %d", captured.PointsToReturn)
	}
	if path := captured.OrderNo; path != "ORD-1" {
		t.Fatalf("unexpected order number %s", path)
	}
}

func TestReverseRefundGrantedPointsInvokesPointsService(t *testing.T) {
	ctx := context.Background()
	restoreConf := configurePointsIntegrationForTest(pointsIntegrationConf{
		Enabled:     true,
		BaseURL:     "http://points",
		Timeout:     3 * time.Second,
	})
	defer restoreConf()
	called := false
	var captured pointsRefundReverseRequest
	doPointsJSONRequest = func(ctx context.Context, conf pointsIntegrationConf, path string, reqBody any, respBody any) error {
		called = true
		captured = *(reqBody.(*pointsRefundReverseRequest))
		if resp, ok := respBody.(*pointsRefundReverseResponse); ok {
			resp.Success = true
			resp.PointsCashOffsetAmount = 123
		}
		return nil
	}
	defer func() { doPointsJSONRequest = defaultDoPointsJSONRequest }()

	s := &sOrder{}
	mainRow := &entity.OrderMain{OrderNo: "ORD-2", UserId: 123}
	subRow := &entity.OrderSub{SubOrderNo: "SUB-2", ShopNo: "SHOP-2"}
	if _, err := s.reverseRefundGrantedPoints(ctx, mainRow, subRow, "REFUND-B", 12900, 10); err != nil {
		t.Fatalf("reverseRefundGrantedPoints failed: %v", err)
	}
	if !called {
		t.Fatalf("expected points reverse request to execute")
	}
	if captured.PointsToReverse != 10 {
		t.Fatalf("expected points to reverse 10, got %d", captured.PointsToReverse)
	}
}

type mockInventoryClient struct {
	t   *testing.T
	req *inventoryv1.BatchAdjustStockByAdminReq
}

func (m *mockInventoryClient) BatchAdjustStockByAdmin(ctx context.Context, req *inventoryv1.BatchAdjustStockByAdminReq, opts ...grpc.CallOption) (*inventoryv1.BatchAdjustStockByAdminRes, error) {
	if len(req.GetItems()) == 0 {
		m.t.Fatalf("expected adjust items, got none")
	}
	m.req = req
	return &inventoryv1.BatchAdjustStockByAdminRes{}, nil
}

func configurePointsIntegrationForTest(conf pointsIntegrationConf) func() {
	original := loadPointsIntegrationConf
	loadPointsIntegrationConf = func(ctx context.Context) pointsIntegrationConf {
		return conf
	}
	return func() {
		loadPointsIntegrationConf = original
	}
}
