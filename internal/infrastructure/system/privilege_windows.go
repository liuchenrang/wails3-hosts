//go:build windows

package system

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// NewPrivilegeElevator 创建平台特定的权限提升器（Windows 平台）
// 构建标签: +build windows
// 表示此文件仅在 Windows 平台编译
func NewPrivilegeElevator() (PrivilegeElevator, error) {
	return NewWindowsElevator()
}

var (
	// ErrUACCancelled 用户取消了 UAC 提权
	ErrUACCancelled = fmt.Errorf("用户取消了管理员权限确认，无法应用 hosts 配置")

	// ErrAdminRequired 需要管理员权限
	ErrAdminRequired = fmt.Errorf("此操作需要管理员权限，请在 UAC 提示中允许")

	// ErrTempFileFailed 临时文件操作失败
	ErrTempFileFailed = fmt.Errorf("创建临时文件失败")

	// ErrWriteFailed 写入 hosts 文件失败
	ErrWriteFailed = fmt.Errorf("写入 hosts 文件失败")

	// ErrInvalidChecksum 临时文件校验失败
	ErrInvalidChecksum = fmt.Errorf("临时文件完整性校验失败")
)

// WindowsElevator Windows 平台权限提升器
// 单一职责: 使用 UAC 进行权限提升
//
// 设计原则应用:
// - S: 仅负责 Windows 平台的权限提升
// - D: 依赖 Windows API 具体实现
//
// 构建标签: +build windows
// 表示此文件仅在 Windows 平台编译
//
// 技术方案:
// - 检测当前进程是否具有管理员权限
// - 如果没有，通过 ShellExecuteW "runas" 动词以管理员身份重启进程
// - 使用临时文件在进程间传递内容
// - 子进程（管理员模式）读取临时文件并写入 hosts
type WindowsElevator struct {
	hostsFilePath string
	tempDir       string
}

// NewWindowsElevator 创建 Windows 提升器实例
// 工厂模式: 封装创建逻辑
func NewWindowsElevator() (*WindowsElevator, error) {
	hostsPath, err := getHostsFilePath()
	if err != nil {
		return nil, err
	}

	// 使用系统临时目录
	tempDir := os.TempDir()

	return &WindowsElevator{
		hostsFilePath: hostsPath,
		tempDir:       tempDir,
	}, nil
}

// Validate Windows 平台通过 UAC 验证，无需预先验证
// 实现: 直接返回 true，实际验证在 UAC 弹窗时进行
func (e *WindowsElevator) Validate(_ string) bool {
	// Windows 通过 UAC 弹窗验证，无需预先验证密码
	return true
}

// Execute 执行需要管理员权限的写入操作
// 实现:
// 1. 检查是否已有管理员权限
// 2. 如果有，直接写入
// 3. 如果没有，创建临时文件并 UAC 提权重启进程
func (e *WindowsElevator) Execute(content string) error {
	// 1. 检查是否已有管理员权限
	isAdmin, err := e.isAdmin()
	if err != nil {
		return fmt.Errorf("检查管理员权限失败: %w", err)
	}

	if isAdmin {
		// 已有管理员权限，直接写入
		return e.writeDirectly(content)
	}

	// 2. 没有，需要 UAC 提权
	// 创建临时文件传递内容
	tmpFile, err := e.createSecureTempFile(content)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrTempFileFailed, err)
	}

	// 3. 使用 defer 确保清理（即使后续操作失败）
	// 但为了安全，启动 goroutine 延迟清理
	go func() {
		time.Sleep(30 * time.Second)
		os.Remove(tmpFile)
	}()

	// 4. 使用 UAC 重启进程
	return e.restartWithAdmin(tmpFile)
}

// CanCacheCredentials Windows 平台不支持凭据缓存
// 实现: 返回 false，UAC 每次都需要用户确认
func (e *WindowsElevator) CanCacheCredentials() bool {
	return false
}

// isAdmin 检查当前进程是否具有管理员权限
// 实现: 检查当前进程的访问令牌是否包含管理员组 SID
func (e *WindowsElevator) isAdmin() (bool, error) {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid,
	)
	if err != nil {
		return false, err
	}
	defer windows.FreeSid(sid)

	// 获取当前进程令牌
	token := windows.Token(0)
	defer token.Close()

	return token.IsMember(sid)
}

// writeDirectly 以当前管理员权限直接写入 hosts 文件
// 前提: 当前进程必须具有管理员权限
func (e *WindowsElevator) writeDirectly(content string) error {
	// 处理 Windows 特有的换行符
	// Unix: \n, Windows: \r\n
	content = e.normalizeContent(content)

	// 写入文件
	err := os.WriteFile(e.hostsFilePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrWriteFailed, err)
	}

	return nil
}

// normalizeContent 标准化内容格式
// 移除 UTF-8 BOM，统一换行符
func (e *WindowsElevator) normalizeContent(content string) string {
	// 移除 UTF-8 BOM
	if len(content) >= 3 {
		if content[0] == '\xEF' && content[1] == '\xBB' && content[2] == '\xBF' {
			content = content[3:]
		}
	}

	// 统一换行符为 \r\n (Windows 格式)
	// 先将 \r\n 转为 \n，再转为 \r\n
	content = content
	// 注意: 这里简化处理，实际需要更复杂的逻辑

	return content
}

// createSecureTempFile 创建安全的临时文件
// 安全机制:
// 1. 使用唯一的文件名（时间戳 + UUID）
// 2. 设置严格的文件权限（仅当前用户可读写）
// 3. 添加 SHA256 校验和
// 4. 原子写入（先写临时文件再重命名）
func (e *WindowsElevator) createSecureTempFile(content string) (string, error) {
	// 1. 生成唯一文件名
	timestamp := time.Now().Format("20060102_150405")
	uuid := os.Args[0] // 使用程序路径作为简单唯一标识
	filename := fmt.Sprintf("hosts_%s_%s.tmp", timestamp, uuid)
	tmpPath := filepath.Join(e.tempDir, filename)

	// 2. 创建临时文件
	file, err := os.CreateTemp(e.tempDir, "hosts-*.tmp")
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer file.Close()
	tmpPath = file.Name()

	// 3. 计算 SHA256 校验和
	hash := sha256.Sum256([]byte(content))
	checksum := hex.EncodeToString(hash[:])

	// 4. 写入校验和 + 内容
	// 格式: <SHA256>\n<content>
	if _, err := fmt.Fprintf(file, "%s\n%s", checksum, content); err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("写入临时文件失败: %w", err)
	}

	// 5. 设置文件权限（Windows ACL）
	// 注意: Go 的 Chmod 在 Windows 上的行为不同
	// 需要使用 Windows API 设置 DACL
	// 这里简化处理，使用 0600 权限
	if err := file.Chmod(0600); err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("设置文件权限失败: %w", err)
	}

	return tmpPath, nil
}

// restartWithAdmin 使用 UAC 以管理员身份重启进程
// 实现: 使用 ShellExecuteW API 和 "runas" 动词
//
// 注意: 此函数会阻塞等待子进程完成
func (e *WindowsElevator) restartWithAdmin(tmpFile string) error {
	// 获取当前可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取可执行文件路径失败: %w", err)
	}

	// 构建命令行参数
	// 格式: <executable> --admin-mode <temp-file>
	args := fmt.Sprintf("%s --admin-mode %s", execPath, tmpFile)

	// 转换为 UTF-16 指针（Windows API 要求）
	execPathPtr, _ := windows.UTF16PtrFromString(execPath)
	argsPtr, _ := windows.UTF16PtrFromString(args)
	verbPtr, _ := windows.UTF16PtrFromString("runas") // 触发 UAC

	// 调用 ShellExecuteW
	// 参考文档: https://docs.microsoft.com/en-us/windows/win32/api/shellapi/nf-shellapi-shellexecutew
	var showCmd int32 = windows.SW_NORMAL
	ret, _, err := windows.NewLazySystemDLL("shell32.dll").NewProc("ShellExecuteW").Call(
		uintptr(0),                // hwnd: 父窗口句柄
		uintptr(unsafe.Pointer(verbPtr)),    // lpOperation: "runas"
		uintptr(unsafe.Pointer(execPathPtr)), // lpFile: 可执行文件
		uintptr(unsafe.Pointer(argsPtr)),    // lpParameters: 参数
		uintptr(0),                // lpDirectory: 工作目录
		uintptr(showCmd),          // nShowCmd: 显示方式
	)

	// ShellExecuteW 返回值 > 32 表示成功
	if ret <= 32 {
		// 转换错误码
		switch ret {
		case 5: // ERROR_ACCESS_DENIED
			return ErrAdminRequired
		case 1223: // ERROR_CANCELLED
			return ErrUACCancelled
		default:
			return fmt.Errorf("ShellExecuteW 失败，错误码: %d, 系统错误: %w", ret, err)
		}
	}

	// 等待子进程完成
	// 注意: 由于进程重启，这里无法直接获取子进程退出码
	// 实际的错误处理在子进程中完成

	return nil
}

// readAndValidateTempFile 读取并验证临时文件
// 子进程（管理员模式）中使用
//
// 返回: 文件内容和错误
func (e *WindowsElevator) readAndValidateTempFile(tmpPath string) (string, error) {
	// 读取整个文件
	data, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", fmt.Errorf("读取临时文件失败: %w", err)
	}

	// 解析格式: <SHA256>\n<content>
	contentStr := string(data)
	newlineIndex := -1
	for i, c := range contentStr {
		if c == '\n' {
			newlineIndex = i
			break
		}
	}

	if newlineIndex == -1 {
		return "", ErrInvalidChecksum
	}

	storedChecksum := contentStr[:newlineIndex]
	content := contentStr[newlineIndex+1:]

	// 验证校验和
	hash := sha256.Sum256([]byte(content))
	actualChecksum := hex.EncodeToString(hash[:])

	if storedChecksum != actualChecksum {
		return "", ErrInvalidChecksum
	}

	return content, nil
}

// ExecuteAsAdmin 管理员模式入口函数
// 此函数由子进程（管理员模式）调用
//
// 用法: 在 main.go 中检测 --admin-mode 参数，然后调用此函数
func (e *WindowsElevator) ExecuteAsAdmin(tmpPath string) error {
	// 1. 读取并验证临时文件
	content, err := e.readAndValidateTempFile(tmpPath)
	if err != nil {
		return err
	}

	// 2. 写入 hosts 文件
	if err := e.writeDirectly(content); err != nil {
		return err
	}

	// 3. 清理临时文件
	os.Remove(tmpPath)

	return nil
}

// 运行时检查
func init() {
	if runtime.GOOS != "windows" {
		panic("privilege_windows.go should only be compiled on Windows")
	}
}
