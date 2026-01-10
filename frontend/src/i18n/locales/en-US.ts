export const enUS = {
  translation: {
    // Common
    common: {
      save: 'Save',
      cancel: 'Cancel',
      delete: 'Delete',
      edit: 'Edit',
      add: 'Add',
      confirm: 'Confirm',
      loading: 'Loading...',
      success: 'Success',
      error: 'Error',
      search: 'Search',
    },
    // App
    app: {
      title: 'Hosts Manager',
      subtitle: 'Cross-platform hosts file management tool',
    },
    // Sidebar
    sidebar: {
      groups: 'Groups',
      groupCount: '{{count}} groups',
      createGroup: 'Create Group',
      deleteGroup: 'Delete Group',
      editGroup: 'Edit Group',
      noGroups: 'No groups yet',
      createFirstGroup: 'Click the button above to create the first group',
      groupActions: 'Group Actions',
      deleteConfirm: 'Are you sure you want to delete group "{{name}}"?',
      moreEntries: '{{count}} more entries...',
    },
    // Main Panel
    mainPanel: {
      selectGroup: 'Please select or create a group',
      entries: 'Hosts Entries',
      addEntry: 'Add Entry',
      preview: 'Preview hosts file',
      apply: 'Apply Configuration',
      applyShortcut: 'Shortcut: ⌘S / Ctrl+S',
      reset: 'Reset',
      validation: {
        title: 'Format Error',
        invalidFormat: 'Invalid format, should be: IP address + hostname',
        invalidIP: 'Invalid IP address format: {{ip}}',
        invalidHostname: 'Invalid hostname format: {{hostname}}',
      },
    },
    // Form
    form: {
      groupName: 'Group Name',
      groupDesc: 'Group Description',
      ipAddress: 'IP Address',
      hostname: 'Hostname',
      comment: 'Comment (Optional)',
      enabled: 'Enabled',
    },
    // Versions
    versions: {
      title: 'Version History',
      rollback: 'Rollback',
      rollbackConfirm: 'Are you sure you want to rollback to this version?',
      versionInfo: 'Version Info',
      timestamp: 'Time',
      description: 'Description',
      source: 'Source',
      source_manual: 'Manual',
      source_auto: 'Auto',
      source_rollback: 'Rollback',
    },
    // Sudo
    sudo: {
      title: 'Administrator Privileges Required',
      description: 'Modifying hosts file requires administrator privileges',
      password: 'Please enter password',
      passwordPlaceholder: 'Enter sudo password',
      passwordCached: 'Cached',
      passwordOptionalPlaceholder: 'Password cached, you can confirm directly or re-enter password',
      passwordCachedHint: 'Sudo password is cached in the system. No need to re-enter, just click confirm',
      validateError: 'Password validation failed',
      cached: 'Password cached ({{seconds}}s)',
      required: 'Sudo password required to write hosts file',
      invalid: 'Sudo password validation failed, please check if password is correct',
    },
    // Theme
    theme: {
      light: 'Light Theme',
      dark: 'Dark Theme',
      toggle: 'Toggle Theme',
    },
    // Errors
    errors: {
      loadFailed: 'Failed to load',
      saveFailed: 'Failed to save',
      invalidIP: 'Invalid IP address format',
      invalidHostname: 'Invalid hostname format',
      duplicateEntry: 'Entry already exists',
      networkError: 'Network error',
      unknownError: 'Unknown error',
    },
    // Conflicts
    conflicts: {
      title: 'Configuration Conflicts',
      description: 'The following hostnames have multiple mappings:',
      hostname: 'Hostname',
      ips: 'IP Addresses',
      ignore: 'Ignore and Continue',
    },
    // About
    about: {
      title: 'About Us',
      version: 'Version',
      email: 'Contact Email',
      description: 'A simple and efficient hosts file management tool',
    },
  },
}
