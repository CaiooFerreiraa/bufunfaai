import type { ReactElement } from 'react';
import { StyleSheet, View } from 'react-native';

import { AppText } from '@/components/ui/AppText';
import { Card } from '@/components/ui/Card';
import { theme } from '@/theme/tokens';

interface ExpenseSlice {
  readonly amountLabel: string;
  readonly color: string;
  readonly name: string;
  readonly percent: number;
  readonly trendLabel: string;
}

const EXPENSE_SLICES: readonly ExpenseSlice[] = [
  {
    amountLabel: 'R$ 2.140',
    color: theme.colors.primary,
    name: 'Moradia',
    percent: 42,
    trendLabel: '+4.2%',
  },
  {
    amountLabel: 'R$ 1.180',
    color: '#78A7FF',
    name: 'Cartão',
    percent: 26,
    trendLabel: '-2.1%',
  },
  {
    amountLabel: 'R$ 890',
    color: '#37D7A3',
    name: 'Alimentação',
    percent: 18,
    trendLabel: '+8.0%',
  },
  {
    amountLabel: 'R$ 520',
    color: '#F8A94C',
    name: 'Mobilidade',
    percent: 14,
    trendLabel: '-1.3%',
  },
] as const;

export function CategoryPieChart(): ReactElement {
  return (
    <Card>
      <View style={styles.wrapper}>
        <View style={styles.header}>
          <View>
            <AppText color={theme.colors.accent} variant="label">
              Analytics
            </AppText>
            <AppText variant="headline">Pressão de saída</AppText>
          </View>
          <View style={styles.deltaBadge}>
            <AppText color={theme.colors.textInverse} variant="label">
              72% mapeado
            </AppText>
          </View>
        </View>
        <View style={styles.chartShell}>
          {EXPENSE_SLICES.map((slice: ExpenseSlice): ReactElement => (
            <View
              key={slice.name}
              style={[
                styles.slice,
                {
                  backgroundColor: slice.color,
                  width: `${slice.percent}%`,
                },
              ]}
            />
          ))}
        </View>
        <View style={styles.legend}>
          <AppText color={theme.colors.textSecondary}>
            Uma leitura compacta das categorias que mais comprimem o saldo neste ciclo.
          </AppText>
          <View style={styles.slicesList}>
            {EXPENSE_SLICES.map((slice: ExpenseSlice): ReactElement => (
              <View key={slice.name} style={styles.row}>
                <View style={styles.rowLabel}>
                  <View style={[styles.dot, { backgroundColor: slice.color }]} />
                  <View style={styles.labelStack}>
                    <AppText>{slice.name}</AppText>
                    <AppText color={theme.colors.textSecondary}>{slice.percent}% do total</AppText>
                  </View>
                </View>
                <View style={styles.rowValues}>
                  <AppText>{slice.amountLabel}</AppText>
                  <AppText color={theme.colors.textSecondary}>{slice.trendLabel}</AppText>
                </View>
              </View>
            ))}
          </View>
        </View>
      </View>
    </Card>
  );
}

const styles = StyleSheet.create({
  wrapper: {
    gap: theme.spacing.lg,
  },
  header: {
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  deltaBadge: {
    backgroundColor: theme.colors.primary,
    borderRadius: theme.radii.pill,
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.sm,
  },
  chartShell: {
    backgroundColor: theme.colors.surfaceMuted,
    borderRadius: theme.radii.pill,
    flexDirection: 'row',
    height: 22,
    overflow: 'hidden',
    width: '100%',
  },
  slice: {
    height: '100%',
  },
  legend: {
    gap: theme.spacing.sm,
  },
  slicesList: {
    gap: theme.spacing.sm,
    marginTop: theme.spacing.sm,
  },
  row: {
    alignItems: 'center',
    gap: theme.spacing.lg,
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  labelStack: {
    gap: 2,
  },
  rowLabel: {
    alignItems: 'center',
    flexDirection: 'row',
    gap: theme.spacing.sm,
  },
  rowValues: {
    alignItems: 'flex-end',
    gap: 2,
  },
  dot: {
    borderRadius: 999,
    height: 10,
    width: 10,
  },
});
