import type { UserProfile } from '@bufunfa/shared-types';
import { useRouter } from 'expo-router';

import type { LoginSchema, RegisterSchema } from '@/features/auth/schemas/authSchemas';
import { loginRequest, logoutRequest, registerRequest } from '@/features/auth/services/authService';
import { destroySession, persistSession } from '@/lib/sessionManager';
import { authenticateWithBiometrics, canUseBiometrics } from '@/services/biometric/biometricService';
import { getBiometricPreference, getBiometricPrompted, setBiometricPreference, setBiometricPrompted } from '@/services/storage/preferences';
import { useSessionStore } from '@/stores/sessionStore';
import type { SessionMetadata, SessionTokens } from '@/types/session';


interface UseAuthResult {
  readonly login: (input: LoginSchema) => Promise<void>;
  readonly logout: () => Promise<void>;
  readonly register: (input: RegisterSchema) => Promise<void>;
}

export function useAuth(): UseAuthResult {
  const router = useRouter();
  const setAuthenticatedSession = useSessionStore((state) => state.setAuthenticatedSession);

  async function routeAfterAuthentication(tokens: SessionTokens, user: UserProfile): Promise<void> {
    const biometricEnabled: boolean = await getBiometricPreference();
    const biometricPrompted: boolean = await getBiometricPrompted();
    const biometricAvailable: boolean = await canUseBiometrics();

    if (!biometricEnabled && !biometricPrompted && biometricAvailable) {
      const authenticated: boolean = await authenticateWithBiometrics();
      await setBiometricPrompted(true);

      if (authenticated) {
        await setBiometricPreference(true);
        await persistSession(tokens, user);
        const metadata: SessionMetadata | null = useSessionStore.getState().metadata;

        if (metadata) {
          setAuthenticatedSession(tokens, user, {
            ...metadata,
            biometricEnabled: true,
          });
        }
      }
    }

    router.replace('/(private)/home');
  }

  async function login(input: LoginSchema): Promise<void> {
    const response = await loginRequest({
      email: input.email,
      password: input.password,
    });

    await persistSession(response.tokens, response.user);
    await routeAfterAuthentication(response.tokens, response.user);
  }

  async function register(input: RegisterSchema): Promise<void> {
    const response = await registerRequest({
      email: input.email,
      fullName: input.fullName,
      password: input.password,
      phone: input.phone,
    });

    await persistSession(response.tokens, response.user);
    await routeAfterAuthentication(response.tokens, response.user);
  }

  async function logout(): Promise<void> {
    try {
      await logoutRequest();
    } finally {
      await destroySession();
      router.replace('/(public)/welcome');
    }
  }

  return {
    login,
    logout,
    register,
  };
}
