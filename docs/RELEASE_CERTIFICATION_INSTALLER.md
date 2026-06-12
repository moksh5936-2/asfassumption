# Release Certification — Installer

## Files Audited

| File | Purpose | Lines |
|------|---------|-------|
| `install.sh` (root) | Primary curl-to-bash installer | 592 |
| `install.ps1` | Windows PowerShell installer | 332 |
| `release/install.sh` | Duplicate of root installer | 592 |
| `asf-tui/install.sh` | Local/dev installer | 234 |

## Issues Found

### Issue 1: Version Fallback Mismatch (MEDIUM)

**File:** `install.sh` (root), `release/install.sh`, `install.ps1`, `asf-tui/install.sh`

**Severity:** MEDIUM

**Description:** The fallback version in all installers is `2.1.1`, but the actual release is `v2.1.2`. If GitHub API is unreachable or rate-limited, the installer will attempt to download a non-existent release.

**Evidence:**
- `install.sh` line 364: `VERSION="2.1.1"`
- `install.ps1` line 82: `$Version = "2.1.1"`
- `asf-tui/install.sh` line 15: `ASF_VERSION="2.1.1"`

**Impact:** New installations will fail if GitHub API is unavailable. Upgrades may downgrade to v2.1.1.

**Reproduction:**
```bash
# Block GitHub API
export GITHUB_TOKEN=""
ASF_VERSION="" ./install.sh --help
# Or simply run without internet
```

### Issue 2: Duplicate Installer (LOW)

**File:** `release/install.sh`

**Severity:** LOW

**Description:** `release/install.sh` is an exact duplicate of `install.sh` (root). This creates maintenance risk — both must be updated in sync.

**Evidence:** `diff install.sh release/install.sh` shows identical files.

**Impact:** Risk of version drift if only one is updated.

### Issue 3: Missing Checksum for Local Install (LOW)

**File:** `asf-tui/install.sh`

**Severity:** LOW

**Description:** The local installer does not verify checksums when using a local binary. It only checks file existence and size.

**Evidence:** `asf-tui/install.sh` line 119-120:
```bash
if [ -n "$LOCAL_BINARY" ] && [ -s "$LOCAL_BINARY" ]; then
    cp "$LOCAL_BINARY" "${ASF_HOME}/asf"
```

### Issue 4: PATH Modification Without User Consent (LOW)

**File:** `install.sh` (root)

**Severity:** LOW

**Description:** The installer appends `export PATH="$PATH:${dir}"` to `.zshrc` or `.bashrc` without explicit confirmation. This is standard practice but may surprise users.

**Evidence:** `install.sh` line 193-195:
```bash
echo "" >> "$rc_file"
echo "# ASF" >> "$rc_file"
echo "export PATH=\"\$PATH:${dir}\"" >> "$rc_file"
```

### Issue 5: No Rollback on Checksum Failure (MEDIUM)

**File:** `install.sh` (root)

**Severity:** MEDIUM

**Description:** If checksum verification fails, the installer exits with `err` but does not clean up the partially downloaded binary. The temp directory is cleaned by `trap`, but the `err` function exits immediately.

**Evidence:** `install.sh` line 429-430:
```bash
if [ "$COMPUTED_HASH" != "$EXPECTED_HASH" ]; then
    err "Checksum mismatch! Expected ${EXPECTED_HASH}, got ${COMPUTED_HASH}. Download may be corrupted."
```

## Installer Capabilities Verified

| Capability | Status | File |
|------------|--------|------|
| Install | ✅ | `install.sh` |
| Upgrade | ✅ | `install.sh` |
| Repair | ✅ | `install.sh` |
| Clean | ✅ | `install.sh` |
| Purge | ✅ | `install.sh` |
| Checksum verification | ✅ | `install.sh` |
| PATH modification | ✅ | `install.sh` |
| Windows install | ✅ | `install.ps1` |
| Config creation | ✅ | All |
| License backup | ✅ | `install.sh` |

## Dry-Run Testing

```bash
# Test help
./install.sh --help                    # ✅ Works

# Test repair (no binary)
./install.sh --repair                  # ✅ Fails gracefully with error

# Test clean (no binary)
./install.sh --clean                   # ✅ Works (no-op)

# Test purge without clean
./install.sh --purge                   # ✅ Fails with usage error
```

## Verdict

⚠️ **WARN** — The installer is functional but has a version mismatch (2.1.1 vs 2.1.2) that could cause installation failures. The checksum failure handling and duplicate installer are also concerns.

**Required Fix:** Update all version references from `2.1.1` to `2.1.2` in all installer files.
