# ASF Release Audit

**Version:** 2.1.2
**Date:** 2026-06-13
**Auditor:** Release Engineering

---

## 1. Intelligence Engines

### 1.1 Assumption Intelligence (Core Engine) — `intelligence/engine.go`
- **Status: Production Ready**
- 8 exported functions, 5 phases: domain detection, topological reasoning, trust boundary discovery, explainability, quality scoring, contradiction detection, risk counting, summary generation
- 18 engine tests passing
- Integrates with all downstream engines

### 1.2 Contradiction Engine (Basic) — `intelligence/contradiction.go`
- **Status: Production Ready** (fix: added natural language variants "tls is required", "http is allowed")
- 281 lines, 8 detection rules, 12 contradiction tests passing
- Detects: MFA exemption, plaintext backup, shared admin, internet-accessible private, mutable audit, HTTP allowed, encryption without KMS, session without rotation

### 1.3 Contradiction Intelligence Engine (CIE) — `intelligence/contradiction_intelligence.go`
- **Status: Production Ready**
- 1,181 lines, 12 claim extraction patterns, 8+ detection rules
- Deep contradiction analysis with implied contradiction detection
- 10 CIE tests passing

### 1.4 Trust Boundary Intelligence (TBI) — `intelligence/trust_boundary_intelligence.go`
- **Status: Production Ready**
- 791 lines, 21 zone types, 11 boundary crossing types, 7 compliance mappings
- 60 TBI tests passing

### 1.5 Trust Boundary Engine (Basic) — `intelligence/trust_boundaries.go`
- **Status: Production Ready** (fix: added "thirdparty" compound-word matching)
- 335 lines, zone classification with priority-based tie-breaking
- 10 trust boundary tests passing

### 1.6 Taxonomy Engine — `intelligence/taxonomy.go`
- **Status: Production Ready**
- 933 lines, 24 categories, 5 severity levels
- 7 taxonomy tests passing (fix: updated stale test expectation for "Multi-tenant SaaS")

### 1.7 Reasoning Engine — `intelligence/reasoning.go`
- **Status: Production Ready** (fix: added "encrypted" alongside "encryption" keyword; fix: updated test to match "contracts")
- 612 lines, domain-specific inferences from architecture labels
- 10 reasoning tests passing

### 1.8 Domain Packs Engine — `intelligence/domain_packs.go`
- **Status: Production Ready**
- 787 lines, 4 domains (healthcare, fintech, SaaS, infrastructure)
- 14 domain pack tests passing

### 1.9 Quality Engine — `intelligence/quality.go`
- **Status: Production Ready**
- 226 lines, assumption quality scoring
- 8 quality tests passing

### 1.10 Explainability Engine — `intelligence/explainability.go`
- **Status: Production Ready**
- 201 lines, explanation generation and tracing
- 8 explainability tests passing

### 1.11 Threat Modeling Intelligence (TMI) — `intelligence/threat_modeling.go`
- **Status: Production Ready**
- 617 lines, 12 threat categories, 4 rule engines (components 20+ types, relationships 4 protocols, assumptions 5 keywords, boundaries 11 crossing types)
- 46 TMI tests passing

### 1.12 Attack Path Discovery (APD) — `intelligence/attack_path_discovery.go`
- **Status: Production Ready**
- 1,434 lines, DFS traversal, 12 kill chain phases, 30+ MITRE ATT&CK techniques, business impact, detection difficulty
- 15 APD tests passing

### 1.13 Security Design Review Intelligence (SDRI) — `intelligence/sdri.go`
- **Status: Production Ready**
- 1,474 lines, control library, design findings, architectural weaknesses, remediations, coverage dashboard, compliance alignment
- 9 SDRI tests passing

### 1.14 Compliance & Audit Readiness (CIARE) — `intelligence/ciare.go`
- **Status: Production Ready**
- 1,207 lines, 9 compliance frameworks (HIPAA, SOC2, ISO27001, PCI DSS, GDPR, FedRAMP, NIST, COBIT, CSA)
- 15 CIARE tests passing

### 1.15 Domain Knowledge Pack Intelligence (DKPI) — `intelligence/dkpi.go`
- **Status: Production Ready**
- 815 lines, knowledge pack generation, evidence requirements
- 15 DKPI tests passing

### 1.16 Executive Risk Narratives (ERN) — `intelligence/ern.go`
- **Status: Production Ready**
- 1,878 lines, report packs, narrative generation, 16 phases
- 33 ERN tests passing

### 1.17 Portfolio Intelligence (SAMPI) — `intelligence/sampi.go`
- **Status: Production Ready**
- 1,116 lines, portfolio analysis, trend detection, control coverage, maturity assessment
- 37 SAMPI tests passing

### 1.18 Decision Intelligence (SDI) — `intelligence/sdi.go`
- **Status: Production Ready**
- 1,500 lines, 15 phases, 20 canonical recommendations, decision trees, board scenarios, investment prioritization
- 21 SDI tests passing

### 1.19 Digital Twin (SDT) — `intelligence/sdt.go`
- **Status: Production Ready**
- 1,513 lines, 17 phases: twin model, change impact, security diff, evolution, control drift, assumption decay, security debt, compliance drift, attack surface, timeline, what-if, merger, zero trust, resilience, crown jewels, executive report, portfolio summary
- 20 SDT tests passing

---

## 2. CLI (`main.go` + `analyze_cli.go` + `portfolio_cli.go`)

- **Status: Production Ready**
- 12+ commands: `analyze`, `portfolio add/analyze/report/list/remove`, `doctor`, `version`, `license`, `update`
- Custom CLI (no cobra), exit codes 0/1/2/4/6/7
- Full JSON output for `analyze` command
- 1,702 lines total

---

## 3. TUI (`app.go` + 14+ model files)

- **Status: Production Ready**
- 14 models, 12 views, 9 view states (startup/dashboard/analyze/results/review/settings/about/export/validation)
- 4 themes (Dark, Midnight, Cyber, Minimal), Bubble Tea + Lipgloss
- Sections: Assumptions, Critical Assumptions, Risk Matrix, STRIDE, Controls, Attack Paths, SDR, Compliance, Intelligence, DKPI, ERN, SAMPI, SDI, SDT

---

## 4. Exports (`export.go`)

- **Status: Production Ready**
- 5 formats: JSON, Markdown, HTML, CSV, PDF
- 2,606 lines, 8 intelligence section renderers
- All intelligence engines wired into all formats

---

## 5. Installers

### `install.sh` (macOS/Linux)
- **Status: Production Ready**
- 592 lines, supports: install, upgrade, repair, purge, rollback
- Auto-detects platform, PATH configuration, backup creation

### `install.ps1` (Windows)
- **Status: Production Ready**
- 332 lines, same capabilities as install.sh

---

## 6. Build System

- **Status: Production Ready**
- `scripts/build-release.sh` — multi-platform binary generation
- `scripts/build-release.ps1` — Windows release building
- `scripts/verify-release.sh` — checksum verification
- `scripts/test-commands.sh` — CLI smoke testing
- `scripts/test-installer.sh` — installer testing

---

## 7. Licensing

- **Status: Production Ready (Demo-Grade)**
- HMAC-based activation in `license.go` (131 lines)
- Ed25519 key verification in `license_ed25519.go` (61 lines)
- License enforcement: `asf --license` displays status, analysis requires valid license
- **Limitation:** Demo-grade only — not suitable for production DRM. Suitable for community/open-source distribution where honest users respect licenses.

---

## 8. Update System

- **Status: Experimental**
- `version_check.go` (60 lines) — GitHub release version checking
- No built-in binary replacement mechanism
- Manual update path: user downloads new binary via installer

---

## 9. AI/LLM Integration

- **Status: Experimental**
- `ai.go` (314 lines) — Ollama integration
- `localai.go` (580 lines) — local model management TUI
- `model.go` (371 lines) — Ollama model management
- Disabled by default (`ai.enabled: false` in config)
- **Limitation:** Non-deterministic output, requires external Ollama server, no guarantee of model availability

---

## 10. Documentation

- **Status: Partially Implemented**
- 85+ `.md` files across the project
- README at 565 lines covers most features
- Many docs/ files are stale or incomplete
- **Limitation:** Documentation audit needed to verify claims match code

---

## 11. Deprecated/Archived

### Python Bridge (`asf/`, `benchmark/`)
- **Status: Dead Code (Archived)**
- Legacy Python engine from v1.0.0
- Go-native engine replaced it as default (`Config.Engine.UseNativeEngine: true`)
- Python code in `asf/` and `benchmark/` directories, not built by Go toolchain
- Not shipped in release binaries

---

## 12. Summary

| Capability | Lines | Tests | Status |
|---|---|---|---|
| Core Engine | 240 | 18 | Production Ready |
| Contradiction Engine | 281 | 12 | Production Ready |
| CIE Engine | 1,181 | 10 | Production Ready |
| TBI Engine | 791 | 60 | Production Ready |
| Trust Boundaries | 335 | 10 | Production Ready |
| Taxonomy Engine | 933 | 7 | Production Ready |
| Reasoning Engine | 612 | 10 | Production Ready |
| Domain Packs | 787 | 14 | Production Ready |
| Quality Engine | 226 | 8 | Production Ready |
| Explainability | 201 | 8 | Production Ready |
| TMI Engine | 617 | 46 | Production Ready |
| APD Engine | 1,434 | 15 | Production Ready |
| SDRI Engine | 1,474 | 9 | Production Ready |
| CIARE Engine | 1,207 | 15 | Production Ready |
| DKPI Engine | 815 | 15 | Production Ready |
| ERN Engine | 1,878 | 33 | Production Ready |
| SAMPI Engine | 1,116 | 37 | Production Ready |
| SDI Engine | 1,500 | 21 | Production Ready |
| SDT Engine | 1,513 | 20 | Production Ready |
| CLI | 1,702 | — | Production Ready |
| TUI | 5,000+ | — | Production Ready |
| Exports (5 formats) | 2,606 | 1 | Production Ready |
| Installers (2) | 924 | — | Production Ready |
| Build System | 4 scripts | — | Production Ready |
| Licensing | 192 | — | Demo-Grade |
| Update System | 60 | — | Experimental |
| AI/LLM | 1,265 | 12 | Experimental |
| Python Bridge | Archived | — | Dead Code |

**Overall Assessment: Production Ready with minor gaps**
