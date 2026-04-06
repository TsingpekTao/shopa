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

export type AdjustReasonCode =
  | "SELLER_REPLENISH"
  | "SELLER_CORRECTION"
  | "ADMIN_MANUAL"
  | "PURCHASE_IN"
  | "RETURN_IN"
  | "DAMAGE_OUT"
  | "RECONCILE";

export interface InventoryAdjustPayload {
  shopNo?: string;
  skuNo: string;
  spuNo?: string;
  delta: number;
  reason?: string;
  reasonCode?: AdjustReasonCode;
  bizNo?: string;
  remark?: string;
  idempotencyKey?: string;
}

export interface FailedItem {
  skuNo: string;
  reason: string;
}

export interface InventoryAdjustResponse {
  success: boolean;
  failedItems: FailedItem[];
  results?: InventoryRecord[];
}
