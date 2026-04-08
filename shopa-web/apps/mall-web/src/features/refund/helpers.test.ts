import { expect, test } from "vitest";
import { formatAfterSaleStatusLabel, hasRefundableSubOrder } from "./helpers";
import type { BuyerOrder, BuyerOrderSub } from "@/features/order/types";

function buildOrderWithSub(subStatus: BuyerOrderSub["subStatus"]): BuyerOrder {
  return {
    orderNo: "ORD-1",
    userId: 1,
    orderStatus: "PAID",
    paymentStatus: "PAID",
    amount: {
      goodsAmount: 0,
      freightAmount: 0,
      discountAmount: 0,
      payableAmount: 0,
      paidAmount: 0,
      pointsDiscountAmount: 0
    },
    address: null,
    reservationNo: "",
    payDeadlineAt: "",
    paidAt: "",
    closedAt: "",
    createdAt: "",
    updatedAt: "",
    version: 1,
    subOrders: [
      {
        subOrderNo: "SUB-1",
        orderNo: "ORD-1",
        shopNo: "SHOP-1",
        subStatus,
        amount: {
          goodsAmount: 0,
          freightAmount: 0,
          discountAmount: 0,
          payableAmount: 0,
          paidAmount: 0,
          pointsDiscountAmount: 0
        },
        sellerRemark: "",
        buyerRemark: "",
        createdAt: "",
        updatedAt: "",
        items: [],
        pointsUsed: 0,
        pointsDiscountAmount: 0
      }
    ],
    buyerRemark: "",
    cancelReasonCode: "",
    pointsReservationNo: "",
    pointsUsed: 0,
    pointsDiscountAmount: 0,
    pointsRuleSnapshotJson: "",
    pointsRuleSnapshotDigest: ""
  };
}

test("formatAfterSaleStatusLabel returns localized strings", () => {
  expect(formatAfterSaleStatusLabel("PENDING_SELLER_REVIEW", true)).toBe("待商家审核");
  expect(formatAfterSaleStatusLabel("REFUNDED", false)).toBe("Refunded");
});

test("hasRefundableSubOrder recognizes WAIT_SHIP and PAID sub orders", () => {
  const paidOrder = buildOrderWithSub("PAID");
  const waitShipOrder = buildOrderWithSub("WAIT_SHIP");
  const notEligibleOrder = buildOrderWithSub("SHIPPED");

  expect(hasRefundableSubOrder(paidOrder)).toBe(true);
  expect(hasRefundableSubOrder(waitShipOrder)).toBe(true);
  expect(hasRefundableSubOrder(notEligibleOrder)).toBe(false);
});
