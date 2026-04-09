package runtime

import (
	"context"
	"testing"
)

type stubOrderSnapshotRepository struct {
	lastFilter OrderOwnershipFilter
	snapshot   *OrderSnapshot
	err        error
}

func (s *stubOrderSnapshotRepository) QueryOrderSnapshot(ctx context.Context, filter OrderOwnershipFilter) (*OrderSnapshot, error) {
	_ = ctx
	s.lastFilter = filter
	return s.snapshot, s.err
}

func TestDetectGuardIntentGreeting(t *testing.T) {
	got := detectGuardIntent(UserMessage{
		ContentText: "在吗",
	})

	if got.IntentCode != intentCodeGreeting {
		t.Fatalf("expected greeting intent, got %q", got.IntentCode)
	}
	if got.GuardResultCode != guardResultGreeting {
		t.Fatalf("expected greeting guard result, got %q", got.GuardResultCode)
	}
}

func TestDetectGuardIntentOutOfScope(t *testing.T) {
	got := detectGuardIntent(UserMessage{
		ContentText: "这件衣服透气吗",
	})

	if got.IntentCode != intentCodeOutOfScope {
		t.Fatalf("expected out_of_scope intent, got %q", got.IntentCode)
	}
	if got.GuardResultCode != guardResultOutOfScope {
		t.Fatalf("expected out_of_scope guard result, got %q", got.GuardResultCode)
	}
}

func TestDetectGuardIntentAttachmentOnly(t *testing.T) {
	got := detectGuardIntent(UserMessage{
		ContentText:     "   ",
		MessageTypeCode: "IMAGE",
		AssetIDs:        []uint64{101},
	})

	if got.IntentCode != intentCodeAttachmentOnly {
		t.Fatalf("expected attachment_only intent, got %q", got.IntentCode)
	}
	if got.GuardResultCode != guardResultAttachmentOnly {
		t.Fatalf("expected attachment_only guard result, got %q", got.GuardResultCode)
	}
}

func TestPlanAfterSaleTurnEscalatesAfterSlotRetryThreshold(t *testing.T) {
	out, err := planAfterSaleTurn(context.Background(), AfterSaleTurnInput{
		Message: UserMessage{
			ContentText: "我这个能退款吗",
		},
		PreviousSession: TaskSessionState{
			ActiveTaskCode:      taskCodeRefundDecision,
			SlotRetryCount:      1,
			LastSlotPromptCode:  slotPromptCodeOrderNo,
			MissingSlots:        []MissingSlot{{SlotCode: slotCodeOrderNo, PromptText: "请提供订单号", Required: true}},
			SlotValues:          map[string]string{},
			UnresolvedTurnCount: 1,
		},
		Security: SecurityContext{
			UserID:         20001,
			ShopNo:         "SHOP1001",
			RequestID:      "REQ202604090001",
			ConversationNo: "ACV202604090001",
			RunNo:          "ARN202604090001",
		},
	})
	if err != nil {
		t.Fatalf("plan after-sale turn returned error: %v", err)
	}

	if out.Session.SlotRetryCount != 2 {
		t.Fatalf("expected slot retry count 2, got %d", out.Session.SlotRetryCount)
	}
	if !out.Session.SlotFillingFailed {
		t.Fatalf("expected slot filling failed to be true")
	}
	if out.Session.EscalationReasonCode != escalationReasonSlotFillingFailed {
		t.Fatalf("expected escalation reason %q, got %q", escalationReasonSlotFillingFailed, out.Session.EscalationReasonCode)
	}
	if !out.Reply.HandoffRecommended {
		t.Fatalf("expected handoff recommended to be true")
	}
	if out.Reply.HandoffReasonCode != escalationReasonSlotFillingFailed {
		t.Fatalf("expected handoff reason %q, got %q", escalationReasonSlotFillingFailed, out.Reply.HandoffReasonCode)
	}
}

func TestSecureQueryOrderSnapshotUsesUserScopedFilter(t *testing.T) {
	repo := &stubOrderSnapshotRepository{}

	out, err := secureQueryOrderSnapshot(context.Background(), repo, SecurityContext{
		UserID:         9527,
		ShopNo:         "SHOP9527",
		RequestID:      "REQ202604090002",
		ConversationNo: "ACV202604090002",
		RunNo:          "ARN202604090002",
	}, OrderQuery{
		OrderNo:    "ORD202604090001",
		SubOrderNo: "SUB202604090001",
	})
	if err != nil {
		t.Fatalf("secure query returned error: %v", err)
	}

	if repo.lastFilter.UserID != 9527 {
		t.Fatalf("expected user scoped filter to carry user_id, got %d", repo.lastFilter.UserID)
	}
	if repo.lastFilter.OrderNo != "ORD202604090001" {
		t.Fatalf("expected order_no to be forwarded, got %q", repo.lastFilter.OrderNo)
	}
	if repo.lastFilter.SubOrderNo != "SUB202604090001" {
		t.Fatalf("expected sub_order_no to be forwarded, got %q", repo.lastFilter.SubOrderNo)
	}
	if out.ErrorCode != orderLookupErrorNotFound {
		t.Fatalf("expected not found error code when repository returns no row, got %q", out.ErrorCode)
	}
}

func TestPlanAfterSaleTurnBuildsRefundDecisionFromOrderFacts(t *testing.T) {
	repo := &stubOrderSnapshotRepository{
		snapshot: &OrderSnapshot{
			OrderNo:            "ORD202604090003",
			MainStatus:         "PAID",
			PaymentStatus:      "PAID",
			FulfillmentStatus:  "UNSHIPPED",
			LogisticsStatus:    "NOT_SHIPPED",
			AfterSaleStatus:    "NONE",
			LatestUpdateTime:   "2026-04-09T12:00:00+08:00",
			OwnershipConfirmed: true,
		},
	}

	out, err := planAfterSaleTurn(context.Background(), AfterSaleTurnInput{
		Message: UserMessage{
			ContentText: "订单号 ORD202604090003 现在能退款吗",
		},
		Security: SecurityContext{
			UserID:         10086,
			ShopNo:         "SHOP10086",
			RequestID:      "REQ202604090003",
			ConversationNo: "ACV202604090003",
			RunNo:          "ARN202604090003",
		},
		OrderRepository: repo,
	})
	if err != nil {
		t.Fatalf("plan after-sale turn returned error: %v", err)
	}

	if out.Reply.IntentCode != intentCodeAfterSaleTask {
		t.Fatalf("expected after-sale intent, got %q", out.Reply.IntentCode)
	}
	if len(out.Reply.DataCards) != 2 {
		t.Fatalf("expected two data cards, got %d", len(out.Reply.DataCards))
	}
	if out.Reply.DataCards[0].OrderSnapshotCard == nil {
		t.Fatalf("expected order snapshot card in first data card")
	}
	if out.Reply.DataCards[0].OrderSnapshotCard.OrderNo != "ORD202604090003" {
		t.Fatalf("expected order snapshot card to keep order_no, got %q", out.Reply.DataCards[0].OrderSnapshotCard.OrderNo)
	}
	if out.Reply.DataCards[1].AfterSaleDecisionCard == nil {
		t.Fatalf("expected decision card in second data card")
	}
	if out.Reply.DataCards[1].AfterSaleDecisionCard.DecisionPathCode != decisionPathRefundOnly {
		t.Fatalf("expected refund decision path, got %q", out.Reply.DataCards[1].AfterSaleDecisionCard.DecisionPathCode)
	}
	if len(out.Reply.SuggestedActions) == 0 {
		t.Fatalf("expected suggested actions to be present")
	}
	if out.Reply.SuggestedActions[0].ActionCode != actionCodeRequestRefund {
		t.Fatalf("expected first suggested action to be refund, got %q", out.Reply.SuggestedActions[0].ActionCode)
	}
	if !out.Reply.SuggestedActions[0].Enabled {
		t.Fatalf("expected refund suggested action to be enabled")
	}
}
