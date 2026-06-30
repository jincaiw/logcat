package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"syslog-alert/internal/models"
	"syslog-alert/internal/repository"
	applogger "syslog-alert/pkg/logger"
)

// ExportParseTemplates 导出解析模板。ids 为空时导出全部。
func (ws *WebServer) ExportParseTemplates(w http.ResponseWriter, r *http.Request) {
	ids := parseIDListFromQuery(r, "ids")

	var templates []models.ParseTemplate
	if len(ids) > 0 {
		all := repository.GetParseTemplates()
		idSet := make(map[uint]bool, len(ids))
		for _, id := range ids {
			idSet[id] = true
		}
		for _, t := range all {
			if idSet[t.ID] {
				templates = append(templates, t)
			}
		}
	} else {
		templates = repository.GetParseTemplates()
	}

	export := models.ConfigExport{
		Version:        "1.0",
		ExportedAt:     time.Now().Format(time.RFC3339),
		Name:           "解析模板导出",
		ParseTemplates: templates,
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		applogger.Error("导出解析模板序列化失败: %v", err)
		JSONError(w, "导出失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// ExportFilterPolicies 导出筛选策略。ids 为空时导出全部。
func (ws *WebServer) ExportFilterPolicies(w http.ResponseWriter, r *http.Request) {
	ids := parseIDListFromQuery(r, "ids")

	var policies []models.FilterPolicy
	if len(ids) > 0 {
		all := repository.GetFilterPolicies()
		idSet := make(map[uint]bool, len(ids))
		for _, id := range ids {
			idSet[id] = true
		}
		for _, p := range all {
			if idSet[p.ID] {
				policies = append(policies, p)
			}
		}
	} else {
		policies = repository.GetFilterPolicies()
	}

	export := models.ConfigExport{
		Version:        "1.0",
		ExportedAt:     time.Now().Format(time.RFC3339),
		Name:           "筛选策略导出",
		FilterPolicies: policies,
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		applogger.Error("导出筛选策略序列化失败: %v", err)
		JSONError(w, "导出失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// ImportParseTemplates 导入解析模板。
func (ws *WebServer) ImportParseTemplates(w http.ResponseWriter, r *http.Request) {
	var export models.ConfigExport
	if !DecodeJSON(w, r, &export) {
		return
	}

	result := models.ImportResult{
		Success: true,
		Errors:  []string{},
	}

	templates := export.ParseTemplates
	if len(templates) == 0 {
		result.Success = false
		result.Message = "未找到解析模板数据"
		JSONResponse(w, result)
		return
	}

	for i := range templates {
		template := templates[i]
		template.ID = 0
		template.CreatedAt = time.Now()
		template.UpdatedAt = time.Now()

		if err := repository.CreateParseTemplate(&template); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("创建模板 %s 失败: %v", template.Name, err))
		} else {
			result.Count++
		}
	}

	result.Message = fmt.Sprintf("成功导入 %d 个解析模板", result.Count)
	if len(result.Errors) > 0 {
		result.Message += fmt.Sprintf("，%d 个失败", len(result.Errors))
	}

	JSONResponse(w, result)
}

// ImportFilterPolicies 导入筛选策略。
func (ws *WebServer) ImportFilterPolicies(w http.ResponseWriter, r *http.Request) {
	var export models.ConfigExport
	if !DecodeJSON(w, r, &export) {
		return
	}

	result := models.ImportResult{
		Success: true,
		Errors:  []string{},
	}

	policies := export.FilterPolicies
	if len(policies) == 0 {
		result.Success = false
		result.Message = "未找到筛选策略数据"
		JSONResponse(w, result)
		return
	}

	for i := range policies {
		policy := policies[i]
		policy.ID = 0
		policy.CreatedAt = time.Now()
		policy.UpdatedAt = time.Now()

		if err := repository.CreateFilterPolicy(&policy); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("创建策略 %s 失败: %v", policy.Name, err))
		} else {
			result.Count++
		}
	}

	result.Message = fmt.Sprintf("成功导入 %d 个筛选策略", result.Count)
	if len(result.Errors) > 0 {
		result.Message += fmt.Sprintf("，%d 个失败", len(result.Errors))
	}

	JSONResponse(w, result)
}

// GetPresetTemplates 返回预设模板列表（暂返回空数组）。
func (ws *WebServer) GetPresetTemplates(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w, []interface{}{})
}

// parseIDListFromQuery 从 URL query 中解析逗号分隔的 uint ID 列表。
func parseIDListFromQuery(r *http.Request, key string) []uint {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	var ids []uint
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id64, err := strconv.ParseUint(p, 10, 32)
		if err != nil {
			continue
		}
		ids = append(ids, uint(id64))
	}
	return ids
}
