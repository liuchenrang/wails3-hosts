# 项目文档目录

本目录包含项目的所有文档和设计说明。

## 📚 文档列表

### Hosts Manager 项目文档

#### 1. [hosts-manager-architecture.md](./hosts-manager-architecture.md)
**DDD 架构设计文档**

- 项目概述和技术栈
- DDD 分层架构设计
- 实体、值对象、领域服务设计
- 应用服务、DTO 设计
- 基础设施层设计

**适用场景**：了解项目整体架构、代码结构、设计模式

#### 2. [PROJECT_SUMMARY.md](./PROJECT_SUMMARY.md)
**项目实现总结**

- 项目概述和实现进度
- 技术亮点和最佳实践
- SOLID 原则应用实例
- 代码质量保证措施
- 下一步开发计划

**适用场景**：快速了解项目现状、技术特色、开发进度

### 其他项目文档

#### 3. [super_pt.md](./super_pt.md)
**超级打印系统设计**

- 打印机监控和订单系统
- 调度系统和 Worker 设计
- 监控报警和面单模块
- 耗材管理和保养模块

**说明**：这是另一个项目的设计文档，与 Hosts Manager 无关

## 📖 使用建议

### 新成员入职
1. 先阅读 `PROJECT_SUMMARY.md` 了解项目全貌
2. 再阅读 `hosts-manager-architecture.md` 理解架构设计
3. 结合代码实际操作，加深理解

### 代码审查
- 参考 `hosts-manager-architecture.md` 中的 DDD 设计原则
- 检查代码是否符合 SOLID、DRY、KISS、YAGNI 原则

### 功能开发
- 查看 `PROJECT_SUMMARY.md` 的"待办事项"和"下一步计划"
- 参考 `hosts-manager-architecture.md` 确保新功能符合架构设计

## 🔄 文档维护

文档更新原则：
- **及时性**：架构变更后立即更新文档
- **准确性**：文档与代码保持一致
- **简洁性**：避免冗余，突出重点

---

**最后更新**: 2025-01-18
