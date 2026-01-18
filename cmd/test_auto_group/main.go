package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chen/wails3-hosts/internal/application/service"
	"github.com/chen/wails3-hosts/internal/infrastructure/persistence"
	"github.com/chen/wails3-hosts/internal/infrastructure/system"
)

func main() {
	fmt.Println("========== 测试自动创建默认分组功能 ==========")

	// 获取配置目录
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("❌ 获取配置目录失败: %v\n", err)
		return
	}
	configDir = filepath.Join(configDir, "hosts-manager")

	// 初始化存储
	storage, err := persistence.NewJSONStorage(configDir)
	if err != nil {
		fmt.Printf("❌ 创建存储失败: %v\n", err)
		return
	}

	// 初始化仓储
	hostsRepo := persistence.NewHostsRepository(storage)
	versionRepo := persistence.NewVersionRepository(storage)

	// 初始化系统操作
	hostsFileOp, err := system.NewHostsFileOperator()
	if err != nil {
		fmt.Printf("❌ 创建 hosts 文件操作器失败: %v\n", err)
		return
	}
	sudoManager := system.NewSudoManager()

	// 创建应用服务
	appService := service.NewHostsApplicationService(
		hostsRepo,
		versionRepo,
		hostsFileOp,
		sudoManager,
	)

	ctx := context.Background()

	// 步骤1: 检查当前有多少个分组
	fmt.Println("\n步骤1: 检查当前分组状态")
	groups, err := appService.GetAllGroups(ctx)
	if err != nil {
		fmt.Printf("❌ 获取分组失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 当前有 %d 个分组\n", len(groups))

	for i, group := range groups {
		fmt.Printf("  %d. %s (启用: %v, 条目数: %d)\n", i+1, group.Name, group.IsEnabled, len(group.Entries))
	}

	// 步骤2: 如果有默认分组，查看其内容
	if len(groups) > 0 {
		for _, group := range groups {
			if group.Name == "默认分组" {
				fmt.Printf("\n✅ 发现默认分组！\n")
				fmt.Printf("  分组名称: %s\n", group.Name)
				fmt.Printf("  描述: %s\n", group.Description)
				fmt.Printf("  启用状态: %v\n", group.IsEnabled)
				fmt.Printf("  条目数量: %d\n", len(group.Entries))

				if len(group.Entries) > 0 {
					fmt.Printf("\n  前5个条目:\n")
					for i, entry := range group.Entries {
						if i >= 5 {
							break
						}
						fmt.Printf("    %d. %s\t%s\t%s\n", i+1, entry.IP, entry.Hostname, entry.Comment)
					}
				}
				break
			}
		}
	}

	fmt.Println("\n========== 测试完成！ ==========")
	fmt.Printf("\n配置文件路径: %s\n", filepath.Join(configDir, "config.json"))
}
