import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import { zhCN } from './locales/zh-CN'
import { enUS } from './locales/en-US'
import { jaJP } from './locales/ja-JP'

// 检测系统语言
const systemLang = navigator.language || 'zh-CN'
const supportedLangs = ['zh-CN', 'en-US', 'ja-JP']
const defaultLang = supportedLangs.includes(systemLang) ? systemLang : 'zh-CN'

i18n
  .use(initReactI18next)
  .init({
    lng: defaultLang,
    fallbackLng: 'zh-CN',
    interpolation: {
      escapeValue: false,
    },
    resources: {
      'zh-CN': zhCN,
      'en-US': enUS,
      'ja-JP': jaJP,
    },
  })

export default i18n
