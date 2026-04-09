package runtime

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/v2/errors/gerror"
)

const (
	nodeLoadConversationContext = "LoadConversationContext"
	nodeLoadRecentMessages      = "LoadRecentMessages"
	nodeLoadSummary             = "LoadSummary"
	nodeIntentRouter            = "IntentRouter"
	nodeKnowledgeRetriever      = "KnowledgeRetriever"
	nodeToolRouter              = "ToolRouter"
	nodeCallMcpTool             = "CallMcpTool"
	nodeGenerateAnswer          = "GenerateAnswer"
	nodePersistCheckpoint       = "PersistCheckpoint"

	runStatusPending         = "PENDING"
	runStatusProcessing      = "PROCESSING"
	runStatusWaitingToolCall = "WAITING_TOOL_CALL"
	runStatusGenerating      = "GENERATING"
	runStatusSuccess         = "SUCCESS"
	runStatusFailed          = "FAILED"
	runStatusAborted         = "ABORTED"
	runStatusEscalated       = "ESCALATED"

	toolResultUnspecified = "UNSPECIFIED"
	toolResultSuccess     = "SUCCESS"
	toolResultDegraded    = "DEGRADED"
	toolResultError       = "ERROR"
)

type CheckpointHook func(ctx context.Context, cp GraphCheckpointSnapshot) error

type RunInput struct {
	SecurityPrompt string
	SceneCode      string
	Summary        string
	Messages       []*schema.Message
	UserQuery      string
	RequestID      string
	Message        UserMessage
	PreviousSession TaskSessionState
	Conversation   ConversationAnchors
	Security       SecurityContext
	CheckpointHook CheckpointHook
}

type RunOutput struct {
	AnswerText          string
	RunStatusCode       string
	RiskDecisionCode    string
	PromptInjection     bool
	AnswerSources       []AnswerSource
	PromptTokenCount    uint32
	CompletionTokens    uint32
	CurrentNodeCode     string
	ToolResultStatus    string
	QueueBlocked        bool
	QueueHintMessage    string
	DegradedReasonCode  string
	GraphStateJSON      string
	ToolWaitTimeoutAt   *time.Time
	SelectedToolName    string
	SelectedAdapterCode string
	SelectedToolScope   string
	ReplyPayload        ReplyPayload
	TaskSession         TaskSessionState
}

type AnswerSource struct {
	SourceTypeCode string `json:"source_type_code"`
	SourceID       string `json:"source_id"`
	SourceVersion  uint64 `json:"source_version"`
	Title          string `json:"title"`
	Snippet        string `json:"snippet"`
}

type GraphCheckpointSnapshot struct {
	CurrentNodeCode    string
	CheckpointStatus   string
	GraphStateJSON     string
	ToolWaitTimeoutAt  *time.Time
	ToolResultStatus   string
	DegradedReasonCode string
	QueueBlocked       bool
	QueueHintMessage   string
}

type Runner interface {
	Run(ctx context.Context, in RunInput) (*RunOutput, error)
	Summarize(ctx context.Context, summary string, recent []*schema.Message) (string, error)
	HandoffSummary(ctx context.Context, summary string, recent []*schema.Message, reasonCode string) (string, error)
}

type RunnerOptions struct {
	OrderRepository OrderSnapshotRepository
}

type einoRunner struct {
	runFlow         compose.Runnable[*runState, *runState]
	adapters        *AdapterRegistry
	orderRepository OrderSnapshotRepository
}

type runState struct {
	Input             RunInput        `json:"-"`
	SecurityPrompt    string          `json:"security_prompt"`
	SceneCode         string          `json:"scene_code"`
	Summary           string          `json:"summary"`
	UserQuery         string          `json:"user_query"`
	NormalizedQuery   string          `json:"normalized_query"`
	Message           UserMessage     `json:"message"`
	TaskSession       TaskSessionState `json:"task_session"`
	ReplyPayload      ReplyPayload    `json:"reply_payload"`
	GuardResultCode   string          `json:"guard_result_code,omitempty"`
	Conversation      ConversationAnchors `json:"conversation"`
	Security          SecurityContext `json:"security"`
	RecentMessages    []messageDigest `json:"recent_messages,omitempty"`
	IntentCode        string          `json:"intent_code"`
	NeedKnowledge     bool            `json:"need_knowledge"`
	NeedTool          bool            `json:"need_tool"`
	SelectedToolName  string          `json:"selected_tool_name,omitempty"`
	SelectedAdapter   string          `json:"selected_adapter_code,omitempty"`
	SelectedScopeCode string          `json:"selected_tool_scope_code,omitempty"`
	ToolTimeoutAt     *time.Time      `json:"tool_wait_timeout_at,omitempty"`
	ToolResult        toolResult      `json:"tool_result"`
	CurrentNodeCode   string          `json:"current_node_code"`
	RunStatusCode     string          `json:"run_status_code"`
	RiskDecisionCode  string          `json:"risk_decision_code"`
	PromptInjection   bool            `json:"prompt_injection"`
	QueueBlocked      bool            `json:"queue_blocked"`
	QueueHintMessage  string          `json:"queue_hint_message,omitempty"`
	DegradedReason    string          `json:"degraded_reason_code,omitempty"`
	AnswerText        string          `json:"answer_text,omitempty"`
	AnswerSources     []AnswerSource  `json:"answer_sources,omitempty"`
}

type messageDigest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type toolResult struct {
	Status             string `json:"status"`
	Message            string `json:"message,omitempty"`
	DataJSON           string `json:"data_json,omitempty"`
	DegradedReasonCode string `json:"degraded_reason_code,omitempty"`
	AdapterCode        string `json:"adapter_code,omitempty"`
	ToolName           string `json:"tool_name,omitempty"`
}

func NewRunner() Runner {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	r := &einoRunner{
		adapters: newDefaultAdapterRegistry(),
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if err := r.init(); err != nil {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		panic(err)
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
return r
}

func NewRunnerWithOptions(opts RunnerOptions) Runner {
	// 初始化带依赖注入的 Runner，让测试和生产都能挂上不同的订单查询实现。
	r := &einoRunner{
		adapters:        newDefaultAdapterRegistry(),
		orderRepository: opts.OrderRepository,
	}
	// 如果运行图初始化失败，就在启动阶段直接暴露问题而不是带病运行。
	if err := r.init(); err != nil {
		panic(err)
	}
	// 返回完成初始化的 Runner 实例，供上层业务流程复用。
	return r
}

func (r *einoRunner) init() error {
	// 构建并编译 Eino 执行链，确保后续 Run 以固定节点顺序推进。
	chain := compose.NewChain[*runState, *runState]()
	// 把当前处理节点接入 Eino 图中，明确每一步执行顺序和职责边界。
	chain.AppendLambda(compose.InvokableLambda(r.loadConversationContext), compose.WithNodeKey(nodeLoadConversationContext))
	// 把当前处理节点接入 Eino 图中，明确每一步执行顺序和职责边界。
	chain.AppendLambda(compose.InvokableLambda(r.loadRecentMessages), compose.WithNodeKey(nodeLoadRecentMessages))
	// 把当前处理节点接入 Eino 图中，明确每一步执行顺序和职责边界。
	chain.AppendLambda(compose.InvokableLambda(r.loadSummary), compose.WithNodeKey(nodeLoadSummary))
	// 把当前处理节点接入 Eino 图中，明确每一步执行顺序和职责边界。
	chain.AppendLambda(compose.InvokableLambda(r.assistantIntentRouter), compose.WithNodeKey(nodeIntentRouter))
	// 把当前处理节点接入 Eino 图中，明确每一步执行顺序和职责边界。
	chain.AppendLambda(compose.InvokableLambda(r.knowledgeRetriever), compose.WithNodeKey(nodeKnowledgeRetriever))
	// 把当前处理节点接入 Eino 图中，明确每一步执行顺序和职责边界。
	chain.AppendLambda(compose.InvokableLambda(r.toolRouter), compose.WithNodeKey(nodeToolRouter))
	// 把当前处理节点接入 Eino 图中，明确每一步执行顺序和职责边界。
	chain.AppendLambda(compose.InvokableLambda(r.callMcpTool), compose.WithNodeKey(nodeCallMcpTool))
	// 把当前处理节点接入 Eino 图中，明确每一步执行顺序和职责边界。
	chain.AppendLambda(compose.InvokableLambda(r.assistantGenerateAnswer), compose.WithNodeKey(nodeGenerateAnswer))
	// 把当前处理节点接入 Eino 图中，明确每一步执行顺序和职责边界。
	chain.AppendLambda(compose.InvokableLambda(r.persistCheckpoint), compose.WithNodeKey(nodePersistCheckpoint))

	// 构建并编译 Eino 执行链，确保后续 Run 以固定节点顺序推进。
	runFlow, err := chain.Compile(context.Background())
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return gerror.Wrap(err, "compile eino run flow failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	r.runFlow = runFlow
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return nil
}

func (r *einoRunner) Run(ctx context.Context, in RunInput) (*RunOutput, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if r.runFlow == nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.New("eino run flow is not initialized")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	state := &runState{
		Input:            in,
		SecurityPrompt:   strings.TrimSpace(in.SecurityPrompt),
		SceneCode:        strings.TrimSpace(in.SceneCode),
		Summary:          strings.TrimSpace(in.Summary),
		UserQuery:        strings.TrimSpace(in.UserQuery),
		NormalizedQuery:  strings.ToLower(strings.TrimSpace(in.UserQuery)),
		Message:          in.Message,
		TaskSession:      cloneTaskSession(in.PreviousSession),
		Conversation:     in.Conversation,
		Security:         in.Security,
		RunStatusCode:    runStatusPending,
		RiskDecisionCode: "PASS",
		ToolResult: toolResult{
			Status: toolResultUnspecified,
		},
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if strings.TrimSpace(state.Message.ContentText) == "" {
		state.Message.ContentText = state.UserQuery
	}
	outState, err := r.runFlow.Invoke(ctx, state)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "invoke eino run flow failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	graphStateJSON, _ := json.Marshal(outState)
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	output := &RunOutput{
		AnswerText:          outState.AnswerText,
		RunStatusCode:       normalizeFinalRunStatus(outState.RunStatusCode),
		RiskDecisionCode:    outState.RiskDecisionCode,
		PromptInjection:     outState.PromptInjection,
		AnswerSources:       outState.AnswerSources,
		PromptTokenCount:    estimateTokenCount(outState.SecurityPrompt, outState.UserQuery, outState.Summary),
		CompletionTokens:    estimateTokenCount(outState.AnswerText),
		CurrentNodeCode:     outState.CurrentNodeCode,
		ToolResultStatus:    outState.ToolResult.Status,
		QueueBlocked:        outState.QueueBlocked,
		QueueHintMessage:    outState.QueueHintMessage,
		DegradedReasonCode:  coalesce(outState.DegradedReason, outState.ToolResult.DegradedReasonCode),
		GraphStateJSON:      string(graphStateJSON),
		ToolWaitTimeoutAt:   outState.ToolTimeoutAt,
		SelectedToolName:    outState.SelectedToolName,
		SelectedAdapterCode: outState.SelectedAdapter,
		SelectedToolScope:   outState.SelectedScopeCode,
		ReplyPayload:        outState.ReplyPayload,
		TaskSession:         outState.TaskSession,
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return output, nil
}

func (r *einoRunner) Summarize(ctx context.Context, summary string, recent []*schema.Message) (string, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	_ = ctx
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if strings.TrimSpace(summary) != "" {
		// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
		return strings.TrimSpace(summary), nil
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	parts := make([]string, 0, min(len(recent), 6))
	// 遍历当前集合或循环条件，逐项推进本段业务处理并累计最终结果。
	for i := max(0, len(recent)-6); i < len(recent); i++ {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		if recent[i] == nil {
			// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
			continue
		}
		// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
		text := strings.TrimSpace(recent[i].Content)
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		if text != "" {
			// 把当前元素追加到结果集合中，逐步汇总出最终返回或落库的数据集。
			parts = append(parts, text)
		}
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return strings.Join(parts, " | "), nil
}

func (r *einoRunner) HandoffSummary(ctx context.Context, summary string, recent []*schema.Message, reasonCode string) (string, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	mergedSummary, err := r.Summarize(ctx, summary, recent)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return "", err
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if mergedSummary == "" {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		mergedSummary = "用户请求转人工，需要人工客服继续接手。"
	}
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	reasonCode = strings.TrimSpace(reasonCode)
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if reasonCode == "" {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		reasonCode = "USER_REQUEST"
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return "转人工原因: " + reasonCode + "；会话摘要: " + mergedSummary, nil
}

func (r *einoRunner) loadConversationContext(ctx context.Context, state *runState) (*runState, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	state.RunStatusCode = runStatusProcessing
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	state.CurrentNodeCode = nodeLoadConversationContext
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return state, state.emitCheckpoint(ctx)
}

func (r *einoRunner) loadRecentMessages(ctx context.Context, state *runState) (*runState, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	state.CurrentNodeCode = nodeLoadRecentMessages
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	state.RecentMessages = make([]messageDigest, 0, len(state.Input.Messages))
	// 遍历当前集合或循环条件，逐项推进本段业务处理并累计最终结果。
	for _, msg := range state.Input.Messages {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		if msg == nil {
			// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
			continue
		}
		// 把当前元素追加到结果集合中，逐步汇总出最终返回或落库的数据集。
		state.RecentMessages = append(state.RecentMessages, messageDigest{
			Role:    string(msg.Role),
			Content: strings.TrimSpace(msg.Content),
		})
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return state, state.emitCheckpoint(ctx)
}

func (r *einoRunner) loadSummary(ctx context.Context, state *runState) (*runState, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	state.CurrentNodeCode = nodeLoadSummary
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return state, state.emitCheckpoint(ctx)
}

func (r *einoRunner) intentRouter(ctx context.Context, state *runState) (*runState, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	state.CurrentNodeCode = nodeIntentRouter

	// 根据当前关键状态或意图进入不同分支，确保每个业务场景按对应规则处理。
	switch {
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case containsAny(state.NormalizedQuery, "ignore previous", "system prompt", "忽略之前", "免费拿", "卡 bug", "赔付保证"):
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.IntentCode = "SAFETY_REFUSAL"
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.PromptInjection = true
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.RiskDecisionCode = "REFUSE"
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case containsAny(state.NormalizedQuery, "人工", "转人工", "客服接入"):
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.IntentCode = "HANDOFF_REQUEST"
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case containsAny(state.NormalizedQuery, "订单", "物流", "发货", "退款", "退货", "换货", "支付"):
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.IntentCode = "ORDER_SERVICE"
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.NeedTool = true
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case containsAny(state.NormalizedQuery, "积分", "账号", "会员", "收货地址"):
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.IntentCode = "ACCOUNT_SERVICE"
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.NeedTool = true
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case containsAny(state.NormalizedQuery, "规则", "平台", "政策", "价保"):
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.IntentCode = "POLICY_QA"
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.NeedTool = true
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case containsAny(state.NormalizedQuery, "商品", "手机", "鞋子", "尺码", "颜色", "搜索"):
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.IntentCode = "CATALOG_GUIDE"
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.NeedKnowledge = true
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	default:
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.IntentCode = "GENERAL_QA"
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.NeedKnowledge = true
	}

	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return state, state.emitCheckpoint(ctx)
}

func (r *einoRunner) knowledgeRetriever(ctx context.Context, state *runState) (*runState, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	state.CurrentNodeCode = nodeKnowledgeRetriever
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if !state.NeedKnowledge {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return state, state.emitCheckpoint(ctx)
	}

	// 根据当前关键状态或意图进入不同分支，确保每个业务场景按对应规则处理。
	switch state.IntentCode {
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case "CATALOG_GUIDE":
		// 把当前元素追加到结果集合中，逐步汇总出最终返回或落库的数据集。
		state.AnswerSources = append(state.AnswerSources, AnswerSource{
			SourceTypeCode: "CATALOG_GUIDE",
			SourceID:       "shopping-guide",
			SourceVersion:  1,
			Title:          "选品导购建议",
			Snippet:        "可以结合预算、品牌、颜色和使用场景继续缩小范围。",
		})
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case "GENERAL_QA":
		// 把当前元素追加到结果集合中，逐步汇总出最终返回或落库的数据集。
		state.AnswerSources = append(state.AnswerSources, AnswerSource{
			SourceTypeCode: "HELP_CENTER",
			SourceID:       "assistant-intro",
			SourceVersion:  1,
			Title:          "智能客服说明",
			Snippet:        "支持商品咨询、规则答疑、订单进度和售后流程说明。",
		})
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return state, state.emitCheckpoint(ctx)
}

func (r *einoRunner) toolRouter(ctx context.Context, state *runState) (*runState, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	state.CurrentNodeCode = nodeToolRouter
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if !state.NeedTool {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return state, state.emitCheckpoint(ctx)
	}

	// 根据当前关键状态或意图进入不同分支，确保每个业务场景按对应规则处理。
	switch state.IntentCode {
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case "ORDER_SERVICE":
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.SelectedAdapter = "assistant-order-mcp"
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.SelectedToolName = "query_order_context"
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.SelectedScopeCode = "USER_PRIVATE_READ"
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case "ACCOUNT_SERVICE":
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.SelectedAdapter = "assistant-account-mcp"
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.SelectedToolName = "query_account_context"
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.SelectedScopeCode = "USER_PRIVATE_READ"
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case "POLICY_QA":
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.SelectedAdapter = "assistant-policy-mcp"
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.SelectedToolName = "query_policy_knowledge"
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.SelectedScopeCode = "PUBLIC_READ"
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	default:
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.SelectedAdapter = "assistant-catalog-mcp"
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.SelectedToolName = "query_catalog_context"
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.SelectedScopeCode = "PUBLIC_READ"
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return state, state.emitCheckpoint(ctx)
}

func (r *einoRunner) callMcpTool(ctx context.Context, state *runState) (*runState, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	state.CurrentNodeCode = nodeCallMcpTool
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if !state.NeedTool || state.SelectedToolName == "" {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return state, state.emitCheckpoint(ctx)
	}

	// 记录当前业务时间，确保状态流转、排序展示和审计字段使用同一时间基线。
	timeoutAt := time.Now().Add(3 * time.Second)
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	state.RunStatusCode = runStatusWaitingToolCall
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	state.ToolTimeoutAt = &timeoutAt
	// 在当前节点输出检查点快照，让运行状态在关键节点前后都能持久化。
	if err := state.emitCheckpoint(ctx); err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}

	// 整理当前工具调用结果，统一成功、降级和错误语义供后续生成节点消费。
	result := toolResult{
		ToolName:    state.SelectedToolName,
		AdapterCode: state.SelectedAdapter,
		Status:      toolResultSuccess,
	}
	// 通过统一 Adapter 边界调用业务域工具，避免在图节点中直接耦合具体下游实现。
	adapterResp, ok, err := r.adapters.Call(ctx, AdapterRequest{
		IntentCode:  state.IntentCode,
		AdapterCode: state.SelectedAdapter,
		ToolName:    state.SelectedToolName,
		ToolScope:   state.SelectedScopeCode,
		UserQuery:   state.UserQuery,
		RequestID:   state.Input.RequestID,
	})
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "call mcp adapter failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if ok && adapterResp != nil {
		// 整理当前工具调用结果，统一成功、降级和错误语义供后续生成节点消费。
		result.Status = normalizeToolResultStatus(adapterResp.Status)
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		result.Message = adapterResp.Message
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		result.DataJSON = adapterResp.DataJSON
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		result.DegradedReasonCode = adapterResp.DegradedReasonCode
	}
	// 整理当前工具调用结果，统一成功、降级和错误语义供后续生成节点消费。
	if result.Status == toolResultDegraded {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.DegradedReason = result.DegradedReasonCode
	}
	// 整理当前工具调用结果，统一成功、降级和错误语义供后续生成节点消费。
	state.ToolResult = result
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return state, state.emitCheckpoint(ctx)
}

func (r *einoRunner) generateAnswer(ctx context.Context, state *runState) (*runState, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	state.CurrentNodeCode = nodeGenerateAnswer
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	state.RunStatusCode = runStatusGenerating
	// 在当前节点输出检查点快照，让运行状态在关键节点前后都能持久化。
	if err := state.emitCheckpoint(ctx); err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}

	// 根据当前关键状态或意图进入不同分支，确保每个业务场景按对应规则处理。
	switch {
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case state.PromptInjection:
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.AnswerText = "这类请求我不能协助。我可以继续帮你查询商品信息、平台规则、订单进度，或者为你转人工客服。"
		// 把当前元素追加到结果集合中，逐步汇总出最终返回或落库的数据集。
		state.AnswerSources = append(state.AnswerSources, AnswerSource{
			SourceTypeCode: "PLATFORM_RULE",
			SourceID:       "agent-safety",
			SourceVersion:  1,
			Title:          "智能客服安全规则",
			Snippet:        "禁止价格承诺、赔付承诺、绕过平台规则或协助异常行为。",
		})
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case state.IntentCode == "HANDOFF_REQUEST":
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.AnswerText = "我可以继续先帮你梳理问题，也可以为你转人工客服。如果你已经确认需要人工，请直接点击转人工或发起转人工请求。"
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case state.ToolResult.Status == toolResultDegraded:
		// 整理当前工具调用结果，统一成功、降级和错误语义供后续生成节点消费。
		state.AnswerText = state.ToolResult.Message
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case containsAny(state.NormalizedQuery, "退款", "退货", "换货"):
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.AnswerText = "如果你是在处理售后问题，我可以先帮你判断路径：未发货通常优先申请退款；已签收更适合走退货退款或换货。你也可以把订单号发给我，我继续帮你梳理。"
		// 把当前元素追加到结果集合中，逐步汇总出最终返回或落库的数据集。
		state.AnswerSources = append(state.AnswerSources, AnswerSource{
			SourceTypeCode: "HELP_CENTER",
			SourceID:       "after-sale-policy",
			SourceVersion:  1,
			Title:          "售后帮助中心",
			Snippet:        "未发货优先退款，已签收按退货退款或换货流程处理。",
		})
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case containsAny(state.NormalizedQuery, "发货", "物流", "快递"):
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.AnswerText = "我可以先帮你看发货与物流问题。一般待发货阶段建议先确认订单状态和店铺时效；若已经发货，再看物流轨迹是否停滞。你把订单号发我，我继续帮你排查。"
		// 把当前元素追加到结果集合中，逐步汇总出最终返回或落库的数据集。
		state.AnswerSources = append(state.AnswerSources, AnswerSource{
			SourceTypeCode: "ORDER_CONTEXT",
			SourceID:       "shipping-faq",
			SourceVersion:  1,
			Title:          "发货与物流 FAQ",
			Snippet:        "待发货先看订单状态，已发货重点看物流轨迹与承运商更新。",
		})
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case containsAny(state.NormalizedQuery, "红色", "手机"):
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.AnswerText = "如果你在找红色手机，告诉我预算、品牌偏好和容量需求会更快；如果你有参考图，也可以直接发图，我再继续帮你筛。"
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	default:
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		state.AnswerText = "我可以帮你处理商品咨询、订单进度、售后流程和平台规则问题。你可以直接发商品链接、订单号，或者上传图片让我继续判断。"
	}

	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	state.RunStatusCode = runStatusSuccess
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return state, state.emitCheckpoint(ctx)
}

func (r *einoRunner) persistCheckpoint(ctx context.Context, state *runState) (*runState, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	state.CurrentNodeCode = nodePersistCheckpoint
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return state, state.emitCheckpoint(ctx)
}

func (s *runState) emitCheckpoint(ctx context.Context) error {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if s.Input.CheckpointHook == nil {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return nil
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	graphStateJSON, err := json.Marshal(s)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return gerror.Wrap(err, "marshal graph state failed")
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return s.Input.CheckpointHook(ctx, GraphCheckpointSnapshot{
		CurrentNodeCode:    s.CurrentNodeCode,
		CheckpointStatus:   normalizeCheckpointStatus(s.RunStatusCode),
		GraphStateJSON:     string(graphStateJSON),
		ToolWaitTimeoutAt:  s.ToolTimeoutAt,
		ToolResultStatus:   s.ToolResult.Status,
		DegradedReasonCode: coalesce(s.DegradedReason, s.ToolResult.DegradedReasonCode),
		QueueBlocked:       s.QueueBlocked,
		QueueHintMessage:   s.QueueHintMessage,
	})
}

func normalizeCheckpointStatus(status string) string {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	status = strings.ToUpper(strings.TrimSpace(status))
	// 根据当前关键状态或意图进入不同分支，确保每个业务场景按对应规则处理。
	switch status {
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case runStatusProcessing, runStatusWaitingToolCall, runStatusGenerating, runStatusSuccess, runStatusFailed, runStatusAborted, runStatusEscalated:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return status
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	default:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return runStatusProcessing
	}
}

func normalizeToolResultStatus(status string) string {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	status = strings.ToUpper(strings.TrimSpace(status))
	// 根据当前关键状态或意图进入不同分支，确保每个业务场景按对应规则处理。
	switch status {
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case toolResultSuccess, toolResultDegraded, toolResultError:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return status
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	default:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return toolResultUnspecified
	}
}

func normalizeFinalRunStatus(status string) string {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	status = strings.ToUpper(strings.TrimSpace(status))
	// 根据当前关键状态或意图进入不同分支，确保每个业务场景按对应规则处理。
	switch status {
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case runStatusSuccess, runStatusFailed, runStatusAborted, runStatusEscalated:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return status
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	default:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return runStatusSuccess
	}
}

func estimateTokenCount(parts ...string) uint32 {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	total := 0
	// 遍历当前集合或循环条件，逐项推进本段业务处理并累计最终结果。
	for _, part := range parts {
		// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
		total += len([]rune(strings.TrimSpace(part)))
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if total <= 0 {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return 0
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return uint32(total)
}

func containsAny(text string, words ...string) bool {
	// 遍历当前集合或循环条件，逐项推进本段业务处理并累计最终结果。
	for _, word := range words {
		// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
		if strings.Contains(text, strings.ToLower(strings.TrimSpace(word))) {
			// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
			return true
		}
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return false
}

func coalesce(values ...string) string {
	// 遍历当前集合或循环条件，逐项推进本段业务处理并累计最终结果。
	for _, value := range values {
		// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
		value = strings.TrimSpace(value)
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		if value != "" {
			// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
			return value
		}
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return ""
}

func min(a, b int) int {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if a < b {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return a
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return b
}

func max(a, b int) int {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if a > b {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return a
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return b
}
