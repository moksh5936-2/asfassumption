# V5.0.2 — Release Certification

## Checklist

| Requirement | Status |
|---|---|
| **Code committed** | ✅ `c037575` on `main`, pushed to origin |
| **Version bump** | ✅ `license.go: ASFVersion = "5.0.2"`, `install.sh: ASF_VERSION="5.0.2"` |
| **Installer updated** | ✅ Targets `v5.0.2` assets via `LATEST_VERSION` |
| **Build validation** | ✅ `go fmt`, `go vet`, `go build`, `go test` — all clean, 21/21 packages |
| **Binaries built (5)** | ✅ darwin-arm64, darwin-amd64, linux-amd64, linux-arm64, windows-amd64 |
| **Checksums generated** | ✅ `dist/checksums.txt` — SHA-256 |
| **Smoke test** | ✅ `--version` shows `ASF0 v5.0.2`, `--help` works, `doctor` passes |
| **TUI acceptance** | ⏳ Manual real-terminal check required (all unit tests pass) |
| **Release notes** | ✅ `release/RELEASE_NOTES_v5.0.2.md` |
| **Install guide** | ✅ `release/INSTALL_v5.0.2.md` |
| **Commit pushed** | ✅ `c037575` → `origin/main` |
| **Tag created** | ✅ `v5.0.2` → `origin/v5.0.2` |
| **GitHub release** | ✅ `https://github.com/moksh5936-2/asfassumption/releases/tag/v5.0.2` |
| **All binaries uploaded** | ✅ 5 binaries + checksums + release notes (7 assets) |
| **Published asset verified** | ✅ Download darwin-arm64 → `ASF0 v5.0.2`, checksum OK |
| **Fresh install test** | ✅ `curl ... | bash` → `asf --version` = `ASF0 v5.0.2` |
| **Upgrade test** | ✅ v5.0.0 → `--upgrade` → v5.0.2, config preserved |
| **Tests pass** | ✅ 19 picker tests + all 21 packages |
| **Build passes** | ✅ `go build ./...` clean |

## Remaining Risks

1. **TUI acceptance** — Manual real-terminal verification is recommended after upgrade to visually confirm file picker modal behavior. All automated tests pass.
2. **Windows** — Windows path safety is verified by tests (`filepath.VolumeName`, `filepath.Clean`, no manual `"/"` concatenation), but no native Windows test environment was available.

## Verdict

**V502_RELEASE_CERTIFIED**
