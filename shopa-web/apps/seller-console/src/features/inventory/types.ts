export type InventoryStatus = "inStock" | "lowStock" | "outOfStock" | "unknown";

export interface InventoryRecord {
  skuNo: string;
  skuName: string;
  availableQty: number;
  lockedQty: number;
  totalQty: number;
  status: InventoryStatus;
  updatedAt: string;
}

export interface InventorySummary {
  totalSkus: number;
  lowStockCount: number;
  outOfStockCount: number;
  lastUpdated: string;
}

export interface InventoryAdjustPayload {
  skuNo: string;
  delta: number;
  reason?: string;
}

export interface FailedItem {
  skuNo: string;
  reason: string;
}

export interface InventoryAdjustResponse {
  success: boolean;
  failedItems: FailedItem[];
}
