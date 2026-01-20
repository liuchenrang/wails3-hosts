# Project Context

## Purpose
一个类似 SwitchHosts 的跨平台 hosts 文件管理工具，采用 DDD（领域驱动设计）架构开发。
提供分组管理、条目编辑、版本历史、sudo 权限管理、国际化支持等功能。

## Tech Stack

### 后端
- **Go 1.24** + **Wails v3** - 跨平台桌面应用框架
- **JSON 存储** - 配置和版本历史使用 JSON 文件持久化
- **依赖管理** - Go Modules (go.mod)

### 前端
- **React 18** + **Vite** - 现代化前端构建工具
- **TypeScript** - 类型安全
- **Tailwind CSS** + **shadcn/ui** - UI 组件库
- **i18next** - 国际化支持（中文、英文、日文）
- **Lucide React** - 图标库

### 开发工具
- **Make** - 构建自动化
- **Task** (可选) - 任务运行器
- **npm scripts** - 前端包管理

## Project Conventions

### Code Style

#### Go 代码规范
- **命名约定**:
  - 包名: 小写单词，不使用下划线
  - 接口: 以 `er` 结尾或描述行为（如 `HostsRepository`）
  - 实体: 使用 `New` 前缀的工厂函数创建（如 `NewHostsGroup()`）
  - 常量: 驼峰命名
- **注释语言**: 中文注释
- **注释要求**:
  - 所有公共函数/方法必须添加注释说明职责
  - 复杂逻辑必须添加行内注释解释
  - 注释中标注遵循的设计原则（如 `// 单一职责:`、`// DRY:`）
- **错误处理**:
  - 使用自定义领域错误类型（`DomainError`）
  - 错误消息使用中文
  - 错误包装使用 `fmt.Errorf` 配合 `%w`
- **代码组织**:
  - 每个文件只包含一个主要类型
  - 文件名与类型名对应
  - 导入顺序: 标准库 -> 第三方库 -> 项目内部

#### React/TypeScript 代码规范
- **组件命名**: PascalCase（如 `Sidebar.tsx`）
- **文件组织**:
  - `components/ui/` - 通用 UI 组件
  - `components/layout/` - 布局组件
  - `components/hosts/` - 业务组件
  - `hooks/` - 自定义 Hooks
  - `api/` - API 调用封装
  - `i18n/` - 国际化资源
  - `utils/` - 工具函数
- **样式约定**:
  - 使用 `cn()` 工具函数合并 Tailwind 类名（来自 `tailwind-merge`）
  - 优先使用 Tailwind 工具类而非内联样式
  - 使用 shadcn/ui 组件保持一致性
- **类型定义**:
  - 所有 props 必须定义接口
  - 使用 TypeScript 类型推断
- **注释**: 与现有代码库保持一致（检测中英文）

### Architecture Patterns

#### DDD 四层架构

```
internal/
├── domain/          # 领域层：核心业务逻辑
│   ├── entity/          # 实体（如 HostsGroup, HostsEntry）
│   ├── valueobject/     # 值对象（如 HostsContent, SudoCredentials）
│   ├── repository/      # 仓储接口（抽象）
│   └── service/         # 领域服务（复杂业务逻辑）
├── application/     # 应用层：用例编排
│   ├── service/         # 应用服务（协调领域对象）
│   └── dto/             # 数据传输对象
├── infrastructure/  # 基础设施层：技术实现
│   ├── persistence/     # 持久化实现
│   └── system/          # 系统操作（hosts 文件、sudo）
└── interface/       # 接口层：外部交互
    └── handler/         # Wails 服务处理器
```

#### 设计原则（严格执行）
- **SOLID**:
  - **S**: 每个模块单一职责，明确注释
  - **O**: 通过仓储接口实现扩展性
  - **L**: 子类型可替换父类型
  - **I**: 接口专一，避免胖接口
  - **D**: 依赖抽象（仓储接口）而非具体实现
- **DRY**: 通过领域服务、仓储模式、工厂函数消除重复
- **KISS**: 简单直接的实现，避免过度设计
- **YAGNI**: 仅实现当前所需功能

#### 分层职责
1. **领域层** (Domain): 纯业务逻辑，不依赖外部
2. **应用层** (Application): 编排用例，依赖领域和基础设施接口
3. **基础设施层** (Infrastructure): 实现仓储接口，处理技术细节
4. **接口层** (Interface): 暴露给前端的 Wails 服务

#### 依赖方向
Interface → Application → Domain ← Infrastructure

### Testing Strategy

#### 测试文件
- 测试文件与源文件同目录，命名: `*_test.go`
- 示例: `hosts_domain_service_test.go`

#### 测试覆盖
- 领域服务: 必须有单元测试
- 仓储实现: 必须有单元测试
- 系统操作: 可选集成测试

#### 运行测试
```bash
make test              # 运行所有测试
make test-cov          # 生成覆盖率报告
```

### Git Workflow

#### 分支策略
- `main`: 主分支，稳定版本
- 功能分支: 基于 `main` 创建，完成后合并回 `main`

#### 提交规范
- **提交消息格式**: 简洁的中文描述
- **示例**:
  - `实现分组拖动排序功能`
  - `修复 Modal title 类型支持`
  - `文档整理：移动项目文档到 doc 目录`

#### 提交前检查
- 运行 `make lint` 进行代码检查
- 运行 `make test` 确保测试通过
- ⚠️ **重要**: 仅在用户明确要求时执行 git 操作

### 文件组织规范

#### Go 目录结构
```
internal/domain/
├── entity/           # 实体定义
├── valueobject/      # 值对象
├── repository/       # 仓储接口
└── service/          # 领域服务

internal/application/
├── service/          # 应用服务
└── dto/              # DTO 定义

internal/infrastructure/
├── persistence/      # 持久化实现
└── system/           # 系统操作

internal/interface/
└── handler/          # Wails 服务处理器
```

#### 前端目录结构
```
frontend/src/
├── components/
│   ├── ui/          # 通用 UI 组件（Button, Modal, Input 等）
│   ├── layout/      # 布局组件（Sidebar, MainPanel）
│   └── hosts/       # 业务组件（ConflictAlert, VersionHistory）
├── hooks/           # 自定义 Hooks
├── api/             # API 调用封装
├── i18n/            # 国际化资源
│   ├── zh/          # 中文
│   ├── en/          # 英文
│   └── ja/          # 日文
├── theme/           # 主题配置
├── types/           # TypeScript 类型定义
└── utils/           # 工具函数（cn, formatDate 等）
```

## Domain Context

### 核心领域概念
- **HostsGroup**: hosts 配置分组，包含多个条目
- **HostsEntry**: 单个 hosts 条目（IP + 主机名 + 注释）
- **HostsVersion**: 版本历史记录
- **SudoCredentials**: sudo 权限凭证（缓存机制）

### 业务规则
1. 分组可以启用/禁用，禁用的分组不生效
2. 条目可以单独启用/禁用
3. 应用 hosts 配置需要 sudo 权限
4. 密码缓存 5 分钟，避免重复输入
5. 支持拖动排序分组
6. 自动保存版本历史

### 用户工作流
1. 创建分组 → 添加条目 → 启用分组 → 应用配置
2. 或: 导入系统 hosts → 编辑 → 应用
3. 回滚: 版本历史 → 选择版本 → 回滚

## Important Constraints

### 技术约束
- **平台支持**: macOS, Linux, Windows（跨平台）
- **Go 版本**: 1.24+
- **Node.js 版本**: 18+
- **Wails 版本**: v3
- **sudo 权限**: 修改系统 hosts 文件需要管理员权限

### 安全约束
- sudo 密码仅在内存中缓存，不写入文件
- 密码缓存时间: 5 分钟
- hosts 文件路径:
  - macOS/Linux: `/etc/hosts`
  - Windows: `C:\Windows\System32\drivers\etc\hosts`

### 配置文件位置
- **macOS**: `~/Library/Application Support/hosts-manager/`
- **Linux**: `~/.config/hosts-manager/`
- **Windows**: `%APPDATA%\hosts-manager\`

## External Dependencies

### 系统依赖
- **sudo 命令**: macOS/Linux 权限提升
- **系统 hosts 文件**: 核心操作对象

### Go 依赖（主要）
- `github.com/google/uuid` - UUID 生成
- `github.com/wailsio/wails/v3` - Wails 框架

### 前端依赖（主要）
- `@wailsio/runtime` - Wails 运行时
- `react` + `react-dom` - React 框架
- `tailwindcss` - CSS 框架
- `i18next` + `react-i18next` - 国际化
- `lucide-react` - 图标库
- `clsx` + `tailwind-merge` - 类名工具

## 开发指南

### 快速开始
```bash
make quick         # 安装依赖 + 启动开发模式
make dev           # 仅启动开发模式
make build         # 构建应用
```

### 代码质量检查
```bash
make lint          # Go 代码检查
make test          # 运行测试
make format        # 格式化代码
```

### 常用命令
- `wails3 dev` - 开发模式（热重载）
- `wails3 build` - 生产构建
- `make help` - 显示所有可用命令

### 新增功能开发流程
1. 在 `internal/domain/entity/` 定义实体
2. 在 `internal/domain/repository/` 定义仓储接口
3. 在 `internal/domain/service/` 实现领域服务
4. 在 `internal/infrastructure/persistence/` 实现仓储
5. 在 `internal/application/service/` 编写应用服务
6. 在 `internal/interface/handler/` 暴露 Wails 服务
7. 在 `frontend/src/api/` 封装 API 调用
8. 在 `frontend/src/components/` 实现 UI 组件

### 代码审查检查点
- [ ] 遵循 SOLID 原则
- [ ] 中文注释完整
- [ ] 函数有职责说明
- [ ] 错误处理完善
- [ ] 类型安全（TypeScript）
- [ ] 国际化支持（i18n）
- [ ] 测试覆盖核心逻辑

## 重要注意事项

### Docker 环境
根据全局配置，新功能需要在 Docker 环境中测试（TODO: 完善 Docker 配置）

### 性能考虑
- JSON 文件读写需要优化（大量条目时）
- 版本历史应限制最大数量
- 前端使用虚拟滚动（长列表优化）

### 待优化项
- [ ] Docker 测试环境
- [ ] 远程同步（GitHub Gist）
- [ ] 导入/导出功能
- [ ] 系统托盘图标
- [ ] 命令行接口
