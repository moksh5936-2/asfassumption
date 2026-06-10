# ASF Installation Architecture

> Version: 1.0.0 | June 2026

## Overview

ASF distributes as a single statically-linked Go binary that acts as the TUI frontend. A Python package provides the assumption extraction backend. Optional dependencies add OCR and AI capabilities.

## Distribution Model

```
┌─────────────────────────────────────────────────────┐
│                  ASF Distribution                    │
├─────────────────┬───────────────────────────────────┤
│  Go Binary      │  asf-tui (11.9MB, static)         │
│  (TUI Frontend) │  Platform: darwin/arm64            │
│                 │  Dependencies: None (self-contained)│
├─────────────────┼───────────────────────────────────┤
│  Python Package │  asf/ (pip install -e .)           │
│  (ASF Engine)   │  Required: Python 3.8+             │
│                 │  Dependencies: See asf/setup.py     │
├─────────────────┼───────────────────────────────────┤
│  Optional       │  Tesseract (OCR)                   │
│  Dependencies   │  Ollama (Local AI)                 │
└─────────────────┴───────────────────────────────────┘
```

## Installation Methods

### Method 1: Quick Install (curl pipe)

```bash
curl -sfL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/asf-tui/install.sh | bash
```

**Process:**
1. Downloads the appropriate binary for the detected platform
2. Places it in `/usr/local/bin/asf`
3. Creates `~/.asf/` config directory
4. Outputs instructions for Python engine setup

**Risks:**
- curl-pipe-bash is a security anti-pattern
- No signature verification
- Single point of failure (GitHub RAW)

### Method 2: Manual Install

```bash
# 1. Download binary
curl -sfL https://github.com/moksh5936-2/asfassumption/releases/download/v1.0.0/asf-darwin-arm64 -o /usr/local/bin/asf
chmod +x /usr/local/bin/asf

# 2. Install Python engine
cd /path/to/asf
pip install -e .

# 3. Verify
asf --version
```

### Method 3: Build from Source

```bash
git clone https://github.com/moksh5936-2/asfassumption.git
cd asf/asf-tui
go build -o asf-tui .
```

**Requires:** Go 1.24+, git

## Python ASF Engine

The Python package at `asf/` provides the core assumption extraction:

```bash
pip install -e .
```

This runs ASF's `setup.py` which registers the `asf` CLI commands. The Go binary calls:

```bash
python3 -m asf.cli.main analyze --json <architecture_text>
```

### Engine Interface

| Aspect | Detail |
|--------|--------|
| Protocol | Subprocess (stdin → JSON) |
| Input | Pre-processed architecture text |
| Output | JSON with assumptions, severity, type |
| Error handling | Timeout-based (10s default) |
| Fallback | "Pipeline fallback" mode on failure |

## Optional Dependencies

### Tesseract (OCR)

```bash
brew install tesseract   # macOS
apt-get install tesseract # Linux
```

Used by `parser.go` when processing `.png`, `.jpg`, `.jpeg` files. Detected at runtime — if missing, image files show a "Tesseract not installed" error.

### Ollama (Local AI)

```bash
brew install ollama
ollama serve
```

Used by `ai.go` + `localai.go` for optional AI enhancement. Communicates via HTTP POST to `http://localhost:11434/api/generate`.

## Platform Support

| Platform | Binary | Status |
|----------|--------|--------|
| macOS ARM64 (Apple Silicon) | `asf-darwin-arm64` | ✅ Available, 11.9MB |
| macOS AMD64 (Intel) | `asf-darwin-amd64` | ❌ Not built |
| Linux AMD64 | `asf-linux-amd64` | ❌ Not built |
| Linux ARM64 | `asf-linux-arm64` | ❌ Not built |
| Windows AMD64 | `asf-windows-amd64.exe` | ❌ Not built |

**Currently only darwin/arm64 is available.** Cross-compilation requires Go toolchain setup.

## File Locations

| Path | Purpose |
|------|---------|
| `/usr/local/bin/asf` | Binary installation (default) |
| `~/.asf/config.yaml` | User configuration |
| `~/.asf/license.key` | Enterprise license key |
| `~/.config/asf/config.yaml` | Legacy config path (auto-migrated) |

## Offline Deployment

ASF is fully offline-capable:

1. Download binary on connected machine
2. Transfer via USB to air-gapped system
3. Copy Python ASF engine package
4. Install Python engine: `pip install -e .`
5. Run `asf`

The only features requiring network are optional: Ollama model downloads and API calls (local network only).

## Upgrade Path

```bash
# Re-run installer
curl -sfL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/asf-tui/install.sh | bash

# Or manual download
curl -sfL https://github.com/moksh5936-2/asfassumption/releases/download/v1.0.1/asf-darwin-arm64 -o /usr/local/bin/asf
chmod +x /usr/local/bin/asf
```

Upgrades preserve `~/.asf/config.yaml`. The Python ASF engine must be updated separately via `pip install -e .` from the updated repository.

## Verification

```bash
# Version check
asf --version
# Expected: ASF v1.0.0

# Checksum verification
shasum -a 256 -c release/checksums.txt
# Expected: asf-darwin-arm64: OK

# macOS code signing (if notarized)
codesign -dv /usr/local/bin/asf
```

## Known Issues

1. **Python dependency not gated** — No error message if Python is missing at startup. Only fails when analysis begins.
2. **No PATH verification** — Install script does not verify `/usr/local/bin` is in PATH.
3. **No uninstall script** — `install.sh` has no `--uninstall` flag.
4. **No version pinning** — Installer always downloads latest, no version selection.
5. **Single binary only** — Other platforms require manual cross-compilation.
