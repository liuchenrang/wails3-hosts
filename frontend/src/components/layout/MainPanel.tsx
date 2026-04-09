import { useState, useMemo, useEffect, useRef } from 'react'
import { useTranslation } from 'react-i18next'
import { Save, Eye, AlertCircle } from 'lucide-react'
import { HostsGroup } from '../../types/hosts'
import { Button } from '../ui/Button'
import { Textarea } from '../ui/Textarea'
import React from "react"
import {t} from "i18next";

interface MainPanelProps {
  group: HostsGroup | null
  onUpdateEntry: (entryId: string, ip: string, hostname: string, comment: string) => void
  onAddEntry: (ip: string, hostname: string, comment: string) => void
  onDeleteEntry: (entryId: string) => void
  onApply: () => void
  onPreview: () => void
  onBatchUpdate: (entries: Array<{ ip: string; hostname: string; comment: string; enabled: boolean }>, silent?: boolean) => Promise<void>
}

// 验证 IP 地址格式
const isValidIP = (ip: string): boolean => {
  const ipPattern = /^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/
  return ipPattern.test(ip)
}

// 验证主机名格式
const isValidHostname = (hostname: string): boolean => {
  if (!hostname || hostname.length > 253) return false
  const hostnamePattern = /^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$/
  return hostnamePattern.test(hostname)
}

// 验证 hosts 条目行格式
const validateHostsLine = (line: string): { valid: boolean; error?: string } => {
  const trimmed = line.trim()
  const enabled = !trimmed.startsWith('#')
  const cleanLine = enabled ? trimmed : trimmed.slice(1).trim()

  // 分离注释
  const commentIndex = cleanLine.indexOf('#')
  let mainContent = cleanLine

  if (commentIndex !== -1) {
    mainContent = cleanLine.substring(0, commentIndex).trim()
  }

  // 解析 IP 和 hostname
  const parts = mainContent.split(/\s+/).filter(p => p)

  if (parts.length < 2) {
    return { valid: false, error: t('mainPanel.validation.invalidFormat') }
  }

  const [ip, ...hostnameParts] = parts
  const hostname = hostnameParts.join(' ')

  if (!isValidIP(ip)) {
    return { valid: false, error: t('mainPanel.validation.invalidIP', { ip }) }
  }

  if (!isValidHostname(hostname)) {
    return { valid: false, error: t('mainPanel.validation.invalidHostname', { hostname }) }
  }

  return { valid: true }
}

// 主面板组件
// 单一职责: 显示选中分组的 hosts 条目，支持编辑和管理
export function MainPanel({
  group,
  onUpdateEntry,
  onAddEntry,
  onDeleteEntry,
  onApply,
  onPreview,
  onBatchUpdate,
}: MainPanelProps) {
  const { t } = useTranslation()
  const [memoContent, setMemoContent] = useState('')
  const [validationErrors, setValidationErrors] = useState<string[]>([])
  const [showErrors, setShowErrors] = useState(false)
  const textareaRef = useRef<HTMLTextAreaElement>(null)

  // 将 entries 转换为 hosts 文件格式
  const entriesToHostsFormat = useMemo(() => {
    if (!group?.entries) return ''
    return group.entries
      .map(entry => {
        const comment = entry.comment ? ` # ${entry.comment}` : ''
        const disabled = entry.enabled ? '' : '# '
        return `${disabled}${entry.ip} ${entry.hostname}${comment}`
      })
      .join('\n')
  }, [group])

  // 初始化 memo 内容
  useEffect(() => {
    setMemoContent(entriesToHostsFormat)
    setValidationErrors([])
    setShowErrors(false)
  }, [entriesToHostsFormat])

  // 解析 memo 内容为 entries
  const parseMemoToEntries = (content: string): Array<{ ip: string; hostname: string; comment: string; enabled: boolean }> => {
    return content
      .split('\n')
      .filter(line => line.trim())
      .map(line => {
        const trimmed = line.trim()
        const enabled = !trimmed.startsWith('#')
        const cleanLine = enabled ? trimmed : trimmed.slice(1).trim()

        // 分离注释
        const commentIndex = cleanLine.indexOf('#')
        let mainContent = cleanLine
        let comment = ''

        if (commentIndex !== -1) {
          mainContent = cleanLine.substring(0, commentIndex).trim()
          comment = cleanLine.substring(commentIndex + 1).trim()
        }

        // 解析 IP 和 hostname
        const parts = mainContent.split(/\s+/)
        const ip = parts[0] || ''
        const hostname = parts.slice(1).join(' ') || ''

        return { ip, hostname, comment, enabled }
      })
      .filter(entry => entry.ip && entry.hostname)
  }

  // 验证所有内容
  const validateContent = (content: string): string[] => {
    const errors: string[] = []
    const lines = content.split('\n').filter(line => line.trim())

    lines.forEach((line, index) => {
      const validation = validateHostsLine(line)
      if (!validation.valid && validation.error) {
        errors.push(`行 ${index + 1}: ${validation.error}`)
      }
    })

    return errors
  }

  // 处理光标离开时的自动保存（静默模式，不显示提示）
  const handleBlur = () => {
    const errors = validateContent(memoContent)

    if (errors.length > 0) {
      setValidationErrors(errors)
      setShowErrors(true)
      return
    }

    setValidationErrors([])
    setShowErrors(false)

    const entries = parseMemoToEntries(memoContent)
    // 静默保存，不触发应用配置
    onBatchUpdate(entries, true)
  }

  if (!group) {
    return (
      <div className="flex h-full items-center justify-center text-muted-foreground w-full">
        {t('mainPanel.selectGroup')}
      </div>
    )
  }

  const handleSaveMemo = async () => {
    const errors = validateContent(memoContent)

    if (errors.length > 0) {
      setValidationErrors(errors)
      setShowErrors(true)
      return
    }

    setValidationErrors([])
    setShowErrors(false)

    // 先保存修改（静默模式，不显示提示）
    const entries = parseMemoToEntries(memoContent)
    await onBatchUpdate(entries, true)

    // 然后应用配置（会显示成功提示）
    await onApply()
  }

  const handleReset = () => {
    setMemoContent(entriesToHostsFormat)
    setValidationErrors([])
    setShowErrors(false)
  }

  return (
    <div className="flex h-full flex-1 flex-col bg-background" data-testid="main-panel">
      {/* 头部 */}
      <div className="flex items-center justify-between border-b px-6 py-4">
        <div className="flex-1">
          <h2 className="text-lg font-semibold">{group.name}</h2>
          {group.description && (
            <p className="text-sm text-muted-foreground">{group.description}</p>
          )}
          {showErrors && validationErrors.length > 0 && (
            <div className="mt-2 flex items-start gap-2 rounded-md bg-destructive/10 p-2 text-sm text-destructive">
              <AlertCircle className="h-4 w-4 flex-shrink-0 mt-0.5" />
              <div className="flex-1">
                <p className="font-medium">{t('mainPanel.validation.title')}</p>
                <ul className="mt-1 list-inside list-disc space-y-0.5">
                  {validationErrors.map((error, index) => (
                    <li key={index}>{error}</li>
                  ))}
                </ul>
              </div>
            </div>
          )}
        </div>
        <div className="flex gap-2">
          <Button
            variant="outline"
            onClick={handleReset}
          >
            {t('mainPanel.reset')}
          </Button>
          <Button
            variant="outline"
            className={"w-[220px]"}
            onClick={onPreview}
          >
            <Eye className="mr-2 h-4 " />
            {t('mainPanel.preview')}
          </Button>
          <Button onClick={handleSaveMemo} className="w-[120px]">
            <Save className="mr-2 h-4 w-4" />
            {t('mainPanel.apply')}
          </Button>
        </div>
      </div>

      {/* Memo 编辑区域 */}
      <div className="flex-1 overflow-auto p-6">
        <Textarea
          ref={textareaRef}
          value={memoContent}
          onChange={e => {
            setMemoContent(e.target.value)
            // 实时验证并隐藏错误（如果用户正在编辑）
            if (showErrors) {
              const errors = validateContent(e.target.value)
              setValidationErrors(errors)
            }
          }}
          onBlur={handleBlur}
          placeholder="127.0.0.1 localhost&#10;192.168.0.1 www.mytest.com"
          className={`h-full w-full resize-none font-mono text-sm ${
            showErrors && validationErrors.length > 0 ? 'border-destructive' : ''
          }`}
          spellCheck={false}
        />
      </div>
    </div>
  )
}
