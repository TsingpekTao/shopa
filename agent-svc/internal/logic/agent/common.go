package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	agentv1 "github.com/TsingpekTao/shopa/agent-svc/api/v1"
	"github.com/TsingpekTao/shopa/agent-svc/internal/dao"
	agentruntime "github.com/TsingpekTao/shopa/agent-svc/internal/logic/agent/runtime"
	"github.com/TsingpekTao/shopa/agent-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/agent-svc/internal/model/entity"
	"github.com/TsingpekTao/shopa/agent-svc/internal/service"
	"github.com/cloudwego/eino/schema"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	defaultPageSize      = 20
	maxPageSize          = 100
	maxRecentMessages    = 10
	defaultKnowledgeSize = 10

	runtimeStatusPending         = "PENDING"
	runtimeStatusProcessing      = "PROCESSING"
	runtimeStatusWaitingToolCall = "WAITING_TOOL_CALL"
	runtimeStatusGenerating      = "GENERATING"
	runtimeStatusSuccess         = "SUCCESS"
	runtimeStatusFailed          = "FAILED"
	runtimeStatusAborted         = "ABORTED"
	runtimeStatusEscalated       = "ESCALATED"

	runtimeToolResultUnspecified = "UNSPECIFIED"
	runtimeToolResultSuccess     = "SUCCESS"
	runtimeToolResultDegraded    = "DEGRADED"
	runtimeToolResultError       = "ERROR"

	conversationStatusActive             = "ACTIVE"
	conversationStatusEscalationPending  = "ESCALATION_PENDING"
	conversationStatusEscalationAccepted = "ESCALATION_ACCEPTED"
	conversationStatusHumanHandover      = "HUMAN_HANDOVER"
	conversationStatusHumanClosed        = "HUMAN_CLOSED"
	conversationStatusArchived           = "ARCHIVED"
)

type sAgent struct {
	runner agentruntime.Runner
}

func New() *sAgent {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	svc := &sAgent{runner: agentruntime.NewRunnerWithOptions(agentruntime.RunnerOptions{
		OrderRepository: newBuyerOrderSnapshotHTTPRepository(func() string {
			value, err := g.Cfg().Get(context.Background(), "services.orderHttp", "")
			if err != nil || value == nil {
				return ""
			}
			return value.String()
		}(), nil),
	})}
	// 启动后台维护协程，定期清理卡死 Run，避免僵尸执行长期占据会话状态。
	svc.startRuntimeMaintenance()
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return svc
}

func init() {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	service.RegisterAgent(New())
}

func (s *sAgent) processRun(ctx context.Context, conv *entity.AgentConversation, run *entity.AgentRun, userMsg *entity.AgentMessage) (*entity.AgentRun, *entity.AgentMessage, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if conv == nil || run == nil || userMsg == nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, nil, gerror.NewCode(gcode.CodeInvalidParameter, "conversation, run and user message are required")
	}
	// 拉取最近消息窗口，为总结、推理和转人工摘要提供最新上下文。
	recentRows, err := s.listRecentMessages(ctx, conv.ConversationNo, maxRecentMessages)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, nil, err
	}
	// 记录当前业务时间，确保状态流转、排序展示和审计字段使用同一时间基线。
	now := gtime.Now()
	// 把当前图执行节点和状态持久化到 Run 真相源，避免 Pod 抖动后前端长期卡在生成中。
	if err = s.updateRunCheckpoint(ctx, run.RunNo, agentruntime.GraphCheckpointSnapshot{
		CurrentNodeCode:  "LoadConversationContext",
		CheckpointStatus: runtimeStatusProcessing,
		GraphStateJSON:   normalizeJSON(`{"phase":"start"}`),
	}); err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, nil, err
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	// 回溯上一轮 Run 的任务会话，为当前轮多轮售后对话恢复上下文。
	previousRun, err := s.getPreviousRun(ctx, run)
	// 如果上一轮查询失败，就直接终止，避免带着不完整上下文继续推理。
	if err != nil {
		return nil, nil, err
	}
	// 先准备当前轮要传给 runtime 的上一轮任务会话。
	var previousSession agentruntime.TaskSessionState
	// 如果上一轮存在图快照，就从 graph_state_json 中恢复任务会话。
	if previousTaskSession := agentruntime.ParseGraphTaskSession(firstNonEmpty(func() string {
		if previousRun == nil {
			return ""
		}
		return previousRun.GraphStateJson
	}(), "")); previousTaskSession != nil {
		previousSession = *previousTaskSession
	}
	// 把运行时输入扩展成带消息、锚点和安全身份的新模型。
	out, err := s.runner.Run(ctx, agentruntime.RunInput{
		SecurityPrompt: g.Cfg().MustGet(ctx, "agent.securityPrompt", "").String(),
		SceneCode:      conv.SceneCode,
		Summary:        conv.SessionSummary,
		Messages:       rowsToMessages(recentRows),
		UserQuery:      strings.TrimSpace(userMsg.ContentText),
		RequestID:      requestIDFromContext(ctx),
		Message: agentruntime.UserMessage{
			ContentText:     strings.TrimSpace(userMsg.ContentText),
			MessageTypeCode: strings.TrimSpace(userMsg.MessageTypeCode),
			AssetIDs:        parseUint64Slice(userMsg.AssetIdsJson),
		},
		PreviousSession: previousSession,
		Conversation: agentruntime.ConversationAnchors{
			OrderNo:    strings.TrimSpace(conv.OrderNo),
			SubOrderNo: strings.TrimSpace(conv.SubOrderNo),
		},
		Security: agentruntime.SecurityContext{
			UserID:         conv.UserId,
			ShopNo:         strings.TrimSpace(conv.ShopNo),
			RequestID:      requestIDFromContext(ctx),
			ConversationNo: conv.ConversationNo,
			RunNo:          run.RunNo,
		},
		CheckpointHook: func(innerCtx context.Context, cp agentruntime.GraphCheckpointSnapshot) error {
			// 把当前图执行节点和状态持久化到 Run 真相源，避免 Pod 抖动后前端长期卡在生成中。
			return s.updateRunCheckpoint(innerCtx, run.RunNo, cp)
		},
	})
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 在数据库模型上追加精确过滤条件，确保只处理当前业务主体可见的数据。
		_, _ = dao.AgentRun.Ctx(ctx).Where(dao.AgentRun.Columns().RunNo, run.RunNo).Data(do.AgentRun{
			RunStatusCode:     runtimeStatusFailed,
			CurrentNodeCode:   firstNonEmpty(run.CurrentNodeCode, "GenerateAnswer"),
			ErrorCode:         "RUNNER_ERROR",
			ErrorMessage:      err.Error(),
			ToolResultStatus:  runtimeToolResultError,
			CheckpointVersion: gdb.Raw("checkpoint_version + 1"),
			FinishedAt:        gtime.Now(),
			UpdatedAt:         gtime.Now(),
		}).Update()
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, nil, gerror.Wrap(err, "run agent failed")
	}
	// 生成业务唯一编号，作为后续跨表关联、审计追踪和幂等定位的稳定主键。
	assistantMessageNo := generateBizNo("AMSG")
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	_, err = dao.AgentMessage.Ctx(ctx).Data(do.AgentMessage{
		MessageNo:       assistantMessageNo,
		ConversationNo:  conv.ConversationNo,
		RunNo:           run.RunNo,
		ReplyToTurnNo:   run.TurnNo,
		SenderTypeCode:  "AGENT",
		MessageTypeCode: "TEXT",
		ContentText:     out.AnswerText,
		Interrupted:     boolToInt(out.PromptInjection),
		SentAt:          now,
	}).Insert()
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, nil, gerror.Wrap(err, "create assistant message failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	runStatusCode := normalizeRunStatus(out.RunStatusCode, out.RiskDecisionCode)
	// 在数据库模型上追加精确过滤条件，确保只处理当前业务主体可见的数据。
	_, err = dao.AgentRun.Ctx(ctx).Where(dao.AgentRun.Columns().RunNo, run.RunNo).Data(do.AgentRun{
		RunStatusCode:           runStatusCode,
		CurrentNodeCode:         out.CurrentNodeCode,
		GraphStateJson:          normalizeJSON(out.GraphStateJSON),
		CheckpointVersion:       gdb.Raw("checkpoint_version + 1"),
		ToolWaitTimeoutAt:       toGTime(out.ToolWaitTimeoutAt),
		ToolResultStatus:        normalizeToolResultStatus(out.ToolResultStatus),
		DegradedReasonCode:      strings.TrimSpace(out.DegradedReasonCode),
		QueueBlocked:            boolToInt(out.QueueBlocked),
		QueueHintMessage:        strings.TrimSpace(out.QueueHintMessage),
		SubjectUserId:           conv.UserId,
		SubjectShopNo:           conv.ShopNo,
		RiskDecisionCode:        out.RiskDecisionCode,
		PromptInjectionFlag:     boolToInt(out.PromptInjection),
		PromptTokenEstimate:     out.PromptTokenCount,
		CompletionTokenEstimate: out.CompletionTokens,
		AnswerSourcesJson:       mustJSON(out.AnswerSources),
		FinishedAt:              now,
		UpdatedAt:               now,
	}).Update()
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, nil, gerror.Wrap(err, "update run result failed")
	}
	// 记录本次工具调用审计日志，保留主体身份、降级原因和返回结果供后续追踪。
	if err = s.logToolCall(ctx, conv, run, userMsg, out); err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, nil, err
	}
	// 在数据库模型上追加精确过滤条件，确保只处理当前业务主体可见的数据。
	_, _ = dao.AgentConversation.Ctx(ctx).Where(dao.AgentConversation.Columns().ConversationNo, conv.ConversationNo).Data(do.AgentConversation{
		LastRunStatusCode: runStatusCode,
		LastMessageAt:     now,
		UpdatedAt:         now,
	}).Update()
	// 按 Run 编号加载执行记录，确保后续查询和状态流转都基于指定运行实例。
	updatedRun, _ := s.getRunByNo(ctx, run.RunNo)
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var assistantMsg entity.AgentMessage
	// 从数据库读取最新真相源，避免后续逻辑继续依赖过期内存快照。
	_ = dao.AgentMessage.Ctx(ctx).Where(dao.AgentMessage.Columns().MessageNo, assistantMessageNo).Scan(&assistantMsg)
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return updatedRun, &assistantMsg, nil
}

func (s *sAgent) createEscalationTicket(ctx context.Context, conv *entity.AgentConversation, runNo, reasonCode, remark string) (*entity.AgentEscalationTicket, error) {
	// 拉取最近消息窗口，为总结、推理和转人工摘要提供最新上下文。
	recentRows, err := s.listRecentMessages(ctx, conv.ConversationNo, maxRecentMessages)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if strings.TrimSpace(reasonCode) == "" {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		reasonCode = "USER_REQUEST"
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	summary, err := s.runner.HandoffSummary(ctx, conv.SessionSummary, rowsToMessages(recentRows), reasonCode)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "generate handoff summary failed")
	}
	// 生成业务唯一编号，作为后续跨表关联、审计追踪和幂等定位的稳定主键。
	ticketNo := generateBizNo("AHT")
	// 记录当前业务时间，确保状态流转、排序展示和审计字段使用同一时间基线。
	now := gtime.Now()
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	_, err = dao.AgentEscalationTicket.Ctx(ctx).Data(do.AgentEscalationTicket{
		TicketNo:                 ticketNo,
		ConversationNo:           conv.ConversationNo,
		RunNo:                    strings.TrimSpace(runNo),
		UserId:                   conv.UserId,
		ShopNo:                   conv.ShopNo,
		EscalationReasonCode:     strings.ToUpper(strings.TrimSpace(reasonCode)),
		StatusCode:               conversationStatusEscalationPending,
		HandoffSummary:           summary,
		HandoffSummaryStatusCode: "SUCCESS",
		Remark:                   remark,
		HandoffGeneratedAt:       now,
		CreatedAt:                now,
		UpdatedAt:                now,
	}).Insert()
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "create escalation ticket failed")
	}
	// 在数据库模型上追加精确过滤条件，确保只处理当前业务主体可见的数据。
	_, _ = dao.AgentConversation.Ctx(ctx).Where(dao.AgentConversation.Columns().ConversationNo, conv.ConversationNo).Data(do.AgentConversation{
		IsHumanHandover:        0,
		ConversationStatusCode: conversationStatusEscalationPending,
		LastQueueNoticeAt:      now,
		UpdatedAt:              now,
	}).Update()
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var ticket entity.AgentEscalationTicket
	// 从数据库读取最新真相源，避免后续逻辑继续依赖过期内存快照。
	if err = dao.AgentEscalationTicket.Ctx(ctx).Where(dao.AgentEscalationTicket.Columns().TicketNo, ticketNo).Scan(&ticket); err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "query escalation ticket failed")
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &ticket, nil
}

func (s *sAgent) findConversation(ctx context.Context, userID uint64, sceneCode, shopNo, orderNo, subOrderNo, spuNo, skuNo string) (*entity.AgentConversation, error) {
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var row entity.AgentConversation
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	err := dao.AgentConversation.Ctx(ctx).
		Where(dao.AgentConversation.Columns().UserId, userID).
		Where(dao.AgentConversation.Columns().SceneCode, sceneCode).
		Where(dao.AgentConversation.Columns().ShopNo, strings.TrimSpace(shopNo)).
		Where(dao.AgentConversation.Columns().OrderNo, strings.TrimSpace(orderNo)).
		Where(dao.AgentConversation.Columns().SubOrderNo, strings.TrimSpace(subOrderNo)).
		Where(dao.AgentConversation.Columns().AnchorSpuNo, strings.TrimSpace(spuNo)).
		Where(dao.AgentConversation.Columns().AnchorSkuNo, strings.TrimSpace(skuNo)).
		OrderDesc(dao.AgentConversation.Columns().Id).
		Limit(1).
		Scan(&row)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "query assistant conversation failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if row.Id == 0 {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, nil
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &row, nil
}

func (s *sAgent) getConversationByNo(ctx context.Context, conversationNo string) (*entity.AgentConversation, error) {
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var row entity.AgentConversation
	// 从数据库读取最新真相源，避免后续逻辑继续依赖过期内存快照。
	if err := dao.AgentConversation.Ctx(ctx).Where(dao.AgentConversation.Columns().ConversationNo, conversationNo).Scan(&row); err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "query conversation failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if row.Id == 0 {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeNotFound, "conversation not found")
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &row, nil
}

func (s *sAgent) getLatestRun(ctx context.Context, conversationNo string) (*entity.AgentRun, error) {
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var row entity.AgentRun
	// 从数据库读取最新真相源，避免后续逻辑继续依赖过期内存快照。
	if err := dao.AgentRun.Ctx(ctx).Where(dao.AgentRun.Columns().ConversationNo, conversationNo).OrderDesc(dao.AgentRun.Columns().Id).Limit(1).Scan(&row); err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "query latest run failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if row.Id == 0 {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, nil
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &row, nil
}

func (s *sAgent) getPreviousRun(ctx context.Context, currentRun *entity.AgentRun) (*entity.AgentRun, error) {
	// 如果当前 Run 不存在，就没有可追溯的上一轮运行状态。
	if currentRun == nil || currentRun.Id == 0 {
		return nil, nil
	}
	// 先集中声明查询结果变量，便于统一做空值判断和收口。
	var row entity.AgentRun
	// 按当前 Run 的主键回溯上一条同会话运行记录，供多轮任务会话复用。
	if err := dao.AgentRun.Ctx(ctx).
		Where(dao.AgentRun.Columns().ConversationNo, currentRun.ConversationNo).
		WhereLT(dao.AgentRun.Columns().Id, currentRun.Id).
		OrderDesc(dao.AgentRun.Columns().Id).
		Limit(1).
		Scan(&row); err != nil {
		// 查询失败时继续向上返回统一错误，避免拿不到上一轮状态却静默降级。
		return nil, gerror.Wrap(err, "query previous run failed")
	}
	// 如果数据库里没有更早的 Run，就返回空值表示当前是首轮。
	if row.Id == 0 {
		return nil, nil
	}
	// 返回上一轮 Run 真相源，供当前轮恢复任务会话。
	return &row, nil
}

func (s *sAgent) getRunByNo(ctx context.Context, runNo string) (*entity.AgentRun, error) {
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var row entity.AgentRun
	// 从数据库读取最新真相源，避免后续逻辑继续依赖过期内存快照。
	if err := dao.AgentRun.Ctx(ctx).Where(dao.AgentRun.Columns().RunNo, runNo).Scan(&row); err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "query run failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if row.Id == 0 {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeNotFound, "run not found")
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &row, nil
}

func (s *sAgent) getLatestAssistantMessage(ctx context.Context, conversationNo string, run *entity.AgentRun) (*entity.AgentMessage, error) {
	// 在数据库模型上追加精确过滤条件，确保只处理当前业务主体可见的数据。
	model := dao.AgentMessage.Ctx(ctx).Where(dao.AgentMessage.Columns().ConversationNo, conversationNo).Where(dao.AgentMessage.Columns().SenderTypeCode, "AGENT")
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if run != nil && run.RunNo != "" {
		// 在数据库模型上追加精确过滤条件，确保只处理当前业务主体可见的数据。
		model = model.Where(dao.AgentMessage.Columns().RunNo, run.RunNo)
	}
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var row entity.AgentMessage
	// 从数据库读取最新真相源，避免后续逻辑继续依赖过期内存快照。
	if err := model.OrderDesc(dao.AgentMessage.Columns().Id).Limit(1).Scan(&row); err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "query latest assistant message failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if row.Id == 0 {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, nil
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &row, nil
}

func (s *sAgent) getLatestBuyerMessage(ctx context.Context, conversationNo string) (*entity.AgentMessage, error) {
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var row entity.AgentMessage
	// 从数据库读取最新真相源，避免后续逻辑继续依赖过期内存快照。
	if err := dao.AgentMessage.Ctx(ctx).Where(dao.AgentMessage.Columns().ConversationNo, conversationNo).Where(dao.AgentMessage.Columns().SenderTypeCode, "BUYER").OrderDesc(dao.AgentMessage.Columns().Id).Limit(1).Scan(&row); err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "query latest buyer message failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if row.Id == 0 {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeNotFound, "buyer message not found")
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &row, nil
}

func (s *sAgent) getMessageByClientNo(ctx context.Context, conversationNo, clientMessageNo string) (*entity.AgentMessage, error) {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if strings.TrimSpace(clientMessageNo) == "" {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, nil
	}
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var row entity.AgentMessage
	// 从数据库读取最新真相源，避免后续逻辑继续依赖过期内存快照。
	if err := dao.AgentMessage.Ctx(ctx).Where(dao.AgentMessage.Columns().ConversationNo, conversationNo).Where(dao.AgentMessage.Columns().ClientMessageNo, clientMessageNo).Scan(&row); err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "query message by client_message_no failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if row.Id == 0 {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, nil
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &row, nil
}

func (s *sAgent) listRecentMessages(ctx context.Context, conversationNo string, limit int) ([]entity.AgentMessage, error) {
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var rows []entity.AgentMessage
	// 从数据库读取最新真相源，避免后续逻辑继续依赖过期内存快照。
	if err := dao.AgentMessage.Ctx(ctx).Where(dao.AgentMessage.Columns().ConversationNo, conversationNo).OrderDesc(dao.AgentMessage.Columns().Id).Limit(limit).Scan(&rows); err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "query recent messages failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	sort.SliceStable(rows, func(i, j int) bool {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return rows[i].Id < rows[j].Id
	})
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return rows, nil
}

func (s *sAgent) getKnowledgeDocByNo(ctx context.Context, knowledgeDocNo string) (*entity.AgentKnowledgeDoc, error) {
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var row entity.AgentKnowledgeDoc
	// 从数据库读取最新真相源，避免后续逻辑继续依赖过期内存快照。
	if err := dao.AgentKnowledgeDoc.Ctx(ctx).Where(dao.AgentKnowledgeDoc.Columns().KnowledgeDocNo, knowledgeDocNo).Scan(&row); err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "query knowledge doc failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if row.Id == 0 {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeNotFound, "knowledge doc not found")
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &row, nil
}

func rowsToMessages(rows []entity.AgentMessage) []*schema.Message {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	out := make([]*schema.Message, 0, len(rows))
	// 遍历当前集合或循环条件，逐项推进本段业务处理并累计最终结果。
	for i := range rows {
		// 统一把编码转换成约定的大写形式，避免状态码和场景码因大小写不一致而失配。
		switch strings.ToUpper(rows[i].SenderTypeCode) {
		// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
		case "BUYER":
			// 把当前元素追加到结果集合中，逐步汇总出最终返回或落库的数据集。
			out = append(out, schema.UserMessage(rows[i].ContentText))
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		default:
			// 把当前元素追加到结果集合中，逐步汇总出最终返回或落库的数据集。
			out = append(out, schema.AssistantMessage(rows[i].ContentText, nil))
		}
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return out
}

func parseAnswerSources(run *entity.AgentRun) []*agentv1.AnswerSource {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if run == nil || strings.TrimSpace(run.AnswerSourcesJson) == "" {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return nil
	}
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var items []agentruntime.AnswerSource
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if err := json.Unmarshal([]byte(run.AnswerSourcesJson), &items); err != nil {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return nil
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	out := make([]*agentv1.AnswerSource, 0, len(items))
	// 遍历当前集合或循环条件，逐项推进本段业务处理并累计最终结果。
	for _, item := range items {
		// 把当前元素追加到结果集合中，逐步汇总出最终返回或落库的数据集。
		out = append(out, &agentv1.AnswerSource{
			SourceTypeCode: item.SourceTypeCode,
			SourceId:       item.SourceID,
			SourceVersion:  item.SourceVersion,
			Title:          item.Title,
			Snippet:        item.Snippet,
		})
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return out
}

func parseReplyPayload(run *entity.AgentRun) *agentv1.AssistantReplyPayload {
	// 如果当前 Run 不存在，就直接返回空 payload。
	if run == nil || strings.TrimSpace(run.GraphStateJson) == "" {
		return nil
	}
	// 从运行图快照里提取结构化回复载荷，避免再维护一份重复存储。
	payload := agentruntime.ParseGraphReplyPayload(run.GraphStateJson)
	// 如果快照里没有 reply payload，就维持空值返回。
	if payload == nil {
		return nil
	}
	// 先把内部 data card 映射成对外 proto 结构，供前端直接渲染。
	dataCards := make([]*agentv1.AssistantDataCard, 0, len(payload.DataCards))
	// 逐项复制卡片内容，避免直接暴露内部运行时结构。
	for _, card := range payload.DataCards {
		dataCards = append(dataCards, &agentv1.AssistantDataCard{
			OrderSnapshotCard:     toProtoOrderSnapshotCard(card.OrderSnapshotCard),
			AfterSaleDecisionCard: toProtoAfterSaleDecisionCard(card.AfterSaleDecisionCard),
		})
	}
	// 继续把建议动作映射成稳定的协议对象。
	actions := make([]*agentv1.SuggestedAction, 0, len(payload.SuggestedActions))
	// 逐项复制动作信息，供前端控制按钮状态和禁用原因。
	for _, action := range payload.SuggestedActions {
		actions = append(actions, &agentv1.SuggestedAction{
			ActionCode:       action.ActionCode,
			Label:            action.Label,
			Enabled:          action.Enabled,
			ReasonIfDisabled: action.ReasonIfDisabled,
		})
	}
	// 也把缺槽信息映射出来，方便前端渲染追问提示。
	missingSlots := make([]*agentv1.MissingSlot, 0, len(payload.MissingSlots))
	// 逐项复制缺槽定义，确保 required 标记和 prompt 一并透出。
	for _, slot := range payload.MissingSlots {
		missingSlots = append(missingSlots, &agentv1.MissingSlot{
			SlotCode:   slot.SlotCode,
			PromptText: slot.PromptText,
			Required:   slot.Required,
		})
	}
	// 返回完整的结构化 reply payload，供 Run/Status/Send 接口复用。
	return &agentv1.AssistantReplyPayload{
		ReplyText:          payload.ReplyText,
		IntentCode:         payload.IntentCode,
		GuardResultCode:    payload.GuardResultCode,
		Confidence:         payload.Confidence,
		DataCards:          dataCards,
		SuggestedActions:   actions,
		MissingSlots:       missingSlots,
		SlotRetryCount:     payload.SlotRetryCount,
		HandoffRecommended: payload.HandoffRecommended,
		HandoffReasonCode:  payload.HandoffReasonCode,
	}
}

func toProtoOrderSnapshotCard(card *agentruntime.OrderSnapshotCard) *agentv1.OrderSnapshotCard {
	// 如果没有订单卡片，就不返回 proto 卡片。
	if card == nil {
		return nil
	}
	// 把内部订单卡片结构映射成对外协议对象。
	return &agentv1.OrderSnapshotCard{
		OrderNo:           card.OrderNo,
		MainStatus:        card.MainStatus,
		PaymentStatus:     card.PaymentStatus,
		FulfillmentStatus: card.FulfillmentStatus,
		LogisticsStatus:   card.LogisticsStatus,
		AfterSaleStatus:   card.AfterSaleStatus,
		LatestUpdateTime:  card.LatestUpdateTime,
	}
}

func toProtoAfterSaleDecisionCard(card *agentruntime.AfterSaleDecisionCard) *agentv1.AfterSaleDecisionCard {
	// 如果没有决策卡片，就不返回 proto 卡片。
	if card == nil {
		return nil
	}
	// 把内部决策卡片映射成对外协议对象。
	return &agentv1.AfterSaleDecisionCard{
		DecisionPathCode: card.DecisionPathCode,
		ReasonText:       card.ReasonText,
		ConstraintText:   card.ConstraintText,
		NextStepText:     card.NextStepText,
		SceneCode:        card.SceneCode,
	}
}

func toProtoConversation(row *entity.AgentConversation) *agentv1.AssistantConversation {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if row == nil || row.Id == 0 {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return nil
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &agentv1.AssistantConversation{
		ConversationNo:          row.ConversationNo,
		BotCode:                 row.BotCode,
		SceneCode:               row.SceneCode,
		UserId:                  row.UserId,
		ShopNo:                  row.ShopNo,
		OrderNo:                 row.OrderNo,
		SubOrderNo:              row.SubOrderNo,
		AnchorSpuNo:             row.AnchorSpuNo,
		AnchorSkuNo:             row.AnchorSkuNo,
		ConversationStatusCode:  toProtoConversationStatus(row.ConversationStatusCode),
		LastRunStatusCode:       toProtoRunStatus(row.LastRunStatusCode),
		IsHumanHandover:         row.IsHumanHandover == 1,
		PendingMessageCount:     row.PendingMessageCount,
		SessionSummary:          row.SessionSummary,
		SessionSummaryVersion:   row.SessionSummaryVersion,
		LastSummarizedMessageNo: row.LastSummarizedMessageNo,
		CreatedAt:               toProtoTs(row.CreatedAt),
		UpdatedAt:               toProtoTs(row.UpdatedAt),
		LastQueueNoticeAt:       toProtoTs(row.LastQueueNoticeAt),
	}
}

func toProtoRun(row *entity.AgentRun) *agentv1.AssistantRun {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if row == nil || row.Id == 0 {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return nil
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &agentv1.AssistantRun{
		RunNo:                   row.RunNo,
		ConversationNo:          row.ConversationNo,
		TurnNo:                  row.TurnNo,
		AcceptedStatusCode:      row.AcceptedStatusCode,
		RunStatusCode:           toProtoRunStatus(row.RunStatusCode),
		RiskDecisionCode:        row.RiskDecisionCode,
		PromptInjectionFlag:     row.PromptInjectionFlag == 1,
		ReplyInterrupted:        row.ReplyInterrupted == 1,
		InterruptReasonCode:     row.InterruptReasonCode,
		MergedMessageCount:      uint32(row.MergedMessageCount),
		ToolIterationCount:      uint32(row.ToolIterationCount),
		ToolCallCount:           uint32(row.ToolCallCount),
		PromptTokenEstimate:     uint32(row.PromptTokenEstimate),
		CompletionTokenEstimate: uint32(row.CompletionTokenEstimate),
		ErrorCode:               row.ErrorCode,
		ErrorMessage:            row.ErrorMessage,
		CreatedAt:               toProtoTs(row.CreatedAt),
		StartedAt:               toProtoTs(row.StartedAt),
		FinishedAt:              toProtoTs(row.FinishedAt),
		CurrentNodeCode:         row.CurrentNodeCode,
		ToolResultStatus:        toProtoToolResultStatus(row.ToolResultStatus),
		QueueBlocked:            row.QueueBlocked == 1,
		QueueHintMessage:        row.QueueHintMessage,
		DegradedReasonCode:      row.DegradedReasonCode,
		ReplyPayload:            parseReplyPayload(row),
	}
}

func toProtoMessage(row *entity.AgentMessage) *agentv1.AssistantMessage {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if row == nil || row.Id == 0 {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return nil
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &agentv1.AssistantMessage{
		MessageNo:       row.MessageNo,
		ConversationNo:  row.ConversationNo,
		RunNo:           row.RunNo,
		SenderTypeCode:  row.SenderTypeCode,
		MessageTypeCode: row.MessageTypeCode,
		ClientMessageNo: row.ClientMessageNo,
		ContentText:     row.ContentText,
		AssetIds:        parseUint64Slice(row.AssetIdsJson),
		ExtJson:         row.ExtJson,
		Interrupted:     row.Interrupted == 1,
		ReplyToTurnNo:   row.ReplyToTurnNo,
		SentAt:          toProtoTs(row.SentAt),
	}
}

func toProtoTicket(row *entity.AgentEscalationTicket) *agentv1.HandoffTicket {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if row == nil || row.Id == 0 {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return nil
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &agentv1.HandoffTicket{
		TicketNo:                 row.TicketNo,
		ConversationNo:           row.ConversationNo,
		UserId:                   row.UserId,
		ShopNo:                   row.ShopNo,
		EscalationReasonCode:     row.EscalationReasonCode,
		StatusCode:               row.StatusCode,
		HandoffSummary:           row.HandoffSummary,
		HandoffSummaryStatusCode: row.HandoffSummaryStatusCode,
		OperatorUserId:           row.OperatorUserId,
		HandoffGeneratedAt:       toProtoTs(row.HandoffGeneratedAt),
		CreatedAt:                toProtoTs(row.CreatedAt),
		UpdatedAt:                toProtoTs(row.UpdatedAt),
		AcceptedAt:               toProtoTs(row.AcceptedAt),
		ClosedAt:                 toProtoTs(row.ClosedAt),
		QueueAppendixJson:        row.QueueAppendixJson,
	}
}

func toProtoKnowledgeDoc(row *entity.AgentKnowledgeDoc) *agentv1.KnowledgeDoc {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if row == nil || row.Id == 0 {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return nil
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &agentv1.KnowledgeDoc{
		KnowledgeDocNo:     row.KnowledgeDocNo,
		Title:              row.Title,
		KnowledgeScopeCode: row.KnowledgeScopeCode,
		SourceTypeCode:     row.SourceTypeCode,
		SourceId:           row.SourceId,
		SourceVersion:      row.SourceVersion,
		ShopNo:             row.ShopNo,
		ContentTypeCode:    row.ContentTypeCode,
		ContentText:        row.ContentText,
		AssetIds:           parseUint64Slice(row.AssetIdsJson),
		Tags:               parseStringSlice(row.TagsJson),
		StatusCode:         row.StatusCode,
		RejectReasonCode:   row.RejectReasonCode,
		RejectComment:      row.RejectComment,
		Version:            row.Version,
		PublishedAt:        toProtoTs(row.PublishedAt),
		CreatedAt:          toProtoTs(row.CreatedAt),
		UpdatedAt:          toProtoTs(row.UpdatedAt),
	}
}

func toProtoKnowledgeChunk(row *entity.AgentKnowledgeChunk) *agentv1.KnowledgeChunk {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if row == nil || row.Id == 0 {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return nil
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &agentv1.KnowledgeChunk{
		ChunkNo:          row.ChunkNo,
		KnowledgeDocNo:   row.KnowledgeDocNo,
		SourceTypeCode:   row.SourceTypeCode,
		SourceId:         row.SourceId,
		SourceVersion:    row.SourceVersion,
		ChunkIndex:       uint32(row.ChunkIndex),
		ChunkText:        row.ChunkText,
		ChunkTextPreview: row.ChunkTextPreview,
		IsDeleted:        row.IsDeleted == 1,
		IndexStatusCode:  row.IndexStatusCode,
		CreatedAt:        toProtoTs(row.CreatedAt),
		UpdatedAt:        toProtoTs(row.UpdatedAt),
	}
}

func (s *sAgent) handleQueueBlockedMessage(ctx context.Context, conv *entity.AgentConversation, userMsg *entity.AgentMessage) (*entity.AgentMessage, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if conv == nil || userMsg == nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "conversation and user message are required")
	}
	// 记录当前业务时间，确保状态流转、排序展示和审计字段使用同一时间基线。
	now := gtime.Now()
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	queueHint := "您已进入人工客服队列，请稍作等待，我们会尽快为您接入人工客服。"

	// 把排队期间补充的用户消息整理进工单附录，方便人工接手后快速补齐上下文。
	appendixJSON, err := s.buildQueueAppendix(ctx, conv.ConversationNo, userMsg, now)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if appendixJSON != "" {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		_, _ = dao.AgentEscalationTicket.Ctx(ctx).
			Where(dao.AgentEscalationTicket.Columns().ConversationNo, conv.ConversationNo).
			Where(dao.AgentEscalationTicket.Columns().StatusCode, conversationStatusEscalationPending).
			OrderDesc(dao.AgentEscalationTicket.Columns().Id).
			Limit(1).
			Data(do.AgentEscalationTicket{
				QueueAppendixJson: appendixJSON,
				UpdatedAt:         now,
			}).
			Update()
	}

	// 生成业务唯一编号，作为后续跨表关联、审计追踪和幂等定位的稳定主键。
	systemMessageNo := generateBizNo("AMSG")
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	_, err = dao.AgentMessage.Ctx(ctx).Data(do.AgentMessage{
		MessageNo:       systemMessageNo,
		ConversationNo:  conv.ConversationNo,
		SenderTypeCode:  "SYSTEM",
		MessageTypeCode: "QUEUE_NOTICE",
		ContentText:     queueHint,
		ExtJson:         normalizeJSON(`{"queue_blocked":true}`),
		SentAt:          now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}).Insert()
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "create queue placeholder message failed")
	}

	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	_, err = dao.AgentConversation.Ctx(ctx).
		Where(dao.AgentConversation.Columns().ConversationNo, conv.ConversationNo).
		Data(do.AgentConversation{
			PendingMessageCount: conv.PendingMessageCount + 1,
			LastQueueNoticeAt:   now,
			LastMessageAt:       now,
			UpdatedAt:           now,
		}).
		Update()
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "update queue blocked conversation failed")
	}

	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var systemMsg entity.AgentMessage
	// 从数据库读取最新真相源，避免后续逻辑继续依赖过期内存快照。
	if err = dao.AgentMessage.Ctx(ctx).Where(dao.AgentMessage.Columns().MessageNo, systemMessageNo).Scan(&systemMsg); err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "query queue placeholder message failed")
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &systemMsg, nil
}

func (s *sAgent) buildQueueAppendix(ctx context.Context, conversationNo string, userMsg *entity.AgentMessage, now *gtime.Time) (string, error) {
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var ticket entity.AgentEscalationTicket
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if err := dao.AgentEscalationTicket.Ctx(ctx).
		Where(dao.AgentEscalationTicket.Columns().ConversationNo, conversationNo).
		OrderDesc(dao.AgentEscalationTicket.Columns().Id).
		Limit(1).
		Scan(&ticket); err != nil {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return "", gerror.Wrap(err, "query escalation ticket failed")
	}

	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	items := make([]queueAppendixItem, 0)
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if strings.TrimSpace(ticket.QueueAppendixJson) != "" {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		_ = json.Unmarshal([]byte(ticket.QueueAppendixJson), &items)
	}
	// 把当前元素追加到结果集合中，逐步汇总出最终返回或落库的数据集。
	items = append(items, queueAppendixItem{
		MessageNo: userMsg.MessageNo,
		Content:   strings.TrimSpace(userMsg.ContentText),
		SentAt:    now.Format("Y-m-d H:i:s"),
	})
	// 把复杂结构序列化成 JSON 字符串，便于在当前表结构下稳定存储和回显。
	return mustJSON(items), nil
}

func userIDFromContext(ctx context.Context) (uint64, error) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if r := g.RequestFromCtx(ctx); r != nil {
		// 遍历当前集合或循环条件，逐项推进本段业务处理并累计最终结果。
		for _, key := range []string{"x-user-id", "X-User-Id", "user_id", "uid"} {
			// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
			if raw := strings.TrimSpace(r.Header.Get(key)); raw != "" {
				// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
				uid, err := strconv.ParseUint(raw, 10, 64)
				// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
				if err == nil && uid > 0 {
					// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
					return uid, nil
				}
			}
		}
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	md := grpcx.Ctx.IncomingMap(ctx)
	// 遍历当前集合或循环条件，逐项推进本段业务处理并累计最终结果。
	for _, key := range []string{"x-user-id", "user_id", "uid", "userid"} {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		if val := md.Get(key); val != nil {
			// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
			uid := gconv.Uint64(val)
			// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
			if uid > 0 {
				// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
				return uid, nil
			}
		}
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return 0, gerror.NewCode(gcode.CodeNotAuthorized, "missing x-user-id")
}

func toProtoTs(t *gtime.Time) *timestamppb.Timestamp {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if t == nil {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return nil
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return timestamppb.New(t.Time)
}

func parseCursor(cursor string) (uint64, error) {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	cursor = strings.TrimSpace(cursor)
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if cursor == "" {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return 0, nil
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	id, err := strconv.ParseUint(cursor, 10, 64)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return 0, gerror.WrapCode(gcode.CodeInvalidParameter, err, "next_cursor must be uint64")
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return id, nil
}

func normalizePageSize(reqSize int32) int {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	size := int(reqSize)
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if size <= 0 {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		size = defaultPageSize
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if size > maxPageSize {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		size = maxPageSize
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return size
}

func normalizeOffsetPage(page, pageSize int32) (int, int) {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	p := int(page)
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if p <= 0 {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		p = 1
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return p, normalizePageSize(pageSize)
}

func normalizeMessageType(value string) string {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	v := strings.ToUpper(strings.TrimSpace(value))
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if v == "" {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return "TEXT"
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return v
}

func normalizeRunStatus(runStatus, riskDecision string) string {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	runStatus = strings.ToUpper(strings.TrimSpace(runStatus))
	// 根据当前关键状态或意图进入不同分支，确保每个业务场景按对应规则处理。
	switch runStatus {
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case runtimeStatusPending, runtimeStatusProcessing, runtimeStatusWaitingToolCall, runtimeStatusGenerating, runtimeStatusSuccess, runtimeStatusFailed, runtimeStatusAborted, runtimeStatusEscalated:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return runStatus
	}
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if strings.EqualFold(strings.TrimSpace(riskDecision), "ESCALATE") {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return runtimeStatusEscalated
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return runtimeStatusSuccess
}

func normalizeToolResultStatus(status string) string {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	status = strings.ToUpper(strings.TrimSpace(status))
	// 根据当前关键状态或意图进入不同分支，确保每个业务场景按对应规则处理。
	switch status {
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case runtimeToolResultSuccess, runtimeToolResultDegraded, runtimeToolResultError:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return status
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	default:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return runtimeToolResultUnspecified
	}
}

func toProtoRunStatus(status string) agentv1.AgentRunStatusCode {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	switch strings.ToUpper(strings.TrimSpace(status)) {
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case runtimeStatusPending:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return agentv1.AgentRunStatusCode_AGENT_RUN_STATUS_CODE_PENDING
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case runtimeStatusProcessing:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return agentv1.AgentRunStatusCode_AGENT_RUN_STATUS_CODE_PROCESSING
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case runtimeStatusWaitingToolCall:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return agentv1.AgentRunStatusCode_AGENT_RUN_STATUS_CODE_WAITING_TOOL_CALL
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case runtimeStatusGenerating:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return agentv1.AgentRunStatusCode_AGENT_RUN_STATUS_CODE_GENERATING
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case runtimeStatusSuccess:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return agentv1.AgentRunStatusCode_AGENT_RUN_STATUS_CODE_SUCCESS
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case runtimeStatusFailed:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return agentv1.AgentRunStatusCode_AGENT_RUN_STATUS_CODE_FAILED
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case runtimeStatusAborted:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return agentv1.AgentRunStatusCode_AGENT_RUN_STATUS_CODE_ABORTED
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case runtimeStatusEscalated:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return agentv1.AgentRunStatusCode_AGENT_RUN_STATUS_CODE_ESCALATED
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	default:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return agentv1.AgentRunStatusCode_AGENT_RUN_STATUS_CODE_UNSPECIFIED
	}
}

func toProtoConversationStatus(status string) agentv1.ConversationStatusCode {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	switch strings.ToUpper(strings.TrimSpace(status)) {
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case conversationStatusActive:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return agentv1.ConversationStatusCode_CONVERSATION_STATUS_CODE_ACTIVE
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case conversationStatusEscalationPending:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return agentv1.ConversationStatusCode_CONVERSATION_STATUS_CODE_ESCALATION_PENDING
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case conversationStatusEscalationAccepted:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return agentv1.ConversationStatusCode_CONVERSATION_STATUS_CODE_ESCALATION_ACCEPTED
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case conversationStatusHumanHandover:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return agentv1.ConversationStatusCode_CONVERSATION_STATUS_CODE_HUMAN_HANDOVER
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case conversationStatusHumanClosed:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return agentv1.ConversationStatusCode_CONVERSATION_STATUS_CODE_HUMAN_CLOSED
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case conversationStatusArchived:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return agentv1.ConversationStatusCode_CONVERSATION_STATUS_CODE_ARCHIVED
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	default:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return agentv1.ConversationStatusCode_CONVERSATION_STATUS_CODE_UNSPECIFIED
	}
}

func toProtoToolResultStatus(status string) agentv1.McpToolResultStatus {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	switch strings.ToUpper(strings.TrimSpace(status)) {
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case runtimeToolResultSuccess:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return agentv1.McpToolResultStatus_MCP_TOOL_RESULT_STATUS_SUCCESS
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case runtimeToolResultDegraded:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return agentv1.McpToolResultStatus_MCP_TOOL_RESULT_STATUS_DEGRADED
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case runtimeToolResultError:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return agentv1.McpToolResultStatus_MCP_TOOL_RESULT_STATUS_ERROR
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	default:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return agentv1.McpToolResultStatus_MCP_TOOL_RESULT_STATUS_UNSPECIFIED
	}
}

func toProtoCheckpoint(row *entity.AgentRun) *agentv1.GraphCheckpoint {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if row == nil || row.Id == 0 {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return nil
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &agentv1.GraphCheckpoint{
		RunNo:             row.RunNo,
		CurrentNodeCode:   row.CurrentNodeCode,
		CheckpointStatus:  toProtoRunStatus(row.RunStatusCode),
		GraphStateJson:    row.GraphStateJson,
		ToolWaitTimeoutAt: toProtoTs(row.ToolWaitTimeoutAt),
	}
}

func requestIDFromContext(ctx context.Context) string {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if r := g.RequestFromCtx(ctx); r != nil {
		// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
		if value := strings.TrimSpace(r.Header.Get("x-request-id")); value != "" {
			// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
			return value
		}
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	md := grpcx.Ctx.IncomingMap(ctx)
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if value := strings.TrimSpace(gconv.String(md.Get("x-request-id"))); value != "" {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return value
	}
	// 生成业务唯一编号，作为后续跨表关联、审计追踪和幂等定位的稳定主键。
	return generateBizNo("REQ")
}

func firstNonEmpty(values ...string) string {
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

func toGTime(value *time.Time) *gtime.Time {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if value == nil {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return nil
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return gtime.NewFromTime(*value)
}

func (s *sAgent) updateRunCheckpoint(ctx context.Context, runNo string, cp agentruntime.GraphCheckpointSnapshot) error {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if strings.TrimSpace(runNo) == "" {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return nil
	}
	// 在数据库模型上追加精确过滤条件，确保只处理当前业务主体可见的数据。
	_, err := dao.AgentRun.Ctx(ctx).Where(dao.AgentRun.Columns().RunNo, runNo).Data(do.AgentRun{
		CheckpointVersion:  gdb.Raw("checkpoint_version + 1"),
		CurrentNodeCode:    strings.TrimSpace(cp.CurrentNodeCode),
		RunStatusCode:      normalizeRunStatus(cp.CheckpointStatus, ""),
		GraphStateJson:     normalizeJSON(cp.GraphStateJSON),
		ToolWaitTimeoutAt:  toGTime(cp.ToolWaitTimeoutAt),
		ToolResultStatus:   normalizeToolResultStatus(cp.ToolResultStatus),
		DegradedReasonCode: strings.TrimSpace(cp.DegradedReasonCode),
		QueueBlocked:       boolToInt(cp.QueueBlocked),
		QueueHintMessage:   strings.TrimSpace(cp.QueueHintMessage),
		UpdatedAt:          gtime.Now(),
	}).Update()
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return gerror.Wrap(err, "update run checkpoint failed")
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return nil
}

func (s *sAgent) logToolCall(ctx context.Context, conv *entity.AgentConversation, run *entity.AgentRun, userMsg *entity.AgentMessage, out *agentruntime.RunOutput) error {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if conv == nil || run == nil || out == nil || strings.TrimSpace(out.SelectedToolName) == "" {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return nil
	}
	// 记录当前业务时间，确保状态流转、排序展示和审计字段使用同一时间基线。
	now := gtime.Now()
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	_, err := dao.AgentToolCallLog.Ctx(ctx).Data(do.AgentToolCallLog{
		ConversationNo:      conv.ConversationNo,
		RunNo:               run.RunNo,
		ToolName:            strings.TrimSpace(out.SelectedToolName),
		AdapterCode:         strings.TrimSpace(out.SelectedAdapterCode),
		ToolScopeCode:       strings.TrimSpace(out.SelectedToolScope),
		SubjectUserId:       conv.UserId,
		SubjectShopNo:       conv.ShopNo,
		RequestPayloadJson:  mustJSON(g.Map{"query": strings.TrimSpace(userMsg.ContentText)}),
		ResponsePayloadJson: mustJSON(g.Map{"answer": out.AnswerText, "degraded_reason_code": out.DegradedReasonCode}),
		ToolResultStatus:    normalizeToolResultStatus(out.ToolResultStatus),
		DegradedReasonCode:  strings.TrimSpace(out.DegradedReasonCode),
		DurationMs:          0,
		CreatedAt:           now,
		UpdatedAt:           now,
	}).Insert()
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return gerror.Wrap(err, "create tool call log failed")
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return nil
}

func isQueueBlockedConversation(status string) bool {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	switch strings.ToUpper(strings.TrimSpace(status)) {
	// 命中当前分支后执行专属处理逻辑，保证状态机和意图路由语义清晰可追踪。
	case conversationStatusEscalationPending, conversationStatusEscalationAccepted, conversationStatusHumanHandover:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return true
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	default:
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return false
	}
}

type queueAppendixItem struct {
	MessageNo string `json:"message_no"`
	Content   string `json:"content"`
	SentAt    string `json:"sent_at"`
}

func boolToInt(v bool) int {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if v {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return 1
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return 0
}

func mustJSON(v any) string {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if v == nil {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return "[]"
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	b, _ := json.Marshal(v)
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return string(b)
}

func normalizeJSON(v string) string {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	v = strings.TrimSpace(v)
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if v == "" {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return "{}"
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return v
}

func parseUint64Slice(raw string) []uint64 {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if strings.TrimSpace(raw) == "" {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return nil
	}
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var out []uint64
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	_ = json.Unmarshal([]byte(raw), &out)
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return out
}

func parseStringSlice(raw string) []string {
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if strings.TrimSpace(raw) == "" {
		// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
		return nil
	}
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var out []string
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	_ = json.Unmarshal([]byte(raw), &out)
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return out
}

func fieldMaskPaths(paths []string) []string {
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	out := make([]string, 0, len(paths))
	// 遍历当前集合或循环条件，逐项推进本段业务处理并累计最终结果。
	for _, item := range paths {
		// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
		item = strings.TrimSpace(item)
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		if item != "" {
			// 把当前元素追加到结果集合中，逐步汇总出最终返回或落库的数据集。
			out = append(out, strings.ToLower(item))
		}
	}
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return out
}

func generateBizNo(prefix string) string {
	// 记录当前业务时间，确保状态流转、排序展示和审计字段使用同一时间基线。
	now := time.Now()
	// 在当前分支完成收口并返回结果，避免后续逻辑继续执行造成状态污染。
	return fmt.Sprintf("%s%s%06d", prefix, now.Format("20060102150405"), now.UnixNano()%1000000)
}
