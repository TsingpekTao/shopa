package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/payment-svc/api/v1"
	"github.com/TsingpekTao/shopa/payment-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

// Controller 汇聚 payment 的 gRPC 接口实现。
type Controller struct {
	v1.UnimplementedBuyerPaymentServiceServer
	v1.UnimplementedInternalPaymentServiceServer
}

// Register 向 gRPC 服务器注册 payment 服务。
func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{}
	v1.RegisterBuyerPaymentServiceServer(s.Server, ctrl)
	v1.RegisterInternalPaymentServiceServer(s.Server, ctrl)
}

// CreatePaymentIntent 创建支付意图。
func (*Controller) CreatePaymentIntent(ctx context.Context, req *v1.CreatePaymentIntentReq) (*v1.CreatePaymentIntentRes, error) {
	return service.Payment().CreatePaymentIntent(ctx, req)
}

// QueryPaymentIntent 查询支付意图。
func (*Controller) QueryPaymentIntent(ctx context.Context, req *v1.QueryPaymentIntentReq) (*v1.QueryPaymentIntentRes, error) {
	return service.Payment().QueryPaymentIntent(ctx, req)
}

// HandleGatewayCallback 处理网关回调。
func (*Controller) HandleGatewayCallback(ctx context.Context, req *v1.HandleGatewayCallbackReq) (*v1.HandleGatewayCallbackRes, error) {
	return service.Payment().HandleGatewayCallback(ctx, req)
}

// ClosePaymentIntent 关闭支付意图。
func (*Controller) ClosePaymentIntent(ctx context.Context, req *v1.ClosePaymentIntentReq) (*v1.ClosePaymentIntentRes, error) {
	return service.Payment().ClosePaymentIntent(ctx, req)
}

// CreateRefundTask 创建退款任务。
func (*Controller) CreateRefundTask(ctx context.Context, req *v1.CreateRefundTaskReq) (*v1.CreateRefundTaskRes, error) {
	return service.Payment().CreateRefundTask(ctx, req)
}

// ExecuteRefundTask 执行退款任务。
func (*Controller) ExecuteRefundTask(ctx context.Context, req *v1.ExecuteRefundTaskReq) (*v1.ExecuteRefundTaskRes, error) {
	return service.Payment().ExecuteRefundTask(ctx, req)
}

// RunDailyReconciliation 触发每日对账。
func (*Controller) RunDailyReconciliation(ctx context.Context, req *v1.RunDailyReconciliationReq) (*v1.RunDailyReconciliationRes, error) {
	return service.Payment().RunDailyReconciliation(ctx, req)
}

// ListReconciliationDiffs 查询对账差异。
func (*Controller) ListReconciliationDiffs(ctx context.Context, req *v1.ListReconciliationDiffsReq) (*v1.ListReconciliationDiffsRes, error) {
	return service.Payment().ListReconciliationDiffs(ctx, req)
}

// ResolveReconciliationDiff 处理对账差异。
func (*Controller) ResolveReconciliationDiff(ctx context.Context, req *v1.ResolveReconciliationDiffReq) (*v1.ResolveReconciliationDiffRes, error) {
	return service.Payment().ResolveReconciliationDiff(ctx, req)
}

// GetPaymentSnapshot 返回支付快照。
func (*Controller) GetPaymentSnapshot(ctx context.Context, req *v1.GetPaymentSnapshotReq) (*v1.GetPaymentSnapshotRes, error) {
	return service.Payment().GetPaymentSnapshot(ctx, req)
}