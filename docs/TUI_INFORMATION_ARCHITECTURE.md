# TUI Information Architecture

## Navigation Model

The TUI uses a sidebar-driven navigation model with 16 primary items. Each item maps to either a standalone view or a tab within the results view.

## Sidebar Items

| Index | Item                        | View            | Tab | Description                          |
|-------|-----------------------------|-----------------|-----|--------------------------------------|
| 0     | Dashboard                   | dashboardView   | -   | Quick actions, recent files          |
| 1     | File Explorer               | fileBrowserView | -   | Filesystem browser and selection     |
| 2     | Analyze                     | analyzeView     | -   | Analysis configuration and execution |
| 3     | Assumptions                 | resultsView     | 1   | Extracted assumptions                |
| 4     | Verification                | resultsView     | 2   | Verification status and evidence     |
| 5     | Contradictions              | resultsView     | 3   | Detected contradictions              |
| 6     | Trust Chains                | resultsView     | 4   | Trust dependency chains              |
| 7     | Single Points of Trust      | resultsView     | 4   | SPOF identification                  |
| 8     | Assumption Impact Analysis  | resultsView     | 5   | Impact & priority analysis           |
| 9     | Blind Spots                 | resultsView     | 6   | Coverage gaps and blind spots        |
| 10    | SDRI                        | resultsView     | 9   | Security Design Review Intelligence  |
| 11    | Recommended Controls        | resultsView     | 7   | Control recommendations              |
| 12    | Security Design Review      | resultsView     | 10  | Detailed SDR findings                |
| 13    | Reports / Exports           | resultsView     | 8   | Export and report options            |
| 14    | Settings                    | settingsView    | -   | Application configuration            |
| 15    | Help                        | helpView        | -   | Keyboard reference                   |

## Results Tabs

The results view has 11 tabs:

| Tab | Name                   | Data Source                                 |
|-----|------------------------|---------------------------------------------|
| 0   | Summary                | Aggregate counts across all outputs         |
| 1   | Assumptions            | `r.Assumptions[]`                           |
| 2   | Verification           | `r.VerificationOutput`                      |
| 3   | Contradictions         | `r.Contradictions[]`                        |
| 4   | Trust                  | `r.TrustOutput` (chains + SPOFs)            |
| 5   | Impact                 | `r.ReviewOutput`, `r.TrustOutput.SPOFs`     |
| 6   | Blind Spots            | `r.CoverageOutput`                          |
| 7   | Controls               | `r.Controls[]`                              |
| 8   | Reports                | `r.NarrativeOutput`, `r.ReviewOutput`       |
| 9   | SDRI                   | `r.SDRISummary`, `r.SDRIControls`, coverage |
| 10  | Security Design Review | `r.SDRIDesignFindings`, weaknesses, rem.    |

## Key Mapping

### Global (always available)
| Key        | Action                          |
|------------|----------------------------------|
| Ctrl+C / Q | Force quit                       |
| q          | Go back / navigate to previous   |
| Esc        | Go back / Cancel                 |
| ?          | Toggle help                      |
| Tab        | Cycle sidebar items (forward)    |
| Shift+Tab  | Cycle sidebar items (reverse)    |
| f          | Open file explorer               |
| r          | Run analysis                     |

### Navigation (on scrollable views)
| Key           | Action          |
|---------------|-----------------|
| ↑ / k         | Scroll up       |
| ↓ / j         | Scroll down     |
| PgUp / b      | Page up         |
| PgDn / Space  | Page down       |
| Home / g      | Go to top       |
| End / G       | Go to bottom    |
| Ctrl+U        | Half page up    |
| Ctrl+D        | Half page down  |

### Results View
| Key         | Action                    |
|-------------|---------------------------|
| Tab         | Next result tab           |
| Shift+Tab   | Previous result tab       |
| /           | Search/filter             |
| n/N         | Next/prev match           |
| e           | Export results            |
| c           | Clear results             |
| r           | Open review mode          |
| v           | Open validation mode      |

## Architecture Principles

1. **One navigation owner**: `mainModel` owns `currentView`, `navigateTo()`, `navigateBack()`, history stack, and sidebar selection.
2. **One focus owner**: `focusManager` (lightweight struct in `mainModel`) tracks active view and sub-focus string.
3. **One active screen renderer**: `renderContent()` switches on `m.currentView` — exactly one view renders at a time.
4. **One global key dispatch layer**: `handleGlobalKey()` returns `(handled, model, cmd)`; unhandled keys fall through to child model dispatch.
5. **One layout manager**: `layoutManager` (lightweight struct in `mainModel`) provides sidebar width, top/bottom bar heights.
