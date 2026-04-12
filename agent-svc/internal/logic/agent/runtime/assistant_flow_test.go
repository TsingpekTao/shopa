package runtime

import (
	"context"
	"testing"
)

type stubOrderSnapshotRepository struct {
	lastFilter   OrderOwnershipFilter
	snapshot     *OrderSnapshot
	err          error
	recentFilter RecentOrderListFilter
	recentOrders []RecentOrderCandidate
	recentErr    error
}

func (s *stubOrderSnapshotRepository) QueryOrderSnapshot(ctx context.Context, filter OrderOwnershipFilter) (*OrderSnapshot, error) {
	_ = ctx
	s.lastFilter = filter
	return s.snapshot, s.err
}

func (s *stubOrderSnapshotRepository) ListRecentOrders(ctx context.Context, filter RecentOrderListFilter) ([]RecentOrderCandidate, error) {
	_ = ctx
	s.recentFilter = filter
	return append([]RecentOrderCandidate(nil), s.recentOrders...), s.recentErr
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
	if out.Reply.DataCards[1].AfterSaleDecisionCard.SceneCode != sceneCodeRefundBeforeShip {
		t.Fatalf("expected refund-before-shipment scene code, got %q", out.Reply.DataCards[1].AfterSaleDecisionCard.SceneCode)
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

func TestPlanAfterSaleTurnEscalatesAfterRepeatedUnresolvedLookup(t *testing.T) {
	repo := &stubOrderSnapshotRepository{}

	out, err := planAfterSaleTurn(context.Background(), AfterSaleTurnInput{
		Message: UserMessage{
			ContentText: "订单号 ORD202604090004 现在能处理吗",
		},
		PreviousSession: TaskSessionState{
			SlotValues: map[string]string{
				slotCodeOrderNo: "ORD202604090004",
				slotCodeProblem: problemTypeRefund,
			},
			UnresolvedTurnCount: 1,
		},
		Security: SecurityContext{
			UserID:         10010,
			ShopNo:         "SHOP10010",
			RequestID:      "REQ202604090004",
			ConversationNo: "ACV202604090004",
			RunNo:          "ARN202604090004",
		},
		OrderRepository: repo,
	})
	if err != nil {
		t.Fatalf("plan after-sale turn returned error: %v", err)
	}

	if !out.Reply.HandoffRecommended {
		t.Fatalf("expected repeated unresolved lookup to recommend handoff")
	}
	if out.Reply.HandoffReasonCode != escalationReasonHighRiskCase {
		t.Fatalf("expected repeated unresolved lookup handoff reason %q, got %q", escalationReasonHighRiskCase, out.Reply.HandoffReasonCode)
	}
	if out.Session.UnresolvedTurnCount != 2 {
		t.Fatalf("expected unresolved turn count 2, got %d", out.Session.UnresolvedTurnCount)
	}
}

func TestPlanAfterSaleTurnReturnsRecentOrderSelectionWhenOrderMissing(t *testing.T) {
	repo := &stubOrderSnapshotRepository{
		recentOrders: []RecentOrderCandidate{
			{
				OrderNo:           "ORD-IN-TRANSIT",
				DisplayTitle:      "云感运动鞋",
				FulfillmentStatus: "SHIPPED",
				LogisticsStatus:   "IN_TRANSIT",
				LatestUpdateTime:  "2026-04-10T09:30:00+08:00",
			},
			{
				OrderNo:           "ORD-UNSHIPPED",
				DisplayTitle:      "亚麻衬衫",
				FulfillmentStatus: "UNSHIPPED",
				LogisticsStatus:   "NOT_SHIPPED",
				LatestUpdateTime:  "2026-04-10T08:30:00+08:00",
			},
		},
	}

	out, err := planAfterSaleTurn(context.Background(), AfterSaleTurnInput{
		Message: UserMessage{
			ContentText: "帮我查一下这个订单现在到哪了",
		},
		Security: SecurityContext{
			UserID:         30001,
			ShopNo:         "SHOP30001",
			RequestID:      "REQ202604100001",
			ConversationNo: "ACV202604100001",
			RunNo:          "ARN202604100001",
		},
		OrderRepository: repo,
	})
	if err != nil {
		t.Fatalf("plan after-sale turn returned error: %v", err)
	}

	if repo.recentFilter.UserID != 30001 {
		t.Fatalf("expected recent order lookup to use user scope, got %d", repo.recentFilter.UserID)
	}
	if repo.recentFilter.Limit != defaultRecentOrderCandidateLimit {
		t.Fatalf("expected recent order lookup limit %d, got %d", defaultRecentOrderCandidateLimit, repo.recentFilter.Limit)
	}
	if len(out.Reply.DataCards) == 0 || out.Reply.DataCards[0].OrderSelectionCard == nil {
		t.Fatalf("expected order selection card when order slot is missing")
	}
	card := out.Reply.DataCards[0].OrderSelectionCard
	if len(card.Candidates) != 2 {
		t.Fatalf("expected two recent order candidates, got %d", len(card.Candidates))
	}
	if card.Candidates[0].OrderNo != "ORD-IN-TRANSIT" {
		t.Fatalf("expected in-transit order to rank first for logistics query, got %q", card.Candidates[0].OrderNo)
	}
	if out.Session.PendingSelectionQuery != "帮我查一下这个订单现在到哪了" {
		t.Fatalf("expected pending selection query to be saved, got %q", out.Session.PendingSelectionQuery)
	}
	if out.Session.PendingSelectionTaskCode != taskCodeLogisticsQuery {
		t.Fatalf("expected pending selection task %q, got %q", taskCodeLogisticsQuery, out.Session.PendingSelectionTaskCode)
	}
}

func TestPlanAfterSaleTurnQuickPromptRequiresFreshOrderSelectionEvenWithPreviousOrder(t *testing.T) {
	repo := &stubOrderSnapshotRepository{
		recentOrders: []RecentOrderCandidate{
			{
				OrderNo:           "ORD-NEW-1",
				DisplayTitle:      "新订单一",
				FulfillmentStatus: "SHIPPED",
				LogisticsStatus:   "IN_TRANSIT",
				LatestUpdateTime:  "2026-04-10T12:30:00+08:00",
			},
			{
				OrderNo:           "ORD-NEW-2",
				DisplayTitle:      "新订单二",
				FulfillmentStatus: "UNSHIPPED",
				LogisticsStatus:   "NOT_SHIPPED",
				LatestUpdateTime:  "2026-04-10T11:30:00+08:00",
			},
		},
		snapshot: &OrderSnapshot{
			OrderNo:            "ORD-OLD-1",
			MainStatus:         "PAID",
			PaymentStatus:      "PAID",
			FulfillmentStatus:  "SHIPPED",
			LogisticsStatus:    "IN_TRANSIT",
			AfterSaleStatus:    "NONE",
			LatestUpdateTime:   "2026-04-10T10:00:00+08:00",
			OwnershipConfirmed: true,
		},
	}

	out, err := planAfterSaleTurn(context.Background(), AfterSaleTurnInput{
		Message: UserMessage{
			ContentText: "帮我查一下这个订单现在到哪了",
		},
		PreviousSession: TaskSessionState{
			ActiveTaskCode:  taskCodeLogisticsQuery,
			AnchoredOrderNo: "ORD-OLD-1",
			SlotValues: map[string]string{
				slotCodeOrderNo: "ORD-OLD-1",
				slotCodeProblem: problemTypeLogistics,
			},
			LatestFactsSummary:    "order_no=ORD-OLD-1",
			LatestDecisionSummary: "旧订单结果",
		},
		Security: SecurityContext{
			UserID:         30005,
			ShopNo:         "SHOP30005",
			RequestID:      "REQ202604100305",
			ConversationNo: "ACV202604100305",
			RunNo:          "ARN202604100305",
		},
		OrderRepository: repo,
	})
	if err != nil {
		t.Fatalf("plan after-sale turn returned error: %v", err)
	}

	if len(out.Reply.DataCards) == 0 || out.Reply.DataCards[0].OrderSelectionCard == nil {
		t.Fatalf("expected quick prompt to force a fresh order selection card")
	}
	if repo.lastFilter.OrderNo != "" {
		t.Fatalf("expected no direct snapshot lookup for previous order, got %q", repo.lastFilter.OrderNo)
	}
	if got := out.Session.SlotValues[slotCodeOrderNo]; got != "" {
		t.Fatalf("expected previous anchored order to be cleared before reselection, got %q", got)
	}
}

func TestPlanAfterSaleTurnQuickPromptRequiresFreshOrderSelectionEvenWithConversationAnchor(t *testing.T) {
	repo := &stubOrderSnapshotRepository{
		recentOrders: []RecentOrderCandidate{
			{
				OrderNo:           "ORD-ANCHOR-NEW-1",
				DisplayTitle:      "候选订单一",
				FulfillmentStatus: "SHIPPED",
				LogisticsStatus:   "IN_TRANSIT",
				LatestUpdateTime:  "2026-04-10T13:30:00+08:00",
			},
		},
		snapshot: &OrderSnapshot{
			OrderNo:            "ORD-ANCHOR-OLD",
			MainStatus:         "PAID",
			PaymentStatus:      "PAID",
			FulfillmentStatus:  "SHIPPED",
			LogisticsStatus:    "IN_TRANSIT",
			AfterSaleStatus:    "NONE",
			LatestUpdateTime:   "2026-04-10T10:00:00+08:00",
			OwnershipConfirmed: true,
		},
	}

	out, err := planAfterSaleTurn(context.Background(), AfterSaleTurnInput{
		Message: UserMessage{
			ContentText: "帮我查一下这个订单现在到哪了",
		},
		Conversation: ConversationAnchors{
			OrderNo: "ORD-ANCHOR-OLD",
		},
		Security: SecurityContext{
			UserID:         30006,
			ShopNo:         "SHOP30006",
			RequestID:      "REQ202604100306",
			ConversationNo: "ACV202604100306",
			RunNo:          "ARN202604100306",
		},
		OrderRepository: repo,
	})
	if err != nil {
		t.Fatalf("plan after-sale turn returned error: %v", err)
	}

	if len(out.Reply.DataCards) == 0 || out.Reply.DataCards[0].OrderSelectionCard == nil {
		t.Fatalf("expected quick prompt to force a fresh order selection card even with conversation anchor")
	}
	if repo.lastFilter.OrderNo != "" {
		t.Fatalf("expected no direct snapshot lookup when quick prompt should reselection, got %q", repo.lastFilter.OrderNo)
	}
}

func TestRankRecentOrdersPrefersUnshippedForUrgeShipment(t *testing.T) {
	candidates := []RecentOrderCandidate{
		{
			OrderNo:           "ORD-SHIPPED",
			FulfillmentStatus: "SHIPPED",
			LogisticsStatus:   "IN_TRANSIT",
			LatestUpdateTime:  "2026-04-10T09:30:00+08:00",
		},
		{
			OrderNo:           "ORD-UNSHIPPED",
			FulfillmentStatus: "UNSHIPPED",
			LogisticsStatus:   "NOT_SHIPPED",
			LatestUpdateTime:  "2026-04-10T08:30:00+08:00",
		},
	}

	ranked := rankRecentOrders(taskCodeLogisticsQuery, problemTypeUrgeShipment, "这个订单怎么还没发货", candidates)
	if len(ranked) != 2 {
		t.Fatalf("expected two ranked candidates, got %d", len(ranked))
	}
	if ranked[0].OrderNo != "ORD-UNSHIPPED" {
		t.Fatalf("expected unshipped order to rank first for urge shipment, got %q", ranked[0].OrderNo)
	}
}

func TestPlanAfterSaleTurnResumesPendingSelectionAfterBuyerChoosesOrder(t *testing.T) {
	repo := &stubOrderSnapshotRepository{
		snapshot: &OrderSnapshot{
			OrderNo:            "ORD202604100088",
			MainStatus:         "PAID",
			PaymentStatus:      "PAID",
			FulfillmentStatus:  "SHIPPED",
			LogisticsStatus:    "IN_TRANSIT",
			AfterSaleStatus:    "NONE",
			LatestUpdateTime:   "2026-04-10T10:00:00+08:00",
			OwnershipConfirmed: true,
		},
	}

	out, err := planAfterSaleTurn(context.Background(), AfterSaleTurnInput{
		Message: UserMessage{
			ContentText: "ORD202604100088",
		},
		PreviousSession: TaskSessionState{
			ActiveTaskCode:           taskCodeLogisticsQuery,
			PendingSelectionQuery:    "帮我查一下这个订单现在到哪了",
			PendingSelectionTaskCode: taskCodeLogisticsQuery,
			PendingSelectionCandidates: []RecentOrderCandidate{
				{OrderNo: "ORD202604100088", FulfillmentStatus: "SHIPPED", LogisticsStatus: "IN_TRANSIT"},
			},
		},
		Security: SecurityContext{
			UserID:         30002,
			ShopNo:         "SHOP30002",
			RequestID:      "REQ202604100002",
			ConversationNo: "ACV202604100002",
			RunNo:          "ARN202604100002",
		},
		OrderRepository: repo,
	})
	if err != nil {
		t.Fatalf("plan after-sale turn returned error: %v", err)
	}

	if repo.lastFilter.OrderNo != "ORD202604100088" {
		t.Fatalf("expected selected order to be queried, got %q", repo.lastFilter.OrderNo)
	}
	if out.Session.PendingSelectionQuery != "" {
		t.Fatalf("expected pending selection query to be cleared after selection, got %q", out.Session.PendingSelectionQuery)
	}
	if len(out.Reply.DataCards) == 0 || out.Reply.DataCards[0].OrderSnapshotCard == nil {
		t.Fatalf("expected order snapshot card after selecting a candidate order")
	}
}

func TestPlanAfterSaleTurnMatchesColloquialQueriesWithAnchoredOrder(t *testing.T) {
	repo := &stubOrderSnapshotRepository{
		snapshot: &OrderSnapshot{
			OrderNo:            "ORD202604100188",
			MainStatus:         "PAID",
			PaymentStatus:      "PAID",
			FulfillmentStatus:  "SHIPPED",
			LogisticsStatus:    "IN_TRANSIT",
			AfterSaleStatus:    "NONE",
			LatestUpdateTime:   "2026-04-10T11:00:00+08:00",
			OwnershipConfirmed: true,
		},
	}

	cases := []struct {
		name          string
		query         string
		wantSceneCode string
	}{
		{
			name:          "logistics colloquial wording",
			query:         "帮我看看这个订单到哪了",
			wantSceneCode: sceneCodeLogisticsWatch,
		},
		{
			name:          "refund colloquial wording",
			query:         "这个单现在还能退吗",
			wantSceneCode: sceneCodeRefundInTransit,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := planAfterSaleTurn(context.Background(), AfterSaleTurnInput{
				Message: UserMessage{
					ContentText: tc.query,
				},
				Conversation: ConversationAnchors{
					OrderNo: "ORD202604100188",
				},
				Security: SecurityContext{
					UserID:         30003,
					ShopNo:         "SHOP30003",
					RequestID:      "REQ202604100188",
					ConversationNo: "ACV202604100188",
					RunNo:          "ARN202604100188",
				},
				OrderRepository: repo,
			})
			if err != nil {
				t.Fatalf("plan after-sale turn returned error: %v", err)
			}
			if len(out.Reply.DataCards) < 2 || out.Reply.DataCards[1].AfterSaleDecisionCard == nil {
				t.Fatalf("expected decision card to be returned")
			}
			gotSceneCode := out.Reply.DataCards[1].AfterSaleDecisionCard.SceneCode
			if gotSceneCode != tc.wantSceneCode {
				t.Fatalf("expected scene code %q, got %q with reply %q", tc.wantSceneCode, gotSceneCode, out.Reply.ReplyText)
			}
		})
	}
}

func TestPlanAfterSaleTurnExplainsUnshippedLogisticsQuery(t *testing.T) {
	repo := &stubOrderSnapshotRepository{
		snapshot: &OrderSnapshot{
			OrderNo:            "ORD202604100288",
			MainStatus:         "PAID",
			PaymentStatus:      "PAID",
			FulfillmentStatus:  "UNSHIPPED",
			LogisticsStatus:    "NOT_SHIPPED",
			AfterSaleStatus:    "NONE",
			LatestUpdateTime:   "2026-04-10T11:30:00+08:00",
			OwnershipConfirmed: true,
		},
	}

	out, err := planAfterSaleTurn(context.Background(), AfterSaleTurnInput{
		Message: UserMessage{
			ContentText: "我想看下这单物流怎么还没更新",
		},
		Conversation: ConversationAnchors{
			OrderNo: "ORD202604100288",
		},
		Security: SecurityContext{
			UserID:         30004,
			ShopNo:         "SHOP30004",
			RequestID:      "REQ202604100288",
			ConversationNo: "ACV202604100288",
			RunNo:          "ARN202604100288",
		},
		OrderRepository: repo,
	})
	if err != nil {
		t.Fatalf("plan after-sale turn returned error: %v", err)
	}
	if len(out.Reply.DataCards) < 2 || out.Reply.DataCards[1].AfterSaleDecisionCard == nil {
		t.Fatalf("expected decision card to be returned")
	}
	if out.Reply.DataCards[1].AfterSaleDecisionCard.SceneCode != "logistics_before_shipment" {
		t.Fatalf("expected unshipped logistics scene code, got %q with reply %q", out.Reply.DataCards[1].AfterSaleDecisionCard.SceneCode, out.Reply.ReplyText)
	}
	if out.Reply.ReplyText == "我已经定位到订单，但当前更适合先结合订单状态继续判断退款、退货退款还是换货。" {
		t.Fatalf("expected a specific unshipped-logistics explanation instead of generic reply")
	}
}
