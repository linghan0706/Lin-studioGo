package domain

import (
	"time"
)

// Category 分类模型
type Category struct {
	ID          uint       `gorm:"primaryKey;column:id" json:"id"`
	Name        string     `gorm:"column:name;size:50;not null" json:"name"`
	Slug        string     `gorm:"column:slug;size:50;uniqueIndex;not null" json:"slug"`
	Description string     `gorm:"column:description;type:text" json:"description"`
	ParentID    *uint      `gorm:"column:parent_id" json:"parent_id"`
	Parent      *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	CreatedAt   time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	// 用于统计的虚拟字段
	ArticleCount int64 `gorm:"-" json:"article_count,omitempty"`
}

// TableName 表名
func (Category) TableName() string {
	return "categories"
}

// CategoryResponse 分类响应结构
type CategoryResponse struct {
	ID           uint               `json:"id"`
	Name         string             `json:"name"`
	Slug         string             `json:"slug"`
	Description  string             `json:"description"`
	ParentID     *uint              `json:"parent_id"`
	Children     []CategoryResponse `json:"children,omitempty"`
	ArticleCount int64              `json:"article_count"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

// ToResponse 将分类模型转换为响应数据
func (c *Category) ToResponse(includeChildren bool) CategoryResponse {
	response := CategoryResponse{
		ID:           c.ID,
		Name:         c.Name,
		Slug:         c.Slug,
		Description:  c.Description,
		ParentID:     c.ParentID,
		ArticleCount: c.ArticleCount,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
	
	// 仅当需要包含子分类且子分类不为空时，转换子分类
	if includeChildren && len(c.Children) > 0 {
		response.Children = make([]CategoryResponse, len(c.Children))
		for i, child := range c.Children {
			childCopy := child // 创建副本避免循环引用问题
			response.Children[i] = childCopy.ToResponse(false) // 子分类不再递归
		}
	}
	
	return response
}

// SimpleCategoryResponse 简化的分类响应数据
type SimpleCategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ToSimpleResponse 将分类模型转换为简化响应数据
func (c *Category) ToSimpleResponse() SimpleCategoryResponse {
	return SimpleCategoryResponse{
		ID:   c.ID,
		Name: c.Name,
	}
} 