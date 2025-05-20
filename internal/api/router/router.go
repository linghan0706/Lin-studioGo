package router

import (
	"Lin_studio/internal/api/handler"
	"Lin_studio/internal/api/middleware"
	"Lin_studio/internal/config"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	categoryHandler *handler.CategoryHandler,
	tagHandler *handler.TagHandler,
	articleHandler *handler.ArticleHandler,
	commentHandler *handler.CommentHandler,
	toolHandler *handler.ToolHandler,
	// 其他处理器...
) *gin.Engine {
	r := gin.Default()

	// 允许跨域请求
	r.Use(corsMiddleware())

	// 根路径接口 - 直接通过IP和端口访问
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "api服务正常",
		})
	})

	// 健康检查端点
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"message": "Service is running",
		})
	})

	// API版本前缀
	api := r.Group("/api/v1")

	// 身份验证路由
	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		
		// 需要管理员权限的路由
		auth.POST("/register", middleware.JWTAuth(), middleware.RequireAdmin(), authHandler.Register)
		
		// 需要认证的路由
		auth.POST("/change-password", middleware.JWTAuth(), authHandler.ChangePassword)
		auth.POST("/logout", middleware.JWTAuth(), authHandler.Logout)
	}

	// 用户路由
	users := api.Group("/users")
	{
		// 需要认证的路由
		users.GET("/profile", middleware.JWTAuth(), userHandler.GetProfile)
		users.PUT("/profile", middleware.JWTAuth(), userHandler.UpdateProfile)
		users.POST("/avatar", middleware.JWTAuth(), userHandler.UploadAvatar)
	}

	// 分类路由
	categories := api.Group("/categories")
	{
		// 公开路由
		categories.GET("", categoryHandler.GetAllCategories)
		categories.GET("/:id", categoryHandler.GetCategoryByID)
		categories.GET("/slug/:slug", categoryHandler.GetCategoryBySlug)

		// 需要管理员权限的路由
		categories.POST("", middleware.JWTAuth(), middleware.RequireAdmin(), categoryHandler.CreateCategory)
	}

	// 标签路由
	tags := api.Group("/tags")
	{
		// 公开路由
		tags.GET("", tagHandler.GetAllTags)
		tags.GET("/:id", tagHandler.GetTagByID)
		tags.GET("/slug/:slug", tagHandler.GetTagBySlug)

		// 需要管理员权限的路由
		tags.POST("", middleware.JWTAuth(), middleware.RequireAdmin(), tagHandler.CreateTag)
	}

	// 文章路由
	articles := api.Group("/articles")
	{
		// 公开路由
		articles.GET("", articleHandler.GetArticles)
		articles.GET("/:slug", articleHandler.GetArticleBySlug)
		articles.GET("/featured", articleHandler.GetFeaturedArticles)
		articles.POST("/like/:id", articleHandler.LikeArticle)

		// 需要认证的路由
		articles.POST("", middleware.JWTAuth(), articleHandler.CreateArticle)
		articles.PUT("/:id", middleware.JWTAuth(), articleHandler.UpdateArticle)
		articles.DELETE("/:id", middleware.JWTAuth(), articleHandler.DeleteArticle)
		articles.POST("/upload-cover", middleware.JWTAuth(), articleHandler.UploadCoverImage)
	}

	// 评论路由
	comments := api.Group("/comments")
	{
		// 公开路由
		comments.GET("", commentHandler.GetComments)
		comments.GET("/:id", commentHandler.GetCommentByID)
		comments.POST("/like/:id", commentHandler.LikeComment)
		
		// 创建评论 - 可以是登录用户，也可以是匿名用户
		comments.POST("", commentHandler.CreateComment)
		
		// 需要认证的路由
		comments.PUT("/:id", middleware.JWTAuth(), commentHandler.UpdateComment)
		comments.DELETE("/:id", middleware.JWTAuth(), commentHandler.DeleteComment)
		
		// 需要管理员权限的路由
		comments.PUT("/:id/approve", middleware.JWTAuth(), middleware.RequireAdmin(), commentHandler.ApproveComment)
		comments.PUT("/:id/spam", middleware.JWTAuth(), middleware.RequireAdmin(), commentHandler.MarkCommentAsSpam)
	}

	// 工具路由
	tools := api.Group("/tools")
	{
		// 公开路由
		tools.GET("", toolHandler.GetTools)
		tools.GET("/categories", toolHandler.GetToolCategories)
		tools.GET("/:slug", toolHandler.GetToolBySlug)
		
		// 需要管理员权限的路由
		tools.POST("", middleware.JWTAuth(), middleware.RequireAdmin(), toolHandler.CreateTool)
		tools.PUT("/:id", middleware.JWTAuth(), middleware.RequireAdmin(), toolHandler.UpdateTool)
		tools.DELETE("/:id", middleware.JWTAuth(), middleware.RequireAdmin(), toolHandler.DeleteTool)
	}

	return r
}

// corsMiddleware 跨域中间件
func corsMiddleware() gin.HandlerFunc {
	// 获取CORS配置
	corsConfig := config.GetConfig().CORS

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowOrigin := "*"
		
		// 检查请求的Origin是否在允许列表中
		for _, allowed := range corsConfig.AllowedOrigins {
			if origin == allowed || allowed == "*" {
				// 如果在允许列表中，设置为具体的Origin，而不是*
				// 这样才能支持带认证的请求
				allowOrigin = origin
				break
			}
		}
		
		// 设置CORS响应头
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		
		// 当允许特定来源而不是通配符*时，才能设置允许凭证
		if allowOrigin != "*" && corsConfig.AllowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}
		
		// 设置允许的请求方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowedMethods, ", "))
		
		// 设置允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowedHeaders, ", "))
		
		// 设置暴露的响应头
		c.Writer.Header().Set("Access-Control-Expose-Headers", strings.Join(corsConfig.ExposedHeaders, ", "))
		
		// 设置预检请求的缓存时间
		c.Writer.Header().Set("Access-Control-Max-Age", strconv.Itoa(corsConfig.MaxAge))

		// 对于OPTIONS预检请求，直接返回成功并结束请求处理
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
} 