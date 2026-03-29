"use client";

import { useI18n } from "@shopa/ui";

export default function ShopPage({ params }: { params: { shopNo: string } }) {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";

  return (
    <main className="tb-home" style={{ display: "grid", gap: 16 }}>
      <h1>
        {isZh ? "店铺 #" : "Shop #"}
        {params.shopNo}
      </h1>
      <p>{isZh ? "店铺首页预览页面。" : "Storefront preview page."}</p>
    </main>
  );
}
