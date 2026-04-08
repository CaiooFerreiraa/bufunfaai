import NetInfo from '@react-native-community/netinfo';

export function subscribeToConnectivity(
  onChange: (isOnline: boolean) => void,
): () => void {
  return NetInfo.addEventListener((state): void => {
    onChange(Boolean(state.isConnected));
  });
}
