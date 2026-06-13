# ASF v3.0.0-RC2 — Installation Guide

---

## macOS

### Apple Silicon (ARM64)

```bash
# Download the binary
curl -LO https://github.com/user/repo/releases/download/v3.0.0-RC2/asf-darwin-arm64

# Make executable
chmod +x asf-darwin-arm64

# (Optional) Move to PATH
sudo mv asf-darwin-arm64 /usr/local/bin/asf

# Verify
asf --version
```

### Intel (AMD64)

```bash
# Download the binary
curl -LO https://github.com/user/repo/releases/download/v3.0.0-RC2/asf-darwin-amd64

# Make executable
chmod +x asf-darwin-amd64

# (Optional) Move to PATH
sudo mv asf-darwin-amd64 /usr/local/bin/asf

# Verify
asf --version
```

---

## Linux

### AMD64

```bash
# Download the binary
curl -LO https://github.com/user/repo/releases/download/v3.0.0-RC2/asf-linux-amd64

# Make executable
chmod +x asf-linux-amd64

# (Optional) Move to PATH
sudo mv asf-linux-amd64 /usr/local/bin/asf

# Verify
asf --version
```

### ARM64

```bash
# Download the binary
curl -LO https://github.com/user/repo/releases/download/v3.0.0-RC2/asf-linux-arm64

# Make executable
chmod +x asf-linux-arm64

# (Optional) Move to PATH
sudo mv asf-linux-arm64 /usr/local/bin/asf

# Verify
asf --version
```

---

## Windows

1. Download `asf-windows-amd64.exe`
2. Rename to `asf.exe`
3. Add to a directory in your `%PATH%` (e.g., `C:\Windows\System32\` or create `C:\asf\` and add to PATH)

Verify:

```cmd
asf --version
```

---

## Verify Checksums

After downloading, verify the binary integrity:

```bash
# Download the checksums file
curl -LO https://github.com/user/repo/releases/download/v3.0.0-RC2/checksums.txt

# On macOS/Linux:
shasum -a 256 -c checksums.txt --ignore-missing

# On Windows (PowerShell):
# Get-FileHash asf.exe -Algorithm SHA256
```

---

## Quick Start

```bash
# Run analysis on an architecture YAML file
asf analyze sample.yaml

# Run analysis with JSON output
asf analyze sample.yaml > results.json

# Run system diagnostics
asf doctor

# Launch TUI
asf
```

### Sample Architecture File

```yaml
name: Sample Architecture
version: "1.0"
components:
  - name: WebApp
    technology: Go
    network: public
  - name: Database
    technology: PostgreSQL
    network: private
assumptions:
  - All traffic uses TLS
  - Database access is restricted
```

---

## Requirements

- **OS:** macOS 11+ (ARM64 or AMD64), Linux (AMD64 or ARM64), Windows 10+ (AMD64)
- **Disk:** ~20 MB per binary
- **Dependencies:** None (statically linked)

---

## Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| `"Bad CPU type"` | Wrong binary for architecture | Download the correct platform binary |
| `"Permission denied"` | Not executable | Run `chmod +x <binary>` |
| `"command not found"` | Not in PATH | Move binary to `/usr/local/bin/` or add to PATH |
| No output from `analyze` | Missing input file | Provide a valid YAML file path |
