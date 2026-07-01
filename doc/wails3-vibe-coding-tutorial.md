# 零基础也能做！用 Wails 3 + Vibe Coding 5分钟极速打造轻量级桌面应用hosts文件管理

> 以 [wails3-hosts](https://github.com/chen/wails3-hosts)（一个跨平台 hosts 文件管理工具）为例，手把手教你用 AI 辅助（Vibe Coding）开发桌面应用。

---

## 什么是 Wails3？

Wails 是一个用 Go 语言写后端逻辑、用 Web 技术（React/Vue/Svelte 等）写界面的桌面应用框架。你可以把它理解为"Go 版本的 Electron"——但它比 Electron 小得多，性能好得多，打包出来的应用只有几 MB。

**Wails3 的核心思路：**

```
┌─────────────────────────────────────┐
│           桌面应用窗口               │
│  ┌──────────────┐ ┌──────────────┐  │
│  │  前端界面     │ │  Go 后端     │  │
│  │  React/TS    │◄──► 业务逻辑   │  │
│  │  Tailwind    │ │  系统操作    │  │
│  └──────────────┘ └──────────────┘  │
└─────────────────────────────────────┘
```

- 前端负责 UI，和写网页一样
- Go 后端负责读写文件、系统调用等"网页做不到"的事
- Wails 自动把 Go 函数生成为前端可以直接调用的 JS 函数

---

## 什么是 Vibe Coding？

Vibe Coding 就是"氛围编程"——你描述想要什么，AI（比如 Claude Code）来写代码，你来审查和调整方向。适合：

- 对技术有概念但不熟悉某个框架的人
- 想快速验证产品想法的人
- 想学习新技术栈的人

**工作流程：**

```
你描述需求 → AI 生成代码 → 你运行测试 → 你反馈问题 → AI 修改 → 重复
```

---

## 环境准备

### 1. 安装 Go

前往 https://go.dev/dl/ 下载安装，验证：

```bash
go version
# 输出: go version go1.21.x ...
```

### 2. 安装 Node.js

前往 https://nodejs.org 下载 LTS 版本，验证：

```bash
node --version  # v20.x.x
npm --version   # 10.x.x
```

### 3. 安装 Wails3 CLI

```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
wails3 --version
```

### 4. 安装 Task（任务运行器，可选但推荐）

```bash
# macOS
brew install go-task

# Windows (via Scoop)
scoop install task

# 验证
task --version
```

### 5. 安装 Claude Code（AI 助手）

```bash
npm install -g @anthropic-ai/claude-code
claude
```

---

## 创建第一个项目

### 用 Wails3 脚手架初始化

```bash
wails3 init -name my-app -template react-ts
cd my-app
```

这会生成这样的结构：

```
my-app/
├── main.go              # 应用入口，创建窗口和注册服务
├── go.mod               # Go 依赖管理
├── frontend/            # 前端代码
│   ├── src/
│   │   ├── App.tsx      # React 主组件
│   │   └── main.tsx     # 前端入口
│   ├── package.json
│   └── vite.config.js
└── build/               # 构建配置和资源
```

### 启动开发模式

```bash
# 安装前端依赖
cd frontend && npm install && cd ..

# 启动开发服务器（热重载）
wails3 dev
```

> 修改前端代码会即时刷新；修改 Go 代码会自动重新编译。

---

## 项目架构详解（以 wails3-hosts 为例）

这个项目实现了一个 hosts 文件管理工具，让你可以在 GUI 中管理系统的 `/etc/hosts`。

### 完整目录结构

```
wails3-hosts/
├── main.go                          # 应用入口
├── internal/
│   ├── domain/                      # 核心业务模型
│   │   ├── entity/                  # 实体定义
│   │   │   ├── hosts_group.go       # 分组实体
│   │   │   └── hosts_entry.go       # 条目实体
│   │   ├── repository/              # 仓储接口（抽象）
│   │   └── service/                 # 领域服务
│   ├── application/                 # 应用层（用例）
│   │   ├── service/hosts_app_service.go
│   │   └── dto/                     # 数据传输对象
│   ├── infrastructure/              # 基础设施层
│   │   ├── persistence/             # 数据持久化（JSON 文件）
│   │   └── system/                  # 系统操作（读写 hosts 文件）
│   └── interface/
│       └── handler/hosts_handler.go # 暴露给前端的 API
└── frontend/
    └── src/
        ├── api/hosts.ts             # 封装 Wails 自动生成的绑定
        ├── components/              # React 组件
        └── App.tsx                  # 主界面
```

> **小白提示：** 不用一开始就搞这么复杂的分层。初学者可以直接把所有 Go 逻辑写在 `main.go` 里，等项目变大再拆分。

---

## 核心概念：Go 函数如何变成前端 API

这是 Wails3 最神奇的地方。

### 第一步：写 Go 服务

```go
// internal/interface/handler/hosts_handler.go

type HostsHandler struct {
    appService *service.HostsApplicationService
}

// 这个函数会自动暴露给前端
func (h *HostsHandler) CreateGroup(name, description string) (*dto.HostsGroupDTO, error) {
    return h.appService.CreateGroup(context.Background(), dto.CreateHostsGroupRequest{
        Name:        name,
        Description: description,
    })
}
```

### 第二步：在 main.go 注册服务

```go
// main.go

app := application.New(application.Options{
    Name: "Hosts Manager",
    Services: []application.Service{
        application.NewService(hostsHandler),  // 注册！
    },
    Assets: application.AssetOptions{
        Handler: application.AssetFileServerFS(assets),
    },
})
```

### 第三步：Wails 自动生成 TypeScript 绑定

启动 `wails3 dev` 后，Wails 会在 `frontend/bindings/` 目录下自动生成：

```typescript
// frontend/bindings/.../hostshandler.ts（自动生成，不要手动编辑）

export function CreateGroup(name: string, description: string): Promise<...> {
    return Call("HostsHandler.CreateGroup", name, description)
}
```

### 第四步：前端直接调用

```typescript
// frontend/src/api/hosts.ts

import * as HostsHandler from '../../bindings/.../hostshandler'

export const hostsApi = {
    async createGroup(name: string, description: string) {
        const result = await HostsHandler.CreateGroup(name, description)
        return result
    }
}
```

### 第五步：在 React 组件里用

```tsx
// frontend/src/App.tsx

const [groups, setGroups] = useState([])

async function loadGroups() {
    const data = await hostsApi.getAllGroups()
    setGroups(data)
}
```

**一句话总结：** Go 里写个方法 → 注册到服务 → Wails 自动生成 TS 类型 → 前端像调用本地函数一样调用。

---

## Vibe Coding 实战：用 AI 开发一个功能

下面演示如何用 Claude Code 从零实现"添加 hosts 分组"功能。

### 对话示例

**你：**
> 我在用 Wails3 + React 做一个 hosts 文件管理工具。我需要实现"创建分组"功能：用户输入分组名称，点击确认后调用 Go 后端创建，然后刷新侧边栏列表。目前已有 `hostsApi.createGroup(name, description)` 方法。帮我写这个 Modal 组件。

**Claude：**
> 好的，我来实现这个创建分组的 Modal 组件...

```tsx
// frontend/src/components/CreateGroupModal.tsx

import { useState } from 'react'
import { hostsApi } from '../api/hosts'

interface Props {
    isOpen: boolean
    onClose: () => void
    onCreated: () => void  // 创建成功后刷新列表
}

export function CreateGroupModal({ isOpen, onClose, onCreated }: Props) {
    const [name, setName] = useState('')
    const [description, setDescription] = useState('')
    const [loading, setLoading] = useState(false)

    async function handleSubmit() {
        if (!name.trim()) return
        setLoading(true)
        try {
            await hostsApi.createGroup(name, description)
            onCreated()  // 通知父组件刷新
            onClose()
            setName('')
            setDescription('')
        } catch (err) {
            console.error('创建分组失败', err)
        } finally {
            setLoading(false)
        }
    }

    if (!isOpen) return null

    return (
        <div className="modal-overlay">
            <div className="modal">
                <h2>创建分组</h2>
                <input
                    value={name}
                    onChange={e => setName(e.target.value)}
                    placeholder="分组名称"
                />
                <input
                    value={description}
                    onChange={e => setDescription(e.target.value)}
                    placeholder="描述（可选）"
                />
                <button onClick={handleSubmit} disabled={loading}>
                    {loading ? '创建中...' : '创建'}
                </button>
                <button onClick={onClose}>取消</button>
            </div>
        </div>
    )
}
```

### AI 辅助的技巧

**给 AI 好的上下文：**

```
❌ 差的提示: "帮我写个添加功能"

✅ 好的提示: "我在用 Wails3 + React + Tailwind 做 hosts 管理工具。
  已有 hostsApi.createGroup(name: string, description: string) 
  会返回 Promise<HostsGroup>，其中 HostsGroup 的类型定义在 
  src/types/hosts.ts。请帮我写一个 Modal 组件用于创建分组，
  样式参考已有的 src/components/ui/Modal.tsx"
```

**让 AI 帮你做系统操作（需要 sudo）：**

这个项目里修改 `/etc/hosts` 需要 root 权限，这种平台相关的代码特别适合让 AI 写：

```
你: 我需要在 macOS 上以 root 权限写文件，但不想把密码存在磁盘上。
    Go 里怎么安全地处理 sudo 密码？

AI: 可以用管道把密码传给 sudo -S，在内存里处理，不写磁盘...
```

---

## 主窗口配置

`main.go` 里的窗口配置是你最常修改的地方：

```go
app.Window.NewWithOptions(application.WebviewWindowOptions{
    Title:  "Hosts Manager",
    Width:  1200,
    Height: 800,
    
    // macOS 特有：毛玻璃效果
    Mac: application.MacWindow{
        InvisibleTitleBarHeight: 50,          // 标题栏区域留给自定义 UI
        Backdrop: application.MacBackdropTranslucent,  // 半透明背景
    },
    
    BackgroundColour: application.NewRGB(27, 38, 54),  // 深色背景色
    URL: "/",  // 加载前端入口
})
```

---

## 前端监听 Go 发出的事件

有时 Go 需要主动推送消息给前端（比如菜单点击），用事件系统：

**Go 端发送事件：**

```go
// 菜单点击时发送事件
aboutMenu.OnClick(func(_ *application.Context) {
    app.Event.Emit("show-about-dialog", map[string]interface{}{
        "version": "1.0.0",
    })
})
```

**前端监听事件：**

```typescript
import { Events } from '@wailsio/runtime'

useEffect(() => {
    const unlisten = Events.On('show-about-dialog', (data) => {
        setAboutInfo(data)
        setShowAbout(true)
    })
    
    return unlisten  // 组件卸载时取消监听
}, [])
```

---

## 构建和打包

### 开发构建（快，包含调试信息）

```bash
task build:dev
# 或
wails3 build -debug -clean
```

### 生产构建（优化，体积小）

```bash
task build
```

### 打包为安装包

```bash
# macOS: 生成 .app 和 .dmg
task darwin:package

# Windows: 生成 .exe 和 NSIS 安装包
task windows:package
```

---

## 常见问题

### Q: 修改了 Go 方法签名，前端报类型错误？

Wails3 在 `dev` 模式下会自动重新生成绑定。如果手动构建，需要：

```bash
task bindings:clean  # 清理并重新生成绑定
```

### Q: 前端调用 Go 函数但没有反应？

检查：
1. `main.go` 里有没有用 `application.NewService()` 注册这个 Handler
2. Go 方法名首字母要大写（Go 的导出规则）
3. 打开浏览器 DevTools 看 console 报错

### Q: 如何调试 Go 后端？

方法一：用 `fmt.Println` 输出，在终端看日志

```go
fmt.Println("[Handler] ApplyHosts 被调用")
```

方法二：`wails3 dev` 模式下，Go 的 `log.Println` 会输出到终端

### Q: 打包后应用图标怎么换？

替换 `build/darwin/` 或 `build/windows/` 里的图标文件，然后重新构建。

---

## 推荐的 Vibe Coding 开发节奏

```
1. 想清楚要做什么（用自然语言描述，不用写代码）
      ↓
2. 告诉 AI 当前项目的技术栈和已有的结构
      ↓
3. 让 AI 生成第一版代码
      ↓
4. task dev 跑起来，在真实应用里测试
      ↓
5. 发现问题 → 描述给 AI → AI 修改
      ↓
6. 确认功能 OK 后提交 git
```

**重要原则：**
- 每次只让 AI 做一件事，不要一次性提太多需求
- 让 AI 写完后自己读一遍，理解它在做什么
- 遇到 AI 一直改不好的问题，换个角度描述或者自己动手改

---

## 延伸阅读

- [Wails3 官方文档](https://v3.wails.io/)
- [本项目源码](https://github.com/chen/wails3-hosts)
- [Go 入门教程](https://tour.golang.org/)
- [React 官方文档](https://react.dev/)
