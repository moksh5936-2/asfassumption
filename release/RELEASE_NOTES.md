# ASF v4.0.0 — Release Notes

**Architecture Security Framework — Terminal-Native Product Experience**

---

## Overview

ASF v4.0.0 is a major release that transforms the TUI from a minimal interface into a full terminal-native product experience. The entire TUI has been rebuilt around Midnight-Commander-style file navigation, Lazygit-inspired keyboard UX, K9s-style structured navigation, and proper scrollable views throughout.

This release also includes all engine improvements from the v3.0.0-RC2 benchmark evaluation: rewritten verification engine, contradiction precision fixes, trust chain CLI exposure, SDRI control awareness, and the full 11-engine production pipeline.

---

## TUI Improvements

- **MC-style file explorer** — column layout, directory navigation, file preview (40 lines), path breadcrumb, hidden files toggle (`.`), search filter (`/`), `h` for help, `Backspace` for parent directory
- **Stable 8-item sidebar** — Dashboard, Analyze, Results, File Explorer, AI Models, Settings, About, Help; Tab/Shift+Tab navigation
- **Per-view scroll tracking** — each view maintains its own scroll position; scroll resets on new analysis; supports mouse wheel, PgUp/PgDn, Home/End, j/k, g/G
- **Dashboard** — recent files with number-key quick re-analysis
- **9-tab results** — Summary, Assumptions, Verification, Contradictions, Trust, Impact, Blind Spots, Controls, Reports; per-tab scroll; count badges
- **Search/filter** — `/` opens search on 4 result tabs (Assumptions, Contradictions, Trust, Controls); `n`/`N` scrolls matches
- **Global key bindings** — `r` run analysis, `q` universal back, `Esc` cancel analysis, `c` clear results, `e` export, `?` help
- **Help screen** — 12-section keyboard reference
- **Settings** — theme (4), fox style, analysis depth, risk threshold, STRIDE toggle, controls toggle, export format, AI model, debug logging toggle, reset to defaults
- **Empty/error states** — clean messages for no-file, no-results, parse-errors, unsupported files, etc.
- **Export** — 7 formats via settings screen (JSON, Markdown, HTML, CSV, PDF, Narrative MD, Narrative HTML)
- **All long-line truncation removed** — content wraps via viewport instead of `[:57]+"..."`

---

## Fixed Issues

- TUI scrolling is now smooth and consistent across all views (was broken/absent)
- No content clipped by panel borders
- No raw log messages inside TUI (all logs to `~/.asf/logs/asf.log`)
- Terminal resize no longer causes layout corruption
- File explorer parent directory navigation no longer conflicts with help (`Backspace` = parent, `h` = help)
- Debug logging toggle in settings persists to config (ASF_DEBUG env var on restart)
- Version constant corrected from `3.0.0-RC2` to `4.0.0`
- Version comparison is now substring-safe (installer no longer reports false "already installed")

---

## Installer Fixes

- `install.sh` — PATH auto-configuration for zsh/bash, `--repair`/`--clean`/`--purge` modes, symlink support
- `install.ps1` — Windows PATH setup, `-Repair`/`-Clean`/`-Purge` modes
- Release recovery: installer resolves correct version tag
- Default version fallback updated
- Checksum verification across all platforms

---

## Export Improvements

- All 7 export formats include structured controls, compliance data, CIE contradictions, TBI zones/boundaries/weaknesses, APD attack paths, TMI threat clusters
- JSON schema extended with trust-chain fields (100 chains, 33 SPOFs, 50 failure cascades)
- All exports accessible via TUI settings screen (no CLI required)

---

## Known Limitations

1. **Tesseract not bundled** — OCR requires `tesseract` installed separately (`brew install tesseract` / `apt install tesseract-ocr`)
2. **AI features require Ollama** — optional; disabled by default
3. **CIARE compliance coverage** remains 0% — framework-mapping model is an architectural limitation outside this release scope
4. **17 UNKNOWN assumptions** on positive verification fixture — generated assumptions with categories like `Compliance` and `Privacy` don't match existing control maps
5. **Installer fallback version** (v3.0.0) used when GitHub API is unreachable — will resolve correctly when online

---

## Build Artifacts

| Platform | Binary | Size |
|----------|--------|------|
| macOS ARM64 | `ASF-v4.0.0-darwin-arm64` | 11 MB |
| macOS AMD64 | `ASF-v4.0.0-darwin-amd64` | 12 MB |
| Linux AMD64 | `ASF-v4.0.0-linux-amd64` | 11 MB |
| Linux ARM64 | `ASF-v4.0.0-linux-arm64` | 11 MB |
| Windows AMD64 | `ASF-v4.0.0-windows-amd64.exe` | 12 MB |

All binaries built with `CGO_ENABLED=0`, `-trimpath`, `-ldflags="-s -w"` (stripped).

---

## Upgrade Notes

- No database migrations required
- No configuration changes required
- Drop-in replacement for any v2.x or v3.0.x installation
- Config backup is automatic on upgrade via installer `--upgrade` flag
- JSON output schema is backward compatible with v3.0.0-RC2 (trust-chain fields added in RC2, retained unchanged)
