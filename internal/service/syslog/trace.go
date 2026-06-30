package syslog

import (
	"time"

	"syslog-alert/internal/models"
	"syslog-alert/internal/repository"
)

// GetTraceInfo 获取日志的全链路追踪信息。
// 优先从内存缓存读取，缓存未命中时从数据库重建。
func (s *Server) GetTraceInfo(logID uint) *models.LogTraceInfo {
	// 优先读内存缓存
	s.traceMu.RLock()
	if trace, ok := s.traceMap[logID]; ok {
		s.traceMu.RUnlock()
		return trace
	}
	s.traceMu.RUnlock()

	// 缓存未命中，从数据库重建
	return s.buildTraceFromDB(logID)
}

// buildTraceFromDB 从数据库重建追踪信息。
func (s *Server) buildTraceFromDB(logID uint) *models.LogTraceInfo {
	log, err := repository.GetLogByID(logID)
	if err != nil {
		return nil
	}

	trace := &models.LogTraceInfo{
		LogID:         log.ID,
		ReceivedAt:    log.ReceivedAt,
		SourceIP:      log.SourceIP,
		RawMessage:    log.RawMessage,
		ReceiveStatus: "success",
		ParseStatus:   "success",
		ParsedData:    log.ParsedData,
		FilterStatus:  log.FilterStatus,
		AlertStatus:   log.AlertStatus,
	}

	if log.MatchedPolicyID > 0 {
		if policy, err := repository.GetFilterPolicyByID(log.MatchedPolicyID); err == nil {
			trace.MatchedPolicy = policy.Name
			trace.FilterEnabled = true
		}
	}

	// 加载告警记录
	s.loadAlertRecordsIntoTrace(logID, trace)
	return trace
}

// loadAlertRecordsIntoTrace 从数据库加载告警记录到追踪信息。
func (s *Server) loadAlertRecordsIntoTrace(logID uint, trace *models.LogTraceInfo) {
	// 使用 AlertRecord 查询（repository 暂无专门方法，直接用 DB）
	// 注意：这里通过 repository 包暴露的能力查询
	var records []models.AlertRecord
	repository.DB().Where("log_id = ?", logID).Find(&records)

	for _, record := range records {
		robotName := ""
		platform := ""
		if robot, err := repository.GetRobotByID(record.RobotID); err == nil {
			robotName = robot.Name
			platform = robot.Platform
		}
		trace.AlertRecords = append(trace.AlertRecords, models.AlertTraceInfo{
			RobotID:   record.RobotID,
			RobotName: robotName,
			Platform:  platform,
			Status:    record.Status,
			ErrorMsg:  record.ErrorMsg,
			SentAt:    record.SentAt,
		})
	}
}

// ---- 追踪信息更新方法 ----

// createTrace 创建日志追踪记录。
func (s *Server) createTrace(logID uint, sourceIP, rawMessage string) {
	s.traceMu.Lock()
	defer s.traceMu.Unlock()
	s.traceMap[logID] = &models.LogTraceInfo{
		LogID:         logID,
		ReceivedAt:    time.Now(),
		SourceIP:      sourceIP,
		RawMessage:    rawMessage,
		ReceiveStatus: "success",
	}
}

// updateTraceParse 更新解析阶段的追踪信息。
func (s *Server) updateTraceParse(logID uint, status, templateName, parsedData, parseError string) {
	s.traceMu.Lock()
	defer s.traceMu.Unlock()
	if trace, ok := s.traceMap[logID]; ok {
		trace.ParseStatus = status
		trace.ParseTemplate = templateName
		trace.ParsedData = parsedData
		trace.ParseError = parseError
	}
}

// updateTraceFilter 更新过滤阶段的追踪信息。
func (s *Server) updateTraceFilter(logID uint, status string, filterEnabled bool, matchedPolicy, filterConditions, filterResult string) {
	s.traceMu.Lock()
	defer s.traceMu.Unlock()
	if trace, ok := s.traceMap[logID]; ok {
		trace.FilterStatus = status
		trace.FilterEnabled = filterEnabled
		trace.MatchedPolicy = matchedPolicy
		trace.FilterConditions = filterConditions
		trace.FilterResult = filterResult
	}
}

// updateTraceAlert 更新告警阶段的追踪信息。
func (s *Server) updateTraceAlert(logID uint, status string) {
	s.traceMu.Lock()
	defer s.traceMu.Unlock()
	if trace, ok := s.traceMap[logID]; ok {
		trace.AlertStatus = status
	}
}

// addTraceAlertRecord 添加告警记录到追踪信息。
func (s *Server) addTraceAlertRecord(logID uint, record models.AlertTraceInfo) {
	s.traceMu.Lock()
	defer s.traceMu.Unlock()
	if trace, ok := s.traceMap[logID]; ok {
		trace.AlertRecords = append(trace.AlertRecords, record)
	}
}

// ClearOldTraces 清理过期的追踪缓存。
func (s *Server) ClearOldTraces(maxAge time.Duration) {
	s.traceMu.Lock()
	defer s.traceMu.Unlock()
	cutoff := time.Now().Add(-maxAge)
	for logID, trace := range s.traceMap {
		if trace.ReceivedAt.Before(cutoff) {
			delete(s.traceMap, logID)
		}
	}
}
