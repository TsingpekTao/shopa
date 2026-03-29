"use client";

import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";
import { TokenPair } from "./types";

type AuthState = {
  tokenPair?: TokenPair;
  userId?: number;
  accountStatusCode?: string;
  roles: string[];
  permissions: string[];
  setTokenPair: (tokenPair: TokenPair) => void;
  setOverview: (payload: { userId: number; accountStatusCode: string; roles: string[]; permissions: string[] }) => void;
  clear: () => void;
};

const storage = typeof window !== "undefined" ? createJSONStorage(() => localStorage) : undefined;

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      tokenPair: undefined,
      userId: undefined,
      accountStatusCode: undefined,
      roles: [],
      permissions: [],
      setTokenPair: (tokenPair) => set({ tokenPair }),
      setOverview: ({ userId, accountStatusCode, roles, permissions }) =>
        set({
          userId,
          accountStatusCode,
          roles: roles || [],
          permissions: permissions || []
        }),
      clear: () =>
        set({
          tokenPair: undefined,
          userId: undefined,
          accountStatusCode: undefined,
          roles: [],
          permissions: []
        })
    }),
    {
      name: "shopa-admin-auth",
      storage
    }
  )
);
