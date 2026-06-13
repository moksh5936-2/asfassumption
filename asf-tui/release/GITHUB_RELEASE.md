# ASF v3.0.0-RC2 — GitHub Pre-Release Package

## Release Configuration

| Field | Value |
|-------|-------|
| **Title** | `ASF v3.0.0-RC2` |
| **Tag** | `v3.0.0-RC2` |
| **Release Type** | Pre-Release |
| **Target** | `main` |

## Assets

| # | File | Source |
|---|------|--------|
| 1 | `asf-darwin-arm64` | `dist/asf-darwin-arm64` |
| 2 | `asf-darwin-amd64` | `dist/asf-darwin-amd64` |
| 3 | `asf-linux-amd64` | `dist/asf-linux-amd64` |
| 4 | `asf-linux-arm64` | `dist/asf-linux-arm64` |
| 5 | `asf-windows-amd64.exe` | `dist/asf-windows-amd64.exe` |
| 6 | `checksums.txt` | `dist/checksums.txt` |
| 7 | `RELEASE_NOTES_v3.0.0-RC2.md` | `release/RELEASE_NOTES_v3.0.0-RC2.md` |

## Release Description (for GitHub UI)

Paste the contents of `release/RELEASE_NOTES_v3.0.0-RC2.md` into the release description field.

## Upload Commands

Using GitHub CLI (`gh`):

```bash
# Create the release
gh release create v3.0.0-RC2 \
  --title "ASF v3.0.0-RC2" \
  --notes-file release/RELEASE_NOTES_v3.0.0-RC2.md \
  --prerelease \
  --target main

# Upload assets
gh release upload v3.0.0-RC2 dist/asf-darwin-arm64
gh release upload v3.0.0-RC2 dist/asf-darwin-amd64
gh release upload v3.0.0-RC2 dist/asf-linux-amd64
gh release upload v3.0.0-RC2 dist/asf-linux-arm64
gh release upload v3.0.0-RC2 dist/asf-windows-amd64.exe
gh release upload v3.0.0-RC2 dist/checksums.txt
gh release upload v3.0.0-RC2 release/RELEASE_NOTES_v3.0.0-RC2.md
```

## Asset Verification

After upload, verify all assets by downloading and checking checksums:

```bash
# Download checksums
curl -LO https://github.com/user/repo/releases/download/v3.0.0-RC2/checksums.txt

# Verify each asset
shasum -a 256 -c checksums.txt --ignore-missing
```

## Post-Release Smoke Test

After the release is published, test on each target platform:

```bash
# macOS ARM64
asf --version
asf --help

# Run a quick analysis
asf analyze path/to/sample.yaml
```
