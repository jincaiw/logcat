package alert

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"syslog-alert/internal/models"
)

// WeworkSender 企业微信推送实现。
type WeworkSender struct{}

func (s *WeworkSender) Send(robot *models.Robot, message string, parsedData map[string]interface{}, log *models.SyslogLog) error {
	return SendWeworkMessage(robot.WeworkWebhookURL, robot.WeworkKey, message)
}

func (s *WeworkSender) Test(robot *models.Robot) (string, error) {
	return SendWeworkTestMessage(robot.WeworkWebhookURL, robot.WeworkKey)
}

// SendWeworkMessage 发送企业微信告警消息（markdown 格式）。
func SendWeworkMessage(webhookURL, key, content string) error {
	webhookURL, err := buildWeworkURL(webhookURL, key)
	if err != nil {
		return err
	}

	message := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": content,
		},
	}

	body, err := postJSON(webhookURL, message)
	if err != nil {
		return err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	if errcode, ok := result["errcode"]; ok && errcode.(float64) != 0 {
		errmsg, _ := result["errmsg"].(string)
		return fmt.Errorf("wework api error: %s", errmsg)
	}

	return nil
}

// SendWeworkTestMessage 发送企业微信测试消息（text 格式）。
func SendWeworkTestMessage(webhookURL, key string) (string, error) {
	webhookURL, err := buildWeworkURL(webhookURL, key)
	if err != nil {
		return "", err
	}

	message := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": "【测试消息】Syslog告警系统连接测试成功！\n\n发送时间: " + time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	body, err := postJSON(webhookURL, message)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if errcode, ok := result["errcode"]; ok && errcode.(float64) != 0 {
		errmsg, _ := result["errmsg"].(string)
		return "", fmt.Errorf("wework api error: %s", errmsg)
	}

	return "测试消息发送成功！", nil
}

// buildWeworkURL 构建企业微信 webhook URL，附加 key 参数，消除 send/test 中的重复逻辑。
func buildWeworkURL(webhookURL, key string) (string, error) {
	if webhookURL == "" {
		return "", fmt.Errorf("webhook URL is empty")
	}

	if key == "" {
		return webhookURL, nil
	}

	if strings.Contains(webhookURL, "?") {
		return fmt.Sprintf("%s&key=%s", webhookURL, key), nil
	}
	return fmt.Sprintf("%s?key=%s", webhookURL, key), nil
}
