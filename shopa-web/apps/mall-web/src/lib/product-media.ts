import { BuyerProductDetail } from "@/features/catalog/types";

export type ProductDetailMap = Record<string, BuyerProductDetail>;

type ProductImageTarget = {
  spuNo?: string;
  skuNo?: string;
  skuImageAssetId?: string;
};

function pushUniqueAssetId(assetIds: string[], seen: Set<string>, raw: string | undefined) {
  const assetId = (raw ?? "").trim();
  if (!assetId || seen.has(assetId)) {
    return;
  }
  seen.add(assetId);
  assetIds.push(assetId);
}

function normalizeProductImageUrl(raw: string | undefined): string {
  const value = (raw ?? "").trim();
  if (!value) {
    return "";
  }
  if (value.startsWith("data:") || value.startsWith("blob:")) {
    return value;
  }
  try {
    const parsed = new URL(value);
    parsed.pathname = parsed.pathname
      .split("/")
      .map((segment, index) => (index === 0 ? segment : normalizeUrlSegment(segment)))
      .join("/");
    return parsed.toString();
  } catch {
    return encodeURI(value);
  }
}

function normalizeUrlSegment(segment: string): string {
  if (!segment) {
    return segment;
  }
  try {
    return encodeURIComponent(decodeURIComponent(segment));
  } catch {
    return encodeURIComponent(segment);
  }
}

export function collectProductDetailAssetIds(detailMap: ProductDetailMap): string[] {
  const assetIds: string[] = [];
  const seen = new Set<string>();

  Object.values(detailMap).forEach((detail) => {
    detail.mainImageAssetIds.forEach((assetId) => {
      pushUniqueAssetId(assetIds, seen, assetId);
    });
    detail.skus.forEach((sku) => {
      pushUniqueAssetId(assetIds, seen, sku.skuImageAssetId);
    });
  });

  return assetIds;
}

export function resolveDetailAssetUrl(
  item: ProductImageTarget,
  detailMap: ProductDetailMap,
  detailAssetUrlMap: Record<string, string>
): string {
  const spuNo = item.spuNo?.trim() ?? "";
  if (!spuNo) {
    return "";
  }

  const detail = detailMap[spuNo];
  if (!detail) {
    return "";
  }

  const candidateAssetIds: string[] = [];
  const seen = new Set<string>();
  const matchedSku = detail.skus.find((sku) => sku.skuNo === item.skuNo);

  pushUniqueAssetId(candidateAssetIds, seen, matchedSku?.skuImageAssetId);
  detail.mainImageAssetIds.forEach((assetId) => {
    pushUniqueAssetId(candidateAssetIds, seen, assetId);
  });
  detail.skus.forEach((sku) => {
    pushUniqueAssetId(candidateAssetIds, seen, sku.skuImageAssetId);
  });

  for (const assetId of candidateAssetIds) {
    const url = normalizeProductImageUrl(detailAssetUrlMap[assetId]);
    if (url) {
      return url;
    }
  }

  return "";
}

export function resolveProductImageUrl(
  item: ProductImageTarget,
  assetUrlMap: Record<string, string>,
  productImageMap: Record<string, string>,
  detailMap: ProductDetailMap,
  detailAssetUrlMap: Record<string, string>
): string {
  const assetId = item.skuImageAssetId?.trim() ?? "";
  if (assetId) {
    const assetUrl = normalizeProductImageUrl(assetUrlMap[assetId]);
    if (assetUrl) {
      return assetUrl;
    }
  }

  const detailAssetUrl = resolveDetailAssetUrl(item, detailMap, detailAssetUrlMap);
  if (detailAssetUrl) {
    return detailAssetUrl;
  }

  const spuNo = item.spuNo?.trim() ?? "";
  if (spuNo) {
    const productImageUrl = normalizeProductImageUrl(productImageMap[spuNo]);
    if (productImageUrl) {
      return productImageUrl;
    }
  }

  return "";
}
