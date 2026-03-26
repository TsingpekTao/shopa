export type SpuStatus = "draft" | "reviewing" | "approved" | "onShelf" | "offShelf";

export interface CatalogProduct {
  spuNo: string;
  title: string;
  mainImageAssetId?: string;
  status: SpuStatus;
  salePrice: number;
}

export interface CatalogDraft {
  title: string;
  summary: string;
  mainImageAssetId?: string;
  attributeValues: Record<string, string>;
}

export interface CatalogSkuDraft {
  skuNo?: string;
  name: string;
  salePrice: number;
  marketPrice: number;
  saleSpecs: Record<string, string>;
}

export interface CatalogDraftPayload {
  draft: CatalogDraft;
  skus: CatalogSkuDraft[];
}
