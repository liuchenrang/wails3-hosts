// 应用相关类型定义

export type Theme = 'light' | 'dark'

export type Language = 'zh-CN' | 'en-US' | 'ja-JP'

export interface AppSettings {
  theme: Theme
  language: Language
  autoApply: boolean
}

export interface ConflictInfo {
  [hostname: string]: string[]
}
