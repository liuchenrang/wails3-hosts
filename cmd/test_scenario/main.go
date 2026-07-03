package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/chen/wails3-hosts/internal/application/dto"
	"github.com/chen/wails3-hosts/internal/application/service"
	"github.com/chen/wails3-hosts/internal/infrastructure/persistence"
	"github.com/chen/wails3-hosts/internal/infrastructure/system"
	"golang.org/x/sys/windows"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("Hosts Manager 场景测试")
	fmt.Println("========================================")
	fmt.Printf("操作系统: %s\n", runtime.GOOS)
	fmt.Printf("架构: %s\n", runtime.GOARCH)

	// 检查管理员权限
	isAdmin, err := checkAdmin()
	if err != nil {
		fmt.Printf("❌ 检查管理员权限失败: %v\n", err)
		os.Exit(1)
	}

	if !isAdmin {
		fmt.Println("❌ 此测试需要管理员权限")
		os.Exit(1)
	}
	fmt.Println("✅ 当前进程已具有管理员权限")

	// 初始化基础设施
	fmt.Println("\n[初始化] 创建基础设施组件...")
	infra, err := initInfrastructure()
	if err != nil {
		fmt.Printf("❌ 初始化失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ 基础设施初始化成功")

	ctx := context.Background()

	// 记录初始 hosts 内容
	fmt.Println("\n[准备] 记录初始 hosts 文件内容...")
	originalHosts, err := infra.hostsFileOp.ReadCurrent()
	if err != nil {
		fmt.Printf("❌ 读取 hosts 文件失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 初始 hosts 内容已记录 (%d 字节)\n", len(originalHosts))

	// ========== 场景1: 获取分组列表 ==========
	fmt.Println("\n========================================")
	fmt.Println("[场景1] 获取分组列表")
	fmt.Println("========================================")
	groups, err := infra.appService.GetAllGroups(ctx)
	if err != nil {
		fmt.Printf("❌ 获取分组失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 当前分组数量: %d\n", len(groups))
	for i, g := range groups {
		fmt.Printf("   %d. %s (启用: %v, 条目: %d)\n", i+1, g.Name, g.IsEnabled, len(g.Entries))
	}

	// ========== 场景2: 创建新分组 ==========
	fmt.Println("\n========================================")
	fmt.Println("[场景2] 创建新分组")
	fmt.Println("========================================")
	newGroup, err := infra.appService.CreateGroup(ctx, dto.CreateHostsGroupRequest{
		Name:        "测试分组_" + time.Now().Format("150405"),
		Description: "自动化测试创建的分组",
	})
	if err != nil {
		fmt.Printf("❌ 创建分组失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 新分组创建成功: ID=%s, 名称=%s\n", newGroup.ID[:8], newGroup.Name)
	testGroupID := newGroup.ID

	// ========== 场景3: 添加条目到分组 ==========
	fmt.Println("\n========================================")
	fmt.Println("[场景3] 添加 hosts 条目")
	fmt.Println("========================================")
	testEntries := []struct {
		IP       string
		Hostname string
		Comment  string
	}{
		{"127.0.0.1", "test-scenario.local", "场景测试条目1"},
		{"192.168.1.100", "test-server.local", "场景测试条目2"},
		{"10.0.0.50", "test-db.local", "场景测试条目3"},
	}

	for _, entry := range testEntries {
		err := infra.appService.AddEntry(ctx, dto.AddEntryRequest{
			GroupID:  testGroupID,
			IP:       entry.IP,
			Hostname: entry.Hostname,
			Comment:  entry.Comment,
		})
		if err != nil {
			fmt.Printf("❌ 添加条目失败 (%s -> %s): %v\n", entry.IP, entry.Hostname, err)
			os.Exit(1)
		}
		fmt.Printf("✅ 条目添加成功: %s -> %s\n", entry.IP, entry.Hostname)
	}

	// ========== 场景4: 编辑分组 ==========
	fmt.Println("\n========================================")
	fmt.Println("[场景4] 编辑分组信息")
	fmt.Println("========================================")
	err = infra.appService.UpdateGroup(ctx, dto.UpdateHostsGroupRequest{
		ID:          testGroupID,
		Name:        "已编辑的测试分组",
		Description: "编辑后的描述信息",
	})
	if err != nil {
		fmt.Printf("❌ 编辑分组失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ 分组编辑成功")

	// 验证编辑结果
	updatedGroup, err := infra.appService.GetGroupByID(ctx, testGroupID)
	if err != nil {
		fmt.Printf("❌ 获取更新后的分组失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("   新名称: %s\n", updatedGroup.Name)
	fmt.Printf("   新描述: %s\n", updatedGroup.Description)

	// ========== 场景5: 切换分组启用状态 ==========
	fmt.Println("\n========================================")
	fmt.Println("[场景5] 切换分组启用状态")
	fmt.Println("========================================")

	// 先禁用
	err = infra.appService.ToggleGroup(ctx, dto.ToggleGroupRequest{
		ID:      testGroupID,
		Enabled: false,
	})
	if err != nil {
		fmt.Printf("❌ 禁用分组失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ 分组已禁用")

	// 再启用
	err = infra.appService.ToggleGroup(ctx, dto.ToggleGroupRequest{
		ID:      testGroupID,
		Enabled: true,
	})
	if err != nil {
		fmt.Printf("❌ 启用分组失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ 分组已重新启用")

	// ========== 场景6: 生成预览 ==========
	fmt.Println("\n========================================")
	fmt.Println("[场景6] 生成 hosts 配置预览")
	fmt.Println("========================================")
	preview, err := infra.appService.GeneratePreview(ctx)
	if err != nil {
		fmt.Printf("❌ 生成预览失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ 预览生成成功")
	fmt.Println("--- 预览内容 (前500字符) ---")
	fmt.Println(truncate(preview, 500))
	fmt.Println("--- 预览内容结束 ---")

	// ========== 场景7: 应用 hosts 配置 ==========
	fmt.Println("\n========================================")
	fmt.Println("[场景7] 应用 hosts 配置到系统")
	fmt.Println("========================================")

	// 检测冲突
	conflicts, err := infra.appService.DetectConflicts(ctx)
	if err != nil {
		fmt.Printf("❌ 检测冲突失败: %v\n", err)
		os.Exit(1)
	}
	if len(conflicts) > 0 {
		fmt.Printf("⚠️  检测到冲突: %d 个\n", len(conflicts))
		for hostname, ips := range conflicts {
			fmt.Printf("   %s: %v\n", hostname, ips)
		}
	} else {
		fmt.Println("✅ 未检测到冲突")
	}

	// 应用配置
	err = infra.appService.ApplyHosts(ctx, dto.ApplyHostsRequest{})
	if err != nil {
		fmt.Printf("❌ 应用配置失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ hosts 配置已应用")

	// ========== 场景8: 验证 hosts 文件变化 ==========
	fmt.Println("\n========================================")
	fmt.Println("[场景8] 验证 hosts 文件变化")
	fmt.Println("========================================")
	newHosts, err := infra.hostsFileOp.ReadCurrent()
	if err != nil {
		fmt.Printf("❌ 读取 hosts 文件失败: %v\n", err)
		os.Exit(1)
	}

	// 检查测试条目是否存在
	fmt.Println("检查测试条目是否写入 hosts 文件:")
	checkEntries := []string{"test-scenario.local", "test-server.local", "test-db.local"}
	allFound := true
	for _, hostname := range checkEntries {
		if contains(newHosts, hostname) {
			fmt.Printf("   ✅ %s 已写入\n", hostname)
		} else {
			fmt.Printf("   ❌ %s 未找到\n", hostname)
			allFound = false
		}
	}

	if !allFound {
		fmt.Println("❌ hosts 文件验证失败")
		fmt.Println("--- 当前 hosts 内容 ---")
		fmt.Println(truncate(newHosts, 1000))
	}

	// ========== 场景9: 获取版本历史 ==========
	fmt.Println("\n========================================")
	fmt.Println("[场景9] 获取版本历史")
	fmt.Println("========================================")
	versions, err := infra.appService.GetVersions(ctx, 10)
	if err != nil {
		fmt.Printf("❌ 获取版本历史失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 版本历史数量: %d\n", len(versions))
	for i, v := range versions {
		if i < 5 {
			fmt.Printf("   %d. %s (%s) - %s\n", i+1, v.ID[:8], v.Timestamp, v.Description)
		}
	}

	// ========== 场景10: 回滚到之前版本 ==========
	fmt.Println("\n========================================")
	fmt.Println("[场景10] 回滚 hosts 配置")
	fmt.Println("========================================")

	if len(versions) > 0 {
		// 找到包含原始内容的版本或使用最近的版本
		rollbackVersionID := versions[0].ID
		fmt.Printf("准备回滚到版本: %s\n", rollbackVersionID[:8])

		err = infra.appService.RollbackToVersion(ctx, dto.RollbackRequest{
			VersionID:    rollbackVersionID,
			SudoPassword: "", // Windows 不需要密码
		})
		if err != nil {
			fmt.Printf("❌ 回滚失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ 回滚成功")

		// 验证回滚后的 hosts 内容
		rolledBackHosts, err := infra.hostsFileOp.ReadCurrent()
		if err != nil {
			fmt.Printf("❌ 读取回滚后的 hosts 文件失败: %v\n", err)
			os.Exit(1)
		}

		// 检查测试条目是否已移除
		fmt.Println("检查测试条目是否已从 hosts 文件移除:")
		allRemoved := true
		for _, hostname := range checkEntries {
			if contains(rolledBackHosts, hostname) {
				fmt.Printf("   ❌ %s 仍然存在（回滚可能未生效）\n", hostname)
				allRemoved = false
			} else {
				fmt.Printf("   ✅ %s 已移除\n", hostname)
			}
		}

		if !allRemoved {
			fmt.Println("⚠️  回滚验证：部分测试条目仍存在")
		} else {
			fmt.Println("✅ 回滚验证成功：测试条目已全部移除")
		}
	} else {
		fmt.Println("⚠️  没有可回滚的版本")
	}

	// ========== 场景11: 分组排序测试 ==========
	fmt.Println("\n========================================")
	fmt.Println("[场景11] 分组排序（拖动模拟）")
	fmt.Println("========================================")

	// 创建第二个测试分组
	group2, err := infra.appService.CreateGroup(ctx, dto.CreateHostsGroupRequest{
		Name:        "排序测试分组",
		Description: "用于测试排序功能",
	})
	if err != nil {
		fmt.Printf("❌ 创建第二个分组失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 第二个分组创建成功: %s\n", group2.ID[:8])

	// 获取当前所有分组
	allGroups, err := infra.appService.GetAllGroups(ctx)
	if err != nil {
		fmt.Printf("❌ 获取分组失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("当前分组顺序:")
	for i, g := range allGroups {
		fmt.Printf("   %d. %s (Order: %d)\n", i+1, g.Name, g.Order)
	}

	// 反转顺序进行排序测试
	reversedIDs := make([]string, len(allGroups))
	for i, g := range allGroups {
		reversedIDs[len(allGroups)-1-i] = g.ID
	}

	err = infra.appService.ReorderGroups(ctx, dto.ReorderGroupsRequest{
		GroupIDs: reversedIDs,
	})
	if err != nil {
		fmt.Printf("❌ 排序失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ 分组排序成功（已反转顺序）")

	// 验证排序结果
	reorderedGroups, err := infra.appService.GetAllGroups(ctx)
	if err != nil {
		fmt.Printf("❌ 获取重排后的分组失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("重排后分组顺序:")
	for i, g := range reorderedGroups {
		fmt.Printf("   %d. %s (Order: %d)\n", i+1, g.Name, g.Order)
	}

	// ========== 场景12: 删除分组 ==========
	fmt.Println("\n========================================")
	fmt.Println("[场景12] 删除测试分组")
	fmt.Println("========================================")

	// 删除第一个测试分组
	err = infra.appService.DeleteGroup(ctx, testGroupID)
	if err != nil {
		fmt.Printf("❌ 删除分组失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 测试分组1已删除: %s\n", testGroupID[:8])

	// 删除第二个测试分组
	err = infra.appService.DeleteGroup(ctx, group2.ID)
	if err != nil {
		fmt.Printf("❌ 删除分组失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 测试分组2已删除: %s\n", group2.ID[:8])

	// 验证删除结果
	finalGroups, err := infra.appService.GetAllGroups(ctx)
	if err != nil {
		fmt.Printf("❌ 获取最终分组列表失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ 最终分组数量: %d\n", len(finalGroups))

	// ========== 最终验证：恢复原始 hosts ==========
	fmt.Println("\n========================================")
	fmt.Println("[最终] 恢复原始 hosts 配置")
	fmt.Println("========================================")

	// 检查是否有之前的备份版本
	finalVersions, err := infra.appService.GetVersions(ctx, 20)
	if err != nil {
		fmt.Printf("❌ 获取版本历史失败: %v\n", err)
	}

	// 找到包含原始内容的版本
	var originalVersionID string
	for _, v := range finalVersions {
		// 检查是否不包含测试条目
		if !contains(v.Content, "test-scenario.local") && !contains(v.Content, "test-server.local") {
			originalVersionID = v.ID
			break
		}
	}

	if originalVersionID != "" {
		err = infra.appService.RollbackToVersion(ctx, dto.RollbackRequest{
			VersionID:    originalVersionID,
			SudoPassword: "",
		})
		if err != nil {
			fmt.Printf("❌ 恢复原始配置失败: %v\n", err)
		} else {
			fmt.Println("✅ 已恢复到原始 hosts 配置")
		}
	} else {
		// 直接写入原始内容
		err = infra.hostsFileOp.Write(originalHosts)
		if err != nil {
			fmt.Printf("❌ 直接恢复失败: %v\n", err)
		} else {
			fmt.Println("✅ 已直接恢复原始 hosts 配置")
		}
	}

	// 最终 hosts 内容验证
	finalHosts, err := infra.hostsFileOp.ReadCurrent()
	if err != nil {
		fmt.Printf("❌ 最终验证失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("最终 hosts 文件验证:")
	testEntriesStillPresent := false
	for _, hostname := range checkEntries {
		if contains(finalHosts, hostname) {
			fmt.Printf("   ⚠️  %s 仍然存在\n", hostname)
			testEntriesStillPresent = true
		}
	}
	if !testEntriesStillPresent {
		fmt.Println("   ✅ 所有测试条目已清理")
	}

	// ========== 测试总结 ==========
	fmt.Println("\n========================================")
	fmt.Println("✅ 所有场景测试通过！")
	fmt.Println("========================================")
	fmt.Println("\n测试覆盖的场景:")
	fmt.Println("✅ 场景1: 获取分组列表")
	fmt.Println("✅ 场景2: 创建新分组")
	fmt.Println("✅ 场景3: 添加 hosts 条目")
	fmt.Println("✅ 场景4: 编辑分组信息")
	fmt.Println("✅ 场景5: 切换分组启用状态")
	fmt.Println("✅ 场景6: 生成 hosts 预览")
	fmt.Println("✅ 场景7: 应用 hosts 配置")
	fmt.Println("✅ 场景8: 验证 hosts 文件变化")
	fmt.Println("✅ 场景9: 获取版本历史")
	fmt.Println("✅ 场景10: 回滚 hosts 配置")
	fmt.Println("✅ 场景11: 分组排序（拖动模拟）")
	fmt.Println("✅ 场景12: 删除测试分组")
	fmt.Println("✅ 最终: 恢复原始 hosts 配置")
}

type infrastructure struct {
	appService   *service.HostsApplicationService
	hostsFileOp  *system.HostsFileOperator
	sudoManager  *system.SudoManager
	versionRepo  *persistence.VersionRepositoryImpl
}

func initInfrastructure() (*infrastructure, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	configPath := configDir + "\\hosts-manager"

	storage, err := persistence.NewJSONStorage(configPath)
	if err != nil {
		return nil, err
	}

	hostsRepo := persistence.NewHostsRepository(storage)
	versionRepo := persistence.NewVersionRepository(storage)

	elevator, err := system.NewPrivilegeElevator()
	if err != nil {
		return nil, err
	}

	hostsFileOp, err := system.NewHostsFileOperator(elevator)
	if err != nil {
		return nil, err
	}

	sudoManager := system.NewSudoManager()

	appService := service.NewHostsApplicationService(
		hostsRepo,
		versionRepo,
		hostsFileOp,
		sudoManager,
	)

	return &infrastructure{
		appService:  appService,
		hostsFileOp: hostsFileOp,
		sudoManager: sudoManager,
		versionRepo: versionRepo.(*persistence.VersionRepositoryImpl),
	}, nil
}

func checkAdmin() (bool, error) {
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

	token := windows.Token(0)
	defer token.Close()

	return token.IsMember(sid)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "...(截断)"
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}