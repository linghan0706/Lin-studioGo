package handler

import (
	"Lin_studio/internal/domain"
	"Lin_studio/internal/service"
	"Lin_studio/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ArticleHandler 文章处理器
type ArticleHandler struct {
	articleService service.ArticleService
}

// NewArticleHandler 创建文章处理器实例
func NewArticleHandler(articleService service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
	}
}

// GetArticles 获取文章列表
func (h *ArticleHandler) GetArticles(c *gin.Context) {
	// 获取查询参数
	page, limit := utils.GetPagination(c)
	
	// 获取过滤参数
	var categoryID, tagID, authorID *uint
	var status, search, sort *string
	
	// 分类ID
	categoryIDStr := c.Query("category")
	if categoryIDStr != "" {
		id, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err == nil {
			catID := uint(id)
			categoryID = &catID
		}
	}
	
	// 标签ID
	tagIDStr := c.Query("tag")
	if tagIDStr != "" {
		id, err := strconv.ParseUint(tagIDStr, 10, 32)
		if err == nil {
			tID := uint(id)
			tagID = &tID
		}
	}
	
	// 作者ID
	authorIDStr := c.Query("author")
	if authorIDStr != "" {
		id, err := strconv.ParseUint(authorIDStr, 10, 32)
		if err == nil {
			aID := uint(id)
			authorID = &aID
		}
	}
	
	// 状态
	statusStr := c.Query("status")
	if statusStr != "" {
		status = &statusStr
	}
	
	// 搜索
	searchStr := c.Query("search")
	if searchStr != "" {
		search = &searchStr
	}
	
	// 排序
	sortStr := c.Query("sort")
	if sortStr != "" {
		sort = &sortStr
	}
	
	// 检查是否需要渲染Markdown
	renderHTML := c.Query("render_html") == "true"
	
	// 查询文章
	articles, pagination, err := h.articleService.GetArticles(
		c.Request.Context(),
		page,
		limit,
		categoryID,
		tagID,
		authorID,
		status,
		search,
		sort,
		renderHTML,
	)
	
	if err != nil {
		utils.InternalServerErrorResponse(c, "获取文章列表失败: "+err.Error())
		return
	}
	
	// 返回结果
	utils.SuccessResponse(c, "获取文章列表成功", gin.H{
		"articles":   articles,
		"pagination": pagination,
	})
}

// GetArticleBySlug 根据Slug获取文章
func (h *ArticleHandler) GetArticleBySlug(c *gin.Context) {
	// 获取slug
	slug := c.Param("slug")
	
	// 检查是否需要渲染Markdown
	renderHTML := c.Query("render_html") == "true"
	
	// 获取文章
	article, err := h.articleService.GetArticleBySlug(c.Request.Context(), slug, renderHTML)
	if err != nil {
		utils.NotFoundResponse(c, "文章不存在")
		return
	}
	
	// 增加浏览量
	go h.articleService.ViewArticle(c.Request.Context(), article.ID)
	
	// 返回文章
	utils.SuccessResponse(c, "获取文章成功", article)
}

// CreateArticle 创建文章
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.UnauthorizedResponse(c, "未授权")
		return
	}
	
	// 请求体结构
	var req struct {
		Title      string   `json:"title" binding:"required"`
		Excerpt    string   `json:"excerpt"`
		Content    string   `json:"content" binding:"required"`
		CategoryID uint     `json:"category_id"`
		Tags       []uint   `json:"tags"`
		CoverImage string   `json:"cover_image"`
		Status     string   `json:"status" binding:"required"`
	}
	
	// 绑定请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "无效的请求数据", err.Error())
		return
	}
	
	// 创建文章
	article, err := h.articleService.CreateArticle(
		c.Request.Context(),
		req.Title,
		req.Excerpt,
		req.Content,
		userID.(uint),
		req.CategoryID,
		req.Tags,
		req.CoverImage,
		req.Status,
	)
	
	if err != nil {
		utils.BadRequestResponse(c, "创建文章失败", err.Error())
		return
	}
	
	// 检查是否需要渲染Markdown
	renderHTML := c.Query("render_html") == "true"
	var fullArticle *domain.Article
	
	if renderHTML {
		// 获取完整文章信息，包含渲染后的HTML
		fullArticle, err = h.articleService.GetArticleByID(c.Request.Context(), article.ID, true)
		if err != nil {
			// 如果获取失败，仍然返回基本创建成功信息
			utils.CreatedResponse(c, "文章创建成功", gin.H{
				"id":    article.ID,
				"title": article.Title,
				"slug":  article.Slug,
			})
			return
		}
		
		// 返回完整文章信息
		utils.CreatedResponse(c, "文章创建成功", fullArticle)
		return
	}
	
	// 返回基本创建成功信息
	utils.CreatedResponse(c, "文章创建成功", gin.H{
		"id":    article.ID,
		"title": article.Title,
		"slug":  article.Slug,
	})
}

// UpdateArticle 更新文章
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	// 获取文章ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的文章ID", err.Error())
		return
	}
	
	// 请求体结构
	var req struct {
		Title      string   `json:"title" binding:"required"`
		Excerpt    string   `json:"excerpt"`
		Content    string   `json:"content" binding:"required"`
		CategoryID uint     `json:"category_id"`
		Tags       []uint   `json:"tags"`
		CoverImage string   `json:"cover_image"`
		Status     string   `json:"status" binding:"required"`
	}
	
	// 绑定请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "无效的请求数据", err.Error())
		return
	}
	
	// 更新文章
	article, err := h.articleService.UpdateArticle(
		c.Request.Context(),
		uint(id),
		req.Title,
		req.Excerpt,
		req.Content,
		req.CategoryID,
		req.Tags,
		req.CoverImage,
		req.Status,
	)
	
	if err != nil {
		utils.BadRequestResponse(c, "更新文章失败", err.Error())
		return
	}
	
	// 检查是否需要渲染Markdown
	renderHTML := c.Query("render_html") == "true"
	var fullArticle *domain.Article
	
	if renderHTML {
		// 获取完整文章信息，包含渲染后的HTML
		fullArticle, err = h.articleService.GetArticleByID(c.Request.Context(), article.ID, true)
		if err != nil {
			// 如果获取失败，仍然返回基本更新成功信息
			utils.SuccessResponse(c, "文章更新成功", gin.H{
				"id":    article.ID,
				"title": article.Title,
				"slug":  article.Slug,
			})
			return
		}
		
		// 返回完整文章信息
		utils.SuccessResponse(c, "文章更新成功", fullArticle)
		return
	}
	
	// 返回基本更新成功信息
	utils.SuccessResponse(c, "文章更新成功", gin.H{
		"id":    article.ID,
		"title": article.Title,
		"slug":  article.Slug,
	})
}

// DeleteArticle 删除文章
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	// 获取文章ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的文章ID", err.Error())
		return
	}
	
	// 删除文章
	err = h.articleService.DeleteArticle(c.Request.Context(), uint(id))
	if err != nil {
		utils.BadRequestResponse(c, "删除文章失败", err.Error())
		return
	}
	
	// 返回结果
	utils.SuccessResponse(c, "文章删除成功", nil)
}

// UploadCoverImage 上传文章封面图片
func (h *ArticleHandler) UploadCoverImage(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("cover")
	if err != nil {
		utils.BadRequestResponse(c, "上传文件失败", err.Error())
		return
	}
	
	// 上传封面图片
	url, err := h.articleService.UploadCoverImage(c.Request.Context(), file)
	if err != nil {
		utils.InternalServerErrorResponse(c, "上传封面图片失败: "+err.Error())
		return
	}
	
	// 返回URL
	utils.SuccessResponse(c, "封面上传成功", gin.H{
		"url": url,
	})
}

// GetFeaturedArticles 获取精选文章
func (h *ArticleHandler) GetFeaturedArticles(c *gin.Context) {
	// 获取限制数量
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 5
	}
	
	// 检查是否需要渲染Markdown
	renderHTML := c.Query("render_html") == "true"
	
	// 获取精选文章
	articles, err := h.articleService.GetFeaturedArticles(c.Request.Context(), limit, renderHTML)
	if err != nil {
		utils.InternalServerErrorResponse(c, "获取精选文章失败: "+err.Error())
		return
	}
	
	// 返回结果
	utils.SuccessResponse(c, "获取精选文章成功", gin.H{
		"articles": articles,
	})
}

// LikeArticle 点赞文章
func (h *ArticleHandler) LikeArticle(c *gin.Context) {
	// 获取文章ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的文章ID", err.Error())
		return
	}
	
	// 点赞文章
	err = h.articleService.LikeArticle(c.Request.Context(), uint(id))
	if err != nil {
		utils.InternalServerErrorResponse(c, "点赞文章失败: "+err.Error())
		return
	}
	
	// 返回结果
	utils.SuccessResponse(c, "点赞成功", nil)
} 