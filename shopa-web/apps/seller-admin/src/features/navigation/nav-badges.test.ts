import assert from "node:assert/strict";
import {
  SELLER_CURRENT_SHOP_KEY,
  countSellerPendingShipments,
  formatSellerNavBadgeCount,
  resolveSellerCurrentShopNo,
  sumSellerConversationUnread
} from "./nav-badges.js";

function runResolveSellerCurrentShopNoTest() {
  assert.equal(SELLER_CURRENT_SHOP_KEY, "seller-current-shop-no");
  assert.equal(resolveSellerCurrentShopNo(["SHOP-2", "SHOP-1"], "SHOP-1"), "SHOP-1");
  assert.equal(resolveSellerCurrentShopNo(["SHOP-2", "SHOP-1"], "SHOP-3"), "SHOP-2");
  assert.equal(resolveSellerCurrentShopNo([], "SHOP-1"), "");
}

function runFormatSellerNavBadgeCountTest() {
  assert.equal(formatSellerNavBadgeCount(0), "");
  assert.equal(formatSellerNavBadgeCount(7), "7");
  assert.equal(formatSellerNavBadgeCount(120), "99+");
}

function runSumSellerConversationUnreadTest() {
  assert.equal(
    sumSellerConversationUnread([
      { unreadCount: 2 } as never,
      { unreadCount: 0 } as never,
      { unreadCount: 5 } as never
    ]),
    7
  );
}

function runCountSellerPendingShipmentsTest() {
  assert.equal(
    countSellerPendingShipments([
      { subStatus: "PAID" } as never,
      { subStatus: "WAIT_SHIP" } as never,
      { subStatus: "SHIPPED" } as never,
      { subStatus: "CLOSED" } as never
    ]),
    2
  );
}

runResolveSellerCurrentShopNoTest();
runFormatSellerNavBadgeCountTest();
runSumSellerConversationUnreadTest();
runCountSellerPendingShipmentsTest();
console.log("seller nav badge tests passed");
