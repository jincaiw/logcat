package api

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"syslog-alert/internal/models"
	"syslog-alert/internal/service/alert"
	"syslog-alert/internal/service/parser"
	applogger "syslog-alert/pkg/logger"
)

// SendTestSyslog 向指定目标发送测试 Syslog 消息。
func (ws *WebServer) SendTestSyslog(w http.ResponseWriter, r *http.Request) {
	var req models.TestSyslogRequest
	if !DecodeJSON(w, r, &req) {
		return
	}

	result := models.TestSyslogResult{
		Success:     true,
		Errors:      []string{},
		SentCount:   0,
		FailedCount: 0,
	}

	if req.Host == "" {
		req.Host = "127.0.0.1"
	}
	if req.Port == 0 {
		req.Port = 5140
	}
	if req.Protocol == "" {
		req.Protocol = "udp"
	}
	if req.Count <= 0 {
		req.Count = 1
	}
	if req.IntervalMs < 0 {
		req.IntervalMs = 0
	}

	addr := net.JoinHostPort(req.Host, fmt.Sprintf("%d", req.Port))

	var conn net.Conn
	var err error

	if req.Protocol == "tcp" {
		conn, err = net.Dial("tcp", addr)
	} else {
		conn, err = net.Dial("udp", addr)
	}

	if err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("连接失败: %v", err)
		JSONResponse(w, result)
		return
	}
	defer conn.Close()

	for i := 0; i < req.Count; i++ {
		_, err := conn.Write([]byte(req.Message))
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d条发送失败: %v", i+1, err))
		} else {
			result.SentCount++
		}

		if req.IntervalMs > 0 && i < req.Count-1 {
			time.Sleep(time.Duration(req.IntervalMs) * time.Millisecond)
		}
	}

	if result.FailedCount > 0 {
		result.Success = false
		result.Message = fmt.Sprintf("发送完成，成功%d条，失败%d条", result.SentCount, result.FailedCount)
	} else {
		result.Message = fmt.Sprintf("成功发送%d条测试日志", result.SentCount)
	}

	JSONResponse(w, result)
}

// testSyslogForwardRequest Syslog 转发测试请求体。
type testSyslogForwardRequest struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Format   string `json:"format"`
}

// TestSyslogForward 测试 Syslog 转发连通性。
func (ws *WebServer) TestSyslogForward(w http.ResponseWriter, r *http.Request) {
	var req testSyslogForwardRequest
	if !DecodeJSON(w, r, &req) {
		return
	}
	if err := alert.TestSyslogForward(req.Host, req.Port, req.Protocol, req.Format); err != nil {
		applogger.Error("测试 Syslog 转发失败: %v", err)
		JSONError(w, "测试失败: "+err.Error(), http.StatusInternalServerError)
		return
	}
	JSONResponse(w, map[string]string{"message": "测试消息发送成功！"})
}

// TestParse 测试解析模板。
func (ws *WebServer) TestParse(w http.ResponseWriter, r *http.Request) {
	var req models.ParseTestRequest
	if !DecodeJSON(w, r, &req) {
		return
	}

	result := models.ParseTestResult{
		Success: false,
		Fields:  []string{},
		Data:    make(map[string]interface{}),
	}

	if req.SampleLog == "" {
		result.Error = "请输入示例日志"
		JSONResponse(w, result)
		return
	}

	template := &models.ParseTemplate{
		ParseType:      req.ParseType,
		HeaderRegex:    req.HeaderRegex,
		FieldMapping:   req.FieldMapping,
		ValueTransform: req.ValueTransform,
	}

	p, err := parser.New(template)
	if err != nil {
		result.Error = "解析器初始化失败: " + err.Error()
		JSONResponse(w, result)
		return
	}

	data, err := p.Parse(req.SampleLog)
	if err != nil {
		result.Error = "解析失败: " + err.Error()
		JSONResponse(w, result)
		return
	}

	result.Success = true
	result.Data = data

	fieldSet := make(map[string]bool)
	for k := range data {
		fieldSet[k] = true
	}
	for k := range fieldSet {
		result.Fields = append(result.Fields, k)
	}

	JSONResponse(w, result)
}
