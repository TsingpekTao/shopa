"use client";

import Link from "next/link";
import { useEffect, useMemo, useRef, useState } from "react";
import { useSearchParams } from "next/navigation";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { message } from "antd";
import { useI18n } from "@shopa/ui";
import { listMyAddresses } from "@/features/address/api";
import { UserAddress } from "@/features/address/types";
import { buildIdempotencyKey, cancelMyOrder, getMyOrderDetail, requestOrderPay, updateMyOrderAddress } from "@/features/order/api";
import { buildOrderAmountBreakdown } from "@/features/order/checkout-points";
import { BuyerOrder } from "@/features/order/types";
import { formatCnyFromCents } from "@/lib/price";
import {
  BUYER_ORDERS_QUERY_KEY,
  getSettlementRefreshQueryKeys,
  hasPaymentRedirectSignal,
  isPaidOrder,
  paymentAttemptStorageKey,
  shouldShowPayCountdown
} from "./payment-state";

function toPositiveInt(value: string | null): number {
  const parsed = Number(value ?? "");
  if (!Number.isFinite(parsed) || parsed <= 0) {
    return 0;
  }
  return Math.trunc(parsed);
}

function buildAddressText(address?: {
  provinceName?: string;
  cityName?: string;
  districtName?: string;
  street?: string;
  detail?: string;
} | null): string {
  if (!address) {
    return "";
  }
  return [address.provinceName, address.cityName, address.districtName, address.street, address.detail]
    .filter((item) => item && item.trim().length > 0)
    .join(" ");
}

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

function formatPayCountdown(remainingMs: number, isZh: boolean): string {
  const totalSeconds = Math.max(0, Math.ceil(remainingMs / 1000));
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const seconds = totalSeconds % 60;

  if (hours > 0) {
    return isZh ? `${hours}\u5c0f\u65f6 ${minutes}\u5206 ${seconds}\u79d2` : `${hours}h ${minutes}m ${seconds}s`;
  }
  return isZh ? `${minutes}\u5206 ${seconds}\u79d2` : `${minutes}m ${seconds}s`;
}

function orderStatusText(order: BuyerOrder | null, isZh: boolean): string {
  if (!order) {
    return "--";
  }
  if (order.orderStatus === "PENDING_PAY") {
    return isZh ? "待付款" : "Pending Payment";
  }
  if (order.orderStatus === "PAID" || order.orderStatus === "FULFILLING") {
    return isZh ? "待发货" : "Pending Shipment";
  }
  if (order.orderStatus === "COMPLETED") {
    return isZh ? "已完成" : "Completed";
  }
  if (order.orderStatus === "REFUNDING") {
    return isZh ? "退款中" : "Refunding";
  }
  if (order.orderStatus === "REFUNDED") {
    return isZh ? "已退款" : "Refunded";
  }
  if (order.orderStatus === "CANCELED") {
    return isZh ? "已取消" : "Canceled";
  }
  if (order.orderStatus === "CLOSED") {
    return isZh ? "已关闭" : "Closed";
  }
  return isZh ? "处理中" : "Processing";
}

function paymentStatusText(order: BuyerOrder | null, isZh: boolean, paymentExpired: boolean): string {
  if (!order) {
    return "--";
  }
  if (paymentExpired && order.orderStatus === "PENDING_PAY") {
    return isZh ? "支付超时" : "Payment Expired";
  }
  switch (order.paymentStatus) {
    case "UNPAID":
      return isZh ? "未支付" : "Unpaid";
    case "PAYING":
      return isZh ? "支付中" : "Paying";
    case "PAID":
      return isZh ? "已支付" : "Paid";
    case "PAY_FAILED":
      return isZh ? "支付失败" : "Payment Failed";
    case "REFUNDED":
      return isZh ? "已退款" : "Refunded";
    default:
      return isZh ? "待确认" : "Pending";
  }
}

function canRequestPay(order: BuyerOrder | null, paymentExpired = false): boolean {
  if (!order || paymentExpired) {
    return false;
  }
  return order.orderStatus === "PENDING_PAY" && (order.paymentStatus === "UNPAID" || order.paymentStatus === "PAY_FAILED" || order.paymentStatus === "PAYING");
}

function canUpdateAddress(order: BuyerOrder | null, paymentExpired = false): boolean {
  if (!order || paymentExpired) {
    return false;
  }
  return order.orderStatus === "PENDING_PAY" && (order.paymentStatus === "UNPAID" || order.paymentStatus === "PAY_FAILED");
}

function canCancelOrder(order: BuyerOrder | null, paymentExpired = false): boolean {
  if (!order || paymentExpired) {
    return false;
  }
  return order.orderStatus === "PENDING_PAY";
}

function findAddress(addresses: UserAddress[], addressId: number): UserAddress | null {
  if (!addressId) {
    return null;
  }
  return addresses.find((item) => item.addressId === addressId) ?? null;
}

async function invalidateSettlementQueries(
  queryClient: ReturnType<typeof useQueryClient>,
  order: BuyerOrder | null | undefined
): Promise<void> {
  await Promise.all(getSettlementRefreshQueryKeys(order).map((queryKey) => queryClient.invalidateQueries({ queryKey })));
}

export default function CheckoutPayPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [messageApi, contextHolder] = message.useMessage();
  const queryClient = useQueryClient();
  const searchParams = useSearchParams();
  const orderNo = (searchParams.get("order_no") ?? "").trim();
  const fallbackAddressId = toPositiveInt(searchParams.get("address_id"));
  const [latestPayLink, setLatestPayLink] = useState("");
  const [selectedAddressId, setSelectedAddressId] = useState(0);
  const [nowMs, setNowMs] = useState(() => Date.now());
  const [hasRecentPayAttempt, setHasRecentPayAttempt] = useState(false);
  const [isRefreshingPaymentResult, setIsRefreshingPaymentResult] = useState(false);
  const hasRefetchedAfterExpire = useRef(false);
  const paymentSignalFromRedirect = hasPaymentRedirectSignal(searchParams);

  const orderQuery = useQuery({
    queryKey: ["checkout-pay-order", orderNo],
    queryFn: () => getMyOrderDetail(orderNo),
    enabled: Boolean(orderNo),
    staleTime: 15_000,
    refetchOnWindowFocus: false
  });

  const addressQuery = useQuery({
    queryKey: ["checkout-pay-addresses"],
    queryFn: () => listMyAddresses(false),
    staleTime: 20_000,
    refetchOnWindowFocus: false
  });

  const order = orderQuery.data ?? null;
  const addresses = useMemo(() => addressQuery.data?.addresses ?? [], [addressQuery.data?.addresses]);
  const orderItems = useMemo(() => order?.subOrders.flatMap((sub) => sub.items) ?? [], [order]);
  const currentAddressId = order?.address?.sourceAddressId ?? fallbackAddressId;

  useEffect(() => {
    if (currentAddressId > 0) {
      setSelectedAddressId(currentAddressId);
      return;
    }
    const defaultAddress = addresses.find((item) => item.isDefault) ?? addresses[0];
    if (defaultAddress?.addressId) {
      setSelectedAddressId(defaultAddress.addressId);
    }
  }, [addresses, currentAddressId]);

  const selectedAddress = useMemo(() => findAddress(addresses, selectedAddressId), [addresses, selectedAddressId]);
  const currentAddress = useMemo(() => findAddress(addresses, currentAddressId), [addresses, currentAddressId]);
  const displayAddress = selectedAddress ?? currentAddress;
  const amountBreakdown = useMemo(
    () =>
      buildOrderAmountBreakdown({
        goodsAmount: order?.amount.goodsAmount ?? 0,
        freightAmount: order?.amount.freightAmount ?? 0,
        discountAmount: order?.amount.discountAmount ?? 0,
        pointsDiscountAmount: order?.amount.pointsDiscountAmount ?? order?.pointsDiscountAmount ?? 0,
        payableAmount: order?.amount.payableAmount ?? 0,
        paidAmount: order?.amount.paidAmount ?? 0
      }),
    [order]
  );

  const payDeadlineMs = useMemo(() => toTimeMs(order?.payDeadlineAt ?? ""), [order?.payDeadlineAt]);
  const remainingPayMs = payDeadlineMs > 0 ? Math.max(0, payDeadlineMs - nowMs) : 0;
  const paymentExpired = Boolean(order && order.orderStatus === "PENDING_PAY" && payDeadlineMs > 0 && remainingPayMs <= 0);
  const showPayCountdown = shouldShowPayCountdown(order, paymentExpired) && !isRefreshingPaymentResult;
  const canRequestPayNow = canRequestPay(order, paymentExpired);
  const canUpdateAddressNow = canUpdateAddress(order, paymentExpired);
  const canCancelOrderNow = canCancelOrder(order, paymentExpired);

  useEffect(() => {
    if (!orderNo || typeof window === "undefined") {
      setHasRecentPayAttempt(false);
      return;
    }
    const timestamp = Number(window.sessionStorage.getItem(paymentAttemptStorageKey(orderNo)) ?? "0");
    if (!Number.isFinite(timestamp) || timestamp <= 0) {
      setHasRecentPayAttempt(false);
      return;
    }
    setHasRecentPayAttempt(Date.now() - timestamp <= 5 * 60_000);
  }, [orderNo]);

  useEffect(() => {
    setNowMs(Date.now());
    if (!order || order.orderStatus !== "PENDING_PAY" || payDeadlineMs <= 0) {
      return;
    }
    const timer = window.setInterval(() => {
      setNowMs(Date.now());
    }, 1000);
    return () => window.clearInterval(timer);
  }, [order?.orderNo, order?.orderStatus, payDeadlineMs, order]);

  useEffect(() => {
    if (!paymentExpired) {
      hasRefetchedAfterExpire.current = false;
      return;
    }
    if (hasRefetchedAfterExpire.current) {
      return;
    }
    hasRefetchedAfterExpire.current = true;
    void orderQuery.refetch();
    void queryClient.invalidateQueries({ queryKey: BUYER_ORDERS_QUERY_KEY });
  }, [paymentExpired, orderQuery, queryClient]);

  useEffect(() => {
    if (!order || !orderNo || paymentExpired) {
      setIsRefreshingPaymentResult(false);
      return;
    }
    if (order.orderStatus !== "PENDING_PAY") {
      setIsRefreshingPaymentResult(false);
      if (typeof window !== "undefined") {
        window.sessionStorage.removeItem(paymentAttemptStorageKey(orderNo));
      }
      if (isPaidOrder(order)) {
        void invalidateSettlementQueries(queryClient, order);
      }
      return;
    }
    if (!paymentSignalFromRedirect && !hasRecentPayAttempt) {
      setIsRefreshingPaymentResult(false);
      return;
    }
    let stopped = false;
    let attempts = 0;
    let timer: number | null = null;
    const maxAttempts = 10;
    setIsRefreshingPaymentResult(true);

    const checkLatestOrder = async () => {
      if (stopped) {
        return;
      }
      attempts += 1;
      const result = await orderQuery.refetch();
      const nextOrder = result.data;
      const paid = isPaidOrder(nextOrder);
      const noLongerPending = Boolean(nextOrder && nextOrder.orderStatus !== "PENDING_PAY");
      if (paid || noLongerPending || attempts >= maxAttempts) {
        stopped = true;
        if (timer !== null) {
          window.clearInterval(timer);
        }
        setIsRefreshingPaymentResult(false);
        if (typeof window !== "undefined" && (paid || noLongerPending)) {
          window.sessionStorage.removeItem(paymentAttemptStorageKey(orderNo));
          setHasRecentPayAttempt(false);
        }
        if (paid || noLongerPending) {
          await invalidateSettlementQueries(queryClient, nextOrder);
        }
      }
    };

    void checkLatestOrder();
    timer = window.setInterval(() => {
      void checkLatestOrder();
    }, 1500);
    return () => {
      stopped = true;
      if (timer !== null) {
        window.clearInterval(timer);
      }
      setIsRefreshingPaymentResult(false);
    };
  }, [hasRecentPayAttempt, order, orderNo, orderQuery, paymentExpired, paymentSignalFromRedirect, queryClient]);

  const deadlineLabel = useMemo(() => {
    if (!order) {
      return "";
    }
    if (order.orderStatus === "CANCELED") {
      return isZh ? "订单已取消" : "Order Canceled";
    }
    if (order.orderStatus === "CLOSED") {
      return isZh ? "订单已关闭" : "Order Closed";
    }
    if (paymentExpired) {
      return isZh ? "支付时间已结束" : "Payment window expired";
    }
    if (!showPayCountdown || payDeadlineMs <= 0) {
      return "";
    }
    return formatPayCountdown(remainingPayMs, isZh);
  }, [isZh, order, payDeadlineMs, paymentExpired, remainingPayMs, showPayCountdown]);

  const updateAddressMutation = useMutation({
    mutationFn: async (addressId: number) => {
      if (!orderNo) {
        throw new Error("order_no is required");
      }
      return updateMyOrderAddress(orderNo, addressId);
    },
    onSuccess: (nextOrder) => {
      queryClient.setQueryData(["checkout-pay-order", orderNo], nextOrder);
      void queryClient.invalidateQueries({ queryKey: BUYER_ORDERS_QUERY_KEY });
      messageApi.success(isZh ? "收货地址已更新" : "Shipping address updated");
    },
    onError: (error) => {
      const errorMessage = error instanceof Error ? error.message : isZh ? "地址更新失败" : "Failed to update address";
      messageApi.error(errorMessage);
      setSelectedAddressId(currentAddressId);
    }
  });

  const payMutation = useMutation({
    mutationFn: async () => {
      if (!orderNo) {
        throw new Error("order_no is required");
      }
      return requestOrderPay({
        orderNo,
        payChannel: "ALIPAY",
        idempotencyKey: buildIdempotencyKey("alipay")
      });
    },
    onSuccess: (result) => {
      if (typeof window !== "undefined" && orderNo) {
        window.sessionStorage.setItem(paymentAttemptStorageKey(orderNo), String(Date.now()));
        setHasRecentPayAttempt(true);
      }
      setLatestPayLink(result.payUrl);
      if (result.payUrl) {
        const popup = window.open(result.payUrl, "_blank", "noopener,noreferrer");
        if (!popup) {
          window.location.href = result.payUrl;
          return;
        }
      }
      messageApi.success(isZh ? "已打开支付宝支付页面" : "Alipay page opened in a new tab");
    },
    onError: (error) => {
      const errorMessage = error instanceof Error ? error.message : isZh ? "发起支付失败" : "Failed to request payment";
      messageApi.error(errorMessage);
    }
  });

  const cancelMutation = useMutation({
    mutationFn: async () => {
      if (!orderNo) {
        throw new Error("order_no is required");
      }
      return cancelMyOrder({
        orderNo,
        reasonCode: 1,
        idempotencyKey: buildIdempotencyKey("cancel_order")
      });
    },
    onSuccess: async () => {
      if (typeof window !== "undefined" && orderNo) {
        window.sessionStorage.removeItem(paymentAttemptStorageKey(orderNo));
      }
      setHasRecentPayAttempt(false);
      setLatestPayLink("");
      await orderQuery.refetch();
      await queryClient.invalidateQueries({ queryKey: BUYER_ORDERS_QUERY_KEY });
      messageApi.success(isZh ? "\u8ba2\u5355\u5df2\u53d6\u6d88" : "Order canceled");
    },
    onError: (error) => {
      const errorMessage = error instanceof Error ? error.message : isZh ? "\u53d6\u6d88\u8ba2\u5355\u5931\u8d25" : "Failed to cancel order";
      messageApi.error(errorMessage);
    }
  });

  if (!orderNo) {
    return (
      <main className="tb-checkout-page">
        {contextHolder}
        <h1>{isZh ? "订单支付" : "Order Payment"}</h1>
        <p>{isZh ? "缺少订单号，请回到订单页重新进入。" : "Missing order number. Please reopen the order from My Orders."}</p>
        <Link href="/me/orders">{isZh ? "返回我的订单" : "Back to My Orders"}</Link>
      </main>
    );
  }

  return (
    <main className="tb-checkout-page">
      {contextHolder}
      <h1>{isZh ? "订单支付" : "Order Payment"}</h1>

      {orderQuery.isLoading ? <p>{isZh ? "正在加载订单..." : "Loading order..."}</p> : null}
      {orderQuery.isError ? <p>{isZh ? "订单加载失败，请稍后重试。" : "Failed to load order detail."}</p> : null}

      {order ? (
        <>
          <section className="tb-checkout-address">
            <h2>{isZh ? "订单信息" : "Order Info"}</h2>
            <div className={`tb-checkout-address-card ${paymentExpired || order.orderStatus === "CLOSED" || order.orderStatus === "CANCELED" ? "is-expired" : ""}`}>
              <p>
                {isZh ? "订单号：" : "Order No: "}
                {order.orderNo}
              </p>
              <small>
                {isZh ? "下单时间：" : "Created At: "}
                {formatTime(order.createdAt, locale)}
              </small>
              <small>
                {isZh ? "订单状态：" : "Order Status: "}
                {orderStatusText(order, isZh)}
              </small>
              <small>
                {isZh ? "支付状态：" : "Payment Status: "}
                {paymentStatusText(order, isZh, paymentExpired)}
              </small>
              {isRefreshingPaymentResult ? (
                <small>{isZh ? "正在确认支付结果，请稍候..." : "Confirming payment result..."}</small>
              ) : null}
              {!isRefreshingPaymentResult && isPaidOrder(order) ? (
                <small>{isZh ? "支付已完成，订单已更新。" : "Payment completed. Order updated."}</small>
              ) : null}
              {deadlineLabel ? (
                <div className={`tb-checkout-deadline ${paymentExpired || order.orderStatus === "CLOSED" || order.orderStatus === "CANCELED" ? "is-expired" : ""}`}>
                  <span>{paymentExpired ? (isZh ? "支付状态" : "Payment Status") : isZh ? "剩余支付时间" : "Time Left"}</span>
                  <strong>{deadlineLabel}</strong>
                  {showPayCountdown && order.payDeadlineAt ? <small>{`${isZh ? "截止时间：" : "Deadline: "}${formatTime(order.payDeadlineAt, locale)}`}</small> : null}
                </div>
              ) : null}
            </div>
          </section>

          <section className="tb-checkout-address">
            <h2>{isZh ? "收货地址" : "Shipping Address"}</h2>
            <div className="tb-checkout-address-card">
              <p>
                {(displayAddress?.receiverName || order.address?.receiverName || "--")} {(displayAddress?.receiverPhone || order.address?.receiverPhone || "--")}
              </p>
              <small>{buildAddressText(displayAddress) || buildAddressText(order.address) || "--"}</small>
            </div>

            {addresses.length > 0 ? (
              <div className="tb-checkout-address-list">
                {addresses.map((item) => {
                  const active = item.addressId === selectedAddressId;
                  const current = item.addressId === currentAddressId;
                  return (
                    <button
                      key={item.addressId}
                      type="button"
                      className={`tb-checkout-address-item ${active ? "active" : ""}`}
                      onClick={() => setSelectedAddressId(item.addressId)}
                    >
                      <strong>
                        {item.receiverName} {item.receiverPhone}
                        {item.isDefault ? <span>{isZh ? "默认" : "Default"}</span> : null}
                        {current ? <span>{isZh ? "当前" : "Current"}</span> : null}
                      </strong>
                      <small>{buildAddressText(item)}</small>
                    </button>
                  );
                })}
              </div>
            ) : (
              <p className="tb-checkout-address-empty">
                <Link href="/me/address">{isZh ? "先去新增地址" : "Add an address first"}</Link>
              </p>
            )}

            {canUpdateAddressNow && addresses.length > 0 ? (
              <p className="tb-checkout-address-empty">
                <button
                  type="button"
                  disabled={selectedAddressId <= 0 || selectedAddressId === currentAddressId || updateAddressMutation.isPending}
                  onClick={() => updateAddressMutation.mutate(selectedAddressId)}
                >
                  {updateAddressMutation.isPending ? (isZh ? "正在更新地址..." : "Updating address...") : isZh ? "更新到当前订单" : "Update Order Address"}
                </button>
                <Link href="/me/address">{isZh ? "管理地址" : "Manage Addresses"}</Link>
              </p>
            ) : null}
          </section>

          <section className="tb-checkout-items">
            <h2>{isZh ? "订单商品" : "Order Items"}</h2>
            <table>
              <thead>
                <tr>
                  <th>{isZh ? "商品" : "Item"}</th>
                  <th>{isZh ? "单价" : "Price"}</th>
                  <th>{isZh ? "数量" : "Qty"}</th>
                  <th>{isZh ? "小计" : "Subtotal"}</th>
                </tr>
              </thead>
              <tbody>
                {orderItems.map((item) => (
                  <tr key={item.itemNo || `${item.subOrderNo}_${item.skuNo}`}>
                    <td>{item.spuTitle || item.skuName || item.skuNo}</td>
                    <td>{`CNY ${formatCnyFromCents(item.salePrice)}`}</td>
                    <td>{item.qty}</td>
                    <td>{`CNY ${formatCnyFromCents(item.salePrice * item.qty)}`}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </section>

          <section className="tb-checkout-submit">
            <div className="tb-checkout-submit-form">
              <div className="tb-checkout-price-stack">
                <p className="tb-checkout-submit-note">
                  {isZh ? "商品金额：" : "Items: "}
                  <span>{`CNY ${formatCnyFromCents(amountBreakdown.goodsAmount)}`}</span>
                </p>
                {amountBreakdown.freightAmount > 0 ? (
                  <p className="tb-checkout-submit-note">
                    {isZh ? "运费：" : "Shipping: "}
                    <span>{`CNY ${formatCnyFromCents(amountBreakdown.freightAmount)}`}</span>
                  </p>
                ) : null}
                {amountBreakdown.discountAmount > 0 ? (
                  <p className="tb-checkout-submit-note">
                    {isZh ? "活动优惠：" : "Promotion: "}
                    <span>{`- CNY ${formatCnyFromCents(amountBreakdown.discountAmount)}`}</span>
                  </p>
                ) : null}
                {amountBreakdown.pointsDiscountAmount > 0 ? (
                  <p className="tb-checkout-submit-note">
                    {isZh ? "积分抵扣：" : "Points Discount: "}
                    <span>{`- CNY ${formatCnyFromCents(amountBreakdown.pointsDiscountAmount)}`}</span>
                  </p>
                ) : null}
              </div>
            </div>
            <div className="tb-checkout-payable">
              {amountBreakdown.totalDiscountAmount > 0 ? (
                <small>{`CNY ${formatCnyFromCents(amountBreakdown.originalPayableAmount)}`}</small>
              ) : null}
              <p>
                {isZh ? "应付金额：" : "Payable: "}
                <strong>{`CNY ${formatCnyFromCents(amountBreakdown.payableAmount)}`}</strong>
              </p>
            </div>
            {deadlineLabel ? (
              <p className={`tb-checkout-submit-note ${paymentExpired || order.orderStatus === "CLOSED" || order.orderStatus === "CANCELED" ? "is-expired" : ""}`}>
                {isZh ? "支付时效：" : "Payment Window: "}
                {deadlineLabel}
              </p>
            ) : null}
            <button type="button" disabled={!canRequestPayNow || payMutation.isPending} onClick={() => payMutation.mutate()}>
              {payMutation.isPending ? (isZh ? "正在拉起支付..." : "Opening payment...") : isZh ? "去支付宝支付" : "Pay with Alipay"}
            </button>
            {canCancelOrderNow ? (
              <button
                type="button"
                style={{ marginTop: 12 }}
                disabled={cancelMutation.isPending}
                onClick={() => cancelMutation.mutate()}
              >
                {cancelMutation.isPending ? (isZh ? "\u6b63\u5728\u53d6\u6d88..." : "Canceling...") : isZh ? "\u53d6\u6d88\u8ba2\u5355" : "Cancel Order"}
              </button>
            ) : null}
            {latestPayLink && canRequestPayNow ? (
              <p style={{ marginTop: 12 }}>
                <a href={latestPayLink} target="_blank" rel="noreferrer">
                  {isZh ? "如果新窗口被拦截，点这里继续支付" : "Continue payment here if the popup was blocked"}
                </a>
              </p>
            ) : null}
            {paymentExpired ? (
              <p className="tb-checkout-submit-note is-expired">
                {isZh ? "支付时间已结束，订单关闭后可在历史订单中查看。" : "Payment time has ended. You can find this order in History after it closes."}
              </p>
            ) : null}
            <p style={{ marginTop: 12 }}>
              <Link href={`/me/orders?order_no=${encodeURIComponent(order.orderNo)}`}>{isZh ? "返回我的订单" : "Back to My Orders"}</Link>
            </p>
          </section>
        </>
      ) : null}
    </main>
  );
}
