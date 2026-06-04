package syslog

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/logcat/logcat/internal/config"
	"github.com/logcat/logcat/internal/models"
)

// ParsedLog is a fully parsed syslog message ready for the pipeline.
type ParsedLog struct {
	SyslogLog *models.SyslogLog
	SourceIP  string
	Protocol  string
}

// ReceiverMetrics exposes operational counters for the receiver.
type ReceiverMetrics struct {
	UDPReceived    int64 `json:"udpReceived"`
	TCPReceived    int64 `json:"tcpReceived"`
	UDPErrors      int64 `json:"udpErrors"`
	TCPErrors      int64 `json:"tcpErrors"`
	ParseErrors    int64 `json:"parseErrors"`
	ChannelDropped int64 `json:"channelDropped"`
	TCPConnections int64 `json:"tcpConnections"`
	LastReceiveAt  int64 `json:"lastReceiveAt"` // unix nano
}

// Receiver handles UDP/TCP syslog message reception with full RFC parsing.
type Receiver struct {
	udpListener *net.UDPConn
	tcpListener net.Listener
	running     bool
	stopCh      chan struct{}
	wg          sync.WaitGroup
	mu          sync.Mutex
	udpPort     int
	tcpPort     int

	metrics ReceiverMetrics

	// parsedCh is the output channel that delivers *ParsedLog to the pipeline.
	parsedCh chan<- *ParsedLog
}

// NewReceiver creates a new syslog Receiver that sends parsed logs to ch.
func NewReceiver(udpPort, tcpPort int, parsedCh chan<- *ParsedLog) *Receiver {
	return &Receiver{
		udpPort:  udpPort,
		tcpPort:  tcpPort,
		stopCh:   make(chan struct{}),
		parsedCh: parsedCh,
	}
}

// Start begins listening for syslog messages on UDP and TCP.
func (r *Receiver) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		return nil
	}

	cfg := config.Get()
	if cfg == nil || !cfg.Syslog.Enabled {
		log.Println("[syslog] Receiver is disabled in configuration")
		return nil
	}

	r.stopCh = make(chan struct{})
	r.running = true
	started := 0
	var startErrs []error

	// Start UDP listener
	if r.udpPort > 0 {
		addr := &net.UDPAddr{Port: r.udpPort, IP: net.ParseIP("0.0.0.0")}
		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			startErrs = append(startErrs, fmt.Errorf("udp %d: %w", r.udpPort, err))
			log.Printf("[syslog] WARNING: Failed to start UDP listener on port %d: %v", r.udpPort, err)
		} else {
			r.udpListener = conn
			started++
			r.wg.Add(1)
			go r.handleUDP()
			log.Printf("[syslog] UDP listener started on port %d", r.udpPort)
		}
	}

	// Start TCP listener
	if r.tcpPort > 0 {
		addr := &net.TCPAddr{Port: r.tcpPort, IP: net.ParseIP("0.0.0.0")}
		listener, err := net.ListenTCP("tcp", addr)
		if err != nil {
			startErrs = append(startErrs, fmt.Errorf("tcp %d: %w", r.tcpPort, err))
			log.Printf("[syslog] WARNING: Failed to start TCP listener on port %d: %v", r.tcpPort, err)
		} else {
			r.tcpListener = listener
			started++
			r.wg.Add(1)
			go r.handleTCP()
			log.Printf("[syslog] TCP listener started on port %d", r.tcpPort)
		}
	}

	if started == 0 && (r.udpPort > 0 || r.tcpPort > 0) {
		r.running = false
		return fmt.Errorf("failed to start syslog receiver: %w", errors.Join(startErrs...))
	}

	return nil
}

// Stop shuts down the syslog receiver gracefully.
func (r *Receiver) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.running {
		return
	}

	r.running = false
	close(r.stopCh)

	if r.udpListener != nil {
		r.udpListener.Close()
		r.udpListener = nil
	}
	if r.tcpListener != nil {
		r.tcpListener.Close()
		r.tcpListener = nil
	}

	r.wg.Wait()
	log.Println("[syslog] Receiver stopped")
}

// IsRunning returns whether the receiver is running.
func (r *Receiver) IsRunning() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.running
}

// Metrics returns a snapshot of receiver metrics.
func (r *Receiver) Metrics() ReceiverMetrics {
	return ReceiverMetrics{
		UDPReceived:    atomic.LoadInt64(&r.metrics.UDPReceived),
		TCPReceived:    atomic.LoadInt64(&r.metrics.TCPReceived),
		UDPErrors:      atomic.LoadInt64(&r.metrics.UDPErrors),
		TCPErrors:      atomic.LoadInt64(&r.metrics.TCPErrors),
		ParseErrors:    atomic.LoadInt64(&r.metrics.ParseErrors),
		ChannelDropped: atomic.LoadInt64(&r.metrics.ChannelDropped),
		TCPConnections: atomic.LoadInt64(&r.metrics.TCPConnections),
		LastReceiveAt:  atomic.LoadInt64(&r.metrics.LastReceiveAt),
	}
}

// ---------------------------------------------------------------------------
// UDP handler
// ---------------------------------------------------------------------------

func (r *Receiver) handleUDP() {
	defer r.wg.Done()

	buffer := make([]byte, 65535)
	for {
		select {
		case <-r.stopCh:
			return
		default:
		}

		r.udpListener.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, remoteAddr, err := r.udpListener.ReadFromUDP(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			if r.running {
				atomic.AddInt64(&r.metrics.UDPErrors, 1)
				log.Printf("[syslog] UDP read error: %v", err)
			}
			continue
		}

		atomic.AddInt64(&r.metrics.UDPReceived, 1)
		atomic.StoreInt64(&r.metrics.LastReceiveAt, time.Now().UnixNano())

		raw := string(buffer[:n])
		parsed := r.parseMessage(raw, remoteAddr.IP.String(), "udp")
		if parsed != nil {
			r.deliver(parsed)
		}
	}
}

// ---------------------------------------------------------------------------
// TCP handler (line-delimited, RFC 6587 octet-framing)
// ---------------------------------------------------------------------------

func (r *Receiver) handleTCP() {
	defer r.wg.Done()

	for {
		select {
		case <-r.stopCh:
			return
		default:
		}

		if tcpListener, ok := r.tcpListener.(*net.TCPListener); ok {
			tcpListener.SetDeadline(time.Now().Add(1 * time.Second))
		}

		conn, err := r.tcpListener.Accept()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			if r.running {
				atomic.AddInt64(&r.metrics.TCPErrors, 1)
				log.Printf("[syslog] TCP accept error: %v", err)
			}
			continue
		}

		atomic.AddInt64(&r.metrics.TCPConnections, 1)
		r.wg.Add(1)
		go r.handleTCPConnection(conn)
	}
}

func (r *Receiver) handleTCPConnection(conn net.Conn) {

	defer r.wg.Done()
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()
	host, _, _ := net.SplitHostPort(remoteAddr)

	scanner := bufio.NewScanner(conn)
	// 64KB buffer per line
	scanner.Buffer(make([]byte, 0, 65535), 65535)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		atomic.AddInt64(&r.metrics.TCPReceived, 1)
		atomic.StoreInt64(&r.metrics.LastReceiveAt, time.Now().UnixNano())

		parsed := r.parseMessage(line, host, "tcp")
		if parsed != nil {
			r.deliver(parsed)
		}
	}

	if err := scanner.Err(); err != nil && r.running {
		atomic.AddInt64(&r.metrics.TCPErrors, 1)
		log.Printf("[syslog] TCP read error from %s: %v", host, err)
	}
}

// ---------------------------------------------------------------------------
// Message parsing
// ---------------------------------------------------------------------------

// priorityToFacilitySeverity decodes a PRI value (facility*8 + severity).
func priorityToFacilitySeverity(pri int) (facility string, severity string) {
	facNum := pri / 8
	sevNum := pri % 8

	facNames := []string{
		"kern", "user", "mail", "daemon", "auth", "syslog",
		"lpr", "news", "uucp", "cron", "authpriv", "ftp",
		"ntp", "audit", "alert", "clock", "local0", "local1",
		"local2", "local3", "local4", "local5", "local6", "local7",
	}
	sevNames := []string{
		"emerg", "alert", "crit", "err", "warning",
		"notice", "info", "debug",
	}

	if facNum >= 0 && facNum < len(facNames) {
		facility = facNames[facNum]
	} else {
		facility = fmt.Sprintf("facility%d", facNum)
	}

	if sevNum >= 0 && sevNum < len(sevNames) {
		severity = sevNames[sevNum]
	} else {
		severity = fmt.Sprintf("severity%d", sevNum)
	}

	return
}

// monthNameToNumber converts abbreviated month to number (Jan=1).
var months = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4,
	"may": 5, "jun": 6, "jul": 7, "aug": 8,
	"sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

// parseMessage takes a raw syslog line and produces a ParsedLog.
// Supports RFC3164 (<PRI>MMM DD HH:MM:SS HOSTNAME MSG) and
// RFC5424 (<PRI>VERSION TIMESTAMP HOSTNAME APP-NAME PROCID MSGID SD MSG) formats.
func (r *Receiver) parseMessage(raw, sourceIP, protocol string) *ParsedLog {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	var (
		facility  string
		severity  string
		hostname  string
		timestamp time.Time
	)

	// ---- Extract PRI ----
	if !strings.HasPrefix(raw, "<") {
		atomic.AddInt64(&r.metrics.ParseErrors, 1)
		// Still create a log entry with raw message
		return &ParsedLog{
			SyslogLog: &models.SyslogLog{
				LogID:      uuid.New().String(),
				SourceIP:   sourceIP,
				RawMessage: raw,
				Severity:   "unknown",
				Facility:   "unknown",
				ReceivedAt: time.Now(),
			},
			SourceIP: sourceIP,
			Protocol: protocol,
		}
	}

	priEnd := strings.IndexByte(raw, '>')
	if priEnd < 2 {
		atomic.AddInt64(&r.metrics.ParseErrors, 1)
		return &ParsedLog{
			SyslogLog: &models.SyslogLog{
				LogID:      uuid.New().String(),
				SourceIP:   sourceIP,
				RawMessage: raw,
				Severity:   "unknown",
				Facility:   "unknown",
				ReceivedAt: time.Now(),
			},
			SourceIP: sourceIP,
			Protocol: protocol,
		}
	}

	priStr := raw[1:priEnd]
	pri, err := strconv.Atoi(priStr)
	if err != nil {
		atomic.AddInt64(&r.metrics.ParseErrors, 1)
		return &ParsedLog{
			SyslogLog: &models.SyslogLog{
				LogID:      uuid.New().String(),
				SourceIP:   sourceIP,
				RawMessage: raw,
				Severity:   "unknown",
				Facility:   "unknown",
				ReceivedAt: time.Now(),
			},
			SourceIP: sourceIP,
			Protocol: protocol,
		}
	}
	facility, severity = priorityToFacilitySeverity(pri)

	rest := strings.TrimSpace(raw[priEnd+1:])

	// ---- Try RFC5424 format: VERSION TIMESTAMP... ----
	if rfc5424, ok := tryRFC5424(rest, sourceIP); ok {
		rfc5424.SyslogLog.LogID = uuid.New().String()
		rfc5424.SyslogLog.SourceIP = sourceIP
		rfc5424.SyslogLog.Facility = facility
		rfc5424.SyslogLog.Severity = severity
		rfc5424.SyslogLog.ReceivedAt = time.Now()
		rfc5424.SourceIP = sourceIP
		rfc5424.Protocol = protocol
		return rfc5424
	}

	// ---- Try RFC3164 format: MMM DD HH:MM:SS HOSTNAME MSG ----
	timestamp, hostname, _ = parseRFC3164Header(rest)
	if hostname == "" {
		hostname = sourceIP
	}

	return &ParsedLog{
		SyslogLog: &models.SyslogLog{
			LogID:      uuid.New().String(),
			SourceIP:   sourceIP,
			RawMessage: raw,
			Facility:   facility,
			Severity:   severity,
			DeviceName: hostname,
			ReceivedAt: timestamp,
		},
		SourceIP: sourceIP,
		Protocol: protocol,
	}
}

// tryRFC5424 attempts to parse the rest of the message after PRI as RFC5424.
// RFC5424: VERSION SP TIMESTAMP SP HOSTNAME SP APP-NAME SP PROCID SP MSGID SP [SD] SP MSG
func tryRFC5424(rest, sourceIP string) (*ParsedLog, bool) {
	// Split on space to get version field
	parts := strings.SplitN(rest, " ", 8)
	if len(parts) < 7 {
		return nil, false
	}

	versionStr := parts[0]
	if _, err := strconv.Atoi(versionStr); err != nil {
		// Version must be numeric
		return nil, false
	}

	tsStr := parts[1]
	hostname := parts[2]

	// Parse timestamp (RFC3339 or the special 'NILVALUE' = "-")
	var ts time.Time
	if tsStr != "" && tsStr != "-" {
		var err error
		ts, err = time.Parse(time.RFC3339Nano, tsStr)
		if err != nil {
			ts, err = time.Parse(time.RFC3339, tsStr)
			if err != nil {
				// Try ISO-like format without TZ
				ts, err = time.Parse("2006-01-02T15:04:05", tsStr)
				if err != nil {
					ts = time.Now()
				}
			}
		}
	} else {
		ts = time.Now()
	}

	// Successfully parsed as RFC5424

	if hostname == "" || hostname == "-" {
		hostname = sourceIP
	}

	return &ParsedLog{
		SyslogLog: &models.SyslogLog{
			RawMessage: rest,
			DeviceName: hostname,
			ReceivedAt: ts,
		},
	}, true
}

// parseRFC3164Header parses the RFC3164 header portion: "MMM DD HH:MM:SS HOSTNAME MSG"
func parseRFC3164Header(rest string) (timestamp time.Time, hostname, body string) {
	now := time.Now()
	year := now.Year()

	fields := strings.Fields(rest)
	if len(fields) < 4 {
		return now, "", rest
	}

	// Month: 3-letter abbreviation
	monthNum, ok := months[strings.ToLower(fields[0])]
	if !ok {
		return now, "", rest
	}

	// Day
	day, err := strconv.Atoi(fields[1])
	if err != nil {
		return now, "", rest
	}

	// Time (HH:MM:SS)
	timeParts := strings.Split(fields[2], ":")
	if len(timeParts) != 3 {
		return now, "", rest
	}
	hour, _ := strconv.Atoi(timeParts[0])
	min, _ := strconv.Atoi(timeParts[1])
	sec, _ := strconv.Atoi(timeParts[2])

	timestamp = time.Date(year, time.Month(monthNum), day, hour, min, sec, 0, time.UTC)

	// Determine hostname and body
	// RFC3164: the first word after timestamp is hostname
	bodyStart := 3
	if len(fields) > 3 {
		hostname = fields[3]
		bodyStart = 4
	}

	if len(fields) > bodyStart {
		body = strings.Join(fields[bodyStart:], " ")
	} else {
		body = rest
	}

	return
}

// deliver sends a parsed log to the pipeline. If the channel is full it drops
// the message and records the drop.
func (r *Receiver) deliver(parsed *ParsedLog) {
	if r.parsedCh == nil {
		return
	}
	select {
	case r.parsedCh <- parsed:
		// delivered
	default:
		atomic.AddInt64(&r.metrics.ChannelDropped, 1)
		log.Printf("[syslog] WARNING: pipeline input channel full, dropping message from %s", parsed.SourceIP)
	}
}

// globalReceiver holds the singleton receiver instance for API access.
var globalReceiver *Receiver

// SetGlobalReceiver stores the receiver for API access.
func SetGlobalReceiver(r *Receiver) {
	globalReceiver = r
}

// GetGlobalReceiver returns the singleton receiver, or nil if not set.
func GetGlobalReceiver() *Receiver {
	return globalReceiver
}
