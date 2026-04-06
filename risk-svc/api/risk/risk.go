package risk

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/risk-svc/api/risk/v1"
)

type IRiskV1 interface {
	PreCreateOrderCheck(ctx context.Context, req *v1.PreCreateOrderCheckReq) (res *v1.PreCreateOrderCheckRes, err error)
	PrePayCheck(ctx context.Context, req *v1.PrePayCheckReq) (res *v1.PrePayCheckRes, err error)
	PreRefundCheck(ctx context.Context, req *v1.PreRefundCheckReq) (res *v1.PreRefundCheckRes, err error)
	IngestRiskEvent(ctx context.Context, req *v1.IngestRiskEventReq) (res *v1.IngestRiskEventRes, err error)
	RebuildUserFeatures(ctx context.Context, req *v1.RebuildUserFeaturesReq) (res *v1.RebuildUserFeaturesRes, err error)
	UpsertRule(ctx context.Context, req *v1.UpsertRuleReq) (res *v1.UpsertRuleRes, err error)
	EnableRule(ctx context.Context, req *v1.EnableRuleReq) (res *v1.EnableRuleRes, err error)
	ListRiskHits(ctx context.Context, req *v1.ListRiskHitsReq) (res *v1.ListRiskHitsRes, err error)
}
