import type { ReactElement } from 'react';
import { useState } from 'react';
import { StyleSheet, View } from 'react-native';

import { FeatureScreen } from '@/components/layout/FeatureScreen';
import { AppText } from '@/components/ui/AppText';
import { Button } from '@/components/ui/Button';
import { TextField } from '@/components/ui/TextField';
import { useAuth } from '@/features/auth/hooks/useAuth';
import { registerSchema } from '@/features/auth/schemas/authSchemas';
import { theme } from '@/theme/tokens';

export default function RegisterScreen(): ReactElement {
  const { register } = useAuth();
  const [fullName, setFullName] = useState<string>('');
  const [email, setEmail] = useState<string>('');
  const [phone, setPhone] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [errorMessage, setErrorMessage] = useState<string>('');

  async function handleRegister(): Promise<void> {
    const parsed = registerSchema.safeParse({ fullName, email, password, phone });
    if (!parsed.success) {
      setErrorMessage('Revise seus dados antes de continuar.');
      return;
    }

    setErrorMessage('');
    try {
      await register(parsed.data);
    } catch {
      setErrorMessage('Não foi possível criar sua conta.');
    }
  }

  return (
    <FeatureScreen description="Crie sua conta para começar a conectar instituições financeiras." title="Cadastro">
      <View style={styles.form}>
        <TextField label="Nome completo" onChangeText={setFullName} value={fullName} />
        <TextField autoCapitalize="none" keyboardType="email-address" label="Email" onChangeText={setEmail} value={email} />
        <TextField keyboardType="phone-pad" label="Telefone" onChangeText={setPhone} value={phone} />
        <TextField label="Senha" onChangeText={setPassword} secureTextEntry value={password} />
        {errorMessage ? <AppText color={theme.colors.error}>{errorMessage}</AppText> : null}
        <Button label="Criar conta" onPress={(): void => void handleRegister()} />
      </View>
    </FeatureScreen>
  );
}

const styles = StyleSheet.create({
  form: {
    gap: theme.spacing.md,
  },
});
