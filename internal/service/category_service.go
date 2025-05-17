package service

import (
	"Lin_studio/internal/domain"
	"Lin_studio/internal/repository"
	"context"
	"errors"
	"strings"
)

// CategoryService 分类服务接口
type CategoryService interface {
	GetAllCategories(ctx context.Context, parentID *uint) ([]domain.CategoryResponse, error)
	CreateCategory(ctx context.Context, name, description string, parentID *uint) (*domain.Category, error)
	GetCategoryByID(ctx context.Context, id uint) (*domain.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error)
}

// CategoryServiceImpl 分类服务实现
type CategoryServiceImpl struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryService 创建分类服务实例
func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &CategoryServiceImpl{
		categoryRepo: categoryRepo,
	}
}

// GetAllCategories 获取所有分类
func (s *CategoryServiceImpl) GetAllCategories(ctx context.Context, parentID *uint) ([]domain.CategoryResponse, error) {
	categories, err := s.categoryRepo.FindAll(ctx, parentID)
	if err != nil {
		return nil, err
	}
	
	var response []domain.CategoryResponse
	for _, category := range categories {
		articleCount, err := s.categoryRepo.CountArticles(ctx, category.ID)
		if err != nil {
			return nil, err
		}
		
		response = append(response, domain.CategoryResponse{
			ID:           category.ID,
			Name:         category.Name,
			Slug:         category.Slug,
			Description:  category.Description,
			ParentID:     category.ParentID,
			ArticleCount: articleCount,
		})
	}
	
	return response, nil
}

// CreateCategory 创建分类
func (s *CategoryServiceImpl) CreateCategory(ctx context.Context, name, description string, parentID *uint) (*domain.Category, error) {
	// 检查分类名是否为空
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("分类名不能为空")
	}
	
	// 生成slug
	slug := s.generateSlug(name)
	
	// 检查slug是否已存在
	existingCategory, err := s.categoryRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	
	if existingCategory != nil {
		return nil, errors.New("分类已存在")
	}
	
	// 如果有父分类，检查是否存在
	if parentID != nil {
		parentCategory, err := s.categoryRepo.FindByID(ctx, *parentID)
		if err != nil {
			return nil, err
		}
		
		if parentCategory == nil {
			return nil, errors.New("父分类不存在")
		}
	}
	
	// 创建新分类
	category := &domain.Category{
		Name:        name,
		Slug:        slug,
		Description: description,
		ParentID:    parentID,
	}
	
	err = s.categoryRepo.Create(ctx, category)
	if err != nil {
		return nil, err
	}
	
	return category, nil
}

// GetCategoryByID 根据ID获取分类
func (s *CategoryServiceImpl) GetCategoryByID(ctx context.Context, id uint) (*domain.Category, error) {
	category, err := s.categoryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if category == nil {
		return nil, errors.New("分类不存在")
	}
	
	return category, nil
}

// GetCategoryBySlug 根据Slug获取分类
func (s *CategoryServiceImpl) GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	category, err := s.categoryRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	
	if category == nil {
		return nil, errors.New("分类不存在")
	}
	
	return category, nil
}

// generateSlug 生成分类别名
func (s *CategoryServiceImpl) generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	// 移除特殊字符
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)
	
	return slug
} 