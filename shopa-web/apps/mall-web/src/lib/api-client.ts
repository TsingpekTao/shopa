"use client";

import { createApiClient } from "@shopa/api-client";

export const apiClient = createApiClient({
  hooks: {
    getToken: () => localStorage.getItem("shopa_mall_access_token"),
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
