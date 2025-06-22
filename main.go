package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"ai-assistant-backend/config"
	"ai-assistant-backend/models"
	"ai-assistant-backend/routes"

	_ "net/http/pprof"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置文件
	if err := config.LoadConfig(""); err != nil {
		log.Fatal("加载配置文件失败:", err)
	}

	// 设置Gin模式
	gin.SetMode(config.GlobalConfig.Server.Mode)

	// 初始化数据库
	if err := config.InitDB(); err != nil {
		log.Fatal("初始化数据库失败:", err)
	}

	// 初始化Redis
	if err := config.InitRedis(); err != nil {
		log.Fatal("初始化Redis失败:", err)
	}

	// 自动迁移数据库表
	config.DB.AutoMigrate(
		&models.User{},
		&models.Agent{},
		&models.AgentCarouselImage{},
		&models.SelfService{},
		&models.DocumentCategory{},
		&models.Document{},
		&models.DocumentTag{},
		&models.Tag{},
		&models.FAQCategory{},
		&models.FAQ{},
	)

	// 创建Gin实例
	router := gin.Default()

	// 配置CORS
	corsConfig := cors.Config{
		AllowOrigins:     config.GlobalConfig.Server.CORS.AllowedOrigins,
		AllowHeaders:     config.GlobalConfig.Server.CORS.AllowedHeaders,
		AllowMethods:     config.GlobalConfig.Server.CORS.AllowedMethods,
		AllowCredentials: true,
	}
	router.Use(cors.New(corsConfig))

	// 设置路由
	routes.SetupAuthRoutes(router)
	routes.SetupAgentRoutes(router)
	routes.SetupDocumentRoutes(router)
	routes.SetupFAQRoutes(router)
	routes.SetupUploadRoutes(router)

	// 健康检查接口
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": config.GlobalConfig.App.Name + " is running",
			"version": config.GlobalConfig.App.Version,
		})
	})

	// 启动pprof性能分析端口
	go func() {
		pprofPort := config.GlobalConfig.Server.PprofPort
		if pprofPort == 0 {
			pprofPort = 6060
		}
		addr := fmt.Sprintf(":%d", pprofPort)
		log.Printf("[pprof] 性能分析端口监听在 %s", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Printf("[pprof] 性能分析端口启动失败: %v", err)
		}
	}()

	// 启动服务器
	port := config.GlobalConfig.Server.Port
	if port == 0 {
		port = 8080
	}

	// 优雅关闭
	go func() {
		log.Printf("服务器启动在端口 %d", port)
		if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
			log.Fatal("启动服务器失败:", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	// 关闭Redis连接
	if err := config.CloseRedis(); err != nil {
		log.Printf("关闭Redis连接失败: %v", err)
	}

	log.Println("服务器已关闭")
}
