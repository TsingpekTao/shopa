package risk

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	v1 "github.com/TsingpekTao/shopa/risk-svc/api/v1"
	"github.com/TsingpekTao/shopa/risk-svc/internal/dao"
	"github.com/TsingpekTao/shopa/risk-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/risk-svc/internal/model/entity"
	"github.com/TsingpekTao/shopa/risk-svc/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// sRisk 是 risk 领域逻辑实现。
type sRisk struct{}

// New 创建 risk 逻辑实例。
func New() *sRisk {
	// 返回无状态逻辑对象。
	return &sRisk{}
}

func init() {
	// 注册 risk 逻辑到 service 层。
	service.RegisterRisk(New())
}

// PreCreateOrderCheck 执行下单前风控检查。
func (s *sRisk) PreCreateOrderCheck(ctx context.Context, req *v1.PreCreateOrderCheckReq) (*v1.RiskDecisionRes, error) {
	// 校验请求不能为空。
	if req == nil || req.GetUserId() == 0 {
		// 返回参数错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "user_id is required")
	}
	// 计算基础风控结果。
	res := s.computeDecision(req.GetPayableAmount(), req.GetContainsVirtualItem(), req.GetCategoryNos())
	// 记录风控决策日志。
	_ = s.saveDecisionLog(ctx, req.GetUserId(), "PRE_CREATE_ORDER", req.GetOrderNo(), res, false)
	// 返回风控结果。
	return res, nil
}

// PrePayCheck 执行支付前风控检查。
func (s *sRisk) PrePayCheck(ctx context.Context, req *v1.PrePayCheckReq) (*v1.RiskDecisionRes, error) {
	// 校验请求不能为空。
	if req == nil || req.GetUserId() == 0 {
		// 返回参数错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "user_id is required")
	}
	// 计算风控决策。
	res := s.computeDecision(req.GetPayableAmount(), false, nil)
	// 记录决策日志。
	_ = s.saveDecisionLog(ctx, req.GetUserId(), "PRE_PAY", req.GetPaymentNo(), res, false)
	// 返回风控结果。
	return res, nil
}

// PreRefundCheck 执行退款前风控检查。
func (s *sRisk) PreRefundCheck(ctx context.Context, req *v1.PreRefundCheckReq) (*v1.RiskDecisionRes, error) {
	// 校验请求不能为空。
	if req == nil || req.GetUserId() == 0 {
		// 返回参数错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "user_id is required")
	}
	// 计算风控决策。
	res := s.computeDecision(req.GetRefundAmount(), false, nil)
	// 记录决策日志。
	_ = s.saveDecisionLog(ctx, req.GetUserId(), "PRE_REFUND", req.GetAfterSaleNo(), res, false)
	// 返回风控结果。
	return res, nil
}

// IngestRiskEvent 写入风险事件用于异步特征预热。
func (s *sRisk) IngestRiskEvent(ctx context.Context, req *v1.IngestRiskEventReq) (*v1.IngestRiskEventRes, error) {
	// 校验请求不能为空。
	if req == nil || strings.TrimSpace(req.GetEventId()) == "" {
		// 返回参数错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "event_id is required")
	}
	// 序列化事件发生时间字符串。
	occurredAt := gtime.Now()
	// 如果传了 occurred_at 则使用传入时间。
	if req.GetOccurredAt() != nil {
		// 把 protobuf 时间转换为 gtime。
		occurredAt = gtime.NewFromTime(req.GetOccurredAt().AsTime())
	}
	// 插入事件收件箱，重复事件直接忽略。
	if _, err := dao.RiskEventInbox.Ctx(ctx).Data(do.RiskEventInbox{
		EventId:     req.GetEventId(),     // 写入事件ID。
		EventType:   req.GetEventType(),   // 写入事件类型。
		UserId:      req.GetUserId(),      // 写入用户ID。
		BizNo:       req.GetBizNo(),       // 写入业务号。
		PayloadJson: req.GetPayloadJson(), // 写入事件载荷。
		OccurredAt:  occurredAt,           // 写入发生时间。
		CreatedAt:   gtime.Now(),          // 写入创建时间。
	}).InsertIgnore(); err != nil {
		// 返回入库错误。
		return nil, gerror.Wrap(err, "insert risk_event_inbox failed")
	}
	// 返回已接受。
	return &v1.IngestRiskEventRes{Accepted: true}, nil
}

// RebuildUserFeatures 重建指定用户的风险特征快照。
func (s *sRisk) RebuildUserFeatures(ctx context.Context, req *v1.RebuildUserFeaturesReq) (*v1.RebuildUserFeaturesRes, error) {
	// 校验用户ID。
	if req == nil || req.GetUserId() == 0 {
		// 返回参数错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "user_id is required")
	}
	// 统计该用户事件数量。
	eventCount, err := dao.RiskEventInbox.Ctx(ctx).Where(dao.RiskEventInbox.Columns().UserId, req.GetUserId()).Count()
	// 处理统计异常。
	if err != nil {
		// 返回统计错误。
		return nil, gerror.Wrap(err, "count risk_event_inbox failed")
	}
	// 按事件数量粗略计算风险分。
	score := uint(eventCount * 2)
	// 风险分做上限保护。
	if score > 100 {
		// 超过上限时截断为100。
		score = 100
	}
	// 组装特征JSON。
	featureJSON, _ := json.Marshal(map[string]any{
		"event_count": eventCount,                      // 写入事件数量。
		"rebuild_at":  time.Now().Format(time.RFC3339), // 写入重建时间。
	})
	// 组装标签JSON。
	tagJSON, _ := json.Marshal([]string{"AUTO_REBUILT"})
	// 先尝试更新已有快照。
	res, err := dao.RiskFeatureSnapshot.Ctx(ctx).
		Where(dao.RiskFeatureSnapshot.Columns().UserId, req.GetUserId()).
		Data(do.RiskFeatureSnapshot{
			RiskScore:    score,               // 更新风险分。
			TagsJson:     string(tagJSON),     // 更新标签。
			FeaturesJson: string(featureJSON), // 更新特征。
			UpdatedAt:    gtime.Now(),         // 更新时间。
		}).Update()
	// 处理更新异常。
	if err != nil {
		// 返回更新错误。
		return nil, gerror.Wrap(err, "update risk_feature_snapshot failed")
	}
	// 获取更新影响行数。
	rows, _ := res.RowsAffected()
	// 如果没有更新到记录则插入新快照。
	if rows == 0 {
		// 插入新快照记录。
		if _, err := dao.RiskFeatureSnapshot.Ctx(ctx).Data(do.RiskFeatureSnapshot{
			UserId:       req.GetUserId(),     // 写入用户ID。
			RiskScore:    score,               // 写入风险分。
			TagsJson:     string(tagJSON),     // 写入标签。
			FeaturesJson: string(featureJSON), // 写入特征。
			UpdatedAt:    gtime.Now(),         // 写入更新时间。
		}).Insert(); err != nil {
			// 返回插入错误。
			return nil, gerror.Wrap(err, "insert risk_feature_snapshot failed")
		}
	}
	// 返回重建结果。
	return &v1.RebuildUserFeaturesRes{
		UserId:    req.GetUserId(),                   // 返回用户ID。
		UpdatedAt: timestamppbFromGTime(gtime.Now()), // 返回更新时间。
	}, nil
}

// UpsertRule 创建或更新风控规则。
func (s *sRisk) UpsertRule(ctx context.Context, req *v1.UpsertRuleReq) (*v1.UpsertRuleRes, error) {
	// 校验规则编码。
	if req == nil || strings.TrimSpace(req.GetRuleCode()) == "" {
		// 返回参数错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "rule_code is required")
	}
	// 先更新已有规则。
	res, err := dao.RiskRule.Ctx(ctx).
		Where(dao.RiskRule.Columns().RuleCode, req.GetRuleCode()).
		Data(do.RiskRule{
			RuleName:     req.GetRuleName(),           // 更新规则名。
			RuleExprJson: req.GetRuleExprJson(),       // 更新表达式。
			Priority:     req.GetPriority(),           // 更新优先级。
			Enabled:      boolToInt(req.GetEnabled()), // 更新启用状态。
			UpdatedAt:    gtime.Now(),                 // 更新时间。
		}).Update()
	// 处理更新异常。
	if err != nil {
		// 返回更新错误。
		return nil, gerror.Wrap(err, "update risk_rule failed")
	}
	// 获取影响行数。
	rows, _ := res.RowsAffected()
	// 如果规则不存在则插入新规则。
	if rows == 0 {
		// 插入新规则。
		if _, err := dao.RiskRule.Ctx(ctx).Data(do.RiskRule{
			RuleCode:     req.GetRuleCode(),           // 写入编码。
			RuleName:     req.GetRuleName(),           // 写入名称。
			RuleExprJson: req.GetRuleExprJson(),       // 写入表达式。
			Priority:     req.GetPriority(),           // 写入优先级。
			Enabled:      boolToInt(req.GetEnabled()), // 写入启用状态。
			CreatedAt:    gtime.Now(),                 // 创建时间。
			UpdatedAt:    gtime.Now(),                 // 更新时间。
		}).Insert(); err != nil {
			// 返回插入错误。
			return nil, gerror.Wrap(err, "insert risk_rule failed")
		}
	}
	// 写入规则版本记录便于追踪。
	_, _ = dao.RiskRuleVersion.Ctx(ctx).Data(do.RiskRuleVersion{
		RuleCode:     req.GetRuleCode(),           // 写入规则编码。
		RuleExprJson: req.GetRuleExprJson(),       // 写入规则表达式。
		Enabled:      boolToInt(req.GetEnabled()), // 写入启用状态。
		CreatedAt:    gtime.Now(),                 // 写入创建时间。
	}).Insert()
	// 返回更新结果。
	return &v1.UpsertRuleRes{RuleCode: req.GetRuleCode()}, nil
}

// EnableRule 切换规则启用状态。
func (s *sRisk) EnableRule(ctx context.Context, req *v1.EnableRuleReq) (*v1.EnableRuleRes, error) {
	// 校验参数。
	if req == nil || strings.TrimSpace(req.GetRuleCode()) == "" {
		// 返回参数错误。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "rule_code is required")
	}
	// 更新规则启用状态。
	if _, err := dao.RiskRule.Ctx(ctx).
		Where(dao.RiskRule.Columns().RuleCode, req.GetRuleCode()).
		Data(do.RiskRule{
			Enabled:   boolToInt(req.GetEnabled()), // 更新启用状态。
			UpdatedAt: gtime.Now(),                 // 更新时间。
		}).Update(); err != nil {
		// 返回更新错误。
		return nil, gerror.Wrap(err, "enable rule failed")
	}
	// 返回切换结果。
	return &v1.EnableRuleRes{
		RuleCode: req.GetRuleCode(), // 返回规则编码。
		Enabled:  req.GetEnabled(),  // 返回启用状态。
	}, nil
}

// ListRiskHits 分页查询命中记录。
func (s *sRisk) ListRiskHits(ctx context.Context, req *v1.ListRiskHitsReq) (*v1.ListRiskHitsRes, error) {
	// 计算分页大小。
	limit := int(req.GetPageSize())
	// 设置分页默认值。
	if limit <= 0 || limit > 100 {
		// 使用默认分页大小。
		limit = 20
	}
	// 查询命中日志。
	var rows []entity.RiskHitLog
	// 按ID倒序分页查询。
	if err := dao.RiskHitLog.Ctx(ctx).OrderDesc(dao.RiskHitLog.Columns().Id).Limit(limit).Scan(&rows); err != nil {
		// 返回查询错误。
		return nil, gerror.Wrap(err, "query risk_hit_log failed")
	}
	// 组装响应列表。
	list := make([]*v1.RiskHit, 0, len(rows))
	// 遍历实体记录。
	for _, row := range rows {
		// 追加一条响应记录。
		list = append(list, &v1.RiskHit{
			HitNo:     row.HitNo,                           // 写入命中号。
			UserId:    row.UserId,                          // 写入用户ID。
			RuleCode:  row.RuleCode,                        // 写入规则编码。
			BizType:   row.BizType,                         // 写入业务类型。
			BizNo:     row.BizNo,                           // 写入业务号。
			Decision:  v1.RiskDecision(row.Decision),       // 写入决策。
			CreatedAt: timestamppbFromGTime(row.CreatedAt), // 写入创建时间。
		})
	}
	// 返回分页结果。
	return &v1.ListRiskHitsRes{
		List:       list,  // 返回命中列表。
		NextCursor: "",    // 基础版不返回游标。
		HasMore:    false, // 基础版不计算更多页。
	}, nil
}

// computeDecision 根据上下文计算基础风险决策。
func (s *sRisk) computeDecision(amount uint64, containsVirtual bool, categories []string) *v1.RiskDecisionRes {
	// 初始化默认决策为允许。
	decision := v1.RiskDecision_RISK_DECISION_ALLOW
	// 初始化风险分为0。
	score := uint32(0)
	// 初始化原因列表。
	reasons := make([]string, 0)
	// 虚拟商品提高基础风险分。
	if containsVirtual {
		// 累加风险分。
		score += 40
		// 追加原因标签。
		reasons = append(reasons, "CONTAINS_VIRTUAL_ITEM")
	}
	// 大额订单提高风险分。
	if amount >= 300000 {
		// 累加风险分。
		score += 35
		// 追加原因标签。
		reasons = append(reasons, "HIGH_AMOUNT")
	}
	// 判断是否命中高危类目。
	for _, c := range categories {
		// 统一处理类目编码大小写。
		code := strings.ToUpper(strings.TrimSpace(c))
		// 检测高危类目关键字。
		if strings.Contains(code, "VIRTUAL") || strings.Contains(code, "CARD") {
			// 累加风险分。
			score += 30
			// 追加原因标签。
			reasons = append(reasons, "HIGH_RISK_CATEGORY")
			// 命中后即可退出循环。
			break
		}
	}
	// 根据风险分映射决策。
	switch {
	case score >= 80:
		// 高风险直接拒绝。
		decision = v1.RiskDecision_RISK_DECISION_REJECT
	case score >= 50:
		// 中风险进入人工复核。
		decision = v1.RiskDecision_RISK_DECISION_REVIEW
	default:
		// 低风险允许放行。
		decision = v1.RiskDecision_RISK_DECISION_ALLOW
	}
	// 返回风控响应。
	return &v1.RiskDecisionRes{
		Decision:         decision,   // 返回决策结果。
		RiskScore:        score,      // 返回风险分。
		MatchedRuleCodes: []string{}, // 基础版先不返回规则命中。
		ReasonCodes:      reasons,    // 返回原因编码。
		Degraded:         false,      // 基础版非降级。
		DegradedFields:   []string{}, // 基础版无降级字段。
	}
}

// saveDecisionLog 保存风控决策日志与命中日志。
func (s *sRisk) saveDecisionLog(ctx context.Context, userID uint64, bizType, bizNo string, res *v1.RiskDecisionRes, degraded bool) error {
	// 生成决策流水号。
	decisionNo := "RDL_" + strings.ToUpper(uuid.NewString())
	// 序列化命中规则数组。
	matched, _ := json.Marshal(res.GetMatchedRuleCodes())
	// 序列化原因数组。
	reasons, _ := json.Marshal(res.GetReasonCodes())
	// 写入决策日志。
	if _, err := dao.RiskDecisionLog.Ctx(ctx).Data(do.RiskDecisionLog{
		DecisionNo:       decisionNo,              // 写入决策号。
		UserId:           userID,                  // 写入用户ID。
		BizType:          bizType,                 // 写入业务类型。
		BizNo:            bizNo,                   // 写入业务号。
		Decision:         uint(res.GetDecision()), // 写入决策值。
		RiskScore:        res.GetRiskScore(),      // 写入风险分。
		MatchedRuleCodes: string(matched),         // 写入命中规则JSON。
		ReasonCodes:      string(reasons),         // 写入原因JSON。
		Degraded:         boolToInt(degraded),     // 写入降级标记。
		CreatedAt:        gtime.Now(),             // 写入创建时间。
	}).Insert(); err != nil {
		// 返回写入错误。
		return gerror.Wrap(err, "insert risk_decision_log failed")
	}
	// 如果决策不是允许则记录命中日志。
	if res.GetDecision() != v1.RiskDecision_RISK_DECISION_ALLOW {
		// 生成命中流水号。
		hitNo := "RHL_" + strings.ToUpper(uuid.NewString())
		// 默认规则编码。
		ruleCode := "AUTO_SCORE_RULE"
		// 若有命中规则则取首个。
		if len(res.GetMatchedRuleCodes()) > 0 {
			// 读取首个命中规则。
			ruleCode = res.GetMatchedRuleCodes()[0]
		}
		// 写入命中日志。
		_, _ = dao.RiskHitLog.Ctx(ctx).Data(do.RiskHitLog{
			HitNo:     hitNo,                   // 写入命中号。
			UserId:    userID,                  // 写入用户ID。
			RuleCode:  ruleCode,                // 写入规则编码。
			BizType:   bizType,                 // 写入业务类型。
			BizNo:     bizNo,                   // 写入业务号。
			Decision:  uint(res.GetDecision()), // 写入决策值。
			CreatedAt: gtime.Now(),             // 写入创建时间。
		}).Insert()
	}
	// 返回成功。
	return nil
}

// boolToInt 把布尔值转换为0/1。
func boolToInt(v bool) int {
	// true 转换为1。
	if v {
		// 返回1。
		return 1
	}
	// false 转换为0。
	return 0
}

// timestamppbFromGTime 把 gtime 时间转换为 protobuf 时间。
func timestamppbFromGTime(t *gtime.Time) *timestamppb.Timestamp {
	// 判空保护。
	if t == nil {
		// 空值返回 nil。
		return nil
	}
	// 返回 protobuf 时间对象。
	return timestamppb.New(t.Time)
}
