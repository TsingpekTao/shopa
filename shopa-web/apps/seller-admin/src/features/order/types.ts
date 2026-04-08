export type SellerOrderStatusCode =
  | "UNSPECIFIED"
  | "PENDING_PAY"
  | "PAID"
  | "CANCELED"
  | "CLOSED"
  | "FULFILLING"
  | "COMPLETED"
  | "REFUNDING"
  | "REFUNDED";

export type SellerPaymentStatusCode = "UNSPECIFIED" | "UNPAID" | "PAYING" | "PAID" | "PAY_FAILED" | "REFUNDED";

export type SellerSubOrderStatusCode =
  | "UNSPECIFIED"
  | "PENDING_PAY"
  | "PAID"
  | "CANCELED"
  | "CLOSED"
  | "WAIT_SHIP"
  | "SHIPPED"
  | "COMPLETED"
  | "REFUNDING"
  | "REFUNDED";

export type SellerOrderAmount = {
  goodsAmount: number;
  freightAmount: number;
  discountAmount: number;
  payableAmount: number;
  paidAmount: number;
  pointsDiscountAmount: number;
};

export type SellerOrderAddress = {
  sourceAddressId: number;
  sourceAddressVersion: number;
  receiverName: string;
  receiverPhone: string;
  countryCode: string;
  provinceCode: string;
  provinceName: string;
  cityCode: string;
  cityName: string;
  districtCode: string;
  districtName: string;
  street: string;
  detail: string;
  postalCode: string;
  latitude: number;
  longitude: number;
};

export type SellerOrderItem = {
  itemNo: string;
  orderNo: string;
  subOrderNo: string;
  shopNo: string;
  spuNo: string;
  skuNo: string;
  spuTitle: string;
  skuName: string;
  skuImageAssetId: string;
  qty: number;
  salePrice: number;
  marketPrice: number;
  saleAttrsJson: string;
};

export type SellerOrderSub = {
  subOrderNo: string;
  orderNo: string;
  shopNo: string;
  subStatus: SellerSubOrderStatusCode;
  amount: SellerOrderAmount;
  sellerRemark: string;
  buyerRemark: string;
  createdAt: string;
  updatedAt: string;
  items: SellerOrderItem[];
  pointsUsed: number;
  pointsDiscountAmount: number;
};

export type SellerOrderMain = {
  orderNo: string;
  userId: number;
  orderStatus: SellerOrderStatusCode;
  paymentStatus: SellerPaymentStatusCode;
  amount: SellerOrderAmount;
  address: SellerOrderAddress | null;
  reservationNo: string;
  payDeadlineAt: string;
  paidAt: string;
  closedAt: string;
  createdAt: string;
  updatedAt: string;
  version: number;
  subOrders: SellerOrderSub[];
  buyerRemark: string;
  cancelReasonCode: string;
  pointsReservationNo: string;
  pointsUsed: number;
  pointsDiscountAmount: number;
  pointsRuleSnapshotJson: string;
  pointsRuleSnapshotDigest: string;
};

export type SellerOrderDetail = {
  order: SellerOrderMain | null;
  subOrder: SellerOrderSub | null;
};

export type SellerListOrdersPayload = {
  shopNo: string;
  pageSize?: number;
  nextCursor?: string;
  statuses?: SellerSubOrderStatusCode[];
  keyword?: string;
};

export type SellerListOrdersResponse = {
  orders: SellerOrderSub[];
  nextCursor: string;
  hasMore: boolean;
};

export type ShipSellerSubOrderPayload = {
  subOrderNo: string;
  logisticsNo?: string;
  shippedNote?: string;
};

export type ShipSellerSubOrderResult = {
  subOrderNo: string;
  subStatus: SellerSubOrderStatusCode;
};

export type SellerShippingListParams = {
  shopNo: string;
  statuses?: string[];
  pageSize?: number;
  nextCursor?: string;
  keyword?: string;
};

export type SellerShippingListResponse = {
  subOrders: SellerOrderSub[];
  nextCursor: string;
  hasMore: boolean;
};
