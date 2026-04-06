import { apiClient } from "@shopa/api-client";
import { DerivedAsset, MediaAsset, UploadJob, UploadStatus } from "./types";

type UploadTicket = {
  method?: string;
  uploadUrl?: string;
  upload_url?: string;
  headers?: Record<string, string>;
};

type AssetResponse = {
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
  derivedAssets?: DerivedAssetResponse[];
  derived_assets?: DerivedAssetResponse[];
};

type DerivedAssetResponse = {
  assetId?: number;
  asset_id?: number;
  derivedKind?: string;
  derived_kind?: string;
  mimeType?: string;
  mime_type?: string;
  publicUrl?: string;
  public_url?: string;
  processStatus?: string | number;
  process_status?: string | number;
};

type BizAssetResponse = {
  asset?: AssetResponse;
};

type GetBizAssetsRes = {
  assets?: BizAssetResponse[];
};

type InitUploadRes = {
  asset?: AssetResponse;
  uploadTicket?: UploadTicket;
  upload_ticket?: UploadTicket;
};

type CompleteUploadRes = {
  asset?: AssetResponse;
};

type GetAssetProcessStatusRes = {
  assetId?: number;
  asset_id?: number;
  processStatus?: string | number;
  process_status?: string | number;
  processProgress?: number;
  process_progress?: number;
  derivedAssets?: DerivedAssetResponse[];
  derived_assets?: DerivedAssetResponse[];
};

type ReplaceBindingsRes = {
  bindings?: Array<{ assetId?: number; asset_id?: number }>;
};

export const OSS_CORS_BLOCKED_ERROR = "OSS_CORS_BLOCKED";

const MEDIA_SCENE_CODE = process.env.NEXT_PUBLIC_MEDIA_SCENE_CODE ?? "seller_media";
const DEFAULT_MEDIA_BIZ_TYPE = process.env.NEXT_PUBLIC_MEDIA_BIZ_TYPE ?? "seller_product_publish";
const DEFAULT_MEDIA_BINDING_FIELD = process.env.NEXT_PUBLIC_MEDIA_BINDING_FIELD ?? "gallery";
const CHECKSUM_PLACEHOLDER = "0000000000000000000000000000000000000000000000000000000000000000";
const UPLOAD_ID_PREFIX = "asset-";
const THUMBNAIL_FALLBACK = "/images/media/storefront.png";
const LOCAL_UPLOAD_HOST = "upload.shopa.local";
const LOCAL_READ_HOST = "read.shopa.local";

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

function toUploadStatus(value?: string | number): UploadStatus {
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

function uploadIdFromAssetId(assetId: number): string {
  return `${UPLOAD_ID_PREFIX}${assetId}`;
}

function assetIdFromUploadId(uploadId: string): number | undefined {
  const matched = uploadId.match(/^asset-(\d+)$/);
  if (!matched) {
    return undefined;
  }
  return Number(matched[1]);
}

function toDerivedAssets(items?: DerivedAssetResponse[]): DerivedAsset[] {
  return (items ?? []).map((item) => ({
    assetId: item.assetId ?? item.asset_id ?? 0,
    kind: item.derivedKind ?? item.derived_kind,
    mimeType: item.mimeType ?? item.mime_type,
    publicUrl: item.publicUrl ?? item.public_url,
    processStatus: item.processStatus ?? item.process_status
  }));
}

function toMediaAsset(asset?: AssetResponse, fallbackPreview?: string): MediaAsset {
  const assetId = asset?.assetId ?? asset?.asset_id ?? 0;
  const publicUrl = asset?.publicUrl ?? asset?.public_url;
  const processStatus = asset?.processStatus ?? asset?.process_status;
  const processProgress = asset?.processProgress ?? asset?.process_progress ?? 0;
  const derivedAssets = toDerivedAssets(asset?.derivedAssets ?? asset?.derived_assets);
  return {
    id: uploadIdFromAssetId(assetId),
    assetId,
    name: asset?.fileName ?? asset?.file_name ?? `Asset-${assetId}`,
    thumbnail: publicUrl ?? fallbackPreview ?? THUMBNAIL_FALLBACK,
    status: toUploadStatus(processStatus),
    progress: processProgress,
    uploadedAt: asset?.processedAt ?? asset?.processed_at ?? asset?.createdAt ?? asset?.created_at,
    publicUrl,
    processStatus,
    derivedAssets
  };
}

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

function buildBindingItems(assetIds: Array<number | string>) {
  return assetIds
    .map((assetId, index) => ({
      assetId: Number(assetId),
      sortOrder: index
    }))
    .filter((item) => item.assetId > 0);
}

export async function replaceMediaBindings(params: {
  bizNo: string;
  assetIds: Array<number | string>;
  bizType?: string;
  bindingField?: string;
}) {
  if (!params.bizNo) {
    return [];
  }

  const response = await apiClient.post<ReplaceBindingsRes>("/v1/media/bindings/replace", {
    sceneCode: MEDIA_SCENE_CODE,
    bizType: params.bizType ?? DEFAULT_MEDIA_BIZ_TYPE,
    bizNo: params.bizNo,
    bindingField: params.bindingField ?? DEFAULT_MEDIA_BINDING_FIELD,
    items: buildBindingItems(params.assetIds)
  });

  return response.bindings ?? [];
}

export async function fetchMediaLibrary(params?: {
  bizNo?: string;
  bizType?: string;
  bindingField?: string;
}): Promise<MediaAsset[]> {
  const bizNo = params?.bizNo ?? "seller-console-library";
  try {
    const response = await apiClient.get<GetBizAssetsRes>("/v1/media/biz-assets", {
      params: {
        sceneCode: MEDIA_SCENE_CODE,
        bizType: params?.bizType ?? DEFAULT_MEDIA_BIZ_TYPE,
        bizNo,
        bindingField: params?.bindingField ?? DEFAULT_MEDIA_BINDING_FIELD
      },
      silentDegraded: true
    });
    return (response.assets ?? []).map((item) => toMediaAsset(item.asset)).filter((item) => item.assetId > 0);
  } catch {
    return [];
  }
}

export async function uploadSellerMediaAsset(params: {
  file: File;
  bizNo: string;
  existingAssetIds?: Array<number | string>;
  bizType?: string;
  bindingField?: string;
}): Promise<MediaAsset> {
  const checksum = await sha256Hex(params.file);
  const init = await apiClient.post<InitUploadRes>("/v1/media/upload/init", {
    sceneCode: MEDIA_SCENE_CODE,
    bizType: params.bizType ?? DEFAULT_MEDIA_BIZ_TYPE,
    bizNo: params.bizNo,
    fileName: params.file.name,
    mimeType: params.file.type || "application/octet-stream",
    sizeBytes: params.file.size,
    checksumSha256: checksum
  });

  const assetId = Number(init.asset?.assetId ?? init.asset?.asset_id ?? 0);
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

  const complete = await apiClient.post<CompleteUploadRes>("/v1/media/upload/complete", {
    assetId,
    etag,
    sizeBytes: params.file.size,
    mimeType: params.file.type || "application/octet-stream",
    checksumSha256: checksum
  });

  await replaceMediaBindings({
    bizNo: params.bizNo,
    bizType: params.bizType,
    bindingField: params.bindingField,
    assetIds: [...(params.existingAssetIds ?? []), assetId]
  });

  return toMediaAsset(complete.asset, createObjectPreview(params.file));
}

export async function startUpload(): Promise<UploadJob> {
  return { uploadId: `asset-${Date.now()}`, status: "failed", progress: 0 };
}

export async function pollUploadStatus(uploadId: string): Promise<UploadJob> {
  const assetId = assetIdFromUploadId(uploadId);
  if (!assetId) {
    return { uploadId, status: "failed", progress: 0 };
  }
  try {
    const response = await apiClient.get<GetAssetProcessStatusRes>(`/v1/media/assets/${assetId}/process-status`, {
      silentDegraded: true
    });
    const processStatus = response.processStatus ?? response.process_status;
    const progress = response.processProgress ?? response.process_progress ?? 0;
    const derivedAssets = toDerivedAssets(response.derivedAssets ?? response.derived_assets);
    return {
      uploadId,
      assetId: response.assetId ?? response.asset_id ?? assetId,
      status: toUploadStatus(processStatus),
      progress,
      derivedAssets
    };
  } catch {
    return { uploadId, assetId, status: "failed", progress: 0 };
  }
}

