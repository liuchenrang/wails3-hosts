#!/bin/bash

# Task 安装脚本
# 用于安装 Task 任务运行器

set -e

echo "=========================================="
echo "Task 安装脚本"
echo "=========================================="
echo ""

# 检测操作系统
OS="$(uname -s)"
echo "检测到操作系统: $OS"

case "$OS" in
  Darwin)
    echo ""
    echo "在 macOS 上安装 Task..."
    echo ""
    echo "请选择安装方式:"
    echo "1) 使用 Homebrew (推荐)"
    echo "2) 使用 Go 安装"
    echo "3) 手动安装"
    echo ""
    read -p "请输入选项 (1-3): " choice

    case $choice in
      1)
        if command -v brew &> /dev/null; then
          echo "正在使用 Homebrew 安装 Task..."
          brew install go-task/tap/go-task
          echo "✅ Task 安装完成!"
          task --version
        else
          echo "❌ 错误: 未找到 Homebrew"
          echo "请先安装 Homebrew: https://brew.sh"
          exit 1
        fi
        ;;
      2)
        if command -v go &> /dev/null; then
          echo "正在使用 Go 安装 Task..."
          go install github.com/go-task/task/v3/cmd/task@latest
          echo "✅ Task 安装完成!"
          echo "请确保 \$GOPATH/bin 在你的 PATH 中"
          export PATH=$PATH:$(go env GOPATH)/bin
          task --version
        else
          echo "❌ 错误: 未找到 Go"
          echo "请先安装 Go: https://go.dev"
          exit 1
        fi
        ;;
      3)
        echo "请手动安装 Task:"
        echo "1. 访问: https://taskfile.dev/installation/"
        echo "2. 下载适合你系统的二进制文件"
        echo "3. 将其添加到 PATH"
        ;;
      *)
        echo "无效的选项"
        exit 1
        ;;
    esac
    ;;

  Linux)
    echo ""
    echo "在 Linux 上安装 Task..."
    echo ""
    echo "请选择安装方式:"
    echo "1) 使用包管理器"
    echo "2) 使用 Go 安装"
    echo "3) 手动安装"
    echo ""
    read -p "请输入选项 (1-3): " choice

    case $choice in
      1)
        # 尝试检测包管理器
        if command -v apt-get &> /dev/null; then
          echo "正在使用 apt 安装 Task..."
          sudo apt-get update
          sudo apt-get install -y task
        elif command -v yum &> /dev/null; then
          echo "正在使用 yum 安装 Task..."
          sudo yum install -y task
        elif command -v pacman &> /dev/null; then
          echo "正在使用 pacman 安装 Task..."
          sudo pacman -S task
        else
          echo "❌ 无法检测包管理器"
          echo "请尝试其他安装方式"
          exit 1
        fi
        ;;
      2)
        if command -v go &> /dev/null; then
          echo "正在使用 Go 安装 Task..."
          go install github.com/go-task/task/v3/cmd/task@latest
          echo "✅ Task 安装完成!"
          export PATH=$PATH:$(go env GOPATH)/bin
          task --version
        else
          echo "❌ 错误: 未找到 Go"
          echo "请先安装 Go: https://go.dev"
          exit 1
        fi
        ;;
      3)
        echo "请手动安装 Task:"
        echo "1. 访问: https://taskfile.dev/installation/"
        echo "2. 下载适合你系统的二进制文件"
        echo "3. 将其添加到 PATH"
        ;;
      *)
        echo "无效的选项"
        exit 1
        ;;
    esac
    ;;

  *)
    echo "❌ 不支持的操作系统: $OS"
    echo "请手动安装 Task: https://taskfile.dev/installation/"
    exit 1
    ;;
esac

echo ""
echo "=========================================="
echo "安装完成!"
echo "=========================================="
echo ""
echo "验证安装:"
task --version || echo "⚠️  请重新加载终端或重启系统"
echo ""
echo "开始使用:"
echo "  task          # 查看所有可用命令"
echo "  task dev      # 启动开发模式"
echo "  task build    # 构建应用"
echo ""
