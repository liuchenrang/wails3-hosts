// UAC 提权测试工具
// 用途: 测试 Windows 平台的 UAC 提权流程
// 编译: go build -o test_uac.exe ./cmd/test_uac
//
// 使用方法:
//   1. 以标准用户身份运行: test_uac.exe
//   2. 观察是否弹出 UAC 提示
//   3. 点击"是"后应以管理员权限运行
//   4. 使用管理员模式: test_uac.exe --admin-mode

package main

import (
	"fmt"
	"os"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	// AdminModeFlag 管理员模式标志
	AdminModeFlag = "--admin-mode"
)

var (
	// ErrUACCancelled UAC 取消错误
	ErrUACCancelled = fmt.Errorf("用户取消了 UAC 提权")
	// ErrAdminRequired 需要管理员权限
	ErrAdminRequired = fmt.Errorf("需要管理员权限")
)

func main() {
	fmt.Println("========================================")
	fmt.Println("Windows UAC 提权测试工具")
	fmt.Println("========================================")
	fmt.Printf("操作系统: %s\n", runtime.GOOS)
	fmt.Printf("架构: %s\n", runtime.GOARCH)

	// 检查命令行参数
	if len(os.Args) > 1 && os.Args[1] == AdminModeFlag {
		// 管理员模式
		runAsAdminMode()
	} else {
		// 标准模式
		runNormalMode()
	}
}

// runNormalMode 标准模式运行
// 检测是否有管理员权限，如果没有则请求 UAC 提权
func runNormalMode() {
	fmt.Println("\n[模式] 标准用户模式")

	// 1. 检查当前权限
	isAdmin, err := isAdmin()
	if err != nil {
		fmt.Printf("❌ 检查管理员权限失败: %v\n", err)
		os.Exit(1)
	}

	if isAdmin {
		fmt.Println("✅ 当前已具有管理员权限")
		fmt.Println("\n提示: 请以标准用户身份运行此程序来测试 UAC 提权")
		os.Exit(0)
	}

	fmt.Println("ℹ️  当前是标准用户权限")
	fmt.Println("\n准备请求 UAC 提权...")

	// 2. 请求 UAC 提权
	err = restartWithAdmin()
	if err != nil {
		fmt.Printf("\n❌ UAC 提权失败: %v\n", err)

		if err == ErrUACCancelled {
			fmt.Println("提示: 您取消了 UAC 提示，这是正常的用户行为")
		} else if err == ErrAdminRequired {
			fmt.Println("提示: 需要管理员权限才能执行此操作")
		}

		os.Exit(1)
	}

	// 注意: 如果提权成功，当前进程会被替换，不会执行到这里
	fmt.Println("✅ UAC 提权成功，程序将在管理员模式下重新启动")
}

// runAsAdminMode 管理员模式运行
// 执行需要管理员权限的操作
func runAsAdminMode() {
	fmt.Println("\n[模式] 管理员模式")
	fmt.Println("ℹ️  此程序以管理员权限运行")

	// 1. 验证管理员权限
	isAdmin, err := isAdmin()
	if err != nil {
		fmt.Printf("❌ 检查管理员权限失败: %v\n", err)
		os.Exit(1)
	}

	if !isAdmin {
		fmt.Println("❌ 错误: 管理员模式下应具有管理员权限")
		os.Exit(1)
	}

	fmt.Println("✅ 验证: 确认具有管理员权限")

	// 2. 执行测试操作
	fmt.Println("\n开始执行管理员权限操作...")

	// 读取系统 hosts 文件
	hostsPath := os.Getenv("SystemRoot") + "\\System32\\drivers\\etc\\hosts"
	fmt.Printf("读取 hosts 文件: %s\n", hostsPath)

	content, err := os.ReadFile(hostsPath)
	if err != nil {
		fmt.Printf("❌ 读取 hosts 文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ 成功读取 hosts 文件 (%d 字节)\n", len(content))
	fmt.Printf("前 100 个字符:\n%s\n", string(content[:min(100, len(content))]))

	// 3. 测试写入（不实际修改，仅测试权限）
	testFile := os.Getenv("TEMP") + "\\uac_test.tmp"
	fmt.Printf("\n测试写入临时文件: %s\n", testFile)

	err = os.WriteFile(testFile, []byte("UAC 提权测试\n"), 0644)
	if err != nil {
		fmt.Printf("❌ 写入测试文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ 成功写入测试文件")

	// 清理测试文件
	os.Remove(testFile)
	fmt.Println("✅ 清理测试文件")

	// 完成
	fmt.Println("\n========================================")
	fmt.Println("✅ 所有测试通过！")
	fmt.Println("========================================")
	fmt.Println("\n结论:")
	fmt.Println("✅ UAC 提权流程正常工作")
	fmt.Println("✅ 管理员权限正确提升")
	fmt.Println("✅ 可以执行需要管理员权限的操作")
}

// isAdmin 检查当前进程是否具有管理员权限
// 返回: bool - 是否是管理员, error - 错误信息
func isAdmin() (bool, error) {
	// 创建管理员组 SID
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

	// 检查是否是管理员组成员
	return token.IsMember(sid)
}

// restartWithAdmin 使用 UAC 以管理员身份重启当前进程
// 返回: error - 错误信息
func restartWithAdmin() error {
	// 获取当前可执行文件路径
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取可执行文件路径失败: %w", err)
	}

	// 构建命令行参数: <executable> --admin-mode
	args := fmt.Sprintf("%s %s", execPath, AdminModeFlag)

	// 转换为 UTF-16 指针（Windows API 要求）
	execPathPtr, err := windows.UTF16PtrFromString(execPath)
	if err != nil {
		return fmt.Errorf("转换路径失败: %w", err)
	}

	argsPtr, err := windows.UTF16PtrFromString(args)
	if err != nil {
		return fmt.Errorf("转换参数失败: %w", err)
	}

	verbPtr, err := windows.UTF16PtrFromString("runas") // 触发 UAC
	if err != nil {
		return fmt.Errorf("转换动词失败: %w", err)
	}

	fmt.Println("\n正在请求 UAC 提权...")
	fmt.Println("提示: 请在弹出的 UAC 窗口中点击【是】")

	// 调用 ShellExecuteW
	showCmd := int32(windows.SW_NORMAL)

	// 加载 shell32.dll
	shell32, err := windows.LoadLibrary("shell32.dll")
	if err != nil {
		return fmt.Errorf("加载 shell32.dll 失败: %w", err)
	}
	defer windows.FreeLibrary(shell32)

	// 获取 ShellExecuteW 函数
	proc, err := windows.GetProcAddress(shell32, "ShellExecuteW")
	if err != nil {
		return fmt.Errorf("获取 ShellExecuteW 函数失败: %w", err)
	}

	// 调用函数
	ret, _, _ := windows.SyscallN(
		uintptr(proc),
		8,
		uintptr(0),                         // hwnd: 父窗口句柄
		uintptr(unsafe.Pointer(verbPtr)),  // lpOperation: "runas"
		uintptr(unsafe.Pointer(execPathPtr)), // lpFile: 可执行文件
		uintptr(unsafe.Pointer(argsPtr)),  // lpParameters: 参数
		uintptr(0),                         // lpDirectory: 工作目录
		uintptr(showCmd),                   // nShowCmd: 显示方式
		uintptr(0),                         // 其他参数
		uintptr(0),                         // 其他参数
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
			return fmt.Errorf("ShellExecuteW 失败，错误码: %d", ret)
		}
	}

	return nil
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
