import { usePathname, useRouter } from 'expo-router';
import type { Href } from 'expo-router';
import {
  ArrowLeftRight,
  House,
  Landmark,
  UserRound,
} from 'lucide-react-native';
import type { ReactElement } from 'react';
import { Pressable, StyleSheet, View } from 'react-native';

import { AppText } from '@/components/ui/AppText';
import { theme } from '@/theme/tokens';

interface NavItem {
  readonly href: Href;
  readonly icon: typeof House;
  readonly label: string;
  readonly match: (pathname: string) => boolean;
}

const NAV_ITEMS: readonly NavItem[] = [
  {
    href: '/(private)/home',
    icon: House,
    label: 'Início',
    match: (pathname: string): boolean => pathname === '/home',
  },
  {
    href: '/(private)/connections',
    icon: Landmark,
    label: 'Bancos',
    match: (pathname: string): boolean => pathname.startsWith('/connections'),
  },
  {
    href: '/(private)/transactions',
    icon: ArrowLeftRight,
    label: 'Despesas',
    match: (pathname: string): boolean => pathname.startsWith('/transactions'),
  },
  {
    href: '/(private)/profile',
    icon: UserRound,
    label: 'Perfil',
    match: (pathname: string): boolean => pathname.startsWith('/profile'),
  },
] as const;

export function BottomNavigation(): ReactElement {
  const pathname: string = usePathname();
  const router = useRouter();

  return (
    <View style={styles.shell}>
      {NAV_ITEMS.map((item: NavItem): ReactElement => {
        const isActive: boolean = item.match(pathname);
        const Icon = item.icon;

        return (
          <Pressable
            accessibilityRole="button"
            key={item.label}
            onPress={(): void => router.replace(item.href)}
            style={({ pressed }): object[] => [
              styles.item,
              isActive ? styles.itemActive : styles.itemIdle,
              pressed ? styles.itemPressed : {},
            ]}
          >
            <View style={[styles.activeIndicator, isActive ? styles.activeIndicatorVisible : styles.activeIndicatorHidden]} />
            <Icon
              color={isActive ? '#25D366' : theme.colors.textSecondary}
              size={20}
              strokeWidth={2}
            />
            <AppText
              color={isActive ? '#25D366' : theme.colors.textSecondary}
              variant="label"
            >
              {item.label}
            </AppText>
          </Pressable>
        );
      })}
    </View>
  );
}

const styles = StyleSheet.create({
  shell: {
    backgroundColor: theme.colors.surface,
    borderTopColor: theme.colors.border,
    borderTopWidth: 1,
    flexDirection: 'row',
    paddingTop: 4,
    paddingHorizontal: theme.spacing.xs,
  },
  item: {
    alignItems: 'center',
    flex: 1,
    gap: 4,
    minHeight: 58,
    justifyContent: 'center',
    paddingHorizontal: theme.spacing.sm,
    paddingBottom: theme.spacing.sm,
    paddingTop: theme.spacing.sm,
  },
  itemActive: {
    backgroundColor: 'transparent',
  },
  itemIdle: {
    backgroundColor: 'transparent',
  },
  itemPressed: {
    opacity: 0.92,
  },
  activeIndicator: {
    borderRadius: theme.radii.pill,
    height: 3,
    marginBottom: 8,
    width: 28,
  },
  activeIndicatorVisible: {
    backgroundColor: '#25D366',
    opacity: 1,
  },
  activeIndicatorHidden: {
    backgroundColor: 'transparent',
    opacity: 0,
  },
});
