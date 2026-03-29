import { apiClient } from "@shopa/api-client";
import { DerivedAsset, MediaAsset, UploadJob, UploadStatus } from "./types";

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

const MEDIA_SCENE_CODE = process.env.NEXT_PUBLIC_MEDIA_SCENE_CODE ?? "seller_media";
const MEDIA_BIZ_TYPE = process.env.NEXT_PUBLIC_MEDIA_BIZ_TYPE ?? "seller_media_library";
const MEDIA_BIZ_NO = process.env.NEXT_PUBLIC_MEDIA_BIZ_NO ?? "default";
const MEDIA_BINDING_FIELD = process.env.NEXT_PUBLIC_MEDIA_BINDING_FIELD ?? "main_images";
const CHECKSUM_PLACEHOLDER = "0000000000000000000000000000000000000000000000000000000000000000";
const UPLOAD_ID_PREFIX = "asset-";
const THUMBNAIL_FALLBACK = "/images/media/preview-brand.jpg";

const mockAssets: MediaAsset[] = [
  {
    id: "asset-1",
    assetId: 1,
    name: "Brand Authorization.pdf",
    thumbnail: "/images/media/preview-brand.jpg",
    status: "processing",
    progress: 58,
    uploadedAt: "2026-03-24T10:33:00Z"
  },
  {
    id: "asset-2",
    assetId: 2,
    name: "Storefront Cover.png",
    thumbnail: "/images/media/storefront.png",
    status: "done",
    progress: 100,
    uploadedAt: "2026-03-24T09:12:00Z"
  }
];

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

function toMediaAsset(asset?: AssetResponse): MediaAsset {
  const assetId = asset?.assetId ?? asset?.asset_id ?? 0;
  const publicUrl = asset?.publicUrl ?? asset?.public_url;
  const processStatus = asset?.processStatus ?? asset?.process_status;
  const processProgress = asset?.processProgress ?? asset?.process_progress ?? 0;
  const derivedAssets = toDerivedAssets(asset?.derivedAssets ?? asset?.derived_assets);
  return {
    id: uploadIdFromAssetId(assetId),
    assetId,
    name: asset?.fileName ?? asset?.file_name ?? `Asset-${assetId}`,
    thumbnail: publicUrl ?? THUMBNAIL_FALLBACK,
    status: toUploadStatus(processStatus),
    progress: processProgress,
    uploadedAt: asset?.processedAt ?? asset?.processed_at ?? asset?.createdAt ?? asset?.created_at,
    publicUrl,
    processStatus,
    derivedAssets
  };
}

export async function fetchMediaLibrary(): Promise<MediaAsset[]> {
  try {
    const response = await apiClient.get<GetBizAssetsRes>("/v1/media/biz-assets", {
      params: {
        sceneCode: MEDIA_SCENE_CODE,
        bizType: MEDIA_BIZ_TYPE,
        bizNo: MEDIA_BIZ_NO,
        bindingField: MEDIA_BINDING_FIELD
      }
    });
    const list = (response.assets ?? []).map((item) => toMediaAsset(item.asset)).filter((item) => item.assetId > 0);
    return list.length > 0 ? list : mockAssets;
  } catch (error) {
    console.warn("media library fallback", error);
    return mockAssets;
  }
}

export async function startUpload(): Promise<UploadJob> {
  try {
    const response = await apiClient.post<InitUploadRes>("/v1/media/upload/init", {
      sceneCode: MEDIA_SCENE_CODE,
      bizType: MEDIA_BIZ_TYPE,
      bizNo: MEDIA_BIZ_NO,
      fileName: `upload-${Date.now()}.bin`,
      mimeType: "application/octet-stream",
      sizeBytes: 128,
      checksumSha256: CHECKSUM_PLACEHOLDER
    });
    const asset = toMediaAsset(response.asset);
    return {
      uploadId: asset.id,
      assetId: asset.assetId,
      status: asset.status,
      progress: asset.progress,
      derivedAssets: asset.derivedAssets
    };
  } catch (error) {
    console.warn("start upload fallback", error);
    return { uploadId: `asset-${Date.now()}`, status: "failed", progress: 0 };
  }
}

export async function pollUploadStatus(uploadId: string): Promise<UploadJob> {
  const assetId = assetIdFromUploadId(uploadId);
  if (!assetId) {
    return { uploadId, status: "failed", progress: 0 };
  }
  try {
    const response = await apiClient.get<GetAssetProcessStatusRes>(`/v1/media/assets/${assetId}/process-status`);
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
  } catch (error) {
    console.warn("poll upload status fallback", error);
    return { uploadId, assetId, status: "failed", progress: 0 };
  }
}
