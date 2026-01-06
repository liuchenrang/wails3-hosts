package handler

import (
	"context"

	"github.com/chen/wails3-hosts/internal/application/dto"
	"github.com/chen/wails3-hosts/internal/application/service"
)

// HostsHandler hosts 管理 Wails 服务处理器
// 单一职责: 将应用服务暴露给 Wails 前端
// DDD: 接口层负责外部接口适配，不包含业务逻辑
type HostsHandler struct {
	appService *service.HostsApplicationService
}

// NewHostsHandler 创建处理器实例
// 依赖注入: 通过构造函数注入应用服务
func NewHostsHandler(appService *service.HostsApplicationService) *HostsHandler {
	return &HostsHandler{
		appService: appService,
	}
}

// ========== 分组管理 ==========

// CreateGroup 创建新的 hosts 分组
func (h *HostsHandler) CreateGroup(name, description string) (*dto.HostsGroupDTO, error) {
	req := dto.CreateHostsGroupRequest{
		Name:        name,
		Description: description,
	}
	return h.appService.CreateGroup(context.Background(), req)
}

// GetAllGroups 获取所有分组
func (h *HostsHandler) GetAllGroups() ([]dto.HostsGroupDTO, error) {
	return h.appService.GetAllGroups(context.Background())
}

// GetGroupByID 根据 ID 获取分组
func (h *HostsHandler) GetGroupByID(id string) (*dto.HostsGroupDTO, error) {
	return h.appService.GetGroupByID(context.Background(), id)
}

// UpdateGroup 更新分组
func (h *HostsHandler) UpdateGroup(id, name, description string) error {
	req := dto.UpdateHostsGroupRequest{
		ID:          id,
		Name:        name,
		Description: description,
	}
	return h.appService.UpdateGroup(context.Background(), req)
}

// DeleteGroup 删除分组
func (h *HostsHandler) DeleteGroup(id string) error {
	return h.appService.DeleteGroup(context.Background(), id)
}

// ToggleGroup 切换分组启用状态
func (h *HostsHandler) ToggleGroup(id string, enabled bool) error {
	req := dto.ToggleGroupRequest{
		ID:      id,
		Enabled: enabled,
	}
	return h.appService.ToggleGroup(context.Background(), req)
}

// ========== 条目管理 ==========

// AddEntry 添加 hosts 条目
func (h *HostsHandler) AddEntry(groupID, ip, hostname, comment string) error {
	req := dto.AddEntryRequest{
		GroupID:  groupID,
		IP:       ip,
		Hostname: hostname,
		Comment:  comment,
	}
	return h.appService.AddEntry(context.Background(), req)
}

// UpdateEntry 更新 hosts 条目
func (h *HostsHandler) UpdateEntry(groupID, entryID, ip, hostname, comment string) error {
	req := dto.UpdateEntryRequest{
		GroupID:  groupID,
		EntryID:  entryID,
		IP:       ip,
		Hostname: hostname,
		Comment:  comment,
	}
	return h.appService.UpdateEntry(context.Background(), req)
}

// DeleteEntry 删除 hosts 条目
func (h *HostsHandler) DeleteEntry(groupID, entryID string) error {
	req := dto.DeleteEntryRequest{
		GroupID: groupID,
		EntryID: entryID,
	}
	return h.appService.DeleteEntry(context.Background(), req)
}

// BatchUpdateEntries 批量更新分组中的所有条目
func (h *HostsHandler) BatchUpdateEntries(groupID string, entries []dto.BatchUpdateEntryRequest) error {
	req := dto.BatchUpdateEntriesRequest{
		GroupID: groupID,
		Entries: entries,
	}
	return h.appService.BatchUpdateEntries(context.Background(), req)
}

// ========== 配置应用 ==========

// GeneratePreview 生成 hosts 配置预览
func (h *HostsHandler) GeneratePreview() (string, error) {
	return h.appService.GeneratePreview(context.Background())
}

// DetectConflicts 检测配置冲突
func (h *HostsHandler) DetectConflicts() (map[string][]string, error) {
	return h.appService.DetectConflicts(context.Background())
}

// ApplyHosts 应用 hosts 配置到系统
func (h *HostsHandler) ApplyHosts(sudoPassword string) error {
	req := dto.ApplyHostsRequest{
		SudoPassword: sudoPassword,
	}
	return h.appService.ApplyHosts(context.Background(), req)
}

// ========== 版本历史 ==========

// GetVersions 获取版本历史
func (h *HostsHandler) GetVersions(limit int) ([]dto.HostsVersionDTO, error) {
	return h.appService.GetVersions(context.Background(), limit)
}

// RollbackToVersion 回滚到指定版本
func (h *HostsHandler) RollbackToVersion(versionID, sudoPassword string) error {
	req := dto.RollbackRequest{
		VersionID:    versionID,
		SudoPassword: sudoPassword,
	}
	return h.appService.RollbackToVersion(context.Background(), req)
}

// ========== Sudo 管理 ==========

// ValidateSudoPassword 验证 sudo 密码
func (h *HostsHandler) ValidateSudoPassword(password string) (bool, string) {
	req := dto.ValidateSudoRequest{
		Password: password,
	}
	resp := h.appService.ValidateSudoPassword(context.Background(), req)
	return resp.Valid, resp.Error
}

// IsSudoPasswordCached 检查 sudo 密码是否已缓存
func (h *HostsHandler) IsSudoPasswordCached() bool {
	// 这里需要从 SudoManager 获取，暂时返回 false
	// TODO: 添加到应用服务
	return false
}
