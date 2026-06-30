package api

import (
	"net/http"

	"syslog-alert/internal/models"
	"syslog-alert/internal/repository"
	applogger "syslog-alert/pkg/logger"
)

// ---- 筛选策略 ----

// ListFilterPolicies 列出全部筛选策略。
func (ws *WebServer) ListFilterPolicies(w http.ResponseWriter, r *http.Request) {
	policies := repository.GetFilterPolicies()
	JSONResponse(w, policies)
}

// CreateFilterPolicy 创建筛选策略。
func (ws *WebServer) CreateFilterPolicy(w http.ResponseWriter, r *http.Request) {
	var policy models.FilterPolicy
	if !DecodeJSON(w, r, &policy) {
		return
	}
	if err := repository.CreateFilterPolicy(&policy); err != nil {
		applogger.Error("创建筛选策略失败: %v", err)
		JSONError(w, "创建筛选策略失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, policy)
}

// GetFilterPolicy 获取单个筛选策略。
func (ws *WebServer) GetFilterPolicy(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	policy, err := repository.GetFilterPolicyByID(id)
	if err != nil {
		applogger.Error("获取筛选策略失败: %v", err)
		JSONError(w, "筛选策略不存在", http.StatusNotFound)
		return
	}
	JSONResponse(w, policy)
}

// UpdateFilterPolicy 更新筛选策略。
func (ws *WebServer) UpdateFilterPolicy(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	var policy models.FilterPolicy
	if !DecodeJSON(w, r, &policy) {
		return
	}
	policy.ID = id
	if err := repository.UpdateFilterPolicy(&policy); err != nil {
		applogger.Error("更新筛选策略失败: %v", err)
		JSONError(w, "更新筛选策略失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, policy)
}

// DeleteFilterPolicy 删除筛选策略。
func (ws *WebServer) DeleteFilterPolicy(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	if err := repository.DeleteFilterPolicy(id); err != nil {
		applogger.Error("删除筛选策略失败: %v", err)
		JSONError(w, "删除筛选策略失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, map[string]bool{"success": true})
}

// ---- 告警策略 ----

// ListAlertPolicies 列出全部告警策略。
func (ws *WebServer) ListAlertPolicies(w http.ResponseWriter, r *http.Request) {
	policies := repository.GetAlertPolicies()
	JSONResponse(w, policies)
}

// CreateAlertPolicy 创建告警策略。
func (ws *WebServer) CreateAlertPolicy(w http.ResponseWriter, r *http.Request) {
	var policy models.AlertPolicy
	if !DecodeJSON(w, r, &policy) {
		return
	}
	if err := repository.CreateAlertPolicy(&policy); err != nil {
		applogger.Error("创建告警策略失败: %v", err)
		JSONError(w, "创建告警策略失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, policy)
}

// GetAlertPolicy 获取单个告警策略。
func (ws *WebServer) GetAlertPolicy(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	policy, err := repository.GetAlertPolicyByID(id)
	if err != nil {
		applogger.Error("获取告警策略失败: %v", err)
		JSONError(w, "告警策略不存在", http.StatusNotFound)
		return
	}
	JSONResponse(w, policy)
}

// UpdateAlertPolicy 更新告警策略。
func (ws *WebServer) UpdateAlertPolicy(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	var policy models.AlertPolicy
	if !DecodeJSON(w, r, &policy) {
		return
	}
	policy.ID = id
	if err := repository.UpdateAlertPolicy(&policy); err != nil {
		applogger.Error("更新告警策略失败: %v", err)
		JSONError(w, "更新告警策略失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, policy)
}

// DeleteAlertPolicy 删除告警策略。
func (ws *WebServer) DeleteAlertPolicy(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	if err := repository.DeleteAlertPolicy(id); err != nil {
		applogger.Error("删除告警策略失败: %v", err)
		JSONError(w, "删除告警策略失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, map[string]bool{"success": true})
}
