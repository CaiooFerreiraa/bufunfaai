import * as LocalAuthentication from 'expo-local-authentication';

export async function canUseBiometrics(): Promise<boolean> {
  const hasHardware: boolean = await LocalAuthentication.hasHardwareAsync();
  const isEnrolled: boolean = await LocalAuthentication.isEnrolledAsync();

  return hasHardware && isEnrolled;
}

export async function authenticateWithBiometrics(): Promise<boolean> {
  const result = await LocalAuthentication.authenticateAsync({
    promptMessage: 'Desbloquear BufunfaAI',
    disableDeviceFallback: false,
  });

  return result.success;
}
