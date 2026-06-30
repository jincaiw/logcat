// Package alert 提供统一的告警推送能力，支持钉钉、飞书、企业微信、邮箱和 Syslog 转发。
package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"syslog-alert/internal/models"
	"syslog-alert/pkg/constants"
)

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
	case constants.PlatformDingTalk:
		return &DingTalkSender{}, nil
	case constants.PlatformFeishu:
		return &FeishuSender{}, nil
	case constants.PlatformWework:
		return &WeworkSender{}, nil
	case constants.PlatformEmail:
		return &EmailSender{}, nil
	case constants.PlatformSyslog:
		return &SyslogForwardSender{}, nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", robot.Platform)
	}
}

// postJSON 发送 JSON POST 请求并返回响应体，供各推送通道复用以消除重复代码。
func postJSON(targetURL string, body interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %v", err)
	}
	resp, err := http.Post(targetURL, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	return responseBody, nil
}
