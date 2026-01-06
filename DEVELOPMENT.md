# Vibe Kanban 开发指南

## 📋 前置要求

在开始开发之前,请确保已安装以下工具:

- **Go** 1.21+ (后端开发)
- **Node.js** 18+ (前端开发)
- **npm** 或 **yarn** (包管理)
- **Wails CLI v3** (应用框架)
- **Task** (任务运行器,推荐) 或 Make
- **Docker** (可选,用于跨平台构建)

### 安装 Task

```bash
# macOS/Linux
brew install go-task/tap/go-task

# 或使用 Go 安装
go install github.com/go-task/task/v3/cmd/task@latest

# 验证安装
task --version
```

## 🚀 快速开始

### 方式一:使用 npm 脚本 (推荐新手)

```bash
# 1. 安装依赖并启动开发模式
npm run quick

# 2. 或者只启动开发模式
npm run dev
```

### 方式二:使用 Task (推荐)

```bash
# 1. 查看所有可用命令
task

# 2. 快速启动
task quick

# 3. 启动开发模式
task dev
```

### 方式三:使用 Make

```bash
# 查看所有可用命令
make help

# 快速启动
make quick

# 启动开发模式
make dev
```

## 📚 常用命令

### 开发相关

```bash
# 启动开发模式 (热重载)
task dev
# 或
npm run dev

# 仅启动前端开发服务器
task dev:frontend
# 或
npm run dev:frontend

# 清理缓存并重新启动
task dev:clean
# 或
npm run dev:clean
```

### 构建相关

```bash
# 构建当前平台应用
task build
# 或
npm run build

# 构建开发版本 (包含调试信息)
task build:dev
# 或
npm run build:dev

# 构建生产版本 (优化)
task build:prod
# 或
npm run build:prod

# 跨平台构建 (需要 Docker)
task build:all
# 或
npm run build:all
```

### 测试相关

```bash
# 运行 Go 测试
task go:test
# 或
npm run test

# 生成测试覆盖率报告
task go:test:coverage
# 或
npm run test:coverage

# 代码检查
task lint
# 或
npm run lint

# 格式化代码
task format
# 或
npm run format
```

### 清理相关

```bash
# 清理构建产物
task clean
# 或
npm run clean

# 深度清理 (包括缓存)
task clean:all
# 或
npm run clean:all
```

### 打包相关

```bash
# 打包当前平台应用
task package
# 或
npm run package

# 跨平台打包
task package:all
# 或
npm run package:all
```

## 🔧 开发工作流

### 1. 首次设置

```bash
# 克隆项目后
git clone <repository-url>
cd vibe-kanban

# 快速启动 (自动安装依赖)
task quick
```

### 2. 日常开发

```bash
# 启动开发模式
task dev

# 在另一个终端窗口中,可以单独操作前端
cd frontend
npm run dev
```

### 3. 提交代码前

```bash
# 格式化代码
task format

# 运行测试
task test

# 代码检查
task lint
```

### 4. 构建生产版本

```bash
# 清理旧构建
task clean

# 构建生产版本
task build:prod

# 运行构建的应用
task run
```

## 🐳 Docker 跨平台构建

### 设置 Docker 环境

```bash
# 首次使用,构建 Docker 镜像 (需要下载约 800MB)
task setup:docker
```

### 跨平台构建

```bash
# 构建所有平台版本
task build:all

# 或者单独构建特定平台
task darwin:build    # macOS
task windows:build   # Windows
task linux:build     # Linux
```

## 🔗 前端绑定生成

当修改了 Go 后端服务接口时,需要重新生成前端绑定:

```bash
# 生成绑定
task bindings

# 清理并重新生成
task bindings:clean
```

## 📦 项目结构

```
vibe-kanban/
├── frontend/              # 前端代码 (React + Vite)
│   ├── src/              # 源代码
│   ├── dist/             # 构建产物
│   ├── bindings/         # Go 绑定 (自动生成)
│   └── package.json      # 前端依赖
├── internal/             # Go 后端代码
│   ├── application/      # 应用服务层
│   ├── domain/           # 领域层
│   ├── infrastructure/   # 基础设施层
│   └── interface/        # 接口层
├── build/                # 构建配置
│   ├── config.yml        # Wails 配置
│   └── Taskfile.yml      # 构建任务
├── main.go               # 应用入口
├── Taskfile.yml          # 主任务文件
├── Makefile              # Make 任务
├── package.json          # npm 脚本
└── go.mod                # Go 依赖
```

## 🎯 环境变量

### Vite 开发服务器端口

默认端口: 9245

```bash
# 自定义端口
WAILS_VITE_PORT=3000 task dev

# 或使用 npm
WAILS_VITE_PORT=3000 npm run dev
```

## 🐛 调试

### 前端调试

开发模式下,前端运行在 `http://localhost:9245`

### 后端调试

使用开发模式 (`task dev:clean`) 会保留调试信息

### 查看日志

```bash
# 开发模式下日志会直接输出到终端
task dev

# 查看 Wails 日志
tail -f /tmp/wails-logs/*.log
```

## 🔍 常见问题

### 1. 依赖安装失败

```bash
# 清理并重新安装
task clean:all
task quick
```

### 2. 前端绑定问题

```bash
# 重新生成绑定
task bindings:clean
```

### 3. 构建失败

```bash
# 深度清理并重试
task clean:all
task build:dev
```

### 4. 端口被占用

```bash
# 使用其他端口
WAILS_VITE_PORT=3000 task dev
```

## 📖 更多资源

- [Wails v3 文档](https://v3.wails.io)
- [Task 文档](https://taskfile.dev)
- [React 文档](https://react.dev)
- [Vite 文档](https://vitejs.dev)

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

提交前请确保:
- 代码已格式化 (`task format`)
- 测试通过 (`task test`)
- 代码检查通过 (`task lint`)

## 📄 许可证

MIT License
