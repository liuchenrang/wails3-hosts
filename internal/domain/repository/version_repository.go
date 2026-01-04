package repository

import (
	"context"
	"time"

	"github.com/chen/wails3-demo/internal/domain/entity"
)

// VersionRepository 定义版本历史的仓储接口
// 单一职责: 定义版本历史的持久化契约
type VersionRepository interface {
	// Save 保存一个版本记录
	Save(ctx context.Context, version *entity.HostsVersion) error

	// FindLatest 查找最新的 N 个版本
	FindLatest(ctx context.Context, limit int) ([]*entity.HostsVersion, error)

	// FindByID 根据 ID 查找版本
	FindByID(ctx context.Context, id string) (*entity.HostsVersion, error)

	// Delete 删除指定 ID 的版本
	Delete(ctx context.Context, id string) error

	// DeleteBefore 删除指定时间之前的所有版本
	DeleteBefore(ctx context.Context, before time.Time) error

	// Count 获取版本总数
	Count(ctx context.Context) (int, error)
}
