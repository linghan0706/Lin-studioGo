package domain

import (
	"database/sql"
	"time"
	
	"gorm.io/gorm"
)

// Article 文章模型
type Article struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Title         string    `gorm:"size:255;not null" json:"title"`
	Slug          string    `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	Excerpt       string    `gorm:"type:text" json:"excerpt,omitempty"`
	Content       string    `gorm:"type:longtext;not null" json:"content"`
	ContentHTML   string    `gorm:"type:longtext" json:"content_html,omitempty"` // 渲染后的HTML内容
	AuthorID      uint      `gorm:"not null" json:"author_id"`
	Author        *User     `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	CategoryID    *uint     `json:"category_id,omitempty"`
	Category      *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	CoverImage    string    `gorm:"size:255" json:"cover_image,omitempty"`
	ReadTime      uint16    `gorm:"default:0" json:"read_time"`
	Views         uint      `gorm:"default:0" json:"views"`
	Likes         uint      `gorm:"default:0" json:"likes"`
	CommentsCount uint      `gorm:"default:0" json:"comments_count"`
	Status        string    `gorm:"type:enum('draft','published','archived');default:'draft'" json:"status"`
	FeaturedOrder *uint8    `json:"featured_order,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	PublishedAt   sql.NullTime `json:"published_at,omitempty"`
	Tags          []Tag     `gorm:"many2many:article_tags;" json:"tags,omitempty"`
}

// TableName 指定表名
func (Article) TableName() string {
	return "articles"
}

// BeforeSave GORM钩子 - 保存前的操作
func (a *Article) BeforeSave(tx *gorm.DB) error {
	// 如果状态变为published且发布时间为空，则设置发布时间
	if a.Status == "published" && !a.PublishedAt.Valid {
		a.PublishedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}
	return nil
}

// ArticleResponse 文章响应数据
type ArticleResponse struct {
	ID            uint             `json:"id"`
	Title         string           `json:"title"`
	Slug          string           `json:"slug"`
	Excerpt       string           `json:"excerpt,omitempty"`
	Content       string           `json:"content,omitempty"` // 根据请求选择是否包含全文
	ContentHTML   string           `json:"content_html,omitempty"` // 渲染后的HTML内容
	Author        UserResponse     `json:"author"`
	Category      *CategoryResponse `json:"category,omitempty"`
	Tags          []SimpleTagResponse `json:"tags,omitempty"`
	CoverImage    string           `json:"cover_image,omitempty"`
	ReadTime      uint16           `json:"read_time"`
	Views         uint             `json:"views"`
	Likes         uint             `json:"likes"`
	CommentsCount uint             `json:"comments_count"`
	Status        string           `json:"status"`
	FeaturedOrder *uint8           `json:"featured_order,omitempty"`
	PublishedAt   *time.Time       `json:"published_at,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

// ToResponse 将文章模型转换为响应数据
func (a *Article) ToResponse(includeContent bool) ArticleResponse {
	var publishedAt *time.Time
	if a.PublishedAt.Valid {
		publishedAt = &a.PublishedAt.Time
	}

	var categoryResponse *CategoryResponse
	if a.Category != nil {
		category := a.Category.ToResponse(false)
		categoryResponse = &category
	}

	// 转换标签
	tags := make([]SimpleTagResponse, 0)
	if len(a.Tags) > 0 {
		for _, tag := range a.Tags {
			tags = append(tags, tag.ToSimpleResponse())
		}
	}

	response := ArticleResponse{
		ID:            a.ID,
		Title:         a.Title,
		Slug:          a.Slug,
		Excerpt:       a.Excerpt,
		Category:      categoryResponse,
		Tags:          tags,
		CoverImage:    a.CoverImage,
		ReadTime:      a.ReadTime,
		Views:         a.Views,
		Likes:         a.Likes,
		CommentsCount: a.CommentsCount,
		Status:        a.Status,
		FeaturedOrder: a.FeaturedOrder,
		PublishedAt:   publishedAt,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}

	// 如果作者信息可用
	if a.Author != nil {
		response.Author = a.Author.ToResponse()
	}

	// 仅在需要时包含内容
	if includeContent {
		response.Content = a.Content
		response.ContentHTML = a.ContentHTML // 添加HTML内容
	}

	return response
}

// ArticleListResponse 文章列表响应数据
type ArticleListResponse struct {
	ID            uint             `json:"id"`
	Title         string           `json:"title"`
	Slug          string           `json:"slug"`
	Excerpt       string           `json:"excerpt,omitempty"`
	Author        SimpleUserResponse `json:"author"`
	Category      *SimpleCategoryResponse `json:"category,omitempty"`
	Tags          []SimpleTagResponse `json:"tags,omitempty"`
	CoverImage    string           `json:"cover_image,omitempty"`
	ReadTime      uint16           `json:"read_time"`
	Views         uint             `json:"views"`
	Likes         uint             `json:"likes"`
	CommentsCount uint             `json:"comments_count"`
	Status        string           `json:"status"`
	PublishedAt   *time.Time       `json:"published_at,omitempty"`
}

// ToListResponse 将文章模型转换为列表响应数据
func (a *Article) ToListResponse() ArticleListResponse {
	var publishedAt *time.Time
	if a.PublishedAt.Valid {
		publishedAt = &a.PublishedAt.Time
	}

	var categoryResponse *SimpleCategoryResponse
	if a.Category != nil {
		categoryResponse = &SimpleCategoryResponse{
			ID:   a.Category.ID,
			Name: a.Category.Name,
		}
	}

	// 转换标签
	tags := make([]SimpleTagResponse, 0)
	if len(a.Tags) > 0 {
		for _, tag := range a.Tags {
			tags = append(tags, tag.ToSimpleResponse())
		}
	}

	// 创建响应
	response := ArticleListResponse{
		ID:            a.ID,
		Title:         a.Title,
		Slug:          a.Slug,
		Excerpt:       a.Excerpt,
		Category:      categoryResponse,
		Tags:          tags,
		CoverImage:    a.CoverImage,
		ReadTime:      a.ReadTime,
		Views:         a.Views,
		Likes:         a.Likes,
		CommentsCount: a.CommentsCount,
		Status:        a.Status,
		PublishedAt:   publishedAt,
	}

	// 如果作者信息可用，则设置作者
	if a.Author != nil {
		response.Author = SimpleUserResponse{
			ID:       a.Author.ID,
			Username: a.Author.Username,
			Avatar:   a.Author.Avatar,
		}
	}

	return response
}

// SimpleUserResponse 简化的用户响应数据
type SimpleUserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar,omitempty"`
}