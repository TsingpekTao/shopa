package runtime

import (
	"context"
	"testing"
)

func TestRunnerRunReturnsStructuredGuardPayload(t *testing.T) {
	runner := NewRunnerWithOptions(RunnerOptions{})

	out, err := runner.Run(context.Background(), RunInput{
		UserQuery: "在吗",
		Message: UserMessage{
			ContentText: "在吗",
		},
	})
	if err != nil {
		t.Fatalf("runner returned error: %v", err)
	}

	if out.ReplyPayload.IntentCode != intentCodeGreeting {
		t.Fatalf("expected greeting reply payload intent, got %q", out.ReplyPayload.IntentCode)
	}
	if out.ReplyPayload.GuardResultCode != guardResultGreeting {
		t.Fatalf("expected greeting guard result, got %q", out.ReplyPayload.GuardResultCode)
	}
	if out.AnswerText != out.ReplyPayload.ReplyText {
		t.Fatalf("expected answer text to mirror reply payload text")
	}
}

func TestRunnerRunReusesPreviousTaskSessionForSlotFuse(t *testing.T) {
	runner := NewRunnerWithOptions(RunnerOptions{})

	out, err := runner.Run(context.Background(), RunInput{
		UserQuery: "我这个能退款吗",
		Message: UserMessage{
			ContentText: "我这个能退款吗",
		},
		PreviousSession: TaskSessionState{
			ActiveTaskCode:     taskCodeRefundDecision,
			SlotRetryCount:     1,
			LastSlotPromptCode: slotPromptCodeOrderNo,
			MissingSlots:       []MissingSlot{{SlotCode: slotCodeOrderNo, PromptText: "请提供订单号", Required: true}},
		},
		Security: SecurityContext{
			UserID:         30001,
			ShopNo:         "SHOP30001",
			RequestID:      "REQ202604090101",
			ConversationNo: "ACV202604090101",
			RunNo:          "ARN202604090101",
		},
	})
	if err != nil {
		t.Fatalf("runner returned error: %v", err)
	}

	if out.ReplyPayload.HandoffReasonCode != escalationReasonSlotFillingFailed {
		t.Fatalf("expected slot filling failed handoff reason, got %q", out.ReplyPayload.HandoffReasonCode)
	}
	if !out.ReplyPayload.HandoffRecommended {
		t.Fatalf("expected handoff to be recommended")
	}
	if out.TaskSession.SlotRetryCount != 2 {
		t.Fatalf("expected slot retry count 2, got %d", out.TaskSession.SlotRetryCount)
	}
}
