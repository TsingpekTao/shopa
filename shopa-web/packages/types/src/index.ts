export type ApiResponse<T> = {
  data: T;
  partial?: boolean;
  degradedFields?: string[];
};

export type SellerShopSummary = {
  shopNo: string;
  shopName: string;
  shopDisplayName: string;
  shopStatusCode: string;
  buyerVisible: boolean;
  updatedAt: string;
};

export type SellerProductSummary = {
  total: number;
  onShelf: number;
  offShelf: number;
  reviewing: number;
  draft: number;
  rejected: number;
};

export type SellerApplicationSummary = {
  applicationNo: string;
  applicationStatusCode: string;
  version: number;
  submittedAt: string;
  updatedAt: string;
};
