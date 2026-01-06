# Vibe Kanban - Wails3 项目 Makefile
# 可以不依赖 Task 命令独立使用

.PHONY: help dev build test clean package info

# 默认目标 - 显示帮助
help:
	@echo "=========================================="
	@echo "Vibe Kanban - 可用命令"
	@echo "=========================================="
	@echo ""
	@echo "开发命令:"
	@echo "  make dev           - 启动开发模式 (热重载)"
	@echo "  make dev-clean     - 清理缓存并启动开发模式"
	@echo "  make quick         - 快速启动 (安装依赖 + 开发)"
	@echo ""
	@echo "构建命令:"
	@echo "  make build         - 构建应用 (当前平台)"
	@echo "  make build-dev     - 构建开发版本 (含调试信息)"
	@echo "  make build-prod    - 构建生产版本 (优化)"
	@echo ""
	@echo "测试命令:"
	@echo "  make test          - 运行 Go 测试"
	@echo "  make test-cov      - 生成测试覆盖率报告"
	@echo "  make lint          - 代码检查"
	@echo "  make format        - 格式化代码"
	@echo ""
	@echo "清理命令:"
	@echo "  make clean         - 清理构建产物"
	@echo "  make clean-all     - 深度清理 (包括缓存)"
	@echo ""
	@echo "打包命令:"
	@echo "  make package       - 打包应用 (当前平台)"
	@echo ""
	@echo "其他命令:"
	@echo "  make info          - 显示项目信息"
	@echo "  make bindings      - 生成前端绑定"
	@echo "  make bindings-clean - 清理并重新生成绑定"
	@echo ""
	@echo "=========================================="
	@echo "提示: 也可以使用 'npm run <command>'"
	@echo "如果安装了 Task,也可以使用 'task <task>'"
	@echo "=========================================="

# 开发命令
dev:
	wails3 dev -config ./build/config.yml -port 9245

dev-clean:
	echo "清理构建缓存..."
	rm -rf frontend/dist
	rm -rf frontend/bindings
	rm -rf bin
	echo "完成! 正在启动开发模式..."
	wails3 dev -config ./build/config.yml -port 9245

quick:
	echo "🚀 快速启动流程..."
	cd frontend && npm install
	go mod tidy
	wails3 dev -config ./build/config.yml -port 9245

# 构建命令
build:
	wails3 build

build-dev:
	wails3 build -debug -clean

build-prod:
	wails3 build

# 测试命令
test:
	go test -v ./...

test-cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	echo "覆盖率报告已生成: coverage.html"

lint:
	echo "检查 Go 代码..."
	go vet ./...
	echo "检查前端代码..."
	cd frontend && npm run lint || echo "前端 lint 未配置"

format:
	echo "格式化 Go 代码..."
	go fmt ./...
	echo "格式化前端代码..."
	cd frontend && npm run format || echo "前端 format 未配置"

# 清理命令
clean:
	rm -rf bin
	rm -rf frontend/dist
	rm -rf frontend/bindings
	rm -rf frontend/node_modules
	rm -f coverage.out coverage.html
	echo "清理完成!"

clean-all:
	rm -rf bin
	rm -rf frontend/dist
	rm -rf frontend/bindings
	rm -rf frontend/node_modules
	rm -rf frontend/.vite
	rm -f coverage.out coverage.html
	go clean -cache -testcache
	echo "深度清理完成!"

# 打包命令
package:
	@echo "请使用 wails3 build 或特定平台脚本"

# 其他命令
info:
	@echo "=========================================="
	@echo "📱 应用名称: Vibe Kanban"
	@echo "🔧 Go 版本: $$(go version)"
	@echo "📦 Node 版本: $$(node --version)"
	@echo "🎨 NPM 版本: $$(npm --version)"
	@echo "💻 操作系统: $$(uname -s)"
	@echo "🔌 Vite 端口: 9245"
	@echo "=========================================="

bindings:
	wails3 generate bindings -f -clean=true

bindings-clean:
	rm -rf frontend/bindings
	wails3 generate bindings -f -clean=true
