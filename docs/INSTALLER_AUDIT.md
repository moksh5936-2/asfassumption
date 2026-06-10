# Installer Audit

## `install.sh` Verification

| Check | Status | Notes |
|-------|--------|-------|
| OS detection | ✅ | `uname -s` → darwin/linux |
| Arch detection | ✅ | `uname -m` → amd64/arm64 |
| Linux amd64 | ✅ | `ASF-v{VERSION}-linux-amd64` |
| Linux arm64 | ✅ | `ASF-v{VERSION}-linux-arm64` |
| macOS amd64 | ✅ | `ASF-v{VERSION}-darwin-amd64` |
| macOS arm64 | ✅ | `ASF-v{VERSION}-darwin-arm64` |
| Binary permissions | ✅ | `chmod +x` after download |
| Checksum verification | ✅ | SHA-256 via `shasum` |
| PATH update | ✅ | Installs to `/usr/local/bin` or `~/.local/bin` |
| Config creation | ✅ | XDG path on Linux/macOS |
| Error handling | ✅ | HTTP 404 → helpful message |
| Rollback | ✅ | Temp dir cleanup via `trap` |
| Private repos | ✅ | `GITHUB_TOKEN` env var support |
| Public repos | ✅ | Direct download URL |
| `--upgrade` flag | ✅ | Replaces binary in-place |
| `--help` flag | ✅ | Shows usage |

## `install.ps1` Verification

| Check | Status | Notes |
|-------|--------|-------|
| OS detection | ✅ | Hardcoded `windows-amd64` |
| Arch detection | ✅ | Only amd64 (no ARM Windows) |
| Windows amd64 | ✅ | `ASF-v{VERSION}-windows-amd64.exe` |
| Binary permissions | ✅ | Windows doesn't use `+x` |
| Checksum verification | ✅ | SHA-256 via `Get-FileHash` |
| PATH update | ✅ | Adds `%LOCALAPPDATA%\ASF\bin` |
| Config creation | ✅ | `%APPDATA%\ASF\config.yaml` |
| Error handling | ✅ | Try/catch with cleanup |
| Rollback | ✅ | Temp dir removed on failure |
| Private repos | ✅ | `$Token` param or env var |
| `-Upgrade` flag | ✅ | Replaces existing |
| `-Help` flag | ✅ | Shows usage |

## Edge Cases

| Scenario | Expected | Actual |
|----------|----------|--------|
| No curl or wget | Error: "Need curl or wget" | ✅ |
| No write to /usr/local/bin | Fallback to ~/.local/bin | ✅ |
| Checksum mismatch | Error with hash values | ✅ |
| HTTP 404 on download | Error with troubleshooting | ✅ |
| Binary version mismatch | Warning, not fatal | ✅ |
| Install dir not in PATH | Warning with fix | ✅ |
| Already installed, no --upgrade | Info message, exit 0 | ✅ |
