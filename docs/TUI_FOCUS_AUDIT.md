# TUI Focus Audit

## Test Method

Focus tracked via:
1. Code analysis of `focusManager` struct usage
2. `handleGlobalKey` routing analysis
3. Child model dispatch flow

## Focus Architecture

The TUI uses a `focusManager` struct defined in `app.go:47-50`:

```go
type focusManager struct {
    activeView view
    subFocus   string
}
```

However, this struct is **declared but never functionally used** — the `activeView` field duplicates `Router.currentView`, and `subFocus` is never set or read by any update/render function. The `focusMgr` field exists on `mainModel` but has no behavioral impact.

## Focus Flow

| Action | Focus Target | Mechanism | Verdict |
|--------|-------------|-----------|---------|
| Tab (non-results, non-fileBrowser) | Next sidebar item + activate | `Router.CycleSidebar(1)` + `Router.ActivateSidebar()` | ✅ PASS |
| Shift+Tab (non-results, non-fileBrowser) | Previous sidebar item + activate | `Router.CycleSidebar(-1)` + `Router.ActivateSidebar()` | ✅ PASS |
| Tab on resultsView | Next result tab | Falls through to `resultsModel.Update` | ✅ PASS |
| Tab on fileBrowserView | Toggle preview | Falls through to `fileBrowserModel.Update` | ✅ PASS |
| Enter on dashboard | Select quick action | Falls through to `dashboardModel.Update` → returns `navigateMsg` | ✅ PASS |
| Enter on analyze | Select field / start | Falls through to `analyzeModel.Update` | ✅ PASS |
| Enter on settings | Start editing | Falls through to `settingsModel.Update` | ✅ PASS |
| Enter on file browser | Open folder / select file | Falls through to `fileBrowserModel.Update` | ✅ PASS |
| Esc on settings (editing) | Cancel edit | Falls through (global handler returns false) | ✅ PASS |
| Esc on export (confirming) | Cancel export | Falls through (global handler returns false) | ✅ PASS |
| / on resultsView | Enter search mode | Global handler sets `m.searchActive = true` | ✅ PASS |

## Focus Issues

| ID | Issue | Severity | Detail |
|----|-------|----------|--------|
| FOC-01 | `focusManager` is unused dead code | Medium | The struct is instantiated in `newMainModel` but its fields are never read or written by any function. `activeView` duplicates `Router.currentView`. `subFocus` is always empty. |
| FOC-02 | No visual focus indicator on sidebar items | Low | Sidebar active item is highlighted but there's no keyboard focus indicator (e.g., cursor) showing which item is _currently selected_ vs _currently active_. The `sidebarSel` tracks selection but the rendering just uses different styles for active vs inactive. |
| FOC-03 | ResultsView tab browsing doesn't update sidebar highlight | Low | When using Tab on resultsView to cycle result tabs, the sidebar highlight stays on whichever results entry was last selected via sidebar. The sidebar "Single Points of Trust" entry (tab 4) and "Trust Chains" entry (also tab 4) map to the same tab — ambiguous. |
| FOC-04 | No focus trap enforcement | Low | Currently not an issue since the TUI doesn't have modal dialogs, but the architecture has no mechanism to trap focus within a modal if one were added. |
