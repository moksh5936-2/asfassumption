# ASF v3.0.0-RC2 — Smoke Test Report

**Date:** 2026-06-13
**Host:** darwin/arm64 (Apple M-series)
**Go Version:** go1.24.0

---

## Binary Verification

| Binary | File Type | Executable | --version | --help | Notes |
|--------|-----------|------------|-----------|--------|-------|
| `asf-darwin-arm64` | Mach-O 64-bit arm64 | ✅ | ✅ v3.0.0-RC2 | ✅ | Native on host |
| `asf-darwin-amd64` | Mach-O 64-bit x86_64 | ✅ (on Intel) | ⛔ (ARM host) | ⛔ (ARM host) | Correct arch; Intel-only runtime |
| `asf-linux-amd64` | ELF 64-bit x86-64, static | ✅ (by file) | ⛔ (cross-platform) | ⛔ (cross-platform) | Static binary, no runtime test |
| `asf-linux-arm64` | ELF 64-bit ARM aarch64, static | ✅ (by file) | ⛔ (cross-platform) | ⛔ (cross-platform) | Static binary, no runtime test |
| `asf-windows-amd64.exe` | PE32+ console x86-64 | ✅ (by file) | ⛔ (cross-platform) | ⛔ (cross-platform) | Windows executable |

## Detailed Results

### asf-darwin-arm64 (native)

```
$ ./asf-darwin-arm64 --version
ASF v3.0.0-RC2

$ ./asf-darwin-arm64 --help
ASF v3.0.0-RC2 — Architecture Security Framework

Usage:
  asf                        Launch the TUI
  asf --version, -v          Show version
  asf --license              Show license status
  asf analyze <file>         Run native analysis (JSON output)
  asf analyze <file> -e <ev> ...   With evidence files/dirs
  asf analyze <file> --graph Include graph in JSON output
  asf doctor                 Run system diagnostics
  asf doctor --verbose       Detailed diagnostics
  asf doctor --fix           Clean stale binaries
  asf --help, -h             Show this help
  asf --version-check       Check for newer version
```

No panic. No crash. Version prints correctly. Help text renders correctly.

### Cross-platform binaries

Linux and Windows binaries are statically linked (CGO_ENABLED=0) with correct ELF/PE formats. They cannot be executed on the darwin/arm64 build host but are structurally valid per `file(1)` inspection.

---

## Verdict

**SMOKE_TEST_PASS** — All 5 binaries are correctly formed. The native binary passes --version and --help without errors. Cross-platform binaries have correct file types for their target platforms.
