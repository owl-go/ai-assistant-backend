package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 配置结构体
type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
	Server   ServerConfig   `yaml:"server"`
	Upload   UploadConfig   `yaml:"upload"`
	MinIO    MinIOConfig    `yaml:"minio"`
	App      AppConfig      `yaml:"app"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	Name         string `yaml:"name"`
	Charset      string `yaml:"charset"`
	ParseTime    bool   `yaml:"parse_time"`
	Loc          string `yaml:"loc"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port      int        `yaml:"port"`
	Mode      string     `yaml:"mode"`
	PprofPort int        `yaml:"pprof_port"`
	CORS      CORSConfig `yaml:"cors"`
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedHeaders []string `yaml:"allowed_headers"`
	AllowedMethods []string `yaml:"allowed_methods"`
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	Path         string   `yaml:"path"`
	MaxFileSize  int64    `yaml:"max_file_size"`
	AllowedTypes []string `yaml:"allowed_types"`
}

// MinIOConfig MinIO配置
type MinIOConfig struct {
	Endpoint       string `yaml:"endpoint"`
	AccessKey      string `yaml:"access_key"`
	SecretKey      string `yaml:"secret_key"`
	Bucket         string `yaml:"bucket"`
	Location       string `yaml:"location"`
	UseSSL         bool   `yaml:"use_ssl"`
	ForcePathStyle bool   `yaml:"force_path_style"`
	URLExpireHours int    `yaml:"url_expire_hours"` // 链接有效时间（小时）
}

// AppConfig 应用配置
type AppConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

// GlobalConfig 全局配置实例
var GlobalConfig *Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) error {
	// 如果未指定配置文件路径，使用默认路径
	if configPath == "" {
		configPath = "config.yaml"
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析YAML
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 设置全局配置
	GlobalConfig = config

	log.Println("配置文件加载成功")
	return nil
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.Charset, c.ParseTime, c.Loc)
}

// GetRedisAddr 获取Redis地址
func (c *RedisConfig) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
