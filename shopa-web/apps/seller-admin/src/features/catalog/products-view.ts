import { SpuStatusCode } from "./types";

export type SellerProductsScope = "all" | "drafts";

function toCurrency(amountInMinorUnits: number): string {
  return `¥${(amountInMinorUnits / 100).toFixed(2)}`;
}

export function formatSellerProductPriceRange(minSalePrice: number, maxSalePrice: number): string {
  if (minSalePrice <= 0 && maxSalePrice <= 0) {
    return "-";
  }
  if (minSalePrice === maxSalePrice || maxSalePrice <= 0) {
    return toCurrency(minSalePrice);
  }
  return `${toCurrency(minSalePrice)} ~ ${toCurrency(maxSalePrice)}`;
}

export function isSellerDraftBoxStatus(status: SpuStatusCode): boolean {
  return status === "SPU_STATUS_DRAFT";
}

export function resolveSellerProductsScope(value?: string | null): SellerProductsScope {
  return value?.trim().toLowerCase() === "drafts" ? "drafts" : "all";
}
