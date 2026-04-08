import assert from "node:assert/strict";
import { buildSellerOrderDetailPath, buildSellerOrderListPath, buildSellerSubOrderShipPath } from "./paths.js";

function runBuildSellerOrderDetailPathTest() {
  assert.equal(buildSellerOrderListPath(), "/v1/order/seller/list");
  assert.equal(buildSellerOrderDetailPath("SUB20260407192826493100"), "/v1/order/seller/sub/SUB20260407192826493100");
  assert.equal(buildSellerOrderDetailPath(" SUB/1001 "), "/v1/order/seller/sub/SUB%2F1001");
  assert.equal(buildSellerSubOrderShipPath(), "/v1/order/seller/sub/ship");
}

runBuildSellerOrderDetailPathTest();
console.log("seller-order api tests passed");
