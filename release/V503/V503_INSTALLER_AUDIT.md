# ASF0 v5.0.3 — Installer Audit

## File: `release/install.sh`

### Binary Name Pattern (line 386)
```
BINARY_NAME="ASF-${LATEST_VERSION}-${OS_FINAL}-${ARCH_FINAL}"
```

Generates: `ASF-v5.0.3-{os}-{arch}` ✓

### Download URL Pattern (line 387)
```
https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${BINARY_NAME}
```

Generates: `https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.3/ASF-v5.0.3-{os}-{arch}` ✓

### Checksums URL Pattern (line 388)
```
https://github.com/${REPO}/releases/download/${LATEST_VERSION}/checksums.txt
```

Generates: `https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.3/checksums.txt` ✓

### Version Detection (lines 360–383)
- Primary: GitHub API (latest non-draft release tag)
- Fallback: `v5.0.3` (updated) ✓

### URL Validation (lines 391–396)
- Whitespace/newline check ✓
- `ASF-` prefix check on binary name ✓

### Platform Detection (lines 48–68)
- darwin (arm64, amd64) ✓
- linux (arm64, amd64) ✓
- windows: stderr message to use install.ps1 ✓

### Upgrade Mode (lines 440+)
- Backs up existing binary ✓
- Downloads new version ✓
- Replaces binary ✓

### Issues
- `install.ps1` for Windows does not exist (pre-existing, not a release blocker)
