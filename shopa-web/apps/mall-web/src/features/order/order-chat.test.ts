import assert from "node:assert/strict";
import { buildAfterSaleConsultHref, resolveAfterSaleConsultParams } from "./order-chat.js";

function createOrder() {
  return {
    orderNo: "ORD1001",
    subOrders: [
      {
        subOrderNo: "SUB1001",
        shopNo: "SHOP1001",
        items: [
          {
            subOrderNo: "",
            shopNo: "",
            spuNo: "SPU1001",
            skuNo: "SKU1001"
          }
        ]
      },
      {
        subOrderNo: "SUB1002",
        shopNo: "SHOP1002",
        items: [
          {
            subOrderNo: "SUB1002",
            shopNo: "SHOP1002",
            spuNo: "SPU1002",
            skuNo: "SKU1002"
          }
        ]
      }
    ]
  } as never;
}

function runResolveAfterSaleConsultParamsTest() {
  assert.deepEqual(resolveAfterSaleConsultParams(createOrder()), {
    shopNo: "SHOP1001",
    subOrderNo: "SUB1001",
    spuNo: "SPU1001",
    skuNo: "SKU1001"
  });
}

function runBuildAfterSaleConsultHrefTest() {
  assert.equal(
    buildAfterSaleConsultHref(createOrder()),
    "/me/messages?shop_no=SHOP1001&spu_no=SPU1001&sku_no=SKU1001&order_no=ORD1001&sub_order_no=SUB1001&scene_code=AFTER_SALE"
  );
  assert.equal(
    buildAfterSaleConsultHref({
      orderNo: "ORD1002",
      subOrders: [
        {
          subOrderNo: "SUB_EMPTY",
          shopNo: "",
          items: [{ subOrderNo: "", shopNo: "", spuNo: "SPU_EMPTY", skuNo: "SKU_EMPTY" }]
        },
        {
          subOrderNo: "SUB2002",
          shopNo: "SHOP2002",
          items: [{ subOrderNo: "SUB2002", shopNo: "SHOP2002", spuNo: "SPU2002", skuNo: "SKU2002" }]
        }
      ]
    } as never),
    "/me/messages?shop_no=SHOP2002&spu_no=SPU2002&sku_no=SKU2002&order_no=ORD1002&sub_order_no=SUB2002&scene_code=AFTER_SALE"
  );
  assert.equal(
    buildAfterSaleConsultHref({
      orderNo: "ORD_EMPTY",
      subOrders: [{ subOrderNo: "", shopNo: "", items: [] }]
    } as never),
    ""
  );
}

runResolveAfterSaleConsultParamsTest();
runBuildAfterSaleConsultHrefTest();
console.log("order-chat tests passed");
