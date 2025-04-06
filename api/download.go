package api

import (
	"encoding/json"
	"net/http"

	"github.com/iawia002/lux/service"
)

// DownloadRequest 下载请求
type DownloadRequest struct {
	URL     string `json:"url"`
	Format  string `json:"format"`
	Quality string `json:"quality"`
}

// DownloadController 下载控制器
type DownloadController struct {
	downloadService *service.DownloadService
}

// NewDownloadController 创建下载控制器
func NewDownloadController(downloadService *service.DownloadService) *DownloadController {
	return &DownloadController{
		downloadService: downloadService,
	}
}

// HandleDownload 处理下载请求
func (c *DownloadController) HandleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持 POST 方法", http.StatusMethodNotAllowed)
		return
	}

	var req DownloadRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "无效的请求数据: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// 验证请求
	if req.URL == "" {
		http.Error(w, "URL不能为空", http.StatusBadRequest)
		return
	}

	// 设置默认格式和质量
	if req.Format == "" {
		req.Format = "mp4"
	}
	if req.Quality == "" {
		req.Quality = "best"
	}

	// 创建下载任务
	task := c.downloadService.CreateTask(req.URL, req.Format, req.Quality)

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// HandleStatus 处理获取状态请求
func (c *DownloadController) HandleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "仅支持 GET 方法", http.StatusMethodNotAllowed)
		return
	}

	// 从URL路径中提取任务ID
	taskID := r.URL.Path[len("/api/status/"):]
	if taskID == "" {
		http.Error(w, "缺少任务ID", http.StatusBadRequest)
		return
	}

	// 获取任务
	task, found := c.downloadService.GetTask(taskID)
	if !found {
		http.Error(w, "任务不存在", http.StatusNotFound)
		return
	}

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
} 