import assert from "node:assert/strict";
import {
  formatSellerProductPriceRange,
  isSellerDraftBoxStatus,
  resolveSellerProductsScope
} from "./products-view.js";

function runFormatSellerProductPriceRangeTest() {
  assert.equal(formatSellerProductPriceRange(0, 0), "-");
  assert.equal(formatSellerProductPriceRange(9900, 9900), "¥99.00");
  assert.equal(formatSellerProductPriceRange(9900, 12900), "¥99.00 ~ ¥129.00");
}

function runIsSellerDraftBoxStatusTest() {
  assert.equal(isSellerDraftBoxStatus("SPU_STATUS_DRAFT"), true);
  assert.equal(isSellerDraftBoxStatus("SPU_STATUS_REJECTED"), false);
  assert.equal(isSellerDraftBoxStatus("SPU_STATUS_REVIEWING"), false);
}

function runResolveSellerProductsScopeTest() {
  assert.equal(resolveSellerProductsScope("drafts"), "drafts");
  assert.equal(resolveSellerProductsScope("all"), "all");
  assert.equal(resolveSellerProductsScope(" anything "), "all");
}

runFormatSellerProductPriceRangeTest();
runIsSellerDraftBoxStatusTest();
runResolveSellerProductsScopeTest();
console.log("seller products view tests passed");
