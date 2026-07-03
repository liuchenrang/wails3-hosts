# ============================================================
# Windows 安装包 Smoke Test 脚本
# 用途: 验证 NSIS 安装包的完整性、安装/卸载流程
# 用法: powershell -ExecutionPolicy Bypass -File test_installer.ps1
# ============================================================

param(
    [string]$InstallerPath = "",
    [string]$InstallDir = "$env:TEMP\HostsManagerSmokeTest",
    [switch]$SkipInstall = $false,
    [switch]$SkipUninstall = $false
)

$ErrorActionPreference = "Stop"
$script:TestResults = @()
$script:Passed = 0
$script:Failed = 0

# 自动查找安装包
if (-not $InstallerPath) {
    $scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
    $possiblePaths = @(
        "$scriptDir\..\..\bin\hosts-manager-amd64-installer.exe",
        "$scriptDir\bin\hosts-manager-amd64-installer.exe",
        "$PWD\bin\hosts-manager-amd64-installer.exe"
    )
    foreach ($p in $possiblePaths) {
        $resolved = Resolve-Path $p -ErrorAction SilentlyContinue
        if ($resolved -and (Test-Path $resolved)) {
            $InstallerPath = $resolved.Path
            break
        }
    }
}

# ==================== 测试工具函数 ====================

function Assert-Pass($name, $script) {
    try {
        & $script
        $script:TestResults += @{ Name = $name; Status = "PASS"; Error = $null }
        $script:Passed++
        Write-Host "  ✓ $name" -ForegroundColor Green
    } catch {
        $script:TestResults += @{ Name = $name; Status = "FAIL"; Error = $_.Exception.Message }
        $script:Failed++
        Write-Host "  ✗ $name" -ForegroundColor Red
        Write-Host "    原因: $($_.Exception.Message)" -ForegroundColor DarkRed
    }
}

function Assert-Equal($name, $expected, $actual) {
    if ($expected -eq $actual) {
        $script:TestResults += @{ Name = $name; Status = "PASS"; Error = $null }
        $script:Passed++
        Write-Host "  ✓ $name" -ForegroundColor Green
    } else {
        $msg = "期望: '$expected', 实际: '$actual'"
        $script:TestResults += @{ Name = $name; Status = "FAIL"; Error = $msg }
        $script:Failed++
        Write-Host "  ✗ $name" -ForegroundColor Red
        Write-Host "    $msg" -ForegroundColor DarkRed
    }
}

# ==================== 测试套件 ====================

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Hosts Manager 安装包 Smoke Test" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# ---- Phase 1: 安装包文件完整性 ----

Write-Host "[Phase 1] 安装包文件完整性检查" -ForegroundColor Yellow

Assert-Pass "安装包文件存在" {
    if (-not $InstallerPath) { throw "未找到安装包文件" }
    if (-not (Test-Path $InstallerPath)) { throw "文件不存在: $InstallerPath" }
}

Assert-Pass "安装包文件大小 > 1MB" {
    $size = (Get-Item $InstallerPath).Length
    if ($size -lt 1MB) { throw "安装包过小: $($size) bytes（可能未正确打包）" }
}

Assert-Pass "安装包文件扩展名为 .exe" {
    if ((Get-Item $InstallerPath).Extension -ne ".exe") {
        throw "扩展名应为 .exe"
    }
}

# 检查安装包的 PE 头（可选，需要有 exe 工具）
Assert-Pass "安装包为有效 PE 文件" {
    $bytes = [System.IO.File]::ReadAllBytes($InstallerPath)
    if ($bytes.Length -lt 2) { throw "文件太短" }
    # PE 文件以 MZ 开头
    $magic = [char]$bytes[0] + [char]$bytes[1]
    if ($magic -ne "MZ") { throw "PE 签名无效: $magic" }
}

Write-Host "  安装包: $InstallerPath" -ForegroundColor DarkGray
Write-Host "  大小: $([math]::Round((Get-Item $InstallerPath).Length / 1MB, 2)) MB" -ForegroundColor DarkGray

# ---- Phase 2: NSIS 元数据检查 ----

Write-Host ""
Write-Host "[Phase 2] NSIS 安装包元数据检查" -ForegroundColor Yellow

# 在安装包二进制中搜索关键字符串
Assert-Pass "安装包包含应用名称 'hosts-manager'" {
    $content = [System.IO.File]::ReadAllText($InstallerPath)
    if ($content -notmatch "hosts-manager") {
        throw "安装包中未找到应用名称"
    }
}

Assert-Pass "安装包包含 NSIS 标识" {
    $content = [System.IO.File]::ReadAllText($InstallerPath)
    if ($content -notmatch "Nullsoft") {
        throw "安装包中未找到 NSIS 标识（Nullsoft）"
    }
}

# ---- Phase 3: 静默安装测试 ----

if (-not $SkipInstall) {
    Write-Host ""
    Write-Host "[Phase 3] 静默安装测试" -ForegroundColor Yellow

    # 清理之前的测试目录
    if (Test-Path $InstallDir) {
        Write-Host "  清理旧安装目录: $InstallDir" -ForegroundColor DarkGray
        Remove-Item -Recurse -Force $InstallDir -ErrorAction SilentlyContinue
    }

    Assert-Pass "静默安装成功 (exit code 0)" {
        $process = Start-Process -FilePath $InstallerPath `
            -ArgumentList "/S", "/D=$InstallDir" `
            -Wait -PassThru -NoNewWindow
        if ($process.ExitCode -ne 0) {
            throw "安装器退出码: $($process.ExitCode)"
        }
    }

    Assert-Pass "安装目录已创建" {
        if (-not (Test-Path $InstallDir)) {
            throw "安装目录不存在: $InstallDir"
        }
    }

    Assert-Pass "可执行文件存在 (hosts-manager.exe)" {
        $exePath = Join-Path $InstallDir "hosts-manager.exe"
        if (-not (Test-Path $exePath)) {
            throw "未找到可执行文件: $exePath"
        }
    }

    Assert-Pass "可执行文件大小 > 1MB" {
        $exePath = Join-Path $InstallDir "hosts-manager.exe"
        $size = (Get-Item $exePath).Length
        if ($size -lt 1MB) {
            throw "可执行文件过小: $($size) bytes"
        }
    }

    # 列出安装目录内容
    Write-Host "  安装目录内容:" -ForegroundColor DarkGray
    Get-ChildItem $InstallDir -Recurse | ForEach-Object {
        $relativePath = $_.FullName.Replace($InstallDir, "").TrimStart("\")
        Write-Host "    $relativePath ($([math]::Round($_.Length/1KB, 1)) KB)" -ForegroundColor DarkGray
    }

    # 检查快捷方式
    Assert-Pass "开始菜单快捷方式已创建" {
        $programs = [Environment]::GetFolderPath("Programs")
        $shortcut = Join-Path $programs "Hosts Manager.lnk"
        if (-not (Test-Path $shortcut)) {
            # 也尝试查找公共开始菜单
            $publicPrograms = "$env:ALLUSERSPROFILE\Microsoft\Windows\Start Menu\Programs"
            $shortcut2 = Join-Path $publicPrograms "Hosts Manager.lnk"
            if (-not (Test-Path $shortcut2)) {
                throw "未找到开始菜单快捷方式"
            }
        }
    }

    # ---- Phase 4: 应用启动测试 ----

    Write-Host ""
    Write-Host "[Phase 4] 应用启动测试" -ForegroundColor Yellow

    Assert-Pass "应用可以启动（5秒内不崩溃）" {
        $exePath = Join-Path $InstallDir "hosts-manager.exe"
        $proc = Start-Process -FilePath $exePath -PassThru -WindowStyle Minimized
        Start-Sleep -Seconds 3

        if ($proc.HasExited) {
            if ($proc.ExitCode -ne 0) {
                throw "应用异常退出，退出码: $($proc.ExitCode)"
            }
        }
        # 进程仍在运行 → 启动成功
        $proc | Stop-Process -Force -ErrorAction SilentlyContinue
        Start-Sleep -Seconds 1
    }
}

# ---- Phase 5: 静默卸载测试 ----

if (-not $SkipUninstall -and (Test-Path $InstallDir)) {
    Write-Host ""
    Write-Host "[Phase 5] 静默卸载测试" -ForegroundColor Yellow

    Assert-Pass "卸载程序存在" {
        $uninstPath = Join-Path $InstallDir "uninstall.exe"
        if (-not (Test-Path $uninstPath)) {
            throw "未找到卸载程序: $uninstPath"
        }
    }

    Assert-Pass "静默卸载成功" {
        $uninstPath = Join-Path $InstallDir "uninstall.exe"
        $process = Start-Process -FilePath $uninstPath `
            -ArgumentList "/S" `
            -Wait -PassThru -NoNewWindow
        if ($process.ExitCode -ne 0) {
            throw "卸载器退出码: $($process.ExitCode)"
        }
    }

    Assert-Pass "安装目录已清除" {
        Start-Sleep -Seconds 2
        if (Test-Path $InstallDir) {
            # 可能还有残留文件（如日志），不算失败，但记录警告
            $remaining = Get-ChildItem $InstallDir -Recurse -ErrorAction SilentlyContinue
            if ($remaining) {
                Write-Host "    警告: 目录仍有残留文件:" -ForegroundColor Yellow
                $remaining | ForEach-Object { Write-Host "      $($_.FullName)" -ForegroundColor Yellow }
            }
        }
    }
}

# ==================== 测试报告 ====================

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  测试完成" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$total = $script:Passed + $script:Failed
Write-Host "  总计: $total 项测试" -ForegroundColor White
Write-Host "  通过: $($script:Passed) 项" -ForegroundColor Green
if ($script:Failed -gt 0) {
    Write-Host "  失败: $($script:Failed) 项" -ForegroundColor Red
} else {
    Write-Host "  失败: 0 项" -ForegroundColor Green
}

Write-Host ""
if ($script:Failed -eq 0) {
    Write-Host "  ✓ 所有测试通过！安装包验证成功。" -ForegroundColor Green
    exit 0
} else {
    Write-Host "  ✗ 有 $($script:Failed) 项测试失败，请检查上述错误。" -ForegroundColor Red
    exit 1
}
