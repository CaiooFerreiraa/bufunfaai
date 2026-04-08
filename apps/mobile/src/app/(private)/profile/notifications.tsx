import type { ReactElement } from 'react';

import { FeatureScreen } from '@/components/layout/FeatureScreen';
import { AppText } from '@/components/ui/AppText';
import { Card } from '@/components/ui/Card';
import { theme } from '@/theme/tokens';

export default function ProfileNotificationsScreen(): ReactElement {
  return (
    <FeatureScreen
      description="O push ficou desativado neste ambiente local até existir um projeto EAS com credenciais válidas."
      title="Notificações"
    >
      <Card>
        <AppText color={theme.colors.accent} variant="label">
          Estado atual
        </AppText>
        <AppText variant="headline">Push desativado no Expo Go local</AppText>
        <AppText color={theme.colors.textSecondary}>
          Para Android e iOS, o registro de push será reativado quando o app tiver `projectId` EAS e build própria de desenvolvimento.
        </AppText>
      </Card>
    </FeatureScreen>
  );
}
