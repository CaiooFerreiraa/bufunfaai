import type { ReactElement } from 'react';
import { FlatList, StyleSheet, View } from 'react-native';

import { CategoryPieChart, type CategorySlice } from '@/components/charts/CategoryPieChart';
import { EmptyState } from '@/components/feedback/EmptyState';
import { ErrorState } from '@/components/feedback/ErrorState';
import { LoadingState } from '@/components/feedback/LoadingState';
import { AppText } from '@/components/ui/AppText';
import { Card } from '@/components/ui/Card';
import { Screen } from '@/components/ui/Screen';
import { useTransactionsFeedQuery } from '@/features/open-finance/hooks/useFinancialData';
import type {
  CategoryBreakdown,
  FinancialTransaction,
} from '@/features/open-finance/types/financialData.types';
import { theme } from '@/theme/tokens';
import { formatCurrency } from '@/utils/formatCurrency';

const CHART_COLORS: readonly string[] = [
  theme.colors.primary,
  '#78A7FF',
  '#37D7A3',
  '#F8A94C',
] as const;

export default function TransactionsScreen(): ReactElement {
  const from = new Date();
  from.setUTCDate(1);
  from.setUTCMonth(from.getUTCMonth() - 1);

  const feedQuery = useTransactionsFeedQuery({
    from: from.toISOString().slice(0, 10),
    limit: 500,
  });

  const feed = feedQuery.data;
  const transactions: readonly FinancialTransaction[] = (feed?.transactions ?? []).filter(
    (transaction: FinancialTransaction): boolean => transaction.signed_amount < 0,
  );
  const chartSlices: readonly CategorySlice[] = buildCategorySlices(feed?.expense_breakdown ?? []);
  const deltaPercentage = buildExpenseDeltaPercentage(
    feed?.month_expense_total ?? 0,
    feed?.previous_month_expense_total ?? 0,
  );

  return (
    <Screen scrollable={false}>
      {feedQuery.isLoading && !feed ? <LoadingState label="Atualizando despesas..." /> : null}
      {feedQuery.isError ? <ErrorState onRetry={(): void => void feedQuery.refetch()} /> : null}

      {!feedQuery.isLoading && !feedQuery.isError ? (
        <FlatList
          contentContainerStyle={styles.content}
          data={transactions}
          keyExtractor={(item: FinancialTransaction): string => item.id}
          ListEmptyComponent={
            <EmptyState
              description="Assim que as movimentações chegarem dos bancos conectados, elas aparecem organizadas aqui."
              title="Nenhuma despesa sincronizada"
            />
          }
          ListHeaderComponent={
            <View style={styles.headerContent}>
              <View style={styles.header}>
                <View style={styles.badge}>
                  <View style={styles.badgeDot} />
                  <AppText color={theme.colors.textPrimary} variant="label">
                    BufunfaAI finance OS
                  </AppText>
                </View>
                <AppText variant="display">Despesas</AppText>
                <AppText color={theme.colors.textSecondary}>
                  Timeline real das saídas e das categorias que mais pesam no período.
                </AppText>
              </View>

              <Card style={styles.heroCard}>
                <View style={styles.heroTop}>
                  <View>
                    <AppText color={theme.colors.accent} variant="label">
                      Saída acumulada
                    </AppText>
                    <AppText variant="display">{formatCurrency(feed?.month_expense_total ?? 0)}</AppText>
                  </View>
                  <View style={styles.heroBadge}>
                    <AppText color={theme.colors.textInverse} variant="label">
                      {deltaPercentage}
                    </AppText>
                  </View>
                </View>
                <AppText color={theme.colors.textSecondary}>
                  {buildExpenseDescription(
                    feed?.month_expense_total ?? 0,
                    feed?.previous_month_expense_total ?? 0,
                    transactions.length,
                  )}
                </AppText>
              </Card>

              <CategoryPieChart
                emptyDescription="As categorias aparecem aqui assim que houver despesas sincronizadas no período."
                slices={chartSlices}
              />

              <Card>
                <View style={styles.summary}>
                  <AppText color={theme.colors.accent} variant="label">
                    Leitura do período
                  </AppText>
                  <AppText variant="headline">Últimas saídas capturadas</AppText>
                </View>
              </Card>
            </View>
          }
          renderItem={({ item }: { item: FinancialTransaction }): ReactElement => (
            <Card style={styles.transactionCard}>
              <View style={styles.row}>
                <View style={styles.meta}>
                  <View style={styles.titleRow}>
                    <View
                      style={[
                        styles.indicator,
                        {
                          backgroundColor:
                            CHART_COLORS[Math.abs(hashString(item.category)) % CHART_COLORS.length] ?? theme.colors.primary,
                        },
                      ]}
                    />
                    <AppText>{item.description}</AppText>
                  </View>
                  <AppText color={theme.colors.textSecondary}>
                    {item.category} · {item.institution_name}
                  </AppText>
                  <AppText color={theme.colors.textSecondary}>
                    {formatTransactionDate(item.date)} · {item.account_name}
                  </AppText>
                </View>
                <AppText color={theme.colors.error}>{formatCurrency(item.signed_amount)}</AppText>
              </View>
            </Card>
          )}
          showsVerticalScrollIndicator={false}
        />
      ) : null}
    </Screen>
  );
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

function buildExpenseDeltaPercentage(current: number, previous: number): string {
  if (current === 0 && previous === 0) {
    return 'Sem variação';
  }
  if (previous === 0) {
    return 'Novo período';
  }

  const delta = ((current - previous) / previous) * 100;
  const prefix = delta > 0 ? '+' : '';
  return `${prefix}${delta.toFixed(1)}% vs. mês anterior`;
}

function buildExpenseDescription(current: number, previous: number, count: number): string {
  if (count === 0) {
    return 'Ainda não há despesas sincronizadas para este recorte.';
  }
  if (previous === 0) {
    return 'Este é o primeiro período comparável com dados suficientes para o painel.';
  }
  if (current > previous) {
    return 'As saídas deste mês estão acima do período anterior.';
  }
  if (current < previous) {
    return 'As saídas deste mês estão abaixo do período anterior.';
  }
  return 'As saídas seguem no mesmo ritmo do mês anterior.';
}

function formatTransactionDate(value: string): string {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return 'Data indisponível';
  }

  return date.toLocaleDateString('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  });
}

function hashString(value: string): number {
  return Array.from(value).reduce((total: number, character: string): number => total + character.charCodeAt(0), 0);
}

const styles = StyleSheet.create({
  content: {
    gap: theme.spacing.lg,
    paddingBottom: theme.spacing['2xl'],
  },
  headerContent: {
    gap: theme.spacing.lg,
  },
  badge: {
    alignItems: 'center',
    alignSelf: 'flex-start',
    backgroundColor: theme.colors.surfaceMuted,
    borderColor: theme.colors.border,
    borderRadius: theme.radii.pill,
    borderWidth: 1,
    flexDirection: 'row',
    gap: theme.spacing.sm,
    paddingHorizontal: theme.spacing.md,
    paddingVertical: 10,
  },
  badgeDot: {
    backgroundColor: theme.colors.primary,
    borderRadius: theme.radii.pill,
    height: 8,
    width: 8,
  },
  header: {
    gap: theme.spacing.sm,
    paddingBottom: theme.spacing.sm,
  },
  heroCard: {
    gap: theme.spacing.md,
  },
  heroTop: {
    gap: theme.spacing.md,
  },
  heroBadge: {
    alignSelf: 'flex-start',
    backgroundColor: theme.colors.primary,
    borderRadius: theme.radii.pill,
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.sm,
  },
  summary: {
    gap: theme.spacing.md,
  },
  transactionCard: {
    paddingVertical: theme.spacing.md,
  },
  row: {
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  meta: {
    flex: 1,
    gap: 4,
  },
  titleRow: {
    alignItems: 'center',
    flexDirection: 'row',
    gap: theme.spacing.sm,
  },
  indicator: {
    borderRadius: theme.radii.pill,
    height: 10,
    width: 10,
  },
});
