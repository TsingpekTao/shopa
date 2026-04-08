import assert from "node:assert/strict";
import {
  buildSellerSalesAnalyticsPath,
  formatSellerSalesRefundRate,
  getSellerSalesMetricValue,
  getSellerStatDisplayValue
} from "./sales-helpers.js";

function runBuildSellerSalesAnalyticsPathTest() {
  assert.equal(buildSellerSalesAnalyticsPath({ range: "7D" }), "/v1/seller/sales/analytics?range=7D");
  assert.equal(
    buildSellerSalesAnalyticsPath({ range: "30D", shopNo: " SHOP/1001 " }),
    "/v1/seller/sales/analytics?range=30D&shopNo=SHOP%2F1001"
  );
}

function runSellerSalesHelpersTest() {
  assert.equal(getSellerSalesMetricValue({ gmv: 12800, paidOrderCount: 9 }, "GMV"), 12800);
  assert.equal(getSellerSalesMetricValue({ gmv: 12800, paidOrderCount: 9 }, "ORDERS"), 9);
  assert.equal(formatSellerSalesRefundRate(0.0525), "5.25%");
  assert.equal(getSellerStatDisplayValue(0, true), "--");
  assert.equal(getSellerStatDisplayValue(0, false), 0);
}

runBuildSellerSalesAnalyticsPathTest();
runSellerSalesHelpersTest();
console.log("seller sales helper tests passed");
