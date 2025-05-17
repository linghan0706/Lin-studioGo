package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Project 项目模型
type Project struct {
	ID             uint           `gorm:"primaryKey;column:id" json:"id"`
	Title          string         `gorm:"column:title;size:255;not null" json:"title"`
	Slug           string         `gorm:"column:slug;size:255;uniqueIndex;not null" json:"slug"`
	Description    string         `gorm:"column:description;type:text" json:"description"`
	Content        string         `gorm:"column:content;type:text;not null" json:"content"`
	Technologies   JSONSlice      `gorm:"column:technologies;type:json" json:"technologies"`
	Features       JSONSlice      `gorm:"column:features;type:json" json:"features"`
	Images         JSONImages     `gorm:"column:images;type:json" json:"images"`
	Links          JSONMap        `gorm:"column:links;type:json" json:"links"`
	Status         string         `gorm:"column:status;type:enum('planning','in-progress','completed','archived');default:'in-progress'" json:"status"`
	Views          uint           `gorm:"column:views;default:0" json:"views"`
	Likes          uint           `gorm:"column:likes;default:0" json:"likes"`
	CommentsCount  uint           `gorm:"column:comments_count;default:0" json:"comments_count"`
	FeaturedOrder  *uint          `gorm:"column:featured_order" json:"featured_order,omitempty"`
	CompletionDate *time.Time     `gorm:"column:completion_date" json:"completion_date,omitempty"`
	CreatedAt      time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	Tags           []Tag          `gorm:"-" json:"tags,omitempty"`
	Members        []ProjectMember `gorm:"-" json:"members,omitempty"`
}

// ProjectMember 项目成员模型
type ProjectMember struct {
	ProjectID uint      `gorm:"primaryKey;column:project_id" json:"project_id"`
	UserID    uint      `gorm:"primaryKey;column:user_id" json:"user_id"`
	Role      string    `gorm:"column:role;default:'member'" json:"role"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	User      *User     `gorm:"-" json:"user,omitempty"`
}

// JSONSlice 自定义JSON数组类型
type JSONSlice []string

// Scan 实现sql.Scanner接口
func (j *JSONSlice) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言错误")
	}
	if len(bytes) == 0 {
		*j = make([]string, 0)
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value 实现driver.Valuer接口
func (j JSONSlice) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

// JSONMap 自定义JSON对象类型
type JSONMap map[string]string

// Scan 实现sql.Scanner接口
func (j *JSONMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言错误")
	}
	if len(bytes) == 0 {
		*j = make(map[string]string)
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value 实现driver.Valuer接口
func (j JSONMap) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

// ImageItem 图片项
type ImageItem struct {
	URL     string `json:"url"`
	Caption string `json:"caption,omitempty"`
}

// JSONImages 自定义JSON图片数组类型
type JSONImages []ImageItem

// Scan 实现sql.Scanner接口
func (j *JSONImages) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言错误")
	}
	if len(bytes) == 0 {
		*j = make([]ImageItem, 0)
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value 实现driver.Valuer接口
func (j JSONImages) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

// TableName 表名
func (Project) TableName() string {
	return "projects"
}

// TableName 表名
func (ProjectMember) TableName() string {
	return "project_members"
}

// ProjectResponse 项目响应结构
type ProjectResponse struct {
	ID             uint                  `json:"id"`
	Title          string                `json:"title"`
	Slug           string                `json:"slug"`
	Description    string                `json:"description"`
	Content        string                `json:"content,omitempty"`
	Technologies   []string              `json:"technologies"`
	Features       []string              `json:"features"`
	Images         []ImageItem           `json:"images"`
	Links          map[string]string     `json:"links"`
	Status         string                `json:"status"`
	Views          uint                  `json:"views"`
	Likes          uint                  `json:"likes"`
	CommentsCount  uint                  `json:"comments_count"`
	FeaturedOrder  *uint                 `json:"featured_order,omitempty"`
	CompletionDate *time.Time            `json:"completion_date,omitempty"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	Tags           []SimpleTagResponse   `json:"tags,omitempty"`
	Members        []MemberResponse      `json:"members,omitempty"`
}

// MemberResponse 成员响应结构
type MemberResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar,omitempty"`
	Role     string `json:"role"`
}

// ToResponse 将项目模型转换为响应数据
func (p *Project) ToResponse(includeContent bool) ProjectResponse {
	resp := ProjectResponse{
		ID:             p.ID,
		Title:          p.Title,
		Slug:           p.Slug,
		Description:    p.Description,
		Technologies:   p.Technologies,
		Features:       p.Features,
		Images:         p.Images,
		Links:          p.Links,
		Status:         p.Status,
		Views:          p.Views,
		Likes:          p.Likes,
		CommentsCount:  p.CommentsCount,
		FeaturedOrder:  p.FeaturedOrder,
		CompletionDate: p.CompletionDate,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}

	if includeContent {
		resp.Content = p.Content
	}

	// 转换标签
	if len(p.Tags) > 0 {
		resp.Tags = make([]SimpleTagResponse, len(p.Tags))
		for i, tag := range p.Tags {
			resp.Tags[i] = tag.ToSimpleResponse()
		}
	}

	// 转换成员
	if len(p.Members) > 0 {
		resp.Members = make([]MemberResponse, len(p.Members))
		for i, member := range p.Members {
			if member.User != nil {
				resp.Members[i] = MemberResponse{
					ID:       member.User.ID,
					Username: member.User.Username,
					Avatar:   member.User.Avatar,
					Role:     member.Role,
				}
			}
		}
	}

	return resp
}

// SimpleProjectResponse 简化的项目响应结构
type SimpleProjectResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToSimpleResponse 将项目模型转换为简化响应数据
func (p *Project) ToSimpleResponse() SimpleProjectResponse {
	return SimpleProjectResponse{
		ID:          p.ID,
		Title:       p.Title,
		Slug:        p.Slug,
		Description: p.Description,
		Status:      p.Status,
		CreatedAt:   p.CreatedAt,
	}
} 