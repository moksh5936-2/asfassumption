# TUI Regression Audit

**Date:** 2026-06-13  
**Version:** v4.0.0  
**Auditor:** AI-assisted code review

---

## Executive Summary

The TUI has a fundamental event-routing bug introduced alongside the MC-style file explorer integration. The `mainModel.handleKeyMsg()` function intercepts **all** `tea.KeyMsg` events and returns `(m, nil)` for both handled and unhandled keys, preventing child models (`dashboardModel`, `analyzeModel`, `fileBrowserModel`, `settingsModel`, etc.) from ever receiving keyboard input. All per-view keyboard navigation is dead code.

---

## Root Cause: Event Routing Architecture

### Problem

In `app.go:165-166`:
```go
case tea.KeyMsg:
    return m.handleKeyMsg(msg)
```

This hard-return means **no KeyMsg ever reaches the child model dispatch** in `app.go:236-261`. Every keyboard event is consumed by `handleKeyMsg`, which handles ~20 keys and silently drops the rest.

### Impact

The following views have their own `Update` methods with keyboard handling — all are **dead code** for KeyMsg:

| View | File | Key handlers that never fire |
|------|------|------------------------------|
| `startupModel` | `startup.go:77-107` | up/k, down/j, enter |
| `dashboardModel` | `dashboard.go:27-55` | up/k, down/j, enter, a, l, s, i |
| `analyzeModel` | `analyze.go:57-100` | up/k, down/j, enter |
| `fileBrowserModel` | `filebrowser.go:107-191` | up/k, down/j, enter, backspace, tab, ., / |
| `settingsModel` | `settings.go:119-169` | up/k, down/j, enter, left/h, right/l, esc |
| `resultsModel` | `results.go:42-53` | tab, shift+tab |
| `localaiModel` | `localai.go:179-258` | up/k, down/j, enter, esc |
| `reviewModel` | `review.go:160-224` | up/k, down/j, enter, s, r, m, n, v |
| `validationModel` | `validation.go:192-213` | up/k, down/j, enter |

### Manifestation

1. **Arrow keys scroll the viewport instead of navigating** — because `handleKeyMsg` intercepts `up/k` and `down/j` for viewport scroll (lines 306-311) in ALL views
2. **Enter never selects anything** — not in dashboard, analyze, file browser, settings, review, validation, or startup
3. **File browser can't navigate** — up/down/enter/backspace all consumed by global handler
4. **Tab doesn't switch result tabs** — the global `tab` handler skips results (line 287) but still returns `(m, nil)`, consuming the key; `resultsModel.Update` never processes it
5. **Sidebar selection doesn't sync** — `sidebarSel` is not updated when navigating via shortcuts (f, r, v, e)
6. **Dashboard shortcut keys don't work** — `a`, `l`, `s`, `i` are not in `handleKeyMsg`'s switch; they silently fall through to `return m, nil`
7. **File browser `q` and `esc` never fire** — global handler catches them first

---

## Viewport Ownership

### Problem

A single `viewport.Model` (`m.vp`) is shared across all 12 views. Every view's content is rendered into the same viewport.

### Code

```go
// app.go:509-511 (in View())
m.vp.Width = m.mainWidth()
m.vp.Height = m.mainHeight()
m.vp.SetContent(content)
```

### Issues

1. **Scroll position cross-contamination** — scrolling in the analyze view corrupts scroll state when switching to dashboard, settings, etc.
2. **`scrollY map[view]int` tracked but not always honored** — `navigateTo()` and `navigateBack()` save/restore via `scrollY`, but `analysisCompleteMsg` (line 219) and `fileSelectedMsg` (line 231) set `m.vp.YOffset = 0` directly, bypassing the scroll map
3. **Results tab scroll tracking is fragile** — `updateResults` (results.go:176-189) saves/restores using `m.vp.YOffset` directly, which aliases with the shared viewport scroll
4. **No per-view viewport means variable-height content is problematic** — sidebar + topbar + bottombar dimensions change with view, but viewport size is fixed at layout time

---

## Sidebar State Management

### Problem

`sidebarSel` is managed independently of `currentView`, and not synced when navigation occurs via non-sidebar mechanisms.

### Code

```go
// app.go:113-122
var sidebarViews = []view{
    dashboardView, analyzeView, resultsView, fileBrowserView,
    localaiView, settingsView, aboutView, helpView,
}
```

### Issues

1. **`navigateTo()` does not update `sidebarSel`** — pressing `f` sets `currentView = fileBrowserView` but `sidebarSel` stays at whatever it was before
2. **Sidebar active highlight is wrong** — `renderSidebar()` (line 562) checks `i == m.sidebarSel && viewForSidebar(i) == m.currentView`, so if sidebarSel is stale, the wrong item appears highlighted
3. **Tab cycling always goes to dashboard first** — `sidebarSel` starts at 0 (Dashboard), so pressing Tab from any non-results view cycles to dashboard regardless of current view

---

## File Explorer Integration

### Problem

The file explorer (`filebrowser.go`) was added as an MC-style model embedded in `mainModel`, but its integration is incomplete.

### Issues

1. **Dual key handling** — `fileBrowserModel.handleKey()` at filebrowser.go:115 handles keys for file navigation, but `handleKeyMsg` at app.go:265 intercepts them first. The file browser model's `Update()` is called from `updateFileBrowser` at filebrowser.go:354, but only for non-KeyMsg
2. **No focus isolation** — there's no concept of "file browser has focus." All keys go through the global handler
3. **Enter on file browser never selects a file** — `handleKeyMsg` doesn't handle enter for file selection
4. **File browser shares the single viewport** — scroll position bleeds between file browser and other views

---

## Missing Focus Model

### Problem

There is no focus tracking in the TUI. Every view implicitly "has focus" when it's the `currentView`, but the key routing system doesn't respect view-specific key handling.

### Issues

1. **Global keys shadow view-specific keys** — `q`, `esc`, `?`, `/` are handled globally and can't be overridden per-view
2. **No way to determine "what has focus within a view"** — e.g., in settings, is the user editing a value or navigating the list? The `editing` flag exists but it's never checked by the global handler
3. **Search mode overrides everything but is poorly scoped** — `searchActive` is a mainModel flag but only affects results view

---

## View History

### Problem

The `viewHistory` stack tracks navigation for "back" navigation (`q` / `esc`), but has edge cases.

### Issues

1. **Startup view is in the history** — if user navigates startup → dashboard → analyze, then presses `q`, they go back to dashboard (correct), but pressing `q` again goes to startup (which is not useful — you can't go back to startup meaningfully)
2. **History not bounded** — could grow unbounded with repeated navigation
3. **`esc` calls `navigateBack()` but `q` also calls `navigateBack()`** — identical behavior, confusing UX

---

## Resize Handling

### Issues

1. **`WindowSizeMsg` updates viewport dimensions but doesn't trigger content re-render** — `m.vp.SetContent(content)` is only called in `View()`, not in response to resize
2. **No minimum-size enforcement in Update** — the minimum size check (60x12) only happens in `View()`
3. **Sidebar width is hardcoded to 23** — no proportional sizing

---

## Startup View Layout

### Issues

1. **Startup view bypasses the viewport** — `viewStartup()` is called from `renderContent()`, which feeds into the viewport. But startup has its own centering logic that uses `m.width - 4`, which doesn't account for sidebar width. The startup view content inside the viewport is narrower than intended
2. **Startup view's `q` behavior differs** — `q` quits from startup but navigates back elsewhere

---

## Summary of Issues by Severity

### Critical
- **E1**: All keyboard input is consumed by `handleKeyMsg`; child models never receive KeyMsg — TUI is non-interactive for navigation, selection, and data entry
- **E2**: Arrow keys only scroll viewport in all views — cannot navigate menus, file lists, settings, or dashboard

### High
- **H1**: Shared single viewport causes scroll position cross-contamination between unrelated views
- **H2**: File browser is completely non-functional (cannot navigate, select, or open files)
- **H3**: Sidebar selection not synced with navigation via keyboard shortcuts
- **H4**: Result tab switching is broken (tab key consumed by global handler, never reaches resultsModel)

### Medium
- **M1**: No focus model — ambiguous key handling, especially for views with editing modes (settings, review)
- **M2**: View history unbounded; startup view pollutes back navigation
- **M3**: Window resize doesn't re-layout content until next `View()` call

### Low
- **L1**: Sidebar width hardcoded; no proportional sizing
- **L2**: Dashboard shortcut keys (`a`, `l`, `s`, `i`) don't work
- **L3**: Review/validation modes inaccessible because Enter never navigates to detail view

---

## Code Map

| Issue | File | Lines |
|-------|------|-------|
| Global key handler intercepts all KeyMsg | `app.go` | 165-166, 265-385 |
| Arrow keys -> viewport scroll in all views | `app.go` | 306-311 |
| Child Update methods never receive KeyMsg | `app.go` | 236-261 (dead dispatch) |
| Shared viewport across all views | `app.go` | 93-95, 509-511 |
| Sidebar selection not synced | `app.go` | 422-427 (`navigateTo`) |
| Results tab key consumed | `app.go` | 286-305 (tab handling) |
| File browser key handling dead code | `filebrowser.go` | 107-191 |
| Dashboard key handling dead code | `dashboard.go` | 27-55 |
| Analyze key handling dead code | `analyze.go` | 57-100 |
| Settings key handling dead code | `settings.go` | 119-169 |
| Review key handling dead code | `review.go` | 160-224 |
| Validation key handling dead code | `validation.go` | 192-213 |
| LocalAI key handling dead code | `localai.go` | 179-258 |
| Startup key handling dead code | `startup.go` | 77-107 |

---

## Conclusion

The TUI regression has a single root cause: **the global key handler (`handleKeyMsg`) returns for all KeyMsg, preventing sub-model dispatch**. This affects every interactive feature. The regression was introduced when the MC-style file explorer was integrated and the key routing was centralized in `handleKeyMsg` without wiring unhandled keys back to sub-models.

The fix requires restructuring the event dispatch so that:
1. Global keys (ctrl+c, q, ?) are still handled at the top level
2. All other keys are forwarded to the current view's Update method
3. Each view owns its own keyboard handling, focus management, and viewport (or at minimum, scroll state)

---

## Per-Screen Viewport Audit

All 11 content screens (startup, dashboard, analyze, results, filebrowser, localai, settings, about, export, review, validation, help) share a **single** `viewport.Model` at `mainModel.vp`.

### Viewport Configuration

- **Owner:** `mainModel.vp`
- **Width:** `m.mainWidth()` = `m.width - m.sidebarWidth() - 2`
- **Height:** `m.mainHeight()` = `m.height - 3`
- **Sidebar width:** 23 (hardcoded)
- **Resize:** `WindowSizeMsg` updates vp dimensions; content re-set in `View()` on next render cycle

### Per-Screen Breakdown

| Screen | Own Viewport? | Scrolling? | Content Length | Resize Behavior |
|--------|---------------|------------|----------------|-----------------|
| **Startup** | No (shared vp) | No (fits in viewport) | Short (menu + fox art) | Re-rendered in View(); centering uses `m.width - 4` (ignores sidebar) |
| **Dashboard** | No (shared vp) | No (list fits viewport) | 8-10 items + recent files | Re-rendered in View() |
| **Analyze** | No (shared vp) | Rare (form fields ~8 items) | Short form + results | Re-rendered in View() |
| **Results** | No (shared vp) | Yes (long assumption lists) | Depends on result size (can be 1000+ lines) | Tab-switching saves/restores scroll via `tabScroll` map |
| **File Browser** | No (shared vp) | Yes (directory listings) | Depends on directory size | Re-rendered in View(); file list scrolls with viewport |
| **LocalAI** | No (shared vp) | Yes (model list can be long) | Depends on installed models | Re-rendered in View() |
| **Settings** | No (shared vp) | Yes (20+ settings) | ~25 items | Re-rendered in View() |
| **About** | No (shared vp) | No (short content) | ~20 lines | Re-rendered in View() |
| **Export** | No (shared vp) | No (14 format items) | ~20 lines | Re-rendered in View() |
| **Review** | No (shared vp) | Yes (long assumption lists) | Depends on result size | Mode-switch (browse/detail) re-renders in View() |
| **Validation** | No (shared vp) | Yes (long assumption lists) | Depends on result size | Re-rendered in View() |
| **Help** | No (shared vp) | Yes (long help text) | ~125 lines of keyboard docs | Re-rendered in View() |

### Viewport Issues Found

1. **Shared viewport means scroll position cross-contamination** — navigating from a long view (results, review) to a short view (about, export) and back loses scroll position. The `scrollY map[view]int` mitigation is incomplete (not all views save/restore).
2. **No per-view viewport** — every view's content is forced into the same viewport model, making it impossible to have viewport-specific settings like max-height, scrollbar visibility, or key bindings.
3. **Startup view center calculation ignores sidebar** — `viewStartup()` uses `m.width - 4` for its centering title style, but the actual content area width is `m.width - sidebarWidth() - 2`. The startup content appears 23 characters narrower than expected when the sidebar is open.
4. **Sidebar scroll bar shows "All" when help/about views are displayed** — because the viewport is shared, the scroll percentage display in the bottom bar reflects the scroll state of the shared viewport, not the current content.
5. **File browser and results both use the viewport for scrolling** — but they have very different content heights and scroll requirements. No isolation.

### Viewport Scroll Key Handling (after fix)

| Action | Scope | Handled In |
|--------|-------|------------|
| `up/k` | Results, Help, About only (scroll) | `handleGlobalKey` |
| `down/j` | Results, Help, About only (scroll) | `handleGlobalKey` |
| `pgup/b` | All views (page up) | `handleGlobalKey` |
| `pgdown/space` | All views except Results (space reserved) | `handleGlobalKey` |
| `ctrl+u` | All views (half page up) | `handleGlobalKey` |
| `ctrl+d` | All views (half page down) | `handleGlobalKey` |
| `home/g` | All views (top) | `handleGlobalKey` |
| `end/G` | All views (bottom) | `handleGlobalKey` |
| `up/k` | Dashboard, Analyze, Settings, FileBrowser, Review, Validation, Startup, LocalAI, Export | Child model (navigation) |

---

## Applied Fix: Event Routing Restructure

### Change Summary

**File:** `app.go`

#### 1. `handleKeyMsg` → `handleGlobalKey`

**Before:**
```go
func (m mainModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
    // returns for ALL keys, blocking child dispatch
}
```

**After:**
```go
func (m mainModel) handleGlobalKey(msg tea.KeyMsg) (bool, tea.Model, tea.Cmd) {
    // returns (true, modified_model, cmd) for keys consumed globally
    // returns (false, m, nil) for keys that should fall through to child
}
```

#### 2. `Update` method fall-through

**Before:**
```go
case tea.KeyMsg:
    return m.handleKeyMsg(msg)  // always returns, blocking child dispatch
```

**After:**
```go
case tea.KeyMsg:
    if m.searchActive {
        return m.handleSearchInput(msg)
    }
    handled, model, cmd := m.handleGlobalKey(msg)
    if handled {
        return model, cmd
    }
    // falls through to child dispatch below
```

#### 3. Esc exception handling

`esc` is handled globally (`navigateBack()`) EXCEPT when the current view has active state that needs esc:

| View | Exception Condition | Effect |
|------|-------------------|--------|
| analyze | `m.analyze.running \|\| m.analyze.inputMode != ""` | Cancel analysis or text input |
| settings | `m.settings.editing` | Cancel editing |
| localai | `m.localai.showActions` | Close action menu |
| export | `m.exportV.showConfirmation \|\| m.exportV.done` | Cancel confirmation |
| review | `m.review.editing` | Cancel note editing |

#### 4. Key routing table (after fix)

| Key | Handled Globally? | Falls Through To Child? | Notes |
|-----|-------------------|------------------------|-------|
| `ctrl+c`, `Q` | Yes (quit) | — | Always quits |
| `q` | Yes (back/quit) | — | Quit from startup, navigate back otherwise |
| `?` | Yes (help) | — | Navigate to help view |
| `esc` | Yes (back) | If view has active state | See esc exception table |
| `tab` | Yes (sidebar) | Results, FileBrowser | Sidebar cycling for all other views |
| `shift+tab` | Yes (sidebar) | Results, FileBrowser | Reverse sidebar cycling |
| `up`, `k` | Content views only | Navigation views | Content: results/help/about scroll. Nav: children handle |
| `down`, `j` | Content views only | Navigation views | Content: results/help/about scroll. Nav: children handle |
| `pgup`, `b` | Yes (scroll) | — | All views |
| `pgdown`, ` ` | Yes (scroll) | Results (space only) | All views; space reserved on results |
| `ctrl+u`, `ctrl+d` | Yes (scroll) | — | All views |
| `home`, `g` | Yes (scroll) | — | All views |
| `end`, `G` | Yes (scroll) | — | All views |
| `f` | Yes (file browser) | — | Navigate to file browser |
| `r` | Yes | ReviewView | Results→review, review rejected, else→analyze |
| `v` | Yes | — | Results/review→validation |
| `c` | Yes | — | Clear results |
| `e` | Yes | — | Export results |
| `s` | SettingsView only | All other views | Settings save; review Accept |
| `/` | Yes (search) | — | Always activates search |
| `enter` | — | All views | Children handle selection/activation |
| `h`, `l` | — | All views | Settings uses left/right for values; file browser doesn't |
| `a`, `i` | — | Dashboard | Dashboard shortcut keys |
| `.` | — | FileBrowser | Toggle hidden files |
| `backspace` | — | FileBrowser, Review | Parent directory, delete note char |
| `n`, `N` | — | All views | Search next/prev (when search active); review note (when not) |

#### 5. Sidebar selection sync

**Before:** `navigateTo()` and `navigateBack()` never updated `m.sidebarSel`.

**After:** Both functions iterate `sidebarViews` to find a matching view and update `m.sidebarSel`. Non-sidebar views (startup, export, review, validation) preserve the previous selection.

#### 6. Window resize propagation

**Before:** `tea.WindowSizeMsg` handler returned early with `return m, nil`, never reaching child dispatch.

**After:** Removed the early return so resize events fall through to the child models. All views now receive `WindowSizeMsg` and can re-layout if needed.

#### 7. Pre-existing bug fix: PgDn key string

**Before:** The handler checked `"pgdn"` but Bubble Tea uses `"pgdown"` for `KeyPgDown.String()`. PgDn key press never triggered a scroll.

**After:** Changed to `"pgdown"` to match Bubble Tea's key string representation.

---

## Remaining Issues (Post-Fix)

| ID | Issue | Severity | Status |
|----|-------|----------|--------|
| E1 | All keys consumed by global handler | Critical | **FIXED** |
| E2 | Arrow keys only scroll in all views | Critical | **FIXED** |
| H1 | Shared single viewport scroll contamination | High | Mitigated (scrollY map) — full fix requires per-view viewport |
| H2 | File browser non-functional | High | **FIXED** (keys now reach child model) |
| H3 | Sidebar selection not synced | High | **FIXED** |
| H4 | Result tab switching broken | High | **FIXED** |
| M1 | No focus model | Medium | Pending — needs formal focus tracking |
| M2 | View history unbounded | Medium | Pending — needs history limit |
| M3 | Window resize doesn't re-layout | Medium | **FIXED** (now propagates to children) |
| L1 | Sidebar width hardcoded | Low | Pending |
| L2 | Dashboard shortcuts don't work | Low | **FIXED** (keys now reach dashboard model) |
| L3 | Review/validation inaccessible | Low | **FIXED** (enter now reaches child models) |
