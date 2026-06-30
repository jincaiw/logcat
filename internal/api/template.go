package api

import (
	"net/http"

	"syslog-alert/internal/models"
	"syslog-alert/internal/repository"
	applogger "syslog-alert/pkg/logger"
)

// ---- 解析模板 ----

// ListParseTemplates 列出全部解析模板。
func (ws *WebServer) ListParseTemplates(w http.ResponseWriter, r *http.Request) {
	templates := repository.GetParseTemplates()
	JSONResponse(w, templates)
}

// CreateParseTemplate 创建解析模板。
func (ws *WebServer) CreateParseTemplate(w http.ResponseWriter, r *http.Request) {
	var template models.ParseTemplate
	if !DecodeJSON(w, r, &template) {
		return
	}
	if err := repository.CreateParseTemplate(&template); err != nil {
		applogger.Error("创建解析模板失败: %v", err)
		JSONError(w, "创建解析模板失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, template)
}

// GetParseTemplate 获取单个解析模板。
func (ws *WebServer) GetParseTemplate(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	template, err := repository.GetParseTemplateByID(id)
	if err != nil {
		applogger.Error("获取解析模板失败: %v", err)
		JSONError(w, "解析模板不存在", http.StatusNotFound)
		return
	}
	JSONResponse(w, template)
}

// UpdateParseTemplate 更新解析模板。
func (ws *WebServer) UpdateParseTemplate(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	var template models.ParseTemplate
	if !DecodeJSON(w, r, &template) {
		return
	}
	template.ID = id
	if err := repository.UpdateParseTemplate(&template); err != nil {
		applogger.Error("更新解析模板失败: %v", err)
		JSONError(w, "更新解析模板失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, template)
}

// DeleteParseTemplate 删除解析模板。
func (ws *WebServer) DeleteParseTemplate(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	if err := repository.DeleteParseTemplate(id); err != nil {
		applogger.Error("删除解析模板失败: %v", err)
		JSONError(w, "删除解析模板失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, map[string]bool{"success": true})
}

// ---- 输出模板 ----

// ListOutputTemplates 列出全部输出模板。
func (ws *WebServer) ListOutputTemplates(w http.ResponseWriter, r *http.Request) {
	templates := repository.GetOutputTemplates()
	JSONResponse(w, templates)
}

// CreateOutputTemplate 创建输出模板。
func (ws *WebServer) CreateOutputTemplate(w http.ResponseWriter, r *http.Request) {
	var template models.OutputTemplate
	if !DecodeJSON(w, r, &template) {
		return
	}
	if platform, ok := normalizeNotificationPlatform(template.Platform); ok {
		template.Platform = platform
	} else {
		JSONError(w, "不支持的通知渠道", http.StatusBadRequest)
		return
	}
	if err := repository.CreateOutputTemplate(&template); err != nil {
		applogger.Error("创建输出模板失败: %v", err)
		JSONError(w, "创建输出模板失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, template)
}

// GetOutputTemplate 获取单个输出模板。
func (ws *WebServer) GetOutputTemplate(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	template, err := repository.GetOutputTemplateByID(id)
	if err != nil {
		applogger.Error("获取输出模板失败: %v", err)
		JSONError(w, "输出模板不存在", http.StatusNotFound)
		return
	}
	JSONResponse(w, template)
}

// UpdateOutputTemplate 更新输出模板。
func (ws *WebServer) UpdateOutputTemplate(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	var template models.OutputTemplate
	if !DecodeJSON(w, r, &template) {
		return
	}
	if platform, ok := normalizeNotificationPlatform(template.Platform); ok {
		template.Platform = platform
	} else {
		JSONError(w, "不支持的通知渠道", http.StatusBadRequest)
		return
	}
	template.ID = id
	if err := repository.UpdateOutputTemplate(&template); err != nil {
		applogger.Error("更新输出模板失败: %v", err)
		JSONError(w, "更新输出模板失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, template)
}

// DeleteOutputTemplate 删除输出模板。
func (ws *WebServer) DeleteOutputTemplate(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	if err := repository.DeleteOutputTemplate(id); err != nil {
		applogger.Error("删除输出模板失败: %v", err)
		JSONError(w, "删除输出模板失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, map[string]bool{"success": true})
}

// ---- 字段映射文档 ----

// ListFieldMappingDocs 列出全部字段映射文档。
func (ws *WebServer) ListFieldMappingDocs(w http.ResponseWriter, r *http.Request) {
	docs := repository.GetFieldMappingDocs()
	JSONResponse(w, docs)
}

// CreateFieldMappingDoc 创建字段映射文档。
func (ws *WebServer) CreateFieldMappingDoc(w http.ResponseWriter, r *http.Request) {
	var doc models.FieldMappingDoc
	if !DecodeJSON(w, r, &doc) {
		return
	}
	if err := repository.CreateFieldMappingDoc(&doc); err != nil {
		applogger.Error("创建字段映射文档失败: %v", err)
		JSONError(w, "创建字段映射文档失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, doc)
}

// GetFieldMappingDoc 获取单个字段映射文档。
func (ws *WebServer) GetFieldMappingDoc(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	doc, err := repository.GetFieldMappingDocByID(id)
	if err != nil {
		applogger.Error("获取字段映射文档失败: %v", err)
		JSONError(w, "字段映射文档不存在", http.StatusNotFound)
		return
	}
	JSONResponse(w, doc)
}

// UpdateFieldMappingDoc 更新字段映射文档。
func (ws *WebServer) UpdateFieldMappingDoc(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	var doc models.FieldMappingDoc
	if !DecodeJSON(w, r, &doc) {
		return
	}
	doc.ID = id
	if err := repository.UpdateFieldMappingDoc(&doc); err != nil {
		applogger.Error("更新字段映射文档失败: %v", err)
		JSONError(w, "更新字段映射文档失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, doc)
}

// DeleteFieldMappingDoc 删除字段映射文档。
func (ws *WebServer) DeleteFieldMappingDoc(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	if err := repository.DeleteFieldMappingDoc(id); err != nil {
		applogger.Error("删除字段映射文档失败: %v", err)
		JSONError(w, "删除字段映射文档失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, map[string]bool{"success": true})
}
