package runtime

import (
	"context"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	intentCodeGreeting       = "Greeting"
	intentCodeOutOfScope     = "OutOfScope"
	intentCodeAttachmentOnly = "AttachmentOnly"
	intentCodeAfterSaleTask  = "AfterSaleTask"

	guardResultGreeting       = "GREETING"
	guardResultOutOfScope     = "OUT_OF_SCOPE"
	guardResultAttachmentOnly = "ATTACHMENT_ONLY"
	guardResultAfterSaleTask  = "AFTER_SALE_TASK"

	taskCodeOrderStatusQuery  = "order_status_query"
	taskCodeLogisticsQuery    = "logistics_query"
	taskCodeRefundDecision    = "aftersale_decision"
	taskCodeActionExplanation = "action_explanation"

	slotCodeOrderNo    = "order_no"
	slotCodeSubOrderNo = "sub_order_no"
	slotCodeProblem    = "problem_type"

	slotPromptCodeOrderNo        = "ask_order_no"
	slotPromptCodeOrderSelection = "select_recent_order"

	orderLookupErrorNotFound            = "NOT_FOUND"
	orderLookupErrorPermissionDenied    = "PERMISSION_DENIED"
	orderLookupErrorDownstream          = "DOWNSTREAM_UNAVAILABLE"
	orderLookupErrorInsufficientContext = "INSUFFICIENT_CONTEXT"

	escalationReasonUserRequested     = "user_requested_handoff"
	escalationReasonSlotFillingFailed = "slot_filling_failed"
	escalationReasonRuleConflict      = "rule_conflict"
	escalationReasonDownstream        = "downstream_unavailable"
	escalationReasonHighRiskCase      = "high_risk_case"

	decisionPathRefundOnly             = "refund_only"
	decisionPathWaitShipment           = "wait_for_shipment"
	decisionPathReturnRefundOrExchange = "return_refund_or_exchange"
	decisionPathWaitExistingAfterSale  = "wait_existing_after_sale"

	actionCodeRequestRefund   = "request_refund"
	actionCodeRequestReturn   = "request_return_refund"
	actionCodeRequestExchange = "request_exchange"
	actionCodeWaitForUpdate   = "wait_for_update"
	actionCodeEscalateToHuman = "escalate_to_human"
	actionCodeUrgeShipment    = "urge_shipment"
	actionCodeViewLogistics   = "view_logistics"
	actionCodeSelectOrder     = "select_order"

	defaultSlotRetryThreshold        = 2
	defaultUnresolvedThreshold       = 2
	defaultRecentOrderCandidateLimit = 5
)

var orderNumberPattern = regexp.MustCompile(`(?i)\b[A-Z]{2,}[A-Z0-9_-]{6,}\b`)

type UserMessage struct {
	ContentText     string   `json:"content_text,omitempty"`
	MessageTypeCode string   `json:"message_type_code,omitempty"`
	AssetIDs        []uint64 `json:"asset_ids,omitempty"`
}

type HiddenAction struct {
	Type  string `json:"type,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type GuardDecision struct {
	IntentCode      string `json:"intent_code,omitempty"`
	GuardResultCode string `json:"guard_result_code,omitempty"`
}

type MissingSlot struct {
	SlotCode   string `json:"slot_code,omitempty"`
	PromptText string `json:"prompt_text,omitempty"`
	Required   bool   `json:"required,omitempty"`
}

type SuggestedAction struct {
	ActionCode       string `json:"action_code,omitempty"`
	Label            string `json:"label,omitempty"`
	Enabled          bool   `json:"enabled,omitempty"`
	ReasonIfDisabled string `json:"reason_if_disabled,omitempty"`
}

type OrderSnapshotCard struct {
	OrderNo           string `json:"order_no,omitempty"`
	MainStatus        string `json:"main_status,omitempty"`
	PaymentStatus     string `json:"payment_status,omitempty"`
	FulfillmentStatus string `json:"fulfillment_status,omitempty"`
	LogisticsStatus   string `json:"logistics_status,omitempty"`
	AfterSaleStatus   string `json:"after_sale_status,omitempty"`
	LatestUpdateTime  string `json:"latest_update_time,omitempty"`
}

type RecentOrderCandidate struct {
	OrderNo           string `json:"order_no,omitempty"`
	SubOrderNo        string `json:"sub_order_no,omitempty"`
	DisplayTitle      string `json:"display_title,omitempty"`
	MainStatus        string `json:"main_status,omitempty"`
	PaymentStatus     string `json:"payment_status,omitempty"`
	FulfillmentStatus string `json:"fulfillment_status,omitempty"`
	LogisticsStatus   string `json:"logistics_status,omitempty"`
	AfterSaleStatus   string `json:"after_sale_status,omitempty"`
	LatestUpdateTime  string `json:"latest_update_time,omitempty"`
	SelectionHint     string `json:"selection_hint,omitempty"`
}

type OrderSelectionCard struct {
	TitleText     string                 `json:"title_text,omitempty"`
	HelperText    string                 `json:"helper_text,omitempty"`
	OriginalQuery string                 `json:"original_query,omitempty"`
	TaskCode      string                 `json:"task_code,omitempty"`
	Candidates    []RecentOrderCandidate `json:"candidates,omitempty"`
}

type AfterSaleDecisionCard struct {
	SceneCode        string `json:"scene_code,omitempty"`
	DecisionPathCode string `json:"decision_path_code,omitempty"`
	ReasonText       string `json:"reason_text,omitempty"`
	ConstraintText   string `json:"constraint_text,omitempty"`
	NextStepText     string `json:"next_step_text,omitempty"`
}

type LogisticsTrackingCard struct {
	OrderNo           string `json:"order_no,omitempty"`
	SubOrderNo        string `json:"sub_order_no,omitempty"`
	FulfillmentStatus string `json:"fulfillment_status,omitempty"`
	LogisticsStatus   string `json:"logistics_status,omitempty"`
	LatestUpdateTime  string `json:"latest_update_time,omitempty"`
	LatestTraceText   string `json:"latest_trace_text,omitempty"`
	TimelineSummary   string `json:"timeline_summary,omitempty"`
}

type ProductRecommendationItem struct {
	SpuNo      string `json:"spu_no,omitempty"`
	Title      string `json:"title,omitempty"`
	CoverURL   string `json:"cover_url,omitempty"`
	PriceText  string `json:"price_text,omitempty"`
	ShopName   string `json:"shop_name,omitempty"`
	ReasonText string `json:"reason_text,omitempty"`
}

type ProductRecommendationCard struct {
	TitleText  string                      `json:"title_text,omitempty"`
	HelperText string                      `json:"helper_text,omitempty"`
	Items      []ProductRecommendationItem `json:"items,omitempty"`
}

type ReplyDataCard struct {
	OrderSnapshotCard         *OrderSnapshotCard         `json:"order_snapshot_card,omitempty"`
	OrderSelectionCard        *OrderSelectionCard        `json:"order_selection_card,omitempty"`
	AfterSaleDecisionCard     *AfterSaleDecisionCard     `json:"after_sale_decision_card,omitempty"`
	LogisticsTrackingCard     *LogisticsTrackingCard     `json:"logistics_tracking_card,omitempty"`
	ProductRecommendationCard *ProductRecommendationCard `json:"product_recommendation_card,omitempty"`
}

type ReplyPayload struct {
	ReplyText          string            `json:"reply_text,omitempty"`
	IntentCode         string            `json:"intent_code,omitempty"`
	GuardResultCode    string            `json:"guard_result_code,omitempty"`
	Confidence         float64           `json:"confidence,omitempty"`
	DataCards          []ReplyDataCard   `json:"data_cards,omitempty"`
	SuggestedActions   []SuggestedAction `json:"suggested_actions,omitempty"`
	MissingSlots       []MissingSlot     `json:"missing_slots,omitempty"`
	SlotRetryCount     uint32            `json:"slot_retry_count,omitempty"`
	HandoffRecommended bool              `json:"handoff_recommended,omitempty"`
	HandoffReasonCode  string            `json:"handoff_reason_code,omitempty"`
}

type TaskSessionState struct {
	ConversationNo             string                 `json:"conversation_no,omitempty"`
	SessionVersion             uint64                 `json:"version,omitempty"`
	ActiveIntent               string                 `json:"active_intent,omitempty"`
	ActiveTaskCode             string                 `json:"active_task_code,omitempty"`
	UserGoal                   string                 `json:"user_goal,omitempty"`
	SelectedOrderNo            string                 `json:"selected_order_no,omitempty"`
	SelectedSubOrderNo         string                 `json:"selected_sub_order_no,omitempty"`
	SlotValues                 map[string]string      `json:"slot_values,omitempty"`
	MissingSlots               []MissingSlot          `json:"missing_slots,omitempty"`
	SlotRetryCount             uint32                 `json:"slot_retry_count,omitempty"`
	LastSlotPromptCode         string                 `json:"last_slot_prompt_code,omitempty"`
	SlotFillingFailed          bool                   `json:"slot_filling_failed,omitempty"`
	EscalationReasonCode       string                 `json:"escalation_reason_code,omitempty"`
	AnchoredOrderNo            string                 `json:"anchored_order_no,omitempty"`
	AnchoredSubOrderNo         string                 `json:"anchored_sub_order_no,omitempty"`
	LatestFactsSummary         string                 `json:"latest_facts_summary,omitempty"`
	LatestDecisionSummary      string                 `json:"latest_decision_summary,omitempty"`
	HandoffRecommended         bool                   `json:"handoff_recommended,omitempty"`
	UnresolvedTurnCount        uint32                 `json:"unresolved_turn_count,omitempty"`
	PendingSelectionQuery      string                 `json:"pending_selection_query,omitempty"`
	PendingSelectionTaskCode   string                 `json:"pending_selection_task_code,omitempty"`
	PendingSelectionCandidates []RecentOrderCandidate `json:"pending_selection_candidates,omitempty"`
}

type SecurityContext struct {
	UserID         uint64
	ShopNo         string
	RequestID      string
	ConversationNo string
	RunNo          string
}

type ConversationAnchors struct {
	OrderNo    string
	SubOrderNo string
}

type OrderQuery struct {
	OrderNo     string
	SubOrderNo  string
	ProblemType string
}

type OrderOwnershipFilter struct {
	UserID         uint64
	ShopNo         string
	OrderNo        string
	SubOrderNo     string
	RequestID      string
	ConversationNo string
	RunNo          string
}

type RecentOrderListFilter struct {
	UserID         uint64
	ShopNo         string
	RequestID      string
	ConversationNo string
	RunNo          string
	Limit          uint32
}

type OrderSnapshot struct {
	OrderNo            string
	SubOrderNo         string
	MainStatus         string
	PaymentStatus      string
	FulfillmentStatus  string
	LogisticsStatus    string
	AfterSaleStatus    string
	LatestUpdateTime   string
	OwnershipConfirmed bool
}

type OrderSnapshotRepository interface {
	QueryOrderSnapshot(ctx context.Context, filter OrderOwnershipFilter) (*OrderSnapshot, error)
	ListRecentOrders(ctx context.Context, filter RecentOrderListFilter) ([]RecentOrderCandidate, error)
}

type OrderLookupError struct {
	Code string
}

func (e *OrderLookupError) Error() string {
	return e.Code
}

type OrderLookupResult struct {
	Snapshot  *OrderSnapshot
	ErrorCode string
}

type AfterSaleTurnInput struct {
	Message         UserMessage
	Conversation    ConversationAnchors
	PreviousSession TaskSessionState
	Security        SecurityContext
	OrderRepository OrderSnapshotRepository
	RuleEngine      AfterSaleRuleEngine
}

type AfterSaleTurnOutput struct {
	Session TaskSessionState
	Reply   ReplyPayload
}

func detectGuardIntent(msg UserMessage) GuardDecision {
	text := strings.TrimSpace(msg.ContentText)
	lowerText := strings.ToLower(text)

	if text == "" && len(msg.AssetIDs) > 0 {
		return GuardDecision{IntentCode: intentCodeAttachmentOnly, GuardResultCode: guardResultAttachmentOnly}
	}
	if hasGreetingIntent(lowerText) {
		return GuardDecision{IntentCode: intentCodeGreeting, GuardResultCode: guardResultGreeting}
	}
	if hasAfterSaleIntent(lowerText) {
		return GuardDecision{IntentCode: intentCodeAfterSaleTask, GuardResultCode: guardResultAfterSaleTask}
	}
	return GuardDecision{IntentCode: intentCodeOutOfScope, GuardResultCode: guardResultOutOfScope}
}

func planAfterSaleTurn(ctx context.Context, in AfterSaleTurnInput) (AfterSaleTurnOutput, error) {
	session := copyTaskSession(in.PreviousSession)
	guard := detectGuardIntent(in.Message)
	if guard.IntentCode != intentCodeAfterSaleTask && shouldResumePendingSelection(session, in.Message.ContentText) {
		guard = GuardDecision{IntentCode: intentCodeAfterSaleTask, GuardResultCode: guardResultAfterSaleTask}
	}
	if guard.IntentCode != intentCodeAfterSaleTask {
		return AfterSaleTurnOutput{
			Session: session,
			Reply:   buildGuardReply(guard),
		}, nil
	}

	activeTaskFallback := session.ActiveTaskCode
	if strings.TrimSpace(activeTaskFallback) == "" {
		activeTaskFallback = session.PendingSelectionTaskCode
	}

	activeTaskCode := detectTaskCode(in.Message.ContentText, activeTaskFallback)
	session.ActiveTaskCode = activeTaskCode
	conversationAnchors := in.Conversation
	if shouldForceFreshOrderSelection(session, in.Conversation, in.Message.ContentText, activeTaskCode) {
		session = clearOrderSelectionContext(session)
		conversationAnchors = ConversationAnchors{}
	}
	session = mergeSlotValues(session, in.Message.ContentText, conversationAnchors, activeTaskCode)

	if containsAny(strings.ToLower(strings.TrimSpace(in.Message.ContentText)), "转人工", "人工", "人工客服") {
		session.HandoffRecommended = true
		session.EscalationReasonCode = escalationReasonUserRequested
		return AfterSaleTurnOutput{
			Session: session,
			Reply: ReplyPayload{
				ReplyText:          "已识别到你希望转人工，我可以把当前问题和已知订单上下文整理给人工客服继续处理。",
				IntentCode:         intentCodeAfterSaleTask,
				GuardResultCode:    guard.GuardResultCode,
				Confidence:         0.98,
				SlotRetryCount:     session.SlotRetryCount,
				HandoffRecommended: true,
				HandoffReasonCode:  escalationReasonUserRequested,
				SuggestedActions: []SuggestedAction{{
					ActionCode: actionCodeEscalateToHuman,
					Label:      "转人工",
					Enabled:    true,
				}},
			},
		}, nil
	}

	if strings.TrimSpace(session.SlotValues[slotCodeOrderNo]) == "" {
		return handleMissingOrderSlot(ctx, session, guard, in.OrderRepository, in.Security, in.Message.ContentText)
	}

	lookupResult, err := secureQueryOrderSnapshot(ctx, in.OrderRepository, in.Security, OrderQuery{
		OrderNo:     session.SlotValues[slotCodeOrderNo],
		SubOrderNo:  session.SlotValues[slotCodeSubOrderNo],
		ProblemType: session.SlotValues[slotCodeProblem],
	})
	if err != nil {
		return AfterSaleTurnOutput{}, err
	}
	if lookupResult.ErrorCode != "" {
		return handleOrderLookupFailure(session, guard, lookupResult.ErrorCode), nil
	}

	ruleEngine := in.RuleEngine
	if ruleEngine == nil {
		ruleEngine = NewAfterSaleRuleEngine(nil)
	}

	ruleOutcome := ruleEngine.Evaluate(AfterSaleRuleContext{
		TaskCode:    activeTaskCode,
		ProblemType: session.SlotValues[slotCodeProblem],
		Snapshot:    lookupResult.Snapshot,
		Session:     session,
	})

	session.LatestFactsSummary = summarizeOrderSnapshot(lookupResult.Snapshot)
	session.LatestDecisionSummary = summarizeDecisionCard(ruleOutcome.DecisionCard, ruleOutcome.ReplyText)
	session.MissingSlots = nil
	session.SlotRetryCount = 0
	session.LastSlotPromptCode = ""
	session.PendingSelectionQuery = ""
	session.PendingSelectionTaskCode = ""
	session.PendingSelectionCandidates = nil
	session.HandoffRecommended = ruleOutcome.HandoffRecommended
	session.SlotFillingFailed = false
	session.EscalationReasonCode = ruleOutcome.HandoffReasonCode
	if ruleOutcome.HandoffRecommended {
		session.UnresolvedTurnCount++
	} else {
		session.UnresolvedTurnCount = 0
	}

	reply := ReplyPayload{
		ReplyText:       ruleOutcome.ReplyText,
		IntentCode:      intentCodeAfterSaleTask,
		GuardResultCode: guard.GuardResultCode,
		Confidence:      0.92,
		DataCards: []ReplyDataCard{
			{OrderSnapshotCard: toOrderSnapshotCard(lookupResult.Snapshot)},
			{AfterSaleDecisionCard: ruleOutcome.DecisionCard},
		},
		SuggestedActions:   ruleOutcome.SuggestedActions,
		SlotRetryCount:     session.SlotRetryCount,
		HandoffRecommended: ruleOutcome.HandoffRecommended,
		HandoffReasonCode:  ruleOutcome.HandoffReasonCode,
	}

	return AfterSaleTurnOutput{Session: session, Reply: reply}, nil
}

func handleMissingOrderSlot(
	ctx context.Context,
	session TaskSessionState,
	guard GuardDecision,
	repo OrderSnapshotRepository,
	security SecurityContext,
	query string,
) (AfterSaleTurnOutput, error) {
	recentOrders, err := secureListRecentOrders(ctx, repo, security, defaultRecentOrderCandidateLimit)
	if err != nil {
		return AfterSaleTurnOutput{}, err
	}
	rankedOrders := rankRecentOrders(session.ActiveTaskCode, session.SlotValues[slotCodeProblem], query, recentOrders)
	if len(rankedOrders) > 0 {
		session.PendingSelectionQuery = strings.TrimSpace(query)
		session.PendingSelectionTaskCode = session.ActiveTaskCode
		session.PendingSelectionCandidates = append([]RecentOrderCandidate(nil), rankedOrders...)
		session.MissingSlots = []MissingSlot{{
			SlotCode:   slotCodeOrderNo,
			PromptText: "请选择具体订单，或者直接补充订单号。",
			Required:   true,
		}}
		session.LastSlotPromptCode = slotPromptCodeOrderSelection

		return AfterSaleTurnOutput{
			Session: session,
			Reply: ReplyPayload{
				ReplyText:       buildOrderSelectionReplyText(session.ActiveTaskCode),
				IntentCode:      intentCodeAfterSaleTask,
				GuardResultCode: guard.GuardResultCode,
				Confidence:      0.9,
				DataCards: []ReplyDataCard{{
					OrderSelectionCard: &OrderSelectionCard{
						TitleText:     "找到你最近的订单了，点一个我继续处理",
						HelperText:    "我会沿着你刚才的问题继续往下查，不需要重新描述。",
						OriginalQuery: strings.TrimSpace(query),
						TaskCode:      session.ActiveTaskCode,
						Candidates:    rankedOrders,
					},
				}},
				MissingSlots:   session.MissingSlots,
				SlotRetryCount: session.SlotRetryCount,
				SuggestedActions: []SuggestedAction{{
					ActionCode: actionCodeEscalateToHuman,
					Label:      "转人工",
					Enabled:    true,
				}},
			},
		}, nil
	}

	if session.LastSlotPromptCode == slotPromptCodeOrderNo {
		session.SlotRetryCount++
	} else {
		session.SlotRetryCount = 1
	}
	session.MissingSlots = []MissingSlot{{
		SlotCode:   slotCodeOrderNo,
		PromptText: "请提供订单号，或者告诉我是哪个订单。",
		Required:   true,
	}}
	session.LastSlotPromptCode = slotPromptCodeOrderNo
	session.UnresolvedTurnCount++

	if session.SlotRetryCount >= defaultSlotRetryThreshold {
		session.SlotFillingFailed = true
		session.EscalationReasonCode = escalationReasonSlotFillingFailed
		session.HandoffRecommended = true
		return AfterSaleTurnOutput{
			Session: session,
			Reply: ReplyPayload{
				ReplyText:          "我暂时还没拿到足够的订单信息，继续反复追问会影响体验，建议直接转人工帮你核实。",
				IntentCode:         intentCodeAfterSaleTask,
				GuardResultCode:    guard.GuardResultCode,
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
			},
		}, nil
	}

	return AfterSaleTurnOutput{
		Session: session,
		Reply: ReplyPayload{
			ReplyText:       "为了准确帮你判断退款、退货退款还是继续等待，我需要先定位到具体订单。",
			IntentCode:      intentCodeAfterSaleTask,
			GuardResultCode: guard.GuardResultCode,
			Confidence:      0.86,
			MissingSlots:    session.MissingSlots,
			SlotRetryCount:  session.SlotRetryCount,
		},
	}, nil
}

func secureQueryOrderSnapshot(
	ctx context.Context,
	repo OrderSnapshotRepository,
	security SecurityContext,
	query OrderQuery,
) (OrderLookupResult, error) {
	if security.UserID == 0 || strings.TrimSpace(query.OrderNo) == "" {
		return OrderLookupResult{ErrorCode: orderLookupErrorInsufficientContext}, nil
	}
	if repo == nil {
		return OrderLookupResult{ErrorCode: orderLookupErrorDownstream}, nil
	}

	filter := OrderOwnershipFilter{
		UserID:         security.UserID,
		ShopNo:         strings.TrimSpace(security.ShopNo),
		OrderNo:        strings.TrimSpace(query.OrderNo),
		SubOrderNo:     strings.TrimSpace(query.SubOrderNo),
		RequestID:      strings.TrimSpace(security.RequestID),
		ConversationNo: strings.TrimSpace(security.ConversationNo),
		RunNo:          strings.TrimSpace(security.RunNo),
	}

	snapshot, err := repo.QueryOrderSnapshot(ctx, filter)
	if err != nil {
		if lookupErr, ok := err.(*OrderLookupError); ok {
			return OrderLookupResult{ErrorCode: lookupErr.Code}, nil
		}
		return OrderLookupResult{ErrorCode: orderLookupErrorDownstream}, nil
	}
	if snapshot == nil {
		return OrderLookupResult{ErrorCode: orderLookupErrorNotFound}, nil
	}
	if !snapshot.OwnershipConfirmed {
		return OrderLookupResult{ErrorCode: orderLookupErrorPermissionDenied}, nil
	}
	return OrderLookupResult{Snapshot: snapshot}, nil
}

func secureListRecentOrders(
	ctx context.Context,
	repo OrderSnapshotRepository,
	security SecurityContext,
	limit uint32,
) ([]RecentOrderCandidate, error) {
	if security.UserID == 0 || repo == nil {
		return nil, nil
	}

	candidates, err := repo.ListRecentOrders(ctx, RecentOrderListFilter{
		UserID:         security.UserID,
		ShopNo:         strings.TrimSpace(security.ShopNo),
		RequestID:      strings.TrimSpace(security.RequestID),
		ConversationNo: strings.TrimSpace(security.ConversationNo),
		RunNo:          strings.TrimSpace(security.RunNo),
		Limit:          limit,
	})
	if err != nil || len(candidates) == 0 {
		return nil, nil
	}
	return candidates, nil
}

func handleOrderLookupFailure(session TaskSessionState, guard GuardDecision, errorCode string) AfterSaleTurnOutput {
	session.UnresolvedTurnCount++

	if errorCode == orderLookupErrorNotFound || errorCode == orderLookupErrorPermissionDenied {
		reply := ReplyPayload{
			ReplyText:       "暂未查询到与你当前账号匹配的订单信息，请核对订单号后再试。",
			IntentCode:      intentCodeAfterSaleTask,
			GuardResultCode: guard.GuardResultCode,
			Confidence:      0.84,
		}
		session, reply = applyUnresolvedTurnEscalation(session, reply)
		return AfterSaleTurnOutput{Session: session, Reply: reply}
	}

	if errorCode == orderLookupErrorDownstream {
		session.HandoffRecommended = true
		session.EscalationReasonCode = escalationReasonDownstream
		return AfterSaleTurnOutput{
			Session: session,
			Reply: ReplyPayload{
				ReplyText:          "订单系统当前暂时不可用，我先不乱给结论，建议转人工继续核实。",
				IntentCode:         intentCodeAfterSaleTask,
				GuardResultCode:    guard.GuardResultCode,
				Confidence:         0.8,
				HandoffRecommended: true,
				HandoffReasonCode:  escalationReasonDownstream,
				SuggestedActions: []SuggestedAction{{
					ActionCode: actionCodeEscalateToHuman,
					Label:      "转人工",
					Enabled:    true,
				}},
			},
		}
	}

	reply := ReplyPayload{
		ReplyText:       "我还缺少足够的信息来准确判断，请补充订单号或更具体的问题描述。",
		IntentCode:      intentCodeAfterSaleTask,
		GuardResultCode: guard.GuardResultCode,
		Confidence:      0.78,
	}
	session, reply = applyUnresolvedTurnEscalation(session, reply)
	return AfterSaleTurnOutput{Session: session, Reply: reply}
}

func buildOrderSelectionReplyText(taskCode string) string {
	switch strings.TrimSpace(taskCode) {
	case taskCodeLogisticsQuery:
		return "我先把你最近的订单拉出来，点一个我马上继续查物流和发货进度。"
	case taskCodeRefundDecision:
		return "我先把你最近的订单拉出来，点一个我继续帮你判断退款、退货退款还是继续等待。"
	default:
		return "我先把你最近的订单拉出来，点一个我继续往下处理。"
	}
}

func shouldResumePendingSelection(session TaskSessionState, query string) bool {
	if strings.TrimSpace(session.PendingSelectionQuery) == "" && len(session.PendingSelectionCandidates) == 0 {
		return false
	}

	orderNo := extractOrderNumber(query, "ord")
	if strings.TrimSpace(orderNo) == "" {
		return false
	}
	if len(session.PendingSelectionCandidates) == 0 {
		return true
	}
	for _, candidate := range session.PendingSelectionCandidates {
		if strings.EqualFold(strings.TrimSpace(candidate.OrderNo), strings.TrimSpace(orderNo)) {
			return true
		}
	}
	return false
}

func rankRecentOrders(taskCode, problemType, query string, candidates []RecentOrderCandidate) []RecentOrderCandidate {
	if len(candidates) == 0 {
		return nil
	}

	type scoredCandidate struct {
		candidate RecentOrderCandidate
		score     int
		updatedAt time.Time
	}

	scored := make([]scoredCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		scored = append(scored, scoredCandidate{
			candidate: enrichRecentOrderCandidateHint(taskCode, problemType, candidate),
			score:     scoreRecentOrderCandidate(taskCode, problemType, query, candidate),
			updatedAt: parseCandidateTime(candidate.LatestUpdateTime),
		})
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].score != scored[j].score {
			return scored[i].score > scored[j].score
		}
		return scored[i].updatedAt.After(scored[j].updatedAt)
	})

	ranked := make([]RecentOrderCandidate, 0, len(scored))
	for _, item := range scored {
		ranked = append(ranked, item.candidate)
	}
	return ranked
}

func scoreRecentOrderCandidate(taskCode, problemType, query string, candidate RecentOrderCandidate) int {
	intent := strings.TrimSpace(problemType)
	if intent == "" || intent == problemTypeUnknown {
		intent = inferProblemType(query, taskCode)
	}

	score := 0
	switch intent {
	case problemTypeUrgeShipment:
		if isRecentOrderUnshipped(candidate) {
			score += 120
		}
		if isRecentOrderInTransit(candidate) {
			score -= 20
		}
	case problemTypeRefund, problemTypeReturnRefund:
		if strings.EqualFold(strings.TrimSpace(candidate.AfterSaleStatus), "PROCESSING") {
			score += 130
		}
		if isRecentOrderUnshipped(candidate) {
			score += 120
		}
		if isRecentOrderDelivered(candidate) {
			score += 90
		}
		if isRecentOrderInTransit(candidate) {
			score += 70
		}
	case problemTypeLogistics, problemTypeOrderStatus:
		if isRecentOrderInTransit(candidate) {
			score += 120
		}
		if isRecentOrderDelivered(candidate) {
			score += 90
		}
		if isRecentOrderUnshipped(candidate) {
			score += 40
		}
	default:
		if isRecentOrderInTransit(candidate) {
			score += 80
		}
		if isRecentOrderUnshipped(candidate) {
			score += 60
		}
	}

	if strings.TrimSpace(candidate.DisplayTitle) != "" {
		score++
	}
	return score
}

func enrichRecentOrderCandidateHint(taskCode, problemType string, candidate RecentOrderCandidate) RecentOrderCandidate {
	if strings.TrimSpace(candidate.SelectionHint) != "" {
		return candidate
	}

	intent := strings.TrimSpace(problemType)
	if intent == "" || intent == problemTypeUnknown {
		intent = inferProblemType("", taskCode)
	}

	switch intent {
	case problemTypeUrgeShipment:
		if isRecentOrderUnshipped(candidate) {
			candidate.SelectionHint = "这单还没发货，适合继续催发货或直接退款。"
		}
	case problemTypeRefund, problemTypeReturnRefund:
		if strings.EqualFold(strings.TrimSpace(candidate.AfterSaleStatus), "PROCESSING") {
			candidate.SelectionHint = "这单已经有售后在处理中，适合先看当前售后进展。"
		} else if isRecentOrderUnshipped(candidate) {
			candidate.SelectionHint = "这单还没发货，更适合直接申请退款。"
		} else if isRecentOrderDelivered(candidate) {
			candidate.SelectionHint = "这单已经签收，更适合判断退货退款或换货。"
		}
	default:
		if isRecentOrderInTransit(candidate) {
			candidate.SelectionHint = "这单正在运输中，适合继续查看物流进度。"
		} else if isRecentOrderUnshipped(candidate) {
			candidate.SelectionHint = "这单还没发货，适合先看发货状态。"
		}
	}

	return candidate
}

func parseCandidateTime(raw string) time.Time {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}
	}

	formats := []string{time.RFC3339, "2006-01-02 15:04:05", time.RFC3339Nano}
	for _, format := range formats {
		if parsed, err := time.Parse(format, raw); err == nil {
			return parsed
		}
	}
	return time.Time{}
}

func isRecentOrderUnshipped(candidate RecentOrderCandidate) bool {
	return strings.EqualFold(strings.TrimSpace(candidate.FulfillmentStatus), "UNSHIPPED") ||
		strings.EqualFold(strings.TrimSpace(candidate.LogisticsStatus), "NOT_SHIPPED")
}

func isRecentOrderInTransit(candidate RecentOrderCandidate) bool {
	return strings.EqualFold(strings.TrimSpace(candidate.FulfillmentStatus), "SHIPPED") ||
		strings.EqualFold(strings.TrimSpace(candidate.LogisticsStatus), "IN_TRANSIT")
}

func isRecentOrderDelivered(candidate RecentOrderCandidate) bool {
	return strings.EqualFold(strings.TrimSpace(candidate.FulfillmentStatus), "DELIVERED") ||
		strings.EqualFold(strings.TrimSpace(candidate.LogisticsStatus), "DELIVERED") ||
		strings.EqualFold(strings.TrimSpace(candidate.LogisticsStatus), "SIGNED")
}

func buildGuardReply(guard GuardDecision) ReplyPayload {
	if guard.IntentCode == intentCodeGreeting {
		return ReplyPayload{
			ReplyText:       "你好，我是订单售后助手，可以帮你查订单状态、物流进度和退款退货建议。",
			IntentCode:      guard.IntentCode,
			GuardResultCode: guard.GuardResultCode,
			Confidence:      0.99,
		}
	}
	if guard.IntentCode == intentCodeAttachmentOnly {
		return ReplyPayload{
			ReplyText:       "我可以帮你处理订单和售后问题，请再补充订单号或问题描述。",
			IntentCode:      guard.IntentCode,
			GuardResultCode: guard.GuardResultCode,
			Confidence:      0.95,
		}
	}
	return ReplyPayload{
		ReplyText:       "抱歉，我目前主要负责订单和售后问题。如果你要查询订单、物流或退款退货，我可以继续帮你处理。",
		IntentCode:      guard.IntentCode,
		GuardResultCode: guard.GuardResultCode,
		Confidence:      0.93,
	}
}

func detectTaskCode(query, fallback string) string {
	lowerText := strings.ToLower(strings.TrimSpace(query))

	if hasShipmentStatusIntent(lowerText) {
		return taskCodeOrderStatusQuery
	}
	if hasLogisticsIntent(lowerText) || hasUrgeShipmentIntent(lowerText) {
		return taskCodeLogisticsQuery
	}
	if hasOrderStatusIntent(lowerText) {
		return taskCodeOrderStatusQuery
	}
	if hasRefundIntent(lowerText) || hasReturnRefundIntent(lowerText) || hasExchangeIntent(lowerText) {
		return taskCodeRefundDecision
	}
	if strings.TrimSpace(fallback) != "" {
		return strings.TrimSpace(fallback)
	}
	return taskCodeActionExplanation
}

func shouldForceFreshOrderSelection(
	session TaskSessionState,
	anchors ConversationAnchors,
	query string,
	taskCode string,
) bool {
	if strings.TrimSpace(extractOrderNumber(query, "ord")) != "" {
		return false
	}
	if strings.TrimSpace(session.SlotValues[slotCodeOrderNo]) == "" &&
		strings.TrimSpace(session.AnchoredOrderNo) == "" &&
		strings.TrimSpace(anchors.OrderNo) == "" &&
		strings.TrimSpace(anchors.SubOrderNo) == "" {
		return false
	}
	return isFreshOrderQuickPrompt(query, taskCode)
}

func isFreshOrderQuickPrompt(query, taskCode string) bool {
	normalized := normalizePromptText(query)
	if normalized == "" {
		return false
	}

	switch strings.TrimSpace(taskCode) {
	case taskCodeLogisticsQuery:
		return matchesAnyNormalizedPrompt(normalized,
			"帮我查一下这个订单现在到哪了",
			"帮我查一下这个订单的物流进度",
			"帮我同步一下这个订单的最新进展",
			"这个订单怎么还没发货",
			"帮我看下这个订单为什么还没发货",
			"checkwherethisorderisnow",
			"checkthelogisticsprogressforthisorder",
			"sharethelatestupdateforthisorder",
			"whyhasthisordernotshippedyet",
		)
	case taskCodeRefundDecision:
		return matchesAnyNormalizedPrompt(normalized,
			"这个订单现在能退款吗",
			"这个订单现在更适合退货退款吗",
			"这个订单现在可以申请换货吗",
			"canthisorderberefundednow",
			"shouldthisorderusereturnandrefundnow",
			"canirequestanexchangeforthisorder",
		)
	case taskCodeOrderStatusQuery:
		return matchesAnyNormalizedPrompt(normalized,
			"帮我同步一下这个订单的最新进展",
			"sharethelatestupdateforthisorder",
		)
	default:
		return false
	}
}

func normalizePromptText(raw string) string {
	replacer := strings.NewReplacer(
		" ", "",
		"\t", "",
		"\n", "",
		"\r", "",
		"，", "",
		",", "",
		"。", "",
		".", "",
		"？", "",
		"?", "",
		"！", "",
		"!", "",
		"：", "",
		":", "",
	)
	return strings.ToLower(strings.TrimSpace(replacer.Replace(raw)))
}

func matchesAnyNormalizedPrompt(normalized string, prompts ...string) bool {
	for _, prompt := range prompts {
		if normalized == normalizePromptText(prompt) {
			return true
		}
	}
	return false
}

func clearOrderSelectionContext(session TaskSessionState) TaskSessionState {
	if session.SlotValues != nil {
		delete(session.SlotValues, slotCodeOrderNo)
		delete(session.SlotValues, slotCodeSubOrderNo)
	}
	session.AnchoredOrderNo = ""
	session.AnchoredSubOrderNo = ""
	session.SelectedOrderNo = ""
	session.SelectedSubOrderNo = ""
	session.LatestFactsSummary = ""
	session.LatestDecisionSummary = ""
	return session
}

func hasGreetingIntent(text string) bool {
	return containsAny(text, "在吗", "你好", "您好", "有人吗", "hello", "hi", "hey")
}

func hasRefundIntent(text string) bool {
	return containsAny(text, "退款", "退了", "能退吗", "还能退吗", "能不能退", "可以退吗", "退一个", "退款吗")
}

func hasReturnRefundIntent(text string) bool {
	return containsAny(text, "退货退款", "退货", "退回去", "寄回", "寄回去", "退货吗")
}

func hasExchangeIntent(text string) bool {
	return containsAny(text, "换货", "换一件", "换一个", "换吗", "换新")
}

func hasUrgeShipmentIntent(text string) bool {
	return containsAny(
		text,
		"催发货", "还没发货", "怎么还不发货", "怎么还没发货", "为什么还没发货",
		"什么时候发货", "啥时候发货", "何时发货", "发没发", "还不发",
	)
}

func hasLogisticsIntent(text string) bool {
	return containsAny(
		text,
		"物流", "快递", "包裹", "到哪", "到哪了", "到哪儿了", "在哪", "走到哪",
		"查物流", "查快递", "配送", "运输", "运到哪", "进度",
	)
}

func hasOrderStatusIntent(text string) bool {
	return containsAny(
		text,
		"订单状态", "订单进度", "什么状态", "现在怎么样", "当前状态", "订单怎么样",
	)
}

func hasAfterSaleIntent(text string) bool {
	if containsAny(text, "订单", "售后", "人工", "客服") {
		return true
	}
	return hasRefundIntent(text) ||
		hasReturnRefundIntent(text) ||
		hasExchangeIntent(text) ||
		hasUrgeShipmentIntent(text) ||
		hasLogisticsIntent(text) ||
		hasOrderStatusIntent(text)
}

func hasShipmentStatusIntent(text string) bool {
	return containsAny(
		text,
		"发货了吗",
		"是否发货",
		"有没有发货",
		"发没发货",
		"发了没有",
		"是否已经发货",
	)
}

func mergeSlotValues(session TaskSessionState, query string, anchors ConversationAnchors, activeTaskCode string) TaskSessionState {
	if session.SlotValues == nil {
		session.SlotValues = make(map[string]string)
	}

	if strings.TrimSpace(session.SlotValues[slotCodeOrderNo]) == "" && strings.TrimSpace(session.SelectedOrderNo) != "" {
		session.SlotValues[slotCodeOrderNo] = strings.TrimSpace(session.SelectedOrderNo)
	}
	if strings.TrimSpace(session.SlotValues[slotCodeSubOrderNo]) == "" && strings.TrimSpace(session.SelectedSubOrderNo) != "" {
		session.SlotValues[slotCodeSubOrderNo] = strings.TrimSpace(session.SelectedSubOrderNo)
	}
	if strings.TrimSpace(session.SlotValues[slotCodeOrderNo]) == "" && strings.TrimSpace(anchors.OrderNo) != "" {
		session.SlotValues[slotCodeOrderNo] = strings.TrimSpace(anchors.OrderNo)
	}
	if strings.TrimSpace(session.SlotValues[slotCodeSubOrderNo]) == "" && strings.TrimSpace(anchors.SubOrderNo) != "" {
		session.SlotValues[slotCodeSubOrderNo] = strings.TrimSpace(anchors.SubOrderNo)
	}
	if orderNo := extractOrderNumber(query, "ord"); orderNo != "" {
		session.SlotValues[slotCodeOrderNo] = orderNo
		session.AnchoredOrderNo = orderNo
		session.SelectedOrderNo = orderNo
	}
	if subOrderNo := extractOrderNumber(query, "sub"); subOrderNo != "" {
		session.SlotValues[slotCodeSubOrderNo] = subOrderNo
		session.AnchoredSubOrderNo = subOrderNo
		session.SelectedSubOrderNo = subOrderNo
	}

	if orderNo := strings.TrimSpace(session.SlotValues[slotCodeOrderNo]); orderNo != "" {
		session.SelectedOrderNo = orderNo
	}
	if subOrderNo := strings.TrimSpace(session.SlotValues[slotCodeSubOrderNo]); subOrderNo != "" {
		session.SelectedSubOrderNo = subOrderNo
	}
	session.SlotValues[slotCodeProblem] = inferProblemType(query, activeTaskCode)
	return session
}

func applyHiddenActionToTaskSession(session TaskSessionState, action HiddenAction) TaskSessionState {
	if session.SlotValues == nil {
		session.SlotValues = make(map[string]string)
	}

	normalizedType := strings.ToUpper(strings.TrimSpace(action.Type))
	normalizedKey := strings.ToLower(strings.TrimSpace(action.Key))
	normalizedValue := strings.TrimSpace(action.Value)
	if normalizedType != "SET_SLOT" || normalizedValue == "" {
		return session
	}

	switch normalizedKey {
	case "selected_order_no":
		session.SelectedOrderNo = normalizedValue
		session.AnchoredOrderNo = normalizedValue
		session.SlotValues[slotCodeOrderNo] = normalizedValue
	case "selected_sub_order_no":
		session.SelectedSubOrderNo = normalizedValue
		session.AnchoredSubOrderNo = normalizedValue
		session.SlotValues[slotCodeSubOrderNo] = normalizedValue
	default:
		return session
	}

	return session
}

func extractOrderNumber(query, prefix string) string {
	matches := orderNumberPattern.FindAllString(strings.ToUpper(strings.TrimSpace(query)), -1)
	for _, item := range matches {
		if strings.HasPrefix(item, strings.ToUpper(prefix)) {
			return item
		}
	}
	if len(matches) == 1 {
		return matches[0]
	}
	return ""
}

func decideAfterSalePath(taskCode string, snapshot *OrderSnapshot) (*AfterSaleDecisionCard, []SuggestedAction, string) {
	outcome := NewAfterSaleRuleEngine(nil).Evaluate(AfterSaleRuleContext{
		TaskCode:    taskCode,
		ProblemType: inferProblemType("", taskCode),
		Snapshot:    snapshot,
	})
	return outcome.DecisionCard, outcome.SuggestedActions, outcome.ReplyText
}

func summarizeDecisionCard(card *AfterSaleDecisionCard, fallback string) string {
	if card == nil {
		return strings.TrimSpace(fallback)
	}
	if strings.TrimSpace(card.ReasonText) != "" {
		return strings.TrimSpace(card.ReasonText)
	}
	return strings.TrimSpace(fallback)
}

func applyUnresolvedTurnEscalation(session TaskSessionState, reply ReplyPayload) (TaskSessionState, ReplyPayload) {
	if session.UnresolvedTurnCount < defaultUnresolvedThreshold || reply.HandoffRecommended {
		return session, reply
	}

	session.HandoffRecommended = true
	if strings.TrimSpace(session.EscalationReasonCode) == "" {
		session.EscalationReasonCode = escalationReasonHighRiskCase
	}
	reply.HandoffRecommended = true
	reply.HandoffReasonCode = session.EscalationReasonCode
	reply.SuggestedActions = append(reply.SuggestedActions, SuggestedAction{
		ActionCode: actionCodeEscalateToHuman,
		Label:      "转人工",
		Enabled:    true,
	})
	return session, reply
}

func summarizeOrderSnapshot(snapshot *OrderSnapshot) string {
	if snapshot == nil {
		return ""
	}
	return strings.Join([]string{
		"order_no=" + strings.TrimSpace(snapshot.OrderNo),
		"main_status=" + strings.TrimSpace(snapshot.MainStatus),
		"payment_status=" + strings.TrimSpace(snapshot.PaymentStatus),
		"fulfillment_status=" + strings.TrimSpace(snapshot.FulfillmentStatus),
		"logistics_status=" + strings.TrimSpace(snapshot.LogisticsStatus),
		"after_sale_status=" + strings.TrimSpace(snapshot.AfterSaleStatus),
	}, ", ")
}

func toOrderSnapshotCard(snapshot *OrderSnapshot) *OrderSnapshotCard {
	if snapshot == nil {
		return nil
	}
	return &OrderSnapshotCard{
		OrderNo:           snapshot.OrderNo,
		MainStatus:        snapshot.MainStatus,
		PaymentStatus:     snapshot.PaymentStatus,
		FulfillmentStatus: snapshot.FulfillmentStatus,
		LogisticsStatus:   snapshot.LogisticsStatus,
		AfterSaleStatus:   snapshot.AfterSaleStatus,
		LatestUpdateTime:  snapshot.LatestUpdateTime,
	}
}

func copyTaskSession(in TaskSessionState) TaskSessionState {
	out := in
	if in.SlotValues != nil {
		out.SlotValues = make(map[string]string, len(in.SlotValues))
		for key, value := range in.SlotValues {
			out.SlotValues[key] = value
		}
	}
	if len(in.MissingSlots) > 0 {
		out.MissingSlots = append([]MissingSlot(nil), in.MissingSlots...)
	}
	if len(in.PendingSelectionCandidates) > 0 {
		out.PendingSelectionCandidates = append([]RecentOrderCandidate(nil), in.PendingSelectionCandidates...)
	}
	return out
}
