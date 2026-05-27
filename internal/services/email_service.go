package services

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
)

// EmailService handles SMTP email sending
type EmailService struct{}

// NewEmailService creates a new EmailService
func NewEmailService() *EmailService {
	return &EmailService{}
}

// EmailResult holds the result of an email send
type EmailResult struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// SendEmail sends an email using the push config's email settings
func (s *EmailService) SendEmail(pushConfigID uint, data map[string]interface{}) (*EmailResult, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var config models.PushConfig
	if err := db.First(&config, pushConfigID).Error; err != nil {
		return nil, err
	}

	if config.Type != "email" {
		return nil, errors.New("push config is not email type")
	}

	tmplSvc := &TemplateService{}
	subject := tmplSvc.Replace(config.SubjectTemplate, data)
	body := tmplSvc.Replace(config.EmailBodyTemplate, data)

	// Parse to addresses
	var toAddresses []string
	if err := json.Unmarshal([]byte(config.ToAddresses), &toAddresses); err != nil {
		toAddresses = strings.Split(config.ToAddresses, ",")
	}

	if len(toAddresses) == 0 {
		return &EmailResult{Success: false, ErrorMessage: "no recipient addresses"}, nil
	}

	// Build email message
	msg := fmt.Sprintf("From: %s\r\n", config.FromAddress) +
		fmt.Sprintf("To: %s\r\n", strings.Join(toAddresses, ",")) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body

	addr := fmt.Sprintf("%s:%d", config.SMTPHost, config.SMTPPort)

	// Use TLS if port is 465
	if config.SMTPPort == 465 {
		tlsConfig := &tls.Config{
			ServerName:         config.SMTPHost,
			InsecureSkipVerify: false,
		}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return &EmailResult{Success: false, ErrorMessage: fmt.Sprintf("TLS connection failed: %v", err)}, nil
		}
		defer conn.Close()

		client, err := smtp.NewClient(conn, config.SMTPHost)
		if err != nil {
			return &EmailResult{Success: false, ErrorMessage: fmt.Sprintf("SMTP client failed: %v", err)}, nil
		}
		defer client.Quit()

		if config.SMTPUsername != "" {
			auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost)
			if err := client.Auth(auth); err != nil {
				return &EmailResult{Success: false, ErrorMessage: fmt.Sprintf("auth failed: %v", err)}, nil
			}
		}

		if err := client.Mail(config.FromAddress); err != nil {
			return &EmailResult{Success: false, ErrorMessage: fmt.Sprintf("MAIL FROM failed: %v", err)}, nil
		}
		for _, to := range toAddresses {
			if err := client.Rcpt(strings.TrimSpace(to)); err != nil {
				return &EmailResult{Success: false, ErrorMessage: fmt.Sprintf("RCPT TO failed: %v", err)}, nil
			}
		}
		w, err := client.Data()
		if err != nil {
			return &EmailResult{Success: false, ErrorMessage: fmt.Sprintf("DATA failed: %v", err)}, nil
		}
		_, err = w.Write([]byte(msg))
		if err != nil {
			return &EmailResult{Success: false, ErrorMessage: fmt.Sprintf("write failed: %v", err)}, nil
		}
		w.Close()

		return &EmailResult{Success: true, Message: "email sent successfully"}, nil
	}

	// Use STARTTLS for port 587
	var auth smtp.Auth
	if config.SMTPUsername != "" {
		auth = smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost)
	}

	if err := smtp.SendMail(addr, auth, config.FromAddress, toAddresses, []byte(msg)); err != nil {
		return &EmailResult{Success: false, ErrorMessage: fmt.Sprintf("send failed: %v", err)}, nil
	}

	return &EmailResult{Success: true, Message: "email sent successfully"}, nil
}

// TestEmailConfig tests an email configuration
func (s *EmailService) TestEmailConfig(smtpHost string, smtpPort int, username, password, fromAddress, toAddress string) *EmailResult {
	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)

	msg := fmt.Sprintf("From: %s\r\n", fromAddress) +
		fmt.Sprintf("To: %s\r\n", toAddress) +
		"Subject: logcat Email Test\r\n" +
		"\r\n" +
		"This is a test email from logcat platform.\r\n"

	var auth smtp.Auth
	if username != "" {
		auth = smtp.PlainAuth("", username, password, smtpHost)
	}

	if err := smtp.SendMail(addr, auth, fromAddress, []string{toAddress}, []byte(msg)); err != nil {
		return &EmailResult{Success: false, ErrorMessage: fmt.Sprintf("send test failed: %v", err)}
	}

	return &EmailResult{Success: true, Message: "test email sent successfully"}
}