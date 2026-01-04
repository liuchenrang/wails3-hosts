import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { Clock, RotateCcw } from 'lucide-react'
import { HostsVersion } from '../../types/hosts'
import { Button } from '../ui/Button'
import { Modal } from '../ui/Modal'
import { Input } from '../ui/Input'

interface VersionHistoryProps {
  isOpen: boolean
  onClose: () => void
  versions: HostsVersion[]
  onRollback: (versionId: string, password: string) => void
}

// 版本历史组件
// 单一职责: 显示和操作版本历史
export function VersionHistory({ isOpen, onClose, versions, onRollback }: VersionHistoryProps) {
  const { t } = useTranslation()
  const [rollbackVersion, setRollbackVersion] = useState<HostsVersion | null>(null)
  const [sudoPassword, setSudoPassword] = useState('')

  const handleRollback = () => {
    if (rollbackVersion && sudoPassword) {
      onRollback(rollbackVersion.id, sudoPassword)
      setRollbackVersion(null)
      setSudoPassword('')
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
        <div className="max-h-96 overflow-auto">
          {versions.length === 0 ? (
            <div className="py-8 text-center text-muted-foreground">
              {t('sidebar.noGroups')}
            </div>
          ) : (
            <div className="space-y-3">
              {versions.map((version, index) => (
                <div
                  key={version.id}
                  className="rounded-lg border bg-card p-4 transition-colors hover:bg-accent"
                >
                  <div className="flex items-start gap-4">
                    <div className="flex-1">
                      <div className="mb-2 flex items-center gap-2">
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
                      <div className="mb-1 text-sm font-medium">{version.description}</div>
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
                    </div>
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => setRollbackVersion(version)}
                    >
                      <RotateCcw className="mr-1 h-3 w-3" />
                      {t('versions.rollback')}
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </Modal>

      {/* 回滚确认模态框 */}
      {rollbackVersion && (
        <Modal
          isOpen={!!rollbackVersion}
          onClose={() => {
            setRollbackVersion(null)
            setSudoPassword('')
          }}
          title={t('versions.rollback')}
          footer={
            <>
              <Button
                variant="outline"
                onClick={() => {
                  setRollbackVersion(null)
                  setSudoPassword('')
                }}
              >
                {t('common.cancel')}
              </Button>
              <Button onClick={handleRollback}>{t('common.confirm')}</Button>
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
          </div>
        </Modal>
      )}
    </>
  )
}
