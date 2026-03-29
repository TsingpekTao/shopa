"use client";

import Link from "next/link";
import { useI18n } from "@shopa/ui";

export default function CartPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";

  return (
    <main className="tb-home" style={{ display: "grid", gap: 16 }}>
      <h1>{isZh ? "购物车" : "Shopping Cart"}</h1>
      <p>{isZh ? "你的购物车商品会显示在这里。" : "Your cart items will be shown here."}</p>
      <Link href="/checkout/confirm" className="tb-btn-main">
        {isZh ? "去结算" : "Go to Checkout"}
      </Link>
    </main>
  );
}
