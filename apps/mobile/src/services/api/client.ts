import axios, { type AxiosInstance } from 'axios';

import { env } from '@/constants/env';

export const apiClient: AxiosInstance = axios.create({
  baseURL: env.apiUrl,
  timeout: 15000,
});

export const authApiClient: AxiosInstance = axios.create({
  baseURL: env.apiUrl,
  timeout: 15000,
});
