import { useConnectionsQuery } from '@/features/connections/hooks/useConnections';
import { useFinancialOverviewQuery } from '@/features/open-finance/hooks/useFinancialData';
import { useCurrentUserQuery } from '@/features/profile/hooks/useProfile';
import { formatCurrency } from '@/utils/formatCurrency';

export interface DashboardSummary {
  readonly activeConnections: number;
  readonly pendingConnections: number;
  readonly totalBalanceLabel: string;
  readonly userName: string;
}

export function useDashboardSummary(): DashboardSummary {
  const { data: user } = useCurrentUserQuery();
  const { data: connections } = useConnectionsQuery();
  const { data: overview } = useFinancialOverviewQuery();

  const connectionList = connections ?? [];
  const firstName: string | undefined = user?.fullName?.split(' ')[0];

  return {
    activeConnections: connectionList.filter((item) => item.status === 'ACTIVE').length,
    pendingConnections: connectionList.filter((item) => item.status !== 'ACTIVE').length,
    totalBalanceLabel: formatCurrency(overview?.total_available ?? 0),
    userName: firstName ?? 'Você',
  };
}
