import { QueryClientProvider } from '@tanstack/react-query';
import { Stack } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import type { ReactElement } from 'react';

import { useSessionBootstrap } from '@/features/auth/hooks/useSessionBootstrap';
import { useAppBootstrap } from '@/hooks/useAppBootstrap';
import { useBiometricAppLock } from '@/hooks/useBiometricAppLock';
import { useInitializeApi } from '@/hooks/useInitializeApi';
import { queryClient } from '@/lib/queryClient';
import { useSessionStore } from '@/stores/sessionStore';
import { theme } from '@/theme/tokens';

export default function RootLayout(): ReactElement {
  useInitializeApi();
  const isReady: boolean = useAppBootstrap();
  const isAuthenticated = useSessionStore((state) => state.isAuthenticated);
  const requiresBiometricUnlock = useSessionStore((state) => state.requiresBiometricUnlock);

  useSessionBootstrap(isReady && isAuthenticated && !requiresBiometricUnlock);
  useBiometricAppLock(isReady);

  if (!isReady) {
    return <></>;
  }

  return (
    <QueryClientProvider client={queryClient}>
      <StatusBar style="light" />
      <Stack
        screenOptions={{
          contentStyle: {
            backgroundColor: theme.colors.background,
          },
          headerShown: false,
        }}
      />
    </QueryClientProvider>
  );
}
