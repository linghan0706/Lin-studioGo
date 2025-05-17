package service

import (
	"Lin_studio/internal/domain"
	"Lin_studio/internal/repository"
	"context"
	"errors"
	"strings"
)

// CommentService 评论服务接口
type CommentService interface {
	GetComments(ctx context.Context, itemType string, itemID uint, parentID *uint, page, limit int, status *string) ([]domain.Comment, domain.PaginationData, error)
	GetCommentByID(ctx context.Context, id uint) (*domain.Comment, error)
	CreateComment(ctx context.Context, content, itemType string, itemID uint, userID *uint, anonymousName, anonymousEmail string, parentID *uint) (*domain.Comment, error)
	UpdateComment(ctx context.Context, id uint, content string, userID uint) (*domain.Comment, error)
	DeleteComment(ctx context.Context, id uint, userID uint, isAdmin bool) error
	LikeComment(ctx context.Context, id uint) error
	ApproveComment(ctx context.Context, id uint) error
	MarkCommentAsSpam(ctx context.Context, id uint) error
}

// CommentServiceImpl 评论服务实现
type CommentServiceImpl struct {
	commentRepo repository.CommentRepository
	userRepo    repository.UserRepository
}

// NewCommentService 创建评论服务实例
func NewCommentService(
	commentRepo repository.CommentRepository,
	userRepo repository.UserRepository,
) CommentService {
	return &CommentServiceImpl{
		commentRepo: commentRepo,
		userRepo:    userRepo,
	}
}

// GetComments 获取评论列表
func (s *CommentServiceImpl) GetComments(
	ctx context.Context,
	itemType string,
	itemID uint,
	parentID *uint,
	page, limit int,
	status *string,
) ([]domain.Comment, domain.PaginationData, error) {
	// 创建过滤条件
	filter := repository.CommentFilter{
		Page:     page,
		Limit:    limit,
		ItemType: itemType,
		ItemID:   itemID,
		ParentID: parentID,
		Status:   status,
	}

	// 查询评论
	comments, total, err := s.commentRepo.FindAll(ctx, filter)
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

	// 加载用户信息和回复数量
	for i := range comments {
		if comments[i].UserID != nil {
			user, err := s.userRepo.FindByID(ctx, *comments[i].UserID)
			if err == nil && user != nil {
				comments[i].User = user
			}
		}

		// 获取回复计数
		replyCount, err := s.commentRepo.CountReplies(ctx, comments[i].ID)
		if err == nil {
			comments[i].ReplyCount = replyCount
		}

		// 如果需要，获取回复列表（通常限制在前几条）
		if replyCount > 0 && replyCount <= 3 { // 只有少量回复时才自动加载
			replies, err := s.commentRepo.FindReplies(ctx, comments[i].ID)
			if err == nil {
				// 加载回复的用户信息
				for j := range replies {
					if replies[j].UserID != nil {
						user, err := s.userRepo.FindByID(ctx, *replies[j].UserID)
						if err == nil && user != nil {
							replies[j].User = user
						}
					}
				}
				comments[i].Replies = replies
			}
		}
	}

	return comments, pagination, nil
}

// GetCommentByID 根据ID获取评论
func (s *CommentServiceImpl) GetCommentByID(ctx context.Context, id uint) (*domain.Comment, error) {
	comment, err := s.commentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if comment == nil {
		return nil, errors.New("评论不存在")
	}

	// 加载用户信息
	if comment.UserID != nil {
		user, err := s.userRepo.FindByID(ctx, *comment.UserID)
		if err == nil && user != nil {
			comment.User = user
		}
	}

	// 获取回复列表
	replies, err := s.commentRepo.FindReplies(ctx, comment.ID)
	if err == nil && len(replies) > 0 {
		// 加载回复的用户信息
		for i := range replies {
			if replies[i].UserID != nil {
				user, err := s.userRepo.FindByID(ctx, *replies[i].UserID)
				if err == nil && user != nil {
					replies[i].User = user
				}
			}
		}
		comment.Replies = replies
		comment.ReplyCount = int64(len(replies))
	}

	return comment, nil
}

// CreateComment 创建评论
func (s *CommentServiceImpl) CreateComment(
	ctx context.Context,
	content, itemType string,
	itemID uint,
	userID *uint,
	anonymousName, anonymousEmail string,
	parentID *uint,
) (*domain.Comment, error) {
	// 检查内容是否为空
	if strings.TrimSpace(content) == "" {
		return nil, errors.New("评论内容不能为空")
	}

	// 检查内容类型是否有效
	if itemType != "article" && itemType != "project" && itemType != "tool" {
		return nil, errors.New("无效的内容类型")
	}

	// 检查用户身份
	if userID == nil && (anonymousName == "" || anonymousEmail == "") {
		return nil, errors.New("需要提供用户身份或匿名信息")
	}

	// 检查父评论是否存在
	if parentID != nil {
		parent, err := s.commentRepo.FindByID(ctx, *parentID)
		if err != nil {
			return nil, err
		}
		if parent == nil {
			return nil, errors.New("父评论不存在")
		}
	}

	// 创建评论
	comment := &domain.Comment{
		Content:        content,
		UserID:         userID,
		AnonymousName:  anonymousName,
		AnonymousEmail: anonymousEmail,
		ItemType:       itemType,
		ItemID:         itemID,
		ParentID:       parentID,
		Status:         "pending", // 默认为待审核状态
	}

	// 如果有用户ID，则加载用户信息
	if userID != nil {
		user, err := s.userRepo.FindByID(ctx, *userID)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, errors.New("用户不存在")
		}
		
		// 如果是管理员或者普通用户，自动批准评论
		if user.Role == "admin" || user.Role == "editor" {
			comment.Status = "approved"
		}
		
		comment.User = user
	}

	// 保存评论
	err := s.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// UpdateComment 更新评论
func (s *CommentServiceImpl) UpdateComment(
	ctx context.Context,
	id uint,
	content string,
	userID uint,
) (*domain.Comment, error) {
	// 检查内容是否为空
	if strings.TrimSpace(content) == "" {
		return nil, errors.New("评论内容不能为空")
	}

	// 获取评论
	comment, err := s.commentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if comment == nil {
		return nil, errors.New("评论不存在")
	}

	// 检查用户权限
	if comment.UserID == nil || *comment.UserID != userID {
		// 获取用户信息，检查是否为管理员
		user, err := s.userRepo.FindByID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if user == nil || user.Role != "admin" {
			return nil, errors.New("无权修改此评论")
		}
	}

	// 更新评论内容
	comment.Content = content

	// 保存更新
	err = s.commentRepo.Update(ctx, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// DeleteComment 删除评论
func (s *CommentServiceImpl) DeleteComment(ctx context.Context, id uint, userID uint, isAdmin bool) error {
	// 获取评论
	comment, err := s.commentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if comment == nil {
		return errors.New("评论不存在")
	}

	// 检查权限
	if !isAdmin && (comment.UserID == nil || *comment.UserID != userID) {
		return errors.New("无权删除此评论")
	}

	// 删除评论
	return s.commentRepo.Delete(ctx, id)
}

// LikeComment 点赞评论
func (s *CommentServiceImpl) LikeComment(ctx context.Context, id uint) error {
	// 检查评论是否存在
	comment, err := s.commentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if comment == nil {
		return errors.New("评论不存在")
	}

	// 增加点赞数
	return s.commentRepo.IncrementLikes(ctx, id)
}

// ApproveComment 批准评论
func (s *CommentServiceImpl) ApproveComment(ctx context.Context, id uint) error {
	// 检查评论是否存在
	comment, err := s.commentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if comment == nil {
		return errors.New("评论不存在")
	}

	// 更新评论状态
	return s.commentRepo.UpdateStatus(ctx, id, "approved")
}

// MarkCommentAsSpam 标记评论为垃圾信息
func (s *CommentServiceImpl) MarkCommentAsSpam(ctx context.Context, id uint) error {
	// 检查评论是否存在
	comment, err := s.commentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if comment == nil {
		return errors.New("评论不存在")
	}

	// 更新评论状态
	return s.commentRepo.UpdateStatus(ctx, id, "spam")
} 