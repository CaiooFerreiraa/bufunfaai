import type { ReactElement } from 'react';

import { ComingSoonScreen } from '@/components/layout/ComingSoonScreen';

export default function AccountsScreen(): ReactElement {
  return (
    <ComingSoonScreen
      description="A listagem de contas vai consumir o read model consolidado da API assim que os endpoints forem expostos."
      title="Contas"
    />
  );
}
