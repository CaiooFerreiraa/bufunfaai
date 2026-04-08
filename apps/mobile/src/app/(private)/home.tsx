import {
  ArrowDownLeft,
  ArrowUpRight,
  Bell,
  CreditCard,
  ShieldCheck,
  Sparkles,
} from 'lucide-react-native';
import type { ReactElement } from 'react';
import { StyleSheet, View } from 'react-native';

import { CategoryPieChart } from '@/components/charts/CategoryPieChart';
import { ErrorState } from '@/components/feedback/ErrorState';
import { LoadingState } from '@/components/feedback/LoadingState';
import { LinkCard } from '@/components/layout/LinkCard';
import { AppText } from '@/components/ui/AppText';
import { Card } from '@/components/ui/Card';
import { Screen } from '@/components/ui/Screen';
import { useConnectionsQuery } from '@/features/connections/hooks/useConnections';
import { useDashboardSummary } from '@/features/dashboard/hooks/useDashboard';
import { theme } from '@/theme/tokens';

interface QuickAction {
  readonly icon: typeof ArrowUpRight;
  readonly label: string;
}

interface InsightItem {
  readonly amount: string;
  readonly label: string;
  readonly tone: string;
}

const QUICK_ACTIONS: readonly QuickAction[] = [
  { icon: ArrowUpRight, label: 'Enviar' },
  { icon: ArrowDownLeft, label: 'Receber' },
  { icon: CreditCard, label: 'Cartão' },
  { icon: ShieldCheck, label: 'Proteção' },
] as const;

const INSIGHT_ITEMS: readonly InsightItem[] = [
  { amount: '+R$ 4.820', label: 'Entradas do ciclo', tone: theme.colors.success },
  { amount: '-R$ 1.840', label: 'Saídas urgentes', tone: '#78A7FF' },
  { amount: '84 pts', label: 'Saúde financeira', tone: theme.colors.primary },
] as const;

const RECENT_ACTIVITY: readonly {
  readonly amount: string;
  readonly counterparty: string;
  readonly time: string;
}[] = [
  { amount: '-R$ 248,90', counterparty: 'iFood', time: 'Agora há pouco' },
  { amount: '-R$ 89,00', counterparty: 'Uber', time: 'Hoje, 14:20' },
  { amount: '+R$ 1.280,00', counterparty: 'Pix recebido', time: 'Hoje, 09:12' },
] as const;

export default function HomeScreen(): ReactElement {
  const { activeConnections, pendingConnections, totalBalanceLabel, userName } = useDashboardSummary();
  const connectionsQuery = useConnectionsQuery();

  return (
    <Screen>
      <View style={styles.content}>
        <View style={styles.header}>
          <View style={styles.headerText}>
            <View style={styles.eyebrow}>
              <View style={styles.eyebrowDot} />
              <AppText variant="label">BufunfaAI wallet</AppText>
            </View>
            <AppText variant="display">Olá, {userName}</AppText>
            <AppText color={theme.colors.textSecondary}>
              Seu painel principal foi redesenhado para leitura rápida de saldo, sinais e movimento.
            </AppText>
          </View>
          <View style={styles.headerAction}>
            <Bell color={theme.colors.textPrimary} size={18} strokeWidth={2} />
          </View>
        </View>

        <Card style={styles.balanceCard}>
          <View pointerEvents="none" style={styles.balanceGlow} />
          <View style={styles.balanceTop}>
            <View>
              <AppText color={theme.colors.textSecondary} variant="label">
                Saldo disponível
              </AppText>
              <AppText variant="display">{totalBalanceLabel}</AppText>
            </View>
            <View style={styles.deltaChip}>
              <Sparkles color={theme.colors.textInverse} size={14} strokeWidth={2} />
              <AppText color={theme.colors.textInverse} variant="label">
                +12.8%
              </AppText>
            </View>
          </View>
          <AppText color={theme.colors.textSecondary}>
            Meta do mês em trilha positiva. Seu caixa está estável e com espaço para antecipar pagamentos.
          </AppText>
          <View style={styles.actionsRow}>
            {QUICK_ACTIONS.map((action: QuickAction): ReactElement => {
              const Icon = action.icon;

              return (
                <View key={action.label} style={styles.actionItem}>
                  <View style={styles.actionOrb}>
                    <Icon color={theme.colors.textPrimary} size={18} strokeWidth={2} />
                  </View>
                  <AppText color={theme.colors.textSecondary}>{action.label}</AppText>
                </View>
              );
            })}
          </View>
        </Card>

        <View style={styles.metricsRow}>
          {INSIGHT_ITEMS.map((item: InsightItem): ReactElement => (
            <Card key={item.label} style={styles.metricCard}>
              <AppText color={item.tone} variant="label">
                {item.label}
              </AppText>
              <AppText variant="headline">{item.amount}</AppText>
            </Card>
          ))}
        </View>

        <Card style={styles.cardPreview}>
          <View pointerEvents="none" style={styles.cardHighlight} />
          <AppText color={theme.colors.textSecondary} variant="label">
            Carteira principal
          </AppText>
          <View style={styles.cardChip} />
          <View style={styles.cardNumbers}>
            <AppText variant="headline">3456 9842 1031 3277</AppText>
            <AppText color={theme.colors.textSecondary}>Expira 07/29</AppText>
          </View>
          <View style={styles.cardFooter}>
            <View>
              <AppText color={theme.colors.textSecondary} variant="label">
                Limite usado
              </AppText>
              <AppText>R$ 2.480,00</AppText>
            </View>
            <AppText variant="headline">VISA</AppText>
          </View>
        </Card>

        {connectionsQuery.isLoading ? <LoadingState label="Carregando conexoes..." /> : null}
        {connectionsQuery.isError ? <ErrorState onRetry={(): void => void connectionsQuery.refetch()} /> : null}
        {!connectionsQuery.isLoading && !connectionsQuery.isError ? (
          <Card style={styles.networkCard}>
            <View style={styles.networkTop}>
              <View>
                <AppText color={theme.colors.accent} variant="label">
                  Open Finance
                </AppText>
                <AppText variant="headline">Rede conectada</AppText>
              </View>
              <View style={styles.networkBadge}>
                <AppText>{String(activeConnections).padStart(2, '0')} ativa(s)</AppText>
              </View>
            </View>
            <AppText color={theme.colors.textSecondary}>
              {pendingConnections > 0
                ? `${pendingConnections} integração(ões) ainda aguardam autorização ou novo sync.`
                : 'Todas as integrações visíveis estão estáveis para leitura de dados.'}
            </AppText>
          </Card>
        ) : null}

        <CategoryPieChart />

        <Card>
          <View style={styles.activityHeader}>
            <View>
              <AppText color={theme.colors.accent} variant="label">
                Movimento
              </AppText>
              <AppText variant="headline">Últimas transações</AppText>
            </View>
            <AppText color={theme.colors.textSecondary}>Hoje</AppText>
          </View>
          <View style={styles.activityList}>
            {RECENT_ACTIVITY.map((item) => (
              <View key={`${item.counterparty}-${item.time}`} style={styles.activityRow}>
                <View style={styles.activityIcon}>
                  <ArrowUpRight color={theme.colors.textInverse} size={14} strokeWidth={2} />
                </View>
                <View style={styles.activityMeta}>
                  <AppText>{item.counterparty}</AppText>
                  <AppText color={theme.colors.textSecondary}>{item.time}</AppText>
                </View>
                <AppText>{item.amount}</AppText>
              </View>
            ))}
          </View>
        </Card>

        <View style={styles.links}>
          <LinkCard description="Gerencie instituições conectadas e sincronização." href="/(private)/connections" title="Bancos conectados" />
          <LinkCard description="Acompanhe timeline, filtros e separação por categoria." href="/(private)/transactions" title="Despesas do mês" />
          <LinkCard description="Segurança, notificações e sessão." href="/(private)/profile" title="Perfil e proteção" />
        </View>
      </View>
    </Screen>
  );
}

const styles = StyleSheet.create({
  content: {
    gap: theme.spacing.lg,
  },
  header: {
    alignItems: 'flex-start',
    flexDirection: 'row',
    gap: theme.spacing.lg,
    justifyContent: 'space-between',
  },
  headerText: {
    flex: 1,
    gap: theme.spacing.sm,
  },
  eyebrow: {
    alignItems: 'center',
    alignSelf: 'flex-start',
    backgroundColor: theme.colors.surfaceMuted,
    borderColor: theme.colors.border,
    borderRadius: theme.radii.pill,
    borderWidth: 1,
    flexDirection: 'row',
    gap: theme.spacing.sm,
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.sm,
  },
  eyebrowDot: {
    backgroundColor: theme.colors.primary,
    borderRadius: theme.radii.pill,
    height: 8,
    width: 8,
  },
  headerAction: {
    alignItems: 'center',
    backgroundColor: theme.colors.surfaceMuted,
    borderColor: theme.colors.border,
    borderRadius: theme.radii.pill,
    borderWidth: 1,
    height: 42,
    justifyContent: 'center',
    width: 42,
  },
  balanceCard: {
    backgroundColor: theme.colors.surfaceInverse,
    gap: theme.spacing.md,
    overflow: 'hidden',
  },
  balanceGlow: {
    backgroundColor: theme.colors.primary,
    borderRadius: 160,
    height: 160,
    opacity: 0.12,
    position: 'absolute',
    right: -30,
    top: -50,
    width: 160,
  },
  balanceTop: {
    alignItems: 'flex-start',
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  deltaChip: {
    alignItems: 'center',
    backgroundColor: theme.colors.primary,
    borderRadius: theme.radii.pill,
    flexDirection: 'row',
    gap: 6,
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.sm,
  },
  actionsRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  actionItem: {
    alignItems: 'center',
    gap: theme.spacing.sm,
  },
  actionOrb: {
    alignItems: 'center',
    backgroundColor: theme.colors.surface,
    borderColor: theme.colors.border,
    borderRadius: theme.radii.pill,
    borderWidth: 1,
    height: 52,
    justifyContent: 'center',
    width: 52,
  },
  metricsRow: {
    gap: theme.spacing.md,
  },
  metricCard: {
    gap: theme.spacing.sm,
    paddingBottom: 18,
  },
  cardPreview: {
    backgroundColor: theme.colors.surface,
    gap: theme.spacing.sm,
    overflow: 'hidden',
  },
  cardHighlight: {
    backgroundColor: '#95BDFF',
    borderRadius: 220,
    height: 170,
    opacity: 0.12,
    position: 'absolute',
    right: -50,
    top: -60,
    width: 170,
  },
  cardChip: {
    backgroundColor: '#DFC27C',
    borderRadius: 8,
    height: 42,
    marginTop: theme.spacing.md,
    width: 56,
  },
  cardNumbers: {
    gap: theme.spacing.sm,
    marginTop: theme.spacing.lg,
  },
  cardFooter: {
    alignItems: 'flex-end',
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginTop: theme.spacing.md,
  },
  networkCard: {
    gap: theme.spacing.sm,
  },
  networkTop: {
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'space-between',
  },
  networkBadge: {
    backgroundColor: theme.colors.surfaceMuted,
    borderColor: theme.colors.border,
    borderRadius: theme.radii.pill,
    borderWidth: 1,
    paddingHorizontal: theme.spacing.md,
    paddingVertical: theme.spacing.sm,
  },
  activityHeader: {
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: theme.spacing.lg,
  },
  activityList: {
    gap: theme.spacing.md,
  },
  activityRow: {
    alignItems: 'center',
    flexDirection: 'row',
    gap: theme.spacing.md,
  },
  activityIcon: {
    alignItems: 'center',
    backgroundColor: theme.colors.primary,
    borderRadius: theme.radii.pill,
    height: 36,
    justifyContent: 'center',
    width: 36,
  },
  activityMeta: {
    flex: 1,
    gap: 2,
  },
  links: {
    gap: theme.spacing.md,
  },
});
