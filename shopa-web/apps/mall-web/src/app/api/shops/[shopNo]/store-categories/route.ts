import { NextRequest, NextResponse } from "next/server";

type GoFrameResponse<T> = {
  code?: number;
  message?: string;
  data?: T;
};

type RawStoreCategory = {
  id?: number | string;
  shop_no?: string;
  shopNo?: string;
  parent_id?: number | string;
  parentId?: number | string;
  name?: string;
  level?: number | string;
  sort_order?: number | string;
  sortOrder?: number | string;
  is_visible?: boolean;
  isVisible?: boolean;
  product_count?: number | string;
  productCount?: number | string;
  children?: RawStoreCategory[];
};

type RawStoreCategoryList = {
  categories?: RawStoreCategory[];
};

type ShopStoreCategoryResponse = {
  id: number;
  shopNo: string;
  parentId: number;
  name: string;
  level: number;
  sortOrder: number;
  isVisible: boolean;
  productCount: number;
  children: ShopStoreCategoryResponse[];
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

function normalizeCategory(raw: RawStoreCategory): ShopStoreCategoryResponse {
  return {
    id: toNumber(raw.id),
    shopNo: toString(raw.shopNo ?? raw.shop_no),
    parentId: toNumber(raw.parentId ?? raw.parent_id),
    name: toString(raw.name),
    level: toNumber(raw.level),
    sortOrder: toNumber(raw.sortOrder ?? raw.sort_order),
    isVisible: Boolean(raw.isVisible ?? raw.is_visible),
    productCount: toNumber(raw.productCount ?? raw.product_count),
    children: Array.isArray(raw.children) ? raw.children.map((child) => normalizeCategory(child)) : []
  };
}

export async function GET(_: NextRequest, context: { params: { shopNo: string } }) {
  const shopNo = context.params.shopNo?.trim() || "";
  if (!shopNo) {
    return NextResponse.json({ message: "shopNo is required" }, { status: 400 });
  }

  const baseUrl = process.env.SELLER_SHOP_SERVICE_BASE_URL || "http://127.0.0.1:8007";
  const response = await fetch(`${baseUrl}/v1/shops/${encodeURIComponent(shopNo)}/store-categories`, {
    method: "GET",
    cache: "no-store"
  });
  if (!response.ok) {
    return NextResponse.json({ message: `store categories request failed: ${response.status}` }, { status: response.status });
  }

  const payload = unwrapGoFrame<RawStoreCategoryList>(await response.json());
  const categories = Array.isArray(payload?.categories) ? payload.categories.map((item) => normalizeCategory(item)) : [];
  return NextResponse.json(categories);
}
