package system

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// HostsFileOperator hosts 文件操作器
// 单一职责: 负责系统 hosts 文件的读写操作
// DDD: 基础设施层提供与外部系统的交互
//
// 设计原则应用:
// - D: 依赖 PrivilegeElevator 抽象接口，而非具体实现
// - O: 通过接口支持不同平台的权限提升方式
type HostsFileOperator struct {
	hostsFilePath string
	backupDir     string
	elevator      PrivilegeElevator // 依赖注入: 平台特定的权限提升器
}

// NewHostsFileOperator 创建操作器实例
//
// 参数:
//   - elevator: 平台特定的权限提升器（由外部注入）
//
// 返回:
//   - *HostsFileOperator: 操作器实例
//   - error: 创建失败时的错误
//
// 使用示例:
//   elevator, _ := system.NewPrivilegeElevator()
//   operator, err := system.NewHostsFileOperator(elevator)
func NewHostsFileOperator(elevator PrivilegeElevator) (*HostsFileOperator, error) {
	hostsPath, err := getHostsFilePath()
	if err != nil {
		return nil, err
	}

	// 创建备份目录
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	backupDir := filepath.Join(configDir, "hosts-manager", "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, err
	}

	return &HostsFileOperator{
		hostsFilePath: hostsPath,
		backupDir:     backupDir,
		elevator:      elevator, // 注入依赖
	}, nil
}

// CanCacheCredentials 检查是否可以缓存凭据
// 实现: 委托给提升器接口
//
// 用途: 前端根据此值决定是否显示"密码已缓存"提示
func (o *HostsFileOperator) CanCacheCredentials() bool {
	return o.elevator.CanCacheCredentials()
}

// GetPrivilegeElevator 获取权限提升器实例
// 实现: 返回内部提升器接口
//
// 用途: 应用服务层需要访问提升器获取平台信息
func (o *HostsFileOperator) GetPrivilegeElevator() PrivilegeElevator {
	return o.elevator
}

// getHostsFilePath 获取系统 hosts 文件路径
// KISS: 根据操作系统返回对应的路径
func getHostsFilePath() (string, error) {
	switch runtime.GOOS {
	case "darwin", "linux":
		return "/etc/hosts", nil
	case "windows":
		return filepath.Join(os.Getenv("SystemRoot"), "System32", "drivers", "etc", "hosts"), nil
	default:
		return "", fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}
}

// ReadCurrent 读取当前 hosts 文件内容
func (o *HostsFileOperator) ReadCurrent() (string, error) {
	data, err := os.ReadFile(o.hostsFilePath)
	if err != nil {
		return "", fmt.Errorf("读取 hosts 文件失败: %w", err)
	}
	return string(data), nil
}

// Write 写入内容到 hosts 文件（需要提升权限）
// 实现: 使用提升器接口执行写入操作
//
// 注意:
// - Unix: 应先调用 ValidateSudoPassword 验证密码
// - Windows: 会自动弹出 UAC 提示
func (o *HostsFileOperator) Write(content string) error {
	// 委托给提升器接口
	return o.elevator.Execute(content)
}

// WriteWithPassword 写入内容到 hosts 文件（使用提供的凭据）
// 使用场景：直接使用用户提供的密码进行 sudo 操作
//
// 平台差异:
// - Unix: 使用密码验证并执行写入
// - Windows: 密码参数被忽略，通过 UAC 提权
//
// 注意: 此方法为向后兼容保留，新代码应优先使用 Validate + Write 组合
func (o *HostsFileOperator) WriteWithPassword(content string, password string) error {
	fmt.Println("[HostsFileOp] WriteWithPassword 开始", "路径:", o.hostsFilePath, "内容长度:", len(content), "密码长度:", len(password))

	// Unix: 先验证密码（缓存到系统）
	// Windows: password 被忽略
	if !o.elevator.Validate(password) {
		return fmt.Errorf("凭据验证失败")
	}

	// 执行写入
	if err := o.elevator.Execute(content); err != nil {
		fmt.Println("[HostsFileOp] 提升器执行失败:", err.Error())
		return fmt.Errorf("写入 hosts 文件失败: %w", err)
	}

	fmt.Println("[HostsFileOp] 写入成功")
	return nil
}

// Backup 备份当前 hosts 文件
func (o *HostsFileOperator) Backup() error {
	// 读取当前内容
	content, err := o.ReadCurrent()
	if err != nil {
		return err
	}

	// 生成备份文件名（带时间戳）
	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(o.backupDir, fmt.Sprintf("hosts_%s.bak", timestamp))

	// 写入备份文件
	if err := os.WriteFile(backupPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("创建备份失败: %w", err)
	}

	// 清理旧备份（保留最近 5 个）
	return o.cleanupOldBackups()
}

// cleanupOldBackups 清理旧备份文件
// YAGNI: 简单的清理逻辑，只保留最近 5 个备份
func (o *HostsFileOperator) cleanupOldBackups() error {
	entries, err := os.ReadDir(o.backupDir)
	if err != nil {
		return err
	}

	// 获取所有备份文件及其修改时间
	type backupInfo struct {
		name    string
		modTime time.Time
	}

	var backups []backupInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		name := entry.Name()
		if len(name) > 10 && name[:6] == "hosts_" && name[len(name)-4:] == ".bak" {
			backups = append(backups, backupInfo{
				name:    name,
				modTime: info.ModTime(),
			})
		}
	}

	// 按修改时间排序（最新的在前）
	for i := 0; i < len(backups)-1; i++ {
		for j := 0; j < len(backups)-1-i; j++ {
			if backups[j].modTime.Before(backups[j+1].modTime) {
				backups[j], backups[j+1] = backups[j+1], backups[j]
			}
		}
	}

	// 删除超过 5 个的旧备份
	const maxBackups = 5
	if len(backups) > maxBackups {
		for i := maxBackups; i < len(backups); i++ {
			backupPath := filepath.Join(o.backupDir, backups[i].name)
			os.Remove(backupPath)
		}
	}

	return nil
}

// RestoreFromBackup 从备份恢复 hosts 文件
func (o *HostsFileOperator) RestoreFromBackup(backupPath string) error {
	content, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("读取备份文件失败: %w", err)
	}

	// 注意: 这需要 sudo 权限，使用系统缓存的凭证
	return o.Write(string(content))
}

// GetBackupList 获取所有备份文件列表
func (o *HostsFileOperator) GetBackupList() ([]string, error) {
	entries, err := os.ReadDir(o.backupDir)
	if err != nil {
		return nil, err
	}

	backups := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if len(name) > 10 && name[:6] == "hosts_" && name[len(name)-4:] == ".bak" {
			backups = append(backups, filepath.Join(o.backupDir, name))
		}
	}

	return backups, nil
}
