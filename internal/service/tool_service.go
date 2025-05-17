package service

import (
	"Lin_studio/internal/domain"
	"Lin_studio/internal/repository"
	"Lin_studio/internal/utils"
	"context"
	"errors"
	"strings"
)

// ToolService 工具服务接口
type ToolService interface {
	GetTools(ctx context.Context, page, limit int, category, status, search *string) ([]domain.Tool, domain.PaginationData, error)
	GetToolByID(ctx context.Context, id uint) (*domain.Tool, error)
	GetToolBySlug(ctx context.Context, slug string) (*domain.Tool, error)
	CreateTool(ctx context.Context, name, description, icon, category, content string, config domain.JSONConfig, status string) (*domain.Tool, error)
	UpdateTool(ctx context.Context, id uint, name, description, icon, category, content string, config domain.JSONConfig, status string) (*domain.Tool, error)
	DeleteTool(ctx context.Context, id uint) error
	ViewTool(ctx context.Context, id uint) error
}

// ToolServiceImpl 工具服务实现
type ToolServiceImpl struct {
	toolRepo repository.ToolRepository
}

// NewToolService 创建工具服务实例
func NewToolService(toolRepo repository.ToolRepository) ToolService {
	return &ToolServiceImpl{
		toolRepo: toolRepo,
	}
}

// GetTools 获取工具列表
func (s *ToolServiceImpl) GetTools(
	ctx context.Context,
	page, limit int,
	category, status, search *string,
) ([]domain.Tool, domain.PaginationData, error) {
	// 创建过滤条件
	filter := repository.ToolFilter{
		Page:     page,
		Limit:    limit,
		Category: category,
		Status:   status,
		Search:   search,
	}

	// 查询工具
	tools, total, err := s.toolRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, domain.PaginationData{}, err
	}

	// 创建分页数据
	pagination := domain.PaginationData{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: (int(total) + limit - 1) / limit,
	}

	return tools, pagination, nil
}

// GetToolByID 根据ID获取工具
func (s *ToolServiceImpl) GetToolByID(ctx context.Context, id uint) (*domain.Tool, error) {
	tool, err := s.toolRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if tool == nil {
		return nil, errors.New("工具不存在")
	}

	return tool, nil
}

// GetToolBySlug 根据Slug获取工具
func (s *ToolServiceImpl) GetToolBySlug(ctx context.Context, slug string) (*domain.Tool, error) {
	tool, err := s.toolRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if tool == nil {
		return nil, errors.New("工具不存在")
	}

	// 增加浏览量
	go s.toolRepo.IncrementViews(context.Background(), tool.ID)

	return tool, nil
}

// CreateTool 创建工具
func (s *ToolServiceImpl) CreateTool(
	ctx context.Context,
	name, description, icon, category, content string,
	config domain.JSONConfig,
	status string,
) (*domain.Tool, error) {
	// 检查名称是否为空
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("工具名称不能为空")
	}

	// 生成slug
	slug := utils.GenerateSlug(name)

	// 检查slug是否已存在
	existingTool, err := s.toolRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if existingTool != nil {
		return nil, errors.New("已存在相同名称的工具")
	}

	// 检查状态是否有效
	if status != "active" && status != "maintenance" && status != "deprecated" {
		status = "active" // 默认为active
	}

	// 创建工具
	tool := &domain.Tool{
		Name:        name,
		Slug:        slug,
		Description: description,
		Icon:        icon,
		Category:    category,
		Content:     content,
		Config:      config,
		Status:      status,
	}

	// 保存工具
	err = s.toolRepo.Create(ctx, tool)
	if err != nil {
		return nil, err
	}

	return tool, nil
}

// UpdateTool 更新工具
func (s *ToolServiceImpl) UpdateTool(
	ctx context.Context,
	id uint,
	name, description, icon, category, content string,
	config domain.JSONConfig,
	status string,
) (*domain.Tool, error) {
	// 检查工具是否存在
	tool, err := s.toolRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tool == nil {
		return nil, errors.New("工具不存在")
	}

	// 检查名称是否为空
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("工具名称不能为空")
	}

	// 如果名称变了，则更新slug并检查是否已存在
	if name != tool.Name {
		newSlug := utils.GenerateSlug(name)
		existingTool, err := s.toolRepo.FindBySlug(ctx, newSlug)
		if err != nil {
			return nil, err
		}
		if existingTool != nil && existingTool.ID != id {
			return nil, errors.New("已存在相同名称的工具")
		}
		tool.Slug = newSlug
	}

	// 检查状态是否有效
	if status != "active" && status != "maintenance" && status != "deprecated" {
		status = "active" // 默认为active
	}

	// 更新工具
	tool.Name = name
	tool.Description = description
	tool.Icon = icon
	tool.Category = category
	tool.Content = content
	tool.Config = config
	tool.Status = status

	// 保存更新
	err = s.toolRepo.Update(ctx, tool)
	if err != nil {
		return nil, err
	}

	return tool, nil
}

// DeleteTool 删除工具
func (s *ToolServiceImpl) DeleteTool(ctx context.Context, id uint) error {
	// 检查工具是否存在
	tool, err := s.toolRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if tool == nil {
		return errors.New("工具不存在")
	}

	// 删除工具
	return s.toolRepo.Delete(ctx, id)
}

// ViewTool 增加工具浏览量
func (s *ToolServiceImpl) ViewTool(ctx context.Context, id uint) error {
	return s.toolRepo.IncrementViews(ctx, id)
} 