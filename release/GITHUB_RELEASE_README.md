# ASF v4.0.0 — GitHub Release Instructions

## Quick Upload (Recommended)

```bash
# From repository root
VERSION="4.0.0"

gh release create "v${VERSION}" \
  --title "ASF v${VERSION}" \
  --notes-file release/RELEASE_NOTES.md \
  --target main

gh release upload "v${VERSION}" "asf-tui/dist/ASF-v${VERSION}-darwin-arm64"
gh release upload "v${VERSION}" "asf-tui/dist/ASF-v${VERSION}-darwin-amd64"
gh release upload "v${VERSION}" "asf-tui/dist/ASF-v${VERSION}-linux-amd64"
gh release upload "v${VERSION}" "asf-tui/dist/ASF-v${VERSION}-linux-arm64"
gh release upload "v${VERSION}" "asf-tui/dist/ASF-v${VERSION}-windows-amd64.exe"
gh release upload "v${VERSION}" "asf-tui/dist/checksums.txt"
gh release upload "v${VERSION}" "release/RELEASE_NOTES.md"
gh release upload "v${VERSION}" "release/INSTALL.md"
```

## Manual Upload via GitHub UI

1. Go to https://github.com/moksh5936-2/asfassumption/releases/new
2. Tag: `v4.0.0`
3. Title: `ASF v4.0.0`
4. Description: Paste contents of `release/RELEASE_NOTES.md`
5. Attach files:
   - `asf-tui/dist/ASF-v4.0.0-darwin-arm64`
   - `asf-tui/dist/ASF-v4.0.0-darwin-amd64`
   - `asf-tui/dist/ASF-v4.0.0-linux-amd64`
   - `asf-tui/dist/ASF-v4.0.0-linux-arm64`
   - `asf-tui/dist/ASF-v4.0.0-windows-amd64.exe`
   - `asf-tui/dist/checksums.txt`
   - `release/RELEASE_NOTES.md`
   - `release/INSTALL.md`

## Post-Release Verification

```bash
# Verify all assets are downloadable
VERSION="4.0.0"
for bin in darwin-arm64 darwin-amd64 linux-amd64 linux-arm64 windows-amd64.exe; do
  url="https://github.com/moksh5936-2/asfassumption/releases/download/v${VERSION}/ASF-v${VERSION}-${bin}"
  echo "Checking: $url"
  curl -sfL -o /dev/null -w "  HTTP %{http_code}\n" "$url"
done

# Verify checksums
curl -sfL "https://github.com/moksh5936-2/asfassumption/releases/download/v${VERSION}/checksums.txt" | shasum -a 256 -c --ignore-missing
```

## Version

**Do not change the version number.** The release is tagged `v4.0.0`, matching `ASFVersion = "4.0.0"` in the source code.
