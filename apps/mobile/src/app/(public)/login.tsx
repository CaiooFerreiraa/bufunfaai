import { Link } from 'expo-router';
import { isAxiosError } from 'axios';
import type { ReactElement } from 'react';
import { useState } from 'react';
import { StyleSheet, View } from 'react-native';

import { FeatureScreen } from '@/components/layout/FeatureScreen';
import { AppText } from '@/components/ui/AppText';
import { Button } from '@/components/ui/Button';
import { TextField } from '@/components/ui/TextField';
import { useAuth } from '@/features/auth/hooks/useAuth';
import { loginSchema } from '@/features/auth/schemas/authSchemas';
import { theme } from '@/theme/tokens';

export default function LoginScreen(): ReactElement {
  const { login } = useAuth();
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [errorMessage, setErrorMessage] = useState<string>('');

  async function handleLogin(): Promise<void> {
    const parsed = loginSchema.safeParse({ email, password });
    if (!parsed.success) {
      setErrorMessage('Preencha email e senha corretamente.');
      return;
    }

    setErrorMessage('');
    try {
      await login(parsed.data);
    } catch (error) {
      if (isAxiosError<{ error?: { message?: string } }>(error)) {
        const message = error.response?.data?.error?.message;
        if (message) {
          setErrorMessage(message);
          return;
        }
      }

      setErrorMessage('Não foi possível entrar agora.');
    }
  }

  return (
    <FeatureScreen description="Entre para retomar seu panorama financeiro, conexões e trilha de segurança." title="Acesse sua mesa financeira">
      <View style={styles.form}>
        <View style={styles.callout}>
          <AppText color={theme.colors.accent} variant="label">
            Sessao protegida
          </AppText>
          <AppText color={theme.colors.textSecondary}>
            Seu refresh token fica em armazenamento seguro no dispositivo e a biometria protege o desbloqueio local.
          </AppText>
        </View>
        <TextField autoCapitalize="none" keyboardType="email-address" label="Email" onChangeText={setEmail} value={email} />
        <TextField label="Senha" onChangeText={setPassword} secureTextEntry value={password} />
        {errorMessage ? <AppText color={theme.colors.error}>{errorMessage}</AppText> : null}
        <Button label="Continuar" onPress={(): void => void handleLogin()} />
        <Link href="/(public)/forgot-password" style={styles.link}>
          <AppText color={theme.colors.primary}>Esqueci minha senha</AppText>
        </Link>
      </View>
    </FeatureScreen>
  );
}

const styles = StyleSheet.create({
  callout: {
    backgroundColor: theme.colors.surfaceMuted,
    borderColor: theme.colors.border,
    borderRadius: theme.radii.lg,
    borderWidth: 1,
    gap: theme.spacing.sm,
    padding: theme.spacing.lg,
  },
  form: {
    gap: theme.spacing.md,
  },
  link: {
    paddingVertical: theme.spacing.sm,
  },
});
