import { useLocalSearchParams } from 'expo-router';
import type { ReactElement } from 'react';
import { StyleSheet, View } from 'react-native';

import { ErrorState } from '@/components/feedback/ErrorState';
import { LoadingState } from '@/components/feedback/LoadingState';
import { FeatureScreen } from '@/components/layout/FeatureScreen';
import { AppText } from '@/components/ui/AppText';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { useConnectionSyncMutation, useSyncStatusQuery } from '@/features/connections/hooks/useConnections';
import { theme } from '@/theme/tokens';

export default function ConnectionDetailsScreen(): ReactElement {
  const params = useLocalSearchParams<{ readonly id?: string }>();
  const connectionId: string = params.id ?? '';
  const syncStatusQuery = useSyncStatusQuery(connectionId);
  const syncMutation = useConnectionSyncMutation();

  async function handleSync(): Promise<void> {
    await syncMutation.mutateAsync(connectionId);
    await syncStatusQuery.refetch();
  }

  return (
    <FeatureScreen description="Acompanhe o status do consentimento e da sincronização desta conexão." title="Detalhe da conexão">
      {syncStatusQuery.isLoading ? <LoadingState label="Carregando status..." /> : null}
      {syncStatusQuery.isError ? <ErrorState onRetry={(): void => void syncStatusQuery.refetch()} /> : null}
      {syncStatusQuery.data ? (
        <>
          <Card>
            <View style={styles.content}>
              <AppText variant="headline">{syncStatusQuery.data.connection.status}</AppText>
              <AppText color={theme.colors.textSecondary}>
                Consentimento: {syncStatusQuery.data.connection.consent_id}
              </AppText>
              <AppText color={theme.colors.textSecondary}>
                Último sucesso: {syncStatusQuery.data.connection.last_successful_sync_at ?? 'ainda não houve'}
              </AppText>
              <Button label={syncMutation.isPending ? 'Sincronizando...' : 'Sincronizar agora'} onPress={(): void => void handleSync()} />
            </View>
          </Card>

          {syncStatusQuery.data.jobs.map((job) => (
            <Card key={job.id}>
              <View style={styles.content}>
                <AppText>{job.resource_type}</AppText>
                <AppText color={theme.colors.textSecondary}>{job.status}</AppText>
                <AppText color={theme.colors.textSecondary}>
                  Tentativas: {job.attempt_count}
                </AppText>
              </View>
            </Card>
          ))}
        </>
      ) : null}
    </FeatureScreen>
  );
}

const styles = StyleSheet.create({
  content: {
    gap: theme.spacing.sm,
  },
});
