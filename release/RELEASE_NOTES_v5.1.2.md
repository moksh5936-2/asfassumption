# ASF0 v5.1.2 — Reports Library Workflow Stabilization

**v5.1.2** is a patch release that fixes real-world usability gaps in the Reports Library introduced in v5.1.1. Search filtering, delete confirmation flow, case name display, and the export key hint are all now functional. No engine behavior was changed.

## What's Fixed

### Reports Library
- **Search filtering** (`/`): Entering search mode now actually filters the report list by filename and caseName. Previously `/` entered search mode but the filter was a no-op that reset the search on exit.
- **Delete confirmation** (`d`): First press sets a confirmation state showing the filename, path, and a warning. `y` confirms deletion, `n`/`esc` cancels. Previously deletion was immediate with no confirmation.
- **Case name in list view**: Each report line now displays `Case: %s` so you can see which case a report belongs to without opening it.
- **Export key inside Reports** (`e`): Sets a status message explaining that export is only available from an active case workspace. Previously the `e` key was a no-op inside Reports.
- **Esc handler**: Now checks `searching` and `confirmDelete` state before navigating back. Previously pressing `Esc` during search or delete confirmation would navigate away.
- **`/` guarded during confirmDelete**: Cannot enter search mode while a delete confirmation is active.

### Docs & Tests
- 2 new tests added for delete cancellation and export key messaging
- Search and delete tests updated for actual filtering / confirmation flow
- Acceptance checklist expanded from 8 to 12 scenarios
- Fix report updated with current status and limitations

## Upgrade Instructions

```bash
# macOS/Linux (curl | bash)
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash

# Or download directly:
# https://github.com/moksh5936-2/asfassumption/releases/download/v5.1.2/ASF-v5.1.2-darwin-arm64
# https://github.com/moksh5936-2/asfassumption/releases/download/v5.1.2/ASF-v5.1.2-darwin-amd64
# https://github.com/moksh5936-2/asfassumption/releases/download/v5.1.2/ASF-v5.1.2-linux-amd64
# https://github.com/moksh5936-2/asfassumption/releases/download/v5.1.2/ASF-v5.1.2-linux-arm64
# https://github.com/moksh5936-2/asfassumption/releases/download/v5.1.2/ASF-v5.1.2-windows-amd64.exe
```

## Checksums (SHA-256)
```
95fe5f3d24b51634e36a20546e5db865398dde1a34b4d5e5f03c5f530b5f45ae  ASF-v5.1.2-darwin-amd64
8ef10a7a5fbfb3232e101fbb9833ecb43b9a991b9fa02092d9919eb74909b1e4  ASF-v5.1.2-darwin-arm64
9c7941e05243922120cff373983553ef5b3de59084f675e3eba47b22f0779435  ASF-v5.1.2-linux-amd64
53c30c96b44230d8db37733cd495f11dcd4533acd8bb2fd840a13b770cf01523  ASF-v5.1.2-linux-arm64
4c4f71c867b6deeae0b88ba7d0b04d5d58944386b8321e91261ce0f859fdcd11  ASF-v5.1.2-windows-amd64.exe
```

## Files Changed

### Fixes
- `asf-tui/app.go` — `handleGlobalKey`: `/` guarded during confirmDelete; `e` sets statusMsg; `esc` checks searching/confirmDelete; `navigateMsg` clears search/confirm states
- `asf-tui/export.go` — `reportsModel`: search filtering, delete confirmation, case name in list, export key message

### Tests
- `asf-tui/reports_library_test.go` — 15 tests (3 new/updated)
  - `TestReportsLibrary_SearchFiltersReports` — updated for real filtering
  - `TestReportsLibrary_DeleteRemovesEntry` — updated for confirmation flow
  - `TestReportsLibrary_DeleteCancelsOnN` — new
  - `TestReportsLibrary_ExportKeyInReportsShowsMessage` — new

### Docs
- `asf-tui/docs/REPORTS_LIBRARY_ACCEPTANCE.md` — expanded to 12 scenarios
- `asf-tui/docs/REPORTS_LIBRARY_WORKFLOW_FIX_REPORT.md` — updated for v5.1.2

### Release Engineering
- All version references bumped to v5.1.2
- Build validation, binary verification, checksums documented

## Known Limitations
- Export output goes to configured `Output.Directory` (or default reports dir)
- `openFile()` and `copyToClipboard()` are stubs returning errors (platform-specific implementations not yet implemented)
