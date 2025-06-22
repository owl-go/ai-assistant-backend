package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	if GlobalConfig == nil {
		return fmt.Errorf("配置未加载，请先调用 LoadConfig")
	}

	var err error

	// 使用新的GORM连接数据库
	DB, err = gorm.Open(mysql.Open(GlobalConfig.Database.GetDSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 启用SQL日志
	})

	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 获取底层的sql.DB对象来设置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(GlobalConfig.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(GlobalConfig.Database.MaxOpenConns)

	log.Println("数据库连接成功")
	return nil
}
