import { Platform } from 'react-native';

const arialFamily: string = Platform.select({
  android: 'sans-serif',
  default: 'Arial',
  ios: 'Arial',
  web: 'Arial, sans-serif',
}) ?? 'Arial';

export const theme = {
  colors: {
    primary: '#C7FF45',
    primaryMuted: '#7AA62C',
    accent: '#8BF03C',
    accentSoft: '#E7FF9B',
    background: '#070A12',
    backgroundStrong: '#03050B',
    surface: '#0F1420',
    surfaceMuted: '#151B2A',
    surfaceInverse: '#1A2233',
    border: '#242C40',
    borderStrong: '#46506A',
    textPrimary: '#F4F7FF',
    textSecondary: '#8D98AF',
    textInverse: '#081018',
    success: '#53E38A',
    warning: '#FFB648',
    error: '#FF6B78',
  },
  fonts: {
    body: arialFamily,
    display: arialFamily,
    label: arialFamily,
    mono: 'monospace',
  },
  spacing: {
    xs: 4,
    sm: 8,
    md: 12,
    lg: 16,
    xl: 24,
    '2xl': 32,
  },
  radii: {
    md: 12,
    lg: 18,
    xl: 28,
    pill: 999,
  },
} as const;
