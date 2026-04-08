import type { BuyerRefundBatchStatus, RefundTone } from "./types";
import type { BuyerOrderSub } from "../order/types";

function normalizeRefundStatus(status: string): BuyerRefundBatchStatus {
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

export function getBuyerRefundStatusMeta(status: string): { label: string; tone: RefundTone } {
  const code = normalizeRefundStatus(status);
  switch (code) {
    case "PENDING_SELLER_REVIEW":
      return { label: "待商家审核", tone: "warning" };
    case "SELLER_REJECTED":
      return { label: "商家已驳回", tone: "danger" };
    case "WAIT_REFUND_TASK":
      return { label: "待退款", tone: "processing" };
    case "REFUND_PROCESSING":
      return { label: "退款处理中", tone: "processing" };
    case "REFUNDED":
      return { label: "已退款", tone: "success" };
    case "CANCELED":
      return { label: "已撤销", tone: "default" };
    case "CLOSED":
      return { label: "已关闭", tone: "default" };
    default:
      return { label: "处理中", tone: "default" };
  }
}

export function canApplyRefundForSubOrder(subOrder: Pick<BuyerOrderSub, "subStatus" | "items">): boolean {
  const normalized = String(subOrder.subStatus ?? "").trim().toUpperCase();
  if (!Array.isArray(subOrder.items) || subOrder.items.length === 0) {
    return false;
  }
  return normalized === "WAIT_SHIP" || normalized === "PAID";
}
