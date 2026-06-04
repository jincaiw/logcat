package main

import (
	"context"
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

func main() {
	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Printf("Configuration loaded: server=%s:%d, database=%s",
		cfg.Server.Host, cfg.Server.Port, cfg.Database.Type)

	if err := database.Initialize(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	if cfg.Database.AutoMigrate {
		if err := database.AutoMigrate(); err != nil {
			log.Fatalf("Failed to auto-migrate database: %v", err)
		}
	}

	if err := database.Seed(cfg); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

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

	cleanupService.Start(1 * time.Hour)

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
	engine.SetGlobalPipeline(pipeline)
	pipeline.Start()

	syslogReceiver := logsyslog.NewReceiver(
		cfg.Syslog.UDPPort,
		cfg.Syslog.TCPPort,
		pipeline.RawChannel(),
	)
	logsyslog.SetGlobalReceiver(syslogReceiver)
	if cfg.Syslog.Enabled {
		if err := syslogReceiver.Start(); err != nil {
			log.Printf("WARNING: Failed to start syslog receiver: %v", err)
		}
	}

	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.Use(middleware.RequestID())
	router.Use(middleware.CORS())
	router.Use(middleware.MaxBodySize(10 * 1024 * 1024))

	serveFrontend(router)

	router.GET("/healthz", handlers.Healthz)
	router.GET("/readyz", handlers.Readyz)
	router.GET("/metrics", handlers.Metrics)

	api := router.Group("/api")
	api.Use(middleware.APIRateLimit())

	requirePerm := func(code string) gin.HandlerFunc {
		return middleware.RequirePermission(code)
	}

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

	api.GET("/metrics/runtime", handlers.RuntimeMetrics)
	api.GET("/metrics", handlers.Metrics)
	api.GET("/healthz", handlers.Healthz)
	api.GET("/readyz", handlers.Readyz)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		log.Printf("logcat server starting on %s", addr)
		log.Printf("Health check: http://%s/healthz", addr)
		log.Printf("API base: http://%s/api", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	syslogReceiver.Stop()
	pipeline.Stop()
	cleanupService.Stop()
	dedupService.Stop()
	highFreqService.Stop()
	middleware.DefaultSessionStore.Stop()
	middleware.StopRateLimiters()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

func serveFrontend(router *gin.Engine) {
	subFS, err := fs.Sub(webAssets, "web/dist")
	if err != nil {
		log.Printf("WARNING: Frontend static files not embedded (web/dist not found): %v", err)
		return
	}

	assetsFS, err := fs.Sub(subFS, "assets")
	if err == nil {
		router.GET("/assets/*filepath", func(c *gin.Context) {
			filePath := c.Param("filepath")
			c.FileFromFS("/"+filePath, http.FS(assetsFS))
		})
	}

	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		if strings.HasPrefix(path, "/api/") ||
			path == "/healthz" || path == "/readyz" || path == "/metrics" {
			c.Status(http.StatusNotFound)
			return
		}

		if strings.HasPrefix(path, "/assets/") {
			c.Status(http.StatusNotFound)
			return
		}

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
