package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

//go:embed web/dist/*
var webAssets embed.FS

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
	dedupService := services.NewDedupService()
	pushService := services.NewPushService()
	emailService := services.NewEmailService()
	syslogForwardService := services.NewSyslogForwardService()
	alertService := services.NewAlertService()
	traceService := services.NewTraceService()
	aggregateService := services.NewAggregateService()
	highFreqService := services.NewHighFreqService()
	desensitizeService := services.NewDesensitizeService()
	statsService := services.NewStatsService()
	auditService := services.NewAuditService()
	cleanupService := services.NewCleanupService()

	// Start periodic cleanup
	cleanupService.Start(1 * time.Hour)

	// Initialize pipeline
	pipeline := engine.NewPipeline(cfg, &engine.PipelineServices{
		DeviceService:        deviceService,
		ParseService:         parseService,
		FilterService:        filterService,
		DedupService:         dedupService,
		AggregateService:     aggregateService,
		HighFreqService:      highFreqService,
		DesensitizeService:   desensitizeService,
		PushService:          pushService,
		EmailService:         emailService,
		SyslogForwardService: syslogForwardService,
		AlertService:         alertService,
		TraceService:         traceService,
	})
	pipeline.Start()

	// Initialize syslog receiver -- wires directly into pipeline input channel
	syslogReceiver := logsyslog.NewReceiver(
		cfg.Syslog.UDPPort,
		cfg.Syslog.TCPPort,
		pipeline.RawChannel(),
	)
	if cfg.Syslog.Enabled {
		if err := syslogReceiver.Start(); err != nil {
			log.Printf("WARNING: Failed to start syslog receiver: %v", err)
		}
	}

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

	// --- Frontend static file serving ---
	serveFrontend(router)

	// --- Public health endpoints ---
	router.GET("/healthz", handlers.Healthz)
	router.GET("/readyz", handlers.Readyz)
	router.GET("/metrics", handlers.Metrics)

	// --- API routes ---
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
	handlers.RegisterHighFreqIPRoutes(api, highFreqService, requirePerm)
	handlers.RegisterDesensitizeRoutes(api, desensitizeService, requirePerm)
	handlers.RegisterStatsRoutes(api, statsService, requirePerm)
	handlers.RegisterDashboardRoutes(api, statsService, requirePerm)
	handlers.RegisterImportExportRoutes(api, requirePerm)
	handlers.RegisterSystemRoutes(api, statsService, requirePerm)
	handlers.RegisterAuditLogRoutes(api, auditService, requirePerm)

	// Runtime metrics
	router.GET("/api/metrics/runtime", handlers.RuntimeMetrics)

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

	// ==============================================================
	// Graceful shutdown
	// ==============================================================
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Step 1: Stop syslog receiver first (stop accepting new logs)
	syslogReceiver.Stop()

	// Step 2: Drain pipeline (process remaining in-flight items, flush DB)
	pipeline.Stop()

	// Step 3: Stop cleanup service
	cleanupService.Stop()

	// Step 4: Shutdown HTTP server with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

// serveFrontend configures the router to serve the embedded frontend
// SPA. It serves /assets/* from the embedded filesystem and falls back
// to index.html for all other paths (SPA client-side routing).
func serveFrontend(router *gin.Engine) {
	// Try to strip the "web/dist" prefix from the embedded filesystem
	subFS, err := fs.Sub(webAssets, "web/dist")
	if err != nil {
		log.Printf("WARNING: Frontend static files not embedded (web/dist not found): %v", err)
		return
	}

	// Serve /assets/* for static assets (JS, CSS, images, etc.)
	router.StaticFS("/assets", http.FS(subFS))

	// SPA fallback: serve index.html for any non-API, non-asset route
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Don't intercept API calls or health/metrics endpoints
		if strings.HasPrefix(path, "/api/") ||
			path == "/healthz" || path == "/readyz" || path == "/metrics" {
			c.Status(http.StatusNotFound)
			return
		}

		// Don't handle /assets/ routes here (already handled by StaticFS)
		if strings.HasPrefix(path, "/assets/") {
			c.Status(http.StatusNotFound)
			return
		}

		// Serve index.html for SPA routing
		data, err := subFS.Open("index.html")
		if err != nil {
			c.String(http.StatusNotFound, "Frontend not available. Run 'cd web && npm run build' to build the frontend.")
			return
		}
		defer data.Close()

		stat, err := data.Stat()
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to stat index.html")
			return
		}

		http.ServeContent(c.Writer, c.Request, "index.html", stat.ModTime(), data.(io.ReadSeeker))
	})
}

// Ensure io.ReadSeeker type-check passes
var _ io.ReadSeeker