package system

import (
	"strings"
	"sync"
	"time"
)

// SudoManager sudo 密码管理器
// 单一职责: 管理 sudo 密码的缓存和验证
// 安全考虑: 密码仅在内存中保存，不持久化到磁盘
type SudoManager struct {
	cachedPassword string
	expiresAt      time.Time
	cacheDuration  time.Duration
	mu             sync.RWMutex
}

// NewSudoManager 创建 sudo 管理器实例
func NewSudoManager() *SudoManager {
	return &SudoManager{
		cacheDuration: 5 * time.Minute, // 默认缓存 5 分钟
	}
}

// ValidatePassword 验证 sudo 密码是否有效并提升权限
// 通过执行 sudo 命令来验证，验证成功后系统会缓存 sudo 凭证
func (m *SudoManager) ValidatePassword(password string) bool {
	// 使用 sudo -v 验证密码
	// sudo 会将凭证缓存在系统中，后续的 sudo 命令不需要密码
	cmd := NewSudoCommand([]string{"-v"})
	cmd.SetPassword(password)

	if err := cmd.Run(); err != nil {
		return false
	}

	// 验证成功，缓存密码到内存（用于状态检查）
	m.CachePassword(password)

	return true
}

// CachePassword 缓存 sudo 密码
func (m *SudoManager) CachePassword(password string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cachedPassword = password
	m.expiresAt = time.Now().Add(m.cacheDuration)
}

// GetCachedPassword 获取缓存的密码
func (m *SudoManager) GetCachedPassword() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.cachedPassword
}

// IsPasswordCached 检查是否有有效的缓存密码
func (m *SudoManager) IsPasswordCached() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.cachedPassword == "" {
		return false
	}

	return time.Now().Before(m.expiresAt)
}

// ClearCache 清除缓存的密码
// 安全考虑: 提供手动清除缓存的方法
func (m *SudoManager) ClearCache() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 安全擦除密码（简单覆写）
	if m.cachedPassword != "" {
		m.cachedPassword = strings.Repeat("*", len(m.cachedPassword))
	}
	m.cachedPassword = ""
	m.expiresAt = time.Time{}
}

// SetCacheDuration 设置缓存时长
func (m *SudoManager) SetCacheDuration(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cacheDuration = duration
}

// GetCacheRemaining 获取缓存剩余时间（秒）
func (m *SudoManager) GetCacheRemaining() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.cachedPassword == "" {
		return 0
	}

	remaining := int(time.Until(m.expiresAt).Seconds())
	if remaining < 0 {
		return 0
	}
	return remaining
}
