import type { ReactElement } from 'react';
import { StyleSheet, TextInput, type TextInputProps, View } from 'react-native';

import { AppText } from '@/components/ui/AppText';
import { theme } from '@/theme/tokens';

interface TextFieldProps extends TextInputProps {
  readonly label: string;
}

export function TextField(props: TextFieldProps): ReactElement {
  const { label, style, ...rest } = props;

  return (
    <View style={styles.wrapper}>
      <AppText variant="label">{label}</AppText>
      <TextInput
        placeholderTextColor={theme.colors.textSecondary}
        style={[styles.input, style]}
        {...rest}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  wrapper: {
    gap: theme.spacing.sm,
  },
  input: {
    backgroundColor: theme.colors.surfaceMuted,
    borderColor: theme.colors.border,
    borderRadius: theme.radii.lg,
    borderWidth: 1,
    color: theme.colors.textPrimary,
    fontFamily: theme.fonts.body,
    minHeight: 52,
    paddingHorizontal: theme.spacing.lg,
    paddingVertical: 14,
  },
});
