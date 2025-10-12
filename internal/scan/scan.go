package scan

import (
	"context"
	"fmt"
	"net"
	"time"
)

// CheckPort tries to open a TCP connection to host:port with a timeout.
// Returns (true, nil) if connection succeeded. If it failed, returns (false, err).
// ctx is respected (cancellable).
func CheckPort(ctx context.Context, host string, port int, timeout time.Duration) (bool, error) {
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	// Create a child context with timeout so Dial respects both ctx and timeout.
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	d := net.Dialer{}
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return false, err
	}
	_ = conn.Close()
	return true, nil
}
