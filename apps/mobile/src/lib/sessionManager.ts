import type { UserProfile } from '@bufunfa/shared-types';

import { getBiometricPreference } from '@/services/storage/preferences';
import { clearSecureSession, getSessionMetadata, getStoredSessionTokens, saveSessionMetadata, saveSessionTokens } from '@/services/storage/secureStore';
import { useSessionStore } from '@/stores/sessionStore';
import type { SessionMetadata, SessionTokens } from '@/types/session';

export async function persistSession(
  tokens: SessionTokens,
  user: UserProfile,
): Promise<void> {
  const biometricEnabled: boolean = await getBiometricPreference();
  const metadata: SessionMetadata = {
    biometricEnabled,
    lastLoginAt: new Date().toISOString(),
  };

  await saveSessionTokens(tokens);
  await saveSessionMetadata(metadata);
  useSessionStore.getState().setAuthenticatedSession(tokens, user, metadata);
}

export async function restoreSession(): Promise<void> {
  const tokens: SessionTokens | null = await getStoredSessionTokens();
  const metadata: SessionMetadata | null = await getSessionMetadata();

  if (!tokens || !metadata) {
    useSessionStore.getState().clearSession();
    return;
  }

  useSessionStore.getState().hydrateSession(tokens, metadata, null);

  if (metadata.biometricEnabled) {
    useSessionStore.getState().setRequiresBiometricUnlock(true);
  }
}

export async function destroySession(): Promise<void> {
  await clearSecureSession();
  useSessionStore.getState().clearSession();
}
