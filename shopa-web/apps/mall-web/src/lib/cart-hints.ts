export type CartItemHint = {
  skuNo: string;
  spuNo: string;
  shopNo: string;
  spuTitle: string;
  subTitle: string;
  skuName: string;
  imageUrl: string;
  salePrice: number;
  marketPrice: number;
  updatedAt: number;
};

const CART_HINT_KEY = "shopa_mall_cart_item_hints_v1";
const CART_HINT_MAX = 500;

function canUseStorage(): boolean {
  return typeof window !== "undefined" && Boolean(window.localStorage);
}

function parseHints(raw: string | null): Record<string, CartItemHint> {
  if (!raw) {
    return {};
  }
  try {
    const parsed = JSON.parse(raw) as Record<string, CartItemHint>;
    if (!parsed || typeof parsed !== "object") {
      return {};
    }
    return parsed;
  } catch {
    return {};
  }
}

function persistHints(hints: Record<string, CartItemHint>): void {
  if (!canUseStorage()) {
    return;
  }
  try {
    window.localStorage.setItem(CART_HINT_KEY, JSON.stringify(hints));
  } catch {
    // ignore quota errors
  }
}

function trimHints(hints: Record<string, CartItemHint>): Record<string, CartItemHint> {
  const entries = Object.entries(hints);
  if (entries.length <= CART_HINT_MAX) {
    return hints;
  }
  entries.sort((a, b) => (b[1]?.updatedAt ?? 0) - (a[1]?.updatedAt ?? 0));
  return Object.fromEntries(entries.slice(0, CART_HINT_MAX));
}

export function getCartItemHints(): Record<string, CartItemHint> {
  if (!canUseStorage()) {
    return {};
  }
  return parseHints(window.localStorage.getItem(CART_HINT_KEY));
}

export function upsertCartItemHint(hint: CartItemHint): void {
  const skuNo = hint.skuNo.trim();
  if (!skuNo || !canUseStorage()) {
    return;
  }
  const hints = getCartItemHints();
  hints[skuNo] = { ...hint, skuNo, updatedAt: Date.now() };
  persistHints(trimHints(hints));
}

