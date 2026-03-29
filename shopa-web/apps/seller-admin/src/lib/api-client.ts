"use client";

import { createApiClient } from "@shopa/api-client";

export const apiClient = createApiClient({
  hooks: {
    getToken: () => localStorage.getItem("shopa_seller_access_token"),
    onUnauthorized: () => {
      localStorage.removeItem("shopa_seller_access_token");
      localStorage.removeItem("shopa_seller_refresh_token");
      localStorage.removeItem("shopa-seller-auth");
      window.location.href = "/login";
    }
  }
});
