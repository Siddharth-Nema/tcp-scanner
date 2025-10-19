package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Siddharth-Nema/tcp-scanner/internal/scan"
)

func main() {
	host := flag.String("host", "", "Host to scan (required)")
	start := flag.Int("start", 1, "Start port")
	end := flag.Int("end", 1024, "End port")
	workers := flag.Int("workers", 100, "Number of concurrent workers")
	timeout := flag.Duration("timeout", 1*time.Second, "Per-connection timeout (e.g., 1s, 500ms)")
	format := flag.String("format", "text", "Output format: text|json")
	flag.Parse()

	if *host == "" {
		fmt.Fprintln(os.Stderr, "Error: --host is required")
		flag.Usage()
		os.Exit(2)
	}

	startTime := time.Now()
	ctx := context.Background()

	results := scan.ScanRangeConcurrent(ctx, *host, *start, *end, *workers, *timeout)
	open := scan.FilterOpen(results)
	elapsed := time.Since(startTime)

	switch *format {
	case "json":
		out, _ := json.MarshalIndent(open, "", "  ")
		fmt.Println(string(out))
	default:
		fmt.Printf("\nOpen Ports for %s:\n", *host)
		if len(open) == 0 {
			fmt.Println("None found.")
		}
		for _, r := range open {
			fmt.Printf("Port %-5d : OPEN (rtt=%s)\n", r.Port, r.Duration)
		}
	}

	fmt.Printf("\nScan complete: %d open port(s) in %.2f seconds.\n",
		len(open), elapsed.Seconds())
}
