package bff

import "testing"

func TestBuyerRefundScopeLabel(t *testing.T) {
	if buyerRefundScopeLabel != "当前按该商品所属整笔子单退款处理" {
		t.Fatalf("unexpected buyer refund scope label: %q", buyerRefundScopeLabel)
	}
}
