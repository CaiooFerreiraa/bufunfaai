import { ClerkProvider, useAuth } from '@clerk/expo';
import { tokenCache } from '@clerk/expo/token-cache';
import { QueryClientProvider } from '@tanstack/react-query';
import { Stack } from 'expo-router';
import { StatusBar } from 'expo-status-bar';
import { useEffect, type ReactElement } from 'react';

import { env } from '@/constants/env';
import { useAppBootstrap } from '@/hooks/useAppBootstrap';
import { useBiometricAppLock } from '@/hooks/useBiometricAppLock';
import { useInitializeApi } from '@/hooks/useInitializeApi';
import { setClerkTokenGetter } from '@/lib/clerkToken';
import { queryClient } from '@/lib/queryClient';
import { theme } from '@/theme/tokens';

export default function RootLayout(): ReactElement {
  if (!env.clerkPublishableKey) {
    throw new Error('Defina EXPO_PUBLIC_CLERK_PUBLISHABLE_KEY no ambiente do mobile.');
  }

  return (
    <ClerkProvider publishableKey={env.clerkPublishableKey} tokenCache={tokenCache}>
      <RootNavigator />
    </ClerkProvider>
  );
}

function ClerkTokenBridge(): null {
  const { getToken, isLoaded } = useAuth();

  useEffect(() => {
    if (!isLoaded) {
      setClerkTokenGetter(null);
      return;
    }

    setClerkTokenGetter(() => getToken());

    return (): void => {
      setClerkTokenGetter(null);
    };
  }, [getToken, isLoaded]);

  return null;
}

function RootNavigator(): ReactElement {
  useInitializeApi();
  const isReady: boolean = useAppBootstrap();
  const { isLoaded } = useAuth();

  useBiometricAppLock(isReady && isLoaded);

  if (!isReady || !isLoaded) {
    return <></>;
  }

  return (
    <QueryClientProvider client={queryClient}>
      <ClerkTokenBridge />
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
