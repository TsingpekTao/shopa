package payment

import (
	"context"

	httpv1 "github.com/TsingpekTao/shopa/payment-svc/api/payment/v1"
	pb "github.com/TsingpekTao/shopa/payment-svc/api/v1"
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
