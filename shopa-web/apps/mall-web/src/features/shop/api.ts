import { MallShopDetail, ShopStoreCategory } from "./types";

export async function getMallShopDetail(shopNo: string): Promise<MallShopDetail> {
  const normalizedShopNo = shopNo.trim();
  if (!normalizedShopNo) {
    throw new Error("shopNo is required");
  }

  const response = await fetch(`/api/shops/${encodeURIComponent(normalizedShopNo)}`, {
    method: "GET",
    cache: "no-store"
  });
  if (!response.ok) {
    throw new Error(`shop request failed: ${response.status}`);
  }

  return (await response.json()) as MallShopDetail;
}

export async function getShopStoreCategories(shopNo: string): Promise<ShopStoreCategory[]> {
  const normalizedShopNo = shopNo.trim();
  if (!normalizedShopNo) {
    return [];
  }

  const response = await fetch(`/api/shops/${encodeURIComponent(normalizedShopNo)}/store-categories`, {
    method: "GET",
    cache: "no-store"
  });
  if (!response.ok) {
    throw new Error(`store categories request failed: ${response.status}`);
  }

  return (await response.json()) as ShopStoreCategory[];
}
