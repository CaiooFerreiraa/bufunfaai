export type AppEnvironment = 'development' | 'staging' | 'production';

export interface ApiHealthResponse {
  readonly status: 'ok';
  readonly service: string;
  readonly version: string;
  readonly environment: AppEnvironment;
}

export interface AuthSession {
  readonly accessToken: string;
  readonly refreshToken: string;
  readonly expiresAt: string;
}

export interface UserProfile {
  readonly id: string;
  readonly email: string;
  readonly fullName: string;
}
