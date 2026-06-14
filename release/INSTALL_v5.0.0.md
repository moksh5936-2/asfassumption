# Installation Guide — ASF0 v5.0.0

## Prerequisites
- A terminal with 24-bit color support
- Minimum 80×24 terminal dimensions
- No Go toolchain required (self-contained binary)

## macOS

### Apple Silicon (M1/M2/M3/M4)
```bash
curl -sfLO https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.0/ASF-v5.0.0-darwin-arm64
chmod +x ASF-v5.0.0-darwin-arm64
mkdir -p ~/.local/bin ~/.asf
cp ASF-v5.0.0-darwin-arm64 ~/.asf/asf
ln -sf ~/.asf/asf ~/.local/bin/asf
```

### Intel
```bash
curl -sfLO https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.0/ASF-v5.0.0-darwin-amd64
chmod +x ASF-v5.0.0-darwin-amd64
mkdir -p ~/.local/bin ~/.asf
cp ASF-v5.0.0-darwin-amd64 ~/.asf/asf
ln -sf ~/.asf/asf ~/.local/bin/asf
```

## Linux

### AMD64
```bash
curl -sfLO https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.0/ASF-v5.0.0-linux-amd64
chmod +x ASF-v5.0.0-linux-amd64
mkdir -p ~/.local/bin ~/.asf
cp ASF-v5.0.0-linux-amd64 ~/.asf/asf
ln -sf ~/.asf/asf ~/.local/bin/asf
```

### ARM64
```bash
curl -sfLO https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.0/ASF-v5.0.0-linux-arm64
chmod +x ASF-v5.0.0-linux-arm64
mkdir -p ~/.local/bin ~/.asf
cp ASF-v5.0.0-linux-arm64 ~/.asf/asf
ln -sf ~/.asf/asf ~/.local/bin/asf
```

## Windows

Open **PowerShell** as Administrator:
```powershell
curl.exe -sfLO https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.0/ASF-v5.0.0-windows-amd64.exe
mkdir -force $env:USERPROFILE\.asf
move .\ASF-v5.0.0-windows-amd64.exe $env:USERPROFILE\.asf\asf.exe
```

Add `%USERPROFILE%\.asf` to your `PATH`:
```powershell
[Environment]::SetEnvironmentVariable("Path", "$env:USERPROFILE\.asf;" + [Environment]::GetEnvironmentVariable("Path", "User"), "User")
```

## Verify Installation
```bash
asf --version
# Expected: ASF0 v5.0.0
```

## Upgrading
```bash
# macOS / Linux
curl -sfL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash

# Windows
powershell -c "irm https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.ps1 | iex"

# Doctor cleanup after upgrade
asf doctor
```

## Uninstall
```bash
rm -f ~/.local/bin/asf ~/.asf/asf
rm -rf ~/.asf

# Windows
rm $env:USERPROFILE\.asf\asf.exe
```
