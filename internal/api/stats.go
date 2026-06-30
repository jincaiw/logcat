package api

import (
	"net/http"

	"syslog-alert/internal/models"
	"syslog-alert/internal/repository"
)

// GetStats 返回系统统计信息。
func (ws *WebServer) GetStats(w http.ResponseWriter, r *http.Request) {
	stats := ws.GetSystemStats()
	JSONResponse(w, stats)
}

// GetFieldStats 按字段聚合统计。
func (ws *WebServer) GetFieldStats(w http.ResponseWriter, r *http.Request) {
	var req models.FieldStatsRequest
	if !DecodeJSON(w, r, &req) {
		return
	}
	result := repository.GetFieldStats(req)
	JSONResponse(w, result)
}

// GetAvailableStatsFields 获取指定筛选策略下可统计的字段列表。
func (ws *WebServer) GetAvailableStatsFields(w http.ResponseWriter, r *http.Request) {
	policyID, ok := ParseUintID(w, r.PathValue("policyId"))
	if !ok {
		return
	}
	fields := repository.GetAvailableStatsFields(policyID)
	if fields == nil {
		fields = []models.StatsField{}
	}
	JSONResponse(w, fields)
}
