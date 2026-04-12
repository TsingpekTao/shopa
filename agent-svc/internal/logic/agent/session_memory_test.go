package agent

import (
	"testing"

	"github.com/TsingpekTao/shopa/agent-svc/internal/model/entity"
)

func TestResolveConversationTaskSessionPrefersRicherStateOnVersionTie(t *testing.T) {
	runs := []entity.AgentRun{
		{
			Id: 101,
			GraphStateJson: `{
				"task_session": {
					"conversation_no": "ACV202604110401",
					"version": 3,
					"active_task_code": "logistics_query",
					"user_goal": "check logistics"
				}
			}`,
		},
		{
			Id: 102,
			GraphStateJson: `{
				"task_session": {
					"conversation_no": "ACV202604110401",
					"version": 3,
					"active_task_code": "logistics_query",
					"selected_order_no": "ORD202604110401",
					"slot_values": {
						"order_no": "ORD202604110401"
					},
					"latest_facts_summary": "order selected"
				}
			}`,
		},
	}

	got := resolveConversationTaskSession(runs)
	if got.SessionVersion != 3 {
		t.Fatalf("expected session version 3, got %d", got.SessionVersion)
	}
	if got.SelectedOrderNo != "ORD202604110401" {
		t.Fatalf("expected richer session with selected order to win, got %q", got.SelectedOrderNo)
	}
	if got.SlotValues["order_no"] != "ORD202604110401" {
		t.Fatalf("expected selected order slot to survive merge, got %+v", got.SlotValues)
	}
}
