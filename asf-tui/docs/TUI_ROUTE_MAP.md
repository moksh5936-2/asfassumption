# TUI Route Map — Sidebar → View → Tab Mapping

## One Authoritative Router

`router.go` defines a single `Router` struct with:

- `currentView view` — active screen
- `sidebarSel int` — selected sidebar index
- `history viewHistory` — navigation stack (max 50)

Navigation flows through `mainModel.navigateTo()` / `mainModel.navigateBack()` which wrap router methods with scroll save/restore.

## Sidebar → View → Tab Mapping

| # | Sidebar Entry | View | Tab | Tab Name |
|---|---------------|------|-----|----------|
| 0 | Dashboard | `dashboardView` | — | — |
| 1 | File Explorer | `fileBrowserView` | — | — |
| 2 | Analyze | `analyzeView` | — | — |
| 3 | Summary | `resultsView` | 0 | Summary |
| 4 | Assumptions | `resultsView` | 1 | Assumptions |
| 5 | Verification | `resultsView` | 2 | Verification |
| 6 | Contradictions | `resultsView` | 3 | Contradictions |
| 7 | Trust Chains | `resultsView` | 4 | Trust Chains |
| 8 | Single Points of Trust | `resultsView` | 11 | SPOFs |
| 9 | Assumption Impact Analysis | `resultsView` | 5 | Impact |
| 10 | Blind Spots | `resultsView` | 6 | Blind Spots |
| 11 | SDRI | `resultsView` | 9 | SDRI |
| 12 | Recommended Controls | `resultsView` | 7 | Controls |
| 13 | Security Design Review | `resultsView` | 10 | Security Design Review |
| 14 | Reports / Exports | `resultsView` | 8 | Reports |
| 15 | Settings | `settingsView` | — | — |
| 16 | Help | `helpView` | — | — |
| 17 | About | `aboutView` | — | — |

## Results Tab Detail (12 tabs)

| Tab | Tab Name | Sidebar Entry |
|-----|----------|---------------|
| 0 | Summary | Summary |
| 1 | Assumptions | Assumptions |
| 2 | Verification | Verification |
| 3 | Contradictions | Contradictions |
| 4 | Trust Chains | Trust Chains |
| 5 | Impact | Assumption Impact Analysis |
| 6 | Blind Spots | Blind Spots |
| 7 | Controls | Recommended Controls |
| 8 | Reports | Reports / Exports |
| 9 | SDRI | SDRI |
| 10 | Security Design Review | Security Design Review |
| 11 | SPOFs | Single Points of Trust |

## Navigation Flow

### Tab / Shift+Tab (global)
```
handleGlobalKey
  → m.saveScroll()
  → m.router.CycleSidebar(1/-1)    // move sidebar selection
  → m.router.ActivateSidebar()      // change currentView to selected entry
  → m.router.ActivateSidebarTab()   // get tab index if resultsView
  → m.results.resultTab = tab       // set result tab
  → m.restoreScroll()
```

Tab passes through (return false) when `currentView` is `resultsView` or `fileBrowserView` — those views handle Tab for their own purposes (tab cycling, preview toggle).

### Enter (per-view)
Each view handles Enter independently:
- **startupView**: Select menu action (Analyze, Results, AI, Settings, About, Exit)
- **dashboardView**: Execute quick action (Analyze, Local AI, Settings, About)
- **analyzeView**: Edit path / select mode / start analysis
- **fileBrowserView**: Open folder / select file
- **settingsView**: Start editing
- **reviewView**: Toggle browse/detail
- **exportView**: Confirm selection
- **localaiView**: Show model actions

Enter is NOT handled globally — this prevents sidebar/global navigation conflicts.

### Direct Keys
| Key | View Switch | Handler |
|-----|-------------|---------|
| f | → fileBrowserView | handleGlobalKey |
| r | → analyzeView (or reviewView) | handleGlobalKey |
| ? | → helpView | handleGlobalKey |
| q | ← navigateBack | handleGlobalKey |
| Esc | ← navigateBack/cancel | handleGlobalKey |
| e | → exportView (results only) | handleGlobalKey |
| c | Clear results → analyzeView | handleGlobalKey |

## Scroll State Management

Each view gets its own scroll state via `scrollY map[view]int`:

```
mainModel.navigateTo(to):
  → m.saveScroll()            // save current view's scroll to scrollY[currentView]
  → m.router.NavigateTo(to)   // push history, set currentView
  → m.restoreScroll()         // load scrollY[to] into vp.YOffset

mainModel.navigateBack():
  → m.saveScroll()            // save current view's scroll
  → m.router.NavigateBack()   // pop history
  → m.restoreScroll()         // load previous view's scroll
```

Results tab changes also save/restore via `results.tabScroll map[int]int`.

## Ownership Model

```
AppModel
├── Router           → currentView, sidebarSel, history
├── Layout           → sidebarWidth, topBarHeight, bottomBarHeight
├── Sidebar          → renderSidebar() (driven by Router.sidebarSel)
├── StatusBar        → renderTopBar(), renderBottomBar()
├── Viewport         → shared vp with scrollY map manager
├── FileExplorer     → fileBrowse model (active when view=fileBrowserView)
├── Screens          → Per-view models (each renders only its content)
└── EngineBridge     → engine field, analysis commands
```

No child screen renders a full app shell. Each screen renders only content within the main viewport.
