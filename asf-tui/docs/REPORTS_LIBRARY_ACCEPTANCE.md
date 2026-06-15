# REPORTS LIBRARY ACCEPTANCE

## Manual Test

### 1. Launch ASF0

```bash
./asf-tui
```

Expected: Startup screen appears.

### 2. Open WORK → Reports

Navigate to sidebar > WORK > Reports (or press `q` to focus sidebar, arrow to Reports, Enter).

Expected: Shows "Exported Reports" header with empty state: "No exported reports yet."

**PASS / FAIL**

### 3. Confirm it does not open export prompt

Expected: No "Select Export Format" prompt shown. No format list. Instead, empty state guidance.

**PASS / FAIL**

### 4. Open a case

Run an analysis or open an existing case file from the sidebar.

Expected: Case workspace with Overview/Assumptions/Verification/Contradictions/Trust/Controls/SDRI tabs appears.

**PASS / FAIL**

### 5. Export report

In case workspace, press `e`.

Expected: Export dialog appears centered with format list. Navigate to a format with up/down, press Enter, press Y to confirm.

Expected: "Export Complete" message with file path.

Press Esc to dismiss.

**PASS / FAIL**

### 6. Return to WORK → Reports

Navigate to sidebar > WORK > Reports.

Expected: The exported report appears in the list with filename, format, date, and size.

**PASS / FAIL**

### 7. Select report and confirm details

Press Enter on the report.

Expected: Details view shows: file name, full path, source case, format, creation time, file size.

**PASS / FAIL**

### 8. Delete report with confirmation

In Reports, select a report, press Enter to view details, press `d`.

Expected: Delete confirmation dialog appears showing filename, path, and warning.

Press `n` or Esc.

Expected: Dialog dismissed, report still appears in list.

Press `d` again, then press `y`.

Expected: Report deleted, removed from list.

**PASS / FAIL**

### 9. Search filters reports

In Reports browse view, press `/`.

Expected: Search bar appears: "Search: █"

Type "security" (or partial filename).

Expected: List filters to show only matching reports.

Press Esc to exit search.

Expected: Search bar disappears, full list restored.

**PASS / FAIL**

### 10. Export key in Reports shows guidance

In Reports view, press `e`.

Expected: Status bar shows "Export is available from an active case..."

**PASS / FAIL**

### 11. Confirm case name appears in list view

In Reports browse view with at least one report.

Expected: Each entry shows "Case: <name>" in the list line.

**PASS / FAIL**

### 12. Confirm file picker still works outside reports directory

Navigate to analysis screen. Press Enter on "+ New Analysis". Select an architecture file from outside the reports directory (e.g., from home directory or /tmp).

Expected: File picker shows files outside the reports directory.

**PASS / FAIL**

---

## Automated Tests

```bash
go test -count=1 -run TestReportsLibrary ./...
```

Expected: All tests pass.

**PASS / FAIL**

---

## Verdict

**REPORTS_LIBRARY_ACCEPTANCE_CERTIFIED**

Only if all manual tests pass and all automated tests pass.
