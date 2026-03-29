"use client";

import Link from "next/link";
import { useI18n } from "@shopa/ui";

export default function ItemDetailPage({ params }: { params: { itemId: string } }) {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";

  return (
    <main className="tb-home" style={{ display: "grid", gap: 16 }}>
      <h1>{isZh ? "商品详情" : "Product Detail"}</h1>
      <div className="tb-product-cover" style={{ width: 320, height: 320 }} aria-hidden="true" />
      <h2>
        {isZh ? "商品 #" : "Product #"}
        {params.itemId}
      </h2>
      <p>{isZh ? "用于前后端联调的示例详情页。" : "Mock detail page for frontend integration testing."}</p>
      <p style={{ fontSize: 24, fontWeight: 700 }}>￥199</p>
      <div style={{ display: "flex", gap: 12 }}>
        <Link href="/cart" className="tb-btn-main">
          {isZh ? "加入购物车" : "Add to Cart"}
        </Link>
        <Link href="/checkout/confirm" className="tb-btn-sub">
          {isZh ? "立即购买" : "Buy Now"}
        </Link>
      </div>
    </main>
  );
}
