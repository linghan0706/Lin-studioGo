package main

import (
	"Lin_studio/internal/config"
	"Lin_studio/internal/domain"
	"Lin_studio/internal/utils"
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

func main() {
	// 初始化配置
	config.GetConfig()

	// 初始化数据库连接
	config.InitDB()

	// 执行迁移
	if err := migrateMarkdownContent(); err != nil {
		log.Fatalf("迁移Markdown内容失败: %v", err)
	}

	log.Println("Markdown内容迁移完成!")
}

// migrateMarkdownContent 将所有文章的Markdown内容渲染为HTML并更新数据库
func migrateMarkdownContent() error {
	db := config.DB
	ctx := context.Background()

	log.Println("开始迁移Markdown内容...")
	startTime := time.Now()

	// 查询所有有内容的文章
	var articles []domain.Article
	if err := db.WithContext(ctx).Where("content IS NOT NULL AND content != ''").Find(&articles).Error; err != nil {
		return fmt.Errorf("查询文章失败: %w", err)
	}

	log.Printf("找到 %d 篇需要处理的文章", len(articles))

	// 批量更新计数
	processed := 0
	errors := 0

	// 为每篇文章生成HTML内容
	for i, article := range articles {
		if article.Content == "" {
			continue
		}

		// 渲染Markdown为HTML
		html, err := utils.RenderMarkdown(article.Content)
		if err != nil {
			log.Printf("警告: 文章ID=%d 渲染失败: %v", article.ID, err)
			errors++
			continue
		}

		// 更新文章HTML内容
		if err := db.WithContext(ctx).Model(&domain.Article{}).
			Where("id = ?", article.ID).
			Update("content_html", html).Error; err != nil {
			log.Printf("警告: 文章ID=%d 更新失败: %v", article.ID, err)
			errors++
			continue
		}

		processed++

		// 每100篇文章输出一次进度
		if (i+1)%100 == 0 || i == len(articles)-1 {
			log.Printf("进度: %d/%d 处理完成", i+1, len(articles))
		}
	}

	duration := time.Since(startTime)
	log.Printf("迁移完成! 处理: %d, 成功: %d, 失败: %d, 用时: %s",
		len(articles), processed, errors, duration)

	return nil
}

// batchUpdate 批量更新文章（可选的批量处理方式）
func batchUpdate(db *gorm.DB, ctx context.Context, updates []map[string]interface{}, batchSize int) error {
	if len(updates) == 0 {
		return nil
	}

	totalBatches := (len(updates) + batchSize - 1) / batchSize
	log.Printf("开始批量更新，共 %d 批", totalBatches)

	for i := 0; i < len(updates); i += batchSize {
		end := i + batchSize
		if end > len(updates) {
			end = len(updates)
		}

		batch := updates[i:end]
		tx := db.WithContext(ctx).Begin()
		
		for _, update := range batch {
			id := update["id"].(uint)
			html := update["content_html"].(string)
			
			if err := tx.Model(&domain.Article{}).
				Where("id = ?", id).
				Update("content_html", html).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("批量更新失败: %w", err)
			}
		}
		
		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("提交事务失败: %w", err)
		}
		
		log.Printf("批次 %d/%d 更新完成", (i/batchSize)+1, totalBatches)
	}

	return nil
} 