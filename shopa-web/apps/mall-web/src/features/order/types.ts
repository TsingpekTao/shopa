export type OrderStatusCode =
  | "UNSPECIFIED"
  | "PENDING_PAY"
  | "PAID"
  | "CANCELED"
  | "CLOSED"
  | "FULFILLING"
  | "COMPLETED"
  | "REFUNDING"
  | "REFUNDED";

export type PaymentStatusCode = "UNSPECIFIED" | "UNPAID" | "PAYING" | "PAID" | "PAY_FAILED" | "REFUNDED";

export type SubOrderStatusCode =
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

export type BuyerOrderAmount = {
  goodsAmount: number;
  freightAmount: number;
  discountAmount: number;
  payableAmount: number;
  paidAmount: number;
  pointsDiscountAmount: number;
};

export type BuyerOrderAddress = {
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

export type BuyerOrderItem = {
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

export type BuyerOrderSub = {
  subOrderNo: string;
  orderNo: string;
  shopNo: string;
  subStatus: SubOrderStatusCode;
  amount: BuyerOrderAmount;
  sellerRemark: string;
  buyerRemark: string;
  createdAt: string;
  updatedAt: string;
  items: BuyerOrderItem[];
  pointsUsed: number;
  pointsDiscountAmount: number;
};

export type BuyerOrder = {
  orderNo: string;
  userId: number;
  orderStatus: OrderStatusCode;
  paymentStatus: PaymentStatusCode;
  amount: BuyerOrderAmount;
  address: BuyerOrderAddress | null;
  reservationNo: string;
  payDeadlineAt: string;
  paidAt: string;
  closedAt: string;
  createdAt: string;
  updatedAt: string;
  version: number;
  subOrders: BuyerOrderSub[];
  buyerRemark: string;
  cancelReasonCode: string;
  pointsReservationNo: string;
  pointsUsed: number;
  pointsDiscountAmount: number;
  pointsRuleSnapshotJson: string;
  pointsRuleSnapshotDigest: string;
};

export type CreateOrderResult = {
  orderNo: string;
  order: BuyerOrder | null;
  idempotentReplay: boolean;
};

export type ListMyOrdersResult = {
  orders: BuyerOrder[];
  nextCursor: string;
  hasMore: boolean;
};

export type RequestOrderPayResult = {
  orderNo: string;
  payNo: string;
  paymentStatus: PaymentStatusCode;
  payUrl: string;
  payPayloadJson: string;
  expireAt: string;
};
