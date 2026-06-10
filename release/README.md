# ASF Release Assets

## Version

**v1.0.0** — June 2026

## Downloads

| File | Platform | Architecture | Size |
|------|----------|-------------|------|
| `asf-linux-amd64` | Linux | AMD64 | ~11MB |
| `asf-linux-arm64` | Linux | ARM64 | ~11MB |
| `asf-darwin-amd64` | macOS | Intel | ~11MB |
| `asf-darwin-arm64` | macOS | Apple Silicon | 11.9MB |
| `asf-windows-amd64.exe` | Windows | AMD64 | ~11MB |
| `install.sh` | All | All | Installer script |

## Verification

### Check SHA-256 checksums

```bash
shasum -a 256 -c checksums.txt
```

Expected output:
```
asf-darwin-arm64: OK
install.sh: OK
```

### Verify binary

```bash
# macOS
codesign -dv asf-darwin-amd64

# Linux
file asf-linux-amd64
# Expected: ELF 64-bit LSB executable, x86-64

# All platforms
./asf-darwin-arm64 --version
# Expected: ASF v1.0.0
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

1. Download `asf-windows-amd64.exe`
2. Rename to `asf.exe`
3. Move to a directory in your PATH
4. Run from Command Prompt or PowerShell: `asf`

### Prerequisites

After installing ASF, you also need:

```bash
# Python ASF engine (REQUIRED)
cd /path/to/asf
pip install -e .

# Ollama (OPTIONAL — for AI features)
brew install ollama
ollama serve

# Tesseract (OPTIONAL — for image OCR)
brew install tesseract
```

## Upgrade

```bash
# Re-run the install script
curl -sfL https://raw.githubusercontent.com/asfsecurity/asf/main/install.sh | bash

# Or download the new binary manually
```

## License Activation

```bash
echo 'ASF-XXXX-XXXX-XXXX-XXXX' > ~/.asf/license.key
asf --license
```

## Offline Usage

ASF is fully offline-capable:

1. Download the binary on a connected machine
2. Transfer via USB to the air-gapped system
3. Copy the Python ASF engine package
4. Run `asf` normally

The only features requiring network:
- Ollama model downloads (optional)
- Ollama API calls (optional, localhost only)

## AI Integration Setup

```bash
# 1. Install Ollama
brew install ollama

# 2. Start Ollama server
ollama serve

# 3. Launch ASF and go to AI Settings
asf
# → AI Settings → Select model → Download → Set as Active
```

## Building from Source

```bash
# Requires Go 1.24+
git clone https://github.com/asfsecurity/asf.git
cd asf/asf-tui
go build -o asf-tui .
```

## Release Assets

- `*.tar.gz` — Binary archives (when available)
- `checksums.txt` — SHA-256 checksums
- `install.sh` — Installer script
- `VERSION` — Version manifest

## Support

- GitHub Issues: https://github.com/asfsecurity/asf/issues
- Security: security@asfsecurity.com
