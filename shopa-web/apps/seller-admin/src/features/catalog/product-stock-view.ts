import { SellerInventoryStock } from "@/features/inventory/types";
import { SellerProductSku, SellerProductSpu, StockStatusCode } from "./types";

export interface SellerProductLiveStockSnapshot {
  spuNo: string;
  stockStatus: StockStatusCode;
  inventoryBacked: boolean;
}

type ResolveSellerProductLiveStockInput = {
  spu: Pick<SellerProductSpu, "spuNo" | "spuStockStatus">;
  skus?: Array<Pick<SellerProductSku, "skuNo" | "stockStatus">>;
  inventoryStocks?: Array<Pick<SellerInventoryStock, "skuNo" | "spuNo" | "availableQty" | "stockStatus">>;
};

export function resolveSellerProductLiveStockStatus(
  input: ResolveSellerProductLiveStockInput
): SellerProductLiveStockSnapshot {
  const skuNos = new Set((input.skus ?? []).map((sku) => sku.skuNo).filter(Boolean));
  const inventoryStocks = (input.inventoryStocks ?? []).filter((stock) => {
    if (stock.spuNo && stock.spuNo === input.spu.spuNo) {
      return true;
    }
    return skuNos.has(stock.skuNo);
  });

  if (inventoryStocks.length > 0) {
    const hasAvailable = inventoryStocks.some(
      (stock) => Number(stock.availableQty ?? 0) > 0 || stock.stockStatus === "STOCK_STATUS_IN_STOCK"
    );
    return {
      spuNo: input.spu.spuNo,
      stockStatus: hasAvailable ? "STOCK_STATUS_IN_STOCK" : "STOCK_STATUS_OUT_OF_STOCK",
      inventoryBacked: true
    };
  }

  const skuStatuses = (input.skus ?? [])
    .map((sku) => sku.stockStatus)
    .filter((status): status is StockStatusCode => Boolean(status) && status !== "STOCK_STATUS_UNSPECIFIED");

  if (skuStatuses.some((status) => status === "STOCK_STATUS_IN_STOCK")) {
    return {
      spuNo: input.spu.spuNo,
      stockStatus: "STOCK_STATUS_IN_STOCK",
      inventoryBacked: false
    };
  }

  if (skuStatuses.length > 0 && skuStatuses.every((status) => status === "STOCK_STATUS_OUT_OF_STOCK")) {
    return {
      spuNo: input.spu.spuNo,
      stockStatus: "STOCK_STATUS_OUT_OF_STOCK",
      inventoryBacked: false
    };
  }

  return {
    spuNo: input.spu.spuNo,
    stockStatus: input.spu.spuStockStatus,
    inventoryBacked: false
  };
}

export function applySellerProductLiveStock(
  products: SellerProductSpu[],
  snapshots: Map<string, SellerProductLiveStockSnapshot>
): SellerProductSpu[] {
  return products.map((product) => {
    const snapshot = snapshots.get(product.spuNo);
    if (!snapshot) {
      return product;
    }
    return {
      ...product,
      spuStockStatus: snapshot.stockStatus
    };
  });
}
