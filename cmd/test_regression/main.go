// 回归测试工具
// 用途: 验证现有功能在重构后仍然正常工作
// 编译: go build -o test_regression ./cmd/test_regression

package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/chen/wails3-hosts/internal/infrastructure/persistence"
	"github.com/chen/wails3-hosts/internal/infrastructure/system"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("Unix 平台回归测试")
	fmt.Println("========================================")
	fmt.Printf("操作系统: %s\n", runtime.GOOS)
	fmt.Printf("架构: %s\n\n", runtime.GOARCH)

	testResults := []TestResult{}

	// 测试 1: 创建权限提升器
	fmt.Println("[测试 1] 创建权限提升器...")
	elevator, err := system.NewPrivilegeElevator()
	if err != nil {
		fmt.Printf("❌ 创建提升器失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 权限提升器创建成功\n")
	testResults = append(testResults, TestResult{Name: "创建权限提升器", Passed: true})

	// 测试 2: 检查凭据缓存能力
	fmt.Println("\n[测试 2] 检查凭据缓存能力...")
	canCache := elevator.CanCacheCredentials()
	if canCache {
		fmt.Printf("✅ 支持凭据缓存 (符合 Unix 平台特性)\n")
		testResults = append(testResults, TestResult{Name: "凭据缓存检查", Passed: true})
	} else {
		fmt.Printf("⚠️  不支持凭据缓存 (Windows 特性，Unix 应该支持)\n")
		testResults = append(testResults, TestResult{Name: "凭据缓存检查", Passed: false})
	}

	// 测试 3: 创建 HostsFileOperator
	fmt.Println("\n[测试 3] 创建 HostsFileOperator...")
	operator, err := system.NewHostsFileOperator(elevator)
	if err != nil {
		fmt.Printf("❌ 创建操作器失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ HostsFileOperator 创建成功\n")
	testResults = append(testResults, TestResult{Name: "创建操作器", Passed: true})

	// 测试 4: 读取 hosts 文件
	fmt.Println("\n[测试 4] 读取 hosts 文件...")
	content, err := operator.ReadCurrent()
	if err != nil {
		fmt.Printf("❌ 读取 hosts 文件失败: %v\n", err)
		testResults = append(testResults, TestResult{Name: "读取hosts", Passed: false})
	} else {
		fmt.Printf("✅ 成功读取 hosts 文件\n")
		fmt.Printf("   内容长度: %d 字节\n", len(content))
		fmt.Printf("   前 100 字符:\n")
		if len(content) > 100 {
			fmt.Printf("   %s...\n", content[:100])
		} else {
			fmt.Printf("   %s\n", content)
		}
		testResults = append(testResults, TestResult{Name: "读取hosts", Passed: true})
	}

	// 测试 5: 检查配置目录
	fmt.Println("\n[测试 5] 检查配置目录...")
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("❌ 获取配置目录失败: %v\n", err)
		testResults = append(testResults, TestResult{Name: "配置目录", Passed: false})
	} else {
		appConfigDir := fmt.Sprintf("%s/hosts-manager", configDir)
		if _, err := os.Stat(appConfigDir); os.IsNotExist(err) {
			fmt.Printf("ℹ️  配置目录不存在（首次运行正常）: %s\n", appConfigDir)
		} else {
			fmt.Printf("✅ 配置目录存在: %s\n", appConfigDir)
		}
		testResults = append(testResults, TestResult{Name: "配置目录", Passed: true})
	}

	// 测试 6: 初始化存储（不实际创建文件）
	fmt.Println("\n[测试 6] 测试 JSON 存储...")
	_, err = persistence.NewJSONStorage(configDir)
	if err != nil {
		fmt.Printf("❌ 创建存储失败: %v\n", err)
		testResults = append(testResults, TestResult{Name: "JSON存储", Passed: false})
	} else {
		fmt.Printf("✅ JSON 存储初始化成功\n")
		testResults = append(testResults, TestResult{Name: "JSON存储", Passed: true})
	}

	// 测试 7: 测试 SudoManager
	fmt.Println("\n[测试 7] 测试 SudoManager...")
	sudoManager := system.NewSudoManager()
	if sudoManager == nil {
		fmt.Printf("❌ 创建 SudoManager 失败\n")
		testResults = append(testResults, TestResult{Name: "SudoManager", Passed: false})
	} else {
		fmt.Printf("✅ SudoManager 创建成功\n")
		isCached := sudoManager.IsPasswordCached()
		fmt.Printf("   密码缓存状态: %v\n", isCached)
		remaining := sudoManager.GetCacheRemaining()
		fmt.Printf("   缓存剩余时间: %d 秒\n", remaining)
		testResults = append(testResults, TestResult{Name: "SudoManager", Passed: true})
	}

	// 测试总结
	fmt.Println("\n========================================")
	fmt.Println("测试总结")
	fmt.Println("========================================")

	passed := 0
	failed := 0
	for _, result := range testResults {
		if result.Passed {
			fmt.Printf("✅ [%s] 通过\n", result.Name)
			passed++
		} else {
			fmt.Printf("❌ [%s] 失败\n", result.Name)
			failed++
		}
	}

	fmt.Printf("\n总计: %d 个测试\n", len(testResults))
	fmt.Printf("✅ 通过: %d 个\n", passed)
	fmt.Printf("❌ 失败: %d 个\n", failed)

	if failed == 0 {
		fmt.Println("\n🎉 所有测试通过！现有功能正常，无回归。")
		os.Exit(0)
	} else {
		fmt.Println("\n⚠️  存在失败的测试，请检查。")
		os.Exit(1)
	}
}

// TestResult 测试结果
type TestResult struct {
	Name   string
	Passed bool
}
