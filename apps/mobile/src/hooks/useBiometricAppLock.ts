import { useAuth } from '@clerk/expo';
import { useEffect, useRef } from 'react';
import { AppState, type AppStateStatus } from 'react-native';

import { useSessionStore } from '@/stores/sessionStore';

export function useBiometricAppLock(enabled: boolean): void {
  const appState = useRef<AppStateStatus>(AppState.currentState);
  const { isSignedIn } = useAuth();
  const biometricEnabled = useSessionStore((state) => state.biometricEnabled);
  const requiresBiometricUnlock = useSessionStore((state) => state.requiresBiometricUnlock);
  const setRequiresBiometricUnlock = useSessionStore((state) => state.setRequiresBiometricUnlock);

  useEffect((): (() => void) | void => {
    if (!enabled) {
      return;
    }

    const subscription = AppState.addEventListener('change', (nextState: AppStateStatus): void => {
      const previousState: AppStateStatus = appState.current;

      if (
        biometricEnabled &&
        isSignedIn &&
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
  }, [biometricEnabled, enabled, isSignedIn, requiresBiometricUnlock, setRequiresBiometricUnlock]);
}
