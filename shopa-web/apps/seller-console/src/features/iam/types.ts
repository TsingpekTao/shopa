export interface LoginRequest {
  identifier: string;
  password: string;
}

export interface TokenPair {
  tokenType: string;
  accessToken: string;
  refreshToken?: string;
  sid: string;
}

export interface LoginResponse {
  userId: number;
  initDisplayName: string;
  auth: { tokenPair: TokenPair };
  session: { userId: number; accountStatus: number; roles: string[] };
}
