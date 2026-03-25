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
