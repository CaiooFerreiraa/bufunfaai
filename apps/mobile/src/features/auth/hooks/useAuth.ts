import { useClerk } from '@clerk/expo';
import { useRouter } from 'expo-router';

import { useSessionStore } from '@/stores/sessionStore';

interface UseAuthResult {
  readonly logout: () => Promise<void>;
}

export function useAuth(): UseAuthResult {
  const { signOut } = useClerk();
  const router = useRouter();
  const clearSession = useSessionStore((state) => state.clearSession);

  async function logout(): Promise<void> {
    await signOut();
    clearSession();
    router.replace('/(public)/welcome');
  }

  return {
    logout,
  };
}
