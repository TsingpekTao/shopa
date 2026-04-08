package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	httpv1 "github.com/TsingpekTao/shopa/payment-svc/api/payment/v1"
	pb "github.com/TsingpekTao/shopa/payment-svc/api/v1"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *ControllerV1) CreatePaymentIntent(ctx context.Context, req *httpv1.CreatePaymentIntentReq) (*httpv1.CreatePaymentIntentRes, error) {
	return c.payment.CreatePaymentIntent(ctx, &req.CreatePaymentIntentReq)
}

func (c *ControllerV1) QueryPaymentIntent(ctx context.Context, req *httpv1.QueryPaymentIntentReq) (*httpv1.QueryPaymentIntentRes, error) {
	return c.payment.QueryPaymentIntent(ctx, &pb.QueryPaymentIntentReq{
		PaymentNo: req.PaymentNo,
		OrderNo:   req.OrderNo,
	})
}

func (c *ControllerV1) HandleGatewayCallback(ctx context.Context, req *httpv1.HandleGatewayCallbackReq) (*httpv1.HandleGatewayCallbackRes, error) {
	if r := g.RequestFromCtx(ctx); r != nil {
		form := r.GetFormMapStrStr()
		if callbackReq, ok, err := buildGatewayCallbackReqFromForm(form, r.GetBodyString()); ok {
			if err != nil {
				g.Log().Errorf(ctx, "parse payment gateway callback failed: %v", err)
				writeGatewayCallbackResponse(r, http.StatusBadRequest, "failure")
				return nil, nil
			}
			if _, err = c.payment.HandleGatewayCallback(ctx, callbackReq); err != nil {
				g.Log().Errorf(ctx, "handle payment gateway callback failed: %v", err)
				writeGatewayCallbackResponse(r, http.StatusBadRequest, "failure")
				return nil, nil
			}
			writeGatewayCallbackResponse(r, http.StatusOK, "success")
			return nil, nil
		}
	}
	return c.payment.HandleGatewayCallback(ctx, &req.HandleGatewayCallbackReq)
}

func (c *ControllerV1) ClosePaymentIntent(ctx context.Context, req *httpv1.ClosePaymentIntentReq) (*httpv1.ClosePaymentIntentRes, error) {
	return c.payment.ClosePaymentIntent(ctx, &req.ClosePaymentIntentReq)
}

func (c *ControllerV1) CreateRefundTask(ctx context.Context, req *httpv1.CreateRefundTaskReq) (*httpv1.CreateRefundTaskRes, error) {
	return c.payment.CreateRefundTask(ctx, &req.CreateRefundTaskReq)
}

func (c *ControllerV1) ExecuteRefundTask(ctx context.Context, req *httpv1.ExecuteRefundTaskReq) (*httpv1.ExecuteRefundTaskRes, error) {
	return c.payment.ExecuteRefundTask(ctx, &req.ExecuteRefundTaskReq)
}

func (c *ControllerV1) RunDailyReconciliation(ctx context.Context, req *httpv1.RunDailyReconciliationReq) (*httpv1.RunDailyReconciliationRes, error) {
	return c.payment.RunDailyReconciliation(ctx, &req.RunDailyReconciliationReq)
}

func (c *ControllerV1) ListReconciliationDiffs(ctx context.Context, req *httpv1.ListReconciliationDiffsReq) (*httpv1.ListReconciliationDiffsRes, error) {
	return c.payment.ListReconciliationDiffs(ctx, &pb.ListReconciliationDiffsReq{
		ReconTaskNo: req.ReconTaskNo,
		PageSize:    req.PageSize,
		NextCursor:  req.NextCursor,
	})
}

func (c *ControllerV1) ResolveReconciliationDiff(ctx context.Context, req *httpv1.ResolveReconciliationDiffReq) (*httpv1.ResolveReconciliationDiffRes, error) {
	return c.payment.ResolveReconciliationDiff(ctx, &req.ResolveReconciliationDiffReq)
}

func (c *ControllerV1) GetPaymentSnapshot(ctx context.Context, req *httpv1.GetPaymentSnapshotReq) (*httpv1.GetPaymentSnapshotRes, error) {
	return c.payment.GetPaymentSnapshot(ctx, &pb.GetPaymentSnapshotReq{
		PaymentNo: req.PaymentNo,
		OrderNo:   req.OrderNo,
	})
}

func buildGatewayCallbackReqFromForm(form map[string]string, rawPayload string) (*pb.HandleGatewayCallbackReq, bool, error) {
	outTradeNo := strings.TrimSpace(form["out_trade_no"])
	tradeStatus := strings.TrimSpace(form["trade_status"])
	tradeNo := strings.TrimSpace(form["trade_no"])
	notifyID := strings.TrimSpace(form["notify_id"])
	if outTradeNo == "" && tradeStatus == "" && tradeNo == "" && notifyID == "" {
		return nil, false, nil
	}
	if outTradeNo == "" {
		return nil, true, gerror.NewCode(gcode.CodeInvalidParameter, "out_trade_no is required")
	}

	paidAmount, err := parseAmountYuanToFen(form["total_amount"])
	if err != nil {
		return nil, true, err
	}
	paidAt, err := parseGatewayPaidAt(form["gmt_payment"])
	if err != nil {
		return nil, true, err
	}

	callbackEventID := notifyID
	if callbackEventID == "" {
		callbackEventID = buildFallbackCallbackEventID(outTradeNo, tradeNo, tradeStatus)
	}
	return &pb.HandleGatewayCallbackReq{
		CallbackEventId:   callbackEventID,
		ExternalTradeNo:   tradeNo,
		PaymentNo:         outTradeNo,
		PayChannel:        pb.PayChannel_PAY_CHANNEL_ALIPAY,
		GatewayStatusCode: tradeStatus,
		PaidAmount:        paidAmount,
		PaidAt:            paidAt,
		RawPayload:        buildGatewayRawPayloadJSON(form, rawPayload),
		Signature:         strings.TrimSpace(form["sign"]),
		IdempotencyKey:    fmt.Sprintf("payment-gateway-callback:%s", callbackEventID),
	}, true, nil
}

func parseAmountYuanToFen(raw string) (uint64, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, nil
	}
	parts := strings.Split(raw, ".")
	if len(parts) > 2 {
		return 0, gerror.NewCodef(gcode.CodeInvalidParameter, "invalid amount %q", raw)
	}
	whole, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return 0, gerror.Wrapf(err, "parse amount whole failed, amount=%q", raw)
	}
	fraction := "00"
	if len(parts) == 2 {
		fraction = strings.TrimSpace(parts[1])
		switch len(fraction) {
		case 0:
			fraction = "00"
		case 1:
			fraction += "0"
		case 2:
		default:
			return 0, gerror.NewCodef(gcode.CodeInvalidParameter, "amount %q has more than 2 decimals", raw)
		}
	}
	fen, err := strconv.ParseUint(fraction, 10, 64)
	if err != nil {
		return 0, gerror.Wrapf(err, "parse amount fraction failed, amount=%q", raw)
	}
	return whole*100 + fen, nil
}

func parseGatewayPaidAt(raw string) (*timestamppb.Timestamp, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	paidAt, err := time.ParseInLocation("2006-01-02 15:04:05", raw, alipayTimeLocation())
	if err != nil {
		return nil, gerror.Wrapf(err, "parse gmt_payment failed, value=%q", raw)
	}
	return timestamppb.New(paidAt), nil
}

func alipayTimeLocation() *time.Location {
	return time.FixedZone("CST", 8*3600)
}

func buildFallbackCallbackEventID(paymentNo, tradeNo, tradeStatus string) string {
	parts := []string{"alipay", strings.TrimSpace(paymentNo), strings.TrimSpace(tradeNo), strings.TrimSpace(tradeStatus)}
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		filtered = append(filtered, part)
	}
	return strings.Join(filtered, ":")
}

func buildGatewayRawPayloadJSON(form map[string]string, rawPayload string) string {
	payload := map[string]any{
		"gateway":  "alipay",
		"form":     form,
		"raw_body": strings.TrimSpace(rawPayload),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		body, _ = json.Marshal(map[string]string{
			"gateway":  "alipay",
			"raw_body": strings.TrimSpace(rawPayload),
		})
		return string(body)
	}
	return string(body)
}

func writeGatewayCallbackResponse(r *ghttp.Request, status int, body string) {
	if r == nil {
		return
	}
	r.Response.Status = status
	r.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	r.Response.WriteOver(body)
}
