# Hosts Manager

一个类似 SwitchHosts 的跨平台 hosts 文件管理工具，采用 DDD（领域驱动设计）架构开发。

## 功能特性

- ✅ **分组管理**: 左侧分组列表，支持创建、编辑、删除、启用/禁用分组
- ✅ **条目管理**: 右侧 hosts 条目列表，支持添加、编辑、删除条目
- ✅ **快捷键保存**: `Cmd+S` (Mac) / `Ctrl+S` (Windows/Linux) 快速应用配置
- ✅ **Sudo 权限管理**: 智能缓存 sudo 密码，避免重复输入
- ✅ **版本历史**: 自动保存历史版本，支持一键回滚
- ✅ **国际化**: 支持中文、英文、日文，根据系统语言自动选择
- ✅ **主题切换**: 支持明亮和暗色两种主题

## 技术栈

### 后端
- **Go 1.24** + **Wails v3**
- **DDD 架构**: 领域层、应用层、基础设施层、接口层清晰分离
- **JSON 存储**: 配置和版本历史使用 JSON 文件持久化

### 前端
- **React 18** + **Vite**
- **Tailwind CSS** + **shadcn/ui**: 现代化的 UI 组件
- **i18next**: 完整的国际化支持
- **Lucide React**: 精美的图标库

## 架构设计

### DDD 分层架构

```
wails3-demo/
├── internal/
│   ├── domain/          # 领域层：实体、值对象、仓储接口、领域服务
│   ├── application/     # 应用层：应用服务、DTO
│   ├── infrastructure/  # 基础设施层：持久化、系统操作
│   └── interface/       # 接口层：Wails 服务处理器
└── frontend/
    └── src/
        ├── components/  # React 组件
        ├── hooks/       # 自定义 Hooks
        ├── i18n/        # 国际化资源
        └── theme/       # 主题配置
```

### 设计原则

- **SOLID**: 所有模块遵循单一职责、开闭原则、里氏替换、接口隔离、依赖倒置
- **DRY**: 通过领域服务、仓储模式消除代码重复
- **KISS**: 简单直接的实现，避免过度设计
- **YAGNI**: 仅实现当前所需功能，不预留未来特性

## 快速开始

### 前置要求

- Go 1.24+
- Node.js 18+
- Wails v3 CLI

### 安装

```bash
# 安装前端依赖
cd frontend
npm install

# 返回项目根目录
cd ..
```

### 开发模式

```bash
# 启动开发服务器（支持热重载）
wails3 dev
```

### 生产构建

```bash
# 构建可执行文件
wails3 build
```

构建产物位于 `build/bin/` 目录。

## 使用说明

### 1. 创建分组

点击左侧边栏的"新建分组"按钮，输入分组名称和描述。

### 2. 添加 hosts 条目

选择一个分组后，在右侧面板点击"添加条目"，填写：
- **IP 地址**: 例如 `127.0.0.1`
- **主机名**: 例如 `localhost.local`
- **注释**: 可选的说明文字

### 3. 应用配置

- 方法 1: 点击"应用配置"按钮
- 方法 2: 使用快捷键 `Cmd+S` (Mac) 或 `Ctrl+S` (Windows/Linux)
- 首次应用需要输入 sudo 密码，密码将缓存 5 分钟

### 4. 版本回滚

点击顶部栏的"版本历史"按钮，选择历史版本进行回滚。

### 5. 切换主题

点击顶部栏的主题切换按钮，在明亮和暗色主题间切换。

## 配置文件

配置文件位置：
- **macOS**: `~/Library/Application Support/hosts-manager/`
- **Linux**: `~/.config/hosts-manager/`
- **Windows**: `%APPDATA%\hosts-manager\`

文件结构：
```
hosts-manager/
├── config.json      # 分组配置
├── versions.json    # 版本历史
└── backups/         # hosts 文件备份
```

## 开发指南

### 后端开发

1. **领域层** (`internal/domain/`): 定义实体、值对象、仓储接口
2. **应用层** (`internal/application/`): 编写应用服务，协调领域对象
3. **基础设施层** (`internal/infrastructure/`): 实现仓储接口、系统操作
4. **接口层** (`internal/interface/`): 暴露 Wails 服务给前端

### 前端开发

1. **组件** (`frontend/src/components/`): 编写 React 组件
2. **API** (`frontend/src/api/`): 封装 Wails API 调用
3. **Hooks** (`frontend/src/hooks/`): 自定义 React Hooks
4. **国际化** (`frontend/src/i18n/`): 添加多语言支持

### 代码规范

- 所有函数和组件添加注释说明职责
- 遵循 SOLID 原则，保持单一职责
- 公共函数需要参数验证和错误处理
- 使用 `cn()` 工具函数合并 Tailwind 类名

## 待办事项

- [ ] 连接 Wails v3 的 API 绑定
- [ ] 完善版本历史 UI
- [ ] 添加导入/导出功能
- [ ] 支持远程同步（GitHub Gist）
- [ ] 添加系统托盘图标
- [ ] 实现命令行接口
- [ ] Docker 环境测试

## 贡献指南

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 致谢

- [Wails](https://wails.io/) - 跨平台桌面应用框架
- [SwitchHosts](https://github.com/oldj/SwitchHosts) - 灵感来源
- [shadcn/ui](https://ui.shadcn.com/) - UI 组件设计参考
