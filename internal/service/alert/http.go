package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"syslog-alert/internal/models"
)

const (
	defaultSendWBGURL        = "http://168.63.6.81:8080/cib-message/public/service/sendwbg.do"
	defaultSendWBGTimeout    = 3
	defaultSendWBGRetryCount = 3
	defaultSendWBGRetryDelay = 2
)

var defaultSendWBGNotesIDs = []string{"420102", "420809"}

// HTTPSender HTTP 接口推送实现。
type HTTPSender struct{}

func (s *HTTPSender) Send(robot *models.Robot, message string, parsedData map[string]interface{}, log *models.SyslogLog) error {
	return SendHTTPMessage(robot, message, parsedData, log)
}

func (s *HTTPSender) Test(robot *models.Robot) (string, error) {
	msg := "【logcat测试消息】HTTP接口通知渠道测试"
	if err := SendHTTPMessage(robot, msg, map[string]interface{}{"test": true}, nil); err != nil {
		return "", err
	}
	return "HTTP接口测试消息发送成功", nil
}

// SendHTTPMessage 按配置向 HTTP 接口发送告警，支持超时、失败重试与接收人 ID 列表。
func SendHTTPMessage(robot *models.Robot, message string, parsedData map[string]interface{}, log *models.SyslogLog) error {
	cfg := normalizeHTTPConfig(robot)
	payload := buildHTTPPayload(message, parsedData, log, cfg.notesIDs)

	var lastErr error
	for attempt := 0; attempt <= cfg.retryCount; attempt++ {
		lastErr = postHTTPPayload(cfg.targetURL, cfg.timeout, payload)
		if lastErr == nil {
			return nil
		}
		if attempt < cfg.retryCount {
			time.Sleep(time.Duration(cfg.retryDelay) * time.Second)
		}
	}
	return fmt.Errorf("HTTP接口发送失败（已重试 %d 次）: %w", cfg.retryCount, lastErr)
}

type httpConfig struct {
	targetURL  string
	timeout    int
	retryCount int
	retryDelay int
	notesIDs   []string
}

func normalizeHTTPConfig(robot *models.Robot) httpConfig {
	cfg := httpConfig{
		targetURL:  defaultSendWBGURL,
		timeout:    defaultSendWBGTimeout,
		retryCount: defaultSendWBGRetryCount,
		retryDelay: defaultSendWBGRetryDelay,
		notesIDs:   append([]string(nil), defaultSendWBGNotesIDs...),
	}
	if robot == nil {
		return cfg
	}
	if strings.TrimSpace(robot.HTTPURL) != "" {
		cfg.targetURL = strings.TrimSpace(robot.HTTPURL)
	}
	if robot.HTTPTimeout > 0 {
		cfg.timeout = robot.HTTPTimeout
	}
	if robot.HTTPRetryCount >= 0 {
		cfg.retryCount = robot.HTTPRetryCount
	}
	if robot.HTTPRetryDelay >= 0 {
		cfg.retryDelay = robot.HTTPRetryDelay
	}
	if notesIDs := parseHTTPNotesIDs(robot.HTTPNotesIDs); len(notesIDs) > 0 {
		cfg.notesIDs = notesIDs
	}
	return cfg
}

func parseHTTPNotesIDs(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	var ids []string
	if strings.HasPrefix(raw, "[") {
		if err := json.Unmarshal([]byte(raw), &ids); err == nil {
			return compactStringList(ids)
		}
	}

	splitter := func(r rune) bool {
		return r == ',' || r == ';' || r == '\n' || r == '\r' || r == '\t' || r == ' '
	}
	return compactStringList(strings.FieldsFunc(raw, splitter))
}

func compactStringList(items []string) []string {
	result := make([]string, 0, len(items))
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		item = strings.Trim(strings.TrimSpace(item), `"'`)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}

func buildHTTPPayload(message string, parsedData map[string]interface{}, log *models.SyslogLog, notesIDs []string) map[string]interface{} {
	payload := map[string]interface{}{
		"notesids":   notesIDs,
		"notesIds":   notesIDs,
		"content":    message,
		"message":    message,
		"text":       message,
		"parsedData": parsedData,
	}
	if log != nil {
		payload["log"] = map[string]interface{}{
			"id":           log.ID,
			"deviceId":     log.DeviceID,
			"deviceName":   log.DeviceName,
			"sourceIp":     log.SourceIP,
			"rawMessage":   log.RawMessage,
			"priority":     log.Priority,
			"facility":     log.Facility,
			"severity":     log.Severity,
			"timestamp":    log.Timestamp,
			"receivedAt":   log.ReceivedAt,
			"parsedData":   log.ParsedData,
			"parsedFields": log.ParsedFields,
		}
	}
	return payload
}

func postHTTPPayload(targetURL string, timeoutSeconds int, payload map[string]interface{}) error {
	if strings.TrimSpace(targetURL) == "" {
		return fmt.Errorf("HTTP接口URL为空")
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal HTTP payload: %w", err)
	}

	client := &http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second}
	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return fmt.Errorf("failed to read HTTP response: %w", readErr)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("HTTP status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}
	return nil
}
