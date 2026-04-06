"use client";

import { useCallback } from "react";
import { login as loginApi, logout as logoutApi, refreshToken as refreshTokenApi } from "./api";
import { useAuthStore } from "./store";

const ACCESS_TOKEN_KEY = "shopa_seller_access_token";
const REFRESH_TOKEN_KEY = "shopa_seller_refresh_token";
const AUTH_STORE_KEY = "shopa-seller-auth";

function persistTokenPair(tokenPair: {
  accessToken: string;
  refreshToken?: string;
  tokenType?: string;
  accessExpiresIn?: number;
  refreshExpiresIn?: number;
  sid: string;
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
  const setTokenPair = useAuthStore((state) => state.setTokenPair);
  const clear = useAuthStore((state) => state.clear);

  const login = useCallback(
    async (identifier: string, password: string) => {
      const response = await loginApi({ identifier, password });
      if (response.auth?.tokenPair) {
        setTokenPair(response.auth.tokenPair);
        persistTokenPair(response.auth.tokenPair);
      }
      return response;
    },
    [setTokenPair]
  );

  const bootstrap = useCallback(async () => {
    const accessToken = localStorage.getItem(ACCESS_TOKEN_KEY);
    const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);

    if (accessToken) {
      if (!tokenPair?.accessToken) {
        setTokenPair({
          tokenType: tokenPair?.tokenType ?? "Bearer",
          accessToken,
          refreshToken: refreshToken || tokenPair?.refreshToken,
          refreshExpiresIn: tokenPair?.refreshExpiresIn,
          accessExpiresIn: tokenPair?.accessExpiresIn,
          sid: tokenPair?.sid ?? "restored"
        });
      }
      return true;
    }

    if (!refreshToken) {
      return false;
    }

    try {
      const nextTokenPair = await refreshTokenApi(refreshToken);
      setTokenPair(nextTokenPair);
      persistTokenPair(nextTokenPair);
      return true;
    } catch {
      localStorage.removeItem(ACCESS_TOKEN_KEY);
      localStorage.removeItem(REFRESH_TOKEN_KEY);
      localStorage.removeItem(AUTH_STORE_KEY);
      clear();
      return false;
    }
  }, [clear, setTokenPair, tokenPair?.accessExpiresIn, tokenPair?.accessToken, tokenPair?.refreshExpiresIn, tokenPair?.refreshToken, tokenPair?.sid, tokenPair?.tokenType]);

  const logout = useCallback(async () => {
    const refresh = tokenPair?.refreshToken ?? localStorage.getItem(REFRESH_TOKEN_KEY) ?? "";
    if (refresh) {
      try {
        await logoutApi(refresh);
      } catch (error) {
        console.warn("logout api failed, fallback to local clear", error);
      }
    }
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    localStorage.removeItem(AUTH_STORE_KEY);
    clear();
  }, [clear, tokenPair?.refreshToken]);

  return { tokenPair, login, bootstrap, logout };
}
