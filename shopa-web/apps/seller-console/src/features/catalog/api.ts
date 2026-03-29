import { apiClient } from "@shopa/api-client";
import { CatalogApiProduct, CatalogApiResponse, CatalogDraftPayload, CatalogProduct, spuStatusCodeMap } from "./types";

const fallbackProducts: CatalogProduct[] = [
  {
    spuNo: "SPU-1001",
    title: "Shopa Smart Widget",
    status: "onShelf",
    salePrice: 199,
    mainImageAssetId: "asset-2"
  },
  {
    spuNo: "SPU-1002",
    title: "Shopa Pro Gadget",
    status: "draft",
    salePrice: 299,
    mainImageAssetId: "asset-3"
  }
];

const DEFAULT_SHOP_NO = process.env.NEXT_PUBLIC_SELLER_SHOP_NO ?? "demo-shop";

function toNumber(value?: number | string): number {
  if (typeof value === "number") {
    return value;
  }
  if (typeof value === "string") {
    const parsed = Number(value);
    return Number.isFinite(parsed) ? parsed : 0;
  }
  return 0;
}

function mapProduct(item: CatalogApiProduct): CatalogProduct {
  const statusCode = item.spuStatusCode ?? item.spuStatus ?? item.spu_status ?? "SPU_STATUS_DRAFT";
  const mainImageAssetId = item.mainImageAssetIds?.[0] ?? item.main_image_asset_ids?.[0];
  return {
    spuNo: item.spuNo ?? item.spu_no ?? "",
    title: item.title,
    status: spuStatusCodeMap[statusCode] ?? "draft",
    salePrice: toNumber(item.minSalePrice ?? item.min_sale_price),
    mainImageAssetId: mainImageAssetId ? String(mainImageAssetId) : undefined
  };
}

export async function fetchProductList(): Promise<CatalogProduct[]> {
  try {
    const response = await apiClient.get<CatalogApiResponse>("/v1/catalog/seller/products", {
      params: {
        page: 1,
        page_size: 20
      }
    });
    return (response.products ?? []).map(mapProduct);
  } catch (error) {
    return fallbackProducts;
  }
}

export async function submitDraft(payload: CatalogDraftPayload): Promise<{ success: boolean }> {
  try {
    await apiClient.post("/v1/catalog/seller/products/draft", {
      shop_no: DEFAULT_SHOP_NO,
      title: payload.draft.title,
      sub_title: payload.draft.summary,
      category_id: 0,
      brand_no: "",
      main_image_asset_ids: payload.draft.mainImageAssetId ? [Number(payload.draft.mainImageAssetId)] : [],
      detail_image_asset_ids: [],
      spu_attrs: Object.entries(payload.draft.attributeValues ?? {}).map(([attrCode, value]) => ({
        attr_code: attrCode,
        attr_name: attrCode,
        scope: 1,
        value: value ?? ""
      }))
    });
    return { success: true };
  } catch (error) {
    return { success: false };
  }
}
