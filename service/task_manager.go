package service

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

// TaskStatus 表示下载任务的状态
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"   // 任务等待中
	TaskStatusProcessing TaskStatus = "processing" // 任务处理中
	TaskStatusCompleted TaskStatus = "completed"  // 任务已完成
	TaskStatusError     TaskStatus = "error"      // 任务出错
)

// Task 表示一个下载任务
type Task struct {
	ID          string      `json:"taskId"`      // 任务ID
	URL         string      `json:"url"`         // 下载URL
	Format      string      `json:"format"`      // 下载格式
	Quality     string      `json:"quality"`     // 视频质量
	Status      TaskStatus  `json:"status"`      // 任务状态
	Progress    int         `json:"progress"`    // 下载进度 (0-100)
	Message     string      `json:"message"`     // 状态消息
	FilePath    string      `json:"-"`           // 下载文件的本地路径
	DownloadURL string      `json:"downloadUrl"` // 下载URL
	Filename    string      `json:"filename"`    // 文件名
	CreatedAt   time.Time   `json:"-"`           // 创建时间
	UpdatedAt   time.Time   `json:"-"`           // 更新时间
	Error       error       `json:"-"`           // 错误信息
}

// TaskManager 管理下载任务
type TaskManager struct {
	tasks       map[string]*Task
	queue       chan *Task
	maxWorkers  int
	downloadDir string
	serverAddr  string
	mu          sync.RWMutex
}

// NewTaskManager 创建新的任务管理器
func NewTaskManager(maxWorkers int, downloadDir string, serverAddr string) *TaskManager {
	tm := &TaskManager{
		tasks:       make(map[string]*Task),
		queue:       make(chan *Task, 100), // 队列大小限制为100个任务
		maxWorkers:  maxWorkers,
		downloadDir: downloadDir,
		serverAddr:  serverAddr,
	}
	
	// 启动工作线程池
	for i := 0; i < maxWorkers; i++ {
		go tm.worker()
	}
	
	// 启动过期文件清理
	go tm.cleanupRoutine()
	
	return tm
}

// CreateTask 创建新的下载任务
func (tm *TaskManager) CreateTask(url, format, quality string) (*Task, error) {
	taskID := uuid.New().String()
	
	task := &Task{
		ID:        taskID,
		URL:       url,
		Format:    format,
		Quality:   quality,
		Status:    TaskStatusPending,
		Progress:  0,
		Message:   "下载任务已创建",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	tm.mu.Lock()
	tm.tasks[taskID] = task
	tm.mu.Unlock()
	
	// 将任务放入队列
	select {
	case tm.queue <- task:
		// 成功将任务放入队列
	default:
		// 队列已满
		return nil, fmt.Errorf("任务队列已满，请稍后再试")
	}
	
	return task, nil
}

// GetTask 获取任务状态
func (tm *TaskManager) GetTask(taskID string) (*Task, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	
	task, exists := tm.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("任务不存在")
	}
	
	return task, nil
}

// worker 是处理下载任务的工作线程
func (tm *TaskManager) worker() {
	for task := range tm.queue {
		// 更新任务状态为处理中
		tm.updateTaskStatus(task, TaskStatusProcessing, "正在下载中...", task.Progress)
		
		// 开始下载
		err := tm.processDownload(task)
		if err != nil {
			tm.updateTaskStatus(task, TaskStatusError, fmt.Sprintf("下载失败: %v", err), task.Progress)
		}
	}
}

// updateTaskStatus 更新任务状态
func (tm *TaskManager) updateTaskStatus(task *Task, status TaskStatus, message string, progress int) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	
	task.Status = status
	task.Message = message
	task.Progress = progress
	task.UpdatedAt = time.Now()
	
	// 如果任务已完成，设置下载链接
	if status == TaskStatusCompleted && task.FilePath != "" {
		filename := filepath.Base(task.FilePath)
		task.Filename = filename
		task.DownloadURL = fmt.Sprintf("http://%s/downloads/%s", tm.serverAddr, filename)
	}
}

// cleanupRoutine 定期清理过期文件
func (tm *TaskManager) cleanupRoutine() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	
	for range ticker.C {
		tm.cleanupExpiredFiles()
	}
}

// cleanupExpiredFiles 清理过期的下载文件
func (tm *TaskManager) cleanupExpiredFiles() {
	// 获取所有已完成的任务
	tm.mu.RLock()
	completedTasks := make([]*Task, 0)
	for _, task := range tm.tasks {
		if task.Status == TaskStatusCompleted {
			completedTasks = append(completedTasks, task)
		}
	}
	tm.mu.RUnlock()
	
	now := time.Now()
	
	// 删除过期任务和文件
	for _, task := range completedTasks {
		if now.Sub(task.UpdatedAt) > 24*time.Hour {
			// 删除文件
			if task.FilePath != "" {
				removeFile(task.FilePath)
			}
			
			// 从任务列表中删除
			tm.mu.Lock()
			delete(tm.tasks, task.ID)
			tm.mu.Unlock()
		}
	}
}