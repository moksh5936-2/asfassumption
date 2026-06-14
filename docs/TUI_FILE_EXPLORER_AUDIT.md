# TUI File Explorer Audit

## Test Method

File explorer tested via:
1. Code analysis of `filebrowser.go` and key routing
2. Regression test suite verification for file browser key routing
3. CLI analysis engine (indirect verification that file opening works)

## Feature Audit

| Feature | Expected | Actual | Evidence | Verdict |
|---------|----------|--------|----------|---------|
| Open via `f` key | Global `f` handler navigates to fileBrowserView | `app.go:399-402` | ✅ PASS |
| Open via sidebar | Tab to "File Explorer" entry, Enter | Sidebar item #1 ("File Explorer") maps to fileBrowserView | ✅ PASS |
| Navigate up/down | ↑/↓/j/k pass through to child model | Global handler returns false for all non-content views | ✅ PASS |
| Select file with Enter | Passes through to child model | Enter not handled globally | ✅ PASS |
| Open folder with Enter | Passes through to child model | Enter not handled globally | ✅ PASS |
| Parent directory with Backspace | Passes through to child model | Backspace not handled globally | ✅ PASS |
| Toggle hidden files with `.` | Passes through to child model | `.` not handled globally | ✅ PASS |
| Search with `/` | Passes through to child model on fileBrowserView | `/` handler now checks `m.router.currentView == fileBrowserView` | ✅ FIXED |
| Toggle preview with Tab | Passes through to child model on fileBrowserView | Tab handler returns false for fileBrowserView | ✅ PASS |
| Current path display | Breadcrumb in `View()` | `fileBrowserView` renders path | ✅ PASS |
| Supported file filter | YAML, JSON, Markdown, Mermaid, Draw.io, SVG, PDF, DOCX, TXT | Listed in help screen | ✅ PASS |

## Key Handlers (fileBrowserModel)

| Key | Handler | Verified | Verdict |
|-----|---------|----------|---------|
| ↑/k | `selectionUp()` | Code | ✅ PASS |
| ↓/j | `selectionDown()` | Code | ✅ PASS |
| Enter | `openSelected()` | Code | ✅ PASS |
| Backspace | `parentDirectory()` | Code | ✅ PASS |
| `.` | Toggle `showHidden` | Code | ✅ PASS |
| Tab | Toggle `showPreview` | Code | ✅ PASS |
| `/` | Set `searchMode = true` | Code | ✅ PASS |
| `q`/`esc` | Clear error | Code | ✅ PASS |

## Potential Issues

| ID | Issue | Severity | Detail |
|----|-------|----------|--------|
| FE-01 | No sidebar highlight for file explorer | Low | When file explorer is active, sidebar doesn't highlight "File Explorer" if navigated via `f` key. `f` handler doesn't sync sidebar — it calls `navigateTo` which does call `Router.NavigateTo` → `syncSidebar()` | ✅ Resolved: `f` handler uses `m.navigateTo(fileBrowserView)` which wraps `Router.NavigateTo` → `syncSidebar()` |
| FE-02 | File selection doesn't restore previous scroll position | Low | `fileSelectedMsg` sets `m.vp.YOffset = 0` directly | Known debt |
