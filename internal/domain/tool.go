package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Tool 工具模型
type Tool struct {
	ID          uint       `gorm:"primaryKey;column:id" json:"id"`
	Name        string     `gorm:"column:name;size:100;not null" json:"name"`
	Slug        string     `gorm:"column:slug;size:100;uniqueIndex;not null" json:"slug"`
	Description string     `gorm:"column:description;type:text" json:"description"`
	Icon        string     `gorm:"column:icon;size:50" json:"icon"`
	Category    string     `gorm:"column:category;size:50" json:"category"`
	Content     string     `gorm:"column:content;type:text" json:"content"`
	URL         string     `gorm:"column:url;size:255" json:"url"`
	Config      JSONConfig `gorm:"column:config;type:json" json:"config"`
	Views       uint       `gorm:"column:views;default:0" json:"views"`
	Status      string     `gorm:"column:status;type:enum('active','maintenance','deprecated');default:'active'" json:"status"`
	CreatedAt   time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// JSONConfig 自定义JSON配置类型
type JSONConfig map[string]interface{}

// Scan 实现sql.Scanner接口
func (j *JSONConfig) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("类型断言错误")
	}
	if len(bytes) == 0 {
		*j = make(map[string]interface{})
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value 实现driver.Valuer接口
func (j JSONConfig) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

// TableName 表名
func (Tool) TableName() string {
	return "tools"
}

// ToolResponse 工具响应结构
type ToolResponse struct {
	ID          uint                `json:"id"`
	Name        string              `json:"name"`
	Slug        string              `json:"slug"`
	Description string              `json:"description"`
	Icon        string              `json:"icon"`
	Category    string              `json:"category"`
	URL         string              `json:"url"`
	Content     string              `json:"content,omitempty"`
	Config      JSONConfig          `json:"config,omitempty"`
	Views       uint                `json:"views"`
	Status      string              `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// ToResponse 将工具模型转换为响应数据
func (t *Tool) ToResponse(includeContent bool) ToolResponse {
	resp := ToolResponse{
		ID:          t.ID,
		Name:        t.Name,
		Slug:        t.Slug,
		Description: t.Description,
		Icon:        t.Icon,
		Category:    t.Category,
		URL:         t.URL,
		Config:      t.Config,
		Views:       t.Views,
		Status:      t.Status,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}

	if includeContent {
		resp.Content = t.Content
	}

	return resp
}

// SimpleToolResponse 简化的工具响应结构
type SimpleToolResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Category    string    `json:"category"`
	URL         string    `json:"url"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// ToSimpleResponse 将工具模型转换为简化响应数据
func (t *Tool) ToSimpleResponse() SimpleToolResponse {
	return SimpleToolResponse{
		ID:          t.ID,
		Name:        t.Name,
		Slug:        t.Slug,
		Description: t.Description,
		Icon:        t.Icon,
		Category:    t.Category,
		URL:         t.URL,
		Status:      t.Status,
		CreatedAt:   t.CreatedAt,
	}
} 