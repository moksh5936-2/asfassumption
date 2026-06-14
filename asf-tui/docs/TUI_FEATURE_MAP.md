# TUI Feature Map — Feature Inventory

## Navigation Features

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 1 | Sidebar with 18 entries | ✅ Done | Dashboard through About |
| 2 | Tab/Shift+Tab cycle sidebar | ✅ Done | Falls through for results/fileBrowser |
| 3 | Arrow keys / j/k move sidebar | ✅ Done | |
| 4 | Enter activates sidebar entry | ✅ Done | Per-view handlers |
| 5 | Direct keys (f, r, ?, q, Esc, e, c) | ✅ Done | handleGlobalKey switch |
| 6 | navigateTo/navigateBack wrappers | ✅ Done | Scroll save/restore |
| 7 | History stack (max 50) | ✅ Done | Router.history |
| 8 | Navigation from dashboard quick actions | ✅ Done | Uses same navigateMsg/router |
| 9 | Help screen with sidebar map | ✅ Done | help.go — 18 items |
| 10 | About view | ✅ Done | sidebar entry 17 |

## Viewport & Scrolling

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 11 | Single shared viewport | ✅ Done | m.vp shared |
| 12 | Per-view scroll state | ✅ Done | scrollY map[view]int |
| 13 | Per-tab scroll state (results) | ✅ Done | results.tabScroll |
| 14 | Scroll reset on analysis complete | ✅ Done | H-03 fix |
| 15 | Viewport resizes with window | ✅ Done | tea.WindowSizeMsg |
| 16 | Scroll indicators | ✅ Done | renderScrollIndicator |

## Analysis

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 17 | Load file/directory from path | ✅ Done | analyzeView |
| 18 | Select analysis mode (full/brief) | ✅ Done | |
| 19 | Run analysis | ✅ Done | sent to engine |
| 20 | Analysis progress bar | ✅ Done | renderProgress |
| 21 | Analysis completion notification | ✅ Done | analysisCompleteMsg |
| 22 | Clear results | ✅ Done | c key, btn |

## Results Display

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 23 | 12 result tabs | ✅ Done | 0-11 |
| 24 | Summary tab | ✅ Done | Tab 0 |
| 25 | Assumptions tab | ✅ Done | Tab 1 |
| 26 | Verification tab | ✅ Done | Tab 2 |
| 27 | Contradictions tab | ✅ Done | Tab 3 |
| 28 | Trust Chains tab | ✅ Done | Tab 4 (split from SPOFs) |
| 29 | Assumption Impact Analysis tab | ✅ Done | Tab 5 |
| 30 | Blind Spots tab | ✅ Done | Tab 6 |
| 31 | Recommended Controls tab | ✅ Done | Tab 7 |
| 32 | Reports/Exports tab | ✅ Done | Tab 8 |
| 33 | SDRI tab | ✅ Done | Tab 9 |
| 34 | Security Design Review tab | ✅ Done | Tab 10 |
| 35 | Single Points of Trust tab | ✅ Done | Tab 11 (new) |
| 36 | Tab cycle with Tab key | ✅ Done | |
| 37 | Tab stats display | ✅ Done | |

## File Explorer

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 38 | Browse directories | ✅ Done | fileBrowserView |
| 39 | Select files | ✅ Done | |
| 40 | Preview panel | ✅ Done | content view |
| 41 | Tab/Shift+Tab cycle views | ✅ Done | |
| 42 | Directory indicator | ✅ Done | / suffix |
| 43 | Esc back to previous view | ✅ Done | |

## Dashboard

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 44 | Dashboard quick actions | ✅ Done | Uses navigateMsg→router |
| 45 | Session stats | ✅ Done | files analyzed, total results |
| 46 | Version display | ✅ Done | |
| 47 | Welcome/actions on startup | ✅ Done | startupView |

## Search

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 48 | Global search toggle (/) | ✅ Done | Search: █ prompt |
| 49 | Search bar rendering | ✅ Done | statusWarn style in View() |
| 50 | Enter/Esc close search | ✅ Done | |
| 51 | Search hint in bottom bar | ✅ Done | resultsView only |

## Settings

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 52 | Settings view | ✅ Done | |
| 53 | Editable settings | ✅ Done | |
| 54 | Save/apply settings | ✅ Done | |

## Help

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 55 | Help view | ✅ Done | |
| 56 | Key bindings listed | ✅ Done | |
| 57 | Sidebar navigation map | ✅ Done | |
| 58 | Results tab reference | ✅ Done | |

## Bottom Bar

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 59 | View-sensitive hints | ✅ Done | Per-view hints |
| 60 | /=Search on results only | ✅ Done | |
| 61 | Scroll indicator | ✅ Done | renderScrollIndicator |
| 62 | Tab/Arrow/Space hints per view | ✅ Done | |

## Logging

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 63 | Log to ~/.asf/logs/asf.log | ✅ Done | Path fix applied |
| 64 | Log directory auto-created | ✅ Done | ensureRuntimeDirs |
| 65 | Log rotation | ✅ Done | |

## Architecture

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 66 | No focusManager dead code | ✅ Done | Removed |
| 67 | No duplicate routers | ✅ Done | Single Router |
| 68 | No global key handler duplication | ✅ Done | Single handleGlobalKey |
| 69 | No nested viewports | ✅ Done | Shared viewport |
| 70 | No competing sidebars | ✅ Done | Single sidebar |
| 71 | No hidden modal state machines | ✅ Done | Removed |
| 72 | No duplicate analysisCompleteMsg handler | ✅ Done | Removed |
| 73 | No dead code paths | ✅ Done | Cleaned |

## Certification Criteria

| # | Feature | Status | Notes |
|---|---------|--------|-------|
| 74 | Real TTY manual acceptance | ❌ Not Done | PTY partial only |
| 75 | TUI_ROUTE_MAP.md | ✅ Done | This doc |
| 76 | TUI_FEATURE_MAP.md | ✅ Done | This doc |
| 77 | TUI_REAL_USER_ACCEPTANCE.md | ❌ Not Done | |
| 78 | TUI_UX_RECOVERY_REPORT.md | ❌ Not Done | |
| 79 | Per-view viewports | ❌ Not Done | Current: scrollY map |
| 80 | Log path ~/.asf/logs/ | ✅ Done | |
| 81 | Trust Chains/SPOFs split | ✅ Done | Tab 4 + Tab 11 |
| 82 | Tests pass | ✅ Done | go test -count=1 ./... |
| 83 | Build passes | ✅ Done | go build -o asf-tui . |

## Not In Scope (by spec)

- Executive Risk Narratives
- Architect Attention Score
- Portfolio Intelligence
- Cross-Project Governance
- Organizational Portfolio Tracking
- Program Management
- Risk Aggregation
- Per-content-area viewport (deferred)
