import '@testing-library/jest-dom';

// Mock Wails runtime for unit tests
const mockHostsApi = {
  getAllGroups: vi.fn().mockResolvedValue([
    { id: '1', name: '默认分组', description: '', isEnabled: true, entries: [] },
  ]),
  createGroup: vi.fn().mockResolvedValue({ id: '2', name: '新分组', description: '' }),
  updateGroup: vi.fn().mockResolvedValue(undefined),
  deleteGroup: vi.fn().mockResolvedValue(undefined),
  toggleGroup: vi.fn().mockResolvedValue(undefined),
  addEntry: vi.fn().mockResolvedValue(undefined),
  updateEntry: vi.fn().mockResolvedValue(undefined),
  deleteEntry: vi.fn().mockResolvedValue(undefined),
  batchUpdateEntries: vi.fn().mockResolvedValue(undefined),
  generatePreview: vi.fn().mockResolvedValue('# Hosts file content\n127.0.0.1 localhost'),
  detectConflicts: vi.fn().mockResolvedValue({}),
  applyHosts: vi.fn().mockResolvedValue(undefined),
  validateSudoPassword: vi.fn().mockResolvedValue({ valid: true }),
  isSudoPasswordCached: vi.fn().mockResolvedValue(false),
  getPlatformInfo: vi.fn().mockResolvedValue({
    os: 'windows',
    arch: 'amd64',
    needsSudo: false,
    canCacheCred: false,
  }),
  getVersions: vi.fn().mockResolvedValue([]),
  rollbackToVersion: vi.fn().mockResolvedValue(undefined),
  reorderGroups: vi.fn().mockResolvedValue(undefined),
};

// Mock window EventsOn for Wails events
if (typeof window !== 'undefined') {
  (window as any).EventsOn = vi.fn();
}

// Export mock for use in tests
global.mockHostsApi = mockHostsApi;