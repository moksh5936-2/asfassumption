# V511 — Source Verification

## Feature Audit

| Feature | Present | Verification |
|---|---|---|
| ASF0 TUI | ✓ | `main.go:28` — "ASF0 v%s starting" |
| Result split-pane layout | ✓ | `results.go:214-340` — `viewResults()` with left/right pane |
| No dropdown/inline details | ✓ | `results.go:569-574` — Enter sets `detailFocus`, no inline expansion |
| Selection-follow viewport | ✓ | `results.go:136-167` — `ensureSelectedVisible()` |
| Per-tab result list state | ✓ | `results.go:11-21` — `tabState` per tab |
| Selected item always visible | ✓ | `ensureSelectedVisible` invariant enforced at every movement |
| SDRI 69+ navigable | ✓ | `regression_test.go:1114-1147` — `Test69SDRIFindingsNavigation` |
| Trust 100+ navigable | ✓ | `regression_test.go:1278-1311` — `Test100TrustChainsNavigation` |
| Enter focuses details pane | ✓ | `results.go:569-574` — `detailFocus = true` on Enter |
| Esc returns focus to list | ✓ | `results.go:585-588` — `detailFocus = false` on Esc |
| Detail pane scrolls independently | ✓ | `results.go:506-510` — detail focus routes ↑↓ to `DetailOffset` |
| Resize preserves selection | ✓ | `app.go` — `ensureSelectedVisible` called on `WindowSizeMsg` |
| Sidebar viewport fix retained | ✓ | Sidebar tests pass |
| Duplicate breadcrumb fix retained | ✓ | Breadcrumb tests pass |
| Modal file picker fix retained | ✓ | File picker tests pass |
| Semantic contradiction engine | ✓ | Contradiction tests pass |
| Local AI | ✓ | Local AI tests pass |

## Source commit

```
b86351a release: ASF0 v5.1.0 — add certification and verification docs
```

## Verdict

Current source contains all intended fixes for v5.1.1.

**SOURCE_VERIFIED** — proceed to version bump.
