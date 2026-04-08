import type { UserProfile } from '@bufunfa/shared-types';
import type { UseQueryResult } from '@tanstack/react-query';
import { useQuery } from '@tanstack/react-query';

import { fetchCurrentUser } from '@/features/auth/services/authService';

export function useCurrentUserQuery(): UseQueryResult<UserProfile, Error> {
  return useQuery({
    queryKey: ['me'],
    queryFn: fetchCurrentUser,
    staleTime: 5 * 60 * 1000,
  });
}
