package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chen/wails3-hosts/internal/application/dto"
	"github.com/chen/wails3-hosts/internal/application/service"
	"github.com/chen/wails3-hosts/internal/infrastructure/persistence"
	"github.com/chen/wails3-hosts/internal/infrastructure/system"
)

func main() {
	// 获取配置目录
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("获取配置目录失败: %v\n", err)
		return
	}
	configDir = filepath.Join(configDir, "hosts-manager")

	// 初始化存储
	storage, err := persistence.NewJSONStorage(configDir)
	if err != nil {
		fmt.Printf("创建存储失败: %v\n", err)
		return
	}

	// 初始化仓储
	hostsRepo := persistence.NewHostsRepository(storage)
	versionRepo := persistence.NewVersionRepository(storage)

	// 初始化系统操作
	hostsFileOp, err := system.NewHostsFileOperator()
	if err != nil {
		fmt.Printf("创建 hosts 文件操作器失败: %v\n", err)
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

	fmt.Println("========== 测试创建分组 ==========")
	// 测试创建分组
	group, err := appService.CreateGroup(ctx, dto.CreateHostsGroupRequest{
		Name:        "测试分组",
		Description: "这是一个测试分组",
	})
	if err != nil {
		fmt.Printf("❌ 创建分组失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 创建分组成功: %s (%s)\n", group.Name, group.ID)

	fmt.Println("\n========== 测试添加条目 ==========")
	// 测试添加条目
	err = appService.AddEntry(ctx, dto.AddEntryRequest{
		GroupID:  group.ID,
		IP:       "127.0.0.1",
		Hostname: "test.local",
		Comment:  "测试条目",
	})
	if err != nil {
		fmt.Printf("❌ 添加条目失败: %v\n", err)
		return
	}
	fmt.Println("✅ 添加条目成功")

	fmt.Println("\n========== 测试获取所有分组 ==========")
	// 测试获取所有分组
	groups, err := appService.GetAllGroups(ctx)
	if err != nil {
		fmt.Printf("❌ 获取分组失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 获取分组成功，共 %d 个分组:\n", len(groups))
	for _, g := range groups {
		fmt.Printf("  - %s: %d 个条目\n", g.Name, len(g.Entries))
	}

	fmt.Println("\n========== 测试批量更新条目 ==========")
	// 测试批量更新
	err = appService.BatchUpdateEntries(ctx, dto.BatchUpdateEntriesRequest{
		GroupID: group.ID,
		Entries: []dto.BatchUpdateEntryRequest{
			{IP: "192.168.1.1", Hostname: "dev1.local", Comment: "开发服务器1", Enabled: true},
			{IP: "192.168.1.2", Hostname: "dev2.local", Comment: "开发服务器2", Enabled: true},
			{IP: "127.0.0.1", Hostname: "localhost", Comment: "本地主机", Enabled: true},
		},
	})
	if err != nil {
		fmt.Printf("❌ 批量更新失败: %v\n", err)
		return
	}
	fmt.Println("✅ 批量更新成功")

	// 再次获取分组查看结果
	updatedGroup, _ := appService.GetGroupByID(ctx, group.ID)
	fmt.Printf("  更新后分组 '%s' 现在有 %d 个条目\n", updatedGroup.Name, len(updatedGroup.Entries))

	fmt.Println("\n========== 测试生成预览 ==========")
	// 测试生成预览
	preview, err := appService.GeneratePreview(ctx)
	if err != nil {
		fmt.Printf("❌ 生成预览失败: %v\n", err)
		return
	}
	fmt.Println("✅ 生成预览成功:")
	fmt.Println(preview)

	fmt.Println("\n========== 测试切换分组状态 ==========")
	err = appService.ToggleGroup(ctx, dto.ToggleGroupRequest{
		ID:      group.ID,
		Enabled: true,
	})
	if err != nil {
		fmt.Printf("❌ 切换分组状态失败: %v\n", err)
		return
	}
	fmt.Println("✅ 切换分组状态成功")

	fmt.Println("\n========== 所有测试通过！ ==========")
	fmt.Printf("\n配置文件路径: %s\n", filepath.Join(configDir, "config.json"))
}
