package entity

import (
	"strings"
	"testing"
)

// ==================== TDD: HostsEntry 单元测试 ====================

// TestNewHostsEntry 测试创建 hosts 条目
func TestNewHostsEntry(t *testing.T) {
	t.Run("创建基本条目", func(t *testing.T) {
		entry := NewHostsEntry("127.0.0.1", "localhost", "本地回环")

		if entry.ID == "" {
			t.Error("ID 不应为空")
		}
		if entry.IP != "127.0.0.1" {
			t.Errorf("IP 应为 '127.0.0.1'，实际为 '%s'", entry.IP)
		}
		if entry.Hostname != "localhost" {
			t.Errorf("Hostname 应为 'localhost'，实际为 '%s'", entry.Hostname)
		}
		if entry.Comment != "本地回环" {
			t.Errorf("Comment 应为 '本地回环'，实际为 '%s'", entry.Comment)
		}
		if entry.Enabled != true {
			t.Error("新创建的条目应默认启用")
		}
	})

	t.Run("自动去除空白字符", func(t *testing.T) {
		entry := NewHostsEntry(" 127.0.0.1 ", " localhost ", " 备注 ")

		if entry.IP != "127.0.0.1" {
			t.Errorf("IP 应去除空白，实际为 '%s'", entry.IP)
		}
		if entry.Hostname != "localhost" {
			t.Errorf("Hostname 应去除空白，实际为 '%s'", entry.Hostname)
		}
		if entry.Comment != "备注" {
			t.Errorf("Comment 应去除空白，实际为 '%s'", entry.Comment)
		}
	})

	t.Run("每个条目的ID应唯一", func(t *testing.T) {
		e1 := NewHostsEntry("127.0.0.1", "a.com", "")
		e2 := NewHostsEntry("127.0.0.1", "b.com", "")

		if e1.ID == e2.ID {
			t.Error("不同条目的 ID 应唯一")
		}
	})
}

// TestHostsEntryValidate 测试条目验证
func TestHostsEntryValidate(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		host    string
		wantErr bool
	}{
		// 有效的 IPv4
		{"有效IPv4", "127.0.0.1", "localhost", false},
		{"有效IPv4-公网", "8.8.8.8", "dns.google", false},
		{"有效IPv4-局域网", "192.168.1.1", "router.local", false},

		// 有效的 IPv6
		{"有效IPv6-完整", "::1", "localhost6", false},
		{"有效IPv6-简写", "fe80::1", "linklocal", false},
		{"有效IPv6-全写", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", "ipv6.test", false},

		// 无效的 IP
		{"无效IP-空字符串", "", "localhost", true},
		{"无效IP-随意字符串", "not-an-ip", "localhost", true},
		{"无效IP-超出范围", "256.256.256.256", "localhost", true},
		{"无效IP-格式错", "192.168.1", "localhost", true},

		// 有效的主机名
		{"有效主机名-简单", "127.0.0.1", "test", false},
		{"有效主机名-带点", "127.0.0.1", "test.local", false},
		{"有效主机名-多级", "127.0.0.1", "api.dev.example.com", false},
		{"有效主机名-带连字符", "127.0.0.1", "my-host.local", false},

		// 无效的主机名
		{"无效主机名-空", "127.0.0.1", "", true},
		{"无效主机名-以连字符开头", "127.0.0.1", "-badhost", true},
		{"无效主机名-以连字符结尾", "127.0.0.1", "badhost-", true},
		{"无效主机名-特殊字符", "127.0.0.1", "bad@host", true},
		{"无效主机名-下划线", "127.0.0.1", "bad_host", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := NewHostsEntry(tt.ip, tt.host, "")
			err := entry.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestHostsEntryToHostsLine 测试转换为 hosts 文件行格式
func TestHostsEntryToHostsLine(t *testing.T) {
	t.Run("带注释的条目", func(t *testing.T) {
		entry := NewHostsEntry("127.0.0.1", "local.dev", "本地开发")
		line := entry.ToHostsLine()
		expected := "127.0.0.1\tlocal.dev\t# 本地开发"
		if line != expected {
			t.Errorf("ToHostsLine() = '%s'，期望 '%s'", line, expected)
		}
	})

	t.Run("无注释的条目", func(t *testing.T) {
		entry := NewHostsEntry("192.168.1.1", "server.local", "")
		line := entry.ToHostsLine()
		expected := "192.168.1.1\tserver.local"
		if line != expected {
			t.Errorf("ToHostsLine() = '%s'，期望 '%s'", line, expected)
		}
	})

	t.Run("禁用的条目返回空", func(t *testing.T) {
		entry := NewHostsEntry("127.0.0.1", "disabled.local", "")
		entry.Enabled = false
		line := entry.ToHostsLine()
		if line != "" {
			t.Errorf("禁用的条目应返回空字符串，实际为 '%s'", line)
		}
	})

	t.Run("IPv6 条目", func(t *testing.T) {
		entry := NewHostsEntry("::1", "ipv6.local", "IPv6测试")
		line := entry.ToHostsLine()
		expected := "::1\tipv6.local\t# IPv6测试"
		if line != expected {
			t.Errorf("ToHostsLine() = '%s'，期望 '%s'", line, expected)
		}
	})
}

// TestHostsEntryValidateErrors 测试返回的具体错误类型
func TestHostsEntryValidateErrors(t *testing.T) {
	t.Run("无效IP返回ErrInvalidIP", func(t *testing.T) {
		entry := NewHostsEntry("bad-ip", "localhost", "")
		err := entry.Validate()
		if err == nil {
			t.Fatal("期望返回错误")
		}
		domErr, ok := err.(*DomainError)
		if !ok {
			t.Fatal("期望返回 DomainError 类型")
		}
		if domErr.Code != "INVALID_IP" {
			t.Errorf("错误码应为 'INVALID_IP'，实际为 '%s'", domErr.Code)
		}
	})

	t.Run("无效主机名返回ErrInvalidHostname", func(t *testing.T) {
		entry := NewHostsEntry("127.0.0.1", "", "")
		err := entry.Validate()
		if err == nil {
			t.Fatal("期望返回错误")
		}
		domErr, ok := err.(*DomainError)
		if !ok {
			t.Fatal("期望返回 DomainError 类型")
		}
		if domErr.Code != "INVALID_HOSTNAME" {
			t.Errorf("错误码应为 'INVALID_HOSTNAME'，实际为 '%s'", domErr.Code)
		}
	})
}

// TestIsValidHostname 测试主机名验证的各种边界情况
func TestIsValidHostname(t *testing.T) {
	tests := []struct {
		name     string
		hostname string
		valid    bool
	}{
		// 有效的
		{"简单名称", "test", true},
		{"带数字", "test123", true},
		{"全限定域名", "api.example.com", true},
		{"63字符标签", strings.Repeat("a", 63), true},
		// 无效的
		{"空字符串", "", false},
		{"以连字符开头", "-test", false},
		{"以连字符结尾", "test-", false},
		{"包含下划线", "test_host", false},
		{"包含特殊字符", "test@host", false},
		{"超长主机名(>253)", string(make([]byte, 254)), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidHostname(tt.hostname)
			if result != tt.valid {
				t.Errorf("isValidHostname(%q) = %v，期望 %v", tt.hostname, result, tt.valid)
			}
		})
	}
}

// TestErrInvalidIPAndHostname 测试错误实现
func TestErrInvalidIPAndHostname(t *testing.T) {
	t.Run("ErrInvalidIP", func(t *testing.T) {
		if ErrInvalidIP.Error() != "IP 地址格式无效" {
			t.Errorf("ErrInvalidIP.Error() = '%s'", ErrInvalidIP.Error())
		}
	})

	t.Run("ErrInvalidHostname", func(t *testing.T) {
		if ErrInvalidHostname.Error() != "主机名格式无效" {
			t.Errorf("ErrInvalidHostname.Error() = '%s'", ErrInvalidHostname.Error())
		}
	})
}
