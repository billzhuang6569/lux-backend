package service

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// FileManager 管理下载文件的存储和访问
type FileManager struct {
	downloadDir  string
	fileExpiryTime time.Duration
}

// NewFileManager 创建新的文件管理器
func NewFileManager(downloadDir string, fileExpiryTime time.Duration) *FileManager {
	return &FileManager{
		downloadDir:  downloadDir,
		fileExpiryTime: fileExpiryTime,
	}
}

// ServeFile 通过HTTP提供文件下载
func (fm *FileManager) ServeFile(w http.ResponseWriter, r *http.Request, filename string) {
	filePath := filepath.Join(fm.downloadDir, filename)
	
	// 检查文件是否存在
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		http.Error(w, "文件不存在", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "无法访问文件", http.StatusInternalServerError)
		return
	}
	
	// 检查文件是否过期
	if time.Since(info.ModTime()) > fm.fileExpiryTime {
		// 文件已过期，删除并返回404
		os.Remove(filePath)
		http.Error(w, "文件已过期", http.StatusNotFound)
		return
	}
	
	// 设置Content-Disposition头，使浏览器将响应作为文件下载处理
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	
	// 设置Content-Type
	contentType := getContentType(filename)
	w.Header().Set("Content-Type", contentType)
	
	// 提供文件下载
	http.ServeFile(w, r, filePath)
}

// SaveFile 保存从外部URL下载的文件
func (fm *FileManager) SaveFile(url, filename string) (string, error) {
	// 创建目标文件
	filePath := filepath.Join(fm.downloadDir, filename)
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("无法创建文件: %w", err)
	}
	defer out.Close()
	
	// 获取内容
	response, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("下载失败: %w", err)
	}
	defer response.Body.Close()
	
	// 检查响应状态
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载失败，服务器返回: %d", response.StatusCode)
	}
	
	// 将响应内容写入文件
	_, err = io.Copy(out, response.Body)
	if err != nil {
		return "", fmt.Errorf("写入文件失败: %w", err)
	}
	
	return filePath, nil
}

// CleanupExpiredFiles 清理过期文件
func (fm *FileManager) CleanupExpiredFiles() error {
	files, err := os.ReadDir(fm.downloadDir)
	if err != nil {
		return err
	}
	
	now := time.Now()
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue // 跳过无法获取信息的文件
		}
		
		// 如果文件已过期，删除它
		if now.Sub(info.ModTime()) > fm.fileExpiryTime {
			os.Remove(filepath.Join(fm.downloadDir, file.Name()))
		}
	}
	
	return nil
}

// 根据文件扩展名获取Content-Type
func getContentType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".mp4":
		return "video/mp4"
	case ".mp3":
		return "audio/mpeg"
	case ".m4a":
		return "audio/mp4"
	case ".webm":
		return "video/webm"
	case ".flv":
		return "video/x-flv"
	default:
		return "application/octet-stream"
	}
}