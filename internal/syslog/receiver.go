package syslog

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/logcat/logcat/internal/config"
)

// Receiver handles UDP/TCP syslog message reception
// This is a placeholder for Phase 4 when syslog receiving is fully implemented
type Receiver struct {
	udpListener *net.UDPConn
	tcpListener net.Listener
	running     bool
	stopCh      chan struct{}
	wg          sync.WaitGroup
	mu          sync.Mutex
	udpPort     int
	tcpPort     int
}

// ReceivedMessage represents a received syslog message
type ReceivedMessage struct {
	Data      []byte
	SourceIP  string
	Protocol  string // "udp" or "tcp"
	Timestamp time.Time
}

// messageChan is the channel for received messages
var messageChan chan ReceivedMessage

// SetMessageChannel sets the channel for received messages
func SetMessageChannel(ch chan ReceivedMessage) {
	messageChan = ch
}

// NewReceiver creates a new syslog Receiver
func NewReceiver(udpPort, tcpPort int) *Receiver {
	return &Receiver{
		udpPort: udpPort,
		tcpPort: tcpPort,
		stopCh:  make(chan struct{}),
	}
}

// Start begins listening for syslog messages on UDP and TCP
func (r *Receiver) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.running {
		return nil
	}

	cfg := config.Get()
	if cfg == nil || !cfg.Syslog.Enabled {
		log.Println("Syslog receiver is disabled in configuration")
		return nil
	}

	r.running = true

	// Start UDP listener
	if r.udpPort > 0 {
		addr := &net.UDPAddr{Port: r.udpPort, IP: net.ParseIP("0.0.0.0")}
		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			log.Printf("WARNING: Failed to start UDP syslog listener on port %d: %v", r.udpPort, err)
		} else {
			r.udpListener = conn
			r.wg.Add(1)
			go r.handleUDP()
			log.Printf("Syslog UDP listener started on port %d", r.udpPort)
		}
	}

	// Start TCP listener
	if r.tcpPort > 0 {
		addr := &net.TCPAddr{Port: r.tcpPort, IP: net.ParseIP("0.0.0.0")}
		listener, err := net.ListenTCP("tcp", addr)
		if err != nil {
			log.Printf("WARNING: Failed to start TCP syslog listener on port %d: %v", r.tcpPort, err)
		} else {
			r.tcpListener = listener
			r.wg.Add(1)
			go r.handleTCP()
			log.Printf("Syslog TCP listener started on port %d", r.tcpPort)
		}
	}

	return nil
}

// Stop shuts down the syslog receiver
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
	}
	if r.tcpListener != nil {
		r.tcpListener.Close()
	}

	r.wg.Wait()
	log.Println("Syslog receiver stopped")
}

// IsRunning returns whether the receiver is running
func (r *Receiver) IsRunning() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.running
}

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
				log.Printf("UDP read error: %v", err)
			}
			continue
		}

		msg := ReceivedMessage{
			Data:      make([]byte, n),
			SourceIP:  remoteAddr.IP.String(),
			Protocol:  "udp",
			Timestamp: time.Now(),
		}
		copy(msg.Data, buffer[:n])

		// Send to processing channel if available
		if messageChan != nil {
			select {
			case messageChan <- msg:
			default:
				log.Println("WARNING: Message channel full, dropping message")
			}
		}
	}
}

func (r *Receiver) handleTCP() {
	defer r.wg.Done()

	for {
		select {
		case <-r.stopCh:
			return
		default:
		}

		// Set accept deadline for graceful shutdown
		if tcpListener, ok := r.tcpListener.(*net.TCPListener); ok {
			tcpListener.SetDeadline(time.Now().Add(1 * time.Second))
		}

		conn, err := r.tcpListener.Accept()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			if r.running {
				log.Printf("TCP accept error: %v", err)
			}
			continue
		}

		r.wg.Add(1)
		go r.handleTCPConnection(conn)
	}
}

func (r *Receiver) handleTCPConnection(conn net.Conn) {
	defer r.wg.Done()
	defer conn.Close()

	buffer := make([]byte, 65535)
	remoteAddr := conn.RemoteAddr().String()
	host, _, _ := net.SplitHostPort(remoteAddr)

	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	n, err := conn.Read(buffer)
	if err != nil {
		if r.running {
			log.Printf("TCP read error from %s: %v", host, err)
		}
		return
	}

	msg := ReceivedMessage{
		Data:      make([]byte, n),
		SourceIP:  host,
		Protocol:  "tcp",
		Timestamp: time.Now(),
	}
	copy(msg.Data, buffer[:n])

	if messageChan != nil {
		select {
		case messageChan <- msg:
		default:
			log.Println("WARNING: Message channel full, dropping message")
		}
	}
}