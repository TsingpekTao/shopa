"use client";

import { useMemo, useState } from "react";
import Link from "next/link";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { message } from "antd";
import { useI18n } from "@shopa/ui";
import { getMyCart, removeCartItems, toggleCartItemChecked, updateCartItemQty } from "@/features/cart/api";
import { getBuyerProductDetail, listBuyerAssetReadUrls, listBuyerProductImages } from "@/features/catalog/api";
import { getCartItemHints } from "@/lib/cart-hints";
import { getDemoMarketPriceCents, getDemoSalePriceCents } from "@/lib/demo-pricing";
import { CartItem } from "@/features/cart/types";
import { formatCnyFromCents } from "@/lib/price";
import { collectProductDetailAssetIds, resolveProductImageUrl, type ProductDetailMap } from "@/lib/product-media";

type SkuEnrich = {
  salePrice: number;
  marketPrice: number;
  skuName: string;
  saleAttrsJson: string;
};

type SpuEnrich = {
  title: string;
  subTitle: string;
  shopNo: string;
  detail: ProductDetailMap[string];
};

function parseSaleAttrsText(raw: string): string {
  const source = (raw || "").trim();
  if (!source) {
    return "";
  }
  try {
    const parsed = JSON.parse(source);
    if (!Array.isArray(parsed)) {
      return "";
    }
    const pairs: string[] = [];
    for (const item of parsed) {
      if (!item || typeof item !== "object") {
        continue;
      }
      const row = item as Record<string, unknown>;
      const key = String(row.attr_name ?? row.name ?? row.key ?? "").trim();
      const value = String(row.attr_value ?? row.value ?? row.val ?? "").trim();
      if (key && value) {
        pairs.push(`${key}: ${value}`);
      } else if (value) {
        pairs.push(value);
      }
    }
    return pairs.join(" · ");
  } catch {
    return "";
  }
}

export default function CartPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [messageApi, contextHolder] = message.useMessage();
  const [updatingSku, setUpdatingSku] = useState<string>("");
  const queryClient = useQueryClient();

  const cartQuery = useQuery({
    queryKey: ["my-cart"],
    queryFn: () => getMyCart(false),
    refetchOnWindowFocus: false,
    staleTime: 20_000
  });

  const items = useMemo(() => cartQuery.data?.items ?? [], [cartQuery.data?.items]);
  const summary = cartQuery.data?.summary;

  const spuNos = useMemo(
    () => Array.from(new Set(items.map((item) => item.spuNo.trim()).filter((item) => item.length > 0))),
    [items]
  );

  const imageMapQuery = useQuery({
    queryKey: ["cart-images", spuNos],
    queryFn: () => listBuyerProductImages(spuNos),
    enabled: spuNos.length > 0,
    staleTime: 60_000,
    refetchOnWindowFocus: false
  });

  const enrichQuery = useQuery({
    queryKey: ["cart-enrich", spuNos],
    enabled: spuNos.length > 0,
    staleTime: 60_000,
    refetchOnWindowFocus: false,
    queryFn: async () => {
      const skuMap: Record<string, SkuEnrich> = {};
      const spuMap: Record<string, SpuEnrich> = {};
      const results = await Promise.allSettled(spuNos.map((spuNo) => getBuyerProductDetail(spuNo)));
      for (const result of results) {
        if (result.status !== "fulfilled") {
          continue;
        }
        const detail = result.value;
        spuMap[detail.spuNo] = {
          title: detail.title,
          subTitle: detail.subTitle,
          shopNo: detail.shopNo,
          detail
        };
        for (const sku of detail.skus) {
          skuMap[sku.skuNo] = {
            salePrice: sku.salePrice,
            marketPrice: sku.marketPrice,
            skuName: sku.skuName,
            saleAttrsJson: sku.saleAttrsJson
          };
        }
      }
      return { skuMap, spuMap };
    }
  });

  const assetIds = useMemo(
    () => Array.from(new Set(items.map((item) => item.skuImageAssetId.trim()).filter((item) => item.length > 0))),
    [items]
  );

  const assetUrlsQuery = useQuery({
    queryKey: ["cart-item-asset-urls", assetIds.join(",")],
    queryFn: () => listBuyerAssetReadUrls(assetIds),
    enabled: assetIds.length > 0,
    staleTime: 60_000,
    refetchOnWindowFocus: false
  });

  const detailMap = useMemo<ProductDetailMap>(() => {
    const spuMap = enrichQuery.data?.spuMap ?? {};
    return Object.values(spuMap).reduce<ProductDetailMap>((acc, item) => {
      if (item.detail?.spuNo) {
        acc[item.detail.spuNo] = item.detail;
      }
      return acc;
    }, {});
  }, [enrichQuery.data?.spuMap]);

  const detailAssetIds = useMemo(() => collectProductDetailAssetIds(detailMap), [detailMap]);

  const detailAssetUrlsQuery = useQuery({
    queryKey: ["cart-detail-asset-urls", detailAssetIds.join(",")],
    queryFn: () => listBuyerAssetReadUrls(detailAssetIds),
    enabled: detailAssetIds.length > 0,
    staleTime: 60_000,
    refetchOnWindowFocus: false
  });

  const displayItems = useMemo(() => {
    const skuMap = enrichQuery.data?.skuMap ?? {};
    const spuMap = enrichQuery.data?.spuMap ?? {};
    const imageMap = imageMapQuery.data ?? {};
    const assetUrlMap = assetUrlsQuery.data ?? {};
    const detailAssetUrlMap = detailAssetUrlsQuery.data ?? {};
    const cartHints = getCartItemHints();

    return items.map((item) => {
      const skuEnrich = skuMap[item.skuNo];
      const spuEnrich = spuMap[item.spuNo];
      const hint = cartHints[item.skuNo];
      const demoSalePrice = getDemoSalePriceCents(item.spuNo, item.skuNo);
      const salePrice =
        item.salePrice > 0
          ? item.salePrice
          : skuEnrich?.salePrice && skuEnrich.salePrice > 0
            ? skuEnrich.salePrice
            : hint?.salePrice && hint.salePrice > 0
              ? hint.salePrice
              : demoSalePrice;
      const marketPrice =
        item.marketPrice > salePrice
          ? item.marketPrice
          : skuEnrich?.marketPrice && skuEnrich.marketPrice > salePrice
            ? skuEnrich.marketPrice
            : hint?.marketPrice && hint.marketPrice > salePrice
              ? hint.marketPrice
              : getDemoMarketPriceCents(salePrice, item.spuNo, item.skuNo);
      const title = item.spuTitle || spuEnrich?.title || hint?.spuTitle || item.skuName || item.skuNo;
      const skuName = item.skuName || skuEnrich?.skuName || hint?.skuName || item.skuNo;
      const subTitle = spuEnrich?.subTitle || hint?.subTitle || "";
      const shopNo = item.shopNo || spuEnrich?.shopNo || hint?.shopNo || "";
      const attrsText = parseSaleAttrsText(item.saleAttrsJson || skuEnrich?.saleAttrsJson || "[]");
      const imageUrl = resolveProductImageUrl(item, assetUrlMap, imageMap, detailMap, detailAssetUrlMap) || encodeURI((hint?.imageUrl ?? "").trim());

      return {
        ...item,
        title,
        skuName,
        subTitle,
        shopNo,
        attrsText,
        salePrice,
        marketPrice,
        imageUrl
      };
    });
  }, [items, enrichQuery.data, imageMapQuery.data, assetUrlsQuery.data, detailMap, detailAssetUrlsQuery.data]);

  const selectedSkuNos = useMemo(
    () => displayItems.filter((item) => item.checked).map((item) => item.skuNo),
    [displayItems]
  );

  const checkedPayableAmount = useMemo(() => {
    const backendValue = summary?.checkedPayableAmount ?? 0;
    if (backendValue > 0) {
      return backendValue;
    }
    return displayItems
      .filter((item) => item.checked)
      .reduce((total, item) => total + item.salePrice * item.qty, 0);
  }, [summary?.checkedPayableAmount, displayItems]);

  async function refreshCart() {
    await queryClient.invalidateQueries({ queryKey: ["my-cart"] });
  }

  async function handleQtyChange(item: CartItem, nextQty: number) {
    if (nextQty < 1 || nextQty > 99) {
      return;
    }
    setUpdatingSku(item.skuNo);
    try {
      await updateCartItemQty(item.skuNo, nextQty);
      await refreshCart();
    } catch (err) {
      console.error(err);
      messageApi.error(isZh ? "修改数量失败" : "Failed to update quantity");
    } finally {
      setUpdatingSku("");
    }
  }

  async function handleToggleChecked(item: CartItem, checked: boolean) {
    setUpdatingSku(item.skuNo);
    try {
      await toggleCartItemChecked(item.skuNo, checked);
      await refreshCart();
    } catch (err) {
      console.error(err);
      messageApi.error(isZh ? "更新勾选状态失败" : "Failed to update checked state");
    } finally {
      setUpdatingSku("");
    }
  }

  async function handleRemove(item: CartItem) {
    setUpdatingSku(item.skuNo);
    try {
      await removeCartItems([item.skuNo]);
      await refreshCart();
      messageApi.success(isZh ? "已移除商品" : "Item removed");
    } catch (err) {
      console.error(err);
      messageApi.error(isZh ? "移除失败" : "Failed to remove item");
    } finally {
      setUpdatingSku("");
    }
  }

  return (
    <main className="tb-cart-page">
      {contextHolder}
      <h1>{isZh ? "购物车" : "Shopping Cart"}</h1>

      <section className="tb-cart-group">
        <header>
          <strong>{isZh ? "全部商品" : "All Items"}</strong>
          <span>{isZh ? `共 ${displayItems.length} 件` : `${displayItems.length} items`}</span>
        </header>

        <div className="tb-cart-list">
          {displayItems.map((item) => {
            const disabled = updatingSku === item.skuNo;
            return (
              <article key={item.skuNo} className="tb-cart-item">
                <label className="tb-cart-check">
                  <input
                    type="checkbox"
                    checked={item.checked}
                    onChange={(event) => handleToggleChecked(item, event.target.checked)}
                    disabled={disabled}
                  />
                </label>
                <div className="tb-cart-cover" style={item.imageUrl ? undefined : { background: "linear-gradient(135deg, #ffe2d1, #ffd1ae)" }}>
                  {item.imageUrl ? (
                    <img
                      src={item.imageUrl}
                      alt={item.title}
                      loading="lazy"
                      decoding="async"
                      referrerPolicy="no-referrer"
                      style={{ width: "100%", height: "100%", objectFit: "cover", display: "block" }}
                    />
                  ) : null}
                </div>
                <div className="tb-cart-info">
                  <h3>{item.title}</h3>
                  <p>{item.skuName}</p>
                  {item.attrsText ? <p className="tb-cart-meta">{item.attrsText}</p> : null}
                  {item.subTitle ? <p className="tb-cart-subtitle">{item.subTitle}</p> : null}
                  <p className="tb-cart-meta">{isZh ? `店铺：${item.shopNo || "--"}` : `Shop: ${item.shopNo || "--"}`}</p>
                  <p className="tb-cart-action">
                    <button
                      type="button"
                      className="tb-btn-sub"
                      style={{ height: 30 }}
                      onClick={() => handleRemove(item)}
                      disabled={disabled}
                    >
                      {isZh ? "移除" : "Remove"}
                    </button>
                  </p>
                </div>
                <div className="tb-cart-price-wrap">
                  <p className="tb-cart-price">{`CNY ${formatCnyFromCents(item.salePrice)}`}</p>
                  {item.marketPrice > item.salePrice && item.salePrice > 0 ? (
                    <p className="tb-cart-price-sub">{`CNY ${formatCnyFromCents(item.marketPrice)}`}</p>
                  ) : null}
                </div>
                <div className="tb-cart-qty">
                  <button type="button" onClick={() => handleQtyChange(item, item.qty - 1)} disabled={disabled}>-</button>
                  <span>{item.qty}</span>
                  <button type="button" onClick={() => handleQtyChange(item, item.qty + 1)} disabled={disabled}>+</button>
                </div>
              </article>
            );
          })}
        </div>

        {cartQuery.isLoading && <p>{isZh ? "加载购物车中..." : "Loading cart..."}</p>}
        {cartQuery.isError && <p>{isZh ? "购物车加载失败，请刷新重试。" : "Failed to load cart."}</p>}
      </section>

      <section className="tb-cart-footer">
        <p>
          {isZh ? "已选" : "Selected"} {summary?.checkedItemCount ?? 0}
        </p>
        <p className="tb-cart-total">
          {isZh ? "合计：" : "Total: "}
          <strong>{`CNY ${formatCnyFromCents(checkedPayableAmount)}`}</strong>
        </p>
        <Link
          href={selectedSkuNos.length > 0 ? "/checkout/confirm" : "#"}
          onClick={(event) => {
            if (selectedSkuNos.length === 0) {
              event.preventDefault();
              messageApi.warning(isZh ? "请先勾选要结算的商品" : "Please select items first");
            }
          }}
        >
          {isZh ? "去结算" : "Checkout"}
        </Link>
      </section>
    </main>
  );
}
