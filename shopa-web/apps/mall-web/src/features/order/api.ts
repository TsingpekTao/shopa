import { apiClient } from "@/lib/api-client";
import { CreateOrderResult } from "./types";

type RawOrderMain = {
  order_no?: string;
  orderNo?: string;
};

type RawCreateOrderRes = {
  order?: RawOrderMain;
  idempotent_replay?: boolean;
  idempotentReplay?: boolean;
};

function randomKey(prefix: string): string {
  return `${prefix}_${Date.now()}_${Math.random().toString(36).slice(2, 10)}`;
}

export function buildIdempotencyKey(prefix = "mall"): string {
  return randomKey(prefix);
}

export async function createOrderFromCart(payload: {
  checkoutToken: string;
  addressId: number;
  buyerRemark?: string;
  expectedSnapshotDigest?: string;
  idempotencyKey?: string;
}): Promise<CreateOrderResult> {
  const response = await apiClient.post<RawCreateOrderRes>("/v1/order/buyer/orders:create-from-cart", {
    checkout_token: payload.checkoutToken,
    address_id: payload.addressId,
    buyer_remark: payload.buyerRemark ?? "",
    expected_snapshot_digest: payload.expectedSnapshotDigest ?? "",
    idempotency_key: payload.idempotencyKey ?? buildIdempotencyKey("order_from_cart")
  });
  return {
    orderNo: String(response?.order?.order_no ?? response?.order?.orderNo ?? ""),
    idempotentReplay: Boolean(response?.idempotent_replay ?? response?.idempotentReplay)
  };
}

export async function createOrderBuyNow(payload: {
  skuNo: string;
  spuNo: string;
  shopNo: string;
  qty: number;
  addressId: number;
  buyerRemark?: string;
  idempotencyKey?: string;
}): Promise<CreateOrderResult> {
  const response = await apiClient.post<RawCreateOrderRes>("/v1/order/buyer/orders:create-buy-now", {
    items: [
      {
        sku_no: payload.skuNo,
        spu_no: payload.spuNo,
        shop_no: payload.shopNo,
        qty: payload.qty
      }
    ],
    address_id: payload.addressId,
    buyer_remark: payload.buyerRemark ?? "",
    idempotency_key: payload.idempotencyKey ?? buildIdempotencyKey("order_buy_now")
  });
  return {
    orderNo: String(response?.order?.order_no ?? response?.order?.orderNo ?? ""),
    idempotentReplay: Boolean(response?.idempotent_replay ?? response?.idempotentReplay)
  };
}

