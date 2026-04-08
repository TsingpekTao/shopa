export function buildSellerOrderListPath(): string {
  return "/v1/order/seller/list";
}

export function buildSellerOrderDetailPath(subOrderNo: string): string {
  return `/v1/order/seller/sub/${encodeURIComponent(subOrderNo.trim())}`;
}

export function buildSellerSubOrderShipPath(): string {
  return "/v1/order/seller/sub/ship";
}
