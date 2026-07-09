package cache

import (
	"fmt"
	"sync"

	"syslog-alert/internal/models"
	"syslog-alert/internal/service/parser"
)

type sliceCache[T any] struct {
	mu     sync.RWMutex
	loaded bool
	items  []T
}

func (c *sliceCache[T]) get(loader func() []T) []T {
	c.mu.RLock()
	if c.loaded {
		items := cloneSlice(c.items)
		c.mu.RUnlock()
		return items
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.loaded {
		c.items = cloneSlice(loader())
		c.loaded = true
	}
	return cloneSlice(c.items)
}

func (c *sliceCache[T]) invalidate() {
	c.mu.Lock()
	c.loaded = false
	c.items = nil
	c.mu.Unlock()
}

func cloneSlice[T any](src []T) []T {
	if len(src) == 0 {
		return make([]T, 0)
	}
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}

var (
	devicesCache         sliceCache[models.Device]
	filterPoliciesCache  sliceCache[models.FilterPolicy]
	parseTemplatesCache  sliceCache[models.ParseTemplate]
	outputTemplatesCache sliceCache[models.OutputTemplate]
	robotsCache          sliceCache[models.Robot]
	alertRulesCache      sliceCache[models.AlertRule]

	systemConfigMu sync.RWMutex
	systemConfig   models.SystemConfig
	systemConfigOK bool

	parserMu    sync.RWMutex
	parserCache = make(map[uint]parser.Parser)
)

// ---- 通用失效 ----

func InvalidateDevices()        { devicesCache.invalidate() }
func InvalidateFilterPolicies() { filterPoliciesCache.invalidate() }
func InvalidateParseTemplates() {
	parseTemplatesCache.invalidate()
	InvalidateParsers()
}
func InvalidateOutputTemplates() { outputTemplatesCache.invalidate() }
func InvalidateRobots()          { robotsCache.invalidate() }
func InvalidateAlertRules()      { alertRulesCache.invalidate() }
func InvalidateSystemConfig() {
	systemConfigMu.Lock()
	systemConfigOK = false
	systemConfig = models.SystemConfig{}
	systemConfigMu.Unlock()
}

func InvalidateAll() {
	InvalidateDevices()
	InvalidateFilterPolicies()
	InvalidateParseTemplates()
	InvalidateOutputTemplates()
	InvalidateRobots()
	InvalidateAlertRules()
	InvalidateSystemConfig()
}

// ---- 系统配置 ----

func GetSystemConfig(loader func() models.SystemConfig) models.SystemConfig {
	systemConfigMu.RLock()
	if systemConfigOK {
		cfg := systemConfig
		systemConfigMu.RUnlock()
		return cfg
	}
	systemConfigMu.RUnlock()

	systemConfigMu.Lock()
	defer systemConfigMu.Unlock()
	if !systemConfigOK {
		systemConfig = loader()
		systemConfigOK = true
	}
	return systemConfig
}

// ---- 设备 ----

func GetDevices(loader func() []models.Device) []models.Device {
	return devicesCache.get(loader)
}

func GetDeviceByID(loader func() []models.Device, id uint) (*models.Device, error) {
	devices := GetDevices(loader)
	for i := range devices {
		if devices[i].ID == id {
			dev := devices[i]
			return &dev, nil
		}
	}
	return nil, fmt.Errorf("device not found: %d", id)
}

func GetDeviceByIP(loader func() []models.Device, ip string) (*models.Device, error) {
	devices := GetDevices(loader)
	for i := range devices {
		if devices[i].IPAddress == ip {
			dev := devices[i]
			return &dev, nil
		}
	}
	return nil, fmt.Errorf("device not found: %s", ip)
}

// ---- 解析模板 / Parser ----

func GetParseTemplates(loader func() []models.ParseTemplate) []models.ParseTemplate {
	return parseTemplatesCache.get(loader)
}

func GetParseTemplateByID(loader func() []models.ParseTemplate, id uint) (*models.ParseTemplate, error) {
	templates := GetParseTemplates(loader)
	for i := range templates {
		if templates[i].ID == id {
			item := templates[i]
			return &item, nil
		}
	}
	return nil, fmt.Errorf("parse template not found: %d", id)
}

func GetParserByTemplateID(loader func() []models.ParseTemplate, id uint) (parser.Parser, error) {
	parserMu.RLock()
	if p, ok := parserCache[id]; ok {
		parserMu.RUnlock()
		return p, nil
	}
	parserMu.RUnlock()

	template, err := GetParseTemplateByID(loader, id)
	if err != nil {
		return nil, err
	}
	p, err := parser.New(template)
	if err != nil {
		return nil, err
	}

	parserMu.Lock()
	parserCache[id] = p
	parserMu.Unlock()
	return p, nil
}

func InvalidateParsers() {
	parserMu.Lock()
	parserCache = make(map[uint]parser.Parser)
	parserMu.Unlock()
}

// ---- 筛选策略 ----

func GetFilterPolicies(loader func() []models.FilterPolicy) []models.FilterPolicy {
	return filterPoliciesCache.get(loader)
}

func GetFilterPoliciesByDeviceID(loader func() []models.FilterPolicy, deviceID uint) []models.FilterPolicy {
	policies := GetFilterPolicies(loader)
	result := make([]models.FilterPolicy, 0, len(policies))
	for _, policy := range policies {
		if policy.DeviceID == 0 || policy.DeviceID == deviceID {
			result = append(result, policy)
		}
	}
	return result
}

func GetFilterPoliciesByDeviceGroupID(loader func() []models.FilterPolicy, groupID uint) []models.FilterPolicy {
	policies := GetFilterPolicies(loader)
	result := make([]models.FilterPolicy, 0, len(policies))
	for _, policy := range policies {
		if policy.DeviceGroupID == 0 || policy.DeviceGroupID == groupID {
			result = append(result, policy)
		}
	}
	return result
}

func GetFilterPolicyByID(loader func() []models.FilterPolicy, id uint) (*models.FilterPolicy, error) {
	policies := GetFilterPolicies(loader)
	for i := range policies {
		if policies[i].ID == id {
			item := policies[i]
			return &item, nil
		}
	}
	return nil, fmt.Errorf("filter policy not found: %d", id)
}

// ---- 输出模板 ----

func GetOutputTemplates(loader func() []models.OutputTemplate) []models.OutputTemplate {
	return outputTemplatesCache.get(loader)
}

func GetOutputTemplateByID(loader func() []models.OutputTemplate, id uint) (*models.OutputTemplate, error) {
	templates := GetOutputTemplates(loader)
	for i := range templates {
		if templates[i].ID == id {
			item := templates[i]
			return &item, nil
		}
	}
	return nil, fmt.Errorf("output template not found: %d", id)
}

func GetOutputTemplateByPlatform(loader func() []models.OutputTemplate, platform string) (*models.OutputTemplate, error) {
	templates := GetOutputTemplates(loader)
	for i := range templates {
		item := templates[i]
		if item.Platform == platform && item.IsActive {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("output template not found for platform: %s", platform)
}

// ---- 机器人 ----

func GetRobots(loader func() []models.Robot) []models.Robot {
	return robotsCache.get(loader)
}

func GetRobotByID(loader func() []models.Robot, id uint) (*models.Robot, error) {
	robots := GetRobots(loader)
	for i := range robots {
		if robots[i].ID == id {
			item := robots[i]
			return &item, nil
		}
	}
	return nil, fmt.Errorf("robot not found: %d", id)
}

// ---- 告警规则 ----

func GetAlertRules(loader func() []models.AlertRule) []models.AlertRule {
	return alertRulesCache.get(loader)
}

func GetAlertRulesByFilterPolicyID(loader func() []models.AlertRule, filterPolicyID uint) []models.AlertRule {
	rules := GetAlertRules(loader)
	result := make([]models.AlertRule, 0, len(rules))
	for _, rule := range rules {
		if rule.FilterPolicyID == filterPolicyID && rule.IsActive {
			result = append(result, rule)
		}
	}
	return result
}

func GetAlertRulesByRobotID(loader func() []models.AlertRule, robotID uint) []models.AlertRule {
	rules := GetAlertRules(loader)
	result := make([]models.AlertRule, 0, len(rules))
	for _, rule := range rules {
		if rule.RobotID == robotID && rule.IsActive {
			result = append(result, rule)
		}
	}
	return result
}
