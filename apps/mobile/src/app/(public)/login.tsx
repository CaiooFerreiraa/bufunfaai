import { useSignIn } from '@clerk/expo';
import { Link } from 'expo-router';
import type { ReactElement } from 'react';
import { useState } from 'react';
import { StyleSheet, View } from 'react-native';

import { FeatureScreen } from '@/components/layout/FeatureScreen';
import { AppText } from '@/components/ui/AppText';
import { Button } from '@/components/ui/Button';
import { TextField } from '@/components/ui/TextField';
import { usePostAuthFlow } from '@/features/auth/hooks/usePostAuthFlow';
import { loginSchema } from '@/features/auth/schemas/authSchemas';
import { getClerkErrorMessage } from '@/features/auth/utils/clerkErrors';
import { theme } from '@/theme/tokens';

export default function LoginScreen(): ReactElement {
  const { signIn, fetchStatus } = useSignIn();
  const { completeAuthentication } = usePostAuthFlow();
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [errorMessage, setErrorMessage] = useState<string>('');

  async function handleLogin(): Promise<void> {
    const parsed = loginSchema.safeParse({ email, password });
    if (!parsed.success) {
      setErrorMessage('Preencha email e senha corretamente.');
      return;
    }

    const result = await signIn.password({
      identifier: parsed.data.email,
      password: parsed.data.password,
    });

    if (result.error) {
      setErrorMessage(getClerkErrorMessage(result.error, ['identifier', 'password'], 'Não foi possível entrar agora.'));
      return;
    }

    if (signIn.status !== 'complete') {
      setErrorMessage('Confirme sua conta antes de continuar.');
      return;
    }

    await signIn.finalize({
      navigate: async (): Promise<void> => {
        await completeAuthentication();
      },
    });
  }

  return (
    <FeatureScreen description="Entre para acompanhar seu dinheiro em um só lugar." title="Acesse sua conta">
      <View style={styles.form}>
        <TextField autoCapitalize="none" keyboardType="email-address" label="Email" onChangeText={setEmail} value={email} />
        <TextField label="Senha" onChangeText={setPassword} secureTextEntry value={password} />
        {errorMessage ? <AppText color={theme.colors.error}>{errorMessage}</AppText> : null}
        <Button
          disabled={fetchStatus === 'fetching'}
          label={fetchStatus === 'fetching' ? 'Entrando...' : 'Continuar'}
          onPress={(): void => void handleLogin()}
        />
        <Link href="/(public)/forgot-password" style={styles.link}>
          <AppText color={theme.colors.textSecondary}>Esqueci minha senha</AppText>
        </Link>
      </View>
    </FeatureScreen>
  );
}

const styles = StyleSheet.create({
  form: {
    gap: theme.spacing.md,
  },
  link: {
    alignSelf: 'center',
    marginTop: theme.spacing.xs,
  },
});
