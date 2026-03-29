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
  return apiClient.post<LoginResponse>("/v1/auth/login/password", payload);
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
  await apiClient.post("/v1/auth/logout", {
    refreshToken,
    allDevices
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
