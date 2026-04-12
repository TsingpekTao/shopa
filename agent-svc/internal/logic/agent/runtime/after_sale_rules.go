package runtime

import "strings"

const (
	sceneCodeExistingAfterSale    = "existing_after_sale_processing"
	sceneCodeRefundBeforeShip     = "refund_before_shipment"
	sceneCodeUrgeShipmentPending  = "urge_shipment_pending"
	sceneCodeLogisticsBeforeShip  = "logistics_before_shipment"
	sceneCodeRefundInTransit      = "refund_in_transit_wait"
	sceneCodeLogisticsWatch       = "logistics_in_transit_watch"
	sceneCodeOrderStatusUnshipped = "order_status_unshipped"
	sceneCodeOrderStatusInTransit = "order_status_in_transit"
	sceneCodeOrderStatusDelivered = "order_status_delivered"
	sceneCodeReturnRefund         = "return_refund_delivered"
	sceneCodeExchangeDelivered    = "exchange_delivered"
	sceneCodeGenericGuidance      = "generic_after_sale_guidance"
	sceneCodeRuleConflict         = "rule_conflict"

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

	var (
		selected  *AfterSaleRule
		conflicts = make([]AfterSaleRule, 0, 2)
	)

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
						NextStepText:     "建议先等待当前售后单处理进展，必要时再转人工继续跟进。",
					},
					SuggestedActions: []SuggestedAction{
						buildAction(actionCodeWaitForUpdate, "等待进展", true, ""),
						buildAction(actionCodeEscalateToHuman, "转人工", true, ""),
					},
					ReplyText: "当前订单已经存在售后处理记录，我建议先看这一笔售后的处理进展，避免重复提交。",
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
						NextStepText:     "你可以优先发起退款；如果商家长时间未处理，再考虑转人工继续核实。",
					},
					SuggestedActions: []SuggestedAction{
						buildAction(actionCodeRequestRefund, "申请退款", true, ""),
						buildAction(actionCodeEscalateToHuman, "转人工", true, ""),
					},
					ReplyText: "根据当前订单事实，这笔订单还没有发货，现在更适合直接申请退款。",
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
						ReasonText:       "订单目前处于待发货阶段，当前更适合先催发货或继续等待商家处理。",
						ConstraintText:   "在未发货阶段，不应把正常等待直接判断成物流异常。",
						NextStepText:     "如果你着急收货，可以先催发货；如果不想继续等，也可以直接退款。",
					},
					SuggestedActions: []SuggestedAction{
						buildAction(actionCodeUrgeShipment, "催发货", true, ""),
						buildAction(actionCodeWaitForUpdate, "继续等待", true, ""),
						buildAction(actionCodeRequestRefund, "申请退款", true, ""),
					},
					ReplyText: "这个订单目前还是待发货状态，我更建议你先催发货，或者根据自己的诉求直接申请退款。",
				}
			},
		},
		{
			SceneCode: sceneCodeLogisticsBeforeShip,
			Priority:  82,
			Match: func(ctx AfterSaleRuleContext) bool {
				return isUnshipped(ctx.Snapshot) && ctx.ProblemType == problemTypeLogistics
			},
			Build: func(ctx AfterSaleRuleContext) AfterSaleRuleOutcome {
				return AfterSaleRuleOutcome{
					DecisionCard: &AfterSaleDecisionCard{
						SceneCode:        sceneCodeLogisticsBeforeShip,
						DecisionPathCode: decisionPathWaitShipment,
						ReasonText:       "订单当前还没发货，所以暂时不会有在途物流轨迹。",
						ConstraintText:   "未发货前查不到运输节点，不代表包裹异常或丢件。",
						NextStepText:     "如果你想尽快收到货，可以先催发货；如果不想继续等，也可以直接申请退款。",
					},
					SuggestedActions: []SuggestedAction{
						buildAction(actionCodeUrgeShipment, "催发货", true, ""),
						buildAction(actionCodeRequestRefund, "申请退款", true, ""),
						buildAction(actionCodeWaitForUpdate, "继续等待", true, ""),
					},
					ReplyText: "我已经帮你定位到这个订单了。它目前还是待发货状态，所以暂时看不到物流轨迹；接下来更适合先催发货，或者直接申请退款。",
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
						NextStepText:     "建议先继续关注物流；如果后续签收或拒收，再根据状态选择退货退款。",
					},
					SuggestedActions: []SuggestedAction{
						buildAction(actionCodeViewLogistics, "查看物流", true, ""),
						buildAction(actionCodeWaitForUpdate, "继续等待", true, ""),
						buildAction(actionCodeEscalateToHuman, "转人工", true, ""),
					},
					ReplyText: "订单已经发货并且还在运输途中，我不建议现在直接给出可退款结论，先观察物流节点会更稳妥。",
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
						ReasonText:       "订单已发货，物流仍在正常运输区间，当前更适合先同步事实并继续观察。",
						ConstraintText:   "没有达到异常阈值前，不应该直接判成丢件或物流异常。",
						NextStepText:     "建议先看最新物流节点；如果后续长时间不更新，再转人工核实。",
					},
					SuggestedActions: []SuggestedAction{
						buildAction(actionCodeViewLogistics, "查看物流", true, ""),
						buildAction(actionCodeWaitForUpdate, "继续等待", true, ""),
						buildAction(actionCodeEscalateToHuman, "转人工", true, ""),
					},
					ReplyText: "目前订单物流还在正常运输途中，我先按已知物流事实给你同步，不会直接把它判断成异常件。",
				}
			},
		},
		{
			SceneCode: sceneCodeOrderStatusUnshipped,
			Priority:  74,
			Match: func(ctx AfterSaleRuleContext) bool {
				return isUnshipped(ctx.Snapshot) && ctx.ProblemType == problemTypeOrderStatus
			},
			Build: func(ctx AfterSaleRuleContext) AfterSaleRuleOutcome {
				return AfterSaleRuleOutcome{
					DecisionCard: &AfterSaleDecisionCard{
						SceneCode:        sceneCodeOrderStatusUnshipped,
						DecisionPathCode: decisionPathWaitShipment,
						ReasonText:       "这笔订单当前是待发货状态。",
						ConstraintText:   "商家未发货前，系统不会出现新的物流轨迹。",
						NextStepText:     "如果你着急，可以先催发货；如果不想继续等待，也可以申请退款。",
					},
					SuggestedActions: []SuggestedAction{
						buildAction(actionCodeUrgeShipment, "催发货", true, ""),
						buildAction(actionCodeRequestRefund, "申请退款", true, ""),
					},
					ReplyText: "我已经同步到订单状态了：这笔订单目前还没发货，接下来更适合催发货或者直接退款。",
				}
			},
		},
		{
			SceneCode: sceneCodeOrderStatusInTransit,
			Priority:  74,
			Match: func(ctx AfterSaleRuleContext) bool {
				return isInTransit(ctx.Snapshot) && ctx.ProblemType == problemTypeOrderStatus
			},
			Build: func(ctx AfterSaleRuleContext) AfterSaleRuleOutcome {
				return AfterSaleRuleOutcome{
					DecisionCard: &AfterSaleDecisionCard{
						SceneCode:        sceneCodeOrderStatusInTransit,
						DecisionPathCode: decisionPathWaitShipment,
						ReasonText:       "这笔订单当前已经发货，正在运输途中。",
						ConstraintText:   "只要物流节点还在正常更新，就不建议直接判成异常件。",
						NextStepText:     "你可以继续查看物流进度；如果后续长时间不更新，再考虑转人工核实。",
					},
					SuggestedActions: []SuggestedAction{
						buildAction(actionCodeViewLogistics, "查看物流", true, ""),
						buildAction(actionCodeWaitForUpdate, "继续等待", true, ""),
					},
					ReplyText: "我已经同步到订单状态了：它现在正在运输途中，你可以继续看物流节点更新。",
				}
			},
		},
		{
			SceneCode: sceneCodeOrderStatusDelivered,
			Priority:  74,
			Match: func(ctx AfterSaleRuleContext) bool {
				return isDelivered(ctx.Snapshot) && ctx.ProblemType == problemTypeOrderStatus
			},
			Build: func(ctx AfterSaleRuleContext) AfterSaleRuleOutcome {
				return AfterSaleRuleOutcome{
					DecisionCard: &AfterSaleDecisionCard{
						SceneCode:        sceneCodeOrderStatusDelivered,
						DecisionPathCode: decisionPathReturnRefundOrExchange,
						ReasonText:       "这笔订单当前已经完成签收。",
						ConstraintText:   "如果接下来要处理商品问题，通常会进入退货退款或换货流程。",
						NextStepText:     "你可以继续告诉我诉求，我会按订单事实帮你判断该走退货退款还是换货。",
					},
					SuggestedActions: []SuggestedAction{
						buildAction(actionCodeRequestReturn, "退货退款", true, ""),
						buildAction(actionCodeRequestExchange, "申请换货", true, ""),
					},
					ReplyText: "我已经同步到订单状态了：这笔订单目前已经签收。如果你要继续处理售后，我可以接着帮你判断该走退货退款还是换货。",
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
						ConstraintText:   "是否支持换货，仍取决于平台和店铺规则。",
						NextStepText:     "建议优先申请换货；如果店铺不支持，再考虑退货退款或转人工。",
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
					ReplyText: "这笔订单已经进入履约完成阶段，如果你想退款，当前更适合走退货退款。",
				}
			},
		},
	}
}

func buildGenericAfterSaleOutcome(ctx AfterSaleRuleContext) AfterSaleRuleOutcome {
	replyText := "我已经定位到订单了，但还缺少更明确的场景线索。你可以继续告诉我是想查物流、催发货，还是判断退款、退货退款或换货。"
	reasonText := "当前订单已经进入售后判断阶段，但还没有命中更细的规则场景。"
	nextStepText := "你可以继续补充更具体的诉求，我会根据订单事实给你更明确的建议。"

	switch ctx.TaskCode {
	case taskCodeLogisticsQuery:
		replyText = "我已经定位到订单了，但当前还需要结合物流和履约状态继续判断。你也可以继续问我物流进度、是否发货，或者接下来该怎么处理。"
		reasonText = "当前订单已经定位成功，但物流场景还需要更多上下文来给出更具体判断。"
		nextStepText = "你可以继续问物流进度、是否发货，或者告诉我你是想催发货还是转人工。"
	case taskCodeOrderStatusQuery:
		replyText = "我已经定位到订单了，先把当前订单事实同步给你；如果你还想继续判断退款、退货退款或换货，也可以直接告诉我。"
		reasonText = "当前订单事实已经可用，但还没有命中更具体的订单状态场景。"
		nextStepText = "你可以继续问我是否能退款、是否适合换货，或者继续查看物流进度。"
	}

	return AfterSaleRuleOutcome{
		DecisionCard: &AfterSaleDecisionCard{
			SceneCode:        sceneCodeGenericGuidance,
			DecisionPathCode: decisionPathReturnRefundOrExchange,
			ReasonText:       reasonText,
			ConstraintText:   "我不会在缺少事实和规则依据时直接猜结论。",
			NextStepText:     nextStepText,
		},
		SuggestedActions: []SuggestedAction{
			buildAction(actionCodeViewLogistics, "查看物流", true, ""),
			buildAction(actionCodeEscalateToHuman, "转人工", true, ""),
		},
		ReplyText: replyText,
	}
}

func buildRuleConflictOutcome() AfterSaleRuleOutcome {
	return AfterSaleRuleOutcome{
		DecisionCard: &AfterSaleDecisionCard{
			SceneCode:        sceneCodeRuleConflict,
			DecisionPathCode: decisionPathReturnRefundOrExchange,
			ReasonText:       "当前订单同时命中了多个同优先级规则，系统无法稳定给出唯一售后建议。",
			ConstraintText:   "为了避免错误引导，我不会在规则冲突时强行给结论。",
			NextStepText:     "建议转人工，由客服结合更完整的上下文继续处理。",
		},
		SuggestedActions: []SuggestedAction{
			buildAction(actionCodeEscalateToHuman, "转人工", true, ""),
		},
		ReplyText:          "这个问题当前命中了冲突规则。为了避免误导，我建议直接转人工继续核实。",
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
	case hasShipmentStatusIntent(lowerText):
		return problemTypeOrderStatus
	case hasExchangeIntent(lowerText):
		return problemTypeExchange
	case hasReturnRefundIntent(lowerText):
		return problemTypeReturnRefund
	case hasRefundIntent(lowerText):
		return problemTypeRefund
	case hasUrgeShipmentIntent(lowerText):
		return problemTypeUrgeShipment
	case hasLogisticsIntent(lowerText):
		return problemTypeLogistics
	case hasOrderStatusIntent(lowerText):
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
