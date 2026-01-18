import { useState, useEffect } from 'react'
import { useTranslation } from 'react-i18next'
import { Clock, RotateCcw } from 'lucide-react'
import { HostsVersion } from '../../types/hosts'
import { Button } from '../ui/Button'
import { Modal } from '../ui/Modal'
import { Input } from '../ui/Input'
import { useToast } from '../../hooks/useToast'
import React from "react"

interface VersionHistoryProps {
  isOpen: boolean
  onClose: () => void
  versions: HostsVersion[]
  onRollback: (versionId: string, password?: string) => Promise<void>
  checkPasswordCache: () => Promise<boolean>
}

// 版本历史组件
// 单一职责: 显示和操作版本历史
export function VersionHistory({ isOpen, onClose, versions, onRollback, checkPasswordCache }: VersionHistoryProps) {
  const { t } = useTranslation()
  const toast = useToast()
  const [rollbackVersion, setRollbackVersion] = useState<HostsVersion | null>(null)
  const [sudoPassword, setSudoPassword] = useState('')
  const [isRollingBack, setIsRollingBack] = useState(false)
  const [rollbackError, setRollbackError] = useState<string>('')
  // 对话框打开时的实时缓存状态（不使用组件状态）
  const [realtimeCacheStatus, setRealtimeCacheStatus] = useState<boolean | null>(null)
  const [isCheckingRealtimeCache, setIsCheckingRealtimeCache] = useState(false)

  // 点击回滚按钮的处理函数
  const handleRollbackClick = async (version: HostsVersion) => {
    console.log('[VersionHistory] 点击回滚按钮', {
      versionId: version.id
    })

    setIsRollingBack(true)
    setIsCheckingRealtimeCache(true)
    setRealtimeCacheStatus(null)

    try {
      // 实时向后端检查密码缓存状态
      console.log('[VersionHistory] 实时检查密码缓存状态')
      const cached = await checkPasswordCache()
      console.log('[VersionHistory] 密码缓存状态:', cached)

      // 保存实时查询结果到临时状态
      setRealtimeCacheStatus(cached)

      // 显示确认对话框
      console.log('[VersionHistory] 显示确认对话框')
      setRollbackVersion(version)
    } catch (error) {
      console.error('[VersionHistory] 检查密码缓存失败', error)
      // 出错时假设密码未缓存
      setRealtimeCacheStatus(false)
      setRollbackVersion(version)
    } finally {
      setIsCheckingRealtimeCache(false)
      setIsRollingBack(false)
    }
  }

  const handleRollback = async () => {
    if (!rollbackVersion) return

    console.log('[VersionHistory] 开始回滚流程', {
      versionId: rollbackVersion.id,
      hasPassword: !!sudoPassword
    })

    setIsRollingBack(true)
    setRollbackError('')

    try {
      console.log('[VersionHistory] 调用 onRollback 前')
      // 传递用户输入的密码(空字符串也传递,后端会使用缓存的密码)
      await onRollback(rollbackVersion.id, sudoPassword)
      console.log('[VersionHistory] 调用 onRollback 后，准备关闭模态框')

      // 回滚成功后，刷新密码缓存状态
      await refreshCacheStatus()

      // 关闭所有模态框
      setRollbackVersion(null)
      setSudoPassword('')
      onClose()
      console.log('[VersionHistory] 模态框已关闭')
    } catch (error) {
      console.error('[VersionHistory] 回滚失败', error)
      // 回滚失败，保持模态框打开，显示错误信息
      setRollbackError(error instanceof Error ? error.message : t('versions.rollback') + ' ' + t('common.error'))
    } finally {
      console.log('[VersionHistory] finally 块，设置 isRollingBack = false')
      setIsRollingBack(false)
    }
  }

  const getSourceText = (source: string) => {
    const sourceMap: Record<string, string> = {
      manual: t('versions.source_manual'),
      auto: t('versions.source_auto'),
      rollback: t('versions.source_rollback'),
    }
    return sourceMap[source] || source
  }

  // 刷新密码缓存状态的辅助函数
  const refreshCacheStatus = async () => {
    console.log('[VersionHistory] 刷新密码缓存状态')
    setIsCheckingRealtimeCache(true)
    try {
      const cached = await checkPasswordCache()
      console.log('[VersionHistory] 刷新后的缓存状态:', cached)
      setRealtimeCacheStatus(cached)
    } catch (error) {
      console.error('[VersionHistory] 刷新缓存状态失败', error)
      setRealtimeCacheStatus(false)
    } finally {
      setIsCheckingRealtimeCache(false)
    }
  }

  // 版本历史列表打开时，检查密码缓存状态
  useEffect(() => {
    if (isOpen) {
      refreshCacheStatus()
    }
  }, [isOpen])

  return (
    <>
      <Modal
        isOpen={isOpen}
        onClose={onClose}

        title={
          <div className="flex items-center gap-2">
            <Clock className="h-5 w-5" />
            {t('versions.title')}
          </div>
        }
        footer={<Button onClick={onClose}>{t('common.cancel')}</Button>}
      >
        <div className="space-y-2">
          <div className="rounded-lg bg-muted/50 px-3 py-2 text-xs text-muted-foreground">
            {t('versions.maxVersions')}
          </div>
          {/* 密码缓存状态提示 */}
          <div className="rounded-lg px-3 py-2 text-xs">
            {isCheckingRealtimeCache ? (
              <div className="flex items-center gap-2 text-muted-foreground">
                <span>正在检查密码缓存状态...</span>
              </div>
            ) : realtimeCacheStatus === true ? (
              <div className="flex items-center gap-2 text-green-600 dark:text-green-400">
                <span>✓ 密码已缓存,可以直接回滚</span>
              </div>
            ) : realtimeCacheStatus === false ? (
              <div className="flex items-center gap-2 text-yellow-600 dark:text-yellow-400">
                <span>⚠ 密码未缓存,回滚时需要输入密码</span>
              </div>
            ) : null}
          </div>
          <div className="max-h-96 overflow-auto">
            {versions.length === 0 ? (
              <div className="py-8 text-center text-muted-foreground">
                {t('versions.noVersions')}
              </div>
            ) : (
            <div className="space-y-3">
              {versions.map((version, index) => (
                <div
                  key={version.id}
                  className="rounded-lg border bg-card p-4 transition-colors hover:bg-accent"
                >
                  <div className="flex flex-col gap-3">
                    <div className="flex items-center gap-2">
                      <Clock className="h-4 w-4 text-muted-foreground" />
                      <span className="text-sm font-mono text-muted-foreground">
                        #{index + 1}
                      </span>
                      <span className="rounded bg-primary/10 px-2 py-0.5 text-xs text-primary">
                        {getSourceText(version.source)}
                      </span>
                      {version.age > 0 && (
                        <span className="text-xs text-muted-foreground">
                          {version.age} {t('versions.timestamp')}前
                        </span>
                      )}
                    </div>
                    <div className="text-sm font-medium">{version.description}</div>
                    <div className="text-xs text-muted-foreground">
                      {t('versions.timestamp')}: {version.timestamp}
                    </div>
                    {version.content && (
                      <details className="mt-2">
                        <summary className="cursor-pointer text-xs text-muted-foreground hover:text-foreground">
                          查看内容
                        </summary>
                        <pre className="mt-2 max-h-32 overflow-auto rounded bg-muted p-2 text-xs">
                          {version.content.substring(0, 500)}
                          {version.content.length > 500 && '...'}
                        </pre>
                      </details>
                    )}
                    <div className="flex justify-end">
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => handleRollbackClick(version)}
                        disabled={isRollingBack}
                        className="min-w-[4em] whitespace-nowrap"
                      >
                        <RotateCcw className="mr-1 h-4 w-4" />
                        {isRollingBack ? t('common.loading') : t('versions.rollback')}
                      </Button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
          </div>
        </div>
      </Modal>

      {/* 回滚确认模态框 */}
      {rollbackVersion && (
        <Modal
          isOpen={!!rollbackVersion}
          onClose={async () => {
            if (!isRollingBack) {
              setRollbackVersion(null)
              setSudoPassword('')
              setRollbackError('')
              // 关闭时刷新密码缓存状态
              await refreshCacheStatus()
            }
          }}
          title={t('versions.rollback')}
          footer={
            <>
              <Button
                variant="outline"
                onClick={async () => {
                  setRollbackVersion(null)
                  setSudoPassword('')
                  setRollbackError('')
                  // 关闭时刷新密码缓存状态
                  await refreshCacheStatus()
                }}
                disabled={isRollingBack}
              >
                {t('common.cancel')}
              </Button>
              <Button
                onClick={handleRollback}
                disabled={
                  isRollingBack ||
                  isCheckingRealtimeCache ||
                  (!sudoPassword && realtimeCacheStatus === false) ||
                  (realtimeCacheStatus === null)
                }
              >
                {isRollingBack || isCheckingRealtimeCache ? t('common.loading') : t('common.confirm')}
              </Button>
            </>
          }
        >
          <div className="space-y-4">
            <p className="text-sm text-muted-foreground">
              {t('versions.rollbackConfirm')}
            </p>
            <div className="rounded-lg bg-muted p-3">
              <div className="text-sm font-medium">{rollbackVersion.description}</div>
              <div className="text-xs text-muted-foreground">{rollbackVersion.timestamp}</div>
            </div>

            {/* 检查密码缓存状态 */}
            {isCheckingRealtimeCache && (
              <div className="rounded-lg bg-muted p-3 text-center">
                <p className="text-sm text-muted-foreground">正在检查密码缓存状态...</p>
              </div>
            )}

            {/* 密码未缓存时显示输入框 */}
            {!isCheckingRealtimeCache && realtimeCacheStatus === false && (
              <div>
                <label className="block text-sm font-medium">{t('sudo.password')}</label>
                <Input
                  type="password"
                  value={sudoPassword}
                  onChange={e => setSudoPassword(e.target.value)}
                  onKeyDown={e => {
                    if (e.key === 'Enter') {
                      handleRollback()
                    }
                  }}
                  placeholder={t('sudo.passwordPlaceholder')}
                  className="mt-1"
                  disabled={isRollingBack}
                  autoFocus
                />
              </div>
            )}

            {/* 密码已缓存时显示提示 */}
            {!isCheckingRealtimeCache && realtimeCacheStatus === true && (
              <div className="rounded-lg bg-blue-50 dark:bg-blue-950/20 border border-blue-200 dark:border-blue-800 p-3">
                <p className="text-sm text-blue-700 dark:text-blue-300">
                  ✓ {t('sudo.passwordCachedHint')}
                </p>
              </div>
            )}

            {rollbackError && (
              <div className="rounded-lg bg-destructive/10 border border-destructive/20 p-3">
                <p className="text-sm text-destructive">{rollbackError}</p>
              </div>
            )}
          </div>
        </Modal>
      )}
    </>
  )
}
