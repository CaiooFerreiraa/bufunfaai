import { Redirect, Stack } from 'expo-router';
import type { ReactElement } from 'react';

import { useSessionStore } from '@/stores/sessionStore';

export default function AuthLayout(): ReactElement {
  const isAuthenticated = useSessionStore((state) => state.isAuthenticated);

  if (!isAuthenticated) {
    return <Redirect href="/(public)/welcome" />;
  }

  return <Stack screenOptions={{ headerShown: false }} />;
}
