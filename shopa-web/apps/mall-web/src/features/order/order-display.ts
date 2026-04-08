import type { BuyerOrder } from "./types";

const AUTO_RECEIVE_WINDOW_MS = 7 * 24 * 60 * 60 * 1000;

function toTimeMs(raw: string): number {
  if (!raw) {
    return 0;
  }
  const timestamp = new Date(raw).getTime();
  return Number.isFinite(timestamp) ? timestamp : 0;
}

export function resolveOrderTab(order: BuyerOrder): "pending-pay" | "waiting-shipment" | "waiting-receipt" | "waiting-review" | "after-sale" | "history" {
  const subStatuses = order.subOrders.map((sub) => sub.subStatus);
  if (order.orderStatus === "CANCELED" || order.orderStatus === "CLOSED") {
    return "history";
  }
  if (order.orderStatus === "REFUNDING" || order.orderStatus === "REFUNDED" || subStatuses.some((status) => status === "REFUNDING" || status === "REFUNDED")) {
    return "after-sale";
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

export function canConfirmReceipt(order: BuyerOrder): boolean {
  if (resolveOrderTab(order) !== "waiting-receipt") {
    return false;
  }
  if (order.subOrders.length === 0) {
    return false;
  }
  return order.subOrders.every((sub) => sub.subStatus === "SHIPPED" || sub.subStatus === "COMPLETED");
}

export function getAutoReceiveDeadlineMs(order: BuyerOrder): number {
  const shippedAtMs = order.subOrders
    .filter((sub) => sub.subStatus === "SHIPPED")
    .map((sub) => toTimeMs(sub.updatedAt))
    .filter((value) => value > 0)
    .reduce((latest, value) => Math.max(latest, value), 0);

  if (shippedAtMs <= 0) {
    return 0;
  }
  return shippedAtMs + AUTO_RECEIVE_WINDOW_MS;
}

export function getAutoReceiveRemainingMs(order: BuyerOrder, nowMs: number): number {
  const deadlineMs = getAutoReceiveDeadlineMs(order);
  if (deadlineMs <= 0) {
    return 0;
  }
  return Math.max(0, deadlineMs - nowMs);
}
