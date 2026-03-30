package review

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/review-svc/api/review/v1"
)

type IReviewV1 interface {
	CreateReview(ctx context.Context, req *v1.CreateReviewReq) (res *v1.CreateReviewRes, err error)
	AppendReview(ctx context.Context, req *v1.AppendReviewReq) (res *v1.AppendReviewRes, err error)
	ListMyReviews(ctx context.Context, req *v1.ListMyReviewsReq) (res *v1.ListMyReviewsRes, err error)
	ListMyReviewsAlias(ctx context.Context, req *v1.ListMyReviewsAliasReq) (res *v1.ListMyReviewsAliasRes, err error)
	ReplyReview(ctx context.Context, req *v1.ReplyReviewReq) (res *v1.ReplyReviewRes, err error)
	ListSpuReviews(ctx context.Context, req *v1.ListSpuReviewsReq) (res *v1.ListSpuReviewsRes, err error)
	GetSpuRatingSummary(ctx context.Context, req *v1.GetSpuRatingSummaryReq) (res *v1.GetSpuRatingSummaryRes, err error)
	ModerateReview(ctx context.Context, req *v1.ModerateReviewReq) (res *v1.ModerateReviewRes, err error)
	RebuildSpuRatingSummary(ctx context.Context, req *v1.RebuildSpuRatingSummaryReq) (res *v1.RebuildSpuRatingSummaryRes, err error)
}
