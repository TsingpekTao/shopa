package points

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/points-svc/api/points/v1"
)

type IPointsV1 interface {
	InitPointsAccountIfAbsent(ctx context.Context, req *v1.InitPointsAccountIfAbsentReq) (res *v1.InitPointsAccountIfAbsentRes, err error)
	GetPointsByUserId(ctx context.Context, req *v1.GetPointsByUserIdReq) (res *v1.GetPointsByUserIdRes, err error)
	PreviewOrder(ctx context.Context, req *v1.PreviewOrderReq) (res *v1.PreviewOrderRes, err error)
	LockOrder(ctx context.Context, req *v1.LockOrderReq) (res *v1.LockOrderRes, err error)
	CancelOrder(ctx context.Context, req *v1.CancelOrderReq) (res *v1.SimpleAckRes, err error)
	ConfirmOrder(ctx context.Context, req *v1.ConfirmOrderReq) (res *v1.SimpleAckRes, err error)
	GrantOrder(ctx context.Context, req *v1.GrantOrderReq) (res *v1.SimpleAckRes, err error)
	ReturnRefund(ctx context.Context, req *v1.ReturnRefundReq) (res *v1.ReturnRefundRes, err error)
	ReverseRefund(ctx context.Context, req *v1.ReverseRefundReq) (res *v1.ReverseRefundRes, err error)
}
