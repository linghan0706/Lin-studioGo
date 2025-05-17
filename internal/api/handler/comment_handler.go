package handler

import (
	"Lin_studio/internal/service"
	"Lin_studio/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CommentHandler 评论处理器
type CommentHandler struct {
	commentService service.CommentService
}

// NewCommentHandler 创建评论处理器实例
func NewCommentHandler(commentService service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// GetComments 获取评论列表
func (h *CommentHandler) GetComments(c *gin.Context) {
	// 获取查询参数
	page, limit := utils.GetPagination(c)
	
	// 获取内容类型和ID
	itemType := c.Query("item_type")
	if itemType == "" {
		utils.BadRequestResponse(c, "缺少内容类型参数", nil)
		return
	}
	
	itemIDStr := c.Query("item_id")
	if itemIDStr == "" {
		utils.BadRequestResponse(c, "缺少内容ID参数", nil)
		return
	}
	
	itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的内容ID", err.Error())
		return
	}
	
	// 获取父评论ID（可选）
	var parentID *uint
	parentIDStr := c.Query("parent_id")
	if parentIDStr != "" {
		id, err := strconv.ParseUint(parentIDStr, 10, 32)
		if err == nil {
			pid := uint(id)
			parentID = &pid
		}
	}
	
	// 获取状态（可选，仅管理员可用）
	var status *string
	userRole, exists := c.Get("userRole")
	if exists && userRole.(string) == "admin" {
		statusStr := c.Query("status")
		if statusStr != "" {
			status = &statusStr
		}
	}
	
	// 查询评论
	comments, pagination, err := h.commentService.GetComments(
		c.Request.Context(),
		itemType,
		uint(itemID),
		parentID,
		page,
		limit,
		status,
	)
	
	if err != nil {
		utils.InternalServerErrorResponse(c, "获取评论列表失败: "+err.Error())
		return
	}
	
	// 转换为响应格式
	commentsResponse := make([]interface{}, len(comments))
	for i, comment := range comments {
		commentsResponse[i] = comment.ToResponse()
	}
	
	// 返回结果
	utils.SuccessResponse(c, "获取评论列表成功", gin.H{
		"comments":   commentsResponse,
		"pagination": pagination,
	})
}

// GetCommentByID 根据ID获取评论
func (h *CommentHandler) GetCommentByID(c *gin.Context) {
	// 获取评论ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的评论ID", err.Error())
		return
	}
	
	// 获取评论
	comment, err := h.commentService.GetCommentByID(c.Request.Context(), uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "评论不存在")
		return
	}
	
	// 返回结果
	utils.SuccessResponse(c, "获取评论成功", comment.ToResponse())
}

// CreateComment 创建评论
func (h *CommentHandler) CreateComment(c *gin.Context) {
	// 请求体结构
	var req struct {
		Content        string `json:"content" binding:"required"`
		ItemType       string `json:"item_type" binding:"required"`
		ItemID         uint   `json:"item_id" binding:"required"`
		ParentID       *uint  `json:"parent_id"`
		AnonymousAuthor *struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"anonymous_author"`
	}
	
	// 绑定请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "无效的请求数据", err.Error())
		return
	}
	
	// 获取用户ID（可选）
	var userID *uint
	userIDInterface, exists := c.Get("userID")
	if exists {
		uid := userIDInterface.(uint)
		userID = &uid
	}
	
	// 获取匿名信息
	var anonymousName, anonymousEmail string
	if req.AnonymousAuthor != nil {
		anonymousName = req.AnonymousAuthor.Name
		anonymousEmail = req.AnonymousAuthor.Email
	}
	
	// 创建评论
	comment, err := h.commentService.CreateComment(
		c.Request.Context(),
		req.Content,
		req.ItemType,
		req.ItemID,
		userID,
		anonymousName,
		anonymousEmail,
		req.ParentID,
	)
	
	if err != nil {
		utils.BadRequestResponse(c, "创建评论失败", err.Error())
		return
	}
	
	// 返回结果
	utils.CreatedResponse(c, "评论创建成功", comment.ToResponse())
}

// UpdateComment 更新评论
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	// 获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.UnauthorizedResponse(c, "未授权")
		return
	}
	
	// 获取评论ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的评论ID", err.Error())
		return
	}
	
	// 请求体结构
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	
	// 绑定请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "无效的请求数据", err.Error())
		return
	}
	
	// 更新评论
	comment, err := h.commentService.UpdateComment(
		c.Request.Context(),
		uint(id),
		req.Content,
		userID.(uint),
	)
	
	if err != nil {
		utils.BadRequestResponse(c, "更新评论失败", err.Error())
		return
	}
	
	// 返回结果
	utils.SuccessResponse(c, "评论更新成功", comment.ToResponse())
}

// DeleteComment 删除评论
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	// 获取用户ID和角色
	userID, exists := c.Get("userID")
	if !exists {
		utils.UnauthorizedResponse(c, "未授权")
		return
	}
	
	// 检查是否为管理员
	isAdmin := false
	userRole, exists := c.Get("userRole")
	if exists && userRole.(string) == "admin" {
		isAdmin = true
	}
	
	// 获取评论ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的评论ID", err.Error())
		return
	}
	
	// 删除评论
	err = h.commentService.DeleteComment(
		c.Request.Context(),
		uint(id),
		userID.(uint),
		isAdmin,
	)
	
	if err != nil {
		utils.BadRequestResponse(c, "删除评论失败", err.Error())
		return
	}
	
	// 返回结果
	utils.SuccessResponse(c, "评论删除成功", nil)
}

// LikeComment 点赞评论
func (h *CommentHandler) LikeComment(c *gin.Context) {
	// 获取评论ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的评论ID", err.Error())
		return
	}
	
	// 点赞评论
	err = h.commentService.LikeComment(c.Request.Context(), uint(id))
	if err != nil {
		utils.BadRequestResponse(c, "点赞评论失败", err.Error())
		return
	}
	
	// 返回结果
	utils.SuccessResponse(c, "点赞成功", nil)
}

// ApproveComment 批准评论
func (h *CommentHandler) ApproveComment(c *gin.Context) {
	// 检查权限
	userRole, exists := c.Get("userRole")
	if !exists || userRole.(string) != "admin" {
		utils.ForbiddenResponse(c, "无权操作")
		return
	}
	
	// 获取评论ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的评论ID", err.Error())
		return
	}
	
	// 批准评论
	err = h.commentService.ApproveComment(c.Request.Context(), uint(id))
	if err != nil {
		utils.BadRequestResponse(c, "批准评论失败", err.Error())
		return
	}
	
	// 返回结果
	utils.SuccessResponse(c, "评论已批准", nil)
}

// MarkCommentAsSpam 标记评论为垃圾信息
func (h *CommentHandler) MarkCommentAsSpam(c *gin.Context) {
	// 检查权限
	userRole, exists := c.Get("userRole")
	if !exists || userRole.(string) != "admin" {
		utils.ForbiddenResponse(c, "无权操作")
		return
	}
	
	// 获取评论ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "无效的评论ID", err.Error())
		return
	}
	
	// 标记为垃圾信息
	err = h.commentService.MarkCommentAsSpam(c.Request.Context(), uint(id))
	if err != nil {
		utils.BadRequestResponse(c, "标记评论失败", err.Error())
		return
	}
	
	// 返回结果
	utils.SuccessResponse(c, "评论已标记为垃圾信息", nil)
} 