// +build !windows

package system

import (
	"runtime"
)

// UnixElevator Unix/Linux/macOS 权限提升器
// 单一职责: 使用 sudo 命令进行权限提升
//
// 设计原则应用:
// - S: 仅负责 Unix 平台的权限提升
// - D: 依赖 SudoCommand 具体实现
//
// 构建标签: +build !windows
// 表示此文件仅在非 Windows 平台编译（darwin, linux 等）
type UnixElevator struct {
	sudoManager *SudoManager
}

// NewUnixElevator 创建 Unix 提升器实例
// 工厂模式: 封装创建逻辑
func NewUnixElevator() *UnixElevator {
	return &UnixElevator{
		sudoManager: NewSudoManager(),
	}
}

// Validate 验证 sudo 密码
// 实现: 使用 sudo -v 验证密码，验证成功后系统会缓存凭证
func (e *UnixElevator) Validate(password string) bool {
	return e.sudoManager.ValidatePassword(password)
}

// Execute 执行需要 sudo 权限的写入操作
// 实现: 使用系统缓存的 sudo 凭证执行写入
//
// 注意: 调用此方法前应先调用 Validate 验证密码
// 或者密码已经被缓存（在 5 分钟有效期内）
func (e *UnixElevator) Execute(content string) error {
	// 使用 SudoCommand 执行写入
	// 注意: 这里不传递密码，使用系统缓存的 sudo 凭证
	cmd := NewSudoCommand([]string{"sh", "-c", "cat > /etc/hosts"})
	cmd.SetStdin([]byte(content))

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// CanCacheCredentials Unix 平台支持 sudo 凭证缓存
// 实现: 返回 true，表示可以缓存凭证
//
// 注意: sudo 系统默认缓存 5 分钟，可由系统配置调整
func (e *UnixElevator) CanCacheCredentials() bool {
	return true
}

// GetOS 获取操作系统名称
func (e *UnixElevator) GetOS() string {
	return runtime.GOOS
}

// GetArch 获取系统架构
func (e *UnixElevator) GetArch() string {
	return runtime.GOARCH
}

// NeedsSudo Unix 平台需要 sudo 密码验证
func (e *UnixElevator) NeedsSudo() bool {
	return true
}
