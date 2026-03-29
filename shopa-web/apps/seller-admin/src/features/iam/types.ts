export interface LoginRequest {
  identifier: string;
  password: string;
}

export interface TokenPair {
  tokenType: string;
  accessToken: string;
  refreshToken?: string;
  refreshExpiresIn?: number;
  accessExpiresIn?: number;
  sid: string;
}

export interface RoleItem {
  roleCode: number;
  scopeType: number;
  scopeId: number;
}

export interface SessionSummary {
  userId: number;
  accountStatus: number;
  roles: RoleItem[];
  membership?: { levelCode: string; points: number; expireAt?: string };
  lastLoginAt?: string;
  lastLoginIp?: string;
}

export interface LoginResponse {
  channel: number;
  auth: {
    tokenPair?: TokenPair;
  };
  session?: SessionSummary;
}

export interface RefreshTokenResponse {
  tokenPair: TokenPair;
}

export interface SendSmsCodeRequest {
  scene: number;
  phone: string;
}

export interface SendSmsCodeResponse {
  resendAfterSeconds: number;
}

export interface RegisterByPasswordRequest {
  phone: string;
  smsCode: string;
  password: string;
  confirmPassword: string;
}

export interface RegisterByPasswordResponse {
  userId: number;
  initDisplayName: string;
  auth?: {
    tokenPair?: TokenPair;
  };
  session?: SessionSummary;
}

export interface ResetPasswordBySmsRequest {
  phone: string;
  smsCode: string;
  newPassword: string;
}

export interface ResetPasswordBySmsResponse {
  updated: boolean;
}
