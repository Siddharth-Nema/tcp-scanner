package scan

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

// Result is the outcome for a single port check.
type Result struct {
	Port     int
	Open     bool
	Err      error
	Duration time.Duration
}

// worker polls port numbers from `ports` and sends Result into `results`.
func worker(ctx context.Context, host string, ports <-chan int, results chan<- Result, wg *sync.WaitGroup, timeout time.Duration) {
	defer wg.Done()
	for p := range ports {
		start := time.Now()

		// create a per-dial context with timeout derived from parent ctx
		dialCtx, cancel := context.WithTimeout(ctx, timeout)
		open, err := CheckPort(dialCtx, host, p, timeout)
		cancel()

		results <- Result{
			Port:     p,
			Open:     open,
			Err:      err,
			Duration: time.Since(start),
		}
	}
}

// ScanRangeConcurrent scans ports in [start..end] using `workers` concurrent workers.
// It shows a progress bar during the scan and returns a slice of Results sorted by Port.
func ScanRangeConcurrent(ctx context.Context, host string, start, end, workers int, timeout time.Duration) []Result {
	if workers <= 0 {
		workers = 100
	}
	if start < 1 {
		start = 1
	}
	if end < start {
		end = start
	}

	totalPorts := end - start + 1

	ports := make(chan int, workers)
	results := make(chan Result, totalPorts)

	var wg sync.WaitGroup

	// start progress bar
	bar := progressbar.NewOptions(totalPorts,
		progressbar.OptionSetDescription(fmt.Sprintf("Scanning %s", host)),
		progressbar.OptionSetWidth(20),
		progressbar.OptionShowCount(),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetPredictTime(false),
	)

	// start worker goroutines
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker(ctx, host, ports, results, &wg, timeout)
	}

	// feed ports
	go func() {
		for p := start; p <= end; p++ {
			select {
			case <-ctx.Done():
				break
			default:
				ports <- p
			}
		}
		close(ports)
	}()

	// wait for workers and close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// collect results with progress updates
	out := make([]Result, 0, totalPorts)
	for r := range results {
		_ = bar.Add(1)
		out = append(out, r)
	}

	// sort by port for deterministic output
	sort.Slice(out, func(i, j int) bool { return out[i].Port < out[j].Port })

	return out
}

// FilterOpen returns only the open results from a slice.
func FilterOpen(results []Result) []Result {
	out := make([]Result, 0, len(results))
	for _, r := range results {
		if r.Open && r.Err == nil {
			out = append(out, r)
		}
	}
	return out
}
