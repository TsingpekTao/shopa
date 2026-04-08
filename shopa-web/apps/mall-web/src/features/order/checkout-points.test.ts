import assert from "node:assert/strict";
import { buildCheckoutPointsSummary, buildOrderAmountBreakdown } from "./checkout-points.js";

function runBuildCheckoutPointsSummaryUsesPointsTest() {
  const summary = buildCheckoutPointsSummary({
    payableAmount: 29900,
    availablePoints: 2800,
    usePoints: true
  });

  assert.equal(summary.maxDiscountAmount, 1495);
  assert.equal(summary.maxUsablePoints, 1495);
  assert.equal(summary.selectedPoints, 1495);
  assert.equal(summary.selectedDiscountAmount, 1495);
  assert.equal(summary.finalPayableAmount, 28405);
  assert.equal(summary.canToggle, true);
}

function runBuildCheckoutPointsSummaryNoPointsTest() {
  const summary = buildCheckoutPointsSummary({
    payableAmount: 29900,
    availablePoints: 0,
    usePoints: false
  });

  assert.equal(summary.maxDiscountAmount, 1495);
  assert.equal(summary.maxUsablePoints, 0);
  assert.equal(summary.selectedDiscountAmount, 0);
  assert.equal(summary.finalPayableAmount, 29900);
  assert.equal(summary.canToggle, false);
}

function runBuildOrderAmountBreakdownTest() {
  const breakdown = buildOrderAmountBreakdown({
    goodsAmount: 29900,
    freightAmount: 0,
    discountAmount: 300,
    pointsDiscountAmount: 500,
    payableAmount: 29100,
    paidAmount: 29100
  });

  assert.equal(breakdown.originalPayableAmount, 29900);
  assert.equal(breakdown.totalDiscountAmount, 800);
}

runBuildCheckoutPointsSummaryUsesPointsTest();
runBuildCheckoutPointsSummaryNoPointsTest();
runBuildOrderAmountBreakdownTest();
console.log("checkout points tests passed");
