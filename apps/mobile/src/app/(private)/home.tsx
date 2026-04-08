import {
  ArrowDownLeft,
  ArrowUpRight,
  Bell,
  CreditCard,
  ShieldCheck,
  Sparkles,
} from 'lucide-react-native';
import type { ReactElement } from 'react';
import { StyleSheet, View } from 'react-native';

import { CategoryPieChart, type CategorySlice } from '@/components/charts/CategoryPieChart';
import { ErrorState } from '@/components/feedback/ErrorState';
import { LoadingState } from '@/components/feedback/LoadingState';
import { LinkCard } from '@/components/layout/LinkCard';
import { AppText } from '@/components/ui/AppText';
import { Card } from '@/components/ui/Card';
import { Screen } from '@/components/ui/Screen';
import { useConnectionsQuery } from '@/features/connections/hooks/useConnections';
import { useDashboardSummary } from '@/features/dashboard/hooks/useDashboard';
import { useFinancialOverviewQuery } from '@/features/open-finance/hooks/useFinancialData';
import type {
  AccountSnapshot,
  CategoryBreakdown,
  FinancialOverview,
  FinancialTransaction,
} from '@/features/open-finance/types/financialData.types';
import { theme } from '@/theme/tokens';
import { formatCurrency } from '@/utils/formatCurrency';

interface QuickAction {
  readonly icon: typeof ArrowUpRight;
  readonly label: string;
}

interface InsightItem {
  readonly amount: string;
  readonly label: string;
  readonly tone: string;
}

const QUICK_ACTIONS: readonly QuickAction[] = [
  { icon: ArrowUpRight, label: 'Enviar' },
  { icon: ArrowDownLeft, label: 'Receber' },
  { icon: CreditCard, label: 'Cartões' },
  { icon: ShieldCheck, label: 'Segurança' },
] as const;

const CHART_COLORS: readonly string[] = [
  theme.colors.primary,
  '#78A7FF',
  '#37D7A3',
  '#F8A94C',
] as const;

export default function HomeScreen(): ReactElement {
  const { activeConnections, pendingConnections, totalBalanceLabel, userName } = useDashboardSummary();
  const connectionsQuery = useConnectionsQuery();
  const overviewQuery = useFinancialOverviewQuery();

  const overview = overviewQuery.data;
  const insightItems: readonly InsightItem[] = buildInsightItems(overview);
  const chartSlices: readonly CategorySlice[] = buildCategorySlices(overview?.expense_breakdown ?? []);
  const firstAccount: AccountSnapshot | undefined = overview?.accounts[0];
  const recentTransactions: readonly FinancialTransaction[] = overview?.recent_transactions ?? [];
  const netFlow: number = (overview?.month_income ?? 0) - (overview?.month_expenses ?? 0);

  return (
    <Screen>
      <View style={styles.content}>
        <View style={styles.header}>
          <View style={styles.headerText}>
            <View style={styles.eyebrow}>
              <View style={styles.eyebrowDot} />
              <AppText variant="label">BufunfaAI wallet</AppText>
            </View>
            <AppText variant="display">Olá, {userName}</AppText>
          </View>
          <View style={styles.headerAction}>
            <Bell color={theme.colors.textPrimary} size={18} strokeWidth={2} />
          </View>
        </View>
        <View style={styles.headerDescription}>
          <AppText color={theme.colors.textSecondary}>
            Seu consolidado agora responde ao que vier dos bancos conectados.
          </AppText>
        </View>

        {overviewQuery.isLoading && !overview ? <LoadingState label="Atualizando seu painel..." /> : null}
        {overviewQuery.isError ? <ErrorState onRetry={(): void => void overviewQuery.refetch()} /> : null}

        {!overviewQuery.isLoading && !overviewQuery.isError ? (
          <>
            <Card style={styles.balanceCard}>
              <View pointerEvents="none" style={styles.balanceGlow} />
              <View style={styles.balanceTop}>
                <View>
                  <AppText color={theme.colors.textSecondary} variant="label">
                    Saldo disponível
                  </AppText>
                  <AppText variant="display">{totalBalanceLabel}</AppText>
                </View>
                <View style={styles.deltaChip}>
                  <Sparkles color={theme.colors.textInverse} size={14} strokeWidth={2} />
                  <AppText color={theme.colors.textInverse} variant="label">
                    {formatFlowLabel(netFlow)}
                  </AppText>
                </View>
              </View>
              <AppText color={theme.colors.textSecondary}>
                {buildBalanceDescription(overview?.connected_accounts ?? 0, netFlow)}
              </AppText>
              <View style={styles.actionsRow}>
                {QUICK_ACTIONS.map((action: QuickAction): ReactElement => {
                  const Icon = action.icon;

                  return (
                    <View key={action.label} style={styles.actionItem}>
                      <View style={styles.actionOrb}>
                        <Icon color={theme.colors.textPrimary} size={18} strokeWidth={2} />
                      </View>
                      <AppText color={theme.colors.textSecondary}>{action.label}</AppText>
                    </View>
                  );
                })}
              </View>
            </Card>

            <View style={styles.metricsRow}>
              {insightItems.map(
                (item: InsightItem): ReactElement => (
                  <Card key={item.label} style={styles.metricCard}>
                    <AppText color={item.tone} variant="label">
                      {item.label}
                    </AppText>
                    <AppText variant="headline">{item.amount}</AppText>
                  </Card>
                ),
              )}
            </View>

            <Card style={styles.cardPreview}>
              <View pointerEvents="none" style={styles.cardHighlight} />
              <AppText color={theme.colors.textSecondary} variant="label">
                Conta principal
              </AppText>
              {firstAccount ? (
                <>
                  <View style={styles.accountHeader}>
                    <View style={styles.cardChip} />
                    <AppText color={theme.colors.textSecondary}>{firstAccount.institution_name}</AppText>
                  </View>
                  <View style={styles.cardNumbers}>
                    <AppText variant="headline">{resolveAccountTitle(firstAccount)}</AppText>
                    <AppText color={theme.colors.textSecondary}>
                      {firstAccount.number ? `Final ${maskAccountNumber(firstAccount.number)}` : 'Conta sincronizada'}
                    </AppText>
                  </View>
                  <View style={styles.cardFooter}>
                    <View>
                      <AppText color={theme.colors.textSecondary} variant="label">
                        Saldo atual
                      </AppText>
                      <AppText>{formatCurrency(firstAccount.balance)}</AppText>
                    </View>
                    <AppText variant="headline">{normalizeAccountType(firstAccount.type)}</AppText>
                  </View>
                </>
              ) : (
                <AppText color={theme.colors.textSecondary}>
                  Conecte um banco para começar a acompanhar saldo, categorias e movimentação sem entrada manual.
                </AppText>
              )}
            </Card>

            {connectionsQuery.isLoading ? <LoadingState label="Carregando conexões..." /> : null}
            {connectionsQuery.isError ? <ErrorState onRetry={(): void => void connectionsQuery.refetch()} /> : null}
            {!connectionsQuery.isLoading && !connectionsQuery.isError ? (
              <Card style={styles.networkCard}>
                <View style={styles.networkTop}>
                  <View>
                    <AppText color={theme.colors.accent} variant="label">
                      Open Finance
                    </AppText>
                    <AppText variant="headline">Rede conectada</AppText>
                  </View>
                  <View style={styles.networkBadge}>
                    <AppText>{String(activeConnections).padStart(2, '0')} ativa(s)</AppText>
                  </View>
                </View>
                <AppText color={theme.colors.textSecondary}>
                  {pendingConnections > 0
                    ? `${pendingConnections} conexão(ões) ainda dependem de nova autorização ou sincronização.`
                    : (overview?.connections_with_data ?? 0) > 0
                      ? 'Os bancos conectados já estão entregando dados para o painel.'
                      : 'Conecte seu primeiro banco para começar o consolidado automático.'}
                </AppText>
              </Card>
            ) : null}

            <CategoryPieChart
              emptyDescription="Suas categorias aparecem aqui assim que houver transações sincronizadas."
              slices={chartSlices}
            />

            <Card>
              <View style={styles.activityHeader}>
                <View>
                  <AppText color={theme.colors.accent} variant="label">
                    Movimento
                  </AppText>
                  <AppText variant="headline">Últimas transações</AppText>
                </View>
                <AppText color={theme.colors.textSecondary}>Atualizado</AppText>
              </View>
              {recentTransactions.length > 0 ? (
                <View style={styles.activityList}>
                  {recentTransactions.map(
                    (item: FinancialTransaction): ReactElement => (
                      <View key={item.id} style={styles.activityRow}>
                        <View
                          style={[
                            styles.activityIcon,
                            item.signed_amount >= 0 ? styles.incomeIcon : styles.expenseIcon,
                          ]}
                        >
                          {item.signed_amount >= 0 ? (
                            <ArrowDownLeft color={theme.colors.textInverse} size={14} strokeWidth={2} />
                          ) : (
                            <ArrowUpRight color={theme.colors.textInverse} size={14} strokeWidth={2} />
                          )}
                        </View>
                        <View style={styles.activityMeta}>
                          <AppText>{item.description}</AppText>
                          <AppText color={theme.colors.textSecondary}>
                            {buildTransactionSubtitle(item)}
                          </AppText>
                        </View>
                        <AppText>{formatCurrency(item.signed_amount)}</AppText>
                      </View>
                    ),
                  )}
                </View>
              ) : (
                <AppText color={theme.colors.textSecondary}>
                  Assim que a sincronização trouxer movimentações, elas aparecem aqui em ordem recente.
                </AppText>
              )}
            </Card>
          </>
        ) : null}

        <View style={styles.links}>
          <LinkCard description="Gerencie instituições conectadas e sincronização." href="/(private)/connections" title="Bancos conectados" />
          <LinkCard description="Veja a timeline real de despesas e categorias do período." href="/(private)/transactions" title="Despesas do mês" />
          <LinkCard description="Segurança, notificações e sessão." href="/(private)/profile" title="Perfil e proteção" />
        </View>
      </View>
    </Screen>
  );
}

function buildInsightItems(overview: FinancialOverview | undefined): readonly InsightItem[] {
  return [
    {
      amount: formatCurrency(overview?.month_income ?? 0),
      label: 'Entradas do mês',
      tone: theme.colors.success,
    },
    {
      amount: formatCurrency(overview?.month_expenses ?? 0),
      label: 'Saídas do mês',
      tone: '#78A7FF',
    },
    {
      amount:
        (overview?.credit_card_balance ?? 0) > 0
          ? formatCurrency(overview?.credit_card_balance ?? 0)
          : String(overview?.connected_accounts ?? 0).padStart(2, '0'),
      label: (overview?.credit_card_balance ?? 0) > 0 ? 'Fatura em aberto' : 'Contas visíveis',
      tone: theme.colors.primary,
    },
  ] as const;
}

function buildCategorySlices(breakdown: readonly CategoryBreakdown[]): readonly CategorySlice[] {
  return breakdown.map(
    (slice: CategoryBreakdown, index: number): CategorySlice => ({
      amountLabel: formatCurrency(slice.amount),
      color: CHART_COLORS[index % CHART_COLORS.length] ?? theme.colors.primary,
      name: slice.category,
      percent: slice.percent,
    }),
  );
}

function formatFlowLabel(value: number): string {
  if (value === 0) {
    return 'Fluxo neutro';
  }

  const prefix = value > 0 ? '+' : '-';
  return `${prefix}${formatCurrency(Math.abs(value))}`;
}

function buildBalanceDescription(connectedAccounts: number, netFlow: number): string {
  if (connectedAccounts === 0) {
    return 'Nenhuma conta conectada ainda. O painel começa a consolidar saldo e histórico assim que seu primeiro banco for autorizado.';
  }

  if (netFlow > 0) {
    return 'O mês está fechando com fluxo positivo até aqui. Continue sincronizando para manter o consolidado atualizado.';
  }

  if (netFlow < 0) {
    return 'As saídas estão acima das entradas no mês. O painel mostra isso em tempo real conforme os bancos sincronizam.';
  }

  return 'As movimentações deste mês ainda estão equilibradas entre entrada e saída.';
}

function resolveAccountTitle(account: AccountSnapshot): string {
  return account.marketing_name ?? account.name;
}

function normalizeAccountType(accountType: string): string {
  switch (accountType) {
    case 'BANK':
      return 'Conta';
    case 'CREDIT':
      return 'Crédito';
    default:
      return 'Saldo';
  }
}

function maskAccountNumber(value: string): string {
  const digits = value.replace(/\D/g, '');
  if (digits.length === 0) {
    return '0000';
  }
  return digits.slice(-4);
}

function buildTransactionSubtitle(transaction: FinancialTransaction): string {
  const date = new Date(transaction.date);
  const label = Number.isNaN(date.getTime())
    ? transaction.institution_name
    : date.toLocaleDateString('pt-BR', {
        day: '2-digit',
        month: '2-digit',
      });

  return `${transaction.category} · ${label}`;
}

const styles = StyleSheet.create({
  content: {
    gap: theme.spacing.lg,
  },
  header: {
    alignItems: 'flex-start',
    flexDirection: 'row',
    gap: theme.spacing.md,
    justifyContent: 'space-between',
  },
  headerText: {
    flex: 1,
    gap: theme.spacing.sm,
  },
  eyebrow: {
    alignItems: 'center',
    alignSelf: 'flex-start',
    backgroundColor: theme.colors.surfaceMuted,
    borderColor: theme.colors.border,
    borderRadius: theme.radii.pill,
    borderWidth: 1,
    flexDirection: 'row',
    gap: theme.spacing.sm,
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.sm,
  },
  eyebrowDot: {
    backgroundColor: theme.colors.primary,
    borderRadius: theme.radii.pill,
    height: 8,
    width: 8,
  },
  headerAction: {
    alignItems: 'center',
    backgroundColor: theme.colors.surfaceMuted,
    borderColor: theme.colors.border,
    borderRadius: theme.radii.pill,
    borderWidth: 1,
    height: 42,
    justifyContent: 'center',
    width: 42,
  },
  headerDescription: {
    marginTop: -theme.spacing.md + theme.spacing.sm,
  },
  balanceCard: {
    backgroundColor: theme.colors.surfaceInverse,
    gap: theme.spacing.md,
    overflow: 'hidden',
  },
  balanceGlow: {
    backgroundColor: theme.colors.primary,
    borderRadius: 160,
    height: 160,
    opacity: 0.12,
    position: 'absolute',
    right: -30,
    top: -50,
    width: 160,
  },
  balanceTop: {
    alignItems: 'flex-start',
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  deltaChip: {
    alignItems: 'center',
    backgroundColor: theme.colors.primary,
    borderRadius: theme.radii.pill,
    flexDirection: 'row',
    gap: 6,
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.sm,
  },
  actionsRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  actionItem: {
    alignItems: 'center',
    gap: theme.spacing.sm,
  },
  actionOrb: {
    alignItems: 'center',
    backgroundColor: theme.colors.surface,
    borderColor: theme.colors.border,
    borderRadius: theme.radii.pill,
    borderWidth: 1,
    height: 52,
    justifyContent: 'center',
    width: 52,
  },
  metricsRow: {
    gap: theme.spacing.md,
  },
  metricCard: {
    gap: theme.spacing.sm,
    paddingBottom: 18,
  },
  cardPreview: {
    backgroundColor: theme.colors.surface,
    gap: theme.spacing.sm,
    overflow: 'hidden',
  },
  cardHighlight: {
    backgroundColor: '#95BDFF',
    borderRadius: 220,
    height: 170,
    opacity: 0.12,
    position: 'absolute',
    right: -50,
    top: -70,
    width: 170,
  },
  accountHeader: {
    alignItems: 'center',
    flexDirection: 'row',
    gap: theme.spacing.sm,
  },
  cardChip: {
    backgroundColor: theme.colors.primary,
    borderRadius: 10,
    height: 28,
    width: 40,
  },
  cardNumbers: {
    gap: theme.spacing.xs,
  },
  cardFooter: {
    alignItems: 'flex-end',
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  networkCard: {
    gap: theme.spacing.sm,
  },
  networkTop: {
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  networkBadge: {
    backgroundColor: theme.colors.surfaceMuted,
    borderRadius: theme.radii.pill,
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.sm,
  },
  activityHeader: {
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: theme.spacing.md,
  },
  activityList: {
    gap: theme.spacing.md,
  },
  activityRow: {
    alignItems: 'center',
    flexDirection: 'row',
    gap: theme.spacing.md,
  },
  activityIcon: {
    alignItems: 'center',
    borderRadius: theme.radii.pill,
    height: 36,
    justifyContent: 'center',
    width: 36,
  },
  incomeIcon: {
    backgroundColor: theme.colors.success,
  },
  expenseIcon: {
    backgroundColor: '#FF8D97',
  },
  activityMeta: {
    flex: 1,
    gap: 2,
  },
  links: {
    gap: theme.spacing.md,
  },
});
