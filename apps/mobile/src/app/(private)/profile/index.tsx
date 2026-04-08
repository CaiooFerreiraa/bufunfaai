import type { ReactElement } from 'react';
import { StyleSheet, View } from 'react-native';

import { ErrorState } from '@/components/feedback/ErrorState';
import { LoadingState } from '@/components/feedback/LoadingState';
import { FeatureScreen } from '@/components/layout/FeatureScreen';
import { LinkCard } from '@/components/layout/LinkCard';
import { AppText } from '@/components/ui/AppText';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { useAuth } from '@/features/auth/hooks/useAuth';
import { useCurrentUserQuery } from '@/features/profile/hooks/useProfile';
import { theme } from '@/theme/tokens';

export default function ProfileScreen(): ReactElement {
  const { logout } = useAuth();
  const userQuery = useCurrentUserQuery();
  const initials: string = userQuery.data?.fullName
    ?.split(' ')
    .slice(0, 2)
    .map((chunk: string): string => chunk.charAt(0).toUpperCase())
    .join('') ?? 'BU';

  return (
    <FeatureScreen description="Gerencie sua conta, preferências de segurança e atalhos operacionais." title="Perfil">
      {userQuery.isLoading ? <LoadingState label="Carregando perfil..." /> : null}
      {userQuery.isError ? <ErrorState onRetry={(): void => void userQuery.refetch()} /> : null}
      {userQuery.data ? (
        <Card style={styles.heroCard}>
          <View style={styles.heroContent}>
            <View style={styles.avatar}>
              <AppText color={theme.colors.textInverse} variant="headline">
                {initials}
              </AppText>
            </View>
            <View style={styles.heroText}>
              <AppText variant="headline">{userQuery.data.fullName}</AppText>
              <AppText color={theme.colors.textSecondary}>{userQuery.data.email}</AppText>
            </View>
          </View>
          <View style={styles.securityBadge}>
            <AppText color={theme.colors.textInverse} variant="label">
              Segurança alta
            </AppText>
          </View>
        </Card>
      ) : null}

      <View style={styles.content}>
        <LinkCard description="Biometria, logout e configurações de sessão." href="/(private)/profile/security" title="Segurança" />
        <LinkCard description="Preferências de alertas e estado do push no ambiente atual." href="/(private)/profile/notifications" title="Notificações" />
      </View>

      <Button label="Sair da conta" onPress={(): void => void logout()} variant="secondary" />
    </FeatureScreen>
  );
}

const styles = StyleSheet.create({
  content: {
    gap: theme.spacing.md,
  },
  heroCard: {
    backgroundColor: theme.colors.surfaceMuted,
    gap: theme.spacing.md,
  },
  heroContent: {
    alignItems: 'center',
    flexDirection: 'row',
    gap: theme.spacing.md,
  },
  heroText: {
    flex: 1,
    gap: theme.spacing.xs,
  },
  avatar: {
    alignItems: 'center',
    backgroundColor: theme.colors.primary,
    borderRadius: theme.radii.pill,
    height: 56,
    justifyContent: 'center',
    width: 56,
  },
  securityBadge: {
    alignSelf: 'flex-start',
    backgroundColor: theme.colors.primary,
    borderRadius: theme.radii.pill,
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.sm,
  },
});
