# Spec: Windows 平台支持

**Capability**: `windows-platform-support`
**Related Capabilities**:
- `hosts-file-management` (现有)
- `privilege-elevation` (新增)

## ADDED Requirements

### Requirement: 平台检测与适配
系统必须能够检测当前运行平台并为 Windows 系统提供适当的权限提升机制。

#### Scenario: Windows 平台自动识别
**Given** 应用在 Windows 系统上运行
**When** 应用初始化基础设施层
**Then** 系统应创建 `WindowsElevator` 实例而非 `UnixElevator`
**And** `HostsFileOperator` 应使用 Windows 提升器

#### Scenario: Unix 平台保持兼容
**Given** 应用在 Unix/Linux/macOS 系统上运行
**When** 应用初始化基础设施层
**Then** 系统应创建 `UnixElevator` 实例
**And** 功能行为与之前完全一致

---

### Requirement: Windows UAC 权限提升
系统必须通过 UAC (User Account Control) 提升权限以修改系统 hosts 文件。

#### Scenario: 首次应用配置时触发 UAC
**Given** 用户在 Windows 系统上首次应用 hosts 配置
**And** 当前进程不具有管理员权限
**When** 用户点击"应用配置"按钮
**Then** 系统应弹出 UAC 提示窗口
**And** 提示内容应包含"需要管理员权限来修改 hosts 文件"
**And** 用户确认后应以管理员权限执行写入操作
**And** 用户取消后应返回友好的取消提示

#### Scenario: UAC 提升后写入 hosts 文件
**Given** UAC 提示已显示
**When** 用户点击"是"确认提权
**Then** 系统应以管理员权限重启进程
**And** 临时文件应包含要写入的 hosts 内容
**And** 内容应成功写入 `C:\Windows\System32\drivers\etc\hosts`
**And** 操作完成后临时文件应被删除
**And** 主窗口应显示"配置已成功应用"

#### Scenario: 用户取消 UAC 提示
**Given** UAC 提示已显示
**When** 用户点击"否"或关闭提示窗口
**Then** 写入操作应中止
**And** 应显示友好提示"您取消了管理员权限确认，无法应用 hosts 配置"
**And** 应用应保持正常运行状态

---

### Requirement: Windows hosts 文件路径正确性
系统必须在 Windows 平台上使用正确的 hosts 文件路径。

#### Scenario: 获取 Windows hosts 文件路径
**Given** 应用在 Windows 10/11 系统上运行
**When** 系统初始化 `HostsFileOperator`
**Then** hosts 文件路径应为 `C:\Windows\System32\drivers\etc\hosts`
**And** 应使用 `SystemRoot` 环境变量动态构建路径
**And** 应支持不同的 Windows 安装目录

#### Scenario: 读取 Windows hosts 文件
**Given** 应用在 Windows 系统上运行
**When** 调用 `HostsFileOperator.ReadCurrent()`
**Then** 应成功读取 `C:\Windows\System32\drivers\etc\hosts` 内容
**And** 返回的内容应与文件实际内容一致
**And** 即使文件包含非 UTF-8 字符也应正确读取

---

### Requirement: Windows 平台备份与恢复
系统必须在 Windows 平台上支持 hosts 文件的备份和恢复功能。

#### Scenario: 创建 Windows hosts 备份
**Given** 应用在 Windows 系统上运行
**When** 调用 `HostsFileOperator.Backup()`
**Then** 应在 `%APPDATA%\hosts-manager\backups\` 创建备份文件
**And** 备份文件名应包含时间戳，格式为 `hosts_YYYYMMDD_HHMMSS.bak`
**And** 备份内容应与当前 hosts 文件内容一致
**And** 应自动清理超过 5 个的旧备份

#### Scenario: 从备份恢复 Windows hosts
**Given** 存在 Windows hosts 备份文件
**And** 备份路径为 `%APPDATA%\hosts-manager\backups\hosts_20250119_120000.bak`
**When** 调用 `HostsFileOperator.RestoreFromBackup(backupPath)`
**Then** 应弹出 UAC 提示
**And** 用户确认后应将备份内容写入系统 hosts 文件
**And** 写入成功后应显示恢复成功提示

---

### Requirement: Windows 错误处理
系统必须针对 Windows 平台提供适当的错误处理和用户提示。

#### Scenario: 处理 UAC 提权失败
**Given** 用户在 Windows 系统上应用配置
**When** UAC 提权失败（用户拒绝或系统错误）
**Then** 应捕获具体错误类型
**And** 显示明确的错误提示："无法获取管理员权限，请在 UAC 提示中允许操作"
**And** 不应导致应用崩溃

#### Scenario: 处理 hosts 文件被占用
**Given** Windows hosts 文件被其他程序占用
**When** 尝试写入 hosts 文件
**Then** 应捕获文件占用错误
**And** 显示友好提示："hosts 文件正在被其他程序使用，请关闭相关程序后重试"
**And** 建议的占位程序包括：杀毒软件、系统优化工具等

#### Scenario: 处理临时文件创建失败
**Given** Windows 系统临时目录不可写或空间不足
**When** 尝试创建临时文件用于 UAC 提权
**Then** 应捕获创建失败错误
**And** 显示提示："无法创建临时文件，请检查磁盘空间和权限"
**And** 写入操作应安全中止

---

### Requirement: Windows 配置文件存储
系统必须在 Windows 平台上使用标准的配置文件存储位置。

#### Scenario: 获取 Windows 配置目录
**Given** 应用在 Windows 系统上运行
**When** 系统初始化存储组件
**Then** 配置目录应为 `%APPDATA%\hosts-manager`
**And** 应自动创建该目录（如果不存在）
**And** 应正确解析 `APPDATA` 环境变量

#### Scenario: 保存配置到 Windows 配置目录
**Given** 用户创建分组和条目
**When** 配置被保存
**Then** `config.json` 应保存在 `%APPDATA%\hosts-manager\`
**And** `versions.json` 应保存在同一目录
**And** 文件编码应为 UTF-8 with BOM（Windows 兼容性）

---

### Requirement: Windows 凭据缓存行为
系统必须正确处理 Windows 平台上不支持凭据缓存的特性。

#### Scenario: Windows 不缓存管理员凭据
**Given** 应用在 Windows 系统上运行
**When** 用户首次成功应用配置
**Then** 不应缓存管理员令牌
**And** 下次应用配置时应再次弹出 UAC 提示
**And** 前端应显示相应提示："Windows 系统每次操作都需要管理员权限确认"

#### Scenario: 前端检测凭据缓存能力
**Given** 应用在 Windows 系统上运行
**When** 前端调用 `IsSudoPasswordCached()` API
**Then** 返回值应为 `false`（Windows 平台）
**And** 前端应据此调整 UI，不显示"密码已缓存"提示

---

### Requirement: 跨平台一致性
系统必须在功能上保持跨平台一致性，仅在实现细节上有所不同。

#### Scenario: 核心功能在 Windows 上正常工作
**Given** 应用在 Windows 系统上运行
**When** 用户执行以下操作：
- 创建分组
- 添加条目
- 启用/禁用分组
- 应用配置
- 查看版本历史
- 回滚到历史版本
**Then** 所有功能应与 Unix/Linux/macOS 平台行为一致
**And** 仅权限提升方式不同（UAC vs sudo）

#### Scenario: API 接口保持不变
**Given** 前端代码
**When** 调用后端 API
**Then** API 接口签名应保持不变
**And** 返回数据格式应保持一致
**And** 前端代码无需因 Windows 支持而修改

---

## MODIFIED Requirements

### Requirement: 权限管理器 (PrivilegeManager)
**原实现**: 仅支持 Unix/Linux/macOS 的 sudo 密码缓存
**修改后**: 支持平台特定的权限管理

#### Scenario: 根据平台创建权限管理器
**Given** 应用启动
**When** 初始化权限管理组件
**Then** Windows 平台应创建不支持缓存的权限管理器
**And** Unix 平台应创建支持 sudo 缓存的权限管理器
**And** 接口方法签名保持一致

---

## DEPRECATED Requirements

无。

---

## REMOVED Requirements

无。

---

## Non-Functional Requirements

### 性能要求
- UAC 提权导致的应用重启应在 3 秒内完成
- Windows 平台上的写入操作不应比 Unix 平台慢超过 2 倍

### 安全要求
- 临时文件应设置适当的访问权限（仅当前用户可读写）
- 临时文件应在操作完成后立即删除
- 不应在日志或临时文件中记录敏感信息

### 兼容性要求
- 支持 Windows 10（版本 1607+）
- 支持 Windows 11（所有版本）
- 支持 Windows Server 2016+
- 支持 64 位系统（不考虑 32 位）

### 可测试性要求
- Windows 实现应包含单元测试（使用构建标签）
- 应提供模拟 Windows API 的测试辅助工具
- 测试覆盖率应达到 80% 以上

---

## Dependencies

### 外部依赖
- `golang.org/x/sys/windows` - Windows API 绑定

### 内部依赖
- `internal/infrastructure/system/hosts_file_operator.go` - 修改以使用接口
- `internal/application/service/hosts_app_service.go` - 可能需要调整错误处理
- `main.go` - 修改工厂函数以创建平台特定实例

---

## Migration Notes

### 数据迁移
无需数据迁移，Windows 平台使用相同的配置文件格式（JSON）。

### 用户迁移
Unix/Linux/macOS 用户无影响。
Windows 用户需要：
1. 以管理员身份运行应用（首次）
2. 允许 UAC 提示（每次写入操作）

### 配置迁移
无需配置迁移，配置文件格式保持兼容。

---

## Testing Strategy

### 单元测试
- 测试 Windows 平台检测逻辑
- 测试 hosts 文件路径构建
- 测试临时文件创建和清理
- 测试错误处理分支

### 集成测试
- 测试 UAC 提权流程（需要真实 Windows 环境）
- 测试完整的写入流程
- 测试备份/恢复功能

### 平台测试
- 在 Windows 10 上测试所有功能
- 在 Windows 11 上测试所有功能
- 在 Unix/Linux/macOS 上回归测试确保无影响

### 模拟测试
- 使用接口模拟 Windows API 行为
- 在非 Windows 平台上运行 Windows 逻辑测试

---

## Rollback Plan

如果 Windows 实现出现问题：
1. 可以通过构建标签禁用 Windows 代码
2. 回退到仅支持 Unix 平台的版本
3. Windows 用户将无法使用此功能

**建议**: 在功能分支上完整测试后再合并到 main 分支。
