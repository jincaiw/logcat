package alert

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"syslog-alert/internal/models"
	applogger "syslog-alert/pkg/logger"
)

// FeishuSender 飞书推送实现。
type FeishuSender struct{}

func (s *FeishuSender) Send(robot *models.Robot, message string, parsedData map[string]interface{}, log *models.SyslogLog) error {
	return SendFeishuMessage(robot.FeishuWebhookURL, robot.FeishuSecret, message)
}

func (s *FeishuSender) Test(robot *models.Robot) (string, error) {
	return SendFeishuTestMessage(robot.FeishuWebhookURL, robot.FeishuSecret)
}

// 飞书消息结构定义
type FeishuMessage struct {
	MsgType string        `json:"msg_type"`
	Content FeishuContent `json:"content"`
}

type FeishuContent struct {
	Text string      `json:"text,omitempty"`
	Post *FeishuPost `json:"post,omitempty"`
}

type FeishuPost struct {
	ZhCN FeishuPostContent `json:"zh_cn"`
}

type FeishuPostContent struct {
	Title   string                `json:"title"`
	Content [][]FeishuPostElement `json:"content"`
}

type FeishuPostElement struct {
	Tag  string `json:"tag"`
	Text string `json:"text,omitempty"`
}

type FeishuResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// SendFeishuMessage 发送飞书告警消息（markdown 转 post 格式）。
func SendFeishuMessage(webhookURL, secret, content string) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL is empty")
	}

	applogger.Debug("SendFeishuMessage - webhookURL: %s, secret set: %t", webhookURL, secret != "")

	webhookURL = buildFeishuSignedURL(webhookURL, secret)
	applogger.Debug("After sign processing - webhookURL: %s", webhookURL)

	title, postContent := markdownToFeishuPost(content)
	if title == "" {
		title = "安全告警"
	}

	message := FeishuMessage{
		MsgType: "post",
		Content: FeishuContent{
			Post: &FeishuPost{
				ZhCN: FeishuPostContent{
					Title:   title,
					Content: postContent,
				},
			},
		},
	}

	applogger.Debug("Final webhookURL: %s", webhookURL)

	body, err := postJSON(webhookURL, message)
	if err != nil {
		return err
	}

	var result FeishuResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	if result.Code != 0 {
		return fmt.Errorf("feishu api error: %s", result.Msg)
	}

	return nil
}

// SendFeishuTestMessage 发送飞书测试消息（text 格式）。
func SendFeishuTestMessage(webhookURL, secret string) (string, error) {
	if webhookURL == "" {
		return "", fmt.Errorf("webhook URL is empty")
	}

	webhookURL = buildFeishuSignedURL(webhookURL, secret)

	message := FeishuMessage{
		MsgType: "text",
		Content: FeishuContent{
			Text: "【测试消息】Syslog告警系统连接测试成功！\n\n发送时间: " + time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	body, err := postJSON(webhookURL, message)
	if err != nil {
		return "", err
	}

	var result FeishuResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if result.Code != 0 {
		return "", fmt.Errorf("feishu api error: %s", result.Msg)
	}

	return "测试消息发送成功！", nil
}

// buildFeishuSignedURL 构建带加签参数的飞书 webhook URL，消除 send/test 中的重复签名逻辑。
func buildFeishuSignedURL(webhookURL, secret string) string {
	if secret == "" {
		return webhookURL
	}
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sign := generateFeishuSign(timestamp, secret)
	if strings.Contains(webhookURL, "?") {
		return fmt.Sprintf("%s&timestamp=%s&sign=%s", webhookURL, timestamp, url.QueryEscape(sign))
	}
	return fmt.Sprintf("%s?timestamp=%s&sign=%s", webhookURL, timestamp, url.QueryEscape(sign))
}

// generateFeishuSign 生成飞书加签签名。
func generateFeishuSign(timestamp, secret string) string {
	stringToSign := timestamp + "\n" + secret

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// markdownToFeishuPost 将 markdown 文本转换为飞书 post 格式的内容结构。
// 提取自原 SendFeishuMessage 中的内联逻辑，便于复用与测试。
func markdownToFeishuPost(content string) (string, [][]FeishuPostElement) {
	lines := strings.Split(content, "\n")
	var title string
	var postContent [][]FeishuPostElement

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "### ") {
			title = strings.TrimPrefix(line, "### ")
			continue
		}

		var elements []FeishuPostElement

		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				fieldName := strings.TrimSpace(parts[0])
				fieldValue := strings.TrimSpace(parts[1])

				fieldName = strings.ReplaceAll(fieldName, "**", "")
				fieldValue = strings.ReplaceAll(fieldValue, "**", "")

				elements = append(elements, FeishuPostElement{Tag: "text", Text: fieldName + ": " + fieldValue + "\n"})
			} else {
				line = strings.ReplaceAll(line, "**", "")
				elements = append(elements, FeishuPostElement{Tag: "text", Text: line + "\n"})
			}
		} else {
			line = strings.ReplaceAll(line, "**", "")
			elements = append(elements, FeishuPostElement{Tag: "text", Text: line + "\n"})
		}

		postContent = append(postContent, elements)
	}

	return title, postContent
}
