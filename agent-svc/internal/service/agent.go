package service

import (
	"context"

	agentv1 "github.com/TsingpekTao/shopa/agent-svc/api/v1"
)

type (
	IAgent interface {
		CreateOrGetAssistantConversation(ctx context.Context, req *agentv1.CreateOrGetAssistantConversationReq) (*agentv1.CreateOrGetAssistantConversationRes, error)
		SendAssistantMessage(ctx context.Context, req *agentv1.SendAssistantMessageReq) (*agentv1.SendAssistantMessageRes, error)
		GetAssistantRunStatus(ctx context.Context, req *agentv1.GetAssistantRunStatusReq) (*agentv1.GetAssistantRunStatusRes, error)
		ListAssistantMessages(ctx context.Context, req *agentv1.ListAssistantMessagesReq) (*agentv1.ListAssistantMessagesRes, error)
		EscalateToHuman(ctx context.Context, req *agentv1.EscalateToHumanReq) (*agentv1.EscalateToHumanRes, error)
		SubmitAnswerFeedback(ctx context.Context, req *agentv1.SubmitAnswerFeedbackReq) (*agentv1.SubmitAnswerFeedbackRes, error)
		CreateKnowledgeDoc(ctx context.Context, req *agentv1.CreateKnowledgeDocReq) (*agentv1.CreateKnowledgeDocRes, error)
		UpdateKnowledgeDoc(ctx context.Context, req *agentv1.UpdateKnowledgeDocReq) (*agentv1.UpdateKnowledgeDocRes, error)
		GetKnowledgeDocDetail(ctx context.Context, req *agentv1.GetKnowledgeDocDetailReq) (*agentv1.GetKnowledgeDocDetailRes, error)
		ListKnowledgeDocs(ctx context.Context, req *agentv1.ListKnowledgeDocsReq) (*agentv1.ListKnowledgeDocsRes, error)
		SubmitKnowledgeDoc(ctx context.Context, req *agentv1.SubmitKnowledgeDocReq) (*agentv1.SubmitKnowledgeDocRes, error)
		ApproveKnowledgeDoc(ctx context.Context, req *agentv1.ApproveKnowledgeDocReq) (*agentv1.ApproveKnowledgeDocRes, error)
		RejectKnowledgeDoc(ctx context.Context, req *agentv1.RejectKnowledgeDocReq) (*agentv1.RejectKnowledgeDocRes, error)
		PublishKnowledgeDoc(ctx context.Context, req *agentv1.PublishKnowledgeDocReq) (*agentv1.PublishKnowledgeDocRes, error)
		OfflineKnowledgeDoc(ctx context.Context, req *agentv1.OfflineKnowledgeDocReq) (*agentv1.OfflineKnowledgeDocRes, error)
		ReindexKnowledgeDoc(ctx context.Context, req *agentv1.ReindexKnowledgeDocReq) (*agentv1.ReindexKnowledgeDocRes, error)
		SearchKnowledgeChunks(ctx context.Context, req *agentv1.SearchKnowledgeChunksReq) (*agentv1.SearchKnowledgeChunksRes, error)
		ProcessUserTurn(ctx context.Context, req *agentv1.ProcessUserTurnReq) (*agentv1.ProcessUserTurnRes, error)
		AbortRun(ctx context.Context, req *agentv1.AbortRunReq) (*agentv1.AbortRunRes, error)
		RefreshSessionSummary(ctx context.Context, req *agentv1.RefreshSessionSummaryReq) (*agentv1.RefreshSessionSummaryRes, error)
		GenerateHandoffSummary(ctx context.Context, req *agentv1.GenerateHandoffSummaryReq) (*agentv1.GenerateHandoffSummaryRes, error)
		UpsertKnowledgeChunks(ctx context.Context, req *agentv1.UpsertKnowledgeChunksReq) (*agentv1.UpsertKnowledgeChunksRes, error)
		MarkKnowledgeSourceDeleted(ctx context.Context, req *agentv1.MarkKnowledgeSourceDeletedReq) (*agentv1.MarkKnowledgeSourceDeletedRes, error)
	}
)

var localAgent IAgent

func Agent() IAgent {
	if localAgent == nil {
		panic("implement not found for interface IAgent, forgot register?")
	}
	return localAgent
}

func RegisterAgent(i IAgent) {
	localAgent = i
}
