# Changelog

## v2.2.0 (2026-06-13)

### Added
- Intelligence pipeline now 100% complete with 11 engines (V3 → CIE → TBI → TMI → APD → SDRI → CIARE → DKPI → ERN → SAMPI → SDI → SDT)
- **Digital Twin Intelligence (SDT)** — Full 17-phase engine with architectural risk, gap, control, and attack path analysis
- **Decision Intelligence (SDI)** — 15-phase engine with 20 canonical security recommendations
- **Portfolio Intelligence (SAMPI)** — Multi-architecture security portfolio analysis
- TUI Section 13 — "Digital Twin" results view
- Export: SDT, SDI, SAMPI sections in Markdown/HTML/PDF

### Fixed
- 5 pre-existing test failures (contradiction keyword matching, reasoning keyword matching, taxonomy category matching, trust boundary matching)
- All 257 tests across 11 packages now pass

### Changed
- Version bumped from 2.1.2 → 2.2.0
- Pipeline progression: V3(65%) → SDT(100%)
- README test count corrected: 257 tests across 11 packages
- Release binaries rebuilt with all engine fixes

### Removed
- Outdated Python bridge (dead code)

## v2.1.2 (2026-06-12)
- Intelligence pipeline up to SDI(98%)
- TUI v2 with section rendering for all engines
- Release hardening sprint (Phases 1-15)

## v2.1.1 (2026-06-10)
- Intelligence pipeline up to SAMPI(96%)
- DKPI/ERN/SAMPI engines integrated into TUI

## v2.0.2 (2026-06-09)
- Full Go-native binary, Python dependency removed
- CIARE, DKPI, ERN engines

## v2.0.1 (2026-06-07)
- Intelligence pipeline up to SDRI(84%)
- 12 intelligence engines integrated

## v2.0.0 (2026-06-05)
- Initial Go-native release
- V1-V3 intelligence engines
