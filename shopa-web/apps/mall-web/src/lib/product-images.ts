"use client";

function parseImagePool(raw: string | undefined): string[] {
  if (!raw) {
    return [];
  }
  return raw
    .split(",")
    .map((item) => item.trim())
    .filter((item) => item.length > 0);
}

const PRODUCT_IMAGE_POOL = parseImagePool(process.env.NEXT_PUBLIC_MALL_OSS_PRODUCT_IMAGES);

function asciiHash(text: string): number {
  let sum = 0;
  for (let i = 0; i < text.length; i += 1) {
    sum += text.charCodeAt(i);
  }
  return sum;
}

export function pickProductImageBySpuNo(spuNo: string): string {
  if (!spuNo || PRODUCT_IMAGE_POOL.length === 0) {
    return "";
  }
  const index = asciiHash(spuNo) % PRODUCT_IMAGE_POOL.length;
  return PRODUCT_IMAGE_POOL[index] ?? "";
}

export function hasProductImagePool(): boolean {
  return PRODUCT_IMAGE_POOL.length > 0;
}

