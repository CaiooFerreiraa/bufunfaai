import type { UserProfile } from '@bufunfa/shared-types';

import type { MeResponseData } from '@/features/auth/types/auth.types';
import { mapMeToProfile } from '@/features/auth/types/auth.types';
import { apiClient } from '@/services/api/client';
import { endpoints } from '@/services/api/endpoints';
import type { ApiResponse } from '@/types/api';

export async function fetchCurrentUser(): Promise<UserProfile> {
  const response = await apiClient.get<ApiResponse<MeResponseData>>(endpoints.users.me);
  return mapMeToProfile(response.data.data.user);
}
