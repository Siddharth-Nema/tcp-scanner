# tcp-scanner

A small TCP port scanner written in Go. This repository currently contains a simple, sequential scanner that tests a range of TCP ports on a target host and reports which ports are open.

This project is intentionally small and educational — it demonstrates basic network dialing with timeouts and use of Go's `context` for cancellation and deadlines. Concurrency and richer output formats are potential future additions.

# tcp-scanner

A small TCP port scanner written in Go. This repository currently contains a simple, sequential scanner that tests a range of TCP ports on a target host and reports which ports are open.

This project is intentionally small and educational — it demonstrates basic network dialing with timeouts and use of Go's `context` for cancellation and deadlines. Concurrency and richer output formats are potential future additions.

## Features (current)

- Sequential port scanning using `net.Dialer` + `DialContext`
- Per-port timeout (configurable)
- Simple CLI flags for host and port range
- Cross-platform (Windows / macOS / Linux)

## Project layout (current)

```
tcp-scanner/
├─ main.go                   # CLI entrypoint (placed at project root)
├─ internal/scan/scan.go     # Single-port check logic
├─ pkg/output/               # (reserved for output/formatting helpers)
├─ go.mod
└─ README.md
```

> Note: the repository places `main.go` at the project root (not under `cmd/`) — the README reflects that current layout.

## Installation

Clone the repository and build the tool:

```powershell
git clone https://github.com/Siddharth-Nema/tcp-scanner.git
cd tcp-scanner
go mod tidy
go build -o tcp-scanner ./
```

Run without building:

```powershell
go run ./ --host=scanme.nmap.org
```

## Usage

Flags supported by the CLI (see `main.go`):

- `--host` (string) — target hostname or IP (required)
- `--start` (int) — starting port (default: 1)
- `--end` (int) — ending port (default: 1024)
- `--workers` (int) — number of concurrent workers (default: 100)
- `--timeout` (duration) — per-connection timeout like `1s` or `500ms` (default: `1s`)
- `--format` (string) — output format: `text` or `json` (default: `text`)

Examples:

```powershell
# scan ports 1..1024 on scanme.nmap.org (sequential-ish with workers)
go run ./ --host=scanme.nmap.org --start=1 --end=1024 --workers=100 --timeout=1s

# build and run the binary and request JSON output
./tcp-scanner --host=127.0.0.1 --start=20 --end=1024 --workers=50 --timeout=500ms --format=json
```

Behavior notes:

- The program uses a per-connection timeout for each port. The current scanner is sequential, so scanning large ranges may take a long time.
- The current implementation uses `context.Background()` at the top level; it does not install an OS signal handler to cancel the scan on Ctrl+C. Adding `signal.NotifyContext` is a recommended improvement.

## Recommended next steps (nice-to-have changes)

- Move `main.go` into `cmd/tcp-scanner/` and update build instructions (if you prefer the conventional layout).
- Add a worker pool / bounded concurrency to speed up large scans and a `--workers` flag.
- Add progress reporting (progress bar) and an output formatter (JSON/CSV).
- Add unit tests for `internal/scan.CheckPort` using a temporary local listener.

## Security & legal notice

Only scan systems you own or have explicit permission to test. Unauthorized scanning can be considered intrusive or illegal. Use responsibly. For a safe testing target, see `scanme.nmap.org`.

## License

MIT © 2025 Siddharth Nema
