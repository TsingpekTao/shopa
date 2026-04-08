import type { BuyerOrder, SubOrderStatusCode } from "@/features/order/types";
import type { AfterSaleStatus } from "./types";

export const REFUNDABLE_SUB_ORDER_STATUSES: SubOrderStatusCode[] = ["PAID", "WAIT_SHIP"];

const statusLabels: Record<AfterSaleStatus, { zh: string; en: string }> = {
  UNSPECIFIED: { zh: "处理中", en: "Processing" },
  PENDING_SELLER_REVIEW: { zh: "待商家审核", en: "Pending seller review" },
  SELLER_REJECTED: { zh: "商家已驳回", en: "Seller rejected" },
  WAIT_REFUND_TASK: { zh: "待退款", en: "Waiting for refund" },
  REFUND_PROCESSING: { zh: "退款中", en: "Refund processing" },
  REFUNDED: { zh: "已退款", en: "Refunded" },
  CANCELED: { zh: "已取消", en: "Canceled" },
  CLOSED: { zh: "已关闭", en: "Closed" }
};

const reasonLabels: Record<string, { zh: string; en: string }> = {
  changed_mind: { zh: "不想要了", en: "Changed my mind" },
  found_a_better_price: { zh: "发现更优惠的价格", en: "Found a better price" },
  shipping_delay: { zh: "发货时间太久", en: "Shipping delay" }
};

const cancelableStatuses: AfterSaleStatus[] = ["PENDING_SELLER_REVIEW", "WAIT_REFUND_TASK"];

export function formatAfterSaleStatusLabel(status: AfterSaleStatus, isZh: boolean): string {
  const label = statusLabels[status] ?? statusLabels.UNSPECIFIED;
  return isZh ? label.zh : label.en;
}

export function getRefundReasonLabel(reasonCode: string, isZh: boolean): string {
  const normalizedCode = reasonCode.trim().toLowerCase();
  const label = reasonLabels[normalizedCode];
  if (label) {
    return isZh ? label.zh : label.en;
  }
  return isZh ? "其他原因" : "Other reason";
}

export function formatRefundReason(reasonCode: string, reasonDesc: string, isZh: boolean): string {
  const normalizedDesc = reasonDesc.trim();
  if (!reasonCode.trim()) {
    return normalizedDesc || (isZh ? "未填写原因" : "No reason provided");
  }
  const mappedLabel = getRefundReasonLabel(reasonCode, isZh);
  return mappedLabel;
}

export function hasRefundableSubOrder(order: BuyerOrder): boolean {
  return order.subOrders.some((sub) => REFUNDABLE_SUB_ORDER_STATUSES.includes(sub.subStatus));
}

export function isRefundBatchCancelable(status: AfterSaleStatus): boolean {
  return cancelableStatuses.includes(status);
}
