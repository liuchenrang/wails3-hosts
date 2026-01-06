import { useTranslation } from 'react-i18next'
import { Modal } from './Modal'
import { Button } from './Button'

interface AboutDialogProps {
  isOpen: boolean
  onClose: () => void
  version?: string
  email?: string
}

export function AboutDialog({ isOpen, onClose, version, email }: AboutDialogProps) {
  const { t } = useTranslation()

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={t('about.title')}
      footer={
        <Button onClick={onClose}>{t('common.cancel')}</Button>
      }
    >
      <div className="space-y-6 py-4">
        {/* 应用信息 */}
        <div className="text-center">
          <h2 className="text-2xl font-bold">{t('app.title')}</h2>
          <p className="mt-2 text-sm text-muted-foreground">{t('about.description')}</p>
        </div>

        {/* 版本信息 */}
        <div className="rounded-lg border bg-muted/30 p-4">
          <div className="space-y-3">
            <div className="flex items-center justify-between border-b pb-3">
              <span className="text-sm font-medium text-muted-foreground">{t('about.version')}</span>
              <span className="text-sm font-mono font-semibold">{version || '1.0.0'}</span>
            </div>
            <div className="flex items-center justify-between pt-1">
              <span className="text-sm font-medium text-muted-foreground">{t('about.email')}</span>
              <a
                href={`mailto:${email || 'support@hostsmanager.com'}`}
                className="text-sm font-mono text-primary hover:underline"
              >
                {email || 'support@hostsmanager.com'}
              </a>
            </div>
          </div>
        </div>

        {/* 技术栈信息 */}
        <div className="text-center text-xs text-muted-foreground">
          <p>基于 Wails v3 + React + TypeScript 构建</p>
          <p className="mt-1">© 2025 Hosts Manager. All rights reserved.</p>
        </div>
      </div>
    </Modal>
  )
}
