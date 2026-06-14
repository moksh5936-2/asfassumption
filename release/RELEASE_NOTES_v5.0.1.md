# ASF0 v5.0.1 — Case-First Navigation & UX Stabilization

## Overview
v5.0.1 restructures ASF0 from feature-centric navigation to case-centric navigation. The application is now organized around CASES, not FEATURES. All engine outputs are accessible from within the case workspace.

## What's New

### Case-First Navigation
- **WORK section removed** from sidebar (Review Queue, Validation Queue, Reports)
- Review, Validation, and Reports are now **key-driven actions** on the current case (`r`, `v`, `e`)
- Sidebar simplified to: CASES > + New Analysis > case files, AI > Local AI, SYSTEM > Settings/Help/About

### Workspace Tabs (7 tabs)
- **Overview** — Case metadata, risk distribution, verification stats, contradiction count, trust analysis, coverage summary
- **Assumptions** — Full assumption list with risk, category, confidence; scrollable and searchable
- **Verification** — Verified/Partial/Unverified/No Evidence breakdown with confidence
- **Contradictions** — All contradictions with severity and affected assumptions
- **Trust** — Trust chains + Single Points of Trust Failure + Priority Queue + CISO View
- **Controls** — Recommended control mappings with coverage
- **SDRI** — Executive summary + control inventory + coverage by category + design findings + architectural weaknesses + remediations + compliance alignments

### Tab Navigation
- `←` / `→` or `h` / `l` to switch workspace tabs
- Tab bar displayed at top of case workspace

### Startup Screen
- Fox branding with ASF0 title
- Enter to start, `?` for help, `q` to quit
- Centered, respects terminal size

### Modal File Picker
- Architecture file selection (single)
- Evidence file selection (multi-add)
- `↑↓/jk` navigate, `Enter` select/open, `Backspace` parent, `Esc` cancel
- `/` search, `.` toggle hidden files
- Tab preview panel

### Theme Polish
- Chiko-inspired palette: orange primary, amber secondary, soft green accents
- Dark charcoal background with off-white text
- Consistent card, sidebar, header, footer styling

## Bug Fixes
- View() ordering bug fixed (startup check after ready check)
- Emoji compatibility fix (`➕` → `+` for terminal display)
- Manual path typing removed from file picker workflow
- Contradiction count now displayed in Overview tab

## Breaking Changes
- WORK sidebar section removed. Use `r`/`v`/`e` keys from case workspace instead
- Old 12-tab workspace consolidated to 7 tabs (merged Trust/Impact/SPOFs, SDRI/Security Design Review, removed Reports/Blind Spots tabs)

## Files Changed
- `asf-tui/router.go` — Removed WORK section from sidebar
- `asf-tui/results.go` — Consolidated 12→7 tabs, merged content
- `asf-tui/app.go` — Tab navigation, key handlers, hints
- `asf-tui/help.go` — Updated sidebar tree, workspace keys
- `asf-tui/tui_test.go` — Updated for new architecture
- `asf-tui/license.go` — Version bump 5.0.0→5.0.1
- `asf-tui/styles.go` — Theme polish
- `install.sh`, `install.ps1`, `asf-tui/install.sh` — Version bump

## Build Information
- **Go version:** 1.22.5
- **Platforms:** darwin/arm64, darwin/amd64, linux/amd64, linux/arm64, windows/amd64
- **Build flags:** `CGO_ENABLED=0`, `-trimpath`, `-ldflags="-s -w"`
