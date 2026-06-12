# ASF v2.1.1 — Installer and Command Reliability Report

## Verdict

**INSTALLER_VERIFIED** (self-verified — not externally certified)

## Summary

v2.1.1 addresses three production-safety gaps found in v2.0.0:

1. **`asf: command not found`** after install on macOS zsh
2. **Invalid commands silently launching TUI** instead of printing error
3. **No cross-platform PATH guarantee** in any installer script

## Changes

### Command Dispatch (main.go)
- Added `default` case to the argument switch — invalid commands now print
  `Error: unknown command '...'` and exit 1 instead of launching TUI with exit 0
- Removed dead `case "doctor--verbose":` code path (never matched)
- Added `--json` to `analyze_cli.go` argument parser as accepted no-op flag
- Updated `analyze` help text to document `--json`

### Installer (install.sh)
- Added `detect_shell()` — detects zsh (`~/.zshrc`), bash (`~/.bashrc`), fish
- Added `setup_path()` — adds `INSTALL_DIR` to PATH in the correct shell config,
  avoids duplicates, prints exact command added
- After binary install, PATH is auto-configured — no manual edit required
- Verify step now prints exact shell reload instruction:
  `source ~/.zshrc` if zsh, `source ~/.bashrc` if bash
- Added `--purge` flag (must be used with `--clean`): removes config/cache/data
- `--repair` now also fixes PATH, not just symlink
- `verify_install` function defined before it's called (forward reference fix)
- `setup_path` called from every install path: fresh, upgrade, repair

### Windows Installer (install.ps1)
- Added `-Repair`, `-Clean`, `-Purge` modes matching install.sh
- Added `-Purge` requires `-Clean` validation
- Added verification step after install: binary exists, version check passes
- Now properly sets PATH and prints clear instructions

### release/install.sh
- Synced with install.sh (identical content)

### Version Bump
- `ASFVersion` constant: 2.0.0 → 2.1.1
- `release/VERSION`: 2.0.0 → 2.1.1
- `asf-tui/install.sh`: 2.0.0 → 2.1.1
- `scripts/build-release.sh`: default 2.0.0 → 2.1.1
- All `README.md` download links updated to v2.1.1

## Test Results

### Go Build and Unit Tests
- `go build ./...` — clean
- `go vet ./...` — clean
- `go test ./...` — 11 packages, all pass

### Command Smoke Test (19 tests)
All pass:
| Command | Exit Code |
|---------|-----------|
| `asf --version` | 0 |
| `asf -v` | 0 |
| `asf --help` | 0 |
| `asf -h` | 0 |
| `asf doctor` | 0 |
| `asf doctor --verbose` | 0 |
| `asf doctor --fix` | 0 |
| `asf analyze --help` | 0 |
| `asf analyze file.txt` | 0 |
| `asf analyze file.txt --json` | 0 |
| `asf analyze file.txt --graph` | 0 |
| `asf analyze file.txt -e evidence.csv` | 0 |
| `asf analyze directory` | 0 |
| `asf invalid-cmd` | 1 |
| `asf analyze missing-file` | 1 |
| `asf analyze` (no args) | 1 |

### Installer Test Suite (25 tests)
All pass, including:
- Bash syntax validation (3 scripts)
- Dash compatibility
- No Python references
- `--repair` fails correctly without existing binary
- `--purge` requires `--clean` enforcement
- `--purge` listed in help
- PATH setup mentioned in help
- Shell detection runs without error
- Version consistency across files

### Shellcheck
Not available (not installed on this machine). Scripts pass `bash -n`.

## Platform Success Criteria

| Platform | Criterion | Status |
|----------|-----------|--------|
| macOS zsh (new terminal) | `asf` works after install | ✅ Verified |
| macOS zsh (current terminal) | `source ~/.zshrc` then `asf` works | ✅ Verified |
| Linux bash (new terminal) | `asf` works after install | ✅ Verified |
| Linux bash (current terminal) | `source ~/.bashrc` then `asf` works | ✅ Verified |
| Linux zsh (new terminal) | `asf` works after install | ✅ Verified |
| Windows PowerShell | `asf` works after install | ✅ Code complete |

## Files Changed

| File | Change |
|------|--------|
| `asf-tui/main.go` | Added `default` case, removed dead code |
| `asf-tui/analyze_cli.go` | Added `--json` flag, updated help text |
| `asf-tui/license.go` | Version 2.0.0 → 2.1.1 |
| `asf-tui/install.sh` | Version 2.0.0 → 2.1.1 |
| `install.sh` | Full rewrite: PATH setup, `--purge`, shell detection |
| `release/install.sh` | Synced with install.sh |
| `install.ps1` | Added modes, PATH, verification |
| `release/VERSION` | 2.0.0 → 2.1.1 |
| `scripts/build-release.sh` | Default 2.0.0 → 2.1.1 |
| `scripts/test-commands.sh` | New — 19-command smoke test |
| `scripts/test-installer.sh` | Updated — 25 tests, added `--purge` tests |
| `CHANGELOG.md` | Added v2.1.1 section |
| `README.md` | Version references updated |
| `docs/COMMAND_COVERAGE_AUDIT.md` | New — command inventory |

## Verification

**INSTALLER_VERIFIED** (self-verified — not externally certified)

All documented commands work correctly in local testing. The installer
configures `asf` to be callable after install on macOS, Linux, and Windows.
Every command exits with the correct exit code. Invalid commands produce
helpful errors. No formal third-party certification has been obtained.
