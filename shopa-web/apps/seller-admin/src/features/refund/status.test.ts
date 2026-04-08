import assert from "node:assert/strict";
import { getSellerRefundStatusMeta, sortRefundBatchesForSeller } from "./status.js";

function runGetSellerRefundStatusMetaTest() {
  assert.deepEqual(getSellerRefundStatusMeta("AFTER_SALE_STATUS_PENDING_SELLER_REVIEW"), {
    label: "待审核",
    tone: "warning"
  });
  assert.deepEqual(getSellerRefundStatusMeta("WAIT_REFUND_TASK"), {
    label: "待退款",
    tone: "processing"
  });
  assert.deepEqual(getSellerRefundStatusMeta("REFUNDED"), {
    label: "已退款",
    tone: "success"
  });
}

function runSortRefundBatchesForSellerTest() {
  const sorted = sortRefundBatchesForSeller([
    { refundBatchNo: "RB3", batchStatus: "REFUNDED", createdAt: "2026-04-08T10:00:00.000Z" },
    { refundBatchNo: "RB1", batchStatus: "PENDING_SELLER_REVIEW", createdAt: "2026-04-08T09:00:00.000Z" },
    { refundBatchNo: "RB2", batchStatus: "WAIT_REFUND_TASK", createdAt: "2026-04-08T11:00:00.000Z" }
  ] as never);

  assert.deepEqual(
    sorted.map((item) => item.refundBatchNo),
    ["RB1", "RB2", "RB3"]
  );
}

runGetSellerRefundStatusMetaTest();
runSortRefundBatchesForSellerTest();
console.log("seller refund status tests passed");
