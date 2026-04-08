import assert from "node:assert/strict";
import {
  getSettlementRefreshQueryKeys,
  hasPaymentRedirectSignal,
  isPaidOrder,
  paymentAttemptStorageKey,
  shouldShowPayCountdown
} from "./payment-state.js";

function runIsPaidOrderTest() {
  assert.equal(
    isPaidOrder({
      paymentStatus: "PAID",
      orderStatus: "PENDING_PAY"
    } as never),
    true
  );
  assert.equal(
    isPaidOrder({
      paymentStatus: "UNPAID",
      orderStatus: "FULFILLING"
    } as never),
    true
  );
  assert.equal(
    isPaidOrder({
      paymentStatus: "UNPAID",
      orderStatus: "CANCELED"
    } as never),
    false
  );
}

function runSettlementRefreshQueryKeysTest() {
  assert.deepEqual(
    getSettlementRefreshQueryKeys({
      paymentStatus: "PAID",
      orderStatus: "PAID"
    } as never),
    [
      ["buyer-orders", "all"],
      ["my-cart"],
      ["shell", "cart"]
    ]
  );
  assert.deepEqual(
    getSettlementRefreshQueryKeys({
      paymentStatus: "UNPAID",
      orderStatus: "CANCELED"
    } as never),
    [["buyer-orders", "all"]]
  );
}

function runHasPaymentRedirectSignalTest() {
  const redirectParams = new URLSearchParams({
    trade_no: "202604080001",
    out_trade_no: "PAY1770000000000000000"
  });
  assert.equal(hasPaymentRedirectSignal(redirectParams), true);
  assert.equal(hasPaymentRedirectSignal(new URLSearchParams({ foo: "bar" })), false);
  assert.equal(hasPaymentRedirectSignal(null), false);
}

function runShouldShowPayCountdownTest() {
  assert.equal(
    shouldShowPayCountdown(
      {
        orderStatus: "PENDING_PAY",
        paymentStatus: "UNPAID",
        payDeadlineAt: "2026-04-08T20:00:00+08:00"
      } as never,
      false
    ),
    true
  );
  assert.equal(
    shouldShowPayCountdown(
      {
        orderStatus: "FULFILLING",
        paymentStatus: "PAID",
        payDeadlineAt: "2026-04-08T20:00:00+08:00"
      } as never,
      false
    ),
    false
  );
  assert.equal(
    shouldShowPayCountdown(
      {
        orderStatus: "PENDING_PAY",
        paymentStatus: "UNPAID",
        payDeadlineAt: "2026-04-08T20:00:00+08:00"
      } as never,
      true
    ),
    false
  );
}

function runPaymentAttemptStorageKeyTest() {
  assert.equal(paymentAttemptStorageKey("ORD202604080001"), "checkout_pay_attempt_ORD202604080001");
}

runIsPaidOrderTest();
runSettlementRefreshQueryKeysTest();
runHasPaymentRedirectSignalTest();
runShouldShowPayCountdownTest();
runPaymentAttemptStorageKeyTest();
console.log("payment-state tests passed");
