import assert from "node:assert/strict";
import { test } from "vitest";
import {
  buildActionPrompt,
  buildOrderSelectionPrompt,
  getIntroMessageText,
  getVariantCopy,
} from "./assistant-copy.js";

test("official assistant copy uses buyer-facing Chinese wording", () => {
  const copy = getVariantCopy("official", true);

  assert.equal(copy.chip, "Shopa 官方客服");
  assert.equal(copy.title, "官方客服");
  assert.equal(
    copy.description,
    "订单、物流、退款退货问题都可以直接聊，我会先理解你的问题，再帮你查订单、看物流或给出下一步建议。",
  );
  assert.deepEqual(copy.quickPrompts, [
    "帮我查一下这个订单现在到哪了",
    "这个订单为什么还没发货",
    "这个订单现在能退款吗",
    "我要转人工客服",
  ]);
  assert.equal(
    getIntroMessageText("official", true),
    "您好，请问我能帮助您什么？我可以帮您查询订单状态、物流进度、退款退货建议，也可以为您转接人工客服。",
  );
});

test("assistant prompt builders return clean Chinese task text", () => {
  assert.equal(
    buildActionPrompt(
      {
        actionCode: "view_logistics",
        label: "查看物流",
        enabled: true,
        reasonIfDisabled: "",
      },
      true,
    ),
    "帮我查一下这个订单的物流进度",
  );

  assert.equal(
    buildOrderSelectionPrompt(
      {
        titleText: "",
        helperText: "",
        originalQuery: "",
        taskCode: "logistics_query",
        candidates: [],
      },
      {
        orderNo: "ORD123",
        subOrderNo: "",
        displayTitle: "",
        mainStatus: "",
        paymentStatus: "",
        fulfillmentStatus: "",
        logisticsStatus: "",
        afterSaleStatus: "",
        latestUpdateTime: "",
        selectionHint: "",
      },
      true,
    ),
    "帮我查一下订单 ORD123 现在到哪了",
  );
});
