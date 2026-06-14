# Installation Guide — ASF0 v5.0.4

## macOS (ARM)

```bash
curl -sfLO https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.4/ASF-v5.0.4-darwin-arm64
chmod +x ASF-v5.0.4-darwin-arm64
mkdir -p ~/.local/bin
cp ASF-v5.0.4-darwin-arm64 ~/.local/bin/asf
```

Verify:

```bash
asf --version
# Expected: ASF0 v5.0.4
```

## macOS (Intel)

```bash
curl -sfLO https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.4/ASF-v5.0.4-darwin-amd64
chmod +x ASF-v5.0.4-darwin-amd64
mkdir -p ~/.local/bin
cp ASF-v5.0.4-darwin-amd64 ~/.local/bin/asf
```

## Linux (AMD64)

```bash
curl -sfLO https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.4/ASF-v5.0.4-linux-amd64
chmod +x ASF-v5.0.4-linux-amd64
mkdir -p ~/.local/bin
cp ASF-v5.0.4-linux-amd64 ~/.local/bin/asf
```

## Linux (ARM64)

```bash
curl -sfLO https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.4/ASF-v5.0.4-linux-arm64
chmod +x ASF-v5.0.4-linux-arm64
mkdir -p ~/.local/bin
cp ASF-v5.0.4-linux-arm64 ~/.local/bin/asf
```

## Windows

```powershell
curl.exe -sfLO https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.4/ASF-v5.0.4-windows-amd64.exe
mkdir $env:USERPROFILE\.asf -Force
move .\ASF-v5.0.4-windows-amd64.exe $env:USERPROFILE\.asf\asf.exe
```

## Upgrade from Older Version

```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash -s -- --upgrade
```

## Stale Binary Removal

```bash
rm -f ~/.local/bin/asf ~/.asf/asf
```

Then reinstall with the curl commands above.

## Verify Install

```bash
asf --version
# Expected: ASF0 v5.0.4

asf doctor
# Should report version 5.0.4 and all paths writable
```

## Troubleshooting

**Wrong version**: Run `which asf` and check for stale installations. Remove all but the latest.

**PATH issues**: Ensure `~/.local/bin` is in your PATH:
```bash
export PATH="$HOME/.local/bin:$PATH"
```

**Permission denied**: Ensure the binary is executable:
```bash
chmod +x ~/.local/bin/asf
```
