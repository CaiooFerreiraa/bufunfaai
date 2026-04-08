import { useAuth } from '@clerk/expo';
import { Redirect, Stack } from 'expo-router';
import type { ReactElement } from 'react';

import { useSessionStore } from '@/stores/sessionStore';

export default function AuthLayout(): ReactElement {
  const { isLoaded, isSignedIn } = useAuth();
  const requiresBiometricUnlock = useSessionStore((state) => state.requiresBiometricUnlock);

  if (!isLoaded) {
    return <></>;
  }

  if (!isSignedIn) {
    return <Redirect href="/(public)/welcome" />;
  }

  if (!requiresBiometricUnlock) {
    return <Redirect href="/(private)/home" />;
  }

  return <Stack screenOptions={{ headerShown: false }} />;
}
