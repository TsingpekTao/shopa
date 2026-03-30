"use client";

import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";
import { TokenPair } from "./types";

type AuthState = {
  tokenPair?: TokenPair;
  setTokenPair: (tokenPair: TokenPair) => void;
  clear: () => void;
};

const storage = typeof window !== "undefined" ? createJSONStorage(() => localStorage) : undefined;

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      tokenPair: undefined,
      setTokenPair: (tokenPair) => {
        if (typeof window !== "undefined") {
          localStorage.setItem("shopa_mall_access_token", tokenPair.accessToken);
          if (tokenPair.refreshToken) {
            localStorage.setItem("shopa_mall_refresh_token", tokenPair.refreshToken);
          }
        }
        set({ tokenPair });
      },
      clear: () => {
        if (typeof window !== "undefined") {
          localStorage.removeItem("shopa_mall_access_token");
          localStorage.removeItem("shopa_mall_refresh_token");
        }
        set({ tokenPair: undefined });
      }
    }),
    {
      name: "shopa-mall-auth",
      storage,
      onRehydrateStorage: () => (state) => {
        if (typeof window === "undefined") {
          return;
        }
        const accessToken = state?.tokenPair?.accessToken?.trim() ?? "";
        const refreshToken = state?.tokenPair?.refreshToken?.trim() ?? "";
        if (accessToken) {
          localStorage.setItem("shopa_mall_access_token", accessToken);
        }
        if (refreshToken) {
          localStorage.setItem("shopa_mall_refresh_token", refreshToken);
        }
      }
    }
  )
);
