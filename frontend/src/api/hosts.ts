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

  async applyHosts(sudoPassword: string): Promise<void> {
    await HostsHandler.ApplyHosts(sudoPassword)
  },

  // ========== 版本历史 ==========

  async getVersions(limit: number = 50): Promise<HostsVersion[]> {
    const result = await HostsHandler.GetVersions(limit)
    return result as unknown as HostsVersion[]
  },

  async rollbackToVersion(versionID: string, sudoPassword: string): Promise<void> {
    await HostsHandler.RollbackToVersion(versionID, sudoPassword)
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
    return await HostsHandler.IsSudoPasswordCached()
  },
}
