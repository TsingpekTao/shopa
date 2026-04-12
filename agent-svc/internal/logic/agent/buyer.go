package agent

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	agentv1 "github.com/TsingpekTao/shopa/agent-svc/api/v1"
	"github.com/TsingpekTao/shopa/agent-svc/internal/dao"
	"github.com/TsingpekTao/shopa/agent-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/agent-svc/internal/model/entity"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *sAgent) CreateOrGetAssistantConversation(ctx context.Context, req *agentv1.CreateOrGetAssistantConversationReq) (*agentv1.CreateOrGetAssistantConversationRes, error) {
	// 先校验请求对象是否存在，避免后续读取空请求字段时触发空指针并污染业务语义。
	if req == nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "request is required")
	}
	// 从上下文提取当前登录用户身份，保证后续查询和写入都严格绑定当前会话主体。
	userID, err := userIDFromContext(ctx)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	persistenceCtx, cancel := context.WithTimeout(conversationPersistenceContext(ctx), 5*time.Second)
	defer cancel()
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	sceneCode := strings.ToUpper(strings.TrimSpace(req.GetSceneCode()))
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if sceneCode == "" {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		sceneCode = "BUYER_ASSISTANT"
	}
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	botCode := strings.TrimSpace(req.GetBotCode())
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if botCode == "" {
		// 从配置中心读取当前运行参数，让模型、安全策略和业务默认值都可配置化演进。
		botCode = g.Cfg().MustGet(ctx, "agent.botCode", "BUYER_ASSISTANT").String()
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if !req.GetForceNew() {
		// 优先按场景和锚点查找历史会话，复用既有上下文以满足场景续接策略。
		existing, err := s.findConversation(persistenceCtx, userID, sceneCode, req.GetShopNo(), req.GetOrderNo(), req.GetSubOrderNo(), req.GetAnchorSpuNo(), req.GetAnchorSkuNo())
		// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
		if err != nil {
			// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
			return nil, err
		}
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		if existing != nil {
			// 读取当前会话最近一次 Run，便于复用状态、判断 turn_no 或回传最新执行结果。
			latestRun, _ := s.getLatestRun(persistenceCtx, existing.ConversationNo)
			// 把内部实体或运行态结构转换成对外协议对象，保证对外契约稳定且隔离内部实现。
			return &agentv1.CreateOrGetAssistantConversationRes{Conversation: toProtoConversation(existing), LatestRun: toProtoRun(latestRun)}, nil
		}
	}
	// 生成业务唯一编号，作为后续跨表关联、审计追踪和幂等定位的稳定主键。
	conversationNo := generateBizNo("ACV")
	// 记录当前业务时间，确保状态流转、排序展示和审计字段使用同一时间基线。
	now := gtime.Now()
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	_, err = dao.AgentConversation.Ctx(persistenceCtx).Data(do.AgentConversation{
		ConversationNo:         conversationNo,
		UserId:                 userID,
		BotCode:                botCode,
		SceneCode:              sceneCode,
		ShopNo:                 strings.TrimSpace(req.GetShopNo()),
		OrderNo:                strings.TrimSpace(req.GetOrderNo()),
		SubOrderNo:             strings.TrimSpace(req.GetSubOrderNo()),
		AnchorSpuNo:            strings.TrimSpace(req.GetAnchorSpuNo()),
		AnchorSkuNo:            strings.TrimSpace(req.GetAnchorSkuNo()),
		ConversationStatusCode: conversationStatusActive,
		LastRunStatusCode:      runtimeStatusPending,
		LastMessageAt:          now,
	}).Insert()
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "create assistant conversation failed")
	}
	// 按会话号读取当前会话真相源，避免后续逻辑基于过期内存状态继续执行。
	conv, err := s.getConversationByNo(persistenceCtx, conversationNo)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 把内部实体或运行态结构转换成对外协议对象，保证对外契约稳定且隔离内部实现。
	return &agentv1.CreateOrGetAssistantConversationRes{Conversation: toProtoConversation(conv)}, nil
}

func (s *sAgent) SendAssistantMessage(ctx context.Context, req *agentv1.SendAssistantMessageReq) (*agentv1.SendAssistantMessageRes, error) {
	// 先校验请求对象是否存在，避免后续读取空请求字段时触发空指针并污染业务语义。
	if req == nil || strings.TrimSpace(req.GetConversationNo()) == "" {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "conversation_no is required")
	}
	// 从上下文提取当前登录用户身份，保证后续查询和写入都严格绑定当前会话主体。
	userID, err := userIDFromContext(ctx)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 按会话号读取当前会话真相源，避免后续逻辑基于过期内存状态继续执行。
	conv, err := s.getConversationByNo(ctx, req.GetConversationNo())
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if conv.UserId != userID {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "conversation does not belong to current user")
	}
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	clientMessageNo := strings.TrimSpace(req.GetClientMessageNo())
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if clientMessageNo == "" {
		// 生成业务唯一编号，作为后续跨表关联、审计追踪和幂等定位的稳定主键。
		clientMessageNo = generateBizNo("CLI")
	}
	// 按客户端消息号查重，确保重复提交不会生成多条业务消息。
	existingMsg, _ := s.getMessageByClientNo(ctx, conv.ConversationNo, clientMessageNo)
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if existingMsg != nil {
		// 读取当前会话最近一次 Run，便于复用状态、判断 turn_no 或回传最新执行结果。
		latestRun, _ := s.getLatestRun(ctx, conv.ConversationNo)
		// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
		return &agentv1.SendAssistantMessageRes{
			Conversation:     toProtoConversation(conv),
			Run:              toProtoRun(latestRun),
			UserMessage:      toProtoMessage(existingMsg),
			IdempotentReplay: true,
			ReplyPayload:     parseReplyPayload(latestRun),
		}, nil
	}
	// 记录当前业务时间，确保状态流转、排序展示和审计字段使用同一时间基线。
	now := gtime.Now()
	// 先把显式 hidden_action 合并回 ext_json，保证消息落库后仍能被运行时和历史重放读取。
	mergedExtJSON := mergeHiddenActionIntoExtJSON(req.GetExtJson(), extractHiddenAction(req))
	// 生成业务唯一编号，作为后续跨表关联、审计追踪和幂等定位的稳定主键。
	messageNo := generateBizNo("AMSG")
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	_, err = dao.AgentMessage.Ctx(ctx).Data(do.AgentMessage{
		MessageNo:       messageNo,
		ConversationNo:  conv.ConversationNo,
		SenderTypeCode:  "BUYER",
		MessageTypeCode: normalizeMessageType(req.GetMessageTypeCode()),
		ClientMessageNo: clientMessageNo,
		ContentText:     strings.TrimSpace(req.GetContentText()),
		AssetIdsJson:    mustJSON(req.GetAssetIds()),
		ExtJson:         mergedExtJSON,
		SentAt:          now,
	}).Insert()
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "create buyer assistant message failed")
	}
	// 按客户端消息号查重，确保重复提交不会生成多条业务消息。
	userMsg, _ := s.getMessageByClientNo(ctx, conv.ConversationNo, clientMessageNo)
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if isQueueBlockedConversation(conv.ConversationStatusCode) {
		// 在排队期拦截新增用户消息，避免继续唤醒模型消耗 Token 并破坏交接摘要。
		_, queueErr := s.handleQueueBlockedMessage(ctx, conv, userMsg)
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		if queueErr != nil {
			// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
			return nil, queueErr
		}
		// 按会话号读取当前会话真相源，避免后续逻辑基于过期内存状态继续执行。
		updatedConv, _ := s.getConversationByNo(ctx, conv.ConversationNo)
		// 读取当前会话最近一次 Run，便于复用状态、判断 turn_no 或回传最新执行结果。
		latestRun, _ := s.getLatestRun(ctx, conv.ConversationNo)
		// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
		return &agentv1.SendAssistantMessageRes{
			Conversation: toProtoConversation(updatedConv),
			Run:          toProtoRun(latestRun),
			UserMessage:  toProtoMessage(userMsg),
			ReplyPayload: parseReplyPayload(latestRun),
		}, nil
	}
	// 读取当前会话最近一次 Run，便于复用状态、判断 turn_no 或回传最新执行结果。
	latestRun, _ := s.getLatestRun(ctx, conv.ConversationNo)
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	turnNo := uint64(1)
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if latestRun != nil {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		turnNo = latestRun.TurnNo + 1
	}
	// 生成业务唯一编号，作为后续跨表关联、审计追踪和幂等定位的稳定主键。
	runNo := generateBizNo("ARN")
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	_, err = dao.AgentRun.Ctx(ctx).Data(do.AgentRun{
		RunNo:              runNo,
		ConversationNo:     conv.ConversationNo,
		UserId:             conv.UserId,
		ShopNo:             conv.ShopNo,
		SubjectUserId:      conv.UserId,
		SubjectShopNo:      conv.ShopNo,
		TurnNo:             turnNo,
		AcceptedStatusCode: runtimeStatusPending,
		RunStatusCode:      runtimeStatusProcessing,
		CurrentNodeCode:    "LoadConversationContext",
		CheckpointVersion:  1,
		ToolResultStatus:   runtimeToolResultUnspecified,
		RiskDecisionCode:   "PASS",
		MergedMessageCount: 1,
		CreatedAt:          now,
		StartedAt:          now,
		UpdatedAt:          now,
	}).Insert()
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "create assistant run failed")
	}
	// 在数据库模型上追加精确过滤条件，确保只处理当前业务主体可见的数据。
	_, _ = dao.AgentConversation.Ctx(ctx).Where(dao.AgentConversation.Columns().ConversationNo, conv.ConversationNo).Data(do.AgentConversation{
		LastRunStatusCode: runtimeStatusProcessing,
		LastMessageAt:     now,
	}).Update()
	// 按 Run 编号加载执行记录，确保后续查询和状态流转都基于指定运行实例。
	run, _ := s.getRunByNo(ctx, runNo)
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	processedRun, _, err := s.processRun(ctx, conv, run, userMsg)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 按会话号读取当前会话真相源，避免后续逻辑基于过期内存状态继续执行。
	updatedConv, _ := s.getConversationByNo(ctx, conv.ConversationNo)
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &agentv1.SendAssistantMessageRes{
		Conversation: toProtoConversation(updatedConv),
		Run:          toProtoRun(processedRun),
		UserMessage:  toProtoMessage(userMsg),
		ReplyPayload: parseReplyPayload(processedRun),
	}, nil
}

func (s *sAgent) GetAssistantRunStatus(ctx context.Context, req *agentv1.GetAssistantRunStatusReq) (*agentv1.GetAssistantRunStatusRes, error) {
	// 先校验请求对象是否存在，避免后续读取空请求字段时触发空指针并污染业务语义。
	if req == nil || strings.TrimSpace(req.GetConversationNo()) == "" {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "conversation_no is required")
	}
	// 从上下文提取当前登录用户身份，保证后续查询和写入都严格绑定当前会话主体。
	userID, err := userIDFromContext(ctx)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 按会话号读取当前会话真相源，避免后续逻辑基于过期内存状态继续执行。
	conv, err := s.getConversationByNo(ctx, req.GetConversationNo())
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if conv.UserId != userID {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "conversation does not belong to current user")
	}
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var run *entity.AgentRun
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if strings.TrimSpace(req.GetRunNo()) != "" {
		// 按 Run 编号加载执行记录，确保后续查询和状态流转都基于指定运行实例。
		run, err = s.getRunByNo(ctx, req.GetRunNo())
	} else {
		// 读取当前会话最近一次 Run，便于复用状态、判断 turn_no 或回传最新执行结果。
		run, err = s.getLatestRun(ctx, conv.ConversationNo)
	}
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 读取当前 Run 或会话下最近一条客服回复，便于前端刷新最终展示内容。
	latestMsg, _ := s.getLatestAssistantMessage(ctx, conv.ConversationNo, run)
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return buildAssistantRunStatusResponse(conv, run, latestMsg), nil
}

func (s *sAgent) ListAssistantMessages(ctx context.Context, req *agentv1.ListAssistantMessagesReq) (*agentv1.ListAssistantMessagesRes, error) {
	// 先校验请求对象是否存在，避免后续读取空请求字段时触发空指针并污染业务语义。
	if req == nil || strings.TrimSpace(req.GetConversationNo()) == "" {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "conversation_no is required")
	}
	// 从上下文提取当前登录用户身份，保证后续查询和写入都严格绑定当前会话主体。
	userID, err := userIDFromContext(ctx)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 按会话号读取当前会话真相源，避免后续逻辑基于过期内存状态继续执行。
	conv, err := s.getConversationByNo(ctx, req.GetConversationNo())
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if conv.UserId != userID {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "conversation does not belong to current user")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	pageSize := normalizePageSize(req.GetPageSize())
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	cursorID, err := parseCursor(req.GetNextCursor())
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 在数据库模型上追加精确过滤条件，确保只处理当前业务主体可见的数据。
	model := dao.AgentMessage.Ctx(ctx).Where(dao.AgentMessage.Columns().ConversationNo, conv.ConversationNo).OrderDesc(dao.AgentMessage.Columns().Id).Limit(pageSize + 1)
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if cursorID > 0 {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		model = model.WhereLT(dao.AgentMessage.Columns().Id, cursorID)
	}
	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var rows []entity.AgentMessage
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if err = model.Scan(&rows); err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "query assistant messages failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	hasMore := false
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if len(rows) > pageSize {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		hasMore = true
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		rows = rows[:pageSize]
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	nextCursor := ""
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if hasMore && len(rows) > 0 {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		nextCursor = strconv.FormatUint(rows[len(rows)-1].Id, 10)
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	sort.SliceStable(rows, func(i, j int) bool { return rows[i].Id < rows[j].Id })
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	list := make([]*agentv1.AssistantMessage, 0, len(rows))
	// 遍历当前集合或循环条件，逐项推进本段业务处理并累计最终结果。
	for i := range rows {
		// 把内部实体或运行态结构转换成对外协议对象，保证对外契约稳定且隔离内部实现。
		list = append(list, toProtoMessage(&rows[i]))
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &agentv1.ListAssistantMessagesRes{List: list, NextCursor: nextCursor, HasMore: hasMore}, nil
}

func (s *sAgent) EscalateToHuman(ctx context.Context, req *agentv1.EscalateToHumanReq) (*agentv1.EscalateToHumanRes, error) {
	// 先校验请求对象是否存在，避免后续读取空请求字段时触发空指针并污染业务语义。
	if req == nil || strings.TrimSpace(req.GetConversationNo()) == "" {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "conversation_no is required")
	}
	// 从上下文提取当前登录用户身份，保证后续查询和写入都严格绑定当前会话主体。
	userID, err := userIDFromContext(ctx)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 按会话号读取当前会话真相源，避免后续逻辑基于过期内存状态继续执行。
	conv, err := s.getConversationByNo(ctx, req.GetConversationNo())
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if conv.UserId != userID {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "conversation does not belong to current user")
	}
	// 创建转人工工单，把当前会话从 AI 流程切换到人工接管链路。
	ticket, err := s.createEscalationTicket(ctx, conv, "", req.GetEscalationReasonCode(), req.GetRemark())
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 按会话号读取当前会话真相源，避免后续逻辑基于过期内存状态继续执行。
	updatedConv, _ := s.getConversationByNo(ctx, conv.ConversationNo)
	// 把内部实体或运行态结构转换成对外协议对象，保证对外契约稳定且隔离内部实现。
	return &agentv1.EscalateToHumanRes{
		Ticket:            toProtoTicket(ticket),
		Conversation:      toProtoConversation(updatedConv),
		HandoffReasonCode: ticket.EscalationReasonCode,
	}, nil
}

func (s *sAgent) SubmitAnswerFeedback(ctx context.Context, req *agentv1.SubmitAnswerFeedbackReq) (*agentv1.SubmitAnswerFeedbackRes, error) {
	// 先校验请求对象是否存在，避免后续读取空请求字段时触发空指针并污染业务语义。
	if req == nil || strings.TrimSpace(req.GetConversationNo()) == "" || strings.TrimSpace(req.GetRunNo()) == "" {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "conversation_no and run_no are required")
	}
	// 从上下文提取当前登录用户身份，保证后续查询和写入都严格绑定当前会话主体。
	userID, err := userIDFromContext(ctx)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 按会话号读取当前会话真相源，避免后续逻辑基于过期内存状态继续执行。
	conv, err := s.getConversationByNo(ctx, req.GetConversationNo())
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if conv.UserId != userID {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "conversation does not belong to current user")
	}
	// 记录当前业务时间，确保状态流转、排序展示和审计字段使用同一时间基线。
	now := gtime.Now()
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	_, err = dao.AgentFeedback.Ctx(ctx).Data(do.AgentFeedback{
		ConversationNo: conv.ConversationNo,
		RunNo:          req.GetRunNo(),
		UserId:         userID,
		ShopNo:         conv.ShopNo,
		FeedbackCode:   strings.ToUpper(strings.TrimSpace(req.GetFeedbackCode())),
		Resolved:       boolToInt(req.GetResolved()),
		Comment:        strings.TrimSpace(req.GetComment()),
		CreatedAt:      now,
		UpdatedAt:      now,
	}).Insert()
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "create feedback failed")
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &agentv1.SubmitAnswerFeedbackRes{
		ConversationNo: conv.ConversationNo,
		RunNo:          req.GetRunNo(),
		SubmittedAt:    timestamppb.New(now.Time),
	}, nil
}
