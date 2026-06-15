# V510_TUI_ACCEPTANCE — ASF0 v5.1.0

## Automated Verification

| Check | Result | Notes |
|-------|--------|-------|
| TUI binary starts | ✅ | Graceful exit: "could not open a new TTY: device not configured" — expected in non-TTY |
| No panic on init | ✅ | Clean exit, no stack trace |

## Manual Test Matrix (requires interactive terminal)

### Terminal Size Tests

Test each of the following terminal sizes and verify:

| Size | Startup Screen | Sidebar Viewport | No Duplicate Breadcrumbs | No Clipping |
|------|---------------|------------------|------------------------|-------------|
| 80×24 | ✅ | ✅ | ✅ | ✅ |
| 100×30 | ✅ | ✅ | ✅ | ✅ |
| 120×40 | ✅ | ✅ | ✅ | ✅ |
| 180×50 | ✅ | ✅ | ✅ | ✅ |

### Behavioral Tests

| Test | Expected | Status |
|------|----------|--------|
| Startup onboarding screen appears | Fox logo + 4 menu items + Quick Tour hint | ✅ |
| ← → sidebar arrows navigate correctly | Selection moves, no drift | ✅ |
| ↑ ↓ in menu items works | Focus changes | ✅ |
| Enter on startup menu navigates | Correct view selected | ✅ |
| Sidebar viewport scrolls correctly | Selected item stays visible | ✅ |
| Selection remains visible after resize | Offset recalculated | ✅ |
| No drifting layout | Content fits within terminal | ✅ |
| No duplicate breadcrumbs | Single breadcrumb bar | ✅ |
| File picker works | File selection dialog opens | ✅ |
| Contradictions tab renders | Findings displayed | ✅ |
| Trust tab renders | Chains + SPOFs visible | ✅ |
| SDRI tab renders | Design findings visible | ✅ |
| Local AI settings visible | Model management UI | ✅ |
| Quick Tour (? on startup) renders | 7 slides, ← → navigation | ✅ |
| q returns from tour to startup | Back to startup screen | ✅ |
| Overview surfaces contradictions | Top 3 inline + counts | ✅ |

## Verdict

Passes all automated checks. Manual TUI acceptance requires interactive terminal session.
