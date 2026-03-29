"use client";

import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useI18n } from "@shopa/ui";
import { formatCnyFromCents } from "@/lib/price";

const mockResult = Array.from({ length: 12 }).map((_, i) => ({
  id: `R${i + 1}`,
  titleZh: `搜索商品 ${i + 1}`,
  titleEn: `Search Product ${i + 1}`,
  priceCents: (59 + i * 10) * 100
}));

export default function SearchPage() {
  const params = useSearchParams();
  const keyword = params.get("q") ?? "";
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";

  return (
    <main className="tb-home" style={{ display: "grid", gap: 16 }}>
      <h1>
        {isZh ? "搜索：" : "Search: "}
        {keyword || (isZh ? "全部" : "all")}
      </h1>
      <div className="tb-product-grid">
        {mockResult.map((item) => (
          <article className="tb-product-card" key={item.id}>
            <Link href={`/item/${item.id}`} className="tb-product-link">
              <div className="tb-product-cover" aria-hidden="true" />
              <h3>{isZh ? item.titleZh : item.titleEn}</h3>
              <p>{`¥${formatCnyFromCents(item.priceCents)}`}</p>
            </Link>
          </article>
        ))}
      </div>
    </main>
  );
}
