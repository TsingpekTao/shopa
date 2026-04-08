"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { message } from "antd";
import { useI18n } from "@shopa/ui";
import { getBuyerProductDetail, listBuyerAssetReadUrls, listBuyerProductImages } from "@/features/catalog/api";
import { buildIdempotencyKey, cancelMyOrder, confirmMyOrderReceived, listMyOrders } from "@/features/order/api";
import { canConfirmReceipt, getAutoReceiveDeadlineMs, getAutoReceiveRemainingMs } from "@/features/order/order-display";
import { BuyerOrder, BuyerOrderItem } from "@/features/order/types";
import { buildAfterSaleConsultHref } from "@/features/order/order-chat";
import { collectProductDetailAssetIds, resolveProductImageUrl, type ProductDetailMap } from "@/lib/product-media";
import { formatCnyFromCents } from "@/lib/price";
import { canApplyRefundForSubOrder } from "@/features/refund/status";
import {
  getSettlementRefreshQueryKeys,
  hasPaymentRedirectSignal,
  isPaidOrder,
  isPendingPayOrder,
  paymentAttemptStorageKey
} from "../../checkout/pay/payment-state";

type OrderTabKey = "pending-pay" | "waiting-shipment" | "waiting-receipt" | "waiting-review" | "after-sale" | "history";

type TabOption = {
  key: OrderTabKey;
  zh: string;
  en: string;
};

const ORDER_TABS: TabOption[] = [
  { key: "pending-pay", zh: "\u5f85\u4ed8\u6b3e", en: "Pending Payment" },
  { key: "waiting-shipment", zh: "\u5f85\u53d1\u8d27", en: "Pending Shipment" },
  { key: "waiting-receipt", zh: "\u5f85\u6536\u8d27", en: "Pending Receipt" },
  { key: "waiting-review", zh: "\u5f85\u8bc4\u4ef7", en: "Pending Review" },
  { key: "after-sale", zh: "\u9000\u6b3e/\u552e\u540e", en: "After-sale" },
  { key: "history", zh: "\u5386\u53f2\u8ba2\u5355", en: "History" }
];

const ORDER_ITEM_PLACEHOLDER = "linear-gradient(135deg, #ffe2d1, #ffd1ae)";

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

function toTimeMs(raw: string): number {
  if (!raw) {
    return 0;
  }
  const timestamp = new Date(raw).getTime();
  return Number.isFinite(timestamp) ? timestamp : 0;
}

function formatDurationMs(remainingMs: number, isZh: boolean): string {
  const totalSeconds = Math.max(0, Math.ceil(remainingMs / 1000));
  const days = Math.floor(totalSeconds / (24 * 60 * 60));
  const hours = Math.floor((totalSeconds % (24 * 60 * 60)) / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;

  const pad = (value: number) => String(value).padStart(2, "0");
  if (isZh) {
    return `${days}天 ${pad(hours)}:${pad(minutes)}:${pad(seconds)}`;
  }
  return `${days}d ${pad(hours)}:${pad(minutes)}:${pad(seconds)}`;
}

function resolveOrderTab(order: BuyerOrder): OrderTabKey {
  const subStatuses = order.subOrders.map((sub) => sub.subStatus);
  if (order.orderStatus === "CANCELED" || order.orderStatus === "CLOSED") {
    return "history";
  }
  if (order.orderStatus === "REFUNDING" || subStatuses.some((status) => status === "REFUNDING")) {
    return "after-sale";
  }
  if (order.orderStatus === "REFUNDED" || subStatuses.some((status) => status === "REFUNDED")) {
    return "history";
  }
  if (order.orderStatus === "COMPLETED" || subStatuses.some((status) => status === "COMPLETED")) {
    return "waiting-review";
  }
  if (subStatuses.some((status) => status === "SHIPPED")) {
    return "waiting-receipt";
  }
  if (
    order.paymentStatus === "PAID" ||
    order.orderStatus === "PAID" ||
    order.orderStatus === "FULFILLING" ||
    subStatuses.some((status) => status === "PAID" || status === "WAIT_SHIP")
  ) {
    return "waiting-shipment";
  }
  return "pending-pay";
}

function orderStatusText(order: BuyerOrder, isZh: boolean, activeTab: OrderTabKey): string {
  if (activeTab === "history") {
    if (order.orderStatus === "CANCELED") {
      return isZh ? "已取消" : "Canceled";
    }
    if (order.orderStatus === "CLOSED") {
      return isZh ? "已关闭" : "Closed";
    }
    if (order.orderStatus === "REFUNDED" || order.subOrders.some((sub) => sub.subStatus === "REFUNDED")) {
      return isZh ? "已退款" : "Refunded";
    }
    if (order.orderStatus === "REFUNDING" || order.subOrders.some((sub) => sub.subStatus === "REFUNDING")) {
      return isZh ? "售后处理中" : "After-sale Processing";
    }
    if (order.orderStatus === "COMPLETED" || order.subOrders.some((sub) => sub.subStatus === "COMPLETED")) {
      return isZh ? "已完成" : "Completed";
    }
    if (order.subOrders.some((sub) => sub.subStatus === "SHIPPED")) {
      return isZh ? "已发货" : "Shipped";
    }
    if (
      order.paymentStatus === "PAID" ||
      order.orderStatus === "PAID" ||
      order.orderStatus === "FULFILLING" ||
      order.subOrders.some((sub) => sub.subStatus === "PAID" || sub.subStatus === "WAIT_SHIP")
    ) {
      return isZh ? "已支付" : "Paid";
    }
  }

  const tab = resolveOrderTab(order);
  switch (tab) {
    case "after-sale":
      return order.orderStatus === "REFUNDED" || order.subOrders.some((sub) => sub.subStatus === "REFUNDED")
        ? isZh
          ? "\u5df2\u9000\u6b3e"
          : "Refunded"
        : isZh
          ? "\u552e\u540e\u5904\u7406\u4e2d"
          : "After-sale Processing";
    case "waiting-review":
      return isZh ? "待评价" : "Pending Review";
    case "waiting-receipt":
      return isZh ? "待收货" : "Pending Receipt";
    case "waiting-shipment":
      return isZh ? "待发货" : "Pending Shipment";
    default:
      return isZh ? "待付款" : "Pending Payment";
  }
}

function sortOrdersForTab(orders: BuyerOrder[], tab: OrderTabKey): BuyerOrder[] {
  const sorted = [...orders];
  sorted.sort((left, right) => {
    const leftSortTime = tab === "history" ? toTimeMs(left.paidAt) || toTimeMs(left.createdAt) : toTimeMs(left.createdAt);
    const rightSortTime = tab === "history" ? toTimeMs(right.paidAt) || toTimeMs(right.createdAt) : toTimeMs(right.createdAt);
    if (rightSortTime !== leftSortTime) {
      return rightSortTime - leftSortTime;
    }
    return right.orderNo.localeCompare(left.orderNo);
  });
  return sorted;
}

function firstOrderItem(order: BuyerOrder): BuyerOrderItem | null {
  return order.subOrders[0]?.items[0] ?? null;
}

function orderSummary(order: BuyerOrder, isZh: boolean): string {
  const items = order.subOrders.flatMap((sub) => sub.items);
  const firstItem = items[0];
  if (!firstItem) {
    return order.orderNo;
  }
  const title = firstItem.spuTitle || firstItem.skuName || firstItem.skuNo;
  if (items.length <= 1) {
    return title;
  }
  return isZh ? `${title} \u7b49${items.length}\u4ef6\u5546\u54c1` : `${title} and ${items.length - 1} more item(s)`;
}

function shopLabel(order: BuyerOrder, isZh: boolean): string {
  const shopNo = firstOrderItem(order)?.shopNo;
  if (!shopNo) {
    return isZh ? "店铺待同步" : "Shop pending";
  }
  return isZh ? `店铺 ${shopNo}` : `Shop ${shopNo}`;
}

function canCancelPendingOrder(order: BuyerOrder): boolean {
  return order.orderStatus === "PENDING_PAY";
}

function getRefundEntry(order: BuyerOrder): { orderNo: string; subOrderNo: string; itemNo: string } | null {
  const targetSubOrder = order.subOrders.find((subOrder) => canApplyRefundForSubOrder(subOrder));
  if (!targetSubOrder) {
    return null;
  }
  const firstItem = targetSubOrder.items[0];
  if (!firstItem) {
    return null;
  }
  return {
    orderNo: order.orderNo,
    subOrderNo: targetSubOrder.subOrderNo,
    itemNo: firstItem.itemNo
  };
}

export default function OrdersPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [messageApi, contextHolder] = message.useMessage();
  const queryClient = useQueryClient();
  const searchParams = useSearchParams();
  const highlightedOrderNo = (searchParams.get("order_no") ?? "").trim();
  const explicitTab = (searchParams.get("tab") ?? "").trim() as OrderTabKey;
  const [activeTab, setActiveTab] = useState<OrderTabKey>("pending-pay");
  const [nowMs, setNowMs] = useState(() => Date.now());
  const [hasRecentPayAttempt, setHasRecentPayAttempt] = useState(false);
  const [isRefreshingPaymentResult, setIsRefreshingPaymentResult] = useState(false);
  const paymentSignalFromRedirect = hasPaymentRedirectSignal(searchParams);

  const ordersQuery = useQuery({
    queryKey: ["buyer-orders", "all"],
    queryFn: () => listMyOrders({ pageSize: 100 }),
    staleTime: 10_000,
    refetchOnWindowFocus: false
  });

  const allOrders = useMemo(() => ordersQuery.data?.orders ?? [], [ordersQuery.data?.orders]);
  const highlightedOrder = useMemo(() => allOrders.find((order) => order.orderNo === highlightedOrderNo) ?? null, [allOrders, highlightedOrderNo]);

  const assetIds = useMemo(() => {
    return Array.from(
      new Set(
        allOrders
          .flatMap((order) => order.subOrders)
          .flatMap((sub) => sub.items)
          .map((item) => item.skuImageAssetId?.trim() ?? "")
          .filter((item) => item.length > 0)
      )
    );
  }, [allOrders]);

  const spuNos = useMemo(() => {
    return Array.from(
      new Set(
        allOrders
          .flatMap((order) => order.subOrders)
          .flatMap((sub) => sub.items)
          .map((item) => item.spuNo?.trim() ?? "")
          .filter((item) => item.length > 0)
      )
    );
  }, [allOrders]);

  const assetUrlsQuery = useQuery({
    queryKey: ["buyer-order-asset-urls", assetIds.join(",")],
    queryFn: () => listBuyerAssetReadUrls(assetIds),
    enabled: assetIds.length > 0,
    staleTime: 60_000,
    refetchOnWindowFocus: false
  });

  const productImagesQuery = useQuery({
    queryKey: ["buyer-order-product-images", spuNos.join(",")],
    queryFn: () => listBuyerProductImages(spuNos),
    enabled: spuNos.length > 0,
    staleTime: 60_000,
    refetchOnWindowFocus: false
  });

  const productDetailQuery = useQuery({
    queryKey: ["buyer-order-product-details", spuNos.join(",")],
    queryFn: async () => {
      const results = await Promise.all(
        spuNos.map(async (spuNo) => {
          try {
            return await getBuyerProductDetail(spuNo);
          } catch {
            return null;
          }
        })
      );
      return results.reduce<ProductDetailMap>((acc, detail) => {
        if (!detail?.spuNo) {
          return acc;
        }
        acc[detail.spuNo] = detail;
        return acc;
      }, {});
    },
    enabled: spuNos.length > 0,
    staleTime: 60_000,
    refetchOnWindowFocus: false
  });

  const detailAssetIds = useMemo(() => collectProductDetailAssetIds(productDetailQuery.data ?? {}), [productDetailQuery.data]);

  const detailAssetUrlsQuery = useQuery({
    queryKey: ["buyer-order-detail-asset-urls", detailAssetIds.join(",")],
    queryFn: () => listBuyerAssetReadUrls(detailAssetIds),
    enabled: detailAssetIds.length > 0,
    staleTime: 60_000,
    refetchOnWindowFocus: false
  });

  const assetUrlMap = assetUrlsQuery.data ?? {};
  const productImageMap = productImagesQuery.data ?? {};
  const productDetailMap = productDetailQuery.data ?? {};
  const detailAssetUrlMap = detailAssetUrlsQuery.data ?? {};

  useEffect(() => {
    if (explicitTab && ORDER_TABS.some((tab) => tab.key === explicitTab)) {
      setActiveTab(explicitTab);
      return;
    }
    if (!highlightedOrderNo) {
      return;
    }
    const target = allOrders.find((order) => order.orderNo === highlightedOrderNo);
    if (target) {
      setActiveTab(resolveOrderTab(target));
    }
  }, [allOrders, explicitTab, highlightedOrderNo]);

  useEffect(() => {
    const timer = window.setInterval(() => {
      setNowMs(Date.now());
    }, 1000);
    return () => window.clearInterval(timer);
  }, []);

  useEffect(() => {
    if (!highlightedOrderNo || typeof window === "undefined") {
      setHasRecentPayAttempt(false);
      return;
    }
    const timestamp = Number(window.sessionStorage.getItem(paymentAttemptStorageKey(highlightedOrderNo)) ?? "0");
    if (!Number.isFinite(timestamp) || timestamp <= 0) {
      setHasRecentPayAttempt(false);
      return;
    }
    setHasRecentPayAttempt(Date.now() - timestamp <= 5 * 60_000);
  }, [highlightedOrderNo]);

  useEffect(() => {
    if (!highlightedOrderNo) {
      setIsRefreshingPaymentResult(false);
      return;
    }
    if (!paymentSignalFromRedirect && !hasRecentPayAttempt) {
      setIsRefreshingPaymentResult(false);
      return;
    }
    if (highlightedOrder && !isPendingPayOrder(highlightedOrder)) {
      setIsRefreshingPaymentResult(false);
      if (typeof window !== "undefined") {
        window.sessionStorage.removeItem(paymentAttemptStorageKey(highlightedOrderNo));
      }
      setHasRecentPayAttempt(false);
      if (isPaidOrder(highlightedOrder)) {
        void Promise.all(
          getSettlementRefreshQueryKeys(highlightedOrder).map((queryKey) => queryClient.invalidateQueries({ queryKey }))
        );
      }
      return;
    }

    let stopped = false;
    let attempts = 0;
    let timer: number | null = null;
    const maxAttempts = 10;
    setIsRefreshingPaymentResult(true);

    const refreshOrders = async () => {
      if (stopped) {
        return;
      }
      attempts += 1;
      const result = await ordersQuery.refetch();
      const nextOrder = result.data?.orders.find((order) => order.orderNo === highlightedOrderNo) ?? null;
      const settled = Boolean(nextOrder && !isPendingPayOrder(nextOrder));
      if (!settled && attempts < maxAttempts) {
        return;
      }
      stopped = true;
      if (timer !== null) {
        window.clearInterval(timer);
      }
      setIsRefreshingPaymentResult(false);
      if (settled && typeof window !== "undefined") {
        window.sessionStorage.removeItem(paymentAttemptStorageKey(highlightedOrderNo));
      }
      if (settled) {
        setHasRecentPayAttempt(false);
        if (nextOrder && isPaidOrder(nextOrder)) {
          await Promise.all(
            getSettlementRefreshQueryKeys(nextOrder).map((queryKey) => queryClient.invalidateQueries({ queryKey }))
          );
        }
      }
    };

    void refreshOrders();
    timer = window.setInterval(() => {
      void refreshOrders();
    }, 1500);

    return () => {
      stopped = true;
      if (timer !== null) {
        window.clearInterval(timer);
      }
      setIsRefreshingPaymentResult(false);
    };
  }, [hasRecentPayAttempt, highlightedOrder, highlightedOrderNo, ordersQuery, paymentSignalFromRedirect, queryClient]);

  const groupedOrders = useMemo(() => {
    const buckets = new Map<OrderTabKey, BuyerOrder[]>();
    ORDER_TABS.forEach((tab) => buckets.set(tab.key, []));

    allOrders.forEach((order) => {
      const primaryTab = resolveOrderTab(order);
      if (primaryTab !== "history") {
        buckets.get(primaryTab)?.push(order);
      }
    });

    ORDER_TABS.forEach((tab) => {
      const sourceOrders = tab.key === "history" ? allOrders : buckets.get(tab.key) ?? [];
      buckets.set(tab.key, sortOrdersForTab(sourceOrders, tab.key));
    });

    return buckets;
  }, [allOrders]);

  const currentOrders = groupedOrders.get(activeTab) ?? [];
  const totalAmount = currentOrders.reduce((sum, order) => sum + order.amount.payableAmount, 0);

  const cancelMutation = useMutation({
    mutationFn: async (orderNo: string) => {
      return cancelMyOrder({
        orderNo,
        reasonCode: 1,
        idempotencyKey: buildIdempotencyKey("cancel_order")
      });
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["buyer-orders", "all"] });
      messageApi.success(isZh ? "\u8ba2\u5355\u5df2\u53d6\u6d88" : "Order canceled");
    },
    onError: (error) => {
      const errorMessage = error instanceof Error ? error.message : isZh ? "\u53d6\u6d88\u8ba2\u5355\u5931\u8d25" : "Failed to cancel order";
      messageApi.error(errorMessage);
    }
  });

  const confirmReceiptMutation = useMutation({
    mutationFn: async (orderNo: string) => {
      return confirmMyOrderReceived({
        orderNo,
        idempotencyKey: buildIdempotencyKey("confirm_receipt")
      });
    },
    onSuccess: async (result) => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ["buyer-orders", "all"] }),
        queryClient.invalidateQueries({ queryKey: ["mall-me-overview"] })
      ]);
      messageApi.success(
        result?.pointsGrantTriggered
          ? isZh
            ? "已确认收货，积分已发放"
            : "Receipt confirmed and points granted"
          : isZh
            ? "已确认收货，积分稍后到账"
            : "Receipt confirmed, points will arrive shortly"
      );
    },
    onError: (error) => {
      const errorMessage = error instanceof Error ? error.message : isZh ? "确认收货失败" : "Failed to confirm receipt";
      messageApi.error(errorMessage);
    }
  });

  return (
    <section className="tb-order-center">
      {contextHolder}
      <header className="tb-order-hero">
        <div className="tb-order-hero-copy">
          <p className="tb-order-hero-kicker">{isZh ? "订单中心" : "Order Center"}</p>
          <h1>{isZh ? "我的订单" : "My Orders"}</h1>
          {highlightedOrderNo && (paymentSignalFromRedirect || hasRecentPayAttempt) ? (
            <p className="tb-order-hero-kicker" style={{ marginTop: 10 }}>
              {isRefreshingPaymentResult
                ? isZh
                  ? "正在同步支付结果..."
                  : "Syncing payment result..."
                : highlightedOrder && isPaidOrder(highlightedOrder)
                  ? isZh
                    ? "支付成功，订单已更新到最新状态。"
                    : "Payment completed and order status is up to date."
                  : isZh
                    ? "订单状态已更新。"
                    : "Order status updated."}
            </p>
          ) : null}
        </div>
        <div className="tb-order-hero-metrics">
          <div>
            <strong>{allOrders.length}</strong>
            <span>{isZh ? "全部订单" : "All Orders"}</span>
          </div>
          <div>
            <strong>{currentOrders.length}</strong>
            <span>{isZh ? "当前标签" : "Current Tab"}</span>
          </div>
          <div>
            <strong>{formatCnyFromCents(totalAmount)}</strong>
            <span>{isZh ? "当前金额" : "Current Amount"}</span>
          </div>
        </div>
      </header>

      <div className="tb-order-tabs">
        {ORDER_TABS.map((tab) => (
          <button
            key={tab.key}
            type="button"
            className={`tb-order-tab ${activeTab === tab.key ? "is-active" : ""}`}
            onClick={() => setActiveTab(tab.key)}
          >
            <span>{isZh ? tab.zh : tab.en}</span>
            <strong>{groupedOrders.get(tab.key)?.length ?? 0}</strong>
          </button>
        ))}
      </div>

      {ordersQuery.isLoading ? <p className="tb-chat-empty">{isZh ? "正在加载订单..." : "Loading orders..."}</p> : null}
      {ordersQuery.isError ? <p className="tb-chat-empty">{isZh ? "订单加载失败，请稍后重试。" : "Failed to load orders."}</p> : null}

      {!ordersQuery.isLoading && currentOrders.length === 0 ? (
        <section className="tb-order-empty-card">
          <strong>{isZh ? "当前状态暂无订单" : "No orders in this tab"}</strong>
          <p>{isZh ? "下单后的订单会显示在这里。" : "Orders will appear here after checkout."}</p>
        </section>
      ) : null}

      <div className="tb-order-list">
        {currentOrders.map((order) => {
          const highlighted = highlightedOrderNo === order.orderNo;
          const items = order.subOrders.flatMap((sub) => sub.items);
          const consultHref = buildAfterSaleConsultHref(order);
          const refundEntry = getRefundEntry(order);
          const autoReceiveDeadlineMs = getAutoReceiveDeadlineMs(order);
          const autoReceiveRemainingMs = getAutoReceiveRemainingMs(order, nowMs);
          const showAutoReceiveCountdown = resolveOrderTab(order) === "waiting-receipt" && autoReceiveDeadlineMs > 0;

          return (
            <article key={order.orderNo} className={`tb-order-card ${highlighted ? "is-highlighted" : ""}`}>
              <div className="tb-order-card-head">
                <div className="tb-order-card-shop">
                  <strong>{shopLabel(order, isZh)}</strong>
                  <span>{formatTime(order.createdAt, locale)}</span>
                </div>
                <div className="tb-order-card-status">
                  <span>{orderStatusText(order, isZh, activeTab)}</span>
                  <small>{order.orderNo}</small>
                </div>
              </div>

              <div className="tb-order-card-body">
                <div className="tb-order-card-items">
                  {items.map((item) => {
                    const itemImage = resolveProductImageUrl(item, assetUrlMap, productImageMap, productDetailMap, detailAssetUrlMap);
                    const showImg = itemImage && itemImage !== ORDER_ITEM_PLACEHOLDER;
                    return (
                      <div key={item.itemNo || `${item.subOrderNo}_${item.skuNo}`} className="tb-order-item-row">
                        <div
                          className="tb-order-item-thumb"
                          style={showImg ? undefined : { backgroundImage: ORDER_ITEM_PLACEHOLDER }}
                        >
                          {showImg ? <img src={itemImage} alt={item.spuTitle || item.skuName || item.skuNo} loading="lazy" decoding="async" referrerPolicy="no-referrer" /> : null}
                        </div>
                        <div className="tb-order-item-copy">
                          <strong>{item.spuTitle || item.skuName || item.skuNo}</strong>
                          <span>{item.skuName || item.skuNo}</span>
                          <small>{`x${item.qty}`}</small>
                        </div>
                        <div className="tb-order-item-price">{`CNY ${formatCnyFromCents(item.salePrice)}`}</div>
                      </div>
                    );
                  })}
                </div>

                <aside className="tb-order-card-side">
                  <div>
                    <p>{isZh ? "订单摘要" : "Order Summary"}</p>
                    <strong>{orderSummary(order, isZh)}</strong>
                  </div>
                  <div>
                    <p>{isZh ? "应付金额" : "Payable"}</p>
                    <strong className="is-price">{`CNY ${formatCnyFromCents(order.amount.payableAmount)}`}</strong>
                  </div>
                  {showAutoReceiveCountdown ? (
                    <div>
                      <p>{isZh ? "自动收货" : "Auto Receipt"}</p>
                      <strong>{autoReceiveRemainingMs > 0 ? formatDurationMs(autoReceiveRemainingMs, isZh) : isZh ? "系统处理中" : "Processing"}</strong>
                      <small>
                        {isZh
                          ? `截止 ${formatTime(new Date(autoReceiveDeadlineMs).toISOString(), locale)}`
                          : `Deadline ${formatTime(new Date(autoReceiveDeadlineMs).toISOString(), locale)}`}
                      </small>
                    </div>
                  ) : null}
                  <div className="tb-order-card-actions">
                    {resolveOrderTab(order) === "pending-pay" ? (
                      <Link href={`/checkout/pay?order_no=${encodeURIComponent(order.orderNo)}`} className="tb-order-btn tb-order-btn-primary">
                        {isZh ? "去支付" : "Pay Now"}
                      </Link>
                    ) : null}
                    {canCancelPendingOrder(order) ? (
                      <button
                        type="button"
                        className="tb-order-btn tb-order-btn-secondary"
                        disabled={cancelMutation.isPending && cancelMutation.variables === order.orderNo}
                        onClick={() => cancelMutation.mutate(order.orderNo)}
                      >
                        {cancelMutation.isPending && cancelMutation.variables === order.orderNo
                          ? isZh
                            ? "\u6b63\u5728\u53d6\u6d88..."
                            : "Canceling..."
                          : isZh
                            ? "\u53d6\u6d88\u8ba2\u5355"
                            : "Cancel Order"}
                      </button>
                    ) : null}
                    {canConfirmReceipt(order) ? (
                      <button
                        type="button"
                        className="tb-order-btn tb-order-btn-primary"
                        disabled={confirmReceiptMutation.isPending && confirmReceiptMutation.variables === order.orderNo}
                        onClick={() => confirmReceiptMutation.mutate(order.orderNo)}
                      >
                        {confirmReceiptMutation.isPending && confirmReceiptMutation.variables === order.orderNo
                          ? isZh
                            ? "正在确认..."
                            : "Confirming..."
                          : isZh
                            ? "确认收货"
                            : "Confirm Receipt"}
                      </button>
                    ) : null}
                    {refundEntry ? (
                      <Link
                        href={`/me/refunds/apply?orderNo=${encodeURIComponent(refundEntry.orderNo)}&subOrderNo=${encodeURIComponent(refundEntry.subOrderNo)}${refundEntry.itemNo ? `&itemNo=${encodeURIComponent(refundEntry.itemNo)}` : ""}`}
                        className="tb-order-btn tb-order-btn-secondary"
                      >
                        {isZh ? "申请退款" : "Request Refund"}
                      </Link>
                    ) : null}
                    {consultHref ? (
                      <Link href={consultHref} className="tb-order-btn tb-order-btn-secondary">
                        {isZh ? "售后咨询" : "After-sale Chat"}
                      </Link>
                    ) : null}
                  </div>
                </aside>
              </div>
            </article>
          );
        })}
      </div>
    </section>
  );
}
