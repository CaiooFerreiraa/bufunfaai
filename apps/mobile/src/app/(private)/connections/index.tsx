import { Link } from 'expo-router';
import type { ReactElement } from 'react';
import { Pressable, StyleSheet, View } from 'react-native';

import { EmptyState } from '@/components/feedback/EmptyState';
import { ErrorState } from '@/components/feedback/ErrorState';
import { LoadingState } from '@/components/feedback/LoadingState';
import { FeatureScreen } from '@/components/layout/FeatureScreen';
import { AppText } from '@/components/ui/AppText';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { useConnectInstitutionMutation, useConnectionsQuery, useInstitutionsQuery } from '@/features/connections/hooks/useConnections';
import { theme } from '@/theme/tokens';

export default function ConnectionsScreen(): ReactElement {
  const institutionsQuery = useInstitutionsQuery();
  const connectionsQuery = useConnectionsQuery();
  const connectMutation = useConnectInstitutionMutation();

  async function handleConnectFirstInstitution(): Promise<void> {
    const firstInstitution = institutionsQuery.data?.[0];
    if (!firstInstitution) {
      return;
    }

    await connectMutation.mutateAsync(firstInstitution.id);
    await connectionsQuery.refetch();
  }

  return (
    <FeatureScreen description="Conecte bancos, acompanhe consentimentos e mantenha a leitura do Open Finance sob controle." title="Bancos conectados">
      <Card style={styles.heroCard}>
        <View style={styles.heroRow}>
          <View style={styles.heroMetric}>
            <AppText color={theme.colors.accent} variant="label">
              Ativas
            </AppText>
            <AppText variant="headline">{String(connectionsQuery.data?.filter((item) => item.status === 'ACTIVE').length ?? 0).padStart(2, '0')}</AppText>
          </View>
          <View style={styles.heroMetric}>
            <AppText color={theme.colors.accent} variant="label">
              Disponíveis
            </AppText>
            <AppText variant="headline">{String(institutionsQuery.data?.length ?? 0).padStart(2, '0')}</AppText>
          </View>
        </View>
        <AppText color={theme.colors.textSecondary}>
          A proposta visual agora trata suas integrações como módulos ativos do cockpit, não como uma lista burocrática.
        </AppText>
      </Card>

      {institutionsQuery.isLoading || connectionsQuery.isLoading ? <LoadingState label="Buscando instituições..." /> : null}
      {institutionsQuery.isError || connectionsQuery.isError ? (
        <ErrorState
          onRetry={(): void => {
            void institutionsQuery.refetch();
            void connectionsQuery.refetch();
          }}
        />
      ) : null}

      {!institutionsQuery.isLoading && !connectionsQuery.isLoading && !institutionsQuery.isError && !connectionsQuery.isError ? (
        <>
          {connectionsQuery.data && connectionsQuery.data.length > 0 ? (
            <View style={styles.list}>
              {connectionsQuery.data.map((connection) => (
                <Link asChild href={`/(private)/connections/${connection.id}`} key={connection.id}>
                  <Pressable>
                    <Card style={styles.connectionCard}>
                      <View style={styles.card}>
                        <View style={styles.connectionTop}>
                          <AppText color={theme.colors.accent} variant="label">
                            Conexão ativa
                          </AppText>
                          <View style={styles.statusChip}>
                            <AppText>{connection.status}</AppText>
                          </View>
                        </View>
                        <AppText variant="headline">{connection.institution_id}</AppText>
                        <AppText color={theme.colors.textSecondary}>Último sync: {connection.last_sync_at ?? 'ainda não executado'}</AppText>
                      </View>
                    </Card>
                  </Pressable>
                </Link>
              ))}
            </View>
          ) : (
            <EmptyState
              actionLabel={connectMutation.isPending ? 'Conectando...' : 'Conectar primeiro banco'}
              description="O fluxo mock de Open Finance já está ligado na API para validar consentimento, callback e sincronização."
              onActionPress={(): void => void handleConnectFirstInstitution()}
              title="Nenhuma conexão ativa"
            />
          )}

          <Card>
            <View style={styles.card}>
              <AppText color={theme.colors.accent} variant="label">
                Cadastro de bancos
              </AppText>
              <AppText variant="headline">Instituições disponíveis</AppText>
              {(institutionsQuery.data ?? []).map((institution) => (
                <View key={institution.id} style={styles.row}>
                  <View style={styles.meta}>
                    <AppText>{institution.display_name}</AppText>
                    <AppText color={theme.colors.textSecondary}>{institution.status}</AppText>
                  </View>
                  <Button label="Conectar" onPress={(): void => void connectMutation.mutateAsync(institution.id)} style={styles.connectButton} />
                </View>
              ))}
            </View>
          </Card>
        </>
      ) : null}
    </FeatureScreen>
  );
}

const styles = StyleSheet.create({
  heroCard: {
    gap: theme.spacing.md,
  },
  heroRow: {
    flexDirection: 'row',
    gap: theme.spacing.md,
  },
  heroMetric: {
    backgroundColor: theme.colors.surfaceMuted,
    borderColor: theme.colors.border,
    borderRadius: theme.radii.lg,
    borderWidth: 1,
    flex: 1,
    gap: theme.spacing.xs,
    padding: theme.spacing.lg,
  },
  list: {
    gap: theme.spacing.md,
  },
  connectionCard: {
    backgroundColor: theme.colors.surfaceMuted,
  },
  card: {
    gap: theme.spacing.sm,
  },
  connectionTop: {
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  statusChip: {
    backgroundColor: theme.colors.surface,
    borderColor: theme.colors.border,
    borderRadius: theme.radii.pill,
    borderWidth: 1,
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.sm,
  },
  row: {
    alignItems: 'center',
    borderTopColor: theme.colors.border,
    borderTopWidth: 1,
    flexDirection: 'row',
    gap: theme.spacing.md,
    justifyContent: 'space-between',
    paddingTop: theme.spacing.md,
  },
  meta: {
    flex: 1,
    gap: theme.spacing.xs,
  },
  connectButton: {
    minWidth: 116,
  },
});
