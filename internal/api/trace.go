package api

import (
	"net/http"

	applogger "syslog-alert/pkg/logger"
)

// GetLogTrace 获取日志的全链路追踪信息。
func (ws *WebServer) GetLogTrace(w http.ResponseWriter, r *http.Request) {
	logID, ok := ParseUintID(w, r.PathValue("logId"))
	if !ok {
		return
	}
	trace := ws.SyslogServer().GetTraceInfo(logID)
	if trace == nil {
		applogger.Warn("未找到日志追踪信息: logID=%d", logID)
		JSONError(w, "未找到追踪信息", http.StatusNotFound)
		return
	}
	JSONResponse(w, trace)
}
