import { SellerApplicationSummary, SellerProductSummary, SellerShopSummary } from "@shopa/types";

export interface SellerWorkbenchResponse {
  userId: number;
  shops: SellerShopSummary[];
  shopsTotal: number;
  shopsTruncated: boolean;
  latestApplications: SellerApplicationSummary[];
  applicationsTotal: number;
  applicationsTruncated: boolean;
  productSummary: SellerProductSummary;
  partial?: boolean;
  degradedFields?: string[];
}

export interface InventoryRiskSummary {
  lowStockSkuCount: number;
  outOfStockSkuCount: number;
  partial?: boolean;
  degradedFields?: string[];
}

export interface ShopDashboardResponse {
  shopNo: string;
  shopName: string;
  shopStatusCode: string;
  productSummary: SellerProductSummary;
  inventoryRisk: InventoryRiskSummary;
  partial?: boolean;
  degradedFields?: string[];
}

export interface SellerEntityProfileInput {
  merchantTypeCode?: string;
  entityName: string;
  contactName: string;
  contactPhone: string;
  contactEmail?: string;
  ext?: Record<string, string>;
}

export interface SellerShopProfileInput {
  shopName: string;
  shopDisplayName: string;
  shopTypeCode?: string;
  servicePhone?: string;
  serviceEmail?: string;
  description?: string;
  ext?: Record<string, string>;
}

export interface SellerApplicationDetail {
  applicationNo: string;
  version: number;
  applicationStatusCode?: string;
  entity?: SellerEntityProfileInput;
  shop?: SellerShopProfileInput;
  submittedAt?: string;
  updatedAt?: string;
}

export interface CreateApplicationDraftPayload {
  entity: SellerEntityProfileInput;
  shop: SellerShopProfileInput;
}

export interface UpdateApplicationDraftPayload {
  expectedVersion: number;
  entity: SellerEntityProfileInput;
  shop: SellerShopProfileInput;
  updateMask?: string[];
}

export interface SubmitApplicationPayload {
  expectedVersion: number;
}

export interface ResubmitApplicationPayload extends SubmitApplicationPayload {
  entity?: SellerEntityProfileInput;
  shop?: SellerShopProfileInput;
  updateMask?: string[];
}

export interface SellerApplicationListItem {
  applicationNo: string;
  version: number;
  applicationStatusCode: string;
  submittedAt?: string;
  updatedAt?: string;
}

export interface ListMyApplicationsPayload {
  page?: number;
  pageSize?: number;
  statuses?: number[];
}

export interface ListMyApplicationsResponse {
  applications: SellerApplicationListItem[];
  page: number;
  pageSize: number;
  total: number;
}

export interface SellerStoreCategory {
  id: number;
  shopNo: string;
  parentId: number;
  name: string;
  level: number;
  sortOrder: number;
  isVisible: boolean;
  isDeleted?: boolean;
  productCount: number;
  children: SellerStoreCategory[];
}

export interface SellerStoreCategorySortItem {
  categoryId: number;
  sortOrder: number;
}

export interface CreateSellerStoreCategoryPayload {
  parentId?: number;
  name: string;
  sortOrder?: number;
  isVisible?: boolean;
}

export interface UpdateSellerStoreCategoryPayload {
  name?: string;
  sortOrder?: number;
  isVisible?: boolean;
}

export interface SellerProductStoreCategoryBinding {
  shopNo: string;
  spuNo: string;
  storeCategoryId: number;
  storeCategoryL1: number;
  storeCategoryL2: number;
  storeCategoryPath: number[];
  storeCategoryName: string;
  updatedAt?: string;
}
