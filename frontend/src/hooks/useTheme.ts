import { useEffect, useState } from 'react'
import type { Theme } from '../types/app'

// 主题 Hook
// 单一职责: 管理主题状态和切换逻辑
export function useTheme() {
  const [theme, setTheme] = useState<Theme>(() => {
    // 从本地存储读取主题偏好
    const savedTheme = localStorage.getItem('theme') as Theme
    if (savedTheme) {
      return savedTheme
    }

    // 跟随系统主题
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
      return 'dark'
    }

    return 'light'
  })

  useEffect(() => {
    // 应用主题到 DOM
    const root = window.document.documentElement
    root.classList.remove('light', 'dark')
    root.classList.add(theme)

    // 保存到本地存储
    localStorage.setItem('theme', theme)
  }, [theme])

  const toggleTheme = () => {
    setTheme(prev => (prev === 'light' ? 'dark' : 'light'))
  }

  return {
    theme,
    toggleTheme,
    setTheme,
  }
}
