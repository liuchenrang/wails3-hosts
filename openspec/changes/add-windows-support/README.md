# Windows 平台支持提案

## 变更概览

**变更ID**: `add-windows-support`
**类型**: 功能增强
**状态**: 待审查
**预计工作量**: 10-15 天

## 快速摘要

为 Hosts Manager 添加完整的 Windows 平台支持，包括：
- ✅ Windows UAC 权限提升机制
- ✅ hosts 文件读写操作
- ✅ 备份/恢复功能
- ✅ 平台抽象架构设计
- ✅ 保持 Unix/Linux/macOS 功能不变

## 文档结构

```
add-windows-support/
├── README.md           # 本文件
├── proposal.md         # 提案概述
├── design.md           # 技术设计文档
├── tasks.md            # 实施任务清单
└── specs/
    └── windows-platform-support/
        └── spec.md     # 需求规格说明
```

## 核心设计

### 架构模式
采用**策略模式**实现平台特定的权限提升：

```
HostsFileOperator
    ↓ 依赖
PrivilegeElevator (接口)
    ↓ 实现
├── UnixElevator    (sudo)
└── WindowsElevator (UAC)
```

### 关键决策

1. **接口抽象**: 定义 `PrivilegeElevator` 接口统一权限提升逻辑
2. **平台隔离**: Unix 和 Windows 实现完全分离，互不影响
3. **依赖注入**: 通过工厂函数创建平台特定的实例
4. **用户体验**: Windows 使用 UAC 提示，Unix 继续使用 sudo 缓存

### 技术亮点

- **SOLID 原则**: 单一职责、依赖倒置、开闭原则
- **平台检测**: 使用 `runtime.GOOS` 和构建标签
- **Windows API**: 通过 `golang.org/x/sys/windows` 调用
- **UAC 提权**: 使用 `ShellExecuteW` 和 "runas" 动词
- **向后兼容**: 现有 Unix 功能完全保留

## 功能范围

### 包含 ✅
- Windows 10/11 平台支持
- UAC 权限提升
- hosts 文件路径处理
- 备份/恢复功能
- 错误处理和用户提示
- 单元测试和集成测试

### 不包含 ❌
- Windows 服务安装
- MSI 安装程序优化
- 系统托盘集成
- 任务计划程序权限缓存（未来优化）

## 影响范围

### 修改的文件
- `internal/infrastructure/system/hosts_file_operator.go` - 使用接口
- `internal/infrastructure/system/sudo_manager.go` - 保留
- `main.go` - 工厂函数修改

### 新增的文件
- `internal/infrastructure/system/privilege.go` - 接口定义
- `internal/infrastructure/system/privilege_unix.go` - Unix 实现
- `internal/infrastructure/system/privilege_windows.go` - Windows 实现
- `internal/infrastructure/system/privilege_windows_test.go` - 测试

### 不受影响
- 前端代码（API 接口保持不变）
- 应用服务层（仅错误处理微调）
- 领域层（完全不变）

## 实施计划

### 阶段划分
1. **接口设计** (1-2天) - 定义抽象接口和工厂函数
2. **Windows 实现** (2-3天) - 实现 UAC 提权和文件操作
3. **集成重构** (1-2天) - 集成到现有代码
4. **测试验证** (2-3天) - 跨平台测试
5. **发布准备** (1天) - 文档和审查

### 关键里程碑
- Day 2: 接口定义完成
- Day 5: Windows 实现完成
- Day 7: 集成测试通过
- Day 9: Windows 平台验证通过
- Day 11: PR 创建并待审查

## 风险评估

### 技术风险 - 中等
- **Windows API 调用复杂性**
  - 缓解: 提前研究，准备示例代码

### 用户体验风险 - 中等
- **UAC 每次需要确认**
  - 缓解: 添加清晰的用户引导文案

### 兼容性风险 - 低
- **Unix 平台功能回归**
  - 缓解: 保持现有代码路径，充分回归测试

## 验收标准

- [ ] Windows 10/11 所有核心功能正常工作
- [ ] Unix/Linux/macOS 功能无回归
- [ ] 单元测试覆盖率 > 80%
- [ ] 跨平台集成测试通过
- [ ] 文档完整且准确
- [ ] 代码审查通过

## 后续工作

未来可以考虑的优化：
- Windows 任务计划程序（权限缓存）
- 性能优化（减少 UAC 确认次数）
- Windows 安装程序优化
- 系统托盘集成

## 相关资源

- [OpenSpec 规范](../../AGENTS.md)
- [项目上下文](../../project.md)
- [Windows API 文档](https://learn.microsoft.com/en-us/windows/win32/api/)
- [Go Windows 包](https://pkg.go.dev/golang.org/x/sys/windows)

## 联系方式

如有问题或建议，请：
1. 创建 GitHub Issue
2. 在代码审查中提出
3. 联系项目负责人

---

**最后更新**: 2025-01-19
**提案版本**: 1.0
