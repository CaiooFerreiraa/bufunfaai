import type { ReactElement } from 'react';

import { EmptyState } from '@/components/feedback/EmptyState';

interface ErrorStateProps {
  readonly onRetry?: () => void;
}

export function ErrorState(props: ErrorStateProps): ReactElement {
  const { onRetry } = props;

  return (
    <EmptyState
      actionLabel={onRetry ? 'Tentar novamente' : undefined}
      description="Não foi possível carregar os dados agora."
      onActionPress={onRetry}
      title="Algo deu errado"
    />
  );
}
