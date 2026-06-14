# TUI Layered Conflict Audit

## Overview

Audit of all duplicated/overlapping UI systems in the ASF TUI, identifying which were removed, which remain, and ownership model.

## Router

| System | Location | Status |
|--------|----------|--------|
| Router type | `router.go:3-7` | Single authoritative router |
| NavigateTo | `router.go:15-19` | Push history + set view + sync sidebar |
| NavigateBack | `router.go:22-33` | Pop history or fallback to dashboard |
| CycleSidebar | `router.go:44-47` | Modulo cycle of sidebar selection |
| ActivateSidebar | `router.go:57-63` | Change view to selected sidebar entry |
| activateSidebarTab | `router.go:49-55` | Return tab index for results |
| syncSidebar | `router.go:73-79` | Find first matching sidebar entry |

**Finding:** ONE Router. No duplication. Dashboard `navigateMsg` goes through the same Router.

## Focus / Active View

| System | Location | Status |
|--------|----------|--------|
| Router.currentView | `router.go:4` | Authoritative active view |
| Router.sidebarSel | `router.go:5` | Authoritative sidebar selection |
| ~~focusManager~~ | ~~`app.go:47-50`~~ | **REMOVED v4.0.2** |
| Child model internal state | Various | Local only (menu selection, form input) |

**Finding:** `focusManager` was dead code with `activeView` duplicating `Router.currentView`. Removed. No duplicate activeView fields remain.

## View Rendering

| System | Location | Status |
|--------|----------|--------|
| app.go View() | `app.go:553-594` | Single View method |
| renderContent() | `app.go:695-723` | Switch dispatch per view |
| renderSidebar() | `app.go:618-643` | Single sidebar renderer |
| renderTopBar() | `app.go:595-616` | Single top bar |
| renderBottomBar() | `app.go:645-674` | Single bottom bar with per-view hints |
| Viewport | Single `vp viewport.Model` | With `scrollY map[view]int` manager |

**Finding:** ONE rendering pipeline. No old/new overlay. Each view renders only its own content. The startup view is a full-screen view with sidebar, not a separate shell.

## Keymap Handling

| Key | Handler | Location | Status |
|-----|---------|----------|--------|
| Tab | handleGlobalKey + results/fileBrowser passthrough | `app.go:341-354` | Single owner |
| Shift+Tab | handleGlobalKey + passthrough | `app.go:355-368` | Single owner |
| Enter | Per-view handlers (dashboard, analyze, etc.) | Per view files | Single owner per view |
| / | handleGlobalKey sets searchActive | `app.go:454-460` | Single owner |
| Up/Down | handleGlobalKey (scroll views) + per-view | `app.go:369-382` | Fall-through pattern |
| r | handleGlobalKey (Analyze/Review) | `app.go:408-423` | Single owner |
| f | handleGlobalKey | `app.go:404-407` | Single owner |
| q | handleGlobalKey (back/quit) | `app.go:306-312` | Single owner |
| Esc | handleGlobalKey (back/cancel) | `app.go:316-340` | Single owner |
| ? | handleGlobalKey | `app.go:313-315` | Single owner |

**Finding:** Each key has ONE owning handler. The fall-through pattern (`return false` → per-view handler) is intentional and clean. Global handlers own sidebar/scroll/search; per-view handlers own local actions.

## Data / State Overlap

| Data | Location | Owner |
|------|----------|-------|
| currentFile | `app.go:107` | AppModel (top bar display) |
| analyze.docPath | `analyze.go:152-159` | AnalyzeModel (analysis input) |
| results.result | `results.go:12` | ResultsModel |
| router.currentView | `router.go:4` | Router |
| router.sidebarSel | `router.go:5` | Router |
| dash.selected | `dashboard.go:13` | DashboardModel (local menu only) |
| analyze.selected | `analyze.go:23` | AnalyzeModel (local menu only) |
| searchActive/searchQuery | `app.go:109-110` | AppModel (search overlay) |

**Finding:** No conflicting state ownership. Dashboard `selected` is local-only and does not fight sidebar selection. Search state is owned by AppModel.

## Sidebar ↔ Results Tab Mapping

| Sidebar Entry | View | Tab | Status |
|---------------|------|-----|--------|
| Dashboard | dashboardView | — | Unique |
| File Explorer | fileBrowserView | — | Unique |
| Analyze | analyzeView | — | Unique |
| Summary | resultsView | 0 | Unique |
| Assumptions | resultsView | 1 | Unique |
| Verification | resultsView | 2 | Unique |
| Contradictions | resultsView | 3 | Unique |
| Trust Chains | resultsView | 4 | Unique (split from SPOFs in v4.0.2) |
| Single Points of Trust | resultsView | 11 | Unique (new in v4.0.2) |
| Assumption Impact Analysis | resultsView | 5 | Unique |
| Blind Spots | resultsView | 6 | Unique |
| SDRI | resultsView | 9 | Unique |
| Recommended Controls | resultsView | 7 | Unique |
| Security Design Review | resultsView | 10 | Unique |
| Reports / Exports | resultsView | 8 | Unique |
| Settings | settingsView | — | Unique |
| Help | helpView | — | Unique |
| About | aboutView | — | Unique |

**Finding:** Every sidebar entry maps to a unique (view, tab) pair. No duplicates. Summary and About were added in v4.0.2. Trust Chains and SPOFs were split into separate tabs.

## Ownership Model (Current)

```
AppModel
├── Router           → currentView, sidebarSel, history
├── Layout           → sidebarWidth, topBarHeight, bottomBarHeight
├── Sidebar          → renderSidebar() (driven by Router.sidebarSel)
├── StatusBar        → renderTopBar(), renderBottomBar()
├── Viewport         → shared vp with scrollY map manager
├── FileExplorer     → fileBrowse model (active when view=fileBrowserView)
├── Screens          → Per-view models (dashboard, analyze, results, etc.)
└── EngineBridge     → engine field, commands
```

No child screen renders a full app shell. Each screen renders only its content area.

## Resolved in v4.0.2

| Conflict | Resolution |
|----------|------------|
| focusManager dead code | Removed |
| Trust Chains + SPOFs same tab | Split into tab 4 + tab 11 |
| Missing Summary, About sidebar items | Added |
| Search bar invisible | Added View-level rendering |
| /=Search hint always visible | Moved to results-only |
| Log directory not created | Added asfLogsDir() to ensureRuntimeDirs |
