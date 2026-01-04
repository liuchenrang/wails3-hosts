package dto

// HostsGroupDTO hosts 分组数据传输对象
// 单一职责: 跨层数据传输，避免直接传递领域实体
// DDD: DTO 用于隔离领域模型与外部接口
type HostsGroupDTO struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	IsEnabled   bool           `json:"is_enabled"`
	Entries     []HostsEntryDTO `json:"entries"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

// HostsEntryDTO hosts 条目数据传输对象
type HostsEntryDTO struct {
	ID       string `json:"id"`
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	Comment  string `json:"comment"`
	Enabled  bool   `json:"enabled"`
}

// CreateHostsGroupRequest 创建分组请求
type CreateHostsGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateHostsGroupRequest 更新分组请求
type UpdateHostsGroupRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ToggleGroupRequest 切换分组状态请求
type ToggleGroupRequest struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
}

// AddEntryRequest 添加条目请求
type AddEntryRequest struct {
	GroupID  string `json:"group_id"`
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	Comment  string `json:"comment"`
}

// UpdateEntryRequest 更新条目请求
type UpdateEntryRequest struct {
	GroupID  string `json:"group_id"`
	EntryID  string `json:"entry_id"`
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	Comment  string `json:"comment"`
}

// DeleteEntryRequest 删除条目请求
type DeleteEntryRequest struct {
	GroupID string `json:"group_id"`
	EntryID string `json:"entry_id"`
}

// BatchUpdateEntriesRequest 批量更新条目请求
type BatchUpdateEntriesRequest struct {
	GroupID string                    `json:"group_id"`
	Entries []BatchUpdateEntryRequest `json:"entries"`
}

// BatchUpdateEntryRequest 批量更新条目单项
type BatchUpdateEntryRequest struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	Comment  string `json:"comment"`
	Enabled  bool   `json:"enabled"`
}
