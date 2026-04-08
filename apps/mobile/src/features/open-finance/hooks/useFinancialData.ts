import type { UseQueryResult } from '@tanstack/react-query';
import { useQuery } from '@tanstack/react-query';

import {
  fetchAccounts,
  fetchFinancialOverview,
  fetchTransactionsFeed,
} from '@/features/open-finance/services/financialDataService';
import type {
  AccountSnapshot,
  FinancialOverview,
  FinancialTransactionFeed,
} from '@/features/open-finance/types/financialData.types';

const overviewQueryKey: readonly ['open-finance', 'overview'] = ['open-finance', 'overview'];
const accountsQueryKey: readonly ['open-finance', 'accounts'] = ['open-finance', 'accounts'];

interface TransactionsFeedParams {
  readonly from?: string;
  readonly limit?: number;
  readonly to?: string;
}

export function useFinancialOverviewQuery(): UseQueryResult<FinancialOverview, Error> {
  return useQuery({
    queryKey: overviewQueryKey,
    queryFn: fetchFinancialOverview,
    staleTime: 30 * 1000,
  });
}

export function useAccountsQuery(): UseQueryResult<AccountSnapshot[], Error> {
  return useQuery({
    queryKey: accountsQueryKey,
    queryFn: fetchAccounts,
    staleTime: 30 * 1000,
  });
}

export function useTransactionsFeedQuery(params: TransactionsFeedParams): UseQueryResult<FinancialTransactionFeed, Error> {
  return useQuery({
    queryKey: ['open-finance', 'transactions-feed', params.from ?? '', params.to ?? '', params.limit ?? 0],
    queryFn: (): Promise<FinancialTransactionFeed> => fetchTransactionsFeed(params),
    staleTime: 30 * 1000,
  });
}
