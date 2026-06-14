# TUI Scrolling Audit

## Test Method

Scrolling tested via:
1. Code analysis of `handleGlobalKey` scroll key routing
2. Regression test suite (TestScrollKeys, TestPgUpPgDn, TestHomeEnd, TestCtrlUD)
3. Viewport model verification

## Scroll Key Dispatch

| Key | Scope | Handler | Verdict |
|-----|-------|---------|---------|
| ↑/k | Global (content views only: results, help, about) | `m.vp.LineUp(1)` | ✅ PASS |
| ↓/j | Global (content views only: results, help, about) | `m.vp.LineDown(1)` | ✅ PASS |
| PgUp/b | Global (all views) | `m.vp.HalfViewUp()` | ✅ PASS |
| PgDn/Space | Global (all views except Space on resultsView) | `m.vp.HalfViewDown()`; Space passes through on results | ✅ PASS |
| Home/g | Global (all views) | `m.vp.GotoTop()` | ✅ PASS |
| End/G | Global (all views) | `m.vp.GotoBottom()` | ✅ PASS |
| Ctrl+U | Global (all views) | `m.vp.ViewUp()` — full page up | ✅ PASS |
| Ctrl+D | Global (all views) | `m.vp.ViewDown()` — full page down | ✅ PASS |
| Mouse wheel | Global | `tea.MouseWheelUp`/`tea.MouseWheelDown` | ✅ PASS |

## Scroll State Management

| Feature | Implementation | Verdict |
|---------|---------------|---------|
| Single shared viewport | `m.vp` shared across all views | ⚠️ Known debt |
| Per-view scroll memory | `m.scrollY map[view]int` — saved/restored in `navigateTo`/`navigateBack` | ✅ PARTIAL |
| Results per-tab scroll | `m.results.tabScroll map[int]int` — saved/restored in `updateResults` | ✅ PASS |
| Scroll percent display | `viewportScrollPercent()` in bottom bar | ✅ PASS |

## Scroll Position Accuracy

| View | Scroll Save | Scroll Restore | Verified | Verdict |
|------|-------------|----------------|----------|---------|
| Dashboard | `saveScroll()` on navigate | `restoreScroll()` on navigate back | Code | ✅ PASS |
| Analyze | `saveScroll()` on navigate | `restoreScroll()` on navigate back | Code | ✅ PASS |
| Results (per-tab) | `tabScroll[resultTab]` save | `tabScroll[resultTab]` restore | Code | ✅ PASS |
| File Browser | `saveScroll()` on navigate | `restoreScroll()` on navigate back | Code | ✅ PASS |
| Settings | `saveScroll()` on navigate | `restoreScroll()` on navigate back | Code | ✅ PASS |
| Help | `saveScroll()` on navigate | `restoreScroll()` on navigate back | Code | ✅ PASS |
| Export | `saveScroll()` on navigate | `restoreScroll()` on navigate back | Code | ✅ PASS |
| Review | `saveScroll()` on navigate | `restoreScroll()` on navigate back | Code | ✅ PASS |
| Validation | `saveScroll()` on navigate | `restoreScroll()` on navigate back | Code | ✅ PASS |
| About | `saveScroll()` on navigate | `restoreScroll()` on navigate back | Code | ✅ PASS |
| Local AI | `saveScroll()` on navigate | `restoreScroll()` on navigate back | Code | ✅ PASS |

## Issues Found

| ID | Issue | Severity |
|----|-------|----------|
| SCR-01 | Single shared viewport means scroll position is clobbered when navigating away and back if the viewport content height changes | Low |
| SCR-02 | `analysisCompleteMsg` sets `m.vp.YOffset = 0` directly instead of using `m.scrollY[resultsView]` | Low |
| SCR-03 | `fileSelectedMsg` sets `m.vp.YOffset = 0` directly instead of using `m.scrollY[analyzeView]` | Low |
