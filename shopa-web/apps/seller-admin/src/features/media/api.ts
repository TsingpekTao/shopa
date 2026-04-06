"use client";

import { apiClient } from "@/lib/api-client";
import { MediaAsset } from "./types";

type UploadTicket = {
  method?: string;
  uploadUrl?: string;
  upload_url?: string;
  headers?: Record<string, string>;
};

type AssetRes = {
  assetId?: number;
  asset_id?: number;
  fileName?: string;
  file_name?: string;
  publicUrl?: string;
  public_url?: string;
  processStatus?: string | number;
  process_status?: string | number;
  processProgress?: number;
  process_progress?: number;
  processedAt?: string;
  processed_at?: string;
  createdAt?: string;
  created_at?: string;
};

type BizAssetRes = {
  asset?: AssetRes;
};

type GetBizAssetsRes = {
  assets?: BizAssetRes[];
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

const ONBOARDING_SCENE_CODE = process.env.NEXT_PUBLIC_MEDIA_SCENE_CODE ?? "seller_onboarding_cert";
const PRODUCT_SCENE_CODE =
  process.env.NEXT_PUBLIC_PRODUCT_MEDIA_SCENE_CODE ??
  process.env.NEXT_PUBLIC_MEDIA_PRODUCT_SCENE_CODE ??
  "seller_media";
const CHECKSUM_PLACEHOLDER = "0000000000000000000000000000000000000000000000000000000000000000";
const ONBOARDING_BIZ_TYPE = "seller_application";
const PRODUCT_BIZ_TYPE =
  process.env.NEXT_PUBLIC_PRODUCT_MEDIA_BIZ_TYPE ??
  process.env.NEXT_PUBLIC_MEDIA_PRODUCT_BIZ_TYPE ??
  "seller_product_publish";
const PRODUCT_BINDING_FIELD =
  process.env.NEXT_PUBLIC_PRODUCT_MEDIA_BINDING_FIELD ??
  process.env.NEXT_PUBLIC_MEDIA_PRODUCT_BINDING_FIELD ??
  "gallery";
export const SELLER_MEDIA_LIBRARY_BIZ_NO =
  process.env.NEXT_PUBLIC_SELLER_MEDIA_LIBRARY_BIZ_NO ?? "seller-console-library";
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

function normalizeProcessStatus(value?: string | number): string {
  if (typeof value === "string") {
    return value;
  }
  if (typeof value === "number") {
    switch (value) {
      case 1:
        return "PROCESS_STATUS_PENDING";
      case 2:
        return "PROCESS_STATUS_PROCESSING";
      case 3:
        return "PROCESS_STATUS_SUCCESS";
      case 4:
        return "PROCESS_STATUS_FAILED";
      case 5:
        return "PROCESS_STATUS_PARTIAL_SUCCESS";
      default:
        return "PROCESS_STATUS_UNSPECIFIED";
    }
  }
  return "PROCESS_STATUS_UNSPECIFIED";
}

function toSellerMediaStatus(value?: string | number): MediaAsset["status"] {
  const processStatus = normalizeProcessStatus(value);
  switch (processStatus) {
    case "PROCESS_STATUS_PENDING":
      return "uploading";
    case "PROCESS_STATUS_PROCESSING":
    case "PROCESS_STATUS_PARTIAL_SUCCESS":
      return "processing";
    case "PROCESS_STATUS_SUCCESS":
      return "done";
    case "PROCESS_STATUS_FAILED":
      return "failed";
    default:
      return "idle";
  }
}

function toSellerMediaAsset(asset?: AssetRes, fallbackPreview?: string): MediaAsset {
  const assetId = Number(asset?.assetId ?? asset?.asset_id ?? 0);
  const publicUrl = String(asset?.publicUrl ?? asset?.public_url ?? "");
  return {
    id: `asset-${assetId}`,
    assetId,
    name: String(asset?.fileName ?? asset?.file_name ?? `Asset-${assetId}`),
    thumbnail: publicUrl || fallbackPreview || "",
    status: toSellerMediaStatus(asset?.processStatus ?? asset?.process_status),
    progress: Number(asset?.processProgress ?? asset?.process_progress ?? 0),
    uploadedAt: String(asset?.processedAt ?? asset?.processed_at ?? asset?.createdAt ?? asset?.created_at ?? "") || undefined,
    publicUrl: publicUrl || undefined
  };
}

async function withReadableUrl(asset: MediaAsset): Promise<MediaAsset> {
  if (!asset.assetId) {
    return asset;
  }
  try {
    const readableUrl = await issueAssetReadUrl(asset.assetId);
    if (!readableUrl) {
      return asset;
    }
    return {
      ...asset,
      publicUrl: readableUrl,
      thumbnail: readableUrl
    };
  } catch {
    return asset;
  }
}

function buildBindingItems(assetIds: Array<number | string>) {
  return assetIds
    .map((assetId, index) => ({
      assetId: Number(assetId),
      sortOrder: index
    }))
    .filter((item) => item.assetId > 0);
}

async function initUpload(params: {
  sceneCode: string;
  bizType: string;
  bizNo: string;
  file: File;
}) {
  const checksum = await sha256Hex(params.file);
  const init = await apiClient.post<InitUploadRes>("/v1/media/upload/init", {
    sceneCode: params.sceneCode,
    bizType: params.bizType,
    bizNo: params.bizNo || undefined,
    fileName: params.file.name,
    mimeType: params.file.type || "application/octet-stream",
    sizeBytes: params.file.size,
    checksumSha256: checksum
  });

  const assetId = String(init.asset?.assetId ?? init.asset?.asset_id ?? "");
  if (!assetId) {
    throw new Error("missing asset id");
  }

  return {
    assetId,
    init,
    checksum
  };
}

async function uploadByTicket(params: {
  uploadTicket?: UploadTicket;
  file: File;
}) {
  const ticket = params.uploadTicket;
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
          "Content-Type": params.file.type || "application/octet-stream"
        },
        body: params.file
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
  return etag;
}

async function completeUpload(params: {
  assetId: string;
  checksum: string;
  file: File;
  etag: string;
}) {
  return apiClient.post<CompleteUploadRes>("/v1/media/upload/complete", {
    assetId: Number(params.assetId),
    etag: params.etag,
    sizeBytes: params.file.size,
    mimeType: params.file.type || "application/octet-stream",
    checksumSha256: params.checksum
  });
}

async function replaceBizBinding(params: {
  sceneCode: string;
  bizType: string;
  bizNo: string;
  bindingField: string;
  items: SimpleBindingItem[];
}) {
  return apiClient.post("/v1/media/bindings/replace", {
    sceneCode: params.sceneCode,
    bizType: params.bizType,
    bizNo: params.bizNo,
    bindingField: params.bindingField,
    items: params.items
  });
}

async function batchUnbindBizBinding(params: {
  sceneCode: string;
  bizType: string;
  bizNo: string;
  bindingField: string;
  assetIds: Array<number | string>;
}) {
  return apiClient.post<{ affectedRows?: number; affected_rows?: number }>("/v1/media/bindings/batch-unbind", {
    sceneCode: params.sceneCode,
    bizType: params.bizType,
    bizNo: params.bizNo,
    bindingField: params.bindingField,
    assetIds: params.assetIds.map((assetId) => Number(assetId)).filter((assetId) => assetId > 0)
  });
}

export async function fetchSellerProductMediaLibrary(params: {
  bizNo: string;
  bizType?: string;
  bindingField?: string;
}): Promise<MediaAsset[]> {
  if (!params.bizNo) {
    return [];
  }

  try {
    const response = await apiClient.get<GetBizAssetsRes>("/v1/media/biz-assets", {
      params: {
        sceneCode: PRODUCT_SCENE_CODE,
        bizType: params.bizType ?? PRODUCT_BIZ_TYPE,
        bizNo: params.bizNo,
        bindingField: params.bindingField ?? PRODUCT_BINDING_FIELD
      },
      silentDegraded: true
    });
    const assets = (response.assets ?? []).map((item) => toSellerMediaAsset(item.asset)).filter((item) => item.assetId > 0);
    return Promise.all(assets.map((asset) => withReadableUrl(asset)));
  } catch {
    return [];
  }
}

export async function replaceSellerProductMediaBindings(params: {
  bizNo: string;
  assetIds: Array<number | string>;
  bizType?: string;
  bindingField?: string;
}) {
  if (!params.bizNo) {
    return [];
  }

  return replaceBizBinding({
    sceneCode: PRODUCT_SCENE_CODE,
    bizType: params.bizType ?? PRODUCT_BIZ_TYPE,
    bizNo: params.bizNo,
    bindingField: params.bindingField ?? PRODUCT_BINDING_FIELD,
    items: buildBindingItems(params.assetIds)
  });
}

export async function uploadSellerProductAsset(params: {
  file: File;
  bizNo: string;
  existingAssetIds?: Array<number | string>;
  bizType?: string;
  bindingField?: string;
}): Promise<MediaAsset> {
  const { assetId, init, checksum } = await initUpload({
    sceneCode: PRODUCT_SCENE_CODE,
    bizType: params.bizType ?? PRODUCT_BIZ_TYPE,
    bizNo: params.bizNo,
    file: params.file
  });

  const etag = await uploadByTicket({
    uploadTicket: init.uploadTicket ?? init.upload_ticket,
    file: params.file
  });

  const complete = await completeUpload({
    assetId,
    checksum,
    file: params.file,
    etag
  });

  await replaceSellerProductMediaBindings({
    bizNo: params.bizNo,
    bizType: params.bizType,
    bindingField: params.bindingField,
    assetIds: [...(params.existingAssetIds ?? []), Number(assetId)]
  });

  const readableUrl = await issueAssetReadUrl(assetId).catch(() => "");
  return toSellerMediaAsset(
    {
      assetId: Number(assetId),
      fileName: params.file.name,
      publicUrl: readableUrl || complete.asset?.publicUrl || complete.asset?.public_url,
      processStatus: "PROCESS_STATUS_SUCCESS",
      processProgress: 100,
      createdAt: new Date().toISOString()
    },
    readableUrl || createObjectPreview(params.file)
  );
}

export async function uploadCertificateAsset(file: File, options?: UploadAssetOptions): Promise<UploadAssetResult> {
  const { assetId, init, checksum } = await initUpload({
    sceneCode: ONBOARDING_SCENE_CODE,
    bizType: options?.bizType || ONBOARDING_BIZ_TYPE,
    bizNo: options?.bizNo || "",
    file
  });

  const etag = await uploadByTicket({
    uploadTicket: init.uploadTicket ?? init.upload_ticket,
    file
  });

  const complete = await completeUpload({
    assetId,
    checksum,
    file,
    etag
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
      sceneCode: ONBOARDING_SCENE_CODE,
      bizType: ONBOARDING_BIZ_TYPE,
      bizNo,
      bindingField: "id_card_front",
      items: toItem(params.idCardFrontAssetId)
    }),
    replaceBizBinding({
      sceneCode: ONBOARDING_SCENE_CODE,
      bizType: ONBOARDING_BIZ_TYPE,
      bizNo,
      bindingField: "id_card_back",
      items: toItem(params.idCardBackAssetId)
    }),
    replaceBizBinding({
      sceneCode: ONBOARDING_SCENE_CODE,
      bizType: ONBOARDING_BIZ_TYPE,
      bizNo,
      bindingField: "business_license",
      items: toItem(params.businessLicenseAssetId)
    })
  ]);
}

export async function fetchMediaLibrary(params?: {
  bizNo?: string;
  bizType?: string;
  bindingField?: string;
}) {
  return fetchSellerProductMediaLibrary({
    bizNo: params?.bizNo ?? SELLER_MEDIA_LIBRARY_BIZ_NO,
    bizType: params?.bizType,
    bindingField: params?.bindingField
  });
}

export async function replaceMediaBindings(params: {
  bizNo: string;
  assetIds: Array<number | string>;
  bizType?: string;
  bindingField?: string;
}) {
  return replaceSellerProductMediaBindings(params);
}

export async function uploadSellerMediaAsset(params: {
  file: File;
  bizNo: string;
  existingAssetIds?: Array<number | string>;
  bizType?: string;
  bindingField?: string;
}) {
  return uploadSellerProductAsset(params);
}

export async function batchUnbindSellerMediaAssets(params: {
  assetIds: Array<number | string>;
  bizNo?: string;
  bizType?: string;
  bindingField?: string;
}) {
  return batchUnbindBizBinding({
    sceneCode: PRODUCT_SCENE_CODE,
    bizType: params.bizType ?? PRODUCT_BIZ_TYPE,
    bizNo: params.bizNo ?? SELLER_MEDIA_LIBRARY_BIZ_NO,
    bindingField: params.bindingField ?? PRODUCT_BINDING_FIELD,
    assetIds: params.assetIds
  });
}
