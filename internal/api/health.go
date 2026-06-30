package api

import (
	"net/http"
	"time"

	"syslog-alert/pkg/constants"
)

// Health 返回公开健康检查信息，供 Docker/Kubernetes/负载均衡探活使用。
func (ws *WebServer) Health(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w, map[string]interface{}{
		"status":    "ok",
		"app":       constants.AppName,
		"version":   constants.AppVersion,
		"uptimeSec": int64(time.Since(ws.startTime).Seconds()),
	})
}
