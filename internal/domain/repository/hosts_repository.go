package repository

import (
	"context"

	"github.com/chen/wails3-demo/internal/domain/entity"
)

// HostsRepository 定义 hosts 分组的仓储接口
// 单一职责: 定义 hosts 分组的持久化契约
// 接口隔离原则 (ISP): 仅包含必要的方法
// 依赖倒置原则 (DIP): 高层模块依赖抽象而非具体实现
type HostsRepository interface {
	// Save 保存一个 hosts 分组（创建或更新）
	Save(ctx context.Context, group *entity.HostsGroup) error

	// FindByID 根据 ID 查找分组
	FindByID(ctx context.Context, id string) (*entity.HostsGroup, error)

	// FindAll 查找所有分组
	FindAll(ctx context.Context) ([]*entity.HostsGroup, error)

	// Delete 删除指定 ID 的分组
	Delete(ctx context.Context, id string) error

	// ExistsByName 检查是否存在指定名称的分组
	ExistsByName(ctx context.Context, name string) (bool, error)
}
