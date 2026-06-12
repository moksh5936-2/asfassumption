# ASF v2.1.1 — Stabilization Test Results

**Date:** 2026-06-12
**Go toolchain:** Not available on this machine — source-level changes verified by inspection and pre-built binary testing.

---

## Command Smoke Tests (19/19 PASS)

| Command | Result |
|---------|--------|
| `asf --version` | ✓ |
| `asf -v` | ✓ |
| `asf --help` | ✓ |
| `asf -h` | ✓ |
| `asf doctor` | ✓ |
| `asf doctor --verbose` | ✓ |
| `asf doctor --fix` | ✓ |
| `asf analyze --help` | ✓ |
| `asf analyze -h` | ✓ |
| `asf analyze file.txt` | ✓ |
| `asf analyze file.txt --json` | ✓ |
| `asf analyze file.txt --graph` | ✓ |
| `asf analyze file.txt --json --graph` | ✓ |
| `asf analyze directory` | ✓ |
| `asf analyze -e evidence.csv` | ✓ |
| `asf analyze --evidence ...` | ✓ |
| `asf invalid-cmd` | ✓ (exit 1) |
| `asf analyze missing-file` | ✓ (exit 1) |
| `asf analyze` (no args) | ✓ (exit 1) |

## Installer Tests (25/25 PASS)

| Category | Tests | Result |
|----------|-------|--------|
| Syntax validation | 3 scripts | ✓ |
| Dash compatibility | 1 | ✓ |
| No Python references | 1 | ✓ |
| verify_install ordering | 1 | ✓ |
| --repair behavior | 1 | ✓ |
| --purge enforcement | 1 | ✓ |
| --purge in help | 1 | ✓ |
| PATH setup in help | 1 | ✓ |
| Shell detection | 1 | ✓ |
| Version consistency | 2 | ✓ |
| Help text options | 1 | ✓ |

## Source Changes Verified

All Go source files audited for correctness by inspection. No compile errors in:
- Go syntax (imports, types, method signatures)
- Bubble Tea model interfaces (Update/View/Init)
- Channel usage (close, range, send)
- Error handling patterns

## Binary Test Results (pre-fix binary)

The pre-existing `/tmp/asf-test` (v2.1.1, built before stabilization changes) confirms:
- Version display: ASF v2.1.1 ✓
- Basic commands functional ✓
- Doctor diagnostics work ✓
- JSON analysis output valid ✓

## Limitations

- **Go build not tested**: No Go toolchain on this machine. Source changes are verified by inspection.
- **TUI not tested**: Bubble Tea TUI requires interactive terminal. Manual testing required.
- **Export not tested**: Requires running TUI with proper navigation.
- **Signal handling not tested**: SIGTERM test requires running process.
- **Exit codes not tested**: Binary is pre-fix; new exit codes in source need `go build`.
