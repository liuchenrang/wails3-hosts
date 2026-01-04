package system

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// SudoCommand sudo 命令封装
// 单一职责: 封装需要 sudo 权限的命令执行
type SudoCommand struct {
	cmd     *exec.Cmd
	password string
	timeout time.Duration
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

// SetStdin 设置标准输入
func (c *SudoCommand) SetStdin(data []byte) {
	c.cmd.Stdin = bytes.NewReader(data)
}

// setTimeout 设置超时
func (c *SudoCommand) setTimeout(timeout time.Duration) {
	c.timeout = timeout
}

// Run 执行命令
func (c *SudoCommand) Run() error {
	// 设置标准输入（密码 + 换行）
	stdin := strings.NewReader(c.password + "\n")
	c.cmd.Stdin = stdin

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
