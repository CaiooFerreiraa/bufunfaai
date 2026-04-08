import type { ReactElement } from 'react';
import { StyleSheet, Text, type TextProps } from 'react-native';

import { theme } from '@/theme/tokens';

type TextVariant = 'display' | 'headline' | 'body' | 'label';

interface AppTextProps extends TextProps {
  readonly variant?: TextVariant;
  readonly color?: string;
}

export function AppText(props: AppTextProps): ReactElement {
  const { children, color = theme.colors.textPrimary, style, variant = 'body', ...rest } = props;

  return (
    <Text {...rest} style={[styles.base, styles[variant], { color }, style]}>
      {children}
    </Text>
  );
}

const styles = StyleSheet.create({
  base: {
    color: theme.colors.textPrimary,
    fontFamily: theme.fonts.body,
  },
  body: {
    fontSize: 15,
    letterSpacing: 0.1,
    lineHeight: 23,
  },
  display: {
    fontFamily: theme.fonts.display,
    fontSize: 36,
    fontWeight: '700',
    letterSpacing: -0.9,
    lineHeight: 40,
  },
  headline: {
    fontFamily: theme.fonts.display,
    fontSize: 24,
    fontWeight: '700',
    letterSpacing: -0.4,
    lineHeight: 30,
  },
  label: {
    fontFamily: theme.fonts.label,
    fontSize: 12,
    fontWeight: '700',
    letterSpacing: 1.5,
    lineHeight: 16,
    textTransform: 'uppercase',
  },
});
