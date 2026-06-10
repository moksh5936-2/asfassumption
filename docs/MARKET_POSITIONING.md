# ASF Market Positioning

> Version: 1.0.0 | June 2026 | Classification: Confidential — Internal Strategy

## 1. Executive Summary

ASF (Architecture Security Framework) is a deterministic, offline-first security assumption discovery tool that finds hidden security assumptions in system architecture diagrams. It occupies a unique niche between threat modeling tools and static analysis tools — focusing on what teams *implicitly assume* about their systems.

**Tagline:** "Discover what your architecture assumes about security."

---

## 2. Competitive Landscape

### Direct Competitors

| Tool | Focus | ASF Advantage |
|------|-------|--------------|
| **Microsoft Threat Modeling Tool** | STRIDE-based threat modeling | ASF is automated, CLI-native, assumption-focused |
| **OWASP Threat Dragon** | Open-source threat modeling | ASF doesn't require manual diagram entry |
| **IriusRisk** | Commercial threat modeling | ASF is free, offline, deterministic |
| **SD Elements** | Security requirements | ASF focuses on assumptions, not requirements |
| **CAI (Cortex XSOAR)** | Automated threat modeling | ASF is simpler, no infra required |

### Adjacent Tools

| Tool | Focus | Relationship |
|------|-------|-------------|
| **Snyk** | Known vulnerability scanning | Complementary — Snyk finds CVEs, ASF finds assumptions |
| **Semgrep** | SAST (static analysis) | Different layer — ASF operates at architecture level |
| **Checkov** | IaC security | Different layer — ASF operates at design level |
| **Lucidchart** | Diagramming | Complementary — ASF analyzes Lucidchart exports |

### Market Gap

```
                    Vulnerability-focused
                    │
    SAST            │   DAST
    (Semgrep)       │   (Burp Suite)
                    │
       ────────────┼───────────────
                    │              Design/Architecture
                    │
    IaC Security    │   ❌ GAP ❌
    (Checkov)       │   ASF fills this
                    │
                    │   Threat Modeling
                    │   (MS TMT, IriusRisk)
```

**The gap:** No tool automatically discovers implicit security assumptions from architecture diagrams. Existing tools find known vulnerabilities (CVEs), misconfigurations (IaC), or require manual threat modeling. ASF automates the discovery phase of threat modeling.

---

## 3. Target Audience

### Primary: Security Architects

| Profile | Value Proposition |
|---------|-------------------|
| Threat modeling at scale | ASF automates assumption discovery across dozens of architectures |
| Reviewing third-party designs | ASF provides independent, reproducible analysis |
| Teaching threat modeling | ASF shows which assumptions are discoverable automatically |

### Secondary: DevOps/Platform Engineers

| Profile | Value Proposition |
|---------|-------------------|
| "Shift left" security | ASF finds design-level issues before code is written |
| Infrastructure review | ASF analyzes IaC-generated architecture diagrams |
| CI/CD integration | ASF CLI supports batch analysis |

### Tertiary: Security Researchers

| Profile | Value Proposition |
|---------|-------------------|
| Studying architecture-level security | ASF provides deterministic, reproducible analysis |
| Comparing architecture patterns | ASF runs the same engine on any architecture |

---

## 4. Differentiation

### Key Differentiators

| Differentiator | ASF | Competitors |
|---------------|-----|-------------|
| **Automated assumption discovery** | ✅ Core feature | ❌ Manual entry required |
| **Deterministic (no AI)** | ✅ Fully reproducible | ❌ Often uses ML/heuristics |
| **Offline-first** | ✅ No internet required | ❌ SaaS-dependent |
| **Single binary** | ✅ 11.9MB, no runtime | ❌ Often require servers, DBs |
| **STRIDE mapping** | ✅ 17 category rules + 30 keyword rules | ⚠️ Usually manual |
| **5×5 risk matrix** | ✅ Deterministic likelihood × impact | ⚠️ Often subjective |
| **Confidence scoring** | ✅ Evidence-based formula | ❌ Not commonly available |
| **Evidence traceability** | ✅ Source component + line tracking | ❌ Not commonly available |
| **5 export formats** | ✅ JSON/CSV/MD/PDF/HTML | ⚠️ Usually 1-2 formats |
| **Cost** | ✅ Free (open core) | ❌ Often expensive per-seat |

### What ASF is NOT

- ASF is **not** a vulnerability scanner (no CVE database)
- ASF is **not** a SAST tool (no code analysis)
- ASF is **not** a threat modeling tool (no manual diagram creation)
- ASF is **not** a compliance tool (no regulatory framework mapping)
- ASF is **not** an AI tool (AI is optional, not core)

---

## 5. Pricing Strategy

### Open Core Model

| Tier | Features | Price |
|------|----------|-------|
| **Community** | Core ASF engine, CLI, TUI, all exports | Free |
| **Enterprise** | License management, priority support, SSO, audit features | Per-seat licensing |

### Enterprise Licensing

- **Individual:** $99/user/year
- **Team (10+):** $79/user/year
- **Enterprise (50+):** Custom pricing

### Current Status

ASF v1.0.0 is fully free (community tier). The enterprise license system exists (`license.go`) but does not gate any features. Monetization strategy is not yet implemented.

---

## 6. Go-to-Market Strategy

### Phase 1: Community Building (Current)

- Open source under permissive license
- Build documentation and validation study
- Engage security researcher community

### Phase 2: Validation (Next)

- Complete expert validation study (10 architects × 20 architectures)
- Publish precision/recall/FPR metrics
- Present at security conferences (BSides, OWASP)

### Phase 3: Enterprise (Future)

- Implement feature gating via license system
- Add SSO and audit logging
- Publish case studies and whitepapers

---

## 7. SWOT Analysis

### Strengths
- Unique value proposition (assumption discovery)
- Deterministic, auditable, offline-first
- Excellent documentation
- Strong technical foundation (Bubble Tea, clean Go)
- Explainability is a differentiator

### Weaknesses
- **No empirical validation data** (zero metrics published)
- Python dependency for core extraction
- Single developer (bus factor = 1)
- Single binary available (darwin/arm64 only)
- No CI/CD, no test coverage for 80% of code
- No Windows testing

### Opportunities
- Growing market for "shift left" security
- Increasing cloud architecture complexity
- AI hype creates interest in automated analysis
- Compliance frameworks (SOC 2, ISO 27001) require design review

### Threats
- Competitors adding assumption discovery features
- Python dependency becomes unmaintained
- Negative study results (ASF performs poorly)
- Enterprise customers require features ASF lacks (SSO, audit, team)

---

## 8. Messaging Framework

### Core Message

> **ASF automatically discovers hidden security assumptions in your architecture diagrams — without AI, without the cloud, without manual effort.**

### Use Case Messages

| Use Case | Message |
|----------|---------|
| Architecture review | "Before you build, know what you're assuming." |
| Threat modeling | "Turn implicit assumptions into explicit security requirements." |
| Third-party audit | "Independently verify the security assumptions in any architecture." |
| Education | "See what assumptions your architecture reveals." |

### Tagline Options

1. "Discover what your architecture assumes."
2. "The assumption discovery engine."
3. "Security assumptions, automated."
4. "What does your architecture assume?"

---

## 9. Channels

| Channel | Priority | Strategy |
|---------|----------|----------|
| GitHub | 🔴 Primary | Open source repository, issues, discussions |
| Security conferences | 🟡 Secondary | Talk proposals, workshops, demos |
| Blog/technical writing | 🟡 Secondary | Architecture analysis case studies |
| Social media (Twitter/LinkedIn) | 🟢 Tertiary | Security community engagement |
| Podcasts | 🟢 Tertiary | Security tooling interviews |

---

## 10. Success Metrics

| Metric | Current | 6-Month Target | 12-Month Target |
|--------|---------|----------------|-----------------|
| GitHub stars | N/A | 500 | 2,000 |
| Weekly active users | N/A | 50 | 500 |
| Expert validation | ❌ None | ✅ Published | ✅ Validated |
| Enterprise customers | 0 | 5 | 25 |
| Precision/recall | Unknown | ≥70%/≥60% | ≥80%/≥70% |
| Download count | N/A | 1,000 | 10,000 |

---

## 11. Current Status

| Item | Status |
|------|--------|
| Core product | ✅ v1.0.0 |
| Documentation | ✅ Comprehensive (12 files) |
| Validation data | ❌ None |
| Market research | ✅ Complete |
| Pricing strategy | ⚠️ Draft |
| Enterprise features | ⚠️ License system exists, no gating |
| CI/CD | ❌ Not implemented |
| Marketing materials | ⚠️ README only |

---

*This document represents the market positioning strategy for ASF v1.0.0. All targets and strategies are preliminary and subject to validation study results.*
