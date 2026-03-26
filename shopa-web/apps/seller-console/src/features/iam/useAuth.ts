"use client";

import { useCallback } from "react";
import { login as loginApi, logout as logoutApi } from "./api";
import { useAuthStore } from "./store";

export function useAuth() {
  const tokenPair = useAuthStore((state) => state.tokenPair);
  const setTokenPair = useAuthStore((state) => state.setTokenPair);
  const clear = useAuthStore((state) => state.clear);

  const login = useCallback(
    async (identifier: string, password: string) => {
      const response = await loginApi({ identifier, password });
      if (response.auth?.tokenPair) {
        setTokenPair(response.auth.tokenPair);
        localStorage.setItem("shopa_access_token", response.auth.tokenPair.accessToken);
        if (response.auth.tokenPair.refreshToken) {
          localStorage.setItem("shopa_refresh_token", response.auth.tokenPair.refreshToken);
        } else {
          localStorage.removeItem("shopa_refresh_token");
        }
      }
      return response;
    },
    [setTokenPair]
  );

  const logout = useCallback(async () => {
    const refresh = tokenPair?.refreshToken ?? "";
    if (refresh) {
      await logoutApi(refresh);
    }
    localStorage.removeItem("shopa_access_token");
    localStorage.removeItem("shopa_refresh_token");
    localStorage.removeItem("shopa-auth");
    clear();
  }, [clear, tokenPair?.refreshToken]);

  return { tokenPair, login, logout };
}
