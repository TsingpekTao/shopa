import assert from "node:assert/strict";
import { applySellerProductLiveStock, resolveSellerProductLiveStockStatus } from "./product-stock-view.js";
import { SellerProductSpu } from "./types.js";

function createProduct(overrides?: Partial<SellerProductSpu>): SellerProductSpu {
  return {
    spuNo: "SPU-1",
    shopNo: "SHOP-1",
    title: "测试商品",
    subTitle: "",
    categoryId: 1,
    brandNo: "",
    mainImageAssetIds: [],
    detailImageAssetIds: [],
    spuStatus: "SPU_STATUS_ON_SHELF",
    spuStockStatus: "STOCK_STATUS_OUT_OF_STOCK",
    minSalePrice: 9900,
    maxSalePrice: 9900,
    minMarketPrice: 12900,
    maxMarketPrice: 12900,
    version: 1,
    ...overrides
  };
}

function runResolveSellerProductLiveStockStatusPrefersInventoryTest() {
  const snapshot = resolveSellerProductLiveStockStatus({
    spu: createProduct(),
    skus: [
      {
        skuNo: "SKU-1",
        stockStatus: "STOCK_STATUS_OUT_OF_STOCK"
      }
    ],
    inventoryStocks: [
      {
        skuNo: "SKU-1",
        spuNo: "SPU-1",
        availableQty: 8,
        stockStatus: "STOCK_STATUS_IN_STOCK"
      }
    ]
  });

  assert.equal(snapshot.stockStatus, "STOCK_STATUS_IN_STOCK");
  assert.equal(snapshot.inventoryBacked, true);
}

function runResolveSellerProductLiveStockStatusFallsBackToSkuTest() {
  const snapshot = resolveSellerProductLiveStockStatus({
    spu: createProduct({ spuNo: "SPU-2", spuStockStatus: "STOCK_STATUS_UNSPECIFIED" }),
    skus: [
      {
        skuNo: "SKU-2",
        stockStatus: "STOCK_STATUS_IN_STOCK"
      }
    ]
  });

  assert.equal(snapshot.stockStatus, "STOCK_STATUS_IN_STOCK");
  assert.equal(snapshot.inventoryBacked, false);
}

function runApplySellerProductLiveStockTest() {
  const products = [createProduct(), createProduct({ spuNo: "SPU-2", title: "第二个商品" })];
  const merged = applySellerProductLiveStock(
    products,
    new Map([
      [
        "SPU-2",
        {
          spuNo: "SPU-2",
          stockStatus: "STOCK_STATUS_IN_STOCK",
          inventoryBacked: true
        }
      ]
    ])
  );

  assert.equal(merged[0].spuStockStatus, "STOCK_STATUS_OUT_OF_STOCK");
  assert.equal(merged[1].spuStockStatus, "STOCK_STATUS_IN_STOCK");
}

runResolveSellerProductLiveStockStatusPrefersInventoryTest();
runResolveSellerProductLiveStockStatusFallsBackToSkuTest();
runApplySellerProductLiveStockTest();
console.log("seller product stock view tests passed");
