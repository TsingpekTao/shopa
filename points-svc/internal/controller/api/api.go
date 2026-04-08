package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/points-svc/api/v1"
	"github.com/TsingpekTao/shopa/points-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

type Controller struct {
	v1.UnimplementedPointsServiceServer
	points service.IPoints
}

func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{points: service.Points()}
	v1.RegisterPointsServiceServer(s.Server, ctrl)
}

func (c *Controller) InitPointsAccountIfAbsent(ctx context.Context, req *v1.InitPointsAccountIfAbsentReq) (*v1.InitPointsAccountIfAbsentRes, error) {
	return c.points.InitPointsAccountIfAbsent(ctx, req)
}

func (c *Controller) GetPointsByUserId(ctx context.Context, req *v1.GetPointsByUserIdReq) (*v1.GetPointsByUserIdRes, error) {
	return c.points.GetPointsByUserId(ctx, req)
}

func (*Controller) GetPointsAccountSummary(ctx context.Context, req *v1.GetPointsAccountSummaryReq) (res *v1.GetPointsAccountSummaryRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) ListMyPointsBills(ctx context.Context, req *v1.ListMyPointsBillsReq) (res *v1.ListMyPointsBillsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) PreviewOrderPointsDeduction(ctx context.Context, req *v1.PreviewOrderPointsDeductionReq) (res *v1.PreviewOrderPointsDeductionRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) LockPointsForOrder(ctx context.Context, req *v1.LockPointsForOrderReq) (res *v1.LockPointsForOrderRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) ConfirmLockedPoints(ctx context.Context, req *v1.ConfirmLockedPointsReq) (res *v1.ConfirmLockedPointsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) CancelLockedPoints(ctx context.Context, req *v1.CancelLockedPointsReq) (res *v1.CancelLockedPointsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) GrantPointsForCompletedOrder(ctx context.Context, req *v1.GrantPointsForCompletedOrderReq) (res *v1.GrantPointsForCompletedOrderRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) ReturnPointsByRefund(ctx context.Context, req *v1.ReturnPointsByRefundReq) (res *v1.ReturnPointsByRefundRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) ReverseGrantedPointsByRefund(ctx context.Context, req *v1.ReverseGrantedPointsByRefundReq) (res *v1.ReverseGrantedPointsByRefundRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) ExpirePointsBatch(ctx context.Context, req *v1.ExpirePointsBatchReq) (res *v1.ExpirePointsBatchRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) GetPointsAccount(ctx context.Context, req *v1.GetPointsAccountReq) (res *v1.GetPointsAccountRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) ListPointsAccounts(ctx context.Context, req *v1.ListPointsAccountsReq) (res *v1.ListPointsAccountsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) ListReservationRecords(ctx context.Context, req *v1.ListReservationRecordsReq) (res *v1.ListReservationRecordsRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) AdjustPointsBalance(ctx context.Context, req *v1.AdjustPointsBalanceReq) (res *v1.AdjustPointsBalanceRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) FreezePointsAccount(ctx context.Context, req *v1.FreezePointsAccountReq) (res *v1.FreezePointsAccountRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) UnfreezePointsAccount(ctx context.Context, req *v1.UnfreezePointsAccountReq) (res *v1.UnfreezePointsAccountRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}

func (*Controller) TriggerExpireSweep(ctx context.Context, req *v1.TriggerExpireSweepReq) (res *v1.TriggerExpireSweepRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
