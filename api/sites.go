package api

import (
	"encoding/json"
	"net/http"

	"github.com/iawia002/lux/service"
)

// SitesController 网站控制器
type SitesController struct {
	downloadService *service.DownloadService
}

// NewSitesController 创建网站控制器
func NewSitesController(downloadService *service.DownloadService) *SitesController {
	return &SitesController{
		downloadService: downloadService,
	}
}

// HandleSupportedSites 处理获取支持的网站请求
func (c *SitesController) HandleSupportedSites(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "仅支持 GET 方法", http.StatusMethodNotAllowed)
		return
	}

	// 获取支持的网站
	sites := c.downloadService.GetSupportedSites()

	// 返回响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"sites": sites,
	})
} 