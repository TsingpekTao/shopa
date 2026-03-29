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
