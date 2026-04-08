import type {
  ConnectSession,
  ConnectionItem,
  InstitutionItem,
  SyncJobItem,
  SyncStatusData,
} from '@/features/connections/types/connections.types';
import { apiClient } from '@/services/api/client';
import { endpoints } from '@/services/api/endpoints';
import type { ApiResponse } from '@/types/api';


export async function fetchInstitutions(): Promise<InstitutionItem[]> {
  const response = await apiClient.get<ApiResponse<{ readonly institutions: InstitutionItem[] }>>(
    endpoints.openFinance.institutions,
  );
  return response.data.data.institutions;
}

export async function fetchConnections(): Promise<ConnectionItem[]> {
  const response = await apiClient.get<ApiResponse<{ readonly connections: ConnectionItem[] }>>(
    endpoints.openFinance.connections,
  );
  return response.data.data.connections;
}

export async function createConnectSession(institutionId: string): Promise<ConnectSession> {
  const consentResponse = await apiClient.post<
    ApiResponse<{ readonly consent: { readonly id: string } }>
  >(endpoints.openFinance.consents, {
    institution_id: institutionId,
    purpose: 'Consolidacao financeira pessoal',
    permissions: ['ACCOUNTS_READ', 'BALANCES_READ', 'TRANSACTIONS_READ'],
    redirect_uri: 'https://bufunfaai-api.caiof.com.br/v1/open-finance/callback',
  });

  const consentId: string = consentResponse.data.data.consent.id;
  const connectTokenResponse = await apiClient.post<ApiResponse<ConnectSession>>(
    endpoints.openFinance.connectToken(consentId),
  );

  return connectTokenResponse.data.data;
}

export async function completeInstitutionConnection(consentId: string, itemId: string): Promise<ConnectionItem> {
  const response = await apiClient.post<
    ApiResponse<{ readonly connection: ConnectionItem }>
  >(endpoints.openFinance.completeConsent(consentId), {
    item_id: itemId,
  });

  return response.data.data.connection;
}

export async function triggerConnectionSync(connectionId: string): Promise<SyncJobItem[]> {
  const response = await apiClient.post<ApiResponse<{ readonly jobs: SyncJobItem[] }>>(
    endpoints.openFinance.syncConnection(connectionId),
  );
  return response.data.data.jobs;
}

export async function fetchSyncStatus(connectionId: string): Promise<SyncStatusData> {
  const response = await apiClient.get<ApiResponse<SyncStatusData>>(
    endpoints.openFinance.syncStatus(connectionId),
  );
  return response.data.data;
}
