import { useState, useMemo, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { Save, Eye } from 'lucide-react'
import { HostsGroup } from '../../types/hosts'
import { Button } from '../ui/Button'
import { Textarea } from '../ui/Textarea'
import React from "react"

interface MainPanelProps {
  group: HostsGroup | null
  onUpdateEntry: (entryId: string, ip: string, hostname: string, comment: string) => void
  onAddEntry: (ip: string, hostname: string, comment: string) => void
  onDeleteEntry: (entryId: string) => void
  onApply: () => void
  onPreview: () => void
  onBatchUpdate: (entries: Array<{ ip: string; hostname: string; comment: string; enabled: boolean }>) => void
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

  if (!group) {
    return (
      <div className="flex h-full items-center justify-center text-muted-foreground">
        {t('mainPanel.selectGroup')}
      </div>
    )
  }

  const handleSaveMemo = () => {
    const entries = parseMemoToEntries(memoContent)
    onBatchUpdate(entries)
    onApply()
  }

  const handleReset = () => {
    setMemoContent(entriesToHostsFormat)
  }

  return (
    <div className="flex h-full flex-1 flex-col bg-background">
      {/* 头部 */}
      <div className="flex items-center justify-between border-b px-6 py-4">
        <div>
          <h2 className="text-lg font-semibold">{group.name}</h2>
          {group.description && (
            <p className="text-sm text-muted-foreground">{group.description}</p>
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
            onClick={onPreview}
          >
            <Eye className="mr-2 h-4 w-4" />
            {t('mainPanel.preview')}
          </Button>
          <Button onClick={handleSaveMemo}>
            <Save className="mr-2 h-4 w-4" />
            {t('mainPanel.apply')}
          </Button>
        </div>
      </div>

      {/* Memo 编辑区域 */}
      <div className="flex-1 overflow-auto p-6">
        <Textarea
          value={memoContent}
          onChange={e => setMemoContent(e.target.value)}
          placeholder="127.0.0.1 localhost&#10;192.168.0.1 www.mytest.com"
          className="h-full w-full resize-none font-mono text-sm"
          spellCheck={false}
        />
      </div>
    </div>
  )
}
