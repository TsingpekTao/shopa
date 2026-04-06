"use client";

import { useCallback } from "react";
import { fetchAdminOverview, login as loginApi, logout as logoutApi, refreshToken as refreshTokenApi } from "./api";
import { useAuthStore } from "./store";

const ACCESS_TOKEN_KEY = "shopa_admin_access_token";
const REFRESH_TOKEN_KEY = "shopa_admin_refresh_token";
const AUTH_STORE_KEY = "shopa-admin-auth";

function persistTokenPair(tokenPair: {
  accessToken: string;
  refreshToken?: string;
}) {
  localStorage.setItem(ACCESS_TOKEN_KEY, tokenPair.accessToken);
  if (tokenPair.refreshToken) {
    localStorage.setItem(REFRESH_TOKEN_KEY, tokenPair.refreshToken);
  } else {
    localStorage.removeItem(REFRESH_TOKEN_KEY);
  }
}

export function useAuth() {
  const tokenPair = useAuthStore((state) => state.tokenPair);
  const permissions = useAuthStore((state) => state.permissions);
  const setTokenPair = useAuthStore((state) => state.setTokenPair);
  const setOverview = useAuthStore((state) => state.setOverview);
  const clear = useAuthStore((state) => state.clear);

  const loadOverview = useCallback(async () => {
    const overview = await fetchAdminOverview();
    setOverview({
      userId: overview.userId,
      accountStatusCode: overview.accountStatusCode,
      roles: overview.roles || [],
      permissions: overview.permissions || []
    });
    return overview;
  }, [setOverview]);

  const login = useCallback(
    async (identifier: string, password: string) => {
      const response = await loginApi({ identifier, password });
      if (response.auth?.mfaChallenge) {
        throw new Error("MFA_REQUIRED");
      }
      if (!response.auth?.tokenPair) {
        throw new Error("missing token pair");
      }
      setTokenPair(response.auth.tokenPair);
      persistTokenPair(response.auth.tokenPair);
      await loadOverview();
      return response;
    },
    [loadOverview, setTokenPair]
  );

  const bootstrap = useCallback(async () => {
    const token = localStorage.getItem(ACCESS_TOKEN_KEY);
    const refresh = localStorage.getItem(REFRESH_TOKEN_KEY);

    if (token) {
      try {
        await loadOverview();
        return true;
      } catch {
        localStorage.removeItem(ACCESS_TOKEN_KEY);
      }
    }

    if (!refresh) {
      return false;
    }

    try {
      const nextTokenPair = await refreshTokenApi(refresh);
      setTokenPair(nextTokenPair);
      persistTokenPair(nextTokenPair);
      await loadOverview();
      return true;
    } catch {
      localStorage.removeItem(ACCESS_TOKEN_KEY);
      localStorage.removeItem(REFRESH_TOKEN_KEY);
      localStorage.removeItem(AUTH_STORE_KEY);
      clear();
      return false;
    }
  }, [clear, loadOverview, setTokenPair]);

  const logout = useCallback(async () => {
    const refresh = tokenPair?.refreshToken ?? localStorage.getItem(REFRESH_TOKEN_KEY) ?? "";
    try {
      if (refresh) {
        await logoutApi(refresh);
      }
    } finally {
      localStorage.removeItem(ACCESS_TOKEN_KEY);
      localStorage.removeItem(REFRESH_TOKEN_KEY);
      localStorage.removeItem(AUTH_STORE_KEY);
      clear();
    }
  }, [clear, tokenPair?.refreshToken]);

  return { tokenPair, permissions, login, bootstrap, logout, loadOverview };
}
