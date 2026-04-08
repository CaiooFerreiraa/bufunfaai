import * as SecureStore from 'expo-secure-store';

const SESSION_KEY: string = 'bufunfaai.session';

export async function getSession(): Promise<string | null> {
  return SecureStore.getItemAsync(SESSION_KEY);
}

export async function saveSession(value: string): Promise<void> {
  await SecureStore.setItemAsync(SESSION_KEY, value);
}

export async function clearSession(): Promise<void> {
  await SecureStore.deleteItemAsync(SESSION_KEY);
}
