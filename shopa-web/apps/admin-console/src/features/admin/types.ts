export interface DashboardMetric {
  key: string;
  label: string;
  value: number;
}

export interface DashboardOverviewResponse {
  metrics: DashboardMetric[];
  partial?: boolean;
  degradedFields?: string[];
}

export interface QualificationDoc {
  docTypeCode?: string;
  assetId?: string;
  validFrom?: string;
  validUntil?: string;
  issuer?: string;
  ext?: Record<string, string>;
}

export interface LegalSubjectProfile {
  subjectTypeCode?: string;
  subjectName?: string;
  subjectCertNo?: string;
  legalRepresentativeName?: string;
  legalRepresentativeCertNo?: string;
  certValidFrom?: string;
  certValidUntil?: string;
  docs?: QualificationDoc[];
  ext?: Record<string, string>;
}

export interface MerchantEntityProfile {
  entityNo?: string;
  ownerUserId?: string;
  merchantTypeCode?: string;
  entityName?: string;
  contactName?: string;
  contactPhone?: string;
  contactEmail?: string;
  legalSubject?: LegalSubjectProfile;
  ext?: Record<string, string>;
}

export interface MerchantShopProfile {
  shopNo?: string;
  shopName?: string;
  shopDisplayName?: string;
  shopTypeCode?: string;
  mainCategoryIds?: string[];
  logoAssetId?: string;
  bannerAssetId?: string;
  servicePhone?: string;
  serviceEmail?: string;
  description?: string;
  ext?: Record<string, string>;
}

export interface RejectInfo {
  rejectReasonCode?: string;
  rejectComment?: string;
  rejectedAt?: string;
}

export interface MerchantApplication {
  applicationNo: string;
  previousApplicationNo?: string;
  nextApplicationNo?: string;
  ownerUserId?: string;
  status?: number;
  version?: number;
  entityDraft?: MerchantEntityProfile;
  shopDraft?: MerchantShopProfile;
  entitySubmitted?: MerchantEntityProfile;
  shopSubmitted?: MerchantShopProfile;
  entity?: MerchantEntityProfile;
  shop?: MerchantShopProfile;
  latestReject?: RejectInfo;
  submittedAt?: string;
  reviewStartedAt?: string;
  reviewedAt?: string;
  reviewerId?: string;
  reviewComment?: string;
  createdAt?: string;
  updatedAt?: string;
}

export interface MerchantApplicationListResponse {
  applications: MerchantApplication[];
  page: number;
  pageSize: number;
  total: number;
}

export interface ProductReviewTask {
  taskNo?: string;
  spuNo: string;
  shopNo?: string;
  title?: string;
  spuStatus?: number;
  submittedAt?: string;
}

export interface ProductReviewDetailSpu {
  spuNo: string;
  shopNo?: string;
  title?: string;
  subTitle?: string;
  categoryId?: number;
  brandNo?: string;
  mainImageAssetIds: string[];
  detailImageAssetIds: string[];
  attributeValues?: Record<string, string>;
  spuStatus?: number;
  spuStockStatus?: number;
  minSalePrice?: number;
  maxSalePrice?: number;
  version?: number;
}

export interface ProductReviewDetailSku {
  skuNo: string;
  spuNo?: string;
  skuName?: string;
  skuImageAssetId?: string;
  salePrice?: number;
  marketPrice?: number;
  saleAttrs?: Record<string, string>;
  stockStatus?: number;
}

export interface ProductInventorySnapshot {
  skuNo: string;
  totalQty?: number;
  availableQty?: number;
  stockStatus?: number;
}

export interface ProductReviewDetail {
  product: {
    spu: ProductReviewDetailSpu | null;
    skus: ProductReviewDetailSku[];
  };
  review?: {
    reviewStatus?: number;
    rejectReasonCode?: string;
    rejectComment?: string;
    reviewerId?: string;
    reviewedAt?: string;
  };
}

export interface ProductReviewTaskListResponse {
  tasks: ProductReviewTask[];
  page: number;
  pageSize: number;
  total: number;
}

export interface ConversationItem {
  conversationNo: string;
  title: string;
  lastMessageAt: string;
  statusCode: string;
}

export interface ConversationListResponse {
  items: ConversationItem[];
  partial?: boolean;
  degradedFields?: string[];
}

export interface ShopInsightsResponse {
  shopNo: string;
  shopName: string;
  shopStatusCode: string;
  partial?: boolean;
  degradedFields?: string[];
}
