package persistence

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/chen/wails3-demo/internal/domain/entity"
)

// JSONStorage JSON 文件存储实现
// 单一职责: 负责配置和版本历史的 JSON 文件持久化
// DDD: 基础设施层实现领域层定义的仓储接口
type JSONStorage struct {
	configPath  string
	versionsPath string
	mu          sync.RWMutex
}

// ConfigData 配置文件数据结构
type ConfigData struct {
	Groups  []*entity.HostsGroup `json:"groups"`
	Settings Settings            `json:"settings"`
}

// Settings 应用设置
type Settings struct {
	Theme     string `json:"theme"`
	Language  string `json:"language"`
	AutoApply bool   `json:"auto_apply"`
}

// VersionsData 版本历史文件数据结构
type VersionsData struct {
	Versions []*entity.HostsVersion `json:"versions"`
}

// NewJSONStorage 创建 JSON 存储实例
// KISS: 使用简单的工厂函数
func NewJSONStorage(configDir string) (*JSONStorage, error) {
	// 确保目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	configPath := filepath.Join(configDir, "config.json")
	versionsPath := filepath.Join(configDir, "versions.json")

	storage := &JSONStorage{
		configPath:   configPath,
		versionsPath: versionsPath,
	}

	// 初始化配置文件（如果不存在）
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := storage.initConfigFile(); err != nil {
			return nil, err
		}
	}

	// 初始化版本文件（如果不存在）
	if _, err := os.Stat(versionsPath); os.IsNotExist(err) {
		if err := storage.initVersionsFile(); err != nil {
			return nil, err
		}
	}

	return storage, nil
}

// initConfigFile 初始化配置文件
func (s *JSONStorage) initConfigFile() error {
	config := ConfigData{
		Groups: []*entity.HostsGroup{},
		Settings: Settings{
			Theme:     "dark",
			Language:  "zh-CN",
			AutoApply: false,
		},
	}
	return s.saveConfig(config)
}

// initVersionsFile 初始化版本文件
func (s *JSONStorage) initVersionsFile() error {
	versions := VersionsData{
		Versions: []*entity.HostsVersion{},
	}
	return s.saveVersions(versions)
}

// LoadGroups 加载所有分组
func (s *JSONStorage) LoadGroups() ([]*entity.HostsGroup, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.configPath)
	if err != nil {
		return nil, err
	}

	var config ConfigData
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config.Groups, nil
}

// SaveGroups 保存所有分组
func (s *JSONStorage) SaveGroups(groups []*entity.HostsGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 读取现有配置
	data, err := os.ReadFile(s.configPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var config ConfigData
	if err == nil {
		if err := json.Unmarshal(data, &config); err != nil {
			return err
		}
	}

	// 更新分组
	config.Groups = groups

	return s.saveConfig(config)
}

// LoadSettings 加载设置
func (s *JSONStorage) LoadSettings() (*Settings, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.configPath)
	if err != nil {
		return nil, err
	}

	var config ConfigData
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config.Settings, nil
}

// SaveSettings 保存设置
func (s *JSONStorage) SaveSettings(settings Settings) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 读取现有配置
	data, err := os.ReadFile(s.configPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var config ConfigData
	if err == nil {
		if err := json.Unmarshal(data, &config); err != nil {
			return err
		}
	}

	// 更新设置
	config.Settings = settings

	return s.saveConfig(config)
}

// LoadVersions 加载版本历史
func (s *JSONStorage) LoadVersions() ([]*entity.HostsVersion, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.versionsPath)
	if err != nil {
		return nil, err
	}

	var versionsData VersionsData
	if err := json.Unmarshal(data, &versionsData); err != nil {
		return nil, err
	}

	return versionsData.Versions, nil
}

// SaveVersions 保存版本历史
func (s *JSONStorage) SaveVersions(versions []*entity.HostsVersion) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.saveVersions(VersionsData{Versions: versions})
}

// saveConfig 保存配置到文件（内部方法，不加锁）
func (s *JSONStorage) saveConfig(config ConfigData) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// 原子写入：先写临时文件，再重命名
	tmpPath := s.configPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, s.configPath)
}

// saveVersions 保存版本到文件（内部方法，不加锁）
func (s *JSONStorage) saveVersions(versions VersionsData) error {
	data, err := json.MarshalIndent(versions, "", "  ")
	if err != nil {
		return err
	}

	// 原子写入
	tmpPath := s.versionsPath + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, s.versionsPath)
}
