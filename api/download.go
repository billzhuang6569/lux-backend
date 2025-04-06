package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/iawia002/lux/service"
)

// DownloadRequest 表示下载请求
type DownloadRequest struct {
	URL     string `json:"url"`
	Format  string `json:"format"`
	Quality string `json:"quality"`
}

// DownloadHandler 处理视频下载请求
type DownloadHandler struct {
	taskManager *service.TaskManager
}

// NewDownloadHandler 创建新的下载处理程序
func NewDownloadHandler(taskManager *service.TaskManager) *DownloadHandler {
	return &DownloadHandler{
		taskManager: taskManager,
	}
}

// HandleDownload 处理下载请求
func (h *DownloadHandler) HandleDownload(w http.ResponseWriter, r *http.Request) {
	// 只允许POST方法
	if r.Method != http.MethodPost {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	// 解析请求正文
	var req DownloadRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		sendJSONError(w, "无效的请求格式", http.StatusBadRequest)
		return
	}

	// 验证必填字段
	if req.URL == "" {
		sendJSONError(w, "URL不能为空", http.StatusBadRequest)
		return
	}

	// 如果format未提供，使用默认值
	if req.Format == "" {
		req.Format = "mp4"
	}

	// 验证格式
	validFormats := map[string]bool{
		"mp4":  true,
		"flv":  true,
		"mp3":  true,
		"webm": true,
		"m4a":  true,
	}
	if !validFormats[req.Format] {
		sendJSONError(w, "不支持的格式", http.StatusBadRequest)
		return
	}

	// 验证质量
	if req.Quality == "" {
		req.Quality = "best"
	}
	validQualities := map[string]bool{
		"best":  true,
		"1080p": true,
		"720p":  true,
		"480p":  true,
		"360p":  true,
	}
	if !validQualities[req.Quality] {
		sendJSONError(w, "不支持的质量", http.StatusBadRequest)
		return
	}

	// 创建下载任务
	task, err := h.taskManager.CreateTask(req.URL, req.Format, req.Quality)
	if err != nil {
		sendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回任务信息
	sendJSON(w, map[string]interface{}{
		"taskId":   task.ID,
		"status":   task.Status,
		"message":  task.Message,
		"progress": task.Progress,
	}, http.StatusOK)

	// 设置超时以避免请求挂起
	r.Context().Done()

	// 如果任务已完成（可能是缓存结果），返回下载URL
	if task.Status == service.TaskStatusCompleted {
		// 不需要做任何事情，状态API将返回下载链接
	}
}

// 辅助函数，发送JSON响应
func sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// 辅助函数，发送错误响应
func sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	sendJSON(w, map[string]string{
		"status":  "error",
		"message": message,
	}, statusCode)
}