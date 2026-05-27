package services

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
)

// DesensitizeService handles field desensitization
type DesensitizeService struct{}

// NewDesensitizeService creates a new DesensitizeService
func NewDesensitizeService() *DesensitizeService {
	return &DesensitizeService{}
}

// Desensitize applies all enabled desensitization rules to the data
func (s *DesensitizeService) Desensitize(data map[string]interface{}) (map[string]interface{}, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var rules []models.DesensitizeRule
	if err := db.Where("enabled = ?", true).Find(&rules).Error; err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for k, v := range data {
		result[k] = v
	}

	for _, rule := range rules {
		if val, exists := result[rule.FieldName]; exists {
			result[rule.FieldName] = s.applyRule(rule.RuleType, fmt.Sprintf("%v", val))
		}
	}

	return result, nil
}

func (s *DesensitizeService) applyRule(ruleType, value string) string {
	switch ruleType {
	case "ip":
		return s.desensitizeIP(value)
	case "account":
		return s.desensitizeAccount(value)
	case "token":
		return s.desensitizeToken(value)
	case "url":
		return s.desensitizeURL(value)
	case "email":
		return s.desensitizeEmail(value)
	case "phone":
		return s.desensitizePhone(value)
	case "regex":
		return s.desensitizeRegex(value)
	default:
		return value
	}
}

func (s *DesensitizeService) desensitizeIP(ip string) string {
	parts := strings.Split(ip, ".")
	if len(parts) == 4 {
		return parts[0] + "." + parts[1] + ".*.*"
	}
	return ip
}

func (s *DesensitizeService) desensitizeAccount(account string) string {
	if len(account) <= 3 {
		return "***"
	}
	return account[:3] + strings.Repeat("*", len(account)-3)
}

func (s *DesensitizeService) desensitizeToken(token string) string {
	if len(token) <= 8 {
		return strings.Repeat("*", len(token))
	}
	return token[:4] + strings.Repeat("*", len(token)-8) + token[len(token)-4:]
}

func (s *DesensitizeService) desensitizeURL(url string) string {
	return "https://***"
}

func (s *DesensitizeService) desensitizeEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "***@***"
	}
	local := parts[0]
	if len(local) <= 2 {
		local = "***"
	} else {
		local = local[:2] + strings.Repeat("*", len(local)-2)
	}
	return local + "@" + parts[1]
}

func (s *DesensitizeService) desensitizePhone(phone string) string {
	if len(phone) < 7 {
		return strings.Repeat("*", len(phone))
	}
	return phone[:3] + strings.Repeat("*", len(phone)-7) + phone[len(phone)-4:]
}

func (s *DesensitizeService) desensitizeRegex(value string) string {
	// Mask all alphanumeric characters
	re := regexp.MustCompile(`[a-zA-Z0-9]`)
	return re.ReplaceAllString(value, "*")
}

// ListRules returns all desensitization rules
func (s *DesensitizeService) ListRules() ([]models.DesensitizeRule, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var rules []models.DesensitizeRule
	if err := db.Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

// CreateRule creates a new desensitization rule
func (s *DesensitizeService) CreateRule(rule *models.DesensitizeRule) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	return db.Create(rule).Error
}

// UpdateRule updates an existing desensitization rule
func (s *DesensitizeService) UpdateRule(id uint, updates map[string]interface{}) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	var rule models.DesensitizeRule
	if err := db.First(&rule, id).Error; err != nil {
		return err
	}

	return db.Model(&rule).Updates(updates).Error
}

// DeleteRule deletes a desensitization rule
func (s *DesensitizeService) DeleteRule(id uint) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	return db.Delete(&models.DesensitizeRule{}, id).Error
}