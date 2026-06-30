// Package platform 封装平台相关能力：本机 IP 探测、浏览器启动等。
package platform

import (
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"syslog-alert/pkg/logger"
)

// GetLocalIP 获取本机首选 IPv4（优先物理网卡）
func GetLocalIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "127.0.0.1"
	}

	var physicalIP, otherIP string

	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		name := strings.ToLower(iface.Name)
		isVirtual := isVirtualInterface(name)
		isPhysical := isPhysicalInterface(name)

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip := ipnet.IP.String()
					if isPhysical && !isVirtual {
						if physicalIP == "" {
							physicalIP = ip
						}
					} else if !isVirtual && otherIP == "" {
						otherIP = ip
					}
				}
			}
		}
	}

	if physicalIP != "" {
		return physicalIP
	}
	if otherIP != "" {
		return otherIP
	}
	return "127.0.0.1"
}

// GetLocalIPs 获取所有非回环 IPv4 地址
func GetLocalIPs() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return []string{"127.0.0.1"}
	}

	var ips []string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}

	if len(ips) == 0 {
		ips = append(ips, "127.0.0.1")
	}
	return ips
}

// OpenBrowser 跨平台打开浏览器
func OpenBrowser(url string) {
	// 延迟一点，确保服务已监听
	time.Sleep(500 * time.Millisecond)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	default:
		return
	}
	if err := cmd.Start(); err != nil {
		logger.Warn("打开浏览器失败: %v", err)
	}
}

// isVirtualInterface 判断是否为虚拟网卡
func isVirtualInterface(name string) bool {
	virtualPrefixes := []string{"vnic", "vmnet", "bridge", "veth", "docker", "vEthernet", "parallels", "utun", "awdl", "llw", "anpi"}
	for _, p := range virtualPrefixes {
		if strings.Contains(name, p) || strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}

// isPhysicalInterface 判断是否为物理网卡
func isPhysicalInterface(name string) bool {
	physicalPrefixes := []string{"en", "eth", "wlan", "wi-fi"}
	for _, p := range physicalPrefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}
