"use client";

import Link from "next/link";
import { useI18n } from "@shopa/ui";

export default function CheckoutConfirmPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";

  return (
    <main className="tb-home" style={{ display: "grid", gap: 16 }}>
      <h1>{isZh ? "确认订单" : "Checkout Confirmation"}</h1>
      <p>{isZh ? "在本页确认订单信息、收货地址与支付方式。" : "Confirm order info, address, and payment method on this page."}</p>
      <Link href="/" className="tb-btn-sub">
        {isZh ? "返回首页" : "Back to Home"}
      </Link>
    </main>
  );
}
