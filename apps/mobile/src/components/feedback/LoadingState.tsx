import type { ReactElement } from 'react';
import { ActivityIndicator, StyleSheet, View } from 'react-native';

import { AppText } from '@/components/ui/AppText';
import { theme } from '@/theme/tokens';

interface LoadingStateProps {
  readonly label?: string;
}

export function LoadingState(props: LoadingStateProps): ReactElement {
  const { label = 'Carregando...' } = props;

  return (
    <View style={styles.container}>
      <ActivityIndicator color={theme.colors.accent} size="large" />
      <AppText color={theme.colors.textSecondary}>{label}</AppText>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    alignItems: 'center',
    gap: theme.spacing.md,
    justifyContent: 'center',
    paddingVertical: theme.spacing['2xl'],
  },
});
