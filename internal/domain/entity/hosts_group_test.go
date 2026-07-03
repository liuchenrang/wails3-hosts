package entity

import (
	"testing"
)

// ==================== TDD: HostsGroup 单元测试 ====================

// TestNewHostsGroup 测试创建 hosts 分组
func TestNewHostsGroup(t *testing.T) {
	t.Run("创建基本分组", func(t *testing.T) {
		group := NewHostsGroup("开发环境", "本地开发配置")

		if group.ID == "" {
			t.Error("ID 不应为空")
		}
		if group.Name != "开发环境" {
			t.Errorf("Name 应为 '开发环境'，实际为 '%s'", group.Name)
		}
		if group.Description != "本地开发配置" {
			t.Errorf("Description 应为 '本地开发配置'，实际为 '%s'", group.Description)
		}
		if group.IsEnabled != false {
			t.Error("新创建的分组应默认不启用（安全默认值）")
		}
		if group.Order != 0 {
			t.Errorf("Order 应为 0，实际为 %d", group.Order)
		}
		if len(group.Entries) != 0 {
			t.Error("新分组的条目列表应为空")
		}
		if group.CreatedAt.IsZero() {
			t.Error("CreatedAt 不应为零值")
		}
		if group.UpdatedAt.IsZero() {
			t.Error("UpdatedAt 不应为零值")
		}
		if !group.CreatedAt.Equal(group.UpdatedAt) {
			t.Error("新创建时 CreatedAt 和 UpdatedAt 应相等")
		}
	})

	t.Run("每个分组的ID应唯一", func(t *testing.T) {
		g1 := NewHostsGroup("分组1", "")
		g2 := NewHostsGroup("分组2", "")

		if g1.ID == g2.ID {
			t.Error("不同分组的 ID 应唯一")
		}
	})
}

// TestHostsGroupAddEntry 测试添加条目的各种情况
func TestHostsGroupAddEntry(t *testing.T) {
	t.Run("添加有效条目", func(t *testing.T) {
		group := NewHostsGroup("测试分组", "")
		entry := *NewHostsEntry("127.0.0.1", "test.local", "测试")

		err := group.AddEntry(entry)

		if err != nil {
			t.Errorf("添加有效条目不应出错: %v", err)
		}
		if len(group.Entries) != 1 {
			t.Errorf("条目数应为 1，实际为 %d", len(group.Entries))
		}
		if group.Entries[0].IP != "127.0.0.1" {
			t.Errorf("条目的 IP 不正确: %s", group.Entries[0].IP)
		}
		if group.UpdatedAt.IsZero() {
			t.Error("添加条目后 UpdatedAt 不应为零值")
		}
	})

	t.Run("添加无效条目应返回错误", func(t *testing.T) {
		group := NewHostsGroup("测试分组", "")
		invalidEntry := *NewHostsEntry("", "test.local", "") // 空 IP

		err := group.AddEntry(invalidEntry)

		if err == nil {
			t.Error("添加无效条目应返回错误")
		}
		if len(group.Entries) != 0 {
			t.Error("无效条目不应被添加到分组中")
		}
	})

	t.Run("添加多个条目", func(t *testing.T) {
		group := NewHostsGroup("测试分组", "")
		e1 := *NewHostsEntry("127.0.0.1", "a.local", "")
		e2 := *NewHostsEntry("127.0.0.2", "b.local", "")
		e3 := *NewHostsEntry("127.0.0.3", "c.local", "")

		_ = group.AddEntry(e1)
		_ = group.AddEntry(e2)
		_ = group.AddEntry(e3)

		if len(group.Entries) != 3 {
			t.Errorf("条目数应为 3，实际为 %d", len(group.Entries))
		}
	})
}

// TestHostsGroupRemoveEntry 测试移除条目
func TestHostsGroupRemoveEntry(t *testing.T) {
	t.Run("移除存在的条目", func(t *testing.T) {
		group := NewHostsGroup("测试分组", "")
		entry := *NewHostsEntry("127.0.0.1", "test.local", "")
		_ = group.AddEntry(entry)

		removed := group.RemoveEntry(entry.ID)

		if !removed {
			t.Error("移除存在的条目应返回 true")
		}
		if len(group.Entries) != 0 {
			t.Errorf("移除后条目数应为 0，实际为 %d", len(group.Entries))
		}
	})

	t.Run("移除不存在的条目", func(t *testing.T) {
		group := NewHostsGroup("测试分组", "")
		entry := *NewHostsEntry("127.0.0.1", "test.local", "")
		_ = group.AddEntry(entry)

		removed := group.RemoveEntry("non-existent-id")

		if removed {
			t.Error("移除不存在的条目应返回 false")
		}
		if len(group.Entries) != 1 {
			t.Errorf("条目数应保持为 1，实际为 %d", len(group.Entries))
		}
	})

	t.Run("从多个条目中移除指定项", func(t *testing.T) {
		group := NewHostsGroup("测试分组", "")
		e1 := *NewHostsEntry("127.0.0.1", "a.local", "")
		e2 := *NewHostsEntry("127.0.0.2", "b.local", "")
		e3 := *NewHostsEntry("127.0.0.3", "c.local", "")
		_ = group.AddEntry(e1)
		_ = group.AddEntry(e2)
		_ = group.AddEntry(e3)

		// 移除中间的
		group.RemoveEntry(e2.ID)

		if len(group.Entries) != 2 {
			t.Errorf("移除后条目数应为 2，实际为 %d", len(group.Entries))
		}
		// 验证剩下的条目
		if group.Entries[0].ID != e1.ID {
			t.Error("e1 应该保留")
		}
		if group.Entries[1].ID != e3.ID {
			t.Error("e3 应该保留")
		}
	})

	t.Run("从空分组移除", func(t *testing.T) {
		group := NewHostsGroup("测试分组", "")

		removed := group.RemoveEntry("any-id")

		if removed {
			t.Error("从空分组移除应返回 false")
		}
	})
}

// TestHostsGroupUpdateEntry 测试更新条目
func TestHostsGroupUpdateEntry(t *testing.T) {
	t.Run("更新存在的条目", func(t *testing.T) {
		group := NewHostsGroup("测试分组", "")
		original := *NewHostsEntry("127.0.0.1", "old.local", "旧")
		_ = group.AddEntry(original)

		updated := *NewHostsEntry("192.168.1.1", "new.local", "新")
		// 保持相同 ID
		updated.ID = original.ID

		err := group.UpdateEntry(original.ID, updated)

		if err != nil {
			t.Errorf("更新应成功: %v", err)
		}
		if group.Entries[0].IP != "192.168.1.1" {
			t.Errorf("IP 应更新: %s", group.Entries[0].IP)
		}
		if group.Entries[0].Hostname != "new.local" {
			t.Errorf("Hostname 应更新: %s", group.Entries[0].Hostname)
		}
	})

	t.Run("更新不存在的条目", func(t *testing.T) {
		group := NewHostsGroup("测试分组", "")
		updated := *NewHostsEntry("127.0.0.1", "test.local", "")

		err := group.UpdateEntry("non-existent", updated)

		if err != ErrEntryNotFound {
			t.Errorf("应返回 ErrEntryNotFound，实际: %v", err)
		}
	})

	t.Run("更新为无效条目", func(t *testing.T) {
		group := NewHostsGroup("测试分组", "")
		original := *NewHostsEntry("127.0.0.1", "old.local", "")
		_ = group.AddEntry(original)

		invalid := *NewHostsEntry("", "new.local", "") // 空 IP
		invalid.ID = original.ID

		err := group.UpdateEntry(original.ID, invalid)

		if err == nil {
			t.Error("更新为无效条目应返回错误")
		}
	})
}

// TestHostsGroupClearEntries 测试清空条目
func TestHostsGroupClearEntries(t *testing.T) {
	t.Run("清空所有条目", func(t *testing.T) {
		group := NewHostsGroup("测试分组", "")
		_ = group.AddEntry(*NewHostsEntry("127.0.0.1", "a.local", ""))
		_ = group.AddEntry(*NewHostsEntry("127.0.0.2", "b.local", ""))

		group.ClearEntries()

		if len(group.Entries) != 0 {
			t.Errorf("清空后条目数应为 0，实际为 %d", len(group.Entries))
		}
	})

	t.Run("清空空分组不会出错", func(t *testing.T) {
		group := NewHostsGroup("测试分组", "")

		// 不应 panic
		group.ClearEntries()

		if len(group.Entries) != 0 {
			t.Error("清空空分组后条目数仍应为 0")
		}
	})
}

// TestHostsGroupToggle 测试切换启用状态
func TestHostsGroupToggle(t *testing.T) {
	tests := []struct {
		name     string
		initial  bool
		expected bool
	}{
		{"禁用→启用", false, true},
		{"启用→禁用", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := NewHostsGroup("测试", "")
			group.IsEnabled = tt.initial

			group.Toggle()

			if group.IsEnabled != tt.expected {
				t.Errorf("Toggle 后 IsEnabled 应为 %v，实际为 %v", tt.expected, group.IsEnabled)
			}
		})
	}
}

// TestHostsGroupGetEnabledEntries 测试获取启用的条目
func TestHostsGroupGetEnabledEntries(t *testing.T) {
	t.Run("分组禁用时返回空", func(t *testing.T) {
		group := NewHostsGroup("测试", "")
		group.IsEnabled = false
		_ = group.AddEntry(*NewHostsEntry("127.0.0.1", "test.local", ""))

		enabled := group.GetEnabledEntries()

		if len(enabled) != 0 {
			t.Errorf("分组禁用时应返回空列表，实际 %d 个条目", len(enabled))
		}
	})

	t.Run("分组启用但条目禁用时返回空", func(t *testing.T) {
		group := NewHostsGroup("测试", "")
		group.IsEnabled = true
		entry := *NewHostsEntry("127.0.0.1", "test.local", "")
		entry.Enabled = false
		_ = group.AddEntry(entry)

		enabled := group.GetEnabledEntries()

		if len(enabled) != 0 {
			t.Errorf("所有条目禁用时应返回空列表，实际 %d 个条目", len(enabled))
		}
	})

	t.Run("分组和条目都启用时返回条目", func(t *testing.T) {
		group := NewHostsGroup("测试", "")
		group.IsEnabled = true
		e1 := *NewHostsEntry("127.0.0.1", "a.local", "")
		e1.Enabled = true
		e2 := *NewHostsEntry("127.0.0.2", "b.local", "")
		e2.Enabled = false // 禁用
		e3 := *NewHostsEntry("127.0.0.3", "c.local", "")
		e3.Enabled = true

		_ = group.AddEntry(e1)
		_ = group.AddEntry(e2)
		_ = group.AddEntry(e3)

		enabled := group.GetEnabledEntries()

		if len(enabled) != 2 {
			t.Errorf("应有 2 个启用的条目，实际 %d 个", len(enabled))
		}
	})

	t.Run("空分组返回空列表", func(t *testing.T) {
		group := NewHostsGroup("测试", "")
		group.IsEnabled = true

		enabled := group.GetEnabledEntries()

		if len(enabled) != 0 {
			t.Errorf("空分组应返回空列表，实际 %d 个", len(enabled))
		}
		if enabled == nil {
			t.Error("返回的空列表不应为 nil")
		}
	})
}

// TestHostsGroupSetEnabled 测试设置启用状态
func TestHostsGroupSetEnabled(t *testing.T) {
	group := NewHostsGroup("测试", "")

	group.SetEnabled(true)

	if !group.IsEnabled {
		t.Error("SetEnabled(true) 后应为启用状态")
	}
	// UpdatedAt 应被更新（非零值），不检查时间顺序（调用可能在同一纳秒完成）
	if group.UpdatedAt.IsZero() {
		t.Error("SetEnabled 后 UpdatedAt 不应为零值")
	}

	// 验证 disable
	group.SetEnabled(false)
	if group.IsEnabled {
		t.Error("SetEnabled(false) 后应为禁用状态")
	}
}

// TestHostsGroupSetOrder 测试设置排序
func TestHostsGroupSetOrder(t *testing.T) {
	group := NewHostsGroup("测试", "")

	group.SetOrder(5)

	if group.Order != 5 {
		t.Errorf("Order 应为 5，实际为 %d", group.Order)
	}
	if group.UpdatedAt.IsZero() {
		t.Error("SetOrder 后 UpdatedAt 不应为零值")
	}

	// 测试负值
	group.SetOrder(-1)
	if group.Order != -1 {
		t.Errorf("Order 应支持负值，实际为 %d", group.Order)
	}

	// 测试零值
	group.SetOrder(0)
	if group.Order != 0 {
		t.Errorf("Order 应为 0，实际为 %d", group.Order)
	}
}

// TestDomainError 测试领域错误类型
func TestDomainError(t *testing.T) {
	err := &DomainError{Code: "TEST", Message: "测试错误"}

	if err.Error() != "测试错误" {
		t.Errorf("Error() 应返回 Message")
	}

	if err.Code != "TEST" {
		t.Errorf("Code 不正确: %s", err.Code)
	}
}
