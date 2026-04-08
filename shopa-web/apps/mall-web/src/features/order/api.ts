import { apiClient } from "@/lib/api-client";
import {
  BuyerOrder,
  BuyerOrderAddress,
  BuyerOrderAmount,
  BuyerOrderItem,
  BuyerOrderSub,
  CreateOrderResult,
  ListMyOrdersResult,
  OrderStatusCode,
  PaymentStatusCode,
  RequestOrderPayResult,
  SubOrderStatusCode
} from "./types";

type RawOrderAmount = {
  goods_amount?: number | string;
  goodsAmount?: number | string;
  freight_amount?: number | string;
  freightAmount?: number | string;
  discount_amount?: number | string;
  discountAmount?: number | string;
  payable_amount?: number | string;
  payableAmount?: number | string;
  paid_amount?: number | string;
  paidAmount?: number | string;
  points_discount_amount?: number | string;
  pointsDiscountAmount?: number | string;
};

type RawOrderAddress = {
  source_address_id?: number | string;
  sourceAddressId?: number | string;
  source_address_version?: number | string;
  sourceAddressVersion?: number | string;
  receiver_name?: string;
  receiverName?: string;
  receiver_phone?: string;
  receiverPhone?: string;
  country_code?: string;
  countryCode?: string;
  province_code?: string;
  provinceCode?: string;
  province_name?: string;
  provinceName?: string;
  city_code?: string;
  cityCode?: string;
  city_name?: string;
  cityName?: string;
  district_code?: string;
  districtCode?: string;
  district_name?: string;
  districtName?: string;
  street?: string;
  detail?: string;
  postal_code?: string;
  postalCode?: string;
  latitude?: number | string;
  longitude?: number | string;
};

type RawOrderItem = {
  item_no?: string;
  itemNo?: string;
  order_no?: string;
  orderNo?: string;
  sub_order_no?: string;
  subOrderNo?: string;
  shop_no?: string;
  shopNo?: string;
  spu_no?: string;
  spuNo?: string;
  sku_no?: string;
  skuNo?: string;
  spu_title?: string;
  spuTitle?: string;
  sku_name?: string;
  skuName?: string;
  sku_image_asset_id?: number | string;
  skuImageAssetId?: number | string;
  qty?: number | string;
  sale_price?: number | string;
  salePrice?: number | string;
  market_price?: number | string;
  marketPrice?: number | string;
  sale_attrs_json?: string;
  saleAttrsJson?: string;
};

type RawOrderSub = {
  sub_order_no?: string;
  subOrderNo?: string;
  order_no?: string;
  orderNo?: string;
  shop_no?: string;
  shopNo?: string;
  sub_status?: number | string;
  subStatus?: number | string;
  amount?: RawOrderAmount;
  seller_remark?: string;
  sellerRemark?: string;
  buyer_remark?: string;
  buyerRemark?: string;
  created_at?: string;
  createdAt?: string;
  updated_at?: string;
  updatedAt?: string;
  items?: RawOrderItem[];
  points_used?: number | string;
  pointsUsed?: number | string;
  points_discount_amount?: number | string;
  pointsDiscountAmount?: number | string;
};

type RawOrderMain = {
  order_no?: string;
  orderNo?: string;
  user_id?: number | string;
  userId?: number | string;
  order_status?: number | string;
  orderStatus?: number | string;
  payment_status?: number | string;
  paymentStatus?: number | string;
  amount?: RawOrderAmount;
  address?: RawOrderAddress | null;
  reservation_no?: string;
  reservationNo?: string;
  pay_deadline_at?: string;
  payDeadlineAt?: string;
  paid_at?: string;
  paidAt?: string;
  closed_at?: string;
  closedAt?: string;
  created_at?: string;
  createdAt?: string;
  updated_at?: string;
  updatedAt?: string;
  version?: number | string;
  sub_orders?: RawOrderSub[];
  subOrders?: RawOrderSub[];
  buyer_remark?: string;
  buyerRemark?: string;
  cancel_reason_code?: number | string;
  cancelReasonCode?: number | string;
  points_reservation_no?: string;
  pointsReservationNo?: string;
  points_used?: number | string;
  pointsUsed?: number | string;
  points_discount_amount?: number | string;
  pointsDiscountAmount?: number | string;
  points_rule_snapshot_json?: string;
  pointsRuleSnapshotJson?: string;
  points_rule_snapshot_digest?: string;
  pointsRuleSnapshotDigest?: string;
};

type RawCreateOrderRes = {
  order?: RawOrderMain;
  idempotent_replay?: boolean;
  idempotentReplay?: boolean;
};

type RawListMyOrdersRes = {
  orders?: RawOrderMain[];
  next_cursor?: string;
  nextCursor?: string;
  has_more?: boolean;
  hasMore?: boolean;
};

type RawGetMyOrderDetailRes = {
  order?: RawOrderMain;
};

type RawUpdateMyOrderAddressRes = {
  order?: RawOrderMain;
};

type RawRequestPayRes = {
  order_no?: string;
  orderNo?: string;
  pay_no?: string;
  payNo?: string;
  payment_status?: number | string;
  paymentStatus?: number | string;
  pay_url?: string;
  payUrl?: string;
  pay_payload_json?: string;
  payPayloadJson?: string;
  expire_at?: string;
  expireAt?: string;
};

type RawCancelMyOrderRes = {
  order_no?: string;
  orderNo?: string;
  order_status?: number | string;
  orderStatus?: number | string;
};

type RawConfirmMyOrderReceivedRes = {
  order_no?: string;
  orderNo?: string;
  order_status?: number | string;
  orderStatus?: number | string;
  points_grant_triggered?: boolean;
  pointsGrantTriggered?: boolean;
};

function randomKey(prefix: string): string {
  return `${prefix}_${Date.now()}_${Math.random().toString(36).slice(2, 10)}`;
}

function toString(value: unknown): string {
  if (value === null || value === undefined) {
    return "";
  }
  if (typeof value === "string") {
    return value;
  }
  if (typeof value === "number" || typeof value === "boolean" || typeof value === "bigint") {
    return String(value);
  }
  if (value instanceof Date) {
    return value.toISOString();
  }
  if (typeof value === "object" && !Array.isArray(value)) {
    const record = value as Record<string, unknown>;
    const seconds = Number(record.seconds ?? 0);
    const nanos = Number(record.nanos ?? 0);
    if (Number.isFinite(seconds) || Number.isFinite(nanos)) {
      const millis = seconds * 1000 + Math.floor(nanos / 1_000_000);
      if (millis > 0) {
        return new Date(millis).toISOString();
      }
    }
    const candidates = [record.text, record.content, record.title, record.name, record.message, record.value];
    for (const item of candidates) {
      if (typeof item === "string" && item.trim()) {
        return item.trim();
      }
    }
  }
  return "";
}

function toNumber(value: unknown): number {
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : 0;
}

function normalizeEnumName(value: unknown): string {
  return toString(value).trim().toUpperCase();
}

function mapOrderStatus(value: unknown): OrderStatusCode {
  const normalized = normalizeEnumName(value);
  switch (normalized) {
    case "1":
    case "ORDER_STATUS_PENDING_PAY":
    case "PENDING_PAY":
      return "PENDING_PAY";
    case "2":
    case "ORDER_STATUS_PAID":
    case "PAID":
      return "PAID";
    case "3":
    case "ORDER_STATUS_CANCELED":
    case "CANCELED":
      return "CANCELED";
    case "4":
    case "ORDER_STATUS_CLOSED":
    case "CLOSED":
      return "CLOSED";
    case "5":
    case "ORDER_STATUS_FULFILLING":
    case "FULFILLING":
      return "FULFILLING";
    case "6":
    case "ORDER_STATUS_COMPLETED":
    case "COMPLETED":
      return "COMPLETED";
    case "7":
    case "ORDER_STATUS_REFUNDING":
    case "REFUNDING":
      return "REFUNDING";
    case "8":
    case "ORDER_STATUS_REFUNDED":
    case "REFUNDED":
      return "REFUNDED";
    default:
      return "UNSPECIFIED";
  }
}

function mapPaymentStatus(value: unknown): PaymentStatusCode {
  const normalized = normalizeEnumName(value);
  switch (normalized) {
    case "1":
    case "PAYMENT_STATUS_UNPAID":
    case "UNPAID":
      return "UNPAID";
    case "2":
    case "PAYMENT_STATUS_PAYING":
    case "PAYING":
      return "PAYING";
    case "3":
    case "PAYMENT_STATUS_PAID":
    case "PAID":
      return "PAID";
    case "4":
    case "PAYMENT_STATUS_PAY_FAILED":
    case "PAY_FAILED":
      return "PAY_FAILED";
    case "5":
    case "PAYMENT_STATUS_REFUNDED":
    case "REFUNDED":
      return "REFUNDED";
    default:
      return "UNSPECIFIED";
  }
}

function mapSubOrderStatus(value: unknown): SubOrderStatusCode {
  const normalized = normalizeEnumName(value);
  switch (normalized) {
    case "1":
    case "SUB_ORDER_STATUS_PENDING_PAY":
    case "PENDING_PAY":
      return "PENDING_PAY";
    case "2":
    case "SUB_ORDER_STATUS_PAID":
    case "PAID":
      return "PAID";
    case "3":
    case "SUB_ORDER_STATUS_CANCELED":
    case "CANCELED":
      return "CANCELED";
    case "4":
    case "SUB_ORDER_STATUS_CLOSED":
    case "CLOSED":
      return "CLOSED";
    case "5":
    case "SUB_ORDER_STATUS_WAIT_SHIP":
    case "WAIT_SHIP":
      return "WAIT_SHIP";
    case "6":
    case "SUB_ORDER_STATUS_SHIPPED":
    case "SHIPPED":
      return "SHIPPED";
    case "7":
    case "SUB_ORDER_STATUS_COMPLETED":
    case "COMPLETED":
      return "COMPLETED";
    case "8":
    case "SUB_ORDER_STATUS_REFUNDING":
    case "REFUNDING":
      return "REFUNDING";
    case "9":
    case "SUB_ORDER_STATUS_REFUNDED":
    case "REFUNDED":
      return "REFUNDED";
    default:
      return "UNSPECIFIED";
  }
}

function normalizeAmount(raw?: RawOrderAmount | null): BuyerOrderAmount {
  return {
    goodsAmount: toNumber(raw?.goods_amount ?? raw?.goodsAmount),
    freightAmount: toNumber(raw?.freight_amount ?? raw?.freightAmount),
    discountAmount: toNumber(raw?.discount_amount ?? raw?.discountAmount),
    payableAmount: toNumber(raw?.payable_amount ?? raw?.payableAmount),
    paidAmount: toNumber(raw?.paid_amount ?? raw?.paidAmount),
    pointsDiscountAmount: toNumber(raw?.points_discount_amount ?? raw?.pointsDiscountAmount)
  };
}

function normalizeAddress(raw?: RawOrderAddress | null): BuyerOrderAddress | null {
  if (!raw) {
    return null;
  }
  return {
    sourceAddressId: toNumber(raw.source_address_id ?? raw.sourceAddressId),
    sourceAddressVersion: toNumber(raw.source_address_version ?? raw.sourceAddressVersion),
    receiverName: toString(raw.receiver_name ?? raw.receiverName),
    receiverPhone: toString(raw.receiver_phone ?? raw.receiverPhone),
    countryCode: toString(raw.country_code ?? raw.countryCode),
    provinceCode: toString(raw.province_code ?? raw.provinceCode),
    provinceName: toString(raw.province_name ?? raw.provinceName),
    cityCode: toString(raw.city_code ?? raw.cityCode),
    cityName: toString(raw.city_name ?? raw.cityName),
    districtCode: toString(raw.district_code ?? raw.districtCode),
    districtName: toString(raw.district_name ?? raw.districtName),
    street: toString(raw.street),
    detail: toString(raw.detail),
    postalCode: toString(raw.postal_code ?? raw.postalCode),
    latitude: toNumber(raw.latitude),
    longitude: toNumber(raw.longitude)
  };
}

function normalizeItem(raw?: RawOrderItem): BuyerOrderItem {
  return {
    itemNo: toString(raw?.item_no ?? raw?.itemNo),
    orderNo: toString(raw?.order_no ?? raw?.orderNo),
    subOrderNo: toString(raw?.sub_order_no ?? raw?.subOrderNo),
    shopNo: toString(raw?.shop_no ?? raw?.shopNo),
    spuNo: toString(raw?.spu_no ?? raw?.spuNo),
    skuNo: toString(raw?.sku_no ?? raw?.skuNo),
    spuTitle: toString(raw?.spu_title ?? raw?.spuTitle),
    skuName: toString(raw?.sku_name ?? raw?.skuName),
    skuImageAssetId: toString(raw?.sku_image_asset_id ?? raw?.skuImageAssetId),
    qty: toNumber(raw?.qty),
    salePrice: toNumber(raw?.sale_price ?? raw?.salePrice),
    marketPrice: toNumber(raw?.market_price ?? raw?.marketPrice),
    saleAttrsJson: toString(raw?.sale_attrs_json ?? raw?.saleAttrsJson)
  };
}

function normalizeSubOrder(raw?: RawOrderSub): BuyerOrderSub {
  return {
    subOrderNo: toString(raw?.sub_order_no ?? raw?.subOrderNo),
    orderNo: toString(raw?.order_no ?? raw?.orderNo),
    shopNo: toString(raw?.shop_no ?? raw?.shopNo),
    subStatus: mapSubOrderStatus(raw?.sub_status ?? raw?.subStatus),
    amount: normalizeAmount(raw?.amount),
    sellerRemark: toString(raw?.seller_remark ?? raw?.sellerRemark),
    buyerRemark: toString(raw?.buyer_remark ?? raw?.buyerRemark),
    createdAt: toString(raw?.created_at ?? raw?.createdAt),
    updatedAt: toString(raw?.updated_at ?? raw?.updatedAt),
    items: Array.isArray(raw?.items) ? raw!.items!.map((item) => normalizeItem(item)) : [],
    pointsUsed: toNumber(raw?.points_used ?? raw?.pointsUsed),
    pointsDiscountAmount: toNumber(raw?.points_discount_amount ?? raw?.pointsDiscountAmount)
  };
}

function normalizeOrder(raw?: RawOrderMain | null): BuyerOrder | null {
  if (!raw) {
    return null;
  }
  const subOrders = Array.isArray(raw.sub_orders ?? raw.subOrders) ? (raw.sub_orders ?? raw.subOrders ?? []).map((item) => normalizeSubOrder(item)) : [];
  return {
    orderNo: toString(raw.order_no ?? raw.orderNo),
    userId: toNumber(raw.user_id ?? raw.userId),
    orderStatus: mapOrderStatus(raw.order_status ?? raw.orderStatus),
    paymentStatus: mapPaymentStatus(raw.payment_status ?? raw.paymentStatus),
    amount: normalizeAmount(raw.amount),
    address: normalizeAddress(raw.address),
    reservationNo: toString(raw.reservation_no ?? raw.reservationNo),
    payDeadlineAt: toString(raw.pay_deadline_at ?? raw.payDeadlineAt),
    paidAt: toString(raw.paid_at ?? raw.paidAt),
    closedAt: toString(raw.closed_at ?? raw.closedAt),
    createdAt: toString(raw.created_at ?? raw.createdAt),
    updatedAt: toString(raw.updated_at ?? raw.updatedAt),
    version: toNumber(raw.version),
    subOrders,
    buyerRemark: toString(raw.buyer_remark ?? raw.buyerRemark),
    cancelReasonCode: toString(raw.cancel_reason_code ?? raw.cancelReasonCode),
    pointsReservationNo: toString(raw.points_reservation_no ?? raw.pointsReservationNo),
    pointsUsed: toNumber(raw.points_used ?? raw.pointsUsed),
    pointsDiscountAmount: toNumber(raw.points_discount_amount ?? raw.pointsDiscountAmount),
    pointsRuleSnapshotJson: toString(raw.points_rule_snapshot_json ?? raw.pointsRuleSnapshotJson),
    pointsRuleSnapshotDigest: toString(raw.points_rule_snapshot_digest ?? raw.pointsRuleSnapshotDigest)
  };
}

function normalizeRequestPayResult(raw?: RawRequestPayRes): RequestOrderPayResult {
  return {
    orderNo: toString(raw?.order_no ?? raw?.orderNo),
    payNo: toString(raw?.pay_no ?? raw?.payNo),
    paymentStatus: mapPaymentStatus(raw?.payment_status ?? raw?.paymentStatus),
    payUrl: toString(raw?.pay_url ?? raw?.payUrl),
    payPayloadJson: toString(raw?.pay_payload_json ?? raw?.payPayloadJson),
    expireAt: toString(raw?.expire_at ?? raw?.expireAt)
  };
}

export function buildIdempotencyKey(prefix = "mall"): string {
  return randomKey(prefix);
}

export async function createOrderFromCart(payload: {
  checkoutToken: string;
  addressId: number;
  buyerRemark?: string;
  expectedSnapshotDigest?: string;
  submitSourceCode?: string;
  usePoints?: boolean;
  intentPoints?: number;
  expectedPointsCashAmount?: number;
  pointsRuleSnapshotDigest?: string;
  idempotencyKey?: string;
}): Promise<CreateOrderResult> {
  const response = await apiClient.post<RawCreateOrderRes>("/v1/order/buyer/orders:create-from-cart", {
    checkout_token: payload.checkoutToken,
    address_id: payload.addressId,
    buyer_remark: payload.buyerRemark ?? "",
    submit_source_code: payload.submitSourceCode ?? "mall-web",
    use_points: Boolean(payload.usePoints),
    intent_points: payload.intentPoints ?? 0,
    expected_points_cash_amount: payload.expectedPointsCashAmount ?? 0,
    points_rule_snapshot_digest: payload.pointsRuleSnapshotDigest ?? "",
    expected_snapshot_digest: payload.expectedSnapshotDigest ?? "",
    idempotency_key: payload.idempotencyKey ?? buildIdempotencyKey("order_from_cart")
  });
  const order = normalizeOrder(response?.order);
  return {
    orderNo: order?.orderNo ?? "",
    order,
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
  submitSourceCode?: string;
  usePoints?: boolean;
  intentPoints?: number;
  expectedPointsCashAmount?: number;
  pointsRuleSnapshotDigest?: string;
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
    submit_source_code: payload.submitSourceCode ?? "mall-web",
    use_points: Boolean(payload.usePoints),
    intent_points: payload.intentPoints ?? 0,
    expected_points_cash_amount: payload.expectedPointsCashAmount ?? 0,
    points_rule_snapshot_digest: payload.pointsRuleSnapshotDigest ?? "",
    idempotency_key: payload.idempotencyKey ?? buildIdempotencyKey("order_buy_now")
  });
  const order = normalizeOrder(response?.order);
  return {
    orderNo: order?.orderNo ?? "",
    order,
    idempotentReplay: Boolean(response?.idempotent_replay ?? response?.idempotentReplay)
  };
}

export async function updateMyOrderAddress(orderNo: string, addressId: number): Promise<BuyerOrder | null> {
  const response = await apiClient.patch<RawUpdateMyOrderAddressRes>(`/v1/order/buyer/orders/${encodeURIComponent(orderNo)}/address`, {
    address_id: addressId
  });
  return normalizeOrder(response?.order);
}

export async function requestOrderPay(payload: {
  orderNo: string;
  payChannel?: "ALIPAY" | "MOCK";
  idempotencyKey?: string;
}): Promise<RequestOrderPayResult> {
  const payChannel = payload.payChannel === "MOCK" ? 1 : 2;
  const response = await apiClient.post<RawRequestPayRes>("/v1/order/buyer/orders:request-pay", {
    order_no: payload.orderNo,
    pay_channel: payChannel,
    idempotency_key: payload.idempotencyKey ?? buildIdempotencyKey("request_pay")
  });
  return normalizeRequestPayResult(response);
}

export async function cancelMyOrder(payload: {
  orderNo: string;
  reasonCode?: number;
  idempotencyKey?: string;
}): Promise<{
  orderNo: string;
  orderStatus: OrderStatusCode;
}> {
  const response = await apiClient.post<RawCancelMyOrderRes>("/v1/order/buyer/orders:cancel", {
    order_no: payload.orderNo,
    reason_code: payload.reasonCode ?? 1,
    idempotency_key: payload.idempotencyKey ?? buildIdempotencyKey("cancel_order")
  });
  return {
    orderNo: toString(response?.order_no ?? response?.orderNo),
    orderStatus: mapOrderStatus(response?.order_status ?? response?.orderStatus)
  };
}

export async function confirmMyOrderReceived(payload: {
  orderNo: string;
  idempotencyKey?: string;
}): Promise<{
  orderNo: string;
  orderStatus: OrderStatusCode;
  pointsGrantTriggered: boolean;
}> {
  const response = await apiClient.post<RawConfirmMyOrderReceivedRes>("/v1/order/buyer/orders:confirm-received", {
    order_no: payload.orderNo,
    completion_source_code: "BUYER_CONFIRM_RECEIPT",
    idempotency_key: payload.idempotencyKey ?? buildIdempotencyKey("confirm_receipt")
  });
  return {
    orderNo: toString(response?.order_no ?? response?.orderNo),
    orderStatus: mapOrderStatus(response?.order_status ?? response?.orderStatus),
    pointsGrantTriggered: Boolean(response?.points_grant_triggered ?? response?.pointsGrantTriggered)
  };
}

export async function getMyOrderDetail(orderNo: string): Promise<BuyerOrder | null> {
  const response = await apiClient.get<RawGetMyOrderDetailRes>(`/v1/order/buyer/orders/${encodeURIComponent(orderNo)}`);
  return normalizeOrder(response?.order);
}

export async function listMyOrders(payload?: {
  pageSize?: number;
  nextCursor?: string;
  keyword?: string;
}): Promise<ListMyOrdersResult> {
  const response = await apiClient.get<RawListMyOrdersRes>("/v1/order/buyer/orders", {
    params: {
      page_size: payload?.pageSize ?? 100,
      next_cursor: payload?.nextCursor ?? "",
      keyword: payload?.keyword ?? ""
    }
  });
  return {
    orders: Array.isArray(response?.orders) ? response.orders.map((item) => normalizeOrder(item)).filter((item): item is BuyerOrder => Boolean(item)) : [],
    nextCursor: toString(response?.next_cursor ?? response?.nextCursor),
    hasMore: Boolean(response?.has_more ?? response?.hasMore)
  };
}
