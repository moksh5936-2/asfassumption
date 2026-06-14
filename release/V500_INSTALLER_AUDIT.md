# Installer Audit — ASF0 v5.0.0

## URL Pattern Verification

Installer scripts use `v5.0.0` fallback when GitHub API is unreachable.

Expected asset URLs resolve correctly:

| Platform | URL |
|----------|-----|
| macOS ARM64 | `https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.0/ASF-v5.0.0-darwin-arm64` |
| macOS AMD64 | `https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.0/ASF-v5.0.0-darwin-amd64` |
| Linux AMD64 | `https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.0/ASF-v5.0.0-linux-amd64` |
| Linux ARM64 | `https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.0/ASF-v5.0.0-linux-arm64` |
| Windows AMD64 | `https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.0/ASF-v5.0.0-windows-amd64.exe` |

## Files Updated

| File | Old Fallback | New Fallback |
|------|-------------|-------------|
| `install.sh` | `v4.0.0` | `v5.0.0` |
| `install.ps1` | `v4.0.0` | `v5.0.0` |
| `release/install.sh` | `v3.0.0` | `v5.0.0` |
| `asf-tui/install.sh` | `3.0.0` | `5.0.0` |

## Notes
- Installers dynamically discover the latest version from GitHub API at runtime
- Fallback version only used when API is unreachable
- v5.0.0 is now the default fallback
