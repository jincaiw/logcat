package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/logcat/logcat/internal/config"
	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/engine"
	"github.com/logcat/logcat/internal/handlers"
	"github.com/logcat/logcat/internal/middleware"
	"github.com/logcat/logcat/internal/services"
	logsyslog "github.com/logcat/logcat/internal/syslog"
)

func main() {
	// Parse command line flags
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Printf("Configuration loaded: server=%s:%d, database=%s",
		cfg.Server.Host, cfg.Server.Port, cfg.Database.Type)

	// Initialize database
	if err := database.Initialize(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Run auto-migration
	if cfg.Database.AutoMigrate {
		if err := database.AutoMigrate(); err != nil {
			log.Fatalf("Failed to auto-migrate database: %v", err)
		}
	}

	// Seed default data
	if err := database.Seed(cfg); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	// Initialize services
	authService := services.NewAuthService()
	userService := services.NewUserService()
	deviceService := services.NewDeviceService()
	parseService := services.NewParseService()
	filterService := services.NewFilterService()
	pushService := services.NewPushService()
	alertService := services.NewAlertService()
	traceService := services.NewTraceService()
	aggregateService := services.NewAggregateService()
	desensitizeService := services.NewDesensitizeService()
	statsService := services.NewStatsService()
	auditService := services.NewAuditService()
	cleanupService := services.NewCleanupService()

	// Start periodic cleanup
	cleanupService.Start(1 * time.Hour)

	// Initialize syslog receiver
	syslogReceiver := logsyslog.NewReceiver(cfg.Syslog.UDPPort, cfg.Syslog.TCPPort)
	if cfg.Syslog.Enabled {
		msgChan := make(chan logsyslog.ReceivedMessage, cfg.Queue.Capacity)
		logsyslog.SetMessageChannel(msgChan)
		if err := syslogReceiver.Start(); err != nil {
			log.Printf("WARNING: Failed to start syslog receiver: %v", err)
		}
	}

	// Initialize pipeline (placeholder for Phase 4)
	pipeline := engine.NewPipeline(
		cfg.Worker.ParseWorkers,
		cfg.Worker.FilterWorkers,
		cfg.Worker.PushWorkers,
	)
	pipeline.Start()

	// Set up Gin router
	if cfg.Theme == "dark" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.RequestID())
	router.Use(middleware.CORS())

	// Public health endpoints
	router.GET("/healthz", handlers.Healthz)
	router.GET("/readyz", handlers.Readyz)
	router.GET("/metrics", handlers.Metrics)

	// API routes
	api := router.Group("/api")

	// Permission middleware factory
	requirePerm := func(code string) gin.HandlerFunc {
		return middleware.RequirePermission(code)
	}

	// Register all routes
	handlers.RegisterRoutes(api, authService, requirePerm)
	handlers.RegisterUserRoutes(api, userService, requirePerm)
	handlers.RegisterRoleRoutes(api, userService, requirePerm)
	handlers.RegisterDeviceRoutes(api, deviceService, requirePerm)
	handlers.RegisterDeviceTemplateRoutes(api, requirePerm)
	handlers.RegisterFieldMappingRoutes(api, requirePerm)
	handlers.RegisterParseTemplateRoutes(api, parseService, requirePerm)
	handlers.RegisterFilterPolicyRoutes(api, filterService, requirePerm)
	handlers.RegisterOutputTemplateRoutes(api, requirePerm)
	handlers.RegisterPushConfigRoutes(api, pushService, requirePerm)
	handlers.RegisterAlertRuleRoutes(api, requirePerm)
	handlers.RegisterLogRoutes(api, traceService, cleanupService, requirePerm)
	handlers.RegisterAlertRecordRoutes(api, alertService, requirePerm)
	handlers.RegisterAggregatedAlertRoutes(api, aggregateService, requirePerm)
	handlers.RegisterHighFreqIPRoutes(api, services.GetHighFreqService(), requirePerm)
	handlers.RegisterDesensitizeRoutes(api, desensitizeService, requirePerm)
	handlers.RegisterStatsRoutes(api, statsService, requirePerm)
	handlers.RegisterDashboardRoutes(api, statsService, requirePerm)
	handlers.RegisterImportExportRoutes(api, requirePerm)
	handlers.RegisterSystemRoutes(api, statsService, requirePerm)
	handlers.RegisterAuditLogRoutes(api, auditService, requirePerm)

	// Runtime metrics
	router.GET("/api/metrics/runtime", handlers.RuntimeMetrics)

	// Serve frontend static files (embedded in Phase 5)
	// router.NoRoute(func(c *gin.Context) {
	//     c.File("web/dist/index.html")
	// })

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("logcat server starting on %s", addr)
		log.Printf("Health check: http://%s/healthz", addr)
		log.Printf("API base: http://%s/api", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Stop syslog receiver
	syslogReceiver.Stop()

	// Stop pipeline
	pipeline.Stop()

	// Stop cleanup service
	cleanupService.Stop()

	// Shutdown HTTP server with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}