import Link from "next/link";

type ItemPageProps = {
  params: {
    itemId: string;
  };
};

export function generateStaticParams() {
  return [{ itemId: "P1001" }, { itemId: "P1002" }, { itemId: "SP-1" }];
}

const specs = ["黑色 128G", "黑色 256G", "银色 128G", "银色 256G"];

export default function ItemPage({ params }: ItemPageProps) {
  return (
    <main className="tb-pdp-page">
      <section className="tb-pdp-main">
        <aside className="tb-pdp-gallery">
          <div className="tb-pdp-main-image" aria-hidden="true" />
          <div className="tb-pdp-thumbs">
            {Array.from({ length: 5 }).map((_, idx) => (
              <button key={idx} type="button" aria-label={`预览图 ${idx + 1}`}>
                <span />
              </button>
            ))}
          </div>
        </aside>

        <article className="tb-pdp-info">
          <h1>Shopa 智能旗舰商品（编号 {params.itemId}）</h1>
          <p className="tb-pdp-subtitle">官方直营，支持 7 天无理由退换，假一赔十。</p>
          <div className="tb-pdp-price-box">
            <span>活动价</span>
            <strong>￥2399</strong>
            <em>￥2999</em>
          </div>
          <div className="tb-pdp-specs">
            <span>规格</span>
            <div>
              {specs.map((spec, idx) => (
                <button key={spec} type="button" className={idx === 0 ? "active" : ""}>
                  {spec}
                </button>
              ))}
            </div>
          </div>
          <div className="tb-pdp-actions">
            <Link href="/cart" className="tb-pdp-cart">
              加入购物车
            </Link>
            <Link href="/checkout/confirm" className="tb-pdp-buy">
              立即购买
            </Link>
          </div>
        </article>
      </section>

      <section className="tb-pdp-shop">
        <h2>店铺信息</h2>
        <p>潮流数码仓 | 描述 4.9 | 服务 4.8 | 物流 4.9</p>
        <Link href="/shop/demo-shop">进店逛逛</Link>
      </section>
    </main>
  );
}
