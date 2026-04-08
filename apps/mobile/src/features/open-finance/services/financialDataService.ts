import type {
  AccountSnapshot,
  FinancialOverview,
  FinancialTransactionFeed,
} from '@/features/open-finance/types/financialData.types';
import { apiClient } from '@/services/api/client';
import { endpoints } from '@/services/api/endpoints';
import type { ApiResponse } from '@/types/api';

interface TransactionsFeedOptions {
  readonly from?: string;
  readonly limit?: number;
  readonly to?: string;
}

export async function fetchFinancialOverview(): Promise<FinancialOverview> {
  const response = await apiClient.get<ApiResponse<{ readonly overview: FinancialOverview }>>(
    endpoints.openFinance.overview,
  );

  return response.data.data.overview;
}

export async function fetchAccounts(): Promise<AccountSnapshot[]> {
  const response = await apiClient.get<ApiResponse<{ readonly accounts: AccountSnapshot[] }>>(
    endpoints.openFinance.accounts,
  );

  return response.data.data.accounts;
}

export async function fetchTransactionsFeed(
  options: TransactionsFeedOptions = {},
): Promise<FinancialTransactionFeed> {
  const params = new URLSearchParams();

  if (options.from) {
    params.set('from', options.from);
  }
  if (options.to) {
    params.set('to', options.to);
  }
  if (options.limit) {
    params.set('limit', String(options.limit));
  }

  const search = params.toString();
  const response = await apiClient.get<ApiResponse<{ readonly feed: FinancialTransactionFeed }>>(
    search.length > 0 ? `${endpoints.openFinance.transactions}?${search}` : endpoints.openFinance.transactions,
  );

  return response.data.data.feed;
}
