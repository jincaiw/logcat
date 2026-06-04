package database

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/logcat/logcat/internal/config"
	"github.com/logcat/logcat/internal/models"

	"golang.org/x/crypto/bcrypt"
)

// Seed initializes default data if not already present (idempotent)
func Seed(cfg *config.Config) error {
	if DB == nil {
		return errors.New("database not initialized")
	}

	log.Println("Seeding default data...")

	// Seed permissions
	if err := seedPermissions(); err != nil {
		return fmt.Errorf("failed to seed permissions: %w", err)
	}

	// Seed roles
	if err := seedRoles(); err != nil {
		return fmt.Errorf("failed to seed roles: %w", err)
	}

	// Assign permissions to roles
	if err := seedRolePermissions(); err != nil {
		return fmt.Errorf("failed to seed role permissions: %w", err)
	}

	// Seed admin user
	if err := seedAdminUser(cfg); err != nil {
		return fmt.Errorf("failed to seed admin user: %w", err)
	}

	log.Println("Seed data initialized successfully")
	return nil
}

// seedPermissions creates default permissions for all API endpoints and menus
func seedPermissions() error {
	permissions := []models.Permission{
		// Auth
		{Name: "用户登录", Code: "auth:login", Type: "api", Resource: "auth", Action: "login"},
		{Name: "用户登出", Code: "auth:logout", Type: "api", Resource: "auth", Action: "logout"},
		{Name: "查看当前用户信息", Code: "auth:me", Type: "api", Resource: "auth", Action: "read"},
		{Name: "修改密码", Code: "auth:change-password", Type: "api", Resource: "auth", Action: "update"},
		{Name: "初始化管理员", Code: "auth:init-admin", Type: "api", Resource: "auth", Action: "create"},
		{Name: "刷新会话", Code: "auth:refresh", Type: "api", Resource: "auth", Action: "refresh"},

		// Users
		{Name: "查看用户列表", Code: "users:list", Type: "api", Resource: "users", Action: "read"},
		{Name: "创建用户", Code: "users:create", Type: "api", Resource: "users", Action: "create"},
		{Name: "编辑用户", Code: "users:update", Type: "api", Resource: "users", Action: "update"},
		{Name: "删除用户", Code: "users:delete", Type: "api", Resource: "users", Action: "delete"},
		{Name: "重置用户密码", Code: "users:reset-password", Type: "api", Resource: "users", Action: "reset"},
		{Name: "解锁用户", Code: "users:unlock", Type: "api", Resource: "users", Action: "unlock"},

		// Roles
		{Name: "查看角色列表", Code: "roles:list", Type: "api", Resource: "roles", Action: "read"},
		{Name: "创建角色", Code: "roles:create", Type: "api", Resource: "roles", Action: "create"},
		{Name: "编辑角色", Code: "roles:update", Type: "api", Resource: "roles", Action: "update"},
		{Name: "删除角色", Code: "roles:delete", Type: "api", Resource: "roles", Action: "delete"},
		{Name: "管理角色权限", Code: "roles:permissions", Type: "api", Resource: "roles", Action: "manage"},

		// Permissions
		{Name: "查看权限列表", Code: "permissions:list", Type: "api", Resource: "permissions", Action: "read"},

		// Devices
		{Name: "查看设备列表", Code: "devices:list", Type: "api", Resource: "devices", Action: "read"},
		{Name: "创建设备", Code: "devices:create", Type: "api", Resource: "devices", Action: "create"},
		{Name: "编辑设备", Code: "devices:update", Type: "api", Resource: "devices", Action: "update"},
		{Name: "删除设备", Code: "devices:delete", Type: "api", Resource: "devices", Action: "delete"},

		// Device Groups
		{Name: "查看设备组列表", Code: "device-groups:list", Type: "api", Resource: "device-groups", Action: "read"},
		{Name: "创建设备组", Code: "device-groups:create", Type: "api", Resource: "device-groups", Action: "create"},
		{Name: "编辑设备组", Code: "device-groups:update", Type: "api", Resource: "device-groups", Action: "update"},
		{Name: "删除设备组", Code: "device-groups:delete", Type: "api", Resource: "device-groups", Action: "delete"},

		// Device Templates
		{Name: "查看设备模板列表", Code: "device-templates:list", Type: "api", Resource: "device-templates", Action: "read"},
		{Name: "创建设备模板", Code: "device-templates:create", Type: "api", Resource: "device-templates", Action: "create"},
		{Name: "编辑设备模板", Code: "device-templates:update", Type: "api", Resource: "device-templates", Action: "update"},
		{Name: "删除设备模板", Code: "device-templates:delete", Type: "api", Resource: "device-templates", Action: "delete"},
		{Name: "应用设备模板", Code: "device-templates:apply", Type: "api", Resource: "device-templates", Action: "apply"},

		// Field Mappings
		{Name: "查看字段映射列表", Code: "field-mappings:list", Type: "api", Resource: "field-mappings", Action: "read"},
		{Name: "创建字段映射", Code: "field-mappings:create", Type: "api", Resource: "field-mappings", Action: "create"},
		{Name: "编辑字段映射", Code: "field-mappings:update", Type: "api", Resource: "field-mappings", Action: "update"},
		{Name: "删除字段映射", Code: "field-mappings:delete", Type: "api", Resource: "field-mappings", Action: "delete"},

		// Parse Templates
		{Name: "查看解析模板列表", Code: "parse-templates:list", Type: "api", Resource: "parse-templates", Action: "read"},
		{Name: "创建解析模板", Code: "parse-templates:create", Type: "api", Resource: "parse-templates", Action: "create"},
		{Name: "编辑解析模板", Code: "parse-templates:update", Type: "api", Resource: "parse-templates", Action: "update"},
		{Name: "删除解析模板", Code: "parse-templates:delete", Type: "api", Resource: "parse-templates", Action: "delete"},
		{Name: "测试解析模板", Code: "parse-templates:test", Type: "api", Resource: "parse-templates", Action: "test"},

		// Filter Policies
		{Name: "查看过滤策略列表", Code: "filter-policies:list", Type: "api", Resource: "filter-policies", Action: "read"},
		{Name: "创建过滤策略", Code: "filter-policies:create", Type: "api", Resource: "filter-policies", Action: "create"},
		{Name: "编辑过滤策略", Code: "filter-policies:update", Type: "api", Resource: "filter-policies", Action: "update"},
		{Name: "删除过滤策略", Code: "filter-policies:delete", Type: "api", Resource: "filter-policies", Action: "delete"},
		{Name: "测试过滤策略", Code: "filter-policies:test", Type: "api", Resource: "filter-policies", Action: "test"},

		// Output Templates
		{Name: "查看输出模板列表", Code: "output-templates:list", Type: "api", Resource: "output-templates", Action: "read"},
		{Name: "创建输出模板", Code: "output-templates:create", Type: "api", Resource: "output-templates", Action: "create"},
		{Name: "编辑输出模板", Code: "output-templates:update", Type: "api", Resource: "output-templates", Action: "update"},
		{Name: "删除输出模板", Code: "output-templates:delete", Type: "api", Resource: "output-templates", Action: "delete"},

		// Push Configs
		{Name: "查看推送配置列表", Code: "push-configs:list", Type: "api", Resource: "push-configs", Action: "read"},
		{Name: "创建推送配置", Code: "push-configs:create", Type: "api", Resource: "push-configs", Action: "create"},
		{Name: "编辑推送配置", Code: "push-configs:update", Type: "api", Resource: "push-configs", Action: "update"},
		{Name: "删除推送配置", Code: "push-configs:delete", Type: "api", Resource: "push-configs", Action: "delete"},
		{Name: "测试推送配置", Code: "push-configs:test", Type: "api", Resource: "push-configs", Action: "test"},

		// Alert Rules
		{Name: "查看告警规则列表", Code: "alert-rules:list", Type: "api", Resource: "alert-rules", Action: "read"},
		{Name: "创建告警规则", Code: "alert-rules:create", Type: "api", Resource: "alert-rules", Action: "create"},
		{Name: "编辑告警规则", Code: "alert-rules:update", Type: "api", Resource: "alert-rules", Action: "update"},
		{Name: "删除告警规则", Code: "alert-rules:delete", Type: "api", Resource: "alert-rules", Action: "delete"},

		// Logs
		{Name: "查看日志列表", Code: "logs:list", Type: "api", Resource: "logs", Action: "read"},
		{Name: "清理日志", Code: "logs:cleanup", Type: "api", Resource: "logs", Action: "delete"},
		{Name: "查看日志追踪", Code: "logs:trace", Type: "api", Resource: "logs", Action: "trace"},

		// Alerts
		{Name: "查看告警记录列表", Code: "alerts:list", Type: "api", Resource: "alerts", Action: "read"},
		{Name: "创建告警处置", Code: "alerts:disposition:create", Type: "api", Resource: "alerts", Action: "create"},
		{Name: "查看告警处置", Code: "alerts:disposition:list", Type: "api", Resource: "alerts", Action: "read"},

		// Aggregated Alerts
		{Name: "查看聚合告警列表", Code: "aggregated-alerts:list", Type: "api", Resource: "aggregated-alerts", Action: "read"},
		{Name: "查看聚合告警关联日志", Code: "aggregated-alerts:logs", Type: "api", Resource: "aggregated-alerts", Action: "read"},
		{Name: "更新聚合告警状态", Code: "aggregated-alerts:update", Type: "api", Resource: "aggregated-alerts", Action: "update"},

		// Stats
		{Name: "查看字段统计", Code: "stats:fields", Type: "api", Resource: "stats", Action: "read"},
		{Name: "查看可用字段", Code: "stats:available-fields", Type: "api", Resource: "stats", Action: "read"},

		// High Freq IPs
		{Name: "查看高频IP列表", Code: "high-freq-ips:list", Type: "api", Resource: "high-freq-ips", Action: "read"},
		{Name: "配置高频IP规则", Code: "high-freq-ips:config", Type: "api", Resource: "high-freq-ips", Action: "update"},

		// Desensitize Rules
		{Name: "查看脱敏规则列表", Code: "desensitize-rules:list", Type: "api", Resource: "desensitize-rules", Action: "read"},
		{Name: "创建脱敏规则", Code: "desensitize-rules:create", Type: "api", Resource: "desensitize-rules", Action: "create"},
		{Name: "编辑脱敏规则", Code: "desensitize-rules:update", Type: "api", Resource: "desensitize-rules", Action: "update"},
		{Name: "删除脱敏规则", Code: "desensitize-rules:delete", Type: "api", Resource: "desensitize-rules", Action: "delete"},

		// Import/Export
		{Name: "导入解析模板", Code: "import:parse-templates", Type: "api", Resource: "import", Action: "import"},
		{Name: "导入过滤策略", Code: "import:filter-policies", Type: "api", Resource: "import", Action: "import"},
		{Name: "导入推送配置", Code: "import:push-configs", Type: "api", Resource: "import", Action: "import"},
		{Name: "导入设备模板", Code: "import:device-templates", Type: "api", Resource: "import", Action: "import"},
		{Name: "导出配置", Code: "export:config", Type: "api", Resource: "export", Action: "export"},

		// System
		{Name: "查看系统配置", Code: "system:config:read", Type: "api", Resource: "system", Action: "read"},
		{Name: "更新系统配置", Code: "system:config:update", Type: "api", Resource: "system", Action: "update"},
		{Name: "查看系统状态", Code: "system:status", Type: "api", Resource: "system", Action: "read"},
		{Name: "启停Syslog接收", Code: "system:syslog", Type: "api", Resource: "system", Action: "manage"},

		// Audit Logs
		{Name: "查看审计日志", Code: "audit-logs:list", Type: "api", Resource: "audit-logs", Action: "read"},

		// Dashboard
		{Name: "查看仪表盘", Code: "dashboard:view", Type: "api", Resource: "dashboard", Action: "read"},

		// Menu permissions
		{Name: "仪表盘菜单", Code: "menu:dashboard", Type: "menu", Resource: "menu", Action: "view"},
		{Name: "日志管理菜单", Code: "menu:logs", Type: "menu", Resource: "menu", Action: "view"},
		{Name: "告警管理菜单", Code: "menu:alerts", Type: "menu", Resource: "menu", Action: "view"},
		{Name: "设备管理菜单", Code: "menu:devices", Type: "menu", Resource: "menu", Action: "view"},
		{Name: "策略管理菜单", Code: "menu:policies", Type: "menu", Resource: "menu", Action: "view"},
		{Name: "模板管理菜单", Code: "menu:templates", Type: "menu", Resource: "menu", Action: "view"},
		{Name: "推送管理菜单", Code: "menu:push", Type: "menu", Resource: "menu", Action: "view"},
		{Name: "用户管理菜单", Code: "menu:users", Type: "menu", Resource: "menu", Action: "view"},
		{Name: "审计日志菜单", Code: "menu:audit", Type: "menu", Resource: "menu", Action: "view"},
		{Name: "系统设置菜单", Code: "menu:system", Type: "menu", Resource: "menu", Action: "view"},
	}

	for _, p := range permissions {
		var existing models.Permission
		result := DB.Where("code = ?", p.Code).First(&existing)
		if result.Error != nil {
			if err := DB.Create(&p).Error; err != nil {
				return fmt.Errorf("failed to create permission %s: %w", p.Code, err)
			}
		}
	}

	return nil
}

// seedRoles creates the 4 built-in roles
func seedRoles() error {
	roles := []models.Role{
		{Name: "管理员", Code: "admin", Description: "系统管理员，拥有所有权限", BuiltIn: true},
		{Name: "运维人员", Code: "operator", Description: "运维操作人员，可管理设备和策略", BuiltIn: true},
		{Name: "查看者", Code: "watcher", Description: "只读用户，可查看所有数据", BuiltIn: true},
		{Name: "审计员", Code: "auditor", Description: "审计人员，可查看日志和审计记录", BuiltIn: true},
	}

	for _, r := range roles {
		var existing models.Role
		result := DB.Where("code = ?", r.Code).First(&existing)
		if result.Error != nil {
			if err := DB.Create(&r).Error; err != nil {
				return fmt.Errorf("failed to create role %s: %w", r.Code, err)
			}
		}
	}

	return nil
}

// seedRolePermissions assigns permissions to roles
func seedRolePermissions() error {
	// Admin gets ALL permissions
	var adminRole models.Role
	if err := DB.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
		return fmt.Errorf("admin role not found: %w", err)
	}

	// Check if admin already has permissions
	var count int64
	DB.Model(&models.RolePermission{}).Where("role_id = ?", adminRole.ID).Count(&count)
	if count == 0 {
		var allPermissions []models.Permission
		if err := DB.Find(&allPermissions).Error; err != nil {
			return fmt.Errorf("failed to find all permissions: %w", err)
		}
		for _, p := range allPermissions {
			rp := models.RolePermission{RoleID: adminRole.ID, PermissionID: p.ID}
			if err := DB.Create(&rp).Error; err != nil {
				return fmt.Errorf("failed to assign permission %d to admin role: %w", p.ID, err)
			}
		}
		log.Printf("Assigned %d permissions to admin role", len(allPermissions))
	}

	// Operator: device management, policy management, template management, push config, alert rules, logs/alerts read
	operatorCodes := []string{
		"auth:login", "auth:logout", "auth:me", "auth:change-password",
		"devices:list", "devices:create", "devices:update", "devices:delete",
		"device-groups:list", "device-groups:create", "device-groups:update", "device-groups:delete",
		"device-templates:list", "device-templates:create", "device-templates:update", "device-templates:delete",
		"field-mappings:list", "field-mappings:create", "field-mappings:update", "field-mappings:delete",
		"parse-templates:list", "parse-templates:create", "parse-templates:update", "parse-templates:delete", "parse-templates:test",
		"filter-policies:list", "filter-policies:create", "filter-policies:update", "filter-policies:delete", "filter-policies:test",
		"output-templates:list", "output-templates:create", "output-templates:update", "output-templates:delete",
		"push-configs:list", "push-configs:create", "push-configs:update", "push-configs:delete", "push-configs:test",
		"alert-rules:list", "alert-rules:create", "alert-rules:update", "alert-rules:delete",
		"logs:list", "logs:trace",
		"alerts:list", "alerts:disposition:list",
		"aggregated-alerts:list", "aggregated-alerts:logs", "aggregated-alerts:update",
		"stats:fields", "stats:available-fields",
		"high-freq-ips:list",
		"desensitize-rules:list", "desensitize-rules:create", "desensitize-rules:update", "desensitize-rules:delete",
		"import:parse-templates", "import:filter-policies", "import:push-configs", "import:device-templates",
		"export:config",
		"menu:dashboard", "menu:logs", "menu:alerts", "menu:devices", "menu:policies", "menu:templates", "menu:push",
	}
	assignPermissionsToRole("operator", operatorCodes)

	// Watcher: read-only access to most resources
	watcherCodes := []string{
		"auth:login", "auth:logout", "auth:me", "auth:change-password",
		"devices:list",
		"device-groups:list",
		"device-templates:list",
		"field-mappings:list",
		"parse-templates:list",
		"filter-policies:list",
		"output-templates:list",
		"push-configs:list",
		"alert-rules:list",
		"logs:list", "logs:trace",
		"alerts:list", "alerts:disposition:list",
		"aggregated-alerts:list", "aggregated-alerts:logs",
		"stats:fields", "stats:available-fields",
		"high-freq-ips:list",
		"desensitize-rules:list",
		"menu:dashboard", "menu:logs", "menu:alerts", "menu:devices", "menu:policies", "menu:templates", "menu:push",
	}
	assignPermissionsToRole("watcher", watcherCodes)

	// Auditor: audit log + log/alerts read-only
	auditorCodes := []string{
		"auth:login", "auth:logout", "auth:me", "auth:change-password",
		"logs:list", "logs:trace",
		"alerts:list", "alerts:disposition:list",
		"aggregated-alerts:list", "aggregated-alerts:logs",
		"audit-logs:list",
		"menu:dashboard", "menu:logs", "menu:alerts", "menu:audit",
	}
	assignPermissionsToRole("auditor", auditorCodes)

	return nil
}

// assignPermissionsToRole assigns a list of permission codes to a role (idempotent)
func assignPermissionsToRole(roleCode string, permCodes []string) {
	var role models.Role
	if err := DB.Where("code = ?", roleCode).First(&role).Error; err != nil {
		log.Printf("WARNING: role %s not found, skipping permission assignment", roleCode)
		return
	}

	var count int64
	DB.Model(&models.RolePermission{}).Where("role_id = ?", role.ID).Count(&count)
	if count > 0 {
		return // already has permissions
	}

	for _, code := range permCodes {
		var perm models.Permission
		if err := DB.Where("code = ?", code).First(&perm).Error; err != nil {
			log.Printf("WARNING: permission code %s not found", code)
			continue
		}
		rp := models.RolePermission{RoleID: role.ID, PermissionID: perm.ID}
		if err := DB.Create(&rp).Error; err != nil {
			log.Printf("failed to assign permission %d to role %d: %v", perm.ID, role.ID, err)
		}
	}
}

// seedAdminUser creates the default admin user
func seedAdminUser(cfg *config.Config) error {
	var count int64
	DB.Model(&models.User{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		return nil // admin already exists
	}

	password := config.GetAdminPassword()
	if password == "" {
		// Generate a random password
		b := make([]byte, 16)
		if _, err := rand.Read(b); err != nil {
			return fmt.Errorf("failed to generate random password: %w", err)
		}
		password = hex.EncodeToString(b)
		log.Printf("==================================================")
		log.Printf("  NO LOGCAT_ADMIN_PASSWORD SET!")
		log.Printf("  Generated admin password written to data/.admin_password")
		log.Printf("  Please change it after first login!")
		log.Printf("==================================================")
	} else {
		log.Printf("Admin password loaded from LOGCAT_ADMIN_PASSWORD environment variable")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	admin := models.User{
		Username:     "admin",
		PasswordHash: string(hashedPassword),
		DisplayName:  "系统管理员",
		Email:        "",
		Status:       models.UserStatusEnabled,
		MustChangePassword: true,
	}

	if err := DB.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	// Assign admin role to the admin user
	var adminRole models.Role
	if err := DB.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
		return fmt.Errorf("admin role not found: %w", err)
	}
	ur := models.UserRole{UserID: admin.ID, RoleID: adminRole.ID}
	if err := DB.Create(&ur).Error; err != nil {
		return fmt.Errorf("failed to assign admin role to admin user: %w", err)
	}

	// Also set admin password in env for reference
	_ = os.WriteFile("data/.admin_password", []byte(password+"\n"), 0600)

	return nil
}