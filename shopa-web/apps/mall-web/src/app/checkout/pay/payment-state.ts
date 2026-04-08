import type { BuyerOrder } from "../../../features/order/types";

export const BUYER_ORDERS_QUERY_KEY = ["buyer-orders", "all"] as const;
export const MY_CART_QUERY_KEY = ["my-cart"] as const;
export const SHELL_CART_QUERY_KEY = ["shell", "cart"] as const;
const PAYMENT_REDIRECT_SIGNAL_KEYS = ["trade_no", "out_trade_no", "trade_status", "app_id", "seller_id", "auth_app_id", "sign"] as const;

type SearchParamReader = {
  get(name: string): string | null;
};

export function isPaidOrder(order: BuyerOrder | null | undefined): boolean {
  if (!order) {
    return false;
  }
  return (
    order.paymentStatus === "PAID" ||
    order.orderStatus === "PAID" ||
    order.orderStatus === "FULFILLING" ||
    order.orderStatus === "COMPLETED"
  );
}

export function isPendingPayOrder(order: BuyerOrder | null | undefined): boolean {
  if (!order) {
    return false;
  }
  return order.orderStatus === "PENDING_PAY" && !isPaidOrder(order);
}

export function hasPaymentRedirectSignal(searchParams: SearchParamReader | null | undefined): boolean {
  if (!searchParams) {
    return false;
  }
  return PAYMENT_REDIRECT_SIGNAL_KEYS.some((key) => {
    const value = searchParams.get(key);
    return Boolean(value && value.trim());
  });
}

export function shouldShowPayCountdown(order: BuyerOrder | null | undefined, paymentExpired = false): boolean {
  if (!order || paymentExpired) {
    return false;
  }
  return isPendingPayOrder(order) && Boolean(order.payDeadlineAt);
}

export function paymentAttemptStorageKey(orderNo: string): string {
  return `checkout_pay_attempt_${orderNo}`;
}

export function getSettlementRefreshQueryKeys(order: BuyerOrder | null | undefined): ReadonlyArray<readonly string[]> {
  if (isPaidOrder(order)) {
    return [BUYER_ORDERS_QUERY_KEY, MY_CART_QUERY_KEY, SHELL_CART_QUERY_KEY];
  }
  return [BUYER_ORDERS_QUERY_KEY];
}
