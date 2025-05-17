package middleware

import (
	"Lin_studio/internal/utils"

	"github.com/gin-gonic/gin"
)

// RequireAdmin 要求管理员角色的中间件
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			utils.ForbiddenResponse(c, "Administrator privileges required")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireEditor 要求编辑者或更高角色的中间件
func RequireEditor() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || (role != "admin" && role != "editor") {
			utils.ForbiddenResponse(c, "Editor privileges required")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireOwnerOrAdmin 要求资源所有者或管理员的中间件
// 需要在路由处理函数中设置资源所有者ID到上下文中: c.Set("owner_id", resourceOwnerID)
func RequireOwnerOrAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, userExists := c.Get("user_id")
		role, roleExists := c.Get("role")
		ownerID, ownerExists := c.Get("owner_id")

		// 管理员可以访问任何资源
		if roleExists && role == "admin" {
			c.Next()
			return
		}

		// 检查用户是否为资源所有者
		if userExists && ownerExists && userID == ownerID {
			c.Next()
			return
		}

		utils.ForbiddenResponse(c, "You don't have permission to access this resource")
		c.Abort()
	}
} 