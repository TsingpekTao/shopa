"use client";

import Link from "next/link";
import { useMemo } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { message } from "antd";
import { useI18n } from "@shopa/ui";
import { formatCnyFromCents } from "@/lib/price";
import { listRefundBatches, cancelRefundBatch } from "@/features/refund/api";
import { formatAfterSaleStatusLabel, isRefundBatchCancelable } from "@/features/refund/helpers";

function formatTime(raw: string, locale: string): string {
  if (!raw) {
    return "--";
  }
  const date = new Date(raw);
  if (Number.isNaN(date.getTime())) {
    return raw;
  }
  return date.toLocaleString(locale === "zh-CN" ? "zh-CN" : "en-US", { hour12: false });
}

export default function RefundsPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [messageApi, contextHolder] = message.useMessage();
  const queryClient = useQueryClient();

  const refundsQuery = useQuery({
    queryKey: ["refund-batches"],
    queryFn: () => listRefundBatches({ pageSize: 50 }),
    staleTime: 30_000,
    refetchOnWindowFocus: false
  });

  const cancelMutation = useMutation({
    mutationFn: (refundBatchNo: string) => cancelRefundBatch(refundBatchNo),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["refund-batches"] });
      messageApi.success(isZh ? "退款申请已取消" : "Refund request canceled");
    },
    onError: (error) => {
      const errorMessage = error instanceof Error ? error.message : isZh ? "取消失败" : "Failed to cancel";
      messageApi.error(errorMessage);
    }
  });

  const refundList = useMemo(() => {
    return (refundsQuery.data?.list ?? []).filter((batch) => batch.batchStatus !== "REFUNDED");
  }, [refundsQuery.data?.list]);

  return (
    <section className="tb-refund-page">
      {contextHolder}

      <header className="tb-refund-page-header tb-refund-panel">
        <div>
          <p className="tb-refund-subhead">{isZh ? "退款/售后" : "After-sale"}</p>
          <h1>{isZh ? "退款处理中" : "Active refunds"}</h1>
          <p>{isZh ? "这里只展示仍在处理中的退款申请，已退款订单会进入历史订单。" : "Only active refund requests are shown here. Refunded orders move to order history."}</p>
        </div>
        <Link href="/me/refunds/apply" className="tb-order-btn tb-order-btn-primary">
          {isZh ? "申请退款" : "Request Refund"}
        </Link>
      </header>

      {refundsQuery.isLoading ? (
        <p className="tb-chat-empty">{isZh ? "正在加载退款记录..." : "Loading refund batches..."}</p>
      ) : refundsQuery.isError ? (
        <p className="tb-chat-empty">{isZh ? "刷新失败，请重试。" : "Failed to load refunds."}</p>
      ) : refundList.length === 0 ? (
        <section className="tb-order-empty-card">
          <strong>{isZh ? "暂无进行中的退款" : "No active refunds"}</strong>
          <p>{isZh ? "已退款订单会在历史订单里查看，新的退款申请可从订单页发起。" : "Refunded orders move to order history. New requests can be started from the order page."}</p>
        </section>
      ) : null}

      <div className="tb-refund-list">
        {refundList.map((batch) => {
          const statusLabel = formatAfterSaleStatusLabel(batch.batchStatus, isZh);
          const createdTime = formatTime(batch.createdAt, locale);
          const reviewDeadline = batch.reviewDeadlineAt ? formatTime(batch.reviewDeadlineAt, locale) : "";
          const refundHref = `/me/refunds/${encodeURIComponent(batch.refundBatchNo)}`;
          const canCancel = isRefundBatchCancelable(batch.batchStatus);
          return (
            <article key={batch.refundBatchNo} className="tb-refund-card">
              <div className="tb-refund-card-head">
                <div>
                  <strong>{isZh ? `退款单 ${batch.refundBatchNo}` : `Refund ${batch.refundBatchNo}`}</strong>
                  <p className="tb-refund-card-meta">{`${isZh ? "创建时间" : "Created"} · ${createdTime}`}</p>
                </div>
                <span className="tb-refund-badge">{statusLabel}</span>
              </div>

              <div className="tb-refund-card-body">
                <div>
                  <p>{isZh ? "申请金额" : "Requested"}</p>
                  <strong className="is-price">{`CNY ${formatCnyFromCents(batch.applyRefundAmount)}`}</strong>
                  <small>
                    {isZh ? "已审批" : "Approved"} {formatCnyFromCents(batch.approvedRefundAmount)}
                  </small>
                </div>
                <div>
                  <p>{isZh ? "关联子单" : "Sub-orders"}</p>
                  <strong>{batch.subOrderNos.length}</strong>
                  <small>{isZh ? `案件数 ${batch.caseCount}` : `${batch.caseCount} cases`}</small>
                </div>
                {reviewDeadline ? (
                  <div>
                    <p>{isZh ? "审核截止" : "Review deadline"}</p>
                    <strong>{reviewDeadline}</strong>
                  </div>
                ) : null}
              </div>

              <div className="tb-refund-card-actions">
                <Link href={refundHref} className="tb-order-btn tb-order-btn-secondary">
                  {isZh ? "查看详情" : "View Detail"}
                </Link>
                {canCancel ? (
                  <button
                    type="button"
                    className="tb-order-btn tb-order-btn-secondary"
                    disabled={cancelMutation.isPending && cancelMutation.variables === batch.refundBatchNo}
                    onClick={() => cancelMutation.mutate(batch.refundBatchNo)}
                  >
                    {cancelMutation.isPending && cancelMutation.variables === batch.refundBatchNo
                      ? isZh
                        ? "取消中..."
                        : "Canceling..."
                      : isZh
                        ? "取消退款"
                        : "Cancel Refund"}
                  </button>
                ) : null}
              </div>
            </article>
          );
        })}
      </div>
    </section>
  );
}
