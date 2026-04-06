package risk

import (
	"context"

	httpv1 "github.com/TsingpekTao/shopa/risk-svc/api/risk/v1"
	pb "github.com/TsingpekTao/shopa/risk-svc/api/v1"
)

func (c *ControllerV1) PreCreateOrderCheck(ctx context.Context, req *httpv1.PreCreateOrderCheckReq) (*httpv1.PreCreateOrderCheckRes, error) {
	return c.risk.PreCreateOrderCheck(ctx, &req.PreCreateOrderCheckReq)
}

func (c *ControllerV1) PrePayCheck(ctx context.Context, req *httpv1.PrePayCheckReq) (*httpv1.PrePayCheckRes, error) {
	return c.risk.PrePayCheck(ctx, &req.PrePayCheckReq)
}

func (c *ControllerV1) PreRefundCheck(ctx context.Context, req *httpv1.PreRefundCheckReq) (*httpv1.PreRefundCheckRes, error) {
	return c.risk.PreRefundCheck(ctx, &req.PreRefundCheckReq)
}

func (c *ControllerV1) IngestRiskEvent(ctx context.Context, req *httpv1.IngestRiskEventReq) (*httpv1.IngestRiskEventRes, error) {
	return c.risk.IngestRiskEvent(ctx, &req.IngestRiskEventReq)
}

func (c *ControllerV1) RebuildUserFeatures(ctx context.Context, req *httpv1.RebuildUserFeaturesReq) (*httpv1.RebuildUserFeaturesRes, error) {
	return c.risk.RebuildUserFeatures(ctx, &req.RebuildUserFeaturesReq)
}

func (c *ControllerV1) UpsertRule(ctx context.Context, req *httpv1.UpsertRuleReq) (*httpv1.UpsertRuleRes, error) {
	return c.risk.UpsertRule(ctx, &req.UpsertRuleReq)
}

func (c *ControllerV1) EnableRule(ctx context.Context, req *httpv1.EnableRuleReq) (*httpv1.EnableRuleRes, error) {
	return c.risk.EnableRule(ctx, &req.EnableRuleReq)
}

func (c *ControllerV1) ListRiskHits(ctx context.Context, req *httpv1.ListRiskHitsReq) (*httpv1.ListRiskHitsRes, error) {
	return c.risk.ListRiskHits(ctx, &pb.ListRiskHitsReq{
		PageSize:   req.PageSize,
		NextCursor: req.NextCursor,
	})
}
