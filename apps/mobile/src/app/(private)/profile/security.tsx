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
  const biometricEnabled = useSessionStore((state) => state.biometricEnabled);
  const setBiometricEnabled = useSessionStore((state) => state.setBiometricEnabled);

  async function toggleBiometricPreference(): Promise<void> {
    const nextValue = !biometricEnabled;
    if (nextValue) {
      const authenticated: boolean = await authenticateWithBiometrics();
      if (!authenticated) {
        return;
      }
    }

    await setBiometricPreference(nextValue);
    setBiometricEnabled(nextValue);
  }

  return (
    <FeatureScreen description="Proteções locais e preferências de acesso do app." title="Segurança">
      <View style={styles.content}>
        <AppText>{`Biometria: ${biometricEnabled ? 'ativada' : 'desativada'}`}</AppText>
        <AppText color={theme.colors.textSecondary}>
          Use sua digital ou reconhecimento facial para voltar ao app com mais rapidez.
        </AppText>
        <Button
          label={biometricEnabled ? 'Desativar biometria' : 'Ativar biometria'}
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
