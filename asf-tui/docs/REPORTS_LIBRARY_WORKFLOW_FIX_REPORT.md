# REPORTS LIBRARY WORKFLOW FIX REPORT

## Previous Behavior

The WORK → Reports sidebar entry acted as an export action. When selected, it showed:

1. A format selection list (14 export formats) as the main screen
2. A confirmation prompt ("Press Y to confirm")
3. An export completion message

This conflated the concept of Reports (a library of already-exported files) with Export (the action of creating a new report from an active case).

## Corrected Behavior

WORK → Reports is now a library/viewer for previously exported reports:

1. Shows "Exported Reports" header with a browsable list of report files
2. Each entry shows: filename, format, date, and size
3. Press Enter to view details: full path, source case, format, creation time, file size
4. Empty state shows guidance explaining how to export from a case workspace
5. No export format prompt is shown as the default screen
6. Reports are discovered by scanning `~/.asf/reports/` and `./reports/` directories
7. Refresh ('r'), delete ('d'), and detail view ('Enter') controls

## Export is Now a Dialog from Case Workspace

The Export action is now triggered by pressing `e` in the case workspace:

1. An overlay dialog appears centered on screen with the format selection list
2. Navigate with up/down, select with Enter, confirm with Y
3. Esc dismisses the dialog at any point
4. After export, the file is saved to the configured output directory
5. The dialog shows completion status
6. Exported files automatically appear in the Reports library

## Files Changed

| File | Change |
|---|---|
| `export.go` | Rewrote `reportsModel` as Reports Library; added `reportEntry`, `scanReportsDirs()`, `scanReportsDir()`, `inferCaseName()`, `reportFormatFromExt()`; changed `asfReportsDir`/`projectReportsDir` from functions to `var` (for testability); removed old export-prompt view/update; added `confirmDelete`/`searching`/`searchQuery` fields; search mode with character input and filtering; delete confirmation (y/n pattern); `esc`/`q` exits detail mode; case name shown in list |
| `app.go` | Added `exportActive`/`exportSelected`/`exportFormat`/`exportDone`/`exportPath`/`exportErr`/`exportConfirm` fields to `mainModel`; changed `'e'` hotkey to set `exportActive` instead of navigating to reportsView; `'e'` in reportsView sets statusMsg guidance; `'/'` in reportsView activates search mode (not global search); esc handler checks `searching`/`confirmDelete` for reportsView; removed old reportsV setup in navigateMsg handler; added export dialog keyboard handling; added `renderExportDialog()` overlay method; updated hints bar for reportsView and caseView |
| `paths.go` | No changes needed (asfReportsDir now defined in export.go) |
| `reports_library_test.go` | 14 acceptance tests (was 12): added `TestReportsLibrary_DeleteCancelsOnN`, `TestReportsLibrary_ExportKeyInReportsShowsMessage`; updated `TestReportsLibrary_SearchFiltersReports` for actual filtering; updated `TestReportsLibrary_DeleteRemovesEntry` for confirmation flow |
| `regression_test.go` | Updated esc handling test for new reportsModel fields (pre-existing) |
| `docs/REPORTS_LIBRARY_ACCEPTANCE.md` | Updated with 4 new manual test steps (delete confirmation, search, export key, case name in list) |

## Export / Reports Separation

| Aspect | Export | Reports |
|---|---|---|
| Purpose | Create new report from active case | Browse, view, manage existing reports |
| Access | Press `e` in case workspace | WORK → Reports sidebar |
| Default screen | Format selection dialog | List of previously exported reports |
| Data source | AnalysisResult + ExportResult() | Filesystem scan of reports dirs |
| After action | File saved to output dir | File appears in Reports list |

## Reports Directory Scanning

Two directories are scanned for report files:

- `~/.asf/reports/` — user reports directory
- `./reports/` — project-local reports directory

Files are listed in reverse chronological order. All file types are included regardless of extension. Format is inferred from the file extension (.pdf → PDF, .json → JSON, etc.). The case name is inferred from the filename by stripping the timestamp suffix.

## Tests Added

| Test | Status |
|---|---|
| `TestReportsLibrary_OpensAsLibrary` | PASS |
| `TestReportsLibrary_NoExportPromptByDefault` | PASS |
| `TestReportsLibrary_EmptyStateExplainsExport` | PASS |
| `TestReportsLibrary_ListsExistingReports` | PASS |
| `TestReportsLibrary_ReportMetadataDisplayed` | PASS |
| `TestReportsLibrary_RefreshReloadsReports` | PASS |
| `TestReportsLibrary_SearchFiltersReports` | PASS |
| `TestReportsLibrary_NavigateUpDown` | PASS |
| `TestReportsLibrary_EnterTogglesDetail` | PASS |
| `TestReportsLibrary_DeleteRemovesEntry` | PASS |
| `TestReportsLibrary_DeleteCancelsOnN` | PASS |
| `TestReportsLibrary_ExportKeyInReportsShowsMessage` | PASS |
| `TestExportAction_StillWorksFromCaseWorkspace` | PASS |
| `TestExportAction_ExportedReportAppearsInReports` | PASS |
| `TestFilePicker_Unrestricted` | PASS |

## Build Validation

| Command | Status |
|---|---|
| `go fmt ./...` | PASS |
| `go vet ./...` | PASS |
| `go test -count=1 ./...` | PASS (19 packages) |
| `go build ./...` | PASS |

## Remaining Limitations

1. **Open ('o')**: Opening a report with the system default app is not implemented in the terminal UI. The `openFile()` function returns an error message. This requires platform-specific code to call `open` (macOS), `xdg-open` (Linux), or `start` (Windows).

2. **Copy Path ('c')**: Copying the report path to clipboard is not implemented. The `copyToClipboard()` function returns an error message. This requires integration with a clipboard library or terminal OSC 52 escape sequence.

3. **Case Name Inference**: The `inferCaseName()` function strips the trailing timestamp pattern (`_YYYYMMDD_HHMMSS`) from filenames. This works for reports generated by the export function but may produce incorrect case names for manually-named files.

4. **Hardcoded Scan**: Reports are only discovered at navigation time (when entering the Reports view) or on explicit refresh ('r'). Reports added externally while the Reports view is open do not appear automatically.

5. **Optional Metadata**: The spec mentions optional checksum/summary fields in the detail pane. These are not stored or displayed currently; only filesystem metadata (name, path, size, mod time) is shown.

## Verdict

**REPORTS_LIBRARY_WORKFLOW_CERTIFIED**

- Reports is a library/viewer ✓
- Reports does not show export prompt by default ✓
- Export still works from active case workspace ✓
- Exported files appear in Reports library ✓
- File picker remains unrestricted ✓
- Search filters reports by name/case name ✓
- Delete requires confirmation ('y'/'n') ✓
- Case name displayed in list view ✓
- Export key in Reports shows guidance message ✓
- All tests pass ✓
- Build passes ✓
