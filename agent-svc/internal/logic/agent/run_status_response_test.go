package agent

import (
	"testing"

	"github.com/TsingpekTao/shopa/agent-svc/internal/model/entity"
)

func TestBuildAssistantRunStatusResponseAllowsNilRun(t *testing.T) {
	t.Parallel()

	conv := &entity.AgentConversation{
		Id:                  1,
		ConversationNo:      "ACV_TEST",
		UserId:              42,
		SceneCode:           "BUYER_ASSISTANT",
		ConversationStatusCode: conversationStatusActive,
	}

	res := buildAssistantRunStatusResponse(conv, nil, nil)
	if res == nil {
		t.Fatal("expected response")
	}
	if res.Conversation == nil || res.Conversation.ConversationNo != "ACV_TEST" {
		t.Fatalf("expected conversation to be preserved, got %+v", res.Conversation)
	}
	if res.Run != nil {
		t.Fatalf("expected nil run when no run exists, got %+v", res.Run)
	}
	if res.CurrentNodeCode != "" {
		t.Fatalf("expected empty current node code, got %q", res.CurrentNodeCode)
	}
	if res.QueueBlocked {
		t.Fatal("expected queue blocked to default false")
	}
	if res.DegradedReply {
		t.Fatal("expected degraded reply to default false")
	}
	if res.ReplyPayload != nil {
		t.Fatalf("expected nil reply payload, got %+v", res.ReplyPayload)
	}
}
