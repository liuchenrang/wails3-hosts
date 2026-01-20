# Windows 平台测试准备完成 ✅

## 📋 概述

已完成 Windows 平台测试环境的准备，包括测试计划、自动化测试脚本、测试工具和 CI/CD 配置。

---

## 📁 创建的文件

### 1. 文档文件

#### `doc/windows-testing-guide.md`
**完整的 Windows 测试指南**，包含：
- ✅ 测试环境要求
- ✅ 环境设置步骤
- ✅ 10 大测试场景（30+ 测试用例）
  - 基础功能测试
  - UAC 提权测试
  - 文件操作测试
  - 错误处理测试
  - 编码和格式测试
  - 性能测试
- ✅ 自动化 PowerShell 测试脚本
- ✅ 常见问题排查
- ✅ 测试报告模板

**使用方法**:
```bash
# 在 Windows 上查看
start doc\windows-testing-guide.md
```

---

#### `doc/windows-test-checklist.md`
**手动测试检查清单**，包含：
- ✅ 10 大测试类别
- ✅ 24 个详细测试用例
- ✅ 每个测试用例包含：步骤、预期结果、实际结果、通过/失败
- ✅ 缺陷报告模板
- ✅ 测试总结和统计
- ✅ 测试人员签名区

**使用方法**:
```bash
# 打印清单
notepad /p doc\windows-test-checklist.md

# 或在电子表格软件中打开
```

---

### 2. 测试工具

#### `test-windows.ps1`
**PowerShell 自动化测试脚本**，功能：
- ✅ 自动检查 Windows 版本
- ✅ 验证管理员权限
- ✅ 检查 Go 和 Node.js 环境
- ✅ 验证 hosts 文件路径和权限
- ✅ 检查应用配置目录和备份目录
- ✅ 尝试编译应用
- ✅ 检查 UAC 状态和提权级别
- ✅ 生成测试报告

**使用方法**:
```powershell
# 以管理员身份运行 PowerShell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# 运行测试（完整模式）
.\test-windows.ps1

# 运行测试（跳过编译）
.\test-windows.ps1 -SkipBuild

# 详细输出模式
.\test-windows.ps1 -Verbose
```

**预期输出**:
```
========================================
Windows Hosts Manager 自动化测试套件
========================================

[测试 1] 环境检查
  1.1 检查 Windows 版本...
  ✅ Windows 版本符合要求 (10+)
  1.2 检查当前权限...
  ✅ 当前具有管理员权限
  ...

========================================
测试总结
========================================
总计: 10 个测试
✅ 通过: 10 个
失败: 0 个

✅ 所有关键测试通过！
```

---

#### `cmd/test_uac/main.go`
**UAC 提权测试工具**，功能：
- ✅ 检测当前用户权限（管理员/标准用户）
- ✅ 标准用户模式下请求 UAC 提权
- ✅ 管理员模式下执行测试操作
- ✅ 验证 hosts 文件读取
- ✅ 测试临时文件写入和清理
- ✅ 完整的错误处理

**编译**:
```bash
go build -o test_uac.exe ./cmd/test_uac
```

**使用方法**:
```powershell
# 以标准用户身份运行（测试 UAC 提权）
.\test_uac.exe

# 直接以管理员身份运行
.\test_uac.exe --admin-mode
```

**预期输出**:
```
========================================
Windows UAC 提权测试工具
========================================
操作系统: windows
架构: amd64

[模式] 标准用户模式
ℹ️  当前是标准用户权限

准备请求 UAC 提权...
正在请求 UAC 提权...
提示: 请在弹出的 UAC 窗口中点击【是】

[弹出 UAC 窗口]

[模式] 管理员模式
ℹ️  此程序以管理员权限运行
✅ 验证: 确认具有管理员权限

开始执行管理员权限操作...
读取 hosts 文件: C:\Windows\System32\drivers\etc\hosts
✅ 成功读取 hosts 文件 (824 字节)
...

========================================
✅ 所有测试通过！
========================================
```

---

### 3. CI/CD 配置

#### `.github/workflows/windows-ci.yml`
**GitHub Actions Windows CI 配置**，包含：
- ✅ 自动构建 Windows 版本
- ✅ 运行 Go 单元测试
- ✅ 上传测试覆盖率到 Codecov
- ✅ Windows 特定功能测试
- ✅ 构建并测试 UAC 测试工具
- ✅ 上传构建产物
- ✅ 系统信息收集

**触发条件**:
- 推送到 `main` 或 `develop` 分支
- 创建 Pull Request
- 手动触发（workflow_dispatch）

**使用方法**:
```yaml
# 自动在以下情况运行
on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]
  workflow_dispatch:  # 手动触发
```

**查看 CI 结果**:
```
GitHub Actions → Windows CI → 查看运行日志
```

---

## 🚀 快速开始

### 方式 1: 本地 Windows 环境测试

#### 步骤 1: 环境准备
```powershell
# 安装 Go (如果未安装)
winget install Golang.Go

# 安装 Node.js (如果未安装)
winget install OpenJS.NodeJS.LTS

# 验证安装
go version
node --version
npm --version
```

#### 步骤 2: 克隆并构建
```powershell
# 克隆仓库
git clone https://github.com/chen/wails3-hosts.git
cd wails3-hosts

# 安装依赖
go mod download
cd frontend && npm install && cd ..

# 构建应用
go build -o wails3-hosts.exe .
```

#### 步骤 3: 运行测试
```powershell
# 运行自动化测试脚本
.\test-windows.ps1

# 运行 UAC 提权测试
.\test_uac.exe

# 启动应用进行手动测试
.\wails3-hosts.exe
```

#### 步骤 4: 填写测试清单
```powershell
# 打开测试检查清单
notepad doc\windows-test-checklist.md
```

---

### 方式 2: 使用 GitHub Actions

#### 步骤 1: 推送代码
```bash
# 推送到 GitHub，自动触发 CI
git add .
git commit -m "测试 Windows 支持"
git push origin main
```

#### 步骤 2: 查看 CI 结果
```
1. 打开 GitHub 仓库
2. 点击 "Actions" 标签
3. 选择 "Windows CI" workflow
4. 查看运行日志和结果
```

#### 步骤 3: 下载构建产物
```
1. 在 workflow 运行页面
2. 滚动到 "Artifacts" 部分
3. 下载 "windows-build"
4. 解压并测试
```

---

## 📊 测试场景覆盖

### 自动化测试（test-windows.ps1）
- ✅ 环境检查（6 项）
- ✅ hosts 文件系统（2 项）
- ✅ 应用配置目录（2 项）
- ✅ 编译测试（1 项）
- ✅ UAC 相关（2 项）
- **总计**: 13 项自动化测试

### UAC 提权测试（test_uac.exe）
- ✅ 权限检测
- ✅ UAC 提权流程
- ✅ hosts 文件读取
- ✅ 临时文件操作
- ✅ 错误处理
- **总计**: 5 项功能测试

### 手动测试（检查清单）
- ✅ 应用启动（3 项）
- ✅ 读取 hosts（3 项）
- ✅ UAC 提权（8 项）
- ✅ 连续应用（3 项）
- ✅ 备份恢复（5 项）
- ✅ 错误处理（3 项）
- ✅ 编码格式（3 项）
- ✅ 性能测试（3 项）
- ✅ 版本历史（4 项）
- ✅ 兼容性（3 项）
- **总计**: 38 项手动测试用例

---

## ⚠️ 重要提示

### 测试权限
- **推荐**: 使用标准用户账户进行测试（以验证 UAC 提权）
- **可选**: 使用管理员账户运行（部分功能可能不需要 UAC）

### UAC 设置
```powershell
# 检查 UAC 状态
Get-ItemProperty -Path HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System -Name EnableLUA

# 如果 UAC 被禁用，启用它（需要重启）
Set-ItemProperty -Path HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System -Name EnableLUA -Value 1
```

### 防火墙/杀毒软件
某些杀毒软件可能会阻止应用修改 hosts 文件，测试时请：
- 暂时禁用杀毒软件实时保护
- 或将应用添加到白名单

---

## 🐛 问题排查

### 常见问题

#### 1. UAC 提示不出现
**原因**: 应用已经以管理员身份运行
**解决**: 以标准用户身份运行应用

#### 2. 编译失败
**原因**: 缺少依赖或 CGO 问题
**解决**:
```powershell
go clean -cache -modcache
go mod download
go build -v -o wails3-hosts.exe .
```

#### 3. hosts 文件被占用
**原因**: 其他程序正在使用 hosts 文件
**解决**: 关闭占用程序（如记事本、杀毒软件等）

---

## 📝 测试报告示例

完成测试后，请按以下格式报告：

```markdown
## Windows 平台测试报告

**测试日期**: 2025-01-20
**测试人员**: 张三
**Windows 版本**: Windows 11 23H2
**应用版本**: 1.0.0

### 自动化测试结果

| 测试项 | 结果 |
|--------|------|
| 环境检查 | ✅ 6/6 通过 |
| hosts 文件 | ✅ 2/2 通过 |
| 配置目录 | ✅ 2/2 通过 |
| 编译测试 | ✅ 通过 |
| UAC 相关 | ✅ 2/2 通过 |

### 手动测试结果

| 类别 | 通过率 |
|------|--------|
| 应用启动 | 3/3 (100%) |
| 读取 hosts | 3/3 (100%) |
| UAC 提权 | 8/8 (100%) |
| ... | ... |

### 发现的问题

无

### 总体评价

✅ 通过 - 所有关键功能正常，可以发布
```

---

## 📚 相关文档

- [Windows 测试指南](./windows-testing-guide.md) - 详细的测试步骤和场景
- [测试检查清单](./windows-test-checklist.md) - 手动测试用例
- [OpenSpec 提案](../openspec/changes/add-windows-support/) - Windows 支持设计文档

---

## ✅ 下一步

### 立即可做
1. 🧪 在 Windows 环境运行 `test-windows.ps1`
2. 🔧 编译并运行 `test_uac.exe`
3. 📝 按照检查清单进行手动测试
4. 🐤 推送代码触发 GitHub Actions CI

### 后续工作
1. 根据测试结果修复发现的问题
2. 优化性能和用户体验
3. 补充单元测试覆盖率
4. 准备发布说明

---

**准备完成！可以开始在 Windows 平台上进行测试了。** 🎉
