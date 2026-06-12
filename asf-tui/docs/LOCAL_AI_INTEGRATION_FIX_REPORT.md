# Local AI/Ollama Integration Fix Report

## Summary

Rewrote the entire Local AI integration layer for ASF's TUI, addressing all 7
defects documented in `LOCAL_AI_BUG_REPORT.md`. The fix covers model
discovery, download persistence, timeout handling, graceful fallback, and
config save reliability.

## Files Changed

| File | Change |
|------|--------|
| `model.go` | Rewritten — Ollama HTTP API client with context timeouts, `CheckRunning`, `GetVersion`, `ListInstalledAPI`, `IsModelInstalled`, `GenerateWithTimeout` |
| `localai.go` | Rewritten — model catalog merge, two-group display, `refreshFromOllama`, config save on download/delete/activate |
| `ai.go` | Rewritten — `Enhance` checks model installed, returns specific errors, uses timeout |
| `engine.go` | Updated — AI failure prepends warning to Summary, sets `ModeASFOnly` |
| `settings.go` | Updated — `applyChange` calls `config.Save()` after every change (auto-save) |
| `doctor.go` | Updated — New "Local AI" diagnostics section with running/version/models |
| `ai_test.go` | New — 12 tests mocking Ollama API with `httptest` |

## Defects Fixed

1. **Model discovery via CLI parsing** → HTTP API (`GET /api/tags`), 30s
   timeout, avoids `ollama list` parsing fragility

2. **No timeout on model list** → Explicit context deadlines on all API calls;
   health checks use 2s, API queries use 30s, generation uses 120s

3. **Download not persisted** → `config.Save()` called immediately after
   `ollama pull` completes

4. **Active model not persisted** → `config.Save()` called immediately on "Set
   as Active" and on model deletion (clears if deleted model was active)

5. **Settings required manual save** → `applyChange()` auto-saves; no more `s`
   key required

6. **No graceful AI fallback** → Three check points before generation (binary
   found, server running, model installed); failure prepends warning to
   summary, switches `AnalysisMode` to `ModeASFOnly`, preserves base results

7. **No AI diagnostics** → Doctor command shows binary found, server running,
   version, installed models, active model status

## Test Coverage (12 new tests)

| Test | What it verifies |
|------|-----------------|
| `TestModelManagerCheckRunning_Offline` | False for unresponsive server |
| `TestListInstalledAPI_Empty` | Zero models from empty server |
| `TestListInstalledAPI_WithModels` | 3 models, custom-model strips `:latest` |
| `TestIsModelInstalled_Recommended` | Catalog model detected |
| `TestIsModelInstalled_NonCatalog` | Non-catalog model detected; nonexistent returns false |
| `TestAIEnhance_Timeout` | Long request doesn't hang forever |
| `TestAIEnhance_FallbackKeepsBaseResults` | TotalAssumptions, assumptions preserved on AI failure |
| `TestGetVersion` | Version parsed from /api/version |
| `TestCheckRunning_Online` | True for responsive mock |
| `TestGenerateWithTimeout` | Error on unreachable server |
| `TestActiveModelPersistence` | Active model and installed list survive save/load cycle |
| `TestConfigSaveOnAIEnable` | AI.Enabled persisted to YAML file |

## Architecture

```
┌──────────────┐   config.Save()    ┌──────────┐
│  localai.go  │ ──────────────────►│ config   │
│  settings.go │                    │ (YAML)   │
└──────┬───────┘                    └──────────┘
       │
       │ HTTP (GET /api/tags, /api/generate, /api/version)
       ▼
┌──────────────┐
│  model.go    │
│  (Ollama     │
│   HTTP API)  │
└──────┬───────┘
       │
       ▼
┌──────────────┐     on error       ┌──────────┐
│  ai.go       │ ──────────────────►│ engine   │
│  (Enhance)   │                    │ (fallback│
└──────────────┘                    │  to base)│
                                    └──────────┘
```

All 12 tests pass, all 11 existing test packages pass, and `go build ./...`
compiles clean.
