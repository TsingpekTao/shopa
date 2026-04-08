package payment

import (
	"encoding/json"
	"testing"
	"time"

	paymentv1 "github.com/TsingpekTao/shopa/payment-svc/api/v1"
)

func TestBuildGatewayCallbackReqFromFormParsesAlipayNotify(t *testing.T) {
	form := map[string]string{
		"notify_id":    "notify-20260407",
		"out_trade_no": "PAY202604070001",
		"trade_no":     "2026040722001499999999999999",
		"trade_status": "TRADE_SUCCESS",
		"total_amount": "299.00",
		"gmt_payment":  "2026-04-07 18:05:14",
		"sign":         "signed-value",
	}

	req, ok, err := buildGatewayCallbackReqFromForm(form, "trade_status=TRADE_SUCCESS&out_trade_no=PAY202604070001")
	if err != nil {
		t.Fatalf("expected alipay form payload to parse, got error: %v", err)
	}
	if !ok {
		t.Fatalf("expected alipay form payload to be detected")
	}
	if req.GetCallbackEventId() != "notify-20260407" {
		t.Fatalf("expected notify_id to become callback_event_id, got %q", req.GetCallbackEventId())
	}
	if req.GetPaymentNo() != "PAY202604070001" {
		t.Fatalf("expected out_trade_no to become payment_no, got %q", req.GetPaymentNo())
	}
	if req.GetExternalTradeNo() != "2026040722001499999999999999" {
		t.Fatalf("expected trade_no to become external_trade_no, got %q", req.GetExternalTradeNo())
	}
	if req.GetGatewayStatusCode() != "TRADE_SUCCESS" {
		t.Fatalf("expected trade_status to propagate, got %q", req.GetGatewayStatusCode())
	}
	if req.GetPayChannel() != paymentv1.PayChannel_PAY_CHANNEL_ALIPAY {
		t.Fatalf("expected alipay pay_channel, got %v", req.GetPayChannel())
	}
	if req.GetPaidAmount() != 29900 {
		t.Fatalf("expected total_amount to convert to fen, got %d", req.GetPaidAmount())
	}
	if req.GetPaidAt() == nil {
		t.Fatalf("expected gmt_payment to become paid_at")
	}
	paidAt := req.GetPaidAt().AsTime().In(alipayTimeLocation())
	if !paidAt.Equal(time.Date(2026, 4, 7, 18, 5, 14, 0, alipayTimeLocation())) {
		t.Fatalf("expected paid_at to preserve callback timestamp, got %s", paidAt.Format(time.RFC3339))
	}
	if req.GetSignature() != "signed-value" {
		t.Fatalf("expected sign to become signature, got %q", req.GetSignature())
	}
	if req.GetRawPayload() == "" {
		t.Fatalf("expected raw payload to be preserved")
	}
	var rawPayload map[string]any
	if err := json.Unmarshal([]byte(req.GetRawPayload()), &rawPayload); err != nil {
		t.Fatalf("expected raw payload to be valid json, got error: %v", err)
	}
	if rawPayload["gateway"] != "alipay" {
		t.Fatalf("expected raw payload gateway marker, got %#v", rawPayload["gateway"])
	}
	formPayload, ok := rawPayload["form"].(map[string]any)
	if !ok {
		t.Fatalf("expected raw payload form object, got %#v", rawPayload["form"])
	}
	if formPayload["out_trade_no"] != "PAY202604070001" {
		t.Fatalf("expected raw payload form to keep out_trade_no, got %#v", formPayload["out_trade_no"])
	}
}

func TestBuildGatewayCallbackReqFromFormGeneratesStableFallbackEventID(t *testing.T) {
	form := map[string]string{
		"out_trade_no": "PAY202604070002",
		"trade_no":     "ALI-TRADE-002",
		"trade_status": "TRADE_SUCCESS",
	}

	req, ok, err := buildGatewayCallbackReqFromForm(form, "")
	if err != nil {
		t.Fatalf("expected fallback event id generation to succeed, got error: %v", err)
	}
	if !ok {
		t.Fatalf("expected alipay form payload to be detected")
	}
	if req.GetCallbackEventId() != "alipay:PAY202604070002:ALI-TRADE-002:TRADE_SUCCESS" {
		t.Fatalf("expected fallback callback_event_id, got %q", req.GetCallbackEventId())
	}
}

func TestBuildGatewayCallbackReqFromFormIgnoresNonGatewayPayload(t *testing.T) {
	req, ok, err := buildGatewayCallbackReqFromForm(map[string]string{"foo": "bar"}, "")
	if err != nil {
		t.Fatalf("expected non-gateway payload to be ignored without error, got %v", err)
	}
	if ok {
		t.Fatalf("expected non-gateway payload to be ignored")
	}
	if req != nil {
		t.Fatalf("expected nil callback request for non-gateway payload")
	}
}

func TestParseAmountYuanToFen(t *testing.T) {
	got, err := parseAmountYuanToFen("299.50")
	if err != nil {
		t.Fatalf("expected amount to parse, got error: %v", err)
	}
	if got != 29950 {
		t.Fatalf("expected 299.50 to become 29950, got %d", got)
	}
}
