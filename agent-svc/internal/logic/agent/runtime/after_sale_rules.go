package runtime

import "strings"

const (
	sceneCodeExistingAfterSale   = "existing_after_sale_processing"
	sceneCodeRefundBeforeShip    = "refund_before_shipment"
	sceneCodeUrgeShipmentPending = "urge_shipment_pending"
	sceneCodeRefundInTransit     = "refund_in_transit_wait"
	sceneCodeLogisticsWatch      = "logistics_in_transit_watch"
	sceneCodeReturnRefund        = "return_refund_delivered"
	sceneCodeExchangeDelivered   = "exchange_delivered"
	sceneCodeGenericGuidance     = "generic_after_sale_guidance"
	sceneCodeRuleConflict        = "rule_conflict"

	problemTypeRefund       = "refund"
	problemTypeReturnRefund = "return_refund"
	problemTypeExchange     = "exchange"
	problemTypeUrgeShipment = "urge_shipment"
	problemTypeLogistics    = "logistics"
	problemTypeOrderStatus  = "order_status"
	problemTypeUnknown      = "unknown"
)

type AfterSaleRuleContext struct {
	TaskCode    string
	ProblemType string
	Snapshot    *OrderSnapshot
	Session     TaskSessionState
}

type AfterSaleRuleOutcome struct {
	DecisionCard       *AfterSaleDecisionCard
	SuggestedActions   []SuggestedAction
	ReplyText          string
	HandoffRecommended bool
	HandoffReasonCode  string
}

type AfterSaleRule struct {
	SceneCode string
	Priority  int
	Match     func(ctx AfterSaleRuleContext) bool
	Build     func(ctx AfterSaleRuleContext) AfterSaleRuleOutcome
}

type AfterSaleRuleEngine interface {
	Evaluate(ctx AfterSaleRuleContext) AfterSaleRuleOutcome
}

type defaultAfterSaleRuleEngine struct {
	rules []AfterSaleRule
}

func NewAfterSaleRuleEngine(rules []AfterSaleRule) AfterSaleRuleEngine {
	if len(rules) == 0 {
		rules = defaultAfterSaleRules()
	}
	return &defaultAfterSaleRuleEngine{rules: rules}
}

func (e *defaultAfterSaleRuleEngine) Evaluate(ctx AfterSaleRuleContext) AfterSaleRuleOutcome {
	if e == nil || len(e.rules) == 0 {
		return buildGenericAfterSaleOutcome(ctx)
	}

	var selected *AfterSaleRule
	conflicts := make([]AfterSaleRule, 0, 2)
	for _, rule := range e.rules {
		if rule.Match == nil || !rule.Match(ctx) {
			continue
		}
		if selected == nil || rule.Priority > selected.Priority {
			ruleCopy := rule
			selected = &ruleCopy
			conflicts = conflicts[:0]
			conflicts = append(conflicts, rule)
			continue
		}
		if selected != nil && rule.Priority == selected.Priority {
			conflicts = append(conflicts, rule)
		}
	}

	if selected == nil {
		return buildGenericAfterSaleOutcome(ctx)
	}
	if hasRuleConflict(conflicts) {
		return buildRuleConflictOutcome()
	}
	if selected.Build == nil {
		return buildGenericAfterSaleOutcome(ctx)
	}
	return selected.Build(ctx)
}

func defaultAfterSaleRules() []AfterSaleRule {
	return []AfterSaleRule{
		{
			SceneCode: sceneCodeExistingAfterSale,
			Priority:  100,
			Match: func(ctx AfterSaleRuleContext) bool {
				return isProcessingAfterSale(ctx.Snapshot)
			},
			Build: func(ctx AfterSaleRuleContext) AfterSaleRuleOutcome {
				return AfterSaleRuleOutcome{
					DecisionCard: &AfterSaleDecisionCard{
						SceneCode:        sceneCodeExistingAfterSale,
						DecisionPathCode: decisionPathWaitExistingAfterSale,
						ReasonText:       "当前订单已经有售后单在处理中，不建议重复发起。",
						ConstraintText:   "重复提交可能导致客服判断分散。",
						NextStepText:     "建议等待当前售后单处理进展，必要时转人工跟进。",
					},
					SuggestedActions: []SuggestedAction{
						buildAction(actionCodeWaitForUpdate, "等待进展", true, ""),
						buildAction(actionCodeEscalateToHuman, "转人工", true, ""),
					},
					ReplyText: "当前订单已经有售后处理记录，我建议先等待现有售后结果，避免重复提交。",
				}
			},
		},
		{
			SceneCode: sceneCodeRefundBeforeShip,
			Priority:  90,
			Match: func(ctx AfterSaleRuleContext) bool {
				return isUnshipped(ctx.Snapshot) && ctx.ProblemType == problemTypeRefund
			},
			Build: func(ctx AfterSaleRuleContext) AfterSaleRuleOutcome {
				return AfterSaleRuleOutcome{
					DecisionCard: &AfterSaleDecisionCard{
						SceneCode:        sceneCodeRefundBeforeShip,
						DecisionPathCode: decisionPathRefundOnly,
						ReasonText:       "订单尚未发货，当前更适合直接申请退款。",
						ConstraintText:   "未发货阶段通常不需要先走退货流程。",
						NextStepText:     "优先发起退款；如果商家长时间未处理，可再考虑转人工。",
					},
					SuggestedActions: []SuggestedAction{
						buildAction(actionCodeRequestRefund, "申请退款", true, ""),
						buildAction(actionCodeEscalateToHuman, "转人工", true, ""),
					},
					ReplyText: "根据当前订单事实，这个订单还没有发货，更适合直接申请退款。",
				}
			},
		},
		{
			SceneCode: sceneCodeUrgeShipmentPending,
			Priority:  85,
			Match: func(ctx AfterSaleRuleContext) bool {
				return isUnshipped(ctx.Snapshot) && ctx.ProblemType == problemTypeUrgeShipment
			},
			Build: func(ctx AfterSaleRuleContext) AfterSaleRuleOutcome {
				return AfterSaleRuleOutcome{
					DecisionCard: &AfterSaleDecisionCard{
						SceneCode:        sceneCodeUrgeShipmentPending,
						DecisionPathCode: decisionPathWaitShipment,
						ReasonText:       "订单还处于待发货阶段，当前更适合先催发货或继续等待商家处理。",
						ConstraintText:   "在未发货阶段，不建议把普通履约等待误判为物流异常。",
						NextStepText:     "可以先催发货；如果超过约定时效仍无进展，再考虑退款或转人工。",
					},
					SuggestedActions: []SuggestedAction{
						buildAction(actionCodeUrgeShipment, "催发货", true, ""),
						buildAction(actionCodeWaitForUpdate, "继续等待", true, ""),
						buildAction(actionCodeRequestRefund, "申请退款", true, ""),
					},
					ReplyText: "这个订单目前还在待发货阶段，我更建议你先催发货或继续观察发货时效。",
				}
			},
		},
		{
			SceneCode: sceneCodeRefundInTransit,
			Priority:  80,
			Match: func(ctx AfterSaleRuleContext) bool {
				return isInTransit(ctx.Snapshot) && ctx.ProblemType == problemTypeRefund
			},
			Build: func(ctx AfterSaleRuleContext) AfterSaleRuleOutcome {
				return AfterSaleRuleOutcome{
					DecisionCard: &AfterSaleDecisionCard{
						SceneCode:        sceneCodeRefundInTransit,
						DecisionPathCode: decisionPathWaitShipment,
						ReasonText:       "订单已经发货且仍在运输途中，当前更适合先观察物流或在签收节点后再走退货退款。",
						ConstraintText:   "物流在途中时，直接退款通常会受履约节点限制。",
						NextStepText:     "建议先继续关注物流；若后续签收或拒收，再根据状态选择退货退款。",
					},
					SuggestedActions: []SuggestedAction{
						buildAction(actionCodeViewLogistics, "查看物流", true, ""),
						buildAction(actionCodeWaitForUpdate, "继续等待", true, ""),
						buildAction(actionCodeEscalateToHuman, "转人工", true, ""),
					},
					ReplyText: "订单已经发货且仍在运输途中，我不建议现在直接给出可退款结论，先观察物流节点会更稳妥。",
				}
			},
		},
		{
			SceneCode: sceneCodeLogisticsWatch,
			Priority:  75,
			Match: func(ctx AfterSaleRuleContext) bool {
				return isInTransit(ctx.Snapshot) && ctx.ProblemType == problemTypeLogistics
			},
			Build: func(ctx AfterSaleRuleContext) AfterSaleRuleOutcome {
				return AfterSaleRuleOutcome{
					DecisionCard: &AfterSaleDecisionCard{
						SceneCode:        sceneCodeLogisticsWatch,
						DecisionPathCode: decisionPathWaitShipment,
						ReasonText:       "订单已发货且物流仍在正常运输区间，当前更适合先陈述事实并继续观察。",
						ConstraintText:   "没有达到异常阈值前，不应误判为丢件或物流异常。",
						NextStepText:     "建议先查看最新物流节点；如果长时间无更新，再转人工进一步核实。",
					},
					SuggestedActions: []SuggestedAction{
						buildAction(actionCodeViewLogistics, "查看物流", true, ""),
						buildAction(actionCodeWaitForUpdate, "继续等待", true, ""),
						buildAction(actionCodeEscalateToHuman, "转人工", true, ""),
					},
					ReplyText: "目前订单物流仍处于运输途中，我先按已知物流事实给你说明，不会把它直接判断成异常件。",
				}
			},
		},
		{
			SceneCode: sceneCodeExchangeDelivered,
			Priority:  70,
			Match: func(ctx AfterSaleRuleContext) bool {
				return isDelivered(ctx.Snapshot) && ctx.ProblemType == problemTypeExchange
			},
			Build: func(ctx AfterSaleRuleContext) AfterSaleRuleOutcome {
				return AfterSaleRuleOutcome{
					DecisionCard: &AfterSaleDecisionCard{
						SceneCode:        sceneCodeExchangeDelivered,
						DecisionPathCode: decisionPathReturnRefundOrExchange,
						ReasonText:       "订单已经履约完成，如果是商品问题或规格不符，更适合走换货路径。",
						ConstraintText:   "换货是否支持仍取决于平台和店铺规则。",
						NextStepText:     "建议优先申请换货；若店铺不支持，再考虑退货退款或转人工。",
					},
					SuggestedActions: []SuggestedAction{
						buildAction(actionCodeRequestExchange, "申请换货", true, ""),
						buildAction(actionCodeRequestReturn, "退货退款", true, ""),
					},
					ReplyText: "订单已经完成履约，如果你更关心换货，这一单现在更适合走换货处理。",
				}
			},
		},
		{
			SceneCode: sceneCodeReturnRefund,
			Priority:  70,
			Match: func(ctx AfterSaleRuleContext) bool {
				return isDelivered(ctx.Snapshot) && (ctx.ProblemType == problemTypeReturnRefund || ctx.ProblemType == problemTypeRefund)
			},
			Build: func(ctx AfterSaleRuleContext) AfterSaleRuleOutcome {
				return AfterSaleRuleOutcome{
					DecisionCard: &AfterSaleDecisionCard{
						SceneCode:        sceneCodeReturnRefund,
						DecisionPathCode: decisionPathReturnRefundOrExchange,
						ReasonText:       "订单已发货并完成履约，当前更适合走退货退款流程。",
						ConstraintText:   "是否支持上门取件、补差价或免运费，仍以平台和店铺规则为准。",
						NextStepText:     "建议先发起退货退款；如果规则解释不清，再转人工核实。",
					},
					SuggestedActions: []SuggestedAction{
						buildAction(actionCodeRequestReturn, "退货退款", true, ""),
						buildAction(actionCodeRequestExchange, "申请换货", true, ""),
						buildAction(actionCodeEscalateToHuman, "转人工", true, ""),
					},
					ReplyText: "这个订单已经进入履约完成阶段，如果你想退款，当前更适合走退货退款。",
				}
			},
		},
	}
}

func buildGenericAfterSaleOutcome(ctx AfterSaleRuleContext) AfterSaleRuleOutcome {
	return AfterSaleRuleOutcome{
		DecisionCard: &AfterSaleDecisionCard{
			SceneCode:        sceneCodeGenericGuidance,
			DecisionPathCode: decisionPathReturnRefundOrExchange,
			ReasonText:       "当前订单已进入售后判断阶段，但缺少更细的规则命中线索。",
			ConstraintText:   "我不会在缺少事实和规则依据时直接猜结论。",
			NextStepText:     "可以继续补充问题类型，我会基于订单事实给出更明确建议。",
		},
		SuggestedActions: []SuggestedAction{
			buildAction(actionCodeViewLogistics, "查看物流", true, ""),
			buildAction(actionCodeEscalateToHuman, "转人工", true, ""),
		},
		ReplyText: "我已经定位到订单，但当前更适合先结合订单状态继续判断退款、退货退款还是换货。",
	}
}

func buildRuleConflictOutcome() AfterSaleRuleOutcome {
	return AfterSaleRuleOutcome{
		DecisionCard: &AfterSaleDecisionCard{
			SceneCode:        sceneCodeRuleConflict,
			DecisionPathCode: decisionPathReturnRefundOrExchange,
			ReasonText:       "当前订单同时命中了多个同优先级规则，系统无法稳定给出唯一售后建议。",
			ConstraintText:   "为了避免错误引导，我不会在规则冲突时强行下结论。",
			NextStepText:     "建议转人工，由客服结合更完整的上下文继续处理。",
		},
		SuggestedActions: []SuggestedAction{
			buildAction(actionCodeEscalateToHuman, "转人工", true, ""),
		},
		ReplyText:          "这个问题当前命中了冲突规则，为了避免误导，我建议直接转人工继续核实。",
		HandoffRecommended: true,
		HandoffReasonCode:  escalationReasonRuleConflict,
	}
}

func hasRuleConflict(rules []AfterSaleRule) bool {
	if len(rules) <= 1 {
		return false
	}
	firstSceneCode := strings.TrimSpace(rules[0].SceneCode)
	for _, rule := range rules[1:] {
		if strings.TrimSpace(rule.SceneCode) != firstSceneCode {
			return true
		}
	}
	return false
}

func buildAction(actionCode, label string, enabled bool, reasonIfDisabled string) SuggestedAction {
	return SuggestedAction{
		ActionCode:       actionCode,
		Label:            label,
		Enabled:          enabled,
		ReasonIfDisabled: reasonIfDisabled,
	}
}

func inferProblemType(query, activeTaskCode string) string {
	lowerText := strings.ToLower(strings.TrimSpace(query))
	switch {
	case containsAny(lowerText, "换货"):
		return problemTypeExchange
	case containsAny(lowerText, "退货退款", "退货"):
		return problemTypeReturnRefund
	case containsAny(lowerText, "退款"):
		return problemTypeRefund
	case containsAny(lowerText, "催发货", "怎么还没发货", "什么时候发货"):
		return problemTypeUrgeShipment
	case containsAny(lowerText, "物流", "快递", "包裹"):
		return problemTypeLogistics
	case containsAny(lowerText, "订单状态", "订单进度"):
		return problemTypeOrderStatus
	}

	switch strings.TrimSpace(activeTaskCode) {
	case taskCodeRefundDecision:
		return problemTypeRefund
	case taskCodeLogisticsQuery:
		return problemTypeLogistics
	case taskCodeOrderStatusQuery:
		return problemTypeOrderStatus
	default:
		return problemTypeUnknown
	}
}

func isProcessingAfterSale(snapshot *OrderSnapshot) bool {
	return snapshot != nil && strings.EqualFold(strings.TrimSpace(snapshot.AfterSaleStatus), "PROCESSING")
}

func isUnshipped(snapshot *OrderSnapshot) bool {
	if snapshot == nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(snapshot.FulfillmentStatus), "UNSHIPPED") ||
		strings.EqualFold(strings.TrimSpace(snapshot.LogisticsStatus), "NOT_SHIPPED")
}

func isInTransit(snapshot *OrderSnapshot) bool {
	if snapshot == nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(snapshot.FulfillmentStatus), "SHIPPED") ||
		strings.EqualFold(strings.TrimSpace(snapshot.LogisticsStatus), "IN_TRANSIT")
}

func isDelivered(snapshot *OrderSnapshot) bool {
	if snapshot == nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(snapshot.FulfillmentStatus), "DELIVERED") ||
		strings.EqualFold(strings.TrimSpace(snapshot.LogisticsStatus), "DELIVERED") ||
		strings.EqualFold(strings.TrimSpace(snapshot.LogisticsStatus), "SIGNED")
}
