package main

import (
	"Lin_studio/internal/api/handler"
	"Lin_studio/internal/api/router"
	"Lin_studio/internal/config"
	"Lin_studio/internal/repository"
	"Lin_studio/internal/service"
	"fmt"
	"log"
)

func main() {
	// 加载配置
	cfg := config.GetConfig()

	// 初始化数据库
	config.InitDB()
	db := config.DB

	// 初始化仓库
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository()
	tagRepo := repository.NewTagRepository()
	articleRepo := repository.NewArticleRepository()
	commentRepo := repository.NewCommentRepository()
	toolRepo := repository.NewToolRepository()

	// 初始化服务
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	tagService := service.NewTagService(tagRepo)
	articleService := service.NewArticleService(articleRepo, tagRepo, userRepo, categoryRepo)
	commentService := service.NewCommentService(commentRepo, userRepo)
	toolService := service.NewToolService(toolRepo)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	tagHandler := handler.NewTagHandler(tagService)
	articleHandler := handler.NewArticleHandler(articleService)
	commentHandler := handler.NewCommentHandler(commentService)
	toolHandler := handler.NewToolHandler(toolService)

	// 设置路由
	r := router.SetupRouter(
		authHandler,
		userHandler,
		categoryHandler,
		tagHandler,
		articleHandler,
		commentHandler,
		toolHandler,
	)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
