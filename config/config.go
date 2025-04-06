package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config 包含应用程序的配置
type Config struct {
	// 服务器配置
	Port            int
	Host            string
	CorsOrigins     []string
	
	// 下载配置
	DownloadDir           string
	MaxConcurrentDownloads int
	DownloadTimeout       time.Duration
	FileExpiryTime        time.Duration
	
	// 安全配置
	APIKey    string
	RateLimit int
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	config := &Config{
		Port:                  getEnvAsInt("PORT", 8080),
		Host:                  getEnv("HOST", "0.0.0.0"),
		CorsOrigins:           getEnvAsSlice("CORS_ORIGINS", []string{"*"}),
		DownloadDir:           getEnv("DOWNLOAD_DIR", "/tmp/lux-downloads"),
		MaxConcurrentDownloads: getEnvAsInt("MAX_CONCURRENT_DOWNLOADS", 5),
		DownloadTimeout:       time.Duration(getEnvAsInt("DOWNLOAD_TIMEOUT", 3600)) * time.Second,
		FileExpiryTime:        time.Duration(getEnvAsInt("FILE_EXPIRY_TIME", 86400)) * time.Second,
		APIKey:                getEnv("API_KEY", ""),
		RateLimit:             getEnvAsInt("RATE_LIMIT", 60),
	}
	
	// 确保下载目录存在
	if _, err := os.Stat(config.DownloadDir); os.IsNotExist(err) {
		err := os.MkdirAll(config.DownloadDir, 0755)
		if err != nil {
			panic("无法创建下载目录: " + err.Error())
		}
	}
	
	return config
}

// 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// 获取环境变量并转换为整数
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// 获取环境变量并转换为字符串切片
func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, ",")
}