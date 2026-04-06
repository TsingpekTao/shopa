import { NextRequest, NextResponse } from "next/server";

type GoFrameResponse<T> = {
  code?: number;
  message?: string;
  data?: T;
};

type RawShop = {
  shop_no?: string;
  shopNo?: string;
  shop_name?: string;
  shopName?: string;
  shop_display_name?: string;
  shopDisplayName?: string;
  status?: number | string;
  buyer_visible?: boolean;
  buyerVisible?: boolean;
};

type RawGetShopByNoRes = {
  shop?: RawShop;
};

function toNumber(value: unknown, fallback = 0): number {
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : fallback;
}

function toString(value: unknown): string {
  if (value === null || value === undefined) {
    return "";
  }
  return String(value);
}

function unwrapGoFrame<T>(payload: unknown): T {
  if (!payload || typeof payload !== "object") {
    return payload as T;
  }

  const response = payload as GoFrameResponse<T>;
  if (typeof response.code === "number" && "data" in response) {
    if (response.code !== 0) {
      throw new Error(response.message || "upstream request failed");
    }
    return response.data as T;
  }
  return payload as T;
}

export async function GET(_: NextRequest, context: { params: { shopNo: string } }) {
  const shopNo = context.params.shopNo?.trim() || "";
  if (!shopNo) {
    return NextResponse.json({ message: "shopNo is required" }, { status: 400 });
  }

  const baseUrl = process.env.SELLER_SHOP_SERVICE_BASE_URL || "http://127.0.0.1:8007";
  const response = await fetch(`${baseUrl}/v1/internal/shops/${encodeURIComponent(shopNo)}`, {
    method: "GET",
    cache: "no-store"
  });
  if (!response.ok) {
    return NextResponse.json({ message: `shop request failed: ${response.status}` }, { status: response.status });
  }

  const payload = unwrapGoFrame<RawGetShopByNoRes>(await response.json());
  const shop = payload?.shop;
  if (!shop) {
    return NextResponse.json({ message: "shop not found" }, { status: 404 });
  }

  return NextResponse.json({
    shopNo: toString(shop.shopNo ?? shop.shop_no),
    shopName: toString(shop.shopName ?? shop.shop_name),
    shopDisplayName: toString(shop.shopDisplayName ?? shop.shop_display_name),
    status: toNumber(shop.status),
    buyerVisible: Boolean(shop.buyerVisible ?? shop.buyer_visible)
  });
}
