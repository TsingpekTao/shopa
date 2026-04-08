import type { BuyerOrderItem, BuyerOrderSub } from "../order/types";

export type BuyerRefundBatchStatus =
  | "UNSPECIFIED"
  | "PENDING_SELLER_REVIEW"
  | "SELLER_REJECTED"
  | "WAIT_REFUND_TASK"
  | "REFUND_PROCESSING"
  | "REFUNDED"
  | "CANCELED"
  | "CLOSED";

export type BuyerRefundTaskStatus = "UNSPECIFIED" | "PENDING" | "PROCESSING" | "SUCCEEDED" | "FAILED" | "DEAD_LETTER";

export type RefundTone = "warning" | "processing" | "success" | "danger" | "default";

export type BuyerRefundPreview = {
  orderNo: string;
  subOrderNo: string;
  itemNo: string;
  shopNo: string;
  userId: number;
  subStatus: string;
  orderStatus: string;
  paymentStatus: string;
  refundableAmount: number;
  pointsReturnAmount: number;
  pointsReverseAmount: number;
  paymentNo: string;
  items: BuyerOrderItem[];
};

export type BuyerRefundBatch = {
  refundBatchNo: string;
  orderNo: string;
  userId: number;
  shopNo: string;
  batchStatus: BuyerRefundBatchStatus;
  applyRefundAmount: number;
  approvedRefundAmount: number;
  caseCount: number;
  subOrderNos: string[];
  reviewDeadlineAt: string;
  autoApprovedAt: string;
  createdAt: string;
  updatedAt: string;
};

export type BuyerRefundCase = {
  afterSaleNo: string;
  refundBatchNo: string;
  orderNo: string;
  subOrderNo: string;
  itemNo: string;
  shopNo: string;
  spuNo: string;
  skuNo: string;
  qty: number;
  afterSaleStatus: BuyerRefundBatchStatus;
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

export type BuyerRefundTask = {
  refundTaskNo: string;
  afterSaleNo: string;
  orderNo: string;
  subOrderNo: string;
  payNo: string;
  refundAmount: number;
  status: BuyerRefundTaskStatus;
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

export type BuyerRefundBatchDetail = {
  batch: BuyerRefundBatch | null;
  cases: BuyerRefundCase[];
  refundTasks: BuyerRefundTask[];
};

export type AfterSaleStatus = BuyerRefundBatchStatus;
export type RefundTaskStatus = BuyerRefundTaskStatus;
export type RefundSubOrderSnapshot = BuyerRefundPreview & {
  selectedItemNos: string[];
};
export type RefundBatch = BuyerRefundBatch;
export type AfterSaleCase = BuyerRefundCase & {
  afterSaleType?: string;
  userId?: number;
  evidenceAssetIds?: number[];
  rejectReasonCode?: number;
  version?: number;
  closedAt?: string;
  cancelReasonCode?: string;
};
export type RefundTask = BuyerRefundTask & {
  nextRetryAt?: string;
};
export type RefundBatchDetail = {
  batch: BuyerRefundBatch;
  cases: AfterSaleCase[];
  refundTasks: RefundTask[];
};
export type PreviewRefundRes = {
  snapshot: RefundSubOrderSnapshot | null;
  scopeCode: string;
  scopeLabel: string;
};

export type BuyerRefundListResult = {
  list: BuyerRefundBatch[];
  nextCursor: string;
  hasMore: boolean;
};

export type BuyerRefundApplyPayload = {
  orderNo: string;
  subOrderNo: string;
  itemNo?: string;
  selectedItemNos?: string[];
  reasonCode: string;
  reasonDesc?: string;
  buyerRemark?: string;
};

export type RefundableSubOrder = BuyerOrderSub & {
  orderNo: string;
};
