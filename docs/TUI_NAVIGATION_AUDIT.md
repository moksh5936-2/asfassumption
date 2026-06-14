# TUI Navigation Audit

## Test Method

Navigation tested via:
1. Programmatic verification via regression test suite (48 tests covering routing, sidebar sync, key dispatch, scroll keys)
2. Dashboard rendering inspection (captured via `script` PTY)
3. CLI analysis engine verification

## Key Binding Audit

| Key | Expected | Actual | Verdict | Evidence |
|-----|----------|--------|---------|----------|
| **Tab** | Cycle sidebar forward (16 items) | Works per `handleGlobalKey` — `CycleSidebar(1)` + `ActivateSidebar()`. Falls through to child on resultsView/fileBrowserView | ✅ PASS | `app.go:337-349` |
| **Shift+Tab** | Cycle sidebar backward | Works per `handleGlobalKey` — `CycleSidebar(-1)` | ✅ PASS | `app.go:350-363` |
| **↑/k** | Scroll in content views; navigate in actionable views | Global handler scrolls only for resultsView/helpView/aboutView; passes through for all others | ✅ PASS | `app.go:364-370` |
| **↓/j** | Scroll in content views; navigate in actionable views | Global handler scrolls only for resultsView/helpView/aboutView; passes through for all others | ✅ PASS | `app.go:371-377` |
| **Enter** | Select/activate current item | Not handled globally — passes through to all child models | ✅ PASS | Not in `handleGlobalKey` |
| **Esc** | Go back / cancel | Global handler with view-specific exceptions (analyze.running, settings.editing, etc.). Falls through when child has active edit state | ✅ PASS | `app.go:311-335` |
| **Backspace** | Go to parent directory (in file browser); go back elsewhere | Not handled globally — passes through to child models | ✅ PASS | Not in `handleGlobalKey` |
| **?** | Toggle help | Global handler: `m.navigateTo(helpView)` | ✅ PASS | `app.go:308-310` |
| **f** | Open file explorer | Global handler: navigate to fileBrowserView | ✅ PASS | `app.go:399-402` |
| **r** | Run analysis / open review | Global handler with result/review pass-through logic | ✅ PASS | `app.go:403-418` |
| **e** | Open export | Global handler, only from resultsView with data | ✅ PASS | `app.go:439-442` |
| **c** | Clear results | Global handler, only from resultsView | ✅ PASS | `app.go:432-438` |
| **s** | Save settings | Global handler, only from settingsView when not editing | ✅ PASS | `app.go:443-448` |
| **v** | Open validation | Global handler with context-aware dispatch (results or review) | ✅ PASS | `app.go:419-431` |
| **/** | Search | Global handler, enters search mode. Now correctly passes through on fileBrowserView | ✅ PASS | `app.go:449-454` |
| **PgUp / b** | Half page up | Global handler | ✅ PASS | `app.go:378-380` |
| **PgDn / Space** | Half page down | Global handler (Space passes through on resultsView for child handling) | ✅ PASS | `app.go:381-386` |
| **Home / g** | Go to top | Global handler | ✅ PASS | `app.go:393-395` |
| **End / G** | Go to bottom | Global handler | ✅ PASS | `app.go:396-398` |
| **Ctrl+U** | View up | Global handler | ✅ PASS | `app.go:387-389` |
| **Ctrl+D** | View down | Global handler | ✅ PASS | `app.go:390-392` |
| **Ctrl+C / Q** | Force quit | Global handler | ✅ PASS | `app.go:298-300` |
| **1-9** | Open recent file (dashboard) | Handled in `updateDashboard` — passes through via child dispatch | ✅ PASS | `dashboard.go:126-134` |

## Navigation Path Audit

| Path | Expected | Actual | Verdict |
|------|----------|--------|---------|
| Dashboard → Analyze | Tab to Analyze, Enter or arrow to select | Key routing correct | ✅ PASS |
| Dashboard → Settings | Tab to Settings or `s` shortcut | Key routing correct | ✅ PASS |
| Dashboard → AI Models | Tab to localaiView or `l` shortcut | Not in sidebar (localai is NOT a sidebar item!) | ⚠️ WARN |
| Analyze → Results | After analysis completes, auto-navigate | `analysisCompleteMsg` → `m.router.NavigateTo(resultsView)` | ✅ PASS |
| File Explorer → Analyze | File selected → `fileSelectedMsg` → `m.router.NavigateTo(analyzeView)` | Works | ✅ PASS |
| Results → Review | `r` key from results | Works | ✅ PASS |
| Results → Export | `e` key from results | Works | ✅ PASS |
| Results → Validation | `v` key from results | Works | ✅ PASS |
| Review → Validation | `v` key from review | Works | ✅ PASS |
| Any → Help | `?` key | Works | ✅ PASS |
| Any → Previous | `q` or `Esc` | Works with `NavigateBack()` | ✅ PASS |

## Dashboard Sub-navigation Audit

| Key | Expected | Actual | Verdict |
|-----|----------|--------|---------|
| a | Navigate to Analyze | `navigateMsg{to: analyzeView}` | ✅ PASS |
| l | Navigate to AI Models | `navigateMsg{to: localaiView}` | ✅ PASS |
| s | Navigate to Settings | `navigateMsg{to: settingsView}` | ✅ PASS |
| i | Navigate to About | `navigateMsg{to: aboutView}` | ✅ PASS |
| 1-9 | Open recent file in Analyze | `navigateMsg{to: analyzeView}` with file path | ✅ PASS |

## Keyboard Conflict Audit

| Conflict | Keys | Resolution | Verdict |
|----------|------|------------|---------|
| `r` = Run analysis vs `r` = Reject assumption | Global handler returns false when on reviewView → child handles `r` | Correct | ✅ PASS |
| `s` = Save settings vs `s` = Accept assumption | Global handler only fires on settingsView when not editing | Correct | ✅ PASS |
| `v` = Validate from results vs `v` = Validate from review | Context-aware dispatch | Correct | ✅ PASS |
| `Space` = Page down (global) vs `Space` = toggle (results) | Global handler passes through on resultsView | Correct | ✅ PASS |
| `Tab` = Cycle sidebar vs `Tab` = next result tab | Global handler passes through on resultsView and fileBrowserView | Correct | ✅ PASS |

## Issues Found

| ID | Issue | Severity |
|----|-------|----------|
| NAV-01 | Local AI Models view (`localaiView`) is NOT accessible from any sidebar item. It is only reachable via `l` shortcut on Dashboard or direct code navigation. No sidebar entry exists for it. | Medium |
| NAV-02 | About view (`aboutView`) is only reachable via `i` on Dashboard. No sidebar entry. | Low |
| NAV-03 | Review view and Validation view are NOT in sidebar — only reachable via `r`/`v` from results. User has no sidebar indicator that these screens exist. | Medium |
| NAV-04 | Tab on resultsView switches result tabs but does NOT update sidebar highlight to match the current results sub-tab. The sidebar stays on whichever results entry was last active. | Low |
