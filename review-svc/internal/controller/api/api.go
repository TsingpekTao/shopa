package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/review-svc/api/v1"
	"github.com/TsingpekTao/shopa/review-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedBuyerReviewServiceServer
	v1.UnimplementedSellerReviewServiceServer
	v1.UnimplementedPublicReviewServiceServer
	v1.UnimplementedInternalReviewServiceServer
}

func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{}
	v1.RegisterBuyerReviewServiceServer(s.Server, ctrl)
	v1.RegisterSellerReviewServiceServer(s.Server, ctrl)
	v1.RegisterPublicReviewServiceServer(s.Server, ctrl)
	v1.RegisterInternalReviewServiceServer(s.Server, ctrl)
}

func (*Controller) CreateReview(ctx context.Context, req *v1.CreateReviewReq) (res *v1.CreateReviewRes, err error) {
	return service.Review().CreateReview(ctx, req)
}

func (*Controller) AppendReview(ctx context.Context, req *v1.AppendReviewReq) (res *v1.AppendReviewRes, err error) {
	return service.Review().AppendReview(ctx, req)
}

func (*Controller) ListMyReviews(ctx context.Context, req *v1.ListMyReviewsReq) (res *v1.ListMyReviewsRes, err error) {
	return service.Review().ListMyReviews(ctx, req)
}

func (*Controller) ReplyReview(ctx context.Context, req *v1.ReplyReviewReq) (res *v1.ReplyReviewRes, err error) {
	return service.Review().ReplyReview(ctx, req)
}

func (*Controller) ListSpuReviews(ctx context.Context, req *v1.ListSpuReviewsReq) (res *v1.ListSpuReviewsRes, err error) {
	return service.Review().ListSpuReviews(ctx, req)
}

func (*Controller) GetSpuRatingSummary(ctx context.Context, req *v1.GetSpuRatingSummaryReq) (res *v1.GetSpuRatingSummaryRes, err error) {
	return service.Review().GetSpuRatingSummary(ctx, req)
}

func (*Controller) ModerateReview(ctx context.Context, req *v1.ModerateReviewReq) (res *v1.ModerateReviewRes, err error) {
	return service.Review().ModerateReview(ctx, req)
}

func (*Controller) RebuildSpuRatingSummary(ctx context.Context, req *v1.RebuildSpuRatingSummaryReq) (res *v1.RebuildSpuRatingSummaryRes, err error) {
	return service.Review().RebuildSpuRatingSummary(ctx, req)
}
