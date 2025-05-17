package domain

import (
	"time"
)

// Tag 标签模型
type Tag struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	Name        string    `gorm:"column:name;size:50;not null" json:"name"`
	Slug        string    `gorm:"column:slug;size:50;uniqueIndex;not null" json:"slug"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	Color       string    `gorm:"column:color;size:7" json:"color"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	// 用于统计的虚拟字段
	ArticleCount int64 `gorm:"-" json:"article_count,omitempty"`
	ProjectCount int64 `gorm:"-" json:"project_count,omitempty"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}

// TagResponse 标签响应数据
type TagResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  string    `json:"description"`
	Color        string    `json:"color"`
	ArticleCount int64     `json:"article_count"`
	ProjectCount int64     `json:"project_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ToResponse 将标签模型转换为响应数据
func (t *Tag) ToResponse() TagResponse {
	return TagResponse{
		ID:           t.ID,
		Name:         t.Name,
		Slug:         t.Slug,
		Description:  t.Description,
		Color:        t.Color,
		ArticleCount: t.ArticleCount,
		ProjectCount: t.ProjectCount,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}
}

// SimpleTagResponse 简化的标签响应数据
type SimpleTagResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ToSimpleResponse 将标签模型转换为简化响应数据
func (t *Tag) ToSimpleResponse() SimpleTagResponse {
	return SimpleTagResponse{
		ID:   t.ID,
		Name: t.Name,
		Slug: t.Slug,
	}
} 