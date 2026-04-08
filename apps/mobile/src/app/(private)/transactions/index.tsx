import type { ReactElement } from 'react';
import { StyleSheet, View } from 'react-native';

import { CategoryPieChart } from '@/components/charts/CategoryPieChart';
import { FeatureScreen } from '@/components/layout/FeatureScreen';
import { AppText } from '@/components/ui/AppText';
import { Card } from '@/components/ui/Card';
import { theme } from '@/theme/tokens';

export default function TransactionsScreen(): ReactElement {
  const transactions: readonly {
    readonly amount: string;
    readonly category: string;
    readonly indicator: string;
    readonly title: string;
  }[] = [
    { amount: '-R$ 482,40', category: 'Alimentação', indicator: '#37D7A3', title: 'Mercado do bairro' },
    { amount: '-R$ 261,20', category: 'Mobilidade', indicator: '#78A7FF', title: 'Corridas e deslocamentos' },
    { amount: '-R$ 134,90', category: 'Assinaturas', indicator: '#F8A94C', title: 'Serviços recorrentes' },
  ] as const;

  return (
    <FeatureScreen
      description="Leitura rápida do que mais pesa no mês, com foco em categorias e ritmo de saída."
      title="Despesas"
    >
      <Card style={styles.heroCard}>
        <View style={styles.heroTop}>
          <View>
            <AppText color={theme.colors.accent} variant="label">
              Saída acumulada
            </AppText>
            <AppText variant="display">R$ 4.730,18</AppText>
          </View>
          <View style={styles.heroBadge}>
            <AppText color={theme.colors.textInverse} variant="label">
              -8.4% vs. mês anterior
            </AppText>
          </View>
        </View>
        <AppText color={theme.colors.textSecondary}>
          A maior compressão continua em moradia, cartão e alimentação. Seu ritmo está mais controlado que no mês passado.
        </AppText>
      </Card>

      <CategoryPieChart />

      <Card>
        <View style={styles.summary}>
          <AppText color={theme.colors.accent} variant="label">
            Leitura do período
          </AppText>
          <AppText variant="headline">Top saídas do mês</AppText>
          {transactions.map(
            (transaction): ReactElement => (
              <View key={transaction.title} style={styles.row}>
                <View style={styles.meta}>
                  <View style={styles.titleRow}>
                    <View style={[styles.indicator, { backgroundColor: transaction.indicator }]} />
                    <AppText>{transaction.title}</AppText>
                  </View>
                  <AppText color={theme.colors.textSecondary}>{transaction.category}</AppText>
                </View>
                <AppText color="#FF9AA2">{transaction.amount}</AppText>
              </View>
            ),
          )}
        </View>
      </Card>
    </FeatureScreen>
  );
}

const styles = StyleSheet.create({
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
  row: {
    alignItems: 'center',
    borderTopColor: theme.colors.border,
    borderTopWidth: 1,
    flexDirection: 'row',
    justifyContent: 'space-between',
    paddingTop: theme.spacing.md,
  },
  meta: {
    flex: 1,
    gap: 2,
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
