import type { UserProfile } from '@bufunfa/shared-types';
import { create } from 'zustand';

import type { SessionMetadata, SessionTokens } from '@/types/session';

interface SessionStoreState {
  readonly isBootstrapping: boolean;
  readonly isAuthenticated: boolean;
  readonly requiresBiometricUnlock: boolean;
  readonly tokens: SessionTokens | null;
  readonly user: UserProfile | null;
  readonly metadata: SessionMetadata | null;
  readonly setBootstrapState: (value: boolean) => void;
  readonly hydrateSession: (tokens: SessionTokens, metadata: SessionMetadata, user: UserProfile | null) => void;
  readonly setAuthenticatedSession: (
    tokens: SessionTokens,
    user: UserProfile,
    metadata: SessionMetadata,
  ) => void;
  readonly clearSession: () => void;
  readonly setRequiresBiometricUnlock: (value: boolean) => void;
}

export const useSessionStore = create<SessionStoreState>((set) => ({
  isBootstrapping: true,
  isAuthenticated: false,
  requiresBiometricUnlock: false,
  tokens: null,
  user: null,
  metadata: null,
  setBootstrapState: (value: boolean): void => set({ isBootstrapping: value }),
  hydrateSession: (
    tokens: SessionTokens,
    metadata: SessionMetadata,
    user: UserProfile | null,
  ): void =>
    set({
      isAuthenticated: true,
      metadata,
      requiresBiometricUnlock: false,
      tokens,
      user,
    }),
  setAuthenticatedSession: (
    tokens: SessionTokens,
    user: UserProfile,
    metadata: SessionMetadata,
  ): void =>
    set({
      isAuthenticated: true,
      metadata,
      requiresBiometricUnlock: false,
      tokens,
      user,
    }),
  clearSession: (): void =>
    set({
      isAuthenticated: false,
      metadata: null,
      requiresBiometricUnlock: false,
      tokens: null,
      user: null,
    }),
  setRequiresBiometricUnlock: (value: boolean): void => set({ requiresBiometricUnlock: value }),
}));
