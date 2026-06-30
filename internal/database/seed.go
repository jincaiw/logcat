package database

import (
	"encoding/json"
	"os"

	"syslog-alert/internal/models"
	"syslog-alert/pkg/constants"
	applogger "syslog-alert/pkg/logger"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// seedDefaults 初始化默认配置和默认数据
func seedDefaults() {
	initDefaultConfig()
	initDefaultDeviceGroup()
	initDefaultFieldMappingDocs()
	initDefaultUser()
}

// initDefaultUser 初始化默认用户（单用户模式）。
// 若用户表为空，则使用默认凭证创建管理员账户。
func initDefaultUser() {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}
	username := os.Getenv(constants.EnvAdminUsername)
	if username == "" {
		username = constants.AuthDefaultUsername
	}
	password := os.Getenv(constants.EnvAdminPassword)
	if password == "" {
		password = constants.AuthDefaultPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), constants.AuthBCryptCost)
	if err != nil {
		applogger.Error("生成默认密码哈希失败: %v", err)
		return
	}
	user := models.User{
		Username:     username,
		PasswordHash: string(hash),
		Nickname:     constants.AuthDefaultNickname,
	}
	if err := db.Create(&user).Error; err != nil {
		applogger.Warn("创建默认用户失败: %v", err)
	} else {
		if os.Getenv(constants.EnvAdminPassword) == "" {
			applogger.Info("已创建默认用户: %s (初始密码: %s，请登录后立即修改)", user.Username, constants.AuthDefaultPassword)
		} else {
			applogger.Info("已创建默认用户: %s (密码来自环境变量 %s)", user.Username, constants.EnvAdminPassword)
		}
	}
}

// initDefaultConfig 初始化默认系统配置
func initDefaultConfig() {
	var config models.SystemConfig
	if err := db.First(&config).Error; err == gorm.ErrRecordNotFound {
		db.Create(&models.SystemConfig{
			ListenPort:            constants.DefaultListenPort,
			LogRetention:          3,
			MaxLogSize:            constants.DefaultMaxLogSize,
			AutoStart:             false,
			MinimizeToTray:        true,
			AlertEnabled:          true,
			AlertInterval:         constants.DefaultAlertInterval,
			UnmatchedLogRetention: 3,
			UnmatchedLogAlert:     true,
			DefaultFilterAction:   constants.ActionKeep,
			Theme:                 constants.DefaultTheme,
			Language:              constants.DefaultLanguage,
		})
		applogger.Info("已初始化默认系统配置")
	}
}

// initDefaultDeviceGroup 初始化默认设备分组
func initDefaultDeviceGroup() {
	var count int64
	db.Model(&models.DeviceGroup{}).Count(&count)
	if count == 0 {
		db.Create(&models.DeviceGroup{
			Name:        constants.DefaultDeviceGroupName,
			Description: "默认设备分组",
			Color:       "#409eff",
			SortOrder:   0,
		})
		applogger.Info("已创建默认设备分组")
	}
}

// initDefaultFieldMappingDocs 初始化默认字段映射文档
// 修复原代码 Bug：天眼字段映射被创建两次（一次在云锁块内，一次在外）
func initDefaultFieldMappingDocs() {
	// 云锁字段映射
	var yunsuoCount int64
	db.Model(&models.FieldMappingDoc{}).Where("device_type = ?", "云锁").Count(&yunsuoCount)
	if yunsuoCount == 0 {
		db.Create(&models.FieldMappingDoc{
			Name:          "云锁字段映射",
			DeviceType:    "云锁",
			Description:   "云锁安全设备Syslog日志字段映射文档",
			FieldMappings: yunsuoFieldMappings,
			IsActive:      true,
		})
		applogger.Info("已创建云锁字段映射文档")
	}

	// 天眼字段映射
	var tianyanCount int64
	db.Model(&models.FieldMappingDoc{}).Where("device_type = ?", "天眼").Count(&tianyanCount)
	if tianyanCount == 0 {
		db.Create(&models.FieldMappingDoc{
			Name:          "天眼字段映射",
			DeviceType:    "天眼",
			Description:   "天眼安全设备Syslog日志字段映射文档",
			FieldMappings: tianyanFieldMappings,
			IsActive:      true,
		})
		applogger.Info("已创建天眼字段映射文档")
	}
}

// loadParseTemplatesFromFile 从 JSON 文件加载解析模板（仅在不存在时创建）
func loadParseTemplatesFromFile(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		applogger.Warn("读取解析模板文件失败: %v", err)
		return
	}

	var data struct {
		Version   string `json:"version"`
		Templates []struct {
			Name           string `json:"name"`
			Description    string `json:"description"`
			ParseType      string `json:"parseType"`
			HeaderRegex    string `json:"headerRegex"`
			FieldMapping   string `json:"fieldMapping"`
			ValueTransform string `json:"valueTransform"`
			DeviceType     string `json:"deviceType"`
			IsActive       bool   `json:"isActive"`
		} `json:"templates"`
	}

	if err := json.Unmarshal(content, &data); err != nil {
		applogger.Warn("解析模板文件 JSON 解析失败: %v", err)
		return
	}

	applogger.Info("从文件发现 %d 个解析模板", len(data.Templates))

	for _, t := range data.Templates {
		var existing models.ParseTemplate
		result := db.Where("name = ?", t.Name).First(&existing)
		if result.Error == gorm.ErrRecordNotFound {
			template := models.ParseTemplate{
				Name:           t.Name,
				Description:    t.Description,
				ParseType:      t.ParseType,
				HeaderRegex:    t.HeaderRegex,
				FieldMapping:   t.FieldMapping,
				ValueTransform: t.ValueTransform,
				DeviceType:     t.DeviceType,
				IsActive:       t.IsActive,
			}
			if err := db.Create(&template).Error; err != nil {
				applogger.Warn("创建解析模板 %s 失败: %v", t.Name, err)
			} else {
				applogger.Info("已加载解析模板: %s", t.Name)
			}
		}
	}
}

// loadFilterPoliciesFromFile 从 JSON 文件加载筛选策略（仅在不存在时创建）
func loadFilterPoliciesFromFile(filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		applogger.Warn("读取筛选策略文件失败: %v", err)
		return
	}

	var data struct {
		Version  string `json:"version"`
		Policies []struct {
			Name              string `json:"name"`
			Description       string `json:"description"`
			ParseTemplateName string `json:"parseTemplateName"`
			Conditions        string `json:"conditions"`
			ConditionLogic    string `json:"conditionLogic"`
			Action            string `json:"action"`
			Priority          int    `json:"priority"`
			IsActive          bool   `json:"isActive"`
			DedupEnabled      bool   `json:"dedupEnabled"`
			DedupWindow       int    `json:"dedupWindow"`
			DropUnmatched     bool   `json:"dropUnmatched"`
		} `json:"policies"`
	}

	if err := json.Unmarshal(content, &data); err != nil {
		applogger.Warn("筛选策略文件 JSON 解析失败: %v", err)
		return
	}

	for _, p := range data.Policies {
		var existing models.FilterPolicy
		if db.Where("name = ?", p.Name).First(&existing).Error == gorm.ErrRecordNotFound {
			// 查找关联的解析模板
			var parseTemplateID uint
			var pt models.ParseTemplate
			if db.Where("name = ?", p.ParseTemplateName).First(&pt).Error == nil {
				parseTemplateID = pt.ID
			}

			policy := models.FilterPolicy{
				Name:            p.Name,
				Description:     p.Description,
				ParseTemplateID: parseTemplateID,
				Conditions:      p.Conditions,
				ConditionLogic:  p.ConditionLogic,
				Action:          p.Action,
				Priority:        p.Priority,
				IsActive:        p.IsActive,
				DedupEnabled:    p.DedupEnabled,
				DedupWindow:     p.DedupWindow,
				DropUnmatched:   p.DropUnmatched,
			}
			if err := db.Create(&policy).Error; err != nil {
				applogger.Warn("创建筛选策略 %s 失败: %v", p.Name, err)
			} else {
				applogger.Info("已加载筛选策略: %s", p.Name)
			}
		}
	}
}
