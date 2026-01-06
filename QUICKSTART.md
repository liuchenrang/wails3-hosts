# 🚀 Vibe Kanban 快速启动指南

## ⚡ 30秒快速启动

```bash
# 方式一:使用 Make (无需额外安装)
make quick

# 方式二:使用 npm (需要 Node.js)
npm run quick

# 方式三:使用 wails3 dev (需要 Wails CLI)
wails3 dev
```

## 📦 前置要求

### 必需
- ✅ **Go** 1.21+ - [下载](https://go.dev/dl/)
- ✅ **Node.js** 18+ - [下载](https://nodejs.org)
- ✅ **Wails CLI v3** - `go install github.com/wailsapp/wails/v3/cmd/wails3@latest`

### 可选 (推荐)
- 🔧 **Task** - 更强大的任务运行器
  ```bash
  # macOS
  brew install go-task/tap/go-task

  # 使用 Go
  go install github.com/go-task/task/v3/cmd/task@latest

  # 或使用提供的安装脚本
  bash scripts/install-task.sh
  ```

## 🎯 开发命令对比

| 功能 | Make | npm | Task (需要安装) |
|------|------|-----|-----------------|
| 开发模式 | `make dev` | `npm run dev` | `task dev` |
| 快速启动 | `make quick` | `npm run quick` | `task quick` |
| 构建 | `make build` | `npm run build` | `task build` |
| 测试 | `make test` | `npm run test` | `task test` |
| 清理 | `make clean` | `npm run clean` | `task clean` |

**推荐**: 使用 `make` 命令,因为它已经在你的系统上可用!

## 💡 常用开发流程

### 1. 首次克隆项目
```bash
git clone <repository-url>
cd vibe-kanban
make quick  # 自动安装所有依赖并启动
```

### 2. 日常开发
```bash
make dev  # 启动热重载开发模式
```

### 3. 构建生产版本
```bash
make clean      # 清理旧构建
make build      # 构建应用
make package    # 打包应用
```

## 🔧 开发模式特性

- ✅ **热重载**: 前端代码修改自动刷新
- ✅ **自动编译**: Go 代码修改自动重新编译
- ✅ **开发工具**: 浏览器开发者工具
- ✅ **绑定生成**: 自动生成前端绑定

## 🐛 常见问题

### 端口被占用
```bash
# 使用自定义端口
WAILS_VITE_PORT=3000 make dev
```

### 依赖问题
```bash
make clean-all  # 深度清理
make quick      # 重新安装并启动
```

### 前端绑定错误
```bash
make bindings:clean  # 清理并重新生成
```

## 📚 详细文档

查看完整的开发文档: [DEVELOPMENT.md](./DEVELOPMENT.md)

## 🎉 开始开发

```bash
# 现在就开始!
make dev
```

访问: http://localhost:9245

---

**提示**: 如果遇到问题,运行 `make info` 查看系统环境信息。
