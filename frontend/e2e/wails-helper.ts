import { _electron as electron } from 'playwright';
import { test as base, expect } from '@playwright/test';
import path from 'path';

type WailsAppFixture = {
  app: Awaited<ReturnType<typeof electron.launch>>;
  window: Awaited<ReturnType<typeof electron.launch>>['windows'][0];
};

// 扩展 Playwright test 以支持 Wails 应用
export const test = base.extend<WailsAppFixture>({
  app: async ({}, use) => {
    // 编译后的应用路径
    const appPath = path.resolve(__dirname, '../../hosts_manager.exe');

    // 启动应用
    const app = await electron.launch({
      path: appPath,
      // 可以添加其他参数
      args: [],
    });

    await use(app);

    // 清理
    await app.close();
  },

  window: async ({ app }, use) => {
    // 获取主窗口
    const window = await app.windows()[0];
    await window.waitForLoadState('domcontentloaded');
    await use(window);
  },
});

export { expect };

/**
 * 使用示例：
 *
 * test('测试应用启动', async ({ window }) => {
 *   await expect(window.locator('h1')).toHaveText('Hosts Manager');
 * });
 */