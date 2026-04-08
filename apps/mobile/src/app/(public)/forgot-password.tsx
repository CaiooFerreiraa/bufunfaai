import type { ReactElement } from 'react';

import { ComingSoonScreen } from '@/components/layout/ComingSoonScreen';

export default function ForgotPasswordScreen(): ReactElement {
  return (
    <ComingSoonScreen
      description="A tela de recuperacao de acesso ja esta reservada para a proxima fase do backend."
      title="Recuperar acesso"
    />
  );
}
