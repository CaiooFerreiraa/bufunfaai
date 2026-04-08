import type { AxiosError, AxiosResponse, InternalAxiosRequestConfig } from 'axios';

import { refreshSessionRequest } from '@/features/auth/services/authService';
import { destroySession, persistSession } from '@/lib/sessionManager';
import { apiClient } from '@/services/api/client';
import { useSessionStore } from '@/stores/sessionStore';

let isRefreshing: boolean = false;
let refreshPromise: Promise<void> | null = null;

export function registerApiInterceptors(): void {
  apiClient.interceptors.request.use((config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
    const accessToken: string | undefined = useSessionStore.getState().tokens?.accessToken;
    if (accessToken) {
      config.headers.Authorization = `Bearer ${accessToken}`;
    }

    return config;
  });

  apiClient.interceptors.response.use(
    (response) => response,
    async (error: AxiosError): Promise<AxiosResponse> => {
      const originalRequest = error.config as (InternalAxiosRequestConfig & { _retry?: boolean }) | undefined;
      if (error.response?.status !== 401 || !originalRequest || originalRequest._retry) {
        throw error;
      }

      originalRequest._retry = true;

      if (!isRefreshing) {
        isRefreshing = true;
        refreshPromise = (async (): Promise<void> => {
          const response = await refreshSessionRequest();
          await persistSession(response.tokens, response.user);
        })().finally((): void => {
          isRefreshing = false;
        });
      }

      try {
        await refreshPromise;
        const nextAccessToken: string | undefined = useSessionStore.getState().tokens?.accessToken;
        if (nextAccessToken) {
          originalRequest.headers.Authorization = `Bearer ${nextAccessToken}`;
        }
        return apiClient(originalRequest);
      } catch (refreshError) {
        await destroySession();
        throw refreshError;
      }
    },
  );
}
