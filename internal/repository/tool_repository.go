package repository

import (
	"Lin_studio/internal/config"
	"Lin_studio/internal/domain"
	"context"
	"errors"

	"gorm.io/gorm"
)

// ToolFilter 工具筛选条件
type ToolFilter struct {
	Page     int
	Limit    int
	Category *string
	Status   *string
	Search   *string
	URL      *string
}

// ToolRepository 工具仓储接口
type ToolRepository interface {
	Create(ctx context.Context, tool *domain.Tool) error
	FindByID(ctx context.Context, id uint) (*domain.Tool, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Tool, error)
	FindAll(ctx context.Context, filter ToolFilter) ([]domain.Tool, int64, error)
	Update(ctx context.Context, tool *domain.Tool) error
	Delete(ctx context.Context, id uint) error
	IncrementViews(ctx context.Context, id uint) error
}

// ToolRepositoryImpl 工具仓储实现
type ToolRepositoryImpl struct {
	db *gorm.DB
}

// NewToolRepository 创建工具仓储实例
func NewToolRepository() ToolRepository {
	return &ToolRepositoryImpl{
		db: config.DB,
	}
}

// Create 创建工具
func (r *ToolRepositoryImpl) Create(ctx context.Context, tool *domain.Tool) error {
	return r.db.WithContext(ctx).Create(tool).Error
}

// FindByID 根据ID查找工具
func (r *ToolRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.Tool, error) {
	var tool domain.Tool
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tool).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tool, nil
}

// FindBySlug 根据Slug查找工具
func (r *ToolRepositoryImpl) FindBySlug(ctx context.Context, slug string) (*domain.Tool, error) {
	var tool domain.Tool
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&tool).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tool, nil
}

// FindAll 查找所有工具
func (r *ToolRepositoryImpl) FindAll(ctx context.Context, filter ToolFilter) ([]domain.Tool, int64, error) {
	var tools []domain.Tool
	var total int64
	
	query := r.db.WithContext(ctx).Model(&domain.Tool{})
	
	// 应用过滤条件
	if filter.Category != nil && *filter.Category != "" {
		query = query.Where("category = ?", *filter.Category)
	}
	
	if filter.Status != nil && *filter.Status != "" {
		query = query.Where("status = ?", *filter.Status)
	} else {
		// 默认只查询激活状态的工具
		query = query.Where("status = ?", "active")
// FindAll 查找所有工具
	}
	
	// 计算总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// 排序和分页
	query = query.Order("name ASC")
	
	// 计算分页
	offset := (filter.Page - 1) * filter.Limit
	query = query.Offset(offset).Limit(filter.Limit)
	
	// 执行查询
	err = query.Find(&tools).Error
	
	return tools, total, err
}

// Update 更新工具
func (r *ToolRepositoryImpl) Update(ctx context.Context, tool *domain.Tool) error {
	return r.db.WithContext(ctx).Save(tool).Error
}

// Delete 删除工具
func (r *ToolRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Tool{}, id).Error
}

// IncrementViews 增加工具浏览量
func (r *ToolRepositoryImpl) IncrementViews(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&domain.Tool{}).
		Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + ?", 1)).
		Error
} 