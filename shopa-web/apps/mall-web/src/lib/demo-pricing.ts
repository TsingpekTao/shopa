function hashString(input: string): number {
  let hash = 0;
  for (let i = 0; i < input.length; i += 1) {
    hash = (hash * 31 + input.charCodeAt(i)) >>> 0;
  }
  return hash;
}

export function getDemoSalePriceCents(spuNo: string, skuNo?: string): number {
  const key = `${spuNo || "SPU"}#${skuNo || "SKU"}`;
  const hash = hashString(key);
  const yuan = 199 + (hash % 800); // 199 ~ 998
  return yuan * 100;
}

export function getDemoMarketPriceCents(salePriceCents: number, spuNo: string, skuNo?: string): number {
  const key = `${spuNo || "SPU"}#${skuNo || "SKU"}#MARKET`;
  const hash = hashString(key);
  const delta = (30 + (hash % 120)) * 100; // +30 ~ +149
  return salePriceCents + delta;
}

