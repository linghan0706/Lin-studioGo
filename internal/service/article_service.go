package service

import (
	"Lin_studio/internal/domain"
	"Lin_studio/internal/repository"
	"Lin_studio/internal/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ArticleService 文章服务接口
type ArticleService interface {
	GetArticles(ctx context.Context, page, limit int, categoryID, tagID, authorID *uint, status, search, sort *string, renderHTML bool) ([]domain.Article, domain.PaginationData, error)
	GetArticleByID(ctx context.Context, id uint, renderHTML bool) (*domain.Article, error)
	GetArticleBySlug(ctx context.Context, slug string, renderHTML bool) (*domain.Article, error)
	CreateArticle(ctx context.Context, title, excerpt, content string, authorID, categoryID uint, tagIDs []uint, coverImage, status string) (*domain.Article, error)
	UpdateArticle(ctx context.Context, id uint, title, excerpt, content string, categoryID uint, tagIDs []uint, coverImage, status string) (*domain.Article, error)
	DeleteArticle(ctx context.Context, id uint) error
	UploadCoverImage(ctx context.Context, file *multipart.FileHeader) (string, error)
	GetFeaturedArticles(ctx context.Context, limit int, renderHTML bool) ([]domain.Article, error)
	ViewArticle(ctx context.Context, id uint) error
	LikeArticle(ctx context.Context, id uint) error
}

// ArticleServiceImpl 文章服务实现
type ArticleServiceImpl struct {
	articleRepo repository.ArticleRepository
	tagRepo     repository.TagRepository
	userRepo    repository.UserRepository
	categoryRepo repository.CategoryRepository
}

// NewArticleService 创建文章服务实例
func NewArticleService(
	articleRepo repository.ArticleRepository,
	tagRepo repository.TagRepository,
	userRepo repository.UserRepository,
	categoryRepo repository.CategoryRepository,
) ArticleService {
	return &ArticleServiceImpl{
		articleRepo: articleRepo,
		tagRepo:     tagRepo,
		userRepo:    userRepo,
		categoryRepo: categoryRepo,
	}
}

// GetArticles 获取文章列表
func (s *ArticleServiceImpl) GetArticles(
	ctx context.Context,
	page, limit int,
	categoryID, tagID, authorID *uint,
	status, search, sort *string,
	renderHTML bool,
) ([]domain.Article, domain.PaginationData, error) {
	// 创建过滤条件
	filter := repository.ArticleFilter{
		Page:       page,
		Limit:      limit,
		CategoryID: categoryID,
		TagID:      tagID,
		AuthorID:   authorID,
		Status:     status,
		Search:     search,
		Sort:       sort,
	}

	// 查询文章
	articles, total, err := s.articleRepo.FindAll(ctx, filter)
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

	// 加载相关数据
	for i := range articles {
		// 加载作者信息
		author, err := s.userRepo.FindByID(ctx, articles[i].AuthorID)
		if err == nil && author != nil {
			articles[i].Author = author
		}

		// 加载分类信息
		if articles[i].CategoryID != nil {
			category, err := s.categoryRepo.FindByID(ctx, *articles[i].CategoryID)
			if err == nil && category != nil {
				articles[i].Category = category
			}
		}

		// 加载标签信息
		tags, err := s.tagRepo.FindByArticleID(ctx, articles[i].ID)
		if err == nil {
			articles[i].Tags = tags
		}
		
		// 如果需要渲染HTML
		if renderHTML && articles[i].Content != "" {
			html, err := utils.RenderMarkdown(articles[i].Content)
			if err == nil {
				articles[i].ContentHTML = html
			}
		}
	}

	return articles, pagination, nil
}

// GetArticleByID 根据ID获取文章
func (s *ArticleServiceImpl) GetArticleByID(ctx context.Context, id uint, renderHTML bool) (*domain.Article, error) {
	article, err := s.articleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if article == nil {
		return nil, errors.New("文章不存在")
	}

	// 加载相关数据
	s.loadArticleRelations(ctx, article)
	
	// 如果需要渲染HTML
	if renderHTML && article.Content != "" {
		html, err := utils.RenderMarkdown(article.Content)
		if err == nil {
			article.ContentHTML = html
		}
	}

	return article, nil
}

// GetArticleBySlug 根据Slug获取文章
func (s *ArticleServiceImpl) GetArticleBySlug(ctx context.Context, slug string, renderHTML bool) (*domain.Article, error) {
	article, err := s.articleRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	if article == nil {
		return nil, errors.New("文章不存在")
	}

	// 加载相关数据
	s.loadArticleRelations(ctx, article)
	
	// 如果需要渲染HTML
	if renderHTML && article.Content != "" {
		html, err := utils.RenderMarkdown(article.Content)
		if err == nil {
			article.ContentHTML = html
		}
	}

	return article, nil
}

// CreateArticle 创建文章
func (s *ArticleServiceImpl) CreateArticle(
	ctx context.Context,
	title, excerpt, content string,
	authorID, categoryID uint,
	tagIDs []uint,
	coverImage, status string,
) (*domain.Article, error) {
	// 检查标题是否为空
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("文章标题不能为空")
	}

	// 检查内容是否为空
	if strings.TrimSpace(content) == "" {
		return nil, errors.New("文章内容不能为空")
	}

	// 检查作者是否存在
	author, err := s.userRepo.FindByID(ctx, authorID)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, errors.New("作者不存在")
	}

	// 检查分类是否存在
	var catID *uint
	if categoryID > 0 {
		category, err := s.categoryRepo.FindByID(ctx, categoryID)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, errors.New("分类不存在")
		}
		catID = &categoryID
	}

	// 生成slug
	slug := s.generateSlug(title)

	// 创建文章
	now := time.Now()
	var publishedAt sql.NullTime
	if status == "published" {
		publishedAt = sql.NullTime{
			Time:  now,
			Valid: true,
		}
	}

	// 渲染Markdown内容为HTML
	contentHTML, err := utils.RenderMarkdown(content)
	if err != nil {
		return nil, fmt.Errorf("渲染Markdown内容失败: %w", err)
	}

	article := &domain.Article{
		Title:       title,
		Slug:        slug,
		Excerpt:     excerpt,
		Content:     content,
		ContentHTML: contentHTML, // 添加渲染后的HTML内容
		AuthorID:    authorID,
		CategoryID:  catID,
		CoverImage:  coverImage,
		Status:      status,
		PublishedAt: publishedAt,
	}

	// 保存文章
	err = s.articleRepo.Create(ctx, article)
	if err != nil {
		return nil, err
	}

	// 添加标签关联
	if len(tagIDs) > 0 {
		for _, tagID := range tagIDs {
			err := s.tagRepo.AddArticleTag(ctx, article.ID, tagID)
			if err != nil {
				// 继续处理其他标签，不返回错误
				continue
			}
		}
	}

	return article, nil
}

// UpdateArticle 更新文章
func (s *ArticleServiceImpl) UpdateArticle(
	ctx context.Context,
	id uint,
	title, excerpt, content string,
	categoryID uint,
	tagIDs []uint,
	coverImage, status string,
) (*domain.Article, error) {
	// 检查文章是否存在
	article, err := s.articleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, errors.New("文章不存在")
	}

	// 检查标题是否为空
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("文章标题不能为空")
	}

	// 检查内容是否为空
	if strings.TrimSpace(content) == "" {
		return nil, errors.New("文章内容不能为空")
	}

	// 检查分类是否存在
	var catID *uint
	if categoryID > 0 {
		category, err := s.categoryRepo.FindByID(ctx, categoryID)
		if err != nil {
			return nil, err
		}
		if category == nil {
			return nil, errors.New("分类不存在")
		}
		catID = &categoryID
	}

	// 更新发布状态
	var publishedAt sql.NullTime
	if status == "published" && article.Status != "published" {
		now := time.Now()
		publishedAt = sql.NullTime{
			Time:  now,
			Valid: true,
		}
	} else if article.Status == "published" {
		publishedAt = article.PublishedAt
	}

	// 检查内容是否有更改，有则重新渲染HTML
	contentHTML := article.ContentHTML
	if article.Content != content {
		html, err := utils.RenderMarkdown(content)
		if err != nil {
			return nil, fmt.Errorf("渲染Markdown内容失败: %w", err)
		}
		contentHTML = html
	}

	// 更新文章字段
	article.Title = title
	article.Excerpt = excerpt
	article.Content = content
	article.ContentHTML = contentHTML // 更新渲染后的HTML内容
	article.CategoryID = catID
	if coverImage != "" {
		article.CoverImage = coverImage
	}
	article.Status = status
	article.PublishedAt = publishedAt

	// 保存更新
	err = s.articleRepo.Update(ctx, article)
	if err != nil {
		return nil, err
	}

	// 更新标签关联
	if len(tagIDs) > 0 {
		// 先清除原有标签
		err = s.tagRepo.ClearArticleTags(ctx, article.ID)
		if err != nil {
			return nil, err
		}

		// 添加新标签
		for _, tagID := range tagIDs {
			err := s.tagRepo.AddArticleTag(ctx, article.ID, tagID)
			if err != nil {
				// 继续处理其他标签
				continue
			}
		}
	}

	return article, nil
}

// DeleteArticle 删除文章
func (s *ArticleServiceImpl) DeleteArticle(ctx context.Context, id uint) error {
	// 检查文章是否存在
	article, err := s.articleRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if article == nil {
		return errors.New("文章不存在")
	}

	// 删除文章
	return s.articleRepo.Delete(ctx, id)
}

// UploadCoverImage 上传文章封面图片
func (s *ArticleServiceImpl) UploadCoverImage(ctx context.Context, file *multipart.FileHeader) (string, error) {
	// 确保上传目录存在
	uploadDir := "./uploads/covers"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	// 创建文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("cover_%d%s", time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// 保存文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// 将源文件内容复制到目标文件
	if _, err = dst.ReadFrom(src); err != nil {
		return "", err
	}

	// 返回文件URL
	return fmt.Sprintf("/uploads/covers/%s", filename), nil
}

// GetFeaturedArticles 获取精选文章
func (s *ArticleServiceImpl) GetFeaturedArticles(ctx context.Context, limit int, renderHTML bool) ([]domain.Article, error) {
	articles, err := s.articleRepo.FindFeatured(ctx, limit)
	if err != nil {
		return nil, err
	}

	// 加载相关数据
	for i := range articles {
		// 只加载作者信息和分类，不加载全部内容
		author, err := s.userRepo.FindByID(ctx, articles[i].AuthorID)
		if err == nil && author != nil {
			articles[i].Author = author
		}

		if articles[i].CategoryID != nil {
			category, err := s.categoryRepo.FindByID(ctx, *articles[i].CategoryID)
			if err == nil && category != nil {
				articles[i].Category = category
			}
		}
		
		// 如果需要渲染HTML
		if renderHTML && articles[i].Content != "" {
			html, err := utils.RenderMarkdown(articles[i].Content)
			if err == nil {
				articles[i].ContentHTML = html
			}
		}
	}

	return articles, nil
}

// ViewArticle 增加文章浏览量
func (s *ArticleServiceImpl) ViewArticle(ctx context.Context, id uint) error {
	return s.articleRepo.IncrementViews(ctx, id)
}

// LikeArticle 点赞文章
func (s *ArticleServiceImpl) LikeArticle(ctx context.Context, id uint) error {
	return s.articleRepo.IncrementLikes(ctx, id)
}

// loadArticleRelations 加载文章关联数据
func (s *ArticleServiceImpl) loadArticleRelations(ctx context.Context, article *domain.Article) {
	// 加载作者信息
	author, err := s.userRepo.FindByID(ctx, article.AuthorID)
	if err == nil && author != nil {
		article.Author = author
	}

	// 加载分类信息
	if article.CategoryID != nil {
		category, err := s.categoryRepo.FindByID(ctx, *article.CategoryID)
		if err == nil && category != nil {
			article.Category = category
		}
	}

	// 加载标签信息
	tags, err := s.tagRepo.FindByArticleID(ctx, article.ID)
	if err == nil {
		article.Tags = tags
	}
}

// generateSlug 生成文章别名
func (s *ArticleServiceImpl) generateSlug(title string) string {
	slug := strings.ToLower(title)
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