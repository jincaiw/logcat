package api

import (
	"net/http"
	"strconv"
	"time"

	"syslog-alert/internal/models"
	"syslog-alert/internal/repository"
	"syslog-alert/internal/service/cache"
	applogger "syslog-alert/pkg/logger"
)

// ListLogs 分页查询日志。
func (ws *WebServer) ListLogs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(query.Get("pageSize"))
	if pageSize <= 0 {
		pageSize = 50
	} else if pageSize > 500 {
		pageSize = 500
	}

	var deviceID *int
	if deviceIDStr := query.Get("deviceId"); deviceIDStr != "" {
		if id, err := strconv.Atoi(deviceIDStr); err == nil && id > 0 {
			deviceID = &id
		}
	}

	startTime := query.Get("startTime")
	endTime := query.Get("endTime")
	keyword := query.Get("keyword")

	logs, total := repository.GetLogs(page, pageSize, deviceID, startTime, endTime, keyword)

	JSONResponse(w, models.LogQueryResult{
		Logs:  logs,
		Total: total,
	})
}

// CleanupLogs 清理指定天数前的日志。
func (ws *WebServer) CleanupLogs(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Days int `json:"days"`
	}
	if !DecodeJSON(w, r, &body) {
		return
	}
	if body.Days <= 0 {
		JSONError(w, "days 必须大于 0", http.StatusBadRequest)
		return
	}
	if err := repository.CleanupOldLogs(body.Days); err != nil {
		applogger.Error("清理旧日志失败: %v", err)
		JSONError(w, "清理旧日志失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	cache.InvalidateStatsCaches()
	ws.SyslogServer().ClearOldTraces(time.Duration(body.Days) * 24 * time.Hour)
	JSONResponse(w, map[string]bool{"success": true})
}

// CleanupAllLogs 清空所有日志。
func (ws *WebServer) CleanupAllLogs(w http.ResponseWriter, r *http.Request) {
	if err := repository.CleanupAllLogs(); err != nil {
		applogger.Error("清空所有日志失败: %v", err)
		JSONError(w, "清空所有日志失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	cache.InvalidateStatsCaches()
	ws.SyslogServer().ClearAllTraces()
	JSONResponse(w, map[string]bool{"success": true})
}

// GetUnmatchedLogsCount 获取未匹配日志数量。
func (ws *WebServer) GetUnmatchedLogsCount(w http.ResponseWriter, r *http.Request) {
	count := repository.GetUnmatchedLogsCount()
	JSONResponse(w, map[string]int64{"count": count})
}

// CleanupUnmatchedLogs 清理指定天数前的未匹配日志。
func (ws *WebServer) CleanupUnmatchedLogs(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Days int `json:"days"`
	}
	if !DecodeJSON(w, r, &body) {
		return
	}
	if body.Days <= 0 {
		JSONError(w, "days 必须大于 0", http.StatusBadRequest)
		return
	}
	if err := repository.CleanupUnmatchedLogs(body.Days); err != nil {
		applogger.Error("清理未匹配日志失败: %v", err)
		JSONError(w, "清理未匹配日志失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	cache.InvalidateStatsCaches()
	ws.SyslogServer().ClearOldTraces(time.Duration(body.Days) * 24 * time.Hour)
	JSONResponse(w, map[string]bool{"success": true})
}
