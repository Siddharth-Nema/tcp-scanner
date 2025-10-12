package scan

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestCheckPort_OpenAndClosed(t *testing.T) {
	ctx := context.Background()
	timeout := 1 * time.Second

	// 1) Start a listener on an OS-assigned free port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	defer ln.Close()

	// extract port from listener address
	_, portStr, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatalf("invalid listener addr: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("invalid port: %v", err)
	}

	// 2) The port should be reported as open
	t.Run("open", func(t *testing.T) {
		open, err := CheckPort(ctx, "127.0.0.1", port, timeout)
		if err != nil {
			t.Fatalf("expected no error checking open port, got: %v", err)
		}
		if !open {
			t.Fatalf("expected port %d to be open", port)
		}
	})

	// 3) Close the listener and the port should be reported closed
	if err := ln.Close(); err != nil {
		t.Fatalf("failed to close listener: %v", err)
	}

	t.Run("closed", func(t *testing.T) {
		open, err := CheckPort(ctx, "127.0.0.1", port, timeout)
		// We expect either an error or open==false (Dial will usually return an error)
		if err == nil && open {
			t.Fatalf("expected port %d to be closed after listener closed", port)
		}
	})
}
