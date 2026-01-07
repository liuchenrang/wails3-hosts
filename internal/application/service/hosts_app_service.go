package service

import (
	"context"
	"fmt"

	"github.com/chen/wails3-hosts/internal/application/dto"
	"github.com/chen/wails3-hosts/internal/domain/entity"
	"github.com/chen/wails3-hosts/internal/domain/repository"
	"github.com/chen/wails3-hosts/internal/domain/service"
	"github.com/chen/wails3-hosts/internal/infrastructure/system"
)

// HostsApplicationService hosts 应用服务
// 单一职责: 协调领域对象和基础设施服务完成业务用例
// DDD: 应用服务是领域模型的门面，处理事务边界
type HostsApplicationService struct {
	hostsRepo     repository.HostsRepository
	versionRepo   repository.VersionRepository
	domainService *service.HostsDomainService
	hostsFileOp   *system.HostsFileOperator
	sudoManager   *system.SudoManager
}

// NewHostsApplicationService 创建应用服务实例
// 依赖注入: 通过构造函数注入所有依赖
func NewHostsApplicationService(
	hostsRepo repository.HostsRepository,
	versionRepo repository.VersionRepository,
	hostsFileOp *system.HostsFileOperator,
	sudoManager *system.SudoManager,
) *HostsApplicationService {
	return &HostsApplicationService{
		hostsRepo:     hostsRepo,
		versionRepo:   versionRepo,
		domainService: service.NewHostsDomainService(),
		hostsFileOp:   hostsFileOp,
		sudoManager:   sudoManager,
	}
}

// CreateGroup 创建一个新的 hosts 分组
// 用例: 用户创建新的配置分组
func (s *HostsApplicationService) CreateGroup(ctx context.Context, req dto.CreateHostsGroupRequest) (*dto.HostsGroupDTO, error) {
	// 验证输入
	if req.Name == "" {
		return nil, entity.ErrInvalidName
	}

	// 检查名称是否已存在
	exists, err := s.hostsRepo.ExistsByName(ctx, req.Name)
	if err != nil {
		return nil, fmt.Errorf("检查分组名称失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("分组名称已存在")
	}

	// 创建领域实体
	group := entity.NewHostsGroup(req.Name, req.Description)

	// 持久化
	if err := s.hostsRepo.Save(ctx, group); err != nil {
		return nil, fmt.Errorf("保存分组失败: %w", err)
	}

	return s.toGroupDTO(group), nil
}

// GetAllGroups 获取所有分组
// 用例: 加载左侧分组列表
func (s *HostsApplicationService) GetAllGroups(ctx context.Context) ([]dto.HostsGroupDTO, error) {
	groups, err := s.hostsRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("查找分组失败: %w", err)
	}

	// 按照order字段排序
	// DRY: 使用Go标准库的sort.Slice进行排序
	sortedGroups := make([]*entity.HostsGroup, len(groups))
	copy(sortedGroups, groups)

	// 简单排序: 按Order字段升序排列
	for i := 0; i < len(sortedGroups)-1; i++ {
		for j := i + 1; j < len(sortedGroups); j++ {
			if sortedGroups[i].Order > sortedGroups[j].Order {
				sortedGroups[i], sortedGroups[j] = sortedGroups[j], sortedGroups[i]
			}
		}
	}

	result := make([]dto.HostsGroupDTO, 0, len(sortedGroups))
	for _, group := range sortedGroups {
		result = append(result, *s.toGroupDTO(group))
	}
	return result, nil
}

// GetGroupByID 根据 ID 获取分组
func (s *HostsApplicationService) GetGroupByID(ctx context.Context, id string) (*dto.HostsGroupDTO, error) {
	group, err := s.hostsRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("查找分组失败: %w", err)
	}
	return s.toGroupDTO(group), nil
}

// UpdateGroup 更新分组信息
func (s *HostsApplicationService) UpdateGroup(ctx context.Context, req dto.UpdateHostsGroupRequest) error {
	group, err := s.hostsRepo.FindByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("查找分组失败: %w", err)
	}

	// 更新字段
	group.Name = req.Name
	group.Description = req.Description

	// 持久化
	return s.hostsRepo.Save(ctx, group)
}

// DeleteGroup 删除分组
func (s *HostsApplicationService) DeleteGroup(ctx context.Context, id string) error {
	return s.hostsRepo.Delete(ctx, id)
}

// ToggleGroup 切换分组的启用状态
// 用例: 用户点击分组前的复选框
func (s *HostsApplicationService) ToggleGroup(ctx context.Context, req dto.ToggleGroupRequest) error {
	group, err := s.hostsRepo.FindByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("查找分组失败: %w", err)
	}

	group.SetEnabled(req.Enabled)
	return s.hostsRepo.Save(ctx, group)
}

// ReorderGroups 重新排序分组
// 用例: 用户拖动分组改变顺序
func (s *HostsApplicationService) ReorderGroups(ctx context.Context, req dto.ReorderGroupsRequest) error {
	// 获取所有分组
	allGroups, err := s.hostsRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("查找分组失败: %w", err)
	}

	// 创建分组ID到分组的映射
	groupMap := make(map[string]*entity.HostsGroup)
	for i := range allGroups {
		groupMap[allGroups[i].ID] = allGroups[i]
	}

	// 按照新顺序更新分组的排序字段
	for i, groupID := range req.GroupIDs {
		if group, exists := groupMap[groupID]; exists {
			group.SetOrder(i)
			if err := s.hostsRepo.Save(ctx, group); err != nil {
				return fmt.Errorf("保存分组排序失败: %w", err)
			}
		}
	}

	return nil
}

// AddEntry 向分组添加一个 hosts 条目
// 用例: 用户在右侧面板添加新条目
func (s *HostsApplicationService) AddEntry(ctx context.Context, req dto.AddEntryRequest) error {
	group, err := s.hostsRepo.FindByID(ctx, req.GroupID)
	if err != nil {
		return fmt.Errorf("查找分组失败: %w", err)
	}

	entry := entity.NewHostsEntry(req.IP, req.Hostname, req.Comment)
	if err := group.AddEntry(*entry); err != nil {
		return err
	}

	return s.hostsRepo.Save(ctx, group)
}

// UpdateEntry 更新分组中的条目
func (s *HostsApplicationService) UpdateEntry(ctx context.Context, req dto.UpdateEntryRequest) error {
	group, err := s.hostsRepo.FindByID(ctx, req.GroupID)
	if err != nil {
		return fmt.Errorf("查找分组失败: %w", err)
	}

	newEntry := entity.NewHostsEntry(req.IP, req.Hostname, req.Comment)
	newEntry.ID = req.EntryID // 保持原有 ID

	if err := group.UpdateEntry(req.EntryID, *newEntry); err != nil {
		return err
	}

	return s.hostsRepo.Save(ctx, group)
}

// DeleteEntry 从分组中删除条目
func (s *HostsApplicationService) DeleteEntry(ctx context.Context, req dto.DeleteEntryRequest) error {
	group, err := s.hostsRepo.FindByID(ctx, req.GroupID)
	if err != nil {
		return fmt.Errorf("查找分组失败: %w", err)
	}

	group.RemoveEntry(req.EntryID)
	return s.hostsRepo.Save(ctx, group)
}

// BatchUpdateEntries 批量更新分组中的所有条目
// 用例: memo 编辑模式下一次性更新所有条目
func (s *HostsApplicationService) BatchUpdateEntries(ctx context.Context, req dto.BatchUpdateEntriesRequest) error {
	group, err := s.hostsRepo.FindByID(ctx, req.GroupID)
	if err != nil {
		return fmt.Errorf("查找分组失败: %w", err)
	}

	// 清空现有条目
	group.ClearEntries()

	// 添加新条目
	for _, entryReq := range req.Entries {
		entry := entity.NewHostsEntry(entryReq.IP, entryReq.Hostname, entryReq.Comment)
		entry.Enabled = entryReq.Enabled
		if err := group.AddEntry(*entry); err != nil {
			return fmt.Errorf("添加条目失败: %w", err)
		}
	}

	return s.hostsRepo.Save(ctx, group)
}

// GeneratePreview 生成 hosts 配置预览
// 用例: 用户查看应用后的 hosts 文件内容
func (s *HostsApplicationService) GeneratePreview(ctx context.Context) (string, error) {
	groups, err := s.hostsRepo.FindAll(ctx)
	if err != nil {
		return "", fmt.Errorf("查找分组失败: %w", err)
	}

	return s.domainService.GenerateHostsContent(groups), nil
}

// DetectConflicts 检测配置冲突
// 用例: 保存前检查是否有重复的主机名映射
func (s *HostsApplicationService) DetectConflicts(ctx context.Context) (map[string][]string, error) {
	groups, err := s.hostsRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("查找分组失败: %w", err)
	}

	conflicts := s.domainService.DetectConflicts(groups)
	return conflicts, nil
}

// ApplyHosts 应用 hosts 配置到系统
// 用例: 用户点击"应用"按钮或按 Cmd+S
// 前置条件: 密码必须先通过 ValidateSudoPassword 验证并缓存
func (s *HostsApplicationService) ApplyHosts(ctx context.Context, req dto.ApplyHostsRequest) error {
	// 1. 检查是否有缓存的 sudo 密码
	if !s.sudoManager.IsPasswordCached() {
		return fmt.Errorf("需要先验证 sudo 密码，请调用 ValidateSudoPassword")
	}

	// 2. 生成 hosts 内容
	content, err := s.GeneratePreview(ctx)
	if err != nil {
		return err
	}

	// 3. 备份当前 hosts 文件
	if err := s.hostsFileOp.Backup(); err != nil {
		return fmt.Errorf("备份 hosts 文件失败: %w", err)
	}

	// 4. 保存当前版本到历史
	currentContent, _ := s.hostsFileOp.ReadCurrent()
	version := entity.NewHostsVersion(
		currentContent,
		"应用 hosts 配置前备份",
		entity.SourceManual,
	)
	if err := s.versionRepo.Save(ctx, version); err != nil {
		return fmt.Errorf("保存版本失败: %w", err)
	}

	// 5. 写入 hosts 文件（使用缓存的系统凭证）
	if err := s.hostsFileOp.Write(content); err != nil {
		return fmt.Errorf("写入 hosts 文件失败: %w", err)
	}

	return nil
}

// GetVersions 获取版本历史
// 用例: 用户查看版本历史界面
func (s *HostsApplicationService) GetVersions(ctx context.Context, limit int) ([]dto.HostsVersionDTO, error) {
	versions, err := s.versionRepo.FindLatest(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("查找版本失败: %w", err)
	}

	result := make([]dto.HostsVersionDTO, 0, len(versions))
	for _, v := range versions {
		result = append(result, dto.HostsVersionDTO{
			ID:          v.ID,
			Timestamp:   v.Timestamp.Format("2006-01-02 15:04:05"),
			Content:     v.Content,
			Description: v.Description,
			Source:      v.Source,
			Age:         v.GetAge(),
		})
	}
	return result, nil
}

// RollbackToVersion 回滚到指定版本
// 用例: 用户在版本历史中选择一个版本并回滚
func (s *HostsApplicationService) RollbackToVersion(ctx context.Context, req dto.RollbackRequest) error {
	// 1. 查找目标版本
	targetVersion, err := s.versionRepo.FindByID(ctx, req.VersionID)
	if err != nil {
		return fmt.Errorf("查找版本失败: %w", err)
	}

	// 2. 备份当前 hosts 文件
	if err := s.hostsFileOp.Backup(); err != nil {
		return fmt.Errorf("备份 hosts 文件失败: %w", err)
	}

	// 3. 写入目标版本内容（使用缓存的系统凭证）
	if err := s.hostsFileOp.Write(targetVersion.Content); err != nil {
		return fmt.Errorf("写入 hosts 文件失败: %w", err)
	}

	// 4. 记录回滚操作
	currentContent, _ := s.hostsFileOp.ReadCurrent()
	rollbackVersion := entity.NewHostsVersion(
		currentContent,
		fmt.Sprintf("回滚到版本 %s", targetVersion.ID[:8]),
		entity.SourceRollback,
	)
	if err := s.versionRepo.Save(ctx, rollbackVersion); err != nil {
		return fmt.Errorf("保存回滚版本失败: %w", err)
	}

	return nil
}

// ValidateSudoPassword 验证 sudo 密码
// 用例: 用户首次输入 sudo 密码
func (s *HostsApplicationService) ValidateSudoPassword(ctx context.Context, req dto.ValidateSudoRequest) dto.ValidateSudoResponse {
	valid := s.sudoManager.ValidatePassword(req.Password)
	if valid {
		// 缓存有效密码
		s.sudoManager.CachePassword(req.Password)
		return dto.ValidateSudoResponse{Valid: true}
	}
	return dto.ValidateSudoResponse{
		Valid: false,
		Error: "sudo 密码无效",
	}
}

// IsSudoPasswordCached 检查 sudo 密码是否已缓存
// 用例: 前端判断是否需要显示密码输入框
func (s *HostsApplicationService) IsSudoPasswordCached(ctx context.Context) bool {
	return s.sudoManager.IsPasswordCached()
}

// toGroupDTO 将领域实体转换为 DTO
// DRY: 统一的转换逻辑
func (s *HostsApplicationService) toGroupDTO(group *entity.HostsGroup) *dto.HostsGroupDTO {
	entries := make([]dto.HostsEntryDTO, 0, len(group.Entries))
	for _, entry := range group.Entries {
		entries = append(entries, dto.HostsEntryDTO{
			ID:       entry.ID,
			IP:       entry.IP,
			Hostname: entry.Hostname,
			Comment:  entry.Comment,
			Enabled:  entry.Enabled,
		})
	}

	return &dto.HostsGroupDTO{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		IsEnabled:   group.IsEnabled,
		Order:       group.Order,
		Entries:     entries,
		CreatedAt:   group.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   group.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
