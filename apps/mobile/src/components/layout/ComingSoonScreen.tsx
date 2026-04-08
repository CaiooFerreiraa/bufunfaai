import type { ReactElement } from 'react';

import { EmptyState } from '@/components/feedback/EmptyState';
import { FeatureScreen } from '@/components/layout/FeatureScreen';

interface ComingSoonScreenProps {
  readonly description: string;
  readonly title: string;
}

export function ComingSoonScreen(props: ComingSoonScreenProps): ReactElement {
  const { description, title } = props;

  return (
    <FeatureScreen description={description} title={title}>
      <EmptyState
        description="A base da tela já está montada, mas os dados dependem da próxima fase do backend."
        title="Em construção"
      />
    </FeatureScreen>
  );
}
