#!/bin/bash
# 实际测试密码是否泄露到文件

set -e

echo "=== 测试密码是否泄露到文件 ==="
echo ""

# 创建临时测试文件
TEST_FILE="/tmp/test_hosts_$$.txt"
TEST_PASSWORD="MySecretPassword123"
TEST_CONTENT="# Test hosts file
127.0.0.1 localhost
127.0.0.1 test.local"

# 清理函数
cleanup() {
    rm -f "$TEST_FILE"
}
trap cleanup EXIT

echo "测试文件: $TEST_FILE"
echo "测试密码: $TEST_PASSWORD"
echo ""

# 测试 sudo 命令
echo "测试 1: 使用 echo | sudo tee 方法 (旧方法 - 会泄露)"
echo "$TEST_PASSWORD" | sudo -S tee "$TEST_FILE" > /dev/null << EOF
$TEST_CONTENT
EOF

echo "文件内容:"
cat "$TEST_FILE"
echo ""

if grep -q "$TEST_PASSWORD" "$TEST_FILE"; then
    echo "❌ 发现密码泄露！"
else
    echo "✓ 未发现密码"
fi

rm -f "$TEST_FILE"

echo ""
echo "测试 2: 使用 cat 方法 (新方法 - 不泄露)"
{
    echo "$TEST_PASSWORD"
    echo "$TEST_CONTENT"
} | sudo -S sh -c "cat > $TEST_FILE"

echo "文件内容:"
cat "$TEST_FILE"
echo ""

if grep -q "$TEST_PASSWORD" "$TEST_FILE"; then
    echo "❌ 发现密码泄露！"
    exit 1
else
    echo "✓ 未发现密码 - 方法正确！"
fi

echo ""
echo "=== 测试完成 ==="
