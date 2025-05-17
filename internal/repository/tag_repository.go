package repository

import (
	"Lin_studio/internal/domain"
	"Lin_studio/internal/config"
	"context"
	"errors"

	"gorm.io/gorm"
)

// TagRepository 标签仓储接口
type TagRepository interface {
	Create(ctx context.Context, tag *domain.Tag) error
	FindByID(ctx context.Context, id uint) (*domain.Tag, error)
	FindBySlug(ctx context.Context, slug string) (*domain.Tag, error)
	FindAll(ctx context.Context) ([]domain.Tag, error)
	Update(ctx context.Context, tag *domain.Tag) error
	Delete(ctx context.Context, id uint) error
	CountArticles(ctx context.Context, tagID uint) (int64, error)
	CountProjects(ctx context.Context, tagID uint) (int64, error)
	FindByArticleID(ctx context.Context, articleID uint) ([]domain.Tag, error)
	AddArticleTag(ctx context.Context, articleID, tagID uint) error
	RemoveArticleTag(ctx context.Context, articleID, tagID uint) error
	ClearArticleTags(ctx context.Context, articleID uint) error
}

// TagRepositoryImpl 标签仓储实现
type TagRepositoryImpl struct {
	db *gorm.DB
}

// NewTagRepository 创建标签仓储实例
func NewTagRepository() TagRepository {
	return &TagRepositoryImpl{
		db: config.DB,
	}
}

// Create 创建标签
func (r *TagRepositoryImpl) Create(ctx context.Context, tag *domain.Tag) error {
	return r.db.WithContext(ctx).Create(tag).Error
}

// FindByID 根据ID查找标签
func (r *TagRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.Tag, error) {
	var tag domain.Tag
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tag, nil
}

// FindBySlug 根据Slug查找标签
func (r *TagRepositoryImpl) FindBySlug(ctx context.Context, slug string) (*domain.Tag, error) {
	var tag domain.Tag
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tag, nil
}

// FindAll 查找所有标签
func (r *TagRepositoryImpl) FindAll(ctx context.Context) ([]domain.Tag, error) {
	var tags []domain.Tag
	err := r.db.WithContext(ctx).Find(&tags).Error
	return tags, err
}

// Update 更新标签
func (r *TagRepositoryImpl) Update(ctx context.Context, tag *domain.Tag) error {
	return r.db.WithContext(ctx).Save(tag).Error
}

// Delete 删除标签
func (r *TagRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Tag{}, id).Error
}

// CountArticles 统计标签下的文章数量
func (r *TagRepositoryImpl) CountArticles(ctx context.Context, tagID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.ArticleTag{}).Where("tag_id = ?", tagID).Count(&count).Error
	return count, err
}

// CountProjects 统计标签下的项目数量
func (r *TagRepositoryImpl) CountProjects(ctx context.Context, tagID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("project_tags").Where("tag_id = ?", tagID).Count(&count).Error
	return count, err
}

// FindByArticleID 查找文章关联的所有标签
func (r *TagRepositoryImpl) FindByArticleID(ctx context.Context, articleID uint) ([]domain.Tag, error) {
	var tags []domain.Tag
	err := r.db.WithContext(ctx).
		Joins("JOIN article_tags ON article_tags.tag_id = tags.id").
		Where("article_tags.article_id = ?", articleID).
		Find(&tags).Error
	return tags, err
}

// AddArticleTag 添加文章标签关联
func (r *TagRepositoryImpl) AddArticleTag(ctx context.Context, articleID, tagID uint) error {
	articleTag := domain.ArticleTag{
		ArticleID: articleID,
		TagID:     tagID,
	}
	return r.db.WithContext(ctx).Create(&articleTag).Error
}

// RemoveArticleTag 移除文章标签关联
func (r *TagRepositoryImpl) RemoveArticleTag(ctx context.Context, articleID, tagID uint) error {
	return r.db.WithContext(ctx).
		Where("article_id = ? AND tag_id = ?", articleID, tagID).
		Delete(&domain.ArticleTag{}).Error
}

// ClearArticleTags 清除文章的所有标签关联
func (r *TagRepositoryImpl) ClearArticleTags(ctx context.Context, articleID uint) error {
	return r.db.WithContext(ctx).
		Where("article_id = ?", articleID).
		Delete(&domain.ArticleTag{}).Error
} 