package v1

import (
	"reflect"
	"testing"
)

func TestBuyerOrderHTTPPathsMatchGatewayRoutes(t *testing.T) {
	t.Helper()

	cases := []struct {
		name string
		req  any
		want string
	}{
		{
			name: "create-from-cart",
			req:  CreateOrderFromCartReq{},
			want: "/v1/order/buyer/orders:create-from-cart",
		},
		{
			name: "create-buy-now",
			req:  CreateOrderBuyNowReq{},
			want: "/v1/order/buyer/orders:create-buy-now",
		},
		{
			name: "request-pay",
			req:  RequestPayReq{},
			want: "/v1/order/buyer/orders:request-pay",
		},
		{
			name: "update-address",
			req:  UpdateMyOrderAddressReq{},
			want: "/v1/order/buyer/orders/{order_no}/address",
		},
		{
			name: "cancel",
			req:  CancelMyOrderReq{},
			want: "/v1/order/buyer/orders:cancel",
		},
		{
			name: "get-detail",
			req:  GetMyOrderDetailReq{},
			want: "/v1/order/buyer/orders/{order_no}",
		},
		{
			name: "list",
			req:  ListMyOrdersReq{},
			want: "/v1/order/buyer/orders",
		},
	}

	for _, tc := range cases {
		got := metaPathTag(t, tc.req)
		if got != tc.want {
			t.Fatalf("%s path mismatch: want %q, got %q", tc.name, tc.want, got)
		}
	}
}

func TestBuyerOrderHTTPMethodsMatchGatewayRoutes(t *testing.T) {
	t.Helper()

	cases := []struct {
		name string
		req  any
		want string
	}{
		{
			name: "create-from-cart",
			req:  CreateOrderFromCartReq{},
			want: "post",
		},
		{
			name: "create-buy-now",
			req:  CreateOrderBuyNowReq{},
			want: "post",
		},
		{
			name: "request-pay",
			req:  RequestPayReq{},
			want: "post",
		},
		{
			name: "update-address",
			req:  UpdateMyOrderAddressReq{},
			want: "patch",
		},
		{
			name: "cancel",
			req:  CancelMyOrderReq{},
			want: "post",
		},
		{
			name: "get-detail",
			req:  GetMyOrderDetailReq{},
			want: "get",
		},
		{
			name: "list",
			req:  ListMyOrdersReq{},
			want: "get",
		},
	}

	for _, tc := range cases {
		got := metaMethodTag(t, tc.req)
		if got != tc.want {
			t.Fatalf("%s method mismatch: want %q, got %q", tc.name, tc.want, got)
		}
	}
}

func metaPathTag(t *testing.T, req any) string {
	return metaTag(t, req, "path")
}

func metaMethodTag(t *testing.T, req any) string {
	return metaTag(t, req, "method")
}

func metaTag(t *testing.T, req any, tag string) string {
	t.Helper()

	reqType := reflect.TypeOf(req)
	if reqType.NumField() == 0 {
		t.Fatalf("%T has no fields", req)
	}

	field := reqType.Field(0)
	got := field.Tag.Get(tag)
	if got == "" {
		t.Fatalf("%T first field has no %s tag", req, tag)
	}
	return got
}
