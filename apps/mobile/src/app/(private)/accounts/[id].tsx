import type { ReactElement } from 'react';

import { ComingSoonScreen } from '@/components/layout/ComingSoonScreen';

export default function AccountDetailsScreen(): ReactElement {
  return (
    <ComingSoonScreen
      description="Detalhe da conta reservado para saldo, limites e historico por instituicao."
      title="Detalhe da conta"
    />
  );
}
