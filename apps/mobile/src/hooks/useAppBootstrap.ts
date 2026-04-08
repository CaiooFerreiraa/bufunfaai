import { useEffect } from 'react';

import { subscribeToConnectivity } from '@/services/network/connectivityService';
import { getBiometricPreference } from '@/services/storage/preferences';
import { initializeSQLite } from '@/services/storage/sqlite';
import { useAppStore } from '@/stores/appStore';
import { useSessionStore } from '@/stores/sessionStore';

export function useAppBootstrap(): boolean {
  const setBootstrapState = useSessionStore((state) => state.setBootstrapState);
  const isBootstrapping = useSessionStore((state) => state.isBootstrapping);
  const setOnline = useAppStore((state) => state.setOnline);

  useEffect(() => {
    let mounted = true;

    async function bootstrap(): Promise<void> {
      await initializeSQLite();
      const biometricEnabled: boolean = await getBiometricPreference();

      if (mounted) {
        useSessionStore.getState().hydrateSession(biometricEnabled);
        setBootstrapState(false);
      }
    }

    const unsubscribe = subscribeToConnectivity((isOnline: boolean): void => {
      setOnline(isOnline);
    });

    void bootstrap();

    return (): void => {
      mounted = false;
      unsubscribe();
    };
  }, [setBootstrapState, setOnline]);

  return !isBootstrapping;
}
