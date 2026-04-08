import { useEffect } from 'react';

import { fetchCurrentUser, refreshSessionRequest } from '@/features/auth/services/authService';
import { destroySession, persistSession } from '@/lib/sessionManager';
import { useSessionStore } from '@/stores/sessionStore';

export function useSessionBootstrap(enabled: boolean): void {
  const tokens = useSessionStore((state) => state.tokens);
  const user = useSessionStore((state) => state.user);
  const hydrateSession = useSessionStore((state) => state.hydrateSession);
  const metadata = useSessionStore((state) => state.metadata);

  useEffect((): void => {
    async function bootstrapAuthenticatedSession(): Promise<void> {
      if (!enabled || !tokens || !metadata) {
        return;
      }

      try {
        if (!user) {
          const refreshed = await refreshSessionRequest();
          await persistSession(refreshed.tokens, refreshed.user);
          return;
        }

        const currentUser = await fetchCurrentUser();
        hydrateSession(tokens, metadata, currentUser);
      } catch {
        await destroySession();
      }
    }

    void bootstrapAuthenticatedSession();
  }, [enabled, hydrateSession, metadata, tokens, user]);
}
