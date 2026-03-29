import { apiClient } from "@/lib/api-client";
import {
  ConversationListResponse,
  DashboardOverviewResponse,
  LegalSubjectProfile,
  MerchantApplication,
  MerchantApplicationListResponse,
  MerchantEntityProfile,
  MerchantShopProfile,
  ProductReviewTaskListResponse,
  QualificationDoc,
  RejectInfo,
  ShopInsightsResponse
} from "./types";

function isRecord(value: unknown): value is Record<string, unknown> {
  return Boolean(value) && typeof value === "object" && !Array.isArray(value);
}

function toStringValue(value: unknown): string | undefined {
  if (value === null || value === undefined) {
    return undefined;
  }
  const normalized = String(value).trim();
  return normalized.length > 0 ? normalized : undefined;
}

function toStringMap(value: unknown): Record<string, string> | undefined {
  if (!isRecord(value)) {
    return undefined;
  }
  const mappedEntries = Object.entries(value)
    .map(([key, item]) => [key, toStringValue(item)] as const)
    .filter(([, item]) => Boolean(item)) as [string, string][];
  if (mappedEntries.length === 0) {
    return undefined;
  }
  return Object.fromEntries(mappedEntries);
}

function toTimestampValue(value: unknown): string | undefined {
  if (value === null || value === undefined) {
    return undefined;
  }
  if (typeof value === "string") {
    const normalized = value.trim();
    return normalized.length > 0 ? normalized : undefined;
  }
  if (typeof value === "number") {
    if (!Number.isFinite(value) || value <= 0) {
      return undefined;
    }
    return new Date(value * 1000).toISOString();
  }
  if (!isRecord(value)) {
    return undefined;
  }
  const seconds = Number(value.seconds ?? 0);
  const nanos = Number(value.nanos ?? 0);
  if (!Number.isFinite(seconds) || seconds <= 0) {
    return undefined;
  }
  const millis = seconds * 1000 + Math.floor((Number.isFinite(nanos) ? nanos : 0) / 1_000_000);
  return new Date(millis).toISOString();
}

function toNumberValue(value: unknown, fallback = 0): number {
  const normalized = Number(value);
  return Number.isFinite(normalized) ? normalized : fallback;
}

function hasAnyValue(values: unknown[]): boolean {
  return values.some((value) => {
    if (value === null || value === undefined) {
      return false;
    }
    if (typeof value === "string") {
      return value.trim().length > 0;
    }
    if (Array.isArray(value)) {
      return value.length > 0;
    }
    if (isRecord(value)) {
      return Object.keys(value).length > 0;
    }
    return true;
  });
}

function normalizeQualificationDoc(raw: unknown): QualificationDoc | undefined {
  if (!isRecord(raw)) {
    return undefined;
  }
  const doc: QualificationDoc = {
    docTypeCode: toStringValue(raw.docTypeCode ?? raw.doc_type_code),
    assetId: toStringValue(raw.assetId ?? raw.asset_id),
    validFrom: toTimestampValue(raw.validFrom ?? raw.valid_from),
    validUntil: toTimestampValue(raw.validUntil ?? raw.valid_until),
    issuer: toStringValue(raw.issuer),
    ext: toStringMap(raw.ext)
  };
  return hasAnyValue(Object.values(doc)) ? doc : undefined;
}

function normalizeLegalSubjectProfile(raw: unknown): LegalSubjectProfile | undefined {
  if (!isRecord(raw)) {
    return undefined;
  }
  const docsRaw = raw.docs;
  const docs = Array.isArray(docsRaw) ? docsRaw.map(normalizeQualificationDoc).filter(Boolean) : undefined;
  const legalSubject: LegalSubjectProfile = {
    subjectTypeCode: toStringValue(raw.subjectTypeCode ?? raw.subject_type_code),
    subjectName: toStringValue(raw.subjectName ?? raw.subject_name),
    subjectCertNo: toStringValue(raw.subjectCertNo ?? raw.subject_cert_no),
    legalRepresentativeName: toStringValue(raw.legalRepresentativeName ?? raw.legal_representative_name),
    legalRepresentativeCertNo: toStringValue(raw.legalRepresentativeCertNo ?? raw.legal_representative_cert_no),
    certValidFrom: toTimestampValue(raw.certValidFrom ?? raw.cert_valid_from),
    certValidUntil: toTimestampValue(raw.certValidUntil ?? raw.cert_valid_until),
    docs: docs as QualificationDoc[] | undefined,
    ext: toStringMap(raw.ext)
  };
  return hasAnyValue(Object.values(legalSubject)) ? legalSubject : undefined;
}

function normalizeEntityProfile(raw: unknown): MerchantEntityProfile | undefined {
  if (!isRecord(raw)) {
    return undefined;
  }
  const entity: MerchantEntityProfile = {
    entityNo: toStringValue(raw.entityNo ?? raw.entity_no),
    ownerUserId: toStringValue(raw.ownerUserId ?? raw.owner_user_id),
    merchantTypeCode: toStringValue(raw.merchantTypeCode ?? raw.merchant_type_code),
    entityName: toStringValue(raw.entityName ?? raw.entity_name),
    contactName: toStringValue(raw.contactName ?? raw.contact_name),
    contactPhone: toStringValue(raw.contactPhone ?? raw.contact_phone),
    contactEmail: toStringValue(raw.contactEmail ?? raw.contact_email),
    legalSubject: normalizeLegalSubjectProfile(raw.legalSubject ?? raw.legal_subject),
    ext: toStringMap(raw.ext)
  };
  return hasAnyValue(Object.values(entity)) ? entity : undefined;
}

function normalizeShopProfile(raw: unknown): MerchantShopProfile | undefined {
  if (!isRecord(raw)) {
    return undefined;
  }
  const categoryIdsRaw = raw.mainCategoryIds ?? raw.main_category_ids;
  const categoryIds = Array.isArray(categoryIdsRaw)
    ? categoryIdsRaw.map((item) => toStringValue(item)).filter(Boolean)
    : undefined;
  const shop: MerchantShopProfile = {
    shopNo: toStringValue(raw.shopNo ?? raw.shop_no),
    shopName: toStringValue(raw.shopName ?? raw.shop_name),
    shopDisplayName: toStringValue(raw.shopDisplayName ?? raw.shop_display_name),
    shopTypeCode: toStringValue(raw.shopTypeCode ?? raw.shop_type_code),
    mainCategoryIds: categoryIds as string[] | undefined,
    logoAssetId: toStringValue(raw.logoAssetId ?? raw.logo_asset_id),
    bannerAssetId: toStringValue(raw.bannerAssetId ?? raw.banner_asset_id),
    servicePhone: toStringValue(raw.servicePhone ?? raw.service_phone),
    serviceEmail: toStringValue(raw.serviceEmail ?? raw.service_email),
    description: toStringValue(raw.description),
    ext: toStringMap(raw.ext)
  };
  return hasAnyValue(Object.values(shop)) ? shop : undefined;
}

function normalizeRejectInfo(raw: unknown): RejectInfo | undefined {
  if (!isRecord(raw)) {
    return undefined;
  }
  const rejectInfo: RejectInfo = {
    rejectReasonCode: toStringValue(raw.rejectReasonCode ?? raw.reject_reason_code),
    rejectComment: toStringValue(raw.rejectComment ?? raw.reject_comment),
    rejectedAt: toTimestampValue(raw.rejectedAt ?? raw.rejected_at)
  };
  return hasAnyValue(Object.values(rejectInfo)) ? rejectInfo : undefined;
}

function normalizeMerchantApplication(raw: unknown): MerchantApplication {
  const record = isRecord(raw) ? raw : {};
  const entityDraft = normalizeEntityProfile(record.entityDraft ?? record.entity_draft);
  const entitySubmitted = normalizeEntityProfile(record.entitySubmitted ?? record.entity_submitted);
  const shopDraft = normalizeShopProfile(record.shopDraft ?? record.shop_draft);
  const shopSubmitted = normalizeShopProfile(record.shopSubmitted ?? record.shop_submitted);
  return {
    applicationNo: toStringValue(record.applicationNo ?? record.application_no) ?? "",
    previousApplicationNo: toStringValue(record.previousApplicationNo ?? record.previous_application_no),
    nextApplicationNo: toStringValue(record.nextApplicationNo ?? record.next_application_no),
    ownerUserId: toStringValue(record.ownerUserId ?? record.owner_user_id),
    status: toNumberValue(record.status, 0),
    version: toNumberValue(record.version, 0),
    entityDraft,
    shopDraft,
    entitySubmitted,
    shopSubmitted,
    entity: entitySubmitted ?? entityDraft,
    shop: shopSubmitted ?? shopDraft,
    latestReject: normalizeRejectInfo(record.latestReject ?? record.latest_reject),
    submittedAt: toTimestampValue(record.submittedAt ?? record.submitted_at),
    reviewStartedAt: toTimestampValue(record.reviewStartedAt ?? record.review_started_at),
    reviewedAt: toTimestampValue(record.reviewedAt ?? record.reviewed_at),
    reviewerId: toStringValue(record.reviewerId ?? record.reviewer_id),
    reviewComment: toStringValue(record.reviewComment ?? record.review_comment),
    createdAt: toTimestampValue(record.createdAt ?? record.created_at),
    updatedAt: toTimestampValue(record.updatedAt ?? record.updated_at)
  };
}

export async function fetchDashboardOverview(): Promise<DashboardOverviewResponse> {
  return apiClient.get<DashboardOverviewResponse>("/v1/admin/dashboard/overview");
}

export async function fetchMerchantApplications(params?: {
  page?: number;
  pageSize?: number;
  keyword?: string;
  statuses?: number[];
}) {
  const page = params?.page ?? 1;
  const pageSize = params?.pageSize ?? 20;
  const query = new URLSearchParams();
  query.set("page", String(page));
  query.set("pageSize", String(pageSize));
  if (params?.keyword?.trim()) {
    query.set("keyword", params.keyword.trim());
  }
  if (params?.statuses?.length) {
    params.statuses.forEach((status) => {
      query.append("statuses", String(status));
    });
  }
  const response = await apiClient.get<MerchantApplicationListResponse>(`/v1/admin/seller/applications?${query.toString()}`);
  return {
    applications: (response.applications ?? []).map((item) => normalizeMerchantApplication(item)).filter((item) => item.applicationNo),
    page: Number(response.page ?? page),
    pageSize: Number(response.pageSize ?? pageSize),
    total: Number(response.total ?? 0)
  };
}

export async function fetchMerchantApplicationDetail(applicationNo: string) {
  const response = await apiClient.get<{ application: MerchantApplication }>(`/v1/admin/seller/applications/${applicationNo}`);
  return {
    application: normalizeMerchantApplication(response.application)
  };
}

type AssetReadUrlResponse = {
  assetId?: number | string;
  asset_id?: number | string;
  url?: string;
  isPublic?: boolean;
  is_public?: boolean;
  expiredAt?: string;
  expired_at?: string;
};

export async function fetchAssetReadUrl(assetId: number | string, ttlSeconds = 3600) {
  const response = await apiClient.get<AssetReadUrlResponse>(`/v1/media/assets/${assetId}/read-url?ttlSeconds=${ttlSeconds}`);
  return {
    assetId: String(response.assetId ?? response.asset_id ?? assetId),
    url: response.url ?? "",
    isPublic: Boolean(response.isPublic ?? response.is_public),
    expiredAt: response.expiredAt ?? response.expired_at ?? ""
  };
}

export async function approveMerchantApplication(payload: { applicationNo: string; expectedVersion: number; reviewComment?: string }) {
  return apiClient.post(`/v1/admin/seller/applications/${payload.applicationNo}/approve`, {
    expectedVersion: payload.expectedVersion,
    reviewComment: payload.reviewComment || ""
  });
}

export async function rejectMerchantApplication(payload: {
  applicationNo: string;
  expectedVersion: number;
  rejectReasonCode: string;
  rejectComment?: string;
}) {
  return apiClient.post(`/v1/admin/seller/applications/${payload.applicationNo}/reject`, {
    expectedVersion: payload.expectedVersion,
    rejectReasonCode: payload.rejectReasonCode,
    rejectComment: payload.rejectComment || ""
  });
}

export async function fetchProductReviewTasks(params?: { page?: number; pageSize?: number }) {
  const page = params?.page ?? 1;
  const pageSize = params?.pageSize ?? 20;
  return apiClient.get<ProductReviewTaskListResponse>(`/v1/catalog/admin/review/tasks?page=${page}&pageSize=${pageSize}`);
}

export async function fetchProductReviewDetail(spuNo: string) {
  return apiClient.get<any>(`/v1/catalog/admin/review/${spuNo}`);
}

export async function approveProduct(payload: { spuNo: string; expectedVersion: number; reviewComment?: string }) {
  return apiClient.post("/v1/catalog/admin/review/approve", payload);
}

export async function rejectProduct(payload: { spuNo: string; expectedVersion: number; rejectReasonCode: string; rejectComment?: string }) {
  return apiClient.post("/v1/catalog/admin/review/reject", payload);
}

export async function freezeProduct(payload: { spuNo: string; expectedVersion: number; reasonCode: string; reason?: string }) {
  return apiClient.post("/v1/catalog/admin/review/freeze", payload);
}

export async function unfreezeProduct(payload: { spuNo: string; expectedVersion: number; reasonCode: string; reason?: string }) {
  return apiClient.post("/v1/catalog/admin/review/unfreeze", payload);
}

export async function forceOffShelf(payload: { spuNo: string; expectedVersion: number; reasonCode: string; reason?: string }) {
  return apiClient.post("/v1/catalog/admin/review/force-off-shelf", payload);
}

export async function fetchConversations(): Promise<ConversationListResponse> {
  return apiClient.get<ConversationListResponse>("/v1/admin/cs/conversations");
}

export async function fetchShopInsights(shopNo: string): Promise<ShopInsightsResponse> {
  return apiClient.get<ShopInsightsResponse>(`/v1/admin/shops/${shopNo}/insights`);
}
