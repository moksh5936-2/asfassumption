# ASF0 Productization Acceptance Checklist

## 1. Startup Screen

- [x] Onboarding screen appears when running `asf` with no arguments
- [x] Displays fox logo, "ASF0" title, slogan lines
- [x] Shows "ENTER Start ASF0", "? Help", "q Quit" key hints
- [x] Shows version number
- [x] Pressing `Enter` launches the main workspace
- [x] Pressing `?` opens Help view
- [x] Pressing `q` or `Q` or `Ctrl+C` quits

## 2. CASES Section

- [x] "CASES" section heading visible in sidebar
- [x] "+ New Analysis" entry visible below CASES
- [x] Completed analysis files appear as case entries (e.g., `📁 payroll.yaml`)
- [x] Case entries show file basename only
- [x] Selected case is highlighted in the sidebar
- [x] Empty state: only "+ New Analysis" when no cases exist

## 3. WORK Section

- [x] "WORK" section heading visible in sidebar
- [x] "Review Queue" entry navigates to review view
- [x] "Validation Queue" entry navigates to validation view
- [x] "Reports" entry navigates to reports view
- [x] All three are separate workflow views (not merged with case tabs)

## 4. Case Workspace Lifecycle

- [x] Breadcrumb bar visible when viewing a case: `ASF0 / filename / TabName`
- [x] Breadcrumb shows current case file and active tab name
- [x] Case workspace shows all analysis result tabs
- [x] New analysis → case appears → open case → explore findings flow works

## 5. Viewport Stabilization

- [x] Sidebar width matches content width (26 chars)
- [x] Main content area fills remaining terminal width
- [x] No content renders outside viewport boundaries
- [x] No clipping of top/bottom content
- [x] Scroll percentage shown in hints bar (`Line N–M / total (P%)`)
- [x] Mouse wheel scrolling works
- [x] `↑↓` / `PgUp` / `PgDn` / `Home` / `End` scrolling works

## 6. Security Researcher Guidance

- [x] Contextual guidance text in hints bar for each view:
  - New Analysis, Case Workspace, Review Queue, Validation Queue, Reports, Settings, Help, About, Local AI
- [x] Review Queue: "Human analyst approval workflow for assumptions"
- [x] Validation Queue: "Evidence-backed verification workflow for assumptions"
- [x] Reports: "Generate and export analysis results"

## 7. Professional Empty States

- [x] Review Queue: "No review items. Analysis is fully reviewed."
- [x] Validation Queue: "No validations pending. All findings have been reviewed."
- [x] Reports: "No reports generated. Run an analysis first."
- [x] All empty states provide guidance on next steps

## 8. Additional UX

- [x] Footer key hints work in all views
- [x] Local AI entry visible in sidebar under AI section
- [x] Help and About views accessible from sidebar
- [x] No visual overlap between sidebar and content
- [x] `Tab` toggles between sidebar and content focus
- [x] `q` navigates back, `Q` quits

## Build & Test

- [x] `go build .` succeeds with no errors
- [x] `go vet ./...` passes with no warnings
- [x] `go test ./... -count=1` passes all tests

## Certification

```
ASF0_PRODUCTIZATION_CERTIFIED
```
