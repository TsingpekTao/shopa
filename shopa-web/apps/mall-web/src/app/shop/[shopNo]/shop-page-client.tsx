"use client";

import { useMemo } from "react";
import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { searchProducts } from "@/features/search/api";
import { getMallShopDetail, getShopStoreCategories } from "@/features/shop/api";
import type { SearchProductCard, SearchProductsResult } from "@/features/search/types";
import type { MallShopDetail, ShopStoreCategory } from "@/features/shop/types";
import { formatCnyFromCents } from "@/lib/price";

type ShopSortKey = "default" | "sales" | "price-asc" | "price-desc";

type ShopPageClientProps = {
  shopNo: string;
  keyword: string;
  fromSearch: boolean;
  isZh: boolean;
  selectedPrimaryCategoryId: number;
  selectedSecondaryCategoryId: number;
  selectedSort: ShopSortKey;
  initialShop: MallShopDetail | null;
  initialCategories: ShopStoreCategory[];
  initialProducts: SearchProductsResult;
};

function normalizeSearchValue(value: string) {
  return value.trim();
}

function formatCompactCount(isZh: boolean, value: number) {
  if (value <= 0) {
    return "0";
  }
  if (!isZh) {
    return new Intl.NumberFormat("en-US", { notation: "compact", maximumFractionDigits: 1 }).format(value);
  }
  if (value >= 10_000) {
    const compact = value >= 100_000 ? (value / 10_000).toFixed(0) : (value / 10_000).toFixed(1);
    return `${compact.replace(/\.0$/, "")}万`;
  }
  return String(value);
}

function formatSalesLabel(isZh: boolean, salesCount: number) {
  return isZh ? `已售 ${formatCompactCount(true, salesCount)}` : `${formatCompactCount(false, salesCount)} sold`;
}

function buildShopHref(params: {
  shopNo: string;
  keyword: string;
  fromSearch: boolean;
  primaryCategoryId?: number;
  secondaryCategoryId?: number;
  sort?: ShopSortKey;
}) {
  const searchParams = new URLSearchParams();
  const keyword = normalizeSearchValue(params.keyword);
  if (keyword) {
    searchParams.set("q", keyword);
  }
  if (params.fromSearch) {
    searchParams.set("from", "search");
  }
  if (params.primaryCategoryId && params.primaryCategoryId > 0) {
    searchParams.set("primary", String(params.primaryCategoryId));
  }
  if (params.secondaryCategoryId && params.secondaryCategoryId > 0) {
    searchParams.set("secondary", String(params.secondaryCategoryId));
  }
  if (params.sort && params.sort !== "default") {
    searchParams.set("sort", params.sort);
  }
  const query = searchParams.toString();
  return `/shop/${encodeURIComponent(params.shopNo)}${query ? `?${query}` : ""}`;
}

function ShopProductCover({ title, coverUrl }: { title: string; coverUrl: string }) {
  if (!coverUrl) {
    return <div className="tb-shop-simple-card-cover" aria-hidden="true" />;
  }

  return (
    <div className="tb-shop-simple-card-cover">
      <img src={coverUrl} alt={title} loading="lazy" decoding="async" referrerPolicy="no-referrer" />
    </div>
  );
}

function HiddenInput({ name, value }: { name: string; value: string | undefined }) {
  if (!value) {
    return null;
  }
  return <input type="hidden" name={name} value={value} />;
}

function findPrimaryCategory(categories: ShopStoreCategory[], primaryId: number, secondaryId: number) {
  if (primaryId > 0) {
    const direct = categories.find((item) => item.id === primaryId);
    if (direct) {
      return direct;
    }
  }
  if (secondaryId > 0) {
    return categories.find((item) => item.children.some((child) => child.id === secondaryId)) ?? null;
  }
  return null;
}

function resolveBreadcrumb(isZh: boolean, primary: ShopStoreCategory | null, secondary: ShopStoreCategory | null) {
  const trail = [isZh ? "全部商品" : "All products"];
  if (primary) {
    trail.push(primary.name);
  }
  if (secondary) {
    trail.push(secondary.name);
  }
  return trail;
}

function countAllStoreProducts(categories: ShopStoreCategory[]) {
  return categories.reduce((sum, category) => sum + category.productCount, 0);
}

export default function ShopPageClient({
  shopNo,
  keyword,
  fromSearch,
  isZh,
  selectedPrimaryCategoryId,
  selectedSecondaryCategoryId,
  selectedSort,
  initialShop,
  initialCategories,
  initialProducts
}: ShopPageClientProps) {
  const shopQuery = useQuery({
    queryKey: ["mall-shop-detail", shopNo],
    queryFn: () => getMallShopDetail(shopNo),
    staleTime: 60_000,
    refetchOnWindowFocus: false,
    initialData: initialShop ?? undefined
  });

  const categoriesQuery = useQuery({
    queryKey: ["mall-shop-store-categories", shopNo],
    queryFn: () => getShopStoreCategories(shopNo),
    staleTime: 60_000,
    refetchOnWindowFocus: false,
    initialData: initialCategories
  });

  const categories = categoriesQuery.data ?? [];
  const activePrimary = useMemo(
    () => findPrimaryCategory(categories, selectedPrimaryCategoryId, selectedSecondaryCategoryId),
    [categories, selectedPrimaryCategoryId, selectedSecondaryCategoryId]
  );
  const activeSecondary = useMemo(
    () => activePrimary?.children.find((item) => item.id === selectedSecondaryCategoryId) ?? null,
    [activePrimary, selectedSecondaryCategoryId]
  );
  const activeStoreCategoryId = activeSecondary?.id ?? activePrimary?.id ?? 0;
  const activeStoreCategoryLevel = activeSecondary ? 2 : activePrimary ? 1 : undefined;

  const productQuery = useQuery({
    queryKey: ["mall-shop-products", shopNo, keyword, activeStoreCategoryId, activeStoreCategoryLevel],
    queryFn: () =>
      searchProducts({
        query: keyword,
        shopNo,
        pageSize: 80,
        storeCategoryId: activeStoreCategoryId || undefined,
        storeCategoryLevel: activeStoreCategoryLevel as 1 | 2 | undefined
      }),
    staleTime: 30_000,
    refetchOnWindowFocus: false,
    initialData: initialProducts
  });

  const shopName = shopQuery.data?.shopDisplayName || shopQuery.data?.shopName || (isZh ? "店铺" : "Shop");
  const items = useMemo(() => {
    const source = [...(productQuery.data?.items ?? [])];
    if (selectedSort === "sales") {
      return source.sort((a, b) => b.salesCount - a.salesCount);
    }
    if (selectedSort === "price-asc") {
      return source.sort((a, b) => a.minPrice - b.minPrice);
    }
    if (selectedSort === "price-desc") {
      return source.sort((a, b) => b.minPrice - a.minPrice);
    }
    return source;
  }, [productQuery.data?.items, selectedSort]);
  const totalSales = items.reduce((sum, item) => sum + item.salesCount, 0);
  const breadcrumb = resolveBreadcrumb(isZh, activePrimary, activeSecondary);
  const categoryProductTotal = useMemo(() => countAllStoreProducts(categories), [categories]);
  const scopeProductCount = activeSecondary?.productCount ?? activePrimary?.productCount ?? categoryProductTotal;
  const allProductsCount = keyword ? items.length : categoryProductTotal || items.length;
  const hasRealCategories = categories.length > 0;

  const sortOptions: Array<{ key: ShopSortKey; label: string }> = [
    { key: "default", label: isZh ? "综合" : "Featured" },
    { key: "sales", label: isZh ? "销量" : "Best selling" },
    { key: "price-asc", label: isZh ? "价格从低到高" : "Price low to high" },
    { key: "price-desc", label: isZh ? "价格从高到低" : "Price high to low" }
  ];

  return (
    <main className="tb-home tb-shop-page tb-shop-simple-page">
      <section className="tb-shop-simple-header">
        <div className="tb-shop-simple-header-main">
          <span>{isZh ? "店内分类选购" : "In-shop categories"}</span>
          <h1>{shopName}</h1>
          <p>
            {isZh
              ? "左边直接切店内分类，右边立即展示对应商品；进入店铺后默认先看全部商品。"
              : "Use the category rail on the left and the product grid on the right. The shop opens on all products by default."}
          </p>
        </div>

        <div className="tb-shop-simple-header-actions">
          {fromSearch ? (
            <Link href="/search?tab=shop" className="tb-shop-simple-backlink">
              {isZh ? "返回搜索" : "Back to search"}
            </Link>
          ) : null}

          <form className="tb-shop-simple-search" action={`/shop/${encodeURIComponent(shopNo)}`} method="get">
            <HiddenInput name="from" value={fromSearch ? "search" : undefined} />
            <HiddenInput name="primary" value={activePrimary ? String(activePrimary.id) : undefined} />
            <HiddenInput name="secondary" value={activeSecondary ? String(activeSecondary.id) : undefined} />
            <HiddenInput name="sort" value={selectedSort !== "default" ? selectedSort : undefined} />
            <input
              name="q"
              defaultValue={keyword}
              placeholder={isZh ? "搜索店内商品" : "Search in this shop"}
              autoComplete="off"
            />
            <button type="submit">{isZh ? "搜索" : "Search"}</button>
          </form>
        </div>
      </section>

      <section className="tb-shop-simple-layout">
        <aside className="tb-shop-simple-sidebar" aria-label={isZh ? "店内一级分类" : "Primary categories"}>
          <div className="tb-shop-simple-sidebar-head">
            <strong>{isZh ? "店内分类" : "Store categories"}</strong>
            <span>{categories.length}</span>
          </div>

          <div className="tb-shop-simple-category-list">
            <Link
              href={buildShopHref({
                shopNo,
                keyword,
                fromSearch,
                sort: selectedSort
              })}
              className={`tb-shop-simple-category-item ${!activePrimary && !activeSecondary ? "is-active" : ""}`}
            >
              <span>{isZh ? "全部商品" : "All products"}</span>
              <small>{allProductsCount}</small>
            </Link>

            {categories.map((category) => (
              <div key={category.id} className="tb-shop-simple-category-group">
                <Link
                  href={buildShopHref({
                    shopNo,
                    keyword,
                    fromSearch,
                    primaryCategoryId: category.id,
                    sort: selectedSort
                  })}
                  className={`tb-shop-simple-category-item ${activePrimary?.id === category.id ? "is-active" : ""}`}
                >
                  <span>{category.name}</span>
                  <small>{category.productCount}</small>
                </Link>

                {category.children.length > 0 ? (
                  <div className={`tb-shop-simple-subcategory-list ${activePrimary?.id === category.id ? "is-open" : ""}`}>
                    {category.children.map((child) => (
                      <Link
                        key={child.id}
                        href={buildShopHref({
                          shopNo,
                          keyword,
                          fromSearch,
                          primaryCategoryId: category.id,
                          secondaryCategoryId: child.id,
                          sort: selectedSort
                        })}
                        className={`tb-shop-simple-subcategory-item ${activeSecondary?.id === child.id ? "is-active" : ""}`}
                      >
                        <span>{child.name}</span>
                        <small>{child.productCount}</small>
                      </Link>
                    ))}
                  </div>
                ) : null}
              </div>
            ))}
          </div>
        </aside>

        <section className="tb-shop-simple-content">
          <div className="tb-shop-simple-toolbar">
            <div className="tb-shop-simple-breadcrumb">
              {breadcrumb.map((segment, index) => (
                <span key={`${segment}-${index}`} className="tb-shop-simple-breadcrumb-item">
                  {index === breadcrumb.length - 1 ? <strong>{segment}</strong> : <span>{segment}</span>}
                  {index < breadcrumb.length - 1 ? <i aria-hidden="true">&gt;</i> : null}
                </span>
              ))}
            </div>

            <div className="tb-shop-simple-toolbar-main">
              <div>
                <h2>{activeSecondary?.name || activePrimary?.name || (isZh ? "全部商品" : "All products")}</h2>
                <p>
                  {keyword
                    ? isZh
                      ? `关键词“${keyword}”，当前匹配 ${items.length} 件商品`
                      : `${items.length} items matched "${keyword}"`
                    : isZh
                      ? `${items.length} 件商品，累计销量 ${formatCompactCount(true, totalSales)}`
                      : `${items.length} items, ${formatCompactCount(false, totalSales)} in total sales`}
                </p>
              </div>

              <div className="tb-shop-simple-sort" aria-label={isZh ? "排序" : "Sorting"}>
                {sortOptions.map((sortOption) => (
                  <Link
                    key={sortOption.key}
                    href={buildShopHref({
                      shopNo,
                      keyword,
                      fromSearch,
                      primaryCategoryId: activePrimary?.id,
                      secondaryCategoryId: activeSecondary?.id,
                      sort: sortOption.key
                    })}
                    className={`tb-shop-simple-sort-link ${selectedSort === sortOption.key ? "is-active" : ""}`}
                  >
                    {sortOption.label}
                  </Link>
                ))}
              </div>
            </div>
          </div>

          <div className="tb-shop-simple-stats" aria-label={isZh ? "当前分类摘要" : "Current category summary"}>
            <article>
              <span>{isZh ? "当前范围" : "Current scope"}</span>
              <strong>{activeSecondary?.name || activePrimary?.name || (isZh ? "全部商品" : "All products")}</strong>
            </article>
            <article>
              <span>{isZh ? "分类商品数" : "Items in scope"}</span>
              <strong>{keyword ? items.length : scopeProductCount || items.length}</strong>
            </article>
            <article>
              <span>{isZh ? "二级分类" : "Subcategories"}</span>
              <strong>{activePrimary ? activePrimary.children.length : hasRealCategories ? categories.length : 0}</strong>
            </article>
          </div>

          {shopQuery.isError ? (
            <div className="tb-shop-simple-empty">
              <h3>{isZh ? "店铺信息暂时没有同步" : "Store details are syncing"}</h3>
              <p>{isZh ? "下方商品列表仍然可以浏览。" : "You can still browse the products below."}</p>
            </div>
          ) : null}

          {productQuery.isError ? (
            <div className="tb-shop-simple-empty">
              <h3>{isZh ? "商品列表暂时未加载" : "Products are temporarily unavailable"}</h3>
              <p>{isZh ? "请稍后再试。" : "Please try again in a moment."}</p>
            </div>
          ) : null}

          {!productQuery.isLoading && !productQuery.isError && items.length === 0 ? (
            <div className="tb-shop-simple-empty">
              <h3>{isZh ? "这个分类暂时没有商品" : "No products in this section yet"}</h3>
              <p>
                {keyword
                  ? isZh
                    ? `没有找到与“${keyword}”相关的店内商品。`
                    : `No in-shop products matched "${keyword}".`
                  : categories.length === 0
                    ? isZh
                      ? "店铺还没有设置店内分类，当前先展示全部商品。"
                      : "This shop has not set store categories yet."
                    : isZh
                      ? "分类会保留展示，商品上架后会出现在这里。"
                      : "The category stays visible and products will appear here once assigned."}
              </p>
              {(keyword || activePrimary || activeSecondary) && (
                <div className="tb-shop-simple-empty-actions">
                  {keyword ? (
                    <Link
                      href={buildShopHref({
                        shopNo,
                        keyword: "",
                        fromSearch,
                        primaryCategoryId: activePrimary?.id,
                        secondaryCategoryId: activeSecondary?.id,
                        sort: selectedSort
                      })}
                      className="tb-shop-simple-backlink"
                    >
                      {isZh ? "清空搜索词" : "Clear search"}
                    </Link>
                  ) : null}
                  {(activePrimary || activeSecondary) ? (
                    <Link
                      href={buildShopHref({
                        shopNo,
                        keyword,
                        fromSearch,
                        sort: selectedSort
                      })}
                      className="tb-shop-simple-backlink"
                    >
                      {isZh ? "回到全部商品" : "Back to all products"}
                    </Link>
                  ) : null}
                </div>
              )}
            </div>
          ) : null}

          {items.length > 0 ? (
            <div className="tb-shop-simple-grid">
              {items.map((item: SearchProductCard) => (
                <article className="tb-shop-simple-card" key={item.spuNo}>
                  <Link href={`/item/${encodeURIComponent(item.spuNo)}`}>
                    <ShopProductCover title={item.title} coverUrl={item.coverUrl} />
                    <div className="tb-shop-simple-card-copy">
                      {activeSecondary || activePrimary ? (
                        <span className="tb-shop-simple-card-tag">{activeSecondary?.name || activePrimary?.name}</span>
                      ) : null}
                      <strong>{item.title}</strong>
                      <p>{formatCnyFromCents(item.minPrice)}</p>
                      <div className="tb-shop-simple-card-meta">
                        <span>{formatSalesLabel(isZh, item.salesCount)}</span>
                      </div>
                    </div>
                  </Link>
                </article>
              ))}
            </div>
          ) : null}
        </section>
      </section>
    </main>
  );
}
