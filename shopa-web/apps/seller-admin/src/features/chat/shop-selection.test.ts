import assert from "node:assert/strict";
import { buildSellerChatShopOptions } from "./shop-selection.js";

function runBuildSellerChatShopOptionsTest() {
  assert.deepEqual(
    buildSellerChatShopOptions({
      workbenchShops: [
        { shopNo: " SHOP1001 ", shopName: "Alpha Shop", shopStatusCode: "ACTIVE" },
        { shopNo: "", shopName: "Ignored" }
      ],
      catalogShopNos: ["SHOP1001", "SHOP1002", " "],
      persistedShopNo: "SHOP1003"
    }),
    [
      { shopNo: "SHOP1001", shopName: "Alpha Shop", shopStatusCode: "ACTIVE" },
      { shopNo: "SHOP1002", shopName: "SHOP1002", shopStatusCode: "" },
      { shopNo: "SHOP1003", shopName: "SHOP1003", shopStatusCode: "" }
    ]
  );
}

runBuildSellerChatShopOptionsTest();
console.log("shop-selection tests passed");
