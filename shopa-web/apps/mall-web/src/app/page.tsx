"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { useI18n } from "@shopa/ui";
import { listBuyerProductImages, listBuyerProducts } from "@/features/catalog/api";
import { BUYER_SORT_BY } from "@/features/catalog/types";
import { pickProductImageBySpuNo } from "@/lib/product-images";
import { formatCnyFromCents } from "@/lib/price";

type HeroSlide = {
  id: string;
  titleZh: string;
  titleEn: string;
  descZh: string;
  descEn: string;
  ctaZh: string;
  ctaEn: string;
  image: string;
};

type FallbackProduct = {
  id: string;
  titleZh: string;
  titleEn: string;
  priceCents: number;
  sold: number;
};

type DisplayProduct = {
  id: string;
  title: string;
  priceCents: number;
  sold: number;
  imageUrl: string;
};

type DisplayProductBase = Omit<DisplayProduct, "imageUrl">;
type MallAccountSnapshot = {
  loggedIn: boolean;
  displayName: string;
  avatarText: string;
};

const AUTH_CHANGED_EVENT = "shopa-mall-auth-changed";

const heroSlides: HeroSlide[] = [
  {
    id: "hero-1",
    titleZh: "春季上新主会场",
    titleEn: "Spring New Arrivals",
    descZh: "服饰、美妆、家居爆品上新，领券立减，限时抢购。",
    descEn: "Fashion, beauty and home new arrivals with instant coupon savings.",
    ctaZh: "马上逛会场",
    ctaEn: "Shop Now",
    image: "https://picsum.photos/seed/shopa-sem-hero-1/1200/560"
  },
  {
    id: "hero-2",
    titleZh: "百亿补贴精选",
    titleEn: "Mega Subsidy Picks",
    descZh: "官方直降叠加跨店优惠，热门单品天天好价。",
    descEn: "Official discounts stacked with cross-store promotions every day.",
    ctaZh: "进入补贴专区",
    ctaEn: "Enter Zone",
    image: "https://picsum.photos/seed/shopa-sem-hero-2/1200/560"
  },
  {
    id: "hero-3",
    titleZh: "品牌旗舰日",
    titleEn: "Brand Flagship Day",
    descZh: "国际品牌旗舰店联动活动，会员专属礼遇。",
    descEn: "Joint campaign from top flagship stores with member-only offers.",
    ctaZh: "领取会员福利",
    ctaEn: "Get Rewards",
    image: "https://picsum.photos/seed/shopa-sem-hero-3/1200/560"
  }
];

const fallbackProductPool: FallbackProduct[] = [
  { id: "P2001", titleZh: "轻量慢跑运动鞋", titleEn: "Lightweight Running Shoes", priceCents: 23900, sold: 9530 },
  { id: "P2002", titleZh: "智能降噪蓝牙耳机", titleEn: "ANC Bluetooth Earbuds", priceCents: 27900, sold: 8210 },
  { id: "P2003", titleZh: "家用空气炸锅 5L", titleEn: "Air Fryer 5L", priceCents: 32900, sold: 5120 },
  { id: "P2004", titleZh: "玻尿酸补水面膜礼盒", titleEn: "Hydrating Mask Set", priceCents: 9900, sold: 16680 },
  { id: "P2005", titleZh: "纯棉柔感四件套", titleEn: "Cotton Bedding Set", priceCents: 26900, sold: 4300 },
  { id: "P2006", titleZh: "65W 双口快充套装", titleEn: "65W Fast Charger Kit", priceCents: 11900, sold: 12090 },
  { id: "P2007", titleZh: "便携手冲咖啡礼盒", titleEn: "Pour-over Coffee Kit", priceCents: 19900, sold: 2950 },
  { id: "P2008", titleZh: "简约皮质通勤包", titleEn: "Leather Commuter Bag", priceCents: 21900, sold: 2710 },
  { id: "P2009", titleZh: "颈椎护托记忆枕", titleEn: "Memory Foam Neck Pillow", priceCents: 13900, sold: 6880 },
  { id: "P2010", titleZh: "智能体脂秤 Pro", titleEn: "Smart Body Scale Pro", priceCents: 12900, sold: 7430 },
  { id: "P2011", titleZh: "防晒透气连帽外套", titleEn: "UV Hooded Jacket", priceCents: 16900, sold: 6080 },
  { id: "P2012", titleZh: "桌面护眼学习灯", titleEn: "Eye-care Desk Lamp", priceCents: 14900, sold: 5320 },
  { id: "P2013", titleZh: "旅行收纳六件套", titleEn: "Travel Organizer Set", priceCents: 5900, sold: 14990 },
  { id: "P2014", titleZh: "儿童益智拼搭积木", titleEn: "Kids Building Blocks", priceCents: 8900, sold: 7940 },
  { id: "P2015", titleZh: "磨砂不粘炒锅 32cm", titleEn: "Non-stick Wok 32cm", priceCents: 17900, sold: 3870 }
];

function ProductCover({
  className,
  imageUrl,
  alt
}: {
  className: string;
  imageUrl: string;
  alt: string;
}) {
  if (!imageUrl) {
    return <div className={className} aria-hidden="true" />;
  }
  return (
    <div className={className}>
      <img src={imageUrl} alt={alt} loading="lazy" decoding="async" referrerPolicy="no-referrer" />
    </div>
  );
}

function tryReadPersistedAccessToken(): string {
  if (typeof window === "undefined") {
    return "";
  }
  const raw = window.localStorage.getItem("shopa-mall-auth");
  if (!raw) {
    return "";
  }
  try {
    const parsed = JSON.parse(raw) as {
      state?: { tokenPair?: { accessToken?: string } };
    };
    return parsed?.state?.tokenPair?.accessToken?.trim() ?? "";
  } catch {
    return "";
  }
}

function tryDecodeIdentityFromToken(token: string): string {
  const parts = token.split(".");
  if (parts.length < 2) {
    return "";
  }
  try {
    const payload = JSON.parse(atob(parts[1])) as Record<string, unknown>;
    const fields = [
      payload.display_name,
      payload.displayName,
      payload.nickname,
      payload.name,
      payload.username,
      payload.account,
      payload.phone,
      payload.email
    ];
    for (const field of fields) {
      if (typeof field === "string" && field.trim()) {
        return field.trim();
      }
    }
    return "";
  } catch {
    return "";
  }
}

function normalizeDisplayName(raw: string): string {
  const value = raw.trim();
  if (!value) {
    return "";
  }
  if (/^\d{11}$/.test(value)) {
    return `${value.slice(0, 3)}****${value.slice(-4)}`;
  }
  if (/^\d+$/.test(value)) {
    return "";
  }
  return value;
}

function readMallAccountSnapshot(locale: "zh-CN" | "en-US"): MallAccountSnapshot {
  if (typeof window === "undefined") {
    return { loggedIn: false, displayName: "", avatarText: "S" };
  }
  const token = window.localStorage.getItem("shopa_mall_access_token")?.trim() || tryReadPersistedAccessToken();
  const loggedIn = token.length > 0;
  const candidates = [
    window.localStorage.getItem("shopa_mall_display_name")?.trim() ?? "",
    window.localStorage.getItem("shopa_mall_last_identifier")?.trim() ?? "",
    tryDecodeIdentityFromToken(token)
  ];
  const displayName = candidates.map((item) => normalizeDisplayName(item)).find((item) => item.length > 0) ?? "";
  const avatarText = (displayName || "S").slice(0, 1).toUpperCase();
  return { loggedIn, displayName, avatarText };
}

export default function HomePage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [hasHydrated, setHasHydrated] = useState(false);
  const [activeHero, setActiveHero] = useState(0);
  const [showBackTop, setShowBackTop] = useState(false);
  const [account, setAccount] = useState<MallAccountSnapshot>({
    loggedIn: false,
    displayName: "",
    avatarText: "S"
  });

  useEffect(() => {
    setHasHydrated(true);
  }, []);

  const productsQuery = useQuery({
    queryKey: ["mall-home-products", 1, 24, BUYER_SORT_BY.SALES_DESC],
    queryFn: async () =>
      listBuyerProducts({
        page: 1,
        pageSize: 24,
        sortBy: BUYER_SORT_BY.SALES_DESC
      }),
    staleTime: 60_000,
    cacheTime: 300_000,
    refetchOnWindowFocus: false
  });

  const categories = useMemo(
    () =>
      isZh
        ? [
            ["女装", "内衣", "配饰"],
            ["男装", "运动", "鞋靴"],
            ["美妆", "个护", "香氛"],
            ["手机", "数码", "电脑"],
            ["母婴", "玩具", "童装"],
            ["零食", "生鲜", "粮油"],
            ["家电", "家装", "家纺"],
            ["图书", "文具", "潮玩"],
            ["箱包", "珠宝", "腕表"],
            ["汽车", "户外", "骑行"]
          ]
        : [
            ["Women", "Innerwear", "Accessories"],
            ["Men", "Sports", "Shoes"],
            ["Beauty", "Care", "Fragrance"],
            ["Phone", "Digital", "Computer"],
            ["Baby", "Toys", "Kidswear"],
            ["Snacks", "Fresh", "Grocery"],
            ["Appliance", "Home", "Textile"],
            ["Books", "Stationery", "Toys"],
            ["Bags", "Jewelry", "Watches"],
            ["Auto", "Outdoor", "Cycling"]
          ],
    [isZh]
  );

  const shortcuts = useMemo(
    () =>
      isZh
        ? ["天猫", "聚划算", "百亿补贴", "品牌馆", "直播", "淘工厂", "同城购", "新人礼"]
        : ["Tmall", "Deals", "Subsidy", "Brands", "Live", "Factory", "Local", "New User"],
    [isZh]
  );

  const baseProducts = useMemo<DisplayProductBase[]>(() => {
    const remoteItems = productsQuery.data?.items ?? [];
    if (remoteItems.length > 0) {
      return remoteItems.map((item) => ({
        id: item.spuNo,
        title: item.title || (isZh ? "未命名商品" : "Untitled Product"),
        priceCents: item.minSalePrice,
        sold: item.soldCount
      }));
    }
    return fallbackProductPool.map((item) => ({
      id: item.id,
      title: isZh ? item.titleZh : item.titleEn,
      priceCents: item.priceCents,
      sold: item.sold
    }));
  }, [isZh, productsQuery.data?.items]);

  const productImageSpuNos = useMemo(
    () => baseProducts.map((item) => item.id).filter((item) => item.trim().length > 0),
    [baseProducts]
  );

  const productImageSpuNosKey = useMemo(() => productImageSpuNos.join(","), [productImageSpuNos]);

  const productImagesQuery = useQuery({
    queryKey: ["mall-home-product-images", productImageSpuNosKey],
    queryFn: async () => listBuyerProductImages(productImageSpuNos),
    enabled: productImageSpuNos.length > 0,
    staleTime: 60_000,
    cacheTime: 300_000,
    refetchOnWindowFocus: false
  });

  const recommendedProducts = useMemo<DisplayProduct[]>(() => {
    const imageMap = productImagesQuery.data ?? {};
    return baseProducts.map((item) => ({
      ...item,
      imageUrl: imageMap[item.id] || pickProductImageBySpuNo(item.id)
    }));
  }, [baseProducts, productImagesQuery.data]);

  const flashProducts = useMemo(() => recommendedProducts.slice(0, 5), [recommendedProducts]);

  useEffect(() => {
    const timer = window.setInterval(() => {
      setActiveHero((prev) => (prev + 1) % heroSlides.length);
    }, 4200);
    return () => window.clearInterval(timer);
  }, []);

  useEffect(() => {
    const onScroll = () => setShowBackTop(window.scrollY > 300);
    onScroll();
    window.addEventListener("scroll", onScroll, { passive: true });
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  useEffect(() => {
    if (!hasHydrated) {
      return;
    }
    const refresh = () => setAccount(readMallAccountSnapshot(locale));
    refresh();
    window.addEventListener("storage", refresh);
    window.addEventListener(AUTH_CHANGED_EVENT, refresh);
    return () => {
      window.removeEventListener("storage", refresh);
      window.removeEventListener(AUTH_CHANGED_EVENT, refresh);
    };
  }, [hasHydrated, locale]);

  const accountLoggedIn = hasHydrated && account.loggedIn;
  const accountName = account.displayName || (isZh ? "已登录用户" : "Signed-in user");

  return (
    <main className="tb-sem-home">
      <section className="tb-sem-top">
        <aside className="tb-sem-category-card">
          <h2>{isZh ? "主题市场" : "Marketplace"}</h2>
          <ul>
            {categories.map((line, idx) => (
              <li key={`${line[0]}-${idx}`}>
                <Link href={`/search?q=${encodeURIComponent(line[0])}`}>{line[0]}</Link>
                <span>/</span>
                <Link href={`/search?q=${encodeURIComponent(line[1])}`}>{line[1]}</Link>
                <span>/</span>
                <Link href={`/search?q=${encodeURIComponent(line[2])}`}>{line[2]}</Link>
              </li>
            ))}
          </ul>
        </aside>

        <section className="tb-sem-stage-card">
          <div className="tb-sem-hero-wrap">
            <div className="tb-sem-hero-track" style={{ transform: `translateX(-${activeHero * 100}%)` }}>
              {heroSlides.map((slide) => (
                <article key={slide.id} className="tb-sem-hero-slide">
                  <div className="tb-sem-hero-media" style={{ backgroundImage: `url(${slide.image})` }} aria-hidden="true" />
                  <div className="tb-sem-hero-copy">
                    <p className="tb-sem-hero-tag">{isZh ? "超级主会场" : "Main Campaign"}</p>
                    <h1>{isZh ? slide.titleZh : slide.titleEn}</h1>
                    <p>{isZh ? slide.descZh : slide.descEn}</p>
                    <Link href="/search" className="tb-sem-hero-cta">
                      {isZh ? slide.ctaZh : slide.ctaEn}
                    </Link>
                  </div>
                </article>
              ))}
            </div>
          </div>
          <div className="tb-sem-hero-dots">
            {heroSlides.map((slide, idx) => (
              <button
                key={slide.id}
                type="button"
                className={idx === activeHero ? "is-active" : ""}
                onClick={() => setActiveHero(idx)}
                aria-label={`${isZh ? "第" : "Slide "}${idx + 1}${isZh ? "张" : ""}`}
              />
            ))}
          </div>

          <div className="tb-sem-shortcuts">
            {shortcuts.map((text) => (
              <Link key={text} href={`/search?q=${encodeURIComponent(text)}`}>
                {text}
              </Link>
            ))}
          </div>
        </section>

        <aside className="tb-sem-account-card">
          <header>
            <div className="tb-sem-avatar" aria-hidden="true">
              {accountLoggedIn ? account.avatarText : "S"}
            </div>
            <div>
              <h3>{accountLoggedIn ? (isZh ? `Hi，${accountName}` : `Hi, ${accountName}`) : isZh ? "Hi，欢迎来到 Shopa" : "Hi, welcome to Shopa"}</h3>
              <p>{accountLoggedIn ? (isZh ? "你可以查看订单、地址和个人资料" : "You can view your orders, addresses and profile.") : isZh ? "登录解锁专属优惠和订单服务" : "Login for offers and order services"}</p>
            </div>
          </header>
          <div className="tb-sem-account-actions">
            {accountLoggedIn ? (
              <>
                <Link href="/me/profile" className="tb-sem-account-action is-primary">{isZh ? "个人中心" : "My Profile"}</Link>
                <Link href="/me/orders" className="tb-sem-account-action is-ghost">{isZh ? "我的订单" : "My Orders"}</Link>
              </>
            ) : (
              <>
                <Link href="/login" className="tb-sem-account-action is-primary">{isZh ? "登录" : "Login"}</Link>
                <Link href="/register" className="tb-sem-account-action is-ghost">{isZh ? "注册" : "Register"}</Link>
              </>
            )}
          </div>

          <div className="tb-sem-news">
            <h4>{isZh ? "公告速递" : "News"}</h4>
            <ul>
              <li>
                <Link href="/help?tab=official">{isZh ? "平台满减活动入口上线" : "Promotion zone now online"}</Link>
              </li>
              <li>
                <Link href="/help?tab=merchant">{isZh ? "商家服务升级，响应更快" : "Merchant support upgraded"}</Link>
              </li>
              <li>
                <Link href="/help?tab=history">{isZh ? "订单中心支持近一年查询" : "1-year order query available"}</Link>
              </li>
            </ul>
          </div>
        </aside>
      </section>

      <section className="tb-sem-flash">
        <div className="tb-sem-flash-title">
          <strong>{isZh ? "限时抢购" : "Flash Deals"}</strong>
          <span>{isZh ? "距结束" : "Ends in"} 03:21:18</span>
        </div>
        <div className="tb-sem-flash-list">
          {flashProducts.map((item) => (
            <Link key={item.id} href={`/item/${item.id}`} className="tb-sem-flash-item">
              <ProductCover className="tb-sem-flash-cover" imageUrl={item.imageUrl} alt={item.title} />
              <h3>{item.title}</h3>
              <p>{`¥${formatCnyFromCents(item.priceCents)}`}</p>
            </Link>
          ))}
        </div>
      </section>

      <section className="tb-sem-floor">
        <header>
          <h2>{isZh ? "猜你喜欢" : "Recommended For You"}</h2>
          <Link href="/search">{isZh ? "查看更多" : "View More"}</Link>
        </header>
        <div className="tb-sem-product-grid">
          {recommendedProducts.map((item) => (
            <article key={item.id} className="tb-sem-product-card">
              <Link href={`/item/${item.id}`}>
                <ProductCover className="tb-sem-product-cover" imageUrl={item.imageUrl} alt={item.title} />
                <h3>{item.title}</h3>
                <p>{`¥${formatCnyFromCents(item.priceCents)}`}</p>
                <small>{isZh ? `近30天已售 ${item.sold}` : `${item.sold} sold in 30d`}</small>
              </Link>
            </article>
          ))}
        </div>
      </section>

      <aside className="tb-sem-tools" aria-label={isZh ? "快捷工具" : "Quick Tools"}>
        <a href="/cart">
          <span>{isZh ? "购物车" : "Cart"}</span>
        </a>
        <a href="/me/orders">
          <span>{isZh ? "订单" : "Orders"}</span>
        </a>
        <a href="/help?tab=official">
          <span>{isZh ? "客服" : "Help"}</span>
        </a>
        <button
          type="button"
          className={`tb-sem-top-btn${showBackTop ? " is-show" : ""}`}
          onClick={() => window.scrollTo({ top: 0, behavior: "smooth" })}
        >
          {isZh ? "顶部" : "Top"}
        </button>
      </aside>
    </main>
  );
}


