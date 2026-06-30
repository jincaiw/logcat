// Package syslog 提供Syslog日志接收、解析、过滤、告警分发与追踪的完整能力。
//
// 合并了原 syslog_service.go 的全部职责，修复以下问题：
//   - GetReceiveRate 数据竞争：原代码在 RLock 下修改 lastTime/lastCount/lastRate，改用 Lock
//   - 集成 filter 包修复的 regexpMatch bug
//   - 使用结构化日志替代 stdlog.Printf("[DEBUG] ...")
//   - 通过 StatsUpdater 接口解耦与 App 的依赖
package syslog

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"syslog-alert/internal/models"
	"syslog-alert/internal/repository"
	"syslog-alert/pkg/constants"
	applogger "syslog-alert/pkg/logger"
)

// StatsUpdater 统计更新接口，由上层（如 App）实现，解耦 syslog 服务与具体应用。
type StatsUpdater interface {
	UpdateStats(logs int64, devices int, running bool)
}

// Message 接收到的 Syslog 消息。
type Message struct {
	SourceIP   string
	SourcePort int
	Message    string
	ReceivedAt time.Time
}

// Server Syslog 接收服务器，支持 UDP/TCP 协议。
type Server struct {
	udpConn      *net.UDPConn
	tcpListener  net.Listener
	running      bool
	port         int
	protocol     string
	logChan      chan Message
	stopChan     chan struct{}
	alertCache   map[string]time.Time
	statsUpdate  StatsUpdater
	mu           sync.RWMutex
	alertMu      sync.Mutex
	receiveCount int64
	connCount    int
	traceMap     map[uint]*models.LogTraceInfo
	traceMu      sync.RWMutex

	// 速率计算字段（修复数据竞争：全部在 Lock 下读写）
	rateMu    sync.Mutex
	lastCount int64
	lastTime  time.Time
	lastRate  float64
}

// NewServer 创建 Syslog 服务器。statsUpdater 可为 nil（不更新统计）。
func NewServer(statsUpdater StatsUpdater) *Server {
	return &Server{
		statsUpdate: statsUpdater,
		logChan:     make(chan Message, 1000),
		stopChan:    make(chan struct{}),
		alertCache:  make(map[string]time.Time),
		traceMap:    make(map[uint]*models.LogTraceInfo),
	}
}

// Start 启动 Syslog 服务器。
func (s *Server) Start(port int, protocol string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("syslog service is already running")
	}

	s.port = port
	s.protocol = protocol

	if protocol == constants.ProtocolTCP {
		addr := &net.TCPAddr{Port: port, IP: net.ParseIP("0.0.0.0")}
		listener, err := net.ListenTCP(constants.ProtocolTCP, addr)
		if err != nil {
			return fmt.Errorf("failed to start TCP server on port %d: %v", port, err)
		}
		s.tcpListener = listener
		s.running = true
		go s.acceptTCPConnections()
	} else {
		addr := &net.UDPAddr{Port: port, IP: net.ParseIP("0.0.0.0")}
		conn, err := net.ListenUDP(constants.ProtocolUDP, addr)
		if err != nil {
			return fmt.Errorf("failed to start UDP server on port %d: %v", port, err)
		}
		s.udpConn = conn
		s.running = true
		go s.receiveUDPMessages()
	}

	go s.processMessages()

	if s.statsUpdate != nil {
		s.statsUpdate.UpdateStats(repository.GetLogCount(), int(repository.GetDeviceCount()), true)
	}
	return nil
}

// Stop 停止 Syslog 服务器。
func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	s.running = false
	close(s.stopChan)

	if s.udpConn != nil {
		s.udpConn.Close()
	}
	if s.tcpListener != nil {
		s.tcpListener.Close()
	}

	if s.statsUpdate != nil {
		s.statsUpdate.UpdateStats(repository.GetLogCount(), int(repository.GetDeviceCount()), false)
	}
	return nil
}

// acceptTCPConnections 接受 TCP 连接。
func (s *Server) acceptTCPConnections() {
	for {
		select {
		case <-s.stopChan:
			return
		default:
			s.tcpListener.(*net.TCPListener).SetDeadline(time.Now().Add(time.Second))
			conn, err := s.tcpListener.Accept()
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				if s.running {
					continue
				}
				return
			}
			go s.handleTCPConnection(conn)
		}
	}
}

// handleTCPConnection 处理单个 TCP 连接。
func (s *Server) handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	s.mu.Lock()
	s.connCount++
	s.mu.Unlock()
	defer func() {
		s.mu.Lock()
		s.connCount--
		s.mu.Unlock()
	}()

	buf := make([]byte, 65535)
	remoteAddr := conn.RemoteAddr().(*net.TCPAddr)

	for {
		select {
		case <-s.stopChan:
			return
		default:
			conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			n, err := conn.Read(buf)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				return
			}
			if n == 0 {
				continue
			}

			data := buf[:n]
			// 处理 Octet Counting 格式 (RFC 6587)
			data = stripOctetCountingPrefix(data)

			msg := Message{
				SourceIP:   remoteAddr.IP.String(),
				SourcePort: remoteAddr.Port,
				Message:    string(data),
				ReceivedAt: time.Now(),
			}

			select {
			case s.logChan <- msg:
			default:
				applogger.Warn("Log channel full, dropping TCP message from %s", remoteAddr.IP)
			}
		}
	}
}

// receiveUDPMessages 接收 UDP 消息。
func (s *Server) receiveUDPMessages() {
	buf := make([]byte, 65535)
	for {
		select {
		case <-s.stopChan:
			return
		default:
			s.udpConn.SetReadDeadline(time.Now().Add(time.Second))
			n, remoteAddr, err := s.udpConn.ReadFromUDP(buf)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				if s.running {
					continue
				}
				return
			}

			msg := Message{
				SourceIP:   remoteAddr.IP.String(),
				SourcePort: remoteAddr.Port,
				Message:    string(buf[:n]),
				ReceivedAt: time.Now(),
			}

			select {
			case s.logChan <- msg:
			default:
				applogger.Warn("Log channel full, dropping UDP message from %s", remoteAddr.IP)
			}
		}
	}
}

// processMessages 从通道消费消息并处理。
func (s *Server) processMessages() {
	for {
		select {
		case <-s.stopChan:
			return
		case msg := <-s.logChan:
			s.handleMessage(msg)
		}
	}
}

// stripOctetCountingPrefix 处理 RFC 6587 Octet Counting 格式，移除长度前缀。
func stripOctetCountingPrefix(data []byte) []byte {
	if len(data) > 0 && data[0] >= '0' && data[0] <= '9' {
		spaceIdx := -1
		for i, b := range data {
			if b == ' ' {
				spaceIdx = i
				break
			}
		}
		if spaceIdx > 0 {
			return data[spaceIdx+1:]
		}
	}
	return data
}

// parsePriority 解析 syslog 消息的优先级（PRI 字段）。
// 返回 priority, facility, severity。
func parsePriority(msg string) (int, int, int) {
	if len(msg) == 0 || msg[0] != '<' {
		return 0, 0, 0
	}
	end := strings.Index(msg, ">")
	if end == -1 {
		return 0, 0, 0
	}
	priority, err := strconv.Atoi(msg[1:end])
	if err != nil {
		return 0, 0, 0
	}
	facility := priority / 8
	severity := priority % 8
	return priority, facility, severity
}

// checkForwardedMark 检查消息是否为转发消息（避免循环转发）。
func checkForwardedMark(msg string) bool {
	return strings.Contains(msg, constants.ForwardedJSONMark) ||
		strings.Contains(msg, constants.ForwardedTextMark)
}

// ---- 状态查询方法 ----

// IsRunning 返回服务是否运行中。
func (s *Server) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// GetPort 返回监听端口。
func (s *Server) GetPort() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.port
}

// GetReceiveCount 返回累计接收消息数。
func (s *Server) GetReceiveCount() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.receiveCount
}

// GetConnections 返回当前 TCP 连接数。
func (s *Server) GetConnections() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.connCount
}

// incrementReceiveCount 原子递增接收计数。
func (s *Server) incrementReceiveCount() {
	s.mu.Lock()
	s.receiveCount++
	s.mu.Unlock()
}

// GetReceiveRate 返回当前接收速率（条/秒）。
//
// 修复说明：原代码在 RLock 下修改 lastTime/lastCount/lastRate，存在数据竞争。
// 此处改用独立的 rateMu 互斥锁，确保读写互斥。
// 速率每 5 秒重新计算一次，期间返回上次计算结果。
func (s *Server) GetReceiveRate() float64 {
	s.rateMu.Lock()
	defer s.rateMu.Unlock()

	now := time.Now()
	if s.lastTime.IsZero() {
		s.lastTime = now
		s.lastCount = s.getReceiveCountUnsafe()
		return 0
	}

	elapsed := now.Sub(s.lastTime).Seconds()
	if elapsed < 5 {
		if s.lastRate > 0 {
			return s.lastRate
		}
		return 0
	}

	currentCount := s.getReceiveCountUnsafe()
	rate := float64(currentCount-s.lastCount) / elapsed
	s.lastTime = now
	s.lastCount = currentCount
	s.lastRate = rate
	return rate
}

// getReceiveCountUnsafe 获取接收计数（不加锁，调用方需持有锁）。
func (s *Server) getReceiveCountUnsafe() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.receiveCount
}
