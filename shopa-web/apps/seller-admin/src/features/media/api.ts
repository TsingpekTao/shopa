"use client";

import { apiClient } from "@/lib/api-client";

type UploadTicket = {
  method?: string;
  uploadUrl?: string;
  upload_url?: string;
  headers?: Record<string, string>;
};

type AssetRes = {
  assetId?: number;
  asset_id?: number;
  publicUrl?: string;
  public_url?: string;
};

type InitUploadRes = {
  asset?: AssetRes;
  uploadTicket?: UploadTicket;
  upload_ticket?: UploadTicket;
};

type CompleteUploadRes = {
  asset?: AssetRes;
};

export type UploadAssetOptions = {
  bizType?: string;
  bizNo?: string;
};

export type UploadAssetResult = {
  assetId: string;
  previewUrl: string;
};

export const OSS_CORS_BLOCKED_ERROR = "OSS_CORS_BLOCKED";

const MEDIA_SCENE_CODE = process.env.NEXT_PUBLIC_MEDIA_SCENE_CODE ?? "seller_onboarding_cert";
const CHECKSUM_PLACEHOLDER = "0000000000000000000000000000000000000000000000000000000000000000";
const ONBOARDING_BIZ_TYPE = "seller_application";
const LOCAL_UPLOAD_HOST = "upload.shopa.local";
const LOCAL_READ_HOST = "read.shopa.local";

type SimpleBindingItem = {
  assetId: number;
  sortOrder: number;
};

function toHex(bytes: Uint8Array): string {
  return Array.from(bytes)
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");
}

async function sha256Hex(file: File): Promise<string> {
  try {
    const buffer = await file.arrayBuffer();
    const digest = await crypto.subtle.digest("SHA-256", buffer);
    return toHex(new Uint8Array(digest));
  } catch {
    return CHECKSUM_PLACEHOLDER;
  }
}

function trimEtag(value?: string | null): string {
  return (value || "").replaceAll('"', "").trim();
}

function isLocalPlaceholderUrl(raw: string): boolean {
  try {
    const url = new URL(raw);
    return url.hostname === LOCAL_UPLOAD_HOST || url.hostname === LOCAL_READ_HOST;
  } catch {
    return false;
  }
}

function createObjectPreview(file: File): string {
  if (typeof URL !== "undefined" && typeof URL.createObjectURL === "function") {
    return URL.createObjectURL(file);
  }
  return "";
}

export async function uploadCertificateAsset(file: File, options?: UploadAssetOptions): Promise<UploadAssetResult> {
  const checksum = await sha256Hex(file);
  const init = await apiClient.post<InitUploadRes>("/v1/media/upload/init", {
    sceneCode: MEDIA_SCENE_CODE,
    bizType: options?.bizType || ONBOARDING_BIZ_TYPE,
    bizNo: options?.bizNo || undefined,
    fileName: file.name,
    mimeType: file.type || "application/octet-stream",
    sizeBytes: file.size,
    checksumSha256: checksum
  });

  const assetId = String(init.asset?.assetId ?? init.asset?.asset_id ?? "");
  if (!assetId) {
    throw new Error("missing asset id");
  }

  const ticket = init.uploadTicket ?? init.upload_ticket;
  const uploadUrl = (ticket?.uploadUrl ?? ticket?.upload_url ?? "").trim();
  const method = (ticket?.method || "PUT").toUpperCase();

  let etag = "";
  if (uploadUrl && !isLocalPlaceholderUrl(uploadUrl)) {
    let resp: Response;
    try {
      resp = await fetch(uploadUrl, {
        method,
        headers: {
          ...(ticket?.headers || {}),
          "Content-Type": file.type || "application/octet-stream"
        },
        body: file
      });
    } catch (err) {
      const isCorsLikely =
        typeof uploadUrl === "string" &&
        uploadUrl.includes(".aliyuncs.com") &&
        err instanceof Error &&
        /failed to fetch|networkerror|network error/i.test(err.message || "");
      if (isCorsLikely) {
        throw new Error(OSS_CORS_BLOCKED_ERROR);
      }
      throw err;
    }
    if (!resp.ok) {
      throw new Error(`upload to object storage failed: ${resp.status}`);
    }
    etag = trimEtag(resp.headers.get("etag"));
  }

  const complete = await apiClient.post<CompleteUploadRes>("/v1/media/upload/complete", {
    assetId: Number(assetId),
    etag,
    sizeBytes: file.size,
    mimeType: file.type || "application/octet-stream",
    checksumSha256: checksum
  });

  const publicUrl = complete.asset?.publicUrl ?? complete.asset?.public_url;
  return {
    assetId,
    previewUrl: publicUrl || createObjectPreview(file)
  };
}

export async function issueAssetReadUrl(assetId: string | number, ttlSeconds = 900): Promise<string> {
  if (!assetId) {
    return "";
  }
  const data = await apiClient.get<{ url?: string }>(`/v1/media/assets/${assetId}/read-url`, {
    params: { ttlSeconds }
  });
  const url = String(data?.url ?? "");
  return isLocalPlaceholderUrl(url) ? "" : url;
}

async function replaceBizBinding(params: {
  bizNo: string;
  bindingField: string;
  items: SimpleBindingItem[];
}) {
  return apiClient.post("/v1/media/bindings/replace", {
    sceneCode: MEDIA_SCENE_CODE,
    bizType: ONBOARDING_BIZ_TYPE,
    bizNo: params.bizNo,
    bindingField: params.bindingField,
    items: params.items
  });
}

export async function syncOnboardingCertificateBindings(params: {
  applicationNo: string;
  idCardFrontAssetId?: string;
  idCardBackAssetId?: string;
  businessLicenseAssetId?: string;
}) {
  const bizNo = (params.applicationNo || "").trim();
  if (!bizNo) {
    return;
  }

  const toItem = (assetId?: string): SimpleBindingItem[] => {
    const id = Number(assetId || 0);
    return id > 0 ? [{ assetId: id, sortOrder: 0 }] : [];
  };

  await Promise.all([
    replaceBizBinding({
      bizNo,
      bindingField: "id_card_front",
      items: toItem(params.idCardFrontAssetId)
    }),
    replaceBizBinding({
      bizNo,
      bindingField: "id_card_back",
      items: toItem(params.idCardBackAssetId)
    }),
    replaceBizBinding({
      bizNo,
      bindingField: "business_license",
      items: toItem(params.businessLicenseAssetId)
    })
  ]);
}
