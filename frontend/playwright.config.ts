import { defineConfig, devices } from '@playwright/test';

/**
 * Playwright 配置 - Wails 桌面应用 E2E 测试
 *
 * 对于 Wails 应用，有两种测试方式：
 * 1. 直接测试编译后的应用（需要先编译应用）
 * 2. 通过开发服务器测试前端（需要 mock 后端 API）
 */
export default defineConfig({
  testDir: './e2e',
  fullyParallel: false, // 桌面应用测试通常不并行
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: 1, // 桌面应用测试通常单进程
  reporter: 'html',

  use: {
    baseURL: 'http://localhost:34115', // Wails 开发服务器默认端口
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
  },

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],

  // 方式1: 启动 Wails 开发模式（wails3 dev）
  // webServer: {
  //   command: 'wails3 dev',
  //   url: 'http://localhost:34115',
  //   reuseExistingServer: !process.env.CI,
  //   timeout: 120000,
  // },

  // 方式2: 使用已构建的前端（需要先 npm run build）
  // webServer: {
  //   command: 'npm run preview',
  //   url: 'http://localhost:4173',
  //   reuseExistingServer: true,
  // },
});