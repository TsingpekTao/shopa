"use client";

import { createApiClient } from "@shopa/api-client";

const ACCESS_TOKEN_KEY = "shopa_admin_access_token";
const REFRESH_TOKEN_KEY = "shopa_admin_refresh_token";
const AUTH_STORE_KEY = "shopa-admin-auth";

export const apiClient = createApiClient({
  baseURL: process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://127.0.0.1:8000",
  hooks: {
    getToken: () => localStorage.getItem(ACCESS_TOKEN_KEY),
    onUnauthorized: () => {
      localStorage.removeItem(ACCESS_TOKEN_KEY);
      localStorage.removeItem(AUTH_STORE_KEY);
      window.location.href = "/login";
    }
  }
});
