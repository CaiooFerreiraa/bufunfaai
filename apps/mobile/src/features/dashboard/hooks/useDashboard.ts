import { useConnectionsQuery } from '@/features/connections/hooks/useConnections';
import { useCurrentUserQuery } from '@/features/profile/hooks/useProfile';

export interface DashboardSummary {
  readonly activeConnections: number;
  readonly pendingConnections: number;
  readonly totalBalanceLabel: string;
  readonly userName: string;
}

export function useDashboardSummary(): DashboardSummary {
  const { data: user } = useCurrentUserQuery();
  const { data: connections } = useConnectionsQuery();

  const connectionList = connections ?? [];
  const firstName: string | undefined = user?.fullName?.split(' ')[0];

  return {
    activeConnections: connectionList.filter((item) => item.status === 'ACTIVE').length,
    pendingConnections: connectionList.filter((item) => item.status !== 'ACTIVE').length,
    totalBalanceLabel: 'R$ 6.324,49',
    userName: firstName ?? 'Você',
  };
}
