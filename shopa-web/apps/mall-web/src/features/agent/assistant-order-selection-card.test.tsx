import React from "react";
import { strict as assert } from "node:assert";
import { renderToStaticMarkup } from "react-dom/server";
import { test } from "vitest";
import { AssistantOrderSelectionCardView } from "./assistant-order-selection-card";

test("assistant order selection card renders selectable recent orders in chat", () => {
  const html = renderToStaticMarkup(
    <AssistantOrderSelectionCardView
      isZh
      onSelect={() => {}}
      card={{
        titleText: "找到你最近的订单了，点一个我继续处理",
        helperText: "我会沿着你刚才的问题继续往下查，不需要重新描述。",
        originalQuery: "帮我查一下这个订单现在到哪了",
        taskCode: "logistics_query",
        candidates: [
          {
            orderNo: "ORD1001",
            subOrderNo: "",
            displayTitle: "云感运动鞋",
            mainStatus: "PAID",
            paymentStatus: "PAID",
            fulfillmentStatus: "SHIPPED",
            logisticsStatus: "IN_TRANSIT",
            afterSaleStatus: "NONE",
            latestUpdateTime: "2026-04-10 13:30:00",
            selectionHint: "这单正在运输中，适合继续查看物流进度。",
          },
        ],
      }}
    />,
  );

  assert.match(html, /找到你最近的订单了/);
  assert.match(html, /云感运动鞋/);
  assert.match(html, /ORD1001/);
  assert.match(html, /这单正在运输中/);
  assert.match(html, /button/);
});
