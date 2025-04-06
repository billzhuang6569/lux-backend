package service

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/iawia002/lux/downloader"
	"github.com/iawia002/lux/extractors"
)

// ProgressCallback 进度回调函数类型
type ProgressCallback func(progress float64)

// DownloadStatus 下载状态
type DownloadStatus string

const (
	// StatusPending 等待下载
	StatusPending DownloadStatus = "pending"
	// StatusProcessing 下载中
	StatusProcessing DownloadStatus = "processing"
	// StatusCompleted 下载完成
	StatusCompleted DownloadStatus = "completed"
	// StatusError 下载错误
	StatusError DownloadStatus = "error"
)

// DownloadTask 下载任务信息
type DownloadTask struct {
	ID           string         `json:"taskId"`
	URL          string         `json:"url"`
	Format       string         `json:"format"`
	Quality      string         `json:"quality"`
	Status       DownloadStatus `json:"status"`
	Progress     float64        `json:"progress"`
	Message      string         `json:"message"`
	DownloadURL  string         `json:"downloadUrl,omitempty"`
	Filename     string         `json:"filename,omitempty"`
	Error        error          `json:"-"`
	CreatedAt    time.Time      `json:"-"`
	CompletedAt  time.Time      `json:"-"`
}

// DownloadService 下载服务
type DownloadService struct {
	tasks         map[string]*DownloadTask
	mutex         sync.RWMutex
	downloadDir   string
	baseURL       string
	maxConcurrent int
	taskChannel   chan *DownloadTask
	taskQueue     []*DownloadTask
	queueMutex    sync.Mutex
}

// NewDownloadService 创建下载服务
func NewDownloadService(downloadDir string, baseURL string, maxConcurrent int) *DownloadService {
	service := &DownloadService{
		tasks:         make(map[string]*DownloadTask),
		downloadDir:   downloadDir,
		baseURL:       baseURL,
		maxConcurrent: maxConcurrent,
		taskChannel:   make(chan *DownloadTask, maxConcurrent),
	}

	// 创建下载目录
	os.MkdirAll(downloadDir, 0755)

	// 启动工作线程
	for i := 0; i < maxConcurrent; i++ {
		go service.worker()
	}

	// 启动队列处理
	go service.queueProcessor()

	return service
}

// queueProcessor 处理队列中的任务
func (s *DownloadService) queueProcessor() {
	for {
		s.queueMutex.Lock()
		if len(s.taskQueue) > 0 && len(s.taskChannel) < s.maxConcurrent {
			task := s.taskQueue[0]
			s.taskQueue = s.taskQueue[1:]
			s.taskChannel <- task
		}
		s.queueMutex.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}

// worker 工作线程
func (s *DownloadService) worker() {
	for task := range s.taskChannel {
		s.processTask(task)
	}
}

// processTask 处理下载任务
func (s *DownloadService) processTask(task *DownloadTask) {
	s.updateTask(task.ID, func(t *DownloadTask) {
		t.Status = StatusProcessing
		t.Message = "正在下载中..."
	})

	// 创建任务目录
	taskDir := filepath.Join(s.downloadDir, task.ID)
	os.MkdirAll(taskDir, 0755)

	// 进度回调
	progressChan := make(chan float64, 10)
	done := make(chan bool)

	// 启动进度监控
	go func() {
		for {
			select {
			case progress := <-progressChan:
				s.updateTask(task.ID, func(t *DownloadTask) {
					t.Progress = progress
				})
			case <-done:
				return
			}
		}
	}()

	// 执行下载
	result, err := s.downloadVideo(task, taskDir, progressChan)
	close(done)

	if err != nil {
		s.updateTask(task.ID, func(t *DownloadTask) {
			t.Status = StatusError
			t.Message = fmt.Sprintf("下载失败: %v", err)
			t.Error = err
		})
		return
	}

	// 计算下载 URL
	downloadURL := fmt.Sprintf("%s/downloads/%s/%s", s.baseURL, task.ID, result.Filename)

	s.updateTask(task.ID, func(t *DownloadTask) {
		t.Status = StatusCompleted
		t.Progress = 100
		t.Message = "下载已完成"
		t.DownloadURL = downloadURL
		t.Filename = result.Filename
		t.CompletedAt = time.Now()
	})
}

// DownloadResult 下载结果
type DownloadResult struct {
	Filename string
	FilePath string
}

// 下载视频
func (s *DownloadService) downloadVideo(task *DownloadTask, taskDir string, progressChan chan float64) (*DownloadResult, error) {
	// 创建自定义下载器
	options := downloader.Options{
		OutputPath:     taskDir,
		ThreadNumber:   10,
		RetryTimes:     5,
		Stream:         task.Format,
		ChunkSizeMB:    10,
	}

	// 根据质量选择流
	switch task.Quality {
	case "best":
		// 使用最佳质量
	case "1080p":
		// 1080p 相关处理
	case "720p":
		// 720p 相关处理
	case "480p":
		// 480p 相关处理
	case "360p":
		// 360p 相关处理
	}

	// 提取视频信息
	extractOptions := extractors.Options{
		Playlist: false,
	}
	data, err := extractors.Extract(task.URL, extractOptions)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 || len(data[0].Streams) == 0 {
		return nil, fmt.Errorf("没有找到可下载的流")
	}

	// 设置文件名
	filename := data[0].Title

	// TODO: 修改 lux 下载器添加进度回调
	// 这里是模拟进度回调的示例，实际生产环境中需要修改 lux 源码
	go func() {
		for i := 0; i <= 100; i += 5 {
			progressChan <- float64(i)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// 执行下载
	d := downloader.New(options)
	err = d.Download(data[0])
	if err != nil {
		return nil, err
	}

	// 查找下载的文件
	var filePath string
	filepath.Walk(taskDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filePath = path
			filename = info.Name()
		}
		return nil
	})

	return &DownloadResult{
		Filename: filename,
		FilePath: filePath,
	}, nil
}

// CreateTask 创建下载任务
func (s *DownloadService) CreateTask(url, format, quality string) *DownloadTask {
	taskID := uuid.New().String()
	task := &DownloadTask{
		ID:        taskID,
		URL:       url,
		Format:    format,
		Quality:   quality,
		Status:    StatusPending,
		Progress:  0,
		Message:   "下载任务已创建",
		CreatedAt: time.Now(),
	}

	s.mutex.Lock()
	s.tasks[taskID] = task
	s.mutex.Unlock()

	// 添加到任务队列
	s.queueMutex.Lock()
	s.taskQueue = append(s.taskQueue, task)
	s.queueMutex.Unlock()

	return task
}

// GetTask 获取任务
func (s *DownloadService) GetTask(taskID string) (*DownloadTask, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	task, found := s.tasks[taskID]
	return task, found
}

// ListTasks 列出所有任务
func (s *DownloadService) ListTasks() []*DownloadTask {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	tasks := make([]*DownloadTask, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// updateTask 更新任务
func (s *DownloadService) updateTask(taskID string, updateFn func(*DownloadTask)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if task, found := s.tasks[taskID]; found {
		updateFn(task)
	}
}

// CleanupTasks 清理过期任务
func (s *DownloadService) CleanupTasks(expiryTime time.Duration) {
	s.mutex.Lock()
	var tasksToDelete []string
	now := time.Now()

	for id, task := range s.tasks {
		if task.Status == StatusCompleted || task.Status == StatusError {
			if task.CompletedAt.Add(expiryTime).Before(now) {
				tasksToDelete = append(tasksToDelete, id)
			}
		}
	}

	for _, id := range tasksToDelete {
		delete(s.tasks, id)
		// 删除文件
		taskDir := filepath.Join(s.downloadDir, id)
		os.RemoveAll(taskDir)
	}
	s.mutex.Unlock()
}

// GetSupportedSites 获取支持的网站
func (s *DownloadService) GetSupportedSites() []string {
	// 这里列出 lux 支持的网站
	// 可以通过分析 extractors 目录获取，这里为了简化，硬编码一些常见网站
	return []string{
		"YouTube",
		"Twitter",
		"Bilibili",
		"抖音",
		"优酷",
		"腾讯视频",
		"爱奇艺",
		"Instagram",
		"Facebook",
		"Vimeo",
		"TikTok",
		"微博",
		"小红书",
		"知乎",
		"AcFun",
		"Pornhub",
		"Reddit",
	}
} 