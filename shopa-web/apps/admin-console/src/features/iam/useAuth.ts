"use client";

import { useCallback } from "react";
import { fetchAdminOverview, login as loginApi, logout as logoutApi } from "./api";
import { useAuthStore } from "./store";

const ACCESS_TOKEN_KEY = "shopa_admin_access_token";
const REFRESH_TOKEN_KEY = "shopa_admin_refresh_token";

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
      localStorage.setItem(ACCESS_TOKEN_KEY, response.auth.tokenPair.accessToken);
      if (response.auth.tokenPair.refreshToken) {
        localStorage.setItem(REFRESH_TOKEN_KEY, response.auth.tokenPair.refreshToken);
      } else {
        localStorage.removeItem(REFRESH_TOKEN_KEY);
      }
      await loadOverview();
      return response;
    },
    [loadOverview, setTokenPair]
  );

  const bootstrap = useCallback(async () => {
    const token = localStorage.getItem(ACCESS_TOKEN_KEY);
    if (!token) {
      return false;
    }
    try {
      await loadOverview();
      return true;
    } catch {
      return false;
    }
  }, [loadOverview]);

  const logout = useCallback(async () => {
    const refresh = tokenPair?.refreshToken ?? localStorage.getItem(REFRESH_TOKEN_KEY) ?? "";
    try {
      if (refresh) {
        await logoutApi(refresh);
      }
    } finally {
      localStorage.removeItem(ACCESS_TOKEN_KEY);
      localStorage.removeItem(REFRESH_TOKEN_KEY);
      localStorage.removeItem("shopa-admin-auth");
      clear();
    }
  }, [clear, tokenPair?.refreshToken]);

  return { tokenPair, permissions, login, bootstrap, logout, loadOverview };
}
