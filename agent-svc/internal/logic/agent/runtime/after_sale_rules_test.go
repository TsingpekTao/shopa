package runtime

import "testing"

func TestDefaultAfterSaleRuleEngineMatchesUrgeShipmentScene(t *testing.T) {
	engine := NewAfterSaleRuleEngine(nil)

	out := engine.Evaluate(AfterSaleRuleContext{
		TaskCode:    taskCodeLogisticsQuery,
		ProblemType: problemTypeUrgeShipment,
		Snapshot: &OrderSnapshot{
			OrderNo:           "ORD202604090301",
			FulfillmentStatus: "UNSHIPPED",
			LogisticsStatus:   "NOT_SHIPPED",
		},
	})

	if out.DecisionCard == nil {
		t.Fatalf("expected decision card")
	}
	if out.DecisionCard.SceneCode != sceneCodeUrgeShipmentPending {
		t.Fatalf("expected urge shipment scene code, got %q", out.DecisionCard.SceneCode)
	}
	if out.DecisionCard.DecisionPathCode != decisionPathWaitShipment {
		t.Fatalf("expected wait shipment decision path, got %q", out.DecisionCard.DecisionPathCode)
	}
	if len(out.SuggestedActions) == 0 || out.SuggestedActions[0].ActionCode != actionCodeUrgeShipment {
		t.Fatalf("expected first suggested action to be urge shipment")
	}
}

func TestDefaultAfterSaleRuleEngineMatchesLogisticsWatchScene(t *testing.T) {
	engine := NewAfterSaleRuleEngine(nil)

	out := engine.Evaluate(AfterSaleRuleContext{
		TaskCode:    taskCodeLogisticsQuery,
		ProblemType: problemTypeLogistics,
		Snapshot: &OrderSnapshot{
			OrderNo:           "ORD202604090302",
			FulfillmentStatus: "SHIPPED",
			LogisticsStatus:   "IN_TRANSIT",
		},
	})

	if out.DecisionCard == nil {
		t.Fatalf("expected decision card")
	}
	if out.DecisionCard.SceneCode != sceneCodeLogisticsWatch {
		t.Fatalf("expected logistics watch scene code, got %q", out.DecisionCard.SceneCode)
	}
	if out.HandoffRecommended {
		t.Fatalf("expected logistics watch scene to avoid direct handoff")
	}
}

func TestCustomRuleEngineEscalatesOnConflict(t *testing.T) {
	engine := NewAfterSaleRuleEngine([]AfterSaleRule{
		{
			SceneCode: "scene_a",
			Priority:  100,
			Match: func(ctx AfterSaleRuleContext) bool {
				return true
			},
			Build: func(ctx AfterSaleRuleContext) AfterSaleRuleOutcome {
				return AfterSaleRuleOutcome{
					DecisionCard: &AfterSaleDecisionCard{SceneCode: "scene_a", DecisionPathCode: decisionPathRefundOnly},
					ReplyText:    "scene a",
				}
			},
		},
		{
			SceneCode: "scene_b",
			Priority:  100,
			Match: func(ctx AfterSaleRuleContext) bool {
				return true
			},
			Build: func(ctx AfterSaleRuleContext) AfterSaleRuleOutcome {
				return AfterSaleRuleOutcome{
					DecisionCard: &AfterSaleDecisionCard{SceneCode: "scene_b", DecisionPathCode: decisionPathWaitShipment},
					ReplyText:    "scene b",
				}
			},
		},
	})

	out := engine.Evaluate(AfterSaleRuleContext{
		TaskCode:    taskCodeRefundDecision,
		ProblemType: problemTypeRefund,
		Snapshot:    &OrderSnapshot{OrderNo: "ORD202604090303"},
	})

	if out.DecisionCard == nil {
		t.Fatalf("expected conflict decision card")
	}
	if out.DecisionCard.SceneCode != sceneCodeRuleConflict {
		t.Fatalf("expected rule conflict scene code, got %q", out.DecisionCard.SceneCode)
	}
	if !out.HandoffRecommended {
		t.Fatalf("expected conflict to recommend handoff")
	}
	if out.HandoffReasonCode != escalationReasonRuleConflict {
		t.Fatalf("expected handoff reason rule_conflict, got %q", out.HandoffReasonCode)
	}
}
