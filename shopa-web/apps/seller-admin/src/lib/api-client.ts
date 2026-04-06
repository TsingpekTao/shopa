"use client";

import { createApiClient } from "@shopa/api-client";
import { useAuthStore } from "@/features/iam/store";

const ACCESS_TOKEN_KEY = "shopa_seller_access_token";
const REFRESH_TOKEN_KEY = "shopa_seller_refresh_token";
const AUTH_STORE_KEY = "shopa-seller-auth";

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

function readPersistedRefreshToken(): string {
  try {
    const raw = localStorage.getItem(AUTH_STORE_KEY);
    if (!raw) {
      return "";
    }
    const parsed = JSON.parse(raw) as PersistedAuthState;
    return parsed?.state?.tokenPair?.refreshToken?.trim() ?? "";
  } catch {
    return "";
  }
}

function clearLocalAuth(shouldRedirect = true) {
  localStorage.removeItem(ACCESS_TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
  localStorage.removeItem(AUTH_STORE_KEY);
  useAuthStore.getState().clear();

  if (!shouldRedirect || window.location.pathname.startsWith("/login")) {
    return;
  }
  const next = new URL("/login", window.location.origin);
  const currentPath = `${window.location.pathname}${window.location.search}${window.location.hash}`;
  if (currentPath && currentPath !== "/") {
    next.searchParams.set("redirect", currentPath);
  }
  window.location.href = next.toString();
}

async function refreshAccessToken(): Promise<string | null> {
  const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY)?.trim() || readPersistedRefreshToken();
  if (!refreshToken) {
    clearLocalAuth();
    return null;
  }

  const baseURL = process.env.NEXT_PUBLIC_API_BASE_URL?.trim() || `${window.location.protocol}//${window.location.hostname}:8000`;
  const response = await fetch(`${baseURL}/v1/auth/token/refresh`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify({ refreshToken })
  });

  const payload = (await response.json().catch(() => null)) as RefreshTokenResponse | null;
  const tokenPair = payload?.data?.tokenPair;
  const accessToken = tokenPair?.accessToken?.trim() ?? "";
  if (!response.ok || payload?.code !== 0 || !accessToken) {
    clearLocalAuth();
    return null;
  }

  useAuthStore.getState().setTokenPair({
    accessToken,
    refreshToken: tokenPair?.refreshToken?.trim() || refreshToken,
    tokenType: tokenPair?.tokenType?.trim() || "Bearer",
    accessExpiresIn: tokenPair?.accessExpiresIn ?? 0,
    refreshExpiresIn: tokenPair?.refreshExpiresIn ?? 0,
    sid: tokenPair?.sid?.trim() || "refreshed"
  });
  localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
  if (tokenPair?.refreshToken?.trim()) {
    localStorage.setItem(REFRESH_TOKEN_KEY, tokenPair.refreshToken.trim());
  }
  return accessToken;
}

export const apiClient = createApiClient({
  hooks: {
    getToken: () => localStorage.getItem(ACCESS_TOKEN_KEY),
    onUnauthorized: async (error) => {
      const requestUrl = error.config?.url ?? "";
      if (requestUrl.includes("/v1/auth/token/refresh")) {
        clearLocalAuth();
        return null;
      }

      const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY)?.trim() || readPersistedRefreshToken();
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
