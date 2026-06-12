# Release Certification — D2C Readiness

## Release Assets

| Asset | Status | Size |
|-------|--------|------|
| ASF-v2.1.2-darwin-amd64 | ✅ | 9.7M |
| ASF-v2.1.2-darwin-arm64 | ✅ | 9.1M |
| ASF-v2.1.2-linux-amd64 | ✅ | 9.4M |
| ASF-v2.1.2-linux-arm64 | ✅ | 8.8M |
| ASF-v2.1.2-windows-amd64.exe | ✅ | 9.8M |
| checksums.txt | ✅ | 453 bytes |
| README.md | ✅ | 513 lines |
| VERSION | ✅ | 2.1.2 |

## Checksum Verification

```
8f73378585bfcf2029ec98b6fa0d26a9637d0fd451b0c4e3021fa1ed8007f008  ASF-v2.1.2-darwin-amd64
a6949f300b43a022650533d5b19dff3d5215cb353219ff8df5f41ec6a91ac2b2  ASF-v2.1.2-darwin-arm64
c315853bf449e5ece3d769162e085c0f2b37f86af82a9cc415e296575deb0cd5  ASF-v2.1.2-linux-amd64
cf955ca552eb79d93616a2616a3da563154566e521a32eb073f652b0c0156a31  ASF-v2.1.2-linux-arm64
f0520182a24021d4246143203e0ab2c69359d6b9862409cf196aef1d57236684  ASF-v2.1.2-windows-amd64.exe
```

## Version Consistency

| File | Version | Status |
|------|---------|--------|
| `license.go` | 2.1.2 | ✅ |
| `release/VERSION` | 2.1.2 | ✅ |
| `README.md` | 2.1.2 | ✅ (fixed from 2.1.1) |
| `install.sh` | 2.1.2 | ✅ (fixed from 2.1.1) |
| `install.ps1` | 2.1.2 | ✅ (fixed from 2.1.1) |
| `asf-tui/install.sh` | 2.1.2 | ✅ (fixed from 2.1.1) |

## Installer References

- Linux: `ASF-v2.1.2-linux-amd64` or `ASF-v2.1.2-linux-arm64`
- macOS: `ASF-v2.1.2-darwin-amd64` or `ASF-v2.1.2-darwin-arm64`
- Windows: `ASF-v2.1.2-windows-amd64.exe`

All installers reference the correct binary names.

## Customer Journey

| Step | Status | Notes |
|------|--------|-------|
| 1. Install | ✅ | `curl ... | bash` or `install.sh` |
| 2. Activate | ✅ | License file at `~/.asf/license.key` (optional for demo) |
| 3. Analyze | ✅ | `asf analyze file.yaml --json` |
| 4. Export | ✅ | TUI export (JSON, Markdown, CSV, HTML, PDF) |
| 5. Upgrade | ✅ | `install.sh --upgrade` |

## README Accuracy

- Features documented ✅
- Installation instructions present ✅
- Version updated to 2.1.2 ✅
- Links to GitHub repo present ✅

## Verdict

✅ **PASS** — Release assets are complete, version consistency achieved, customer journey is supported.
