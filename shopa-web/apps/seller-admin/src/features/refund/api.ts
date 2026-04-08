import { apiClient } from "@/lib/api-client";
import type { SellerRefundBatch, SellerRefundBatchDetail, SellerRefundBatchStatus, SellerRefundListResult, SellerRefundTask, SellerRefundTaskStatus, SellerRefundCase } from "./types";
import { normalizeSellerRefundStatus, normalizeSellerRefundTaskStatus } from "./status";

type RawTimestamp = string | { seconds?: number | string; nanos?: number | string } | null | undefined;
type RawRecord = Record<string, unknown>;

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
  return "";
}

function toNumber(value: unknown): number {
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : 0;
}

function toTimestamp(value: RawTimestamp): string {
  if (!value) {
    return "";
  }
  if (typeof value === "string") {
    return value;
  }
  const seconds = Number(value.seconds ?? 0);
  const nanos = Number(value.nanos ?? 0);
  const ms = seconds * 1000 + Math.floor(nanos / 1_000_000);
  return ms > 0 ? new Date(ms).toISOString() : "";
}

function normalizeBatchStatus(value: unknown): SellerRefundBatchStatus {
  return normalizeSellerRefundStatus(value);
}

function normalizeRefundTaskStatus(value: unknown): SellerRefundTaskStatus {
  return normalizeSellerRefundTaskStatus(value);
}

function normalizeRefundBatch(raw?: RawRecord | null): SellerRefundBatch | null {
  if (!raw) {
    return null;
  }
  const subOrderNos = Array.isArray(raw.sub_order_nos ?? raw.subOrderNos) ? ((raw.sub_order_nos ?? raw.subOrderNos) as unknown[]) : [];
  return {
    refundBatchNo: toString(raw.refund_batch_no ?? raw.refundBatchNo),
    orderNo: toString(raw.order_no ?? raw.orderNo),
    userId: toNumber(raw.user_id ?? raw.userId),
    shopNo: toString(raw.shop_no ?? raw.shopNo),
    batchStatus: normalizeBatchStatus(raw.batch_status ?? raw.batchStatus),
    applyRefundAmount: toNumber(raw.apply_refund_amount ?? raw.applyRefundAmount),
    approvedRefundAmount: toNumber(raw.approved_refund_amount ?? raw.approvedRefundAmount),
    caseCount: toNumber(raw.case_count ?? raw.caseCount),
    subOrderNos: subOrderNos.map((item) => toString(item)).filter(Boolean),
    reviewDeadlineAt: toTimestamp((raw.review_deadline_at ?? raw.reviewDeadlineAt) as RawTimestamp),
    autoApprovedAt: toTimestamp((raw.auto_approved_at ?? raw.autoApprovedAt) as RawTimestamp),
    createdAt: toTimestamp((raw.created_at ?? raw.createdAt) as RawTimestamp),
    updatedAt: toTimestamp((raw.updated_at ?? raw.updatedAt) as RawTimestamp)
  };
}

function normalizeRefundCase(raw?: RawRecord | null): SellerRefundCase | null {
  if (!raw) {
    return null;
  }
  const selectedItemNos = Array.isArray(raw.selected_item_nos ?? raw.selectedItemNos)
    ? ((raw.selected_item_nos ?? raw.selectedItemNos) as unknown[])
    : [];
  return {
    afterSaleNo: toString(raw.after_sale_no ?? raw.afterSaleNo),
    refundBatchNo: toString(raw.refund_batch_no ?? raw.refundBatchNo),
    orderNo: toString(raw.order_no ?? raw.orderNo),
    subOrderNo: toString(raw.sub_order_no ?? raw.subOrderNo),
    itemNo: toString(raw.item_no ?? raw.itemNo),
    shopNo: toString(raw.shop_no ?? raw.shopNo),
    spuNo: toString(raw.spu_no ?? raw.spuNo),
    skuNo: toString(raw.sku_no ?? raw.skuNo),
    qty: toNumber(raw.qty),
    afterSaleStatus: normalizeBatchStatus(raw.after_sale_status ?? raw.afterSaleStatus),
    applyRefundAmount: toNumber(raw.apply_refund_amount ?? raw.applyRefundAmount),
    approvedRefundAmount: toNumber(raw.approved_refund_amount ?? raw.approvedRefundAmount),
    reasonCode: toString(raw.reason_code ?? raw.reasonCode),
    reasonDesc: toString(raw.reason_desc ?? raw.reasonDesc),
    buyerRemark: toString(raw.buyer_remark ?? raw.buyerRemark),
    sellerReply: toString(raw.seller_reply ?? raw.sellerReply),
    scopeCode: toString(raw.scope_code ?? raw.scopeCode),
    selectedItemNos: selectedItemNos.map((item) => toString(item)).filter(Boolean),
    paymentNo: toString(raw.payment_no ?? raw.paymentNo),
    reviewDeadlineAt: toTimestamp((raw.review_deadline_at ?? raw.reviewDeadlineAt) as RawTimestamp),
    autoApprovedAt: toTimestamp((raw.auto_approved_at ?? raw.autoApprovedAt) as RawTimestamp),
    createdAt: toTimestamp((raw.created_at ?? raw.createdAt) as RawTimestamp),
    updatedAt: toTimestamp((raw.updated_at ?? raw.updatedAt) as RawTimestamp)
  };
}

function normalizeRefundTask(raw?: RawRecord | null): SellerRefundTask | null {
  if (!raw) {
    return null;
  }
  return {
    refundTaskNo: toString(raw.refund_task_no ?? raw.refundTaskNo),
    afterSaleNo: toString(raw.after_sale_no ?? raw.afterSaleNo),
    orderNo: toString(raw.order_no ?? raw.orderNo),
    subOrderNo: toString(raw.sub_order_no ?? raw.subOrderNo),
    payNo: toString(raw.pay_no ?? raw.payNo),
    refundAmount: toNumber(raw.refund_amount ?? raw.refundAmount),
    status: normalizeRefundTaskStatus(raw.status),
    retryCount: toNumber(raw.retry_count ?? raw.retryCount),
    lastErrorCode: toString(raw.last_error_code ?? raw.lastErrorCode),
    lastErrorMessage: toString(raw.last_error_message ?? raw.lastErrorMessage),
    createdAt: toTimestamp((raw.created_at ?? raw.createdAt) as RawTimestamp),
    updatedAt: toTimestamp((raw.updated_at ?? raw.updatedAt) as RawTimestamp),
    pointsReturnAmount: toNumber(raw.points_return_amount ?? raw.pointsReturnAmount),
    pointsReverseAmount: toNumber(raw.points_reverse_amount ?? raw.pointsReverseAmount),
    pointsCashOffsetAmount: toNumber(raw.points_cash_offset_amount ?? raw.pointsCashOffsetAmount),
    finalCashRefundAmount: toNumber(raw.final_cash_refund_amount ?? raw.finalCashRefundAmount),
    accountDebtAfter: toNumber(raw.account_debt_after ?? raw.accountDebtAfter)
  };
}

function normalizeRefundDetail(raw?: RawRecord | null): SellerRefundBatchDetail {
  const detail = (raw?.detail ?? raw) as RawRecord | null | undefined;
  return {
    batch: normalizeRefundBatch((detail?.batch ?? null) as RawRecord | null),
    cases: Array.isArray(detail?.cases) ? (detail.cases as RawRecord[]).map((item) => normalizeRefundCase(item)).filter(Boolean) as SellerRefundCase[] : [],
    refundTasks: Array.isArray(detail?.refund_tasks ?? detail?.refundTasks)
      ? ((detail?.refund_tasks ?? detail?.refundTasks ?? []) as RawRecord[]).map((item) => normalizeRefundTask(item)).filter(Boolean) as SellerRefundTask[]
      : []
  };
}

export async function listSellerRefundBatches(shopNo: string, statuses: string[] = []): Promise<SellerRefundListResult> {
  const response = await apiClient.get<RawRecord>(`/v1/aftersale/seller/shops/${encodeURIComponent(shopNo)}/refund-batches`, {
    params: {
      page_size: 100,
      next_cursor: "",
      statuses
    }
  });
  return {
    list: Array.isArray(response?.list) ? (response.list as RawRecord[]).map((item) => normalizeRefundBatch(item)).filter(Boolean) as SellerRefundBatch[] : [],
    nextCursor: toString(response?.next_cursor ?? response?.nextCursor),
    hasMore: Boolean(response?.has_more ?? response?.hasMore)
  };
}

export async function getSellerRefundBatchDetail(refundBatchNo: string, shopNo: string): Promise<SellerRefundBatchDetail> {
  const response = await apiClient.get<RawRecord>(`/v1/aftersale/seller/refund-batches/${encodeURIComponent(refundBatchNo)}`, {
    params: { shopNo }
  });
  return normalizeRefundDetail(response);
}

export async function approveSellerRefundBatch(refundBatchNo: string, shopNo: string, sellerReply = ""): Promise<SellerRefundBatchDetail> {
  const response = await apiClient.post<RawRecord>("/v1/aftersale/seller/refund-batches:approve", {
    refund_batch_no: refundBatchNo,
    shop_no: shopNo,
    seller_reply: sellerReply
  });
  return normalizeRefundDetail(response);
}

export async function rejectSellerRefundBatch(refundBatchNo: string, shopNo: string, sellerReply = ""): Promise<SellerRefundBatchDetail> {
  const response = await apiClient.post<RawRecord>("/v1/aftersale/seller/refund-batches:reject", {
    refund_batch_no: refundBatchNo,
    shop_no: shopNo,
    reject_reason_code: 4,
    seller_reply: sellerReply
  });
  return normalizeRefundDetail(response);
}
