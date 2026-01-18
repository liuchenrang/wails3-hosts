import { useState, useRef } from 'react'
import { useTranslation } from 'react-i18next'
import { Plus, Power, Trash2, Edit, AlertTriangle, GripVertical } from 'lucide-react'
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
  onUpdateGroup: (id: string, name: string, description: string) => void
  onDeleteGroup: (id: string) => void
  onToggleGroup: (id: string, enabled: boolean) => void
  onDoubleClickGroup: (group: HostsGroup) => void
  onReorderGroups: (groupIds: string[]) => void
}
// 侧边栏组件
export function Sidebar({
  groups,
  selectedGroupId,
  onSelectGroup,
  onCreateGroup,
  onUpdateGroup,
  onDeleteGroup,
  onToggleGroup,
  onDoubleClickGroup,
  onReorderGroups,
}: SidebarProps) {
  const { t } = useTranslation()
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)
  const [isEditModalOpen, setIsEditModalOpen] = useState(false)
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false)
  const [editingGroup, setEditingGroup] = useState<HostsGroup | null>(null)
  const [deletingGroup, setDeletingGroup] = useState<HostsGroup | null>(null)
  const [newGroupName, setNewGroupName] = useState('')
  const [newGroupDesc, setNewGroupDesc] = useState('')

  // 拖动相关状态
  const [draggedId, setDraggedId] = useState<string | null>(null)
  const [dragOverId, setDragOverId] = useState<string | null>(null)
  const dragStartPos = useRef<{ x: number; y: number }>({ x: 0, y: 0 })
  const isDragging = useRef(false)

  // 拖动开始
  const handleDragStart = (e: React.DragEvent, groupId: string) => {
    setDraggedId(groupId)
    isDragging.current = true
    dragStartPos.current = { x: e.clientX, y: e.clientY }

    // 设置拖动效果
    e.dataTransfer.effectAllowed = 'move'
    e.dataTransfer.setData('text/plain', groupId)

    // 设置拖动时的透明度
    setTimeout(() => {
      e.target.closest('.group-item')?.classList.add('opacity-50')
    }, 0)
  }

  // 拖动结束
  const handleDragEnd = (e: React.DragEvent) => {
    isDragging.current = false
    setDraggedId(null)
    setDragOverId(null)

    // 移除透明度
    document.querySelectorAll('.group-item.opacity-50').forEach(el => {
      el.classList.remove('opacity-50')
    })
  }

  // 拖动经过
  const handleDragOver = (e: React.DragEvent, targetId: string) => {
    e.preventDefault()
    e.dataTransfer.dropEffect = 'move'

    if (draggedId && draggedId !== targetId) {
      setDragOverId(targetId)
    }
  }

  // 拖动离开
  const handleDragLeave = (e: React.DragEvent) => {
    // 只有真正离开目标元素时才清除
    const rect = (e.currentTarget as HTMLElement).getBoundingClientRect()
    const x = e.clientX
    const y = e.clientY

    if (x < rect.left || x > rect.right || y < rect.top || y > rect.bottom) {
      setDragOverId(null)
    }
  }

  // 放置
  const handleDrop = async (e: React.DragEvent, targetId: string) => {
    e.preventDefault()

    if (!draggedId || draggedId === targetId) {
      setDragOverId(null)
      return
    }

    // 计算新顺序
    const oldIndex = groups.findIndex(g => g.id === draggedId)
    const newIndex = groups.findIndex(g => g.id === targetId)

    // 创建新的顺序数组
    const newGroups = [...groups]
    const [removed] = newGroups.splice(oldIndex, 1)
    newGroups.splice(newIndex, 0, removed)

    // 调用父组件的排序函数
    onReorderGroups(newGroups.map(g => g.id))

    setDragOverId(null)
  }

  // 点击处理（区分点击和拖动）
  const handleGroupClick = (e: React.MouseEvent, group: HostsGroup) => {
    // 如果正在拖动，不触发点击
    if (isDragging.current) {
      isDragging.current = false
      return
    }

    // 计算鼠标移动距离
    const moveDistance = Math.sqrt(
      Math.pow(e.clientX - dragStartPos.current.x, 2) +
      Math.pow(e.clientY - dragStartPos.current.y, 2)
    )

    // 如果移动距离小于 5px，视为点击
    if (moveDistance < 5) {
      onSelectGroup(group)
    }
  }

  const handleCreateGroup = () => {
    if (newGroupName.trim()) {
      onCreateGroup(newGroupName.trim(), newGroupDesc.trim())
      setNewGroupName('')
      setNewGroupDesc('')
      setIsCreateModalOpen(false)
    }
  }

  const handleEditGroup = (group: HostsGroup) => {
    setEditingGroup(group)
    setNewGroupName(group.name)
    setNewGroupDesc(group.description)
    setIsEditModalOpen(true)
  }

  const handleUpdateGroup = () => {
    if (newGroupName.trim() && editingGroup) {
      onUpdateGroup(editingGroup.id, newGroupName.trim(), newGroupDesc.trim())
      setNewGroupName('')
      setNewGroupDesc('')
      setEditingGroup(null)
      setIsEditModalOpen(false)
    }
  }

  const handleCloseEditModal = () => {
    setNewGroupName('')
    setNewGroupDesc('')
    setEditingGroup(null)
    setIsEditModalOpen(false)
  }

  const handleDeleteClick = (group: HostsGroup) => {
    setDeletingGroup(group)
    setIsDeleteModalOpen(true)
  }

  const handleConfirmDelete = () => {
    if (deletingGroup) {
      onDeleteGroup(deletingGroup.id)
      setDeletingGroup(null)
      setIsDeleteModalOpen(false)
    }
  }

  const handleCloseDeleteModal = () => {
    setDeletingGroup(null)
    setIsDeleteModalOpen(false)
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
            {groups.map((group) => {
              const isSelected = selectedGroupId === group.id
              const entryCount = getEntryCount(group)
              const isDragged = draggedId === group.id
              const isDragOver = dragOverId === group.id

              return (
                <div
                  key={group.id}
                  className={cn(
                    'group relative rounded-lg border transition-all cursor-pointer group-item',
                    isSelected
                      ? 'border-primary bg-accent shadow-sm'
                      : 'border-border bg-card hover:border-primary/50 hover:bg-accent/50',
                    // 拖动视觉反馈
                    isDragged && 'opacity-50',
                    isDragOver && 'border-primary border-dashed bg-accent/50'
                  )}
                  onClick={(e) => handleGroupClick(e, group)}
                  onDoubleClick={() => onDoubleClickGroup(group)}
                  onDragOver={(e) => handleDragOver(e, group.id)}
                  onDragLeave={handleDragLeave}
                  onDrop={(e) => handleDrop(e, group.id)}
                >
                  {/* 分组头部 */}
                  <div className="flex items-center p-3">
                    {/* 拖动手柄 */}
                    <div
                      draggable={true}
                      onDragStart={(e) => handleDragStart(e, group.id)}
                      onDragEnd={handleDragEnd}
                      className="mr-2 cursor-grab active:cursor-grabbing opacity-0 group-hover:opacity-100 hover:bg-accent/50 rounded w-[24px] h-[24px] flex items-center justify-center transition-opacity"
                    >
                      <GripVertical className="h-4 w-4 text-muted-foreground" />
                    </div>

                    {/* 状态图标 */}
                    <button
                      draggable={false}
                      onDragStart={(e) => e.stopPropagation()}
                      className="mr-3 flex-shrink-0 w-[30px] flex justify-center items-center rounded hover:bg-accent/50"
                      onClick={(e) => {
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
                      <div className="flex items-center gap-3">
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
                    <div className="flex items-center gap-0.5 opacity-0 transition-opacity group-hover:opacity-100 ml-2">
                      {/* 编辑按钮 */}
                      <button
                        draggable={false}
                        onDragStart={(e) => e.stopPropagation()}
                        className="rounded hover:bg-accent w-[30px]"
                        onClick={(e) => {
                          e.stopPropagation()
                          handleEditGroup(group)
                        }}
                        title="编辑分组"
                      >
                        <Edit className="h-3.5 w-3.5 text-muted-foreground hover:text-foreground" />
                      </button>

                      {/* 删除按钮 */}
                      <button
                        draggable={false}
                        onDragStart={(e) => e.stopPropagation()}
                        className="rounded hover:bg-destructive/10  w-[30px] flex justify-center items-center"
                        onClick={(e) => {
                          e.stopPropagation()
                          handleDeleteClick(group)
                        }}
                        title="删除分组"
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

      {/* 编辑分组模态框 */}
      <Modal
        isOpen={isEditModalOpen}
        onClose={handleCloseEditModal}
        title="编辑分组"
        footer={
          <>
            <Button
              variant="outline"
              onClick={handleCloseEditModal}
            >
              {t('common.cancel')}
            </Button>
            <Button onClick={handleUpdateGroup}>{t('common.confirm')}</Button>
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

      {/* 删除确认模态框 */}
      <Modal
        isOpen={isDeleteModalOpen}
        onClose={handleCloseDeleteModal}
        title="确认删除分组"
        footer={
          <>
            <Button
              variant="outline"
              
              onClick={handleCloseDeleteModal}
            >
              {t('common.cancel')}
            </Button>
            <Button
              variant="destructive"
              onClick={handleConfirmDelete}
            >
              确认删除
            </Button>
          </>
        }
      >
        <div className="space-y-4">
          <div className="flex items-start gap-3">
            <AlertTriangle className="h-5 w-5 text-destructive flex-shrink-0 mt-0.5" />
            <div className="flex-1">
              <p className="text-sm font-medium text-foreground">
                确定要删除分组"{deletingGroup?.name}"吗?
              </p>
              <p className="mt-2 text-sm text-muted-foreground">
                此操作将永久删除该分组及其所有条目,此操作不可恢复。
              </p>
              {deletingGroup && deletingGroup.entries.length > 0 && (
                <p className="mt-2 text-sm text-destructive">
                  该分组包含 {deletingGroup.entries.length} 个条目
                </p>
              )}
            </div>
          </div>
        </div>
      </Modal>
    </div>
  )
}
