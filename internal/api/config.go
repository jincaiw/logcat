package api

import (
	"net/http"

	"syslog-alert/internal/models"
	"syslog-alert/internal/repository"
	applogger "syslog-alert/pkg/logger"
)

// GetConfig 获取系统配置。
func (ws *WebServer) GetConfig(w http.ResponseWriter, r *http.Request) {
	cfg := repository.GetSystemConfig()
	JSONResponse(w, cfg)
}

// SaveConfig 保存系统配置。
func (ws *WebServer) SaveConfig(w http.ResponseWriter, r *http.Request) {
	var cfg models.SystemConfig
	if !DecodeJSON(w, r, &cfg) {
		return
	}
	if err := repository.UpdateSystemConfig(cfg); err != nil {
		applogger.Error("保存系统配置失败: %v", err)
		JSONError(w, "保存系统配置失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, cfg)
}
