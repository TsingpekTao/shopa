import { apiClient } from "@shopa/api-client";
import { InventoryAdjustPayload, InventoryAdjustResponse, InventoryRecord, InventorySummary } from "./types";

const mockList: InventoryRecord[] = [
  {
    skuNo: "SKU-1001",
    skuName: "Shopa Smart Widget",
    availableQty: 120,
    lockedQty: 0,
    totalQty: 120,
    status: "inStock",
    updatedAt: new Date().toISOString()
  },
  {
    skuNo: "SKU-1002",
    skuName: "Shopa Pro Gadget",
    availableQty: 8,
    lockedQty: 2,
    totalQty: 10,
    status: "lowStock",
    updatedAt: new Date().toISOString()
  }
];

export async function fetchInventoryList(): Promise<InventoryRecord[]> {
  try {
    return await apiClient.get<InventoryRecord[]>("/v1/seller/inventory/stocks");
  } catch (error) {
    return mockList;
  }
}

export async function adjustInventory(payload: InventoryAdjustPayload): Promise<InventoryAdjustResponse> {
  try {
    return await apiClient.post<InventoryAdjustResponse>("/v1/inventory/seller/stock:adjust", payload);
  } catch (error) {
    return {
      success: false,
      failedItems: [{ skuNo: payload.skuNo, reason: "Mock fallback: inventory check endpoint unavailable" }]
    };
  }
}

export async function fetchInventorySummary(): Promise<InventorySummary> {
  try {
    return await apiClient.get<InventorySummary>("/v1/seller/inventory/summary");
  } catch (error) {
    return {
      totalSkus: mockList.length,
      lowStockCount: mockList.filter((item) => item.status === "lowStock").length,
      outOfStockCount: mockList.filter((item) => item.status === "outOfStock").length,
      lastUpdated: new Date().toISOString()
    };
  }
}
