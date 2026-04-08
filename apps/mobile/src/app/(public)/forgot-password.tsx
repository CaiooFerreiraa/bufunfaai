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
import { getClerkErrorMessage } from '@/features/auth/utils/clerkErrors';
import { theme } from '@/theme/tokens';

export default function ForgotPasswordScreen(): ReactElement {
  const { signIn, fetchStatus } = useSignIn();
  const { completeAuthentication } = usePostAuthFlow();
  const [emailAddress, setEmailAddress] = useState<string>('');
  const [code, setCode] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [codeSent, setCodeSent] = useState<boolean>(false);
  const [errorMessage, setErrorMessage] = useState<string>('');

  async function sendCode(): Promise<void> {
    const createResult = await signIn.create({ identifier: emailAddress });
    if (createResult.error) {
      setErrorMessage(getClerkErrorMessage(createResult.error, ['identifier'], 'Não foi possível encontrar essa conta.'));
      return;
    }

    const sendCodeResult = await signIn.resetPasswordEmailCode.sendCode();
    if (sendCodeResult.error) {
      setErrorMessage(getClerkErrorMessage(sendCodeResult.error, ['identifier'], 'Não foi possível enviar o código.'));
      return;
    }

    setErrorMessage('');
    setCodeSent(true);
  }

  async function verifyCode(): Promise<void> {
    const result = await signIn.resetPasswordEmailCode.verifyCode({ code });
    if (result.error) {
      setErrorMessage(getClerkErrorMessage(result.error, ['code'], 'Confira o código e tente novamente.'));
      return;
    }

    setErrorMessage('');
  }

  async function submitNewPassword(): Promise<void> {
    const result = await signIn.resetPasswordEmailCode.submitPassword({ password });
    if (result.error) {
      setErrorMessage(getClerkErrorMessage(result.error, ['password'], 'Não foi possível salvar sua nova senha.'));
      return;
    }

    if (signIn.status !== 'complete') {
      setErrorMessage('Confira seus dados e tente novamente.');
      return;
    }

    await signIn.finalize({
      navigate: async (): Promise<void> => {
        await completeAuthentication();
      },
    });
  }

  const isNewPasswordStep = signIn.status === 'needs_new_password';
  const isCodeStep = codeSent && !isNewPasswordStep;

  return (
    <FeatureScreen
      description={
        isNewPasswordStep
          ? 'Escolha sua nova senha.'
          : isCodeStep
            ? 'Digite o código enviado para o seu e-mail.'
            : 'Informe o e-mail da sua conta para continuar.'
      }
      title="Recuperar acesso"
    >
      <View style={styles.form}>
        {!isCodeStep && !isNewPasswordStep ? (
          <TextField
            autoCapitalize="none"
            keyboardType="email-address"
            label="Email"
            onChangeText={setEmailAddress}
            value={emailAddress}
          />
        ) : null}

        {isCodeStep ? <TextField keyboardType="number-pad" label="Código" onChangeText={setCode} value={code} /> : null}
        {isNewPasswordStep ? <TextField label="Nova senha" onChangeText={setPassword} secureTextEntry value={password} /> : null}

        {errorMessage ? <AppText color={theme.colors.error}>{errorMessage}</AppText> : null}

        <Button
          disabled={fetchStatus === 'fetching' || (!emailAddress && !isCodeStep && !isNewPasswordStep)}
          label={
            fetchStatus === 'fetching'
              ? 'Aguarde...'
              : isNewPasswordStep
                ? 'Salvar nova senha'
                : isCodeStep
                  ? 'Confirmar código'
                  : 'Enviar código'
          }
          onPress={(): void => void (isNewPasswordStep ? submitNewPassword() : isCodeStep ? verifyCode() : sendCode())}
        />

        <Link href="/(public)/login" style={styles.link}>
          <AppText color={theme.colors.textSecondary}>Voltar para entrar</AppText>
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
