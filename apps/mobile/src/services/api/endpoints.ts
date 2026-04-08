export const endpoints = {
  auth: {
    login: '/v1/auth/login',
    logout: '/v1/auth/logout',
    logoutAll: '/v1/auth/logout-all',
    refresh: '/v1/auth/refresh',
    register: '/v1/auth/register',
  },
  devices: {
    list: '/v1/devices',
  },
  openFinance: {
    authorizeConsent: (consentId: string): string => `/v1/open-finance/consents/${consentId}/authorize`,
    callback: '/v1/open-finance/callback',
    completeConsent: (consentId: string): string => `/v1/open-finance/consents/${consentId}/complete`,
    connectToken: (consentId: string): string => `/v1/open-finance/consents/${consentId}/connect-token`,
    connections: '/v1/open-finance/connections',
    consentById: (consentId: string): string => `/v1/open-finance/consents/${consentId}`,
    consents: '/v1/open-finance/consents',
    institutions: '/v1/open-finance/institutions',
    revokeConsent: (consentId: string): string => `/v1/open-finance/consents/${consentId}/revoke`,
    syncConnection: (connectionId: string): string => `/v1/open-finance/connections/${connectionId}/sync`,
    syncStatus: (connectionId: string): string => `/v1/open-finance/connections/${connectionId}/sync-status`,
  },
  users: {
    me: '/v1/users/me',
  },
} as const;
