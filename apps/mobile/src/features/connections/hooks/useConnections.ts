import type {
  UseMutationResult,
  UseQueryResult,
} from '@tanstack/react-query';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';

import {
  completeInstitutionConnection,
  createConnectSession,
  fetchConnections,
  fetchInstitutions,
  fetchSyncStatus,
  triggerConnectionSync,
} from '@/features/connections/services/connectionsService';
import type {
  ConnectSession,
  ConnectionItem,
  InstitutionItem,
  SyncStatusData,
  SyncJobItem,
} from '@/features/connections/types/connections.types';
import { saveSnapshot } from '@/services/storage/sqlite';

const institutionsQueryKey: readonly ['open-finance', 'institutions'] = ['open-finance', 'institutions'];
const connectionsQueryKey: readonly ['open-finance', 'connections'] = ['open-finance', 'connections'];

export function useInstitutionsQuery(): UseQueryResult<InstitutionItem[], Error> {
  return useQuery({
    queryKey: institutionsQueryKey,
    queryFn: fetchInstitutions,
    staleTime: 5 * 60 * 1000,
  });
}

export function useConnectionsQuery(): UseQueryResult<ConnectionItem[], Error> {
  return useQuery({
    queryKey: connectionsQueryKey,
    queryFn: async (): Promise<ConnectionItem[]> => {
      const data = await fetchConnections();
      await saveSnapshot<ConnectionItem[]>({
        key: 'connections',
        payload: data,
        syncedAt: new Date().toISOString(),
      });
      return data;
    },
    staleTime: 30 * 1000,
  });
}

export function useCreateConnectSessionMutation(): UseMutationResult<ConnectSession, Error, string> {
  return useMutation({
    mutationFn: createConnectSession,
  });
}

export function useCompleteInstitutionConnectionMutation(): UseMutationResult<
  ConnectionItem,
  Error,
  { readonly consentId: string; readonly itemId: string }
> {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      consentId,
      itemId,
    }): Promise<ConnectionItem> => completeInstitutionConnection(consentId, itemId),
    onSuccess: async (): Promise<void> => {
      await queryClient.invalidateQueries({ queryKey: connectionsQueryKey });
    },
  });
}

export function useConnectionSyncMutation(): UseMutationResult<SyncJobItem[], Error, string> {
  return useMutation({
    mutationFn: triggerConnectionSync,
  });
}

export function useSyncStatusQuery(connectionId: string): UseQueryResult<SyncStatusData, Error> {
  return useQuery({
    queryKey: ['open-finance', 'sync-status', connectionId],
    queryFn: (): Promise<SyncStatusData> => fetchSyncStatus(connectionId),
    enabled: connectionId.length > 0,
    staleTime: 15 * 1000,
  });
}
