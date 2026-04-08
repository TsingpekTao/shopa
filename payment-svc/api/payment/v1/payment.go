package v1

import (
	pb "github.com/TsingpekTao/shopa/payment-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type CreatePaymentIntentReq struct {
	g.Meta `path:"/v1/payment/buyer/payment-intents:create" method:"post" tags:"Payment" summary:"Create payment intent"`
	pb.CreatePaymentIntentReq
}

type CreatePaymentIntentRes = pb.CreatePaymentIntentRes

type QueryPaymentIntentReq struct {
	g.Meta    `path:"/v1/payment/buyer/payment-intents:query" method:"get" tags:"Payment" summary:"Query payment intent"`
	PaymentNo string `json:"paymentNo" in:"query"`
	OrderNo   string `json:"orderNo" in:"query"`
}

type QueryPaymentIntentRes = pb.QueryPaymentIntentRes

type HandleGatewayCallbackReq struct {
	g.Meta `path:"/v1/payment/internal/gateway:callback" method:"post" tags:"PaymentInternal" summary:"Handle payment gateway callback"`
	pb.HandleGatewayCallbackReq
}

type HandleGatewayCallbackRes = pb.HandleGatewayCallbackRes

type ClosePaymentIntentReq struct {
	g.Meta `path:"/v1/internal/payments/close" method:"post" tags:"PaymentInternal" summary:"Close payment intent"`
	pb.ClosePaymentIntentReq
}

type ClosePaymentIntentRes = pb.ClosePaymentIntentRes

type CreateRefundTaskReq struct {
	g.Meta `path:"/v1/internal/payments/refund-tasks" method:"post" tags:"PaymentInternal" summary:"Create refund task"`
	pb.CreateRefundTaskReq
}

type CreateRefundTaskRes = pb.CreateRefundTaskRes

type ExecuteRefundTaskReq struct {
	g.Meta `path:"/v1/internal/payments/refund-tasks/execute" method:"post" tags:"PaymentInternal" summary:"Execute refund task"`
	pb.ExecuteRefundTaskReq
}

type ExecuteRefundTaskRes = pb.ExecuteRefundTaskRes

type RunDailyReconciliationReq struct {
	g.Meta `path:"/v1/internal/payments/reconciliation/run" method:"post" tags:"PaymentInternal" summary:"Run daily reconciliation"`
	pb.RunDailyReconciliationReq
}

type RunDailyReconciliationRes = pb.RunDailyReconciliationRes

type ListReconciliationDiffsReq struct {
	g.Meta      `path:"/v1/internal/payments/reconciliation/diffs" method:"get" tags:"PaymentInternal" summary:"List reconciliation diffs"`
	ReconTaskNo string `json:"reconTaskNo" in:"query" v:"required#reconTaskNo is required"`
	PageSize    int32  `json:"pageSize" in:"query"`
	NextCursor  string `json:"nextCursor" in:"query"`
}

type ListReconciliationDiffsRes = pb.ListReconciliationDiffsRes

type ResolveReconciliationDiffReq struct {
	g.Meta `path:"/v1/internal/payments/reconciliation/diffs/resolve" method:"post" tags:"PaymentInternal" summary:"Resolve reconciliation diff"`
	pb.ResolveReconciliationDiffReq
}

type ResolveReconciliationDiffRes = pb.ResolveReconciliationDiffRes

type GetPaymentSnapshotReq struct {
	g.Meta    `path:"/v1/internal/payments/snapshot" method:"get" tags:"PaymentInternal" summary:"Get payment snapshot"`
	PaymentNo string `json:"paymentNo" in:"query"`
	OrderNo   string `json:"orderNo" in:"query"`
}

type GetPaymentSnapshotRes = pb.GetPaymentSnapshotRes
