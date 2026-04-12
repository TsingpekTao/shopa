import assert from "node:assert/strict";
import { test } from "vitest";
import { apiClient } from "../../lib/api-client.js";
import { getAssistantRunStatus, sendAssistantMessage } from "./api.js";

test("getAssistantRunStatus sanitizes object-like text fields", async () => {
  const originalGet = apiClient.get;

  try {
    (
      apiClient as typeof apiClient & {
        get: typeof apiClient.get;
      }
    ).get = (async () =>
      ({
        conversation: {
          conversation_no: "ACV202604101700000001",
        },
        run: {
          run_no: "ARN202604101700000001",
          run_status_code: "SUCCESS",
        },
        latest_assistant_message: {
          message_no: "AMSG202604101700000001",
          sender_type_code: "AGENT",
          content_text: { rich: "unexpected-object" },
          sent_at: undefined,
        },
        answer_sources: [
          {
            source_type_code: "ORDER_CONTEXT",
            source_id: "shipping-faq",
            source_version: 1,
            title: { text: "发货与物流 FAQ" },
            snippet: { text: "待发货先看订单状态" },
          },
        ],
      }) as never) as typeof apiClient.get;

    const result = await getAssistantRunStatus("ACV202604101700000001");

    assert.equal(result.latestAssistantMessage?.contentText, "");
    assert.equal(result.latestAssistantMessage?.sentAt, "");
    assert.equal(result.answerSources[0]?.title, "");
    assert.equal(result.answerSources[0]?.snippet, "");
  } finally {
    (
      apiClient as typeof apiClient & {
        get: typeof apiClient.get;
      }
    ).get = originalGet;
  }
});

test("sendAssistantMessage forwards hidden action for slot setting", async () => {
  const originalPost = apiClient.post;
  let capturedBody: Record<string, unknown> | undefined;

  try {
    (
      apiClient as typeof apiClient & {
        post: typeof apiClient.post;
      }
    ).post = (async (_url: string, body?: unknown) => {
      capturedBody = body as Record<string, unknown>;
      return {
        conversation: {
          conversation_no: "ACV202604102359000001",
        },
        run: {
          run_no: "ARN202604102359000001",
          run_status_code: "PROCESSING",
        },
        user_message: {
          message_no: "AMSG202604102359000001",
          sender_type_code: "BUYER",
          content_text: "查询这笔订单",
        },
      } as never;
    }) as typeof apiClient.post;

    await sendAssistantMessage({
      conversationNo: "ACV202604102359000001",
      contentText: "查询这笔订单",
      hiddenAction: {
        type: "SET_SLOT",
        key: "selected_order_no",
        value: "ORD202604100001",
      },
    });

    assert.equal(capturedBody?.conversation_no, "ACV202604102359000001");
    assert.deepEqual(capturedBody?.hidden_action, {
      type: "SET_SLOT",
      key: "selected_order_no",
      value: "ORD202604100001",
    });
  } finally {
    (
      apiClient as typeof apiClient & {
        post: typeof apiClient.post;
      }
    ).post = originalPost;
  }
});

test("getAssistantRunStatus normalizes logistics and product cards", async () => {
  const originalGet = apiClient.get;

  try {
    (
      apiClient as typeof apiClient & {
        get: typeof apiClient.get;
      }
    ).get = (async () =>
      ({
        conversation: {
          conversation_no: "ACV202604110001",
        },
        run: {
          run_no: "ARN202604110001",
          run_status_code: "SUCCESS",
          reply_payload: {
            reply_text: "订单正在派送中。",
            intent_code: "LOGISTICS",
            data_cards: [
              {
                logistics_tracking_card: {
                  order_no: "ORD202604110001",
                  logistics_status: "OUT_FOR_DELIVERY",
                  latest_trace_text: "快递员正在派送",
                },
              },
              {
                product_recommendation_card: {
                  title_text: "你可能还想看看",
                  items: [
                    {
                      spu_no: "SPU202604110001",
                      title: "春季风衣",
                      price_text: "¥199",
                    },
                  ],
                },
              },
            ],
          },
        },
      }) as never) as typeof apiClient.get;

    const result = await getAssistantRunStatus("ACV202604110001");

    assert.equal(
      result.run.replyPayload?.dataCards[0]?.logisticsTrackingCard?.orderNo,
      "ORD202604110001",
    );
    assert.equal(
      result.run.replyPayload?.dataCards[0]?.logisticsTrackingCard
        ?.latestTraceText,
      "快递员正在派送",
    );
    assert.equal(
      result.run.replyPayload?.dataCards[1]?.productRecommendationCard?.items[0]
        ?.spuNo,
      "SPU202604110001",
    );
  } finally {
    (
      apiClient as typeof apiClient & {
        get: typeof apiClient.get;
      }
    ).get = originalGet;
  }
});
