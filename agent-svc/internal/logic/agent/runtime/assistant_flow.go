package runtime

import (
	"context"
	"regexp"
	"strings"
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

	slotPromptCodeOrderNo = "ask_order_no"

	orderLookupErrorNotFound            = "NOT_FOUND"
	orderLookupErrorPermissionDenied    = "PERMISSION_DENIED"
	orderLookupErrorDownstream          = "DOWNSTREAM_UNAVAILABLE"
	orderLookupErrorInsufficientContext = "INSUFFICIENT_CONTEXT"

	escalationReasonUserRequested       = "user_requested_handoff"
	escalationReasonSlotFillingFailed   = "slot_filling_failed"
	escalationReasonRuleConflict        = "rule_conflict"
	escalationReasonDownstream          = "downstream_unavailable"
	escalationReasonHighRiskCase        = "high_risk_case"

	decisionPathRefundOnly            = "refund_only"
	decisionPathWaitShipment          = "wait_for_shipment"
	decisionPathReturnRefundOrExchange = "return_refund_or_exchange"
	decisionPathWaitExistingAfterSale = "wait_existing_after_sale"

	actionCodeRequestRefund    = "request_refund"
	actionCodeRequestReturn    = "request_return_refund"
	actionCodeRequestExchange  = "request_exchange"
	actionCodeWaitForUpdate    = "wait_for_update"
	actionCodeEscalateToHuman  = "escalate_to_human"

	defaultSlotRetryThreshold = 2
)

var orderNumberPattern = regexp.MustCompile(`(?i)\b[A-Z]{2,}[A-Z0-9_-]{6,}\b`)

type UserMessage struct {
	ContentText     string   `json:"content_text,omitempty"`
	MessageTypeCode string   `json:"message_type_code,omitempty"`
	AssetIDs        []uint64 `json:"asset_ids,omitempty"`
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

type AfterSaleDecisionCard struct {
	DecisionPathCode string `json:"decision_path_code,omitempty"`
	ReasonText       string `json:"reason_text,omitempty"`
	ConstraintText   string `json:"constraint_text,omitempty"`
	NextStepText     string `json:"next_step_text,omitempty"`
}

type ReplyDataCard struct {
	OrderSnapshotCard     *OrderSnapshotCard     `json:"order_snapshot_card,omitempty"`
	AfterSaleDecisionCard *AfterSaleDecisionCard `json:"after_sale_decision_card,omitempty"`
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
	ActiveTaskCode        string            `json:"active_task_code,omitempty"`
	SlotValues            map[string]string `json:"slot_values,omitempty"`
	MissingSlots          []MissingSlot     `json:"missing_slots,omitempty"`
	SlotRetryCount        uint32            `json:"slot_retry_count,omitempty"`
	LastSlotPromptCode    string            `json:"last_slot_prompt_code,omitempty"`
	SlotFillingFailed     bool              `json:"slot_filling_failed,omitempty"`
	EscalationReasonCode  string            `json:"escalation_reason_code,omitempty"`
	AnchoredOrderNo       string            `json:"anchored_order_no,omitempty"`
	AnchoredSubOrderNo    string            `json:"anchored_sub_order_no,omitempty"`
	LatestFactsSummary    string            `json:"latest_facts_summary,omitempty"`
	LatestDecisionSummary string            `json:"latest_decision_summary,omitempty"`
	HandoffRecommended    bool              `json:"handoff_recommended,omitempty"`
	UnresolvedTurnCount   uint32            `json:"unresolved_turn_count,omitempty"`
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
	Message          UserMessage
	Conversation     ConversationAnchors
	PreviousSession  TaskSessionState
	Security         SecurityContext
	OrderRepository  OrderSnapshotRepository
}

type AfterSaleTurnOutput struct {
	Session TaskSessionState
	Reply   ReplyPayload
}

func detectGuardIntent(msg UserMessage) GuardDecision {
	// 先清洗文本输入，避免空白字符影响 Guard 分类。
	text := strings.TrimSpace(msg.ContentText)
	// 再统一成小写文本，便于后续做稳定的关键词判断。
	lowerText := strings.ToLower(text)
	// 如果只有附件没有有效文本，就直接命中附件补充引导。
	if text == "" && len(msg.AssetIDs) > 0 {
		return GuardDecision{IntentCode: intentCodeAttachmentOnly, GuardResultCode: guardResultAttachmentOnly}
	}
	// 如果文本是典型寒暄，就优先返回欢迎引导而不是进入售后主链路。
	if containsAny(lowerText, "在吗", "你好", "您好", "有人吗", "hello", "hi") {
		return GuardDecision{IntentCode: intentCodeGreeting, GuardResultCode: guardResultGreeting}
	}
	// 如果文本命中售后关键词，就进入查询编排型售后任务。
	if containsAny(lowerText, "订单", "物流", "发货", "退款", "退货", "换货", "催发货", "售后", "包裹", "快递", "人工") {
		return GuardDecision{IntentCode: intentCodeAfterSaleTask, GuardResultCode: guardResultAfterSaleTask}
	}
	// 如果文本没有任何售后线索，就按超纲问题安全兜底。
	return GuardDecision{IntentCode: intentCodeOutOfScope, GuardResultCode: guardResultOutOfScope}
}

func planAfterSaleTurn(ctx context.Context, in AfterSaleTurnInput) (AfterSaleTurnOutput, error) {
	// 先复制上一轮任务会话，保证多轮对话状态可以被本轮安全复用。
	session := cloneTaskSession(in.PreviousSession)
	// 再对本轮用户消息做 Guard 判断，确保非售后输入不被硬套进主链路。
	guard := detectGuardIntent(in.Message)
	// 如果命中前置 Guard，就直接返回兜底回复，不再查询订单事实。
	if guard.IntentCode != intentCodeAfterSaleTask {
		return AfterSaleTurnOutput{
			Session: session,
			Reply:   buildGuardReply(guard),
		}, nil
	}
	// 识别当前售后任务类型，保证规则决策只在有限任务集合内展开。
	activeTaskCode := detectTaskCode(in.Message.ContentText, session.ActiveTaskCode)
	// 把本轮识别出的任务写回会话，供后续多轮追问和复用。
	session.ActiveTaskCode = activeTaskCode
	// 继续把输入中的订单号、子单号和问题类型合并进任务会话。
	session = mergeSlotValues(session, in.Message.ContentText, in.Conversation, activeTaskCode)
	// 如果用户本轮明确要人工，就直接走正式转人工出口。
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
	// 如果当前仍然缺少关键订单定位信息，就先走槽位追问或熔断。
	if strings.TrimSpace(session.SlotValues[slotCodeOrderNo]) == "" {
		return handleMissingOrderSlot(session, guard)
	}
	// 通过安全身份和业务参数做强归属订单查询，不把模型提取结果直接下推到数据库。
	lookupResult, err := secureQueryOrderSnapshot(ctx, in.OrderRepository, in.Security, OrderQuery{
		OrderNo:     session.SlotValues[slotCodeOrderNo],
		SubOrderNo:  session.SlotValues[slotCodeSubOrderNo],
		ProblemType: session.SlotValues[slotCodeProblem],
	})
	// 如果安全查询流程本身异常，就继续向上返回统一错误。
	if err != nil {
		return AfterSaleTurnOutput{}, err
	}
	// 如果工具层返回统一错误码，就按安全语义回复而不是继续规则推断。
	if lookupResult.ErrorCode != "" {
		return handleOrderLookupFailure(session, guard, lookupResult.ErrorCode), nil
	}
	// 用实时订单事实做结构化规则判断，给出当前最适合的售后路径。
	decisionCard, actions, replyText := decideAfterSalePath(activeTaskCode, lookupResult.Snapshot)
	// 把最新事实和决策摘要写入会话，供下一轮“那我现在能退款吗”复用。
	session.LatestFactsSummary = summarizeOrderSnapshot(lookupResult.Snapshot)
	// 把最新规则结论沉淀到会话状态中，便于 handoff 和多轮解释复用。
	session.LatestDecisionSummary = decisionCard.ReasonText
	// 成功查到订单后，清空缺槽和重试计数，避免后续错误熔断。
	session.MissingSlots = nil
	// 成功定位订单后，重置追问计数，避免一次成功后仍被视为失败状态。
	session.SlotRetryCount = 0
	// 成功定位订单后，清空上一轮追问码，避免后续重复 prompt 判断失真。
	session.LastSlotPromptCode = ""
	// 成功完成规则决策后，明确当前会话不需要立刻转人工。
	session.HandoffRecommended = false
	// 当前轮顺利完成查询和建议后，清空熔断标记。
	session.SlotFillingFailed = false
	// 当前轮结论稳定后，清空升级原因码。
	session.EscalationReasonCode = ""
	// 当前问题已被系统处理一轮，就把未解决轮次重置掉。
	session.UnresolvedTurnCount = 0
	// 组装固定格式的结构化回复载荷，供前端直接渲染文本、卡片和动作。
	reply := ReplyPayload{
		ReplyText:       replyText,
		IntentCode:      intentCodeAfterSaleTask,
		GuardResultCode: guard.GuardResultCode,
		Confidence:      0.92,
		DataCards: []ReplyDataCard{
			{OrderSnapshotCard: toOrderSnapshotCard(lookupResult.Snapshot)},
			{AfterSaleDecisionCard: decisionCard},
		},
		SuggestedActions: actions,
		SlotRetryCount:   session.SlotRetryCount,
	}
	// 把当前轮生成的结构化结果和任务会话一起返回给上层流程。
	return AfterSaleTurnOutput{Session: session, Reply: reply}, nil
}

func handleMissingOrderSlot(session TaskSessionState, guard GuardDecision) (AfterSaleTurnOutput, error) {
	// 如果上一轮已经在追问订单号，本轮仍没补齐，就递增追问次数。
	if session.LastSlotPromptCode == slotPromptCodeOrderNo {
		session.SlotRetryCount++
	} else {
		// 首次进入关键槽位追问时，把重试次数重置为第一轮。
		session.SlotRetryCount = 1
	}
	// 把当前缺失的关键槽位写回会话，驱动前端和下一轮追问。
	session.MissingSlots = []MissingSlot{{SlotCode: slotCodeOrderNo, PromptText: "请提供订单号，或告诉我是哪个订单。", Required: true}}
	// 记录本轮追问码，避免后续重复文案无限循环。
	session.LastSlotPromptCode = slotPromptCodeOrderNo
	// 缺少关键槽位时先计入未解决轮次，为连续未解决熔断做准备。
	session.UnresolvedTurnCount++
	// 如果追问次数已经达到阈值，就直接触发熔断和转人工建议。
	if session.SlotRetryCount >= defaultSlotRetryThreshold {
		// 到达阈值后标记槽位补齐失败，阻止系统继续复读。
		session.SlotFillingFailed = true
		// 熔断后记录统一升级原因码，供 handoff summary 复用。
		session.EscalationReasonCode = escalationReasonSlotFillingFailed
		// 熔断后把转人工建议写入任务会话。
		session.HandoffRecommended = true
		// 熔断时直接返回正式体验出口，而不是继续追问订单号。
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
	// 未达到阈值时，只做一次明确追问，不进入订单查询。
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

func secureQueryOrderSnapshot(ctx context.Context, repo OrderSnapshotRepository, security SecurityContext, query OrderQuery) (OrderLookupResult, error) {
	// 如果缺少安全身份或订单号，就直接返回统一上下文不足错误。
	if security.UserID == 0 || strings.TrimSpace(query.OrderNo) == "" {
		return OrderLookupResult{ErrorCode: orderLookupErrorInsufficientContext}, nil
	}
	// 如果当前环境还没有挂上订单仓储，就直接按下游不可用做安全降级。
	if repo == nil {
		return OrderLookupResult{ErrorCode: orderLookupErrorDownstream}, nil
	}
	// 组装带安全身份的强过滤条件，确保底层查询永远附带 user_id。
	filter := OrderOwnershipFilter{
		UserID:         security.UserID,
		ShopNo:         strings.TrimSpace(security.ShopNo),
		OrderNo:        strings.TrimSpace(query.OrderNo),
		SubOrderNo:     strings.TrimSpace(query.SubOrderNo),
		RequestID:      strings.TrimSpace(security.RequestID),
		ConversationNo: strings.TrimSpace(security.ConversationNo),
		RunNo:          strings.TrimSpace(security.RunNo),
	}
	// 通过仓储接口执行订单快照查询，把安全过滤真正压到工具层边界。
	snapshot, err := repo.QueryOrderSnapshot(ctx, filter)
	// 如果下游显式返回统一错误码，就把它原样映射回工具层语义。
	if err != nil {
		if lookupErr, ok := err.(*OrderLookupError); ok {
			return OrderLookupResult{ErrorCode: lookupErr.Code}, nil
		}
		// 其他未知错误统一折叠成下游不可用，避免泄露内部实现细节。
		return OrderLookupResult{ErrorCode: orderLookupErrorDownstream}, nil
	}
	// 如果仓储没有返回任何订单，就按安全 not found 处理。
	if snapshot == nil {
		return OrderLookupResult{ErrorCode: orderLookupErrorNotFound}, nil
	}
	// 如果仓储返回的结果未确认归属，就按权限失败兜底拦截。
	if !snapshot.OwnershipConfirmed {
		return OrderLookupResult{ErrorCode: orderLookupErrorPermissionDenied}, nil
	}
	// 查询成功后把订单事实返回给上层规则决策。
	return OrderLookupResult{Snapshot: snapshot}, nil
}

func handleOrderLookupFailure(session TaskSessionState, guard GuardDecision, errorCode string) AfterSaleTurnOutput {
	// 先把当前未解决轮次递增，便于统一处理连续未解决问题。
	session.UnresolvedTurnCount++
	// 如果是安全 not found 或 permission denied，就返回不泄露存在性的安全话术。
	if errorCode == orderLookupErrorNotFound || errorCode == orderLookupErrorPermissionDenied {
		return AfterSaleTurnOutput{
			Session: session,
			Reply: ReplyPayload{
				ReplyText:       "暂未查询到与你当前账号匹配的订单信息，请核对订单号后再试。",
				IntentCode:      intentCodeAfterSaleTask,
				GuardResultCode: guard.GuardResultCode,
				Confidence:      0.84,
			},
		}
	}
	// 如果下游不可用，就建议转人工但不编造订单事实。
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
	// 其他情况统一按上下文不足回复，继续引导补充信息。
	return AfterSaleTurnOutput{
		Session: session,
		Reply: ReplyPayload{
			ReplyText:       "我还缺少足够的信息来准确判断，请补充订单号或更具体的问题描述。",
			IntentCode:      intentCodeAfterSaleTask,
			GuardResultCode: guard.GuardResultCode,
			Confidence:      0.78,
		},
	}
}

func buildGuardReply(guard GuardDecision) ReplyPayload {
	// Greeting 直接返回欢迎语和业务引导。
	if guard.IntentCode == intentCodeGreeting {
		return ReplyPayload{
			ReplyText:       "你好，我是订单售后助手，可以帮你查订单状态、物流进度和退款退货建议。",
			IntentCode:      guard.IntentCode,
			GuardResultCode: guard.GuardResultCode,
			Confidence:      0.99,
		}
	}
	// AttachmentOnly 只引导补充有效文本或订单信息。
	if guard.IntentCode == intentCodeAttachmentOnly {
		return ReplyPayload{
			ReplyText:       "我可以帮你处理订单和售后问题，请再补充订单号或问题描述。",
			IntentCode:      guard.IntentCode,
			GuardResultCode: guard.GuardResultCode,
			Confidence:      0.95,
		}
	}
	// 其余 Guard 都按超纲问题安全拒答。
	return ReplyPayload{
		ReplyText:       "抱歉，我目前主要负责订单和售后问题。如果你要查询订单、物流或退款退货，我可以继续帮你处理。",
		IntentCode:      guard.IntentCode,
		GuardResultCode: guard.GuardResultCode,
		Confidence:      0.93,
	}
}

func detectTaskCode(query, fallback string) string {
	// 统一清洗和降级文本，避免大小写和空白影响任务分类。
	lowerText := strings.ToLower(strings.TrimSpace(query))
	// 退款、退货、换货优先进入售后路径判断任务。
	if containsAny(lowerText, "退款", "退货", "换货") {
		return taskCodeRefundDecision
	}
	// 发货、物流、快递优先进入物流查询任务。
	if containsAny(lowerText, "物流", "发货", "快递", "催发货") {
		return taskCodeLogisticsQuery
	}
	// 问订单进度时优先进入订单状态查询任务。
	if containsAny(lowerText, "订单状态", "订单进度", "什么时候到") {
		return taskCodeOrderStatusQuery
	}
	// 如果当前轮没有新线索，就尽量复用上一轮任务。
	if strings.TrimSpace(fallback) != "" {
		return strings.TrimSpace(fallback)
	}
	// 默认收敛到动作解释型售后任务。
	return taskCodeActionExplanation
}

func mergeSlotValues(session TaskSessionState, query string, anchors ConversationAnchors, activeTaskCode string) TaskSessionState {
	// 先确保会话里的槽位 map 一定可写，避免 nil map 写入异常。
	if session.SlotValues == nil {
		session.SlotValues = make(map[string]string)
	}
	// 如果会话里还没有锚定订单，就优先复用会话创建时的订单锚点。
	if strings.TrimSpace(session.SlotValues[slotCodeOrderNo]) == "" && strings.TrimSpace(anchors.OrderNo) != "" {
		session.SlotValues[slotCodeOrderNo] = strings.TrimSpace(anchors.OrderNo)
	}
	// 如果会话里还没有锚定子订单，就复用会话创建时的子订单锚点。
	if strings.TrimSpace(session.SlotValues[slotCodeSubOrderNo]) == "" && strings.TrimSpace(anchors.SubOrderNo) != "" {
		session.SlotValues[slotCodeSubOrderNo] = strings.TrimSpace(anchors.SubOrderNo)
	}
	// 从自然语言里提取订单号，用本轮用户补充覆盖旧值。
	if orderNo := extractOrderNumber(query, "ord"); orderNo != "" {
		session.SlotValues[slotCodeOrderNo] = orderNo
		session.AnchoredOrderNo = orderNo
	}
	// 从自然语言里提取子订单号，便于更细粒度售后定位。
	if subOrderNo := extractOrderNumber(query, "sub"); subOrderNo != "" {
		session.SlotValues[slotCodeSubOrderNo] = subOrderNo
		session.AnchoredSubOrderNo = subOrderNo
	}
	// 把本轮识别到的问题类型写入槽位，供规则库和 handoff 复用。
	session.SlotValues[slotCodeProblem] = activeTaskCode
	// 返回合并后的任务会话。
	return session
}

func extractOrderNumber(query, prefix string) string {
	// 把文本里的候选编号全部抽出来，尽量适配自然语言中的订单号。
	matches := orderNumberPattern.FindAllString(strings.ToUpper(strings.TrimSpace(query)), -1)
	// 逐个检查候选编号，优先挑出符合指定前缀的编号。
	for _, item := range matches {
		if strings.HasPrefix(item, strings.ToUpper(prefix)) {
			return item
		}
	}
	// 如果没指定前缀命中 but 只有一个候选编号，也接受它作为订单号。
	if len(matches) == 1 {
		return matches[0]
	}
	// 没找到合适编号时返回空串，交给缺槽流程处理。
	return ""
}

func decideAfterSalePath(taskCode string, snapshot *OrderSnapshot) (*AfterSaleDecisionCard, []SuggestedAction, string) {
	// 如果当前已存在售后单，就优先建议等待现有流程，避免重复申请。
	if strings.EqualFold(strings.TrimSpace(snapshot.AfterSaleStatus), "PROCESSING") {
		return &AfterSaleDecisionCard{
				DecisionPathCode: decisionPathWaitExistingAfterSale,
				ReasonText:       "当前订单已经有售后单在处理中，不建议重复发起。",
				ConstraintText:   "重复提交可能导致客服判断分散。",
				NextStepText:     "建议等待当前售后单处理进展，必要时转人工跟进。",
			},
			[]SuggestedAction{{
				ActionCode: actionCodeWaitForUpdate,
				Label:      "等待进展",
				Enabled:    true,
			}, {
				ActionCode: actionCodeEscalateToHuman,
				Label:      "转人工",
				Enabled:    true,
			}},
			"当前订单已经有售后处理记录，我建议先等待现有售后结果，避免重复提交。"
	}
	// 未发货场景优先建议直接退款，这是电商售后最稳定的决策路径。
	if strings.EqualFold(strings.TrimSpace(snapshot.FulfillmentStatus), "UNSHIPPED") || strings.EqualFold(strings.TrimSpace(snapshot.LogisticsStatus), "NOT_SHIPPED") {
		return &AfterSaleDecisionCard{
				DecisionPathCode: decisionPathRefundOnly,
				ReasonText:       "订单尚未发货，当前更适合直接申请退款。",
				ConstraintText:   "未发货阶段通常不需要先走退货流程。",
				NextStepText:     "优先发起退款；如果商家长时间未处理，可再考虑转人工。",
			},
			[]SuggestedAction{{
				ActionCode: actionCodeRequestRefund,
				Label:      "申请退款",
				Enabled:    true,
			}, {
				ActionCode: actionCodeEscalateToHuman,
				Label:      "转人工",
				Enabled:    true,
			}},
			"根据当前订单事实，这个订单还没有发货，更适合直接申请退款。"
	}
	// 已发货或已签收场景，优先建议退货退款或换货。
	return &AfterSaleDecisionCard{
			DecisionPathCode: decisionPathReturnRefundOrExchange,
			ReasonText:       "订单已发货或已进入履约阶段，更适合走退货退款或换货。",
			ConstraintText:   "如果物流仍在途中，部分动作可能需要等待签收或拒收节点。",
			NextStepText:     "先确认商品状态和物流节点，再选择退货退款或换货。",
		},
		[]SuggestedAction{{
			ActionCode: actionCodeRequestReturn,
			Label:      "退货退款",
			Enabled:    taskCode == taskCodeRefundDecision,
			ReasonIfDisabled: "当前问题更偏向物流或订单进度，建议先补充场景再决定动作。",
		}, {
			ActionCode: actionCodeRequestExchange,
			Label:      "申请换货",
			Enabled:    true,
		}},
		"这个订单已经进入履约阶段，我更建议你根据商品状态选择退货退款或换货。"
}

func summarizeOrderSnapshot(snapshot *OrderSnapshot) string {
	// 如果没有订单快照，就返回空摘要，避免生成伪事实。
	if snapshot == nil {
		return ""
	}
	// 把关键事实压缩成短摘要，便于后续多轮复用和 handoff。
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
	// 如果没有快照，就不返回卡片，避免前端渲染空结构。
	if snapshot == nil {
		return nil
	}
	// 把内部订单快照转成前端直接可用的结构化卡片。
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

func cloneTaskSession(in TaskSessionState) TaskSessionState {
	// 先浅拷贝基础字段，保留上一轮任务语义。
	out := in
	// 再深拷贝槽位 map，避免不同轮次共享同一底层 map。
	if in.SlotValues != nil {
		out.SlotValues = make(map[string]string, len(in.SlotValues))
		for key, value := range in.SlotValues {
			out.SlotValues[key] = value
		}
	}
	// 也复制缺槽切片，避免后续写入污染上一轮状态。
	if len(in.MissingSlots) > 0 {
		out.MissingSlots = append([]MissingSlot(nil), in.MissingSlots...)
	}
	// 返回可安全写入的新会话对象。
	return out
}
