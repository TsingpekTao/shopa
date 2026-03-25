import { apiClient } from "@shopa/api-client";
import { LoginRequest, LoginResponse } from "./types";

export async function login(payload: LoginRequest): Promise<LoginResponse> {
  return apiClient.post<LoginResponse>("/v1/auth/login/password", payload);
}

export async function refreshToken(refreshToken: string) {
  return apiClient.post<{ accessToken: string; refreshToken?: string }>("/v1/auth/token/refresh", {
    refreshToken
  });
}

export async function logout(refreshToken: string) {
  return apiClient.post<{ success: boolean }>("/v1/auth/logout", {
    refreshToken,
    allDevices: true
  });
}
