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

export interface LoginResponse {
  channel: number;
  auth: {
    tokenPair?: TokenPair;
    mfaChallenge?: unknown;
  };
}

export interface RefreshTokenResponse {
  tokenPair: TokenPair;
}

export interface AdminOverviewResponse {
  userId: number;
  accountStatusCode: string;
  roles: string[];
  permissions: string[];
}
