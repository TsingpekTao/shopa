package payment

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/payment-svc/api/payment/v1"
)

type IPaymentV1 interface {
	CreatePaymentIntent(ctx context.Context, req *v1.CreatePaymentIntentReq) (res *v1.CreatePaymentIntentRes, err error)
	QueryPaymentIntent(ctx context.Context, req *v1.QueryPaymentIntentReq) (res *v1.QueryPaymentIntentRes, err error)
	HandleGatewayCallback(ctx context.Context, req *v1.HandleGatewayCallbackReq) (res *v1.HandleGatewayCallbackRes, err error)
	ClosePaymentIntent(ctx context.Context, req *v1.ClosePaymentIntentReq) (res *v1.ClosePaymentIntentRes, err error)
	CreateRefundTask(ctx context.Context, req *v1.CreateRefundTaskReq) (res *v1.CreateRefundTaskRes, err error)
	ExecuteRefundTask(ctx context.Context, req *v1.ExecuteRefundTaskReq) (res *v1.ExecuteRefundTaskRes, err error)
	RunDailyReconciliation(ctx context.Context, req *v1.RunDailyReconciliationReq) (res *v1.RunDailyReconciliationRes, err error)
	ListReconciliationDiffs(ctx context.Context, req *v1.ListReconciliationDiffsReq) (res *v1.ListReconciliationDiffsRes, err error)
	ResolveReconciliationDiff(ctx context.Context, req *v1.ResolveReconciliationDiffReq) (res *v1.ResolveReconciliationDiffRes, err error)
	GetPaymentSnapshot(ctx context.Context, req *v1.GetPaymentSnapshotReq) (res *v1.GetPaymentSnapshotRes, err error)
}
