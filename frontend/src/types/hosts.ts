// Hosts 相关类型定义

export interface HostsGroup {
  id: string
  name: string
  description: string
  is_enabled: boolean
  entries: HostsEntry[]
  created_at: string
  updated_at: string
}

export interface HostsEntry {
  id: string
  ip: string
  hostname: string
  comment: string
  enabled: boolean
}

export interface HostsVersion {
  id: string
  timestamp: string
  content: string
  description: string
  source: string
  age: number
}

export type VersionSource = 'manual' | 'auto' | 'rollback'
