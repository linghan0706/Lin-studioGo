package config

import (
	"os"
	"strings"
	"time"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Upload   UploadConfig
	CORS     CORSConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret           string
	ExpirationHours  int
	RefreshSecret    string
	RefreshExpHours  int
	Issuer           string
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	UploadDir    string
	MaxSize      int64
	AllowedTypes []string
}

// CORSConfig CORS配置
type CORSConfig struct {
	AllowedOrigins     []string // 允许的域名列表
	AllowCredentials   bool     // 是否允许带凭证的请求
	AllowedMethods     []string // 允许的HTTP方法
	AllowedHeaders     []string // 允许的HTTP头
	ExposedHeaders     []string // 暴露给客户端的HTTP头
	MaxAge             int      // 预检请求缓存时间(秒)
}

// GetConfig 获取配置
func GetConfig() Config {
	return Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  time.Second * 60,
			WriteTimeout: time.Second * 60,
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "101.126.146.84"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "mysql_123"),
			DBName:   getEnv("DB_NAME", "Nuxt_admin"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "your-secret-key"),
			ExpirationHours: 24,
			RefreshSecret:   getEnv("JWT_REFRESH_SECRET", "your-refresh-secret-key"),
			RefreshExpHours: 168, // 7 days
			Issuer:          "lin-studio",
		},
		Upload: UploadConfig{
			UploadDir:    getEnv("UPLOAD_DIR", "./uploads"),
			MaxSize:      10 * 1024 * 1024, // 10MB
			AllowedTypes: []string{"image/jpeg", "image/png", "image/gif"},
		},
		CORS: CORSConfig{
			// 默认允许的域名列表，可以通过CORS_ALLOWED_ORIGINS环境变量覆盖
			// 格式为逗号分隔的域名列表，例如：http://localhost:3000,https://example.com
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{
				"http://localhost:3000",
				"http://localhost:8080",
			}),
			AllowCredentials: true,
			AllowedMethods: []string{
				"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH",
			},
			AllowedHeaders: []string{
				"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token",
				"Authorization", "accept", "origin", "Cache-Control", "X-Requested-With",
				"Token", "Refresh-Token",
			},
			ExposedHeaders: []string{
				"Content-Length", "Authorization", "Token", "Refresh-Token",
			},
			MaxAge: 86400, // 24小时
		},
	}
}

// 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// 获取环境变量并转换为字符串切片，以逗号分隔
func getEnvAsSlice(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
} 