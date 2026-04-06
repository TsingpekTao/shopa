package agent

import (
	"context"
	"strings"

	agentv1 "github.com/TsingpekTao/shopa/agent-svc/api/v1"
	"github.com/TsingpekTao/shopa/agent-svc/internal/dao"
	"github.com/TsingpekTao/shopa/agent-svc/internal/model/do"
	"github.com/TsingpekTao/shopa/agent-svc/internal/model/entity"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

func (s *sAgent) ProcessUserTurn(ctx context.Context, req *agentv1.ProcessUserTurnReq) (*agentv1.ProcessUserTurnRes, error) {
	// 先校验请求对象是否存在，避免后续读取空请求字段时触发空指针并污染业务语义。
	if req == nil || strings.TrimSpace(req.GetConversationNo()) == "" {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "conversation_no is required")
	}
	// 按会话号读取当前会话真相源，避免后续逻辑基于过期内存状态继续执行。
	conv, err := s.getConversationByNo(ctx, req.GetConversationNo())
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}

	// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
	var run *entity.AgentRun
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if strings.TrimSpace(req.GetRunNo()) != "" {
		// 按 Run 编号加载执行记录，确保后续查询和状态流转都基于指定运行实例。
		run, err = s.getRunByNo(ctx, req.GetRunNo())
		// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
		if err != nil {
			// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
			return nil, err
		}
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		if run.ConversationNo != conv.ConversationNo {
			// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
			return nil, gerror.NewCode(gcode.CodeInvalidParameter, "run does not belong to conversation")
		}
	} else {
		// 读取当前会话最近一次 Run，便于复用状态、判断 turn_no 或回传最新执行结果。
		run, err = s.getLatestRun(ctx, conv.ConversationNo)
		// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
		if err != nil {
			// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
			return nil, err
		}
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if run == nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeNotFound, "run not found")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if strings.EqualFold(run.RunStatusCode, runtimeStatusSuccess) || strings.EqualFold(run.RunStatusCode, runtimeStatusEscalated) || strings.EqualFold(run.RunStatusCode, runtimeStatusAborted) {
		// 把内部实体或运行态结构转换成对外协议对象，保证对外契约稳定且隔离内部实现。
		return &agentv1.ProcessUserTurnRes{Run: toProtoRun(run)}, nil
	}

	// 读取买家最近一条消息，作为本轮 Agent 推理的最新用户输入。
	userMsg, err := s.getLatestBuyerMessage(ctx, conv.ConversationNo)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	processedRun, _, err := s.processRun(ctx, conv, run, userMsg)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 把内部实体或运行态结构转换成对外协议对象，保证对外契约稳定且隔离内部实现。
	return &agentv1.ProcessUserTurnRes{Run: toProtoRun(processedRun)}, nil
}

func (s *sAgent) AbortRun(ctx context.Context, req *agentv1.AbortRunReq) (*agentv1.AbortRunRes, error) {
	// 先校验请求对象是否存在，避免后续读取空请求字段时触发空指针并污染业务语义。
	if req == nil || strings.TrimSpace(req.GetConversationNo()) == "" || strings.TrimSpace(req.GetRunNo()) == "" {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "conversation_no and run_no are required")
	}
	// 按 Run 编号加载执行记录，确保后续查询和状态流转都基于指定运行实例。
	run, err := s.getRunByNo(ctx, req.GetRunNo())
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	if run.ConversationNo != strings.TrimSpace(req.GetConversationNo()) {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "run does not belong to conversation")
	}

	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	aborted := false
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if !strings.EqualFold(run.RunStatusCode, runtimeStatusSuccess) && !strings.EqualFold(run.RunStatusCode, runtimeStatusEscalated) && !strings.EqualFold(run.RunStatusCode, runtimeStatusAborted) {
		// 记录当前业务时间，确保状态流转、排序展示和审计字段使用同一时间基线。
		now := gtime.Now()
		// 在数据库模型上追加精确过滤条件，确保只处理当前业务主体可见的数据。
		_, err = dao.AgentRun.Ctx(ctx).Where(dao.AgentRun.Columns().RunNo, run.RunNo).Data(do.AgentRun{
			RunStatusCode:       runtimeStatusAborted,
			CurrentNodeCode:     "AbortRun",
			ReplyInterrupted:    1,
			InterruptReasonCode: strings.ToUpper(strings.TrimSpace(req.GetReasonCode())),
			ToolResultStatus:    normalizeToolResultStatus(run.ToolResultStatus),
			FinishedAt:          now,
			UpdatedAt:           now,
		}).Update()
		// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
		if err != nil {
			// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
			return nil, gerror.Wrap(err, "abort run failed")
		}
		// 在数据库模型上追加精确过滤条件，确保只处理当前业务主体可见的数据。
		_, _ = dao.AgentConversation.Ctx(ctx).Where(dao.AgentConversation.Columns().ConversationNo, run.ConversationNo).Data(do.AgentConversation{
			LastRunStatusCode: runtimeStatusAborted,
			UpdatedAt:         now,
		}).Update()
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		aborted = true
	}

	// 按 Run 编号加载执行记录，确保后续查询和状态流转都基于指定运行实例。
	updatedRun, err := s.getRunByNo(ctx, req.GetRunNo())
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 把内部实体或运行态结构转换成对外协议对象，保证对外契约稳定且隔离内部实现。
	return &agentv1.AbortRunRes{Aborted: aborted, Run: toProtoRun(updatedRun)}, nil
}

func (s *sAgent) RefreshSessionSummary(ctx context.Context, req *agentv1.RefreshSessionSummaryReq) (*agentv1.RefreshSessionSummaryRes, error) {
	// 先校验请求对象是否存在，避免后续读取空请求字段时触发空指针并污染业务语义。
	if req == nil || strings.TrimSpace(req.GetConversationNo()) == "" {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "conversation_no is required")
	}
	// 按会话号读取当前会话真相源，避免后续逻辑基于过期内存状态继续执行。
	conv, err := s.getConversationByNo(ctx, req.GetConversationNo())
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 拉取最近消息窗口，为总结、推理和转人工摘要提供最新上下文。
	recentRows, err := s.listRecentMessages(ctx, conv.ConversationNo, maxRecentMessages)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	summary, err := s.runner.Summarize(ctx, conv.SessionSummary, rowsToMessages(recentRows))
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "refresh session summary failed")
	}

	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	summaryUpdated := strings.TrimSpace(summary) != strings.TrimSpace(conv.SessionSummary)
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if summaryUpdated {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		lastMessageNo := ""
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		if len(recentRows) > 0 {
			// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
			lastMessageNo = recentRows[len(recentRows)-1].MessageNo
		}
		// 记录当前业务时间，确保状态流转、排序展示和审计字段使用同一时间基线。
		now := gtime.Now()
		// 在数据库模型上追加精确过滤条件，确保只处理当前业务主体可见的数据。
		_, err = dao.AgentConversation.Ctx(ctx).Where(dao.AgentConversation.Columns().ConversationNo, conv.ConversationNo).Data(do.AgentConversation{
			SessionSummary:          summary,
			SessionSummaryVersion:   conv.SessionSummaryVersion + 1,
			LastSummarizedMessageNo: lastMessageNo,
			UpdatedAt:               now,
		}).Update()
		// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
		if err != nil {
			// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
			return nil, gerror.Wrap(err, "update session summary failed")
		}
	}

	// 按会话号读取当前会话真相源，避免后续逻辑基于过期内存状态继续执行。
	updatedConv, err := s.getConversationByNo(ctx, conv.ConversationNo)
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 把内部实体或运行态结构转换成对外协议对象，保证对外契约稳定且隔离内部实现。
	return &agentv1.RefreshSessionSummaryRes{Conversation: toProtoConversation(updatedConv), SummaryUpdated: summaryUpdated}, nil
}

func (s *sAgent) GenerateHandoffSummary(ctx context.Context, req *agentv1.GenerateHandoffSummaryReq) (*agentv1.GenerateHandoffSummaryRes, error) {
	// 先校验请求对象是否存在，避免后续读取空请求字段时触发空指针并污染业务语义。
	if req == nil || strings.TrimSpace(req.GetConversationNo()) == "" {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "conversation_no is required")
	}
	// 按会话号读取当前会话真相源，避免后续逻辑基于过期内存状态继续执行。
	conv, err := s.getConversationByNo(ctx, req.GetConversationNo())
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 创建转人工工单，把当前会话从 AI 流程切换到人工接管链路。
	ticket, err := s.createEscalationTicket(ctx, conv, req.GetRunNo(), req.GetEscalationReasonCode(), "")
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}
	// 把内部实体或运行态结构转换成对外协议对象，保证对外契约稳定且隔离内部实现。
	return &agentv1.GenerateHandoffSummaryRes{Ticket: toProtoTicket(ticket)}, nil
}

func (s *sAgent) UpsertKnowledgeChunks(ctx context.Context, req *agentv1.UpsertKnowledgeChunksReq) (*agentv1.UpsertKnowledgeChunksRes, error) {
	// 先校验请求对象是否存在，避免后续读取空请求字段时触发空指针并污染业务语义。
	if req == nil || strings.TrimSpace(req.GetKnowledgeDocNo()) == "" {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "knowledge_doc_no is required")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	doc, err := s.getKnowledgeDocByNo(ctx, req.GetKnowledgeDocNo())
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, err
	}

	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	sourceType := strings.ToUpper(strings.TrimSpace(req.GetSourceTypeCode()))
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if sourceType == "" {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		sourceType = doc.SourceTypeCode
	}
	// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
	sourceID := strings.TrimSpace(req.GetSourceId())
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if sourceID == "" {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		sourceID = doc.SourceId
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	sourceVersion := req.GetSourceVersion()
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if sourceVersion == 0 {
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		sourceVersion = doc.SourceVersion
	}

	// 记录当前业务时间，确保状态流转、排序展示和审计字段使用同一时间基线。
	now := gtime.Now()
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	_, err = dao.AgentKnowledgeChunk.Ctx(ctx).
		Where(dao.AgentKnowledgeChunk.Columns().SourceTypeCode, sourceType).
		Where(dao.AgentKnowledgeChunk.Columns().SourceId, sourceID).
		WhereLT(dao.AgentKnowledgeChunk.Columns().SourceVersion, sourceVersion).
		Data(do.AgentKnowledgeChunk{IsDeleted: 1, IndexStatusCode: "STALE", UpdatedAt: now}).
		Update()
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "mark old knowledge chunks deleted failed")
	}

	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	_, err = dao.AgentKnowledgeChunk.Ctx(ctx).
		Where(dao.AgentKnowledgeChunk.Columns().KnowledgeDocNo, doc.KnowledgeDocNo).
		Data(do.AgentKnowledgeChunk{IsDeleted: 1, IndexStatusCode: "STALE", UpdatedAt: now}).
		Update()
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "mark current knowledge chunks stale failed")
	}

	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	upserted := uint64(0)
	// 遍历当前集合或循环条件，逐项推进本段业务处理并累计最终结果。
	for _, item := range req.GetItems() {
		// 先清洗字符串输入中的空白字符，避免参数脏值影响后续状态判断、查询或落库。
		if item == nil || strings.TrimSpace(item.GetChunkText()) == "" {
			// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
			continue
		}
		// 先集中声明这一段流程会复用的变量，便于后续按顺序填充和统一收口。
		var existing entity.AgentKnowledgeChunk
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		err = dao.AgentKnowledgeChunk.Ctx(ctx).
			Where(dao.AgentKnowledgeChunk.Columns().KnowledgeDocNo, doc.KnowledgeDocNo).
			Where(dao.AgentKnowledgeChunk.Columns().ChunkIndex, item.GetChunkIndex()).
			Scan(&existing)
		// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
		if err != nil {
			// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
			return nil, gerror.Wrap(err, "query knowledge chunk failed")
		}
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		data := do.AgentKnowledgeChunk{
			KnowledgeDocNo:   doc.KnowledgeDocNo,
			SourceTypeCode:   sourceType,
			SourceId:         sourceID,
			SourceVersion:    sourceVersion,
			ChunkIndex:       item.GetChunkIndex(),
			ChunkText:        strings.TrimSpace(item.GetChunkText()),
			ChunkTextPreview: strings.TrimSpace(item.GetChunkTextPreview()),
			VectorDocumentId: strings.TrimSpace(item.GetVectorDocumentId()),
			MetadataJson:     normalizeJSON(item.GetMetadataJson()),
			IsDeleted:        0,
			IndexStatusCode:  "READY",
			UpdatedAt:        now,
		}
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		if existing.Id == 0 {
			// 生成业务唯一编号，作为后续跨表关联、审计追踪和幂等定位的稳定主键。
			data.ChunkNo = generateBizNo("AKC")
			// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
			data.CreatedAt = now
			// 把当前业务实体持久化到数据库，先落真相源再继续后续流程。
			_, err = dao.AgentKnowledgeChunk.Ctx(ctx).Data(data).Insert()
		} else {
			// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
			_, err = dao.AgentKnowledgeChunk.Ctx(ctx).
				Where(dao.AgentKnowledgeChunk.Columns().Id, existing.Id).
				Data(data).
				Update()
		}
		// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
		if err != nil {
			// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
			return nil, gerror.Wrap(err, "upsert knowledge chunk failed")
		}
		// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
		upserted++
	}

	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if sourceVersion != doc.SourceVersion {
		// 在数据库模型上追加精确过滤条件，确保只处理当前业务主体可见的数据。
		_, _ = dao.AgentKnowledgeDoc.Ctx(ctx).Where(dao.AgentKnowledgeDoc.Columns().KnowledgeDocNo, doc.KnowledgeDocNo).Data(do.AgentKnowledgeDoc{
			SourceVersion: sourceVersion,
			Version:       doc.Version + 1,
			UpdatedAt:     now,
		}).Update()
	}
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &agentv1.UpsertKnowledgeChunksRes{UpsertedCount: upserted}, nil
}

func (s *sAgent) MarkKnowledgeSourceDeleted(ctx context.Context, req *agentv1.MarkKnowledgeSourceDeletedReq) (*agentv1.MarkKnowledgeSourceDeletedRes, error) {
	// 先校验请求对象是否存在，避免后续读取空请求字段时触发空指针并污染业务语义。
	if req == nil || strings.TrimSpace(req.GetSourceTypeCode()) == "" || strings.TrimSpace(req.GetSourceId()) == "" {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "source_type_code and source_id are required")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	model := dao.AgentKnowledgeChunk.Ctx(ctx).
		Where(dao.AgentKnowledgeChunk.Columns().SourceTypeCode, strings.ToUpper(strings.TrimSpace(req.GetSourceTypeCode()))).
		Where(dao.AgentKnowledgeChunk.Columns().SourceId, strings.TrimSpace(req.GetSourceId()))
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	if req.GetSourceVersion() > 0 {
		// 在数据库模型上追加精确过滤条件，确保只处理当前业务主体可见的数据。
		model = model.Where(dao.AgentKnowledgeChunk.Columns().SourceVersion, req.GetSourceVersion())
	}
	// 记录当前业务时间，确保状态流转、排序展示和审计字段使用同一时间基线。
	result, err := model.Data(do.AgentKnowledgeChunk{IsDeleted: 1, IndexStatusCode: "DELETED", UpdatedAt: gtime.Now()}).Update()
	// 如果上一步已经出现错误，这里立即中断并向上返回，避免带着脏状态继续推进链路。
	if err != nil {
		// 把当前错误继续向上返回，让调用方通过统一错误链路感知失败原因。
		return nil, gerror.Wrap(err, "mark knowledge source deleted failed")
	}
	// 执行当前业务语句，把本步骤产出的状态或数据继续传递给后续流程。
	affected, _ := result.RowsAffected()
	// 组装当前步骤的响应结构，把内部结果转换成稳定的对外返回格式。
	return &agentv1.MarkKnowledgeSourceDeletedRes{AffectedRows: uint64(affected)}, nil
}
