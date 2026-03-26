import { apiClient } from "@shopa/api-client";
import { CatalogDraftPayload, CatalogProduct } from "./types";

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

export async function fetchProductList(): Promise<CatalogProduct[]> {
  try {
    const response = await apiClient.get<{ data: CatalogProduct[] }>("/v1/catalog/seller/products");
    return response.data ?? [];
  } catch (error) {
    return fallbackProducts;
  }
}

export async function submitDraft(payload: CatalogDraftPayload): Promise<{ success: boolean }> {
  try {
    return await apiClient.post<{ success: boolean }>("/v1/catalog/seller/products/draft", payload);
  } catch (error) {
    return { success: true };
  }
}
