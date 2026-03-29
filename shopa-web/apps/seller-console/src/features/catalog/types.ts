export type SpuStatus = "draft" | "reviewing" | "approved" | "onShelf" | "offShelf";

export const spuStatusCodeMap: Record<string, SpuStatus> = {
  SPU_STATUS_DRAFT: "draft",
  SPU_STATUS_REVIEWING: "reviewing",
  SPU_STATUS_REVIEWED: "reviewing",
  SPU_STATUS_APPROVED: "approved",
  SPU_STATUS_ON_SHELF: "onShelf",
  SPU_STATUS_OFF_SHELF: "offShelf",
  SPU_STATUS_REJECTED: "draft",
  SPU_STATUS_FROZEN: "offShelf",
  SPU_STATUS_DELETED: "offShelf",
};

export interface CatalogApiProduct {
  spuNo?: string;
  spu_no?: string;
  title: string;
  mainImageAssetIds?: Array<number | string>;
  main_image_asset_ids?: Array<number | string>;
  spuStatus?: string;
  spu_status?: string;
  spuStatusCode?: string;
  minSalePrice?: number | string;
  min_sale_price?: number | string;
  maxSalePrice?: number | string;
  max_sale_price?: number | string;
}

export interface CatalogApiResponse {
  products: CatalogApiProduct[];
  page?: number;
  page_size?: number;
  total?: number;
}

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

export interface CatalogSpuAttr {
  attrCode: string;
  attrName: string;
  value: string;
}
