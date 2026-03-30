// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/risk-svc/api/v1"
)

type (
	IRisk interface {
		// PreCreateOrderCheck 执行下单前风控检查。
		PreCreateOrderCheck(ctx context.Context, req *v1.PreCreateOrderCheckReq) (*v1.RiskDecisionRes, error)
		// PrePayCheck 执行支付前风控检查。
		PrePayCheck(ctx context.Context, req *v1.PrePayCheckReq) (*v1.RiskDecisionRes, error)
		// PreRefundCheck 执行退款前风控检查。
		PreRefundCheck(ctx context.Context, req *v1.PreRefundCheckReq) (*v1.RiskDecisionRes, error)
		// IngestRiskEvent 写入风险事件用于异步特征预热。
		IngestRiskEvent(ctx context.Context, req *v1.IngestRiskEventReq) (*v1.IngestRiskEventRes, error)
		// RebuildUserFeatures 重建指定用户的风险特征快照。
		RebuildUserFeatures(ctx context.Context, req *v1.RebuildUserFeaturesReq) (*v1.RebuildUserFeaturesRes, error)
		// UpsertRule 创建或更新风控规则。
		UpsertRule(ctx context.Context, req *v1.UpsertRuleReq) (*v1.UpsertRuleRes, error)
		// EnableRule 切换规则启用状态。
		EnableRule(ctx context.Context, req *v1.EnableRuleReq) (*v1.EnableRuleRes, error)
		// ListRiskHits 分页查询命中记录。
		ListRiskHits(ctx context.Context, req *v1.ListRiskHitsReq) (*v1.ListRiskHitsRes, error)
	}
)

var (
	localRisk IRisk
)

func Risk() IRisk {
	if localRisk == nil {
		panic("implement not found for interface IRisk, forgot register?")
	}
	return localRisk
}

func RegisterRisk(i IRisk) {
	localRisk = i
}
