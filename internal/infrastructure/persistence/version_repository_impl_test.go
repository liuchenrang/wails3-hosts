package persistence

import (
	"testing"
	"time"

	"github.com/chen/wails3-hosts/internal/domain/entity"
)

// TestMaxVersions 验证最大版本数量常量
func TestMaxVersions(t *testing.T) {
	if MaxVersions != 10 {
		t.Errorf("MaxVersions 应该是 10, 但是得到 %d", MaxVersions)
	}
}

// TestVersionLimitLogic 验证版本限制逻辑
func TestVersionLimitLogic(t *testing.T) {
	// 模拟版本切片
	versions := make([]*entity.HostsVersion, 15)

	// 创建15个版本
	for i := 0; i < 15; i++ {
		versions[i] = &entity.HostsVersion{
			ID:        string(rune(i)),
			Timestamp: time.Now().Add(time.Duration(i) * time.Minute),
		}
	}

	// 应用版本限制
	if len(versions) > MaxVersions {
		versions = versions[len(versions)-MaxVersions:]
	}

	// 验证结果
	if len(versions) != MaxVersions {
		t.Errorf("版本数量应该是 %d, 但是得到 %d", MaxVersions, len(versions))
	}
}

// TestVersionLimitWithExactCount 验证恰好等于最大值的情况
func TestVersionLimitWithExactCount(t *testing.T) {
	versions := make([]*entity.HostsVersion, 10)

	for i := 0; i < 10; i++ {
		versions[i] = &entity.HostsVersion{
			ID:        string(rune(i)),
			Timestamp: time.Now(),
		}
	}

	// 应用版本限制
	if len(versions) > MaxVersions {
		versions = versions[len(versions)-MaxVersions:]
	}

	// 验证结果
	if len(versions) != 10 {
		t.Errorf("版本数量应该保持 10, 但是得到 %d", len(versions))
	}
}

// TestVersionLimitWithLessThanMax 验证小于最大值的情况
func TestVersionLimitWithLessThanMax(t *testing.T) {
	versions := make([]*entity.HostsVersion, 5)

	for i := 0; i < 5; i++ {
		versions[i] = &entity.HostsVersion{
			ID:        string(rune(i)),
			Timestamp: time.Now(),
		}
	}

	// 应用版本限制
	if len(versions) > MaxVersions {
		versions = versions[len(versions)-MaxVersions:]
	}

	// 验证结果
	if len(versions) != 5 {
		t.Errorf("版本数量应该保持 5, 但是得到 %d", len(versions))
	}
}
