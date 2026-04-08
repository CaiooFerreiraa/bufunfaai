import type { ReactElement } from 'react';

import { ComingSoonScreen } from '@/components/layout/ComingSoonScreen';

export default function SettingsScreen(): ReactElement {
  return (
    <ComingSoonScreen
      description="A tela de configuracoes vai concentrar preferencias visuais, cache e suporte offline."
      title="Configuracoes"
    />
  );
}
