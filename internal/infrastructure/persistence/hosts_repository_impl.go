package persistence

import (
	"context"
	"fmt"

	"github.com/chen/wails3-hosts/internal/domain/entity"
	"github.com/chen/wails3-hosts/internal/domain/repository"
)

// HostsRepositoryImpl hosts 分组仓储实现
// 单一职责: 实现 HostsRepository 接口
// DDD: 基础设施层提供接口的具体实现
type HostsRepositoryImpl struct {
	storage *JSONStorage
}

// NewHostsRepository 创建仓储实例
func NewHostsRepository(storage *JSONStorage) repository.HostsRepository {
	return &HostsRepositoryImpl{
		storage: storage,
	}
}

// Save 保存分组（创建或更新）
// DRY: 统一的保存逻辑
func (r *HostsRepositoryImpl) Save(ctx context.Context, group *entity.HostsGroup) error {
	groups, err := r.storage.LoadGroups()
	if err != nil {
		return err
	}

	// 查找并更新或添加
	found := false
	for i, g := range groups {
		if g.ID == group.ID {
			groups[i] = group
			found = true
			break
		}
	}

	if !found {
		groups = append(groups, group)
	}

	return r.storage.SaveGroups(groups)
}

// FindByID 根据 ID 查找分组
func (r *HostsRepositoryImpl) FindByID(ctx context.Context, id string) (*entity.HostsGroup, error) {
	groups, err := r.storage.LoadGroups()
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		if group.ID == id {
			return group, nil
		}
	}

	return nil, fmt.Errorf("分组不存在")
}

// FindAll 查找所有分组
func (r *HostsRepositoryImpl) FindAll(ctx context.Context) ([]*entity.HostsGroup, error) {
	return r.storage.LoadGroups()
}

// Delete 删除分组
func (r *HostsRepositoryImpl) Delete(ctx context.Context, id string) error {
	groups, err := r.storage.LoadGroups()
	if err != nil {
		return err
	}

	// 过滤掉要删除的分组
	newGroups := make([]*entity.HostsGroup, 0, len(groups))
	for _, group := range groups {
		if group.ID != id {
			newGroups = append(newGroups, group)
		}
	}

	return r.storage.SaveGroups(newGroups)
}

// ExistsByName 检查是否存在指定名称的分组
func (r *HostsRepositoryImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	groups, err := r.storage.LoadGroups()
	if err != nil {
		return false, err
	}

	for _, group := range groups {
		if group.Name == name {
			return true, nil
		}
	}

	return false, nil
}
