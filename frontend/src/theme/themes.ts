// 主题配置
export const themes = {
  light: {
    background: 'ffffff',
    foreground: '09090b',
  },
  dark: {
    background: '09090b',
    foreground: 'fafafa',
  },
} as const

export type Theme = keyof typeof themes
