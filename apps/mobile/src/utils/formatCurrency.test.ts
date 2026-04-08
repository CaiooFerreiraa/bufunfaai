import { describe, expect, it } from 'vitest';

import { formatCurrency } from '@/utils/formatCurrency';

describe('formatCurrency', (): void => {
  it('formats BRL currency values', (): void => {
    expect(formatCurrency(12.5)).toBe('R$\u00a012,50');
  });
});
