import { apiClient } from "@/lib/api-client";
import { CartItem, CartSummary, CheckoutSnapshot, CheckoutSnapshotItem, MyCart } from "./types";

type RawCartItem = {
  user_id?: number | string;
  userId?: number | string;
  sku_no?: string;
  skuNo?: string;
  spu_no?: string;
  spuNo?: string;
  shop_no?: string;
  shopNo?: string;
  qty?: number | string;
  checked?: boolean;
  status?: number | string;
  invalid_reason_code?: string;
  invalidReasonCode?: string;
  spu_title?: string;
  spuTitle?: string;
  sku_name?: string;
  skuName?: string;
  sku_image_asset_id?: number | string;
  skuImageAssetId?: number | string;
  sale_price?: number | string;
  salePrice?: number | string;
  market_price?: number | string;
  marketPrice?: number | string;
  sale_attrs_json?: string;
  saleAttrsJson?: string;
};

type RawCartSummary = {
  total_item_count?: number | string;
  totalItemCount?: number | string;
  checked_item_count?: number | string;
  checkedItemCount?: number | string;
  checked_goods_amount?: number | string;
  checkedGoodsAmount?: number | string;
  checked_payable_amount?: number | string;
  checkedPayableAmount?: number | string;
};

type RawGetMyCartRes = {
  items?: RawCartItem[];
  summary?: RawCartSummary;
  loaded_from_backup?: boolean;
  loadedFromBackup?: boolean;
};

type RawAddCartItemRes = {
  item?: RawCartItem;
  summary?: RawCartSummary;
};

type RawCheckoutSnapshotItem = {
  sku_no?: string;
  skuNo?: string;
  spu_no?: string;
  spuNo?: string;
  shop_no?: string;
  shopNo?: string;
  qty?: number | string;
  settle_price?: number | string;
  settlePrice?: number | string;
  market_price?: number | string;
  marketPrice?: number | string;
  spu_title?: string;
  spuTitle?: string;
  sku_name?: string;
  skuName?: string;
  sku_image_asset_id?: number | string;
  skuImageAssetId?: number | string;
  sale_attrs_json?: string;
  saleAttrsJson?: string;
};

type RawCheckoutSnapshot = {
  checkout_token?: string;
  checkoutToken?: string;
  user_id?: number | string;
  userId?: number | string;
  items?: RawCheckoutSnapshotItem[];
  goods_amount?: number | string;
  goodsAmount?: number | string;
  freight_amount?: number | string;
  freightAmount?: number | string;
  payable_amount?: number | string;
  payableAmount?: number | string;
  snapshot_digest?: string;
  snapshotDigest?: string;
  expire_at?: string;
  expireAt?: string;
  created_at?: string;
  createdAt?: string;
};

type RawPrepareCheckoutRes = {
  snapshot?: RawCheckoutSnapshot;
};

function toFormData(payload: Record<string, string | number | boolean | string[] | undefined>): URLSearchParams {
  const formData = new URLSearchParams();
  Object.entries(payload).forEach(([key, value]) => {
    if (value === undefined) {
      return;
    }
    if (Array.isArray(value)) {
      value.forEach((item) => formData.append(key, item));
      return;
    }
    formData.set(key, String(value));
  });
  return formData;
}

function toNumber(value: unknown, fallback = 0): number {
  const num = Number(value);
  return Number.isFinite(num) ? num : fallback;
}

function toString(value: unknown): string {
  if (value === undefined || value === null) {
    return "";
  }
  return String(value);
}

function normalizeCartItem(raw: RawCartItem): CartItem {
  const saleAttrs = raw.sale_attrs_json ?? raw.saleAttrsJson;
  return {
    userId: toNumber(raw.user_id ?? raw.userId, 0),
    skuNo: toString(raw.sku_no ?? raw.skuNo),
    spuNo: toString(raw.spu_no ?? raw.spuNo),
    shopNo: toString(raw.shop_no ?? raw.shopNo),
    qty: toNumber(raw.qty, 0),
    checked: Boolean(raw.checked),
    status: toNumber(raw.status, 0),
    invalidReasonCode: toString(raw.invalid_reason_code ?? raw.invalidReasonCode),
    spuTitle: toString(raw.spu_title ?? raw.spuTitle),
    skuName: toString(raw.sku_name ?? raw.skuName),
    skuImageAssetId: toString(raw.sku_image_asset_id ?? raw.skuImageAssetId),
    salePrice: toNumber(raw.sale_price ?? raw.salePrice, 0),
    marketPrice: toNumber(raw.market_price ?? raw.marketPrice, 0),
    saleAttrsJson: toString((saleAttrs ?? "[]"))
  };
}

function normalizeSummary(raw?: RawCartSummary): CartSummary {
  return {
    totalItemCount: toNumber(raw?.total_item_count ?? raw?.totalItemCount, 0),
    checkedItemCount: toNumber(raw?.checked_item_count ?? raw?.checkedItemCount, 0),
    checkedGoodsAmount: toNumber(raw?.checked_goods_amount ?? raw?.checkedGoodsAmount, 0),
    checkedPayableAmount: toNumber(raw?.checked_payable_amount ?? raw?.checkedPayableAmount, 0)
  };
}

function normalizeCheckoutItem(raw: RawCheckoutSnapshotItem): CheckoutSnapshotItem {
  const saleAttrs = raw.sale_attrs_json ?? raw.saleAttrsJson;
  return {
    skuNo: toString(raw.sku_no ?? raw.skuNo),
    spuNo: toString(raw.spu_no ?? raw.spuNo),
    shopNo: toString(raw.shop_no ?? raw.shopNo),
    qty: toNumber(raw.qty, 0),
    settlePrice: toNumber(raw.settle_price ?? raw.settlePrice, 0),
    marketPrice: toNumber(raw.market_price ?? raw.marketPrice, 0),
    spuTitle: toString(raw.spu_title ?? raw.spuTitle),
    skuName: toString(raw.sku_name ?? raw.skuName),
    skuImageAssetId: toString(raw.sku_image_asset_id ?? raw.skuImageAssetId),
    saleAttrsJson: toString((saleAttrs ?? "[]"))
  };
}

function normalizeSnapshot(raw?: RawCheckoutSnapshot): CheckoutSnapshot {
  return {
    checkoutToken: toString(raw?.checkout_token ?? raw?.checkoutToken),
    userId: toNumber(raw?.user_id ?? raw?.userId, 0),
    items: Array.isArray(raw?.items) ? raw!.items!.map((item) => normalizeCheckoutItem(item)) : [],
    goodsAmount: toNumber(raw?.goods_amount ?? raw?.goodsAmount, 0),
    freightAmount: toNumber(raw?.freight_amount ?? raw?.freightAmount, 0),
    payableAmount: toNumber(raw?.payable_amount ?? raw?.payableAmount, 0),
    snapshotDigest: toString(raw?.snapshot_digest ?? raw?.snapshotDigest),
    expireAt: toString(raw?.expire_at ?? raw?.expireAt),
    createdAt: toString(raw?.created_at ?? raw?.createdAt)
  };
}

export async function getMyCart(onlyChecked = false): Promise<MyCart> {
  const response = await apiClient.get<RawGetMyCartRes>("/v1/cart/me", {
    params: { only_checked: onlyChecked }
  });
  return {
    items: Array.isArray(response?.items) ? response.items.map((item) => normalizeCartItem(item)) : [],
    summary: normalizeSummary(response?.summary),
    loadedFromBackup: Boolean(response?.loaded_from_backup ?? response?.loadedFromBackup)
  };
}

export async function addCartItem(payload: {
  skuNo: string;
  spuNo: string;
  shopNo: string;
  qty: number;
  checked?: boolean;
}): Promise<{
  item?: CartItem;
  summary: CartSummary;
}> {
  const response = await apiClient.post<RawAddCartItemRes>(
    "/v1/cart/items:add",
    toFormData({
      sku_no: payload.skuNo,
      spu_no: payload.spuNo,
      shop_no: payload.shopNo,
      qty: payload.qty,
      checked: payload.checked ?? true
    }),
    {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    }
  );
  return {
    item: response?.item ? normalizeCartItem(response.item) : undefined,
    summary: normalizeSummary(response?.summary)
  };
}

export async function updateCartItemQty(skuNo: string, qty: number): Promise<MyCart> {
  const response = await apiClient.post<RawGetMyCartRes>(
    "/v1/cart/items:qty",
    toFormData({
      sku_no: skuNo,
      qty
    }),
    {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    }
  );
  return {
    items: response?.items?.map((item) => normalizeCartItem(item)) ?? [],
    summary: normalizeSummary(response?.summary),
    loadedFromBackup: false
  };
}

export async function toggleCartItemChecked(skuNo: string, checked: boolean): Promise<MyCart> {
  const response = await apiClient.post<RawGetMyCartRes>(
    "/v1/cart/items:check",
    toFormData({
      sku_no: skuNo,
      checked
    }),
    {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    }
  );
  return {
    items: response?.items?.map((item) => normalizeCartItem(item)) ?? [],
    summary: normalizeSummary(response?.summary),
    loadedFromBackup: false
  };
}

export async function removeCartItems(skuNos: string[]): Promise<MyCart> {
  const response = await apiClient.post<RawGetMyCartRes>(
    "/v1/cart/items:remove",
    toFormData({
      sku_nos: skuNos
    }),
    {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    }
  );
  return {
    items: response?.items?.map((item) => normalizeCartItem(item)) ?? [],
    summary: normalizeSummary(response?.summary),
    loadedFromBackup: false
  };
}

export async function prepareCheckout(payload: {
  scope: 1 | 2;
  skuNos?: string[];
  addressId: number;
}): Promise<CheckoutSnapshot> {
  const response = await apiClient.post<RawPrepareCheckoutRes>(
    "/v1/cart/checkout:prepare",
    toFormData({
      scope: payload.scope,
      sku_nos: payload.skuNos ?? [],
      address_id: payload.addressId
    }),
    {
      headers: {
        "Content-Type": "application/x-www-form-urlencoded"
      }
    }
  );
  return normalizeSnapshot(response?.snapshot);
}
