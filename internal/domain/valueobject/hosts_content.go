package valueobject

import (
	"strings"

	"github.com/chen/wails3-hosts/internal/domain/entity"
)

// HostsContent hosts 文件内容值对象
// 单一职责: 表示不可变的 hosts 文件内容
// DDD: 值对象通过其属性值来标识，没有唯一标识符
type HostsContent struct {
	RawContent    string
	ParsedEntries []entity.HostsEntry
}

// NewHostsContent 创建 hosts 内容值对象
// KISS: 简单的工厂函数
func NewHostsContent(rawContent string) *HostsContent {
	entries := parseHostsContent(rawContent)
	return &HostsContent{
		RawContent:    rawContent,
		ParsedEntries: entries,
	}
}

// parseHostsContent 解析 hosts 文件内容
// 单一职责: 解析 hosts 文件格式
func parseHostsContent(content string) []entity.HostsEntry {
	entries := make([]entity.HostsEntry, 0)
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			// 跳过空行和注释行
			continue
		}

		// 解析 IP 和主机名
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		ip := parts[0]
		hostname := parts[1]
		comment := ""

		// 提取注释
		if idx := strings.Index(line, "#"); idx > 0 {
			comment = strings.TrimSpace(line[idx+1:])
		}

		entry := entity.NewHostsEntry(ip, hostname, comment)
		entries = append(entries, *entry)
	}

	return entries
}

// GetEntryCount 获取条目数量
func (c *HostsContent) GetEntryCount() int {
	return len(c.ParsedEntries)
}

// ContainsHostname 检查是否包含指定主机名
func (c *HostsContent) ContainsHostname(hostname string) bool {
	for _, entry := range c.ParsedEntries {
		if entry.Hostname == hostname {
			return true
		}
	}
	return false
}

// ToString 转换为字符串格式
func (c *HostsContent) ToString() string {
	return c.RawContent
}
