package v1

import (
	pb "github.com/TsingpekTao/shopa/risk-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type PreCreateOrderCheckReq struct {
	g.Meta `path:"/v1/risk/decisions/pre-create-order" method:"post" tags:"Risk" summary:"Run pre-create-order risk check"`
	pb.PreCreateOrderCheckReq
}

type PreCreateOrderCheckRes = pb.RiskDecisionRes

type PrePayCheckReq struct {
	g.Meta `path:"/v1/risk/decisions/pre-pay" method:"post" tags:"Risk" summary:"Run pre-pay risk check"`
	pb.PrePayCheckReq
}

type PrePayCheckRes = pb.RiskDecisionRes

type PreRefundCheckReq struct {
	g.Meta `path:"/v1/risk/decisions/pre-refund" method:"post" tags:"Risk" summary:"Run pre-refund risk check"`
	pb.PreRefundCheckReq
}

type PreRefundCheckRes = pb.RiskDecisionRes

type IngestRiskEventReq struct {
	g.Meta `path:"/v1/internal/risk/events" method:"post" tags:"RiskInternal" summary:"Ingest risk event"`
	pb.IngestRiskEventReq
}

type IngestRiskEventRes = pb.IngestRiskEventRes

type RebuildUserFeaturesReq struct {
	g.Meta `path:"/v1/internal/risk/features/rebuild" method:"post" tags:"RiskInternal" summary:"Rebuild user risk features"`
	pb.RebuildUserFeaturesReq
}

type RebuildUserFeaturesRes = pb.RebuildUserFeaturesRes

type UpsertRuleReq struct {
	g.Meta `path:"/v1/admin/risk/rules/upsert" method:"post" tags:"RiskAdmin" summary:"Create or update risk rule"`
	pb.UpsertRuleReq
}

type UpsertRuleRes = pb.UpsertRuleRes

type EnableRuleReq struct {
	g.Meta `path:"/v1/admin/risk/rules/enable" method:"post" tags:"RiskAdmin" summary:"Enable or disable risk rule"`
	pb.EnableRuleReq
}

type EnableRuleRes = pb.EnableRuleRes

type ListRiskHitsReq struct {
	g.Meta     `path:"/v1/admin/risk/hits" method:"get" tags:"RiskAdmin" summary:"List risk hit logs"`
	PageSize   int32  `json:"pageSize" in:"query"`
	NextCursor string `json:"nextCursor" in:"query"`
}

type ListRiskHitsRes = pb.ListRiskHitsRes
