"use client";

import Link from "next/link";
import { useEffect, useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { useI18n } from "@shopa/ui";
import { useAuth } from "@/features/iam/useAuth";
import { apiClient } from "@/lib/api-client";

type ProfileImage = {
  url?: string;
};

type RawProfile = {
  display_name?: string;
  displayName?: string;
  updated_at?: string | { seconds?: number | string; nanos?: number | string };
  updatedAt?: string | { seconds?: number | string; nanos?: number | string };
  avatar?: ProfileImage;
  ext?: Record<string, string>;
};

type RawAddress = {
  is_default?: boolean;
};

type RawProfileResponse = {
  profile?: RawProfile;
  addresses?: RawAddress[];
};

type LocalIdentitySnapshot = {
  displayName: string;
  identifier: string;
};

const AUTH_CHANGED_EVENT = "shopa-mall-auth-changed";

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
    return { displayName: "", identifier: "" };
  }
  return {
    displayName: window.localStorage.getItem("shopa_mall_display_name") ?? "",
    identifier: window.localStorage.getItem("shopa_mall_last_identifier") ?? ""
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

export default function ProfilePage() {
  const { locale } = useI18n();
  const isZh = locale === "zh-CN";
  const { tokenPair, logout } = useAuth();
  const [accessToken, setAccessToken] = useState<string>("");
  const [localIdentity, setLocalIdentity] = useState<LocalIdentitySnapshot>({
    displayName: "",
    identifier: ""
  });

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

  const isLoggedIn = accessToken.length > 0;

  const profileQuery = useQuery({
    queryKey: ["mall-me-profile"],
    queryFn: async () => apiClient.get<RawProfileResponse>("/v1/me/profile", { params: { include_addresses: true } }),
    enabled: isLoggedIn,
    staleTime: 60_000,
    retry: 1
  });

  const profile = profileQuery.data?.profile;
  const displayName = useMemo(() => {
    return (
      normalizeDisplayName(profile?.display_name ?? profile?.displayName ?? "") ||
      normalizeDisplayName(localIdentity.displayName) ||
      normalizeDisplayName(localIdentity.identifier) ||
      (isZh ? "Shopa 用户" : "Shopa User")
    );
  }, [isZh, localIdentity.displayName, localIdentity.identifier, profile?.displayName, profile?.display_name]);

  const avatarUrl = profile?.avatar?.url ?? "";
  const avatarText = displayName.slice(0, 1).toUpperCase();
  const ext = profile?.ext ?? {};
  const phone = maskPhone(ext.phone);
  const email = ext.email || "-";
  const bio = ext.bio || (isZh ? "暂未填写个人简介" : "No bio yet");
  const hasDefaultAddress = Boolean(profileQuery.data?.addresses?.some((item) => item.is_default));
  const addressCount = profileQuery.data?.addresses?.length ?? 0;
  const updatedAt = formatTime(profile?.updated_at ?? profile?.updatedAt, locale);

  useEffect(() => {
    if (!displayName || typeof window === "undefined") {
      return;
    }
    window.localStorage.setItem("shopa_mall_display_name", displayName);
    window.dispatchEvent(new Event(AUTH_CHANGED_EVENT));
  }, [displayName]);

  if (!isLoggedIn) {
    return (
      <section className="tb-profile-page">
        <div className="tb-profile-empty" role="status" aria-live="polite">
          <h2>{isZh ? "你还没有登录" : "You're not signed in"}</h2>
          <p>{isZh ? "登录后可查看完整个人资料、地址和订单信息。" : "Sign in to access your full profile, addresses, and orders."}</p>
          <div className="tb-profile-empty-actions">
            <Link href="/login">{isZh ? "去登录" : "Go to Login"}</Link>
            <Link href="/register">{isZh ? "去注册" : "Create Account"}</Link>
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
          <p className="tb-profile-hero-tag">{isZh ? "个人中心" : "Profile Center"}</p>
          <h2>{displayName}</h2>
          <p>{isZh ? "管理账号资料、收货地址与订单偏好" : "Manage account details, addresses and purchase preferences."}</p>
          <div className="tb-profile-hero-actions">
            <Link href="/me/orders">{isZh ? "查看订单" : "My Orders"}</Link>
            <Link href="/me/address">{isZh ? "地址管理" : "Addresses"}</Link>
            <button
              type="button"
              onClick={async () => {
                await logout();
                window.location.href = "/login";
              }}
            >
              {isZh ? "退出登录" : "Logout"}
            </button>
          </div>
        </div>
      </header>

      <div className="tb-profile-metrics" aria-label={isZh ? "账户概览" : "Account Overview"}>
        <article>
          <strong>{addressCount}</strong>
          <span>{isZh ? "收货地址" : "Addresses"}</span>
        </article>
        <article>
          <strong>{hasDefaultAddress ? (isZh ? "已设置" : "Configured") : isZh ? "未设置" : "Not set"}</strong>
          <span>{isZh ? "默认地址" : "Default Address"}</span>
        </article>
        <article>
          <strong>{updatedAt}</strong>
          <span>{isZh ? "最近更新" : "Last Updated"}</span>
        </article>
      </div>

      <div className="tb-profile-grid">
        <article className="tb-profile-card">
          <header>
            <h3>{isZh ? "基础资料" : "Basic Information"}</h3>
            <Link href="/settings">{isZh ? "编辑资料" : "Edit"}</Link>
          </header>
          {profileQuery.isLoading ? (
            <p>{isZh ? "正在加载资料..." : "Loading profile..."}</p>
          ) : (
            <dl className="tb-profile-list">
              <div>
                <dt>{isZh ? "用户名" : "Username"}</dt>
                <dd>{displayName}</dd>
              </div>
              <div>
                <dt>{isZh ? "联系电话" : "Phone"}</dt>
                <dd>{phone}</dd>
              </div>
              <div>
                <dt>{isZh ? "邮箱" : "Email"}</dt>
                <dd>{email}</dd>
              </div>
              <div>
                <dt>{isZh ? "个人简介" : "Bio"}</dt>
                <dd>{bio}</dd>
              </div>
            </dl>
          )}
          {profileQuery.isError && <p className="tb-profile-tip">{isZh ? "资料加载失败，已显示本地信息。" : "Profile request failed, showing local data."}</p>}
        </article>

        <article className="tb-profile-card">
          <header>
            <h3>{isZh ? "账号与安全" : "Account & Security"}</h3>
          </header>
          <ul className="tb-profile-security">
            <li>
              <span>{isZh ? "登录密码" : "Password"}</span>
              <Link href="/forgot-password">{isZh ? "立即修改" : "Reset now"}</Link>
            </li>
            <li>
              <span>{isZh ? "收货地址" : "Address Book"}</span>
              <Link href="/me/address">{isZh ? "管理地址" : "Manage"}</Link>
            </li>
            <li>
              <span>{isZh ? "订单中心" : "Order Center"}</span>
              <Link href="/me/orders">{isZh ? "查看订单" : "Open"}</Link>
            </li>
          </ul>
        </article>
      </div>
    </section>
  );
}
