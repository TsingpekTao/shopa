package v1

import (
	pb "github.com/TsingpekTao/shopa/agent-svc/api/v1"
	"github.com/gogf/gf/v2/frame/g"
)

type CreateOrGetAssistantConversationReq struct {
	g.Meta `path:"/v1/agent/buyer/conversations:get-or-create" method:"post" tags:"Agent-Buyer" summary:"Create or get assistant conversation"`
	pb.CreateOrGetAssistantConversationReq
}

type CreateOrGetAssistantConversationRes = pb.CreateOrGetAssistantConversationRes

type SendAssistantMessageReq struct {
	g.Meta `path:"/v1/agent/buyer/messages:send" method:"post" tags:"Agent-Buyer" summary:"Send message to assistant"`
	pb.SendAssistantMessageReq
}

type SendAssistantMessageRes = pb.SendAssistantMessageRes

type GetAssistantRunStatusReq struct {
	g.Meta         `path:"/v1/agent/buyer/conversations/{conversation_no}/run-status" method:"get" tags:"Agent-Buyer" summary:"Get assistant run status"`
	ConversationNo string `json:"conversation_no" in:"path" v:"required#conversation_no is required"`
	RunNo          string `json:"run_no" in:"query"`
}

type GetAssistantRunStatusRes = pb.GetAssistantRunStatusRes

type ListAssistantMessagesReq struct {
	g.Meta         `path:"/v1/agent/buyer/conversations/{conversation_no}/messages" method:"get" tags:"Agent-Buyer" summary:"List assistant messages"`
	ConversationNo string `json:"conversation_no" in:"path" v:"required#conversation_no is required"`
	PageSize       int32  `json:"page_size" in:"query"`
	NextCursor     string `json:"next_cursor" in:"query"`
}

type ListAssistantMessagesRes = pb.ListAssistantMessagesRes

type EscalateToHumanReq struct {
	g.Meta               `path:"/v1/agent/buyer/conversations/{conversation_no}:escalate" method:"post" tags:"Agent-Buyer" summary:"Escalate conversation to human"`
	ConversationNo       string `json:"conversation_no" in:"path" v:"required#conversation_no is required"`
	EscalationReasonCode string `json:"escalation_reason_code"`
	Remark               string `json:"remark"`
}

type EscalateToHumanRes = pb.EscalateToHumanRes

type SubmitAnswerFeedbackReq struct {
	g.Meta `path:"/v1/agent/buyer/feedback:submit" method:"post" tags:"Agent-Buyer" summary:"Submit assistant answer feedback"`
	pb.SubmitAnswerFeedbackReq
}

type SubmitAnswerFeedbackRes = pb.SubmitAnswerFeedbackRes

type CreateKnowledgeDocReq struct {
	g.Meta `path:"/v1/agent/admin/knowledge-docs" method:"post" tags:"Agent-Knowledge" summary:"Create knowledge document"`
	pb.CreateKnowledgeDocReq
}

type CreateKnowledgeDocRes = pb.CreateKnowledgeDocRes

type UpdateKnowledgeDocReq struct {
	g.Meta          `path:"/v1/agent/admin/knowledge-docs/{knowledge_doc_no}" method:"patch" tags:"Agent-Knowledge" summary:"Update knowledge document"`
	KnowledgeDocNo  string                `json:"knowledge_doc_no" in:"path" v:"required#knowledge_doc_no is required"`
	Doc             *pb.KnowledgeDocDraft `json:"doc"`
	UpdateMask      any                   `json:"-"`
	ExpectedVersion uint64                `json:"expected_version"`
}

type UpdateKnowledgeDocRes = pb.UpdateKnowledgeDocRes

type GetKnowledgeDocDetailReq struct {
	g.Meta         `path:"/v1/agent/admin/knowledge-docs/{knowledge_doc_no}" method:"get" tags:"Agent-Knowledge" summary:"Get knowledge document detail"`
	KnowledgeDocNo string `json:"knowledge_doc_no" in:"path" v:"required#knowledge_doc_no is required"`
}

type GetKnowledgeDocDetailRes = pb.GetKnowledgeDocDetailRes

type ListKnowledgeDocsReq struct {
	g.Meta `path:"/v1/agent/admin/knowledge-docs" method:"get" tags:"Agent-Knowledge" summary:"List knowledge documents"`
	pb.ListKnowledgeDocsReq
}

type ListKnowledgeDocsRes = pb.ListKnowledgeDocsRes

type SubmitKnowledgeDocReq struct {
	g.Meta          `path:"/v1/agent/admin/knowledge-docs/{knowledge_doc_no}:submit" method:"post" tags:"Agent-Knowledge" summary:"Submit knowledge document"`
	KnowledgeDocNo  string `json:"knowledge_doc_no" in:"path" v:"required#knowledge_doc_no is required"`
	ExpectedVersion uint64 `json:"expected_version"`
}

type SubmitKnowledgeDocRes = pb.SubmitKnowledgeDocRes

type ApproveKnowledgeDocReq struct {
	g.Meta          `path:"/v1/agent/admin/knowledge-docs/{knowledge_doc_no}:approve" method:"post" tags:"Agent-Knowledge" summary:"Approve knowledge document"`
	KnowledgeDocNo  string `json:"knowledge_doc_no" in:"path" v:"required#knowledge_doc_no is required"`
	ExpectedVersion uint64 `json:"expected_version"`
	Remark          string `json:"remark"`
}

type ApproveKnowledgeDocRes = pb.ApproveKnowledgeDocRes

type RejectKnowledgeDocReq struct {
	g.Meta           `path:"/v1/agent/admin/knowledge-docs/{knowledge_doc_no}:reject" method:"post" tags:"Agent-Knowledge" summary:"Reject knowledge document"`
	KnowledgeDocNo   string `json:"knowledge_doc_no" in:"path" v:"required#knowledge_doc_no is required"`
	ExpectedVersion  uint64 `json:"expected_version"`
	RejectReasonCode string `json:"reject_reason_code"`
	RejectComment    string `json:"reject_comment"`
}

type RejectKnowledgeDocRes = pb.RejectKnowledgeDocRes

type PublishKnowledgeDocReq struct {
	g.Meta          `path:"/v1/agent/admin/knowledge-docs/{knowledge_doc_no}:publish" method:"post" tags:"Agent-Knowledge" summary:"Publish knowledge document"`
	KnowledgeDocNo  string `json:"knowledge_doc_no" in:"path" v:"required#knowledge_doc_no is required"`
	ExpectedVersion uint64 `json:"expected_version"`
}

type PublishKnowledgeDocRes = pb.PublishKnowledgeDocRes

type OfflineKnowledgeDocReq struct {
	g.Meta          `path:"/v1/agent/admin/knowledge-docs/{knowledge_doc_no}:offline" method:"post" tags:"Agent-Knowledge" summary:"Offline knowledge document"`
	KnowledgeDocNo  string `json:"knowledge_doc_no" in:"path" v:"required#knowledge_doc_no is required"`
	ExpectedVersion uint64 `json:"expected_version"`
	ReasonCode      string `json:"reason_code"`
}

type OfflineKnowledgeDocRes = pb.OfflineKnowledgeDocRes

type ReindexKnowledgeDocReq struct {
	g.Meta          `path:"/v1/agent/admin/knowledge-docs/{knowledge_doc_no}:reindex" method:"post" tags:"Agent-Knowledge" summary:"Reindex knowledge document"`
	KnowledgeDocNo  string `json:"knowledge_doc_no" in:"path" v:"required#knowledge_doc_no is required"`
	ExpectedVersion uint64 `json:"expected_version"`
	ReasonCode      string `json:"reason_code"`
}

type ReindexKnowledgeDocRes = pb.ReindexKnowledgeDocRes

type SearchKnowledgeChunksReq struct {
	g.Meta `path:"/v1/agent/admin/knowledge-chunks:search" method:"get" tags:"Agent-Knowledge" summary:"Search knowledge chunks"`
	pb.SearchKnowledgeChunksReq
}

type SearchKnowledgeChunksRes = pb.SearchKnowledgeChunksRes

type ProcessUserTurnReq struct {
	g.Meta `path:"/v1/internal/agent/runs:process-turn" method:"post" tags:"Agent-Internal" summary:"Process user turn"`
	pb.ProcessUserTurnReq
}

type ProcessUserTurnRes = pb.ProcessUserTurnRes

type AbortRunReq struct {
	g.Meta `path:"/v1/internal/agent/runs:abort" method:"post" tags:"Agent-Internal" summary:"Abort run"`
	pb.AbortRunReq
}

type AbortRunRes = pb.AbortRunRes

type RefreshSessionSummaryReq struct {
	g.Meta `path:"/v1/internal/agent/conversations:refresh-summary" method:"post" tags:"Agent-Internal" summary:"Refresh session summary"`
	pb.RefreshSessionSummaryReq
}

type RefreshSessionSummaryRes = pb.RefreshSessionSummaryRes

type GenerateHandoffSummaryReq struct {
	g.Meta `path:"/v1/internal/agent/handoffs:generate" method:"post" tags:"Agent-Internal" summary:"Generate handoff summary"`
	pb.GenerateHandoffSummaryReq
}

type GenerateHandoffSummaryRes = pb.GenerateHandoffSummaryRes

type UpsertKnowledgeChunksReq struct {
	g.Meta `path:"/v1/internal/agent/knowledge-chunks:upsert" method:"post" tags:"Agent-Internal" summary:"Upsert knowledge chunks"`
	pb.UpsertKnowledgeChunksReq
}

type UpsertKnowledgeChunksRes = pb.UpsertKnowledgeChunksRes

type MarkKnowledgeSourceDeletedReq struct {
	g.Meta `path:"/v1/internal/agent/knowledge-sources:mark-deleted" method:"post" tags:"Agent-Internal" summary:"Mark knowledge source deleted"`
	pb.MarkKnowledgeSourceDeletedReq
}

type MarkKnowledgeSourceDeletedRes = pb.MarkKnowledgeSourceDeletedRes
