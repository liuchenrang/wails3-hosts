package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/chen/wails3-hosts/internal/domain/entity"
	"github.com/chen/wails3-hosts/internal/domain/repository"
)

const (
	// MaxVersions 最大版本数量
	MaxVersions = 10
	// VersionExpirationDays 版本过期天数
	VersionExpirationDays = 30
)

// VersionRepositoryImpl 版本历史仓储实现
// 单一职责: 实现 VersionRepository 接口
type VersionRepositoryImpl struct {
	storage *JSONStorage
}

// NewVersionRepository 创建仓储实例
func NewVersionRepository(storage *JSONStorage) repository.VersionRepository {
	return &VersionRepositoryImpl{
		storage: storage,
	}
}

// Save 保存版本记录
func (r *VersionRepositoryImpl) Save(ctx context.Context, version *entity.HostsVersion) error {
	versions, err := r.storage.LoadVersions()
	if err != nil {
		return err
	}

	// 添加新版本
	versions = append(versions, version)

	// 清理过期版本
	versions = r.cleanupExpiredVersions(versions)

	// 限制版本数量
	if len(versions) > MaxVersions {
		versions = versions[len(versions)-MaxVersions:]
	}

	return r.storage.SaveVersions(versions)
}

// FindLatest 查找最新的 N 个版本
func (r *VersionRepositoryImpl) FindLatest(ctx context.Context, limit int) ([]*entity.HostsVersion, error) {
	versions, err := r.storage.LoadVersions()
	if err != nil {
		return nil, err
	}

	// 按时间倒序排序
	sortedVersions := make([]*entity.HostsVersion, len(versions))
	copy(sortedVersions, versions)

	// 简单的冒泡排序（版本数量不大，性能可接受）
	for i := 0; i < len(sortedVersions)-1; i++ {
		for j := 0; j < len(sortedVersions)-1-i; j++ {
			if sortedVersions[j].Timestamp.Before(sortedVersions[j+1].Timestamp) {
				sortedVersions[j], sortedVersions[j+1] = sortedVersions[j+1], sortedVersions[j]
			}
		}
	}

	// 应用限制
	if limit > 0 && len(sortedVersions) > limit {
		sortedVersions = sortedVersions[:limit]
	}

	return sortedVersions, nil
}

// FindByID 根据 ID 查找版本
func (r *VersionRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.HostsVersion, error) {
	versions, err := r.storage.LoadVersions()
	if err != nil {
		return nil, err
	}

	for _, version := range versions {
		if version.ID == id {
			return version, nil
		}
	}

	return nil, fmt.Errorf("版本不存在")
}

// Delete 删除指定版本
func (r *VersionRepositoryImpl) Delete(ctx context.Context, id string) error {
	versions, err := r.storage.LoadVersions()
	if err != nil {
		return err
	}

	// 过滤掉要删除的版本
	newVersions := make([]*entity.HostsVersion, 0, len(versions))
	for _, version := range versions {
		if version.ID != id {
			newVersions = append(newVersions, version)
		}
	}

	return r.storage.SaveVersions(newVersions)
}

// DeleteBefore 删除指定时间之前的所有版本
func (r *VersionRepositoryImpl) DeleteBefore(ctx context.Context, before time.Time) error {
	versions, err := r.storage.LoadVersions()
	if err != nil {
		return err
	}

	// 过滤
	newVersions := make([]*entity.HostsVersion, 0, len(versions))
	for _, version := range versions {
		if version.Timestamp.After(before) || version.Timestamp.Equal(before) {
			newVersions = append(newVersions, version)
		}
	}

	return r.storage.SaveVersions(newVersions)
}

// Count 获取版本总数
func (r *VersionRepositoryImpl) Count(ctx context.Context) (int, error) {
	versions, err := r.storage.LoadVersions()
	if err != nil {
		return 0, err
	}
	return len(versions), nil
}

// cleanupExpiredVersions 清理过期版本
// YAGNI: 简单的清理逻辑，不处理复杂的过期策略
func (r *VersionRepositoryImpl) cleanupExpiredVersions(versions []*entity.HostsVersion) []*entity.HostsVersion {
	cleaned := make([]*entity.HostsVersion, 0, len(versions))
	for _, version := range versions {
		if !version.IsExpired(VersionExpirationDays) {
			cleaned = append(cleaned, version)
		}
	}

	return cleaned
}
