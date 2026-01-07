# 🔐 新的密码管理策略 - 前端集成指南

## 策略概述

**核心改变**: 不再通过 `ApplyHosts()` 传递密码，而是在前端分两步处理：

1. **验证密码阶段**: 调用 `ValidateSudoPassword(password)` 验证并缓存密码
2. **应用配置阶段**: 调用 `ApplyHosts()` 应用配置（无需密码参数）

## 新的 API 设计

### 后端 API 变更

#### 1. `ValidateSudoPassword(password)` - 验证密码
```javascript
// 验证密码并提升 sudo 权限
const [valid, error] = await ValidateSudoPassword(password);

if (valid) {
  // 密码验证成功，sudo 权限已提升
  // 后续的 ApplyHosts 不需要密码
}
```

**返回值**:
- `valid`: boolean - 密码是否有效
- `error`: string - 错误信息（如果验证失败）

#### 2. `IsSudoPasswordCached()` - 检查缓存状态
```javascript
const isCached = await IsSudoPasswordCached();

if (!isCached) {
  // 显示密码输入框
  showPasswordDialog();
}
```

#### 3. `ApplyHosts()` - 应用配置（无密码参数）
```javascript
// 注意：不再传递密码！
try {
  await ApplyHosts();
  showSuccess('配置已应用');
} catch (error) {
  if (error.message.includes('需要先验证 sudo 密码')) {
    // 提示用户需要先验证密码
    showPasswordDialog();
  }
}
```

## 前端实现流程

### 完整的使用流程

```typescript
// 1. 用户点击"应用配置"按钮
async function handleApplyConfig() {
  // 2. 检查是否有缓存的 sudo 凭证
  const isCached = await window.hosts.IsSudoPasswordCached();

  if (!isCached) {
    // 3. 如果没有缓存，弹出密码输入框
    const password = await showPasswordDialog();

    if (!password) {
      return; // 用户取消
    }

    // 4. 验证密码
    const [valid, error] = await window.hosts.ValidateSudoPassword(password);

    if (!valid) {
      showError('sudo 密码验证失败: ' + error);
      return;
    }

    showSuccess('sudo 密码验证成功');
  }

  // 5. 应用配置（不需要密码）
  try {
    await window.hosts.ApplyHosts();
    showSuccess('配置已应用到系统');
  } catch (error) {
    showError('应用配置失败: ' + error.message);
  }
}
```

### React 组件示例

```tsx
import { useState } from 'react';
import { hosts } from '../api/hosts';

export function ApplyConfigButton() {
  const [showPasswordDialog, setShowPasswordDialog] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const handleApply = async () => {
    setIsLoading(true);

    try {
      // 检查密码缓存状态
      const isCached = await hosts.IsSudoPasswordCached();

      if (!isCached) {
        // 需要密码
        setShowPasswordDialog(true);
        return;
      }

      // 直接应用配置
      await hosts.ApplyHosts();
      alert('配置已成功应用！');

    } catch (error: any) {
      if (error.message.includes('需要先验证 sudo 密码')) {
        setShowPasswordDialog(true);
      } else {
        alert('应用配置失败: ' + error.message);
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handlePasswordSubmit = async (password: string) => {
    setIsLoading(true);
    setShowPasswordDialog(false);

    try {
      // 验证密码
      const [valid, error] = await hosts.ValidateSudoPassword(password);

      if (!valid) {
        alert('密码验证失败: ' + error);
        return;
      }

      // 密码验证成功，应用配置
      await hosts.ApplyHosts();
      alert('配置已成功应用！');

    } catch (error: any) {
      alert('操作失败: ' + error.message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div>
      <button
        onClick={handleApply}
        disabled={isLoading}
      >
        {isLoading ? '处理中...' : '应用配置'}
      </button>

      {showPasswordDialog && (
        <PasswordDialog
          onSubmit={handlePasswordSubmit}
          onCancel={() => setShowPasswordDialog(false)}
        />
      )}
    </div>
  );
}

// 密码输入对话框组件
function PasswordDialog({
  onSubmit,
  onCancel
}: {
  onSubmit: (password: string) => void;
  onCancel: () => void;
}) {
  const [password, setPassword] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(password);
  };

  return (
    <div className="modal-overlay">
      <div className="modal">
        <h2>需要管理员权限</h2>
        <p>修改 hosts 文件需要 sudo 权限，请输入密码：</p>

        <form onSubmit={handleSubmit}>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="sudo 密码"
            autoFocus
            required
          />

          <div className="button-group">
            <button type="button" onClick={onCancel}>
              取消
            </button>
            <button type="submit">
              确认
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
```

## 优势

### 1. 安全性提升
- ✅ 密码只在验证时传递一次
- ✅ 不会通过多个函数调用传递密码
- ✅ 利用系统的 sudo 缓存机制
- ✅ 避免密码在应用层多次传递

### 2. 代码简洁
- ✅ `ApplyHosts()` 不再需要密码参数
- ✅ 密码验证逻辑独立
- ✅ 更符合单一职责原则

### 3. 用户体验
- ✅ 可以在应用配置前提前验证密码
- ✅ 减少重复输入密码的次数
- ✅ 明确的错误提示

## 错误处理

### 常见错误及处理

```typescript
try {
  await hosts.ApplyHosts();
} catch (error) {
  if (error.message.includes('需要先验证 sudo 密码')) {
    // 密码未验证或已过期
    showPasswordDialog();
  } else if (error.message.includes('sudo 密码验证失败')) {
    // 密码错误
    showError('密码错误，请重试');
  } else if (error.message.includes('写入 hosts 文件失败')) {
    // 权限或其他系统错误
    showError('系统错误，请检查权限');
  } else {
    // 其他错误
    showError('操作失败: ' + error.message);
  }
}
```

## 测试清单

### 功能测试

- [ ] 未验证密码时点击应用，显示密码输入框
- [ ] 输入错误密码，显示错误提示
- [ ] 输入正确密码，配置成功应用
- [ ] 验证密码后5分钟内再次应用，不需要输入密码
- [ ] 5分钟后再次应用，提示重新输入密码

### 安全测试

- [ ] 密码不会写入 /etc/hosts 文件
- [ ] 密码不会出现在日志中
- [ ] 密码不会被保存到配置文件

## 迁移指南

### 从旧版本迁移

**旧代码**:
```typescript
// ❌ 旧方式：传递密码
await hosts.ApplyHosts(password);
```

**新代码**:
```typescript
// ✅ 新方式：先验证，再应用
const [valid] = await hosts.ValidateSudoPassword(password);
if (valid) {
  await hosts.ApplyHosts();
}
```

---

**版本**: v3.0
**更新日期**: 2026-01-07
**状态**: ✅ 生产就绪
