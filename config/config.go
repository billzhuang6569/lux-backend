package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config 应用程序配置
type Config struct {
	Server   ServerConfig
	Download DownloadConfig
	Security SecurityConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port        int
	Host        string
	CORSOrigins []string
}

// DownloadConfig 下载配置
type DownloadConfig struct {
	Directory       string
	MaxConcurrent   int
	Timeout         time.Duration
	FileExpiryTime  time.Duration
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	APIKey    string
	RateLimit int
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	cfg := &Config{
		Server: ServerConfig{
			Port:        getEnvAsInt("PORT", 8080),
			Host:        getEnv("HOST", "0.0.0.0"),
			CORSOrigins: getEnvAsSlice("CORS_ORIGINS", []string{"*"}),
		},
		Download: DownloadConfig{
			Directory:      getEnv("DOWNLOAD_DIR", "/tmp/lux-downloads"),
			MaxConcurrent:  getEnvAsInt("MAX_CONCURRENT_DOWNLOADS", 5),
			Timeout:        time.Duration(getEnvAsInt("DOWNLOAD_TIMEOUT", 3600)) * time.Second,
			FileExpiryTime: time.Duration(getEnvAsInt("FILE_EXPIRY_TIME", 86400)) * time.Second,
		},
		Security: SecurityConfig{
			APIKey:    getEnv("API_KEY", ""),
			RateLimit: getEnvAsInt("RATE_LIMIT", 60),
		},
	}

	// 确保下载目录存在
	os.MkdirAll(cfg.Download.Directory, 0755)

	return cfg
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

// FakeHeaders fake http headers
var FakeHeaders = map[string]string{
	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	"Accept-Charset":  "UTF-8,*;q=0.5",
	"Accept-Encoding": "gzip,deflate,sdch",
	"Accept-Language": "en-US,en;q=0.8",
	"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.81 Safari/537.36",
}
