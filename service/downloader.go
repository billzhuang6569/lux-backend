package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iawia002/lux/downloader"
	"github.com/iawia002/lux/extractors"
	"github.com/iawia002/lux/utils"
)

// ProgressCallback 是下载进度回调函数类型
type ProgressCallback func(percent int)

// 移除文件的辅助函数
func removeFile(path string) error {
	return os.Remove(path)
}

// processDownload 处理单个下载任务
func (tm *TaskManager) processDownload(task *Task) error {
	// 创建进度回调函数
	progressCallback := func(percent int) {
		tm.updateTaskStatus(task, TaskStatusProcessing, "正在下载中...", percent)
	}

	// 创建唯一的文件名
	fileName := fmt.Sprintf("%s_%s", task.ID, utils.FileName(filepath.Base(task.URL), "", 255))
	outputPath := tm.downloadDir

	// 设置下载选项
	options := downloader.Options{
		Silent:         true,
		OutputPath:     outputPath,
		OutputName:     fileName,
		FileNameLength: 255,
		RetryTimes:     10,
		ThreadNumber:   10,
		ChunkSizeMB:    1,
	}

	// 设置格式和质量
	if task.Format != "" {
		options.Stream = task.Format
	}
	if task.Quality == "best" {
		// 使用最佳质量
	} else if task.Quality != "" {
		// 例如如果 quality 是 "720p"，我们将在后面尝试选择匹配此分辨率的流
	}

	// 处理音频请求
	if task.Format == "mp3" || task.Format == "m4a" {
		options.AudioOnly = true
	}

	// 创建下载器实例
	luxDownloader := downloader.New(options)

	// 为下载器添加进度回调（需要修改 lux 包）
	luxDownloader.SetProgressCallback(progressCallback)

	// 提取视频信息
	extractOptions := extractors.Options{
		// 这里可以设置提取器选项，如cookie等
	}
	data, err := extractors.Extract(task.URL, extractOptions)
	if err != nil {
		return fmt.Errorf("无法提取视频信息: %w", err)
	}

	if len(data) == 0 || len(data[0].Streams) == 0 {
		return fmt.Errorf("没有找到可下载的流")
	}

	// 如果指定了质量但不是"best"，尝试找到匹配的质量
	// 注意：这里的实现可能需要根据 lux 实际的质量命名规则进行调整
	if task.Quality != "" && task.Quality != "best" {
		for _, d := range data {
			for streamName, stream := range d.Streams {
				if strings.Contains(strings.ToLower(stream.Quality), strings.ToLower(task.Quality)) {
					options.Stream = streamName
					break
				}
			}
		}
	}

	// 开始下载
	for _, item := range data {
		if item.Err != nil {
			return fmt.Errorf("提取数据错误: %w", item.Err)
		}

		err = luxDownloader.Download(item)
		if err != nil {
			return fmt.Errorf("下载错误: %w", err)
		}

		// 获取下载的文件路径
		ext := item.Streams[options.Stream].Ext
		if options.AudioOnly {
			// 确保音频格式正确
			if task.Format == "mp3" && ext != "mp3" {
				// 需要转换格式，但这需要额外实现
			}
		}

		filePath, err := utils.FilePath(fileName, ext, options.FileNameLength, options.OutputPath, false)
		if err != nil {
			return fmt.Errorf("获取文件路径错误: %w", err)
		}

		// 更新任务文件路径
		task.FilePath = filePath
	}

	// 更新任务状态为已完成
	tm.updateTaskStatus(task, TaskStatusCompleted, "下载已完成", 100)
	return nil
}

// 获取支持的网站列表
func GetSupportedSites() []string {
	// 这个列表应该从 lux 的注册表中获取，但为了简单，我们这里直接硬编码
	// 后续可以通过分析 extractors 目录自动生成
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
		"微博",
		"小红书",
		"TikTok",
		"Reddit",
		"PornHub",
		"AcFun",
		"Tumblr",
		"芒果TV",
		"爱拍原创",
		"喜马拉雅",
		"映客直播",
		"斗鱼",
		"虎牙",
		"快手",
	}
}