package agent

import (
	"testing"

	"github.com/TsingpekTao/shopa/agent-svc/internal/model/entity"
)

func TestToProtoRunIncludesReplyPayloadFromGraphState(t *testing.T) {
	row := &entity.AgentRun{
		Id:    1,
		RunNo: "ARN202604090201",
		GraphStateJson: `{
			"reply_payload": {
				"reply_text": "订单尚未发货，建议直接申请退款。",
				"intent_code": "AfterSaleTask",
				"guard_result_code": "AFTER_SALE_TASK",
				"confidence": 0.92,
				"slot_retry_count": 0,
				"handoff_recommended": false,
				"data_cards": [
					{
						"order_snapshot_card": {
							"order_no": "ORD202604090201",
							"main_status": "PAID",
							"payment_status": "PAID",
							"fulfillment_status": "UNSHIPPED",
							"logistics_status": "NOT_SHIPPED",
							"after_sale_status": "NONE",
							"latest_update_time": "2026-04-09T12:00:00+08:00"
						}
					},
					{
						"after_sale_decision_card": {
							"decision_path_code": "refund_only",
							"reason_text": "订单尚未发货，当前更适合直接申请退款。",
							"constraint_text": "未发货阶段通常不需要先走退货流程。",
							"next_step_text": "优先发起退款。",
							"scene_code": "refund_before_shipment"
						}
					}
				],
				"suggested_actions": [
					{
						"action_code": "request_refund",
						"label": "申请退款",
						"enabled": true
					}
				]
			}
		}`,
	}

	got := toProtoRun(row)
	if got == nil {
		t.Fatalf("expected proto run to be non-nil")
	}
	if got.GetReplyPayload() == nil {
		t.Fatalf("expected reply payload to be parsed from graph state")
	}
	if got.GetReplyPayload().GetIntentCode() != "AfterSaleTask" {
		t.Fatalf("expected reply payload intent AfterSaleTask, got %q", got.GetReplyPayload().GetIntentCode())
	}
	if len(got.GetReplyPayload().GetDataCards()) != 2 {
		t.Fatalf("expected two data cards, got %d", len(got.GetReplyPayload().GetDataCards()))
	}
	if got.GetReplyPayload().GetDataCards()[0].GetOrderSnapshotCard().GetOrderNo() != "ORD202604090201" {
		t.Fatalf("expected order snapshot card to keep order_no, got %q", got.GetReplyPayload().GetDataCards()[0].GetOrderSnapshotCard().GetOrderNo())
	}
	if got.GetReplyPayload().GetDataCards()[1].GetAfterSaleDecisionCard().GetDecisionPathCode() != "refund_only" {
		t.Fatalf("expected decision path refund_only, got %q", got.GetReplyPayload().GetDataCards()[1].GetAfterSaleDecisionCard().GetDecisionPathCode())
	}
	if got.GetReplyPayload().GetDataCards()[1].GetAfterSaleDecisionCard().GetSceneCode() != "refund_before_shipment" {
		t.Fatalf("expected scene code refund_before_shipment, got %q", got.GetReplyPayload().GetDataCards()[1].GetAfterSaleDecisionCard().GetSceneCode())
	}
	if len(got.GetReplyPayload().GetSuggestedActions()) != 1 {
		t.Fatalf("expected one suggested action, got %d", len(got.GetReplyPayload().GetSuggestedActions()))
	}
	if got.GetReplyPayload().GetSuggestedActions()[0].GetActionCode() != "request_refund" {
		t.Fatalf("expected suggested action request_refund, got %q", got.GetReplyPayload().GetSuggestedActions()[0].GetActionCode())
	}
}

func TestToProtoRunIncludesOrderSelectionCardFromGraphState(t *testing.T) {
	row := &entity.AgentRun{
		Id:    2,
		RunNo: "ARN202604100301",
		GraphStateJson: `{
			"reply_payload": {
				"reply_text": "我先把你最近的订单拉出来，点一个我继续查。",
				"intent_code": "AfterSaleTask",
				"guard_result_code": "AFTER_SALE_TASK",
				"confidence": 0.9,
				"data_cards": [
					{
						"order_selection_card": {
							"title_text": "找到你最近的订单，点一个我继续处理",
							"helper_text": "我会沿着你刚才的问题继续查，不需要重新描述。",
							"original_query": "帮我查一下这个订单现在到哪了",
							"task_code": "logistics_query",
							"candidates": [
								{
									"order_no": "ORD202604100301",
									"display_title": "云感运动鞋",
									"fulfillment_status": "SHIPPED",
									"logistics_status": "IN_TRANSIT",
									"latest_update_time": "2026-04-10T09:30:00+08:00",
									"selection_hint": "运输中，适合继续查物流"
								}
							]
						}
					}
				]
			}
		}`,
	}

	got := toProtoRun(row)
	if got == nil || got.GetReplyPayload() == nil {
		t.Fatalf("expected reply payload to be parsed from graph state")
	}
	if len(got.GetReplyPayload().GetDataCards()) != 1 {
		t.Fatalf("expected one data card, got %d", len(got.GetReplyPayload().GetDataCards()))
	}
	card := got.GetReplyPayload().GetDataCards()[0].GetOrderSelectionCard()
	if card == nil {
		t.Fatalf("expected order selection card to be mapped")
	}
	if card.GetTaskCode() != "logistics_query" {
		t.Fatalf("expected task code logistics_query, got %q", card.GetTaskCode())
	}
	if len(card.GetCandidates()) != 1 {
		t.Fatalf("expected one order candidate, got %d", len(card.GetCandidates()))
	}
	if card.GetCandidates()[0].GetOrderNo() != "ORD202604100301" {
		t.Fatalf("expected candidate order no ORD202604100301, got %q", card.GetCandidates()[0].GetOrderNo())
	}
}

func TestToProtoRunIncludesLogisticsAndProductCardsFromGraphState(t *testing.T) {
	row := &entity.AgentRun{
		Id:    3,
		RunNo: "ARN202604101122",
		GraphStateJson: `{
			"reply_payload": {
				"reply_text": "订单正在派送中，我也给你补了两款相近商品。",
				"intent_code": "LOGISTICS",
				"guard_result_code": "AFTER_SALE_TASK",
				"confidence": 0.95,
				"data_cards": [
					{
						"logistics_tracking_card": {
							"order_no": "ORD202604101122",
							"sub_order_no": "SUB202604101122",
							"fulfillment_status": "SHIPPED",
							"logistics_status": "OUT_FOR_DELIVERY",
							"latest_update_time": "2026-04-10T16:00:00+08:00",
							"latest_trace_text": "快递员正在派送",
							"timeline_summary": "已发货，正在派送"
						}
					},
					{
						"product_recommendation_card": {
							"title_text": "你可能还想看看",
							"helper_text": "这些商品和你刚才问的风格接近。",
							"items": [
								{
									"spu_no": "SPU1001",
									"title": "春季风衣",
									"cover_url": "https://img.example.com/spu1001.png",
									"price_text": "¥199",
									"shop_name": "Shopa 官方店",
									"reason_text": "轻薄通勤，适合现在的天气"
								}
							]
						}
					}
				]
			}
		}`,
	}

	got := toProtoRun(row)
	if got == nil || got.GetReplyPayload() == nil {
		t.Fatalf("expected reply payload to be parsed from graph state")
	}
	if len(got.GetReplyPayload().GetDataCards()) != 2 {
		t.Fatalf("expected two data cards, got %d", len(got.GetReplyPayload().GetDataCards()))
	}
	logisticsCard := got.GetReplyPayload().GetDataCards()[0].GetLogisticsTrackingCard()
	if logisticsCard == nil {
		t.Fatalf("expected logistics tracking card to be mapped")
	}
	if logisticsCard.GetOrderNo() != "ORD202604101122" {
		t.Fatalf("expected logistics card order no ORD202604101122, got %q", logisticsCard.GetOrderNo())
	}
	productCard := got.GetReplyPayload().GetDataCards()[1].GetProductRecommendationCard()
	if productCard == nil {
		t.Fatalf("expected product recommendation card to be mapped")
	}
	if len(productCard.GetItems()) != 1 {
		t.Fatalf("expected one product recommendation item, got %d", len(productCard.GetItems()))
	}
	if productCard.GetItems()[0].GetSpuNo() != "SPU1001" {
		t.Fatalf("expected spu no SPU1001, got %q", productCard.GetItems()[0].GetSpuNo())
	}
}
