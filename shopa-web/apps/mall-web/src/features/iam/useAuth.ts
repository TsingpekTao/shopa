"use client";

import { useCallback } from "react";
import { login as loginApi, logout as logoutApi } from "./api";
import { useAuthStore } from "./store";

const AUTH_CHANGED_EVENT = "shopa-mall-auth-changed";

function emitAuthChanged() {
  if (typeof window !== "undefined") {
    window.dispatchEvent(new Event(AUTH_CHANGED_EVENT));
  }
}

export function useAuth() {
  const tokenPair = useAuthStore((state) => state.tokenPair);
  const setTokenPair = useAuthStore((state) => state.setTokenPair);
  const clear = useAuthStore((state) => state.clear);

  const login = useCallback(
    async (identifier: string, password: string) => {
      const response = await loginApi({ identifier, password });
      if (response.auth?.tokenPair) {
        setTokenPair(response.auth.tokenPair);
        localStorage.setItem("shopa_mall_access_token", response.auth.tokenPair.accessToken);
        localStorage.setItem("shopa_mall_last_identifier", identifier);
        localStorage.setItem("shopa_mall_display_name", identifier);
        if (response.session?.userId) {
          localStorage.setItem("shopa_mall_user_id", String(response.session.userId));
        }
        if (response.auth.tokenPair.refreshToken) {
          localStorage.setItem("shopa_mall_refresh_token", response.auth.tokenPair.refreshToken);
        } else {
          localStorage.removeItem("shopa_mall_refresh_token");
        }
        emitAuthChanged();
      }
      return response;
    },
    [setTokenPair]
  );

  const logout = useCallback(async () => {
    const refresh = tokenPair?.refreshToken ?? "";
    try {
      if (refresh) {
        await logoutApi(refresh);
      }
    } catch (_err) {
      // Ignore remote logout failure; always clear local session.
    }
    localStorage.removeItem("shopa_mall_access_token");
    localStorage.removeItem("shopa_mall_refresh_token");
    localStorage.removeItem("shopa_mall_last_identifier");
    localStorage.removeItem("shopa_mall_display_name");
    localStorage.removeItem("shopa_mall_user_id");
    localStorage.removeItem("shopa-mall-auth");
    clear();
    emitAuthChanged();
  }, [clear, tokenPair?.refreshToken]);

  return { tokenPair, login, logout };
}
