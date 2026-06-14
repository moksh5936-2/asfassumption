# ASF0 v5.0.3 — Install Guide

## macOS & Linux

### Fresh Install
```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash
```

### Pin Specific Version
```bash
ASF_VERSION=5.0.3 curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash
```

### Upgrade from Older Version
```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash -s -- --upgrade
```

### Custom Install Directory
```bash
ASF_INSTALL_DIR=/usr/local/bin curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash
```

## Windows

Download the binary from the GitHub release assets:
```
https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.3/ASF-v5.0.3-windows-amd64.exe
```

Place it in a directory on your PATH (e.g., `C:\Windows\System32\` or `C:\Tools\`).

## Manual Download

### darwin-arm64 (Apple Silicon)
```bash
curl -fsSL https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.3/ASF-v5.0.3-darwin-arm64 -o asf && chmod +x asf && ./asf --version
```

### darwin-amd64 (Intel Mac)
```bash
curl -fsSL https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.3/ASF-v5.0.3-darwin-amd64 -o asf && chmod +x asf && ./asf --version
```

### linux-amd64
```bash
curl -fsSL https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.3/ASF-v5.0.3-linux-amd64 -o asf && chmod +x asf && ./asf --version
```

### linux-arm64
```bash
curl -fsSL https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.3/ASF-v5.0.3-linux-arm64 -o asf && chmod +x asf && ./asf --version
```

## Verify Install

```bash
asf --version
# Expected: ASF0 v5.0.3
```

## Troubleshooting

### "command not found" after install
Ensure `~/.local/bin` is on your PATH:
```bash
export PATH="$HOME/.local/bin:$PATH"
```
Add to `~/.bashrc`, `~/.zshrc`, or equivalent.

### "path is wrong" in installer
Restart your terminal or run:
```bash
source ~/.bashrc  # or source ~/.zshrc
```

### Still seeing old version
Remove stale binary:
```bash
rm -f ~/.asf/asf
```
Then reinstall.

### Stale Install Script
If using a locally cached install script, re-fetch:
```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash -s -- --upgrade
```

### Checksum Verification
```bash
curl -fsSL https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.3/checksums.txt
sha256sum ASF-v5.0.3-*
```
