package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/payment-svc/api/v1"
)

// IPayment 定义 payment 领域服务接口。
type IPayment interface {
	CreatePaymentIntent(ctx context.Context, req *v1.CreatePaymentIntentReq) (*v1.CreatePaymentIntentRes, error)
	QueryPaymentIntent(ctx context.Context, req *v1.QueryPaymentIntentReq) (*v1.QueryPaymentIntentRes, error)
	HandleGatewayCallback(ctx context.Context, req *v1.HandleGatewayCallbackReq) (*v1.HandleGatewayCallbackRes, error)
	ClosePaymentIntent(ctx context.Context, req *v1.ClosePaymentIntentReq) (*v1.ClosePaymentIntentRes, error)
	CreateRefundTask(ctx context.Context, req *v1.CreateRefundTaskReq) (*v1.CreateRefundTaskRes, error)
	ExecuteRefundTask(ctx context.Context, req *v1.ExecuteRefundTaskReq) (*v1.ExecuteRefundTaskRes, error)
	RunDailyReconciliation(ctx context.Context, req *v1.RunDailyReconciliationReq) (*v1.RunDailyReconciliationRes, error)
	ListReconciliationDiffs(ctx context.Context, req *v1.ListReconciliationDiffsReq) (*v1.ListReconciliationDiffsRes, error)
	ResolveReconciliationDiff(ctx context.Context, req *v1.ResolveReconciliationDiffReq) (*v1.ResolveReconciliationDiffRes, error)
	GetPaymentSnapshot(ctx context.Context, req *v1.GetPaymentSnapshotReq) (*v1.GetPaymentSnapshotRes, error)
}

var localPayment IPayment

// Payment 返回已注册的 payment 服务实现。
func Payment() IPayment {
	if localPayment == nil {
		panic("implement not found for interface IPayment, forgot register?")
	}
	return localPayment
}

// RegisterPayment 注册 payment 服务实现。
func RegisterPayment(i IPayment) {
	localPayment = i
}