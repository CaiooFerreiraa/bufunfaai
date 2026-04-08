import type { ReactElement } from 'react';

import { ComingSoonScreen } from '@/components/layout/ComingSoonScreen';

export default function GoalsScreen(): ReactElement {
  return (
    <ComingSoonScreen
      description="As metas financeiras vao usar snapshots locais e regras de sincronizacao simples."
      title="Metas"
    />
  );
}
