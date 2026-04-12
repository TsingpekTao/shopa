package agent

import (
	"testing"

	agentv1 "github.com/TsingpekTao/shopa/agent-svc/api/v1"
)

func TestExtractHiddenActionPrefersExplicitField(t *testing.T) {
	t.Parallel()

	req := &agentv1.SendAssistantMessageReq{
		ExtJson: `{"hidden_action":{"type":"SET_SLOT","key":"selected_order_no","value":"ORD_FROM_EXT"}}`,
		HiddenAction: &agentv1.AssistantHiddenAction{
			Type:  "SET_SLOT",
			Key:   "selected_order_no",
			Value: "ORD_FROM_FIELD",
		},
	}

	got := extractHiddenAction(req)
	if got == nil {
		t.Fatal("expected hidden action to be parsed")
	}
	if got.Type != "SET_SLOT" {
		t.Fatalf("expected hidden action type SET_SLOT, got %q", got.Type)
	}
	if got.Key != "selected_order_no" {
		t.Fatalf("expected hidden action key selected_order_no, got %q", got.Key)
	}
	if got.Value != "ORD_FROM_FIELD" {
		t.Fatalf("expected explicit field to win, got %q", got.Value)
	}
}

func TestExtractHiddenActionFallsBackToExtJSON(t *testing.T) {
	t.Parallel()

	req := &agentv1.SendAssistantMessageReq{
		ExtJson: `{"hidden_action":{"type":"SET_SLOT","key":"selected_sub_order_no","value":"SUB202604100001"}}`,
	}

	got := extractHiddenAction(req)
	if got == nil {
		t.Fatal("expected hidden action from ext_json to be parsed")
	}
	if got.Type != "SET_SLOT" {
		t.Fatalf("expected hidden action type SET_SLOT, got %q", got.Type)
	}
	if got.Key != "selected_sub_order_no" {
		t.Fatalf("expected hidden action key selected_sub_order_no, got %q", got.Key)
	}
	if got.Value != "SUB202604100001" {
		t.Fatalf("expected hidden action value from ext_json, got %q", got.Value)
	}
}
