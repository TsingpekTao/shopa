import { apiClient } from "@/lib/api-client";

type RawUploadTicket = {
  method?: string;
  uploadUrl?: string;
  upload_url?: string;
  headers?: Record<string, string>;
};

type RawAsset = {
  assetId?: number | string;
  asset_id?: number | string;
  publicUrl?: string;
  public_url?: string;
};

type RawInitUploadRes = {
  asset?: RawAsset;
  uploadTicket?: RawUploadTicket;
  upload_ticket?: RawUploadTicket;
};

type RawCompleteUploadRes = {
  asset?: RawAsset;
};

type RawReadUrlRes = {
  url?: string;
};

const AVATAR_SCENE_CODE = process.env.NEXT_PUBLIC_MEDIA_AVATAR_SCENE_CODE ?? "buyer_avatar";
const AVATAR_BIZ_TYPE = process.env.NEXT_PUBLIC_MEDIA_AVATAR_BIZ_TYPE ?? "buyer_profile_avatar";
const CHECKSUM_PLACEHOLDER = "0000000000000000000000000000000000000000000000000000000000000000";
const LOCAL_UPLOAD_HOST = "upload.shopa.local";
const LOCAL_READ_HOST = "read.shopa.local";

function toNumber(value: unknown, fallback = 0): number {
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : fallback;
}

function toHex(bytes: Uint8Array): string {
  return Array.from(bytes)
    .map((item) => item.toString(16).padStart(2, "0"))
    .join("");
}

function trimEtag(raw?: string | null): string {
  return (raw ?? "").replaceAll('"', "").trim();
}

function isLocalPlaceholderUrl(raw: string): boolean {
  try {
    const parsed = new URL(raw);
    return parsed.hostname === LOCAL_UPLOAD_HOST || parsed.hostname === LOCAL_READ_HOST;
  } catch {
    return false;
  }
}

function toFriendlyAvatarUploadError(error: unknown): Error {
  const message = error instanceof Error ? error.message : String(error ?? "");
  const normalized = message.toLowerCase();
  if (normalized.includes("media_scene_policy") || normalized.includes("scene policy")) {
    return new Error("avatar_upload_policy_missing");
  }
  if (normalized.includes("upload failed")) {
    return new Error("avatar_upload_transport_failed");
  }
  return new Error("avatar_upload_failed");
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

async function issueAssetReadUrl(assetId: number, ttlSeconds = 3600): Promise<string> {
  if (!assetId) {
    return "";
  }
  const response = await apiClient.get<RawReadUrlRes>(`/v1/media/assets/${assetId}/read-url`, {
    params: { ttlSeconds }
  });
  const url = String(response?.url ?? "");
  return isLocalPlaceholderUrl(url) ? "" : url;
}

export async function uploadBuyerAvatar(file: File): Promise<{ assetId: number; url: string }> {
  try {
    const checksumSha256 = await sha256Hex(file);
    const init = await apiClient.post<RawInitUploadRes>("/v1/media/upload/init", {
      sceneCode: AVATAR_SCENE_CODE,
      bizType: AVATAR_BIZ_TYPE,
      fileName: file.name,
      mimeType: file.type || "application/octet-stream",
      sizeBytes: file.size,
      checksumSha256
    });

    const assetId = toNumber(init.asset?.assetId ?? init.asset?.asset_id);
    if (!assetId) {
      throw new Error("missing asset id from upload init");
    }

    const uploadTicket = init.uploadTicket ?? init.upload_ticket;
    const uploadUrl = (uploadTicket?.uploadUrl ?? uploadTicket?.upload_url ?? "").trim();
    const uploadMethod = (uploadTicket?.method || "PUT").toUpperCase();
    let etag = "";

    if (uploadUrl && !isLocalPlaceholderUrl(uploadUrl)) {
      const uploadResp = await fetch(uploadUrl, {
        method: uploadMethod,
        headers: {
          ...(uploadTicket?.headers || {}),
          "Content-Type": file.type || "application/octet-stream"
        },
        body: file
      });
      if (!uploadResp.ok) {
        throw new Error(`upload failed: ${uploadResp.status}`);
      }
      etag = trimEtag(uploadResp.headers.get("etag"));
    }

    const complete = await apiClient.post<RawCompleteUploadRes>("/v1/media/upload/complete", {
      assetId,
      etag,
      sizeBytes: file.size,
      mimeType: file.type || "application/octet-stream",
      checksumSha256
    });

    const maybePublicUrl = String(complete.asset?.publicUrl ?? complete.asset?.public_url ?? "");
    const readableUrl = maybePublicUrl || (await issueAssetReadUrl(assetId));

    return {
      assetId,
      url: readableUrl
    };
  } catch (error) {
    throw toFriendlyAvatarUploadError(error);
  }
}
