#!/bin/bash
# 诊断 hosts 文件写入问题

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║          🔍 hosts 文件写入问题诊断工具                          ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

# 颜色
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}步骤 1: 检查配置文件位置${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

CONFIG_DIR="$HOME/Library/Application Support/hosts-manager"
CONFIG_FILE="$CONFIG_DIR/config.json"

if [ -f "$CONFIG_FILE" ]; then
    echo -e "${GREEN}✓ 配置文件存在: $CONFIG_FILE${NC}"
    echo ""
    echo "配置文件内容:"
    cat "$CONFIG_FILE" | python3 -m json.tool 2>/dev/null || cat "$CONFIG_FILE"
else
    echo -e "${YELLOW}⚠️  配置文件不存在: $CONFIG_FILE${NC}"
    echo "这可能表示您还没有创建任何分组。"
fi

echo ""
echo -e "${BLUE}步骤 2: 检查当前 hosts 文件${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ -f "/etc/hosts" ]; then
    echo "当前 /etc/hosts 文件最后 20 行:"
    echo "..."
    tail -20 /etc/hosts
else
    echo -e "${RED}✗ hosts 文件不存在${NC}"
fi

echo ""
echo -e "${BLUE}步骤 3: 检查 hosts 文件备份${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

BACKUP_DIR="$HOME/Library/Application Support/hosts-manager/backups"
if [ -d "$BACKUP_DIR" ]; then
    BACKUP_COUNT=$(ls -1 "$BACKUP_DIR"/*.bak 2>/dev/null | wc -l)
    echo -e "${GREEN}✓ 找到 $BACKUP_COUNT 个备份文件${NC}"

    if [ $BACKUP_COUNT -gt 0 ]; then
        echo ""
        echo "最新的备份文件:"
        ls -lt "$BACKUP_DIR"/*.bak | head -1 | awk '{print $NF}'
        echo ""
        echo "内容预览:"
        ls -lt "$BACKUP_DIR"/*.bak | head -1 | awk '{print $NF}' | xargs tail -20
    fi
else
    echo -e "${YELLOW}⚠️  备份目录不存在${NC}"
fi

echo ""
echo -e "${BLUE}步骤 4: 常见问题检查${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 检查配置文件中的分组启用状态
if [ -f "$CONFIG_FILE" ]; then
    # 检查是否有启用的分组
    ENABLED_GROUPS=$(cat "$CONFIG_FILE" | grep -o '"is_enabled":true' | wc -l)
    TOTAL_GROUPS=$(cat "$CONFIG_FILE" | grep -o '"groups":\[' | wc -l)

    echo "分组统计:"
    echo "  总分组数: $(cat "$CONFIG_FILE" | grep -o '"id"' | wc -l)"
    echo "  启用的分组: $ENABLED_GROUPS"

    if [ $ENABLED_GROUPS -eq 0 ]; then
        echo ""
        echo -e "${RED}✗ 问题发现: 没有启用的分组！${NC}"
        echo ""
        echo "这就是为什么 hosts 文件没有内容的原因。"
        echo ""
        echo "解决方法:"
        echo "  1. 打开应用"
        echo "  2. 在左侧分组列表中，点击分组左侧的复选框"
        echo "  3. 确保分组显示为✓选中状态"
        echo "  4. 然后点击'应用配置'按钮"
        echo ""
        echo "或者:"
        echo "  在代码中设置 group.IsEnabled = true"
    else
        echo -e "${GREEN}✓ 有 $ENABLED_GROUPS 个分组已启用${NC}"

        # 检查启用的分组是否有条目
        echo ""
        echo "检查启用分组的条目数..."
        cat "$CONFIG_FILE" | grep -A 20 '"is_enabled":true' | grep -o '"entry"' | wc -l
    fi
fi

echo ""
echo -e "${BLUE}步骤 5: 手动测试命令${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "如果需要手动测试 sudo 权限:"
echo "  sudo echo 'test' > /dev/null && echo 'sudo 权限正常' || echo 'sudo 权限失败'"

echo ""
echo "╔════════════════════════════════════════════════════════════════╗"
echo -e "${GREEN}║                  诊断完成                                    ║${NC}"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""
echo "最可能的原因:"
echo "  1. ✓ 分组已创建"
echo "  2. ✗ 分组未启用 (is_enabled = false)"
echo "  3. ✓ hosts 条目已添加"
echo ""
echo "解决步骤:"
echo "  1. 在应用中点击分组左侧的复选框启用分组"
echo "  2. 确保复选框显示为选中状态 (✓)"
echo "  3. 点击'应用配置'按钮"
echo "  4. 输入 sudo 密码"
echo ""
