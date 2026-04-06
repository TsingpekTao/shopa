import { apiClient } from "@shopa/api-client";
import { SellerWorkbenchResponse, ShopDashboardResponse } from "./types";

const mockWorkbench: SellerWorkbenchResponse = {
  userId: 123,
  shops: [],
  shopsTotal: 0,
  shopsTruncated: false,
  latestApplications: [],
  applicationsTotal: 0,
  applicationsTruncated: false,
  productSummary: { total: 0, onShelf: 0, offShelf: 0, reviewing: 0, draft: 0, rejected: 0 },
  partial: true,
  degradedFields: ["seller_shops"]
};

export async function fetchSellerWorkbench(): Promise<SellerWorkbenchResponse> {
  try {
    return await apiClient.get<SellerWorkbenchResponse>("/v1/seller/workbench", {
      silentDegraded: true
    });
  } catch (error) {
    console.warn("workbench fallback", error);
    return mockWorkbench;
  }
}

export async function fetchShopDashboard(shopNo: string): Promise<ShopDashboardResponse> {
  try {
    return await apiClient.get<ShopDashboardResponse>(`/v1/seller/shops/${shopNo}/dashboard`, {
      silentDegraded: true
    });
  } catch (error) {
    console.warn("shop dashboard fallback", error);
    return {
      shopNo,
      shopName: "Demo Shop",
      shopStatusCode: "ACTIVE",
      productSummary: { total: 0, onShelf: 0, offShelf: 0, reviewing: 0, draft: 0, rejected: 0 },
      inventoryRisk: { lowStockSkuCount: 0, outOfStockSkuCount: 0 },
      partial: true,
      degradedFields: ["inventory_risk"]
    };
  }
}
