package runtime

import (
	"context"
	"fmt"
	"strings"
)

const (
	officialIntentGreeting      = "GREETING"
	officialIntentAttachment    = "ATTACHMENT_ONLY"
	officialIntentOrderQuery    = "ORDER_QUERY"
	officialIntentAfterSale     = "AFTER_SALE"
	officialIntentLogistics     = "LOGISTICS"
	officialIntentPolicyQA      = "POLICY_QA"
	officialIntentCatalogGuide  = "CATALOG_GUIDE"
	officialIntentHandoff       = "HANDOFF_REQUEST"
	officialIntentOutOfScope    = "OUT_OF_SCOPE"
	officialActionToolCall      = "TOOL_CALL"
	officialActionClarification = "ASK_CLARIFICATION"
	officialActionDirectReply   = "DIRECT_RESPOND"
)

type ModelCompletionRequest struct {
	SystemPrompt string
	UserPrompt   string
}

func hasGenericOrderQueryIntent(text string) bool {
	if hasLogisticsIntent(text) || hasUrgeShipmentIntent(text) || hasRefundIntent(text) || hasReturnRefundIntent(text) || hasExchangeIntent(text) || hasPolicyIntent(text) {
		return false
	}
	return containsAny(
		text,
		"帮我查一下这个订单",
		"帮我看看这个订单",
		"查一下订单",
		"订单详情",
		"订单信息",
		"这个订单",
		"check this order",
		"order details",
		"order info",
	)
}

type ModelClient interface {
	Complete(ctx context.Context, req ModelCompletionRequest) (string, error)
}

type PlannerDecision struct {
	Thought            string            `json:"thought,omitempty"`
	Intent             string            `json:"intent,omitempty"`
	UserGoal           string            `json:"user_goal,omitempty"`
	Confidence         float64           `json:"confidence,omitempty"`
	SlotValues         map[string]string `json:"slot_values,omitempty"`
	MissingSlots       []string          `json:"missing_slots,omitempty"`
	NextAction         PlannerNextAction `json:"next_action,omitempty"`
	HandoffRecommended bool              `json:"handoff_recommended,omitempty"`
	HandoffReasonCode  string            `json:"handoff_reason_code,omitempty"`
}

type PlannerNextAction struct {
	Type     string            `json:"type,omitempty"`
	ToolName string            `json:"tool_name,omitempty"`
	ToolArgs map[string]string `json:"tool_args,omitempty"`
	Reason   string            `json:"reason,omitempty"`
}

func (r *einoRunner) generateOfficialReply(ctx context.Context, state *runState) (*ReplyPayload, TaskSessionState, error) {
	if reply, session, ok := maybeFuseSlotFilling(state); ok {
		return reply, session, nil
	}

	decision := r.planOfficialDecision(ctx, state)
	session := copyTaskSession(state.TaskSession)
	session.ActiveIntent = decision.Intent
	session.UserGoal = strings.TrimSpace(decision.UserGoal)

	if strings.TrimSpace(decision.NextAction.Type) == "" {
		decision.NextAction.Type = officialActionDirectReply
	}

	switch decision.NextAction.Type {
	case officialActionToolCall:
		return r.executeOfficialToolDecision(ctx, state, session, decision)
	case officialActionClarification:
		reply := ReplyPayload{
			ReplyText:       coalesce(decision.NextAction.Reason, "请告诉我您想查哪一笔订单，或者直接补充订单号。"),
			IntentCode:      decision.Intent,
			GuardResultCode: state.GuardResultCode,
			Confidence:      maxFloat(decision.Confidence, 0.82),
		}
		return &reply, session, nil
	default:
		reply := buildDirectOfficialReply(state, session, decision)
		return &reply, session, nil
	}
}

func (r *einoRunner) planOfficialDecision(ctx context.Context, state *runState) PlannerDecision {
	if r.modelClient != nil {
		if decision, ok := r.tryModelPlannerDecision(ctx, state); ok {
			return decision
		}
	}
	return heuristicPlannerDecision(state)
}

func (r *einoRunner) tryModelPlannerDecision(ctx context.Context, state *runState) (PlannerDecision, bool) {
	var (
		systemPrompt = buildPlannerSystemPrompt(state)
		userPrompt   = strings.TrimSpace(state.Message.ContentText)
	)

	for attempt := 0; attempt <= maxPlannerRepairAttempts; attempt++ {
		raw, err := r.modelClient.Complete(ctx, ModelCompletionRequest{
			SystemPrompt: systemPrompt,
			UserPrompt:   userPrompt,
		})
		if err != nil {
			return PlannerDecision{}, false
		}

		decision, validationErr := decodePlannerDecision(raw)
		if validationErr == nil {
			return decision, true
		}
		if !validationErr.Retryable || attempt >= maxPlannerRepairAttempts {
			return PlannerDecision{}, false
		}
		userPrompt = buildPlannerRepairUserPrompt(strings.TrimSpace(state.Message.ContentText), raw, validationErr)
	}

	return PlannerDecision{}, false
}

func heuristicPlannerDecision(state *runState) PlannerDecision {
	var (
		query           = strings.TrimSpace(state.Message.ContentText)
		selectedOrderNo = resolveSelectedOrderNo(state.TaskSession, state.Conversation)
		intent          = inferOfficialIntent(query, state.IntentCode)
	)
	if selectedOrderNo == "" {
		selectedOrderNo = extractOrderNumber(query, "ord")
	}

	if strings.TrimSpace(selectedOrderNo) == "" && hasExplicitOrderSelectionRequest(query) {
		intent = officialIntentLogistics
	}

	decision := PlannerDecision{
		Intent:     intent,
		UserGoal:   query,
		Confidence: 0.88,
		NextAction: PlannerNextAction{
			Type: officialActionDirectReply,
		},
	}

	switch intent {
	case officialIntentGreeting:
		decision.NextAction.Reason = "您好，请问我能帮您什么？您可以直接问订单、物流、退款退货、商品推荐或转人工。"
	case officialIntentAttachment:
		decision.NextAction.Reason = "我可以帮您处理订单和售后问题，请补充订单号、截图说明，或直接描述您的问题。"
	case officialIntentHandoff:
		decision.HandoffRecommended = true
		decision.HandoffReasonCode = escalationReasonUserRequested
		decision.NextAction.Reason = "我可以先帮您整理当前问题，也可以直接为您转接人工客服。"
	case officialIntentPolicyQA:
		decision.NextAction.Type = officialActionToolCall
		decision.NextAction.ToolName = "search_knowledge_chunks"
		decision.NextAction.ToolArgs = map[string]string{"query": query}
		decision.NextAction.Reason = "先查规则知识，再组织成用户容易理解的话术。"
	case officialIntentCatalogGuide:
		decision.NextAction.Type = officialActionToolCall
		decision.NextAction.ToolName = "search_products"
		decision.NextAction.ToolArgs = map[string]string{"query": query}
		decision.NextAction.Reason = "先搜索商品，再给出推荐结果。"
	case officialIntentOrderQuery:
		decision.NextAction.Type = officialActionToolCall
		decision.NextAction.ToolName = "get_order_snapshot"
		decision.NextAction.Reason = "需要先确认订单事实，再同步当前订单状态。"
		if strings.TrimSpace(selectedOrderNo) == "" {
			decision.MissingSlots = []string{slotCodeOrderNo}
			decision.NextAction.ToolName = "list_recent_orders"
		}
	case officialIntentLogistics:
		decision.NextAction.Type = officialActionToolCall
		decision.NextAction.ToolName = "query_logistics"
		decision.NextAction.Reason = "需要先查询订单物流和发货进度。"
		if strings.TrimSpace(selectedOrderNo) == "" {
			decision.MissingSlots = []string{slotCodeOrderNo}
			decision.NextAction.ToolName = "list_recent_orders"
		}
	case officialIntentAfterSale:
		decision.NextAction.Type = officialActionToolCall
		decision.NextAction.ToolName = "get_order_snapshot"
		decision.NextAction.Reason = "需要先看订单事实，再判断退款、退货退款或换货建议。"
		if strings.TrimSpace(selectedOrderNo) == "" {
			decision.MissingSlots = []string{slotCodeOrderNo}
			decision.NextAction.ToolName = "list_recent_orders"
		}
	default:
		decision.NextAction.Reason = "抱歉，我目前主要负责订单、物流、售后、规则说明和商品推荐。"
	}

	return normalizePlannerDecision(decision)
}

func normalizePlannerDecision(decision PlannerDecision) PlannerDecision {
	decision.Intent = strings.ToUpper(strings.TrimSpace(decision.Intent))
	decision.UserGoal = strings.TrimSpace(decision.UserGoal)
	decision.HandoffReasonCode = strings.TrimSpace(decision.HandoffReasonCode)
	decision.NextAction.Type = strings.ToUpper(strings.TrimSpace(decision.NextAction.Type))
	decision.NextAction.ToolName = strings.TrimSpace(decision.NextAction.ToolName)
	decision.NextAction.Reason = strings.TrimSpace(decision.NextAction.Reason)
	if decision.SlotValues == nil {
		decision.SlotValues = make(map[string]string)
	}
	if decision.NextAction.ToolArgs == nil {
		decision.NextAction.ToolArgs = make(map[string]string)
	}
	return decision
}

func isValidPlannerDecision(decision PlannerDecision) bool {
	decision = normalizePlannerDecision(decision)
	if decision.Intent == "" {
		return false
	}
	switch decision.NextAction.Type {
	case officialActionToolCall:
		return isAllowedPlannerTool(decision.NextAction.ToolName)
	case officialActionClarification, officialActionDirectReply:
		return true
	default:
		return false
	}
}

func isAllowedPlannerTool(toolName string) bool {
	switch strings.TrimSpace(toolName) {
	case "list_recent_orders", "get_order_snapshot", "query_logistics", "search_knowledge_chunks", "search_products", "batch_get_spu_cards", "build_handoff_summary":
		return true
	default:
		return false
	}
}

func buildPlannerSystemPrompt(state *runState) string {
	return strings.Join([]string{
		"你是商城官方客服的 Planner，只能输出 JSON。",
		"如果当前已选订单存在，优先复用，不要重复追问。",
		"缺少 order_no 时，可以调用 list_recent_orders 让用户选订单。",
		"可用工具：list_recent_orders,get_order_snapshot,query_logistics,search_knowledge_chunks,search_products,batch_get_spu_cards,build_handoff_summary。",
		"当前 selected_order_no=" + resolveSelectedOrderNo(state.TaskSession, state.Conversation),
		"当前 selected_sub_order_no=" + resolveSelectedSubOrderNo(state.TaskSession, state.Conversation),
		"最近事实摘要=" + strings.TrimSpace(state.TaskSession.LatestFactsSummary),
	}, "\n")
}

func inferOfficialIntent(query string, guardIntent string) string {
	lowerQuery := strings.ToLower(strings.TrimSpace(query))
	switch {
	case strings.EqualFold(strings.TrimSpace(guardIntent), intentCodeGreeting), hasGreetingIntent(lowerQuery):
		return officialIntentGreeting
	case strings.EqualFold(strings.TrimSpace(guardIntent), intentCodeAttachmentOnly):
		return officialIntentAttachment
	case containsAny(lowerQuery, "人工", "转人工", "人工客服"):
		return officialIntentHandoff
	case hasPolicyIntent(lowerQuery):
		return officialIntentPolicyQA
	case hasCatalogGuideIntent(lowerQuery):
		return officialIntentCatalogGuide
	case hasGenericOrderQueryIntent(lowerQuery):
		return officialIntentOrderQuery
	case hasLogisticsIntent(lowerQuery) || hasUrgeShipmentIntent(lowerQuery):
		return officialIntentLogistics
	case hasOrderStatusIntent(lowerQuery):
		return officialIntentOrderQuery
	case hasRefundIntent(lowerQuery) || hasReturnRefundIntent(lowerQuery) || hasExchangeIntent(lowerQuery):
		return officialIntentAfterSale
	case strings.EqualFold(strings.TrimSpace(guardIntent), intentCodeAfterSaleTask):
		return officialIntentOrderQuery
	case strings.EqualFold(strings.TrimSpace(guardIntent), intentCodeOutOfScope):
		return officialIntentOutOfScope
	default:
		return officialIntentOutOfScope
	}
}

func hasExplicitOrderSelectionRequest(query string) bool {
	normalized := normalizePromptText(query)
	return matchesAnyNormalizedPrompt(normalized,
		"帮我查一下这个订单现在到哪了",
		"帮我看下这单到哪了",
		"查询这笔订单的物流",
		"这笔订单到哪了",
		"check where this order is now",
	)
}

func resolveSelectedOrderNo(session TaskSessionState, anchors ConversationAnchors) string {
	for _, value := range []string{
		session.SlotValues[slotCodeOrderNo],
		session.SelectedOrderNo,
		session.AnchoredOrderNo,
		anchors.OrderNo,
	} {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func maybeFuseSlotFilling(state *runState) (*ReplyPayload, TaskSessionState, bool) {
	session := copyTaskSession(state.TaskSession)
	if strings.TrimSpace(resolveSelectedOrderNo(session, state.Conversation)) != "" || extractOrderNumber(state.Message.ContentText, "ord") != "" {
		return nil, session, false
	}
	if session.LastSlotPromptCode != slotPromptCodeOrderNo {
		return nil, session, false
	}

	intent := inferOfficialIntent(strings.TrimSpace(state.Message.ContentText), state.IntentCode)
	if intent != officialIntentAfterSale && intent != officialIntentLogistics {
		return nil, session, false
	}

	session.SlotRetryCount++
	if session.SlotRetryCount < defaultSlotRetryThreshold {
		return nil, session, false
	}

	session.SlotFillingFailed = true
	session.HandoffRecommended = true
	session.EscalationReasonCode = escalationReasonSlotFillingFailed
	session.MissingSlots = []MissingSlot{{
		SlotCode:   slotCodeOrderNo,
		PromptText: "请选择具体订单，或者直接补充订单号。",
		Required:   true,
	}}

	reply := &ReplyPayload{
		ReplyText:          "我暂时还没有拿到足够的订单信息，继续反复追问会影响体验，建议我先帮您转人工核实。",
		IntentCode:         intent,
		GuardResultCode:    state.GuardResultCode,
		Confidence:         0.88,
		MissingSlots:       session.MissingSlots,
		SlotRetryCount:     session.SlotRetryCount,
		HandoffRecommended: true,
		HandoffReasonCode:  escalationReasonSlotFillingFailed,
		SuggestedActions: []SuggestedAction{{
			ActionCode: actionCodeEscalateToHuman,
			Label:      "转人工",
			Enabled:    true,
		}},
	}
	return reply, session, true
}

func resolveSelectedSubOrderNo(session TaskSessionState, anchors ConversationAnchors) string {
	for _, value := range []string{
		session.SlotValues[slotCodeSubOrderNo],
		session.SelectedSubOrderNo,
		session.AnchoredSubOrderNo,
		anchors.SubOrderNo,
	} {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func (r *einoRunner) executeOfficialToolDecision(
	ctx context.Context,
	state *runState,
	session TaskSessionState,
	decision PlannerDecision,
) (*ReplyPayload, TaskSessionState, error) {
	state.SelectedToolName = strings.TrimSpace(decision.NextAction.ToolName)
	state.SelectedAdapter = "official-tool-executor"
	state.SelectedScopeCode = "OFFICIAL_SUPPORT"

	switch decision.NextAction.ToolName {
	case "list_recent_orders":
		return r.buildRecentOrderSelectionReply(ctx, state, session, decision)
	case "query_logistics":
		return r.buildOrderToolReply(ctx, state, session, decision, taskCodeLogisticsQuery, officialIntentLogistics)
	case "get_order_snapshot":
		return r.buildOrderToolReply(ctx, state, session, decision, deriveOfficialSnapshotTaskCode(state, session, decision), decision.Intent)
	case "search_knowledge_chunks":
		return r.buildPolicyKnowledgeReply(ctx, state, session, decision)
	case "search_products", "batch_get_spu_cards":
		return r.buildCatalogGuideReply(ctx, state, session, decision)
	default:
		reply := buildDirectOfficialReply(state, session, heuristicPlannerDecision(state))
		return &reply, session, nil
	}
}

func (r *einoRunner) buildRecentOrderSelectionReply(
	ctx context.Context,
	state *runState,
	session TaskSessionState,
	decision PlannerDecision,
) (*ReplyPayload, TaskSessionState, error) {
	recentOrders, err := secureListRecentOrders(ctx, r.orderRepository, state.Security, defaultRecentOrderCandidateLimit)
	if err != nil {
		return nil, session, err
	}
	state.ToolResult = toolResult{
		ToolName:    state.SelectedToolName,
		AdapterCode: state.SelectedAdapter,
		Status:      toolResultSuccess,
	}

	taskCode := taskCodeLogisticsQuery
	if decision.Intent == officialIntentAfterSale {
		taskCode = detectTaskCode(state.Message.ContentText, taskCodeRefundDecision)
	}

	session.ActiveTaskCode = taskCode
	session.MissingSlots = []MissingSlot{{
		SlotCode:   slotCodeOrderNo,
		PromptText: "请选择具体订单，或者直接补充订单号。",
		Required:   true,
	}}
	session.PendingSelectionQuery = strings.TrimSpace(state.Message.ContentText)
	session.PendingSelectionTaskCode = taskCode
	session.PendingSelectionCandidates = rankRecentOrders(taskCode, inferProblemType(state.Message.ContentText, taskCode), state.Message.ContentText, recentOrders)

	reply := ReplyPayload{
		ReplyText:       buildOrderSelectionReplyText(taskCode),
		IntentCode:      decision.Intent,
		GuardResultCode: state.GuardResultCode,
		Confidence:      maxFloat(decision.Confidence, 0.9),
		DataCards: []ReplyDataCard{{
			OrderSelectionCard: &OrderSelectionCard{
				TitleText:     "找到您最近的订单了，点一个我继续处理",
				HelperText:    "我会沿着您刚才的问题继续往下查，不需要重新描述。",
				OriginalQuery: strings.TrimSpace(state.Message.ContentText),
				TaskCode:      taskCode,
				Candidates:    session.PendingSelectionCandidates,
			},
		}},
		MissingSlots:   session.MissingSlots,
		SlotRetryCount: session.SlotRetryCount,
		SuggestedActions: []SuggestedAction{{
			ActionCode: actionCodeEscalateToHuman,
			Label:      "转人工",
			Enabled:    true,
		}},
	}
	return &reply, session, nil
}

func (r *einoRunner) buildOrderToolReply(
	ctx context.Context,
	state *runState,
	session TaskSessionState,
	decision PlannerDecision,
	taskCode string,
	intentCode string,
) (*ReplyPayload, TaskSessionState, error) {
	orderNo := resolveSelectedOrderNo(session, state.Conversation)
	subOrderNo := resolveSelectedSubOrderNo(session, state.Conversation)
	if strings.TrimSpace(orderNo) == "" {
		orderNo = strings.TrimSpace(decision.NextAction.ToolArgs["order_no"])
	}
	if strings.TrimSpace(subOrderNo) == "" {
		subOrderNo = strings.TrimSpace(decision.NextAction.ToolArgs["sub_order_no"])
	}
	if strings.TrimSpace(orderNo) == "" {
		orderNo = extractOrderNumber(state.Message.ContentText, "ord")
	}
	if strings.TrimSpace(subOrderNo) == "" {
		subOrderNo = extractOrderNumber(state.Message.ContentText, "sub")
	}
	if strings.TrimSpace(orderNo) == "" {
		return r.buildRecentOrderSelectionReply(ctx, state, session, decision)
	}

	session.SelectedOrderNo = orderNo
	session.SelectedSubOrderNo = subOrderNo
	session.ActiveTaskCode = taskCode

	lookupResult, err := secureQueryOrderSnapshot(ctx, r.orderRepository, state.Security, OrderQuery{
		OrderNo:     orderNo,
		SubOrderNo:  subOrderNo,
		ProblemType: inferProblemType(state.Message.ContentText, taskCode),
	})
	if err != nil {
		return nil, session, err
	}
	if lookupResult.ErrorCode != "" {
		state.ToolResult = toolResult{
			ToolName:           state.SelectedToolName,
			AdapterCode:        state.SelectedAdapter,
			Status:             toolResultDegraded,
			DegradedReasonCode: lookupResult.ErrorCode,
			Message:            lookupResult.ErrorCode,
		}
		out := handleOrderLookupFailure(session, GuardDecision{GuardResultCode: state.GuardResultCode}, lookupResult.ErrorCode)
		out.Reply.IntentCode = intentCode
		return &out.Reply, out.Session, nil
	}

	state.ToolResult = toolResult{
		ToolName:    state.SelectedToolName,
		AdapterCode: state.SelectedAdapter,
		Status:      toolResultSuccess,
	}

	ruleEngine := r.ruleEngine
	if ruleEngine == nil {
		ruleEngine = NewAfterSaleRuleEngine(nil)
	}
	ruleOutcome := ruleEngine.Evaluate(AfterSaleRuleContext{
		TaskCode:    taskCode,
		ProblemType: inferProblemType(state.Message.ContentText, taskCode),
		Snapshot:    lookupResult.Snapshot,
		Session:     session,
	})

	session.SlotValues = ensureSessionSlots(session.SlotValues)
	session.SlotValues[slotCodeOrderNo] = orderNo
	if strings.TrimSpace(subOrderNo) != "" {
		session.SlotValues[slotCodeSubOrderNo] = subOrderNo
	}
	session.LatestFactsSummary = summarizeOrderSnapshot(lookupResult.Snapshot)
	session.LatestDecisionSummary = summarizeDecisionCard(ruleOutcome.DecisionCard, ruleOutcome.ReplyText)
	session.MissingSlots = nil
	session.PendingSelectionQuery = ""
	session.PendingSelectionTaskCode = ""
	session.PendingSelectionCandidates = nil
	session.HandoffRecommended = ruleOutcome.HandoffRecommended
	session.EscalationReasonCode = ruleOutcome.HandoffReasonCode

	reply := ReplyPayload{
		ReplyText:       ruleOutcome.ReplyText,
		IntentCode:      intentCode,
		GuardResultCode: state.GuardResultCode,
		Confidence:      maxFloat(decision.Confidence, 0.92),
		DataCards: []ReplyDataCard{
			{OrderSnapshotCard: toOrderSnapshotCard(lookupResult.Snapshot)},
		},
		SuggestedActions:   ruleOutcome.SuggestedActions,
		HandoffRecommended: ruleOutcome.HandoffRecommended,
		HandoffReasonCode:  ruleOutcome.HandoffReasonCode,
	}
	if intentCode == officialIntentLogistics {
		reply.DataCards = append(reply.DataCards, ReplyDataCard{
			LogisticsTrackingCard: buildLogisticsTrackingCard(lookupResult.Snapshot, ruleOutcome.ReplyText),
		})
	} else {
		reply.DataCards = append(reply.DataCards, ReplyDataCard{
			AfterSaleDecisionCard: ruleOutcome.DecisionCard,
		})
	}
	return &reply, session, nil
}

func (r *einoRunner) buildPolicyKnowledgeReply(
	ctx context.Context,
	state *runState,
	session TaskSessionState,
	decision PlannerDecision,
) (*ReplyPayload, TaskSessionState, error) {
	query := strings.TrimSpace(decision.NextAction.ToolArgs["query"])
	if query == "" {
		query = strings.TrimSpace(state.Message.ContentText)
	}

	session.ActiveIntent = officialIntentPolicyQA
	session.UserGoal = strings.TrimSpace(decision.UserGoal)

	if r.policyKnowledgeRepository == nil {
		state.ToolResult = toolResult{
			ToolName:           state.SelectedToolName,
			AdapterCode:        state.SelectedAdapter,
			Status:             toolResultDegraded,
			DegradedReasonCode: orderLookupErrorDownstream,
		}
		reply := ReplyPayload{
			ReplyText:          "我暂时无法查询到规则依据，建议稍后再试，或者我帮您转人工继续确认。",
			IntentCode:         officialIntentPolicyQA,
			GuardResultCode:    state.GuardResultCode,
			Confidence:         maxFloat(decision.Confidence, 0.7),
			HandoffRecommended: true,
			HandoffReasonCode:  escalationReasonDownstream,
			SuggestedActions: []SuggestedAction{{
				ActionCode: actionCodeEscalateToHuman,
				Label:      "转人工",
				Enabled:    true,
			}},
		}
		return &reply, session, nil
	}

	hits, err := r.policyKnowledgeRepository.SearchPolicyKnowledge(ctx, query, 3)
	if err != nil {
		state.ToolResult = toolResult{
			ToolName:           state.SelectedToolName,
			AdapterCode:        state.SelectedAdapter,
			Status:             toolResultDegraded,
			DegradedReasonCode: orderLookupErrorDownstream,
		}
		reply := ReplyPayload{
			ReplyText:          "规则查询当前暂时不可用，我先不乱给结论，建议稍后再试，或者我帮您转人工确认。",
			IntentCode:         officialIntentPolicyQA,
			GuardResultCode:    state.GuardResultCode,
			Confidence:         maxFloat(decision.Confidence, 0.7),
			HandoffRecommended: true,
			HandoffReasonCode:  escalationReasonDownstream,
			SuggestedActions: []SuggestedAction{{
				ActionCode: actionCodeEscalateToHuman,
				Label:      "转人工",
				Enabled:    true,
			}},
		}
		return &reply, session, nil
	}

	state.ToolResult = toolResult{
		ToolName:    state.SelectedToolName,
		AdapterCode: state.SelectedAdapter,
		Status:      toolResultSuccess,
	}
	state.AnswerSources = make([]AnswerSource, 0, len(hits))
	for _, hit := range hits {
		state.AnswerSources = append(state.AnswerSources, AnswerSource{
			SourceTypeCode: hit.SourceTypeCode,
			SourceID:       hit.SourceID,
			SourceVersion:  hit.SourceVersion,
			Title:          hit.Title,
			Snippet:        hit.Snippet,
		})
	}

	reply := ReplyPayload{
		ReplyText:       buildPolicyReplyText(hits),
		IntentCode:      officialIntentPolicyQA,
		GuardResultCode: state.GuardResultCode,
		Confidence:      maxFloat(decision.Confidence, 0.9),
	}
	if len(hits) == 0 {
		reply.SuggestedActions = []SuggestedAction{{
			ActionCode: actionCodeEscalateToHuman,
			Label:      "转人工",
			Enabled:    true,
		}}
	}
	return &reply, session, nil
}

func (r *einoRunner) buildCatalogGuideReply(
	ctx context.Context,
	state *runState,
	session TaskSessionState,
	decision PlannerDecision,
) (*ReplyPayload, TaskSessionState, error) {
	query := strings.TrimSpace(decision.NextAction.ToolArgs["query"])
	if query == "" {
		query = strings.TrimSpace(state.Message.ContentText)
	}

	session.ActiveIntent = officialIntentCatalogGuide
	session.UserGoal = strings.TrimSpace(decision.UserGoal)

	if r.productSearchRepository == nil {
		state.ToolResult = toolResult{
			ToolName:           state.SelectedToolName,
			AdapterCode:        state.SelectedAdapter,
			Status:             toolResultDegraded,
			DegradedReasonCode: orderLookupErrorDownstream,
		}
		reply := ReplyPayload{
			ReplyText:       "商品搜索当前有点忙，我暂时无法给您准确推荐。您可以稍后再试，或者补充更具体的需求给我。",
			IntentCode:      officialIntentCatalogGuide,
			GuardResultCode: state.GuardResultCode,
			Confidence:      maxFloat(decision.Confidence, 0.7),
		}
		return &reply, session, nil
	}

	items, err := r.productSearchRepository.SearchProducts(ctx, ProductSearchFilter{
		Query:  query,
		ShopNo: strings.TrimSpace(state.Security.ShopNo),
		Limit:  5,
	})
	if err != nil {
		state.ToolResult = toolResult{
			ToolName:           state.SelectedToolName,
			AdapterCode:        state.SelectedAdapter,
			Status:             toolResultDegraded,
			DegradedReasonCode: orderLookupErrorDownstream,
		}
		reply := ReplyPayload{
			ReplyText:       "商品搜索当前暂时不可用，我暂时无法给您准确推荐。您可以稍后再试，或者补充更具体的需求给我。",
			IntentCode:      officialIntentCatalogGuide,
			GuardResultCode: state.GuardResultCode,
			Confidence:      maxFloat(decision.Confidence, 0.7),
		}
		return &reply, session, nil
	}

	state.ToolResult = toolResult{
		ToolName:    state.SelectedToolName,
		AdapterCode: state.SelectedAdapter,
		Status:      toolResultSuccess,
	}

	card := &ProductRecommendationCard{
		TitleText:  "先给您挑了几款更贴近需求的商品",
		HelperText: "您可以先看看这几款，如果想更偏通勤、轻薄、价格友好或某个颜色，我可以继续缩小范围。",
		Items:      make([]ProductRecommendationItem, 0, len(items)),
	}
	for _, item := range items {
		card.Items = append(card.Items, ProductRecommendationItem{
			SpuNo:      item.SpuNo,
			Title:      item.Title,
			CoverURL:   item.CoverURL,
			PriceText:  formatPriceRangeText(item.MinPrice, item.MaxPrice),
			ShopName:   item.ShopName,
			ReasonText: item.ReasonText,
		})
	}

	replyText := "我先按您这次的需求筛了几款更接近的商品，您可以先看看这几款。"
	if len(card.Items) == 0 {
		replyText = "我暂时没有搜到特别贴近您这次需求的商品。您可以告诉我预算、颜色、尺码或更偏好的风格，我继续帮您缩小范围。"
	}

	reply := ReplyPayload{
		ReplyText:       replyText,
		IntentCode:      officialIntentCatalogGuide,
		GuardResultCode: state.GuardResultCode,
		Confidence:      maxFloat(decision.Confidence, 0.9),
		DataCards: []ReplyDataCard{{
			ProductRecommendationCard: card,
		}},
	}
	return &reply, session, nil
}

func buildDirectOfficialReply(state *runState, session TaskSessionState, decision PlannerDecision) ReplyPayload {
	replyText := strings.TrimSpace(decision.NextAction.Reason)
	if replyText == "" {
		replyText = "我可以帮您处理订单、物流、退款退货、规则说明和转人工问题。"
	}
	return ReplyPayload{
		ReplyText:          replyText,
		IntentCode:         decision.Intent,
		GuardResultCode:    state.GuardResultCode,
		Confidence:         maxFloat(decision.Confidence, 0.8),
		HandoffRecommended: decision.HandoffRecommended,
		HandoffReasonCode:  decision.HandoffReasonCode,
		SlotRetryCount:     session.SlotRetryCount,
	}
}

func deriveOfficialSnapshotTaskCode(state *runState, session TaskSessionState, decision PlannerDecision) string {
	taskCode := detectTaskCode(state.Message.ContentText, "")
	if taskCode != taskCodeActionExplanation {
		return taskCode
	}
	if strings.TrimSpace(session.PendingSelectionTaskCode) != "" {
		return strings.TrimSpace(session.PendingSelectionTaskCode)
	}
	if strings.TrimSpace(session.ActiveTaskCode) != "" {
		return strings.TrimSpace(session.ActiveTaskCode)
	}
	if decision.Intent == officialIntentOrderQuery || strings.TrimSpace(resolveSelectedOrderNo(session, state.Conversation)) != "" {
		return taskCodeOrderStatusQuery
	}
	return taskCodeOrderStatusQuery
}

func buildLogisticsTrackingCard(snapshot *OrderSnapshot, replyText string) *LogisticsTrackingCard {
	if snapshot == nil {
		return nil
	}
	return &LogisticsTrackingCard{
		OrderNo:           snapshot.OrderNo,
		SubOrderNo:        snapshot.SubOrderNo,
		FulfillmentStatus: snapshot.FulfillmentStatus,
		LogisticsStatus:   snapshot.LogisticsStatus,
		LatestUpdateTime:  snapshot.LatestUpdateTime,
		LatestTraceText:   strings.TrimSpace(replyText),
		TimelineSummary:   summarizeOrderSnapshot(snapshot),
	}
}

func ensureSessionSlots(slotValues map[string]string) map[string]string {
	if slotValues == nil {
		return make(map[string]string)
	}
	return slotValues
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func buildPolicyReplyText(hits []PolicyKnowledgeHit) string {
	if len(hits) == 0 {
		return "我暂时没有检索到足够明确的规则依据。您可以换个说法继续问我，或者我帮您转人工确认。"
	}
	primary := strings.TrimSpace(hits[0].Snippet)
	if primary == "" {
		primary = strings.TrimSpace(hits[0].Title)
	}
	if primary == "" {
		return "我先帮您查到了相关规则。"
	}
	return "目前查到的规则是：" + primary
}

func formatPriceRangeText(minPrice, maxPrice uint64) string {
	switch {
	case minPrice == 0 && maxPrice == 0:
		return ""
	case maxPrice == 0 || minPrice == maxPrice:
		return fmt.Sprintf("¥%.2f", float64(minPrice)/100)
	default:
		return fmt.Sprintf("¥%.2f - ¥%.2f", float64(minPrice)/100, float64(maxPrice)/100)
	}
}

func hasPolicyIntent(text string) bool {
	return containsAny(
		text,
		"规则", "政策", "平台规则", "售后规则", "退款规则", "退货规则", "换货规则",
		"退款流程", "退货流程", "怎么退款", "怎么退货", "怎么换货",
		"rule", "policy", "refund policy", "return policy",
	)
}

func hasCatalogGuideIntent(text string) bool {
	if containsAny(text, "订单", "物流", "退款", "退货", "换货", "售后", "人工客服") {
		return false
	}
	return containsAny(
		text,
		"推荐", "想买", "有没有", "适合", "选哪款", "帮我挑", "帮我推荐",
		"外套", "衣服", "裤子", "鞋", "包", "裙子",
		"recommend", "suggest", "buy", "outfit", "jacket",
	)
}
