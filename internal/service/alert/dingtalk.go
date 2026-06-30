package alert

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"syslog-alert/internal/models"
)

// DingTalkSender 钉钉推送实现。
type DingTalkSender struct{}

func (s *DingTalkSender) Send(robot *models.Robot, message string, parsedData map[string]interface{}, log *models.SyslogLog) error {
	return SendDingTalkMessage(robot.WebhookURL, robot.Secret, message)
}

func (s *DingTalkSender) Test(robot *models.Robot) (string, error) {
	return SendDingTalkTestMessage(robot.WebhookURL, robot.Secret)
}

// DingTalk 消息结构定义
type DingTalkMessage struct {
	MsgType  string            `json:"msgtype"`
	Markdown *DingTalkMarkdown `json:"markdown,omitempty"`
	Text     *DingTalkText     `json:"text,omitempty"`
}

type DingTalkMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type DingTalkText struct {
	Content string `json:"content"`
}

type DingTalkResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// SendDingTalkMessage 发送钉钉告警消息（markdown 格式）。
func SendDingTalkMessage(webhookURL, secret, content string) error {
	webhookURL = buildDingTalkSignedURL(webhookURL, secret)

	message := DingTalkMessage{
		MsgType: "markdown",
		Markdown: &DingTalkMarkdown{
			Title: "Syslog告警",
			Text:  content,
		},
	}

	body, err := postJSON(webhookURL, message)
	if err != nil {
		return err
	}

	var result DingTalkResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("dingtalk api error: %s", result.ErrMsg)
	}

	return nil
}

// SendDingTalkTestMessage 发送钉钉测试消息（text 格式）。
func SendDingTalkTestMessage(webhookURL, secret string) (string, error) {
	webhookURL = buildDingTalkSignedURL(webhookURL, secret)

	message := DingTalkMessage{
		MsgType: "text",
		Text: &DingTalkText{
			Content: "【测试消息】Syslog告警系统连接测试成功！\n\n发送时间: " + time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	body, err := postJSON(webhookURL, message)
	if err != nil {
		return "", err
	}

	var result DingTalkResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("dingtalk api error: %s", result.ErrMsg)
	}

	return "测试消息发送成功！", nil
}

// buildDingTalkSignedURL 构建带加签参数的钉钉 webhook URL，消除 send/test 中的重复签名逻辑。
func buildDingTalkSignedURL(webhookURL, secret string) string {
	if secret == "" {
		return webhookURL
	}
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	sign := generateSign(timestamp, secret)
	return fmt.Sprintf("%s&timestamp=%s&sign=%s", webhookURL, timestamp, url.QueryEscape(sign))
}

// generateSign 生成钉钉加签签名。
func generateSign(timestamp, secret string) string {
	stringToSign := timestamp + "\n" + secret

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
