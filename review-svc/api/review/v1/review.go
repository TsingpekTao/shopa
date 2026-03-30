package v1

import (
	pb "github.com/TsingpekTao/shopa/review-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type CreateReviewReq struct {
	g.Meta `path:"/v1/review/buyer/reviews:create" method:"post" tags:"Review-Buyer" summary:"Create review"`
	pb.CreateReviewReq
}

type CreateReviewRes = pb.CreateReviewRes

type AppendReviewReq struct {
	g.Meta `path:"/v1/review/buyer/reviews:append" method:"post" tags:"Review-Buyer" summary:"Append review"`
	pb.AppendReviewReq
}

type AppendReviewRes = pb.AppendReviewRes

type ListMyReviewsReq struct {
	g.Meta `path:"/v1/review/buyer/reviews:list" method:"post" tags:"Review-Buyer" summary:"List my reviews"`
	pb.ListMyReviewsReq
}

type ListMyReviewsRes = pb.ListMyReviewsRes

type ListMyReviewsAliasReq struct {
	g.Meta     `path:"/v1/review/buyer/reviews" method:"get" tags:"Review-Buyer" summary:"List my reviews (alias)"`
	PageSize   int32             `json:"page_size"`
	NextCursor string            `json:"next_cursor"`
	Statuses   []pb.ReviewStatus `json:"statuses"`
}

type ListMyReviewsAliasRes = pb.ListMyReviewsRes

type ReplyReviewReq struct {
	g.Meta `path:"/v1/review/seller/reviews:reply" method:"post" tags:"Review-Seller" summary:"Reply review"`
	pb.ReplyReviewReq
}

type ReplyReviewRes = pb.ReplyReviewRes

type ListSpuReviewsReq struct {
	g.Meta        `path:"/v1/review/public/spu/{spu_no}/reviews" method:"get" tags:"Review-Public" summary:"List spu reviews"`
	SpuNo         string            `json:"spu_no" v:"required#spu_no is required"`
	PageSize      int32             `json:"page_size"`
	NextCursor    string            `json:"next_cursor"`
	SortCode      pb.ReviewSortCode `json:"sort_code"`
	WithMediaOnly bool              `json:"with_media_only"`
}

type ListSpuReviewsRes = pb.ListSpuReviewsRes

type GetSpuRatingSummaryReq struct {
	g.Meta `path:"/v1/review/public/spu/{spu_no}/summary" method:"get" tags:"Review-Public" summary:"Get spu rating summary"`
	SpuNo  string `json:"spu_no" v:"required#spu_no is required"`
}

type GetSpuRatingSummaryRes = pb.GetSpuRatingSummaryRes

type ModerateReviewReq struct {
	g.Meta `path:"/v1/review/internal/moderate" method:"post" tags:"Review-Internal" summary:"Moderate review"`
	pb.ModerateReviewReq
}

type ModerateReviewRes = pb.ModerateReviewRes

type RebuildSpuRatingSummaryReq struct {
	g.Meta `path:"/v1/review/internal/rebuild-summary" method:"post" tags:"Review-Internal" summary:"Rebuild spu rating summary"`
	pb.RebuildSpuRatingSummaryReq
}

type RebuildSpuRatingSummaryRes = pb.RebuildSpuRatingSummaryRes
