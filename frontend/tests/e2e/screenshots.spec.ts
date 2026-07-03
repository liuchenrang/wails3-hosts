import { test } from '@playwright/test';
import path from 'path';

test('截图 - 主界面初始状态', async ({ page }) => {
  await page.goto('/');

  // 等待页面加载
  await page.waitForTimeout(2000);

  // 截取主界面
  await page.screenshot({
    path: path.join(process.cwd(), '../doc/articles/screenshots/01-main-interface.png'),
    fullPage: true
  });
});

test('截图 - 模拟添加分组对话框', async ({ page }) => {
  await page.goto('/');

  // 等待页面加载
  await page.waitForTimeout(2000);

  // 使用 JavaScript 直接注入模拟对话框
  await page.evaluate(() => {
    const dialogHTML = `
      <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" id="mock-dialog">
        <div class="bg-white dark:bg-slate-900 rounded-lg shadow-lg p-6 w-[400px]">
          <h2 class="text-lg font-semibold mb-4">新建分组</h2>
          <div class="space-y-4">
            <div>
              <label class="text-sm font-medium mb-1 block">分组名称</label>
              <input type="text" placeholder="例如：开发环境" class="w-full px-3 py-2 border rounded-md focus:ring-2 focus:ring-blue-500" value="开发环境">
            </div>
            <div>
              <label class="text-sm font-medium mb-1 block">分组描述</label>
              <input type="text" placeholder="例如：本地开发域名映射" class="w-full px-3 py-2 border rounded-md focus:ring-2 focus:ring-blue-500" value="本地开发域名映射配置">
            </div>
          </div>
          <div class="flex justify-end gap-2 mt-6">
            <button class="px-4 py-2 text-sm border rounded-md hover:bg-slate-100 dark:hover:bg-slate-800">取消</button>
            <button class="px-4 py-2 text-sm bg-blue-500 text-white rounded-md hover:bg-blue-600">确认</button>
          </div>
        </div>
      </div>
    `;
    document.body.insertAdjacentHTML('beforeend', dialogHTML);
  });

  // 截取对话框界面
  await page.screenshot({
    path: path.join(process.cwd(), '../doc/articles/screenshots/02-add-group-dialog.png'),
    fullPage: true
  });
});

test('截图 - 模拟版本历史对话框', async ({ page }) => {
  await page.goto('/');

  // 等待页面加载
  await page.waitForTimeout(2000);

  // 使用 JavaScript 直接注入模拟版本历史对话框
  await page.evaluate(() => {
    const dialogHTML = `
      <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" id="mock-version-dialog">
        <div class="bg-white dark:bg-slate-900 rounded-lg shadow-lg p-6 w-[500px] max-h-[400px]">
          <h2 class="text-lg font-semibold mb-4">版本历史</h2>
          <div class="space-y-3 overflow-y-auto">
            <div class="flex items-center justify-between p-3 bg-slate-50 dark:bg-slate-800 rounded">
              <div>
                <div class="font-medium">2026-04-10 11:15</div>
                <div class="text-xs text-slate-500">应用开发环境配置</div>
              </div>
              <button class="px-3 py-1 text-xs bg-blue-500 text-white rounded hover:bg-blue-600">回滚</button>
            </div>
            <div class="flex items-center justify-between p-3 bg-slate-50 dark:bg-slate-800 rounded">
              <div>
                <div class="font-medium">2026-04-09 18:30</div>
                <div class="text-xs text-slate-500">添加测试环境分组</div>
              </div>
              <button class="px-3 py-1 text-xs bg-blue-500 text-white rounded hover:bg-blue-600">回滚</button>
            </div>
            <div class="flex items-center justify-between p-3 bg-slate-50 dark:bg-slate-800 rounded">
              <div>
                <div class="font-medium">2026-04-08 09:00</div>
                <div class="text-xs text-slate-500">初始配置</div>
              </div>
              <button class="px-3 py-1 text-xs bg-blue-500 text-white rounded hover:bg-blue-600">回滚</button>
            </div>
          </div>
          <div class="flex justify-end mt-4">
            <button class="px-4 py-2 text-sm border rounded-md hover:bg-slate-100 dark:hover:bg-slate-800">关闭</button>
          </div>
        </div>
      </div>
    `;
    document.body.insertAdjacentHTML('beforeend', dialogHTML);
  });

  // 截取版本历史界面
  await page.screenshot({
    path: path.join(process.cwd(), '../doc/articles/screenshots/03-version-history.png'),
    fullPage: true
  });
});

test('截图 - 模拟完整工作界面', async ({ page }) => {
  await page.goto('/');

  // 等待页面加载
  await page.waitForTimeout(2000);

  // 使用 JavaScript 模拟完整界面状态
  await page.evaluate(() => {
    // 找到分组列表区域并添加模拟分组
    const mockGroups = [
      { name: '开发环境', desc: '本地开发域名映射', enabled: true, entries: ['127.0.0.1 local.api', '127.0.0.1 local.web'] },
      { name: '测试环境', desc: '测试服务器配置', enabled: true, entries: ['192.168.1.100 test.server'] },
      { name: '生产环境', desc: '生产服务器域名', enabled: false, entries: ['10.0.0.1 prod.server'] },
    ];

    // 模拟分组卡片
    const groupCardsHTML = mockGroups.map(g => `
      <div class="flex items-center gap-3 p-3 rounded-lg bg-slate-50 dark:bg-slate-800 cursor-pointer hover:bg-slate-100 dark:hover:bg-slate-700 transition-colors border border-slate-200 dark:border-slate-700 mb-2">
        <input type="checkbox" ${g.enabled ? 'checked' : ''} class="w-4 h-4 rounded accent-blue-500">
        <div class="flex-1 min-w-0">
          <div class="font-medium text-sm truncate">${g.name}</div>
          <div class="text-xs text-slate-500 truncate">${g.desc}</div>
        </div>
        <span class="text-xs px-2 py-1 rounded ${g.enabled ? 'bg-green-100 text-green-700 dark:bg-green-800 dark:text-green-200' : 'bg-slate-100 text-slate-500'}">${g.enabled ? '已启用' : '未启用'}</span>
      </div>
    `).join('');

    // 找到主容器并注入内容
    const mainContainer = document.querySelector('.flex.h-screen');
    if (mainContainer) {
      // 添加应用按钮区域
      const actionBarHTML = `
        <div class="flex items-center justify-between px-4 py-2 bg-slate-100 dark:bg-slate-800 border-b border-slate-200 dark:border-slate-700 mb-4">
          <div class="text-sm text-slate-500">共 ${mockGroups.filter(g => g.enabled).length} 个分组已启用</div>
          <button class="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 text-sm font-medium">
            应用配置
          </button>
        </div>
      `;

      // 添加 hosts 编辑区域
      const editorHTML = `
        <div class="p-4">
          <div class="mb-2 text-sm font-medium">hosts 条目编辑</div>
          <textarea class="w-full h-32 p-3 border rounded-md bg-slate-50 dark:bg-slate-900 text-sm font-mono resize-none focus:ring-2 focus:ring-blue-500" placeholder="输入 hosts 条目...">${mockGroups[0].entries.join('\n')}</textarea>
          <div class="mt-2 flex justify-end gap-2">
            <button class="px-3 py-1.5 text-sm border rounded hover:bg-slate-100 dark:hover:bg-slate-800">添加条目</button>
            <button class="px-3 py-1.5 text-sm bg-green-500 text-white rounded hover:bg-green-600">保存</button>
          </div>
        </div>
      `;
    }

    // 直接在分组列表区域添加内容
    const sidebarContent = document.querySelector('[class*="flex-col"][class*="gap"]');
    if (sidebarContent) {
      const listArea = sidebarContent.querySelector('div:nth-child(2)');
      if (listArea) {
        listArea.innerHTML = groupCardsHTML;
      }
    }

    // 更新分组计数
    const countText = document.querySelector('[class*="text-sm"][class*="text-slate"]');
    if (countText) {
      countText.textContent = `${mockGroups.length} 个分组`;
    }
  });

  await page.waitForTimeout(500);

  // 截取完整工作界面
  await page.screenshot({
    path: path.join(process.cwd(), '../doc/articles/screenshots/04-full-work-interface.png'),
    fullPage: true
  });
});