import type { ReactElement } from 'react';
import { StyleSheet, View, type ViewProps } from 'react-native';

import { theme } from '@/theme/tokens';

export function Card(props: ViewProps): ReactElement {
  const { children, style, ...rest } = props;

  return (
    <View {...rest} style={[styles.card, style]}>
      {children}
    </View>
  );
}

const styles = StyleSheet.create({
  card: {
    backgroundColor: theme.colors.surface,
    borderColor: theme.colors.border,
    borderRadius: theme.radii.xl,
    borderWidth: 1,
    padding: theme.spacing.lg,
    shadowColor: '#000000',
    shadowOffset: {
      width: 0,
      height: 18,
    },
    shadowOpacity: 0.32,
    shadowRadius: 32,
  },
});
