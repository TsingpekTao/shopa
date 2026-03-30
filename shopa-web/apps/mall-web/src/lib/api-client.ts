"use client";

import { createApiClient } from "@shopa/api-client";

function readTokenFromPersistStore(): string {
  try {
    const raw = localStorage.getItem("shopa-mall-auth");
    if (!raw) {
      return "";
    }
    const parsed = JSON.parse(raw) as {
      state?: { tokenPair?: { accessToken?: string } };
    };
    return parsed?.state?.tokenPair?.accessToken?.trim() ?? "";
  } catch {
    return "";
  }
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
    onUnauthorized: () => {
      localStorage.removeItem("shopa_mall_access_token");
      localStorage.removeItem("shopa_mall_refresh_token");
      localStorage.removeItem("shopa_mall_last_identifier");
      localStorage.removeItem("shopa_mall_display_name");
      localStorage.removeItem("shopa_mall_user_id");
      localStorage.removeItem("shopa-mall-auth");
      window.dispatchEvent(new Event("shopa-mall-auth-changed"));
      window.location.href = "/login";
    }
  }
});
