import type { BuyerOrderItem } from "./types";

function toString(value: unknown): string {
  if (value === null || value === undefined) {
    return "";
  }
  if (typeof value === "string") {
    return value;
  }
  if (typeof value === "number" || typeof value === "boolean" || typeof value === "bigint") {
    return String(value);
  }
  return "";
}

function toNumber(value: unknown): number {
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : 0;
}

export function normalizeOrderItem(raw?: Record<string, unknown>): BuyerOrderItem {
  return {
    itemNo: toString(raw?.item_no ?? raw?.itemNo),
    orderNo: toString(raw?.order_no ?? raw?.orderNo),
    subOrderNo: toString(raw?.sub_order_no ?? raw?.subOrderNo),
    shopNo: toString(raw?.shop_no ?? raw?.shopNo),
    spuNo: toString(raw?.spu_no ?? raw?.spuNo),
    skuNo: toString(raw?.sku_no ?? raw?.skuNo),
    spuTitle: toString(raw?.spu_title ?? raw?.spuTitle),
    skuName: toString(raw?.sku_name ?? raw?.skuName),
    skuImageAssetId: toString(raw?.sku_image_asset_id ?? raw?.skuImageAssetId),
    qty: toNumber(raw?.qty),
    salePrice: toNumber(raw?.sale_price ?? raw?.salePrice),
    marketPrice: toNumber(raw?.market_price ?? raw?.marketPrice),
    saleAttrsJson: toString(raw?.sale_attrs_json ?? raw?.saleAttrsJson)
  };
}
