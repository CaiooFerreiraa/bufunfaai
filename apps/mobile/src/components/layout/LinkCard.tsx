import { Link } from 'expo-router';
import type { Href } from 'expo-router';
import { ArrowUpRight } from 'lucide-react-native';
import type { ReactElement } from 'react';
import { Pressable, StyleSheet, View } from 'react-native';

import { AppText } from '@/components/ui/AppText';
import { theme } from '@/theme/tokens';

interface LinkCardProps {
  readonly description: string;
  readonly href: Href;
  readonly title: string;
}

export function LinkCard(props: LinkCardProps): ReactElement {
  const { description, href, title } = props;

  return (
    <Link asChild href={href}>
      <Pressable style={({ pressed }): object[] => [styles.card, pressed ? styles.pressed : {}]}>
        <View style={styles.content}>
          <View style={styles.topRow}>
            <View style={styles.badge}>
              <View style={styles.badgeDot} />
              <AppText color={theme.colors.textPrimary} variant="label">
                Atalho
              </AppText>
            </View>
            <View style={styles.iconWrap}>
              <ArrowUpRight color={theme.colors.textPrimary} size={16} strokeWidth={2} />
            </View>
          </View>
          <AppText variant="headline">{title}</AppText>
          <AppText color={theme.colors.textSecondary}>{description}</AppText>
        </View>
      </Pressable>
    </Link>
  );
}

const styles = StyleSheet.create({
  badge: {
    alignItems: 'center',
    backgroundColor: theme.colors.surfaceInverse,
    borderColor: theme.colors.border,
    borderRadius: theme.radii.pill,
    borderWidth: 1,
    flexDirection: 'row',
    gap: theme.spacing.sm,
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.sm,
  },
  badgeDot: {
    backgroundColor: theme.colors.primary,
    borderRadius: theme.radii.pill,
    height: 8,
    width: 8,
  },
  card: {
    backgroundColor: theme.colors.surface,
    borderColor: theme.colors.border,
    borderRadius: theme.radii.xl,
    borderWidth: 1,
    overflow: 'hidden',
  },
  content: {
    gap: theme.spacing.sm,
    padding: theme.spacing.lg,
  },
  iconWrap: {
    alignItems: 'center',
    backgroundColor: theme.colors.surfaceInverse,
    borderColor: theme.colors.border,
    borderRadius: theme.radii.pill,
    borderWidth: 1,
    height: 34,
    justifyContent: 'center',
    width: 34,
  },
  pressed: {
    opacity: 0.92,
    transform: [{ scale: 0.99 }],
  },
  topRow: {
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
});
