import type { AxiosError, AxiosResponse, InternalAxiosRequestConfig } from 'axios';

import { getClerkToken } from '@/lib/clerkToken';
import { apiClient } from '@/services/api/client';

export function registerApiInterceptors(): void {
  apiClient.interceptors.request.use(async (config: InternalAxiosRequestConfig): Promise<InternalAxiosRequestConfig> => {
    const accessToken: string | null = await getClerkToken();
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
      throw error;
    },
  );
}
