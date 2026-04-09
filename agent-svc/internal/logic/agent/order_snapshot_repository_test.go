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
