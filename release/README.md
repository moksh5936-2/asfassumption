# ASF Release Assets

## Version

**v2.0.0** — June 2026

> Go-native single binary — no Python runtime required.

## Downloads

| File | Platform | Architecture | Size |
|------|----------|-------------|------|
| `ASF-v2.0.0-linux-amd64` | Linux | AMD64 | 8.9MB |
| `ASF-v2.0.0-linux-arm64` | Linux | ARM64 | 8.3MB |
| `ASF-v2.0.0-darwin-amd64` | macOS | Intel | 9.1MB |
| `ASF-v2.0.0-darwin-arm64` | macOS | Apple Silicon | 8.6MB |
| `ASF-v2.0.0-windows-amd64.exe` | Windows | AMD64 | 9.2MB |
| `checksums.txt` | All | All | SHA-256 checksums |

## Verification

### Check SHA-256 checksums

```bash
shasum -a 256 -c checksums.txt
```

Expected output:
```
ASF-v2.0.0-darwin-amd64: OK
ASF-v2.0.0-darwin-arm64: OK
ASF-v2.0.0-linux-amd64: OK
ASF-v2.0.0-linux-arm64: OK
ASF-v2.0.0-windows-amd64.exe: OK
```

### Verify binary

```bash
./ASF-v2.0.0-darwin-arm64 --version
# Expected: ASF v2.0.0
```

## Installation

### macOS / Linux

```bash
# Using the installer
chmod +x install.sh
./install.sh

# Manual installation
chmod +x asf-darwin-arm64
sudo mv asf-darwin-arm64 /usr/local/bin/asf
```

### Windows

```powershell
# Not yet available — see BUILD_SYSTEM.md for cross-compilation
```

### Prerequisites

After installing ASF, you also need the Python ASF engine (REQUIRED):

```bash
cd /path/to/asf
pip install -e .
```

### Optional Dependencies

```bash
# Ollama (OPTIONAL — for AI features)
brew install ollama
ollama serve

# Tesseract (OPTIONAL — for image OCR)
brew install tesseract
```

## Building from Source

```bash
# Requires Go 1.24+
git clone https://github.com/moksh5936-2/asfassumption.git
cd asf/asf-tui
go build -o asf-tui .
```

## Release Assets

- `checksums.txt` — SHA-256 checksums
- `install.sh` — Installer script
- `VERSION` — Version manifest

## Architecture Support

| Capability | Status |
|------------|--------|
| Architecture diagram analysis | ✅ Full support |
| STRIDE threat mapping | ✅ 17 category rules + 30 keyword rules |
| Risk assessment | ✅ 5×5 deterministic risk matrix |
| Explainability | ✅ Evidence, justification, confidence |
| Export (JSON, MD, CSV, PDF, HTML) | ✅ All 5 formats |
| Architect review mode | ✅ Accept/Reject/Modified workflow |
| Validation mode | ✅ TUI-based precision/recall tracking |
| Local AI enhancement | ✅ Optional Ollama integration |
| Enterprise licensing | ✅ HMAC-signed license keys |

## Known Limitations

1. **Single platform binary** — Only macOS Apple Silicon is currently available
2. **Python dependency** — ASF engine requires separate `pip install`
3. **No code signing** — Binary is not notarized for macOS
4. **No CI/CD** — Release builds are manual
5. **No validation study** — Precision/recall metrics not yet established

## Support

- GitHub Issues: https://github.com/moksh5936-2/asfassumption/issues
