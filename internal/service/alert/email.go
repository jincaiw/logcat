package alert

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"syslog-alert/internal/models"
)

// EmailSender 邮箱推送实现。
type EmailSender struct{}

func (s *EmailSender) Send(robot *models.Robot, message string, parsedData map[string]interface{}, log *models.SyslogLog) error {
	return SendEmailMessage(robot.SMTPHost, robot.SMTPPort, robot.SMTPUsername, robot.SMTPPassword, robot.SMTPFrom, robot.SMTPTo, "【Syslog告警】安全告警通知", message)
}

func (s *EmailSender) Test(robot *models.Robot) (string, error) {
	return SendEmailTestMessage(robot.SMTPHost, robot.SMTPPort, robot.SMTPUsername, robot.SMTPPassword, robot.SMTPFrom, robot.SMTPTo)
}

// SendEmailMessage 发送告警邮件，支持 SSL（465）和 STARTTLS（587）。
func SendEmailMessage(host string, port int, username, password, from, to, subject, content string) error {
	if host == "" || from == "" || to == "" {
		return fmt.Errorf("email configuration is incomplete")
	}

	recipients := strings.Split(to, ",")
	for i, r := range recipients {
		recipients[i] = strings.TrimSpace(r)
	}

	fromAddr, err := mail.ParseAddress(from)
	if err != nil {
		fromAddr = &mail.Address{Address: from}
	}

	msg := fmt.Sprintf("From: %s\r\n", fromAddr.String())
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "MIME-version: 1.0;\r\nContent-Type: text/plain; charset=\"UTF-8\";\r\n\r\n"
	msg += strings.ReplaceAll(content, "**", "")

	addr := fmt.Sprintf("%s:%d", host, port)

	var client *smtp.Client

	if port == 465 {
		tlsConfig := &tls.Config{
			ServerName: host,
		}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server (SSL): %v", err)
		}
		client, err = smtp.NewClient(conn, host)
		if err != nil {
			conn.Close()
			return fmt.Errorf("failed to create SMTP client: %v", err)
		}
	} else {
		client, err = smtp.Dial(addr)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %v", err)
		}

		if port == 587 {
			if err := client.StartTLS(nil); err != nil {
				client.Close()
				return fmt.Errorf("failed to start TLS: %v", err)
			}
		}
	}
	defer client.Close()

	var auth smtp.Auth
	if username != "" && password != "" {
		auth = smtp.PlainAuth("", username, password, host)
	}

	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("authentication failed: %v", err)
		}
	}

	if err := client.Mail(fromAddr.Address); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}

	for _, recipient := range recipients {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to add recipient %s: %v", recipient, err)
		}
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to send data: %v", err)
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %v", err)
	}

	return nil
}

// SendEmailTestMessage 发送测试邮件。
func SendEmailTestMessage(host string, port int, username, password, from, to string) (string, error) {
	if host == "" || from == "" || to == "" {
		return "", fmt.Errorf("email configuration is incomplete")
	}

	subject := "【测试消息】Syslog告警系统连接测试"
	content := fmt.Sprintf("Syslog告警系统连接测试成功！\n\n发送时间: %s", time.Now().Format("2006-01-02 15:04:05"))

	err := SendEmailMessage(host, port, username, password, from, to, subject, content)
	if err != nil {
		return "", err
	}

	return "测试邮件发送成功！", nil
}
