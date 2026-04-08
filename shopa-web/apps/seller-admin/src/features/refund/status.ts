import type { RefundTone, SellerRefundBatch, SellerRefundBatchStatus, SellerRefundTaskStatus } from "./types";

export function normalizeSellerRefundStatus(status: unknown): SellerRefundBatchStatus {
  const numeric = Number(status);
  if (Number.isFinite(numeric)) {
    switch (numeric) {
      case 1:
        return "PENDING_SELLER_REVIEW";
      case 2:
        return "SELLER_REJECTED";
      case 3:
        return "WAIT_REFUND_TASK";
      case 4:
        return "REFUND_PROCESSING";
      case 5:
        return "REFUNDED";
      case 6:
        return "CANCELED";
      case 7:
        return "CLOSED";
      default:
        return "UNSPECIFIED";
    }
  }

  const normalized = String(status ?? "").trim().toUpperCase();
  switch (normalized) {
    case "AFTER_SALE_STATUS_PENDING_SELLER_REVIEW":
    case "PENDING_SELLER_REVIEW":
      return "PENDING_SELLER_REVIEW";
    case "AFTER_SALE_STATUS_SELLER_REJECTED":
    case "SELLER_REJECTED":
      return "SELLER_REJECTED";
    case "AFTER_SALE_STATUS_WAIT_REFUND_TASK":
    case "WAIT_REFUND_TASK":
      return "WAIT_REFUND_TASK";
    case "AFTER_SALE_STATUS_REFUND_PROCESSING":
    case "REFUND_PROCESSING":
      return "REFUND_PROCESSING";
    case "AFTER_SALE_STATUS_REFUNDED":
    case "REFUNDED":
      return "REFUNDED";
    case "AFTER_SALE_STATUS_CANCELED":
    case "CANCELED":
      return "CANCELED";
    case "AFTER_SALE_STATUS_CLOSED":
    case "CLOSED":
      return "CLOSED";
    default:
      return "UNSPECIFIED";
  }
}

export function normalizeSellerRefundTaskStatus(status: unknown): SellerRefundTaskStatus {
  const numeric = Number(status);
  if (Number.isFinite(numeric)) {
    switch (numeric) {
      case 1:
        return "PENDING";
      case 2:
        return "PROCESSING";
      case 3:
        return "SUCCEEDED";
      case 4:
        return "FAILED";
      case 5:
        return "DEAD_LETTER";
      default:
        return "UNSPECIFIED";
    }
  }

  const normalized = String(status ?? "").trim().toUpperCase();
  switch (normalized) {
    case "REFUND_TASK_STATUS_PENDING":
    case "PENDING":
      return "PENDING";
    case "REFUND_TASK_STATUS_PROCESSING":
    case "PROCESSING":
      return "PROCESSING";
    case "REFUND_TASK_STATUS_SUCCEEDED":
    case "SUCCEEDED":
      return "SUCCEEDED";
    case "REFUND_TASK_STATUS_FAILED":
    case "FAILED":
      return "FAILED";
    case "REFUND_TASK_STATUS_DEAD_LETTER":
    case "DEAD_LETTER":
      return "DEAD_LETTER";
    default:
      return "UNSPECIFIED";
  }
}

const sellerRefundPriority: Record<SellerRefundBatchStatus, number> = {
  PENDING_SELLER_REVIEW: 0,
  WAIT_REFUND_TASK: 1,
  REFUND_PROCESSING: 2,
  REFUNDED: 3,
  SELLER_REJECTED: 4,
  CANCELED: 5,
  CLOSED: 6,
  UNSPECIFIED: 7
};

const sellerRefundReasonLabels: Record<string, string> = {
  changed_mind: "不想要了",
  found_a_better_price: "发现更优惠的价格",
  shipping_delay: "发货时间太久"
};

export function getSellerRefundStatusMeta(status: unknown): { label: string; tone: RefundTone } {
  const code = normalizeSellerRefundStatus(status);
  switch (code) {
    case "PENDING_SELLER_REVIEW":
      return { label: "待商家审核", tone: "warning" };
    case "WAIT_REFUND_TASK":
      return { label: "待退款", tone: "processing" };
    case "REFUND_PROCESSING":
      return { label: "退款处理中", tone: "processing" };
    case "REFUNDED":
      return { label: "已退款", tone: "success" };
    case "SELLER_REJECTED":
      return { label: "已驳回", tone: "danger" };
    case "CANCELED":
      return { label: "已取消", tone: "default" };
    case "CLOSED":
      return { label: "已关闭", tone: "default" };
    default:
      return { label: "处理中", tone: "default" };
  }
}

export function getSellerRefundTaskStatusLabel(status: unknown): string {
  switch (normalizeSellerRefundTaskStatus(status)) {
    case "PENDING":
      return "待执行";
    case "PROCESSING":
      return "执行中";
    case "SUCCEEDED":
      return "已完成";
    case "FAILED":
      return "执行失败";
    case "DEAD_LETTER":
      return "重试失败";
    default:
      return "处理中";
  }
}

export function formatSellerRefundReasonLabel(reasonCode: string, reasonDesc: string): string {
  const normalizedCode = reasonCode.trim().toLowerCase();
  if (normalizedCode && sellerRefundReasonLabels[normalizedCode]) {
    return sellerRefundReasonLabels[normalizedCode];
  }
  return reasonDesc.trim() || "买家未填写";
}

function toTimeMs(value: string): number {
  const ms = new Date(value).getTime();
  return Number.isFinite(ms) ? ms : 0;
}

export function sortRefundBatchesForSeller(list: SellerRefundBatch[]): SellerRefundBatch[] {
  return [...list].sort((left, right) => {
    const priorityDiff = sellerRefundPriority[left.batchStatus] - sellerRefundPriority[right.batchStatus];
    if (priorityDiff !== 0) {
      return priorityDiff;
    }
    const timeDiff = toTimeMs(right.createdAt) - toTimeMs(left.createdAt);
    if (timeDiff !== 0) {
      return timeDiff;
    }
    return right.refundBatchNo.localeCompare(left.refundBatchNo);
  });
}
