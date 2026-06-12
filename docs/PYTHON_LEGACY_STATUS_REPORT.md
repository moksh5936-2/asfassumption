# Python API Legacy Status Report

## Decision Rationale

ASF v2.x has been a Go-native single binary since v2.0.0 (June 2026).
The Python REST API at `asf/` is legacy code — it was the original v1 engine
before the Go migration. Retaining it as reference only allows:
- Traceability for migration decisions
- Study of original architecture
- No accidental use alongside the production Go binary

## What Was Documented

| Item | Action |
|------|--------|
| `docs/LEGACY_PYTHON_REFERENCE.md` | Created — documents the archived status, config separation, and reference-only intent |
| `README.md` (project structure) | Updated `asf/` entry from "Python ASF engine (v1)" to "Python ASF engine (v1, archived)" |
| `GUIDE.md` | Added legacy banner at top of document; original content preserved |
| Installer scripts | No changes needed — all installers (`install.sh`, `release/install.sh`, `asf-tui/install.sh`, `install.ps1`) already install only the Go binary with no Python references |

## Impact on README

The README already correctly:
- Listed `asf/` as "Python ASF engine (v1)" in the project structure
- Stated in FAQ that ASF v2.1.1+ is a self-contained Go binary with no Python dependency
- Described the Go TUI as the sole application entry point

The only change was strengthening the project structure note from "(v1)" to "(v1, archived)".

## Config Separation Confirmed

- **Go TUI config** — `~/.asf/config.yaml` (production)
- **Python API config** — `asf.config.yaml` in CWD (legacy)

These are documented as separate, non-interoperable config systems in
`LEGACY_PYTHON_REFERENCE.md`.

## Conclusion

No Python API claims remain in current product docs without legacy qualification.
All installer scripts target only the Go binary. The Python source at `asf/`
is preserved in place and marked as archived reference material.
