import * as SecureStore from 'expo-secure-store';

const BIOMETRIC_ENABLED_KEY: string = 'bufunfaai.biometric_enabled';
const BIOMETRIC_PROMPTED_KEY: string = 'bufunfaai.biometric_prompted';

export async function setBiometricPreference(enabled: boolean): Promise<void> {
  await SecureStore.setItemAsync(BIOMETRIC_ENABLED_KEY, JSON.stringify(enabled));
}

export async function getBiometricPreference(): Promise<boolean> {
  const rawValue: string | null = await SecureStore.getItemAsync(BIOMETRIC_ENABLED_KEY);
  return rawValue === 'true';
}

export async function setBiometricPrompted(value: boolean): Promise<void> {
  await SecureStore.setItemAsync(BIOMETRIC_PROMPTED_KEY, JSON.stringify(value));
}

export async function getBiometricPrompted(): Promise<boolean> {
  const rawValue: string | null = await SecureStore.getItemAsync(BIOMETRIC_PROMPTED_KEY);
  return rawValue === 'true';
}
