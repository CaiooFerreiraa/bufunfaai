import { useEffect } from 'react';

import { registerApiInterceptors } from '@/services/api/interceptors';

let initialized: boolean = false;

export function useInitializeApi(): void {
  useEffect((): void => {
    if (initialized) {
      return;
    }

    registerApiInterceptors();
    initialized = true;
  }, []);
}
