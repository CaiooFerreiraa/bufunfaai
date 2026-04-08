import type { UserProfile } from '@bufunfa/shared-types';

export interface SessionTokens {
  readonly accessToken: string;
  readonly refreshToken: string;
  readonly expiresAt: string;
}

export interface SessionMetadata {
  readonly biometricEnabled: boolean;
  readonly lastLoginAt: string;
}

export interface SessionState {
  readonly tokens: SessionTokens | null;
  readonly user: UserProfile | null;
  readonly metadata: SessionMetadata | null;
}
