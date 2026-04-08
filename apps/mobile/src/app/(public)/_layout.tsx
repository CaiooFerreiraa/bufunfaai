import { Redirect, Stack } from 'expo-router';
import type { ReactElement } from 'react';

import { useSessionStore } from '@/stores/sessionStore';

export default function PublicLayout(): ReactElement {
  const isAuthenticated = useSessionStore((state) => state.isAuthenticated);
  const requiresBiometricUnlock = useSessionStore((state) => state.requiresBiometricUnlock);

  if (isAuthenticated && requiresBiometricUnlock) {
    return <Redirect href="/(auth)/biometric-gate" />;
  }

  if (isAuthenticated && !requiresBiometricUnlock) {
    return <Redirect href="/(private)/home" />;
  }

  return <Stack screenOptions={{ headerShown: false }} />;
}
