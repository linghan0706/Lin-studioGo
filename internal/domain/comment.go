package domain

import (
	"time"
)

// Comment 评论模型
type Comment struct {
	ID             uint      `gorm:"primaryKey;column:id" json:"id"`
	Content        string    `gorm:"column:content;type:text;not null" json:"content"`
	UserID         *uint     `gorm:"column:user_id" json:"user_id,omitempty"`
	User           *User     `gorm:"-" json:"user,omitempty"`
	AnonymousName  string    `gorm:"column:anonymous_name;size:100" json:"anonymous_name,omitempty"`
	AnonymousEmail string    `gorm:"column:anonymous_email;size:100" json:"anonymous_email,omitempty"`
	ItemType       string    `gorm:"column:item_type;type:enum('article','project','tool');not null" json:"item_type"`
	ItemID         uint      `gorm:"column:item_id;not null" json:"item_id"`
	ParentID       *uint     `gorm:"column:parent_id" json:"parent_id,omitempty"`
	Likes          uint      `gorm:"column:likes;default:0" json:"likes"`
	Status         string    `gorm:"column:status;type:enum('pending','approved','spam','deleted');default:'pending'" json:"status"`
	CreatedAt      time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	// 非数据库字段
	Replies     []Comment `gorm:"-" json:"replies,omitempty"`
	ReplyCount  int64     `gorm:"-" json:"reply_count,omitempty"`
}

// TableName 表名
func (Comment) TableName() string {
	return "comments"
}

// CommentResponse 评论响应结构
type CommentResponse struct {
	ID         uint             `json:"id"`
	Content    string           `json:"content"`
	Author     *CommentAuthor   `json:"author,omitempty"`
	ItemType   string           `json:"item_type"`
	ItemID     uint             `json:"item_id"`
	ParentID   *uint            `json:"parent_id,omitempty"`
	Likes      uint             `json:"likes"`
	Status     string           `json:"status"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	Replies    []CommentResponse `json:"replies,omitempty"`
	ReplyCount int64            `json:"reply_count,omitempty"`
}

// CommentAuthor 评论作者信息
type CommentAuthor struct {
	ID       uint   `json:"id,omitempty"`
	Username string `json:"username"`
	Avatar   string `json:"avatar,omitempty"`
	IsAnonymous bool `json:"is_anonymous"`
}

// ToResponse 将评论模型转换为响应数据
func (c *Comment) ToResponse() CommentResponse {
	var author *CommentAuthor

	// 处理作者信息
	if c.UserID != nil && c.User != nil {
		// 注册用户
		author = &CommentAuthor{
			ID:       c.User.ID,
			Username: c.User.Username,
			Avatar:   c.User.Avatar,
			IsAnonymous: false,
		}
	} else if c.AnonymousName != "" {
		// 匿名用户
		author = &CommentAuthor{
			Username: c.AnonymousName,
			IsAnonymous: true,
		}
	}

	// 处理回复
	replies := make([]CommentResponse, 0)
	if len(c.Replies) > 0 {
		for _, reply := range c.Replies {
			replies = append(replies, reply.ToResponse())
		}
	}

	return CommentResponse{
		ID:         c.ID,
		Content:    c.Content,
		Author:     author,
		ItemType:   c.ItemType,
		ItemID:     c.ItemID,
		ParentID:   c.ParentID,
		Likes:      c.Likes,
		Status:     c.Status,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
		Replies:    replies,
		ReplyCount: c.ReplyCount,
	}
}

// SimpleCommentResponse 简化的评论响应结构
type SimpleCommentResponse struct {
	ID        uint          `json:"id"`
	Content   string        `json:"content"`
	Author    *CommentAuthor `json:"author,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
}

// ToSimpleResponse 将评论模型转换为简化响应数据
func (c *Comment) ToSimpleResponse() SimpleCommentResponse {
	var author *CommentAuthor

	// 处理作者信息
	if c.UserID != nil && c.User != nil {
		// 注册用户
		author = &CommentAuthor{
			ID:       c.User.ID,
			Username: c.User.Username,
			Avatar:   c.User.Avatar,
			IsAnonymous: false,
		}
	} else if c.AnonymousName != "" {
		// 匿名用户
		author = &CommentAuthor{
			Username: c.AnonymousName,
			IsAnonymous: true,
		}
	}

	return SimpleCommentResponse{
		ID:        c.ID,
		Content:   c.Content,
		Author:    author,
		CreatedAt: c.CreatedAt,
	}
} 