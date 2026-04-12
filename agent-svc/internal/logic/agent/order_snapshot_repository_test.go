package agent

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	agentruntime "github.com/TsingpekTao/shopa/agent-svc/internal/logic/agent/runtime"
)

func TestBuyerOrderSnapshotHTTPRepositoryUsesBuyerScopedEndpoint(t *testing.T) {
	var (
		gotPath   string
		gotUserID string
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotUserID = r.Header.Get("X-User-Id")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":0,"message":"ok","data":{"order":{"order_no":"ORD202604090301","order_status":"PAID","payment_status":"PAID","delivery_status":"UNSHIPPED","after_sale_status":"NONE","updated_at":"2026-04-09T12:00:00+08:00"}}}`))
	}))
	defer server.Close()

	repo := newBuyerOrderSnapshotHTTPRepository(server.URL, server.Client())
	got, err := repo.QueryOrderSnapshot(context.Background(), agentruntime.OrderOwnershipFilter{
		UserID:         90001,
		OrderNo:        "ORD202604090301",
		RequestID:      "REQ202604090301",
		ConversationNo: "ACV202604090301",
		RunNo:          "ARN202604090301",
	})
	if err != nil {
		t.Fatalf("query order snapshot returned error: %v", err)
	}

	if gotPath != "/v1/order/buyer/orders/ORD202604090301" {
		t.Fatalf("expected buyer scoped order detail path, got %q", gotPath)
	}
	if gotUserID != "90001" {
		t.Fatalf("expected X-User-Id header 90001, got %q", gotUserID)
	}
	if got == nil {
		t.Fatalf("expected snapshot to be returned")
	}
	if !got.OwnershipConfirmed {
		t.Fatalf("expected ownership to be confirmed for buyer scoped endpoint")
	}
	if got.OrderNo != "ORD202604090301" {
		t.Fatalf("expected order_no to be parsed, got %q", got.OrderNo)
	}
	if got.FulfillmentStatus != "UNSHIPPED" {
		t.Fatalf("expected fulfillment status to be parsed, got %q", got.FulfillmentStatus)
	}
}

func TestBuyerOrderSnapshotHTTPRepositoryListsRecentOrders(t *testing.T) {
	var (
		gotListPath   string
		detailLookups int
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/order/buyer/orders":
			gotListPath = r.URL.RequestURI()
			_, _ = w.Write([]byte(`{"code":0,"message":"ok","data":{"orders":[{"order_no":"ORD202604100401","order_status":"PAID","payment_status":"PAID","updated_at":"2026-04-10T09:30:00+08:00","sub_orders":[{"items":[{"spu_title":"云感运动鞋"}]}]},{"order_no":"ORD202604100402","order_status":"PAID","payment_status":"PAID","updated_at":"2026-04-10T08:30:00+08:00","sub_orders":[{"items":[{"spu_title":"亚麻衬衫"}]}]}]}}`))
		case "/v1/order/buyer/orders/ORD202604100401":
			detailLookups++
			_, _ = w.Write([]byte(`{"code":0,"message":"ok","data":{"order":{"order_no":"ORD202604100401","delivery_status":"IN_TRANSIT","after_sale_status":"NONE","updated_at":"2026-04-10T09:30:00+08:00"}}}`))
		case "/v1/order/buyer/orders/ORD202604100402":
			detailLookups++
			_, _ = w.Write([]byte(`{"code":0,"message":"ok","data":{"order":{"order_no":"ORD202604100402","delivery_status":"UNSHIPPED","after_sale_status":"NONE","updated_at":"2026-04-10T08:30:00+08:00"}}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	repo := newBuyerOrderSnapshotHTTPRepository(server.URL, server.Client())
	got, err := repo.ListRecentOrders(context.Background(), agentruntime.RecentOrderListFilter{
		UserID: 90002,
		Limit:  2,
	})
	if err != nil {
		t.Fatalf("list recent orders returned error: %v", err)
	}

	if gotListPath != "/v1/order/buyer/orders?page_size=2" {
		t.Fatalf("expected recent order path with page_size, got %q", gotListPath)
	}
	if len(got) != 2 {
		t.Fatalf("expected two recent orders, got %d", len(got))
	}
	if detailLookups != 2 {
		t.Fatalf("expected two detail lookups to enrich recent orders, got %d", detailLookups)
	}
	if got[0].DisplayTitle != "云感运动鞋" {
		t.Fatalf("expected first recent order title to be parsed, got %q", got[0].DisplayTitle)
	}
	if got[0].LogisticsStatus != "IN_TRANSIT" {
		t.Fatalf("expected detail snapshot to enrich logistics status, got %q", got[0].LogisticsStatus)
	}
}
