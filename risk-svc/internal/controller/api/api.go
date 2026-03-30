package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/risk-svc/api/v1"
	"github.com/TsingpekTao/shopa/risk-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

// Controller 聚合 risk 所有 gRPC 接口。
type Controller struct {
	v1.UnimplementedRiskDecisionServiceServer
	v1.UnimplementedInternalRiskServiceServer
	v1.UnimplementedAdminRiskServiceServer
}

// Register 注册 risk gRPC 服务。
func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{}
	v1.RegisterRiskDecisionServiceServer(s.Server, ctrl)
	v1.RegisterInternalRiskServiceServer(s.Server, ctrl)
	v1.RegisterAdminRiskServiceServer(s.Server, ctrl)
}

// PreCreateOrderCheck 下单前检查。
func (c *Controller) PreCreateOrderCheck(ctx context.Context, req *v1.PreCreateOrderCheckReq) (*v1.RiskDecisionRes, error) {
	return service.Risk().PreCreateOrderCheck(ctx, req)
}

// PrePayCheck 支付前检查。
func (c *Controller) PrePayCheck(ctx context.Context, req *v1.PrePayCheckReq) (*v1.RiskDecisionRes, error) {
	return service.Risk().PrePayCheck(ctx, req)
}

// PreRefundCheck 退款前检查。
func (c *Controller) PreRefundCheck(ctx context.Context, req *v1.PreRefundCheckReq) (*v1.RiskDecisionRes, error) {
	return service.Risk().PreRefundCheck(ctx, req)
}

// IngestRiskEvent 摄入异步风险事件。
func (c *Controller) IngestRiskEvent(ctx context.Context, req *v1.IngestRiskEventReq) (*v1.IngestRiskEventRes, error) {
	return service.Risk().IngestRiskEvent(ctx, req)
}

// RebuildUserFeatures 重建用户特征。
func (c *Controller) RebuildUserFeatures(ctx context.Context, req *v1.RebuildUserFeaturesReq) (*v1.RebuildUserFeaturesRes, error) {
	return service.Risk().RebuildUserFeatures(ctx, req)
}

// UpsertRule 新增或更新规则。
func (c *Controller) UpsertRule(ctx context.Context, req *v1.UpsertRuleReq) (*v1.UpsertRuleRes, error) {
	return service.Risk().UpsertRule(ctx, req)
}

// EnableRule 启停规则。
func (c *Controller) EnableRule(ctx context.Context, req *v1.EnableRuleReq) (*v1.EnableRuleRes, error) {
	return service.Risk().EnableRule(ctx, req)
}

// ListRiskHits 查询命中记录。
func (c *Controller) ListRiskHits(ctx context.Context, req *v1.ListRiskHitsReq) (*v1.ListRiskHitsRes, error) {
	return service.Risk().ListRiskHits(ctx, req)
}
