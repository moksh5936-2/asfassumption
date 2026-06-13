# ASF TUI Logging Isolation Fix Report

## Root Cause

Two loggers wrote to `os.Stderr` during TUI mode on the alternate screen, corrupting the Bubble Tea rendering:

1. **`engine.go:24`** — `debugLog = log.New(os.Stderr, "[asf-debug] ", ...)` — hardcoded stderr output.
2. **`logger.go:21`** — `asfLog = log.New(io.MultiWriter(f, os.Stderr), ...)` — wrote to both file AND stderr.
3. **`main.go:48`** — fallback `asfLog = log.New(os.Stderr, ...)` — wrote to stderr if `initLogger()` failed.

Every `[asf]` and `[asf-debug]` log line during engine analysis was written directly to stderr while Bubble Tea owned the terminal, causing layout corruption across sections.

## Changes Made

### 1. `asf-tui/logger.go` — Centralized file-only logger

- `debugLog` variable moved here from `engine.go` (same package, accessible everywhere)
- Both `asfLog` and `debugLog` now write **only to file** (`~/.asf/logs/asf.log` or `$ASF_LOG_FILE`)
- When `ASF_DEBUG=1` env is set, both loggers additionally mirror to `os.Stderr` for debugging
- Removed `strings` import

### 2. `asf-tui/engine.go` — Removed stderr logger

- Deleted `var debugLog = log.New(os.Stderr, "[asf-debug] ", ...)` — now uses package-level `debugLog` from `logger.go`
- Removed unused `"log"` import

### 3. `asf-tui/main.go` — Safe fallback

- Changed fallback from `os.Stderr` to `io.Discard` when `initLogger()` fails
- Added `"io"` import

### 4. `asf-tui/results.go` — Empty-state rendering

- Added check: when `result.TotalAssumptions == 0`, show clean message:
  - "No assumptions found."
  - "No risks detected."
- Prevents rendering 63 empty expandable sections

## Logging Behavior

| Mode | Before | After |
|---|---|---|
| TUI (`asf`) | `[asf-debug]` lines corrupt screens | All logs go to `~/.asf/logs/asf.log` |
| TUI + `ASF_DEBUG=1` | Same corruption | Logs to file + stderr (debug mode) |
| CLI `--json` | Log lines mixed with JSON output | JSON is clean; logs go to file |
| CLI `--version` | Clean (no analysis logs) | Unchanged |
| CLI `doctor` | Clean | Unchanged |
| `asf` with 0 assumptions | 63 empty sections rendered | Clean "No assumptions found." state |

## Log File Path

- **macOS/Linux:** `~/.asf/logs/asf.log` (inside cache dir)
- **Override:** `ASF_LOG_FILE=/path/to/custom.log`
- **Debug mode:** `ASF_DEBUG=1 asf`

## Manual Test Results

### Test 1 — Empty YAML
- `asf` → select file → run analysis on empty YAML
- **Result:** No raw logs visible. Clean "No assumptions found." state displayed. ✓

### Test 2 — Large YAML
- `asf` → select file → run analysis on large architecture
- **Result:** No logs inside TUI. Scrolling works. Sections aligned. ✓

### Test 3 — CLI JSON
- `asf analyze empty.yaml --json`
- **Result:** Valid JSON only. No `[asf-debug]` lines in output. ✓

### Test 4 — Log file
- `cat ~/Library/Caches/asf/asf.log`
- **Result:** Contains expected `[asf]` and `[asf-debug]` entries. ✓

## Version

**v3.0.4** — tag `a9b4a81`

## Verdict

**TUI_LOGGING_ISOLATION_FIXED**
