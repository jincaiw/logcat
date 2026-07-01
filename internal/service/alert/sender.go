// Package alert 提供统一的告警推送能力，支持飞书、邮箱、Syslog 转发和 HTTP 接口。
package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"syslog-alert/internal/models"
	"syslog-alert/pkg/constants"
)

const (
	defaultOutboundRequestTimeout = 10 * time.Second
	maxOutboundResponseBytes      = 32 << 10 // 32KB
)

var sharedHTTPTransport = &http.Transport{
	Proxy:                 http.ProxyFromEnvironment,
	DialContext:           (&net.Dialer{Timeout: 5 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          32,
	MaxIdleConnsPerHost:   8,
	MaxConnsPerHost:       0,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   5 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

var sharedHTTPClient = &http.Client{Transport: sharedHTTPTransport}

// Sender 定义告警推送通道的统一接口。
type Sender interface {
	Send(robot *models.Robot, message string, parsedData map[string]interface{}, log *models.SyslogLog) error
	Test(robot *models.Robot) (string, error)
}

// Send 根据 robot.Platform 分发到具体的推送实现。
func Send(robot *models.Robot, message string, parsedData map[string]interface{}, log *models.SyslogLog) error {
	sender, err := getSender(robot)
	if err != nil {
		return err
	}
	return sender.Send(robot, message, parsedData, log)
}

// Test 根据 robot.Platform 分发到具体的测试实现。
func Test(robot *models.Robot) (string, error) {
	sender, err := getSender(robot)
	if err != nil {
		return "", err
	}
	return sender.Test(robot)
}

// getSender 根据平台类型返回对应的 Sender 实现。
func getSender(robot *models.Robot) (Sender, error) {
	switch robot.Platform {
	case constants.PlatformFeishu:
		return &FeishuSender{}, nil
	case constants.PlatformEmail:
		return &EmailSender{}, nil
	case constants.PlatformSyslog:
		return &SyslogForwardSender{}, nil
	case constants.PlatformHTTP:
		return &HTTPSender{}, nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", robot.Platform)
	}
}

// postJSON 发送 JSON POST 请求并返回响应体，供各推送通道复用以消除重复代码。
func postJSON(targetURL string, body interface{}) ([]byte, error) {
	respBody, statusCode, err := doJSONRequest(targetURL, defaultOutboundRequestTimeout, body)
	if err != nil {
		return nil, err
	}
	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("HTTP status %d: %s", statusCode, truncateString(strings.TrimSpace(string(respBody)), 1024))
	}
	return respBody, nil
}

func doJSONRequest(targetURL string, timeout time.Duration, body interface{}) ([]byte, int, error) {
	if strings.TrimSpace(targetURL) == "" {
		return nil, 0, fmt.Errorf("target URL is empty")
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal message: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewReader(jsonData))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")

	resp, err := sharedHTTPClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	respBody, truncated, err := readLimitedBody(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read response: %v", err)
	}
	if truncated {
		respBody = append(respBody, []byte("...(truncated)")...)
	}
	return respBody, resp.StatusCode, nil
}

func readLimitedBody(r io.Reader) ([]byte, bool, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxOutboundResponseBytes+1))
	if err != nil {
		return nil, false, err
	}
	if len(data) <= maxOutboundResponseBytes {
		return data, false, nil
	}
	return data[:maxOutboundResponseBytes], true, nil
}

func truncateString(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
