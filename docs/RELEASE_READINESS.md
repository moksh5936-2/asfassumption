# Release Readiness

## Can ASF run without a repo checkout?

**Yes.** ASF is a single static Go binary. No runtime dependencies on:
- Repository checkout ❌ not needed
- Source code ❌ not needed  
- Development folders ❌ not needed
- Developer machine assumptions ❌ all removed

## Install Locations

| Platform | Binary | Config | Cache |
|----------|--------|--------|-------|
| Linux | `/usr/local/bin/asf` | `~/.config/asf/config.yaml` | `~/.cache/asf/` |
| macOS | `/usr/local/bin/asf` | `~/Library/Application Support/asf/config.yaml` | `~/Library/Caches/asf/` |
| Windows | `%LOCALAPPDATA%\ASF\bin\asf.exe` | `%APPDATA%\ASF\config.yaml` | `%LOCALAPPDATA%\ASF\cache` |

## Runtime Dependencies

| Component | Required | Discovery |
|-----------|----------|-----------|
| Go binary | ✅ Self-contained | Static binary (~8MB) |
| Python 3.8+ | ✅ For engine | PATH search + config override |
| ASF Python package | ✅ For engine | `python3 -m asf.cli.main` |
| Tesseract | 🔶 Optional OCR | System paths + PATH |
| Ollama | 🔶 Optional AI | System paths + PATH |

## Path Independence

All hardcoded developer paths have been removed from the Go source. The binary uses:
- `os.UserConfigDir()` for config (XDG/macOS AppData/Windows AppData)
- `os.UserCacheDir()` for cache 
- `os.Executable()` for self-discovery
- `exec.LookPath()` for dependency discovery

## Verification

```bash
asf doctor
```

Reports OS, arch, binary path, config/cache/license paths, Python engine status, and dependency availability.

## Build Verification

- `go vet ./...` — clean
- `go test ./...` — 20/20 PASS
- `go build -ldflags="-s -w"` — produces static binary
- Cross-compiles for linux/darwin/windows, amd64/arm64
