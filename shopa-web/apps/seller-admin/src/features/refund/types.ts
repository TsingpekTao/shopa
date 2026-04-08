export type SellerRefundBatchStatus =
  | "UNSPECIFIED"
  | "PENDING_SELLER_REVIEW"
  | "SELLER_REJECTED"
  | "WAIT_REFUND_TASK"
  | "REFUND_PROCESSING"
  | "REFUNDED"
  | "CANCELED"
  | "CLOSED";

export type SellerRefundTaskStatus = "UNSPECIFIED" | "PENDING" | "PROCESSING" | "SUCCEEDED" | "FAILED" | "DEAD_LETTER";

export type RefundTone = "warning" | "processing" | "success" | "danger" | "default";

export type SellerRefundBatch = {
  refundBatchNo: string;
  orderNo: string;
  userId: number;
  shopNo: string;
  batchStatus: SellerRefundBatchStatus;
  applyRefundAmount: number;
  approvedRefundAmount: number;
  caseCount: number;
  subOrderNos: string[];
  reviewDeadlineAt: string;
  autoApprovedAt: string;
  createdAt: string;
  updatedAt: string;
};

export type SellerRefundCase = {
  afterSaleNo: string;
  refundBatchNo: string;
  orderNo: string;
  subOrderNo: string;
  itemNo: string;
  shopNo: string;
  spuNo: string;
  skuNo: string;
  qty: number;
  afterSaleStatus: SellerRefundBatchStatus;
  applyRefundAmount: number;
  approvedRefundAmount: number;
  reasonCode: string;
  reasonDesc: string;
  buyerRemark: string;
  sellerReply: string;
  scopeCode: string;
  selectedItemNos: string[];
  paymentNo: string;
  reviewDeadlineAt: string;
  autoApprovedAt: string;
  createdAt: string;
  updatedAt: string;
};

export type SellerRefundTask = {
  refundTaskNo: string;
  afterSaleNo: string;
  orderNo: string;
  subOrderNo: string;
  payNo: string;
  refundAmount: number;
  status: SellerRefundTaskStatus;
  retryCount: number;
  lastErrorCode: string;
  lastErrorMessage: string;
  createdAt: string;
  updatedAt: string;
  pointsReturnAmount: number;
  pointsReverseAmount: number;
  pointsCashOffsetAmount: number;
  finalCashRefundAmount: number;
  accountDebtAfter: number;
};

export type SellerRefundBatchDetail = {
  batch: SellerRefundBatch | null;
  cases: SellerRefundCase[];
  refundTasks: SellerRefundTask[];
};

export type SellerRefundListResult = {
  list: SellerRefundBatch[];
  nextCursor: string;
  hasMore: boolean;
};
