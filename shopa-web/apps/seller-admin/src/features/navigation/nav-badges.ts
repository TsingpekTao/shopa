import type { SellerChatConversation } from "../chat/types";
import type { SellerOrderSub } from "../order/types";

export const SELLER_CURRENT_SHOP_KEY = "seller-current-shop-no";
export const SELLER_CURRENT_SHOP_EVENT = "seller-current-shop-change";

export function resolveSellerCurrentShopNo(availableShopNos: string[], preferredShopNo: string): string {
  const normalizedPreferredShopNo = preferredShopNo.trim();
  const normalizedAvailableShopNos = availableShopNos.map((shopNo) => shopNo.trim()).filter(Boolean);
  if (normalizedPreferredShopNo && normalizedAvailableShopNos.includes(normalizedPreferredShopNo)) {
    return normalizedPreferredShopNo;
  }
  return normalizedAvailableShopNos[0] ?? "";
}

export function formatSellerNavBadgeCount(count: number): string {
  if (!Number.isFinite(count) || count <= 0) {
    return "";
  }
  return count > 99 ? "99+" : String(Math.trunc(count));
}

export function getSellerCurrentShopNo(): string {
  if (typeof window === "undefined") {
    return "";
  }
  return window.localStorage.getItem(SELLER_CURRENT_SHOP_KEY)?.trim() ?? "";
}

export function setSellerCurrentShopNo(shopNo: string): void {
  if (typeof window === "undefined") {
    return;
  }
  const normalizedShopNo = shopNo.trim();
  if (normalizedShopNo) {
    window.localStorage.setItem(SELLER_CURRENT_SHOP_KEY, normalizedShopNo);
  } else {
    window.localStorage.removeItem(SELLER_CURRENT_SHOP_KEY);
  }
  window.dispatchEvent(new CustomEvent(SELLER_CURRENT_SHOP_EVENT, { detail: normalizedShopNo }));
}

export function sumSellerConversationUnread(conversations: SellerChatConversation[]): number {
  return conversations.reduce((sum, conversation) => sum + Math.max(0, Number(conversation.unreadCount ?? 0)), 0);
}

export function countSellerPendingShipments(subOrders: SellerOrderSub[]): number {
  return subOrders.reduce((sum, subOrder) => {
    return subOrder.subStatus === "PAID" || subOrder.subStatus === "WAIT_SHIP" ? sum + 1 : sum;
  }, 0);
}
