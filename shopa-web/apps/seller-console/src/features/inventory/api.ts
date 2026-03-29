import { apiClient } from "@shopa/api-client";
import { InventoryAdjustPayload, InventoryAdjustResponse, InventoryRecord, InventorySummary } from "./types";

const shopNo = process.env.NEXT_PUBLIC_INVENTORY_SHOP_NO ?? "demo-shop";
const fallbackStocks: InventoryRecord[] = [
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

type InventoryStockResponse = {
  sku_no: string;
  spu_no?: string;
  shop_no?: string;
  total_qty: number;
  locked_qty: number;
  available_qty: number;
  stock_version?: number;
  stock_status?: number;
  updated_at?: string | { seconds?: number; nanos?: number };
};

const stockStatusMap: Record<number, InventoryRecord["status"]> = {
  1: "inStock",
  2: "outOfStock"
};

function toDateString(raw?: string | { seconds?: number; nanos?: number }): string {
  if (!raw) {
    return new Date().toISOString();
  }
  if (typeof raw === "string") {
    return raw;
  }
  const seconds = raw.seconds ?? 0;
  const millis = (raw.nanos ?? 0) / 1e6;
  return new Date(seconds * 1000 + millis).toISOString();
}

function normalizeStock(stock: InventoryStockResponse): InventoryRecord {
  const baseStatus = stock.stock_status ? stockStatusMap[stock.stock_status] ?? "unknown" : "unknown";
  const status =
    baseStatus === "inStock" && stock.available_qty > 0 && stock.available_qty <= 10
      ? "lowStock"
      : baseStatus;
  return {
    skuNo: stock.sku_no,
    skuName: stock.sku_no,
    availableQty: stock.available_qty,
    lockedQty: stock.locked_qty,
    totalQty: stock.total_qty,
    status,
    updatedAt: toDateString(stock.updated_at)
  };
}

export async function fetchInventoryList(): Promise<InventoryRecord[]> {
  const skuNos = fallbackStocks.map((item) => item.skuNo);
  try {
    const response = await apiClient.post<{ stocks: InventoryStockResponse[] }>("/v1/inventory/sku:batch-get", {
      sku_nos: skuNos
    });
    return (response?.stocks ?? []).map(normalizeStock);
  } catch {
    return fallbackStocks;
  }
}

function buildAdjustRequest(payload: InventoryAdjustPayload) {
  return {
    shop_no: shopNo,
    items: [
      {
        sku_no: payload.skuNo,
        delta_total_qty: payload.delta,
        biz_no: payload.bizNo ?? "",
        reason_code: payload.reasonCode ?? "SELLER_CORRECTION",
        remark: payload.remark ?? payload.reason ?? ""
      }
    ]
  };
}

export async function adjustInventory(payload: InventoryAdjustPayload): Promise<InventoryAdjustResponse> {
  const idempotencyKey = payload.idempotencyKey ?? `inventory-adjust:${payload.skuNo}:${Date.now()}`;
  try {
    const response = await apiClient.post<{ results: InventoryStockResponse[] }>("/v1/inventory/seller/stock:adjust", buildAdjustRequest(payload), {
      headers: { "x-idempotency-key": idempotencyKey }
    });
    const results = (response?.results ?? []).map(normalizeStock);
    return {
      success: results.length > 0,
      failedItems: [],
      results
    };
  } catch (error) {
    return {
      success: false,
      failedItems: [{ skuNo: payload.skuNo, reason: "Unable to reach inventory service" }]
    };
  }
}

export async function fetchInventorySummary(): Promise<InventorySummary> {
  try {
    const records = await fetchInventoryList();
    const lowStockCount = records.filter((item) => item.status === "lowStock").length;
    const outOfStockCount = records.filter((item) => item.status === "outOfStock").length;
    const lastUpdated = records.reduce((acc, record) => (record.updatedAt > acc ? record.updatedAt : acc), records[0]?.updatedAt ?? new Date().toISOString());
    return {
      totalSkus: records.length,
      lowStockCount,
      outOfStockCount,
      lastUpdated
    };
  } catch {
    return {
      totalSkus: fallbackStocks.length,
      lowStockCount: fallbackStocks.filter((item) => item.status === "lowStock").length,
      outOfStockCount: fallbackStocks.filter((item) => item.status === "outOfStock").length,
      lastUpdated: new Date().toISOString()
    };
  }
}
