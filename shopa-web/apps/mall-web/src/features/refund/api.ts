import { apiClient } from "@/lib/api-client";
import { buildIdempotencyKey } from "@/features/order/api";
import type {
  AfterSaleCase,
  AfterSaleStatus,
  PreviewRefundRes,
  RefundBatch,
  RefundBatchDetail,
  RefundSubOrderSnapshot,
  RefundTask,
  RefundTaskStatus
} from "./types";

type RawTimestamp = {
  seconds?: number | string;
  nanos?: number | string;
  value?: string;
};

type RawRefundBatch = {
  refund_batch_no?: string;
  order_no?: string;
  user_id?: number | string;
  shop_no?: string;
  batch_status?: string;
  apply_refund_amount?: number | string;
  approved_refund_amount?: number | string;
  case_count?: number | string;
  sub_order_nos?: string[];
  review_deadline_at?: RawTimestamp | string;
  auto_approved_at?: RawTimestamp | string;
  created_at?: RawTimestamp | string;
  updated_at?: RawTimestamp | string;
};

type RawAfterSaleCase = {
  after_sale_no?: string;
  order_no?: string;
  sub_order_no?: string;
  item_no?: string;
  user_id?: number | string;
  shop_no?: string;
  spu_no?: string;
  sku_no?: string;
  qty?: number | string;
  after_sale_type?: string;
  after_sale_status?: string;
  apply_refund_amount?: number | string;
  approved_refund_amount?: number | string;
  reason_code?: string;
  reason_desc?: string;
  evidence_asset_ids?: Array<number | string>;
  buyer_remark?: string;
  seller_reply?: string;
  reject_reason_code?: number | string;
  version?: number | string;
  created_at?: RawTimestamp | string;
  updated_at?: RawTimestamp | string;
  closed_at?: RawTimestamp | string;
  cancel_reason_code?: string;
  refund_batch_no?: string;
  scope_code?: string;
  review_deadline_at?: RawTimestamp | string;
  auto_approved_at?: RawTimestamp | string;
  selected_item_nos?: string[];
  payment_no?: string;
};

type RawRefundTask = {
  refund_task_no?: string;
  after_sale_no?: string;
  order_no?: string;
  sub_order_no?: string;
  pay_no?: string;
  refund_amount?: number | string;
  status?: string;
  retry_count?: number | string;
  next_retry_at?: RawTimestamp | string;
  last_error_code?: string;
  last_error_message?: string;
  created_at?: RawTimestamp | string;
  updated_at?: RawTimestamp | string;
  points_return_amount?: number | string;
  points_reverse_amount?: number | string;
  points_cash_offset_amount?: number | string;
  final_cash_refund_amount?: number | string;
  account_debt_after?: number | string;
};

type RawRefundBatchDetail = {
  batch?: RawRefundBatch;
  cases?: RawAfterSaleCase[];
  refund_tasks?: RawRefundTask[];
};

type RawOrderItemSnapshot = {
  item_no?: string;
  order_no?: string;
  sub_order_no?: string;
  shop_no?: string;
  spu_no?: string;
  sku_no?: string;
  spu_title?: string;
  sku_name?: string;
  sku_image_asset_id?: number | string;
  qty?: number | string;
  sale_price?: number | string;
  market_price?: number | string;
  sale_attrs_json?: string;
};

type RawRefundSubOrderSnapshot = {
  order_no?: string;
  sub_order_no?: string;
  item_no?: string;
  shop_no?: string;
  user_id?: number | string;
  sub_status?: string;
  refundable_amount?: number | string;
  points_return_amount?: number | string;
  points_reverse_amount?: number | string;
  payment_no?: string;
  items?: RawOrderItemSnapshot[];
  order_status?: string;
  payment_status?: string;
};

type RawPreviewRefundRes = {
  snapshot?: RawRefundSubOrderSnapshot;
  scope_code?: string;
  scope_label?: string;
};

type RawApplyRefundBatchRes = {
  detail?: RawRefundBatchDetail;
};

type RawListRefundBatchesRes = {
  list?: RawRefundBatch[];
  next_cursor?: string;
  has_more?: boolean;
};

type RawGetRefundBatchDetailRes = {
  detail?: RawRefundBatchDetail;
};

type RawCancelRefundBatchRes = {
  refund_batch_no?: string;
  batch_status?: string;
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

function mapAfterSaleStatus(value: unknown): AfterSaleStatus {
  const numeric = Number(value);
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
  const normalized = toString(value).trim().toUpperCase();
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

function mapRefundTaskStatus(value: unknown): RefundTaskStatus {
  const numeric = Number(value);
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
  const normalized = toString(value).trim().toUpperCase();
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

function normalizeRefundTask(raw?: RawRefundTask): RefundTask | null {
  if (!raw) {
    return null;
  }
  return {
    refundTaskNo: toString(raw.refund_task_no),
    afterSaleNo: toString(raw.after_sale_no),
    orderNo: toString(raw.order_no),
    subOrderNo: toString(raw.sub_order_no),
    payNo: toString(raw.pay_no),
    refundAmount: toNumber(raw.refund_amount),
    status: mapRefundTaskStatus(raw.status),
    retryCount: toNumber(raw.retry_count),
    nextRetryAt: toString(raw.next_retry_at),
    lastErrorCode: toString(raw.last_error_code),
    lastErrorMessage: toString(raw.last_error_message),
    createdAt: toString(raw.created_at),
    updatedAt: toString(raw.updated_at),
    pointsReturnAmount: toNumber(raw.points_return_amount),
    pointsReverseAmount: toNumber(raw.points_reverse_amount),
    pointsCashOffsetAmount: toNumber(raw.points_cash_offset_amount),
    finalCashRefundAmount: toNumber(raw.final_cash_refund_amount),
    accountDebtAfter: toNumber(raw.account_debt_after)
  };
}

function normalizeAfterSaleCase(raw?: RawAfterSaleCase): AfterSaleCase | null {
  if (!raw) {
    return null;
  }
  return {
    afterSaleNo: toString(raw.after_sale_no),
    orderNo: toString(raw.order_no),
    subOrderNo: toString(raw.sub_order_no),
    itemNo: toString(raw.item_no),
    userId: toNumber(raw.user_id),
    shopNo: toString(raw.shop_no),
    spuNo: toString(raw.spu_no),
    skuNo: toString(raw.sku_no),
    qty: toNumber(raw.qty),
    afterSaleType: toString(raw.after_sale_type),
    afterSaleStatus: mapAfterSaleStatus(raw.after_sale_status),
    applyRefundAmount: toNumber(raw.apply_refund_amount),
    approvedRefundAmount: toNumber(raw.approved_refund_amount),
    reasonCode: toString(raw.reason_code),
    reasonDesc: toString(raw.reason_desc),
    evidenceAssetIds: Array.isArray(raw.evidence_asset_ids) ? raw.evidence_asset_ids.filter(Boolean).map((value) => toNumber(value)) : [],
    buyerRemark: toString(raw.buyer_remark),
    sellerReply: toString(raw.seller_reply),
    rejectReasonCode: toNumber(raw.reject_reason_code),
    version: toNumber(raw.version),
    createdAt: toString(raw.created_at),
    updatedAt: toString(raw.updated_at),
    closedAt: toString(raw.closed_at),
    cancelReasonCode: toString(raw.cancel_reason_code),
    refundBatchNo: toString(raw.refund_batch_no),
    scopeCode: toString(raw.scope_code),
    reviewDeadlineAt: toString(raw.review_deadline_at),
    autoApprovedAt: toString(raw.auto_approved_at),
    selectedItemNos: Array.isArray(raw.selected_item_nos) ? raw.selected_item_nos.filter(Boolean) : [],
    paymentNo: toString(raw.payment_no)
  };
}

function normalizeRefundBatch(raw?: RawRefundBatch): RefundBatch | null {
  if (!raw) {
    return null;
  }
  return {
    refundBatchNo: toString(raw.refund_batch_no),
    orderNo: toString(raw.order_no),
    userId: toNumber(raw.user_id),
    shopNo: toString(raw.shop_no),
    batchStatus: mapAfterSaleStatus(raw.batch_status),
    applyRefundAmount: toNumber(raw.apply_refund_amount),
    approvedRefundAmount: toNumber(raw.approved_refund_amount),
    caseCount: toNumber(raw.case_count),
    subOrderNos: Array.isArray(raw.sub_order_nos) ? raw.sub_order_nos.filter(Boolean).map((value) => toString(value)) : [],
    reviewDeadlineAt: toString(raw.review_deadline_at),
    autoApprovedAt: toString(raw.auto_approved_at),
    createdAt: toString(raw.created_at),
    updatedAt: toString(raw.updated_at)
  };
}

function normalizeRefundBatchDetail(raw?: RawRefundBatchDetail): RefundBatchDetail | null {
  if (!raw || !raw.batch) {
    return null;
  }
  const batch = normalizeRefundBatch(raw.batch);
  if (!batch) {
    return null;
  }
  return {
    batch,
    cases: Array.isArray(raw.cases)
      ? raw.cases.map((item) => normalizeAfterSaleCase(item)).filter((item): item is AfterSaleCase => Boolean(item))
      : [],
    refundTasks: Array.isArray(raw.refund_tasks)
      ? raw.refund_tasks.map((item) => normalizeRefundTask(item)).filter((item): item is RefundTask => Boolean(item))
      : []
  };
}

function normalizeOrderItemSnapshot(raw?: RawOrderItemSnapshot) {
  return {
    itemNo: toString(raw?.item_no),
    orderNo: toString(raw?.order_no),
    subOrderNo: toString(raw?.sub_order_no),
    shopNo: toString(raw?.shop_no),
    spuNo: toString(raw?.spu_no),
    skuNo: toString(raw?.sku_no),
    spuTitle: toString(raw?.spu_title),
    skuName: toString(raw?.sku_name),
    skuImageAssetId: toString(raw?.sku_image_asset_id),
    qty: toNumber(raw?.qty),
    salePrice: toNumber(raw?.sale_price),
    marketPrice: toNumber(raw?.market_price),
    saleAttrsJson: toString(raw?.sale_attrs_json)
  };
}

function normalizeRefundSubOrderSnapshot(raw?: RawRefundSubOrderSnapshot): RefundSubOrderSnapshot | null {
  if (!raw) {
    return null;
  }
  return {
    orderNo: toString(raw.order_no),
    subOrderNo: toString(raw.sub_order_no),
    itemNo: toString(raw.item_no),
    shopNo: toString(raw.shop_no),
    userId: toNumber(raw.user_id),
    subStatus: toString(raw.sub_status),
    refundableAmount: toNumber(raw.refundable_amount),
    pointsReturnAmount: toNumber(raw.points_return_amount),
    pointsReverseAmount: toNumber(raw.points_reverse_amount),
    paymentNo: toString(raw.payment_no),
    items: Array.isArray(raw.items) ? raw.items.map((item) => normalizeOrderItemSnapshot(item)) : [],
    selectedItemNos: Array.isArray(raw.items) ? raw.items.map((item) => toString(item.item_no)).filter(Boolean) : [],
    orderStatus: toString(raw.order_status),
    paymentStatus: toString(raw.payment_status)
  };
}

function toProtoStatus(status: AfterSaleStatus): string | null {
  if (status === "UNSPECIFIED") {
    return null;
  }
  return `AFTER_SALE_STATUS_${status}`;
}

export async function previewRefund(payload: {
  orderNo: string;
  subOrderNo: string;
  itemNo?: string;
  selectedItemNos?: string[];
}): Promise<PreviewRefundRes> {
  const response = await apiClient.post<RawPreviewRefundRes>("/v1/aftersale/buyer/refunds:preview", {
    order_no: payload.orderNo,
    sub_order_no: payload.subOrderNo,
    item_no: payload.itemNo ?? "",
    selected_item_nos: payload.selectedItemNos ?? []
  });
  return {
    snapshot: normalizeRefundSubOrderSnapshot(response?.snapshot),
    scopeCode: toString(response?.scope_code),
    scopeLabel: toString(response?.scope_label)
  };
}

export async function applyRefundBatch(payload: {
  targets: {
    orderNo: string;
    subOrderNo: string;
    itemNo?: string;
    selectedItemNos?: string[];
  }[];
  reasonCode: string;
  reasonDesc: string;
  buyerRemark?: string;
  idempotencyKey?: string;
}): Promise<RefundBatchDetail | null> {
  const response = await apiClient.post<RawApplyRefundBatchRes>("/v1/aftersale/buyer/refunds:apply", {
    targets: payload.targets.map((target) => ({
      order_no: target.orderNo,
      sub_order_no: target.subOrderNo,
      item_no: target.itemNo ?? "",
      selected_item_nos: target.selectedItemNos ?? []
    })),
    reason_code: payload.reasonCode,
    reason_desc: payload.reasonDesc,
    buyer_remark: payload.buyerRemark ?? "",
    idempotency_key: payload.idempotencyKey ?? buildIdempotencyKey("refund_apply")
  });
  return normalizeRefundBatchDetail(response?.detail);
}

export async function listRefundBatches(payload?: { pageSize?: number; nextCursor?: string; statuses?: AfterSaleStatus[] }): Promise<{
  list: RefundBatch[];
  nextCursor: string;
  hasMore: boolean;
}> {
  const protoStatuses = payload?.statuses?.map((status) => toProtoStatus(status)).filter((value): value is string => Boolean(value)) ?? [];
  const response = await apiClient.get<RawListRefundBatchesRes>("/v1/aftersale/buyer/refund-batches", {
    params: {
      page_size: payload?.pageSize ?? 20,
      next_cursor: payload?.nextCursor ?? "",
      statuses: protoStatuses
    }
  });
  const list = Array.isArray(response?.list)
    ? response.list
        .map((item) => normalizeRefundBatch(item))
        .filter((item): item is RefundBatch => Boolean(item))
    : [];
  return {
    list,
    nextCursor: toString(response?.next_cursor),
    hasMore: Boolean(response?.has_more)
  };
}

export async function getRefundBatchDetail(refundBatchNo: string): Promise<RefundBatchDetail | null> {
  const response = await apiClient.get<RawGetRefundBatchDetailRes>(`/v1/aftersale/buyer/refund-batches/${encodeURIComponent(refundBatchNo)}`);
  return normalizeRefundBatchDetail(response?.detail);
}

export async function cancelRefundBatch(refundBatchNo: string): Promise<boolean> {
  const response = await apiClient.post<RawCancelRefundBatchRes>("/v1/aftersale/buyer/refund-batches:cancel", {
    refund_batch_no: refundBatchNo,
    idempotency_key: buildIdempotencyKey("cancel_refund_batch")
  });
  return toString(response?.batch_status) !== "" && mapAfterSaleStatus(response?.batch_status) === "CANCELED";
}
