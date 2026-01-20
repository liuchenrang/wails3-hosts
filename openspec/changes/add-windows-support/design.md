# 设计文档: Windows 权限提升机制

## 架构设计

### 整体架构
采用**策略模式**实现平台特定的权限提升逻辑：

```
┌─────────────────────────────────────────┐
│     HostsFileOperator                   │
│     (基础设施层 - 文件操作)              │
└──────────────┬──────────────────────────┘
               │ 依赖
               ↓
┌─────────────────────────────────────────┐
│     PrivilegeElevator (接口)            │
│     - Execute(content) error            │
│     - Validate(password) bool           │
└──────────────┬──────────────────────────┘
               │
      ┌────────┴────────┐
      ↓                 ↓
┌─────────────┐  ┌──────────────┐
│UnixElevator │  │WindowsElevator│
│(sudo)       │  │(UAC)         │
└─────────────┘  └──────────────┘
```

### 设计原则应用
1. **单一职责 (S)**: 每个提升器只负责一个平台的权限提升
2. **开闭原则 (O)**: 通过接口扩展，无需修改 `HostsFileOperator`
3. **依赖倒置 (D)**: 依赖抽象接口，而非具体实现
4. **接口隔离 (I)**: 接口最小化，只包含必要方法

## 核心接口设计

### PrivilegeElevator 接口
```go
// PrivilegeElevator 权限提升器接口
// 单一职责: 定义平台无关的权限提升操作
type PrivilegeElevator interface {
    // Validate 验证凭据是否有效
    // Unix: 验证 sudo 密码
    // Windows: 验证管理员令牌（通常不需要）
    Validate(credentials string) bool

    // Execute 执行需要提升权限的操作
    // content: 要写入的内容
    // 返回: 操作结果或错误
    Execute(content string) error

    // CanCacheCredentials 是否可以缓存凭据
    // Unix: true (sudo 可以缓存)
    // Windows: false (UAC 每次需要确认)
    CanCacheCredentials() bool
}
```

## 平台实现

### UnixElevator (Unix/Linux/macOS)

**职责**: 使用 sudo 命令进行权限提升

**关键特性**:
- 复用现有的 `SudoCommand` 逻辑
- 支持 sudo 密码缓存
- 超时控制

**实现要点**:
```go
type UnixElevator struct {
    sudoCmd *SudoCommand
}

func (e *UnixElevator) Validate(password string) bool {
    // 使用 sudo -v 验证密码
}

func (e *UnixElevator) Execute(content string) error {
    // 使用 sudo 执行写入操作
}

func (e *UnixElevator) CanCacheCredentials() bool {
    return true
}
```

### WindowsElevator

**职责**: 使用 UAC 进行权限提升

**关键特性**:
- 使用 Windows API 重新启动进程并提升权限
- 通过进程间通信传递内容
- 不支持凭据缓存（每次需要 UAC 确认）

**技术方案**:

#### 方案 1: 进程重启 + 命令行参数
```
主进程 (普通权限)
    ↓ 检测需要写入
    ↓ 重新启动自身 (带管理员令牌)
    ↓ 传递 --admin-mode 参数
子进程 (管理员权限)
    ↓ 执行写入操作
    ↓ 退出
```

**优点**: 实现简单，Windows 标准做法
**缺点**: 需要进程重启开销

#### 方案 2: COM 对象 + IPC
```
主进程
    ↓ 创建 COM 服务器 (如果不存在)
    ↓ 通过 IPC 调用管理员进程
管理员辅助进程
    ↓ 持续运行，监听请求
```

**优点**: 性能更好，可以缓存权限
**缺点**: 实现复杂，需要进程间通信

**决策**: 采用**方案 1**（KISS 原则）

### WindowsElevator 实现细节

#### 提权流程
1. **检测当前权限**: 检查是否已具有管理员权限
2. **生成临时文件**: 将要写入的内容保存到临时文件
3. **重启进程**: 使用 `ShellExecuteW` 以管理员身份重启
4. **执行写入**: 子进程读取临时文件并写入 hosts
5. **清理**: 删除临时文件

#### 代码框架
```go
type WindowsElevator struct {
    hostsFilePath string
}

func (e *WindowsElevator) Validate(_ string) bool {
    // Windows 通过 UAC 弹窗验证，无需预先验证
    return true
}

func (e *WindowsElevator) Execute(content string) error {
    // 1. 检查是否已有管理员权限
    if e.isAdmin() {
        return e.writeDirectly(content)
    }

    // 2. 生成临时文件
    tmpFile, err := e.createTempFile(content)
    if err != nil {
        return err
    }
    defer os.Remove(tmpFile)

    // 3. 使用 UAC 重启进程
    return e.restartWithAdmin(tmpFile)
}

func (e *WindowsElevator) CanCacheCredentials() bool {
    // Windows UAC 不支持缓存
    return false
}
```

#### Windows API 调用
```go
import "golang.org/x/sys/windows"

// 检查管理员权限
func (e *WindowsElevator) isAdmin() bool {
    var sid *windows.SID
    err := windows.AllocateAndInitializeSid(
        &windows.SECURITY_NT_AUTHORITY,
        2,
        windows.SECURITY_BUILTIN_DOMAIN_RID,
        windows.DOMAIN_ALIAS_RID_ADMINS,
        0, 0, 0, 0, 0, 0,
        &sid,
    )
    if err != nil {
        return false
    }
    defer windows.FreeSid(sid)

    token := windows.Token(0)
    return token.IsMember(sid)
}

// UAC 提权重启
func (e *WindowsElevator) restartWithAdmin(tmpFile string) error {
    // 使用 ShellExecuteW "runas" 动词
    // 实现: 见完整代码
}
```

## 集成到现有代码

### HostsFileOperator 修改

#### 修改前
```go
type HostsFileOperator struct {
    hostsFilePath string
    backupDir     string
}

func (o *HostsFileOperator) Write(content string) error {
    // 直接调用 sudo
    cmd := exec.Command("sudo", "sh", "-c", script)
    // ...
}
```

#### 修改后
```go
type HostsFileOperator struct {
    hostsFilePath string
    backupDir     string
    elevator      PrivilegeElevator  // 新增
}

func NewHostsFileOperator(elevator PrivilegeElevator) (*HostsFileOperator, error) {
    // elevator 由外部注入
}

func (o *HostsFileOperator) Write(content string) error {
    return o.elevator.Execute(content)
}
```

### 工厂函数修改

```go
// initializeInfrastructure 中的修改
func initializeInfrastructure() (*infrastructure, error) {
    // ...

    // 创建平台特定的权限提升器
    var elevator PrivilegeElevator
    switch runtime.GOOS {
    case "windows":
        elevator = system.NewWindowsElevator()
    default: // darwin, linux
        elevator = system.NewUnixElevator()
    }

    // 传递给 HostsFileOperator
    hostsFileOp, err := system.NewHostsFileOperator(elevator)
    // ...
}
```

## 错误处理

### Windows 特定错误
```go
var (
    ErrUACCancelled   = errors.New("用户取消了 UAC 提权")
    ErrAdminRequired  = errors.New("需要管理员权限")
    ErrTempFileFailed = errors.New("创建临时文件失败")
)
```

### 用户友好提示
- UAC 取消: "您取消了管理员权限确认，无法应用 hosts 配置"
- 权限不足: "此操作需要管理员权限，请在 UAC 提示中允许"
- 写入失败: "写入 hosts 文件失败，请检查文件是否被占用"

## 测试策略

### 单元测试
```go
// 模拟 Windows 环境
func TestWindowsElevator_Execute(t *testing.T) {
    // 使用 mock 或构建标签测试
}
```

### 平台测试
- **macOS/Linux**: 现有测试继续运行
- **Windows**: 需要在真实 Windows 环境测试

### 构建标签
```go
// +build windows

package system

// Windows 特定实现
```

## 性能考虑

### Unix/Linux/macOS
- 首次: ~0.5s (sudo 验证)
- 缓存期: ~0.1s (无需密码)
- 5分钟后: ~0.5s (重新验证)

### Windows
- 每次: ~1-2s (UAC 提示 + 进程重启)
- 无法缓存: 每次都需要用户确认

**优化方向**: 未来可考虑使用 Windows 任务计划程序实现权限缓存

## 安全考虑

### Unix/Linux/macOS
- ✅ 密码仅在内存中，不写入文件
- ✅ 使用系统 sudo 缓存机制
- ✅ 超时自动清除

### Windows
- ✅ 使用系统 UAC 机制
- ✅ 临时文件使用安全权限
- ✅ 进程间通信使用文件而非网络

## 用户体验

### Unix/Linux/macOS 流程
1. 首次应用: 输入密码 → 缓存 5 分钟
2. 5分钟内: 直接应用，无需密码
3. 5分钟后: 再次输入密码

### Windows 流程
1. 每次应用: 弹出 UAC 提示 → 用户确认
2. 应用完成: 返回主窗口

**提示文案优化**:
```
Windows: "需要管理员权限来修改 hosts 文件"
         "请在 UAC 提示窗口中点击【是】"
Unix:    "需要输入 sudo 密码来修改 hosts 文件"
        "密码将缓存 5 分钟"
```

## 依赖项
- 新增: `golang.org/x/sys/windows` - Windows API 绑定
- 现有: `runtime` - 平台检测

## 兼容性矩阵

| 功能 | Unix | macOS | Windows |
|-----|------|-------|---------|
| 读取 hosts | ✅ | ✅ | ✅ |
| 写入 hosts | ✅ | ✅ | ✅ |
| 权限缓存 | ✅ | ✅ | ❌ |
| 备份/恢复 | ✅ | ✅ | ✅ |
| UAC 提示 | N/A | N/A | ✅ |

## 文件清单

### 新增文件
- `internal/infrastructure/system/privilege.go` - 接口定义
- `internal/infrastructure/system/privilege_unix.go` - Unix 实现
- `internal/infrastructure/system/privilege_windows.go` - Windows 实现
- `internal/infrastructure/system/privilege_windows_test.go` - Windows 测试

### 修改文件
- `internal/infrastructure/system/hosts_file_operator.go` - 使用接口
- `internal/infrastructure/system/sudo_manager.go` - 保留（缓存管理）
- `main.go` - 工厂函数修改

### 删除文件
- 无（向后兼容）

## 实施顺序
1. 定义接口 (`privilege.go`)
2. 实现 UnixElevator (重构现有代码)
3. 实现 WindowsElevator (新增)
4. 修改 HostsFileOperator (集成)
5. 更新工厂函数 (main.go)
6. 编写测试
7. Windows 平台验证

## 回滚计划
如果 Windows 实现出现问题：
1. 保留 UnixElevator 作为默认
2. WindowsElevator 可以被禁用
3. 回退到原有 sudo 逻辑（Unix/Win 混用可能有问题）

**建议**: 在功能分支完整开发并测试后再合并
