package entity

import (
	"time"

	"github.com/google/uuid"
)

// HostsGroup 表示一个 hosts 配置分组
// 单一职责: 管理一组相关的 hosts 条目
type HostsGroup struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	IsEnabled   bool        `json:"is_enabled"`
	Entries     []HostsEntry `json:"entries"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// NewHostsGroup 创建一个新的 hosts 分组
// KISS: 使用简单的工厂函数创建实体
func NewHostsGroup(name, description string) *HostsGroup {
	now := time.Now()
	return &HostsGroup{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		IsEnabled:   false, // 默认不启用，遵循最小惊讶原则
		Entries:     make([]HostsEntry, 0),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// AddEntry 添加一个新的 hosts 条目到分组
// 单一职责: 仅负责添加条目并更新时间戳
func (g *HostsGroup) AddEntry(entry HostsEntry) error {
	// 验证条目
	if err := entry.Validate(); err != nil {
		return err
	}

	g.Entries = append(g.Entries, entry)
	g.UpdatedAt = time.Now()
	return nil
}

// RemoveEntry 从分组中移除指定的条目
// DRY: 统一的条目查找和移除逻辑
func (g *HostsGroup) RemoveEntry(entryID string) bool {
	for i, entry := range g.Entries {
		if entry.ID == entryID {
			// 使用切片技巧删除元素
			g.Entries = append(g.Entries[:i], g.Entries[i+1:]...)
			g.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

// UpdateEntry 更新分组中的指定条目
func (g *HostsGroup) UpdateEntry(entryID string, newEntry HostsEntry) error {
	if err := newEntry.Validate(); err != nil {
		return err
	}

	for i, entry := range g.Entries {
		if entry.ID == entryID {
			g.Entries[i] = newEntry
			g.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrEntryNotFound
}

// ClearEntries 清空分组中的所有条目
// 用例: 批量更新前清空现有条目
func (g *HostsGroup) ClearEntries() {
	g.Entries = make([]HostsEntry, 0)
	g.UpdatedAt = time.Now()
}

// Toggle 切换分组的启用状态
func (g *HostsGroup) Toggle() {
	g.IsEnabled = !g.IsEnabled
	g.UpdatedAt = time.Now()
}

// GetEnabledEntries 获取所有启用的条目
func (g *HostsGroup) GetEnabledEntries() []HostsEntry {
	if !g.IsEnabled {
		return []HostsEntry{}
	}

	result := make([]HostsEntry, 0, len(g.Entries))
	for _, entry := range g.Entries {
		if entry.Enabled {
			result = append(result, entry)
		}
	}
	return result
}

// SetEnabled 设置分组的启用状态
func (g *HostsGroup) SetEnabled(enabled bool) {
	g.IsEnabled = enabled
	g.UpdatedAt = time.Now()
}

// 错误定义
var (
	ErrEntryNotFound = &DomainError{Code: "ENTRY_NOT_FOUND", Message: "条目不存在"}
	ErrInvalidName   = &DomainError{Code: "INVALID_NAME", Message: "分组名称无效"}
)

// DomainError 领域错误类型
type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}
