"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useMutation, useQuery } from "@tanstack/react-query";
import { message } from "antd";
import { useI18n } from "@shopa/ui";
import { prepareCheckout } from "@/features/cart/api";
import { CheckoutSnapshot, CheckoutSnapshotItem } from "@/features/cart/types";
import { getBuyerProductDetail } from "@/features/catalog/api";
import { listMyAddresses } from "@/features/address/api";
import { UserAddress } from "@/features/address/types";
import { buildIdempotencyKey, createOrderBuyNow, createOrderFromCart } from "@/features/order/api";
import { buildCheckoutPointsSummary, POINTS_PER_YUAN } from "@/features/order/checkout-points";
import { getMyOverview } from "@/features/overview/api";
import { getCartItemHints } from "@/lib/cart-hints";
import { getDemoSalePriceCents } from "@/lib/demo-pricing";
import { formatCnyFromCents } from "@/lib/price";

type DisplayItem = CheckoutSnapshotItem & {
  settlePrice: number;
  spuTitle: string;
  skuName: string;
};

function toPositiveInt(raw: string | null, fallback: number): number {
  const value = Number(raw ?? "");
  if (!Number.isFinite(value) || value <= 0) {
    return fallback;
  }
  return Math.trunc(value);
}

function getDefaultAddressId(fromQuery: string | null): number {
  const queryAddressId = toPositiveInt(fromQuery, 0);
  if (queryAddressId > 0) {
    return queryAddressId;
  }
  if (typeof window === "undefined") {
    return 0;
  }
  return toPositiveInt(window.localStorage.getItem("shopa_mall_default_address_id"), 0);
}

function buildAddressText(address?: UserAddress | null): string {
  if (!address) {
    return "";
  }
  return [address.provinceName, address.cityName, address.districtName, address.street, address.detail]
    .filter((item) => item && item.trim().length > 0)
    .join(" ");
}

export default function CheckoutConfirmPage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const [messageApi, contextHolder] = message.useMessage();
  const router = useRouter();
  const searchParams = useSearchParams();

  const mode = searchParams.get("mode") === "buy-now" ? "buy-now" : "cart";
  const skuNo = searchParams.get("sku_no") ?? "";
  const spuNo = searchParams.get("spu_no") ?? "";
  const shopNo = searchParams.get("shop_no") ?? "";
  const qty = toPositiveInt(searchParams.get("qty"), 1);

  const [buyerRemark, setBuyerRemark] = useState("");
  const [selectedAddressId, setSelectedAddressId] = useState<number>(() => getDefaultAddressId(searchParams.get("address_id")));
  const [usePoints, setUsePoints] = useState(false);

  const addressQuery = useQuery({
    queryKey: ["checkout-addresses"],
    queryFn: () => listMyAddresses(false),
    staleTime: 20_000,
    refetchOnWindowFocus: false
  });

  const addresses = useMemo(() => addressQuery.data?.addresses ?? [], [addressQuery.data?.addresses]);

  const overviewQuery = useQuery({
    queryKey: ["mall-me-overview", "checkout-confirm"],
    queryFn: getMyOverview,
    staleTime: 20_000,
    refetchOnWindowFocus: false
  });

  useEffect(() => {
    if (addresses.length === 0) {
      return;
    }
    const hasSelected = selectedAddressId > 0 && addresses.some((item) => item.addressId === selectedAddressId);
    if (hasSelected) {
      return;
    }
    const next = addresses.find((item) => item.isDefault)?.addressId ?? addresses[0].addressId;
    setSelectedAddressId(next);
  }, [addresses, selectedAddressId]);

  useEffect(() => {
    if (selectedAddressId > 0 && typeof window !== "undefined") {
      window.localStorage.setItem("shopa_mall_default_address_id", String(selectedAddressId));
    }
  }, [selectedAddressId]);

  const checkoutQuery = useQuery({
    queryKey: ["checkout-snapshot", mode, selectedAddressId],
    queryFn: () =>
      prepareCheckout({
        scope: 1,
        addressId: selectedAddressId
      }),
    enabled: mode === "cart" && selectedAddressId > 0,
    staleTime: 15_000,
    refetchOnWindowFocus: false
  });

  const snapshot: CheckoutSnapshot | null = checkoutQuery.data ?? null;

  const buyNowDetailQuery = useQuery({
    queryKey: ["checkout-buy-now-detail", mode, spuNo],
    queryFn: () => getBuyerProductDetail(spuNo),
    enabled: mode === "buy-now" && Boolean(spuNo),
    staleTime: 60_000,
    refetchOnWindowFocus: false
  });

  const displayItems = useMemo<DisplayItem[]>(() => {
    const hints = getCartItemHints();

    if (mode === "buy-now") {
      const detail = buyNowDetailQuery.data;
      const sku = detail?.skus?.find((item) => item.skuNo === skuNo) ?? detail?.skus?.[0];
      if (!sku) {
        return [];
      }
      const fallbackPrice = getDemoSalePriceCents(sku.spuNo || detail?.spuNo || spuNo, sku.skuNo);
      return [
        {
          skuNo: sku.skuNo,
          spuNo: sku.spuNo || detail?.spuNo || spuNo,
          shopNo: sku.shopNo || detail?.shopNo || shopNo,
          qty,
          settlePrice: sku.salePrice > 0 ? sku.salePrice : fallbackPrice,
          marketPrice: sku.marketPrice,
          spuTitle: detail?.title || hints[sku.skuNo]?.spuTitle || sku.skuNo,
          skuName: sku.skuName || hints[sku.skuNo]?.skuName || sku.skuNo,
          skuImageAssetId: sku.skuImageAssetId,
          saleAttrsJson: sku.saleAttrsJson
        }
      ];
    }

    const source = snapshot?.items ?? [];
    return source.map((item) => {
      const hint = hints[item.skuNo];
      const fallbackPrice = getDemoSalePriceCents(item.spuNo, item.skuNo);
      return {
        ...item,
        settlePrice: item.settlePrice > 0 ? item.settlePrice : hint?.salePrice || fallbackPrice,
        spuTitle: item.spuTitle || hint?.spuTitle || item.skuName || item.skuNo,
        skuName: item.skuName || hint?.skuName || item.skuNo
      };
    });
  }, [mode, buyNowDetailQuery.data, skuNo, spuNo, shopNo, qty, snapshot?.items]);

  const computedPayableAmount = useMemo(
    () => displayItems.reduce((sum, item) => sum + item.settlePrice * item.qty, 0),
    [displayItems]
  );

  const payableAmount = mode === "cart" ? snapshot?.payableAmount || computedPayableAmount : computedPayableAmount;
  const availablePoints = overviewQuery.data?.points ?? 0;
  const pointsSummary = useMemo(
    () =>
      buildCheckoutPointsSummary({
        payableAmount,
        availablePoints,
        usePoints
      }),
    [availablePoints, payableAmount, usePoints]
  );
  const maxPointsDiscountAmount = pointsSummary.maxDiscountAmount;
  const maxUsablePoints = pointsSummary.maxUsablePoints;
  const selectedPoints = pointsSummary.selectedPoints;
  const selectedPointsDiscountAmount = pointsSummary.selectedDiscountAmount;
  const finalPayableAmount = pointsSummary.finalPayableAmount;

  useEffect(() => {
    if (usePoints && maxUsablePoints === 0) {
      setUsePoints(false);
    }
  }, [maxUsablePoints, usePoints]);

  const selectedAddress = useMemo(
    () => addresses.find((item) => item.addressId === selectedAddressId) ?? null,
    [addresses, selectedAddressId]
  );

  const effectiveBuyNowSkuNo = skuNo || displayItems[0]?.skuNo || "";
  const effectiveBuyNowSpuNo = spuNo || buyNowDetailQuery.data?.spuNo || displayItems[0]?.spuNo || "";
  const effectiveBuyNowShopNo = shopNo || buyNowDetailQuery.data?.shopNo || displayItems[0]?.shopNo || "";

  const loading = mode === "buy-now" ? buyNowDetailQuery.isLoading : checkoutQuery.isLoading;
  const loadFailed = mode === "buy-now" ? buyNowDetailQuery.isError : checkoutQuery.isError;

  const canSubmit =
    selectedAddressId > 0 &&
    displayItems.length > 0 &&
    (mode === "buy-now"
      ? Boolean(effectiveBuyNowSkuNo && effectiveBuyNowSpuNo && effectiveBuyNowShopNo)
      : Boolean(snapshot?.checkoutToken));

  const createOrderMutation = useMutation({
    mutationFn: async () => {
      if (selectedAddressId <= 0) {
        throw new Error(isZh ? "请选择收货地址" : "Please select an address");
      }

      if (mode === "buy-now" && effectiveBuyNowSkuNo && effectiveBuyNowSpuNo && effectiveBuyNowShopNo) {
        return createOrderBuyNow({
          skuNo: effectiveBuyNowSkuNo,
          spuNo: effectiveBuyNowSpuNo,
          shopNo: effectiveBuyNowShopNo,
          qty,
          addressId: selectedAddressId,
          buyerRemark,
          submitSourceCode: "mall-web",
          usePoints,
          intentPoints: selectedPoints,
          expectedPointsCashAmount: selectedPointsDiscountAmount,
          idempotencyKey: buildIdempotencyKey("buy_now")
        });
      }

      if (!snapshot) {
        throw new Error("snapshot is not ready");
      }

      return createOrderFromCart({
        checkoutToken: snapshot.checkoutToken,
        addressId: selectedAddressId,
        buyerRemark,
        submitSourceCode: "mall-web",
        usePoints,
        intentPoints: selectedPoints,
        expectedPointsCashAmount: selectedPointsDiscountAmount,
        expectedSnapshotDigest: snapshot.snapshotDigest,
        idempotencyKey: buildIdempotencyKey("from_cart")
      });
    },
    onSuccess: (result) => {
      if (!result.orderNo) {
        messageApi.error(isZh ? "下单成功但缺少订单号，请到订单页刷新查看" : "Order created but order number is missing");
        router.push("/me/orders");
        return;
      }
      messageApi.success(isZh ? `订单已创建：${result.orderNo}` : `Order created: ${result.orderNo}`);
      router.push(`/checkout/pay?order_no=${encodeURIComponent(result.orderNo)}&address_id=${selectedAddressId}`);
    },
    onError: (err) => {
      console.error(err);
      messageApi.error(isZh ? "下单失败，请稍后重试" : "Failed to create order");
    }
  });

  return (
    <main className="tb-checkout-page">
      {contextHolder}
      <h1>{isZh ? "确认订单" : "Checkout Confirmation"}</h1>

      <section className="tb-checkout-address">
        <h2>{isZh ? "收货地址" : "Shipping Address"}</h2>

        {selectedAddress ? (
          <div className="tb-checkout-address-card">
            <p>
              {selectedAddress.receiverName || "--"} {selectedAddress.receiverPhone || "--"}
              {selectedAddress.isDefault ? <span>{isZh ? "默认" : "Default"}</span> : null}
            </p>
            <small>{buildAddressText(selectedAddress) || "--"}</small>
          </div>
        ) : (
          <div className="tb-checkout-address-card">
            <p>{isZh ? "暂无收货地址" : "No address available"}</p>
            <small>{isZh ? "请先新增地址后再结算。" : "Please create an address before checkout."}</small>
          </div>
        )}

        {addresses.length > 0 ? (
          <div className="tb-checkout-address-list">
            {addresses.map((address) => (
              <button
                key={address.addressId}
                type="button"
                className={`tb-checkout-address-item ${selectedAddressId === address.addressId ? "active" : ""}`}
                onClick={() => setSelectedAddressId(address.addressId)}
              >
                <strong>{address.receiverName}</strong>
                <span>{address.receiverPhone}</span>
                <small>{buildAddressText(address)}</small>
              </button>
            ))}
          </div>
        ) : (
          <p className="tb-checkout-address-empty">
            {isZh ? "还没有地址，去地址管理添加。" : "No address yet, go to address book to add one."}
            <Link href="/me/address">{isZh ? "地址管理" : "Address Book"}</Link>
          </p>
        )}
      </section>

      <section className="tb-checkout-items">
        <h2>{isZh ? "商品清单" : "Items"}</h2>
        <table>
          <thead>
            <tr>
              <th>{isZh ? "商品" : "Item"}</th>
              <th>{isZh ? "单价" : "Price"}</th>
              <th>{isZh ? "数量" : "Qty"}</th>
              <th>{isZh ? "小计" : "Subtotal"}</th>
            </tr>
          </thead>
          <tbody>
            {displayItems.map((item) => (
              <tr key={item.skuNo}>
                <td>{item.spuTitle || item.skuName || item.skuNo}</td>
                <td>{`CNY ${formatCnyFromCents(item.settlePrice)}`}</td>
                <td>{item.qty}</td>
                <td>{`CNY ${formatCnyFromCents(item.settlePrice * item.qty)}`}</td>
              </tr>
            ))}
          </tbody>
        </table>
        {displayItems.length === 0 && !loading && <p>{isZh ? "暂无可结算商品" : "No items to checkout."}</p>}
        {loading && <p>{isZh ? "正在加载结算信息..." : "Loading checkout data..."}</p>}
        {loadFailed && <p>{isZh ? "结算信息加载失败，请返回购物车重试。" : "Failed to load checkout data."}</p>}
      </section>

      <section className="tb-checkout-points">
        <div className="tb-checkout-points-head">
          <div>
            <h2>{isZh ? "积分抵扣" : "Points"}</h2>
            <small>{isZh ? "100 积分可抵 1 元，最多抵订单金额的 5%" : "Use points for up to 5% off this order."}</small>
          </div>
          <button
            type="button"
            className={`tb-checkout-points-toggle ${usePoints ? "is-active" : ""}`}
            aria-pressed={usePoints}
            disabled={!pointsSummary.canToggle}
            onClick={() => setUsePoints((current) => !current)}
          >
            <span>{isZh ? "使用积分抵扣" : "Apply Points"}</span>
            <strong>
              {pointsSummary.canToggle
                ? `${maxUsablePoints} ${isZh ? "积分" : "pts"}`
                : isZh
                  ? "暂无可用积分"
                  : "No Points Available"}
            </strong>
          </button>
        </div>
        <div className="tb-checkout-points-grid">
          <div>
            <span>{isZh ? "当前积分" : "Available"}</span>
            <strong>{availablePoints}</strong>
          </div>
          <div>
            <span>{isZh ? "本单最多可抵" : "Max Discount"}</span>
            <strong>{`CNY ${formatCnyFromCents(maxPointsDiscountAmount)}`}</strong>
          </div>
          <div>
            <span>{isZh ? "本次抵扣" : "Selected"}</span>
            <strong>{usePoints ? `${selectedPoints} / CNY ${formatCnyFromCents(selectedPointsDiscountAmount)}` : isZh ? "未使用" : "Not Used"}</strong>
          </div>
        </div>
        <p className="tb-checkout-points-tip">
          {isZh
            ? `当前按照 ${POINTS_PER_YUAN} 积分抵 1 元计算，确认下单后会锁定本次抵扣。`
            : `Points are redeemed at ${POINTS_PER_YUAN} points per CNY 1 and will be reserved after order submission.`}
        </p>
      </section>

      <section className="tb-checkout-submit">
        <div className="tb-checkout-submit-form">
          <p>
            {isZh ? "买家留言：" : "Buyer Remark:"}
            <input
              value={buyerRemark}
              onChange={(event) => setBuyerRemark(event.target.value)}
              placeholder={isZh ? "选填，给商家留言" : "Optional note to seller"}
              style={{ marginLeft: 8, width: 280, maxWidth: "70vw", height: 34, borderRadius: 6, border: "1px solid #ddd", padding: "0 10px" }}
            />
          </p>
          <div className="tb-checkout-price-stack">
            <p className="tb-checkout-submit-note">
              {isZh ? "订单金额：" : "Order Amount: "}
              <span>{`CNY ${formatCnyFromCents(payableAmount)}`}</span>
            </p>
            <p className="tb-checkout-submit-note">
              {isZh ? "积分抵扣：" : "Points Discount: "}
              <span>{usePoints ? `- CNY ${formatCnyFromCents(selectedPointsDiscountAmount)}` : isZh ? "未使用" : "Not Used"}</span>
            </p>
          </div>
        </div>
        <div className="tb-checkout-payable">
          {usePoints ? <small>{`CNY ${formatCnyFromCents(payableAmount)}`}</small> : null}
          <p>
            {isZh ? "应付：" : "Payable: "}
            <strong>{`CNY ${formatCnyFromCents(finalPayableAmount)}`}</strong>
          </p>
        </div>
        <button type="button" disabled={!canSubmit || createOrderMutation.isLoading} onClick={() => createOrderMutation.mutate()}>
          {createOrderMutation.isLoading ? (isZh ? "提交中..." : "Submitting...") : isZh ? "提交订单" : "Submit Order"}
        </button>
      </section>
    </main>
  );
}
