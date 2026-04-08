import { useRouter } from 'expo-router';
import { Fingerprint } from 'lucide-react-native';
import type { ReactElement } from 'react';
import { useCallback, useEffect, useState } from 'react';
import { StyleSheet, View } from 'react-native';

import { FeatureScreen } from '@/components/layout/FeatureScreen';
import { AppText } from '@/components/ui/AppText';
import { Button } from '@/components/ui/Button';
import { authenticateWithBiometrics, canUseBiometrics } from '@/services/biometric/biometricService';
import { setBiometricPreference, setBiometricPrompted } from '@/services/storage/preferences';
import { useSessionStore } from '@/stores/sessionStore';
import { theme } from '@/theme/tokens';

export default function SetupBiometricScreen(): ReactElement {
  const router = useRouter();
  const metadata = useSessionStore((state) => state.metadata);
  const tokens = useSessionStore((state) => state.tokens);
  const user = useSessionStore((state) => state.user);
  const setAuthenticatedSession = useSessionStore((state) => state.setAuthenticatedSession);
  const [isPrompting, setIsPrompting] = useState<boolean>(false);
  const [errorMessage, setErrorMessage] = useState<string>('');

  const handleEnableBiometrics = useCallback(async (): Promise<void> => {
    if (isPrompting) {
      return;
    }

    setIsPrompting(true);
    const available = await canUseBiometrics();
    if (!available || !tokens || !user || !metadata) {
      setIsPrompting(false);
      router.replace('/(private)/home');
      return;
    }

    const authenticated: boolean = await authenticateWithBiometrics();
    if (!authenticated) {
      setErrorMessage('Não foi possível validar sua biometria agora.');
      setIsPrompting(false);
      return;
    }

    await setBiometricPreference(true);
    await setBiometricPrompted(true);
    setAuthenticatedSession(tokens, user, {
      ...metadata,
      biometricEnabled: true,
    });
    setIsPrompting(false);
    router.replace('/(private)/home');
  }, [isPrompting, metadata, router, setAuthenticatedSession, tokens, user]);

  async function handleSkip(): Promise<void> {
    await setBiometricPreference(false);
    await setBiometricPrompted(true);
    router.replace('/(private)/home');
  }

  useEffect(() => {
    void handleEnableBiometrics();
  }, [handleEnableBiometrics]);

  return (
    <FeatureScreen description="Valide sua identidade local agora e deixe o retorno ao app mais rápido e mais seguro." title="Ative sua biometria">
      <View style={styles.content}>
        <View style={styles.badge}>
          <Fingerprint color={theme.colors.accent} size={20} strokeWidth={2} />
          <AppText color={theme.colors.accent} variant="label">
            Segurança imediata
          </AppText>
        </View>
        <AppText color={theme.colors.textSecondary}>
          A biometria não substitui o login da API. Ela protege apenas o desbloqueio local do aplicativo.
        </AppText>
        {errorMessage ? <AppText color={theme.colors.error}>{errorMessage}</AppText> : null}
        <Button label={isPrompting ? 'Validando biometria...' : 'Tentar novamente'} onPress={(): void => void handleEnableBiometrics()} />
        <Button label="Agora não" onPress={(): void => void handleSkip()} style={styles.secondary} variant="secondary" />
      </View>
    </FeatureScreen>
  );
}

const styles = StyleSheet.create({
  badge: {
    alignItems: 'center',
    backgroundColor: theme.colors.surfaceMuted,
    borderColor: theme.colors.border,
    borderRadius: theme.radii.pill,
    borderWidth: 1,
    flexDirection: 'row',
    gap: theme.spacing.sm,
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.sm,
    alignSelf: 'flex-start',
  },
  content: {
    gap: theme.spacing.md,
  },
  secondary: {
    marginTop: theme.spacing.xs,
  },
});
