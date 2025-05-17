package handler

import (
	"Lin_studio/internal/service"
	"Lin_studio/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CategoryHandler 分类处理器
type CategoryHandler struct {
	categoryService service.CategoryService
}

// NewCategoryHandler 创建分类处理器实例
func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

// GetAllCategories 获取所有分类
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	// 获取父分类ID参数
	var parentID *uint
	parentIDStr := c.Query("parent_id")
	if parentIDStr != "" {
		id, err := strconv.ParseUint(parentIDStr, 10, 32)
		if err != nil {
			utils.BadRequestResponse(c, "无效的父分类ID", err.Error())
			return
		}
		parentIDUint := uint(id)
		parentID = &parentIDUint
	}

	// 获取分类列表
	categories, err := h.categoryService.GetAllCategories(c.Request.Context(), parentID)
	if err != nil {
		utils.InternalServerErrorResponse(c, "获取分类列表失败: "+err.Error())
		return
	}

	// 返回分类列表
	utils.SuccessResponse(c, "获取分类列表成功", gin.H{
		"categories": categories,
	})
}

// CreateCategory 创建分类
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	// 请求体结构
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		ParentID    *uint  `json:"parent_id"`
	}

	// 绑定请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "无效的请求数据", err.Error())
		return
	}

	// 创建分类
	category, err := h.categoryService.CreateCategory(
		c.Request.Context(),
		req.Name,
		req.Description,
		req.ParentID,
	)
	if err != nil {
		utils.BadRequestResponse(c, "创建分类失败", err.Error())
		return
	}

	// 返回创建的分类
	utils.CreatedResponse(c, "分类创建成功", gin.H{
		"id":   category.ID,
		"name": category.Name,
		"slug": category.Slug,
	})
}

// GetCategoryByID 根据ID获取分类
func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	// 获取分类ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的分类ID", err.Error())
		return
	}

	// 获取分类
	category, err := h.categoryService.GetCategoryByID(c.Request.Context(), uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "分类不存在")
		return
	}

	// 返回分类
	utils.SuccessResponse(c, "获取分类成功", category)
}

// GetCategoryBySlug 根据Slug获取分类
func (h *CategoryHandler) GetCategoryBySlug(c *gin.Context) {
	// 获取分类Slug
	slug := c.Param("slug")

	// 获取分类
	category, err := h.categoryService.GetCategoryBySlug(c.Request.Context(), slug)
	if err != nil {
		utils.NotFoundResponse(c, "分类不存在")
		return
	}

	// 返回分类
	utils.SuccessResponse(c, "获取分类成功", category)
} 