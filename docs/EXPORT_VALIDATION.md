# Export Validation Report

**Version:** 2.1.2
**Date:** 2026-06-13

---

## Formats Verified

| Format | Function | File | Status |
|---|---|---|---|
| JSON | `exportJSON()` | `export.go:63` | PASS |
| Markdown | `exportMarkdown()` | `export.go:71` | PASS |
| HTML | `exportHTML()` | `export.go:674` | PASS |
| CSV | `exportCSV()` | `export.go:1291` | PASS |
| PDF | `exportPDF()` | `export.go:1331` | PASS |

## Intelligence Sections in Markdown/HTML/PDF

| Section | Function | Status |
|---|---|---|
| Portfolio Intelligence (SAMPI) | `renderSAMPIReportMarkdown()` | PASS |
| Decision Intelligence (SDI) | `renderSDIReportMarkdown()` | PASS |
| Digital Twin (SDT) | `renderSDTReportMarkdown()` | PASS |
| Executive Narratives (ERN) | (in main export) | PASS |

## JSON Export — Section Presence

- Version, architecture, summary: PASS
- Assumptions, verifications, gaps: PASS
- Contradictions, CIE contradictions: PASS
- Trust boundaries, TBI zones/boundaries/weaknesses: PASS
- Threats, threat clusters, threat model summary: PASS
- Attack paths, threat chains, attack path summary: PASS
- SDRI controls, findings, weaknesses, remediations, coverage: PASS
- CIARE framework coverages, audit readiness, evidence, gaps: PASS
- DKPI domain, confidence, recommendations, threats: PASS
- ERN executive risks, reports, exposure: PASS
- SAMPI, SDI, SDT (via AnalysisResult serialization): PASS

## Output Validation

- JSON output validated: 716KB valid JSON for healthcare benchmark
- All sections present with non-empty data
- No corrupted or truncated output

## Backward Compatibility

- All new fields use `omitempty`
- Existing parsers will not break on field addition
- CLI output maintains backward-compatible flat format

## Conclusion

**All export formats produce valid, complete output.**
