// +build !windows

package system

// NewPrivilegeElevator 创建平台特定的权限提升器（非 Windows 平台）
// 构建标签: +build !windows
// 表示此文件仅在非 Windows 平台编译
func NewPrivilegeElevator() (PrivilegeElevator, error) {
	// 非 Windows 平台: 返回 UnixElevator
	return NewUnixElevator(), nil
}
