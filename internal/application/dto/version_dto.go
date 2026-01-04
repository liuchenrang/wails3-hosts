package dto

// HostsVersionDTO 版本历史数据传输对象
type HostsVersionDTO struct {
	ID          string `json:"id"`
	Timestamp   string `json:"timestamp"`
	Content     string `json:"content"`
	Description string `json:"description"`
	Source      string `json:"source"`
	Age         int    `json:"age"` // 版本年龄（天数）
}

// RollbackRequest 回滚请求
type RollbackRequest struct {
	VersionID string `json:"version_id"`
	SudoPassword string `json:"sudo_password"` // 需要sudo权限
}

// ApplyHostsRequest 应用 hosts 配置请求
type ApplyHostsRequest struct {
	SudoPassword string `json:"sudo_password"` // 可选，如果已有缓存则不需要
}

// ValidateSudoRequest 验证 sudo 密码请求
type ValidateSudoRequest struct {
	Password string `json:"password"`
}

// ValidateSudoResponse 验证 sudo 密码响应
type ValidateSudoResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}
