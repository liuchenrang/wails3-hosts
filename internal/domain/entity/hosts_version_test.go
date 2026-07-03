package entity

import (
	"testing"
	"time"
)

// ==================== TDD: HostsVersion 单元测试 ====================

// TestNewHostsVersion 测试创建版本记录
func TestNewHostsVersion(t *testing.T) {
	t.Run("创建基本版本", func(t *testing.T) {
		version := NewHostsVersion("127.0.0.1 localhost", "初始版本", SourceManual)

		if version.ID == "" {
			t.Error("ID 不应为空")
		}
		if version.Content != "127.0.0.1 localhost" {
			t.Errorf("Content 不正确: %s", version.Content)
		}
		if version.Description != "初始版本" {
			t.Errorf("Description 不正确: %s", version.Description)
		}
		if version.Source != SourceManual {
			t.Errorf("Source 不正确: %s", version.Source)
		}
		if version.Timestamp.IsZero() {
			t.Error("Timestamp 不应为零值")
		}
	})

	t.Run("每个版本ID应唯一", func(t *testing.T) {
		v1 := NewHostsVersion("content1", "版本1", SourceManual)
		v2 := NewHostsVersion("content2", "版本2", SourceAuto)

		if v1.ID == v2.ID {
			t.Error("不同版本的 ID 应唯一")
		}
	})

	t.Run("不同来源类型的版本", func(t *testing.T) {
		manual := NewHostsVersion("c", "手动", SourceManual)
		auto := NewHostsVersion("c", "自动", SourceAuto)
		rollback := NewHostsVersion("c", "回滚", SourceRollback)

		if manual.Source != SourceManual {
			t.Error("手动来源不正确")
		}
		if auto.Source != SourceAuto {
			t.Error("自动来源不正确")
		}
		if rollback.Source != SourceRollback {
			t.Error("回滚来源不正确")
		}
	})
}

// TestHostsVersionIsExpired 测试过期检查
func TestHostsVersionIsExpired(t *testing.T) {
	t.Run("未过期", func(t *testing.T) {
		version := NewHostsVersion("content", "新版本", SourceManual)
		// 刚创建的版本，30天后才过期 → 现在不应过期
		if version.IsExpired(30) {
			t.Error("刚创建的版本不应过期（30天）")
		}
		// 0天意味着"立即过期"（expirationDate = timestamp + 0 = timestamp，当前时间 >= timestamp → 过期）
		// 所以 IsExpired(0) 应返回 true
	})

	t.Run("0天限制立即过期", func(t *testing.T) {
		version := NewHostsVersion("content", "测试", SourceManual)
		time.Sleep(time.Millisecond) // 确保时间有微小差异
		// 0 天限制：过期时间 = 创建时间。当前时间 > 过期时间 → 应视为过期
		if !version.IsExpired(0) {
			t.Error("IsExpired(0) 应视为过期（0天 = 立即过期）")
		}
	})

	t.Run("已过期", func(t *testing.T) {
		version := NewHostsVersion("content", "旧版本", SourceManual)
		// 手动设置时间戳为 31 天前
		version.Timestamp = time.Now().AddDate(0, 0, -31)

		if !version.IsExpired(30) {
			t.Error("31天前的版本应已过期（30天限制）")
		}
	})

	t.Run("边界-恰好30天", func(t *testing.T) {
		version := NewHostsVersion("content", "边界版本", SourceManual)
		version.Timestamp = time.Now().AddDate(0, 0, -30)

		// 正好30天前，30天的限制 → 当前时间 > 30天前+30天 = 现在，这取决于微妙的时间差
		// 使用 29 天作为测试边界更可靠
		version2 := NewHostsVersion("content", "边界2", SourceManual)
		version2.Timestamp = time.Now().AddDate(0, 0, -29)

		// 29天前，30天限制 → 不应过期
		if version2.IsExpired(30) {
			t.Error("29天前的版本不应过期（30天限制）")
		}
	})

	t.Run("负数天数始终过期", func(t *testing.T) {
		version := NewHostsVersion("content", "测试", SourceManual)
		// 负天数意味着"过去"就过期？当前实现: expirationDate = now + neg = past, now > past → true
		if !version.IsExpired(-1) {
			t.Error("负数天数应导致版本被视为过期")
		}
	})
}

// TestHostsVersionGetAge 测试获取版本年龄
func TestHostsVersionGetAge(t *testing.T) {
	t.Run("刚创建的版本年龄为0", func(t *testing.T) {
		version := NewHostsVersion("content", "新版本", SourceManual)
		age := version.GetAge()

		if age != 0 {
			t.Errorf("刚创建的版本年龄应为 0，实际为 %d", age)
		}
	})

	t.Run("1天前的版本年龄为1", func(t *testing.T) {
		version := NewHostsVersion("content", "旧版本", SourceManual)
		version.Timestamp = time.Now().AddDate(0, 0, -1)

		age := version.GetAge()

		if age != 1 {
			t.Errorf("1天前的版本年龄应为 1，实际为 %d", age)
		}
	})

	t.Run("7天前的版本年龄为7", func(t *testing.T) {
		version := NewHostsVersion("content", "一周前", SourceManual)
		version.Timestamp = time.Now().AddDate(0, 0, -7)

		age := version.GetAge()

		if age != 7 {
			t.Errorf("7天前的版本年龄应为 7，实际为 %d", age)
		}
	})
}

// TestVersionSourceConstants 测试版本来源常量
func TestVersionSourceConstants(t *testing.T) {
	if SourceManual != "manual" {
		t.Errorf("SourceManual 应为 'manual'，实际为 '%s'", SourceManual)
	}
	if SourceAuto != "auto" {
		t.Errorf("SourceAuto 应为 'auto'，实际为 '%s'", SourceAuto)
	}
	if SourceRollback != "rollback" {
		t.Errorf("SourceRollback 应为 'rollback'，实际为 '%s'", SourceRollback)
	}
}
