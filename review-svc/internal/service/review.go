// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/review-svc/api/v1"
)

type (
	IReview interface {
		CreateReview(ctx context.Context, req *v1.CreateReviewReq) (*v1.CreateReviewRes, error)
		AppendReview(ctx context.Context, req *v1.AppendReviewReq) (*v1.AppendReviewRes, error)
		ListMyReviews(ctx context.Context, req *v1.ListMyReviewsReq) (*v1.ListMyReviewsRes, error)
		ReplyReview(ctx context.Context, req *v1.ReplyReviewReq) (*v1.ReplyReviewRes, error)
		ListSpuReviews(ctx context.Context, req *v1.ListSpuReviewsReq) (*v1.ListSpuReviewsRes, error)
		GetSpuRatingSummary(ctx context.Context, req *v1.GetSpuRatingSummaryReq) (*v1.GetSpuRatingSummaryRes, error)
		ModerateReview(ctx context.Context, req *v1.ModerateReviewReq) (*v1.ModerateReviewRes, error)
		RebuildSpuRatingSummary(ctx context.Context, req *v1.RebuildSpuRatingSummaryReq) (*v1.RebuildSpuRatingSummaryRes, error)
	}
)

var (
	localReview IReview
)

func Review() IReview {
	if localReview == nil {
		panic("implement not found for interface IReview, forgot register?")
	}
	return localReview
}

func RegisterReview(i IReview) {
	localReview = i
}
