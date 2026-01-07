package system

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// TestSudoCommandStdin 测试 sudo 命令的 stdin 处理
func TestSudoCommandStdin(t *testing.T) {
	t.Run("验证stdin内容构建", func(t *testing.T) {
		cmd := NewSudoCommand([]string{"echo", "test"})

		// 设置密码
		testPassword := "mysecretpassword"
		cmd.SetPassword(testPassword)

		// 设置要传递给命令的内容
		testContent := "127.0.0.1\tlocalhost\n127.0.0.1\ttest.local"
		cmd.SetStdin([]byte(testContent))

		// 检查内部状态
		if cmd.password != testPassword {
			t.Errorf("密码设置错误: got %v, want %v", cmd.password, testPassword)
		}

		if cmd.stdin == nil {
			t.Error("stdin 不应为 nil")
		}
	})

	t.Run("验证密码和内容分离", func(t *testing.T) {
		// 这个测试验证密码不会出现在要写入的内容中
		testPassword := "testpass123"
		testContent := "127.0.0.1\texample.com"

		cmd := NewSudoCommand([]string{"tee", "/tmp/test_hosts"})
		cmd.SetPassword(testPassword)
		cmd.SetStdin([]byte(testContent))

		// 验证密码不会与内容混合
		// 我们不实际执行命令，只验证数据设置正确
		if cmd.password != testPassword {
			t.Errorf("密码不匹配: got %v, want %v", cmd.password, testPassword)
		}

		// 验证 stdin 中不包含密码
		buf := new(bytes.Buffer)
		buf.ReadFrom(cmd.stdin)
		stdinContent := buf.String()

		if strings.Contains(stdinContent, testPassword) {
			t.Error("安全错误: 密码出现在 stdin 内容中!")
		}

		if stdinContent != testContent {
			t.Errorf("stdin 内容不匹配: got %v, want %v", stdinContent, testContent)
		}
	})

	t.Run("验证空内容处理", func(t *testing.T) {
		cmd := NewSudoCommand([]string{"echo", "test"})
		cmd.SetPassword("password")
		// 不设置 stdin

		if cmd.password != "password" {
			t.Error("密码设置失败")
		}
	})

	t.Run("验证空密码处理", func(t *testing.T) {
		cmd := NewSudoCommand([]string{"echo", "test"})
		// 不设置密码
		testContent := "test content"
		cmd.SetStdin([]byte(testContent))

		if cmd.stdin == nil {
			t.Error("stdin 不应为 nil")
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(cmd.stdin)
		stdinContent := buf.String()

		if stdinContent != testContent {
			t.Errorf("stdin 内容不匹配: got %v, want %v", stdinContent, testContent)
		}
	})

	t.Run("验证Run方法stdin构建", func(t *testing.T) {
		// 测试 Run 方法中 stdin 的正确构建
		testPassword := "mypassword"
		testContent := "# Test hosts file\n127.0.0.1 localhost"

		// 创建一个模拟命令来检查 stdin 构建
		// 注意：这里我们只测试数据准备，不实际执行 sudo 命令
		cmd := NewSudoCommand([]string{"cat"})
		cmd.SetPassword(testPassword)
		cmd.SetStdin([]byte(testContent))

		// 验证数据结构
		if cmd.password != testPassword {
			t.Error("密码未正确设置")
		}

		if cmd.stdin == nil {
			t.Fatal("stdin 未设置")
		}

		// 读取 stdin 内容
		contentBuf := new(bytes.Buffer)
		io.Copy(contentBuf, cmd.stdin)
		actualContent := contentBuf.String()

		// 验证内容不包含密码
		if strings.Contains(actualContent, testPassword) {
			t.Errorf("安全错误: 密码泄露到内容中! 内容: %s", actualContent)
		}

		// 验证内容正确
		if actualContent != testContent {
			t.Errorf("内容不匹配\n期望: %s\n实际: %s", testContent, actualContent)
		}
	})
}

