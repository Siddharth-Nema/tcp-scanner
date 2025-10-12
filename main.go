package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Siddharth-Nema/tcp-scanner/internal/scan" // replace with your module path
)

func main() {
	host := flag.String("host", "", "Host to scan (required)")
	start := flag.Int("start", 1, "Start port")
	end := flag.Int("end", 1024, "End port")
	workers := flag.Int("workers", 100, "Number of concurrent workers")
	timeout := flag.Duration("timeout", 1*time.Second, "Per-connection timeout (e.g., 1s, 500ms)")
	flag.Parse()

	if *host == "" {
		fmt.Fprintln(os.Stderr, "Error: --host is required")
		flag.Usage()
		os.Exit(2)
	}

	fmt.Printf("Scanning %s ports %d–%d with %d workers (timeout=%s)...\n", *host, *start, *end, *workers, timeout.String())

	ctx := context.Background()
	results := scan.ScanRangeConcurrent(ctx, *host, *start, *end, *workers, *timeout)

	open := scan.FilterOpen(results)
	for _, r := range open {
		fmt.Printf("Port %d: OPEN (rtt=%s)\n", r.Port, r.Duration)
	}

	fmt.Printf("\nScan complete. %d open port(s) found.\n", len(open))
}
