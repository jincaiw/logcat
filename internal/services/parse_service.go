package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
)

// ParseService handles log parsing logic
type ParseService struct{}

// NewParseService creates a new ParseService
func NewParseService() *ParseService {
	return &ParseService{}
}

// ParseResult holds the result of parsing
type ParseResult struct {
	Success bool                   `json:"success"`
	Fields  map[string]interface{} `json:"fields"`
	Error   string                 `json:"error,omitempty"`
}

// TestParse tests a parse template against sample log data
func (s *ParseService) TestParse(templateID uint, rawLog string) (*ParseResult, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var template models.ParseTemplate
	if err := db.First(&template, templateID).Error; err != nil {
		return nil, err
	}

	return s.parse(&template, rawLog), nil
}

// ParseByTemplateID parses a raw log using a specific template
func (s *ParseService) ParseByTemplateID(templateID uint, rawLog string) (*ParseResult, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var template models.ParseTemplate
	if err := db.First(&template, templateID).Error; err != nil {
		return nil, err
	}

	return s.parse(&template, rawLog), nil
}

// parse performs the actual parsing based on the template
func (s *ParseService) parse(template *models.ParseTemplate, rawLog string) *ParseResult {
	fields := make(map[string]interface{})
	fields["raw_message"] = rawLog

	switch template.ParseType {
	case "json":
		return s.parseJSON(rawLog, template)
	case "delimiter":
		return s.parseDelimiter(rawLog, template)
	case "kv":
		return s.parseKV(rawLog, template)
	case "regex":
		return s.parseRegex(rawLog, template)
	case "syslog_json":
		return s.parseSyslogJSON(rawLog, template)
	case "sub_template":
		return s.parseSubTemplate(rawLog, template)
	default:
		return &ParseResult{
			Success: false,
			Fields:  fields,
			Error:   fmt.Sprintf("unsupported parse type: %s", template.ParseType),
		}
	}
}

func (s *ParseService) parseJSON(rawLog string, template *models.ParseTemplate) *ParseResult {
	fields := make(map[string]interface{})
	fields["raw_message"] = rawLog

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(rawLog), &data); err != nil {
		return &ParseResult{Success: false, Fields: fields, Error: fmt.Sprintf("JSON parse failed: %v", err)}
	}

	// Apply field mapping if configured
	if template.FieldMapping != "" {
		var mapping map[string]string
		if err := json.Unmarshal([]byte(template.FieldMapping), &mapping); err == nil {
			for srcField, dstField := range mapping {
				if val, ok := data[srcField]; ok {
					fields[dstField] = val
				}
			}
		}
	} else {
		for k, v := range data {
			fields[k] = v
		}
	}

	return &ParseResult{Success: true, Fields: fields}
}

func (s *ParseService) parseDelimiter(rawLog string, template *models.ParseTemplate) *ParseResult {
	fields := make(map[string]interface{})
	fields["raw_message"] = rawLog

	if template.Delimiter == "" {
		return &ParseResult{Success: false, Fields: fields, Error: "delimiter is not set"}
	}

	parts := strings.Split(rawLog, template.Delimiter)
	if template.FieldMapping != "" {
		var mapping []string
		if err := json.Unmarshal([]byte(template.FieldMapping), &mapping); err == nil {
			for i, fieldName := range mapping {
				if i < len(parts) {
					fields[fieldName] = strings.TrimSpace(parts[i])
				}
			}
		}
	} else {
		for i, part := range parts {
			fields[fmt.Sprintf("field_%d", i+1)] = strings.TrimSpace(part)
		}
	}

	return &ParseResult{Success: true, Fields: fields}
}

func (s *ParseService) parseKV(rawLog string, template *models.ParseTemplate) *ParseResult {
	fields := make(map[string]interface{})
	fields["raw_message"] = rawLog

	if template.Delimiter == "" {
		return &ParseResult{Success: false, Fields: fields, Error: "delimiter is not set"}
	}

	pairs := strings.Split(rawLog, template.Delimiter)
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		idx := strings.Index(pair, "=")
		if idx > 0 {
			key := strings.TrimSpace(pair[:idx])
			value := strings.TrimSpace(pair[idx+1:])
			// Remove quotes if present
			value = strings.Trim(value, "\"'")
			fields[key] = value
		}
	}

	return &ParseResult{Success: true, Fields: fields}
}

func (s *ParseService) parseRegex(rawLog string, template *models.ParseTemplate) *ParseResult {
	fields := make(map[string]interface{})
	fields["raw_message"] = rawLog

	if template.HeaderRegex == "" {
		return &ParseResult{Success: false, Fields: fields, Error: "header regex is not set"}
	}

	// Parse regex and field mappings
	if template.FieldMapping != "" {
		// Store the regex pattern for later processing
		fields["_regex_pattern"] = template.HeaderRegex
		fields["_parsed"] = rawLog
		fields["_status"] = "regex pattern stored, full processing in pipeline"
		return &ParseResult{Success: true, Fields: fields}
	}

	return &ParseResult{Success: true, Fields: fields}
}

func (s *ParseService) parseSyslogJSON(rawLog string, template *models.ParseTemplate) *ParseResult {
	// Parse as syslog header + JSON body
	fields := make(map[string]interface{})
	fields["raw_message"] = rawLog

	// Find JSON body - look for first '{'
	jsonStart := strings.Index(rawLog, "{")
	if jsonStart >= 0 {
		// Header part
		header := strings.TrimSpace(rawLog[:jsonStart])
		fields["syslog_header"] = header

		// JSON body
		body := rawLog[jsonStart:]
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(body), &data); err != nil {
			return &ParseResult{Success: false, Fields: fields, Error: fmt.Sprintf("JSON body parse failed: %v", err)}
		}
		for k, v := range data {
			fields[k] = v
		}
	}

	return &ParseResult{Success: true, Fields: fields}
}

func (s *ParseService) parseSubTemplate(rawLog string, template *models.ParseTemplate) *ParseResult {
	fields := make(map[string]interface{})
	fields["raw_message"] = rawLog

	if template.SubTemplates != "" {
		fields["_sub_templates"] = template.SubTemplates
		fields["_status"] = "sub_templates stored for pipeline processing"
		return &ParseResult{Success: true, Fields: fields}
	}

	return &ParseResult{Success: true, Fields: fields}
}