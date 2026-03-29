import { apiClient } from "@/lib/api-client";
import {
  CreateApplicationDraftPayload,
  ListMyApplicationsPayload,
  ListMyApplicationsResponse,
  SellerApplicationListItem,
  SellerApplicationDetail,
  SellerWorkbenchResponse,
  ShopDashboardResponse,
  SubmitApplicationPayload,
  UpdateApplicationDraftPayload
} from "./types";

const mockWorkbench: SellerWorkbenchResponse = {
  userId: 0,
  shops: [],
  shopsTotal: 0,
  shopsTruncated: false,
  latestApplications: [],
  applicationsTotal: 0,
  applicationsTruncated: false,
  productSummary: { total: 0, onShelf: 0, offShelf: 0, reviewing: 0, draft: 0, rejected: 0 },
  partial: true,
  degradedFields: ["seller_shops"]
};

function pickString(...values: unknown[]): string {
  for (const value of values) {
    if (typeof value === "string" && value.trim()) {
      return value.trim();
    }
    if (typeof value === "number" && Number.isFinite(value)) {
      return String(value);
    }
  }
  return "";
}

function normalizeExt(raw: unknown): Record<string, string> {
  if (!raw || typeof raw !== "object") {
    return {};
  }
  const out: Record<string, string> = {};
  for (const [key, value] of Object.entries(raw as Record<string, unknown>)) {
    if (value === undefined || value === null) {
      continue;
    }
    out[String(key)] = String(value);
  }
  return out;
}

function hasNonEmptyField(raw: unknown): boolean {
  if (!raw || typeof raw !== "object") {
    return false;
  }
  return Object.values(raw as Record<string, unknown>).some((value) => {
    if (value === undefined || value === null) {
      return false;
    }
    if (typeof value === "string") {
      return value.trim().length > 0;
    }
    if (typeof value === "object") {
      return Object.keys(value as Record<string, unknown>).length > 0;
    }
    return true;
  });
}

function normalizeEntityProfile(raw: any) {
  if (!raw) {
    return undefined;
  }
  const legalSubject = raw?.legalSubject ?? raw?.legal_subject ?? {};
  const ext = normalizeExt(raw?.ext);
  return {
    merchantTypeCode: pickString(raw?.merchantTypeCode, raw?.merchant_type_code),
    entityName: pickString(raw?.entityName, raw?.entity_name, legalSubject?.subjectName, legalSubject?.subject_name),
    contactName: pickString(
      raw?.contactName,
      raw?.contact_name,
      legalSubject?.legalRepresentativeName,
      legalSubject?.legal_representative_name
    ),
    contactPhone: pickString(raw?.contactPhone, raw?.contact_phone),
    contactEmail: pickString(raw?.contactEmail, raw?.contact_email),
    ext
  };
}

function normalizeShopProfile(raw: any) {
  if (!raw) {
    return undefined;
  }
  return {
    shopName: pickString(raw?.shopName, raw?.shop_name),
    shopDisplayName: pickString(raw?.shopDisplayName, raw?.shop_display_name, raw?.shopName, raw?.shop_name),
    shopTypeCode: pickString(raw?.shopTypeCode, raw?.shop_type_code),
    servicePhone: pickString(raw?.servicePhone, raw?.service_phone),
    serviceEmail: pickString(raw?.serviceEmail, raw?.service_email),
    description: pickString(raw?.description),
    ext: normalizeExt(raw?.ext)
  };
}

function normalizeApplication(raw: any): SellerApplicationDetail {
  const statusCode = normalizeApplicationStatus(raw?.applicationStatusCode, raw?.status);
  const directEntity = raw?.entity;
  const draftEntity = raw?.entityDraft ?? raw?.entity_draft;
  const submittedEntity = raw?.entitySubmitted ?? raw?.entity_submitted;
  const directShop = raw?.shop;
  const draftShop = raw?.shopDraft ?? raw?.shop_draft;
  const submittedShop = raw?.shopSubmitted ?? raw?.shop_submitted;
  const preferSubmittedSnapshot = ["SUBMITTED", "REVIEWING", "APPROVED"].includes(statusCode);
  const entityRaw = directEntity || (preferSubmittedSnapshot ? submittedEntity || draftEntity : draftEntity || submittedEntity);
  const shopRaw = directShop || (preferSubmittedSnapshot ? submittedShop || draftShop : draftShop || submittedShop);
  const normalizedEntity = normalizeEntityProfile(entityRaw);
  const normalizedShop = normalizeShopProfile(shopRaw);

  return {
    applicationNo: String(raw?.applicationNo ?? raw?.application_no ?? ""),
    version: Number(raw?.version ?? 0),
    applicationStatusCode: statusCode,
    entity: normalizedEntity && hasNonEmptyField(normalizedEntity) ? normalizedEntity : undefined,
    shop: normalizedShop && hasNonEmptyField(normalizedShop) ? normalizedShop : undefined,
    submittedAt: normalizeTimestamp(raw?.submittedAt ?? raw?.submitted_at),
    updatedAt: normalizeTimestamp(raw?.updatedAt ?? raw?.updated_at)
  };
}

function normalizeApplicationStatus(rawStatus: unknown, fallback?: unknown): string {
  const value = rawStatus ?? fallback;
  if (typeof value === "string") {
    const upper = value.toUpperCase();
    return upper.startsWith("APPLICATION_STATUS_") ? upper.replace("APPLICATION_STATUS_", "") : upper;
  }
  if (typeof value === "number") {
    switch (value) {
      case 1:
        return "DRAFT";
      case 2:
        return "SUBMITTED";
      case 3:
        return "REVIEWING";
      case 4:
        return "APPROVED";
      case 5:
        return "REJECTED";
      case 6:
        return "CANCELLED";
      case 7:
        return "SYSTEM_REJECTED";
      default:
        return "UNSPECIFIED";
    }
  }
  return "UNSPECIFIED";
}

function normalizeTimestamp(raw: any): string | undefined {
  if (!raw) {
    return undefined;
  }
  if (typeof raw === "string") {
    return raw;
  }
  const seconds = Number(raw?.seconds ?? 0);
  if (!seconds) {
    return undefined;
  }
  return new Date(seconds * 1000).toISOString();
}

function normalizeApplicationItem(raw: any): SellerApplicationListItem {
  return {
    applicationNo: String(raw?.applicationNo ?? raw?.application_no ?? ""),
    version: Number(raw?.version ?? 0),
    applicationStatusCode: normalizeApplicationStatus(raw?.applicationStatusCode, raw?.status),
    submittedAt: normalizeTimestamp(raw?.submittedAt ?? raw?.submitted_at),
    updatedAt: normalizeTimestamp(raw?.updatedAt ?? raw?.updated_at)
  };
}

export async function fetchSellerWorkbench(): Promise<SellerWorkbenchResponse> {
  try {
    return await apiClient.get<SellerWorkbenchResponse>("/v1/seller/workbench");
  } catch (error) {
    console.warn("workbench fallback", error);
    return mockWorkbench;
  }
}

export async function fetchShopDashboard(shopNo: string): Promise<ShopDashboardResponse> {
  try {
    return await apiClient.get<ShopDashboardResponse>(`/v1/seller/shops/${shopNo}/dashboard`);
  } catch (error) {
    console.warn("shop dashboard fallback", error);
    return {
      shopNo,
      shopName: "Demo Shop",
      shopStatusCode: "ACTIVE",
      productSummary: { total: 0, onShelf: 0, offShelf: 0, reviewing: 0, draft: 0, rejected: 0 },
      inventoryRisk: { lowStockSkuCount: 0, outOfStockSkuCount: 0 },
      partial: true,
      degradedFields: ["inventory_risk"]
    };
  }
}

export async function createApplicationDraft(payload: CreateApplicationDraftPayload): Promise<SellerApplicationDetail> {
  const res = await apiClient.post<{ application?: any }>("/v1/seller/applications/draft", payload);
  return normalizeApplication(res.application ?? {});
}

export async function updateApplicationDraft(
  applicationNo: string,
  payload: UpdateApplicationDraftPayload
): Promise<SellerApplicationDetail> {
  const res = await apiClient.patch<{ application?: any }>(`/v1/seller/applications/${applicationNo}/draft`, payload);
  return normalizeApplication(res.application ?? {});
}

export async function submitApplication(
  applicationNo: string,
  payload: SubmitApplicationPayload
): Promise<SellerApplicationDetail> {
  const res = await apiClient.post<{ application?: any }>(`/v1/seller/applications/${applicationNo}/submit`, payload);
  return normalizeApplication(res.application ?? {});
}

export async function getMyApplication(applicationNo: string): Promise<SellerApplicationDetail> {
  const res = await apiClient.get<{ application?: any }>(`/v1/seller/applications/${applicationNo}`);
  return normalizeApplication(res.application ?? {});
}

export async function listMyApplications(payload?: ListMyApplicationsPayload): Promise<ListMyApplicationsResponse> {
  const page = payload?.page ?? 1;
  const pageSize = payload?.pageSize ?? 20;
  const res = await apiClient.get<{
    applications?: any[];
    page?: number;
    pageSize?: number;
    page_size?: number;
    total?: number;
  }>("/v1/seller/applications", {
    params: {
      page,
      pageSize,
      statuses: payload?.statuses
    }
  });

  return {
    applications: (res.applications ?? []).map(normalizeApplicationItem).filter((item) => item.applicationNo),
    page: Number(res.page ?? page),
    pageSize: Number(res.pageSize ?? res.page_size ?? pageSize),
    total: Number(res.total ?? 0)
  };
}
