import { InstrumentSerif_400Regular } from '@expo-google-fonts/instrument-serif';
import { Manrope_500Medium, Manrope_600SemiBold } from '@expo-google-fonts/manrope';
import { QueryClientProvider } from '@tanstack/react-query';
import { useFonts } from 'expo-font';
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
  const [fontsLoaded]: [boolean, Error | null] = useFonts({
    InstrumentSerif_400Regular,
    Manrope_500Medium,
    Manrope_600SemiBold,
  });
  const isReady: boolean = useAppBootstrap();
  const isAuthenticated = useSessionStore((state) => state.isAuthenticated);
  const requiresBiometricUnlock = useSessionStore((state) => state.requiresBiometricUnlock);

  useSessionBootstrap(isReady && isAuthenticated && !requiresBiometricUnlock);
  useBiometricAppLock(isReady);

  if (!fontsLoaded || !isReady) {
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
