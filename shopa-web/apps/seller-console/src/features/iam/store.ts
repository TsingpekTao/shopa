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
      setTokenPair: (tokenPair) => set({ tokenPair }),
      clear: () => set({ tokenPair: undefined })
    }),
    {
      name: "shopa-auth",
      storage
    }
  )
);
