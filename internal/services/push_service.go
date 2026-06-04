package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
	"github.com/logcat/logcat/pkg/crypto"
)

// PushService handles HTTP push execution
type PushService struct{}

// NewPushService creates a new PushService
func NewPushService() *PushService {
	return &PushService{}
}

// PushResult holds the result of a push operation
type PushResult struct {
	Success      bool   `json:"success"`
	Channel      string `json:"channel"`
	StatusCode   int    `json:"statusCode"`
	ResponseBody string `json:"responseBody"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	Summary      string `json:"summary,omitempty"`
}

// ExecutePush executes a push configuration against the provided data
func (s *PushService) ExecutePush(pushConfigID uint, data map[string]interface{}) (*PushResult, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var config models.PushConfig
	if err := db.First(&config, pushConfigID).Error; err != nil {
		return nil, err
	}

	switch config.Type {
	case "http":
		return s.executeHTTP(&config, data)
	case "email":
		emailSvc := NewEmailService()
		result, err := emailSvc.SendEmail(pushConfigID, data)
		return s.fromEmailResult(&config, result, err), err
	case "syslog":
		syslogSvc := NewSyslogForwardService()
		result, err := syslogSvc.Forward(pushConfigID, data)
		return s.fromSyslogResult(&config, result, err), err
	default:
		return nil, fmt.Errorf("unsupported push type: %s", config.Type)
	}
}

// TestPush tests a push configuration
func (s *PushService) TestPush(pushConfigID uint) (*PushResult, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var config models.PushConfig
	if err := db.First(&config, pushConfigID).Error; err != nil {
		return nil, err
	}

	// Use test data
	testData := map[string]interface{}{
		"source_ip":      "192.168.1.100",
		"destination_ip": "10.0.0.1",
		"event_type":     "test",
		"severity":       "info",
		"message":        "Test push from logcat",
		"timestamp":      time.Now().Format(time.RFC3339),
	}

	switch config.Type {
	case "http":
		return s.executeHTTP(&config, testData)
	case "email":
		emailSvc := NewEmailService()
		result, err := emailSvc.SendEmail(pushConfigID, testData)
		return s.fromEmailResult(&config, result, err), err
	case "syslog":
		syslogSvc := NewSyslogForwardService()
		result, err := syslogSvc.Forward(pushConfigID, testData)
		return s.fromSyslogResult(&config, result, err), err
	default:
		return nil, fmt.Errorf("unsupported push type: %s", config.Type)
	}
}

func (s *PushService) executeHTTP(config *models.PushConfig, data map[string]interface{}) (*PushResult, error) {
	// Build request body
	tmplSvc := &TemplateService{}
	bodyStr := tmplSvc.Replace(config.BodyTemplate, data)

	// Parse headers
	headers := make(map[string]string)
	if config.Headers != "" {
		json.Unmarshal([]byte(config.Headers), &headers)
	}
	headers = tmplSvc.ReplaceAll(headers, data)

	// Set content type
	contentType := config.ContentType
	if contentType == "" {
		contentType = "application/json"
	}
	headers["Content-Type"] = contentType

	// Create request
	var body io.Reader
	if bodyStr != "" {
		body = bytes.NewBufferString(bodyStr)
	}

	method := config.Method
	if method == "" {
		method = "POST"
	}

	req, err := http.NewRequest(strings.ToUpper(method), config.URL, body)
	if err != nil {
		return &PushResult{
			Success:      false,
			Channel:      config.Type,
			ErrorMessage: fmt.Sprintf("failed to create request: %v", err),
			Summary:      "failed to create HTTP request",
		}, err
	}

	// Set headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Set authentication
	token := config.Token
	if token != "" {
		if dec, err := crypto.Decrypt(token); err == nil {
			token = dec
		}
	}
	switch config.AuthType {
	case "bearer":
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
	case "basic":
		if token != "" {
			req.SetBasicAuth("", token)
		}
	case "custom_header":
		if token != "" {
			req.Header.Set("X-Auth-Token", token)
		}
	}

	// Execute request with timeout
	timeout := config.Timeout
	if timeout <= 0 {
		timeout = 30
	}

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return &PushResult{
			Success:      false,
			Channel:      config.Type,
			ErrorMessage: fmt.Sprintf("request failed: %v", err),
			Summary:      "HTTP request failed",
		}, err
	}
	defer resp.Body.Close()

	limit := config.MaxResponseLogSize
	if limit <= 0 {
		limit = 4096
	}
	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, int64(limit+1)))
	respBodyStr := truncateText(string(respBody), limit)

	// Check success criteria
	success := resp.StatusCode >= 200 && resp.StatusCode < 300

	// Check custom success status codes
	if config.SuccessStatusCodes != "" {
		var codes []int
		if json.Unmarshal([]byte(config.SuccessStatusCodes), &codes) == nil {
			success = false
			for _, code := range codes {
				if resp.StatusCode == code {
					success = true
					break
				}
			}
		}
	}

	// Check success body keyword
	if success && config.SuccessBodyKeyword != "" {
		keyword := tmplSvc.Replace(config.SuccessBodyKeyword, data)
		if !strings.Contains(respBodyStr, keyword) {
			success = false
		}
	}

	return &PushResult{
		Success:      success,
		Channel:      config.Type,
		StatusCode:   resp.StatusCode,
		ResponseBody: respBodyStr,
		Summary:      summarizeText(respBodyStr),
	}, nil
}

func (s *PushService) fromEmailResult(config *models.PushConfig, result *EmailResult, err error) *PushResult {
	if result == nil && err == nil {
		return &PushResult{
			Success:      false,
			Channel:      config.Type,
			ErrorMessage: "email test returned empty result",
			Summary:      "email test returned empty result",
		}
	}

	pushResult := &PushResult{
		Channel:    config.Type,
		StatusCode: 200,
	}

	if result != nil {
		pushResult.Success = result.Success
		pushResult.ResponseBody = truncateText(result.Message, config.MaxResponseLogSize)
		pushResult.ErrorMessage = result.ErrorMessage
	}
	if err != nil && pushResult.ErrorMessage == "" {
		pushResult.ErrorMessage = err.Error()
	}
	if !pushResult.Success {
		pushResult.StatusCode = 500
	}
	if pushResult.ErrorMessage != "" {
		pushResult.Summary = summarizeText(pushResult.ErrorMessage)
	} else {
		pushResult.Summary = summarizeText(pushResult.ResponseBody)
	}

	return pushResult
}

func (s *PushService) fromSyslogResult(config *models.PushConfig, result *ForwardResult, err error) *PushResult {
	if result == nil && err == nil {
		return &PushResult{
			Success:      false,
			Channel:      config.Type,
			ErrorMessage: "syslog test returned empty result",
			Summary:      "syslog test returned empty result",
		}
	}

	pushResult := &PushResult{
		Channel:    config.Type,
		StatusCode: 200,
	}

	if result != nil {
		pushResult.Success = result.Success
		pushResult.ResponseBody = truncateText(result.Message, config.MaxResponseLogSize)
		pushResult.ErrorMessage = result.ErrorMessage
	}
	if err != nil && pushResult.ErrorMessage == "" {
		pushResult.ErrorMessage = err.Error()
	}
	if !pushResult.Success {
		pushResult.StatusCode = 500
	}
	if pushResult.ErrorMessage != "" {
		pushResult.Summary = summarizeText(pushResult.ErrorMessage)
	} else {
		pushResult.Summary = summarizeText(pushResult.ResponseBody)
	}

	return pushResult
}

func summarizeText(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	runes := []rune(text)
	if len(runes) <= 200 {
		return text
	}
	return string(runes[:200]) + "...(truncated)"
}

func truncateText(text string, limit int) string {
	text = strings.TrimSpace(text)
	if limit <= 0 {
		limit = 4096
	}
	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}
	return string(runes[:limit]) + "...(truncated)"
}

// RetryPush retries a push operation with retry count and delay
func (s *PushService) RetryPush(configID uint, data map[string]interface{}, maxRetries int, retryDelay int) *PushResult {
	for i := 0; i <= maxRetries; i++ {
		result, err := s.ExecutePush(configID, data)
		if err == nil && result.Success {
			return result
		}
		if i < maxRetries {
			time.Sleep(time.Duration(retryDelay) * time.Second)
		}
		if i == maxRetries {
			return result
		}
	}
	return &PushResult{Success: false, ErrorMessage: "all retries exhausted"}
}
