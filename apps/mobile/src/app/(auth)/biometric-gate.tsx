import { useClerk } from '@clerk/expo';
import { useRouter } from 'expo-router';
import { Fingerprint } from 'lucide-react-native';
import type { ReactElement } from 'react';
import { useCallback, useEffect } from 'react';
import { StyleSheet, View } from 'react-native';

import { FeatureScreen } from '@/components/layout/FeatureScreen';
import { AppText } from '@/components/ui/AppText';
import { Button } from '@/components/ui/Button';
import { authenticateWithBiometrics } from '@/services/biometric/biometricService';
import { useSessionStore } from '@/stores/sessionStore';
import { theme } from '@/theme/tokens';

export default function BiometricGateScreen(): ReactElement {
  const router = useRouter();
  const { signOut } = useClerk();
  const setRequiresBiometricUnlock = useSessionStore((state) => state.setRequiresBiometricUnlock);
  const clearSession = useSessionStore((state) => state.clearSession);

  const handleUnlock = useCallback(async (): Promise<void> => {
    const succeeded = await authenticateWithBiometrics();
    if (!succeeded) {
      return;
    }

    setRequiresBiometricUnlock(false);
    router.replace('/(private)/home');
  }, [router, setRequiresBiometricUnlock]);

  async function handleLogout(): Promise<void> {
    await signOut();
    clearSession();
    router.replace('/(public)/welcome');
  }

  useEffect((): void => {
    void handleUnlock();
  }, [handleUnlock]);

  return (
    <FeatureScreen description="Confirme sua biometria para continuar." title="Desbloqueie sua conta">
      <View style={styles.content}>
        <View style={styles.badge}>
          <Fingerprint color={theme.colors.accent} size={20} strokeWidth={2} />
          <AppText color={theme.colors.accent} variant="label">
            Confirmação rápida
          </AppText>
        </View>
        <Button label="Tentar novamente" onPress={(): void => void handleUnlock()} />
        <Button label="Sair da conta" onPress={(): void => void handleLogout()} style={styles.secondary} variant="secondary" />
        <AppText color={theme.colors.textSecondary}>
          Use sua digital ou reconhecimento facial para entrar no app.
        </AppText>
      </View>
    </FeatureScreen>
  );
}

const styles = StyleSheet.create({
  badge: {
    alignItems: 'center',
    alignSelf: 'flex-start',
    backgroundColor: theme.colors.surfaceMuted,
    borderColor: theme.colors.border,
    borderRadius: theme.radii.pill,
    borderWidth: 1,
    flexDirection: 'row',
    gap: theme.spacing.sm,
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.sm,
  },
  content: {
    gap: theme.spacing.md,
  },
  secondary: {
    marginTop: theme.spacing.xs,
  },
});
