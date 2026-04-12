package runtime

import (
	"context"
	"strings"
	"testing"
)

type stubModelClient struct {
	responses []string
	err       error
}

func TestRunnerRunRepairsPlannerToolArgumentsOnce(t *testing.T) {
	knowledgeRepo := &stubPolicyKnowledgeRepository{
		hits: []PolicyKnowledgeHit{{
			SourceTypeCode: "PLATFORM_RULE",
			SourceID:       "refund-policy",
			SourceVersion:  1,
			Title:          "Refund Policy",
			Snippet:        "Unshipped orders can usually be refunded directly.",
		}},
	}
	runner := NewRunnerWithOptions(RunnerOptions{
		PolicyKnowledgeRepository: knowledgeRepo,
		ModelClient: &stubModelClient{
			responses: []string{
				`{"thought":"Find the rule","intent":"POLICY_QA","user_goal":"Explain refund rules","confidence":0.8,"slot_values":{},"missing_slots":[],"next_action":{"type":"TOOL_CALL","tool_name":"search_knowledge_chunks","tool_args":{"query":123},"reason":"search policy"},"handoff_recommended":false,"handoff_reason_code":""}`,
				`{"thought":"Repair the query argument","intent":"POLICY_QA","user_goal":"Explain refund rules","confidence":0.86,"slot_values":{},"missing_slots":[],"next_action":{"type":"TOOL_CALL","tool_name":"search_knowledge_chunks","tool_args":{"query":"refund policy"},"reason":"search policy"},"handoff_recommended":false,"handoff_reason_code":""}`,
			},
		},
	})

	out, err := runner.Run(context.Background(), RunInput{
		UserQuery: "What are the refund rules?",
		Message: UserMessage{
			ContentText: "What are the refund rules?",
		},
		Security: SecurityContext{
			UserID:         31024,
			ShopNo:         "SHOP31024",
			RequestID:      "REQ202604110304",
			ConversationNo: "ACV202604110304",
			RunNo:          "ARN202604110304",
		},
	})
	if err != nil {
		t.Fatalf("runner returned error: %v", err)
	}

	if knowledgeRepo.lastQuery != "refund policy" {
		t.Fatalf("expected repaired planner query to be used, got %q", knowledgeRepo.lastQuery)
	}
	if out.ReplyPayload.IntentCode != officialIntentPolicyQA {
		t.Fatalf("expected policy intent after repair, got %q", out.ReplyPayload.IntentCode)
	}
}

func (s *stubModelClient) Complete(ctx context.Context, req ModelCompletionRequest) (string, error) {
	_ = ctx
	_ = req
	if s.err != nil {
		return "", s.err
	}
	if len(s.responses) == 0 {
		return "", nil
	}
	response := s.responses[0]
	s.responses = s.responses[1:]
	return response, nil
}

type stubPolicyKnowledgeRepository struct {
	lastQuery string
	lastLimit uint32
	hits      []PolicyKnowledgeHit
	err       error
}

func (s *stubPolicyKnowledgeRepository) SearchPolicyKnowledge(ctx context.Context, query string, limit uint32) ([]PolicyKnowledgeHit, error) {
	_ = ctx
	s.lastQuery = query
	s.lastLimit = limit
	return append([]PolicyKnowledgeHit(nil), s.hits...), s.err
}

type stubProductSearchRepository struct {
	lastFilter ProductSearchFilter
	items      []ProductSearchItem
	err        error
}

func (s *stubProductSearchRepository) SearchProducts(ctx context.Context, filter ProductSearchFilter) ([]ProductSearchItem, error) {
	_ = ctx
	s.lastFilter = filter
	return append([]ProductSearchItem(nil), s.items...), s.err
}

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

	if out.ReplyPayload.IntentCode != officialIntentGreeting {
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

func TestRunnerRunAppliesHiddenSelectedOrderBeforePlanning(t *testing.T) {
	repo := &stubOrderSnapshotRepository{
		snapshot: &OrderSnapshot{
			OrderNo:            "ORD202604110001",
			MainStatus:         "PAID",
			PaymentStatus:      "PAID",
			FulfillmentStatus:  "SHIPPED",
			LogisticsStatus:    "IN_TRANSIT",
			AfterSaleStatus:    "NONE",
			LatestUpdateTime:   "2026-04-11T09:00:00+08:00",
			OwnershipConfirmed: true,
		},
	}
	runner := NewRunnerWithOptions(RunnerOptions{
		OrderRepository: repo,
	})

	out, err := runner.Run(context.Background(), RunInput{
		UserQuery: "查询这笔订单的物流",
		Message: UserMessage{
			ContentText: "查询这笔订单的物流",
		},
		HiddenAction: HiddenAction{
			Type:  "SET_SLOT",
			Key:   "selected_order_no",
			Value: "ORD202604110001",
		},
		Security: SecurityContext{
			UserID:         31001,
			ShopNo:         "SHOP31001",
			RequestID:      "REQ202604110001",
			ConversationNo: "ACV202604110001",
			RunNo:          "ARN202604110001",
		},
	})
	if err != nil {
		t.Fatalf("runner returned error: %v", err)
	}

	if repo.lastFilter.OrderNo != "ORD202604110001" {
		t.Fatalf("expected hidden selected order to be queried directly, got %q", repo.lastFilter.OrderNo)
	}
	if out.TaskSession.SelectedOrderNo != "ORD202604110001" {
		t.Fatalf("expected selected order to persist in session, got %q", out.TaskSession.SelectedOrderNo)
	}
	if len(out.ReplyPayload.DataCards) == 0 || out.ReplyPayload.DataCards[0].OrderSnapshotCard == nil {
		t.Fatalf("expected direct order query to return structured order facts")
	}
}

func TestRunnerRunFallsBackWhenPlannerReturnsInvalidJSON(t *testing.T) {
	repo := &stubOrderSnapshotRepository{
		recentOrders: []RecentOrderCandidate{
			{
				OrderNo:           "ORD202604110101",
				DisplayTitle:      "春季风衣",
				FulfillmentStatus: "SHIPPED",
				LogisticsStatus:   "IN_TRANSIT",
				LatestUpdateTime:  "2026-04-11T10:00:00+08:00",
			},
		},
	}
	runner := NewRunnerWithOptions(RunnerOptions{
		OrderRepository: repo,
		ModelClient: &stubModelClient{
			responses: []string{"not-json"},
		},
	})

	out, err := runner.Run(context.Background(), RunInput{
		UserQuery: "帮我查一下这个订单现在到哪了",
		Message: UserMessage{
			ContentText: "帮我查一下这个订单现在到哪了",
		},
		Security: SecurityContext{
			UserID:         31011,
			ShopNo:         "SHOP31011",
			RequestID:      "REQ202604110101",
			ConversationNo: "ACV202604110101",
			RunNo:          "ARN202604110101",
		},
	})
	if err != nil {
		t.Fatalf("runner returned error: %v", err)
	}

	if len(out.ReplyPayload.DataCards) == 0 || out.ReplyPayload.DataCards[0].OrderSelectionCard == nil {
		t.Fatalf("expected invalid planner json to fall back to recent order selection")
	}
	if repo.recentFilter.UserID != 31011 {
		t.Fatalf("expected fallback planner to still use user scoped recent-order lookup, got %d", repo.recentFilter.UserID)
	}
}

func TestRunnerRunUsesModelPlannerToolDecision(t *testing.T) {
	repo := &stubOrderSnapshotRepository{
		snapshot: &OrderSnapshot{
			OrderNo:            "ORD202604110202",
			MainStatus:         "PAID",
			PaymentStatus:      "PAID",
			FulfillmentStatus:  "SHIPPED",
			LogisticsStatus:    "IN_TRANSIT",
			AfterSaleStatus:    "NONE",
			LatestUpdateTime:   "2026-04-11T10:30:00+08:00",
			OwnershipConfirmed: true,
		},
	}
	runner := NewRunnerWithOptions(RunnerOptions{
		OrderRepository: repo,
		ModelClient: &stubModelClient{
			responses: []string{`{"thought":"用户想查物流","intent":"LOGISTICS","user_goal":"查询物流进度","confidence":0.98,"slot_values":{},"missing_slots":[],"next_action":{"type":"TOOL_CALL","tool_name":"query_logistics","tool_args":{},"reason":"已有选中订单"},"handoff_recommended":false,"handoff_reason_code":""}`},
		},
	})

	out, err := runner.Run(context.Background(), RunInput{
		UserQuery: "这笔订单到哪了",
		Message: UserMessage{
			ContentText: "这笔订单到哪了",
		},
		HiddenAction: HiddenAction{
			Type:  "SET_SLOT",
			Key:   "selected_order_no",
			Value: "ORD202604110202",
		},
		Security: SecurityContext{
			UserID:         31012,
			ShopNo:         "SHOP31012",
			RequestID:      "REQ202604110202",
			ConversationNo: "ACV202604110202",
			RunNo:          "ARN202604110202",
		},
	})
	if err != nil {
		t.Fatalf("runner returned error: %v", err)
	}

	if repo.lastFilter.OrderNo != "ORD202604110202" {
		t.Fatalf("expected model planner to drive logistics lookup for selected order, got %q", repo.lastFilter.OrderNo)
	}
	if out.ReplyPayload.IntentCode != "LOGISTICS" {
		t.Fatalf("expected logistics intent from model planner, got %q", out.ReplyPayload.IntentCode)
	}
}

func TestRunnerRunUsesModelPlannerPolicyKnowledgeDecision(t *testing.T) {
	knowledgeRepo := &stubPolicyKnowledgeRepository{
		hits: []PolicyKnowledgeHit{{
			SourceTypeCode: "PLATFORM_RULE",
			SourceID:       "refund-policy",
			SourceVersion:  3,
			Title:          "退款规则",
			Snippet:        "未发货订单通常可以直接申请退款；已发货订单需要结合物流与售后状态判断。",
		}},
	}
	runner := NewRunnerWithOptions(RunnerOptions{
		PolicyKnowledgeRepository: knowledgeRepo,
		ModelClient: &stubModelClient{
			responses: []string{`{"thought":"用户在问退款规则","intent":"POLICY_QA","user_goal":"了解退款规则","confidence":0.97,"slot_values":{},"missing_slots":[],"next_action":{"type":"TOOL_CALL","tool_name":"search_knowledge_chunks","tool_args":{"query":"退款规则"},"reason":"先查规则知识"},"handoff_recommended":false,"handoff_reason_code":""}`},
		},
	})

	out, err := runner.Run(context.Background(), RunInput{
		UserQuery: "退款规则是什么",
		Message: UserMessage{
			ContentText: "退款规则是什么",
		},
		Security: SecurityContext{
			UserID:         31021,
			ShopNo:         "SHOP31021",
			RequestID:      "REQ202604110301",
			ConversationNo: "ACV202604110301",
			RunNo:          "ARN202604110301",
		},
	})
	if err != nil {
		t.Fatalf("runner returned error: %v", err)
	}

	if knowledgeRepo.lastQuery != "退款规则" {
		t.Fatalf("expected planner tool args to drive policy search query, got %q", knowledgeRepo.lastQuery)
	}
	if out.ReplyPayload.IntentCode != officialIntentPolicyQA {
		t.Fatalf("expected policy intent from model planner, got %q", out.ReplyPayload.IntentCode)
	}
	if len(out.AnswerSources) == 0 || out.AnswerSources[0].SourceID != "refund-policy" {
		t.Fatalf("expected policy answer sources to be attached, got %+v", out.AnswerSources)
	}
	if strings.Contains(out.ReplyPayload.ReplyText, "search_knowledge_chunks") {
		t.Fatalf("expected user-facing reply text without tool names, got %q", out.ReplyPayload.ReplyText)
	}
}

func TestRunnerRunUsesModelPlannerCatalogDecision(t *testing.T) {
	productRepo := &stubProductSearchRepository{
		items: []ProductSearchItem{
			{
				SpuNo:      "SPU202604110001",
				Title:      "轻薄通勤风衣",
				CoverURL:   "https://cdn.shopa.test/spu1.jpg",
				MinPrice:   19900,
				MaxPrice:   25900,
				ShopName:   "轻简女装店",
				ReasonText: "轻薄面料更适合夏天通勤场景。",
			},
			{
				SpuNo:      "SPU202604110002",
				Title:      "短款防晒外套",
				CoverURL:   "https://cdn.shopa.test/spu2.jpg",
				MinPrice:   12900,
				MaxPrice:   15900,
				ShopName:   "城市衣橱",
				ReasonText: "更偏日常通勤，适合空调房和早晚通勤。",
			},
		},
	}
	runner := NewRunnerWithOptions(RunnerOptions{
		ProductSearchRepository: productRepo,
		ModelClient: &stubModelClient{
			responses: []string{`{"thought":"用户想找适合夏天通勤的外套","intent":"CATALOG_GUIDE","user_goal":"推荐适合夏天通勤的外套","confidence":0.96,"slot_values":{},"missing_slots":[],"next_action":{"type":"TOOL_CALL","tool_name":"search_products","tool_args":{"query":"夏天通勤外套"},"reason":"先搜索商品"},"handoff_recommended":false,"handoff_reason_code":""}`},
		},
	})

	out, err := runner.Run(context.Background(), RunInput{
		UserQuery: "推荐一款适合夏天通勤的外套",
		Message: UserMessage{
			ContentText: "推荐一款适合夏天通勤的外套",
		},
		Security: SecurityContext{
			UserID:         31022,
			ShopNo:         "SHOP31022",
			RequestID:      "REQ202604110302",
			ConversationNo: "ACV202604110302",
			RunNo:          "ARN202604110302",
		},
	})
	if err != nil {
		t.Fatalf("runner returned error: %v", err)
	}

	if productRepo.lastFilter.Query != "夏天通勤外套" {
		t.Fatalf("expected planner tool args to drive product search query, got %q", productRepo.lastFilter.Query)
	}
	if out.ReplyPayload.IntentCode != officialIntentCatalogGuide {
		t.Fatalf("expected catalog intent from model planner, got %q", out.ReplyPayload.IntentCode)
	}
	if len(out.ReplyPayload.DataCards) == 0 || out.ReplyPayload.DataCards[0].ProductRecommendationCard == nil {
		t.Fatalf("expected product recommendation card from catalog guide reply")
	}
	if got := len(out.ReplyPayload.DataCards[0].ProductRecommendationCard.Items); got != 2 {
		t.Fatalf("expected 2 recommended products, got %d", got)
	}
	if strings.Contains(out.ReplyPayload.ReplyText, "search_products") {
		t.Fatalf("expected user-facing reply text without tool names, got %q", out.ReplyPayload.ReplyText)
	}
}

func TestRunnerRunUsesSpecificOrderStatusReplyForGenericOrderQuestion(t *testing.T) {
	repo := &stubOrderSnapshotRepository{
		snapshot: &OrderSnapshot{
			OrderNo:            "ORD202604110303",
			MainStatus:         "PAID",
			PaymentStatus:      "PAID",
			FulfillmentStatus:  "SHIPPED",
			LogisticsStatus:    "IN_TRANSIT",
			AfterSaleStatus:    "NONE",
			LatestUpdateTime:   "2026-04-11T11:30:00+08:00",
			OwnershipConfirmed: true,
		},
	}
	runner := NewRunnerWithOptions(RunnerOptions{
		OrderRepository: repo,
		ModelClient: &stubModelClient{
			responses: []string{`{"thought":"用户想了解这笔订单信息","intent":"AFTER_SALE","user_goal":"同步订单当前情况","confidence":0.93,"slot_values":{},"missing_slots":[],"next_action":{"type":"TOOL_CALL","tool_name":"get_order_snapshot","tool_args":{},"reason":"先查订单事实"},"handoff_recommended":false,"handoff_reason_code":""}`},
		},
	})

	out, err := runner.Run(context.Background(), RunInput{
		UserQuery: "帮我查询一下这个订单的信息",
		Message: UserMessage{
			ContentText: "帮我查询一下这个订单的信息",
		},
		HiddenAction: HiddenAction{
			Type:  "SET_SLOT",
			Key:   "selected_order_no",
			Value: "ORD202604110303",
		},
		Security: SecurityContext{
			UserID:         31023,
			ShopNo:         "SHOP31023",
			RequestID:      "REQ202604110303",
			ConversationNo: "ACV202604110303",
			RunNo:          "ARN202604110303",
		},
	})
	if err != nil {
		t.Fatalf("runner returned error: %v", err)
	}

	if repo.lastFilter.OrderNo != "ORD202604110303" {
		t.Fatalf("expected selected order to be queried, got %q", repo.lastFilter.OrderNo)
	}
	if len(out.ReplyPayload.DataCards) < 2 || out.ReplyPayload.DataCards[1].AfterSaleDecisionCard == nil {
		t.Fatalf("expected structured order decision card, got %+v", out.ReplyPayload.DataCards)
	}
	if out.ReplyPayload.DataCards[1].AfterSaleDecisionCard.SceneCode != sceneCodeOrderStatusInTransit {
		t.Fatalf("expected in-transit order status scene, got %q with reply %q", out.ReplyPayload.DataCards[1].AfterSaleDecisionCard.SceneCode, out.ReplyPayload.ReplyText)
	}
	if strings.Contains(out.ReplyPayload.ReplyText, "还缺少更明确的场景线索") {
		t.Fatalf("expected a specific order-status reply instead of generic after-sale guidance, got %q", out.ReplyPayload.ReplyText)
	}
}
