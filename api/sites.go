package api

import (
	"net/http"

	"github.com/iawia002/lux/service"
)

// SitesHandler 处理支持的网站请求
type SitesHandler struct{}

// NewSitesHandler 创建新的网站处理程序
func NewSitesHandler() *SitesHandler {
	return &SitesHandler{}
}

// HandleSupportedSites 处理支持的网站请求
func (h *SitesHandler) HandleSupportedSites(w http.ResponseWriter, r *http.Request) {
	// 只允许GET方法
	if r.Method != http.MethodGet {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	// 获取支持的网站列表
	sites := service.GetSupportedSites()

	// 发送响应
	sendJSON(w, map[string]interface{}{
		"sites": sites,
	}, http.StatusOK)
}