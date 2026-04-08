package payment

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"encoding/pem"
	"net/url"
	"strings"
	"testing"
	"time"

	orderv1 "github.com/TsingpekTao/shopa/order-svc/api/v1"
	v1 "github.com/TsingpekTao/shopa/payment-svc/api/v1"
	"github.com/TsingpekTao/shopa/payment-svc/internal/model/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestIsIdempotencyMissErrorTreatsSQLNoRowsAsMiss(t *testing.T) {
	if !isIdempotencyMissError(sql.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows to be treated as idempotency miss")
	}
}

func TestBuildPaymentLaunchForAlipayGeneratesSandboxURL(t *testing.T) {
	privateKeyPEM := mustMakePKCS8PrivateKeyPEM(t)
	intent := &v1.PaymentIntent{
		PaymentNo:       "PAY202604060001",
		OrderNo:         "ORD202604060001",
		PayChannel:      v1.PayChannel_PAY_CHANNEL_ALIPAY,
		PayableAmount:   1999,
		GatewayExpireAt: tsPtr(time.Date(2026, 4, 6, 14, 0, 0, 0, time.FixedZone("CST", 8*3600))),
	}
	req := &v1.CreatePaymentIntentReq{
		OrderNo:   "ORD202604060001",
		Subject:   "Shopa sandbox order",
		ReturnUrl: "http://localhost:3100/pay/result",
		NotifyUrl: "https://example.ngrok.app/v1/payment/internal/gateway:callback",
	}
	cfg := alipayConfig{
		GatewayURL:    "https://openapi-sandbox.dl.alipaydev.com/gateway.do",
		AppID:         "2021000000000001",
		PrivateKey:    privateKeyPEM,
		ReturnURL:     "http://localhost:3100/default-return",
		NotifyURL:     "https://example.ngrok.app/default-notify",
		SubjectPrefix: "Shopa",
	}

	launch, err := buildAlipayPagePaymentLaunch(cfg, req, intent)
	if err != nil {
		t.Fatalf("expected alipay launch to build, got error: %v", err)
	}
	if !strings.HasPrefix(launch.PayURL, "https://openapi-sandbox.dl.alipaydev.com/gateway.do?") {
		t.Fatalf("expected sandbox gateway url, got %q", launch.PayURL)
	}

	parsed, err := url.Parse(launch.PayURL)
	if err != nil {
		t.Fatalf("expected pay_url to be parseable, got error: %v", err)
	}
	query := parsed.Query()
	if query.Get("app_id") != "2021000000000001" {
		t.Fatalf("expected app_id in pay_url, got %q", query.Get("app_id"))
	}
	if query.Get("method") != "alipay.trade.page.pay" {
		t.Fatalf("expected alipay page pay method, got %q", query.Get("method"))
	}
	if query.Get("sign_type") != "RSA2" {
		t.Fatalf("expected RSA2 sign_type, got %q", query.Get("sign_type"))
	}
	if query.Get("return_url") != "http://localhost:3100/pay/result" {
		t.Fatalf("expected request return_url to win, got %q", query.Get("return_url"))
	}
	if query.Get("notify_url") != "https://example.ngrok.app/v1/payment/internal/gateway:callback" {
		t.Fatalf("expected request notify_url to win, got %q", query.Get("notify_url"))
	}
	if query.Get("sign") == "" {
		t.Fatalf("expected signed pay_url")
	}

	var bizContent map[string]string
	if err := json.Unmarshal([]byte(query.Get("biz_content")), &bizContent); err != nil {
		t.Fatalf("expected biz_content json, got error: %v", err)
	}
	if bizContent["time_expire"] != "2026-04-06 14:00:00" {
		t.Fatalf("expected alipay time_expire in CST, got %q", bizContent["time_expire"])
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(launch.PayPayloadJSON), &payload); err != nil {
		t.Fatalf("expected pay payload json to be valid, got error: %v", err)
	}
	if payload["gateway"] != "alipay" {
		t.Fatalf("expected alipay gateway payload, got %#v", payload["gateway"])
	}
	if payload["url"] != launch.PayURL {
		t.Fatalf("expected payload url to mirror pay_url")
	}
}

func TestBuildPaymentLaunchForAlipayRequiresConfig(t *testing.T) {
	intent := &v1.PaymentIntent{
		PaymentNo:     "PAY202604060002",
		OrderNo:       "ORD202604060002",
		PayChannel:    v1.PayChannel_PAY_CHANNEL_ALIPAY,
		PayableAmount: 1999,
	}

	_, err := buildAlipayPagePaymentLaunch(alipayConfig{}, &v1.CreatePaymentIntentReq{OrderNo: "ORD202604060002"}, intent)
	if err == nil {
		t.Fatalf("expected missing config to fail")
	}
	if !strings.Contains(err.Error(), "alipay") {
		t.Fatalf("expected alipay config error, got %v", err)
	}
}

func TestCfgEnvStringPrefersEnvironmentVariable(t *testing.T) {
	const envKey = "SHOPA_TEST_ALIPAY_APP_ID"
	t.Setenv(envKey, "2021000000009999")

	got := cfgEnvString(context.Background(), "payment.alipay.appId", "", envKey)
	if got != "2021000000009999" {
		t.Fatalf("expected env override, got %q", got)
	}
}

func TestFormatAlipayTimestampUsesChinaTimezone(t *testing.T) {
	ts := timestamppb.New(time.Date(2026, 4, 7, 13, 37, 44, 0, time.FixedZone("CST", 8*3600)))
	if got := formatAlipayTimestamp(ts); got != "2026-04-07 13:37:44" {
		t.Fatalf("expected CST formatted timestamp, got %q", got)
	}
}

func TestCanReuseExistingIntentForCreate(t *testing.T) {
	tests := []struct {
		name   string
		status v1.PaymentIntentStatus
		want   bool
	}{
		{name: "created", status: v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_CREATED, want: true},
		{name: "paying", status: v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_PAYING, want: true},
		{name: "paid", status: v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_PAID, want: false},
		{name: "closed", status: v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_CLOSED, want: false},
		{name: "failed", status: v1.PaymentIntentStatus_PAYMENT_INTENT_STATUS_FAILED, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			row := &entity.PaymentIntent{Status: uint(tc.status)}
			if got := canReuseExistingIntentForCreate(row); got != tc.want {
				t.Fatalf("expected reusable=%v, got %v", tc.want, got)
			}
		})
	}
}

func TestBuildOrderPayCallbackRequestMapsGatewayFields(t *testing.T) {
	req := &v1.HandleGatewayCallbackReq{
		CallbackEventId:   "cb-1",
		PaymentNo:         "PAY202604070001",
		OrderNo:           "ORD202604070001",
		PayChannel:        v1.PayChannel_PAY_CHANNEL_ALIPAY,
		GatewayStatusCode: "TRADE_SUCCESS",
		ExternalTradeNo:   "ALI-TRADE-001",
		PaidAmount:        29900,
		PaidAt:            timestamppb.New(time.Date(2026, 4, 7, 15, 30, 0, 0, time.FixedZone("CST", 8*3600))),
		RawPayload:        "{\"trade_status\":\"TRADE_SUCCESS\"}",
	}
	intent := &v1.PaymentIntent{
		PaymentNo: "PAY202604070001",
		OrderNo:   "ORD202604070001",
	}

	got := buildOrderPayCallbackRequest(req, intent)
	if got == nil {
		t.Fatalf("expected callback request to be built")
	}
	if got.GetPaymentEventId() != "cb-1" {
		t.Fatalf("expected callback_event_id to map to payment_event_id, got %q", got.GetPaymentEventId())
	}
	if got.GetPayNo() != "PAY202604070001" || got.GetOrderNo() != "ORD202604070001" {
		t.Fatalf("expected order/payment numbers to map, got pay=%q order=%q", got.GetPayNo(), got.GetOrderNo())
	}
	if got.GetPayChannel() != orderv1.PayChannel_PAY_CHANNEL_ALIPAY {
		t.Fatalf("expected pay channel to map to order proto, got %v", got.GetPayChannel())
	}
	if got.GetPayStatusCode() != "TRADE_SUCCESS" {
		t.Fatalf("expected gateway status to propagate, got %q", got.GetPayStatusCode())
	}
	if got.GetChannelTradeNo() != "ALI-TRADE-001" {
		t.Fatalf("expected external trade no to propagate, got %q", got.GetChannelTradeNo())
	}
	if got.GetPaidAmount() != 29900 {
		t.Fatalf("expected paid amount to propagate, got %d", got.GetPaidAmount())
	}
	if got.GetIdempotencyKey() != "payment-callback:cb-1" {
		t.Fatalf("expected derived idempotency key, got %q", got.GetIdempotencyKey())
	}
}

func mustMakePKCS8PrivateKeyPEM(t *testing.T) string {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatalf("generate rsa key: %v", err)
	}
	body, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("marshal pkcs8 key: %v", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: body}))
}

func tsPtr(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}
