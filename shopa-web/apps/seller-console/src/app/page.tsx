"use client";

import Link from "next/link";

const categoryItems = [
  "女装",
  "男装",
  "鞋包配饰",
  "美妆个护",
  "手机数码",
  "家电家居",
  "食品生鲜",
  "母婴亲子",
  "运动户外",
  "图书文创",
  "珠宝首饰",
  "宠物生活"
];

const secKillItems = [
  { id: "S1001", name: "夏季防晒服", price: "79", market: "129", progress: 71 },
  { id: "S1002", name: "无线降噪耳机", price: "139", market: "249", progress: 64 },
  { id: "S1003", name: "轻便通勤跑鞋", price: "189", market: "299", progress: 52 },
  { id: "S1004", name: "不粘炒锅套装", price: "169", market: "269", progress: 83 },
  { id: "S1005", name: "小型投影仪", price: "329", market: "499", progress: 46 }
];

const recommendItems = [
  { id: "P1001", title: "简约短袖 T 恤", badge: "新品", price: "69" },
  { id: "P1002", title: "复古帆布托特包", badge: "包邮", price: "129" },
  { id: "P1003", title: "速干运动短裤", badge: "满减", price: "88" },
  { id: "P1004", title: "家用空气炸锅", badge: "热销", price: "299" },
  { id: "P1005", title: "玻尿酸保湿面霜", badge: "店铺券", price: "118" },
  { id: "P1006", title: "桌面收纳抽屉盒", badge: "精品", price: "45" },
  { id: "P1007", title: "机械键盘 87 键", badge: "好评高", price: "259" },
  { id: "P1008", title: "陶瓷早餐碗套装", badge: "限量", price: "59" },
  { id: "P1009", title: "抗皱连衣裙", badge: "爆款", price: "199" },
  { id: "P1010", title: "儿童拼装玩具", badge: "人气", price: "149" }
];

export default function HomePage() {
  return (
    <main className="tb-home">
      <section className="tb-home-hero">
        <aside className="tb-home-cates">
          <h2>全部分类</h2>
          <ul>
            {categoryItems.map((item) => (
              <li key={item}>
                <Link href="/search">{item}</Link>
              </li>
            ))}
          </ul>
        </aside>
        <article className="tb-home-banner">
          <span className="tb-chip">超级品牌日</span>
          <h1>春季焕新会场</h1>
          <p>满 300 减 40，跨店可叠加店铺券，今晚 20:00 继续放量。</p>
          <div className="tb-home-banner-actions">
            <Link href="/search" className="tb-btn-main">
              去抢优惠
            </Link>
            <Link href="/shop/demo-shop" className="tb-btn-sub">
              逛品牌店
            </Link>
          </div>
        </article>
        <aside className="tb-home-side">
          <h3>今日推荐店铺</h3>
          <ul>
            <li>
              <strong>原生家居馆</strong>
              <span>满 199 减 30</span>
            </li>
            <li>
              <strong>轻运动实验室</strong>
              <span>新品低至 8 折</span>
            </li>
            <li>
              <strong>潮流数码仓</strong>
              <span>24 期免息</span>
            </li>
          </ul>
        </aside>
      </section>

      <section className="tb-floor">
        <header className="tb-floor-header">
          <h2>限时秒杀</h2>
          <Link href="/search">查看更多</Link>
        </header>
        <div className="tb-sec-grid">
          {secKillItems.map((item) => (
            <article className="tb-sec-card" key={item.id}>
              <div className="tb-sec-cover" aria-hidden="true" />
              <h3>{item.name}</h3>
              <p className="tb-price-row">
                <strong>￥{item.price}</strong>
                <span>￥{item.market}</span>
              </p>
              <p className="tb-sale-progress">已抢 {item.progress}%</p>
            </article>
          ))}
        </div>
      </section>

      <section className="tb-floor">
        <header className="tb-floor-header">
          <h2>猜你喜欢</h2>
          <Link href="/search">换一批</Link>
        </header>
        <div className="tb-product-grid">
          {recommendItems.map((item) => (
            <article className="tb-product-card" key={item.id}>
              <Link href={`/item/${item.id}`} className="tb-product-link">
                <div className="tb-product-cover" aria-hidden="true" />
                <span className="tb-product-badge">{item.badge}</span>
                <h3>{item.title}</h3>
                <p>￥{item.price}</p>
              </Link>
            </article>
          ))}
        </div>
      </section>
    </main>
  );
}
