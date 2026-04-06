package agent

import (
	"context"

	httpv1 "github.com/TsingpekTao/shopa/agent-svc/api/agent/v1"
	pb "github.com/TsingpekTao/shopa/agent-svc/api/v1"
	fieldmaskpb "google.golang.org/protobuf/types/known/fieldmaskpb"
)

func (c *ControllerV1) CreateOrGetAssistantConversation(ctx context.Context, req *httpv1.CreateOrGetAssistantConversationReq) (*httpv1.CreateOrGetAssistantConversationRes, error) {
	return c.svc.CreateOrGetAssistantConversation(withRequestMetadata(ctx), &req.CreateOrGetAssistantConversationReq)
}

func (c *ControllerV1) SendAssistantMessage(ctx context.Context, req *httpv1.SendAssistantMessageReq) (*httpv1.SendAssistantMessageRes, error) {
	return c.svc.SendAssistantMessage(withRequestMetadata(ctx), &req.SendAssistantMessageReq)
}

func (c *ControllerV1) GetAssistantRunStatus(ctx context.Context, req *httpv1.GetAssistantRunStatusReq) (*httpv1.GetAssistantRunStatusRes, error) {
	return c.svc.GetAssistantRunStatus(withRequestMetadata(ctx), &pb.GetAssistantRunStatusReq{ConversationNo: req.ConversationNo, RunNo: req.RunNo})
}

func (c *ControllerV1) ListAssistantMessages(ctx context.Context, req *httpv1.ListAssistantMessagesReq) (*httpv1.ListAssistantMessagesRes, error) {
	return c.svc.ListAssistantMessages(withRequestMetadata(ctx), &pb.ListAssistantMessagesReq{ConversationNo: req.ConversationNo, PageSize: req.PageSize, NextCursor: req.NextCursor})
}

func (c *ControllerV1) EscalateToHuman(ctx context.Context, req *httpv1.EscalateToHumanReq) (*httpv1.EscalateToHumanRes, error) {
	return c.svc.EscalateToHuman(withRequestMetadata(ctx), &pb.EscalateToHumanReq{ConversationNo: req.ConversationNo, EscalationReasonCode: req.EscalationReasonCode, Remark: req.Remark})
}

func (c *ControllerV1) SubmitAnswerFeedback(ctx context.Context, req *httpv1.SubmitAnswerFeedbackReq) (*httpv1.SubmitAnswerFeedbackRes, error) {
	return c.svc.SubmitAnswerFeedback(withRequestMetadata(ctx), &req.SubmitAnswerFeedbackReq)
}

func (c *ControllerV1) CreateKnowledgeDoc(ctx context.Context, req *httpv1.CreateKnowledgeDocReq) (*httpv1.CreateKnowledgeDocRes, error) {
	return c.svc.CreateKnowledgeDoc(withRequestMetadata(ctx), &req.CreateKnowledgeDocReq)
}

func (c *ControllerV1) UpdateKnowledgeDoc(ctx context.Context, req *httpv1.UpdateKnowledgeDocReq) (*httpv1.UpdateKnowledgeDocRes, error) {
	return c.svc.UpdateKnowledgeDoc(withRequestMetadata(ctx), &pb.UpdateKnowledgeDocReq{
		KnowledgeDocNo:  req.KnowledgeDocNo,
		Doc:             req.Doc,
		UpdateMask:      &fieldmaskpb.FieldMask{},
		ExpectedVersion: req.ExpectedVersion,
	})
}

func (c *ControllerV1) GetKnowledgeDocDetail(ctx context.Context, req *httpv1.GetKnowledgeDocDetailReq) (*httpv1.GetKnowledgeDocDetailRes, error) {
	return c.svc.GetKnowledgeDocDetail(withRequestMetadata(ctx), &pb.GetKnowledgeDocDetailReq{KnowledgeDocNo: req.KnowledgeDocNo})
}

func (c *ControllerV1) ListKnowledgeDocs(ctx context.Context, req *httpv1.ListKnowledgeDocsReq) (*httpv1.ListKnowledgeDocsRes, error) {
	return c.svc.ListKnowledgeDocs(withRequestMetadata(ctx), &req.ListKnowledgeDocsReq)
}

func (c *ControllerV1) SubmitKnowledgeDoc(ctx context.Context, req *httpv1.SubmitKnowledgeDocReq) (*httpv1.SubmitKnowledgeDocRes, error) {
	return c.svc.SubmitKnowledgeDoc(withRequestMetadata(ctx), &pb.SubmitKnowledgeDocReq{KnowledgeDocNo: req.KnowledgeDocNo, ExpectedVersion: req.ExpectedVersion})
}

func (c *ControllerV1) ApproveKnowledgeDoc(ctx context.Context, req *httpv1.ApproveKnowledgeDocReq) (*httpv1.ApproveKnowledgeDocRes, error) {
	return c.svc.ApproveKnowledgeDoc(withRequestMetadata(ctx), &pb.ApproveKnowledgeDocReq{KnowledgeDocNo: req.KnowledgeDocNo, ExpectedVersion: req.ExpectedVersion, Remark: req.Remark})
}

func (c *ControllerV1) RejectKnowledgeDoc(ctx context.Context, req *httpv1.RejectKnowledgeDocReq) (*httpv1.RejectKnowledgeDocRes, error) {
	return c.svc.RejectKnowledgeDoc(withRequestMetadata(ctx), &pb.RejectKnowledgeDocReq{KnowledgeDocNo: req.KnowledgeDocNo, ExpectedVersion: req.ExpectedVersion, RejectReasonCode: req.RejectReasonCode, RejectComment: req.RejectComment})
}

func (c *ControllerV1) PublishKnowledgeDoc(ctx context.Context, req *httpv1.PublishKnowledgeDocReq) (*httpv1.PublishKnowledgeDocRes, error) {
	return c.svc.PublishKnowledgeDoc(withRequestMetadata(ctx), &pb.PublishKnowledgeDocReq{KnowledgeDocNo: req.KnowledgeDocNo, ExpectedVersion: req.ExpectedVersion})
}

func (c *ControllerV1) OfflineKnowledgeDoc(ctx context.Context, req *httpv1.OfflineKnowledgeDocReq) (*httpv1.OfflineKnowledgeDocRes, error) {
	return c.svc.OfflineKnowledgeDoc(withRequestMetadata(ctx), &pb.OfflineKnowledgeDocReq{KnowledgeDocNo: req.KnowledgeDocNo, ExpectedVersion: req.ExpectedVersion, ReasonCode: req.ReasonCode})
}

func (c *ControllerV1) ReindexKnowledgeDoc(ctx context.Context, req *httpv1.ReindexKnowledgeDocReq) (*httpv1.ReindexKnowledgeDocRes, error) {
	return c.svc.ReindexKnowledgeDoc(withRequestMetadata(ctx), &pb.ReindexKnowledgeDocReq{KnowledgeDocNo: req.KnowledgeDocNo, ExpectedVersion: req.ExpectedVersion, ReasonCode: req.ReasonCode})
}

func (c *ControllerV1) SearchKnowledgeChunks(ctx context.Context, req *httpv1.SearchKnowledgeChunksReq) (*httpv1.SearchKnowledgeChunksRes, error) {
	return c.svc.SearchKnowledgeChunks(withRequestMetadata(ctx), &req.SearchKnowledgeChunksReq)
}

func (c *ControllerV1) ProcessUserTurn(ctx context.Context, req *httpv1.ProcessUserTurnReq) (*httpv1.ProcessUserTurnRes, error) {
	return c.svc.ProcessUserTurn(withRequestMetadata(ctx), &req.ProcessUserTurnReq)
}

func (c *ControllerV1) AbortRun(ctx context.Context, req *httpv1.AbortRunReq) (*httpv1.AbortRunRes, error) {
	return c.svc.AbortRun(withRequestMetadata(ctx), &req.AbortRunReq)
}

func (c *ControllerV1) RefreshSessionSummary(ctx context.Context, req *httpv1.RefreshSessionSummaryReq) (*httpv1.RefreshSessionSummaryRes, error) {
	return c.svc.RefreshSessionSummary(withRequestMetadata(ctx), &req.RefreshSessionSummaryReq)
}

func (c *ControllerV1) GenerateHandoffSummary(ctx context.Context, req *httpv1.GenerateHandoffSummaryReq) (*httpv1.GenerateHandoffSummaryRes, error) {
	return c.svc.GenerateHandoffSummary(withRequestMetadata(ctx), &req.GenerateHandoffSummaryReq)
}

func (c *ControllerV1) UpsertKnowledgeChunks(ctx context.Context, req *httpv1.UpsertKnowledgeChunksReq) (*httpv1.UpsertKnowledgeChunksRes, error) {
	return c.svc.UpsertKnowledgeChunks(withRequestMetadata(ctx), &req.UpsertKnowledgeChunksReq)
}

func (c *ControllerV1) MarkKnowledgeSourceDeleted(ctx context.Context, req *httpv1.MarkKnowledgeSourceDeletedReq) (*httpv1.MarkKnowledgeSourceDeletedRes, error) {
	return c.svc.MarkKnowledgeSourceDeleted(withRequestMetadata(ctx), &req.MarkKnowledgeSourceDeletedReq)
}
