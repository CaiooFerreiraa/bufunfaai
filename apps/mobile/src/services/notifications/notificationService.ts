import Constants from 'expo-constants';
import * as Device from 'expo-device';
import * as Notifications from 'expo-notifications';
import { Platform } from 'react-native';

Notifications.setNotificationHandler({
  handleNotification: async () => ({
    shouldPlaySound: true,
    shouldSetBadge: true,
    shouldShowBanner: true,
    shouldShowList: true,
  }),
});

function resolveProjectId(): string | null {
  const configProjectId =
    Constants.expoConfig?.extra?.eas?.projectId ?? Constants.easConfig?.projectId;

  return process.env.EXPO_PUBLIC_EAS_PROJECT_ID ?? configProjectId ?? null;
}

export function canRegisterNotifications(): boolean {
  return resolveProjectId() !== null;
}

export async function registerNotifications(): Promise<string | null> {
  if (!Device.isDevice) {
    return null;
  }

  if (Platform.OS === 'android') {
    await Notifications.setNotificationChannelAsync('default', {
      name: 'default',
      importance: Notifications.AndroidImportance.MAX,
      vibrationPattern: [0, 250, 250, 250],
      lightColor: '#C67622',
    });
  }

  const permissions = await Notifications.getPermissionsAsync();
  if (!permissions.granted) {
    const requestedPermissions = await Notifications.requestPermissionsAsync();
    if (!requestedPermissions.granted) {
      return null;
    }
  }

  const projectId = resolveProjectId();
  if (!projectId) {
    return null;
  }

  const pushToken = await Notifications.getExpoPushTokenAsync({ projectId });
  return pushToken.data;
}
