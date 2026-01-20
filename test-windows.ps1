# Windows 平台自动化测试脚本
# 需要: 管理员权限运行
# 用途: 验证 Windows 环境配置和应用功能

param(
    [switch]$Verbose,
    [switch]$SkipBuild
)

# 颜色输出函数
function Write-ColorOutput($ForegroundColor) {
    $fc = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    if ($args) {
        Write-Output $args
    }
    $host.UI.RawUI.ForegroundColor = $fc
}

function Write-Success { Write-ColorOutput Green @"+ $args" }
function Write-Error { Write-ColorOutput Red @"- $args" }
function Write-Warning { Write-ColorOutput Yellow @"! $args" }
function Write-Info { Write-ColorOutput Cyan @"i $args" }

Write-Info "========================================"
Write-Info "Windows Hosts Manager 自动化测试套件"
Write-Info "========================================"

$testResults = @()
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path

# 测试 1: 环境检查
Write-Host "`n[测试 1] 环境检查" -ForegroundColor Yellow

# 1.1 检查 Windows 版本
Write-Host "`n  1.1 检查 Windows 版本..." -ForegroundColor Cyan
$osVersion = [System.Environment]::OSVersion.Version
$osName = (Get-WmiObject -Class Win32_OperatingSystem).Caption
Write-Host "     操作系统: $osName" -ForegroundColor White
Write-Host "     版本号: $osVersion" -ForegroundColor White

if ($osVersion.Major -ge 10) {
    Write-Success "Windows 版本符合要求 (10+)"
    $testResults += @{ Test = "Windows 版本"; Result = "通过" }
} else {
    Write-Error "Windows 版本不符合要求，需要 Windows 10+"
    $testResults += @{ Test = "Windows 版本"; Result = "失败" }
}

# 1.2 检查管理员权限
Write-Host "`n  1.2 检查当前权限..." -ForegroundColor Cyan
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if ($isAdmin) {
    Write-Success "当前具有管理员权限"
    $testResults += @{ Test = "管理员权限"; Result = "通过" }
} else {
    Write-Warning "当前是标准用户（测试 UAC 提权需要）"
    $testResults += @{ Test = "管理员权限"; Result = "标准用户" }
}

# 1.3 检查 Go 环境
Write-Host "`n  1.3 检查 Go 环境..." -ForegroundColor Cyan
try {
    $goVersionOutput = go version 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Go 已安装: $goVersionOutput"
        $testResults += @{ Test = "Go 环境"; Result = "通过" }
    } else {
        Write-Error "Go 未正确安装"
        $testResults += @{ Test = "Go 环境"; Result = "失败" }
    }
} catch {
    Write-Error "Go 未安装或不在 PATH 中"
    $testResults += @{ Test = "Go 环境"; Result = "失败" }
}

# 1.4 检查 Node.js 环境
Write-Host "`n  1.4 检查 Node.js 环境..." -ForegroundColor Cyan
try {
    $nodeVersion = node --version 2>&1
    if ($LASTEXITCODE -eq 0) {
        $npmVersion = npm --version 2>&1
        Write-Success "Node.js 已安装: $nodeVersion, npm: $npmVersion"
        $testResults += @{ Test = "Node.js 环境"; Result = "通过" }
    } else {
        Write-Error "Node.js 未正确安装"
        $testResults += @{ Test = "Node.js 环境"; Result = "失败" }
    }
} catch {
    Write-Error "Node.js 未安装或不在 PATH 中"
    $testResults += @{ Test = "Node.js 环境"; Result = "失败" }
}

# 测试 2: hosts 文件系统检查
Write-Host "`n[测试 2] hosts 文件系统检查" -ForegroundColor Yellow

# 2.1 检查 hosts 文件路径
Write-Host "`n  2.1 检查 hosts 文件路径..." -ForegroundColor Cyan
$hostsPath = "$env:SystemRoot\System32\drivers\etc\hosts"
Write-Host "     路径: $hostsPath" -ForegroundColor White

if (Test-Path $hostsPath) {
    Write-Success "hosts 文件存在"
    $fileInfo = Get-Item $hostsPath
    Write-Host "     大小: $($fileInfo.Length) 字节" -ForegroundColor White
    Write-Host "     修改时间: $($fileInfo.LastWriteTime)" -ForegroundColor White

    # 读取内容前几行
    $content = Get-Content $hostsPath -First 5
    Write-Host "     前 5 行内容:" -ForegroundColor White
    $content | ForEach-Object { Write-Host "       $_" -ForegroundColor Gray }

    $testResults += @{ Test = "hosts 文件存在"; Result = "通过" }
} else {
    Write-Error "hosts 文件不存在"
    $testResults += @{ Test = "hosts 文件存在"; Result = "失败" }
}

# 2.2 检查 hosts 文件权限
Write-Host "`n  2.2 检查 hosts 文件权限..." -ForegroundColor Cyan
try {
    $acl = Get-Acl $hostsPath
    Write-Info "  文件所有者: $($acl.Owner)"
    Write-Success "成功读取文件 ACL"

    # 检查当前用户是否有写入权限
    $currentUser = [Security.Principal.WindowsIdentity]::GetCurrent().Name
    $hasWriteAccess = $false

    foreach ($access in $acl.Access) {
        if ($access.IdentityReference -like "*$currentUser*" -or $access.IdentityReference -like "*Administrators*") {
            if ($access.FileSystemRights -match "Write") {
                $hasWriteAccess = $true
                break
            }
        }
    }

    if ($hasWriteAccess -or $isAdmin) {
        Write-Success "当前用户具有写入权限（或可通过 UAC 提权）"
        $testResults += @{ Test = "hosts 文件权限"; Result = "通过" }
    } else {
        Write-Warning "当前用户没有直接写入权限，需要 UAC 提权"
        $testResults += @{ Test = "hosts 文件权限"; Result = "需要 UAC" }
    }
} catch {
    Write-Error "无法读取文件权限: $_"
    $testResults += @{ Test = "hosts 文件权限"; Result = "失败" }
}

# 测试 3: 应用配置目录检查
Write-Host "`n[测试 3] 应用配置目录检查" -ForegroundColor Yellow

# 3.1 检查配置目录
Write-Host "`n  3.1 检查配置目录..." -ForegroundColor Cyan
$configDir = "$env:APPDATA\hosts-manager"
Write-Host "     路径: $configDir" -ForegroundColor White

if (Test-Path $configDir) {
    Write-Success "配置目录存在"
    $testResults += @{ Test = "配置目录"; Result = "存在" }

    # 列出目录内容
    $items = Get-ChildItem $configDir
    Write-Host "     包含文件: $($items.Count) 个" -ForegroundColor White
    if ($Verbose) {
        $items | ForEach-Object {
            Write-Host "       $($_.Name) ($($_.Length) 字节)" -ForegroundColor Gray
        }
    }
} else {
    Write-Warning "配置目录不存在（首次运行时正常）"
    $testResults += @{ Test = "配置目录"; Result = "不存在" }
}

# 3.2 检查备份目录
Write-Host "`n  3.2 检查备份目录..." -ForegroundColor Cyan
$backupDir = "$configDir\backups"
Write-Host "     路径: $backupDir" -ForegroundColor White

if (Test-Path $backupDir) {
    Write-Success "备份目录存在"
    $backups = Get-ChildItem $backupDir -Filter "hosts_*.bak"
    Write-Host "     备份文件数量: $($backups.Count)" -ForegroundColor White

    if ($backups.Count -gt 0) {
        Write-Host "     最近的 3 个备份:" -ForegroundColor White
        $backups | Sort-Object LastWriteTime -Descending | Select-Object -First 3 | ForEach-Object {
            Write-Host "       $($_.Name) - $($_.LastWriteTime)" -ForegroundColor Gray
        }
    }

    $testResults += @{ Test = "备份目录"; Result = "存在" }
} else {
    Write-Warning "备份目录不存在（首次运行时正常）"
    $testResults += @{ Test = "备份目录"; Result = "不存在" }
}

# 测试 4: 编译测试
Write-Host "`n[测试 4] 编译测试" -ForegroundColor Yellow

if (-not $SkipBuild) {
    Write-Host "`n  4.1 尝试编译应用..." -ForegroundColor Cyan

    $buildDir = Join-Path $scriptPath "build"
    if (-not (Test-Path $buildDir)) {
        New-Item -ItemType Directory -Path $buildDir | Out-Null
    }

    $outputExe = Join-Path $buildDir "wails3-hosts-test.exe"

    Write-Host "     编译命令: go build -o $outputExe ." -ForegroundColor Gray

    Push-Location $scriptPath
    $buildOutput = go build -o $outputExe . 2>&1
    $buildExitCode = $LASTEXITCODE
    Pop-Location

    if ($buildExitCode -eq 0) {
        if (Test-Path $outputExe) {
            $fileInfo = Get-Item $outputExe
            Write-Success "编译成功: $($fileInfo.Name) ($($fileInfo.Length) 字节)"
            $testResults += @{ Test = "应用编译"; Result = "通过" }
        } else {
            Write-Error "编译命令执行成功但输出文件不存在"
            $testResults += @{ Test = "应用编译"; Result = "失败" }
        }
    } else {
        Write-Error "编译失败"
        if ($Verbose) {
            Write-Host "     编译输出:" -ForegroundColor Red
            $buildOutput | ForEach-Object { Write-Host "       $_" -ForegroundColor Red }
        }
        $testResults += @{ Test = "应用编译"; Result = "失败" }
    }
} else {
    Write-Warning "跳过编译测试（-SkipBuild 参数）"
}

# 测试 5: UAC 相关检查
Write-Host "`n[测试 5] UAC 相关检查" -ForegroundColor Yellow

# 5.1 检查 UAC 状态
Write-Host "`n  5.1 检查 UAC 状态..." -ForegroundColor Cyan
try {
    $uacEnabled = (Get-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System" -ErrorAction Stop).EnableLUA
    if ($uacEnabled -eq 1) {
        Write-Success "UAC 已启用"
        $testResults += @{ Test = "UAC 状态"; Result = "启用" }
    } else {
        Write-Warning "UAC 已禁用（测试提权功能需要启用 UAC）"
        $testResults += @{ Test = "UAC 状态"; Result = "禁用" }
    }
} catch {
    Write-Error "无法读取 UAC 状态"
    $testResults += @{ Test = "UAC 状态"; Result = "未知" }
}

# 5.2 检查 UAC 提权级别
Write-Host "`n  5.2 检查 UAC 提权提示级别..." -ForegroundColor Cyan
try {
    $consentPrompt = (Get-ItemProperty -Path "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System" -ErrorAction Stop).ConsentPromptBehaviorAdmin
    switch ($consentPrompt) {
        0 { Write-Info "  提权级别: 不提示（不安全）" }
        1 { Write-Info "  提权级别: 在安全桌面提示" }
        2 { Write-Info "  提权级别: 在用户桌面提示" }
        3 { Write-Info "  提权级别: 不提示（不安全）" }
        4 { Write-Info "  提权级别: 在安全桌面提示（需凭据）" }
        5 { Write-Info "  提权级别: 在用户桌面提示（需凭据）" }
        default { Write-Info "  提权级别: 未知 ($consentPrompt)" }
    }
    $testResults += @{ Test = "UAC 提权级别"; Result = "$consentPrompt" }
} catch {
    Write-Error "无法读取 UAC 提权级别"
    $testResults += @{ Test = "UAC 提权级别"; Result = "未知" }
}

# 测试总结
Write-Host "`n========================================" -ForegroundColor Cyan
Write-Info "测试总结"
Write-Host "========================================" -ForegroundColor Cyan

$passed = ($testResults | Where-Object { $_.Result -eq "通过" }).Count
$failed = ($testResults | Where-Object { $_.Result -eq "失败" }).Count
$total = $testResults.Count

Write-Host "`n总计: $total 个测试" -ForegroundColor White
Write-Success "通过: $passed 个"
if ($failed -gt 0) {
    Write-Error "失败: $failed 个"
} else {
    Write-Host "失败: 0 个" -ForegroundColor Green
}

Write-Host "`n详细结果:" -ForegroundColor White
$testResults | ForEach-Object {
    $resultStr = $_.Result
    $color = if ($resultStr -eq "通过") { "Green" } elseif ($resultStr -eq "失败") { "Red" } else { "Yellow" }
    Write-ColorOutput $color "  [$($resultStr)] $($_.Test)"
}

Write-Host "`n========================================" -ForegroundColor Cyan

if ($failed -eq 0) {
    Write-Success "所有关键测试通过！"
    Write-Info "可以开始手动功能测试"
    exit 0
} else {
    Write-Error "存在失败的测试，请检查环境配置"
    exit 1
}
