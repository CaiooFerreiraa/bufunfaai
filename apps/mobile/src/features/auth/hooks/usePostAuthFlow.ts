import { useRouter } from 'expo-router';

import { canUseBiometrics } from '@/services/biometric/biometricService';
import { getBiometricPreference, getBiometricPrompted } from '@/services/storage/preferences';
import { useSessionStore } from '@/stores/sessionStore';

interface UsePostAuthFlowResult {
  readonly completeAuthentication: () => Promise<void>;
}

export function usePostAuthFlow(): UsePostAuthFlowResult {
  const router = useRouter();
  const setRequiresBiometricUnlock = useSessionStore((state) => state.setRequiresBiometricUnlock);
  const setBiometricEnabled = useSessionStore((state) => state.setBiometricEnabled);

  async function completeAuthentication(): Promise<void> {
    const biometricEnabled: boolean = await getBiometricPreference();
    const biometricPrompted: boolean = await getBiometricPrompted();
    const biometricAvailable: boolean = await canUseBiometrics();

    setBiometricEnabled(biometricEnabled);
    setRequiresBiometricUnlock(false);

    if (!biometricEnabled && !biometricPrompted && biometricAvailable) {
      router.replace('/(auth)/setup-biometric');
      return;
    }

    router.replace('/(private)/home');
  }

  return {
    completeAuthentication,
  };
}
