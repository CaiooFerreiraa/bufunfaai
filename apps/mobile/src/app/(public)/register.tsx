import { useSignUp } from '@clerk/expo';
import type { ReactElement } from 'react';
import { useState } from 'react';
import { StyleSheet, View } from 'react-native';

import { FeatureScreen } from '@/components/layout/FeatureScreen';
import { AppText } from '@/components/ui/AppText';
import { Button } from '@/components/ui/Button';
import { TextField } from '@/components/ui/TextField';
import { usePostAuthFlow } from '@/features/auth/hooks/usePostAuthFlow';
import { registerSchema } from '@/features/auth/schemas/authSchemas';
import { getClerkErrorMessage } from '@/features/auth/utils/clerkErrors';
import { theme } from '@/theme/tokens';

function splitFullName(fullName: string): { firstName?: string; lastName?: string } {
  const chunks = fullName.trim().split(/\s+/).filter(Boolean);
  if (chunks.length === 0) {
    return {};
  }

  return {
    firstName: chunks[0],
    lastName: chunks.length > 1 ? chunks.slice(1).join(' ') : undefined,
  };
}

export default function RegisterScreen(): ReactElement {
  const { signUp, fetchStatus } = useSignUp();
  const { completeAuthentication } = usePostAuthFlow();
  const [fullName, setFullName] = useState<string>('');
  const [email, setEmail] = useState<string>('');
  const [phone, setPhone] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [code, setCode] = useState<string>('');
  const [errorMessage, setErrorMessage] = useState<string>('');

  async function handleRegister(): Promise<void> {
    const parsed = registerSchema.safeParse({ fullName, email, password, phone });
    if (!parsed.success) {
      setErrorMessage('Revise seus dados antes de continuar.');
      return;
    }

    const names = splitFullName(parsed.data.fullName);
    const result = await signUp.password({
      emailAddress: parsed.data.email,
      password: parsed.data.password,
      firstName: names.firstName,
      lastName: names.lastName,
    });

    if (result.error) {
      setErrorMessage(
        getClerkErrorMessage(result.error, ['emailAddress', 'password', 'firstName', 'lastName'], 'Não foi possível criar sua conta.'),
      );
      return;
    }

    const verificationResult = await signUp.verifications.sendEmailCode();
    if (verificationResult.error) {
      setErrorMessage(getClerkErrorMessage(verificationResult.error, ['emailAddress'], 'Não foi possível enviar o código.'));
    }
  }

  async function handleVerifyEmail(): Promise<void> {
    const result = await signUp.verifications.verifyEmailCode({ code });
    if (result.error) {
      setErrorMessage(getClerkErrorMessage(result.error, ['code'], 'Não foi possível confirmar seu e-mail.'));
      return;
    }

    if (signUp.status !== 'complete') {
      setErrorMessage('Confira o código e tente novamente.');
      return;
    }

    await signUp.finalize({
      navigate: async (): Promise<void> => {
        await completeAuthentication();
      },
    });
  }

  const isVerificationStep =
    signUp.status === 'missing_requirements' &&
    signUp.unverifiedFields.includes('email_address') &&
    signUp.missingFields.length === 0;

  return (
    <FeatureScreen
      description={
        isVerificationStep
          ? 'Enviamos um código para o seu e-mail.'
          : 'Crie sua conta para acompanhar tudo em um só lugar.'
      }
      title={isVerificationStep ? 'Confirme seu e-mail' : 'Criar conta'}
    >
      <View style={styles.form}>
        {isVerificationStep ? (
          <TextField keyboardType="number-pad" label="Código" onChangeText={setCode} value={code} />
        ) : (
          <>
            <TextField label="Nome completo" onChangeText={setFullName} value={fullName} />
            <TextField autoCapitalize="none" keyboardType="email-address" label="Email" onChangeText={setEmail} value={email} />
            <TextField keyboardType="phone-pad" label="Telefone" onChangeText={setPhone} value={phone} />
            <TextField label="Senha" onChangeText={setPassword} secureTextEntry value={password} />
          </>
        )}

        {errorMessage ? <AppText color={theme.colors.error}>{errorMessage}</AppText> : null}

        <Button
          disabled={fetchStatus === 'fetching'}
          label={
            fetchStatus === 'fetching'
              ? 'Aguarde...'
              : isVerificationStep
                ? 'Confirmar e entrar'
                : 'Criar conta'
          }
          onPress={(): void => void (isVerificationStep ? handleVerifyEmail() : handleRegister())}
        />

        {isVerificationStep ? (
          <Button
            disabled={fetchStatus === 'fetching'}
            label="Enviar novo código"
            onPress={(): void => {
              void signUp.verifications.sendEmailCode();
            }}
            variant="secondary"
          />
        ) : null}
      </View>
    </FeatureScreen>
  );
}

const styles = StyleSheet.create({
  form: {
    gap: theme.spacing.md,
  },
});
