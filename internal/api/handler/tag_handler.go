package handler

import (
	"Lin_studio/internal/service"
	"Lin_studio/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TagHandler 标签处理器
type TagHandler struct {
	tagService service.TagService
}

// NewTagHandler 创建标签处理器实例
func NewTagHandler(tagService service.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// GetAllTags 获取所有标签
func (h *TagHandler) GetAllTags(c *gin.Context) {
	// 获取标签列表
	tags, err := h.tagService.GetAllTags(c.Request.Context())
	if err != nil {
		utils.InternalServerErrorResponse(c, "获取标签列表失败: "+err.Error())
		return
	}

	// 返回标签列表
	utils.SuccessResponse(c, "获取标签列表成功", gin.H{
		"tags": tags,
	})
}

// CreateTag 创建标签
func (h *TagHandler) CreateTag(c *gin.Context) {
	// 请求体结构
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Color       string `json:"color"`
	}

	// 绑定请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "无效的请求数据", err.Error())
		return
	}

	// 创建标签
	tag, err := h.tagService.CreateTag(
		c.Request.Context(),
		req.Name,
		req.Description,
		req.Color,
	)
	if err != nil {
		utils.BadRequestResponse(c, "创建标签失败", err.Error())
		return
	}

	// 返回创建的标签
	utils.CreatedResponse(c, "标签创建成功", gin.H{
		"id":   tag.ID,
		"name": tag.Name,
		"slug": tag.Slug,
	})
}

// GetTagByID 根据ID获取标签
func (h *TagHandler) GetTagByID(c *gin.Context) {
	// 获取标签ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的标签ID", err.Error())
		return
	}

	// 获取标签
	tag, err := h.tagService.GetTagByID(c.Request.Context(), uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "标签不存在")
		return
	}

	// 返回标签
	utils.SuccessResponse(c, "获取标签成功", tag)
}

// GetTagBySlug 根据Slug获取标签
func (h *TagHandler) GetTagBySlug(c *gin.Context) {
	// 获取标签Slug
	slug := c.Param("slug")

	// 获取标签
	tag, err := h.tagService.GetTagBySlug(c.Request.Context(), slug)
	if err != nil {
		utils.NotFoundResponse(c, "标签不存在")
		return
	}

	// 返回标签
	utils.SuccessResponse(c, "获取标签成功", tag)
} 