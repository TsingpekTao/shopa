import { apiClient } from "../../lib/api-client";
import { SellerSalesAnalyticsParams, SellerSalesAnalyticsResponse } from "./types";
import { buildSellerSalesAnalyticsPath } from "./sales-helpers";

export * from "./sales-helpers";

export async function fetchSellerSalesAnalytics(params: SellerSalesAnalyticsParams): Promise<SellerSalesAnalyticsResponse> {
  return apiClient.get<SellerSalesAnalyticsResponse>(buildSellerSalesAnalyticsPath(params));
}
