package handler

import (
	"Lin_studio/internal/service"
	"Lin_studio/internal/utils"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile 获取当前用户信息
func (h *UserHandler) GetProfile(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.UnauthorizedResponse(c, "未授权")
		return
	}

	// 获取用户信息
	user, err := h.userService.GetProfile(c.Request.Context(), userID.(uint))
	if err != nil {
		utils.InternalServerErrorResponse(c, "获取用户信息失败: "+err.Error())
		return
	}

	// 返回用户信息
	utils.SuccessResponse(c, "获取用户信息成功", user)
}

// UpdateProfile 更新用户信息
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.UnauthorizedResponse(c, "未授权")
		return
	}

	// 请求体结构
	var req struct {
		Bio         string                 `json:"bio"`
		SocialLinks map[string]interface{} `json:"social_links"`
		ContactInfo map[string]interface{} `json:"contact_info"`
	}

	// 绑定请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "无效的请求数据", err.Error())
		return
	}

	// 更新用户信息
	user, err := h.userService.UpdateProfile(
		c.Request.Context(),
		userID.(uint),
		req.Bio,
		req.SocialLinks,
		req.ContactInfo,
	)
	if err != nil {
		utils.InternalServerErrorResponse(c, "更新用户信息失败: "+err.Error())
		return
	}

	// 返回更新后的用户信息
	utils.SuccessResponse(c, "用户信息更新成功", user)
}

// UploadAvatar 上传用户头像
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.UnauthorizedResponse(c, "未授权")
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		utils.BadRequestResponse(c, "上传文件失败", err.Error())
		return
	}

	// 上传头像
	avatarURL, err := h.userService.UploadAvatar(c.Request.Context(), userID.(uint), file)
	if err != nil {
		utils.InternalServerErrorResponse(c, "上传头像失败: "+err.Error())
		return
	}

	// 返回头像URL
	utils.SuccessResponse(c, "头像上传成功", gin.H{
		"avatar": avatarURL,
	})
} 