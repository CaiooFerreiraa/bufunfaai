import type { ReactElement } from 'react';
import { StyleSheet, View } from 'react-native';

import { FeatureScreen } from '@/components/layout/FeatureScreen';
import { AppText } from '@/components/ui/AppText';
import { Button } from '@/components/ui/Button';
import { authenticateWithBiometrics } from '@/services/biometric/biometricService';
import { setBiometricPreference } from '@/services/storage/preferences';
import { useSessionStore } from '@/stores/sessionStore';
import { theme } from '@/theme/tokens';

export default function ProfileSecurityScreen(): ReactElement {
  const metadata = useSessionStore((state) => state.metadata);
  const tokens = useSessionStore((state) => state.tokens);
  const user = useSessionStore((state) => state.user);
  const setAuthenticatedSession = useSessionStore((state) => state.setAuthenticatedSession);

  async function toggleBiometricPreference(): Promise<void> {
    if (!metadata || !tokens || !user) {
      return;
    }

    const nextValue = !metadata.biometricEnabled;
    if (nextValue) {
      const authenticated: boolean = await authenticateWithBiometrics();
      if (!authenticated) {
        return;
      }
    }

    await setBiometricPreference(nextValue);
    setAuthenticatedSession(tokens, user, {
      ...metadata,
      biometricEnabled: nextValue,
    });
  }

  return (
    <FeatureScreen description="Proteções locais e preferências sensíveis da sessão do app." title="Segurança">
      <View style={styles.content}>
        <AppText>{`Biometria: ${metadata?.biometricEnabled ? 'ativada' : 'desativada'}`}</AppText>
        <AppText color={theme.colors.textSecondary}>
          O desbloqueio biométrico protege a reabertura local do aplicativo.
        </AppText>
        <Button
          label={metadata?.biometricEnabled ? 'Desativar biometria' : 'Ativar biometria'}
          onPress={(): void => void toggleBiometricPreference()}
        />
      </View>
    </FeatureScreen>
  );
}

const styles = StyleSheet.create({
  content: {
    gap: theme.spacing.md,
  },
});
