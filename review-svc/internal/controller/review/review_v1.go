package review

import (
	"context"

	reviewv1 "github.com/TsingpekTao/shopa/review-svc/api/review/v1"
	pb "github.com/TsingpekTao/shopa/review-svc/api/v1"
)

func (c *ControllerV1) CreateReview(ctx context.Context, req *reviewv1.CreateReviewReq) (*reviewv1.CreateReviewRes, error) {
	return c.review.CreateReview(ctx, &req.CreateReviewReq)
}

func (c *ControllerV1) AppendReview(ctx context.Context, req *reviewv1.AppendReviewReq) (*reviewv1.AppendReviewRes, error) {
	return c.review.AppendReview(ctx, &req.AppendReviewReq)
}

func (c *ControllerV1) ListMyReviews(ctx context.Context, req *reviewv1.ListMyReviewsReq) (*reviewv1.ListMyReviewsRes, error) {
	return c.review.ListMyReviews(ctx, &req.ListMyReviewsReq)
}

func (c *ControllerV1) ListMyReviewsAlias(ctx context.Context, req *reviewv1.ListMyReviewsAliasReq) (*reviewv1.ListMyReviewsAliasRes, error) {
	return c.review.ListMyReviews(ctx, &pb.ListMyReviewsReq{
		PageSize:   req.PageSize,
		NextCursor: req.NextCursor,
		Statuses:   req.Statuses,
	})
}

func (c *ControllerV1) ReplyReview(ctx context.Context, req *reviewv1.ReplyReviewReq) (*reviewv1.ReplyReviewRes, error) {
	return c.review.ReplyReview(ctx, &req.ReplyReviewReq)
}

func (c *ControllerV1) ListSpuReviews(ctx context.Context, req *reviewv1.ListSpuReviewsReq) (*reviewv1.ListSpuReviewsRes, error) {
	return c.review.ListSpuReviews(ctx, &pb.ListSpuReviewsReq{
		SpuNo:         req.SpuNo,
		PageSize:      req.PageSize,
		NextCursor:    req.NextCursor,
		SortCode:      req.SortCode,
		WithMediaOnly: req.WithMediaOnly,
	})
}

func (c *ControllerV1) GetSpuRatingSummary(ctx context.Context, req *reviewv1.GetSpuRatingSummaryReq) (*reviewv1.GetSpuRatingSummaryRes, error) {
	return c.review.GetSpuRatingSummary(ctx, &pb.GetSpuRatingSummaryReq{SpuNo: req.SpuNo})
}

func (c *ControllerV1) ModerateReview(ctx context.Context, req *reviewv1.ModerateReviewReq) (*reviewv1.ModerateReviewRes, error) {
	return c.review.ModerateReview(ctx, &req.ModerateReviewReq)
}

func (c *ControllerV1) RebuildSpuRatingSummary(ctx context.Context, req *reviewv1.RebuildSpuRatingSummaryReq) (*reviewv1.RebuildSpuRatingSummaryRes, error) {
	return c.review.RebuildSpuRatingSummary(ctx, &req.RebuildSpuRatingSummaryReq)
}
