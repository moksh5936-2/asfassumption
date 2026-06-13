# ASF v4.0.0 — Installation Guide

## macOS

### Option 1: Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash
```

After install, open a new terminal and run:

```bash
asf
```

### Option 2: Manual Download

1. Download the binary for your architecture:
   - **Apple Silicon (M1/M2/M3/M4):** `ASF-v4.0.0-darwin-arm64`
   - **Intel Mac:** `ASF-v4.0.0-darwin-amd64`

2. Make it executable and install:

```bash
chmod +x ASF-v4.0.0-darwin-arm64
mkdir -p ~/.local/bin
cp ASF-v4.0.0-darwin-arm64 ~/.local/bin/asf
echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.zshrc
source ~/.zshrc
```

### Option 3: Build from Source

```bash
git clone https://github.com/moksh5936-2/asfassumption.git
cd asf-tui
go build -o asf .
cp asf ~/.local/bin/
```

---

## Linux

### Option 1: Quick Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash
```

### Option 2: Manual Download

```bash
# AMD64
wget https://github.com/moksh5936-2/asfassumption/releases/download/v4.0.0/ASF-v4.0.0-linux-amd64
chmod +x ASF-v4.0.0-linux-amd64
sudo mv ASF-v4.0.0-linux-amd64 /usr/local/bin/asf

# ARM64 (e.g., Raspberry Pi)
wget https://github.com/moksh5936-2/asfassumption/releases/download/v4.0.0/ASF-v4.0.0-linux-arm64
chmod +x ASF-v4.0.0-linux-arm64
sudo mv ASF-v4.0.0-linux-arm64 /usr/local/bin/asf
```

---

## Windows

### Option 1: Quick Install (PowerShell)

```powershell
powershell -c "irm https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.ps1 | iex"
```

### Option 2: Manual Download

1. Download `ASF-v4.0.0-windows-amd64.exe` from the releases page
2. Create a folder: `mkdir %LOCALAPPDATA%\ASF\bin`
3. Move the exe there: `move ASF-v4.0.0-windows-amd64.exe %LOCALAPPDATA%\ASF\bin\asf.exe`
4. Add to PATH: `[Environment]::SetEnvironmentVariable("Path", "$env:LOCALAPPDATA\ASF\bin;$env:PATH", "User")`
5. Open a new PowerShell window and run: `asf`

---

## Upgrade Existing Installation

### macOS / Linux

```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash -s -- --upgrade
```

This backs up your config and license before upgrading.

### Windows (PowerShell)

```powershell
powershell -c "irm https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.ps1 | iex" -Upgrade
```

### Manual Upgrade

Replace the existing binary at `~/.asf/asf` (or `%LOCALAPPDATA%\ASF\asf.exe`) with the new version.

---

## Verify Installation

```bash
asf --version
```

Expected output:

```
ASF v4.0.0
```

If you see a "newer version available" message, it means a later release exists — you are still running v4.0.0 successfully.

```bash
asf doctor
```

This runs system diagnostics and confirms the binary is working correctly.

```bash
asf
```

This launches the TUI.

---

## Prerequisites (Optional)

For full functionality:

| Feature | Requirement | Install |
|---------|-------------|---------|
| OCR (image parsing) | Tesseract | `brew install tesseract` / `apt install tesseract-ocr` |
| AI enhancement | Ollama | `brew install ollama` / https://ollama.com |

Both are optional. ASF works fully without them.
