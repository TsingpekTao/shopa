"use client";

import { useEffect, useMemo, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Alert, Button, Empty, Input, Modal, Select, Skeleton, message } from "antd";
import { buildSellerChatShopOptions } from "@/features/chat/shop-selection";
import { getSellerCurrentShopNo, setSellerCurrentShopNo } from "@/features/navigation/nav-badges";
import {
  approveSellerRefundBatch,
  getSellerRefundBatchDetail,
  listSellerRefundBatches,
  rejectSellerRefundBatch
} from "@/features/refund/api";
import {
  formatSellerRefundReasonLabel,
  getSellerRefundStatusMeta,
  getSellerRefundTaskStatusLabel,
  sortRefundBatchesForSeller
} from "@/features/refund/status";
import { fetchSellerWorkbench } from "@/features/seller-shop/api";
import { listSellerProducts } from "@/features/catalog/api";

const SELLER_AFTER_SALES_LAST_SHOP_KEY = "seller-after-sales-last-shop-no";

const FILTER_OPTIONS = [
  { key: "all", label: "全部" },
  { key: "PENDING_SELLER_REVIEW", label: "待审核" },
  { key: "WAIT_REFUND_TASK", label: "待退款" },
  { key: "REFUND_PROCESSING", label: "退款中" },
  { key: "REFUNDED", label: "已退款" },
  { key: "SELLER_REJECTED", label: "已驳回" }
] as const;

function formatTime(raw?: string): string {
  if (!raw) {
    return "--";
  }
  const date = new Date(raw);
  if (Number.isNaN(date.getTime())) {
    return raw;
  }
  return date.toLocaleString("zh-CN", { hour12: false });
}

function formatCnyFromCents(amount: number): string {
  return (Math.max(0, amount) / 100).toFixed(2);
}

function getStatuses(filterKey: string): string[] {
  return filterKey === "all" ? [] : [filterKey];
}

function getReviewHint(batchStatus: string): string {
  return batchStatus === "PENDING_SELLER_REVIEW" ? "填写给买家的说明后即可处理退款申请。" : "当前退款状态无需审核。";
}

export default function SellerAfterSalesPage() {
  const queryClient = useQueryClient();
  const [messageApi, contextHolder] = message.useMessage();
  const [selectedShopNo, setSelectedShopNo] = useState("");
  const [persistedShopNo, setPersistedShopNo] = useState("");
  const [filterKey, setFilterKey] = useState<(typeof FILTER_OPTIONS)[number]["key"]>("all");
  const [activeBatchNo, setActiveBatchNo] = useState("");
  const [detailOpen, setDetailOpen] = useState(false);
  const [sellerReply, setSellerReply] = useState("");

  const workbenchQuery = useQuery({
    queryKey: ["seller-after-sales-workbench"],
    queryFn: fetchSellerWorkbench,
    staleTime: 60_000
  });

  useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }
    setPersistedShopNo(getSellerCurrentShopNo() || window.localStorage.getItem(SELLER_AFTER_SALES_LAST_SHOP_KEY)?.trim() || "");
  }, []);

  const workbenchShops = useMemo(() => {
    return (workbenchQuery.data?.shops ?? [])
      .map((shop) => ({
        shopNo: String(shop.shopNo ?? "").trim(),
        shopName: String(shop.shopDisplayName || shop.shopName || shop.shopNo || "").trim(),
        shopStatusCode: String(shop.shopStatusCode ?? "").trim()
      }))
      .filter((shop) => shop.shopNo);
  }, [workbenchQuery.data?.shops]);

  const fallbackCatalogShopsQuery = useQuery({
    queryKey: ["seller-after-sales-fallback-shops"],
    queryFn: async () => {
      const result = await listSellerProducts({ page: 1, pageSize: 100 });
      return Array.from(new Set(result.products.map((product) => product.shopNo.trim()).filter(Boolean)));
    },
    enabled: workbenchShops.length === 0 || Boolean(workbenchQuery.data?.partial),
    staleTime: 60_000,
    refetchOnWindowFocus: false
  });

  const shops = useMemo(() => {
    return buildSellerChatShopOptions({
      workbenchShops,
      catalogShopNos: fallbackCatalogShopsQuery.data ?? [],
      persistedShopNo
    });
  }, [fallbackCatalogShopsQuery.data, persistedShopNo, workbenchShops]);

  const normalizedSelectedShopNo = selectedShopNo.trim();

  useEffect(() => {
    if (!shops.length) {
      setSelectedShopNo("");
      return;
    }
    setSelectedShopNo((current) => {
      const normalizedCurrent = current.trim();
      if (normalizedCurrent && shops.some((shop) => shop.shopNo === normalizedCurrent)) {
        return normalizedCurrent;
      }
      return shops[0].shopNo;
    });
  }, [shops]);

  useEffect(() => {
    if (!normalizedSelectedShopNo || typeof window === "undefined") {
      return;
    }
    setSellerCurrentShopNo(normalizedSelectedShopNo);
    window.localStorage.setItem(SELLER_AFTER_SALES_LAST_SHOP_KEY, normalizedSelectedShopNo);
    setPersistedShopNo(normalizedSelectedShopNo);
  }, [normalizedSelectedShopNo]);

  useEffect(() => {
    setDetailOpen(false);
    setActiveBatchNo("");
    setSellerReply("");
  }, [normalizedSelectedShopNo]);

  const refundListQuery = useQuery({
    queryKey: ["seller-refund-batches", normalizedSelectedShopNo, filterKey],
    queryFn: () => listSellerRefundBatches(normalizedSelectedShopNo, getStatuses(filterKey)),
    enabled: Boolean(normalizedSelectedShopNo),
    staleTime: 10_000,
    refetchOnWindowFocus: false
  });

  const refundBatches = useMemo(() => sortRefundBatchesForSeller(refundListQuery.data?.list ?? []), [refundListQuery.data?.list]);

  const refundStats = useMemo(() => {
    return refundBatches.reduce(
      (summary, batch) => {
        summary.total += 1;
        summary.amount += batch.applyRefundAmount;
        if (batch.batchStatus === "PENDING_SELLER_REVIEW") {
          summary.pending += 1;
        }
        if (batch.batchStatus === "WAIT_REFUND_TASK" || batch.batchStatus === "REFUND_PROCESSING") {
          summary.processing += 1;
        }
        if (batch.batchStatus === "REFUNDED") {
          summary.completed += 1;
        }
        return summary;
      },
      { total: 0, pending: 0, processing: 0, completed: 0, amount: 0 }
    );
  }, [refundBatches]);

  const detailQuery = useQuery({
    queryKey: ["seller-refund-batch-detail", normalizedSelectedShopNo, activeBatchNo],
    queryFn: () => getSellerRefundBatchDetail(activeBatchNo, normalizedSelectedShopNo),
    enabled: Boolean(normalizedSelectedShopNo && activeBatchNo && detailOpen),
    staleTime: 10_000,
    refetchOnWindowFocus: false
  });

  useEffect(() => {
    if (!detailOpen) {
      return;
    }
    setSellerReply("");
  }, [activeBatchNo, detailOpen]);

  useEffect(() => {
    const existingReply = detailQuery.data?.cases.find((item) => item.sellerReply)?.sellerReply ?? "";
    if (existingReply) {
      setSellerReply(existingReply);
    }
  }, [detailQuery.data?.batch?.refundBatchNo]);

  const reviewMutation = useMutation({
    mutationFn: async (action: "approve" | "reject") => {
      if (!activeBatchNo || !normalizedSelectedShopNo) {
        throw new Error("退款批次信息不完整");
      }
      if (action === "approve") {
        return approveSellerRefundBatch(activeBatchNo, normalizedSelectedShopNo, sellerReply);
      }
      return rejectSellerRefundBatch(activeBatchNo, normalizedSelectedShopNo, sellerReply);
    },
    onSuccess: async (_, action) => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ["seller-refund-batches", normalizedSelectedShopNo] }),
        queryClient.invalidateQueries({ queryKey: ["seller-refund-batch-detail", normalizedSelectedShopNo, activeBatchNo] })
      ]);
      messageApi.success(action === "approve" ? "已同意退款申请" : "已驳回退款申请");
    },
    onError: (error) => {
      messageApi.error(error instanceof Error ? error.message : "售后处理失败");
    }
  });

  const selectedDetail = detailQuery.data;
  const selectedBatch = selectedDetail?.batch ?? null;
  const selectedBatchStatus = selectedBatch?.batchStatus ?? "UNSPECIFIED";
  const canReview = selectedBatchStatus === "PENDING_SELLER_REVIEW";

  return (
    <section className="seller-page seller-after-sales-page">
      {contextHolder}

      <header className="seller-after-sales-head">
        <div className="seller-after-sales-head-copy">
          <span>售后中心</span>
          <h1>退款列表</h1>
          <p>按店铺查看退款申请，待审核的批次会优先展示。</p>
        </div>
        <div className="seller-after-sales-head-actions">
          <Select
            value={normalizedSelectedShopNo || undefined}
            placeholder="选择店铺"
            style={{ minWidth: 240 }}
            loading={workbenchQuery.isLoading}
            options={shops.map((shop) => ({
              label: shop.shopName,
              value: shop.shopNo
            }))}
            onChange={(value) => setSelectedShopNo(value.trim())}
          />
        </div>
      </header>

      {refundListQuery.isError ? (
        <Alert
          type="error"
          showIcon
          message="退款列表加载失败"
          description={refundListQuery.error instanceof Error ? refundListQuery.error.message : "请稍后重试"}
        />
      ) : null}

      {!shops.length && !workbenchQuery.isLoading ? (
        <div className="seller-after-sales-empty">
          <Empty description="当前账号下还没有可处理售后的店铺" />
        </div>
      ) : (
        <>
          <section className="seller-after-sales-summary">
            <article>
              <span>当前列表</span>
              <strong>{refundStats.total}</strong>
              <small>退款批次</small>
            </article>
            <article>
              <span>待处理</span>
              <strong>{refundStats.pending}</strong>
              <small>待商家审核</small>
            </article>
            <article>
              <span>处理中</span>
              <strong>{refundStats.processing}</strong>
              <small>退款执行中</small>
            </article>
            <article>
              <span>申请金额</span>
              <strong>{`CNY ${formatCnyFromCents(refundStats.amount)}`}</strong>
              <small>当前筛选范围</small>
            </article>
          </section>

          <section className="seller-after-sales-board">
            <div className="seller-after-sales-toolbar">
              <div className="seller-after-sales-filter-row">
                {FILTER_OPTIONS.map((item) => (
                  <button key={item.key} type="button" className={filterKey === item.key ? "is-active" : ""} onClick={() => setFilterKey(item.key)}>
                    {item.label}
                  </button>
                ))}
              </div>
              <p>{normalizedSelectedShopNo ? `店铺：${normalizedSelectedShopNo}` : "请选择店铺"}</p>
            </div>

            <div className="seller-after-sales-list-grid">
              {refundListQuery.isLoading ? (
                <Skeleton active paragraph={{ rows: 6 }} />
              ) : refundBatches.length === 0 ? (
                <div className="seller-after-sales-empty seller-after-sales-empty-inline">
                  <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="当前筛选下还没有退款申请" />
                </div>
              ) : (
                refundBatches.map((batch) => {
                  const statusMeta = getSellerRefundStatusMeta(batch.batchStatus);
                  return (
                    <article key={batch.refundBatchNo} className="seller-after-sales-card">
                      <div className="seller-after-sales-card-main">
                        <div className="seller-after-sales-card-head">
                          <div>
                            <strong>{batch.orderNo}</strong>
                            <p>{batch.refundBatchNo}</p>
                          </div>
                          <span className={`seller-after-sales-pill is-${statusMeta.tone}`}>{statusMeta.label}</span>
                        </div>
                        <div className="seller-after-sales-card-meta">
                          <span>{`申请金额：CNY ${formatCnyFromCents(batch.applyRefundAmount)}`}</span>
                          <span>{`子单数量：${batch.caseCount}`}</span>
                          <span>{`申请时间：${formatTime(batch.createdAt)}`}</span>
                          <span>{`审核截止：${formatTime(batch.reviewDeadlineAt)}`}</span>
                        </div>
                      </div>
                      <div className="seller-after-sales-card-side">
                        <small>{batch.subOrderNos.join(" / ") || "暂无子单号"}</small>
                        <Button
                          type="primary"
                          ghost
                          onClick={() => {
                            setActiveBatchNo(batch.refundBatchNo);
                            setDetailOpen(true);
                          }}
                        >
                          查看详情
                        </Button>
                      </div>
                    </article>
                  );
                })
              )}
            </div>
          </section>
        </>
      )}

      <Modal
        open={detailOpen}
        footer={null}
        width={980}
        destroyOnClose
        onCancel={() => {
          setDetailOpen(false);
          setSellerReply("");
        }}
        className="seller-after-sales-modal"
      >
        {detailQuery.isLoading ? (
          <Skeleton active paragraph={{ rows: 10 }} />
        ) : detailQuery.isError ? (
          <Alert
            type="error"
            showIcon
            message="退款详情加载失败"
            description={detailQuery.error instanceof Error ? detailQuery.error.message : "请稍后重试"}
          />
        ) : selectedBatch ? (
          <div className="seller-after-sales-modal-content">
            <header className="seller-after-sales-modal-head">
              <div className="seller-after-sales-modal-copy">
                <span>退款详情</span>
                <h2>{selectedBatch.orderNo}</h2>
                <p>{`退款批次 ${selectedBatch.refundBatchNo}`}</p>
              </div>
              <div className="seller-after-sales-detail-side">
                <span className={`seller-after-sales-pill is-${getSellerRefundStatusMeta(selectedBatch.batchStatus).tone}`}>
                  {getSellerRefundStatusMeta(selectedBatch.batchStatus).label}
                </span>
                <small>{`审核截止：${formatTime(selectedBatch.reviewDeadlineAt)}`}</small>
              </div>
            </header>

            <section className="seller-after-sales-metrics">
              <article>
                <span>申请金额</span>
                <strong>{`CNY ${formatCnyFromCents(selectedBatch.applyRefundAmount)}`}</strong>
              </article>
              <article>
                <span>已确认金额</span>
                <strong>{`CNY ${formatCnyFromCents(selectedBatch.approvedRefundAmount)}`}</strong>
              </article>
              <article>
                <span>店铺编号</span>
                <strong>{selectedBatch.shopNo || "--"}</strong>
              </article>
              <article>
                <span>关联子单</span>
                <strong>{selectedBatch.subOrderNos.join(" / ") || "--"}</strong>
              </article>
            </section>

            <section className="seller-after-sales-section">
              <div className="seller-after-sales-section-head">
                <h3>退款申请</h3>
                <span>{`${selectedDetail?.cases.length ?? 0} 条`}</span>
              </div>
              <div className="seller-after-sales-case-list">
                {selectedDetail?.cases.map((item) => (
                  <article key={item.afterSaleNo} className="seller-after-sales-case-card">
                    <div className="seller-after-sales-case-head">
                      <strong>{item.subOrderNo || item.afterSaleNo}</strong>
                      <span className={`seller-after-sales-pill is-${getSellerRefundStatusMeta(item.afterSaleStatus).tone}`}>
                        {getSellerRefundStatusMeta(item.afterSaleStatus).label}
                      </span>
                    </div>
                    <div className="seller-after-sales-case-grid">
                      <p>{`退款原因：${formatSellerRefundReasonLabel(item.reasonCode, item.reasonDesc)}`}</p>
                      <p>{`退款金额：CNY ${formatCnyFromCents(item.applyRefundAmount)}`}</p>
                      <p>{`商品项：${item.selectedItemNos.join(" / ") || item.itemNo || "--"}`}</p>
                      <p>{`支付单号：${item.paymentNo || "--"}`}</p>
                    </div>
                    <small>{`买家备注：${item.buyerRemark || "无"}`}</small>
                    {item.sellerReply ? <small>{`商家回复：${item.sellerReply}`}</small> : null}
                  </article>
                ))}
              </div>
            </section>

            {selectedDetail?.refundTasks.length ? (
              <section className="seller-after-sales-section">
                <div className="seller-after-sales-section-head">
                  <h3>退款进度</h3>
                  <span>{`${selectedDetail.refundTasks.length} 条`}</span>
                </div>
                <div className="seller-after-sales-task-list">
                  {selectedDetail.refundTasks.map((task) => (
                    <article key={task.refundTaskNo} className="seller-after-sales-task-card">
                      <div className="seller-after-sales-task-head">
                        <strong>{task.refundTaskNo}</strong>
                        <span>{getSellerRefundTaskStatusLabel(task.status)}</span>
                      </div>
                      <p>{`退款金额：CNY ${formatCnyFromCents(task.finalCashRefundAmount || task.refundAmount)}`}</p>
                      <p>{`积分返还：${task.pointsReturnAmount || 0}`}</p>
                      {task.lastErrorMessage ? <small>{`异常信息：${task.lastErrorMessage}`}</small> : <small>{`更新时间：${formatTime(task.updatedAt)}`}</small>}
                    </article>
                  ))}
                </div>
              </section>
            ) : null}

            <section className="seller-after-sales-review-box">
              <div className="seller-after-sales-section-head">
                <h3>审核处理</h3>
                <span>{getReviewHint(selectedBatchStatus)}</span>
              </div>
              <Input.TextArea
                value={sellerReply}
                rows={4}
                placeholder={canReview ? "填写处理说明后，可同意或驳回退款申请" : "当前状态无需审核处理"}
                onChange={(event) => setSellerReply(event.target.value)}
                disabled={!canReview || reviewMutation.isPending}
              />
              <div className="seller-after-sales-review-actions">
                <Button danger disabled={!canReview || reviewMutation.isPending} onClick={() => reviewMutation.mutate("reject")}>
                  {reviewMutation.isPending ? "处理中..." : "驳回退款"}
                </Button>
                <Button type="primary" disabled={!canReview || reviewMutation.isPending} loading={reviewMutation.isPending} onClick={() => reviewMutation.mutate("approve")}>
                  同意退款
                </Button>
              </div>
            </section>
          </div>
        ) : (
          <div className="seller-after-sales-empty seller-after-sales-empty-inline">
            <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="退款详情暂不可用" />
          </div>
        )}
      </Modal>
    </section>
  );
}
