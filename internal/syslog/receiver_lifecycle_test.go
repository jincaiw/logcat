package syslog

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/logcat/logcat/internal/config"
)

func TestReceiverCanRestartAfterStop(t *testing.T) {
	port := freeTCPPort(t)
	restore := config.SetForTest(&config.Config{
		Syslog: config.SyslogConfig{
			Enabled: true,
			TCPPort: port,
		},
	})
	defer restore()

	ch := make(chan *ParsedLog, 1)
	receiver := NewReceiver(0, port, ch)

	if err := receiver.Start(); err != nil {
		t.Fatalf("first start failed: %v", err)
	}
	if !receiver.IsRunning() {
		t.Fatal("receiver should be running after first start")
	}
	receiver.Stop()
	if receiver.IsRunning() {
		t.Fatal("receiver should be stopped after Stop")
	}

	if err := receiver.Start(); err != nil {
		t.Fatalf("second start failed: %v", err)
	}
	defer receiver.Stop()
	if !receiver.IsRunning() {
		t.Fatal("receiver should be running after restart")
	}

	conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", fmt.Sprintf("%d", port)), time.Second)
	if err != nil {
		t.Fatalf("failed to connect to restarted tcp listener: %v", err)
	}
	if _, err := conn.Write([]byte("<34>Oct 11 22:14:15 host app: restart smoke\n")); err != nil {
		_ = conn.Close()
		t.Fatalf("failed to write smoke syslog: %v", err)
	}
	_ = conn.Close()
}

func freeTCPPort(t *testing.T) int {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to allocate tcp port: %v", err)
	}
	defer listener.Close()

	return listener.Addr().(*net.TCPAddr).Port
}
