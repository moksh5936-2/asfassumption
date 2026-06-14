# Local AI Navigation Restore

## Decision Rationale

Local AI was removed as a first-class sidebar tab during the v4.0.0 TUI rebuild (commit `6a25b6c`). The original `localai.go` (580 lines) with full model management, Ollama detection, download progress, and curated UI was kept in git history but excluded from the current navigation tree.

**Product decision:** Local AI is not merely a hidden setting. It is a supported ASF0 capability with intentional curation. Restoring it as a dedicated sidebar tab under the AI section preserves:

- Model discovery and catalog browsing
- Ollama health status
- Download/pull workflow
- Model activation and deletion
- Local AI mode integration with analysis engine

## Sidebar Structure

```
🦊 ASF0
CASES
  + New Analysis
  payroll.yaml
  aws-prod.yaml
  healthcare.yaml
────────────────
WORK
  Review Queue
  Validation Queue
  Reports
────────────────
AI
  🧠 Local AI          ← restored
────────────────
SYSTEM
  Settings
  Help
  About
```

## Route Map

| View | Constant | Sidebar Label | Tab |
|------|----------|---------------|-----|
| Analysis | `analyzeView` | ➕ New Analysis | -1 |
| Case | `caseView` | 📁 {filename} | 0..n |
| Review | `reviewView` | 📋 Review Queue | -1 |
| Validation | `validationView` | ✓ Validation Queue | -1 |
| Reports | `reportsView` | 📦 Reports | -1 |
| **Local AI** | **`localAIView`** | **🧠 Local AI** | **-1** |
| Settings | `settingsView` | ⚙ Settings | -1 |
| Help | `helpView` | ❓ Help | -1 |
| About | `aboutView` | ℹ About | -1 |

## Preserved Local AI Capabilities

- **Ollama detection** — `CheckAvailable()`, `CheckRunning()` on model manager
- **Model list** — catalog of supported ASF models (`SupportedModels` in `model.go`)
- **Model download** — `StartDownload()` with progress polling via `syncProgress`
- **Model selection** — "Set as Active" action updates `config.AI.ActiveModel`
- **Model deletion** — `DeleteModel()` via Ollama CLI
- **Model health** — connection status indicator (Connected / Disconnected)
- **Local AI status** — status bar shows "AI Enhanced" when enabled
- **ASF Engine + Local AI mode** — `ModeASFAndAI` constant, selectable in New Analysis
- **AI-enhanced reports** — `ai.go` prompt builder and response parser
- **AI-enhanced narratives** — narrative generation with AI content
- **AI-enhanced explanations** — explainability engine with AI enrichment
- **Local-only AI mode indicators** — `●` marker on selected mode in analyze view

## Files Modified

| File | Change |
|------|--------|
| `app.go` | Added `localAIView` to view enum, `m.localai` field, init in `newMainModel`, routing in `renderContent()` and `Update()` dispatch, hints bar entries |
| `router.go` | Added `AI` section and `🧠 Local AI` item to `sidebarTreeBase`, added `localAIView` to `NavigateBack` fallback |
| `localai.go` | **New file** — recovered from git history with upgraded visuals (PremiumHeader, Card system) |
| `tui_test.go` | Added 7 Local AI tests, updated sidebar node count from 10→12 |
| `analyze.go` | (No changes needed — `ModeASFAndAI` already present) |

## Tests Added

| Test | What It Verifies |
|------|-----------------|
| `TestNewLocalAIModel` | Model initialization with correct active model, catalog entries, installed/active flags |
| `TestLocalAIViewRender` | `renderContent()` produces non-empty output for `localAIView` |
| `TestLocalAIViewSwitch` | `router.SetView(localAIView)` works, scroll state is tracked |
| `TestLocalAISidebarEntry` | `localAIView` is present in `sidebarVisibleNodes()` |
| `TestLocalAIAnalysisMode` | `ModeASFAndAI` option exists in analyze menu |
| `TestLocalAISidebarNavigation` | `sidebarActivate()` navigates to `localAIView` with correct tab |
| `TestLocalAICasesWorkNavigation` | All CASES/WORK/SYSTEM routes still reachable via sidebar |
| `TestLocalAIRouteDoesNotConflict` | `localAIView` is distinct from all other view constants |

## Preserved AI Infrastructure (Not Modified)

- `model.go` — `ModelManager`, `SupportedModels`, `syncProgress`, Ollama HTTP client
- `ai.go` — `AIEnhancer`, prompt builder, response parser, `mergeAIResults()`
- `config.go` — `AI` config section (Enabled, ActiveModel, InstalledModels)
- `settings.go` — AI Enhancement toggle, Active Model field
- `engine.go` — `ModeASFAndAI` analysis mode, AI integration in `RunAnalysis()`
- `doctor.go` — Ollama diagnostics

## Validation Evidence

```
go build ./...  → clean
go vet ./...    → clean
go test ./...   → all packages pass
```

## Final Verdict

**LOCAL_AI_NAVIGATION_RESTORED**
