package agent

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/agent-svc/api/agent/v1"
)

// IAgentV1 定义 agent-svc 的 HTTP API。
type IAgentV1 interface {
	// 买家智能客服接口。
	CreateOrGetAssistantConversation(ctx context.Context, req *v1.CreateOrGetAssistantConversationReq) (res *v1.CreateOrGetAssistantConversationRes, err error)
	SendAssistantMessage(ctx context.Context, req *v1.SendAssistantMessageReq) (res *v1.SendAssistantMessageRes, err error)
	GetAssistantRunStatus(ctx context.Context, req *v1.GetAssistantRunStatusReq) (res *v1.GetAssistantRunStatusRes, err error)
	ListAssistantMessages(ctx context.Context, req *v1.ListAssistantMessagesReq) (res *v1.ListAssistantMessagesRes, err error)
	EscalateToHuman(ctx context.Context, req *v1.EscalateToHumanReq) (res *v1.EscalateToHumanRes, err error)
	SubmitAnswerFeedback(ctx context.Context, req *v1.SubmitAnswerFeedbackReq) (res *v1.SubmitAnswerFeedbackRes, err error)

	// 知识库管理接口。
	CreateKnowledgeDoc(ctx context.Context, req *v1.CreateKnowledgeDocReq) (res *v1.CreateKnowledgeDocRes, err error)
	UpdateKnowledgeDoc(ctx context.Context, req *v1.UpdateKnowledgeDocReq) (res *v1.UpdateKnowledgeDocRes, err error)
	GetKnowledgeDocDetail(ctx context.Context, req *v1.GetKnowledgeDocDetailReq) (res *v1.GetKnowledgeDocDetailRes, err error)
	ListKnowledgeDocs(ctx context.Context, req *v1.ListKnowledgeDocsReq) (res *v1.ListKnowledgeDocsRes, err error)
	SubmitKnowledgeDoc(ctx context.Context, req *v1.SubmitKnowledgeDocReq) (res *v1.SubmitKnowledgeDocRes, err error)
	ApproveKnowledgeDoc(ctx context.Context, req *v1.ApproveKnowledgeDocReq) (res *v1.ApproveKnowledgeDocRes, err error)
	RejectKnowledgeDoc(ctx context.Context, req *v1.RejectKnowledgeDocReq) (res *v1.RejectKnowledgeDocRes, err error)
	PublishKnowledgeDoc(ctx context.Context, req *v1.PublishKnowledgeDocReq) (res *v1.PublishKnowledgeDocRes, err error)
	OfflineKnowledgeDoc(ctx context.Context, req *v1.OfflineKnowledgeDocReq) (res *v1.OfflineKnowledgeDocRes, err error)
	ReindexKnowledgeDoc(ctx context.Context, req *v1.ReindexKnowledgeDocReq) (res *v1.ReindexKnowledgeDocRes, err error)
	SearchKnowledgeChunks(ctx context.Context, req *v1.SearchKnowledgeChunksReq) (res *v1.SearchKnowledgeChunksRes, err error)

	// 内部任务接口。
	ProcessUserTurn(ctx context.Context, req *v1.ProcessUserTurnReq) (res *v1.ProcessUserTurnRes, err error)
	AbortRun(ctx context.Context, req *v1.AbortRunReq) (res *v1.AbortRunRes, err error)
	RefreshSessionSummary(ctx context.Context, req *v1.RefreshSessionSummaryReq) (res *v1.RefreshSessionSummaryRes, err error)
	GenerateHandoffSummary(ctx context.Context, req *v1.GenerateHandoffSummaryReq) (res *v1.GenerateHandoffSummaryRes, err error)
	UpsertKnowledgeChunks(ctx context.Context, req *v1.UpsertKnowledgeChunksReq) (res *v1.UpsertKnowledgeChunksRes, err error)
	MarkKnowledgeSourceDeleted(ctx context.Context, req *v1.MarkKnowledgeSourceDeletedReq) (res *v1.MarkKnowledgeSourceDeletedRes, err error)
}
