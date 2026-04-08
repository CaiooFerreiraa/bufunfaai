import { Redirect } from 'expo-router';
import type { ReactElement } from 'react';

import { useSessionStore } from '@/stores/sessionStore';

export default function IndexScreen(): ReactElement {
  const isAuthenticated = useSessionStore((state) => state.isAuthenticated);
  const requiresBiometricUnlock = useSessionStore((state) => state.requiresBiometricUnlock);

  if (!isAuthenticated) {
    return <Redirect href="/(public)/welcome" />;
  }

  if (requiresBiometricUnlock) {
    return <Redirect href="/(auth)/biometric-gate" />;
  }

  return <Redirect href="/(private)/home" />;
}
