export const jaJP = {
  translation: {
    // 共通
    common: {
      save: '保存',
      cancel: 'キャンセル',
      delete: '削除',
      edit: '編集',
      add: '追加',
      confirm: '確認',
      loading: '読み込み中...',
      success: '成功',
      error: 'エラー',
      search: '検索',
    },
    // アプリ
    app: {
      title: 'Hosts Manager',
      subtitle: 'クロスプラットフォーム hosts ファイル管理ツール',
    },
    // サイドバー
    sidebar: {
      groups: 'グループ',
      groupCount: '{{count}} グループ',
      createGroup: 'グループ作成',
      deleteGroup: 'グループ削除',
      editGroup: 'グループ編集',
      noGroups: 'グループがありません',
      createFirstGroup: '上のボタンをクリックして最初のグループを作成',
      groupActions: 'グループ操作',
      deleteConfirm: 'グループ"{{name}}"を削除しますか？',
      moreEntries: '他 {{count}} エントリ...',
    },
    // メインパネル
    mainPanel: {
      selectGroup: 'グループを選択または作成してください',
      entries: 'Hosts エントリ',
      addEntry: 'エントリ追加',
      preview: 'hosts ファイルプレビュー',
      apply: '設定を適用',
      applyShortcut: 'ショートカット: ⌘S / Ctrl+S',
      reset: 'リセット',
      validation: {
        title: 'フォーマットエラー',
        invalidFormat: '無効なフォーマットです。IP アドレス + ホスト名である必要があります',
        invalidIP: '無効な IP アドレス形式: {{ip}}',
        invalidHostname: '無効なホスト名形式: {{hostname}}',
      },
    },
    // フォーム
    form: {
      groupName: 'グループ名',
      groupDesc: 'グループ説明',
      ipAddress: 'IP アドレス',
      hostname: 'ホスト名',
      comment: 'コメント（オプション）',
      enabled: '有効',
    },
    // バージョン履歴
    versions: {
      title: 'バージョン履歴',
      rollback: 'ロールバック',
      rollbackConfirm: 'このバージョンにロールバックしますか？',
      versionInfo: 'バージョン情報',
      timestamp: '時間',
      description: '説明',
      source: 'ソース',
      source_manual: '手動',
      source_auto: '自動',
      source_rollback: 'ロールバック',
    },
    // Sudo
    sudo: {
      title: '管理者権限が必要です',
      description: 'hosts ファイルの変更には管理者権限が必要です',
      password: 'パスワードを入力してください',
      passwordPlaceholder: 'sudo パスワードを入力',
      passwordCached: 'キャッシュ済み',
      passwordOptionalPlaceholder: 'パスワードはキャッシュされています。確認をクリックするか、再度パスワードを入力してください',
      passwordCachedHint: 'sudo パスワードはシステムにキャッシュされています。再入力不要です。確認をクリックしてください',
      validateError: 'パスワード検証に失敗しました',
      cached: 'パスワードキャッシュ済み ({{seconds}}秒)',
      required: 'hosts ファイルを書き込むには sudo パスワードが必要です',
      invalid: 'sudo パスワードの検証に失敗しました。パスワードが正しいか確認してください',
      windowsUACHint: 'この操作では UAC プロンプトが表示されます。「許可」をクリックして続行してください',
      windowsPlatformHint: 'Windows プラットフォーム、UAC プロンプトが表示されます',
    },
    // テーマ
    theme: {
      light: 'ライトテーマ',
      dark: 'ダークテーマ',
      toggle: 'テーマ切り替え',
    },
    // エラー
    errors: {
      loadFailed: '読み込みに失敗しました',
      saveFailed: '保存に失敗しました',
      invalidIP: '無効な IP アドレス形式',
      invalidHostname: '無効なホスト名形式',
      duplicateEntry: 'エントリが既に存在します',
      networkError: 'ネットワークエラー',
      unknownError: '不明なエラー',
    },
    // 競合検出
    conflicts: {
      title: '設定の競合',
      description: '以下のホスト名に複数のマッピングがあります：',
      hostname: 'ホスト名',
      ips: 'IP アドレス',
      ignore: '無視して続行',
    },
    // 私たちについて
    about: {
      title: '私たちについて',
      version: 'バージョン',
      email: '連絡先メール',
      description: 'シンプルで効率的な hosts ファイル管理ツール',
      website: 'ウェブサイトを訪問',
    },
  },
}
