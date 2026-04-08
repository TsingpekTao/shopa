package proxy

import (
	"net/http"
	"testing"
)

func TestMatchBuiltinRoute_OrderSellerRoutes(t *testing.T) {
	svc := &sProxy{}

	testCases := []struct {
		name   string
		method string
		path   string
		want   bool
	}{
		{
			name:   "seller list path is proxied",
			method: http.MethodPost,
			path:   "/v1/order/seller/list",
			want:   true,
		},
		{
			name:   "seller sub detail path is proxied",
			method: http.MethodGet,
			path:   "/v1/order/seller/sub/SUB20260407192826493100",
			want:   true,
		},
		{
			name:   "seller ship path is proxied",
			method: http.MethodPost,
			path:   "/v1/order/seller/sub/ship",
			want:   true,
		},
		{
			name:   "legacy seller list path is not used anymore",
			method: http.MethodGet,
			path:   "/v1/order/seller/shops/SHOP1001/orders",
			want:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := svc.matchBuiltinRoute(tc.method, tc.path)
			if (got != nil) != tc.want {
				t.Fatalf("matchBuiltinRoute(%s, %s) = %v, want route existence %v", tc.method, tc.path, got != nil, tc.want)
			}
		})
	}
}

func TestMatchBuiltinRoute_OrderBuyerConfirmReceiptRoute(t *testing.T) {
	svc := &sProxy{}

	route := svc.matchBuiltinRoute(http.MethodPost, "/v1/order/buyer/orders:confirm-received")
	if route == nil {
		t.Fatalf("expected buyer confirm receipt route to be proxied")
	}
	if route.UpstreamPathTemplate != "/v1/order/buyer/orders:confirm-received" {
		t.Fatalf("unexpected upstream path template: %s", route.UpstreamPathTemplate)
	}
}
