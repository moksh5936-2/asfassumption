# ASF — Executive Summary

## Project Vision

Make hidden security assumptions in software architecture visible, auditable, and actionable — without requiring AI or cloud connectivity.

## What ASF Is

ASF (Architecture Security Framework) is a **deterministic, offline-first security review engine** that analyzes system architecture diagrams and documents to discover implicit security assumptions. It processes Draw.io diagrams, Mermaid markdown, YAML/JSON definitions, SVG diagrams, images (via OCR), and plain text documents through a structured pipeline: parsing → assumption extraction → STRIDE threat mapping → risk assessment → confidence scoring → human review.

The output is a fully explainable, reproducible security analysis with evidence traceability for every finding.

## What ASF Is Not

- **Not a vulnerability scanner** — it does not find CVEs, misconfigurations, or known exploits
- **Not an AI application** — AI is entirely optional and fully local (Ollama); core analysis is deterministic
- **Not a compliance checker** — though it can inform compliance assessment
- **Not a cloud service** — everything runs locally, no data leaves the machine
- **Not a real-time monitor** — it is a design-time / review-time analysis tool

## Problem Statement

Most security tools scan for known vulnerabilities. But catastrophic breaches often exploit **implicit assumptions** that architects and engineers make about their systems:

- "The database is in a private subnet" → Reality: misconfigured route table
- "MFA is enforced everywhere" → Reality: legacy app bypasses IdP
- "Vendor handles security" → Reality: vendor SOC 2 covers different scope
- "Audit logs are immutable" → Reality: DBAs can modify logs
- "Backups protect us" → Reality: backups never tested

ASF systematizes the discovery of these hidden dependencies.

## Target Users

| User | Use Case |
|------|----------|
| Security Architects | Pre-deployment review, threat modeling |
| Software Architects | Design validation, assumption audit |
| Security Engineers | Gap analysis, control recommendation |
| Penetration Testers | Target selection, attack surface mapping |
| Compliance Officers | Evidence collection, audit preparation |
| Researchers | Security assumption ontology, validation studies |

## Value Proposition

1. **Deterministic & Auditable** — every risk score, STRIDE mapping, and confidence level is reproducible. No black-box AI.
2. **Fully Offline** — no internet required. All processing is local. Air-gapped environments supported.
3. **Multi-Format Input** — Draw.io, Mermaid, YAML, JSON, SVG, images, PDF, DOCX, TXT.
4. **Explainable Output** — every assumption includes evidence traceability, STRIDE justification, risk decomposition, and confidence scoring.
5. **5 Export Formats** — JSON, Markdown, CSV, PDF, HTML — suitable for reports, dashboards, and further analysis.
6. **Human-in-the-Loop Review** — architect review mode with Accept/Reject/Modified status tracking.
7. **Optional Local AI** — Ollama integration for enhanced analysis, no cloud dependency.

## Architecture Review Workflow

```
1. Open ASF TUI
2. Select "Analyze Architecture"
3. Provide architecture (drawio/mmd/yaml/json/svg/png/jpg/txt/pdf/docx)
4. Optional: provide evidence file (CSV/JSON)
5. ASF processes through 7-stage pipeline
6. Review results in TUI
7. Export to any of 5 formats
8. Optional: human review mode with status tracking
```

## Competitive Advantages

| Feature | ASF | Traditional Threat Modeling Tools |
|---------|-----|-----------------------------------|
| Offline | ✓ Full | Often requires cloud |
| Deterministic | ✓ Pure logic | Usually manual or AI |
| Explainable | ✓ 6 engines | Often black-box |
| Multi-format parse | ✓ 10+ formats | Usually manual input |
| Free | ✓ Open source | Usually expensive |
| Local AI option | ✓ Ollama | Usually cloud AI |

## Current Capabilities

- ✓ Parse 10+ architecture formats (Draw.io, Mermaid, YAML, JSON, SVG, images, PDF, DOCX, TXT)
- ✓ Extract security assumptions via Python ASF engine bridge
- ✓ Map assumptions to STRIDE categories (17 category rules + 33 keyword patterns)
- ✓ Risk assessment via 5×5 deterministic matrix (3 likelihood factors × 3 impact factors)
- ✓ Confidence scoring from 4 metrics (evidence, rules, components, relationships)
- ✓ Evidence traceability to source components and relationships
- ✓ 5 export formats with full explainability data
- ✓ 4 themes (Dark, Midnight, Cyber, Minimal)
- ✓ Architect review mode with Accept/Reject/Modified tracking
- ✓ Validation data collection for precision/recall studies
- ✓ Ollama model management (download, list, delete, activate)
- ✓ Optional AI enhancement layer (local only)
- ✓ License system (HMAC-signed enterprise keys)
- ✓ Auto-config migration from legacy paths

## Current Limitations

- ✗ Python ASF CLI dependency (`python -m asf.cli.main`)
- ✗ Tesseract required for image OCR
- ✗ Ollama required for AI features
- ✗ No cloud-based AI option (by design)
- ✗ STRIDE accuracy not validated against expert review
- ✗ Precision and recall not measured (no validation study yet)
- ✗ False positive rate unknown
- ✗ No CI/CD pipeline
- ✗ No code signing or notarization

## Future Validation Roadmap

| Priority | Study | Status | Impact |
|----------|-------|--------|--------|
| P0 | Expert validation study (10 architects × 20 architectures) | Not started | Highest — proves methodology value |
| P1 | Precision/recall measurement | Design ready | Quantifies accuracy |
| P2 | False positive rate analysis | Design ready | Quantifies noise |
| P3 | STRIDE accuracy vs expert mapping | Design ready | Validates rule engine |
| P4 | Inter-rater reliability (Cohen's kappa) | Design ready | Validates reproducibility |
| P5 | Commercial pilot with enterprise customer | Not started | Validates market fit |

## Known Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Python dependency complicates distribution | High | Medium | Bundle Python or rewrite extraction in Go |
| Validation study may show poor precision | Medium | High | Improve rule engine; adjust confidence thresholds |
| Limited adoption without empirical proof | High | High | Prioritize validation study |
| Ollama dependency limits AI adoption | Medium | Low | AI is optional; core works without it |
| Tesseract OCR quality varies | Low | Medium | Suggest Draw.io/Mermaid for best results |

## Commercial Readiness Assessment

| Criterion | Score | Notes |
|-----------|-------|-------|
| Code quality | 8/10 | Clean Go, good separation of concerns |
| Documentation | 5/10 | Being written now |
| Testing | 6/10 | 20 unit tests, no integration tests |
| CI/CD | 0/10 | None |
| Binary distribution | 5/10 | Single-platform builds, no CI |
| License enforcement | 7/10 | HMAC-based, works offline |
| User experience | 8/10 | Keyboard TUI, no memorized commands |
| **Overall** | **6/10** | Requires validation evidence for commercial |

## Technical Readiness Assessment

| Criterion | Score | Notes |
|-----------|-------|-------|
| Architecture | 9/10 | Clean separation, pipelined design |
| Compilation | 10/10 | `go vet` clean, 11.85MB binary |
| Error handling | 8/10 | Explicit error returns, TUI error display |
| Performance | 8/10 | ~2158 assumptions in benchmark run |
| Security | 7/10 | HMAC license, no hardcoded secrets (one test HMAC key) |
| Dependencies | 7/10 | 6 direct dependencies, all well-maintained |
| Portability | 6/10 | Linux/macOS only; Windows TUI needs testing |
| **Overall** | **8/10** | Solid engineering, minor platform gaps |

## Research Readiness Assessment

| Criterion | Score | Notes |
|-----------|-------|-------|
| Reproducibility | 10/10 | Deterministic pipeline, same input → same output |
| Explainability | 9/10 | 6 engines provide full traceability |
| Validation data | 8/10 | CollectValidationData() ready for studies |
| Empirical evidence | 2/10 | No expert validation study yet |
| Methodology | 7/10 | Structured but unvalidated against human experts |
| **Overall** | **7/10** | Engine is research-ready; needs validation study |

## Overall Maturity Assessment

**Current: 7/10**

ASF is a well-engineered tool with strong technical foundations. The core analysis pipeline is deterministic, explainable, and produces actionable output. The gap is not in the technology — it is in the **empirical validation**. Without a study proving ASF finds assumptions that human architects miss, the product lacks defensible evidence of its value proposition.

**Next milestone**: Expert validation study (20 architectures, 10 security architects, precision/recall/novelty measurement).
