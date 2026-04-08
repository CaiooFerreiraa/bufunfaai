import type { PropsWithChildren, ReactElement } from 'react';
import {
  ScrollView,
  StyleSheet,
  View,
  type StyleProp,
  type ViewStyle,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';

import { theme } from '@/theme/tokens';

interface ScreenProps extends PropsWithChildren {
  readonly contentContainerStyle?: StyleProp<ViewStyle>;
  readonly scrollable?: boolean;
}

export function Screen(props: ScreenProps): ReactElement {
  const { children, contentContainerStyle, scrollable = true } = props;

  if (scrollable) {
    return (
      <SafeAreaView style={styles.screen}>
        <View pointerEvents="none" style={styles.atmosphereTop} />
        <View pointerEvents="none" style={styles.atmosphereBottom} />
        <ScrollView
          contentContainerStyle={[styles.scrollContent, contentContainerStyle]}
          keyboardShouldPersistTaps="handled"
          showsVerticalScrollIndicator={false}
        >
          {children}
        </ScrollView>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.screen}>
      <View pointerEvents="none" style={styles.atmosphereTop} />
      <View pointerEvents="none" style={styles.atmosphereBottom} />
      <View style={[styles.staticContent, contentContainerStyle]}>{children}</View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  screen: {
    backgroundColor: theme.colors.background,
    flex: 1,
    overflow: 'hidden',
  },
  scrollContent: {
    flexGrow: 1,
    paddingHorizontal: theme.spacing.lg,
    paddingTop: theme.spacing.lg,
    paddingBottom: 128,
  },
  staticContent: {
    flex: 1,
    paddingHorizontal: theme.spacing.lg,
    paddingTop: theme.spacing.lg,
    paddingBottom: 128,
  },
  atmosphereTop: {
    backgroundColor: '#C7FF45',
    borderRadius: 260,
    height: 260,
    opacity: 0.12,
    position: 'absolute',
    right: -90,
    top: -70,
    width: 260,
  },
  atmosphereBottom: {
    backgroundColor: '#2B5BFF',
    borderRadius: 320,
    bottom: -180,
    height: 320,
    left: -150,
    opacity: 0.1,
    position: 'absolute',
    width: 320,
  },
});
