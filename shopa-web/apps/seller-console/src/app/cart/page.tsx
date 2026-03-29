import Link from "next/link";

const cartItems = [
  { sku: "SKU1001", name: "轻便跑鞋", spec: "黑色 42", price: "189", qty: 1 },
  { sku: "SKU1002", name: "无线耳机", spec: "白色 标准版", price: "139", qty: 2 }
];

export default function CartPage() {
  return (
    <main className="tb-cart-page">
      <h1>购物车</h1>
      <section className="tb-cart-group">
        <header>
          <strong>潮流数码仓</strong>
          <span>店铺优惠可与平台券叠加</span>
        </header>
        <div className="tb-cart-list">
          {cartItems.map((item) => (
            <article key={item.sku} className="tb-cart-item">
              <label className="tb-cart-check">
                <input type="checkbox" defaultChecked aria-label={`选择 ${item.name}`} />
                <span className="tb-sr-only">选择</span>
              </label>
              <div className="tb-cart-cover" aria-hidden="true" />
              <div className="tb-cart-info">
                <h3>{item.name}</h3>
                <p>{item.spec}</p>
              </div>
              <p className="tb-cart-price">￥{item.price}</p>
              <div className="tb-cart-qty">
                <button type="button" aria-label="减少数量">
                  -
                </button>
                <span>{item.qty}</span>
                <button type="button" aria-label="增加数量">
                  +
                </button>
              </div>
            </article>
          ))}
        </div>
      </section>
      <footer className="tb-cart-footer">
        <p>已选 3 件商品</p>
        <p className="tb-cart-total">
          合计：<strong>￥467</strong>
        </p>
        <Link href="/checkout/confirm">去结算</Link>
      </footer>
    </main>
  );
}
