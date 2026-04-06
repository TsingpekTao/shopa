package api

import (
	"context"

	v1 "github.com/TsingpekTao/shopa/agent-svc/api/v1"
	"github.com/TsingpekTao/shopa/agent-svc/internal/service"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedBuyerAgentServiceServer
	v1.UnimplementedKnowledgeAdminServiceServer
	v1.UnimplementedAgentInternalServiceServer
}

func Register(s *grpcx.GrpcServer) {
	ctrl := &Controller{}
	v1.RegisterBuyerAgentServiceServer(s.Server, ctrl)
	v1.RegisterKnowledgeAdminServiceServer(s.Server, ctrl)
	v1.RegisterAgentInternalServiceServer(s.Server, ctrl)
}

func (*Controller) CreateOrGetAssistantConversation(ctx context.Context, req *v1.CreateOrGetAssistantConversationReq) (*v1.CreateOrGetAssistantConversationRes, error) {
	return service.Agent().CreateOrGetAssistantConversation(ctx, req)
}

func (*Controller) SendAssistantMessage(ctx context.Context, req *v1.SendAssistantMessageReq) (*v1.SendAssistantMessageRes, error) {
	return service.Agent().SendAssistantMessage(ctx, req)
}

func (*Controller) GetAssistantRunStatus(ctx context.Context, req *v1.GetAssistantRunStatusReq) (*v1.GetAssistantRunStatusRes, error) {
	return service.Agent().GetAssistantRunStatus(ctx, req)
}

func (*Controller) ListAssistantMessages(ctx context.Context, req *v1.ListAssistantMessagesReq) (*v1.ListAssistantMessagesRes, error) {
	return service.Agent().ListAssistantMessages(ctx, req)
}

func (*Controller) EscalateToHuman(ctx context.Context, req *v1.EscalateToHumanReq) (*v1.EscalateToHumanRes, error) {
	return service.Agent().EscalateToHuman(ctx, req)
}

func (*Controller) SubmitAnswerFeedback(ctx context.Context, req *v1.SubmitAnswerFeedbackReq) (*v1.SubmitAnswerFeedbackRes, error) {
	return service.Agent().SubmitAnswerFeedback(ctx, req)
}

func (*Controller) CreateKnowledgeDoc(ctx context.Context, req *v1.CreateKnowledgeDocReq) (*v1.CreateKnowledgeDocRes, error) {
	return service.Agent().CreateKnowledgeDoc(ctx, req)
}

func (*Controller) UpdateKnowledgeDoc(ctx context.Context, req *v1.UpdateKnowledgeDocReq) (*v1.UpdateKnowledgeDocRes, error) {
	return service.Agent().UpdateKnowledgeDoc(ctx, req)
}

func (*Controller) GetKnowledgeDocDetail(ctx context.Context, req *v1.GetKnowledgeDocDetailReq) (*v1.GetKnowledgeDocDetailRes, error) {
	return service.Agent().GetKnowledgeDocDetail(ctx, req)
}

func (*Controller) ListKnowledgeDocs(ctx context.Context, req *v1.ListKnowledgeDocsReq) (*v1.ListKnowledgeDocsRes, error) {
	return service.Agent().ListKnowledgeDocs(ctx, req)
}

func (*Controller) SubmitKnowledgeDoc(ctx context.Context, req *v1.SubmitKnowledgeDocReq) (*v1.SubmitKnowledgeDocRes, error) {
	return service.Agent().SubmitKnowledgeDoc(ctx, req)
}

func (*Controller) ApproveKnowledgeDoc(ctx context.Context, req *v1.ApproveKnowledgeDocReq) (*v1.ApproveKnowledgeDocRes, error) {
	return service.Agent().ApproveKnowledgeDoc(ctx, req)
}

func (*Controller) RejectKnowledgeDoc(ctx context.Context, req *v1.RejectKnowledgeDocReq) (*v1.RejectKnowledgeDocRes, error) {
	return service.Agent().RejectKnowledgeDoc(ctx, req)
}

func (*Controller) PublishKnowledgeDoc(ctx context.Context, req *v1.PublishKnowledgeDocReq) (*v1.PublishKnowledgeDocRes, error) {
	return service.Agent().PublishKnowledgeDoc(ctx, req)
}

func (*Controller) OfflineKnowledgeDoc(ctx context.Context, req *v1.OfflineKnowledgeDocReq) (*v1.OfflineKnowledgeDocRes, error) {
	return service.Agent().OfflineKnowledgeDoc(ctx, req)
}

func (*Controller) ReindexKnowledgeDoc(ctx context.Context, req *v1.ReindexKnowledgeDocReq) (*v1.ReindexKnowledgeDocRes, error) {
	return service.Agent().ReindexKnowledgeDoc(ctx, req)
}

func (*Controller) SearchKnowledgeChunks(ctx context.Context, req *v1.SearchKnowledgeChunksReq) (*v1.SearchKnowledgeChunksRes, error) {
	return service.Agent().SearchKnowledgeChunks(ctx, req)
}

func (*Controller) ProcessUserTurn(ctx context.Context, req *v1.ProcessUserTurnReq) (*v1.ProcessUserTurnRes, error) {
	return service.Agent().ProcessUserTurn(ctx, req)
}

func (*Controller) AbortRun(ctx context.Context, req *v1.AbortRunReq) (*v1.AbortRunRes, error) {
	return service.Agent().AbortRun(ctx, req)
}

func (*Controller) RefreshSessionSummary(ctx context.Context, req *v1.RefreshSessionSummaryReq) (*v1.RefreshSessionSummaryRes, error) {
	return service.Agent().RefreshSessionSummary(ctx, req)
}

func (*Controller) GenerateHandoffSummary(ctx context.Context, req *v1.GenerateHandoffSummaryReq) (*v1.GenerateHandoffSummaryRes, error) {
	return service.Agent().GenerateHandoffSummary(ctx, req)
}

func (*Controller) UpsertKnowledgeChunks(ctx context.Context, req *v1.UpsertKnowledgeChunksReq) (*v1.UpsertKnowledgeChunksRes, error) {
	return service.Agent().UpsertKnowledgeChunks(ctx, req)
}

func (*Controller) MarkKnowledgeSourceDeleted(ctx context.Context, req *v1.MarkKnowledgeSourceDeletedReq) (*v1.MarkKnowledgeSourceDeletedRes, error) {
	return service.Agent().MarkKnowledgeSourceDeleted(ctx, req)
}
