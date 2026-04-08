import { Redirect, Stack } from 'expo-router';
import type { ReactElement } from 'react';
import { StyleSheet, View } from 'react-native';
import { useSafeAreaInsets } from 'react-native-safe-area-context';

import { BottomNavigation } from '@/components/layout/BottomNavigation';
import { useSessionStore } from '@/stores/sessionStore';
import { theme } from '@/theme/tokens';

export default function PrivateLayout(): ReactElement {
  const isAuthenticated = useSessionStore((state) => state.isAuthenticated);
  const requiresBiometricUnlock = useSessionStore((state) => state.requiresBiometricUnlock);
  const insets = useSafeAreaInsets();

  if (!isAuthenticated) {
    return <Redirect href="/(public)/welcome" />;
  }

  if (requiresBiometricUnlock) {
    return <Redirect href="/(auth)/biometric-gate" />;
  }

  return (
    <View style={styles.shell}>
      <Stack
        screenOptions={{
          contentStyle: {
            backgroundColor: theme.colors.background,
          },
          headerShown: false,
        }}
      />
      <View style={[styles.navWrap, { paddingBottom: Math.max(insets.bottom, 8) }]}>
        <BottomNavigation />
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  shell: {
    backgroundColor: theme.colors.background,
    flex: 1,
  },
  navWrap: {
    backgroundColor: theme.colors.surface,
    bottom: 0,
    left: 0,
    position: 'absolute',
    right: 0,
  },
});
