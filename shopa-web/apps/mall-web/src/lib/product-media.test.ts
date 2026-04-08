import assert from "node:assert/strict";
import { collectProductDetailAssetIds, resolveProductImageUrl, type ProductDetailMap } from "./product-media.js";

function createDetailMap(): ProductDetailMap {
  return {
    SPU1001: {
      spuNo: "SPU1001",
      shopNo: "SHOP1",
      title: "Product A",
      subTitle: "",
      minSalePrice: 100,
      maxSalePrice: 100,
      minMarketPrice: 100,
      maxMarketPrice: 100,
      mainImageAssetIds: ["32", "32", "46"],
      detailImageAssetIds: [],
      skus: [
        {
          skuNo: "SKU1",
          spuNo: "SPU1001",
          shopNo: "SHOP1",
          skuName: "SKU1",
          skuImageAssetId: "46",
          salePrice: 100,
          marketPrice: 100,
          stockStatus: 1,
          saleAttrsJson: "[]"
        }
      ]
    }
  };
}

function runCollectProductDetailAssetIdsTest() {
  assert.deepEqual(collectProductDetailAssetIds(createDetailMap()), ["32", "46"]);
}

function runResolveProductImageUrlPrefersNormalizedAssetUrlTest() {
  const imageUrl = resolveProductImageUrl(
    {
      spuNo: "SPU1001",
      skuNo: "SKU1",
      skuImageAssetId: "46"
    },
    {
      "46": "https://cdn.example.com/media/folder/示例 图片.png?token=a b"
    },
    {},
    createDetailMap(),
    {}
  );

  assert.equal(imageUrl, "https://cdn.example.com/media/folder/%E7%A4%BA%E4%BE%8B%20%E5%9B%BE%E7%89%87.png?token=a%20b");
}

function runResolveProductImageUrlFallsBackToNormalizedProductImageMapTest() {
  const imageUrl = resolveProductImageUrl(
    {
      spuNo: "SPU1001",
      skuNo: "SKU1",
      skuImageAssetId: ""
    },
    {},
    {
      SPU1001: "https://cdn.example.com/fallback/箱包 1.png"
    },
    createDetailMap(),
    {}
  );

  assert.equal(imageUrl, "https://cdn.example.com/fallback/%E7%AE%B1%E5%8C%85%201.png");
}

function runResolveProductImageUrlPrefersDetailAssetOverLegacyProductImageTest() {
  const imageUrl = resolveProductImageUrl(
    {
      spuNo: "SPU1001",
      skuNo: "SKU1",
      skuImageAssetId: ""
    },
    {},
    {
      SPU1001: "https://cdn.example.com/legacy/old-image.png"
    },
    createDetailMap(),
    {
      "46": "https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/signed/folder/example.png?token=abc"
    }
  );

  assert.equal(imageUrl, "https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/signed/folder/example.png?token=abc");
}

function runResolveProductImageUrlKeepsSignedEncodedAssetUrlStableTest() {
  const imageUrl = resolveProductImageUrl(
    {
      spuNo: "SPU1001",
      skuNo: "SKU1",
      skuImageAssetId: "46"
    },
    {
      "46": "https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/seller_media%2F20260403%2Fexample_%E5%9B%BE.png?Expires=1775539886&OSSAccessKeyId=test&Signature=jwxre4%2FaAg5JbClInKeGMfPMPNs%3D"
    },
    {},
    createDetailMap(),
    {}
  );

  assert.equal(
    imageUrl,
    "https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/seller_media%2F20260403%2Fexample_%E5%9B%BE.png?Expires=1775539886&OSSAccessKeyId=test&Signature=jwxre4%2FaAg5JbClInKeGMfPMPNs%3D"
  );
}

function runResolveProductImageUrlFallsBackToNormalizedDetailAssetTest() {
  const detailMap: ProductDetailMap = {
    SPU2026: {
      spuNo: "SPU2026",
      shopNo: "SHOP2",
      title: "Product B",
      subTitle: "",
      minSalePrice: 29900,
      maxSalePrice: 29900,
      minMarketPrice: 29900,
      maxMarketPrice: 29900,
      mainImageAssetIds: ["32"],
      detailImageAssetIds: [],
      skus: [
        {
          skuNo: "SKU2026",
          spuNo: "SPU2026",
          shopNo: "SHOP2",
          skuName: "SKU2026",
          skuImageAssetId: "",
          salePrice: 29900,
          marketPrice: 29900,
          stockStatus: 1,
          saleAttrsJson: "[]"
        }
      ]
    }
  };

  const imageUrl = resolveProductImageUrl(
    {
      spuNo: "SPU2026",
      skuNo: "SKU2026",
      skuImageAssetId: ""
    },
    {},
    {},
    detailMap,
    {
      "32": "https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/example dir/example 图.png"
    }
  );

  assert.equal(imageUrl, "https://shopa-catalog-1.oss-cn-beijing.aliyuncs.com/example%20dir/example%20%E5%9B%BE.png");
}

runCollectProductDetailAssetIdsTest();
runResolveProductImageUrlPrefersNormalizedAssetUrlTest();
runResolveProductImageUrlFallsBackToNormalizedProductImageMapTest();
runResolveProductImageUrlPrefersDetailAssetOverLegacyProductImageTest();
runResolveProductImageUrlKeepsSignedEncodedAssetUrlStableTest();
runResolveProductImageUrlFallsBackToNormalizedDetailAssetTest();
console.log("product-media tests passed");
