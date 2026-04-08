import type { UserProfile } from '@bufunfa/shared-types';

import type {
  AuthResponseData,
  LoginRequest,
  MeResponseData,
  RegisterRequest,
} from '@/features/auth/types/auth.types';
import {
  mapAuthTokens,
  mapAuthUserToProfile,
  mapMeToProfile,
} from '@/features/auth/types/auth.types';
import { apiClient, authApiClient } from '@/services/api/client';
import { endpoints } from '@/services/api/endpoints';
import { getRefreshToken } from '@/services/storage/secureStore';
import type { ApiResponse } from '@/types/api';
import type { SessionTokens } from '@/types/session';

export interface AuthSessionResult {
  readonly tokens: SessionTokens;
  readonly user: UserProfile;
}

export async function loginRequest(input: LoginRequest): Promise<AuthSessionResult> {
  const response = await authApiClient.post<ApiResponse<AuthResponseData>>(endpoints.auth.login, input);
  return {
    tokens: mapAuthTokens(response.data.data.session),
    user: mapAuthUserToProfile(response.data.data.user),
  };
}

export async function registerRequest(input: RegisterRequest): Promise<AuthSessionResult> {
  const response = await authApiClient.post<ApiResponse<AuthResponseData>>(endpoints.auth.register, {
    email: input.email,
    full_name: input.fullName,
    password: input.password,
    phone: input.phone,
  });

  return {
    tokens: mapAuthTokens(response.data.data.session),
    user: mapAuthUserToProfile(response.data.data.user),
  };
}

export async function refreshSessionRequest(): Promise<AuthSessionResult> {
  const refreshToken: string | null = await getRefreshToken();
  if (!refreshToken) {
    throw new Error('Missing refresh token');
  }

  const response = await authApiClient.post<ApiResponse<AuthResponseData>>(endpoints.auth.refresh, {
    refresh_token: refreshToken,
  });

  return {
    tokens: mapAuthTokens(response.data.data.session),
    user: mapAuthUserToProfile(response.data.data.user),
  };
}

export async function fetchCurrentUser(): Promise<UserProfile> {
  const response = await apiClient.get<ApiResponse<MeResponseData>>(endpoints.users.me);
  return mapMeToProfile(response.data.data.user);
}

export async function logoutRequest(): Promise<void> {
  const refreshToken: string | null = await getRefreshToken();
  if (!refreshToken) {
    return;
  }

  await apiClient.post(endpoints.auth.logout, {
    refresh_token: refreshToken,
  });
}
