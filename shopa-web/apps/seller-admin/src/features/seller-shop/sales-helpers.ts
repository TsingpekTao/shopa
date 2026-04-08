import { SellerSalesAnalyticsParams, SellerSalesBucket } from "./types";

export type SellerSalesMetricType = "GMV" | "ORDERS";

export function buildSellerSalesAnalyticsPath(params: SellerSalesAnalyticsParams): string {
  const query = new URLSearchParams();
  query.set("range", params.range);
  const shopNo = params.shopNo?.trim();
  if (shopNo) {
    query.set("shopNo", shopNo);
  }
  return `/v1/seller/sales/analytics?${query.toString()}`;
}

export function getSellerSalesMetricValue(
  bucket: Pick<SellerSalesBucket, "gmv" | "paidOrderCount">,
  metric: SellerSalesMetricType
): number {
  return metric === "ORDERS" ? Number(bucket.paidOrderCount ?? 0) : Number(bucket.gmv ?? 0);
}

export function formatSellerSalesRefundRate(value?: number | null): string {
  const rate = Number(value ?? 0);
  if (!Number.isFinite(rate) || rate <= 0) {
    return "0.00%";
  }
  return `${(rate * 100).toFixed(2)}%`;
}

export function getSellerStatDisplayValue(value: number | string | null | undefined, degraded = false): number | string {
  if (degraded) {
    return "--";
  }
  if (typeof value === "string") {
    return value;
  }
  return Number.isFinite(Number(value)) ? Number(value) : 0;
}

export function hasSellerDegradedField(fields: string[] | undefined, field: string): boolean {
  return Array.isArray(fields) && fields.includes(field);
}
