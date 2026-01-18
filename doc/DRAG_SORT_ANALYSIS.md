# 分组拖动排序功能分析

## 📋 现有功能梳理

### 已实现的后端 API
✅ **ReorderGroups API** (`internal/application/service/hosts_app_service.go:188`)
```go
func (s *HostsApplicationService) ReorderGroups(ctx context.Context, req dto.ReorderGroupsRequest) error
```
- 功能：按新顺序更新分组的 Order 字段
- 实现：遍历 GroupIDs 数组，设置索引作为新的 Order 值

### 前端交互事件
当前 Sidebar 组件的交互层：

1. **onClick** → `onSelectGroup(group)` - 选择分组
2. **onDoubleClick** → `onDoubleClickGroup(group)` - 双击分组
3. **Power 按钮点击** → `onToggleGroup(group.id, !group.is_enabled)` - 切换启用状态
4. **编辑按钮** → `handleEditGroup(group)` - 编辑分组
5. **删除按钮** → `handleDeleteClick(group)` - 删除分组

---

## ⚠️ 潜在冲突分析

### 1. **拖动与点击事件的冲突**

**问题描述：**
- 用户开始拖动时，鼠标按下 → 移动 → 释放
- 这个过程会被识别为一次完整的点击事件
- 导致拖动结束后触发 `onSelectGroup`

**解决方案：**
```typescript
const [isDragging, setIsDragging] = useState(false)
const [dragStartPos, setDragStartPos] = useState({ x: 0, y: 0 })

// 在 onMouseDown 时记录起始位置
const handleMouseDown = (e: React.MouseEvent, groupId: string) => {
  setDragStartPos({ x: e.clientX, y: e.clientY })
  setIsDragging(false)
}

// 在 onMouseUp 时判断是否为拖动
const handleMouseUp = (e: React.MouseEvent, group: HostsGroup) => {
  const moveDistance = Math.sqrt(
    Math.pow(e.clientX - dragStartPos.x, 2) +
    Math.pow(e.clientY - dragStartPos.y, 2)
  )

  // 如果移动距离小于 5px，视为点击，否则视为拖动
  if (moveDistance < 5 && !isDragging) {
    onSelectGroup(group)
  }
}

// 在 onDragStart 时设置拖动标志
const handleDragStart = () => {
  setIsDragging(true)
}
```

---

### 2. **拖动与按钮操作的冲突**

**问题描述：**
- 点击 Power、编辑、删除按钮时，可能误触发拖动
- 拖动分组时，可能误点击按钮

**解决方案：**
```tsx
{/* 在所有按钮上阻止拖动事件 */}
<button
  draggable={false}
  onDragStart={(e) => e.stopPropagation()}
  onClick={(e) => {
    e.stopPropagation() // 已有
    // ... 按钮逻辑
  }}
>
  <Power className="h-4 w-4" />
</button>
```

---

### 3. **拖动与双击事件的冲突**

**问题描述：**
- 双击包含两次单击，第一次单击可能被误判为拖动开始

**解决方案：**
```typescript
// 使用拖动手柄（Grip）代替整个卡片可拖动
<div className="flex items-center p-3">
  {/* 拖动手柄 - 仅此区域可拖动 */}
  <div
    draggable={true}
    onDragStart={handleDragStart}
    onDragOver={handleDragOver}
    onDrop={handleDrop}
    className="cursor-grab active:cursor-grabbing mr-2"
  >
    <GripVertical className="h-4 w-4 text-muted-foreground" />
  </div>

  {/* 原有内容 */}
</div>
```

---

### 4. **移动端触摸事件兼容性**

**问题描述：**
- 拖动使用 HTML5 Drag and Drop API
- 移动端不支持此 API，需要使用 Touch Events

**解决方案：**
```typescript
// 使用 @dnd-kit/core 库（推荐）
// 优点：
// - 支持 HTML5 Drag and Drop
// - 支持触摸屏
// - 无障碍访问友好
// - 性能优化

import { DndContext, closestCenter } from '@dnd-kit/core'
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable'
import { useSortable } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
```

---

### 5. **性能问题**

**问题描述：**
- 拖动过程中频繁更新 UI
- 每次拖动结束调用 API 保存

**解决方案：**
```typescript
// 使用乐观更新（Optimistic UI）
const handleDragEnd = async (event: DragEndEvent) => {
  const { active, over } = event

  if (over && active.id !== over.id) {
    // 1. 立即更新本地状态（乐观更新）
    const oldIndex = groups.findIndex(g => g.id === active.id)
    const newIndex = groups.findIndex(g => g.id === over.id)
    const newGroups = arrayMove(groups, oldIndex, newIndex)
    setGroups(newGroups)

    // 2. 异步调用 API 保存
    try {
      await hostsApi.reorderGroups(newGroups.map(g => g.id))
    } catch (error) {
      // 3. 失败时回滚状态
      setGroups(groups)
      toast.error('排序失败，请重试')
    }
  }
}
```

---

## 🎯 推荐实现方案

### 方案 A：最小改动方案（HTML5 Drag & Drop）

**优点：**
- 无需额外依赖
- 改动最小
- 后端 API 已就绪

**缺点：**
- 移动端支持差
- 需要手动处理事件冲突

**实现步骤：**

1. **添加拖动手柄**
```tsx
import { GripVertical } from 'lucide-react'

<div className="flex items-center p-3">
  {/* 拖动手柄 */}
  <div
    draggable={true}
    onDragStart={(e) => handleDragStart(e, group.id)}
    className="cursor-grab active:cursor-grabbing mr-2 opacity-0 group-hover:opacity-100"
  >
    <GripVertical className="h-4 w-4 text-muted-foreground" />
  </div>

  {/* 原有内容 */}
</div>
```

2. **实现拖动事件处理**
```typescript
const [draggedId, setDraggedId] = useState<string | null>(null)

const handleDragStart = (e: React.DragEvent, groupId: string) => {
  setDraggedId(groupId)
  e.dataTransfer.effectAllowed = 'move'
}

const handleDragOver = (e: React.DragEvent) => {
  e.preventDefault() // 允许放置
  e.dataTransfer.dropEffect = 'move'
}

const handleDrop = async (e: React.DragEvent, targetId: string) => {
  e.preventDefault()
  if (!draggedId || draggedId === targetId) return

  const oldIndex = groups.findIndex(g => g.id === draggedId)
  const newIndex = groups.findIndex(g => g.id === targetId)
  const newGroups = arrayMove(groups, oldIndex, newIndex)

  // 乐观更新
  setGroups(newGroups)

  // 调用 API
  try {
    await hostsApi.reorderGroups(newGroups.map(g => g.id))
    toast.success('排序已更新')
  } catch (error) {
    setGroups(groups) // 回滚
    toast.error('排序失败')
  }

  setDraggedId(null)
}
```

---

### 方案 B：专业方案（@dnd-kit）

**优点：**
- 完美的跨平台支持
- 优秀的性能
- 内置碰撞检测

**缺点：**
- 需要安装依赖（~15KB gzipped）
- 学习曲线稍陡

**实现步骤：**

1. **安装依赖**
```bash
npm install @dnd-kit/core @dnd-kit/sortable @dnd-kit/utilities
```

2. **实现可排序组件**
```tsx
import { DndContext } from '@dnd-kit/core'
import { SortableContext, verticalListSortingStrategy, useSortable } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'

function SortableGroup({ group, ...props }) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: group.id })

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  }

  return (
    <div ref={setNodeRef} style={style} {...props}>
      {/* 拖动手柄 */}
      <div {...attributes} {...listeners} className="cursor-grab">
        <GripVertical className="h-4 w-4" />
      </div>
      {/* 原有内容 */}
    </div>
  )
}
```

---

## 📊 方案对比

| 特性 | 方案 A (HTML5) | 方案 B (@dnd-kit) |
|------|----------------|-------------------|
| 依赖大小 | 0 KB | ~15 KB |
| 移动端支持 | ❌ | ✅ |
| 开发时间 | 2-3 小时 | 1-2 小时 |
| 性能 | 良好 | 优秀 |
| 可维护性 | 一般 | 优秀 |
| 无障碍访问 | 需手动处理 | 内置支持 |

---

## 🚀 推荐决策

### 如果你的应用：
- ✅ **仅桌面端使用** → 选择方案 A（HTML5）
- ✅ **需要移动端支持** → 选择方案 B（@dnd-kit）
- ✅ **追求最佳体验** → 选择方案 B（@dnd-kit）

### MVP 建议
先实现**方案 A**，快速验证功能需求，再根据用户反馈决定是否升级到方案 B。

---

## ✅ 实现清单（方案 A）

- [ ] 添加拖动手柄图标（GripVertical）
- [ ] 实现 dragStart、dragOver、drop 事件
- [ ] 优化拖动时的视觉反馈（透明度、阴影）
- [ ] 实现乐观更新逻辑
- [ ] 处理错误回滚
- [ ] 添加成功/失败提示
- [ ] 阻止按钮上的拖动事件
- [ ] 区分点击和拖动
- [ ] 测试边界情况（快速拖动、拖动到边界外）

---

## 🧪 测试用例

1. **基本拖动**
   - 拖动分组A到分组B上方 → A应该在B之前
   - 拖动分组A到分组B下方 → A应该在B之后

2. **事件冲突**
   - 拖动分组 → 不应该触发选中
   - 点击Power按钮 → 不应该触发拖动
   - 双击分组 → 不应该触发拖动

3. **边界情况**
   - 拖动到列表顶部/底部
   - 快速连续拖动
   - 拖动过程中切换到其他应用
   - API 失败时的回滚

4. **性能测试**
   - 50+ 分组的拖动性能
   - 拖动过程中的帧率
   - API 调用的响应时间

---

## 📚 参考资料

- [HTML5 Drag and Drop API](https://developer.mozilla.org/en-US/docs/Web/API/HTML_Drag_and_Drop_API)
- [@dnd-kit 文档](https://docs.dndkit.com/)
- [React DnD](https://react-dnd.github.io/react-dnd/)

---

**最后更新**: 2025-01-18
