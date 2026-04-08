import assert from "node:assert/strict";
import { canApplyRefundForSubOrder, getBuyerRefundStatusMeta } from "./status.js";

function runGetBuyerRefundStatusMetaTest() {
  assert.deepEqual(getBuyerRefundStatusMeta("AFTER_SALE_STATUS_PENDING_SELLER_REVIEW"), {
    label: "待商家审核",
    tone: "warning"
  });
  assert.deepEqual(getBuyerRefundStatusMeta("REFUNDED"), {
    label: "已退款",
    tone: "success"
  });
  assert.deepEqual(getBuyerRefundStatusMeta(""), {
    label: "处理中",
    tone: "default"
  });
}

function runCanApplyRefundForSubOrderTest() {
  assert.equal(
    canApplyRefundForSubOrder({
      subStatus: "WAIT_SHIP",
      items: [{ itemNo: "ITEM1" }]
    } as never),
    true
  );
  assert.equal(
    canApplyRefundForSubOrder({
      subStatus: "PAID",
      items: [{ itemNo: "ITEM1" }]
    } as never),
    true
  );
  assert.equal(
    canApplyRefundForSubOrder({
      subStatus: "SHIPPED",
      items: [{ itemNo: "ITEM1" }]
    } as never),
    false
  );
  assert.equal(
    canApplyRefundForSubOrder({
      subStatus: "WAIT_SHIP",
      items: []
    } as never),
    false
  );
}

runGetBuyerRefundStatusMetaTest();
runCanApplyRefundForSubOrderTest();
console.log("buyer refund status tests passed");
