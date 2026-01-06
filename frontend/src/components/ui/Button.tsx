import { ButtonHTMLAttributes, forwardRef } from 'react'
import { cn } from '../../utils/cn'
import React from "react"

// 按钮变体类型
export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'default' | 'outline' | 'destructive' | 'ghost'
  size?: 'default' | 'sm' | 'lg'
}

// 按钮组件
// 单一职责: 提供统一的按钮样式和变体
export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = 'default', size = 'default', children, ...props }, ref) => {
    return (
      <button
        ref={ref}
        className={cn(
          // 基础样式
          'inline-flex items-center justify-center rounded-md text-sm font-medium',
          'transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
          'disabled:pointer-events-none disabled:opacity-50',
          'w-auto', // 明确设置宽度自适应，避免继承拉伸样式

          // 变体样式
          {
            'bg-primary text-primary-foreground hover:bg-primary/90': variant === 'default',
            'border border-input bg-background hover:bg-accent hover:text-accent-foreground':
              variant === 'outline',
            'bg-destructive text-destructive-foreground hover:bg-destructive/90':
              variant === 'destructive',
            'hover:bg-accent hover:text-accent-foreground': variant === 'ghost',
          },

          // 尺寸样式
          {
            'h-10 px-4 py-2': size === 'default',
            'h-9 px-3 text-xs': size === 'sm',
            'h-11 px-8': size === 'lg',
          },

          className
        )}
        {...props}
      >
        {children}
      </button>
    )
  }
)

Button.displayName = 'Button'
