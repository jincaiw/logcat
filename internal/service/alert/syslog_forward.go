package alert

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"syslog-alert/internal/models"
	"syslog-alert/pkg/constants"
)

// SyslogForwardSender Syslog 转发推送实现。
//
// 注意：Sender 接口的 Send 方法仅包含基础参数，无法传递字段映射等高级配置。
// 需要字段映射/值转换/字段筛选的调用方应直接调用 SendSyslogForward。
type SyslogForwardSender struct{}

func (s *SyslogForwardSender) Send(robot *models.Robot, message string, parsedData map[string]interface{}, log *models.SyslogLog) error {
	return SendSyslogForward(robot.SyslogHost, robot.SyslogPort, robot.SyslogProtocol, robot.SyslogFormat, message, parsedData, log, "", nil, nil, "")
}

func (s *SyslogForwardSender) Test(robot *models.Robot) (string, error) {
	if err := TestSyslogForward(robot.SyslogHost, robot.SyslogPort, robot.SyslogProtocol, robot.SyslogFormat); err != nil {
		return "", err
	}
	return "测试消息发送成功！", nil
}

// SyslogForwardMessage Syslog 转发的 JSON 消息结构。
type SyslogForwardMessage struct {
	Timestamp   string                 `json:"timestamp"`
	SourceIP    string                 `json:"sourceIp"`
	DeviceName  string                 `json:"deviceName"`
	Facility    int                    `json:"facility"`
	Severity    int                    `json:"severity"`
	Message     string                 `json:"message"`
	RawLog      string                 `json:"rawLog"`
	ParsedData  map[string]interface{} `json:"parsedData,omitempty"`
	Forwarded   bool                   `json:"forwarded"`
	ForwardedBy string                 `json:"forwardedBy,omitempty"`
}

// FieldMappingItem 字段映射配置项。
type FieldMappingItem struct {
	SourceField string `json:"sourceField"`
	DisplayName string `json:"displayName"`
}

// SendSyslogForward 将告警消息以 Syslog 协议转发到指定目标，支持 JSON/RFC3164/RFC5424 格式。
// 该函数保留原始签名以兼容需要字段映射等高级配置的调用方。
func SendSyslogForward(host string, port int, protocol string, format string, message string, parsedData map[string]interface{}, log *models.SyslogLog, fieldMapping string, fieldNameMapping map[string]string, selectedFields []string, valueTransform string) error {
	if host == "" || port == 0 {
		return fmt.Errorf("syslog host or port is empty")
	}

	mappedData := applyFieldMapping(parsedData, fieldMapping, fieldNameMapping)
	filteredData := filterFieldsBySelection(mappedData, selectedFields)

	if len(selectedFields) > 0 && len(filteredData) < len(selectedFields) {
		reverseMapping := make(map[string]string)
		for eng, chn := range fieldNameMapping {
			reverseMapping[chn] = eng
		}
		if fieldMapping != "" {
			var fieldMappings map[string]string
			if err := json.Unmarshal([]byte(fieldMapping), &fieldMappings); err == nil {
				for eng, chn := range fieldMappings {
					reverseMapping[chn] = eng
				}
			}
		}

		for _, field := range selectedFields {
			if _, exists := filteredData[field]; exists {
				continue
			}
			if value, exists := mappedData[field]; exists {
				filteredData[field] = value
			} else if engField, hasMapping := reverseMapping[field]; hasMapping {
				if value, exists := parsedData[engField]; exists {
					filteredData[field] = value
				}
			} else if value, exists := parsedData[field]; exists {
				filteredData[field] = value
			}
		}
	}

	if len(filteredData) > 0 && valueTransform != "" {
		filteredData = applyValueTransform(filteredData, valueTransform)
	} else if len(mappedData) > 0 && valueTransform != "" {
		transformedMappedData := applyValueTransform(mappedData, valueTransform)
		if len(filteredData) > 0 {
			for k, v := range transformedMappedData {
				if _, existsInFiltered := filteredData[k]; existsInFiltered {
					filteredData[k] = v
				}
			}
		}
	}

	var payload []byte
	var err error

	hostname, _ := os.Hostname()

	switch format {
	case constants.FormatJSON:
		payload, err = buildSyslogJSONPayload(filteredData, selectedFields, message, log, hostname)
		if err != nil {
			return err
		}
	case constants.FormatRFC3164:
		payload = buildSyslogRFC3164Payload(filteredData, mappedData, message, log)
	case constants.FormatRFC5424:
		payload = buildSyslogRFC5424Payload(filteredData, mappedData, message, log)
	default:
		payload, err = buildSyslogJSONPayload(filteredData, selectedFields, message, log, hostname)
		if err != nil {
			return err
		}
	}

	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	protocol = strings.ToLower(protocol)
	if protocol == "" {
		protocol = constants.ProtocolUDP
	}

	if protocol == constants.ProtocolTCP {
		conn, err := net.Dial(constants.ProtocolTCP, address)
		if err != nil {
			return fmt.Errorf("failed to connect to %s: %v", address, err)
		}
		defer conn.Close()
		_, err = conn.Write(payload)
		if err != nil {
			return fmt.Errorf("failed to send tcp message: %v", err)
		}
	} else {
		if len(payload) > constants.MaxUDPPacketSize {
			truncated, truncErr := truncateLargeFields(payload, filteredData, selectedFields, constants.MaxUDPPacketSize)
			if truncErr != nil {
				return fmt.Errorf("UDP数据包大小超限（当前 %d 字节，最大 %d 字节），截断后仍超限。建议：1. 使用TCP协议；2. 减少推送字段；3. 使用更短的字段名称", len(payload), constants.MaxUDPPacketSize)
			}
			payload = truncated
		}

		conn, err := net.Dial(constants.ProtocolUDP, address)
		if err != nil {
			return fmt.Errorf("failed to connect to %s: %v", address, err)
		}
		defer conn.Close()
		_, err = conn.Write(payload)
		if err != nil {
			return fmt.Errorf("failed to send udp message: %v", err)
		}
	}

	return nil
}

// buildSyslogJSONPayload 构建 JSON 格式的 Syslog 转发 payload。
func buildSyslogJSONPayload(filteredData map[string]interface{}, selectedFields []string, message string, log *models.SyslogLog, hostname string) ([]byte, error) {
	if len(selectedFields) > 0 && len(filteredData) > 0 {
		return json.Marshal(filteredData)
	}
	msg := SyslogForwardMessage{
		Timestamp:   time.Now().Format(time.RFC3339),
		SourceIP:    log.SourceIP,
		DeviceName:  log.DeviceName,
		Facility:    log.Facility,
		Severity:    log.Severity,
		Message:     message,
		RawLog:      log.RawMessage,
		ParsedData:  filteredData,
		Forwarded:   true,
		ForwardedBy: hostname,
	}
	return json.Marshal(msg)
}

// buildSyslogRFC3164Payload 构建 RFC3164 格式的 Syslog 转发 payload。
func buildSyslogRFC3164Payload(filteredData, mappedData map[string]interface{}, message string, log *models.SyslogLog) []byte {
	ts := time.Now().Format("Jan 2 15:04:05")
	hostname := log.SourceIP
	if hostname == "" {
		hostname = "unknown"
	}

	var msgContent string
	if len(filteredData) > 0 {
		msgContent = formatFieldsAsKeyValue(filteredData)
	} else if len(mappedData) > 0 {
		dataJSON, _ := json.Marshal(mappedData)
		msgContent = fmt.Sprintf("%s | Data: %s", message, string(dataJSON))
	} else {
		msgContent = message
	}
	return []byte(fmt.Sprintf("<134>%s %s logcat: [FORWARDED] %s", ts, hostname, msgContent))
}

// buildSyslogRFC5424Payload 构建 RFC5424 格式的 Syslog 转发 payload。
func buildSyslogRFC5424Payload(filteredData, mappedData map[string]interface{}, message string, log *models.SyslogLog) []byte {
	ts := time.Now().Format(time.RFC3339)
	hostname := log.SourceIP
	if hostname == "" {
		hostname = "unknown"
	}

	var msgContent string
	if len(filteredData) > 0 {
		msgContent = formatFieldsAsKeyValue(filteredData)
	} else if len(mappedData) > 0 {
		dataJSON, _ := json.Marshal(mappedData)
		msgContent = fmt.Sprintf("%s | Data: %s", message, string(dataJSON))
	} else {
		msgContent = message
	}
	return []byte(fmt.Sprintf("<134>1 %s %s logcat - - - [FORWARDED] %s", ts, hostname, msgContent))
}

// TestSyslogForward 发送 Syslog 测试消息。
func TestSyslogForward(host string, port int, protocol string, format string) error {
	if host == "" || port == 0 {
		return fmt.Errorf("syslog host or port is empty")
	}

	var payload []byte
	var err error

	switch format {
	case constants.FormatRFC3164:
		ts := time.Now().Format("Jan 2 15:04:05")
		payload = []byte(fmt.Sprintf("<134>%s 127.0.0.1 logcat: 【测试消息】logcat连接测试成功！", ts))
	case constants.FormatRFC5424:
		ts := time.Now().Format(time.RFC3339)
		payload = []byte(fmt.Sprintf("<134>1 %s 127.0.0.1 logcat - - - 【测试消息】logcat连接测试成功！", ts))
	default:
		testMessage := SyslogForwardMessage{
			Timestamp:  time.Now().Format(time.RFC3339),
			SourceIP:   "127.0.0.1",
			DeviceName: "logcat",
			Facility:   1,
			Severity:   6,
			Message:    "【测试消息】logcat连接测试成功！",
		}
		payload, err = json.Marshal(testMessage)
		if err != nil {
			return fmt.Errorf("failed to marshal json: %v", err)
		}
	}

	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	protocol = strings.ToLower(protocol)
	if protocol == "" {
		protocol = constants.ProtocolUDP
	}

	if protocol == constants.ProtocolTCP {
		conn, err := net.Dial(constants.ProtocolTCP, address)
		if err != nil {
			return fmt.Errorf("failed to connect to %s: %v", address, err)
		}
		defer conn.Close()
		_, err = conn.Write(payload)
		if err != nil {
			return fmt.Errorf("failed to send tcp message: %v", err)
		}
	} else {
		conn, err := net.Dial(constants.ProtocolUDP, address)
		if err != nil {
			return fmt.Errorf("failed to connect to %s: %v", address, err)
		}
		defer conn.Close()
		_, err = conn.Write(payload)
		if err != nil {
			return fmt.Errorf("failed to send udp message: %v", err)
		}
	}

	return nil
}

// applyFieldMapping 根据字段映射配置重命名字段。
func applyFieldMapping(parsedData map[string]interface{}, fieldMappingJSON string, fieldNameMapping map[string]string) map[string]interface{} {
	if fieldMappingJSON == "" && len(fieldNameMapping) == 0 {
		return parsedData
	}

	result := make(map[string]interface{})

	for k, v := range parsedData {
		result[k] = v
	}

	if len(fieldNameMapping) > 0 {
		for sourceField, displayName := range fieldNameMapping {
			if value, exists := parsedData[sourceField]; exists {
				result[displayName] = value
				delete(result, sourceField)
			}
		}
	}

	if fieldMappingJSON != "" {
		var fieldMappings map[string]string
		if err := json.Unmarshal([]byte(fieldMappingJSON), &fieldMappings); err == nil && len(fieldMappings) > 0 {
			for sourceField, displayName := range fieldMappings {
				if value, exists := parsedData[sourceField]; exists {
					result[displayName] = value
					delete(result, sourceField)
				}
			}
		}
	}

	return result
}

// filterFieldsBySelection 按选定字段列表过滤数据。
func filterFieldsBySelection(mappedData map[string]interface{}, selectedFields []string) map[string]interface{} {
	if len(selectedFields) == 0 {
		return mappedData
	}

	result := make(map[string]interface{})
	for _, field := range selectedFields {
		if value, exists := mappedData[field]; exists {
			result[field] = value
		}
	}
	return result
}

// formatFieldsAsKeyValue 将字段数据格式化为 "key：value | key：value" 形式。
func formatFieldsAsKeyValue(data map[string]interface{}) string {
	if len(data) == 0 {
		return ""
	}
	var parts []string
	for k, v := range data {
		parts = append(parts, fmt.Sprintf("%s：%v", k, v))
	}
	return strings.Join(parts, " | ")
}

// applyValueTransform 根据值转换配置替换字段值。
func applyValueTransform(data map[string]interface{}, valueTransformJSON string) map[string]interface{} {
	if valueTransformJSON == "" {
		return data
	}

	var transforms map[string]map[string]string
	if err := json.Unmarshal([]byte(valueTransformJSON), &transforms); err != nil {
		return data
	}

	for field, transformMap := range transforms {
		if value, exists := data[field]; exists {
			strValue := fmt.Sprintf("%v", value)
			if newValue, ok := transformMap[strValue]; ok {
				data[field] = newValue
			}
		}
	}

	return data
}

// truncateLargeFields 在 UDP 包超限时截断大字段值，尝试将 payload 压缩到 maxSize 以内。
func truncateLargeFields(payload []byte, filteredData map[string]interface{}, selectedFields []string, maxSize int) ([]byte, error) {
	truncatedData := make(map[string]interface{})
	for k, v := range filteredData {
		truncatedData[k] = v
	}

	const truncateThreshold = 500
	for k, v := range truncatedData {
		strVal := fmt.Sprintf("%v", v)
		if len(strVal) > truncateThreshold {
			truncatedData[k] = strVal[:truncateThreshold] + "...[截断]"
		}
	}

	newPayload, err := json.Marshal(truncatedData)
	if err != nil {
		return nil, err
	}

	if len(newPayload) <= maxSize {
		return newPayload, nil
	}

	if len(selectedFields) > 0 && len(selectedFields) > 3 {
		minimalData := make(map[string]interface{})
		for i, field := range selectedFields {
			if i >= len(selectedFields)/2 {
				break
			}
			if value, exists := truncatedData[field]; exists {
				minimalData[field] = value
			}
		}
		minPayload, err := json.Marshal(minimalData)
		if err != nil {
			return nil, err
		}
		if len(minPayload) <= maxSize {
			return minPayload, nil
		}
	}

	return nil, fmt.Errorf("truncated payload still exceeds limit")
}
