package main

import (
	"log"

	"ai-assistant-backend/config"
	"ai-assistant-backend/models"
	"ai-assistant-backend/utils"
)

func main() {
	// 加载配置文件
	if err := config.LoadConfig("../config.yaml"); err != nil {
		log.Fatal("加载配置文件失败:", err)
	}

	// 初始化数据库
	if err := config.InitDB(); err != nil {
		log.Fatal("初始化数据库失败:", err)
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

	// 检查是否已存在默认用户
	var existingUser models.User
	if err := config.DB.Where("username = ?", "admin").First(&existingUser).Error; err == nil {
		log.Println("默认用户已存在，跳过创建")
		return
	}

	// 创建默认用户
	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		log.Fatal("密码加密失败:", err)
	}

	defaultUser := models.User{
		Username: "admin",
		Password: hashedPassword,
		Email:    "admin@example.com",
	}

	if err := config.DB.Create(&defaultUser).Error; err != nil {
		log.Fatal("创建默认用户失败:", err)
	}

	log.Println("数据库初始化完成")
	log.Println("默认用户:")
	log.Println("  用户名: admin")
	log.Println("  密码: admin123")
	log.Println("  邮箱: admin@example.com")
}
