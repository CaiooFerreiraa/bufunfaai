import { useEffect, useRef } from 'react';
import { AppState, type AppStateStatus } from 'react-native';

import { useSessionStore } from '@/stores/sessionStore';

export function useBiometricAppLock(enabled: boolean): void {
  const appState = useRef<AppStateStatus>(AppState.currentState);
  const isAuthenticated = useSessionStore((state) => state.isAuthenticated);
  const metadata = useSessionStore((state) => state.metadata);
  const requiresBiometricUnlock = useSessionStore((state) => state.requiresBiometricUnlock);
  const setRequiresBiometricUnlock = useSessionStore((state) => state.setRequiresBiometricUnlock);

  useEffect((): (() => void) | void => {
    if (!enabled) {
      return;
    }

    const subscription = AppState.addEventListener('change', (nextState: AppStateStatus): void => {
      const previousState: AppStateStatus = appState.current;

      if (
        metadata?.biometricEnabled &&
        isAuthenticated &&
        !requiresBiometricUnlock &&
        (previousState === 'background' || previousState === 'inactive') &&
        nextState === 'active'
      ) {
        setRequiresBiometricUnlock(true);
      }

      appState.current = nextState;
    });

    return (): void => {
      subscription.remove();
    };
  }, [enabled, isAuthenticated, metadata?.biometricEnabled, requiresBiometricUnlock, setRequiresBiometricUnlock]);
}
