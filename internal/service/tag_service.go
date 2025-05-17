package service

import (
	"Lin_studio/internal/domain"
	"Lin_studio/internal/repository"
	"context"
	"errors"
	"strings"
)

// TagService 标签服务接口
type TagService interface {
	GetAllTags(ctx context.Context) ([]domain.TagResponse, error)
	CreateTag(ctx context.Context, name, description, color string) (*domain.Tag, error)
	GetTagByID(ctx context.Context, id uint) (*domain.Tag, error)
	GetTagBySlug(ctx context.Context, slug string) (*domain.Tag, error)
}

// TagServiceImpl 标签服务实现
type TagServiceImpl struct {
	tagRepo repository.TagRepository
}

// NewTagService 创建标签服务实例
func NewTagService(tagRepo repository.TagRepository) TagService {
	return &TagServiceImpl{
		tagRepo: tagRepo,
	}
}

// GetAllTags 获取所有标签
func (s *TagServiceImpl) GetAllTags(ctx context.Context) ([]domain.TagResponse, error) {
	tags, err := s.tagRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	
	var response []domain.TagResponse
	for _, tag := range tags {
		articleCount, err := s.tagRepo.CountArticles(ctx, tag.ID)
		if err != nil {
			return nil, err
		}
		
		projectCount, err := s.tagRepo.CountProjects(ctx, tag.ID)
		if err != nil {
			return nil, err
		}
		
		response = append(response, domain.TagResponse{
			ID:           tag.ID,
			Name:         tag.Name,
			Slug:         tag.Slug,
			Description:  tag.Description,
			Color:        tag.Color,
			ArticleCount: articleCount,
			ProjectCount: projectCount,
		})
	}
	
	return response, nil
}

// CreateTag 创建标签
func (s *TagServiceImpl) CreateTag(ctx context.Context, name, description, color string) (*domain.Tag, error) {
	// 检查标签名是否为空
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("标签名不能为空")
	}
	
	// 验证颜色格式
	if color != "" && !strings.HasPrefix(color, "#") {
		color = "#" + color
	}
	
	// 生成slug
	slug := s.generateSlug(name)
	
	// 检查slug是否已存在
	existingTag, err := s.tagRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	
	if existingTag != nil {
		return nil, errors.New("标签已存在")
	}
	
	// 创建新标签
	tag := &domain.Tag{
		Name:        name,
		Slug:        slug,
		Description: description,
		Color:       color,
	}
	
	err = s.tagRepo.Create(ctx, tag)
	if err != nil {
		return nil, err
	}
	
	return tag, nil
}

// GetTagByID 根据ID获取标签
func (s *TagServiceImpl) GetTagByID(ctx context.Context, id uint) (*domain.Tag, error) {
	tag, err := s.tagRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if tag == nil {
		return nil, errors.New("标签不存在")
	}
	
	return tag, nil
}

// GetTagBySlug 根据Slug获取标签
func (s *TagServiceImpl) GetTagBySlug(ctx context.Context, slug string) (*domain.Tag, error) {
	tag, err := s.tagRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	
	if tag == nil {
		return nil, errors.New("标签不存在")
	}
	
	return tag, nil
}

// generateSlug 生成标签别名
func (s *TagServiceImpl) generateSlug(name string) string {
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