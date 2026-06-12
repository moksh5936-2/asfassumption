# ASF v2.1.1 — Production Readiness Audit

**Date:** 2026-06-12
**Scope:** Full 13-section audit covering installer, commands, architecture, TUI, AI, performance, claims, and D2C readiness.
**Method:** Independent code review, binary analysis, test execution, claim verification. No source of truth is trusted without verification.

---

## Verdict: REJECTED

**Score: 52/100** (D-)

| Section | Weight | Score | Grade |
|---------|--------|-------|-------|
| 1. Go Native Validation | 10% | 100% | A |
| 2. Single Binary Validation | 5% | 100% | A |
| 3. Installer Audit | 10% | 88% | B+ |
| 4. Version Consistency | 5% | 20% | F |
| 5. Release Pipeline | 10% | 35% | F |
| 6. TUI Audit | 10% | 45% | F |
| 7. Local AI Audit | 5% | 50% | F |
| 8. AI Execution Audit | 10% | 35% | F |
| 9. Engine/Test Coverage | 10% | 25% | F |
| 10. Claim Quality | 10% | 55% | F |
| 11. Performance | 10% | 30% | F |
| 12. D2C Readiness | 10% | 35% | F |
| 13. Critical Bug Hunt | 5% | 70% | C |

**Overall: 52/100 — REJECTED**

9 BLOCKING issues found. 15+ HIGH issues found. 2 FALSE claims in documentation. No production deployment should proceed without addressing all BLOCKING and HIGH items.

---

## BLOCKING Issues (Must Fix Before Any Release)

### B1. Dashboard shows wrong version
**File:** `asf-tui/dashboard.go:57`
**Detail:** `version := "v0.1.0"` hardcoded. TUI dashboard displays v0.1.0 instead of actual v2.1.1. Constant `ASFVersion` (v2.1.1) defined in `license.go` exists but is unused here.
**Fix:** Replace with `version := ASFVersion` (after import).

### B2. Goroutine leak — progress channel never closed
**File:** `asf-tui/analyze.go:191`, `asf-tui/engine.go:183`
**Detail:** `runAnalysisCmd` creates `for range progress` goroutine. Producer (`RunAnalysis`) sends final progress update but never calls `close(progress)`. With buffer of 10, the 11th update blocks the producer permanently. After every analysis run, one goroutine leaks.
**Fix:** Add `close(progress)` as deferred final action in `RunAnalysis`. Or use `sync.WaitGroup`. Or redesign with context cancellation.

### B3. Stale results survive re-analysis
**File:** `asf-tui/results.go:425`
**Detail:** `updateResults` gates data transfer on `m.results.result == nil`. After first analysis, `result` is non-nil, so re-analysis never updates the display. User must quit and restart to see fresh analysis.
**Fix:** Remove the nil gate or set `m.results.result = nil` when analysis starts.

### B4. Export confirmation does not export
**File:** `asf-tui/export.go:771`
**Detail:** Pressing 'y' in export confirmation sets `m.done = true` and shows "Export Complete" with empty `exportPath`. No call to `ExportResult()` is made. The export is entirely fake.
**Fix:** Wire the confirmation to actually call `ExportResult()` and populate `exportPath`.

### B5. HMAC license secret hardcoded in binary
**File:** `asf-tui/license.go:57`
**Detail:** `hmac.New(sha256.New, []byte("asf-enterprise-secret-2024"))` — HMAC key is in source, extractable from any binary via `strings`. Anyone can forge licenses.
**Fix:** Use asymmetric signatures (RSA/Ed25519) with public key in binary, private key held offline. Or accept that this is obfuscation only.

### B6. Python API and Go TUI use separate config systems
**File:** `asf/settings.py:22`, `asf-tui/config.go:11`
**Detail:** Python REST API loads `asf.config.yaml` from CWD. Go TUI loads `~/.config/asf/config.yaml`. Completely different paths and schemas. User configuring via TUI gets zero benefit on API side. Splits the product into two independent tools that share nothing but the name.
**Fix:** Unify config path and schema. Or document this as a known limitation if the Python backend is deprecated.

### B7. PDF/DOCX parsing claim is FALSE
**File:** `README.md:51`, `parser.go:61-62`
**Detail:** README claims "PDF documents (.pdf) — Raw text analysis" and "Word documents (.docx) — Raw text analysis". Implementation at `parser.go:61-62` dispatches `.pdf`/`.docx` to `parseTextFile()`, which calls `os.ReadFile()` and treats raw bytes as a string. PDF and DOCX are binary formats; raw reads produce garbage text. No PDF or DOCX parsing library is used.
**Fix:** Remove the claim, implement actual PDF/DOCX parsing (e.g., `unidoc`, `pdfcpu`, `apache/tika`), or clearly document that PDF/DOCX support is binary-only.

### B8. AI risk refinement is a stub
**File:** `ai.go:160-162`
**Detail:** `parseRiskRefinements()` returns `nil` unconditionally. The AI prompt asks the LLM for risk refinements and the response parsing code exists, but the critical function is completely unimplemented. AI-powered risk refinement is advertised but non-functional.
**Fix:** Implement the parsing, or remove the claim and the prompt section.

### B9. Progress bar oscillates infinitely
**File:** `asf-tui/analyze.go:84`
**Detail:** When `progress >= 100`, progress is reset to 90. Next tick hits 100 again, reset to 90. Creates infinite 90→100→90→100 loop that never resolves until `analysisCompleteMsg` arrives. During long analyses, user sees a broken progress bar oscillating forever at the end.
**Fix:** Cap progress at 99 until `analysisCompleteMsg` arrives, then jump to 100.

---

## HIGH Issues (Fix Before Production)

### Claims & Documentation
- **C1** — README claims "20 passing tests": actual count is 53. Outdated.
- **C2** — COMMAND_COVERAGE_AUDIT.md still references issues fixed in v2.1.1 as open.
- **C3** — INSTALLER_AND_COMMAND_RELIABILITY_REPORT.md still references v2.0.1 throughout.
- **C4** — README.md:212 has stale example `# Expected: ASF v2.0.0`.
- **C5** — README.md:469 has stale reference `ASF v2.0.0+`.

### TUI
- **T1** — AI assumptions section (results.go:131) is joined to body string BEFORE it is appended to sectionViews at line 136. AI-enhanced findings are computed but never rendered.
- **T2** — View history stack (`app.go:196`) is pushed to but never popped. `handleBack()` uses hardcoded parent views instead of stack. Dead code + wrong back navigation.
- **T3** — `analyzeModel` uses value receivers that mutate state (`analyze.go:56`). Works by accident because caller reassigns, but any missed reassignment silently loses state.
- **T4** — `formatTime` (localai.go:145) can slice-panic on strings shorter than 10 chars after parse failure.
- **T5** — `ensureRuntimeDirs()` error discarded at `main.go:86`. Runtime dirs may not exist.
- **T6** — Config path migration (`config.go:58`) silently swallows read/write errors. Could lose config silently.

### AI & Local Models
- **A1** — `ListInstalledAPI` error swallowed (`localai.go:82`). If Ollama is unreachable, all models show as not installed with no error to user.
- **A2** — `InstalledModels` config field written to YAML but never read on startup (`localai.go:202`). Model state not restored from config after restart.
- **A3** — Gzip read error silently discarded (`parser.go:78`). Corrupt `.drawio` files produce confusing XML parse errors instead of "file corrupt" message.
- **A4** — Gzip decompression failure falls through silently (`parser.go:76`). Original compressed data passed to XML parser.

### Performance
- **P1** — Unbounded memory growth: assumptions, verifications, gaps append without limit (`engine.go:212,331`). Large documents can OOM.
- **P2** — Unnecessary temp file I/O: architecture text is written to temp file then re-read (`engine.go:133`). Doubles disk I/O per analysis.
- **P3** — Blocking startup: TUI blocks on Ollama HTTP call for up to 30s with no progress indicator (`main.go:88` via `localai.go:58`).

### D2C Readiness
- **D1** — Config not persisted until clean exit (`main.go:95`). Crash loses all session changes.
- **D2** — No file-based logging (`engine.go:16`). Debug output only goes to stderr. Post-mortem debugging impossible.
- **D3** — No upgrade/update path. Users on v2.0.1 cannot discover v2.1.1 without checking GitHub manually.
- **D4** — No SIGTERM handler (`app.go:141`). SIGTERM kills process with no cleanup during long operations.
- **D5** — No differentiated exit codes (`main.go:91`). CI/CD cannot distinguish failure modes.
- **D6** — Review notes is fake (`review.go:191`). Toggles between "" and "Review pending". Actual text input fields declared but unused.

### Engine & Tests
- **E1** — Zero test coverage for drawio, mermaid, svg, image parsers.
- **E2** — Zero tests for empty file handling across ALL parsers.
- **E3** — Zero tests for binary garbage input across ALL parsers.
- **E4** — Evidence loading errors silently skipped (`analyzer/analyzer.go:127`).
- **E5** — Document parsing errors silently skipped (`analyzer/analyzer.go:105`).

---

## MEDIUM Issues (Fix Before Next Release)

- `settings.go:63` — export_dir and active_model settings cannot be edited via TUI.
- `config.go:11` — Duplicate Theme/FoxStyle fields (General vs Appearance).
- `localai.go:145` — formatTime slice-panic on strings shorter than 10 chars.
- `export.go:31` — No overwrite confirmation for same-second exports.
- `model.go:331` — HTTP client 30s timeout vs context timeout semantics confusing.
- `model.go:267` — Download progress heuristic stalls at 90%.
- `model.go:145` — N+1 HTTP pattern risk in IsModelInstalled.
- `ai.go:63` — Unbounded prompt string for LLM (multi-MB for large analyses).
- `engine.go:167` — AI and native analysis could run in parallel but are serial.
- `engine.go:132` — archDesc pointer prevents GC between runs.
- `analyze.go:182` — Progress tick keeps firing after completion.
- Doctor checks tesseract but no graceful fallback if missing.
- No `--data-dir` flag or env var override for cache directory.
- First run creates no config file on disk (user sees nothing at `~/.config/asf/config.yaml`).
- `_ = strings.Contains` dead code in `asf/verification/engine.go:329`.
- `_ = i` unused loop variable in `engine.go:371`.

---

## PASS Sections

### Section 1: Go Native Validation — PASS (100%)
Zero Python exec calls in Go code. No `pip`, `venv`, `python` in binary path. `exec.Command` only calls `ollama`, `tesseract`, or `asf` binary.

### Section 2: Single Binary Validation — PASS (100%)
8.7MB Mach-O arm64 binary. Links only 4 macOS system frameworks. `CGO_ENABLED=0` safe. No Python/Node/Java runtime required.

### Section 3: Installer — PASS with note (88%)
14/15 scenarios pass. SHELL unset + `set -e` causes hard abort when `setup_path` returns 1. Installer test script passes 25/25.

### Section 4: Version Consistency — FAIL (20%)
`dashboard.go:57` hardcodes `"v0.1.0"`. Plus 3 stale references in README, 2 in docs.

### Section 13: Critical Bug Hunt — PASS (70%)
No `panic()` calls, no `recover()`, no TODO/FIXME/HACK in production code. Ignored errors found are already covered by other sections.

---

## Summary Recommendation

**Do not ship v2.1.1 as-is.**

The installer and command dispatch fixes from the previous sprint are solid, but the production readiness audit reveals that the TUI, AI enhancement, export workflow, and documentation claims are not production-grade.

### Minimum fix list to reach APPROVED WITH CONDITIONS (score ≥75):

1. Fix B1 (dashboard version) — 5 minutes
2. Fix B3 (stale results) — 30 minutes
3. Fix B4 (export fake) — 1 hour
4. Fix B7 (PDF/DOCX claim) — remove claims from README, 5 minutes
5. Fix B8 (AI risk refinement) — remove claim, 5 minutes
6. Fix B9 (progress bar oscillation) — 15 minutes
7. Remove stale doc references (C1-C5) — 15 minutes
8. Add SIGTERM handler (D4) — 30 minutes
9. Add file logging (D2) — 1 hour
10. Fix config persist on crash (D1) — save on every change, 30 minutes

Total: ~4 hours to reach APPROVED WITH CONDITIONS.

### Additional effort for full APPROVED (score ≥90):

- B2 (goroutine leak) — major refactor of progress reporting
- B5 (license secret) — switch to asymmetric crypto
- B6 (config unification) — product architecture decision
- T1 (AI assumptions rendering) — fix rendering order
- P1-P2 (memory + I/O) — significant refactor
- E1-E3 (test coverage) — weeks of work
- P3 (blocking startup) — async initialization
