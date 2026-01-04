package entity

import (
	"time"

	"github.com/google/uuid"
)

// HostsVersion 表示 hosts 文件的一个历史版本
// 单一职责: 记录 hosts 文件的某个历史状态
type HostsVersion struct {
	ID          string    `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	Content     string    `json:"content"`
	Description string    `json:"description"`
	Source      string    `json:"source"` // "manual" | "auto" | "rollback"
}

// VersionSource 版本来源类型
// 常量定义避免魔法字符串
const (
	SourceManual   = "manual"   // 手动应用
	SourceAuto     = "auto"     // 自动创建（启动时检测变化）
	SourceRollback = "rollback" // 回滚操作
)

// NewHostsVersion 创建一个新的版本记录
// KISS: 简单的工厂函数
func NewHostsVersion(content, description, source string) *HostsVersion {
	return &HostsVersion{
		ID:          uuid.New().String(),
		Timestamp:   time.Now(),
		Content:     content,
		Description: description,
		Source:      source,
	}
}

// IsExpired 检查版本是否过期（超过指定天数）
// YAGNI: 仅实现简单的过期检查
func (v *HostsVersion) IsExpired(days int) bool {
	expirationDate := v.Timestamp.AddDate(0, 0, days)
	return time.Now().After(expirationDate)
}

// GetAge 获取版本年龄（天数）
func (v *HostsVersion) GetAge() int {
	return int(time.Since(v.Timestamp).Hours() / 24)
}
