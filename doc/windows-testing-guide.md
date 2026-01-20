# Windows 平台测试指南

## 测试环境要求

### 硬件要求
- **处理器**: x64 (AMD64)
- **内存**: 至少 4GB RAM
- **磁盘**: 至少 1GB 可用空间

### 软件要求
- **操作系统**: Windows 10 (版本 1607+) 或 Windows 11
- **Go**: 1.24 或更高版本
- **Node.js**: 18+ (用于前端构建)
- **Build Tools**: (可选) Visual Studio Build Tools

### 权限要求
- **标准用户账户** (用于测试 UAC 提权)
- **管理员账户** (用于某些测试场景)

---

## 测试准备

### 1. 环境设置

#### 安装 Go
```powershell
# 下载并安装 Go
# https://golang.org/dl/

# 验证安装
go version
```

#### 安装 Node.js
```powershell
# 使用 winget 安装
winget install OpenJS.NodeJS.LTS

# 验证安装
node --version
npm --version
```

#### 克隆并构建项目
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

---

## 测试场景

### 场景 1: 基础功能测试

#### 1.1 读取系统 hosts 文件
**目的**: 验证应用能正确读取 Windows hosts 文件路径

**步骤**:
1. 以标准用户身份运行应用
2. 查看当前 hosts 内容
3. 验证路径为 `C:\Windows\System32\drivers\etc\hosts`

**预期结果**:
- ✅ 成功读取 hosts 文件内容
- ✅ 路径正确
- ✅ 内容显示正常（包括 Windows 特有的注释）

---

### 场景 2: UAC 提权测试

#### 2.1 首次应用配置（UAC 提示）
**目的**: 验证 UAC 提权流程正常工作

**前置条件**:
- 当前用户不是管理员
- 或以标准用户身份运行应用

**步骤**:
1. 创建一个新分组
2. 添加一个 hosts 条目（例如：`127.0.0.1 test.local`）
3. 启用分组
4. 点击"应用配置"按钮
5. 观察 UAC 提示窗口

**预期结果**:
- ✅ 弹出 UAC 提示窗口
- ✅ 提示内容清晰："需要管理员权限来修改 hosts 文件"
- ✅ 显示应用的数字签名（如果有）

**如果点击"是"**:
- ✅ 应用以管理员权限执行写入
- ✅ 写入成功后返回主窗口
- ✅ 显示"配置已成功应用"提示
- ✅ hosts 文件实际被修改

**如果点击"否"**:
- ✅ 写入操作中止
- ✅ 显示友好提示："您取消了管理员权限确认，无法应用 hosts 配置"
- ✅ 应用保持正常运行

---

#### 2.2 连续多次应用配置
**目的**: 验证每次操作都需要 UAC 确认（无缓存）

**步骤**:
1. 应用配置后，立即再次点击"应用配置"
2. 观察 UAC 提示

**预期结果**:
- ✅ 再次弹出 UAC 提示（Windows 不支持凭据缓存）
- ✅ 前端不显示"密码已缓存"提示

---

### 场景 3: 文件操作测试

#### 3.1 备份功能
**目的**: 验证 Windows 平台的备份功能正常

**步骤**:
1. 应用配置
2. 检查备份目录：`%APPDATA%\hosts-manager\backups\`

**预期结果**:
- ✅ 备份目录自动创建
- ✅ 备份文件命名格式：`hosts_YYYYMMDD_HHMMSS.bak`
- ✅ 备份内容与原 hosts 文件一致
- ✅ 自动清理超过 5 个的旧备份

---

#### 3.2 从备份恢复
**目的**: 验证从备份恢复功能

**步骤**:
1. 修改 hosts 文件
2. 应用配置（创建备份）
3. 再次修改 hosts
4. 选择之前的备份进行恢复

**预期结果**:
- ✅ 弹出 UAC 提示
- ✅ 确认后 hosts 文件恢复到备份状态
- ✅ 显示恢复成功提示

---

### 场景 4: 错误处理测试

#### 4.1 hosts 文件被占用
**目的**: 验证文件占用时的错误处理

**步骤**:
1. 使用记事本打开 `C:\Windows\System32\drivers\etc\hosts`（以管理员身份）
2. 保持文件打开
3. 尝试应用配置

**预期结果**:
- ✅ 弹出 UAC 提示
- ✅ 确认后显示错误："hosts 文件正在被其他程序使用"
- ✅ 建议关闭占用程序（如记事本、杀毒软件等）
- ✅ 应用不崩溃

---

#### 4.2 临时文件创建失败
**目的**: 验证临时目录不可用时的错误处理

**步骤**:
1. 模拟临时目录不可写（如磁盘已满）
2. 尝试应用配置

**预期结果**:
- ✅ 显示错误："无法创建临时文件，请检查磁盘空间和权限"
- ✅ 操作安全中止
- ✅ 不留下垃圾临时文件

---

#### 4.3 取消 UAC 提示
**目的**: 验证用户取消操作的处理

**步骤**:
1. 点击"应用配置"
2. 在 UAC 提示中点击"否"

**预期结果**:
- ✅ 显示取消提示："您取消了管理员权限确认"
- ✅ hosts 文件未被修改
- ✅ 应用状态正常

---

### 场景 5: 编码和格式测试

#### 5.1 UTF-8 BOM 处理
**目的**: 验证正确处理 Windows hosts 文件的 UTF-8 BOM

**步骤**:
1. 使用带 UTF-8 BOM 的编辑器修改 hosts 文件
2. 在应用中读取文件

**预期结果**:
- ✅ 正确读取并移除 BOM 标记
- ✅ 内容显示正常
- ✅ 写入时可以选择是否添加 BOM

---

#### 5.2 换行符处理
**目的**: 验证正确处理 Windows 换行符（\r\n）

**步骤**:
1. 查看读取的 hosts 内容
2. 应用配置
3. 用记事本打开 hosts 文件

**预期结果**:
- ✅ 显示的换行符正确
- ✅ 写入的文件使用 \r\n (Windows 标准格式)
- ✅ 记事本显示正常

---

### 场景 6: 性能测试

#### 6.1 UAC 提权耗时
**目的**: 测量 UAC 提权和进程重启的开销

**步骤**:
1. 点击"应用配置"
2. 从 UAC 提示出现到写入完成计时

**预期结果**:
- ✅ 总耗时在 1-3 秒内
- ✅ 用户体验可接受

---

#### 6.2 大文件写入
**目的**: 测试大量 hosts 条目的写入性能

**步骤**:
1. 创建包含 1000+ 条目的大文件
2. 应用配置
3. 测量写入时间

**预期结果**:
- ✅ 写入成功
- ✅ 性能可接受（不超过 5 秒）

---

## 自动化测试

### PowerShell 测试脚本

保存为 `test-windows.ps1`:

```powershell
# Windows 平台自动化测试脚本
# 需要: 管理员权限运行

Write-Host "================================" -ForegroundColor Cyan
Write-Host "Windows Hosts Manager 测试套件" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan

# 测试 1: 检查 hosts 文件路径
Write-Host "`n[测试 1] 检查 hosts 文件路径..." -ForegroundColor Yellow
$hostsPath = "$env:SystemRoot\System32\drivers\etc\hosts"
if (Test-Path $hostsPath) {
    Write-Host "✅ hosts 文件路径正确: $hostsPath" -ForegroundColor Green
} else {
    Write-Host "❌ hosts 文件路径不正确" -ForegroundColor Red
}

# 测试 2: 检查应用配置目录
Write-Host "`n[测试 2] 检查应用配置目录..." -ForegroundColor Yellow
$configDir = "$env:APPDATA\hosts-manager"
if (Test-Path $configDir) {
    Write-Host "✅ 配置目录存在: $configDir" -ForegroundColor Green
} else {
    Write-Host "⚠️  配置目录不存在（首次运行正常）" -ForegroundColor Yellow
}

# 测试 3: 检查备份目录
Write-Host "`n[测试 3] 检查备份目录..." -ForegroundColor Yellow
$backupDir = "$configDir\backups"
if (Test-Path $backupDir) {
    $backups = Get-ChildItem $backupDir -Filter "hosts_*.bak"
    Write-Host "✅ 备份目录存在，包含 $($backups.Count) 个备份" -ForegroundColor Green
} else {
    Write-Host "⚠️  备份目录不存在（首次运行正常）" -ForegroundColor Yellow
}

# 测试 4: 检查管理员权限
Write-Host "`n[测试 4] 检查当前权限..." -ForegroundColor Yellow
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if ($isAdmin) {
    Write-Host "✅ 当前具有管理员权限" -ForegroundColor Green
} else {
    Write-Host "ℹ️  当前是标准用户（UAC 提权需要）" -ForegroundColor Cyan
}

# 测试 5: 检查 Go 环境
Write-Host "`n[测试 5] 检查 Go 环境..." -ForegroundColor Yellow
try {
    $goVersion = go version
    Write-Host "✅ Go 已安装: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Go 未安装或不在 PATH 中" -ForegroundColor Red
}

# 测试 6: 检查 Node.js 环境
Write-Host "`n[测试 6] 检查 Node.js 环境..." -ForegroundColor Yellow
try {
    $nodeVersion = node --version
    Write-Host "✅ Node.js 已安装: $nodeVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Node.js 未安装或不在 PATH 中" -ForegroundColor Red
}

Write-Host "`n================================" -ForegroundColor Cyan
Write-Host "环境检查完成" -ForegroundColor Cyan
Write-Host "================================" -ForegroundColor Cyan
```

运行测试：
```powershell
# 以管理员身份运行 PowerShell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
.\test-windows.ps1
```

---

## 常见问题排查

### 问题 1: UAC 提示不出现
**可能原因**:
- 应用已经以管理员身份运行
- UAC 被禁用

**解决方法**:
```powershell
# 检查 UAC 状态
Get-ItemProperty -Path HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System -Name EnableLUA

# 如果被禁用，启用 UAC（需要重启）
Set-ItemProperty -Path HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System -Name EnableLUA -Value 1
```

---

### 问题 2: 写入失败
**可能原因**:
- hosts 文件被占用
- 权限不足
- 临时文件创建失败

**排查步骤**:
1. 检查日志输出
2. 使用 Process Explorer 查找占用文件的进程
3. 检查临时目录权限

---

### 问题 3: 编译错误
**可能原因**:
- 缺少 Windows API 绑定
- CGO 相关问题

**解决方法**:
```powershell
# 清理缓存
go clean -cache -modcache

# 重新下载依赖
go mod download

# 重新编译
go build -v -o wails3-hosts.exe .
```

---

## 测试报告模板

完成测试后，请填写以下报告：

```markdown
## Windows 平台测试报告

**测试日期**: YYYY-MM-DD
**测试人员**: [姓名]
**Windows 版本**: [版本号]
**应用版本**: [版本号]

### 测试结果摘要

| 场景 | 结果 | 备注 |
|------|------|------|
| 读取 hosts 文件 | ✅/❌ | |
| UAC 提权 | ✅/❌ | |
| 备份功能 | ✅/❌ | |
| 从备份恢复 | ✅/❌ | |
| 错误处理 | ✅/❌ | |
| 编码处理 | ✅/❌ | |
| 性能测试 | ✅/❌ | |

### 发现的问题

1. **问题描述**
   - 重现步骤:
   - 预期结果:
   - 实际结果:
   - 截图/日志:

### 建议和改进

- [ ] 建议内容

### 总体评价

- [ ] 通过 / 不通过
```

---

## 附录

### A. Windows API 参考
- [ShellExecuteW 函数](https://docs.microsoft.com/en-us/windows/win32/api/shellapi/nf-shellapi-shellexecutew)
- [UAC 技术参考](https://docs.microsoft.com/en-us/windows/win32/secauth/user-account-control)

### B. 相关工具
- **Process Explorer**: 查看文件占用情况
- **Process Monitor**: 监控文件系统活动
- **Resource Monitor**: 查看资源使用情况

### C. 测试数据
示例 hosts 条目用于测试：
```
# Windows Hosts Manager 测试条目
127.0.0.1 test.local
127.0.0.1 dev.local
127.0.0.1 staging.local

::1 localhost
```
