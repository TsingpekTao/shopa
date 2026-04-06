import { apiClient } from "@/lib/api-client";
import {
  BatchAdjustSellerStockInput,
  BatchAdjustSellerStockResultItem,
  InventoryStockStatus,
  SellerInventoryStock
} from "./types";

type TimestampLike = {
  seconds?: number | string;
  nanos?: number | string;
};

type RawInventoryStock = {
  sku_no?: string;
  skuNo?: string;
  spu_no?: string;
  spuNo?: string;
  shop_no?: string;
  shopNo?: string;
  total_qty?: number | string;
  totalQty?: number | string;
  locked_qty?: number | string;
  lockedQty?: number | string;
  available_qty?: number | string;
  availableQty?: number | string;
  stock_version?: number | string;
  stockVersion?: number | string;
  stock_status?: number | string;
  stockStatus?: number | string;
  updated_at?: string | TimestampLike;
  updatedAt?: string | TimestampLike;
};

type RawBatchGetResponse = {
  stocks?: RawInventoryStock[];
};

type RawAdjustResultItem = {
  sku_no?: string;
  skuNo?: string;
  spu_no?: string;
  spuNo?: string;
  shop_no?: string;
  shopNo?: string;
  total_qty?: number | string;
  totalQty?: number | string;
  locked_qty?: number | string;
  lockedQty?: number | string;
  available_qty?: number | string;
  availableQty?: number | string;
  stock_version?: number | string;
  stockVersion?: number | string;
  stock_status?: number | string;
  stockStatus?: number | string;
};

type RawBatchAdjustResponse = {
  results?: RawAdjustResultItem[];
};

const stockStatusByNumber: Record<number, InventoryStockStatus> = {
  0: "STOCK_STATUS_UNSPECIFIED",
  1: "STOCK_STATUS_IN_STOCK",
  2: "STOCK_STATUS_OUT_OF_STOCK"
};

const reasonCodeByName: Record<string, number> = {
  ADJUST_REASON_CODE_UNSPECIFIED: 0,
  ADJUST_REASON_CODE_SELLER_REPLENISH: 1,
  ADJUST_REASON_CODE_SELLER_CORRECTION: 2,
  ADJUST_REASON_CODE_ADMIN_MANUAL: 3,
  ADJUST_REASON_CODE_PURCHASE_IN: 4,
  ADJUST_REASON_CODE_RETURN_IN: 5,
  ADJUST_REASON_CODE_DAMAGE_OUT: 6,
  ADJUST_REASON_CODE_RECONCILE: 7
};

function toNumber(value: unknown, fallback = 0): number {
  const n = Number(value);
  return Number.isFinite(n) ? n : fallback;
}

function toString(value: unknown): string {
  if (value === null || value === undefined) {
    return "";
  }
  return String(value);
}

function toTimestamp(value: unknown): string | undefined {
  if (!value) {
    return undefined;
  }
  if (typeof value === "string") {
    return value;
  }
  if (typeof value !== "object") {
    return undefined;
  }
  const raw = value as TimestampLike;
  const seconds = toNumber(raw.seconds, 0);
  if (!seconds) {
    return undefined;
  }
  const nanos = toNumber(raw.nanos, 0);
  return new Date(seconds * 1000 + Math.floor(nanos / 1_000_000)).toISOString();
}

function toStockStatus(value: unknown): InventoryStockStatus {
  if (typeof value === "number") {
    return stockStatusByNumber[value] ?? "STOCK_STATUS_UNSPECIFIED";
  }
  if (typeof value === "string") {
    const trimmed = value.trim().toUpperCase();
    if (trimmed in (stockStatusByNumber as unknown as Record<string, unknown>)) {
      return stockStatusByNumber[toNumber(trimmed, 0)] ?? "STOCK_STATUS_UNSPECIFIED";
    }
    if (trimmed.startsWith("STOCK_STATUS_")) {
      return trimmed as InventoryStockStatus;
    }
    return `STOCK_STATUS_${trimmed}` as InventoryStockStatus;
  }
  return "STOCK_STATUS_UNSPECIFIED";
}

function normalizeStock(raw: RawInventoryStock): SellerInventoryStock {
  return {
    skuNo: toString(raw.sku_no ?? raw.skuNo),
    spuNo: toString(raw.spu_no ?? raw.spuNo),
    shopNo: toString(raw.shop_no ?? raw.shopNo),
    totalQty: toNumber(raw.total_qty ?? raw.totalQty, 0),
    lockedQty: toNumber(raw.locked_qty ?? raw.lockedQty, 0),
    availableQty: toNumber(raw.available_qty ?? raw.availableQty, 0),
    stockVersion: toNumber(raw.stock_version ?? raw.stockVersion, 0),
    stockStatus: toStockStatus(raw.stock_status ?? raw.stockStatus),
    updatedAt: toTimestamp(raw.updated_at ?? raw.updatedAt)
  };
}

function normalizeAdjustResult(raw: RawAdjustResultItem): BatchAdjustSellerStockResultItem {
  return {
    skuNo: toString(raw.sku_no ?? raw.skuNo),
    spuNo: toString(raw.spu_no ?? raw.spuNo),
    shopNo: toString(raw.shop_no ?? raw.shopNo),
    totalQty: toNumber(raw.total_qty ?? raw.totalQty, 0),
    lockedQty: toNumber(raw.locked_qty ?? raw.lockedQty, 0),
    availableQty: toNumber(raw.available_qty ?? raw.availableQty, 0),
    stockVersion: toNumber(raw.stock_version ?? raw.stockVersion, 0),
    stockStatus: toStockStatus(raw.stock_status ?? raw.stockStatus)
  };
}

export async function batchGetSkuInventory(skuNos: string[]): Promise<SellerInventoryStock[]> {
  const uniqueSkuNos = Array.from(new Set(skuNos.map((item) => item.trim()).filter(Boolean)));
  if (!uniqueSkuNos.length) {
    return [];
  }
  const response = await apiClient.post<RawBatchGetResponse>("/v1/inventory/sku:batch-get", {
    sku_nos: uniqueSkuNos
  });
  return (response.stocks ?? []).map((item) => normalizeStock(item)).filter((item) => item.skuNo);
}

export async function batchAdjustSellerStock(
  payload: BatchAdjustSellerStockInput
): Promise<BatchAdjustSellerStockResultItem[]> {
  const response = await apiClient.post<RawBatchAdjustResponse>(
    "/v1/inventory/seller/stock:adjust",
    {
      shop_no: payload.shopNo,
      items: payload.items.map((item) => ({
        sku_no: item.skuNo,
        spu_no: item.spuNo ?? "",
        delta_total_qty: item.deltaTotalQty,
        biz_no: item.bizNo ?? "",
        reason_code: reasonCodeByName[item.reasonCode ?? "ADJUST_REASON_CODE_SELLER_CORRECTION"] ?? 2,
        remark: item.remark ?? ""
      }))
    },
    {
      headers: payload.idempotencyKey ? { "x-idempotency-key": payload.idempotencyKey } : undefined
    }
  );
  return (response.results ?? []).map((item) => normalizeAdjustResult(item)).filter((item) => item.skuNo);
}

export async function adjustInventory(payload: {
  shopNo?: string;
  skuNo: string;
  spuNo?: string;
  delta: number;
  bizNo?: string;
  reasonCode?: string;
  remark?: string;
  idempotencyKey?: string;
}): Promise<{
  success: boolean;
  failedItems: Array<{ skuNo: string; reason: string }>;
  results?: BatchAdjustSellerStockResultItem[];
}> {
  try {
    const results = await batchAdjustSellerStock({
      shopNo: payload.shopNo ?? "",
      idempotencyKey: payload.idempotencyKey,
      items: [
        {
          skuNo: payload.skuNo,
          spuNo: payload.spuNo,
          deltaTotalQty: payload.delta,
          bizNo: payload.bizNo,
          reasonCode: (payload.reasonCode as never) ?? "ADJUST_REASON_CODE_SELLER_CORRECTION",
          remark: payload.remark
        }
      ]
    });
    return {
      success: results.length > 0,
      failedItems: [],
      results
    };
  } catch (error) {
    return {
      success: false,
      failedItems: [
        {
          skuNo: payload.skuNo,
          reason: error instanceof Error ? error.message : "Unable to reach inventory service"
        }
      ]
    };
  }
}
