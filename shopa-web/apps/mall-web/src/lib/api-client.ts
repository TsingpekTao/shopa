"use client";

import { createApiClient } from "@shopa/api-client";
import { useAuthStore } from "@/features/iam/store";

const AUTH_CHANGED_EVENT = "shopa-mall-auth-changed";

type PersistedAuthState = {
  state?: {
    tokenPair?: {
      accessToken?: string;
      refreshToken?: string;
      tokenType?: string;
      accessExpiresIn?: number;
      refreshExpiresIn?: number;
      sid?: string;
    };
  };
};

type RefreshTokenResponse = {
  code: number;
  message?: string;
  data?: {
    tokenPair?: {
      accessToken?: string;
      refreshToken?: string;
      tokenType?: string;
      accessExpiresIn?: number;
      refreshExpiresIn?: number;
      sid?: string;
    };
  };
};

let refreshPromise: Promise<string | null> | null = null;

function readTokenFromPersistStore(): string {
  try {
    const raw = localStorage.getItem("shopa-mall-auth");
    if (!raw) {
      return "";
    }
    const parsed = JSON.parse(raw) as PersistedAuthState;
    return parsed?.state?.tokenPair?.accessToken?.trim() ?? "";
  } catch {
    return "";
  }
}

function readRefreshTokenFromPersistStore(): string {
  try {
    const raw = localStorage.getItem("shopa-mall-auth");
    if (!raw) {
      return "";
    }
    const parsed = JSON.parse(raw) as PersistedAuthState;
    return parsed?.state?.tokenPair?.refreshToken?.trim() ?? "";
  } catch {
    return "";
  }
}

function getApiBaseURL(): string {
  const envBaseURL = process.env.NEXT_PUBLIC_API_BASE_URL?.trim();
  if (envBaseURL) {
    return envBaseURL;
  }
  const { protocol, hostname } = window.location;
  return `${protocol}//${hostname}:8000`;
}

function emitAuthChanged() {
  window.dispatchEvent(new Event(AUTH_CHANGED_EVENT));
}

function clearLocalAuth(shouldRedirect = true) {
  localStorage.removeItem("shopa_mall_access_token");
  localStorage.removeItem("shopa_mall_refresh_token");
  localStorage.removeItem("shopa_mall_last_identifier");
  localStorage.removeItem("shopa_mall_display_name");
  localStorage.removeItem("shopa_mall_user_id");
  localStorage.removeItem("shopa-mall-auth");
  useAuthStore.getState().clear();
  emitAuthChanged();

  if (!shouldRedirect) {
    return;
  }
  const currentPath = `${window.location.pathname}${window.location.search}${window.location.hash}`;
  const isLoginPage = window.location.pathname.startsWith("/login");
  if (isLoginPage) {
    return;
  }
  const next = new URL("/login", window.location.origin);
  if (currentPath && currentPath !== "/") {
    next.searchParams.set("redirect", currentPath);
  }
  window.location.href = next.toString();
}

async function refreshAccessToken(): Promise<string | null> {
  const refreshToken =
    localStorage.getItem("shopa_mall_refresh_token")?.trim() ?? readRefreshTokenFromPersistStore();
  if (!refreshToken) {
    clearLocalAuth();
    return null;
  }

  const response = await fetch(`${getApiBaseURL()}/v1/auth/token/refresh`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({
      refreshToken
    })
  });

  const payload = (await response.json().catch(() => null)) as RefreshTokenResponse | null;
  const nextTokenPair = payload?.data?.tokenPair;
  const nextAccessToken = nextTokenPair?.accessToken?.trim() ?? "";
  if (!response.ok || payload?.code !== 0 || !nextAccessToken) {
    clearLocalAuth();
    return null;
  }

  useAuthStore.getState().setTokenPair({
    accessToken: nextAccessToken,
    refreshToken: nextTokenPair?.refreshToken?.trim() || refreshToken,
    tokenType: nextTokenPair?.tokenType?.trim() || "Bearer",
    accessExpiresIn: nextTokenPair?.accessExpiresIn ?? 0,
    refreshExpiresIn: nextTokenPair?.refreshExpiresIn ?? 0,
    sid: nextTokenPair?.sid?.trim() ?? ""
  });
  emitAuthChanged();
  return nextAccessToken;
}

export const apiClient = createApiClient({
  hooks: {
    getToken: () => {
      const direct = localStorage.getItem("shopa_mall_access_token")?.trim();
      if (direct) {
        return direct;
      }
      const fromStore = readTokenFromPersistStore();
      if (fromStore) {
        localStorage.setItem("shopa_mall_access_token", fromStore);
      }
      return fromStore;
    },
    onUnauthorized: async (error) => {
      const requestUrl = error.config?.url ?? "";
      if (requestUrl.includes("/v1/auth/token/refresh")) {
        clearLocalAuth();
        return null;
      }

      const refreshToken =
        localStorage.getItem("shopa_mall_refresh_token")?.trim() ?? readRefreshTokenFromPersistStore();
      if (!refreshToken) {
        clearLocalAuth();
        return null;
      }

      if (!refreshPromise) {
        refreshPromise = refreshAccessToken().finally(() => {
          refreshPromise = null;
        });
      }
      const nextAccessToken = await refreshPromise;
      if (!nextAccessToken) {
        clearLocalAuth(false);
      }
      return nextAccessToken;
    }
  }
});
