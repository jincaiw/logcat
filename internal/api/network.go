package api

import (
	"net/http"
	"syslog-alert/internal/platform"
)

// GetLocalIPs 返回本机全部非回环 IPv4 地址。
func (ws *WebServer) GetLocalIPs(w http.ResponseWriter, r *http.Request) {
	ips := platform.GetLocalIPs()
	JSONResponse(w, ips)
}

// GetServerIP 返回本机首选 IPv4 地址。
func (ws *WebServer) GetServerIP(w http.ResponseWriter, r *http.Request) {
	ip := platform.GetLocalIP()
	JSONResponse(w, map[string]string{"ip": ip})
}
