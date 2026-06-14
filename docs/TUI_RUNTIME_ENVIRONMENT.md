# TUI Runtime Environment

## Test Session Metadata

| Field | Value |
|-------|-------|
| **Date** | 2026-06-13 |
| **Time** | 21:00–22:00 UTC+5:30 |
| **Tester** | Automated QA (opencode) |
| **ASF Version** | v4.0.1 |
| **Go Version** | go1.24.2 darwin/arm64 |
| **OS** | macOS 26.5.1 (Build 25F80) |
| **Architecture** | arm64 (Apple M5) |
| **Shell** | /bin/zsh |
| **Terminal Emulator** | None (non-interactive CI-like environment) |
| **Terminal Size** | 80×24 (simulated) |
| **TERM** | (empty — no real TTY available) |

## Testing Method

The TUI binary (`asf-tui`) was tested via:

1. **`script` PTY**: Launched with `script -q /dev/null ./asf-tui-test` to capture rendered TUI output. Terminal size defaulted to 0×0 in PTY (confirmed limitation).
2. **CLI mode**: `asf analyze <file>` with JSON output for engine verification.
3. **`asf doctor --verbose`**: System diagnostics.
4. **`go test -count=1 -v ./...`**: Full programmatic test suite.
5. **`go vet ./...` / `go fmt ./...` / `go build ./...`**: Static analysis.

## Constraints

- Bubble Tea TUI requires a real TTY for keyboard input. Piped stdin (`echo "Q" | ./asf-tui`) does not work because Bubble Tea reads raw terminal input, not line-buffered stdin.
- The `script` PTY reported terminal size 0×0, causing the TUI's minimum-size check to trigger ("Terminal too small. Minimum: 60×12 Current: 0×0") on every launch. This is a PTY limitation, not a TUI defect.
- `expect` (v5.x) was available but could not set TIOCSWINSZ on the spawned PTY before the TUI queried it.
- Despite the size check, the TUI rendered the full Dashboard with sidebar, top bar, and bottom bar — indicating the size check is advisory, not a hard block.

## Verified Capabilities

| Capability | Status |
|------------|--------|
| Binary startup without crash | ✅ PASS |
| CLI `--help` / `--version` | ✅ PASS |
| CLI `analyze` (5 threat model files) | ✅ PASS |
| CLI `doctor --verbose` | ✅ PASS |
| Engine analysis (69–85 assumptions per file) | ✅ PASS |
| Error handling (no file, bad path, malformed YAML) | ✅ PASS |
| Test suite (21 packages, 48 regressions) | ✅ PASS |
| `go vet` / `go fmt` / `go build` | ✅ PASS |
| Dashboard rendering (sidebar, top bar, bottom bar) | ✅ Visible |
| Bottom bar hints | ✅ Visible |
| ASF version display | ✅ v4.0.1 |

## Cannot Verify Interactively

The following require a real TTY and interactive keyboard input, which this environment does not provide:

- Keyboard navigation (Tab, arrows, Enter, Esc, etc.)
- Sidebar cycling (all 16 items)
- File explorer (folder navigation, file selection, search)
- Results tab switching
- Export flow
- Mouse wheel scrolling
- Resize behavior
- Help screen toggling
