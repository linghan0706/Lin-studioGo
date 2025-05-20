package domain

import (
	"database/sql"
	"encoding/json"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password     string         `gorm:"size:255;not null" json:"-"`
	Avatar       string         `gorm:"size:255" json:"avatar,omitempty"`
	Role         string         `gorm:"type:enum('admin','editor','user');default:'user'" json:"role"`
	Bio          sql.NullString `gorm:"type:text" json:"bio,omitempty"`
	SocialLinks  datatypes.JSON `gorm:"type:json" json:"social_links,omitempty"`
	ContactInfo  datatypes.JSON `gorm:"type:json" json:"contact_info,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	LastLogin    sql.NullTime   `json:"last_login,omitempty"`
	Status       string         `gorm:"type:enum('active','suspended','deleted');default:'active'" json:"status"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeSave GORM钩子 - 保存前的操作
func (u *User) BeforeSave(tx *gorm.DB) error {
	// 明文密码不需要特殊处理
	return nil
}

// CheckPassword 验证密码是否正确
func (u *User) CheckPassword(password string) bool {
	// 直接比较明文密码
	return u.Password == password
}

// SetPassword 设置密码
func (u *User) SetPassword(password string) error {
	// 直接存储明文密码
	u.Password = password
	return nil
}

// UserResponse 用户响应数据
type UserResponse struct {
	ID           uint              `json:"id"`
	Username     string            `json:"username"`
	Email        string            `json:"email"`
	Avatar       string            `json:"avatar,omitempty"`
	Role         string            `json:"role"`
	Bio          string            `json:"bio,omitempty"`
	SocialLinks  map[string]string `json:"social_links,omitempty"`
	ContactInfo  map[string]string `json:"contact_info,omitempty"`
	LastLogin    *time.Time        `json:"last_login,omitempty"`
	Status       string            `json:"status"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// ToResponse 将用户模型转换为响应数据
func (u *User) ToResponse() UserResponse {
	var bio string
	if u.Bio.Valid {
		bio = u.Bio.String
	}

	var lastLogin *time.Time
	if u.LastLogin.Valid {
		lastLogin = &u.LastLogin.Time
	}

	var socialLinks map[string]string
	// 这里可以添加社交链接的逻辑
	socialLinks = make(map[string]string)
	json.Unmarshal(u.SocialLinks, &socialLinks)

	var contactInfo map[string]string
	// 这里可以添加联系信息的逻辑
	contactInfo = make(map[string]string)
	json.Unmarshal(u.ContactInfo, &contactInfo)

	return UserResponse{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		Avatar:       u.Avatar,
		Role:         u.Role,
		Bio:          bio,
		SocialLinks:  socialLinks,
		ContactInfo:  contactInfo,
		LastLogin:    lastLogin,
		Status:       u.Status,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

// UpdateLastLogin 更新用户最后登录时间
func (u *User) UpdateLastLogin(db *gorm.DB) error {
	now := time.Now()
	u.LastLogin = sql.NullTime{
		Time:  now,
		Valid: true,
	}
	return db.Model(u).Update("last_login", now).Error
} 