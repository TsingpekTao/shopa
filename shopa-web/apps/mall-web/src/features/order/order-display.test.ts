import assert from "node:assert/strict";
import { canConfirmReceipt, getAutoReceiveDeadlineMs, getAutoReceiveRemainingMs, resolveOrderTab } from "./order-display.js";

function createOrder(overrides?: Record<string, unknown>) {
  return {
    orderStatus: "FULFILLING",
    paymentStatus: "PAID",
    subOrders: [
      {
        subStatus: "SHIPPED",
        updatedAt: "2026-04-07T12:00:00.000Z"
      }
    ],
    ...overrides
  } as never;
}

function runResolveOrderTabTest() {
  assert.equal(resolveOrderTab(createOrder()), "waiting-receipt");
  assert.equal(
    resolveOrderTab(
      createOrder({
        orderStatus: "COMPLETED",
        subOrders: [{ subStatus: "COMPLETED" }]
      })
    ),
    "waiting-review"
  );
}

function runCanConfirmReceiptTest() {
  assert.equal(canConfirmReceipt(createOrder()), true);
  assert.equal(
    canConfirmReceipt(
      createOrder({
        subOrders: [
          { subStatus: "SHIPPED" },
          { subStatus: "WAIT_SHIP" }
        ]
      })
    ),
    false
  );
  assert.equal(
    canConfirmReceipt(
      createOrder({
        orderStatus: "COMPLETED",
        subOrders: [{ subStatus: "COMPLETED" }]
      })
    ),
    false
  );
}

function runAutoReceiveCountdownTest() {
  const order = createOrder();
  assert.equal(getAutoReceiveDeadlineMs(order), new Date("2026-04-14T12:00:00.000Z").getTime());
  assert.equal(
    getAutoReceiveRemainingMs(order, new Date("2026-04-10T12:00:00.000Z").getTime()),
    4 * 24 * 60 * 60 * 1000
  );
  assert.equal(
    getAutoReceiveRemainingMs(
      createOrder({
        subOrders: [{ subStatus: "WAIT_SHIP", updatedAt: "2026-04-07T12:00:00.000Z" }]
      }),
      Date.now()
    ),
    0
  );
}

runResolveOrderTabTest();
runCanConfirmReceiptTest();
runAutoReceiveCountdownTest();
console.log("order-display tests passed");
