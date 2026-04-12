package runtime

import "testing"

func TestDetectGuardIntentRecognizesChineseShipmentQuestion(t *testing.T) {
	got := detectGuardIntent(UserMessage{
		ContentText: "这个订单发货了吗",
	})

	if got.IntentCode != intentCodeAfterSaleTask {
		t.Fatalf("expected shipment question to enter after-sale guard, got %q", got.IntentCode)
	}
	if got.GuardResultCode != guardResultAfterSaleTask {
		t.Fatalf("expected shipment question guard result %q, got %q", guardResultAfterSaleTask, got.GuardResultCode)
	}
}

func TestDetectTaskCodeRecognizesChineseShipmentQuestion(t *testing.T) {
	got := detectTaskCode("这个订单发货了吗", "")
	if got != taskCodeOrderStatusQuery {
		t.Fatalf("expected shipment question to map to order status task, got %q", got)
	}
}

func TestInferProblemTypeRecognizesChineseShipmentQuestion(t *testing.T) {
	got := inferProblemType("这个订单发货了吗", "")
	if got != problemTypeOrderStatus {
		t.Fatalf("expected shipment question to map to order_status problem type, got %q", got)
	}
}
