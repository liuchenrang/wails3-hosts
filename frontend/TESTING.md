# UI 自动化测试指南

本项目支持两种测试方式：

## 1. 单元测试 (Vitest + React Testing Library)

用于测试单个 React 组件的渲染和行为。

### 运行测试

```bash
cd frontend

# 运行所有测试
npm run test

# 监听模式（开发时使用）
npm run test:watch

# UI 模式（可视化调试）
npm run test:ui
```

### 测试文件位置

- `frontend/test/**/*.test.tsx` - 单元测试
- `frontend/test/setup.ts` - 测试环境配置

### 示例测试

```tsx
import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { Button } from '../../src/components/ui/Button';

describe('Button 组件', () => {
  it('应该正确渲染', () => {
    render(<Button>点击我</Button>);
    expect(screen.getByRole('button')).toHaveTextContent('点击我');
  });

  it('应该支持 onClick', async () => {
    const handleClick = vi.fn();
    const user = userEvent.setup();
    render(<Button onClick={handleClick}>点击</Button>);
    await user.click(screen.getByRole('button'));
    expect(handleClick).toHaveBeenCalled();
  });
});
```

## 2. E2E 测试 (Playwright)

用于测试完整的用户交互流程。

### 安装浏览器

```bash
cd frontend
npm run playwright:install
```

### 运行 E2E 测试

```bash
# 方式1: 测试前端开发服务器（需要 mock 后端）
npm run test:e2e

# 方式2: 测试完整应用（需要先编译）
# 1. 构建前端: npm run build
# 2. 编译应用: go build -o hosts_manager.exe .
# 3. 运行测试: npm run test:e2e
```

### 测试文件位置

- `frontend/e2e/**/*.spec.ts` - E2E 测试
- `frontend/playwright.config.ts` - Playwright 配置

### 测试 data-testid

以下组件添加了 `data-testid` 便于定位：

- `[data-testid="sidebar"]` - 侧边栏
- `[data-testid="group-item-{id}"]` - 分组项
- `[data-testid="main-panel"]` - 主面板

## 测试策略

### 推荐测试内容

| 类型 | 测试内容 | 工具 |
|------|---------|------|
| 单元测试 | UI 组件渲染 | Vitest + RTL |
| 单元测试 | 工具函数 | Vitest |
| E2E 测试 | 关键用户流程 | Playwright |
| E2E 测试 | 跨组件交互 | Playwright |

### 测试覆盖的功能

1. **分组管理**
   - 创建、编辑、删除分组
   - 切换分组启用状态
   - 分组排序（拖拽）

2. **条目管理**
   - 添加、编辑、删除 hosts 条目
   - 批量更新条目

3. **配置应用**
   - 生成预览
   - 应用配置（含权限验证）
   - 版本历史与回滚

4. **UI 交互**
   - 主题切换
   - 模态框操作
   - Toast 提示

## Mock 后端 API

单元测试时需要 mock Wails 后端调用：

```typescript
// frontend/test/setup.ts
vi.mock('./api/hosts', () => ({
  hostsApi: {
    getAllGroups: vi.fn().mockResolvedValue([...]),
    createGroup: vi.fn().mockResolvedValue({...}),
    // ...
  }
}));
```