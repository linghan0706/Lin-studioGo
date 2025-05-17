package repository

import (
	"Lin_studio/internal/config"
	"Lin_studio/internal/domain"
	"context"
	"errors"

	"gorm.io/gorm"
)

// CommentFilter 评论筛选条件
type CommentFilter struct {
	Page     int
	Limit    int
	ItemType string
	ItemID   uint
	ParentID *uint
	UserID   *uint
	Status   *string
}

// CommentRepository 评论仓储接口
type CommentRepository interface {
	Create(ctx context.Context, comment *domain.Comment) error
	FindByID(ctx context.Context, id uint) (*domain.Comment, error)
	FindAll(ctx context.Context, filter CommentFilter) ([]domain.Comment, int64, error)
	FindReplies(ctx context.Context, parentID uint) ([]domain.Comment, error)
	CountReplies(ctx context.Context, parentID uint) (int64, error)
	Update(ctx context.Context, comment *domain.Comment) error
	UpdateStatus(ctx context.Context, id uint, status string) error
	Delete(ctx context.Context, id uint) error
	IncrementLikes(ctx context.Context, id uint) error
}

// CommentRepositoryImpl 评论仓储实现
type CommentRepositoryImpl struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论仓储实例
func NewCommentRepository() CommentRepository {
	return &CommentRepositoryImpl{
		db: config.DB,
	}
}

// Create 创建评论
func (r *CommentRepositoryImpl) Create(ctx context.Context, comment *domain.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

// FindByID 根据ID查找评论
func (r *CommentRepositoryImpl) FindByID(ctx context.Context, id uint) (*domain.Comment, error) {
	var comment domain.Comment
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

// FindAll 查找所有评论
func (r *CommentRepositoryImpl) FindAll(ctx context.Context, filter CommentFilter) ([]domain.Comment, int64, error) {
	var comments []domain.Comment
	var total int64
	
	query := r.db.WithContext(ctx).Model(&domain.Comment{})
	
	// 应用过滤条件
	if filter.ItemType != "" {
		query = query.Where("item_type = ?", filter.ItemType)
	}
	
	if filter.ItemID > 0 {
		query = query.Where("item_id = ?", filter.ItemID)
	}
	
	if filter.ParentID != nil {
		query = query.Where("parent_id = ?", *filter.ParentID)
	} else {
		// 默认查询顶级评论（无父评论）
		query = query.Where("parent_id IS NULL")
	}
	
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	} else {
		// 默认只查询已批准的评论
		query = query.Where("status = ?", "approved")
	}
	
	// 计算总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// 排序和分页
	query = query.Order("created_at DESC")
	
	// 计算分页
	offset := (filter.Page - 1) * filter.Limit
	query = query.Offset(offset).Limit(filter.Limit)
	
	// 执行查询
	err = query.Find(&comments).Error
	
	return comments, total, err
}

// FindReplies 查找指定评论的回复
func (r *CommentRepositoryImpl) FindReplies(ctx context.Context, parentID uint) ([]domain.Comment, error) {
	var replies []domain.Comment
	err := r.db.WithContext(ctx).
		Where("parent_id = ? AND status = ?", parentID, "approved").
		Order("created_at ASC").
		Find(&replies).Error
	return replies, err
}

// CountReplies 统计指定评论的回复数量
func (r *CommentRepositoryImpl) CountReplies(ctx context.Context, parentID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Comment{}).
		Where("parent_id = ? AND status = ?", parentID, "approved").
		Count(&count).Error
	return count, err
}

// Update 更新评论
func (r *CommentRepositoryImpl) Update(ctx context.Context, comment *domain.Comment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

// UpdateStatus 更新评论状态
func (r *CommentRepositoryImpl) UpdateStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).
		Model(&domain.Comment{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete 删除评论
func (r *CommentRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Comment{}, id).Error
}

// IncrementLikes 增加评论点赞数
func (r *CommentRepositoryImpl) IncrementLikes(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&domain.Comment{}).
		Where("id = ?", id).
		UpdateColumn("likes", gorm.Expr("likes + ?", 1)).
		Error
} 