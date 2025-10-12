package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Siddharth-Nema/tcp-scanner/internal/scan"
)

func main() {
	host := flag.String("host", "", "Host to scan (required)")
	port := flag.Int("port", 80, "Port to check")
	timeout := flag.Duration("timeout", 1*time.Second, "Connection timeout (e.g., 1s, 500ms)")
	flag.Parse()

	if *host == "" {
		fmt.Fprintln(os.Stderr, "Error: --host is required")
		flag.Usage()
		os.Exit(2)
	}

	fmt.Printf("Checking %s:%d ...\n", *host, *port)

	open, err := scan.CheckPort(context.Background(), *host, *port, *timeout)
	if err != nil {
		fmt.Printf("Port %d on %s is closed or unreachable (%v)\n", *port, *host, err)
		return
	}
	if open {
		fmt.Printf("Port %d on %s is OPEN\n", *port, *host)
	} else {
		fmt.Printf("Port %d on %s is CLOSED\n", *port, *host)
	}
}
