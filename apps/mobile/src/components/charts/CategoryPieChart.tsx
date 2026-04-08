import type { ReactElement } from 'react';
import { StyleSheet, View } from 'react-native';

import { AppText } from '@/components/ui/AppText';
import { Card } from '@/components/ui/Card';
import { theme } from '@/theme/tokens';

export interface CategorySlice {
  readonly amountLabel: string;
  readonly color: string;
  readonly name: string;
  readonly percent: number;
}

interface CategoryPieChartProps {
  readonly emptyDescription?: string;
  readonly slices: readonly CategorySlice[];
  readonly title?: string;
}

export function CategoryPieChart(props: CategoryPieChartProps): ReactElement {
  const {
    emptyDescription = 'Conecte um banco e aguarde a primeira sincronização para ver a composição das saídas.',
    slices,
    title = 'Pressão de saída',
  } = props;

  const mappedPercent: number = Math.round(slices.reduce((total, slice) => total + slice.percent, 0));

  return (
    <Card>
      <View style={styles.wrapper}>
        <View style={styles.header}>
          <View>
            <AppText color={theme.colors.accent} variant="label">
              Analytics
            </AppText>
            <AppText variant="headline">{title}</AppText>
          </View>
          <View style={styles.deltaBadge}>
            <AppText color={theme.colors.textInverse} variant="label">
              {mappedPercent}% mapeado
            </AppText>
          </View>
        </View>
        {slices.length > 0 ? (
          <>
            <View style={styles.chartShell}>
              {slices.map(
                (slice: CategorySlice): ReactElement => (
                  <View
                    key={slice.name}
                    style={[
                      styles.slice,
                      {
                        backgroundColor: slice.color,
                        width: `${Math.max(slice.percent, 6)}%`,
                      },
                    ]}
                  />
                ),
              )}
            </View>
            <View style={styles.legend}>
              <AppText color={theme.colors.textSecondary}>
                Leitura direta das categorias que mais estão pesando no período.
              </AppText>
              <View style={styles.slicesList}>
                {slices.map(
                  (slice: CategorySlice): ReactElement => (
                    <View key={slice.name} style={styles.row}>
                      <View style={styles.rowLabel}>
                        <View style={[styles.dot, { backgroundColor: slice.color }]} />
                        <View style={styles.labelStack}>
                          <AppText>{slice.name}</AppText>
                          <AppText color={theme.colors.textSecondary}>
                            {Math.round(slice.percent)}% do total
                          </AppText>
                        </View>
                      </View>
                      <View style={styles.rowValues}>
                        <AppText>{slice.amountLabel}</AppText>
                      </View>
                    </View>
                  ),
                )}
              </View>
            </View>
          </>
        ) : (
          <AppText color={theme.colors.textSecondary}>{emptyDescription}</AppText>
        )}
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
    flexDirection: 'row',
    gap: theme.spacing.lg,
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
