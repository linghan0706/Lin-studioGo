package handler

import (
	"Lin_studio/internal/domain"
	"Lin_studio/internal/service"
	"Lin_studio/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ToolHandler 工具处理器
type ToolHandler struct {
	toolService service.ToolService
}

// NewToolHandler 创建工具处理器实例
func NewToolHandler(toolService service.ToolService) *ToolHandler {
	return &ToolHandler{
		toolService: toolService,
	}
}

// GetTools 获取工具列表
func (h *ToolHandler) GetTools(c *gin.Context) {
	// 获取查询参数
	page, limit := utils.GetPagination(c)
	
	// 获取过滤参数
	var category, status, search *string
	
	// 分类
	categoryStr := c.Query("category")
	if categoryStr != "" {
		category = &categoryStr
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
	
	// 查询工具
	tools, pagination, err := h.toolService.GetTools(
		c.Request.Context(),
		page,
		limit,
		category,
		status,
		search,
	)
	
	if err != nil {
		utils.InternalServerErrorResponse(c, "获取工具列表失败: "+err.Error())
		return
	}
	
	// 转换为响应格式
	toolsResponse := make([]interface{}, len(tools))
	for i, tool := range tools {
		toolsResponse[i] = tool.ToSimpleResponse()
	}
	
	// 返回结果
	utils.SuccessResponse(c, "获取工具列表成功", gin.H{
		"tools":      toolsResponse,
		"pagination": pagination,
	})
}

// GetToolBySlug 根据Slug获取工具
func (h *ToolHandler) GetToolBySlug(c *gin.Context) {
	// 获取slug
	slug := c.Param("slug")
	
	// 获取工具
	tool, err := h.toolService.GetToolBySlug(c.Request.Context(), slug)
	if err != nil {
		utils.NotFoundResponse(c, "工具不存在")
		return
	}
	
	// 返回工具
	utils.SuccessResponse(c, "获取工具成功", tool.ToResponse(true))
}

// GetToolCategories 获取工具分类列表
func (h *ToolHandler) GetToolCategories(c *gin.Context) {
	// 由于工具分类是直接存储在工具表中的，我们需要使用聚合查询
	// 这里可以通过原生SQL或者其他方式实现
	// 为简化处理，先返回一些常用分类
	categories := []string{
		"开发工具",
		"设计工具",
		"效率工具",
		"文本处理",
		"图像处理",
		"网络工具",
		"其他",
	}
	
	utils.SuccessResponse(c, "获取工具分类成功", categories)
}

// CreateTool 创建工具
func (h *ToolHandler) CreateTool(c *gin.Context) {
	// 检查权限
	userRole, exists := c.Get("userRole")
	if !exists || userRole.(string) != "admin" {
		utils.ForbiddenResponse(c, "无权操作")
		return
	}
	
	// 请求体结构
	var req struct {
		Name        string            `json:"name" binding:"required"`
		Description string            `json:"description"`
		Icon        string            `json:"icon"`
		Category    string            `json:"category"`
		Content     string            `json:"content"`
		Config      domain.JSONConfig `json:"config"`
		Status      string            `json:"status"`
	}
	
	// 绑定请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "无效的请求数据", err.Error())
		return
	}
	
	// 创建工具
	tool, err := h.toolService.CreateTool(
		c.Request.Context(),
		req.Name,
		req.Description,
		req.Icon,
		req.Category,
		req.Content,
		req.Config,
		req.Status,
	)
	
	if err != nil {
		utils.BadRequestResponse(c, "创建工具失败", err.Error())
		return
	}
	
	// 返回结果
	utils.CreatedResponse(c, "工具创建成功", gin.H{
		"id":   tool.ID,
		"name": tool.Name,
		"slug": tool.Slug,
	})
}

// UpdateTool 更新工具
func (h *ToolHandler) UpdateTool(c *gin.Context) {
	// 检查权限
	userRole, exists := c.Get("userRole")
	if !exists || userRole.(string) != "admin" {
		utils.ForbiddenResponse(c, "无权操作")
		return
	}
	
	// 获取工具ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的工具ID", err.Error())
		return
	}
	
	// 请求体结构
	var req struct {
		Name        string            `json:"name" binding:"required"`
		Description string            `json:"description"`
		Icon        string            `json:"icon"`
		Category    string            `json:"category"`
		Content     string            `json:"content"`
		Config      domain.JSONConfig `json:"config"`
		Status      string            `json:"status"`
	}
	
	// 绑定请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "无效的请求数据", err.Error())
		return
	}
	
	// 更新工具
	tool, err := h.toolService.UpdateTool(
		c.Request.Context(),
		uint(id),
		req.Name,
		req.Description,
		req.Icon,
		req.Category,
		req.Content,
		req.Config,
		req.Status,
	)
	
	if err != nil {
		utils.BadRequestResponse(c, "更新工具失败", err.Error())
		return
	}
	
	// 返回结果
	utils.SuccessResponse(c, "工具更新成功", gin.H{
		"id":   tool.ID,
		"name": tool.Name,
		"slug": tool.Slug,
	})
}

// DeleteTool 删除工具
func (h *ToolHandler) DeleteTool(c *gin.Context) {
	// 检查权限
	userRole, exists := c.Get("userRole")
	if !exists || userRole.(string) != "admin" {
		utils.ForbiddenResponse(c, "无权操作")
		return
	}
	
	// 获取工具ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的工具ID", err.Error())
		return
	}
	
	// 删除工具
	err = h.toolService.DeleteTool(c.Request.Context(), uint(id))
	if err != nil {
		utils.BadRequestResponse(c, "删除工具失败", err.Error())
		return
	}
	
	// 返回结果
	utils.SuccessResponse(c, "工具删除成功", nil)
} 