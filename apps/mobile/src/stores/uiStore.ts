import { create } from 'zustand';

interface UIStoreState {
  readonly transactionSearch: string;
  readonly setTransactionSearch: (value: string) => void;
}

export const useUIStore = create<UIStoreState>((set) => ({
  transactionSearch: '',
  setTransactionSearch: (value: string): void => set({ transactionSearch: value }),
}));
