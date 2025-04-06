package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/iawia002/lux/api"
	"github.com/iawia002/lux/config"
	"github.com/iawia002/lux/middleware"
	"github.com/iawia002/lux/service"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()

	// 创建下载服务
	baseURL := fmt.Sprintf("http://%s:%d", cfg.Server.Host, cfg.Server.Port)
	downloadService := service.NewDownloadService(
		cfg.Download.Directory,
		baseURL,
		cfg.Download.MaxConcurrent,
	)

	// 创建控制器
	downloadController := api.NewDownloadController(downloadService)
	sitesController := api.NewSitesController(downloadService)

	// 设置路由
	mux := http.NewServeMux()
	
	// API路由
	mux.HandleFunc("/api/download", downloadController.HandleDownload)
	mux.HandleFunc("/api/status/", downloadController.HandleStatus)
	mux.HandleFunc("/api/supported-sites", sitesController.HandleSupportedSites)

	// 静态文件服务
	downloadsDir := http.StripPrefix("/downloads/", http.FileServer(http.Dir(cfg.Download.Directory)))
	mux.Handle("/downloads/", downloadsDir)

	// 创建中间件
	corsMiddleware := middleware.NewCORSMiddleware(cfg.Server.CORSOrigins)
	handler := corsMiddleware.Middleware(mux)

	// 启动文件清理协程
	go func() {
		for {
			downloadService.CleanupTasks(cfg.Download.FileExpiryTime)
			time.Sleep(1 * time.Hour)
		}
	}()

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("服务器正在监听 %s\n", addr)
	log.Printf("下载目录: %s\n", cfg.Download.Directory)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
		os.Exit(1)
	}
}
