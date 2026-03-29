import Link from "next/link";

type ShopPageProps = {
  params: {
    shopNo: string;
  };
};

export function generateStaticParams() {
  return [{ shopNo: "demo-shop" }, { shopNo: "digital-lab" }];
}

const items = Array.from({ length: 10 }).map((_, index) => ({
  id: `SP-${index + 1}`,
  title: `店铺精选商品 ${index + 1}`,
  price: `${89 + index * 12}`
}));

export default function ShopPage({ params }: ShopPageProps) {
  return (
    <main className="tb-shop-page">
      <section className="tb-shop-hero">
        <h1>店铺：{params.shopNo}</h1>
        <p>品牌直营 | 粉丝 12.9 万 | 近 30 天复购率 34%</p>
        <div>
          <button type="button">关注店铺</button>
          <Link href="/search">查看全部商品</Link>
        </div>
      </section>
      <section className="tb-shop-grid">
        {items.map((item) => (
          <article key={item.id} className="tb-shop-card">
            <Link href={`/item/${item.id}`}>
              <div className="tb-shop-cover" aria-hidden="true" />
              <h3>{item.title}</h3>
              <p>￥{item.price}</p>
            </Link>
          </article>
        ))}
      </section>
    </main>
  );
}
