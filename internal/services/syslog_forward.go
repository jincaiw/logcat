package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
)

// SyslogForwardService handles syslog forwarding
type SyslogForwardService struct{}

// NewSyslogForwardService creates a new SyslogForwardService
func NewSyslogForwardService() *SyslogForwardService {
	return &SyslogForwardService{}
}

// ForwardResult holds the result of syslog forwarding
type ForwardResult struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// Forward sends a syslog message using the push config's syslog settings
func (s *SyslogForwardService) Forward(pushConfigID uint, data map[string]interface{}) (*ForwardResult, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var config models.PushConfig
	if err := db.First(&config, pushConfigID).Error; err != nil {
		return nil, err
	}

	if config.Type != "syslog" {
		return nil, errors.New("push config is not syslog type")
	}

	// Build syslog message
	tmplSvc := &TemplateService{}
	message := s.buildSyslogMessage(&config, data, tmplSvc)

	addr := fmt.Sprintf("%s:%d", config.SyslogHost, config.SyslogPort)

	switch strings.ToLower(config.SyslogProtocol) {
	case "tcp":
		return s.forwardTCP(addr, message)
	case "udp":
		fallthrough
	default:
		return s.forwardUDP(addr, message)
	}
}

func (s *SyslogForwardService) buildSyslogMessage(config *models.PushConfig, data map[string]interface{}, tmplSvc *TemplateService) string {
	format := config.SyslogFormat
	if format == "" {
		format = "rfc3164"
	}

	now := time.Now()

	switch format {
	case "json":
		// Build JSON syslog
		fields := make(map[string]interface{})
		if config.SyslogFields != "" {
			json.Unmarshal([]byte(config.SyslogFields), &fields)
		}
		for k, v := range fields {
			if s, ok := v.(string); ok {
				fields[k] = tmplSvc.Replace(s, data)
			}
		}
		// Add standard fields
		fields["timestamp"] = now.Format(time.RFC3339)
		fields["hostname"] = config.SyslogHost
		if body, err := json.Marshal(fields); err == nil {
			return string(body)
		}
		return fmt.Sprintf(`{"timestamp":"%s","message":"%v"}`, now.Format(time.RFC3339), data)

	case "rfc5424":
		return s.buildRFC5424(now, data)

	case "rfc3164":
		fallthrough
	default:
		return s.buildRFC3164(now, data)
	}
}

func (s *SyslogForwardService) buildRFC3164(now time.Time, data map[string]interface{}) string {
	hostname := fmt.Sprintf("%v", data["source_ip"])
	timestamp := now.Format("Jan 2 15:04:05")
	message := fmt.Sprintf("%v", data["message"])
	if message == "" {
		message = fmt.Sprintf("%v", data)
	}
	return fmt.Sprintf("<%d>%s %s %s", getPriority(data), timestamp, hostname, message)
}

func (s *SyslogForwardService) buildRFC5424(now time.Time, data map[string]interface{}) string {
	pri := getPriority(data)
	version := 1
	timestamp := now.Format(time.RFC3339Nano)
	hostname := fmt.Sprintf("%v", data["source_ip"])
	appName := "logcat"
	procID := "-"
	msgID := "-"
	structuredData := "-"
	message := fmt.Sprintf("%v", data["message"])

	return fmt.Sprintf("<%d>%d %s %s %s %s %s %s %s",
		pri, version, timestamp, hostname, appName, procID, msgID, structuredData, message)
}

func getPriority(data map[string]interface{}) int {
	facility := 1  // user-level messages
	severity := 6  // informational

	if f, ok := data["facility"].(float64); ok {
		facility = int(f)
	}
	if s, ok := data["severity"].(float64); ok {
		severity = int(s)
	}

	return facility*8 + severity
}

func (s *SyslogForwardService) forwardUDP(addr, message string) (*ForwardResult, error) {
	conn, err := net.DialTimeout("udp", addr, 5*time.Second)
	if err != nil {
		return &ForwardResult{Success: false, ErrorMessage: fmt.Sprintf("UDP connection failed: %v", err)}, nil
	}
	defer conn.Close()

	_, err = conn.Write([]byte(message))
	if err != nil {
		return &ForwardResult{Success: false, ErrorMessage: fmt.Sprintf("UDP send failed: %v", err)}, nil
	}

	return &ForwardResult{Success: true, Message: "syslog forwarded via UDP"}, nil
}

func (s *SyslogForwardService) forwardTCP(addr, message string) (*ForwardResult, error) {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return &ForwardResult{Success: false, ErrorMessage: fmt.Sprintf("TCP connection failed: %v", err)}, nil
	}
	defer conn.Close()

	// RFC 6587: add newline for octet framing
	_, err = fmt.Fprintf(conn, "%s\n", message)
	if err != nil {
		return &ForwardResult{Success: false, ErrorMessage: fmt.Sprintf("TCP send failed: %v", err)}, nil
	}

	return &ForwardResult{Success: true, Message: "syslog forwarded via TCP"}, nil
}