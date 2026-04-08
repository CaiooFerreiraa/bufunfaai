import { create } from 'zustand';

interface AppStoreState {
  readonly isOnline: boolean;
  readonly lastSyncAt: string | null;
  readonly setOnline: (value: boolean) => void;
  readonly setLastSyncAt: (value: string | null) => void;
}

export const useAppStore = create<AppStoreState>((set) => ({
  isOnline: true,
  lastSyncAt: null,
  setOnline: (value: boolean): void => set({ isOnline: value }),
  setLastSyncAt: (value: string | null): void => set({ lastSyncAt: value }),
}));
