import { AlertTriangle, X } from 'lucide-react'
import { ConflictInfo } from '../../types/app'
import { Button } from '../ui/Button'
import { cn } from '../../utils/cn'

interface ConflictAlertProps {
  conflicts: ConflictInfo
  onIgnore: () => void
  onResolve: () => void
}

// 冲突警告组件
// 单一职责: 显示 hosts 配置冲突信息
export function ConflictAlert({ conflicts, onIgnore, onResolve }: ConflictAlertProps) {
  const hostnames = Object.keys(conflicts)

  if (hostnames.length === 0) {
    return null
  }

  return (
    <div className="mx-4 mb-4 rounded-lg border border-destructive/50 bg-destructive/10 p-4">
      <div className="flex items-start gap-3">
        <AlertTriangle className="h-5 w-5 flex-shrink-0 text-destructive" />
        <div className="flex-1">
          <h3 className="mb-2 font-semibold text-destructive">配置冲突检测</h3>
          <p className="mb-3 text-sm text-destructive-foreground">
            以下主机名存在多个 IP 映射，可能会导致 DNS 解析混乱：
          </p>
          <div className="space-y-2">
            {hostnames.map(hostname => (
              <div
                key={hostname}
                className="rounded border border-destructive/20 bg-destructive/5 p-2"
              >
                <div className="mb-1 font-mono text-sm font-medium">{hostname}</div>
                <div className="flex flex-wrap gap-1 text-xs">
                  {conflicts[hostname].map(ip => (
                    <span
                      key={ip}
                      className="rounded bg-destructive/20 px-2 py-0.5 font-mono"
                    >
                      {ip}
                    </span>
                  ))}
                </div>
              </div>
            ))}
          </div>
          <div className="mt-3 flex gap-2">
            <Button
              size="sm"
              variant="outline"
              onClick={onIgnore}
            >
              忽略并继续
            </Button>
            <Button
              size="sm"
              variant="destructive"
              onClick={onResolve}
            >
              解决冲突
            </Button>
          </div>
        </div>
        <Button
          size="sm"
          variant="ghost"
          onClick={onIgnore}
          className="h-6 w-6 p-0"
        >
          <X className="h-4 w-4" />
        </Button>
      </div>
    </div>
  )
}
