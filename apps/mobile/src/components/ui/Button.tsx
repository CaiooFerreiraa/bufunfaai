import type { ReactElement } from 'react';
import { Pressable, StyleSheet, type ViewStyle } from 'react-native';

import { AppText } from '@/components/ui/AppText';
import { theme } from '@/theme/tokens';

interface ButtonProps {
  readonly label: string;
  readonly onPress: () => void;
  readonly disabled?: boolean;
  readonly style?: ViewStyle;
  readonly variant?: 'primary' | 'secondary';
}

export function Button(props: ButtonProps): ReactElement {
  const { disabled = false, label, onPress, style, variant = 'primary' } = props;
  const isSecondary: boolean = variant === 'secondary';

  return (
    <Pressable
      accessibilityRole="button"
      disabled={disabled}
      onPress={onPress}
      style={({ pressed }): ViewStyle[] => [
        styles.base,
        isSecondary ? styles.secondary : styles.primary,
        ...(pressed ? [isSecondary ? styles.pressedSecondary : styles.pressedPrimary] : []),
        disabled ? styles.disabled : styles.enabled,
        style ?? {},
      ]}
    >
      <AppText color={isSecondary ? theme.colors.textPrimary : theme.colors.textInverse} variant="label">
        {label}
      </AppText>
    </Pressable>
  );
}

const styles = StyleSheet.create({
  base: {
    alignItems: 'center',
    borderRadius: theme.radii.pill,
    borderWidth: 1,
    paddingHorizontal: theme.spacing.lg,
    paddingVertical: 15,
  },
  primary: {
    backgroundColor: theme.colors.primary,
    borderColor: theme.colors.primary,
  },
  secondary: {
    backgroundColor: theme.colors.surfaceMuted,
    borderColor: theme.colors.border,
  },
  disabled: {
    opacity: 0.45,
  },
  enabled: {
    opacity: 1,
  },
  pressedPrimary: {
    backgroundColor: theme.colors.accent,
    borderColor: theme.colors.accent,
    transform: [{ scale: 0.985 }],
  },
  pressedSecondary: {
    backgroundColor: theme.colors.surfaceInverse,
    transform: [{ scale: 0.985 }],
  },
});
