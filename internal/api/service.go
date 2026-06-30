package api

import (
	"net/http"

	applogger "syslog-alert/pkg/logger"
)

// serviceStartRequest 服务启动请求体。
type serviceStartRequest struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

// GetServiceStatus 返回 Syslog 服务运行状态。
func (ws *WebServer) GetServiceStatus(w http.ResponseWriter, r *http.Request) {
	server := ws.SyslogServer()
	JSONResponse(w, map[string]interface{}{
		"serviceRunning": server.IsRunning(),
		"listenPort":     server.GetPort(),
		"receiveCount":   server.GetReceiveCount(),
		"receiveRate":    server.GetReceiveRate(),
		"connections":    server.GetConnections(),
	})
}

// StartService 启动 Syslog 服务。
func (ws *WebServer) StartService(w http.ResponseWriter, r *http.Request) {
	var req serviceStartRequest
	if !DecodeJSON(w, r, &req) {
		return
	}
	if err := ws.SyslogServer().Start(req.Port, req.Protocol); err != nil {
		applogger.Error("启动 Syslog 服务失败: %v", err)
		JSONError(w, "启动服务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, map[string]bool{"success": true})
}

// StopService 停止 Syslog 服务。
func (ws *WebServer) StopService(w http.ResponseWriter, r *http.Request) {
	if err := ws.SyslogServer().Stop(); err != nil {
		applogger.Error("停止 Syslog 服务失败: %v", err)
		JSONError(w, "停止服务失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, map[string]bool{"success": true})
}
