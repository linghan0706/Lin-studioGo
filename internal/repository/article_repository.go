package repository

import (
	"Lin_studio/internal/domain"
	"Lin_studio/internal/config"
	"context"
	"errors"

	"gorm.io/gorm"
)

// ArticleFilter 文章筛选条件
type ArticleFilter struct {
	Page       int
	Limit      int
	CategoryID *uint
	TagID      *uint
	AuthorID   *uint
	Status     *string
	Search     *string
	Sort       *string
}

// ArticleRepository 文章仓储接口
type ArticleRepository interface {
	Create(ctx context.Context, article *domain.Article) error
	FindByID(ctx context.Context, id uint) (*domain.Article, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Article, error)
	FindAll(ctx context.Context, filter ArticleFilter) ([]domain.Article, int64, error)
	FindFeatured(ctx context.Context, limit int) ([]domain.Article, error)
	Update(ctx context.Context, article *domain.Article) error
	Delete(ctx context.Context, id uint) error
	IncrementViews(ctx context.Context, id uint) error
	IncrementLikes(ctx context.Context, id uint) error
}

// ArticleRepositoryImpl 文章仓储实现
type ArticleRepositoryImpl struct {
	db *gorm.DB
}

// NewArticleRepository 创建文章仓储实例
func NewArticleRepository() ArticleRepository {
	return &ArticleRepositoryImpl{
		db: config.DB,
	}
}

// Create 创建文章
func (r *ArticleRepositoryImpl) Create(ctx context.Context, article *domain.Article) error {
	return r.db.WithContext(ctx).Create(article).Error
}

// FindByID 根据ID查找文章
func (r *ArticleRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.Article, error) {
	var article domain.Article
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &article, nil
}

// FindBySlug 根据Slug查找文章
func (r *ArticleRepositoryImpl) FindBySlug(ctx context.Context, slug string) (*domain.Article, error) {
	var article domain.Article
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&article).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &article, nil
}

// FindAll 查找所有文章
func (r *ArticleRepositoryImpl) FindAll(ctx context.Context, filter ArticleFilter) ([]domain.Article, int64, error) {
	var articles []domain.Article
	var total int64
	
	query := r.db.WithContext(ctx).Model(&domain.Article{})
	
	// 应用过滤条件
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	
	if filter.AuthorID != nil {
		query = query.Where("author_id = ?", *filter.AuthorID)
	}
	
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	} else {
		// 默认只查询已发布文章
		query = query.Where("status = ?", "published")
	}
	
	if filter.Search != nil && *filter.Search != "" {
		searchTerm := "%" + *filter.Search + "%"
		query = query.Where("title LIKE ? OR excerpt LIKE ?", searchTerm, searchTerm)
	}
	
	// 如果有标签过滤，需要Join
	if filter.TagID != nil {
		query = query.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Where("article_tags.tag_id = ?", *filter.TagID)
	}
	
	// 计算总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// 应用排序
	if filter.Sort != nil && *filter.Sort != "" {
		sortField := *filter.Sort
		if sortField[0] == '-' {
			query = query.Order(sortField[1:] + " DESC")
		} else {
			query = query.Order(sortField)
		}
	} else {
		// 默认按发布时间倒序排序
		query = query.Order("published_at DESC")
	}
	
	// 计算分页
	offset := (filter.Page - 1) * filter.Limit
	query = query.Offset(offset).Limit(filter.Limit)
	
	// 执行查询
	err = query.Find(&articles).Error
	
	return articles, total, err
}

// FindFeatured 查找精选文章
func (r *ArticleRepositoryImpl) FindFeatured(ctx context.Context, limit int) ([]domain.Article, error) {
	var articles []domain.Article
	err := r.db.WithContext(ctx).
		Where("featured_order IS NOT NULL AND status = ?", "published").
		Order("featured_order ASC").
		Limit(limit).
		Find(&articles).Error
	return articles, err
}

// Update 更新文章
func (r *ArticleRepositoryImpl) Update(ctx context.Context, article *domain.Article) error {
	return r.db.WithContext(ctx).Save(article).Error
}

// Delete 删除文章
func (r *ArticleRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Article{}, id).Error
}

// IncrementViews 增加文章浏览量
func (r *ArticleRepositoryImpl) IncrementViews(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&domain.Article{}).
		Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + ?", 1)).
		Error
}

// IncrementLikes 增加文章点赞数
func (r *ArticleRepositoryImpl) IncrementLikes(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&domain.Article{}).
		Where("id = ?", id).
		UpdateColumn("likes", gorm.Expr("likes + ?", 1)).
		Error
} 