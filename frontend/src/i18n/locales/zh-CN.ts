export const zhCN = {
  translation: {
    // 通用
    common: {
      save: '保存',
      cancel: '取消',
      delete: '删除',
      edit: '编辑',
      add: '添加',
      confirm: '确认',
      loading: '加载中...',
      success: '操作成功',
      error: '操作失败',
      search: '搜索',
    },
    // 应用标题
    app: {
      title: 'Hosts Manager',
      subtitle: '跨平台 hosts 文件管理工具',
    },
    // 侧边栏
    sidebar: {
      groups: '分组列表',
      groupCount: '{{count}} 个分组',
      createGroup: '新建分组',
      deleteGroup: '删除分组',
      editGroup: '编辑分组',
      noGroups: '暂无分组',
      createFirstGroup: '点击上方按钮创建第一个分组',
      groupActions: '分组操作',
      deleteConfirm: '确定要删除分组"{{name}}"吗？',
      moreEntries: '还有 {{count}} 个条目...',
    },
    // 主面板
    mainPanel: {
      selectGroup: '请选择或创建一个分组',
      entries: 'Hosts 条目',
      addEntry: '添加条目',
      preview: '预览生成的 hosts 文件',
      apply: '应用配置',
      applyShortcut: '快捷键: ⌘S / Ctrl+S',
      reset: '重置',
      validation: {
        title: '格式错误',
        invalidFormat: '格式无效，应为：IP 地址 + 域名',
        invalidIP: 'IP 地址格式无效: {{ip}}',
        invalidHostname: '域名格式无效: {{hostname}}',
      },
    },
    // 表单
    form: {
      groupName: '分组名称',
      groupDesc: '分组描述',
      ipAddress: 'IP 地址',
      hostname: '主机名',
      comment: '注释（可选）',
      enabled: '启用',
    },
    // 版本历史
    versions: {
      title: '版本历史',
      rollback: '回滚',
      rollbackConfirm: '确定要回滚到此版本吗？',
      versionInfo: '版本信息',
      timestamp: '时间',
      description: '描述',
      source: '来源',
      source_manual: '手动应用',
      source_auto: '自动创建',
      source_rollback: '回滚操作',
    },
    // Sudo 密码
    sudo: {
      title: '需要管理员权限',
      description: '修改 hosts 文件需要管理员权限',
      password: '请输入密码',
      passwordPlaceholder: '输入 sudo 密码',
      validateError: '密码验证失败',
      cached: '密码已缓存 ({{seconds}}秒)',
      required: '需要 sudo 密码才能写入 hosts 文件',
      invalid: 'sudo 密码验证失败，请检查密码是否正确',
    },
    // 主题
    theme: {
      light: '明亮主题',
      dark: '暗色主题',
      toggle: '切换主题',
    },
    // 错误信息
    errors: {
      loadFailed: '加载失败',
      saveFailed: '保存失败',
      invalidIP: 'IP 地址格式无效',
      invalidHostname: '主机名格式无效',
      duplicateEntry: '条目已存在',
      networkError: '网络错误',
      unknownError: '未知错误',
    },
    // 冲突检测
    conflicts: {
      title: '配置冲突',
      description: '以下主机名存在多个映射，可能会导致问题：',
      hostname: '主机名',
      ips: 'IP 地址',
      ignore: '忽略并继续',
    },
    // 关于我们
    about: {
      title: '关于我们',
      version: '版本',
      email: '联系邮箱',
      description: '一个简单高效的 hosts 文件管理工具',
    },
  },
}
