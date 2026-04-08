import type { ReactElement } from 'react';

import { ComingSoonScreen } from '@/components/layout/ComingSoonScreen';

export default function InsightsScreen(): ReactElement {
  return (
    <ComingSoonScreen
      description="Os insights serao gerados pela API a partir de dados minimizados e snapshots consolidados."
      title="Insights"
    />
  );
}
