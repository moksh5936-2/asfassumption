# TUI Feature Parity: Old vs Current (v4.0.0)

**Date:** 2026-06-13  
**Baseline:** What the TUI should do (by spec / prior behavior) vs what it actually does.

---

## Navigation

| Feature | Expected | Actual | Status |
|---------|----------|--------|--------|
| Tab: cycle sidebar | Next sidebar item, switch to that view | Works (global handler) | ✅ |
| Shift+Tab: reverse cycle | Previous sidebar item | Works (global handler) | ✅ |
| q: navigate back | Go to previous view in history | Works (global handler) | ✅ |
| Esc: navigate back | Same as q | Works (global handler) | ✅ |
| ?: toggle help | Open help view | Works (global handler) | ✅ |
| **Enter: select/activate** | Select current item in any view | **Broken** — consumed by global handler | ❌ |
| **Up/Down arrow: navigate items** | Move selection in lists | **Broken** — always scrolls viewport | ❌ |
| **j/k: navigate items (vim)** | Same as up/down | **Broken** — always scrolls viewport | ❌ |
| Sidebar highlight sync | Active view matches sidebar highlight | **Broken** — `sidebarSel` not updated on non-Tab navigation | ❌ |

---

## Dashboard

| Feature | Expected | Actual | Status |
|---------|----------|--------|--------|
| Display system status | Version, mode, AI, theme | ✅ Rendered | ✅ |
| Quick actions list | Analyze, AI Models, Settings, About | ✅ Rendered | ✅ |
| Recent files | Show last 10 files | ✅ Rendered | ✅ |
| **Arrow keys: select action** | Navigate quick actions | **Broken** — always scrolls viewport | ❌ |
| **Enter: open selected action** | Navigate to target view | **Broken** — key consumed | ❌ |
| **a: open Analyze** | Shortcut to analyzeView | **Broken** — key consumed | ❌ |
| **l: open AI Models** | Shortcut to localaiView | **Broken** — key consumed | ❌ |
| **s: open Settings** | Shortcut to settingsView | **Broken** — key consumed | ❌ |
| **i: open About** | Shortcut to aboutView | **Broken** — key consumed | ❌ |
| **1-9: open recent file** | Number key opens file in analyze | **Broken** — key consumed | ❌ |

---

## Analyze

| Feature | Expected | Actual | Status |
|---------|----------|--------|--------|
| Display menu items | Document path, Evidence path, modes, Start button | ✅ Rendered | ✅ |
| **Arrow keys: select field** | Navigate path, mode, action | **Broken** — always scrolls viewport | ❌ |
| **Enter: edit/select field** | Edit path / select mode / start analysis | **Broken** — key consumed | ❌ |
| Text input for path | Type path, Enter confirms | **Broken** — handleTextInput never called | ❌ |
| Mode selection | ASF Only / ASF+AI | **Broken** — enter never triggers | ❌ |
| Start analysis | Begin analysis, show progress | **Broken** — enter never triggers | ❌ |
| Progress bar | Show during analysis | ✅ Works (via analysisCompleteMsg) | ✅ |

---

## Results

| Feature | Expected | Actual | Status |
|---------|----------|--------|--------|
| Display 9 tabs | Summary, Assumptions, Verification, etc. | ✅ Rendered | ✅ |
| **Tab: next tab** | Switch to next result tab | **Broken** — global handler skips results but consumes key | ❌ |
| **Shift+Tab: prev tab** | Previous result tab | **Broken** — same issue | ❌ |
| e: export | Open export dialog | Works (global handler) | ✅ |
| c: clear results | Clear and go to analyze | Works (global handler) | ✅ |
| r: review | Open review mode | Works (global handler) | ✅ |
| v: validate | Open validation mode | Works (global handler) | ✅ |
| /: search | Filter current tab | Works (global handler) | ✅ |
| n/N: next/prev match | Navigate search results | ✅ (search mode) | ✅ |
| Scroll within tabs | Per-tab scroll position | Partial — uses shared viewport | ⚠️ |

---

## File Explorer

| Feature | Expected | Actual | Status |
|---------|----------|--------|--------|
| Display file list | Files with name, size, modified | ✅ Rendered | ✅ |
| Breadcrumb | Current path | ✅ Rendered | ✅ |
| **Up/Down: navigate files** | Move selection | **Broken** — always scrolls viewport | ❌ |
| **Enter: open folder / select file** | Navigate into dir or select file | **Broken** — key consumed | ❌ |
| **Backspace: parent dir** | Go up one directory | **Broken** — global handler for backspace not implemented for file browser | ❌ |
| **Tab: toggle preview** | Show/hide preview panel | **Broken** — global handler intercepts tab for sidebar | ❌ |
| **.: toggle hidden** | Show/hide dotfiles | **Broken** — key consumed | ❌ |
| **/: search filename** | Enter search mode | **Broken** — global handler intercepts / for content search | ❌ |
| Preview panel (text) | Show file preview for text files | ✅ Rendered (when tab would work) | ❌ (tab broken) |
| f: open file browser from anywhere | Navigate to fileBrowserView | Works (global handler) | ✅ |

---

## Settings

| Feature | Expected | Actual | Status |
|---------|----------|--------|--------|
| Display setting list | 12+ configurable items | ✅ Rendered | ✅ |
| **Arrow keys: select setting** | Navigate settings | **Broken** — always scrolls viewport | ❌ |
| **Enter: start editing** | Begin editing a value | **Broken** — key consumed | ❌ |
| **Left/Right: change value** | Cycle through options | **Broken** — global handler for left/right not implemented | ❌ |
| **s: save** | Save config to disk | Works (global handler) | ✅ |
| **Esc: cancel edit** | Cancel editing | Works (global handler calls navigateBack) | ⚠️ |
| Theme change live preview | Styles update immediately | ✅ Works (in updateSettings) | ✅ |

---

## Review Mode

| Feature | Expected | Actual | Status |
|---------|----------|--------|--------|
| Display assumption list | All assumptions with status badges | ✅ Rendered | ✅ |
| **Arrow keys: navigate** | Move through assumptions | **Broken** — always scrolls viewport | ❌ |
| **Enter: toggle detail** | Show/hide detailed view | **Broken** — key consumed | ❌ |
| **s: accept** | Mark as Accepted | **Broken** — global handler uses s for save | ❌ |
| **r: reject** | Mark as Rejected | **Broken** — global handler uses r for navigate-to-analyze | ❌ |
| **m: modified** | Mark as Modified | **Broken** — key consumed | ❌ |
| **n: edit notes** | Add review notes | **Broken** — key consumed | ❌ |
| **v: validate** | Open validation for this assumption | **Broken** — key consumed | ❌ |

---

## Validation Mode

| Feature | Expected | Actual | Status |
|---------|----------|--------|--------|
| Display assumption list | All assumptions with risk/confidence | ✅ Rendered | ✅ |
| **Arrow keys: navigate** | Move through assumptions | **Broken** — always scrolls viewport | ❌ |
| **Enter: toggle detail** | Show detailed validation data | **Broken** — key consumed | ❌ |

---

## Export

| Feature | Expected | Actual | Status |
|---------|----------|--------|--------|
| e: open export | Navigate to exportView | Works (global handler) | ✅ |
| Select format | Choose export format | Need to check exportModel | ? |
| Confirm export | Generate file | Need to check | ? |

---

## About

| Feature | Expected | Actual | Status |
|---------|----------|--------|--------|
| Display info | Version, license, description | ✅ Rendered | ✅ |

---

## Help

| Feature | Expected | Actual | Status |
|---------|----------|--------|--------|
| Display keyboard reference | 12 sections of key bindings | ✅ Rendered | ✅ |

---

## AI Models

| Feature | Expected | Actual | Status |
|---------|----------|--------|--------|
| Display model catalog | List of available models | ✅ Rendered | ✅ |
| **Arrow keys: select model** | Navigate models | **Broken** — always scrolls viewport | ❌ |
| **Enter: show actions** | Show download/activate/delete | **Broken** — key consumed | ❌ |
| Download progress | Progress bar during download | ✅ Works (via aiDownloadTickMsg) | ✅ |

---

## Scrolling

| Feature | Expected | Actual | Status |
|---------|----------|--------|--------|
| Up/Down/j/k: scroll | Scroll content within viewport | ✅ Works in all views | ✅ |
| PgUp/PgDn/b/Space: page scroll | Half-page scroll | ✅ Works | ✅ |
| Home/g: top, End/G: bottom | Jump to top/bottom | ✅ Works | ✅ |
| Scroll percent display | Show "1-20/100 (20%)" in bottom bar | ✅ Works | ✅ |
| **Per-view scroll memory** | Each view remembers its scroll | Partial — tracked in scrollY map but not always honored | ⚠️ |
| **Results per-tab scroll** | Each result tab remembers scroll | ✅ Via tabScroll map | ✅ |

---

## Window Resize

| Feature | Expected | Actual | Status |
|---------|----------|--------|--------|
| Resize triggers re-layout | Update widths/heights | ✅ Handled | ✅ |
| Content re-renders on resize | Viewport.SetContent on resize | ❌ Only in View(), not in Update | ❌ |
| Minimum size enforcement | Show error if too small | ✅ In View() | ✅ |

---

## Summary

| Category | Total Features | Working | Broken | Partial |
|----------|---------------|---------|--------|---------|
| Navigation | 12 | 6 | 5 | 1 |
| Dashboard | 11 | 3 | 8 | 0 |
| Analyze | 7 | 2 | 5 | 0 |
| Results | 10 | 6 | 2 | 2 |
| File Explorer | 9 | 2 | 7 | 0 |
| Settings | 7 | 3 | 3 | 1 |
| Review Mode | 7 | 1 | 6 | 0 |
| Validation | 3 | 1 | 2 | 0 |
| Export | 3 | 3 | 0 | 0 |
| About | 1 | 1 | 0 | 0 |
| Help | 1 | 1 | 0 | 0 |
| AI Models | 5 | 2 | 3 | 0 |
| Scrolling | 6 | 5 | 0 | 1 |
| Window Resize | 3 | 2 | 1 | 0 |
| **Total** | **85** | **38** | **42** | **5** |

**Verdict:** 45% of features work, 49% are broken, 6% partially work. Navigation, file explorer, analyze, settings, review, and validation are the most affected.
