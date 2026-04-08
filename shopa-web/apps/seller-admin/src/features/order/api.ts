import { apiClient } from "@/lib/api-client";
import { buildSellerOrderDetailPath, buildSellerOrderListPath, buildSellerSubOrderShipPath } from "./paths";
import type {
  SellerOrderAddress,
  SellerOrderAmount,
  SellerOrderDetail,
  SellerOrderItem,
  SellerOrderMain,
  SellerOrderStatusCode,
  SellerOrderSub,
  SellerPaymentStatusCode,
  SellerShippingListParams,
  SellerShippingListResponse,
  SellerSubOrderStatusCode
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
  amount?: RawOrderAmount | null;
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
  amount?: RawOrderAmount | null;
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

type RawGetSellerOrderDetailRes = {
  order?: RawOrderMain | null;
  sub_order?: RawOrderSub | null;
  subOrder?: RawOrderSub | null;
};

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

function mapOrderStatus(value: unknown): SellerOrderStatusCode {
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

function mapPaymentStatus(value: unknown): SellerPaymentStatusCode {
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

function mapSubOrderStatus(value: unknown): SellerSubOrderStatusCode {
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

function normalizeAmount(raw?: RawOrderAmount | null): SellerOrderAmount {
  return {
    goodsAmount: toNumber(raw?.goods_amount ?? raw?.goodsAmount),
    freightAmount: toNumber(raw?.freight_amount ?? raw?.freightAmount),
    discountAmount: toNumber(raw?.discount_amount ?? raw?.discountAmount),
    payableAmount: toNumber(raw?.payable_amount ?? raw?.payableAmount),
    paidAmount: toNumber(raw?.paid_amount ?? raw?.paidAmount),
    pointsDiscountAmount: toNumber(raw?.points_discount_amount ?? raw?.pointsDiscountAmount)
  };
}

function normalizeAddress(raw?: RawOrderAddress | null): SellerOrderAddress | null {
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

function normalizeItem(raw?: RawOrderItem): SellerOrderItem {
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

function normalizeSubOrder(raw?: RawOrderSub | null): SellerOrderSub | null {
  if (!raw) {
    return null;
  }
  return {
    subOrderNo: toString(raw.sub_order_no ?? raw.subOrderNo),
    orderNo: toString(raw.order_no ?? raw.orderNo),
    shopNo: toString(raw.shop_no ?? raw.shopNo),
    subStatus: mapSubOrderStatus(raw.sub_status ?? raw.subStatus),
    amount: normalizeAmount(raw.amount),
    sellerRemark: toString(raw.seller_remark ?? raw.sellerRemark),
    buyerRemark: toString(raw.buyer_remark ?? raw.buyerRemark),
    createdAt: toString(raw.created_at ?? raw.createdAt),
    updatedAt: toString(raw.updated_at ?? raw.updatedAt),
    items: Array.isArray(raw.items) ? raw.items.map((item) => normalizeItem(item)) : [],
    pointsUsed: toNumber(raw.points_used ?? raw.pointsUsed),
    pointsDiscountAmount: toNumber(raw.points_discount_amount ?? raw.pointsDiscountAmount)
  };
}

function normalizeOrder(raw?: RawOrderMain | null): SellerOrderMain | null {
  if (!raw) {
    return null;
  }
  const rawSubOrders = raw.sub_orders ?? raw.subOrders ?? [];
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
    subOrders: Array.isArray(rawSubOrders)
      ? rawSubOrders
          .map((item) => normalizeSubOrder(item))
          .filter((item): item is SellerOrderSub => Boolean(item))
      : [],
    buyerRemark: toString(raw.buyer_remark ?? raw.buyerRemark),
    cancelReasonCode: toString(raw.cancel_reason_code ?? raw.cancelReasonCode),
    pointsReservationNo: toString(raw.points_reservation_no ?? raw.pointsReservationNo),
    pointsUsed: toNumber(raw.points_used ?? raw.pointsUsed),
    pointsDiscountAmount: toNumber(raw.points_discount_amount ?? raw.pointsDiscountAmount),
    pointsRuleSnapshotJson: toString(raw.points_rule_snapshot_json ?? raw.pointsRuleSnapshotJson),
    pointsRuleSnapshotDigest: toString(raw.points_rule_snapshot_digest ?? raw.pointsRuleSnapshotDigest)
  };
}

export async function getSellerOrderDetail(subOrderNo: string): Promise<SellerOrderDetail> {
  const response = await apiClient.get<RawGetSellerOrderDetailRes>(buildSellerOrderDetailPath(subOrderNo));
  return {
    order: normalizeOrder(response?.order),
    subOrder: normalizeSubOrder(response?.sub_order ?? response?.subOrder)
  };
}

type RawListShopOrdersRes = {
  orders?: RawOrderSub[];
  next_cursor?: string;
  nextCursor?: string;
  has_more?: boolean;
  hasMore?: boolean;
};

export async function listSellerSubOrders(params: SellerShippingListParams): Promise<SellerShippingListResponse> {
  const normalizedShopNo = params.shopNo.trim();
  if (!normalizedShopNo) {
    return {
      subOrders: [],
      nextCursor: "",
      hasMore: false
    };
  }
  const payload: Record<string, unknown> = {
    shop_no: normalizedShopNo,
    page_size: params.pageSize ?? 20,
    next_cursor: params.nextCursor
  };
  if (params.keyword?.trim()) {
    payload.keyword = params.keyword.trim();
  }
  const statuses = params.statuses ?? ["SUB_ORDER_STATUS_WAIT_SHIP"];
  if (statuses.length > 0) {
    payload.statuses = statuses;
  }
  const response = await apiClient.post<RawListShopOrdersRes>(buildSellerOrderListPath(), payload);
  const subOrders = (response?.orders ?? [])
    .map((raw) => normalizeSubOrder(raw))
    .filter((item): item is SellerOrderSub => Boolean(item));
  return {
    subOrders,
    nextCursor: response?.next_cursor ?? response?.nextCursor ?? "",
    hasMore: Boolean(response?.has_more ?? response?.hasMore)
  };
}

type RawMarkSubOrderShippedRes = {
  sub_order_no?: string;
  subOrderNo?: string;
  sub_status?: string;
  subStatus?: string;
};

type MarkShippedOptions = {
  logisticsCompanyCode?: string;
  logisticsNo?: string;
  shippedNote?: string;
};

export async function markSubOrderShipped(subOrderNo: string, options?: MarkShippedOptions) {
  const payload: Record<string, unknown> = {
    sub_order_no: subOrderNo
  };
  if (options?.logisticsCompanyCode) {
    payload.logistics_company_code = options.logisticsCompanyCode;
  }
  if (options?.logisticsNo) {
    payload.logistics_no = options.logisticsNo;
  }
  if (options?.shippedNote) {
    payload.shipped_note = options.shippedNote;
  }
  const response = await apiClient.post<RawMarkSubOrderShippedRes>(buildSellerSubOrderShipPath(), payload);
  return {
    subOrderNo: response?.sub_order_no ?? response?.subOrderNo ?? subOrderNo,
    subStatus: mapSubOrderStatus(response?.sub_status ?? response?.subStatus)
  };
}
