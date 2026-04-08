type BaseShopOption = {
  shopNo: string;
  shopName?: string;
  shopStatusCode?: string;
};

export type SellerChatShopOption = {
  shopNo: string;
  shopName: string;
  shopStatusCode: string;
};

function normalizeShopNo(value: unknown): string {
  return String(value ?? "").trim();
}

function normalizeShopOption(shop: BaseShopOption | null | undefined): SellerChatShopOption | null {
  const shopNo = normalizeShopNo(shop?.shopNo);
  if (!shopNo) {
    return null;
  }
  const shopName = String(shop?.shopName ?? "").trim() || shopNo;
  const shopStatusCode = String(shop?.shopStatusCode ?? "").trim();
  return {
    shopNo,
    shopName,
    shopStatusCode
  };
}

export function buildSellerChatShopOptions(input: {
  workbenchShops?: BaseShopOption[];
  catalogShopNos?: string[];
  persistedShopNo?: string;
}): SellerChatShopOption[] {
  const ordered: SellerChatShopOption[] = [];
  const seen = new Set<string>();

  const pushShop = (shop: BaseShopOption | null | undefined) => {
    const normalized = normalizeShopOption(shop);
    if (!normalized || seen.has(normalized.shopNo)) {
      return;
    }
    seen.add(normalized.shopNo);
    ordered.push(normalized);
  };

  (input.workbenchShops ?? []).forEach((shop) => pushShop(shop));
  (input.catalogShopNos ?? []).forEach((shopNo) => pushShop({ shopNo }));
  pushShop({ shopNo: input.persistedShopNo ?? "" });

  return ordered;
}
