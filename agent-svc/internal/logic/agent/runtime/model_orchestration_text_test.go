package runtime

import (
	"strings"
	"testing"
)

func TestBuildPlannerSystemPromptUsesReadableChinese(t *testing.T) {
	prompt := buildPlannerSystemPrompt(&runState{
		TaskSession: TaskSessionState{
			SelectedOrderNo:       "ORD202604110001",
			SelectedSubOrderNo:    "SUB202604110001",
			LatestFactsSummary:    "order facts",
			LatestDecisionSummary: "decision facts",
		},
	})

	if !strings.Contains(prompt, "商城官方客服的 Planner") {
		t.Fatalf("expected readable planner role text, got %q", prompt)
	}
	if !strings.Contains(prompt, "当前 selected_order_no=ORD202604110001") {
		t.Fatalf("expected readable selected order text, got %q", prompt)
	}
}

func TestHeuristicPlannerDecisionGreetingUsesReadableReply(t *testing.T) {
	decision := heuristicPlannerDecision(&runState{
		Message: UserMessage{ContentText: "你好"},
	})

	if !strings.Contains(decision.NextAction.Reason, "您好") {
		t.Fatalf("expected readable greeting reply, got %q", decision.NextAction.Reason)
	}
}

func TestBuildPolicyReplyTextUsesReadableChinese(t *testing.T) {
	reply := buildPolicyReplyText(nil)
	if !strings.Contains(reply, "规则") {
		t.Fatalf("expected readable policy empty-state reply, got %q", reply)
	}
}
