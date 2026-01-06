import { useEffect, useState, useCallback } from 'react'
import { useTranslation } from 'react-i18next'
import { Moon, Sun, History } from 'lucide-react'
import { HostsGroup, HostsVersion } from './types/hosts'
import { hostsApi } from './api/hosts'
import { useTheme } from './hooks/useTheme'
import { useHotkey } from './hooks/useHotkeys'
import { Sidebar } from './components/layout/Sidebar'
import { MainPanel } from './components/layout/MainPanel'
import { VersionHistory } from './components/hosts/VersionHistory'
import { ConflictAlert } from './components/hosts/ConflictAlert'
import { AboutDialog } from './components/ui/AboutDialog'
import { Button } from './components/ui/Button'
import { Modal } from './components/ui/Modal'
import { Input } from './components/ui/Input'
import { VibeKanbanWebCompanion } from 'vibe-kanban-web-companion'
import './i18n'
import './index.css'
import React from "react"

// 调试：检查 Wails 服务是否可用
if (typeof window !== 'undefined') {
  console.log('🔍 检查 window 对象属性:')

  // 检查所有 Wails 相关的属性
  const wailsKeys = Object.keys(window).filter(k =>
    k.toLowerCase().includes('wails') ||
    k.toLowerCase().includes('call') ||
    k.toLowerCase().includes('invoke') ||
    k.toLowerCase().includes('host')
  )

  console.log('找到 Wails 相关属性:', wailsKeys)
  wailsKeys.forEach(key => {
    console.log(`  - ${key}:`, typeof (window as any)[key])
  })

  // 检查特定的 Wails 对象
  console.log('window.wails:', (window as any).wails)
  console.log('window.Wails:', (window as any).Wails)

  // 列出所有全局函数
  const allKeys = Object.keys(window).filter(k => typeof (window as any)[k] === 'function')
  console.log('所有全局函数 (前20个):', allKeys.slice(0, 20))
}

function App() {
  const { t } = useTranslation()
  const { theme, toggleTheme } = useTheme()
  const [groups, setGroups] = useState<HostsGroup[]>([])
  const [selectedGroupId, setSelectedGroupId] = useState<string | null>(null)
  const [previewContent, setPreviewContent] = useState('')
  const [showPreview, setShowPreview] = useState(false)
  const [showVersions, setShowVersions] = useState(false)
  const [sudoPassword, setSudoPassword] = useState('')
  const [showSudoPrompt, setShowSudoPrompt] = useState(false)
  const [versions, setVersions] = useState<HostsVersion[]>([])
  const [conflicts, setConflicts] = useState({})
  const [showConflicts, setShowConflicts] = useState(false)
  const [showAbout, setShowAbout] = useState(false)
  const [aboutInfo, setAboutInfo] = useState<{ version?: string; email?: string }>({})

  const selectedGroup = groups.find(g => g.id === selectedGroupId) || null

  // 快捷键: Cmd+S / Ctrl+S 保存
  const handleApply = useCallback(() => {
    setShowSudoPrompt(true)
  }, [])

  useHotkey('s', handleApply)

  // 加载分组列表
  useEffect(() => {
    loadGroups()
    loadVersions()

    // 监听来自 Wails 后端的事件
    if (typeof window !== 'undefined' && (window as any).EventsOn) {
      // 监听"关于我们"对话框事件
      ;(window as any).EventsOn('show-about-dialog', (data: { version?: string; email?: string }) => {
        setAboutInfo(data)
        setShowAbout(true)
      })

      // 监听版本历史事件
      ;(window as any).EventsOn('show-version-history', () => {
        setShowVersions(true)
      })
    }
  }, [])

  const loadGroups = async () => {
    try {
      const data = await hostsApi.getAllGroups()
      setGroups(data)
    } catch (error) {
      console.error('Failed to load groups:', error)
    }
  }

  const handleCreateGroup = async (name: string, description: string) => {
    try {
      await hostsApi.createGroup(name, description)
      await loadGroups()
    } catch (error) {
      console.error('Failed to create group:', error)
      alert('创建分组失败')
    }
  }

  const handleDeleteGroup = async (id: string) => {
    try {
      await hostsApi.deleteGroup(id)
      if (selectedGroupId === id) {
        setSelectedGroupId(null)
      }
      await loadGroups()
    } catch (error) {
      console.error('Failed to delete group:', error)
      alert('删除分组失败')
    }
  }

  const handleToggleGroup = async (id: string, enabled: boolean) => {
    try {
      await hostsApi.toggleGroup(id, enabled)
      await loadGroups()
    } catch (error) {
      console.error('Failed to toggle group:', error)
      alert('切换分组状态失败')
    }
  }

  const handleAddEntry = async (ip: string, hostname: string, comment: string) => {
    if (!selectedGroupId) return
    try {
      await hostsApi.addEntry(selectedGroupId, ip, hostname, comment)
      await loadGroups()
    } catch (error) {
      console.error('Failed to add entry:', error)
      alert('添加条目失败')
    }
  }

  const handleUpdateEntry = async (entryId: string, ip: string, hostname: string, comment: string) => {
    if (!selectedGroupId) return
    try {
      await hostsApi.updateEntry(selectedGroupId, entryId, ip, hostname, comment)
      await loadGroups()
    } catch (error) {
      console.error('Failed to update entry:', error)
      alert('更新条目失败')
    }
  }

  const handleDeleteEntry = async (entryId: string) => {
    if (!selectedGroupId) return
    try {
      await hostsApi.deleteEntry(selectedGroupId, entryId)
      await loadGroups()
    } catch (error) {
      console.error('Failed to delete entry:', error)
      alert('删除条目失败')
    }
  }

  const handlePreview = async () => {
    try {
      const content = await hostsApi.generatePreview()
      setPreviewContent(content)
      setShowPreview(true)
    } catch (error) {
      console.error('Failed to generate preview:', error)
      alert('生成预览失败')
    }
  }

  const handleConfirmApply = async () => {
    try {
      // 检测冲突
      const detectedConflicts = await hostsApi.detectConflicts()
      if (Object.keys(detectedConflicts).length > 0) {
        setConflicts(detectedConflicts)
        setShowConflicts(true)
        return
      }

      await hostsApi.applyHosts(sudoPassword)
      setShowSudoPrompt(false)
      setSudoPassword('')
      await loadVersions()
      alert(t('common.success'))
    } catch (error) {
      console.error('Failed to apply hosts:', error)
      alert(t('common.error'))
    }
  }

  const loadVersions = async () => {
    try {
      const data = await hostsApi.getVersions(50)
      setVersions(data)
    } catch (error) {
      console.error('Failed to load versions:', error)
    }
  }

  const handleRollback = async (versionId: string, password: string) => {
    try {
      await hostsApi.rollbackToVersion(versionId, password)
      alert('回滚成功')
      await loadVersions()
      await loadGroups()
    } catch (error) {
      console.error('Failed to rollback:', error)
      alert('回滚失败')
    }
  }

  const handleIgnoreConflicts = () => {
    setShowConflicts(false)
    // 继续应用
    handleConfirmApply()
  }

  const handleBatchUpdateEntries = async (entries: Array<{ ip: string; hostname: string; comment: string; enabled: boolean }>) => {
    if (!selectedGroupId) return
    try {
      await hostsApi.batchUpdateEntries(selectedGroupId, entries)
      await loadGroups()
    } catch (error) {
      console.error('Failed to batch update entries:', error)
      alert('批量更新失败')
      throw error
    }
  }

  return (
    <>
      <VibeKanbanWebCompanion />
      <div className="flex h-screen flex-col bg-background text-foreground">
        {/* 顶部栏 */}
        <div className="flex items-center justify-between border-b px-6 py-3">
          <div>
            <h1 className="text-xl font-bold">{t('app.title')}</h1>
            <p className="text-sm text-muted-foreground">{t('app.subtitle')}</p>
          </div>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setShowVersions(true)}
            >
              <History className="mr-2 h-4 w-4" />
              {t('versions.title')}
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={toggleTheme}
            >
              {theme === 'light' ? (
                <Moon className="h-4 w-4" />
              ) : (
                <Sun className="h-4 w-4" />
              )}
            </Button>
          </div>
        </div>

        {/* 主内容区 */}
        <div className="flex flex-1 overflow-hidden">
          {/* 左侧分组列表 - 紧贴左侧 */}
          <Sidebar
              groups={groups}
              selectedGroupId={selectedGroupId}
              // onSelectGroup={setSelectedGroupId}
              onCreateGroup={handleCreateGroup}
              onDeleteGroup={handleDeleteGroup}
              onToggleGroup={handleToggleGroup} onSelectGroup={function (group: HostsGroup): void {
            throw new Error("Function not implemented.")
          }}          />

          {/* 右侧主面板 */}
          <MainPanel
            group={selectedGroup}
            onUpdateEntry={handleUpdateEntry}
            onAddEntry={handleAddEntry}
            onDeleteEntry={handleDeleteEntry}
            onApply={handleApply}
            onPreview={handlePreview}
            onBatchUpdate={handleBatchUpdateEntries}
          />
        </div>

        {/* 预览模态框 */}
        <Modal
          isOpen={showPreview}
          onClose={() => setShowPreview(false)}
          title={t('mainPanel.preview')}
          footer={
            <Button onClick={() => setShowPreview(false)}>{t('common.cancel')}</Button>
          }
        >
          <pre className="max-h-96 overflow-auto rounded bg-muted p-4 text-sm">
            {previewContent}
          </pre>
        </Modal>

        {/* Sudo 密码模态框 */}
        <Modal
          isOpen={showSudoPrompt}
          onClose={() => setShowSudoPrompt(false)}
          title={t('sudo.title')}
          footer={
            <>
              <Button
                variant="outline"
                onClick={() => setShowSudoPrompt(false)}
              >
                {t('common.cancel')}
              </Button>
              <Button onClick={handleConfirmApply}>{t('mainPanel.apply')}</Button>
            </>
          }
        >
          <div className="space-y-4">
            <p className="text-sm text-muted-foreground">{t('sudo.description')}</p>
            <div>
              <label className="block text-sm font-medium">{t('sudo.password')}</label>
              <Input
                type="password"
                value={sudoPassword}
                onChange={e => setSudoPassword(e.target.value)}
                placeholder={t('sudo.passwordPlaceholder')}
                className="mt-1"
              />
            </div>
            <p className="text-xs text-muted-foreground">{t('mainPanel.applyShortcut')}</p>
          </div>
        </Modal>

        {/* 冲突警告 */}
        {showConflicts && (
          <ConflictAlert
            conflicts={conflicts}
            onIgnore={handleIgnoreConflicts}
            onResolve={() => setShowConflicts(false)}
          />
        )}

        {/* 版本历史 */}
        <VersionHistory
          isOpen={showVersions}
          onClose={() => setShowVersions(false)}
          versions={versions}
          onRollback={handleRollback}
        />

        {/* 关于我们对话框 */}
        <AboutDialog
          isOpen={showAbout}
          onClose={() => setShowAbout(false)}
          version={aboutInfo.version}
          email={aboutInfo.email}
        />
      </div>
    </>
  )
}

export default App
