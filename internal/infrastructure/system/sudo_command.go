package system

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"time"
)

// SudoCommand sudo 命令封装
// 单一职责: 封装需要 sudo 权限的命令执行
type SudoCommand struct {
	cmd      *exec.Cmd
	password string
	stdin    io.Reader
	timeout  time.Duration
}

// NewSudoCommand 创建 sudo 命令
func NewSudoCommand(args []string) *SudoCommand {
	cmd := exec.Command("sudo", "-S", args[0])
	if len(args) > 1 {
		cmd.Args = append(cmd.Args, args[1:]...)
	}

	return &SudoCommand{
		cmd:     cmd,
		timeout: 30 * time.Second,
	}
}

// SetPassword 设置 sudo 密码
func (c *SudoCommand) SetPassword(password string) {
	c.password = password
}

// SetStdin 设置标准输入（用于传递给命令的内容）
func (c *SudoCommand) SetStdin(data []byte) {
	c.stdin = bytes.NewReader(data)
}

// setTimeout 设置超时
func (c *SudoCommand) setTimeout(timeout time.Duration) {
	c.timeout = timeout
}

// Run 执行命令
func (c *SudoCommand) Run() error {
	// 构建标准输入：密码 + 换行 + 实际内容
	// sudo -S 会先从 stdin 读取密码，剩余内容会传递给实际执行的命令
	var stdinContent bytes.Buffer

	// 首先写入密码（sudo -S 会读取第一行作为密码）
	if c.password != "" {
		stdinContent.WriteString(c.password)
		stdinContent.WriteString("\n")
	}

	// 然后写入实际要传递给命令的内容
	if c.stdin != nil {
		stdinContent.ReadFrom(c.stdin)
	}

	c.cmd.Stdin = &stdinContent

	// 捕获输出
	var stdout, stderr bytes.Buffer
	c.cmd.Stdout = &stdout
	c.cmd.Stderr = &stderr

	// 启动命令
	if err := c.cmd.Start(); err != nil {
		return fmt.Errorf("启动命令失败: %w", err)
	}

	// 设置超时
	done := make(chan error, 1)
	go func() {
		done <- c.cmd.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("命令执行失败: %w, stderr: %s", err, stderr.String())
		}
		return nil
	case <-time.After(c.timeout):
		c.cmd.Process.Kill()
		return fmt.Errorf("命令执行超时")
	}
}
