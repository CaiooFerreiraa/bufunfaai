import type { PropsWithChildren, ReactElement } from 'react';
import { StyleSheet, View } from 'react-native';

import { AppText } from '@/components/ui/AppText';
import { Screen } from '@/components/ui/Screen';
import { theme } from '@/theme/tokens';

interface FeatureScreenProps extends PropsWithChildren {
  readonly description?: string;
  readonly title: string;
}

export function FeatureScreen(props: FeatureScreenProps): ReactElement {
  const { children, description, title } = props;

  return (
    <Screen>
      <View style={styles.content}>
        <View style={styles.header}>
          <View style={styles.badge}>
            <View style={styles.badgeDot} />
            <AppText color={theme.colors.textPrimary} variant="label">
              BufunfaAI finance OS
            </AppText>
          </View>
          <AppText variant="display">{title}</AppText>
          {description ? <AppText color={theme.colors.textSecondary}>{description}</AppText> : null}
        </View>
        {children}
      </View>
    </Screen>
  );
}

const styles = StyleSheet.create({
  content: {
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
});
