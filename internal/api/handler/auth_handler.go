package handler

import (
	"Lin_studio/internal/domain"
	"Lin_studio/internal/service"
	"Lin_studio/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err.Error())
		return
	}

	userResponse, token, refreshToken, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		utils.UnauthorizedResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, "登录成功", gin.H{
		"token":         token,
		"refresh_token": refreshToken,
		"user":          userResponse,
	})
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role" binding:"required,oneof=admin editor user"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err.Error())
		return
	}

	// 创建用户对象
	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     req.Role,
	}

	// 设置密码
	if err := user.SetPassword(req.Password); err != nil {
		utils.InternalServerErrorResponse(c, "Failed to set password")
		return
	}

	// 注册用户
	userResponse, err := h.authService.Register(user)
	if err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.CreatedResponse(c, "用户创建成功", gin.H{
		"user": userResponse,
	})
}

// RefreshToken 刷新令牌
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err.Error())
		return
	}

	token, refreshToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.UnauthorizedResponse(c, err.Error())
		return
	}

	utils.SuccessResponse(c, "令牌刷新成功", gin.H{
		"token":         token,
		"refresh_token": refreshToken,
	})
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err.Error())
		return
	}

	// 从上下文中获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "用户未认证")
		return
	}

	if err := h.authService.ChangePassword(userID.(uint), req.CurrentPassword, req.NewPassword); err != nil {
		utils.BadRequestResponse(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, "密码修改成功", nil)
}

// Logout 退出登录
func (h *AuthHandler) Logout(c *gin.Context) {
	// 客户端只需要清除本地存储的令牌
	// 服务端无需处理
	utils.SuccessResponse(c, "登出成功", nil)
} 