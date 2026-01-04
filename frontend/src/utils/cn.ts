import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

// 合并 Tailwind 类名的工具函数
// DRY: 统一的类名合并逻辑
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
