import { create } from 'zustand';

interface SessionStoreState {
  readonly isBootstrapping: boolean;
  readonly requiresBiometricUnlock: boolean;
  readonly biometricEnabled: boolean;
  readonly setBootstrapState: (value: boolean) => void;
  readonly hydrateSession: (biometricEnabled: boolean) => void;
  readonly setBiometricEnabled: (value: boolean) => void;
  readonly clearSession: () => void;
  readonly setRequiresBiometricUnlock: (value: boolean) => void;
}

export const useSessionStore = create<SessionStoreState>((set) => ({
  isBootstrapping: true,
  requiresBiometricUnlock: false,
  biometricEnabled: false,
  setBootstrapState: (value: boolean): void => set({ isBootstrapping: value }),
  hydrateSession: (biometricEnabled: boolean): void =>
    set({
      biometricEnabled,
      requiresBiometricUnlock: biometricEnabled,
    }),
  setBiometricEnabled: (value: boolean): void =>
    set({
      biometricEnabled: value,
    }),
  clearSession: (): void =>
    set({
      biometricEnabled: false,
      requiresBiometricUnlock: false,
    }),
  setRequiresBiometricUnlock: (value: boolean): void => set({ requiresBiometricUnlock: value }),
}));
