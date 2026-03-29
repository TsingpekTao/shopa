"use client";

import Link from "next/link";
import { useSearchParams } from "next/navigation";

const resultItems = Array.from({ length: 12 }).map((_, idx) => {
  const id = `P${1200 + idx}`;
  return {
    id,
    title: `Shopa 推荐商品 ${idx + 1} 号`,
    shop: idx % 2 === 0 ? "原生家居馆" : "潮流数码仓",
    price: `${69 + idx * 9}`,
    sold: `${120 + idx * 17}`
  };
});

export default function SearchPage() {
  const searchParams = useSearchParams();
  const keyword = searchParams.get("q")?.trim() || "热销";

  return (
    <main className="tb-search-page">
      <section className="tb-search-head">
        <p>
          搜索关键词：<strong>{keyword}</strong>
        </p>
        <div className="tb-search-filters" aria-label="筛选条件">
          <button type="button">综合排序</button>
          <button type="button">销量优先</button>
          <button type="button">价格从低到高</button>
          <button type="button">价格从高到低</button>
          <button type="button">发货地</button>
        </div>
      </section>

      <section className="tb-search-grid" aria-label="搜索结果商品列表">
        {resultItems.map((item) => (
          <article key={item.id} className="tb-search-card">
            <Link href={`/item/${item.id}`}>
              <div className="tb-search-cover" aria-hidden="true" />
              <h3>{item.title}</h3>
              <p className="tb-search-price">￥{item.price}</p>
              <p className="tb-search-meta">
                <span>{item.shop}</span>
                <span>{item.sold} 人付款</span>
              </p>
            </Link>
          </article>
        ))}
      </section>
    </main>
  );
}
