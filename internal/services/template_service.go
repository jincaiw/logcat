package services

import (
	"fmt"
	"strings"
	"time"
)

// TemplateService handles template variable replacement
type TemplateService struct{}

// NewTemplateService creates a new TemplateService
func NewTemplateService() *TemplateService {
	return &TemplateService{}
}

// Replace replaces template variables in a string with actual values
// Supported variables: ${field_name}, ${timestamp}, ${date}, ${time}, ${source_ip}, etc.
func (s *TemplateService) Replace(template string, data map[string]interface{}) string {
	result := template

	// Replace all ${variable} patterns
	for key, val := range data {
		placeholder := fmt.Sprintf("${%s}", key)
		strVal := fmt.Sprintf("%v", val)
		result = strings.ReplaceAll(result, placeholder, strVal)
	}

	// Standard variables
	now := time.Now()
	result = strings.ReplaceAll(result, "${timestamp}", now.Format(time.RFC3339))
	result = strings.ReplaceAll(result, "${date}", now.Format("2006-01-02"))
	result = strings.ReplaceAll(result, "${time}", now.Format("15:04:05"))
	result = strings.ReplaceAll(result, "${datetime}", now.Format("2006-01-02 15:04:05"))

	return result
}

// ReplaceAll replaces variables in all values of a map
func (s *TemplateService) ReplaceAll(templates map[string]string, data map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for key, tmpl := range templates {
		result[key] = s.Replace(tmpl, data)
	}
	return result
}