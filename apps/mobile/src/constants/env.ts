import Constants from 'expo-constants';

import type { ExpoPublicEnv } from '@/types/env';

const extra: ExpoPublicEnv = (Constants.expoConfig?.extra ?? {}) as ExpoPublicEnv;

export const env = {
  apiUrl: process.env.EXPO_PUBLIC_API_URL ?? extra.EXPO_PUBLIC_API_URL ?? 'http://localhost:8080',
  appEnv: process.env.EXPO_PUBLIC_ENV ?? extra.EXPO_PUBLIC_ENV ?? 'development',
} as const;
