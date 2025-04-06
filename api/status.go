package api

import (
	"net/http"
	"strings"

	"github.com/iawia002/lux/service"
)

// StatusHandler 处理任务状态请求
type StatusHandler struct {
	taskManager *service.TaskManager
}

// NewStatusHandler 创建新的状态处理程序
func NewStatusHandler(taskManager *service.TaskManager) *StatusHandler {
	return &StatusHandler{
		taskManager: taskManager,
	}
}

// HandleStatus 处理状态请求
func (h *StatusHandler) HandleStatus(w http.ResponseWriter, r *http.Request) {
	// 只允许GET方法
	if r.Method != http.MethodGet {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	// 从URL中提取任务ID
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		sendJSONError(w, "无效的请求路径", http.StatusBadRequest)
		return
	}
	taskID := parts[len(parts)-1]

	// 获取任务状态
	task, err := h.taskManager.GetTask(taskID)
	if err != nil {
		sendJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	// 构建响应
	response := map[string]interface{}{
		"taskId":   task.ID,
		"status":   task.Status,
		"progress": task.Progress,
		"message":  task.Message,
	}

	// 如果任务已完成，添加下载URL和文件名
	if task.Status == service.TaskStatusCompleted {
		response["downloadUrl"] = task.DownloadURL
		response["filename"] = task.Filename
	}

	// 发送响应
	sendJSON(w, response, http.StatusOK)
}