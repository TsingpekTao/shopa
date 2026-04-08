package v1

import (
	"reflect"
	"testing"
)

func TestPaymentHTTPPathsMatchGatewayRoutes(t *testing.T) {
	t.Helper()

	cases := []struct {
		name string
		req  any
		want string
	}{
		{
			name: "create-payment-intent",
			req:  CreatePaymentIntentReq{},
			want: "/v1/payment/buyer/payment-intents:create",
		},
		{
			name: "query-payment-intent",
			req:  QueryPaymentIntentReq{},
			want: "/v1/payment/buyer/payment-intents:query",
		},
		{
			name: "gateway-callback",
			req:  HandleGatewayCallbackReq{},
			want: "/v1/payment/internal/gateway:callback",
		},
	}

	for _, tc := range cases {
		got := paymentMetaPathTag(t, tc.req)
		if got != tc.want {
			t.Fatalf("%s path mismatch: want %q, got %q", tc.name, tc.want, got)
		}
	}
}

func paymentMetaPathTag(t *testing.T, req any) string {
	t.Helper()

	reqType := reflect.TypeOf(req)
	if reqType.NumField() == 0 {
		t.Fatalf("%T has no fields", req)
	}

	field := reqType.Field(0)
	got := field.Tag.Get("path")
	if got == "" {
		t.Fatalf("%T first field has no path tag", req)
	}
	return got
}
