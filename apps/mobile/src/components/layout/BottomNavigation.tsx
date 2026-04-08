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
            <Icon
              color={isActive ? theme.colors.textInverse : theme.colors.textSecondary}
              size={18}
              strokeWidth={2}
            />
            <AppText
              color={isActive ? theme.colors.textInverse : theme.colors.textSecondary}
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
    backgroundColor: '#0B101B',
    borderColor: theme.colors.border,
    borderRadius: 30,
    borderWidth: 1,
    flexDirection: 'row',
    gap: theme.spacing.sm,
    padding: 10,
    shadowColor: '#000000',
    shadowOffset: {
      width: 0,
      height: 18,
    },
    shadowOpacity: 0.4,
    shadowRadius: 30,
  },
  item: {
    alignItems: 'center',
    borderRadius: 20,
    flex: 1,
    gap: 6,
    minHeight: 56,
    justifyContent: 'center',
    paddingHorizontal: theme.spacing.sm,
  },
  itemActive: {
    backgroundColor: theme.colors.primary,
  },
  itemIdle: {
    backgroundColor: 'transparent',
  },
  itemPressed: {
    opacity: 0.92,
  },
});
