import type { ReactElement } from 'react';
import { StyleSheet, View } from 'react-native';

import { AppText } from '@/components/ui/AppText';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { theme } from '@/theme/tokens';

interface EmptyStateProps {
  readonly actionLabel?: string;
  readonly description: string;
  readonly onActionPress?: () => void;
  readonly title: string;
}

export function EmptyState(props: EmptyStateProps): ReactElement {
  const { actionLabel, description, onActionPress, title } = props;

  return (
    <Card>
      <View style={styles.content}>
        <AppText color={theme.colors.accent} variant="label">
          Estado inicial
        </AppText>
        <AppText variant="headline">{title}</AppText>
        <AppText color={theme.colors.textSecondary}>{description}</AppText>
        {actionLabel && onActionPress ? <Button label={actionLabel} onPress={onActionPress} style={styles.action} /> : null}
      </View>
    </Card>
  );
}

const styles = StyleSheet.create({
  content: {
    gap: theme.spacing.md,
  },
  action: {
    marginTop: theme.spacing.sm,
  },
});
