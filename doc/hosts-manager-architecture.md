# Hosts Manager - DDD 架构设计

## 项目概述
一个类似 SwitchHosts 的跨平台 hosts 文件管理工具，支持分组管理、版本历史、主题切换等功能。

## 技术栈
- **后端**: Go 1.24 + Wails v3
- **前端**: React 18 + Tailwind CSS + shadcn/ui
- **国际化**: i18next
- **存储**: JSON 文件（配置和版本历史）

---

## DDD 分层架构

```
wails3-demo/
├── main.go                          # 应用入口
├── wails.json                       # Wails 配置
│
├── internal/                        # 内部包
│   ├── domain/                      # 领域层
│   │   ├── entity/                  # 实体
│   │   │   ├── hosts_group.go       # Hosts 分组实体
│   │   │   ├── hosts_entry.go       # Hosts 条目实体
│   │   │   └── hosts_version.go     # 版本历史实体
│   │   ├── valueobject/             # 值对象
│   │   │   ├── hosts_content.go     # Hosts 内容值对象
│   │   │   └── sudo_credentials.go  # Sudo 凭证值对象
│   │   ├── repository/              # 仓储接口
│   │   │   ├── hosts_repository.go  # Hosts 仓储接口
│   │   │   └── version_repository.go # 版本仓储接口
│   │   └── service/                 # 领域服务
│   │       └── hosts_domain_service.go
│   │
│   ├── application/                 # 应用层
│   │   ├── service/                 # 应用服务
│   │   │   └── hosts_app_service.go
│   │   └── dto/                     # 数据传输对象
│   │       ├── hosts_group_dto.go
│   │       ├── hosts_entry_dto.go
│   │       └── version_dto.go
│   │
│   ├── infrastructure/              # 基础设施层
│   │   ├── persistence/             # 持久化
│   │   │   ├── json_storage.go      # JSON 存储
│   │   │   └── hosts_repository_impl.go
│   │   └── system/                  # 系统操作
│   │       ├── hosts_file.go        # Hosts 文件操作
│   │       └── sudo_manager.go      # Sudo 权限管理
│   │
│   └── interface/                   # 接口层
│       └── handler/                 # Wails 服务处理器
│           └── hosts_handler.go
│
├── frontend/                        # 前端
│   ├── src/
│   │   ├── api/                     # API 调用
│   │   ├── components/              # React 组件
│   │   ├── hooks/                   # 自定义 Hooks
│   │   ├── i18n/                    # 国际化资源
│   │   ├── theme/                   # 主题配置
│   │   └── types/                   # TypeScript 类型
│   └── ...
│
└── design/                          # 设计文档
```

---

## 核心领域模型

### 1. 实体 (Entities)

#### HostsGroup (Hosts 分组)
```go
type HostsGroup struct {
    ID          string
    Name        string
    Description string
    IsEnabled   bool
    Entries     []HostsEntry
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

#### HostsEntry (Hosts 条目)
```go
type HostsEntry struct {
    ID       string
    IP       string
    Hostname string
    Comment  string
    Enabled  bool
}
```

#### HostsVersion (版本历史)
```go
type HostsVersion struct {
    ID          string
    Timestamp   time.Time
    Content     string      // 完整的 hosts 文件内容
    Description string      // 版本描述（自动生成或手动添加）
    Source      string      // "manual" | "auto" | "rollback"
}
```

### 2. 值对象 (Value Objects)

#### HostsContent
```go
type HostsContent struct {
    RawContent string
    ParsedEntries []HostsEntry
}
```

#### SudoCredentials
```go
type SudoCredentials struct {
    Password string
    Cached  bool
    ExpiresAt time.Time
}
```

### 3. 仓储接口 (Repository Interfaces)

#### HostsRepository
```go
type HostsRepository interface {
    Save(group HostsGroup) error
    FindByID(id string) (*HostsGroup, error)
    FindAll() ([]HostsGroup, error)
    Delete(id string) error
    UpdateEnabled(id string, enabled bool) error
}
```

#### VersionRepository
```go
type VersionRepository interface {
    Save(version HostsVersion) error
    FindLatest(limit int) ([]HostsVersion, error)
    FindByID(id string) (*HostsVersion, error)
    Delete(id string) error
    DeleteBefore(timestamp time.Time) error
}
```

### 4. 领域服务 (Domain Services)

#### HostsDomainService
- 验证 IP 地址格式
- 验证主机名格式
- 合并多个分组的内容
- 生成最终 hosts 文件内容
- 检测 hosts 文件冲突

---

## 应用服务层

### HostsApplicationService
职责：
- 协调领域对象完成业务用例
- 处理事务边界
- 调用基础设施服务

主要方法：
```go
type HostsApplicationService struct {
    hostsRepo    domain.HostsRepository
    versionRepo  domain.VersionRepository
    hostsFileOp  system.HostsFileOperator
    sudoManager  system.SudoManager
}

// 创建分组
func (s *HostsApplicationService) CreateGroup(name, description string) (*dto.HostsGroupDTO, error)

// 添加条目
func (s *HostsApplicationService) AddEntry(groupID, ip, hostname, comment string) error

// 切换分组启用状态
func (s *HostsApplicationService) ToggleGroup(groupID string) error

// 应用 hosts 配置（需要 sudo 权限）
func (s *HostsApplicationService) ApplyHosts() error

// 获取版本历史
func (s *HostsApplicationService) GetVersions(limit int) ([]dto.VersionDTO, error)

// 回滚到指定版本
func (s *HostsApplicationService) RollbackToVersion(versionID string) error

// 验证/缓存 sudo 密码
func (s *HostsApplicationService) ValidateSudoPassword(password string) (bool, error)
```

---

## 基础设施层

### 1. 持久化

#### JSONStorage
```go
type JSONStorage struct {
    configPath string
    versionsPath string
}

func (s *JSONStorage) LoadGroups() ([]domain.HostsGroup, error)
func (s *JSONStorage) SaveGroups(groups []domain.HostsGroup) error
func (s *JSONStorage) LoadVersions() ([]domain.HostsVersion, error)
func (s *JSONStorage) SaveVersions(versions []domain.HostsVersion) error
```

### 2. 系统操作

#### HostsFileOperator
```go
type HostsFileOperator struct {
    hostsFilePath string
}

func (o *HostsFileOperator) ReadCurrent() (string, error)
func (o *HostsFileOperator) Write(content string, sudoPassword string) error
func (o *HostsFileOperator) Backup() error
```

#### SudoManager
```go
type SudoManager struct {
    cachedPassword string
    expiresAt      time.Time
    cacheDuration  time.Duration // 5分钟
}

func (m *SudoManager) ValidatePassword(password string) bool
func (m *SudoManager) WriteWithSudo(filePath, content, password string) error
func (m *SudoManager) IsPasswordCached() bool
```

---

## 前端架构

### 目录结构
```
frontend/src/
├── api/
│   └── hosts.ts                    # Wails API 调用封装
├── components/
│   ├── layout/
│   │   ├── Sidebar.jsx             # 左侧分组列表
│   │   └── MainPanel.jsx           # 右侧主面板
│   ├── hosts/
│   │   ├── HostsGroupList.jsx      # 分组列表
│   │   ├── HostsEntryList.jsx      # 条目列表
│   │   ├── HostsEditor.jsx         # 编辑器
│   │   └── VersionHistory.jsx      # 版本历史
│   └── ui/                         # shadcn/ui 组件
├── hooks/
│   ├── useHostsGroups.ts           # 分组管理 Hook
│   ├── useVersions.ts              # 版本历史 Hook
│   ├── useTheme.ts                 # 主题 Hook
│   └── useSudo.ts                  # Sudo 密码 Hook
├── i18n/
│   ├── index.ts                    # i18n 配置
│   └── locales/
│       ├── zh-CN.ts                # 简体中文
│       ├── en-US.ts                # 英文
│       └── ja-JP.ts                # 日文（示例）
├── theme/
│   ├── themes.ts                   # 主题定义
│   └── index.ts                    # 主题配置
├── types/
│   ├── hosts.ts                    # Hosts 类型定义
│   └── app.ts                      # 应用类型
├── utils/
│   ├── hotkeys.ts                  # 快捷键处理
│   └── validators.ts               # 验证函数
└── App.jsx
```

### 主要组件

#### Sidebar (左侧分组管理)
- 显示所有分组列表
- 支持拖拽排序（可选）
- 启用/禁用切换
- 新建/删除分组

#### MainPanel (右侧主面板)
- 显示选中分组的 hosts 条目
- 支持添加/编辑/删除条目
- IP 地址和主机名验证
- 实时预览生成的 hosts 内容

#### VersionHistory (版本历史)
- 时间线展示历史版本
- 显示版本描述和来源
- 支持回滚操作
- 版本对比（可选）

---

## 功能实现细节

### 1. 国际化 (i18n)

**语言检测顺序**：
1. 读取系统语言
2. 查找是否支持该语言
3. 不支持则使用中文（zh-CN）作为默认

**实现**：
```typescript
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

const systemLang = navigator.language || 'zh-CN';
const supportedLangs = ['zh-CN', 'en-US', 'ja-JP'];
const defaultLang = supportedLangs.includes(systemLang) ? systemLang : 'zh-CN';

i18n.use(initReactI18next).init({
  lng: defaultLang,
  fallbackLng: 'zh-CN',
  resources: {
    'zh-CN': { translation: zhCN },
    'en-US': { translation: enUS },
    'ja-JP': { translation: jaJP },
  },
});
```

### 2. 主题系统

**主题类型**：
- Light（明亮）
- Dark（暗色）

**实现**：
```typescript
// theme/themes.ts
export const themes = {
  light: {
    background: 'ffffff',
    foreground: '09090b',
    primary: '3b82f6',
    // ...
  },
  dark: {
    background: '09090b',
    foreground: 'fafafa',
    primary: '3b82f6',
    // ...
  },
};

// hooks/useTheme.ts
const useTheme = () => {
  const [theme, setTheme] = useState('light');

  useEffect(() => {
    // 从本地存储读取主题偏好
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme) {
      setTheme(savedTheme);
    } else {
      // 跟随系统主题
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
      setTheme(prefersDark ? 'dark' : 'light');
    }
  }, []);

  const toggleTheme = () => {
    setTheme(prev => {
      const newTheme = prev === 'light' ? 'dark' : 'light';
      localStorage.setItem('theme', newTheme);
      return newTheme;
    });
  };

  return { theme, toggleTheme };
};
```

### 3. 快捷键保存

**快捷键定义**：
- `Cmd+S` / `Ctrl+S`: 保存并应用 hosts

**实现**：
```typescript
// utils/hotkeys.ts
export const useHotkey = (key: string, callback: () => void) => {
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      const isMac = navigator.platform.toUpperCase().indexOf('MAC') >= 0;
      const modifierKey = isMac ? e.metaKey : e.ctrlKey;

      if (modifierKey && e.key === key) {
        e.preventDefault();
        callback();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [key, callback]);
};

// 使用
useHotkey('s', () => {
  handleSaveAndApply();
});
```

### 4. Sudo 密码管理

**流程**：
1. 用户点击"应用 hosts"按钮
2. 检查是否有缓存的密码
3. 没有缓存则弹出密码输入框
4. 验证密码有效性
5. 使用 sudo 写入 hosts 文件
6. 缓存密码 5 分钟（可选）

**Go 实现**：
```go
func (m *SudoManager) WriteWithSudo(filePath, content, password string) error {
    cmd := exec.Command("sudo", "-S", "tee", filePath)
    cmd.Stdin = strings.NewReader(password + "\n")
    out, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("sudo failed: %w", err)
    }
    return nil
}
```

### 5. 版本历史

**自动创建版本的场景**：
1. 应用 hosts 配置成功后
2. 启动应用时（检测到系统 hosts 文件被外部修改）

**版本数量限制**：
- 最多保留 50 个版本
- 超过限制自动删除最旧的版本

---

## 数据存储

### 配置文件位置

**macOS**:
```
~/Library/Application Support/hosts-manager/
  ├── config.json          # 分组配置
  └── versions.json        # 版本历史
```

**Linux**:
```
~/.config/hosts-manager/
  ├── config.json
  └── versions.json
```

**Windows**:
```
%APPDATA%/hosts-manager/
  ├── config.json
  └── versions.json
```

### config.json 格式
```json
{
  "groups": [
    {
      "id": "uuid-1",
      "name": "开发环境",
      "description": "本地开发域名映射",
      "enabled": true,
      "entries": [
        {
          "id": "uuid-2",
          "ip": "127.0.0.1",
          "hostname": "localhost.local",
          "comment": "本地开发",
          "enabled": true
        }
      ],
      "created_at": "2025-12-26T10:00:00Z",
      "updated_at": "2025-12-26T10:00:00Z"
    }
  ],
  "settings": {
    "theme": "dark",
    "language": "zh-CN",
    "auto_apply": false
  }
}
```

### versions.json 格式
```json
{
  "versions": [
    {
      "id": "uuid-3",
      "timestamp": "2025-12-26T10:00:00Z",
      "content": "# Hosts Database\n127.0.0.1 localhost\n...",
      "description": "应用开发环境配置",
      "source": "manual"
    }
  ]
}
```

---

## 开发和测试计划

### 阶段 1: 基础架构搭建
- [ ] 创建 DDD 目录结构
- [ ] 实现领域实体和值对象
- [ ] 定义仓储接口
- [ ] 配置前端开发环境（Tailwind + shadcn/ui）

### 阶段 2: 核心功能实现
- [ ] 实现 JSON 持久化
- [ ] 实现 hosts 文件操作
- [ ] 实现 sudo 权限管理
- [ ] 实现应用服务层

### 阶段 3: 前端开发
- [ ] 实现基础布局（Sidebar + MainPanel）
- [ ] 实现分组管理功能
- [ ] 实现 hosts 条目编辑
- [ ] 实现版本历史界面

### 阶段 4: 高级功能
- [ ] 国际化支持
- [ ] 主题切换
- [ ] 快捷键支持
- [ ] 版本回滚

### 阶段 5: 测试和优化
- [ ] Docker 环境测试
- [ ] 跨平台测试（macOS, Linux, Windows）
- [ ] 性能优化
- [ ] 用户文档

---

## 依赖项

### Go 依赖
- `github.com/wailsapp/wails/v3` - Wails v3 框架
- `github.com/google/uuid` - UUID 生成

### 前端依赖
```json
{
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "@wailsio/runtime": "latest",
    "i18next": "^23.0.0",
    "react-i18next": "^13.0.0",
    "lucide-react": "^0.300.0",      // 图标库
    "clsx": "^2.0.0",                // 条件类名
    "tailwind-merge": "^2.0.0"       // Tailwind 类名合并
  },
  "devDependencies": {
    "@types/react": "^18.2.0",
    "@types/react-dom": "^18.2.0",
    "@vitejs/plugin-react": "^4.2.0",
    "autoprefixer": "^10.4.0",
    "postcss": "^8.4.0",
    "tailwindcss": "^3.4.0",
    "vite": "^5.0.0"
  }
}
```

---

## 安全考虑

1. **Sudo 密码安全**：
   - 密码仅保存在内存中
   - 不明文存储密码到配置文件
   - 缓存超时后自动清除

2. **Hosts 文件备份**：
   - 修改前自动备份系统 hosts 文件
   - 保留最近 5 个备份

3. **输入验证**：
   - 严格验证 IP 地址格式
   - 验证主机名格式（防止注入攻击）
   - 限制单个分组最大条目数（防止性能问题）

---

## 性能优化

1. **版本历史清理**：
   - 启动时清理过期版本（> 30 天）
   - 限制版本总数（最多 50 个）

2. **前端优化**：
   - 使用 React.memo 优化组件渲染
   - 使用虚拟滚动处理大量条目
   - 防抖保存（避免频繁写入）

3. **缓存策略**：
   - 缓存 sudo 密码 5 分钟
   - 缓存系统 hosts 文件内容

---

## 未来扩展

- [ ] 支持远程同步（GitHub Gist, Gitee）
- [ ] 支持导入/导出配置
- [ ] 支持定时切换（工作时间/开发环境）
- [ ] 支持命令行接口（CLI）
- [ ] 支持系统托盘图标
- [ ] 支持 hosts 模板和预设
