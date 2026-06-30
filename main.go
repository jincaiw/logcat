// logcat 入口：Web 服务器模式。
//
// 启动流程：
//  1. 初始化配置与数据库
//  2. 创建 WebServer（内含 Syslog 服务器）
//  3. 注册 API 路由与静态文件服务
//  4. 启动 HTTP 服务器并等待信号优雅关闭
package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"syslog-alert/internal/api"
	"syslog-alert/internal/config"
	"syslog-alert/internal/database"
	"syslog-alert/internal/platform"
	"syslog-alert/internal/repository"
	"syslog-alert/pkg/constants"
	applogger "syslog-alert/pkg/logger"
)

//go:embed all:frontend/dist
var webAssets embed.FS

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		printUsage()
		return
	}

	// 初始化配置与数据库
	config.Init()
	database.Get()

	// 创建 Web 服务器
	webServer := api.NewWebServer()

	// 注册路由
	mux := api.NewRouter(webServer)

	// 静态文件服务
	setupStaticFiles(mux)

	// 如配置了自动启动，则启动 Syslog 服务
	if cfg := repository.GetSystemConfig(); cfg.AutoStart {
		port := cfg.ListenPort
		if port == 0 {
			port = constants.DefaultListenPort
		}
		proto := cfg.Protocol
		if proto == "" {
			proto = "udp"
		}
		if err := webServer.SyslogServer().Start(port, proto); err != nil {
			applogger.Warn("自动启动 Syslog 服务失败: %v", err)
		} else {
			applogger.Info("Syslog 服务已自动启动 (端口: %d, 协议: %s)", port, proto)
		}
	}

	// 启动 HTTP 服务器
	port := parsePort()
	addr := fmt.Sprintf("0.0.0.0:%d", port)
	server := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	ip := platform.GetLocalIP()
	url := fmt.Sprintf("http://%s:%d", ip, port)
	printBanner(url)

	if shouldOpenBrowser() {
		go platform.OpenBrowser(url)
	}

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- server.ListenAndServe()
	}()

	select {
	case <-quit:
		fmt.Println("\n正在停止服务...")
		webServer.SyslogServer().Stop()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		server.Shutdown(ctx)
		cancel()
		fmt.Println("服务已停止")
	case err := <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			applogger.Error("服务器错误: %v", err)
		}
	}
}

// setupStaticFiles 设置前端静态文件服务。
func setupStaticFiles(mux *http.ServeMux) {
	distFS, err := fs.Sub(webAssets, "frontend/dist")
	if err != nil {
		applogger.Warn("前端静态文件目录未找到: %v", err)
		return
	}
	fileServer := http.FileServer(http.FS(distFS))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// API 请求走 API 路由（已注册），其他请求回退到前端 SPA
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		// SPA 回退：不存在的路径返回 index.html
		path := r.URL.Path
		if _, err := fs.Stat(distFS, strings.TrimPrefix(path, "/")); err != nil {
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})
}

// parsePort 从命令行参数解析端口号，默认 8080。
func parsePort() int {
	portArg := ""
	if len(os.Args) > 1 {
		if os.Args[1] == "-p" && len(os.Args) > 2 {
			portArg = os.Args[2]
		} else if !strings.HasPrefix(os.Args[1], "-") {
			portArg = os.Args[1]
		}
	}
	if portArg == "" {
		return constants.DefaultWebPort
	}
	port, err := strconv.Atoi(portArg)
	if err != nil || port < 1 || port > 65535 {
		applogger.Warn("无效 Web 端口 %q，使用默认端口 %d", portArg, constants.DefaultWebPort)
		return constants.DefaultWebPort
	}
	return port
}

func shouldOpenBrowser() bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(constants.EnvOpenBrowser)))
	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	}
	return runtime.GOOS == "darwin" || runtime.GOOS == "windows"
}

func printUsage() {
	fmt.Printf("logcat Web Server v%s\n", constants.AppVersion)
	fmt.Println()
	fmt.Println("用法:")
	fmt.Printf("  logcat [端口号]     启动 Web 服务器（默认端口 %d）\n", constants.DefaultWebPort)
	fmt.Println("  logcat -p <端口号>  指定端口启动")
	fmt.Println("  logcat --help       显示帮助信息")
	fmt.Printf("环境变量 %s=1 可在 Linux/容器中启动后自动打开浏览器\n", constants.EnvOpenBrowser)
	fmt.Println()
	fmt.Printf("默认端口: %d\n", constants.DefaultWebPort)
}

func printBanner(url string) {
	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Printf("║  %s v%-33s║\n", constants.AppName, constants.AppVersion)
	fmt.Println("╠══════════════════════════════════════════╣")
	fmt.Printf("║  访问地址: %-31s║\n", url)
	fmt.Println("║  按 Ctrl+C 停止服务                      ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()
}
