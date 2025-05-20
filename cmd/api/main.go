package main

import (
	"Lin_studio/internal/api/handler"
	"Lin_studio/internal/api/router"
	"Lin_studio/internal/config"
	"Lin_studio/internal/repository"
	"Lin_studio/internal/service"
	"fmt"
	"log"
	"os"
)

func main() {
	// 设置日志输出
	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.Println("应用程序启动...")

	// 加载配置
	cfg := config.GetConfig()
	log.Println("配置已加载")

	// 初始化数据库
	config.InitDB()
	db := config.DB
	log.Println("数据库已初始化")

	// 初始化仓库
	log.Println("初始化仓库...")
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository()
	tagRepo := repository.NewTagRepository()
	articleRepo := repository.NewArticleRepository()
	commentRepo := repository.NewCommentRepository()
	toolRepo := repository.NewToolRepository()
	log.Println("仓库初始化完成")

	// 初始化服务
	log.Println("初始化服务...")
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	tagService := service.NewTagService(tagRepo)
	articleService := service.NewArticleService(articleRepo, tagRepo, userRepo, categoryRepo)
	commentService := service.NewCommentService(commentRepo, userRepo)
	toolService := service.NewToolService(toolRepo)
	log.Println("服务初始化完成")

	// 初始化处理器
	log.Println("初始化处理器...")
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	categoryHandler := handler.NewCategoryHandler(categoryService)
	tagHandler := handler.NewTagHandler(tagService)
	articleHandler := handler.NewArticleHandler(articleService)
	commentHandler := handler.NewCommentHandler(commentService)
	toolHandler := handler.NewToolHandler(toolService)
	log.Println("处理器初始化完成")

	// 设置路由
	log.Println("设置路由...")
	r := router.SetupRouter(
		authHandler,
		userHandler,
		categoryHandler,
		tagHandler,
		articleHandler,
		commentHandler,
		toolHandler,
	)
	log.Println("路由设置完成")

	// 启动服务器
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("服务器即将在 %s 启动", serverAddr)
	
	// 使用defer和recover捕获可能的panic
	defer func() {
		if r := recover(); r != nil {
			log.Printf("程序发生panic: %v", r)
			os.Exit(1)
		}
	}()
	
	if err := r.Run(serverAddr); err != nil {
		log.Printf("服务器启动失败: %v", err)
		os.Exit(1)
	}
}
