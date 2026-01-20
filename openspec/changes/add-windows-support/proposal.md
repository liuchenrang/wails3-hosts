# 提案: 实现 Windows 版本的 hosts 兼容处理

## 变更ID
`add-windows-support`

## 概述
为 Windows 平台实现完整的 hosts 文件管理功能，包括权限提升机制、文件操作和用户交互流程。

## 问题陈述

### 当前问题
1. **Unix 依赖**: 当前实现使用 `sudo` 命令进行权限提升，Windows 不支持
2. **权限模型差异**:
   - Unix/Linux: 使用 sudo 提权
   - Windows: 使用 UAC (User Account Control) 和管理员令牌
3. **文件路径差异**: 已部分支持，但缺少完整测试
4. **缺失功能**: Windows 平台无完整的提权流程和错误处理

### 影响范围
- Windows 用户无法应用 hosts 配置
- 回滚功能在 Windows 上不可用
- 备份/恢复功能受限

## 解决方案概述

### 核心策略
采用**平台抽象**模式，在基础设施层实现平台特定的权限提升器：

1. **权限提升器接口**: 定义统一的 `PrivilegeElevator` 接口
2. **平台实现**:
   - `UnixElevator`: 使用 sudo
   - `WindowsElevator`: 使用 UAC 提权
3. **依赖注入**: `HostsFileOperator` 依赖抽象接口而非具体实现

### 技术路径
- **Go 1.24+** 利用 `runtime.GOOS` 进行平台判断
- **Windows API** 使用 `golang.org/x/sys/windows` 包调用 Windows API
- **UAC 提权** 通过 ShellExecuteW 运行具有管理员权限的进程

## 架构影响

### 新增组件
```
infrastructure/system/
├── privilege.go          # 权限提升器接口
├── privilege_unix.go     # Unix/Linux/macOS 实现
├── privilege_windows.go  # Windows 实现
└── hosts_file_operator.go  # 修改: 使用接口而非直接调用 sudo
```

### 依赖关系
```
HostsFileOperator
    ↓ 依赖
PrivilegeElevator (接口)
    ↓ 实现
UnixElevator / WindowsElevator
```

## 功能范围

### 包含功能
1. ✅ Windows hosts 文件路径正确识别
2. ✅ UAC 权限提升机制
3. ✅ Windows 平台的文件写入操作
4. ✅ 备份/恢复功能支持
5. ✅ 错误处理和用户提示
6. ✅ 单元测试（模拟 Windows 环境）

### 不包含功能
- ❌ Windows 服务安装
- ❌ 系统托盘集成
- ❌ MSI 安装程序优化
- ❌ Windows 特定的 UI 定制

## 兼容性考虑

### 平台支持
- ✅ Windows 10/11 (64-bit)
- ✅ Windows Server 2016+
- ✅ 保持 macOS/Linux 功能不受影响

### 向后兼容
- ✅ 现有 Unix/Linux/macOS 功能完全保留
- ✅ API 接口保持不变
- ✅ 前端代码无需修改

## 风险评估

### 技术风险
- **中等**: Windows UAC 提权复杂性
  - 缓解: 使用成熟的 Windows API 和测试覆盖
- **低**: 文件路径处理差异
  - 缓解: 已有部分支持，增强测试即可

### 用户体验风险
- **中等**: Windows 用户首次使用需确认 UAC 提示
  - 缓解: 添加清晰的用户引导文案

## 性能影响
- **Windows**: UAC 提权每次需要用户确认，无法缓存（与 sudo 不同）
- **Unix/Linux/macOS**: 无影响，保持现有性能

## 替代方案

### 方案 A: 要求用户以管理员身份运行应用
**优点**: 实现简单
**缺点**: 用户体验差，不符合最佳实践
**决策**: ❌ 不采用

### 方案 B: 使用 Windows 任务计划程序
**优点**: 可以缓存权限
**缺点**: 实现复杂，需要安装时配置
**决策**: ❌ 不采用（未来可考虑）

### 方案 C: UAC 按需提权（当前方案）
**优点**: 符合 Windows 最佳实践，安全且用户友好
**缺点**: 每次操作需要确认
**决策**: ✅ 采用

## 成功标准
1. ✅ Windows 平台能成功应用 hosts 配置
2. ✅ 所有核心功能在 Windows 上正常工作
3. ✅ 单元测试覆盖率 > 80%
4. ✅ 现有平台功能无回归

## 后续工作
- [ ] Windows 安装程序优化
- [ ] 性能优化（考虑任务计划程序缓存权限）
- [ ] 添加 Windows 集成测试
- [ ] 文档更新（Windows 用户指南）

## 相关变更
- 无前置依赖
- 不阻塞其他功能开发
- 可独立开发和测试
