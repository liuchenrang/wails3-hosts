package entity

import (
	"net"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// HostsEntry 表示一个 hosts 文件条目
// 单一职责: 表示 IP 到主机名的映射
type HostsEntry struct {
	ID       string `json:"id"`
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	Comment  string `json:"comment"`
	Enabled  bool   `json:"enabled"`
}

// NewHostsEntry 创建一个新的 hosts 条目
// KISS: 简单的工厂函数
func NewHostsEntry(ip, hostname, comment string) *HostsEntry {
	return &HostsEntry{
		ID:       uuid.New().String(),
		IP:       strings.TrimSpace(ip),
		Hostname: strings.TrimSpace(hostname),
		Comment:  strings.TrimSpace(comment),
		Enabled:  true, // 默认启用
	}
}

// Validate 验证 hosts 条目的有效性
// 单一职责: 仅负责验证条目数据
func (e *HostsEntry) Validate() error {
	// 验证 IP 地址
	if e.IP == "" {
		return ErrInvalidIP
	}
	if net.ParseIP(e.IP) == nil {
		return ErrInvalidIP
	}

	// 验证主机名
	if e.Hostname == "" {
		return ErrInvalidHostname
	}
	if !isValidHostname(e.Hostname) {
		return ErrInvalidHostname
	}

	return nil
}

// ToHostsLine 将条目转换为 hosts 文件行格式
// DRY: 统一的格式化逻辑
func (e *HostsEntry) ToHostsLine() string {
	if !e.Enabled {
		return ""
	}

	line := e.IP + "\t" + e.Hostname
	if e.Comment != "" {
		line += "\t# " + e.Comment
	}
	return line
}

// isValidHostname 验证主机名格式
// YAGNI: 仅实现基本验证，不处理所有 RFC 标准
func isValidHostname(hostname string) bool {
	// 基本验证: 允许字母、数字、连字符和点
	// 主机名不能以连字符开头或结尾
	// 各部分长度不超过 63 字符
	pattern := `^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`
	matched, _ := regexp.MatchString(pattern, hostname)
	return matched && len(hostname) <= 253
}

// 错误定义
var (
	ErrInvalidIP       = &DomainError{Code: "INVALID_IP", Message: "IP 地址格式无效"}
	ErrInvalidHostname = &DomainError{Code: "INVALID_HOSTNAME", Message: "主机名格式无效"}
)
