package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/iawia002/lux/api"
	"github.com/iawia002/lux/config"
	"github.com/iawia002/lux/middleware"
	"github.com/iawia002/lux/service"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 创建服务地址
	serverAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	// 创建任务管理器
	taskManager := service.NewTaskManager(
		cfg.MaxConcurrentDownloads,
		cfg.DownloadDir,
		serverAddr,
	)

	// 创建文件管理器
	fileManager := service.NewFileManager(
		cfg.DownloadDir,
		cfg.FileExpiryTime,
	)

	// 创建API处理程序
	downloadHandler := api.NewDownloadHandler(taskManager)
	statusHandler := api.NewStatusHandler(taskManager)
	sitesHandler := api.NewSitesHandler()

	// 创建中间件
	corsMiddleware := middleware.NewCorsMiddleware(cfg.CorsOrigins)
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit)

	// 创建路由
	mux := http.NewServeMux()

	// 注册API路由
	mux.HandleFunc("/api/download", downloadHandler.HandleDownload)
	mux.HandleFunc("/api/status/", statusHandler.HandleStatus)
	mux.HandleFunc("/api/supported-sites", sitesHandler.HandleSupportedSites)

	// 注册文件下载路由
	mux.HandleFunc("/downloads/", func(w http.ResponseWriter, r *http.Request) {
		filename := filepath.Base(r.URL.Path)
		fileManager.ServeFile(w, r, filename)
	})

	// 应用中间件
	var handler http.Handler = mux
	handler = corsMiddleware.Middleware(handler)
	handler = rateLimiter.Middleware(handler)

	// 启动HTTP服务器
	log.Printf("启动服务器 %s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, handler))
}