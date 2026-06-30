// Package api 提供 HTTP API 层，包含路由注册、请求处理与响应封装。
//
// 设计要点：
//   - 使用 Go 1.22 增强的 http.ServeMux（支持方法路由与路径参数），零外部依赖
//   - Handler 按实体分文件，降低单文件复杂度
//   - 统一 JSON 响应格式与错误处理
//   - 直接调用 repository/service 层，消除原 App 方法的薄封装
package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// JSONResponse 发送 JSON 响应，设置正确的 Content-Type。
func JSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// JSONError 发送 JSON 格式的错误响应。
func JSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// DecodeJSON 解码请求体 JSON 到目标结构。
// 解码失败时自动发送 400 错误响应，返回 false。
func DecodeJSON(w http.ResponseWriter, r *http.Request, target interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		JSONError(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return false
	}
	return true
}

// ParseUintID 从请求路径参数解析 uint ID。
// 解析失败时自动发送 400 错误响应，返回 0 和 false。
func ParseUintID(w http.ResponseWriter, idStr string) (uint, bool) {
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		JSONError(w, "invalid ID: "+idStr, http.StatusBadRequest)
		return 0, false
	}
	return uint(id64), true
}
