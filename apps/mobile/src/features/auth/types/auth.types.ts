import type { UserProfile } from '@bufunfa/shared-types';

import type { SessionTokens } from '@/types/session';

export interface LoginRequest {
  readonly email: string;
  readonly password: string;
}

export interface RegisterRequest {
  readonly email: string;
  readonly fullName: string;
  readonly password: string;
  readonly phone?: string;
}

export interface AuthUser {
  readonly id: string;
  readonly email: string;
  readonly full_name: string;
  readonly status: string;
}

export interface AuthSessionPayload {
  readonly access_token: string;
  readonly refresh_token: string;
  readonly expires_at: string;
}

export interface AuthResponseData {
  readonly session: AuthSessionPayload;
  readonly user: AuthUser;
}

export interface MeResponseData {
  readonly user: {
    readonly id: string;
    readonly email: string;
    readonly full_name: string;
    readonly phone?: string;
    readonly status: string;
  };
}

export function mapAuthUserToProfile(user: AuthUser): UserProfile {
  return {
    id: user.id,
    email: user.email,
    fullName: user.full_name,
  };
}

export function mapMeToProfile(data: MeResponseData['user']): UserProfile {
  return {
    id: data.id,
    email: data.email,
    fullName: data.full_name,
  };
}

export function mapAuthTokens(session: AuthSessionPayload): SessionTokens {
  return {
    accessToken: session.access_token,
    refreshToken: session.refresh_token,
    expiresAt: session.expires_at,
  };
}
