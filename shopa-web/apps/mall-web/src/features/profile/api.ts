import { apiClient } from "@/lib/api-client";
import { MallAddress, MallProfile, MallProfileBundle, ProfileImageRef } from "./types";

type RawProfileImage = {
  asset_id?: number | string;
  assetId?: number | string;
  url?: string;
};

type RawProfile = {
  display_name?: string;
  displayName?: string;
  profile_version?: number | string;
  profileVersion?: number | string;
  avatar?: RawProfileImage;
  updated_at?: string | { seconds?: number | string; nanos?: number | string };
  updatedAt?: string | { seconds?: number | string; nanos?: number | string };
  ext?: Record<string, string>;
};

type RawAddress = {
  is_default?: boolean;
  isDefault?: boolean;
};

type RawGetMyProfileRes = {
  profile?: RawProfile;
  addresses?: RawAddress[];
};

type RawUpdateMyProfileRes = {
  profile?: RawProfile;
};

function toNumber(value: unknown, fallback = 0): number {
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : fallback;
}

function toString(value: unknown): string {
  if (value === null || value === undefined) {
    return "";
  }
  return String(value);
}

async function issueAssetReadUrl(assetId: number, ttlSeconds = 900): Promise<string> {
  if (!assetId || typeof fetch !== "function") {
    return "";
  }

  try {
    const params = new URLSearchParams({
      assetIds: String(assetId),
      ttlSeconds: String(ttlSeconds)
    });
    const response = await fetch(`/api/media/read-urls?${params.toString()}`, {
      method: "GET",
      cache: "no-store"
    });
    if (!response.ok) {
      return "";
    }

    const payload = (await response.json()) as {
      items?: Array<{ assetId?: string | number; asset_id?: string | number; url?: string }>;
    };
    const item = (payload.items ?? []).find((candidate) => toString(candidate.assetId ?? candidate.asset_id) === String(assetId));
    return toString(item?.url);
  } catch {
    return "";
  }
}

async function normalizeImage(raw?: RawProfileImage): Promise<ProfileImageRef> {
  const assetId = toNumber(raw?.asset_id ?? raw?.assetId);
  const fallbackUrl = toString(raw?.url);
  const readableUrl = await issueAssetReadUrl(assetId);
  return {
    assetId,
    url: readableUrl || fallbackUrl
  };
}

async function normalizeProfile(raw?: RawProfile): Promise<MallProfile | undefined> {
  if (!raw) {
    return undefined;
  }
  return {
    displayName: toString(raw.display_name ?? raw.displayName),
    avatar: await normalizeImage(raw.avatar),
    profileVersion: toNumber(raw.profile_version ?? raw.profileVersion),
    ext: raw.ext ?? {},
    updatedAt: (raw.updated_at ?? raw.updatedAt) as MallProfile["updatedAt"]
  };
}

function normalizeAddresses(rows?: RawAddress[]): MallAddress[] {
  return (rows ?? []).map((item) => ({
    isDefault: Boolean(item.is_default ?? item.isDefault)
  }));
}

export async function getMyProfile(includeAddresses = true): Promise<MallProfileBundle> {
  const response = await apiClient.get<RawGetMyProfileRes>("/v1/me/profile", {
    params: {
      include_addresses: includeAddresses
    }
  });
  return {
    profile: await normalizeProfile(response.profile),
    addresses: normalizeAddresses(response.addresses)
  };
}

export async function updateMyProfile(payload: {
  expectedProfileVersion: number;
  displayName?: string;
  avatar?: ProfileImageRef;
}): Promise<MallProfile> {
  const profilePatch: Record<string, unknown> = {};
  const updateMask: string[] = [];

  if (typeof payload.displayName === "string") {
    profilePatch.display_name = payload.displayName;
    updateMask.push("display_name");
  }

  if (payload.avatar) {
    profilePatch.avatar = {
      asset_id: payload.avatar.assetId,
      url: payload.avatar.url
    };
    updateMask.push("avatar");
  }

  if (updateMask.length === 0) {
    throw new Error("no profile fields to update");
  }

  const response = await apiClient.patch<RawUpdateMyProfileRes>("/v1/me/profile", {
    profile: profilePatch,
    update_mask: updateMask,
    expected_profile_version: Math.max(1, toNumber(payload.expectedProfileVersion, 1))
  });

  const normalized = await normalizeProfile(response.profile);
  if (!normalized) {
    throw new Error("profile response is empty");
  }
  return normalized;
}
