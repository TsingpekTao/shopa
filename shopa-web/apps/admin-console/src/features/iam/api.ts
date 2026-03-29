import { apiClient } from "@/lib/api-client";
import { AdminOverviewResponse, LoginRequest, LoginResponse } from "./types";

export async function login(payload: LoginRequest): Promise<LoginResponse> {
  const body = new URLSearchParams();
  body.set("identifier", payload.identifier);
  body.set("password", payload.password);
  return apiClient.post<LoginResponse>("/v1/auth/login/password", body.toString(), {
    headers: {
      "Content-Type": "application/x-www-form-urlencoded"
    }
  });
}

export async function fetchAdminOverview(): Promise<AdminOverviewResponse> {
  return apiClient.get<AdminOverviewResponse>("/v1/admin/me/overview");
}

export async function logout(refreshToken: string, allDevices = false): Promise<void> {
  const body = new URLSearchParams();
  body.set("refreshToken", refreshToken);
  body.set("allDevices", allDevices ? "true" : "false");
  await apiClient.post("/v1/auth/logout", body.toString(), {
    headers: {
      "Content-Type": "application/x-www-form-urlencoded"
    }
  });
}
