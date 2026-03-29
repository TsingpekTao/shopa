export type InventoryStockStatus =
  | "STOCK_STATUS_UNSPECIFIED"
  | "STOCK_STATUS_IN_STOCK"
  | "STOCK_STATUS_OUT_OF_STOCK";

export type AdjustReasonCode =
  | "ADJUST_REASON_CODE_UNSPECIFIED"
  | "ADJUST_REASON_CODE_SELLER_REPLENISH"
  | "ADJUST_REASON_CODE_SELLER_CORRECTION"
  | "ADJUST_REASON_CODE_ADMIN_MANUAL"
  | "ADJUST_REASON_CODE_PURCHASE_IN"
  | "ADJUST_REASON_CODE_RETURN_IN"
  | "ADJUST_REASON_CODE_DAMAGE_OUT"
  | "ADJUST_REASON_CODE_RECONCILE";

export interface SellerInventoryStock {
  skuNo: string;
  spuNo: string;
  shopNo: string;
  totalQty: number;
  lockedQty: number;
  availableQty: number;
  stockVersion: number;
  stockStatus: InventoryStockStatus;
  updatedAt?: string;
}

export interface BatchAdjustSellerStockItemInput {
  skuNo: string;
  spuNo?: string;
  deltaTotalQty: number;
  bizNo?: string;
  reasonCode?: AdjustReasonCode;
  remark?: string;
}

export interface BatchAdjustSellerStockInput {
  shopNo: string;
  items: BatchAdjustSellerStockItemInput[];
  idempotencyKey?: string;
}

export interface BatchAdjustSellerStockResultItem {
  skuNo: string;
  spuNo: string;
  shopNo: string;
  totalQty: number;
  lockedQty: number;
  availableQty: number;
  stockVersion: number;
  stockStatus: InventoryStockStatus;
}
