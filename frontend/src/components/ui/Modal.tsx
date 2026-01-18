import { ReactNode } from 'react'
import { cn } from '../../utils/cn'
import React from "react"

interface ModalProps {
  isOpen: boolean
  onClose: () => void
  title: string|ReactNode
  children: ReactNode
  footer?: ReactNode
  maxWidth?: string
}

// 模态框组件
// 单一职责: 提供统一的模态框容器
export function Modal({ isOpen, onClose, title, children, footer, maxWidth = 'max-w-md' }: ModalProps) {
  if (!isOpen) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* 背景遮罩 */}
      <div
        className="fixed inset-0 bg-black/50"
        onClick={onClose}
      />

      {/* 模态框内容 */}
      <div className={cn("relative z-50 w-full rounded-lg border bg-background p-6 shadow-lg", maxWidth)}>
        <h2 className="text-lg font-semibold">{title}</h2>
        <div className="mt-4">{children}</div>
        {footer && <div className="mt-6 flex justify-end gap-2">{footer}</div>}
      </div>
    </div>
  )
}
