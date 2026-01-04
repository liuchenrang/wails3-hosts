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
type HostsFileOperator struct {
	hostsFilePath string
	backupDir     string
}

// NewHostsFileOperator 创建操作器实例
func NewHostsFileOperator() (*HostsFileOperator, error) {
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
	}, nil
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

// Write 写入内容到 hosts 文件（需要 sudo 权限）
func (o *HostsFileOperator) Write(content, sudoPassword string) error {
	// 使用 sudo tee 命令写入文件
	// 这比直接使用 sudo 更安全，因为不需要将密码通过命令行参数传递
	cmd := NewSudoCommand([]string{"tee", o.hostsFilePath})
	cmd.SetStdin([]byte(content + "\n"))
	cmd.SetPassword(sudoPassword)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("写入 hosts 文件失败: %w", err)
	}

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

	// 注意: 这需要 sudo 权限，调用者需要提供密码
	return o.Write(string(content), "")
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
