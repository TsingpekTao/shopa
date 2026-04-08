import assert from "node:assert/strict";
import { buildGlobalSearchHref, normalizeGlobalSearchTab } from "./search-form.js";

function runNormalizeGlobalSearchTabTest() {
  assert.equal(normalizeGlobalSearchTab("shop"), "shop");
  assert.equal(normalizeGlobalSearchTab("item"), "item");
  assert.equal(normalizeGlobalSearchTab(""), "item");
  assert.equal(normalizeGlobalSearchTab("unknown"), "item");
}

function runBuildGlobalSearchHrefTest() {
  assert.equal(buildGlobalSearchHref("羽绒服", "item"), "/search?q=%E7%BE%BD%E7%BB%92%E6%9C%8D&tab=item");
  assert.equal(buildGlobalSearchHref("店铺A", "shop"), "/search?q=%E5%BA%97%E9%93%BAA&tab=shop");
}

runNormalizeGlobalSearchTabTest();
runBuildGlobalSearchHrefTest();
console.log("AppShell search tests passed");
