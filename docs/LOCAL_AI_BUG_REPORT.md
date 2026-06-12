# Local AI / Ollama Integration Bug Report

## Root Cause Summary

The ASF Local AI/Ollama integration has 7 distinct defects that together make it unreliable. The core issues are:

1. **Installed models never persisted to disk** — `config.Save()` is never called after download completes or model is deleted. The `InstalledModels` list lives only in memory.

2. **No Ollama state sync on startup** — ASF trusts the stale `config.InstalledModels` list without querying `/api/tags`. Models installed externally (via `ollama pull`) are invisible until the user downloads them through ASF.

3. **No timeouts on any Ollama calls** — HTTP client uses `http.DefaultClient` (no timeout), `exec.Command` has no context. AI generation can hang forever.

4. **No model availability check before AI enhancement** — If the configured model doesn't exist in Ollama, `Enhance()` makes a failing HTTP call with no clear error message.

5. **AI failure silently discarded** — When AI enhancement fails, the error is logged but not surfaced to the user. Base results are returned but the UI shows no warning.

6. **Settings don't auto-save** — User must press `s` to persist AI enabled/disabled or active model changes.

7. **Doctor shows no Ollama health info** — Only checks if `ollama` binary exists, but doesn't check if the Ollama server is running, what version, or how many models are installed.

## Detailed File-by-File Findings

### `model.go` — Ollama wrapper

| Issue | Location | Detail |
|-------|----------|--------|
| No timeout on HTTP calls | `Generate()` line 200 | Uses `http.Post` → `http.DefaultClient` with no timeout |
| No streaming for generation | `Generate()` line 196 | `stream: false` — blocks until full response |
| API tags never used | nowhere | `GET /api/tags` for querying installed models is never called |
| CLI list unused | `ListInstalled()` line 51 | Function exists but is never called |
| No download timeout | `StartDownload()` line 120 | `exec.Command` with no context/cancel |
| Download progress fragile | `StartDownload()` | Parses CLI output which may vary by Ollama version |

### `localai.go` — TUI model management

| Issue | Location | Detail |
|-------|----------|--------|
| No config save after download | lines 97-101 | `installedModels` updated in-memory only |
| No config save after delete | lines 156-161 | `installedModels` updated in-memory only |
| No Ollama refresh on view load | `newLocalAIModel()` line 28 | Uses `cfg.AI.InstalledModels` as-is |
| No "other installed" models section | entire file | Only shows ASF catalog, not user's other Ollama models |
| Download progress polling | `pollDownloadCmd()` | Works but no timeout if goroutine hangs |

### `ai.go` — AI enhancement pipeline

| Issue | Location | Detail |
|-------|----------|--------|
| No model check before call | `Enhance()` line 38 | Only checks `CheckAvailable()` (binary exists), not model installed |
| No timeout wrapper | entire file | No context propagation for cancellation |

### `engine.go` — Analysis orchestrator

| Issue | Location | Detail |
|-------|----------|--------|
| AI failure silently discarded | lines 170-174 | `err` is checked but no warning preserved in result |
| No timeout on AI enhancement | line 170 | `enhancer.Enhance()` could hang forever |
| No fallback message | lines 170-174 | Base results returned without indication AI was attempted |

### `settings.go` — Config UI

| Issue | Location | Detail |
|-------|----------|--------|
| No auto-save | entire file | User must press `s` to save changes |
| Active model not validated | `applyChange()` line 184 | Sets any string as active model without checking if installed |

### `doctor.go` — System diagnostics

| Issue | Location | Detail |
|-------|----------|--------|
| No Ollama server check | `checkDep()` line 257 | Only checks if `ollama` binary exists with `--version`, not if server is running |
| No model count | doctor output | No query of `/api/tags` for installed models |
| No active model validation | doctor output | Shows `cfg.AI.ActiveModel` without checking if actually installed |

### `main.go` / `config.go` — Startup and persistence

| Issue | Location | Detail |
|-------|----------|--------|
| No startup Ollama sync | `main()` line 87 | Config loaded but no `/api/tags` query to reconcile installed models |
| config.Save called once | `main()` line 103 | Only on graceful TUI exit; crash loses all AI state |

## Bug Manifestations

### Bug 1: "ASF Engine + Local AI" hangs
**Causes:**
- `Generate()` uses `http.Post` with no timeout → if Ollama is running but unresponsive (e.g., model not loaded), the HTTP call blocks forever
- `progress` channel in `RunAnalysis()` sends progress updates but the main goroutine waits on the blocking HTTP call
- No `context.WithTimeout` or `context.WithCancel` anywhere in the call chain

### Bug 2: Downloaded model disappears after restart
**Causes:**
- `localai.go` never calls `config.Save()` after download completes
- `localai.go` never calls `config.Save()` after model deletion
- `config.AI.InstalledModels` is updated in memory (local slice) but never written to disk
- On restart, `config.Load()` returns the old (empty or stale) installed models list

### Bug 3: Externally installed models not detected
**Causes:**
- ASF only shows models from `config.AI.InstalledModels`
- `ModelManager.ListInstalled()` exists but is never called
- No startup or periodic sync with `ollama list` or `GET /api/tags`

### Bug 4: AI failure loses base results (partial)
**Causes:**
- When `Enhance()` returns an error, `mergeAIResults` is skipped
- Base result is returned but with no indication that AI was attempted
- User may think AI ran successfully when it was silently skipped

## Severity Assessment

| Bug | Severity | Impact |
|-----|----------|--------|
| AI analysis hangs | **Critical** | Blocks all analysis in ASF+AI mode; user must force-quit |
| Models lost on restart | **High** | Makes model management unusable; must re-download every session |
| Externally installed models invisible | **High** | Users who use `ollama pull` directly can't select those models in ASF |
| AI failure silent | **Medium** | User doesn't know AI enhancement failed |
| Settings not auto-saved | **Medium** | Loses AI config changes if app crashes |
| Doctor missing AI diagnostics | **Low** | Support/debugging harder |
