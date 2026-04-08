import type {
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

export async function connectInstitution(institutionId: string): Promise<ConnectionItem> {
  const redirectUri: string = 'https://app.bufunfa.ai/open-finance/callback';

  const consentResponse = await apiClient.post<
    ApiResponse<{ readonly consent: { readonly id: string } }>
  >(endpoints.openFinance.consents, {
    institution_id: institutionId,
    purpose: 'Consolidacao financeira pessoal',
    permissions: ['ACCOUNTS_READ', 'BALANCES_READ', 'TRANSACTIONS_READ'],
    redirect_uri: redirectUri,
  });

  const consentId: string = consentResponse.data.data.consent.id;
  const authorizationResponse = await apiClient.post<
    ApiResponse<{ readonly authorization_url: string }>
  >(endpoints.openFinance.authorizeConsent(consentId));
  const authorizationUrl: string = authorizationResponse.data.data.authorization_url;

  const parsedUrl = new URL(authorizationUrl);
  const state: string = parsedUrl.searchParams.get('state') ?? '';
  const code: string = parsedUrl.searchParams.get('code') ?? '';

  const callbackResponse = await apiClient.post<
    ApiResponse<{ readonly connection: ConnectionItem }>
  >(endpoints.openFinance.callback, {
    state,
    code,
  });

  return callbackResponse.data.data.connection;
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
