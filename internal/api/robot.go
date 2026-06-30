package api

import (
	"net/http"

	"syslog-alert/internal/models"
	"syslog-alert/internal/repository"
	"syslog-alert/internal/service/alert"
	applogger "syslog-alert/pkg/logger"
)

// ---- 机器人/推送通道 ----

// ListRobots 列出全部机器人。
func (ws *WebServer) ListRobots(w http.ResponseWriter, r *http.Request) {
	robots := repository.GetRobots()
	JSONResponse(w, robots)
}

// CreateRobot 创建机器人。
func (ws *WebServer) CreateRobot(w http.ResponseWriter, r *http.Request) {
	var robot models.Robot
	if !DecodeJSON(w, r, &robot) {
		return
	}
	if err := repository.CreateRobot(&robot); err != nil {
		applogger.Error("创建机器人失败: %v", err)
		JSONError(w, "创建机器人失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, robot)
}

// GetRobot 获取单个机器人。
func (ws *WebServer) GetRobot(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	robot, err := repository.GetRobotByID(id)
	if err != nil {
		applogger.Error("获取机器人失败: %v", err)
		JSONError(w, "机器人不存在", http.StatusNotFound)
		return
	}
	JSONResponse(w, robot)
}

// UpdateRobot 更新机器人。
func (ws *WebServer) UpdateRobot(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	var robot models.Robot
	if !DecodeJSON(w, r, &robot) {
		return
	}
	robot.ID = id
	if err := repository.UpdateRobot(&robot); err != nil {
		applogger.Error("更新机器人失败: %v", err)
		JSONError(w, "更新机器人失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, robot)
}

// DeleteRobot 删除机器人。
func (ws *WebServer) DeleteRobot(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	if err := repository.DeleteRobot(id); err != nil {
		applogger.Error("删除机器人失败: %v", err)
		JSONError(w, "删除机器人失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, map[string]bool{"success": true})
}

// TestRobot 测试机器人推送通道连通性。
func (ws *WebServer) TestRobot(w http.ResponseWriter, r *http.Request) {
	var robot models.Robot
	if !DecodeJSON(w, r, &robot) {
		return
	}
	result, err := alert.Test(&robot)
	if err != nil {
		applogger.Error("测试机器人失败: %v", err)
		JSONError(w, "测试失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, map[string]string{"message": result})
}

// ---- 告警规则 ----

// ListAlertRulesByRobot 列出指定机器人下的告警规则。
func (ws *WebServer) ListAlertRulesByRobot(w http.ResponseWriter, r *http.Request) {
	robotID, ok := ParseUintID(w, r.PathValue("robotId"))
	if !ok {
		return
	}
	rules := repository.GetAlertRulesByRobotID(robotID)
	JSONResponse(w, rules)
}

// CreateAlertRule 创建告警规则。
func (ws *WebServer) CreateAlertRule(w http.ResponseWriter, r *http.Request) {
	var rule models.AlertRule
	if !DecodeJSON(w, r, &rule) {
		return
	}
	if err := repository.CreateAlertRule(&rule); err != nil {
		applogger.Error("创建告警规则失败: %v", err)
		JSONError(w, "创建告警规则失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, rule)
}

// GetAlertRule 获取单个告警规则。
func (ws *WebServer) GetAlertRule(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	rule, err := repository.GetAlertRuleByID(id)
	if err != nil {
		applogger.Error("获取告警规则失败: %v", err)
		JSONError(w, "告警规则不存在", http.StatusNotFound)
		return
	}
	JSONResponse(w, rule)
}

// UpdateAlertRule 更新告警规则。
func (ws *WebServer) UpdateAlertRule(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	var rule models.AlertRule
	if !DecodeJSON(w, r, &rule) {
		return
	}
	rule.ID = id
	if err := repository.UpdateAlertRule(&rule); err != nil {
		applogger.Error("更新告警规则失败: %v", err)
		JSONError(w, "更新告警规则失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, rule)
}

// DeleteAlertRule 删除告警规则。
func (ws *WebServer) DeleteAlertRule(w http.ResponseWriter, r *http.Request) {
	id, ok := ParseUintID(w, r.PathValue("id"))
	if !ok {
		return
	}
	if err := repository.DeleteAlertRule(id); err != nil {
		applogger.Error("删除告警规则失败: %v", err)
		JSONError(w, "删除告警规则失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, map[string]bool{"success": true})
}

// DeleteAlertRulesByRobot 删除指定机器人下的全部告警规则。
func (ws *WebServer) DeleteAlertRulesByRobot(w http.ResponseWriter, r *http.Request) {
	robotID, ok := ParseUintID(w, r.PathValue("robotId"))
	if !ok {
		return
	}
	if err := repository.DeleteAlertRulesByRobotID(robotID); err != nil {
		applogger.Error("删除机器人告警规则失败: %v", err)
		JSONError(w, "删除机器人告警规则失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, map[string]bool{"success": true})
}
