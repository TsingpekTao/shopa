"use client";

import Link from "next/link";
import { useMemo } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { message } from "antd";
import { useI18n } from "@shopa/ui";
import { formatCnyFromCents } from "@/lib/price";
import { getRefundBatchDetail, cancelRefundBatch } from "@/features/refund/api";
import { formatAfterSaleStatusLabel, formatRefundReason, isRefundBatchCancelable } from "@/features/refund/helpers";

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

type RefundDetailPageProps = {
  params: {
    refundBatchNo: string;
  };
};

export default function RefundDetailPage({ params }: RefundDetailPageProps) {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [messageApi, contextHolder] = message.useMessage();
  const queryClient = useQueryClient();

  const detailQuery = useQuery({
    queryKey: ["refund-batch-detail", params.refundBatchNo],
    queryFn: () => getRefundBatchDetail(params.refundBatchNo),
    enabled: Boolean(params.refundBatchNo),
    staleTime: 60_000
  });

  const cancelMutation = useMutation({
    mutationFn: () => cancelRefundBatch(params.refundBatchNo),
    onSuccess: async () => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ["refund-batches"] }),
        detailQuery.refetch()
      ]);
      messageApi.success(isZh ? "退款已取消" : "Refund canceled");
    },
    onError: (error) => {
      const errorMessage = error instanceof Error ? error.message : isZh ? "取消失败" : "Failed to cancel";
      messageApi.error(errorMessage);
    }
  });

  const detail = detailQuery.data;
  const batchStatusLabel = detail ? formatAfterSaleStatusLabel(detail.batch.batchStatus, isZh) : "";
  const showCases = useMemo(() => detail?.cases ?? [], [detail]);
  const showTasks = useMemo(() => detail?.refundTasks ?? [], [detail]);

  return (
    <section className="tb-refund-page">
      {contextHolder}

      <header className="tb-refund-detail-header tb-refund-panel">
        <div>
          <p className="tb-refund-subhead">{isZh ? "退款详情" : "Refund detail"}</p>
          <h1>{isZh ? `退款单 ${params.refundBatchNo}` : `Refund ${params.refundBatchNo}`}</h1>
          <p>{isZh ? "查看退款申请、金额和处理进度。" : "Track refund cases, amounts, and processing progress."}</p>
        </div>
        <div className="tb-refund-header-actions">
          <Link href="/me/refunds" className="tb-order-btn tb-order-btn-secondary">
            {isZh ? "返回列表" : "Back to list"}
          </Link>
          <Link href="/me/orders?tab=history" className="tb-order-btn tb-order-btn-primary">
            {isZh ? "查看历史订单" : "View order history"}
          </Link>
        </div>
      </header>

      {detailQuery.isLoading ? (
        <p className="tb-chat-empty">{isZh ? "正在加载退款详情..." : "Loading refund detail..."}</p>
      ) : !detail ? (
        <p className="tb-chat-empty">{isZh ? "未找到该退款批次。" : "Refund not found."}</p>
      ) : (
        <>
          <article className="tb-refund-card tb-refund-detail-card">
            <div className="tb-refund-card-head">
              <div>
                <strong>{batchStatusLabel}</strong>
                <p className="tb-refund-card-meta">{isZh ? "当前状态" : "Current status"}</p>
              </div>
              <div>
                <p>{isZh ? "申请金额" : "Requested"}</p>
                <strong className="is-price">{`CNY ${formatCnyFromCents(detail.batch.applyRefundAmount)}`}</strong>
                <small>
                  {isZh ? "已审批" : "Approved"} {formatCnyFromCents(detail.batch.approvedRefundAmount)}
                </small>
              </div>
              <div>
                <p>{isZh ? "创建时间" : "Created"}</p>
                <strong>{formatTime(detail.batch.createdAt, locale)}</strong>
              </div>
            </div>

            <div className="tb-refund-card-actions">
              {isRefundBatchCancelable(detail.batch.batchStatus) ? (
                <button
                  type="button"
                  className="tb-order-btn tb-order-btn-secondary"
                  disabled={cancelMutation.isPending}
                  onClick={() => cancelMutation.mutate()}
                >
                  {cancelMutation.isPending ? (isZh ? "取消中..." : "Canceling...") : isZh ? "取消退款" : "Cancel Refund"}
                </button>
              ) : null}
            </div>
          </article>

          <section className="tb-refund-section">
            <header className="tb-refund-section-head">
              <strong>{isZh ? "售后案件" : "After-sale cases"}</strong>
              <span>{showCases.length} {isZh ? "条" : "items"}</span>
            </header>
            <div className="tb-refund-case-list">
              {showCases.map((refundCase) => (
                <article key={refundCase.afterSaleNo} className="tb-refund-case-card">
                  <div className="tb-refund-case-head">
                    <strong>{isZh ? `案件 ${refundCase.afterSaleNo}` : `Case ${refundCase.afterSaleNo}`}</strong>
                    <span>{formatAfterSaleStatusLabel(refundCase.afterSaleStatus, isZh)}</span>
                  </div>
                  <p className="tb-refund-card-meta">{`${refundCase.orderNo} · ${refundCase.subOrderNo}`}</p>
                  <p>{formatRefundReason(refundCase.reasonCode, refundCase.reasonDesc, isZh)}</p>
                  {refundCase.buyerRemark ? <small>{`${isZh ? "补充说明" : "Additional note"}：${refundCase.buyerRemark}`}</small> : null}
                  <div className="tb-refund-case-meta">
                    <span>
                      {isZh ? "申请金额" : "Requested"}: {formatCnyFromCents(refundCase.applyRefundAmount)}
                    </span>
                    <span>
                      {isZh ? "审批金额" : "Approved"}: {formatCnyFromCents(refundCase.approvedRefundAmount)}
                    </span>
                  </div>
                </article>
              ))}
            </div>
          </section>

          <section className="tb-refund-section">
            <header className="tb-refund-section-head">
              <strong>{isZh ? "退款任务" : "Refund tasks"}</strong>
              <span>{showTasks.length} {isZh ? "步" : "steps"}</span>
            </header>
            <div className="tb-refund-task-list">
              {showTasks.map((task) => (
                <div key={task.refundTaskNo} className="tb-refund-task-row">
                  <div>
                    <strong>{task.status}</strong>
                    <p className="tb-refund-card-meta">{formatTime(task.createdAt, locale)}</p>
                  </div>
                  <div>
                    <p>{isZh ? "金额" : "Amount"}</p>
                    <strong>{formatCnyFromCents(task.refundAmount)}</strong>
                  </div>
                  {task.lastErrorMessage ? (
                    <div className="tb-refund-task-error">
                      <p>{isZh ? "错误信息" : "Error"}</p>
                      <small>{task.lastErrorMessage}</small>
                    </div>
                  ) : null}
                </div>
              ))}
            </div>
          </section>
        </>
      )}
    </section>
  );
}
