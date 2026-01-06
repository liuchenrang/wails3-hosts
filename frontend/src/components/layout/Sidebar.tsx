import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Plus, Power, Trash2, Edit } from 'lucide-react'
import { HostsGroup } from '../../types/hosts'
import { Button } from '../ui/Button'
import { Modal } from '../ui/Modal'
import { Input } from '../ui/Input'
import { cn } from '../../utils/cn'
import React from "react"

interface SidebarProps {
  groups: HostsGroup[]
  selectedGroupId: string | null
  onSelectGroup: (group: HostsGroup) => void
  onCreateGroup: (name: string, description: string) => void
  onDeleteGroup: (id: string) => void
  onToggleGroup: (id: string, enabled: boolean) => void
  onDoubleClickGroup: (group: HostsGroup) => void
}
// 侧边栏组件
export function Sidebar({
  groups,
  selectedGroupId,
  onSelectGroup,
  onCreateGroup,
  onDeleteGroup,
  onToggleGroup,
  onDoubleClickGroup,
}: SidebarProps) {
  const { t } = useTranslation()
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)
  const [newGroupName, setNewGroupName] = useState('')
  const [newGroupDesc, setNewGroupDesc] = useState('')

  const handleCreateGroup = () => {
    if (newGroupName.trim()) {
      onCreateGroup(newGroupName.trim(), newGroupDesc.trim())
      setNewGroupName('')
      setNewGroupDesc('')
      setIsCreateModalOpen(false)
    }
  }

  const getEntryCount = (group: HostsGroup) => {
    return group.entries.filter(e => e.enabled).length
  }

  return (
    <div className="flex h-full w-80 flex-shrink-0 flex-col border-r">
      {/* 头部 */}
      <div className="border-b bg-card px-4 py-3">
        <h2 className="text-base font-semibold">{t('sidebar.groups')}</h2>
        <p className="mt-1 text-xs text-muted-foreground">
          {t('sidebar.groupCount', { count: groups.length })}
        </p>
        <Button
          className="mt-3 w-full"
          size="sm"
          onClick={() => setIsCreateModalOpen(true)}
        >
          <Plus className="mr-2 h-3.5 w-3.5" />
          {t('sidebar.createGroup')}
        </Button>
      </div>

      {/* 分组列表 */}
      <div className="flex-1 overflow-y-auto bg-muted/30">
        {groups.length === 0 ? (
          <div className="flex h-full w-full items-center justify-center p-3 text-muted-foreground">
            <div className="text-center">
              <p className="text-sm">{t('sidebar.noGroups')}</p>
              <p className="mt-1 text-xs opacity-60">{t('sidebar.createFirstGroup')}</p>
            </div>
          </div>
        ) : (
          <div className="px-3 py-2 space-y-2">
            {groups.map(group => {
              const isSelected = selectedGroupId === group.id
              const entryCount = getEntryCount(group)

              return (
                <div
                  key={group.id}
                  className={cn(
                    'group relative rounded-lg border transition-all cursor-pointer',
                    isSelected
                      ? 'border-primary bg-accent shadow-sm'
                      : 'border-border bg-card hover:border-primary/50 hover:bg-accent/50'
                  )}
                  onClick={() => onSelectGroup(group)}
                  onDoubleClick={() => onDoubleClickGroup(group)}
                >
                  {/* 分组头部 */}
                  <div className="flex items-center p-3">
                    {/* 状态图标 */}
                    <button
                      className="mr-3 flex-shrink-0"
                      onClick={e => {
                        e.stopPropagation()
                        onToggleGroup(group.id, !group.is_enabled)
                      }}
                    >
                      <Power
                        className={cn(
                          'h-4 w-4 transition-colors',
                          group.is_enabled ? 'text-green-500' : 'text-muted-foreground'
                        )}
                      />
                    </button>

                    {/* 分组信息 */}
                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2">
                        <span className="font-medium truncate text-sm">{group.name}</span>
                        {/* 条目数量徽章 */}
                        {entryCount > 0 && (
                          <span className="flex h-5 flex-shrink-0 items-center rounded-full bg-primary/10 px-2 text-xs text-primary">
                            {entryCount}
                          </span>
                        )}
                      </div>
                      {group.description && (
                        <p className="mt-0.5 truncate text-xs text-muted-foreground">
                          {group.description}
                        </p>
                      )}
                    </div>

                    {/* 操作按钮组 */}
                    <div className="flex gap-1 opacity-0 transition-opacity group-hover:opacity-100 ml-2">
                      <button
                        className="rounded p-1.5 hover:bg-accent"
                        onClick={e => {
                          e.stopPropagation()
                        }}
                      >
                        <Edit className="h-3.5 w-3.5 text-muted-foreground hover:text-foreground" />
                      </button>
                      <button
                        className="rounded p-1.5 hover:bg-destructive/10"
                        onClick={e => {
                          e.stopPropagation()
                          if (confirm(t('sidebar.deleteConfirm', { name: group.name }))) {
                            onDeleteGroup(group.id)
                          }
                        }}
                      >
                        <Trash2 className="h-3.5 w-3.5 text-muted-foreground hover:text-destructive" />
                      </button>
                    </div>
                  </div>
                </div>
              )
            })}
          </div>
        )}
      </div>

      {/* 创建分组模态框 */}
      <Modal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        title={t('sidebar.createGroup')}
        footer={
          <>
            <Button
              variant="outline"
              onClick={() => setIsCreateModalOpen(false)}
            >
              {t('common.cancel')}
            </Button>
            <Button onClick={handleCreateGroup}>{t('common.confirm')}</Button>
          </>
        }
      >
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium">{t('form.groupName')}</label>
            <Input
              value={newGroupName}
              onChange={e => setNewGroupName(e.target.value)}
              placeholder="例如：开发环境"
              className="mt-1"
            />
          </div>
          <div>
            <label className="block text-sm font-medium">{t('form.groupDesc')}</label>
            <Input
              value={newGroupDesc}
              onChange={e => setNewGroupDesc(e.target.value)}
              placeholder="例如：本地开发域名映射"
              className="mt-1"
            />
          </div>
        </div>
      </Modal>
    </div>
  )
}
