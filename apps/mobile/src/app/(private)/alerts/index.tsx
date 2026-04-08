import type { ReactElement } from 'react';

import { ComingSoonScreen } from '@/components/layout/ComingSoonScreen';

export default function AlertsScreen(): ReactElement {
  return (
    <ComingSoonScreen
      description="Alertas de gasto, consentimento e sincronizacao entram quando o backend publicar regras e eventos."
      title="Alertas"
    />
  );
}
