import { HostsGroup, HostsVersion } from '../types/hosts'
// Wails v3 自动生成的绑定
import * as HostsHandler from '../../bindings/github.com/chen/wails3-hosts/internal/interface/handler/hostshandler'

// Wails v3 API 调用封装
// 单一职责: 封装所有与后端的 API 调用
// DDD: 前端 API 层，通过 Wails 自动生成的绑定调用 Go 服务

export const hostsApi = {
  // ========== 分组管理 ==========

  async createGroup(name: string, description: string): Promise<HostsGroup> {
    const result = await HostsHandler.CreateGroup(name, description)
    if (!result) {
      throw new Error('Failed to create group')
    }
    return result as unknown as HostsGroup
  },

  async getAllGroups(): Promise<HostsGroup[]> {
    const result = await HostsHandler.GetAllGroups()
    return result as unknown as HostsGroup[]
  },

  async getGroupByID(id: string): Promise<HostsGroup> {
    const result = await HostsHandler.GetGroupByID(id)
    if (!result) {
      throw new Error('Group not found')
    }
    return result as unknown as HostsGroup
  },

  async updateGroup(id: string, name: string, description: string): Promise<void> {
    await HostsHandler.UpdateGroup(id, name, description)
  },

  async deleteGroup(id: string): Promise<void> {
    await HostsHandler.DeleteGroup(id)
  },

  async toggleGroup(id: string, enabled: boolean): Promise<void> {
    await HostsHandler.ToggleGroup(id, enabled)
  },

  // TODO: Wails绑定文件将在运行时自动生成ReorderGroups
  async reorderGroups(groupIds: string[]): Promise<void> {
    // 临时使用动态调用
    const handler = HostsHandler as any
    if (handler.ReorderGroups) {
      await handler.ReorderGroups(groupIds)
    } else {
      console.warn('ReorderGroups method not available in binding yet')
      // 如果方法不存在,使用API调用
      await fetch('/wails/runtime', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          method: 'github.com/chen/wails3-hosts/internal/interface/handler.HostsHandler.ReorderGroups',
          args: [groupIds]
        })
      })
    }
  },

  // ========== 条目管理 ==========

  async addEntry(groupID: string, ip: string, hostname: string, comment: string): Promise<void> {
    await HostsHandler.AddEntry(groupID, ip, hostname, comment)
  },

  async updateEntry(
    groupID: string,
    entryID: string,
    ip: string,
    hostname: string,
    comment: string
  ): Promise<void> {
    await HostsHandler.UpdateEntry(groupID, entryID, ip, hostname, comment)
  },

  async deleteEntry(groupID: string, entryID: string): Promise<void> {
    await HostsHandler.DeleteEntry(groupID, entryID)
  },

  async batchUpdateEntries(
    groupID: string,
    entries: Array<{ ip: string; hostname: string; comment: string; enabled: boolean }>
  ): Promise<void> {
    await HostsHandler.BatchUpdateEntries(groupID, entries as any)
  },

  // ========== 配置应用 ==========

  async generatePreview(): Promise<string> {
    return await HostsHandler.GeneratePreview()
  },

  async detectConflicts(): Promise<Record<string, string[]>> {
    return await HostsHandler.DetectConflicts()
  },

  async applyHosts(): Promise<void> {
    // 注意：不再传递密码参数
    // 密码必须先通过 validateSudoPassword() 验证
    await HostsHandler.ApplyHosts()
  },

  // ========== 版本历史 ==========

  async getVersions(limit: number = 50): Promise<HostsVersion[]> {
    const result = await HostsHandler.GetVersions(limit)
    return result as unknown as HostsVersion[]
  },

  async rollbackToVersion(versionID: string, sudoPassword: string): Promise<void> {
    console.log('[hostsApi] rollbackToVersion 调用', { versionID, hasPassword: sudoPassword !== '', passwordLength: sudoPassword.length })
    try {
      const result = await HostsHandler.RollbackToVersion(versionID, sudoPassword)
      console.log('[hostsApi] rollbackToVersion 完成', { result })
      return result
    } catch (error) {
      console.error('[hostsApi] rollbackToVersion 失败', error)
      throw error
    }
  },

  // ========== Sudo 管理 ==========

  async validateSudoPassword(password: string): Promise<{ valid: boolean; error: string }> {
    const result = await HostsHandler.ValidateSudoPassword(password)
    // Wails v3 返回数组 [valid, error]
    if (Array.isArray(result)) {
      return { valid: result[0], error: result[1] }
    }
    return result as unknown as { valid: boolean; error: string }
  },

  async isSudoPasswordCached(): Promise<boolean> {
    console.log('[hostsApi] isSudoPasswordCached 调用')
    const result = await HostsHandler.IsSudoPasswordCached()
    console.log('[hostsApi] isSudoPasswordCached 结果', { cached: result })
    return result
  },

  async getPlatformInfo(): Promise<{
    os: string
    arch: string
    needsSudo: boolean
    canCacheCred: boolean
  }> {
    console.log('[hostsApi] getPlatformInfo 调用')
    const result = await HostsHandler.GetPlatformInfo()
    console.log('[hostsApi] getPlatformInfo 结果', result)
    return result as unknown as {
      os: string
      arch: string
      needsSudo: boolean
      canCacheCred: boolean
    }
  },
}
