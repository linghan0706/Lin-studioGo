package domain

import "time"

// ArticleTag 文章标签关联
type ArticleTag struct {
	ArticleID uint      `gorm:"primaryKey;column:article_id"`
	TagID     uint      `gorm:"primaryKey;column:tag_id"`
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
}

// TableName 表名
func (ArticleTag) TableName() string {
	return "article_tags"
} 