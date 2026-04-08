import * as SecureStore from 'expo-secure-store';

import type { SessionMetadata, SessionTokens } from '@/types/session';

const REFRESH_TOKEN_KEY: string = 'bufunfaai.refresh_token';
const SESSION_METADATA_KEY: string = 'bufunfaai.session_metadata';

export async function saveSessionTokens(tokens: SessionTokens): Promise<void> {
  await SecureStore.setItemAsync(REFRESH_TOKEN_KEY, tokens.refreshToken);
  await SecureStore.setItemAsync(SESSION_METADATA_KEY, JSON.stringify({ expiresAt: tokens.expiresAt }));
}

export async function getStoredSessionTokens(): Promise<SessionTokens | null> {
  const refreshToken: string | null = await SecureStore.getItemAsync(REFRESH_TOKEN_KEY);
  const metadataRaw: string | null = await SecureStore.getItemAsync(SESSION_METADATA_KEY);

  if (!refreshToken || !metadataRaw) {
    return null;
  }

  const metadata: { readonly expiresAt: string } = JSON.parse(metadataRaw) as { readonly expiresAt: string };
  return {
    accessToken: '',
    refreshToken,
    expiresAt: metadata.expiresAt,
  };
}

export async function saveSessionMetadata(metadata: SessionMetadata): Promise<void> {
  await SecureStore.setItemAsync('bufunfaai.local_metadata', JSON.stringify(metadata));
}

export async function getSessionMetadata(): Promise<SessionMetadata | null> {
  const rawValue: string | null = await SecureStore.getItemAsync('bufunfaai.local_metadata');
  if (!rawValue) {
    return null;
  }

  return JSON.parse(rawValue) as SessionMetadata;
}

export async function saveRefreshToken(refreshToken: string): Promise<void> {
  await SecureStore.setItemAsync(REFRESH_TOKEN_KEY, refreshToken);
}

export async function getRefreshToken(): Promise<string | null> {
  return SecureStore.getItemAsync(REFRESH_TOKEN_KEY);
}

export async function clearSecureSession(): Promise<void> {
  await SecureStore.deleteItemAsync(REFRESH_TOKEN_KEY);
  await SecureStore.deleteItemAsync(SESSION_METADATA_KEY);
  await SecureStore.deleteItemAsync('bufunfaai.local_metadata');
}
