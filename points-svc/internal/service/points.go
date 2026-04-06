// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	pb "github.com/TsingpekTao/shopa/points-svc/api/v1"
)

type (
	IPoints interface {
		PreviewOrderPointsDeduction(ctx context.Context, req *pb.PreviewOrderPointsDeductionReq) (*pb.PreviewOrderPointsDeductionRes, error)
		LockPointsForOrder(ctx context.Context, req *pb.LockPointsForOrderReq) (*pb.LockPointsForOrderRes, error)
		ConfirmLockedPoints(ctx context.Context, req *pb.ConfirmLockedPointsReq) (*pb.ConfirmLockedPointsRes, error)
		CancelLockedPoints(ctx context.Context, req *pb.CancelLockedPointsReq) (*pb.CancelLockedPointsRes, error)
		GrantPointsForCompletedOrder(ctx context.Context, req *pb.GrantPointsForCompletedOrderReq) (*pb.GrantPointsForCompletedOrderRes, error)
		ReturnPointsByRefund(ctx context.Context, req *pb.ReturnPointsByRefundReq) (*pb.ReturnPointsByRefundRes, error)
		ReverseGrantedPointsByRefund(ctx context.Context, req *pb.ReverseGrantedPointsByRefundReq) (*pb.ReverseGrantedPointsByRefundRes, error)
		InitPointsAccountIfAbsent(ctx context.Context, req *pb.InitPointsAccountIfAbsentReq) (*pb.InitPointsAccountIfAbsentRes, error)
		GetPointsByUserId(ctx context.Context, req *pb.GetPointsByUserIdReq) (*pb.GetPointsByUserIdRes, error)
		InitAccountFromRegisterEvent(ctx context.Context, userID uint64) error
	}
)

var (
	localPoints IPoints
)

func Points() IPoints {
	if localPoints == nil {
		panic("implement not found for interface IPoints, forgot register?")
	}
	return localPoints
}

func RegisterPoints(i IPoints) {
	localPoints = i
}
