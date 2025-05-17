package repository

import (
	"Lin_studio/internal/domain"
	"Lin_studio/internal/config"
	"context"
	"errors"

	"gorm.io/gorm"
)

// CategoryRepository 分类仓储接口
type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	FindByID(ctx context.Context, id uint) (*domain.Category, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Category, error)
	FindAll(ctx context.Context, parentID *uint) ([]domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id uint) error
	CountArticles(ctx context.Context, categoryID uint) (int64, error)
}

// CategoryRepositoryImpl 分类仓储实现
type CategoryRepositoryImpl struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓储实例
func NewCategoryRepository() CategoryRepository {
	return &CategoryRepositoryImpl{
		db: config.DB,
	}
}

// Create 创建分类
func (r *CategoryRepositoryImpl) Create(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

// FindByID 根据ID查找分类
func (r *CategoryRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.Category, error) {
	var category domain.Category
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// FindBySlug 根据Slug查找分类
func (r *CategoryRepositoryImpl) FindBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

// FindAll 查找所有分类
func (r *CategoryRepositoryImpl) FindAll(ctx context.Context, parentID *uint) ([]domain.Category, error) {
	var categories []domain.Category
	query := r.db.WithContext(ctx)
	
	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	
	err := query.Find(&categories).Error
	return categories, err
}

// Update 更新分类
func (r *CategoryRepositoryImpl) Update(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

// Delete 删除分类
func (r *CategoryRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Category{}, id).Error
}

// CountArticles 统计分类下的文章数量
func (r *CategoryRepositoryImpl) CountArticles(ctx context.Context, categoryID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Article{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
} 