"use client";

import Link from "next/link";
import { ChangeEvent, FormEvent, useEffect, useMemo, useRef, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useI18n } from "@shopa/ui";
import { useAuth } from "@/features/iam/useAuth";
import { uploadBuyerAvatar } from "@/features/media/api";
import { getMyOverview } from "@/features/overview/api";
import { getMyProfile, updateMyProfile } from "@/features/profile/api";
import { MallProfile, MallProfileBundle } from "@/features/profile/types";

type LocalIdentitySnapshot = {
  displayName: string;
  identifier: string;
  avatarUrl: string;
};

const AUTH_CHANGED_EVENT = "shopa-mall-auth-changed";
const DISPLAY_NAME_STORAGE_KEY = "shopa_mall_display_name";
const AVATAR_URL_STORAGE_KEY = "shopa_mall_avatar_url";

function readAccessToken(): string {
  if (typeof window === "undefined") {
    return "";
  }
  const direct = window.localStorage.getItem("shopa_mall_access_token")?.trim();
  if (direct) {
    return direct;
  }
  const persisted = window.localStorage.getItem("shopa-mall-auth");
  if (!persisted) {
    return "";
  }
  try {
    const parsed = JSON.parse(persisted) as {
      state?: { tokenPair?: { accessToken?: string } };
    };
    return parsed?.state?.tokenPair?.accessToken?.trim() ?? "";
  } catch {
    return "";
  }
}

function readLocalIdentitySnapshot(): LocalIdentitySnapshot {
  if (typeof window === "undefined") {
    return { displayName: "", identifier: "", avatarUrl: "" };
  }
  return {
    displayName: window.localStorage.getItem(DISPLAY_NAME_STORAGE_KEY) ?? "",
    identifier: window.localStorage.getItem("shopa_mall_last_identifier") ?? "",
    avatarUrl: window.localStorage.getItem(AVATAR_URL_STORAGE_KEY) ?? ""
  };
}

function normalizeDisplayName(raw: string): string {
  const value = raw.trim();
  if (!value) {
    return "";
  }
  if (/^\d{11}$/.test(value)) {
    return `${value.slice(0, 3)}****${value.slice(-4)}`;
  }
  if (/^\d+$/.test(value)) {
    return "";
  }
  return value;
}

function maskPhone(raw?: string): string {
  if (!raw) {
    return "-";
  }
  const value = raw.trim();
  if (/^\d{11}$/.test(value)) {
    return `${value.slice(0, 3)}****${value.slice(-4)}`;
  }
  return value;
}

function toDate(value?: string | { seconds?: number | string; nanos?: number | string }): Date | null {
  if (!value) {
    return null;
  }
  if (typeof value === "string") {
    const date = new Date(value);
    return Number.isNaN(date.getTime()) ? null : date;
  }

  const seconds = Number(value.seconds ?? 0);
  const nanos = Number(value.nanos ?? 0);
  if (!Number.isFinite(seconds) && !Number.isFinite(nanos)) {
    return null;
  }
  const milliseconds = seconds * 1000 + Math.floor(nanos / 1_000_000);
  const date = new Date(milliseconds);
  return Number.isNaN(date.getTime()) ? null : date;
}

function formatTime(raw?: string | { seconds?: number | string; nanos?: number | string }, locale?: string): string {
  const dt = toDate(raw);
  if (!dt) {
    return "-";
  }
  return dt.toLocaleString(locale === "zh-CN" ? "zh-CN" : "en-US", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit"
  });
}

function syncLocalIdentity(displayName: string, avatarUrl: string) {
  if (typeof window === "undefined") {
    return;
  }
  if (displayName.trim()) {
    window.localStorage.setItem(DISPLAY_NAME_STORAGE_KEY, displayName.trim());
  } else {
    window.localStorage.removeItem(DISPLAY_NAME_STORAGE_KEY);
  }
  if (avatarUrl.trim()) {
    window.localStorage.setItem(AVATAR_URL_STORAGE_KEY, avatarUrl.trim());
  } else {
    window.localStorage.removeItem(AVATAR_URL_STORAGE_KEY);
  }
  window.dispatchEvent(new Event(AUTH_CHANGED_EVENT));
}

function toFriendlyProfileNotice(error: unknown, isZh: boolean): string {
  const raw = error instanceof Error ? error.message : String(error ?? "");
  switch (raw) {
    case "avatar_upload_policy_missing":
      return isZh ? "\u5934\u50cf\u4e0a\u4f20\u914d\u7f6e\u7f3a\u5931\uff0c\u8bf7\u7a0d\u540e\u518d\u8bd5" : "Avatar upload policy is unavailable. Please try again later.";
    case "avatar_upload_transport_failed":
      return isZh ? "\u5934\u50cf\u4e0a\u4f20\u5931\u8d25\uff0c\u8bf7\u68c0\u67e5\u7f51\u7edc\u540e\u91cd\u8bd5" : "Avatar upload failed. Please check your network and try again.";
    case "avatar_upload_failed":
      return isZh ? "\u5934\u50cf\u4e0a\u4f20\u5931\u8d25\uff0c\u8bf7\u7a0d\u540e\u91cd\u8bd5" : "Avatar upload failed. Please try again later.";
    default:
      return raw || (isZh ? "\u64cd\u4f5c\u5931\u8d25\uff0c\u8bf7\u7a0d\u540e\u91cd\u8bd5" : "Something went wrong. Please try again later.");
  }
}

function mergeProfile(prev: MallProfile | undefined, next: MallProfile): MallProfile {
  if (!prev) {
    return next;
  }
  return {
    ...prev,
    ...next,
    avatar: {
      ...prev.avatar,
      ...next.avatar
    },
    ext: {
      ...(prev.ext ?? {}),
      ...(next.ext ?? {})
    }
  };
}

export default function ProfilePage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const { tokenPair, logout } = useAuth();
  const queryClient = useQueryClient();
  const avatarInputRef = useRef<HTMLInputElement | null>(null);
  const blobAvatarPreviewRef = useRef("");

  const [accessToken, setAccessToken] = useState("");
  const [localIdentity, setLocalIdentity] = useState<LocalIdentitySnapshot>({
    displayName: "",
    identifier: "",
    avatarUrl: ""
  });
  const [displayNameDraft, setDisplayNameDraft] = useState("");
  const [profileNotice, setProfileNotice] = useState("");
  const [profileNoticeError, setProfileNoticeError] = useState(false);
  const [avatarPreviewOverride, setAvatarPreviewOverride] = useState("");

  useEffect(() => {
    const refresh = () => {
      setAccessToken(readAccessToken());
      setLocalIdentity(readLocalIdentitySnapshot());
    };

    refresh();
    window.addEventListener("storage", refresh);
    window.addEventListener(AUTH_CHANGED_EVENT, refresh);
    return () => {
      window.removeEventListener("storage", refresh);
      window.removeEventListener(AUTH_CHANGED_EVENT, refresh);
    };
  }, []);

  useEffect(() => {
    if (!tokenPair?.accessToken) {
      return;
    }
    setAccessToken(tokenPair.accessToken);
  }, [tokenPair?.accessToken]);

  useEffect(() => {
    return () => {
      if (blobAvatarPreviewRef.current.startsWith("blob:")) {
        URL.revokeObjectURL(blobAvatarPreviewRef.current);
      }
    };
  }, []);

  const isLoggedIn = accessToken.length > 0;

  const profileQuery = useQuery({
    queryKey: ["mall-me-profile"],
    queryFn: () => getMyProfile(true),
    enabled: isLoggedIn,
    staleTime: 60_000,
    retry: 1
  });
  const overviewQuery = useQuery({
    queryKey: ["mall-me-overview"],
    queryFn: getMyOverview,
    enabled: isLoggedIn,
    staleTime: 30_000,
    retry: 1
  });

  const profile = profileQuery.data?.profile;

  useEffect(() => {
    const nextDraft = normalizeDisplayName(profile?.displayName ?? localIdentity.displayName);
    if (nextDraft) {
      setDisplayNameDraft(nextDraft);
    }
  }, [localIdentity.displayName, profile?.displayName]);

  const displayNameMutation = useMutation({
    mutationFn: async (nextDisplayName: string) =>
      updateMyProfile({
        expectedProfileVersion: profile?.profileVersion ?? 1,
        displayName: nextDisplayName
      }),
    onSuccess: (updatedProfile) => {
      queryClient.setQueryData<MallProfileBundle>(["mall-me-profile"], (prev) => ({
        profile: mergeProfile(prev?.profile, updatedProfile),
        addresses: prev?.addresses ?? []
      }));
      syncLocalIdentity(updatedProfile.displayName, updatedProfile.avatar.url || profile?.avatar.url || localIdentity.avatarUrl);
      setProfileNotice(isZh ? "\u6635\u79f0\u5df2\u66f4\u65b0" : "Display name updated");
      setProfileNoticeError(false);
    },
    onError: (error) => {
      console.error("display name update failed", error);
      const errorMessage = error instanceof Error ? error.message : isZh ? "\u6635\u79f0\u66f4\u65b0\u5931\u8d25" : "Failed to update display name";
      setProfileNotice(errorMessage);
      setProfileNoticeError(true);
    }
  });

  const avatarMutation = useMutation({
    mutationFn: async (file: File) => {
      const uploaded = await uploadBuyerAvatar(file);
      const updatedProfile = await updateMyProfile({
        expectedProfileVersion: profile?.profileVersion ?? 1,
        avatar: {
          assetId: uploaded.assetId,
          url: uploaded.url
        }
      });

      let previewUrl = "";
      if (!uploaded.url && typeof URL !== "undefined" && typeof URL.createObjectURL === "function") {
        previewUrl = URL.createObjectURL(file);
      }

      return {
        profile: updatedProfile,
        previewUrl
      };
    },
    onSuccess: (result) => {
      queryClient.setQueryData<MallProfileBundle>(["mall-me-profile"], (prev) => ({
        profile: mergeProfile(prev?.profile, result.profile),
        addresses: prev?.addresses ?? []
      }));

      if (blobAvatarPreviewRef.current.startsWith("blob:")) {
        URL.revokeObjectURL(blobAvatarPreviewRef.current);
        blobAvatarPreviewRef.current = "";
      }

      if (result.previewUrl) {
        blobAvatarPreviewRef.current = result.previewUrl;
        setAvatarPreviewOverride(result.previewUrl);
      } else {
        setAvatarPreviewOverride("");
      }

      syncLocalIdentity(result.profile.displayName, result.profile.avatar.url || result.previewUrl);
      setProfileNotice(isZh ? "\u5934\u50cf\u5df2\u66f4\u65b0" : "Avatar updated");
      setProfileNoticeError(false);
    },
    onError: (error) => {
      console.error("avatar update failed", error);
      setProfileNotice(toFriendlyProfileNotice(error, isZh));
      setProfileNoticeError(true);
    }
  });

  const displayName = useMemo(() => {
    return (
      normalizeDisplayName(profile?.displayName ?? "") ||
      normalizeDisplayName(localIdentity.displayName) ||
      normalizeDisplayName(localIdentity.identifier) ||
      (isZh ? "Shopa \u7528\u6237" : "Shopa User")
    );
  }, [isZh, localIdentity.displayName, localIdentity.identifier, profile?.displayName]);

  const avatarUrl = avatarPreviewOverride || profile?.avatar.url || localIdentity.avatarUrl || "";
  const avatarText = displayName.slice(0, 1).toUpperCase();
  const ext = profile?.ext ?? {};
  const phone = maskPhone(ext.phone);
  const email = ext.email || "-";
  const bio = ext.bio || (isZh ? "\u8fd8\u6ca1\u6709\u586b\u5199\u4e2a\u4eba\u7b80\u4ecb" : "No bio yet");
  const hasDefaultAddress = Boolean(profileQuery.data?.addresses?.some((item) => item.isDefault));
  const addressCount = profileQuery.data?.addresses?.length ?? 0;
  const updatedAt = formatTime(profile?.updatedAt, locale);
  const pointsUnavailable = Boolean(overviewQuery.data?.partial && overviewQuery.data?.degradedFields?.includes("points"));
  const pointsBalance = pointsUnavailable ? null : (overviewQuery.data?.points ?? 0);

  useEffect(() => {
    if (!displayName || typeof window === "undefined") {
      return;
    }
    syncLocalIdentity(displayName, avatarUrl);
  }, [avatarUrl, displayName]);

  function handleDisplayNameSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const normalized = normalizeDisplayName(displayNameDraft);
    if (!normalized) {
      setProfileNotice(isZh ? "\u6635\u79f0\u4e0d\u80fd\u4e3a\u7a7a\uff0c\u4e5f\u4e0d\u80fd\u662f\u7eaf\u6570\u5b57" : "Display name cannot be empty or only numbers");
      setProfileNoticeError(true);
      return;
    }
    if (normalized === displayName) {
      setProfileNotice(isZh ? "\u6635\u79f0\u6ca1\u6709\u53d8\u5316" : "No display name changes");
      setProfileNoticeError(false);
      return;
    }
    displayNameMutation.mutate(normalized);
  }

  function handleAvatarFileChange(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    event.target.value = "";
    if (!file) {
      return;
    }
    if (!file.type.startsWith("image/")) {
      setProfileNotice(isZh ? "\u8bf7\u4e0a\u4f20\u56fe\u7247\u6587\u4ef6" : "Please upload an image file");
      setProfileNoticeError(true);
      return;
    }
    avatarMutation.mutate(file);
  }

  if (!isLoggedIn) {
    return (
      <section className="tb-profile-page">
        <div className="tb-profile-empty" role="status" aria-live="polite">
          <h2>{isZh ? "\u4f60\u8fd8\u6ca1\u6709\u767b\u5f55" : "You're not signed in"}</h2>
          <p>{isZh ? "\u767b\u5f55\u540e\u53ef\u4ee5\u67e5\u770b\u5b8c\u6574\u8d44\u6599\u3001\u6536\u8d27\u5730\u5740\u548c\u8ba2\u5355\u8bb0\u5f55\u3002" : "Sign in to manage your profile, addresses, and orders."}</p>
          <div className="tb-profile-empty-actions">
            <Link href="/login">{isZh ? "\u53bb\u767b\u5f55" : "Go to Login"}</Link>
            <Link href="/register">{isZh ? "\u53bb\u6ce8\u518c" : "Create Account"}</Link>
          </div>
        </div>
      </section>
    );
  }

  return (
    <section className="tb-profile-page">
      <header className="tb-profile-hero">
        <div className="tb-profile-avatar-wrap">
          {avatarUrl ? <img src={avatarUrl} alt={displayName} className="tb-profile-avatar-img" /> : <span>{avatarText}</span>}
        </div>

        <div className="tb-profile-hero-copy">
          <p className="tb-profile-hero-tag">{isZh ? "\u4e2a\u4eba\u4e2d\u5fc3" : "Profile Center"}</p>
          <h2>{displayName}</h2>
          <p>{isZh ? "\u5728\u8fd9\u91cc\u4fee\u6539\u6635\u79f0\u3001\u5934\u50cf\u548c\u8d26\u6237\u57fa\u7840\u4fe1\u606f\u3002" : "Manage your display name, avatar, and account details here."}</p>

          <div className="tb-profile-hero-actions">
            <Link href="/me/orders">{isZh ? "\u67e5\u770b\u8ba2\u5355" : "My Orders"}</Link>
            <Link href="/me/address">{isZh ? "\u5730\u5740\u7ba1\u7406" : "Addresses"}</Link>
            <Link href="/me/points">{isZh ? "\u6211\u7684\u79ef\u5206" : "My Points"}</Link>
            <button type="button" onClick={() => avatarInputRef.current?.click()} disabled={avatarMutation.isPending}>
              {avatarMutation.isPending ? (isZh ? "\u4e0a\u4f20\u5934\u50cf\u4e2d\u2026" : "Uploading avatar\u2026") : isZh ? "\u66f4\u6362\u5934\u50cf" : "Change Avatar"}
            </button>
            <button
              type="button"
              onClick={async () => {
                await logout();
                window.location.href = "/login";
              }}
            >
              {isZh ? "\u9000\u51fa\u767b\u5f55" : "Logout"}
            </button>
          </div>

          <input
            ref={avatarInputRef}
            type="file"
            accept="image/jpeg,image/png,image/webp"
            hidden
            onChange={handleAvatarFileChange}
          />
        </div>
      </header>

      <div className="tb-profile-metrics" aria-label={isZh ? "\u8d26\u6237\u6982\u89c8" : "Account Overview"}>
        <article>
          <strong>{addressCount}</strong>
          <span>{isZh ? "\u6536\u8d27\u5730\u5740" : "Addresses"}</span>
        </article>
        <article>
          <strong>{hasDefaultAddress ? (isZh ? "\u5df2\u8bbe\u7f6e" : "Configured") : isZh ? "\u672a\u8bbe\u7f6e" : "Not set"}</strong>
          <span>{isZh ? "\u9ed8\u8ba4\u5730\u5740" : "Default Address"}</span>
        </article>
        <article>
          <strong>{pointsBalance ?? "--"}</strong>
          <span>{isZh ? "\u5f53\u524d\u79ef\u5206" : "Points"}</span>
        </article>
        <article>
          <strong>{updatedAt}</strong>
          <span>{isZh ? "\u6700\u8fd1\u66f4\u65b0" : "Last Updated"}</span>
        </article>
      </div>

      {pointsUnavailable ? <p className="tb-profile-tip is-error">{isZh ? "积分服务暂时不可用，稍后会自动刷新。" : "Points are temporarily unavailable and will refresh later."}</p> : null}
      <div className="tb-profile-grid">
        <article className="tb-profile-card">
          <header>
            <h3>{isZh ? "\u57fa\u7840\u8d44\u6599" : "Basic Information"}</h3>
            <span className="tb-profile-card-meta">{isZh ? "\u652f\u6301\u4fee\u6539\u6635\u79f0\u4e0e\u5934\u50cf" : "Edit display name and avatar"}</span>
          </header>

          {profileQuery.isLoading ? (
            <p>{isZh ? "\u6b63\u5728\u52a0\u8f7d\u8d44\u6599\u2026" : "Loading profile\u2026"}</p>
          ) : (
            <>
              <form className="tb-profile-edit-form" onSubmit={handleDisplayNameSubmit}>
                <label htmlFor="profile-display-name">
                  {isZh ? "\u6635\u79f0" : "Display Name"}
                  <input
                    id="profile-display-name"
                    type="text"
                    value={displayNameDraft}
                    maxLength={32}
                    onChange={(event) => setDisplayNameDraft(event.target.value)}
                    placeholder={isZh ? "\u8bf7\u8f93\u5165\u6635\u79f0\uff0c\u6700\u591a 32 \u4e2a\u5b57\u7b26" : "Enter a display name up to 32 characters"}
                    disabled={displayNameMutation.isPending}
                  />
                </label>
                <div className="tb-profile-edit-actions">
                  <button type="submit" disabled={displayNameMutation.isPending}>
                    {displayNameMutation.isPending ? (isZh ? "\u4fdd\u5b58\u4e2d\u2026" : "Saving\u2026") : isZh ? "\u4fdd\u5b58\u6635\u79f0" : "Save Name"}
                  </button>
                  <button type="button" onClick={() => setDisplayNameDraft(displayName)} disabled={displayNameMutation.isPending}>
                    {isZh ? "\u91cd\u7f6e" : "Reset"}
                  </button>
                </div>
              </form>

              <dl className="tb-profile-list">
                <div>
                  <dt>{isZh ? "\u5f53\u524d\u6635\u79f0" : "Current Name"}</dt>
                  <dd>{displayName}</dd>
                </div>
                <div>
                  <dt>{isZh ? "\u8054\u7cfb\u7535\u8bdd" : "Phone"}</dt>
                  <dd>{phone}</dd>
                </div>
                <div>
                  <dt>{isZh ? "\u90ae\u7bb1" : "Email"}</dt>
                  <dd>{email}</dd>
                </div>
                <div>
                  <dt>{isZh ? "\u4e2a\u4eba\u7b80\u4ecb" : "Bio"}</dt>
                  <dd>{bio}</dd>
                </div>
              </dl>
            </>
          )}

          {profileNotice ? <p className={`tb-profile-tip ${profileNoticeError ? "is-error" : "is-success"}`}>{profileNotice}</p> : null}
          {profileQuery.isError ? (
            <p className="tb-profile-tip is-error">{isZh ? "\u8d44\u6599\u8bf7\u6c42\u5931\u8d25\uff0c\u5f53\u524d\u663e\u793a\u672c\u5730\u7f13\u5b58\u4fe1\u606f\u3002" : "Profile request failed, showing cached local data."}</p>
          ) : null}
        </article>

        <article className="tb-profile-card">
          <header>
            <h3>{isZh ? "\u8d26\u6237\u4e0e\u5b89\u5168" : "Account & Security"}</h3>
          </header>
          <ul className="tb-profile-security">
            <li>
              <span>{isZh ? "\u767b\u5f55\u5bc6\u7801" : "Password"}</span>
              <Link href="/forgot-password">{isZh ? "\u7acb\u5373\u4fee\u6539" : "Reset now"}</Link>
            </li>
            <li>
              <span>{isZh ? "\u6536\u8d27\u5730\u5740" : "Address Book"}</span>
              <Link href="/me/address">{isZh ? "\u7ba1\u7406\u5730\u5740" : "Manage"}</Link>
            </li>
            <li>
              <span>{isZh ? "\u8ba2\u5355\u4e2d\u5fc3" : "Order Center"}</span>
              <Link href="/me/orders">{isZh ? "\u67e5\u770b\u8ba2\u5355" : "Open"}</Link>
            </li>
            <li>
              <span>{isZh ? "\u79ef\u5206\u4f59\u989d" : "Points Balance"}</span>
              <Link href="/me/points">{isZh ? "\u67e5\u770b\u79ef\u5206" : "View"}</Link>
            </li>
          </ul>
        </article>
      </div>
    </section>
  );
}
