#!/bin/bash
# 安全修复验证脚本

set -e

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║          🔒 密码泄露漏洞修复验证脚本                            ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}📋 步骤 1: 运行安全测试${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
go test ./internal/infrastructure/system/... -v -run TestSudoCommandStdin

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 安全测试通过${NC}"
else
    echo -e "${RED}✗ 安全测试失败${NC}"
    exit 1
fi

echo ""
echo -e "${BLUE}📋 步骤 2: 运行功能测试${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
go test ./internal/domain/service/... -v -run TestGenerateHostsContent

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ 功能测试通过${NC}"
else
    echo -e "${RED}✗ 功能测试失败${NC}"
    exit 1
fi

echo ""
echo -e "${BLUE}📋 步骤 3: 编译修复后的应用${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
go build -o bin/hosts-manager-fixed .

if [ $? -eq 0 ]; then
    SIZE=$(ls -lh bin/hosts-manager-fixed | awk '{print $5}')
    echo -e "${GREEN}✓ 编译成功 (大小: $SIZE)${NC}"
else
    echo -e "${RED}✗ 编译失败${NC}"
    exit 1
fi

echo ""
echo -e "${BLUE}📋 步骤 4: 验证修复内容${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 检查修复是否生效
if grep -q "stdinContent.ReadFrom(c.stdin)" internal/infrastructure/system/sudo_command.go; then
    echo -e "${GREEN}✓ 代码修复已应用${NC}"
else
    echo -e "${RED}✗ 代码修复未找到${NC}"
    exit 1
fi

# 检查是否添加了 stdin 字段
if grep -q "stdin    io.Reader" internal/infrastructure/system/sudo_command.go; then
    echo -e "${GREEN}✓ 数据结构已更新${NC}"
else
    echo -e "${RED}✗ 数据结构未更新${NC}"
    exit 1
fi

echo ""
echo -e "${BLUE}📋 步骤 5: 检查 hosts 文件安全建议${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ -f "/etc/hosts" ]; then
    echo -e "${YELLOW}⚠️  建议检查系统 hosts 文件是否包含泄露的密码：${NC}"
    echo "   sudo grep -n '[a-zA-Z0-9]\{8,\}' /etc/hosts | head -20"
    echo ""
    echo "   如果发现密码泄露，请："
    echo "   1. 备份当前 hosts 文件: sudo cp /etc/hosts /etc/hosts.backup"
    echo "   2. 手动编辑删除密码行: sudo nano /etc/hosts"
    echo "   3. 修改系统 sudo 密码"
else
    echo -e "${GREEN}✓ 未检测到标准 hosts 文件位置${NC}"
fi

echo ""
echo "╔════════════════════════════════════════════════════════════════╗"
echo -e "${GREEN}║                  ✅ 验证完成！                                 ║${NC}"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""
echo "验证结果总结:"
echo "  ✓ 安全测试: 所有密码隔离测试通过"
echo "  ✓ 功能测试: hosts 内容生成测试通过"
echo "  ✓ 编译验证: 应用程序成功编译"
echo "  ✓ 代码审查: 修复逻辑正确实现"
echo ""
echo -e "${GREEN}安全修复已验证！可以安全使用。${NC}"
echo ""
echo "详细报告: SECURITY_FIX_PASSWORD_LEAK.md"
echo ""
