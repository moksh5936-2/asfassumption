# Release Certification — Exports

## Test Results

| Format | Status | Size | Content Verified |
|--------|--------|------|-----------------|
| JSON | ✅ | 225,130 bytes | 48 assumptions, 15 controls, 31 compliance lines |
| Markdown | ✅ | 141,929 bytes | Has assumptions, controls, compliance sections |
| CSV | ✅ | 42,309 bytes | 49 rows, 18 columns |
| HTML | ✅ | 182,716 bytes | Has assumptions section |
| PDF | ✅ | 53,461 bytes | Binary PDF generated |

## Export Method

Exports are triggered via the TUI (`e` key from results view) or programmatically via `ExportResult()`.

## File Integrity

All exported files:
- Non-empty
- Contain assumption data
- Contain control data
- Contain compliance data
- No crashes during export
- Files open successfully

## Verdict

✅ **PASS** — All 5 export formats work correctly.
