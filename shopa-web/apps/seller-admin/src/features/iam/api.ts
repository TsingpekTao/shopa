import { apiClient } from "@/lib/api-client";
import {
  LoginRequest,
  LoginResponse,
  ResetPasswordBySmsRequest,
  ResetPasswordBySmsResponse,
  RefreshTokenResponse,
  RegisterByPasswordRequest,
  RegisterByPasswordResponse,
  SendSmsCodeRequest,
  SendSmsCodeResponse,
  TokenPair
} from "./types";

type LoginPayload = LoginRequest & {
  risk?: unknown;
};

export async function login(payload: LoginPayload): Promise<LoginResponse> {
  const body = new URLSearchParams();
  body.set("identifier", payload.identifier);
  body.set("password", payload.password);
  if (payload.risk !== undefined) {
    body.set("risk", JSON.stringify(payload.risk));
  }
  return apiClient.post<LoginResponse>("/v1/auth/login/password", body.toString(), {
    headers: {
      "Content-Type": "application/x-www-form-urlencoded"
    }
  });
}

export async function refreshToken(refreshToken: string, risk?: unknown): Promise<TokenPair> {
  const response = await apiClient.post<RefreshTokenResponse>("/v1/auth/token/refresh", {
    refreshToken,
    risk
  });
  if (!response?.tokenPair) {
    throw new Error("refresh token response missing tokenPair");
  }
  return response.tokenPair;
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

export async function sendSmsCode(payload: SendSmsCodeRequest): Promise<SendSmsCodeResponse> {
  return apiClient.post<SendSmsCodeResponse>("/v1/auth/sms/send", payload);
}

export async function registerByPassword(
  payload: RegisterByPasswordRequest
): Promise<RegisterByPasswordResponse> {
  return apiClient.post<RegisterByPasswordResponse>("/v1/auth/register/password", payload);
}

export async function resetPasswordBySms(
  payload: ResetPasswordBySmsRequest
): Promise<ResetPasswordBySmsResponse> {
  return apiClient.post<ResetPasswordBySmsResponse>("/v1/auth/password/reset/sms", payload);
}
