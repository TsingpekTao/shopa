import type { BuyerOrder, BuyerOrderItem, BuyerOrderSub } from "./types";

function firstNonEmptySubOrder(order: BuyerOrder): BuyerOrderSub | null {
  return (
    order.subOrders.find((sub) => sub.shopNo.trim() || sub.items.some((item) => item.shopNo.trim())) ??
    order.subOrders.find((sub) => sub.subOrderNo.trim() || sub.items.length > 0) ??
    order.subOrders[0] ??
    null
  );
}

function firstNonEmptyItem(order: BuyerOrder, preferredSubOrder: BuyerOrderSub | null): BuyerOrderItem | null {
  const preferredItem = preferredSubOrder?.items.find((item) => item.shopNo.trim() || item.spuNo.trim() || item.skuNo.trim());
  if (preferredItem) {
    return preferredItem;
  }
  return order.subOrders.flatMap((sub) => sub.items).find((item) => item.shopNo.trim() || item.spuNo.trim() || item.skuNo.trim()) ?? null;
}

export function resolveAfterSaleConsultParams(order: BuyerOrder): {
  shopNo: string;
  subOrderNo: string;
  spuNo: string;
  skuNo: string;
} | null {
  const preferredSubOrder = firstNonEmptySubOrder(order);
  const preferredItem = firstNonEmptyItem(order, preferredSubOrder);
  const shopNo =
    preferredItem?.shopNo.trim() ||
    preferredSubOrder?.shopNo.trim() ||
    order.subOrders.find((sub) => sub.shopNo.trim())?.shopNo.trim() ||
    "";
  if (!shopNo) {
    return null;
  }
  return {
    shopNo,
    subOrderNo:
      preferredItem?.subOrderNo.trim() ||
      preferredSubOrder?.subOrderNo.trim() ||
      order.subOrders.find((sub) => sub.subOrderNo.trim())?.subOrderNo.trim() ||
      "",
    spuNo: preferredItem?.spuNo.trim() || order.subOrders.flatMap((sub) => sub.items).find((item) => item.spuNo.trim())?.spuNo.trim() || "",
    skuNo: preferredItem?.skuNo.trim() || order.subOrders.flatMap((sub) => sub.items).find((item) => item.skuNo.trim())?.skuNo.trim() || ""
  };
}

export function buildAfterSaleConsultHref(order: BuyerOrder): string {
  const params = resolveAfterSaleConsultParams(order);
  if (!params) {
    return "";
  }
  return `/me/messages?shop_no=${encodeURIComponent(params.shopNo)}&spu_no=${encodeURIComponent(params.spuNo)}&sku_no=${encodeURIComponent(params.skuNo)}&order_no=${encodeURIComponent(order.orderNo)}&sub_order_no=${encodeURIComponent(params.subOrderNo)}&scene_code=AFTER_SALE`;
}
