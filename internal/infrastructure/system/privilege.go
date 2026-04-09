package system

// PrivilegeElevator 权限提升器接口
// 单一职责: 定义平台无关的权限提升操作
//
// 设计原则应用:
// - S: 接口职责单一，仅定义权限提升操作
// - O: 通过接口扩展，支持不同平台的实现
// - L: 所有平台实现可互换使用
// - I: 接口最小化，仅包含必要方法
// - D: HostsFileOperator 依赖此抽象接口，而非具体实现
//
// 使用场景:
// - Unix/Linux/macOS: 使用 sudo 提权
// - Windows: 使用 UAC 提权
type PrivilegeElevator interface {
	// Validate 验证凭据是否有效
	// 参数:
	//   - credentials: 凭据内容（Unix 为密码，Windows 通常不需要）
	// 返回:
	//   - bool: 验证是否成功
	//
	// 平台差异:
	// - Unix: 验证 sudo 密码，成功后系统会缓存凭证
	// - Windows: 通常返回 true，通过 UAC 弹窗验证
	Validate(credentials string) bool

	// Execute 执行需要提升权限的操作
	// 参数:
	//   - content: 要写入 hosts 文件的内容
	// 返回:
	//   - error: 操作失败时的错误信息
	//
	// 平台差异:
	// - Unix: 使用 sudo 执行写入，利用系统缓存的凭证
	// - Windows: 通过 UAC 重启进程并提权执行
	//
	// 错误处理:
	// - 权限不足: 返回明确的错误提示
	// - 超时: 操作超过预期时间应中止
	// - 文件占用: hosts 文件被其他程序占用
	Execute(content string) error

	// CanCacheCredentials 是否可以缓存凭据
	// 返回:
	//   - bool: 是否支持凭据缓存
	//
	// 平台差异:
	// - Unix: true (sudo 可以缓存 5 分钟)
	// - Windows: false (UAC 每次需要用户确认)
	//
	// 用途:
	// - 前端根据此值决定是否显示"密码已缓存"提示
	// - Unix 平台可以避免重复输入密码
	// - Windows 平台每次都会弹出 UAC 提示
	CanCacheCredentials() bool

	// GetOS 获取操作系统名称
	// 返回:
	//   - string: 操作系统标识 ("windows", "darwin", "linux")
	GetOS() string

	// GetArch 获取系统架构
	// 返回:
	//   - string: 架构标识 ("amd64", "arm64", etc.)
	GetArch() string

	// NeedsSudo 是否需要 sudo 密码验证
	// 返回:
	//   - bool: 是否需要密码验证
	//
	// 平台差异:
	// - Unix: true (需要 sudo 密码)
	// - Windows: false (使用 UAC，不需要密码)
	NeedsSudo() bool
}
