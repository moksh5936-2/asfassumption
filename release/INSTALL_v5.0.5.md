# Installation Guide — ASF0 v5.0.5

## macOS (arm64)
```bash
curl -sfLO https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.5/ASF-v5.0.5-darwin-arm64
chmod +x ASF-v5.0.5-darwin-arm64
mkdir -p ~/.local/bin
cp ASF-v5.0.5-darwin-arm64 ~/.local/bin/asf
export PATH="$PATH:$HOME/.local/bin"
asf --version
# Expected: ASF0 v5.0.5
```

## macOS (amd64)
```bash
curl -sfLO https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.5/ASF-v5.0.5-darwin-amd64
chmod +x ASF-v5.0.5-darwin-amd64
mkdir -p ~/.local/bin
cp ASF-v5.0.5-darwin-amd64 ~/.local/bin/asf
export PATH="$PATH:$HOME/.local/bin"
asf --version
```

## Linux (amd64)
```bash
curl -sfLO https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.5/ASF-v5.0.5-linux-amd64
chmod +x ASF-v5.0.5-linux-amd64
mkdir -p ~/.local/bin
cp ASF-v5.0.5-linux-amd64 ~/.local/bin/asf
export PATH="$PATH:$HOME/.local/bin"
asf --version
```

## Linux (arm64)
```bash
curl -sfLO https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.5/ASF-v5.0.5-linux-arm64
chmod +x ASF-v5.0.5-linux-arm64
mkdir -p ~/.local/bin
cp ASF-v5.0.5-linux-arm64 ~/.local/bin/asf
export PATH="$PATH:$HOME/.local/bin"
asf --version
```

## Windows (PowerShell)
```powershell
curl.exe -sfLO https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.5/ASF-v5.0.5-windows-amd64.exe
mkdir $env:USERPROFILE\.asf 2>$null
move .\ASF-v5.0.5-windows-amd64.exe $env:USERPROFILE\.asf\asf.exe
$env:Path += ";$env:USERPROFILE\.asf"
asf --version
```

## Quick Install (curl | bash)
```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash
```

## Upgrade from v5.0.4
```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash -s -- --upgrade
```

## Clean Install (removes old binaries)
```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash -s -- --clean
```

## Verify Installation
```bash
asf --version
# Expected: ASF0 v5.0.5

asf doctor
# Should report version 5.0.5 and all paths writable
```

## Troubleshooting
- **Wrong version**: Run the upgrade command above
- **"command not found"**: Ensure `~/.local/bin` is in your PATH
- **Permission denied**: Run `chmod +x ~/.asf/asf`
- **Stale binary**: Run `asf doctor --fix` or `curl ... | bash -s -- --clean`
