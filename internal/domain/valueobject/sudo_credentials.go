package valueobject

import (
	"time"
)

// SudoCredentials sudo 凭证值对象
// 单一职责: 表示 sudo 密码及其缓存状态
// DDD: 值对象是不可变的，这里使用不可变结构
// 安全考虑: 密码仅在内存中保存，不持久化
type SudoCredentials struct {
	Password    string
	IsCached    bool
	ExpiresAt   time.Time
	CacheDuration time.Duration
}

// NewSudoCredentials 创建新的 sudo 凭证
// KISS: 简单的工厂函数，默认缓存 5 分钟
func NewSudoCredentials(password string) *SudoCredentials {
	cacheDuration := 5 * time.Minute
	return &SudoCredentials{
		Password:      password,
		IsCached:      false, // 初始不缓存
		ExpiresAt:     time.Time{},
		CacheDuration: cacheDuration,
	}
}

// NewCachedSudoCredentials 创建带缓存的 sudo 凭证
func NewCachedSudoCredentials(password string, cacheDuration time.Duration) *SudoCredentials {
	return &SudoCredentials{
		Password:      password,
		IsCached:      true,
		ExpiresAt:     time.Now().Add(cacheDuration),
		CacheDuration: cacheDuration,
	}
}

// IsExpired 检查凭证是否过期
func (c *SudoCredentials) IsExpired() bool {
	if !c.IsCached {
		return true
	}
	return time.Now().After(c.ExpiresAt)
}

// IsValid 检查凭证是否有效
func (c *SudoCredentials) IsValid() bool {
	return c.Password != "" && !c.IsExpired()
}

// GetPassword 获取密码
func (c *SudoCredentials) GetPassword() string {
	return c.Password
}

// Clear 清除密码（用于安全清理）
// YAGNI: 目前仅清空字符串，未来可添加安全擦除
func (c *SudoCredentials) Clear() {
	c.Password = ""
	c.IsCached = false
	c.ExpiresAt = time.Time{}
}
