"use client";

import { MouseEvent, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";
import { message } from "antd";
import { useI18n } from "@shopa/ui";
import { addCartItem } from "@/features/cart/api";
import { getBuyerProductDetail, listBuyerAssetReadUrls, listBuyerProductImages } from "@/features/catalog/api";
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

function readMallAccessToken(): string {
  if (typeof window === "undefined") {
    return "";
  }
  const direct = window.localStorage.getItem("shopa_mall_access_token")?.trim() ?? "";
  if (direct) {
    return direct;
  }
  try {
    const raw = window.localStorage.getItem("shopa-mall-auth");
    if (!raw) {
      return "";
    }
    const parsed = JSON.parse(raw) as {
      state?: { tokenPair?: { accessToken?: string } };
    };
    return parsed?.state?.tokenPair?.accessToken?.trim() ?? "";
  } catch {
    return "";
  }
}

export default function ItemDetailPage({ params }: { params: { itemId: string } }) {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const router = useRouter();
  const [messageApi, contextHolder] = message.useMessage();
  const [qty, setQty] = useState(1);
  const [selectedSkuNo, setSelectedSkuNo] = useState("");
  const [selectedImageIndex, setSelectedImageIndex] = useState(0);
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

  const imageAssetIds = useMemo(() => {
    const assetIds = [
      selectedSku?.skuImageAssetId ?? "",
      ...(detailQuery.data?.mainImageAssetIds ?? []),
      ...(detailQuery.data?.detailImageAssetIds ?? [])
    ];
    return Array.from(new Set(assetIds.map((item) => item.trim()).filter((item) => item.length > 0)));
  }, [detailQuery.data?.detailImageAssetIds, detailQuery.data?.mainImageAssetIds, selectedSku?.skuImageAssetId]);

  useEffect(() => {
    setSelectedImageIndex(0);
  }, [selectedSku?.skuImageAssetId, detailQuery.data?.spuNo]);

  const imageAssetUrlsQuery = useQuery({
    queryKey: ["buyer-product-detail-asset-urls", imageAssetIds.join(",")],
    queryFn: () => listBuyerAssetReadUrls(imageAssetIds),
    enabled: imageAssetIds.length > 0,
    staleTime: 60_000,
    refetchOnWindowFocus: false
  });

  const resolvedSpuNo = selectedSku?.spuNo || detailQuery.data?.spuNo || "";
  const resolvedShopNo = selectedSku?.shopNo || detailQuery.data?.shopNo || "";
  const resolvedSkuNo = selectedSku?.skuNo || "";

  const rawSalePrice = selectedSku?.salePrice ?? detailQuery.data?.minSalePrice ?? 0;
  const demoSalePrice = getDemoSalePriceCents(resolvedSpuNo, resolvedSkuNo);
  const currentPrice = rawSalePrice > 0 ? rawSalePrice : demoSalePrice;
  const rawMarketPrice = selectedSku?.marketPrice ?? detailQuery.data?.maxMarketPrice ?? 0;
  const marketPrice =
    rawMarketPrice > currentPrice
      ? rawMarketPrice
      : getDemoMarketPriceCents(currentPrice, resolvedSpuNo, resolvedSkuNo);
  const galleryImages = useMemo(() => {
    const assetUrlMap = imageAssetUrlsQuery.data ?? {};
    const urls = imageAssetIds.map((assetId) => assetUrlMap[assetId]).filter((item): item is string => Boolean(item));
    const fallbackImage = imageQuery.data?.[imageSpuNo] || pickProductImageBySpuNo(detailQuery.data?.spuNo ?? params.itemId);
    if (urls.length > 0) {
      return urls;
    }
    return fallbackImage ? [fallbackImage] : [];
  }, [detailQuery.data?.spuNo, imageAssetIds, imageAssetUrlsQuery.data, imageQuery.data, imageSpuNo, params.itemId]);
  const mainImage = galleryImages[selectedImageIndex] || galleryImages[0] || "";

  const hasResolvedCartPayload = Boolean(resolvedSkuNo && resolvedSpuNo && resolvedShopNo);
  const canSubmit = hasResolvedCartPayload && !detailQuery.isLoading && !submitting;

  async function handleAddToCart() {
    if (!hasResolvedCartPayload) {
      messageApi.warning(isZh ? "当前商品暂无可购买规格" : "No available sku for this product");
      return;
    }
    if (!readMallAccessToken()) {
      messageApi.info(isZh ? "请先登录后再加入购物车" : "Please sign in before adding items to your cart");
      router.push("/login");
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
    if (!hasResolvedCartPayload) {
      messageApi.warning(isZh ? "当前商品暂无可购买规格" : "No available sku for this product");
      return;
    }
    if (!readMallAccessToken()) {
      messageApi.info(isZh ? "请先登录后再立即购买" : "Please sign in before buying now");
      router.push("/login");
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

  function buildConsultHref(): string {
    const query = new URLSearchParams();
    query.set("shop_no", resolvedShopNo);
    query.set("spu_no", resolvedSpuNo);
    query.set("sku_no", resolvedSkuNo);
    query.set("scene_code", "PRE_SALE");
    query.set("mode", "pre-sale");
    const productTitle = detailQuery.data?.title || (isZh ? `商品 #${params.itemId}` : `Product #${params.itemId}`);
    const productSubTitle = detailQuery.data?.subTitle || "";
    const productPrice = `CNY ${formatCnyFromCents(currentPrice)}`;
    if (productTitle.trim()) {
      query.set("product_title", productTitle.trim());
    }
    if (productSubTitle.trim()) {
      query.set("product_subtitle", productSubTitle.trim());
    }
    if (productPrice.trim()) {
      query.set("product_price", productPrice.trim());
    }
    if (mainImage.trim()) {
      query.set("product_image", mainImage.trim());
    }
    return `/me/messages?${query.toString()}`;
  }

  function handleConsultSeller(event: MouseEvent<HTMLAnchorElement>) {
    if (!resolvedShopNo) {
      event.preventDefault();
      messageApi.warning(isZh ? "当前商品暂时无法发起咨询" : "This product is not ready for chat yet");
      return;
    }
    if (!readMallAccessToken()) {
      event.preventDefault();
      messageApi.info(isZh ? "请先登录后再咨询商家" : "Please sign in before chatting with the seller");
      router.push("/login");
    }
  }

  return (
    <main className="tb-pdp-page">
      {contextHolder}
      <section className="tb-pdp-main">
        <div>
          <div className="tb-pdp-main-image" style={mainImage ? { backgroundImage: `url(${mainImage})`, backgroundSize: "cover" } : undefined} />
          <div className="tb-pdp-thumbs">
            {(galleryImages.length > 0 ? galleryImages : [""]).map((imageUrl, index) => (
              <button
                key={`${imageUrl || "empty"}-${index}`}
                type="button"
                aria-label={`preview-${index + 1}`}
                className={index === selectedImageIndex ? "active" : ""}
                onClick={() => setSelectedImageIndex(index)}
              >
                <span style={imageUrl ? { backgroundImage: `url(${imageUrl})`, backgroundSize: "cover" } : undefined} />
              </button>
            ))}
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
            <Link href={buildConsultHref()} className="tb-pdp-chat" onClick={handleConsultSeller}>
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
          <Link href={buildConsultHref()} onClick={handleConsultSeller}>
            {isZh ? "售前咨询" : "Pre-sale Chat"}
          </Link>
          <Link href={`/service/assistant?shop_no=${encodeURIComponent(resolvedShopNo)}&spu_no=${encodeURIComponent(resolvedSpuNo)}&sku_no=${encodeURIComponent(resolvedSkuNo)}`}>
            {isZh ? "AI 导购" : "AI Assistant"}
          </Link>
        </div>
      </section>
    </main>
  );
}
