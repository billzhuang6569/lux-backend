package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// DebugController 调试控制器
type DebugController struct {}

// NewDebugController 创建调试控制器
func NewDebugController() *DebugController {
	return &DebugController{}
}

// HandleDebug 处理调试请求
func (c *DebugController) HandleDebug(w http.ResponseWriter, r *http.Request) {
	result := make(map[string]interface{})
	
	// 收集环境信息
	result["env"] = collectEnvironment()
	
	// 检查程序
	result["commands"] = checkCommands()
	
	// 检查工作目录
	result["directories"] = checkDirectories()
	
	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// 收集环境信息
func collectEnvironment() map[string]string {
	env := make(map[string]string)
	
	env["os"] = runtime.GOOS
	env["arch"] = runtime.GOARCH
	env["goversion"] = runtime.Version()
	env["path"] = os.Getenv("PATH")
	env["pwd"] = getPwd()
	env["whoami"] = getWhoami()
	
	return env
}

// 检查命令
func checkCommands() map[string]interface{} {
	commands := make(map[string]interface{})
	
	// 检查lux
	commands["lux"] = checkCommand("lux", "--help")
	commands["lux_bin"] = checkCommandPath("/bin/lux", "--help")
	commands["lux_usr_bin"] = checkCommandPath("/usr/local/bin/lux", "--help")
	
	// 检查其他命令
	commands["ffmpeg"] = checkCommand("ffmpeg", "-version")
	commands["go"] = checkCommand("go", "version")
	
	return commands
}

// 检查目录
func checkDirectories() map[string]interface{} {
	dirs := make(map[string]interface{})
	
	// 检查当前目录
	dirs["pwd"] = listDir(".")
	
	// 检查bin目录
	dirs["bin"] = listDir("/bin")
	dirs["usr_local_bin"] = listDir("/usr/local/bin")
	
	// 检查下载目录
	dirs["download_dir"] = listDir(os.Getenv("DOWNLOAD_DIR"))
	
	return dirs
}

// 获取当前目录
func getPwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return dir
}

// 获取当前用户
func getWhoami() string {
	cmd := exec.Command("whoami")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return strings.TrimSpace(out.String())
}

// 检查命令
func checkCommand(name string, args ...string) map[string]string {
	result := make(map[string]string)
	
	path, err := exec.LookPath(name)
	if err != nil {
		result["exists"] = "false"
		result["error"] = err.Error()
		return result
	}
	
	result["exists"] = "true"
	result["path"] = path
	
	cmd := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	
	if err != nil {
		result["runnable"] = "false"
		result["run_error"] = err.Error()
	} else {
		result["runnable"] = "true"
	}
	
	result["stdout"] = stdout.String()
	result["stderr"] = stderr.String()
	
	return result
}

// 检查特定路径的命令
func checkCommandPath(path string, args ...string) map[string]string {
	result := make(map[string]string)
	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result["exists"] = "false"
		result["error"] = "file not found"
		return result
	}
	
	result["exists"] = "true"
	result["path"] = path
	
	cmd := exec.Command(path, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	
	if err != nil {
		result["runnable"] = "false"
		result["run_error"] = err.Error()
	} else {
		result["runnable"] = "true"
	}
	
	result["stdout"] = stdout.String()
	result["stderr"] = stderr.String()
	
	return result
}

// 列出目录内容
func listDir(dir string) []map[string]string {
	var files []map[string]string
	
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return []map[string]string{
			{"error": "directory not found"},
		}
	}
	
	entries, err := os.ReadDir(dir)
	if err != nil {
		return []map[string]string{
			{"error": err.Error()},
		}
	}
	
	for _, entry := range entries {
		file := make(map[string]string)
		file["name"] = entry.Name()
		
		info, err := entry.Info()
		if err != nil {
			file["type"] = "unknown"
		} else {
			if info.IsDir() {
				file["type"] = "directory"
			} else {
				file["type"] = "file"
				file["size"] = fmt.Sprintf("%d", info.Size())
				file["mode"] = info.Mode().String()
			}
		}
		
		// 如果是可执行文件，获取文件类型
		if dir == "/bin" || dir == "/usr/local/bin" {
			filePath := filepath.Join(dir, entry.Name())
			file["file_type"] = getFileType(filePath)
		}
		
		files = append(files, file)
	}
	
	return files
}

// 获取文件类型
func getFileType(path string) string {
	cmd := exec.Command("file", path)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return strings.TrimSpace(out.String())
} 