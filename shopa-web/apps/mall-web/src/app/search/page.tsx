"use client";

import { useMemo } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { useI18n } from "@shopa/ui";
import { searchProducts } from "@/features/search/api";
import { formatCnyFromCents } from "@/lib/price";

function ProductCover({ title, coverUrl }: { title: string; coverUrl: string }) {
  if (!coverUrl) {
    return <div className="tb-product-cover" aria-hidden="true" />;
  }

  return (
    <div className="tb-product-cover">
      <img src={coverUrl} alt={title} loading="lazy" decoding="async" referrerPolicy="no-referrer" />
    </div>
  );
}

function ShopAvatar({ label, coverUrl }: { label: string; coverUrl: string }) {
  const avatarText = (label || "SHOP").trim().slice(0, 2).toUpperCase();

  return (
    <div className="tb-shop-search-avatar">
      {coverUrl ? <img src={coverUrl} alt={label} loading="lazy" decoding="async" referrerPolicy="no-referrer" /> : null}
      {!coverUrl ? <span>{avatarText}</span> : null}
    </div>
  );
}

function ShopPreviewTile({
  title,
  coverUrl,
  spuNo,
  minPrice,
  metaText
}: {
  title: string;
  coverUrl: string;
  spuNo: string;
  minPrice: number;
  metaText: string;
}) {
  return (
    <Link href={`/item/${encodeURIComponent(spuNo)}`} className="tb-shop-search-preview-card">
      <div className="tb-shop-search-preview-media">
        {coverUrl ? <img src={coverUrl} alt={title} loading="lazy" decoding="async" referrerPolicy="no-referrer" /> : null}
      </div>
      <p>{formatCnyFromCents(minPrice)}</p>
      <span>{metaText}</span>
    </Link>
  );
}

function normalizeKeyword(value: string): string {
  return value.trim().toLowerCase();
}

export default function SearchPage() {
  const params = useSearchParams();
  const keyword = params.get("q")?.trim() ?? "";
  const searchTab = params.get("tab") === "shop" ? "shop" : "item";
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";

  const searchQuery = useQuery({
    queryKey: ["mall-search-products", keyword],
    queryFn: () =>
      searchProducts({
        query: keyword,
        pageSize: 24
      }),
    staleTime: 30_000,
    refetchOnWindowFocus: false
  });

  const items = searchQuery.data?.items ?? [];
  const shops = searchQuery.data?.shops ?? [];
  const source = searchQuery.data?.source ?? "search-svc";

  const preferredShop = useMemo(() => {
    const normalizedKeyword = normalizeKeyword(keyword);
    if (!normalizedKeyword) {
      return undefined;
    }

    const exactMatch = shops.find((shop) => {
      const displayName = shop.shopDisplayName || shop.shopName;
      return (
        normalizeKeyword(displayName) === normalizedKeyword ||
        normalizeKeyword(shop.shopName) === normalizedKeyword ||
        normalizeKeyword(shop.shopNo) === normalizedKeyword
      );
    });
    if (exactMatch) {
      return exactMatch;
    }
    if (shops.length === 1) {
      return shops[0];
    }
    return undefined;
  }, [keyword, shops]);

  const titleText = keyword
    ? isZh
      ? `\u641c\u7d22\u201c${keyword}\u201d`
      : `Search: ${keyword}`
    : searchTab === "shop"
      ? isZh
        ? "\u5e97\u94fa\u641c\u7d22"
        : "Shop search"
      : isZh
        ? "\u5168\u90e8\u5546\u54c1"
        : "All products";

  const showProductsAsPrimary = searchTab === "item";
  const showShopEmpty = searchTab === "shop" && !searchQuery.isLoading && !searchQuery.isError && shops.length === 0;
  const showGlobalEmpty =
    searchTab === "item" && !searchQuery.isLoading && !searchQuery.isError && items.length === 0 && shops.length === 0;

  return (
    <main className="tb-home tb-search-page">
      <header className="tb-search-head">
        <h1>{titleText}</h1>
        <p style={{ color: "#666", fontSize: 14 }}>
          {searchQuery.isLoading
            ? isZh
              ? "\u6b63\u5728\u641c\u7d22\u5546\u54c1\u548c\u5e97\u94fa..."
              : "Searching products and shops..."
            : searchTab === "shop"
              ? isZh
                ? `\u627e\u5230 ${shops.length} \u5bb6\u76f8\u5173\u5e97\u94fa`
                : `${shops.length} related shops found`
              : isZh
                ? `\u627e\u5230 ${items.length} \u4ef6\u5546\u54c1`
                : `${items.length} products found`}
          {searchTab === "item" && shops.length > 0
            ? isZh
              ? `\uff0c\u76f8\u5173\u5e97\u94fa ${shops.length} \u5bb6`
              : `, ${shops.length} related shops`
            : ""}
          {source === "catalog-fallback"
            ? isZh
              ? "\uff0c\u5f53\u524d\u5c55\u793a\u7684\u662f\u76ee\u5f55\u56de\u9000\u7ed3\u679c"
              : ", currently using catalog fallback"
            : ""}
        </p>

        <div className="tb-search-mode-summary">
          <Link href={`/search?q=${encodeURIComponent(keyword)}&tab=item`} className={showProductsAsPrimary ? "is-active" : ""}>
            {isZh ? "\u5546\u54c1" : "Items"}
          </Link>
          <Link href={`/search?q=${encodeURIComponent(keyword)}&tab=shop`} className={!showProductsAsPrimary ? "is-active" : ""}>
            {isZh ? "\u5e97\u94fa" : "Shops"}
          </Link>
        </div>

        {preferredShop && !showProductsAsPrimary ? (
          <section className="tb-search-shop-hero">
            <div className="tb-search-shop-hero-main">
              <ShopAvatar
                label={preferredShop.shopDisplayName || preferredShop.shopName || preferredShop.shopNo}
                coverUrl={preferredShop.coverUrl}
              />
              <div className="tb-search-shop-hero-copy">
                <span className="tb-search-shop-hero-kicker">{isZh ? "\u5df2\u4e3a\u4f60\u627e\u5230\u76f8\u5173\u5e97\u94fa" : "Shop result found"}</span>
                <strong>{preferredShop.shopDisplayName || preferredShop.shopName || preferredShop.shopNo}</strong>
                <p>
                  {isZh
                    ? `\u5171\u5339\u914d\u5230 ${preferredShop.matchedProductCount} \u4ef6\u76f8\u5173\u5546\u54c1\uff0c\u53ef\u4ee5\u76f4\u63a5\u8fdb\u5e97\u6216\u5148\u770b\u4e3b\u63a8\u5546\u54c1\u3002`
                    : `${preferredShop.matchedProductCount} matching products found in this shop.`}
                </p>
                <div className="tb-search-shop-stats">
                  <span>{isZh ? "\u76f8\u5173\u5546\u54c1" : "Matched items"}: {preferredShop.matchedProductCount}</span>
                  <span>
                    {isZh ? "\u5e97\u94fa\u7cbe\u9009" : "Featured picks"}: {preferredShop.previewItems.length}
                  </span>
                  <span>{isZh ? "\u5e97\u94fa\u7f16\u53f7" : "Shop No"}: {preferredShop.shopNo}</span>
                </div>
              </div>
            </div>
            <div className="tb-search-shop-actions">
              <Link
                href={`/shop/${encodeURIComponent(preferredShop.shopNo)}?from=search`}
                className="tb-search-shop-action-primary"
              >
                {isZh ? "\u8fdb\u5165\u5e97\u94fa" : "Visit shop"}
              </Link>
              {preferredShop.sampleSpuNo ? (
                <Link
                  href={`/item/${encodeURIComponent(preferredShop.sampleSpuNo)}`}
                  className="tb-search-shop-action-secondary"
                >
                  {isZh ? "\u770b\u4e3b\u63a8\u5546\u54c1" : "View featured item"}
                </Link>
              ) : null}
            </div>
            {preferredShop.previewItems.length > 0 ? (
              <div className="tb-search-shop-hero-preview">
                {preferredShop.previewItems.slice(0, 3).map((item) => (
                  <ShopPreviewTile
                    key={item.spuNo}
                    title={item.title}
                    coverUrl={item.coverUrl}
                    spuNo={item.spuNo}
                    minPrice={item.minPrice}
                    metaText={item.salesCount > 0 ? (isZh ? `\u5df2\u552e ${item.salesCount}` : `${item.salesCount} sold`) : item.title}
                  />
                ))}
              </div>
            ) : null}
          </section>
        ) : null}
      </header>

      {searchQuery.isError ? (
        <section className="tb-empty-state">
          <h2>{isZh ? "\u641c\u7d22\u6682\u65f6\u5931\u8d25" : "Search failed"}</h2>
          <p>
            {isZh
              ? "\u8bf7\u7a0d\u540e\u91cd\u8bd5\uff0c\u6216\u786e\u8ba4 search-svc\u3001catalog-svc\u3001seller-shop-svc \u5df2\u6b63\u5e38\u542f\u52a8\u3002"
              : "Please try again later or make sure search-svc, catalog-svc and seller-shop-svc are running."}
          </p>
        </section>
      ) : null}

      {!searchQuery.isLoading && !searchQuery.isError && !showProductsAsPrimary && shops.length > 0 ? (
        <section className="tb-search-head">
          <h2 style={{ margin: "0 0 10px" }}>{isZh ? "\u5e97\u94fa\u7ed3\u679c" : "Shop results"}</h2>
          <div className="tb-shop-search-list">
            {shops.map((shop) => (
              <article key={shop.shopNo} className="tb-shop-search-row">
                <div className="tb-shop-search-row-main">
                  <div className="tb-shop-search-row-brand">
                    <ShopAvatar label={shop.shopDisplayName || shop.shopName || shop.shopNo} coverUrl={shop.coverUrl} />
                    <div className="tb-shop-search-row-copy">
                      <div className="tb-shop-search-row-head">
                        <strong>{shop.shopDisplayName || shop.shopName || shop.shopNo}</strong>
                        <span className="tb-shop-search-row-badge">{isZh ? "\u5e97\u94fa" : "Shop"}</span>
                      </div>
                      <div className="tb-shop-search-row-metrics">
                        <span>{isZh ? `\u5173\u952e\u8bcd\u547d\u4e2d ${shop.matchedProductCount} \u6b3e` : `${shop.matchedProductCount} keyword matches`}</span>
                        <span>{isZh ? `\u7cbe\u9009 ${shop.previewItems.length} \u6b3e` : `${shop.previewItems.length} featured picks`}</span>
                        <span>{isZh ? `\u5e97\u94fa\u7f16\u53f7 ${shop.shopNo}` : `Shop No ${shop.shopNo}`}</span>
                      </div>
                      <p>
                        {isZh
                          ? shop.previewItems.length > 0
                            ? "\u5de6\u4fa7\u5c55\u793a\u5e97\u94fa\u6838\u5fc3\u4fe1\u606f\uff0c\u53f3\u4fa7\u76f4\u63a5\u9884\u89c8\u4e3b\u63a8\u5546\u54c1\uff0c\u6d4f\u89c8\u8def\u5f84\u66f4\u50cf\u6dd8\u5b9d\u5e97\u94fa\u641c\u7d22\u9875\u3002"
                            : "\u5df2\u627e\u5230\u76f8\u5173\u5e97\u94fa\uff0c\u53ef\u76f4\u63a5\u8fdb\u5e97\u7ee7\u7eed\u6d4f\u89c8\u3002"
                          : shop.previewItems.length > 0
                            ? "Shop identity stays on the left, with featured product previews on the right for faster scanning."
                            : "Matching shop found. Enter the shop to continue browsing."}
                      </p>
                      <div className="tb-shop-search-row-actions">
                        <Link href={`/shop/${encodeURIComponent(shop.shopNo)}?from=search`}>
                          {isZh ? "\u8fdb\u5165\u5e97\u94fa" : "Visit shop"}
                        </Link>
                        {shop.sampleSpuNo ? (
                          <Link href={`/item/${encodeURIComponent(shop.sampleSpuNo)}`}>
                            {isZh ? "\u770b\u4e3b\u63a8\u5546\u54c1" : "View featured item"}
                          </Link>
                        ) : null}
                      </div>
                    </div>
                  </div>
                </div>
                <div className="tb-shop-search-row-preview">
                  {shop.previewItems.length > 0 ? (
                    shop.previewItems.map((item) => (
                      <ShopPreviewTile
                        key={item.spuNo}
                        title={item.title}
                        coverUrl={item.coverUrl}
                        spuNo={item.spuNo}
                        minPrice={item.minPrice}
                        metaText={item.salesCount > 0 ? (isZh ? `\u5df2\u552e ${item.salesCount}` : `${item.salesCount} sold`) : item.title}
                      />
                    ))
                  ) : (
                    <div className="tb-shop-search-preview-empty">
                      <strong>{isZh ? "\u6682\u65f6\u6ca1\u6709\u83b7\u53d6\u5230\u5546\u54c1\u9884\u89c8" : "No previews yet"}</strong>
                      <span>{isZh ? "\u53ef\u4ee5\u76f4\u63a5\u8fdb\u5e97\u67e5\u770b\u5168\u90e8\u5546\u54c1" : "Enter the shop to browse all products."}</span>
                    </div>
                  )}
                </div>
              </article>
            ))}
          </div>
        </section>
      ) : null}

      {showProductsAsPrimary && shops.length > 0 ? (
        <section className="tb-search-head">
          <div className="tb-search-section-head">
            <h2>{isZh ? "\u76f8\u5173\u5e97\u94fa" : "Related shops"}</h2>
            <span>
              {isZh
                ? "\u641c\u7d22\u5546\u54c1\u65f6\uff0c\u4ecd\u7136\u4f1a\u4fdd\u7559\u8fd9\u4e9b\u5e97\u94fa\u5165\u53e3\uff0c\u65b9\u4fbf\u76f4\u63a5\u8fdb\u5e97\u6bd4\u8f83\u3002"
                : "Shop entrances stay visible in item search so buyers can jump straight into the storefront."}
            </span>
          </div>
          <div className="tb-shop-search-related">
            {shops.slice(0, 3).map((shop) => (
              <article key={`related-${shop.shopNo}`} className="tb-shop-search-related-card">
                <div className="tb-shop-search-related-head">
                  <ShopAvatar label={shop.shopDisplayName || shop.shopName || shop.shopNo} coverUrl={shop.coverUrl} />
                  <div>
                    <strong>{shop.shopDisplayName || shop.shopName || shop.shopNo}</strong>
                    <span>{isZh ? `\u547d\u4e2d ${shop.matchedProductCount} \u6b3e` : `${shop.matchedProductCount} matches`}</span>
                  </div>
                </div>
                <div className="tb-shop-search-related-actions">
                  <Link href={`/shop/${encodeURIComponent(shop.shopNo)}?from=search`}>
                    {isZh ? "\u8fdb\u5165\u5e97\u94fa" : "Visit shop"}
                  </Link>
                </div>
              </article>
            ))}
          </div>
        </section>
      ) : null}

      {showShopEmpty ? (
        <section className="tb-empty-state">
          <h2>{isZh ? "\u6ca1\u6709\u627e\u5230\u76f8\u5173\u5e97\u94fa" : "No shops found"}</h2>
          <p>
            {keyword
              ? isZh
                ? `\u6682\u65f6\u6ca1\u6709\u627e\u5230\u548c\u201c${keyword}\u201d\u76f8\u5173\u7684\u5e97\u94fa\uff0c\u4f60\u53ef\u4ee5\u5207\u6362\u5230\u5546\u54c1\u641c\u7d22\u7ee7\u7eed\u770b\u3002`
                : `No shops matched "${keyword}". You can switch back to item search.`
              : isZh
                ? "\u8bf7\u8f93\u5165\u5e97\u94fa\u540d\u79f0\u6216\u5e97\u94fa\u7f16\u53f7\u3002"
                : "Try a shop name or shop number."}
          </p>
          {keyword ? (
            <Link href={`/search?q=${encodeURIComponent(keyword)}&tab=item`} className="tb-primary-link">
              {isZh ? "\u5207\u6362\u5230\u5546\u54c1\u641c\u7d22" : "Switch to item search"}
            </Link>
          ) : null}
        </section>
      ) : null}

      {showGlobalEmpty ? (
        <section className="tb-empty-state">
          <h2>{isZh ? "\u6ca1\u6709\u627e\u5230\u76f8\u5173\u5546\u54c1\u6216\u5e97\u94fa" : "No products or shops found"}</h2>
          <p>
            {keyword
              ? isZh
                ? `\u6ca1\u6709\u627e\u5230\u548c\u201c${keyword}\u201d\u76f8\u5173\u7684\u7ed3\u679c\uff0c\u6362\u4e2a\u5173\u952e\u8bcd\u518d\u8bd5\u8bd5\u3002`
                : `No results matched "${keyword}". Try another keyword.`
              : isZh
                ? "\u8bf7\u8f93\u5165\u66f4\u660e\u786e\u7684\u5173\u952e\u8bcd\u8bd5\u8bd5\u3002"
                : "Try a more specific keyword."}
          </p>
          <Link href="/" className="tb-primary-link">
            {isZh ? "\u8fd4\u56de\u9996\u9875" : "Back to home"}
          </Link>
        </section>
      ) : null}

      {items.length > 0 ? (
        <section className="tb-search-product-section">
          {!showProductsAsPrimary ? (
            <div className="tb-search-section-head">
              <h2>{isZh ? "\u76f8\u5173\u5546\u54c1" : "Related products"}</h2>
              <span>{isZh ? "\u9009\u62e9\u5546\u54c1\u641c\u7d22\u65f6\u4f1a\u5c06\u8fd9\u4e9b\u7ed3\u679c\u4f5c\u4e3a\u4e3b\u7ed3\u679c\u5c55\u793a" : "These results become the primary list in item mode."}</span>
            </div>
          ) : null}

          <div className="tb-product-grid">
            {items.map((item) => (
              <article className="tb-product-card" key={item.spuNo}>
                <Link href={`/item/${encodeURIComponent(item.spuNo)}`} className="tb-product-link">
                  <ProductCover title={item.title} coverUrl={item.coverUrl} />
                  <h3>{item.title}</h3>
                  <p>{formatCnyFromCents(item.minPrice)}</p>
                </Link>
                <div className="tb-product-card-footer">
                  {item.shopNo ? (
                    <Link
                      href={`/shop/${encodeURIComponent(item.shopNo)}?from=search`}
                      className="tb-product-shop-link"
                    >
                      {isZh ? `\u5e97\u94fa\uff1a${item.shopName || item.shopNo}` : `Shop: ${item.shopName || item.shopNo}`}
                    </Link>
                  ) : (
                    <small>{isZh ? `\u9500\u91cf ${item.salesCount}` : `${item.salesCount} sold`}</small>
                  )}
                  {item.shopNo ? <small>{isZh ? `\u9500\u91cf ${item.salesCount}` : `${item.salesCount} sold`}</small> : null}
                </div>
              </article>
            ))}
          </div>
        </section>
      ) : null}
    </main>
  );
}
