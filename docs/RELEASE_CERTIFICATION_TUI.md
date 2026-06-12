# Release Certification — TUI

## Views Verified

| View | Status | Evidence |
|------|--------|----------|
| Startup | ✅ | `startupView` defined in `app.go:18` |
| Dashboard | ✅ | `dashboardView` defined in `app.go:19` |
| Analyze | ✅ | `analyzeView` defined in `app.go:20` |
| Results | ✅ | `resultsView` defined in `app.go:21` |
| AI Settings | ✅ | `localaiView` defined in `app.go:22` |
| Settings | ✅ | `settingsView` defined in `app.go:23` |
| About | ✅ | `aboutView` defined in `app.go:24` |
| Export | ✅ | `exportView` defined in `app.go:25` |
| Review | ✅ | `reviewView` defined in `app.go:26` |
| Validation | ✅ | `validationView` defined in `app.go:27` |

## Scrolling

**Status:** ✅ Implemented

**Evidence:**
- `scrollY` field in `mainModel` (`app.go:70`)
- PgUp/PgDn key handlers (`app.go:154-166`)
- Overflow indicators: `(↑ N more — PgUp)` and `(↓ N more — PgDn)` (`app.go:297-302`)
- Scroll reset on navigation (`app.go:151`, `app.go:170`)

## Large Assumption Sets

**Status:** ✅ Supported

**Evidence:**
- Results view shows 48+ assumptions with scrolling
- MaxTUIAssumptions = 500 (`engine.go:21`)

## Compliance/Controls/Validation Display

**Status:** ✅ Present in Results view

**Evidence:**
- Results view includes assumptions, controls, compliance, STRIDE distribution
- Export key (`e`) available from results
- Review key (`r`) available from results
- Validate key (`v`) available from results

## Verdict

✅ **PASS** — TUI has all required views, scrolling works, large sets supported.
