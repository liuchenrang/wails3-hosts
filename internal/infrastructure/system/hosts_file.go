package system

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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
// 注意：此方法不接收密码参数，密码必须提前通过 SudoManager 缓存
func (o *HostsFileOperator) Write(content string) error {
	// 此方法不再需要密码参数
	// 密码应该已经通过 ValidateSudoPassword 验证并缓存
	// 直接调用 sudo，使用缓存的凭证
	script := fmt.Sprintf("cat > %s", o.hostsFilePath)

	// 创建 sudo 命令，不设置密码（使用系统缓存的凭证）
	cmd := exec.Command("sudo", "sh", "-c", script)
	cmd.Stdin = strings.NewReader(content)

	// 捕获输出
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("写入 hosts 文件失败: %w, stderr: %s", err, stderr.String())
	}

	return nil
}

// WriteWithPassword 写入内容到 hosts 文件（使用提供的 sudo 密码）
// 使用场景：直接使用用户提供的密码进行 sudo 操作
func (o *HostsFileOperator) WriteWithPassword(content string, password string) error {
	fmt.Println("[HostsFileOp] WriteWithPassword 开始", "路径:", o.hostsFilePath, "内容长度:", len(content), "密码长度:", len(password))

	script := fmt.Sprintf("cat > %s", o.hostsFilePath)
	fmt.Println("[HostsFileOp] 执行脚本:", script)

	// 使用 SudoCommand 包装器，它会自动处理密码输入
	cmd := NewSudoCommand([]string{"sh", "-c", script})
	cmd.SetPassword(password)
	cmd.SetStdin([]byte(content))

	fmt.Println("[HostsFileOp] 开始执行 SudoCommand")
	if err := cmd.Run(); err != nil {
		fmt.Println("[HostsFileOp] SudoCommand 执行失败:", err.Error())
		return fmt.Errorf("写入 hosts 文件失败: %w", err)
	}

	fmt.Println("[HostsFileOp] SudoCommand 执行成功")
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
