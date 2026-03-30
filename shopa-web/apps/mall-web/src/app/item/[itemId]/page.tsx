"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { message } from "antd";
import { useI18n } from "@shopa/ui";
import { addCartItem } from "@/features/cart/api";
import { getBuyerProductDetail, listBuyerProductImages } from "@/features/catalog/api";
import { upsertCartItemHint } from "@/lib/cart-hints";
import { getDemoMarketPriceCents, getDemoSalePriceCents } from "@/lib/demo-pricing";
import { formatCnyFromCents } from "@/lib/price";
import { pickProductImageBySpuNo } from "@/lib/product-images";

function getDefaultAddressId(): number {
  if (typeof window === "undefined") {
    return 1;
  }
  const raw = window.localStorage.getItem("shopa_mall_default_address_id");
  const value = Number(raw ?? "1");
  if (!Number.isFinite(value) || value <= 0) {
    return 1;
  }
  return Math.trunc(value);
}

export default function ItemDetailPage({ params }: { params: { itemId: string } }) {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const router = useRouter();
  const [messageApi, contextHolder] = message.useMessage();
  const [qty, setQty] = useState(1);
  const [selectedSkuNo, setSelectedSkuNo] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const detailQuery = useQuery({
    queryKey: ["buyer-product-detail", params.itemId],
    queryFn: () => getBuyerProductDetail(params.itemId),
    staleTime: 60_000,
    refetchOnWindowFocus: false
  });

  const imageSpuNo = detailQuery.data?.spuNo || params.itemId;
  const imageQuery = useQuery({
    queryKey: ["buyer-product-detail-image", imageSpuNo],
    queryFn: async () => listBuyerProductImages([imageSpuNo]),
    staleTime: 60_000,
    refetchOnWindowFocus: false
  });

  useEffect(() => {
    const firstSku = detailQuery.data?.skus?.[0]?.skuNo ?? "";
    if (firstSku && !selectedSkuNo) {
      setSelectedSkuNo(firstSku);
    }
  }, [detailQuery.data?.skus, selectedSkuNo]);

  const selectedSku = useMemo(() => {
    const skus = detailQuery.data?.skus ?? [];
    if (skus.length === 0) {
      return undefined;
    }
    return skus.find((item) => item.skuNo === selectedSkuNo) ?? skus[0];
  }, [detailQuery.data?.skus, selectedSkuNo]);

  const resolvedSpuNo = selectedSku?.spuNo || detailQuery.data?.spuNo || params.itemId;
  const resolvedShopNo = selectedSku?.shopNo || detailQuery.data?.shopNo || "DEMO_SHOP";
  const resolvedSkuNo = selectedSku?.skuNo || `SKU_${resolvedSpuNo}`;

  const rawSalePrice = selectedSku?.salePrice ?? detailQuery.data?.minSalePrice ?? 0;
  const demoSalePrice = getDemoSalePriceCents(resolvedSpuNo, resolvedSkuNo);
  const currentPrice = rawSalePrice > 0 ? rawSalePrice : demoSalePrice;
  const rawMarketPrice = selectedSku?.marketPrice ?? detailQuery.data?.maxMarketPrice ?? 0;
  const marketPrice =
    rawMarketPrice > currentPrice
      ? rawMarketPrice
      : getDemoMarketPriceCents(currentPrice, resolvedSpuNo, resolvedSkuNo);
  const mainImage = imageQuery.data?.[imageSpuNo] || pickProductImageBySpuNo(detailQuery.data?.spuNo ?? params.itemId);

  const canSubmit = Boolean(resolvedSkuNo) && !submitting;

  async function handleAddToCart() {
    if (!resolvedSkuNo) {
      messageApi.warning(isZh ? "当前商品暂无可购买规格" : "No available sku for this product");
      return;
    }
    setSubmitting(true);
    try {
      await addCartItem({
        skuNo: resolvedSkuNo,
        spuNo: resolvedSpuNo,
        shopNo: resolvedShopNo,
        qty,
        checked: true
      });
      upsertCartItemHint({
        skuNo: resolvedSkuNo,
        spuNo: resolvedSpuNo,
        shopNo: resolvedShopNo,
        spuTitle: detailQuery.data?.title || (isZh ? `商品 #${params.itemId}` : `Product #${params.itemId}`),
        subTitle: detailQuery.data?.subTitle || "",
        skuName: selectedSku?.skuName || resolvedSkuNo,
        imageUrl: mainImage || pickProductImageBySpuNo(resolvedSpuNo),
        salePrice: currentPrice,
        marketPrice,
        updatedAt: Date.now()
      });
      messageApi.success(isZh ? "已加入购物车" : "Added to cart");
    } catch (err) {
      console.error(err);
      messageApi.error(isZh ? "加入购物车失败，请稍后重试" : "Failed to add cart item");
    } finally {
      setSubmitting(false);
    }
  }

  function handleBuyNow() {
    if (!resolvedSkuNo) {
      messageApi.warning(isZh ? "当前商品暂无可购买规格" : "No available sku for this product");
      return;
    }
    const query = new URLSearchParams({
      mode: "buy-now",
      sku_no: resolvedSkuNo,
      spu_no: resolvedSpuNo,
      shop_no: resolvedShopNo,
      qty: String(qty),
      address_id: String(getDefaultAddressId())
    });
    router.push(`/checkout/confirm?${query.toString()}`);
  }

  return (
    <main className="tb-pdp-page">
      {contextHolder}
      <section className="tb-pdp-main">
        <div>
          <div className="tb-pdp-main-image" style={mainImage ? { backgroundImage: `url(${mainImage})`, backgroundSize: "cover" } : undefined} />
          <div className="tb-pdp-thumbs">
            <button type="button" aria-label="preview">
              <span style={mainImage ? { backgroundImage: `url(${mainImage})`, backgroundSize: "cover" } : undefined} />
            </button>
          </div>
        </div>

        <div className="tb-pdp-info">
          <h1>{detailQuery.data?.title || (isZh ? `商品 #${params.itemId}` : `Product #${params.itemId}`)}</h1>
          <p className="tb-pdp-subtitle">
            {detailQuery.data?.subTitle || (isZh ? "商品详情支持加入购物车与立即购买。" : "Product detail with add-to-cart and buy-now.")}
          </p>

          <div className="tb-pdp-price-box">
            <span>{isZh ? "到手价" : "Price"}</span>
            <strong>{`CNY ${formatCnyFromCents(currentPrice)}`}</strong>
            <em>{`CNY ${formatCnyFromCents(marketPrice)}`}</em>
          </div>

          <div className="tb-pdp-specs">
            <span>{isZh ? "规格" : "Sku"}</span>
            <div>
              {(detailQuery.data?.skus ?? []).map((sku) => {
                const active = sku.skuNo === (selectedSku?.skuNo ?? "");
                return (
                  <button
                    key={sku.skuNo}
                    type="button"
                    className={active ? "active" : ""}
                    onClick={() => setSelectedSkuNo(sku.skuNo)}
                  >
                    {sku.skuName || sku.skuNo}
                  </button>
                );
              })}
            </div>
          </div>

          <div className="tb-pdp-specs">
            <span>{isZh ? "数量" : "Quantity"}</span>
            <div>
              <button type="button" onClick={() => setQty((prev) => Math.max(1, prev - 1))}>-</button>
              <button type="button" className="active">{qty}</button>
              <button type="button" onClick={() => setQty((prev) => Math.min(99, prev + 1))}>+</button>
            </div>
          </div>

          <div className="tb-pdp-actions">
            <button type="button" className="tb-pdp-cart" onClick={handleAddToCart} disabled={!canSubmit}>
              {submitting ? (isZh ? "处理中..." : "Processing...") : isZh ? "加入购物车" : "Add to Cart"}
            </button>
            <button type="button" className="tb-pdp-buy" onClick={handleBuyNow} disabled={!canSubmit}>
              {isZh ? "立即购买" : "Buy Now"}
            </button>
            <Link
              href={`/me/messages?shop_no=${encodeURIComponent(resolvedShopNo)}&spu_no=${encodeURIComponent(resolvedSpuNo)}&sku_no=${encodeURIComponent(resolvedSkuNo)}`}
              className="tb-pdp-chat"
            >
              {isZh ? "咨询商家" : "Consult Seller"}
            </Link>
          </div>

          {detailQuery.isLoading && <p className="tb-pdp-subtitle">{isZh ? "正在加载商品信息..." : "Loading product..."}</p>}
          {detailQuery.isError && <p className="tb-pdp-subtitle">{isZh ? "商品信息加载失败，请稍后重试。" : "Failed to load product."}</p>}
        </div>
      </section>

      <section className="tb-pdp-shop">
        <h2>{isZh ? "店铺信息" : "Shop"}</h2>
        <p>{detailQuery.data?.shopNo || "UNKNOWN_SHOP"}</p>
        <div className="tb-pdp-shop-actions">
          <Link href={`/shop/${encodeURIComponent(detailQuery.data?.shopNo || "UNKNOWN_SHOP")}`}>
            {isZh ? "进入店铺" : "Visit shop"}
          </Link>
          <Link href={`/me/messages?shop_no=${encodeURIComponent(resolvedShopNo)}&spu_no=${encodeURIComponent(resolvedSpuNo)}`}>
            {isZh ? "售前咨询" : "Pre-sale Chat"}
          </Link>
        </div>
      </section>
    </main>
  );
}
