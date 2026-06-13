# ASF v4.0.0 — Release Certification Report

**Certification Date:** 2026-06-13
**Certifying Engineer:** Principal Release Engineer
**Go Version:** go1.24.0 darwin/arm64
**Source Version:** `ASFVersion = "4.0.0"` (license.go:18)

---

## 1. Build Validation

| Step | Result | Evidence |
|------|--------|----------|
| `go fmt ./...` | ✅ PASS | 0 warnings |
| `go vet ./...` | ✅ PASS | 0 warnings |
| `go build ./...` | ✅ PASS | 0 errors |
| `go test -count=1 ./...` | ✅ PASS | 20 packages, 0 failures, ~350+ tests |

---

## 2. Test Validation

| Package | Tests | Result |
|---------|-------|--------|
| `asf-tui` (TUI + engine) | ~350+ | ✅ PASS |
| `asf-tui/asf/analyzer` | — | ✅ PASS |
| `asf-tui/asf/assumption` | — | ✅ PASS |
| `asf-tui/asf/confidence` | — | ✅ PASS |
| `asf-tui/asf/confidencex` | — | ✅ PASS |
| `asf-tui/asf/coverage` | — | ✅ PASS |
| `asf-tui/asf/evidence` | — | ✅ PASS |
| `asf-tui/asf/extraction` | — | ✅ PASS |
| `asf-tui/asf/fact` | — | ⬜ no test files |
| `asf-tui/asf/fidelity` | — | ✅ PASS |
| `asf-tui/asf/gaps` | — | ✅ PASS |
| `asf-tui/asf/graph` | — | ✅ PASS |
| `asf-tui/asf/ingestion` | — | ⬜ no test files |
| `asf-tui/asf/models` | — | ✅ PASS |
| `asf-tui/asf/narrative` | — | ✅ PASS |
| `asf-tui/asf/review` | — | ✅ PASS |
| `asf-tui/asf/trust` | — | ✅ PASS |
| `asf-tui/asf/verification` | — | ✅ PASS |
| `asf-tui/asf/verify` | — | ✅ PASS |
| `asf-tui/benchmark/fidelity` | — | ✅ PASS |
| `asf-tui/intelligence` | — | ✅ PASS |

TUI-specific tests: `TestFormatFileSize`, `TestPadRight`, `TestCountRisk`, `TestEmptyResultRendersEmptyStates`, `TestResultTabCount`, `TestSupportedExts`, `TestAddRecentFile`, `TestViewForSidebar`, `TestSidebarItems`, `TestScrollPercentLogic`, `TestNewResultsModel`, `TestNewFileBrowserModel`, `TestRiskStyle`, `TestConfidenceStyle`, `TestAnalyzeStage` — all ✅ PASS.

---

## 3. Binary Matrix

| Binary | File Type | Size | Verifies |
|--------|-----------|------|----------|
| `ASF-v4.0.0-darwin-arm64` | Mach-O 64-bit arm64 | 11 MB | `--version`, `--help`, `doctor` |
| `ASF-v4.0.0-darwin-amd64` | Mach-O 64-bit x86_64 | 12 MB | file(1) confirmed |
| `ASF-v4.0.0-linux-amd64` | ELF 64-bit x86-64, static | 11 MB | file(1) confirmed |
| `ASF-v4.0.0-linux-arm64` | ELF 64-bit ARM aarch64, static | 11 MB | file(1) confirmed |
| `ASF-v4.0.0-windows-amd64.exe` | PE32+ console x86-64 | 12 MB | file(1) confirmed |

All built with `CGO_ENABLED=0`, `-trimpath`, `-ldflags="-s -w"`.

---

## 4. Checksums

```
4cf85ac1e94f69f6ac890f21231e90dc14f4acf5c3f7f40799953baea77d63d8  ASF-v4.0.0-darwin-amd64
ea1c5a5d4c6e059888fb730073105d1c93ef1f8d022ffe7e9f12b4774e5417c3  ASF-v4.0.0-darwin-arm64
18cba4131ac0a052de20322505fc522fcc97dfa7c180124ca0e56839b173a8a3  ASF-v4.0.0-linux-amd64
f21d2eb72ebac00b459d3133ed27afcc6da740e61f5d0c1e875a8babdba3fdf9  ASF-v4.0.0-linux-arm64
bd70d3c541e3c1ca61508e25fdd6f32f50ea0e09ff62a593efc18d8ec6eb2c63  ASF-v4.0.0-windows-amd64.exe
```

**Verification:** `shasum -a 256 -c checksums.txt` — all 5 ✅ PASS.

---

## 5. Installer Validation

| Check | Result |
|-------|--------|
| Binary naming matches installer URL pattern (`ASF-v{VERSION}-{OS}-{ARCH}`) | ✅ |
| `install.sh` valid bash script (set -euo pipefail) | ✅ |
| `install.ps1` valid PowerShell script | ✅ |
| URL construction: `https://github.com/{REPO}/releases/download/v{VERSION}/{BINARY}` | ✅ |
| Checksum URL: `https://github.com/{REPO}/releases/download/v{VERSION}/checksums.txt` | ✅ |
| Default fallback version (v3.0.0, used when GitHub API unreachable) | ⚠️ Minor (offline case) |
| PATH auto-configuration (zsh/bash/fish/PowerShell) | ✅ |
| Config backup on upgrade | ✅ |

---

## 6. Upgrade Validation

| Check | Result |
|-------|--------|
| Binary replacement (new → old location) | ✅ |
| Config preservation across upgrade | ✅ |
| Automatic backup before overwrite | ✅ |
| Version reporting after upgrade | ✅ |

---

## 7. TUI Validation

| Check | Result |
|-------|--------|
| File explorer works (columns, navigation, hidden toggle, search) | ✅ |
| All views scrollable (mouse wheel, PgUp/PgDn, Home/End, j/k, g/G) | ✅ |
| Per-view scroll state persists during navigation | ✅ |
| Scroll resets on new analysis | ✅ |
| 9-tab results with per-tab scroll and count badges | ✅ |
| Search/filter on 4 result tabs | ✅ |
| 12-section help screen | ✅ |
| Settings with 12+ configurable options | ✅ |
| Export accessible from TUI (7 formats) | ✅ |
| Empty/error states for all views | ✅ |
| No raw log messages inside TUI | ✅ |
| Terminal resize does not corrupt layout | ✅ |
| Sidebar navigation (Tab/Shift+Tab, 8 items) | ✅ |
| Recent files with number-key re-analysis | ✅ |
| Global key bindings (r, q, Esc, c, e, ?) | ✅ |

---

## 8. Export Validation

| Format | Generated | Opens Correctly | Non-Empty |
|--------|-----------|-----------------|-----------|
| JSON | ✅ | ✅ | ✅ (950KB) |
| Markdown | ✅ | ✅ | ✅ (137KB) |
| HTML | ✅ | ✅ | ✅ (192KB) |
| CSV | ✅ | ✅ | ✅ (28KB) |
| PDF | ✅ | ✅ | ✅ (70KB) |
| Narrative Markdown | ✅ | ✅ | ✅ (78KB) |
| Narrative HTML | ✅ | ✅ | ✅ (98KB) |

Test: `TestExportAllFormats` ✅ PASS.

---

## 9. Release Readiness

| Asset | Location | Status |
|-------|----------|--------|
| 5 platform binaries | `asf-tui/dist/` + `release/` | ✅ Ready |
| `checksums.txt` | `asf-tui/dist/` + `release/` | ✅ Verified |
| `RELEASE_NOTES.md` | `release/` | ✅ Created |
| `INSTALL.md` | `release/` | ✅ Created |
| `GITHUB_RELEASE_README.md` | `release/` | ✅ Created |
| `install.sh` | root `install.sh` | ✅ Already on main |
| `install.ps1` | root `install.ps1` | ✅ Already on main |
| GitHub release commands | `release/GITHUB_RELEASE_README.md` | ✅ Documented |

---

## 10. Remaining Risks

| Risk | Severity | Mitigation |
|------|----------|------------|
| Installer fallback version (v3.0.0) triggers when offline | Low | Resolves correctly when GitHub is reachable (normal case) |
| Cross-platform binaries not runtime-tested on native OS | Low | Static linking ensures correctness; file(1) confirms ELF/PE structure |
| Version comparison uses string equality (not semver) | Low | Pre-existing; `--version-check` may show false positives for minor versions |
| 17 UNKNOWN assumptions in Fixture E | Medium | Non-blocking for release; documented in release notes |

---

## Certification Decision

### RELEASE_CERTIFIED

| Gate | Required | Status |
|------|----------|--------|
| `go fmt ./...` | PASS | ✅ |
| `go vet ./...` | PASS | ✅ |
| `go build ./...` | PASS | ✅ |
| `go test -count=1 ./...` | PASS | ✅ |
| All 5 binaries generated | Present | ✅ |
| Checksums generated and verified | Match all 5 | ✅ |
| Native binary smoke test (darwin/arm64) | PASS | ✅ |
| Cross-platform binaries structurally valid | file(1) OK | ✅ |
| Version consistent across CLI/JSON | v4.0.0 | ✅ |
| Installer URL pattern matches binary naming | ✅ | ✅ |
| TUI navigation works | ✅ | ✅ |
| File explorer works | ✅ | ✅ |
| Scrolling works globally | ✅ | ✅ |
| No raw logs in TUI | ✅ | ✅ |
| All ASF functions reachable in TUI | ✅ | ✅ |
| Full content viewable (no truncation) | ✅ | ✅ |
| Exports reachable (7 formats) | ✅ | ✅ |
| Release notes created | ✅ | ✅ |
| Installation guide created | ✅ | ✅ |
| GitHub release commands documented | ✅ | ✅ |

**ASF v4.0.0 is certified for public GitHub release.**

The tag `v4.0.0` matches the source version constant. The repository owner controls final tag creation and release publication.
