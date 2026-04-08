import { useRouter } from 'expo-router';
import type { ReactElement } from 'react';
import { Image, Pressable, StyleSheet, View } from 'react-native';

import { AppText } from '@/components/ui/AppText';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { Screen } from '@/components/ui/Screen';
import { theme } from '@/theme/tokens';

export default function WelcomeScreen(): ReactElement {
  const router = useRouter();

  return (
    <Screen>
      <View style={styles.container}>
        <View style={styles.hero}>
          <View style={styles.brandRow}>
            <Image source={require('../../../assets/logo-mark.png')} style={styles.brandMark} />
            <View style={styles.brandText}>
              <AppText color={theme.colors.accent} variant="label">
                BufunfaAI
              </AppText>
              <AppText color={theme.colors.textSecondary}>Financial operating system</AppText>
            </View>
          </View>
          <AppText color={theme.colors.accent} variant="label">
            Private finance cockpit
          </AppText>
          <AppText variant="display">Seu dinheiro com leitura clara, ritmo e contexto.</AppText>
          <AppText color={theme.colors.textSecondary}>
            BufunfaAI organiza conexoes bancarias, caixa mensal e sinais de risco em uma interface feita para decisao, nao para ruído.
          </AppText>
        </View>

        <Card style={styles.highlightCard}>
          <View pointerEvents="none" style={styles.highlightGlow} />
          <View style={styles.highlightHeader}>
            <AppText color={theme.colors.textPrimary} variant="label">
              Patrimonio visivel
            </AppText>
            <AppText color={theme.colors.accentSoft}>Open Finance + IA controlada</AppText>
          </View>
          <View style={styles.metrics}>
            <View style={styles.metricBlock}>
              <AppText variant="headline">
                3 camadas
              </AppText>
              <AppText color={theme.colors.textSecondary}>conta, fluxo e insight</AppText>
            </View>
            <View style={styles.metricBlock}>
              <AppText variant="headline">
                1 origem
              </AppText>
              <AppText color={theme.colors.textSecondary}>sua API como centro</AppText>
            </View>
          </View>
        </Card>

        <Card>
          <View style={styles.content}>
            <Button label="Entrar" onPress={(): void => router.push('/(public)/login')} />
            <Button label="Criar conta" onPress={(): void => router.push('/(public)/register')} variant="secondary" />
            <Pressable onPress={(): void => router.push('/(public)/login')} style={styles.secondaryAction}>
              <AppText color={theme.colors.textSecondary}>Ja tenho acesso e quero continuar do ultimo ponto.</AppText>
            </Pressable>
          </View>
        </Card>
      </View>
    </Screen>
  );
}

const styles = StyleSheet.create({
  container: {
    gap: theme.spacing.lg,
    justifyContent: 'center',
  },
  hero: {
    gap: theme.spacing.md,
    paddingTop: theme.spacing.md,
  },
  brandRow: {
    alignItems: 'center',
    flexDirection: 'row',
    gap: theme.spacing.md,
  },
  brandMark: {
    height: 56,
    width: 56,
  },
  brandText: {
    gap: 2,
  },
  highlightCard: {
    backgroundColor: theme.colors.surfaceInverse,
    overflow: 'hidden',
  },
  highlightGlow: {
    backgroundColor: theme.colors.primary,
    borderRadius: 180,
    height: 180,
    opacity: 0.1,
    position: 'absolute',
    right: -60,
    top: -80,
    width: 180,
  },
  highlightHeader: {
    gap: theme.spacing.sm,
  },
  metrics: {
    flexDirection: 'row',
    gap: theme.spacing.md,
    marginTop: theme.spacing.lg,
  },
  metricBlock: {
    flex: 1,
    gap: theme.spacing.xs,
  },
  content: {
    gap: theme.spacing.md,
  },
  secondaryAction: {
    alignItems: 'center',
    paddingHorizontal: theme.spacing.sm,
    paddingVertical: theme.spacing.sm,
  },
});
