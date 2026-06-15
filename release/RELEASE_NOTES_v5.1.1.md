# ASF0 v5.1.1

Result-tab UX stabilization release — split detail pane, selection-follow viewport, large-list navigation fix.

## Fixed

- **Split detail pane** — result tabs now show a left list + right detail pane instead of inline dropdown expansion
- **Selected item always visible** — selection-follow viewport algorithm enforces `viewportOffset <= selectedIndex < viewportOffset + visibleHeight` at every movement
- **SDRI large-list navigation** — 69+ findings fully navigable with selection always visible
- **Trust large-list navigation** — 100+ trust chains fully navigable with selection always visible
- **No inline dropdown expansion** — details no longer expand inline under the list, preventing list height mutation
- **Independent detail pane scroll** — detail pane scrolls independently when focused
- **Resize-safe selection** — terminal resize preserves selected item visibility
- **Ghost/stale row cleanup** — panes are padded/truncated to fixed height, no stale content
- **Viewport offset fix** — item-space ViewportOffset correctly mapped to line-space rendering (header lines no longer skew viewport position)

## Retained

- Semantic contradiction engine
- File picker fixes
- Sidebar viewport fixes
- Duplicate breadcrumb fix
- Local AI
- Case-first workflow

## Verification

- Build passed (`go fmt`, `go vet`, `go build`, `go test -count=1`, `go test -race`)
- All 19 packages pass
- 11 acceptance tests for the split-pane viewport fix all pass
- Native binary verified (darwin-arm64): `ASF0 v5.1.1`
- Result split-pane layout verified in TUI
- SDRI + Trust large-list navigation verified

## Upgrade Notes

Users should upgrade to v5.1.1 for the improved result-tab UX with split detail pane and always-visible selection.

```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash
```

Or download the binary directly from the GitHub release page.

## Known Limitations

- Controls tab list items lack `[Gap]`/`[Control]`/`[Coverage]` categorization prefixes (data model limitation)
- Verification tab group labels may cause minor header line estimation imprecision (acceptable for typical small verification lists)
