# ASF — Assumption Security Framework

**ASF is a research framework for discovering hidden security assumptions in software architecture.**

Most security tools scan for known vulnerabilities (CVEs, misconfigurations, SAST findings). ASF targets what they miss: the *implicit assumptions* that architects and engineers make about their systems — assumptions that, when violated, become the root cause of catastrophic breaches.

This repository contains two things:
1. **ASF v1** (`asf/`) — a working CLI/API tool that extracts assumptions from policy documents and verifies them against evidence
2. **Phase 6 Research Validation** (`benchmark/`) — experimental evidence across 20 reference architectures, multi-LLM evaluation, and a human-survey-ready validation package

---

## Key Research Findings (Phase 6)

| Metric | Value |
|--------|-------|
| Mean recall across 20 architectures | 69.9% |
| Mean precision | 43.3% |
| Mean novel assumptions per architecture | 40.3 |
| ASF predictions validated by ≥1 AI (Tier A+B) | 69.1% |
| ASF-unique predictions (Tier C) | 31.6% |
| AI hallucination rate | 2.6% |
| Number of architectures evaluated | 20 |
| Total ASF predictions scored | 1,409 |

ASF consistently discovers assumptions that *no human or AI independently lists* — about monitoring infrastructure, identity lifecycle governance, third-party dependency risk, compensating controls, and vendor exit strategy. These are not false positives; they are blind spots.

**The most important output:** 446 Tier C (ASF-unique) assumptions across 20 architectures — assumptions no other reviewer produced.

---

## Repository Structure

```
asf/                          # ASF v1 tool — CLI, API, analysis engine
  cli/                        # CLI entry point
  api/                        # FastAPI web server + UI
  extraction/                 # Claim extraction from documents
  verification/               # Evidence-based verification engine
  assumption/                 # Assumption models and patterns
  evidence/                   # Evidence file parsing and matching
  gaps/                       # Gap analysis between assumptions and evidence
  confidence/                 # Confidence scoring (verification, freshness, coverage)
  graph/                      # Relationship graph
  ingestion/                  # Document parsing (PDF, DOCX, TXT)
  llm/                        # Optional LLM integration
  db/                         # SQLite persistence
  models/                     # Pydantic data models
  config.py / settings.py     # Configuration

benchmark/                    # Phase 6 Research Validation
  experiments/                # 20 architecture simulations, multi-LLM campaign, diagnostics
    architecture_001-020_simulation.md  # Per-architecture human vs ASF comparison
    multi_llm_campaign_*      # Multi-LLM evaluation results
    AGGREGATE_RESULTS.md      # Cross-architecture aggregate metrics
    MULTI_LLM_CAMPAIGN_GRAND_AGGREGATE.md  # Unified multi-LLM results
    asf_diagnostic_tests.md   # 5 diagnostic tests against ASF v1
    WOULD_YOU_PAY_FOR_THIS_STUDY.md  # Human survey package (ready to run)
    INDEPENDENT_DERIVATION_TEST.md    # AI blind reproduction test (ready to run)
  report/                     # Analyst guide, evaluation framework, ontology
    ASF_v1_Analyst_Guide.md   # Comprehensive analyst training document
    LLM_as_Judge_Framework.md # AUS rubric, consensus matrix, multi-judge voting
    ontology/                 # ASF assumption ontology, trust graph model
  assumption_knowledge_base/  # Gold standard data, architecture patterns, experiment kit
    architecture_patterns.md  # 20 reference architectures (source for all experiments)
    assumption_gold_standard.csv  # 1,229 pre-generated ASF predictions
    experiment_kit.md         # 12-page printable participant packet

sample_data/                  # Sample policies and evidence for testing
tests/                        # Unit and integration tests
scripts/                      # Utility scripts
```

---

## ASF v1 Tool: Quick Start

```bash
# Install
pip install -e .

# Analyze a policy document against evidence
asf analyze policy.txt -e evidence.csv

# JSON output
asf analyze policy.txt -e evidence.csv --json

# Persist results to SQLite database
asf analyze policy.txt -e evidence.csv --persist

# Start the web UI
uvicorn asf.api.app:app --reload
```

The ASF v1 tool extracts security claims from policy documents, classifies them by assumption type (ACCESS, NETWORK, IDENTITY, CONFIGURATION, PROCESS, DEPENDENCY, GOVERNANCE, DOCUMENTATION), and verifies each claim against structured evidence files (CSV, JSON).

**Note:** ASF v1 uses regex-based extraction and has a ~11.7% recall ceiling. It is a proof-of-concept for the verification engine, not the discovery methodology. The discovery methodology (which achieves ~70% recall) is a structured *human-led* process documented in the benchmark experiments.

---

## The Discovery Problem ASF Solves

Traditional security tools answer: *"What vulnerabilities exist in my system?"*

ASF answers: *"What did my team assume about the system that an attacker would exploit?"*

Real-world breaches rarely exploit unknown CVEs. They exploit:
- **Assumption:** "The database is in a private subnet" → Reality: misconfigured route table
- **Assumption:** "MFA is enforced everywhere" → Reality: legacy app bypasses IdP MFA
- **Assumption:** "Vendor handles security" → Reality: vendor SOC 2 covers different scope
- **Assumption:** "Audit logs are immutable" → Reality: DBAs can modify logs
- **Assumption:** "Backups protect us" → Reality: backups never tested, restore fails

ASF systematizes the discovery of these hidden dependencies.

---

## Running the Research Validation

The benchmark experiments are markdown documents — no code execution required.

### To reproduce the multi-LLM campaign:
1. Read `architecture_patterns.md` for the 20 reference architectures
2. For each architecture, run the ASF structured discovery process (documented in `ASF_v1_Analyst_Guide.md`)
3. Score results using the AUS rubric in `LLM_as_Judge_Framework.md`
4. Compare against the pre-generated gold standard in `assumption_gold_standard.csv`

### To run the independent derivation test:
1. Use the prompts in `INDEPENDENT_DERIVATION_TEST.md`
2. Feed each architecture-only prompt to GPT-4o, Claude 4 Sonnet, and Gemini 2.5 Pro in *fresh* chats
3. Map output against the comparison matrices
4. Measure ASF-unique rate

### To run the human study:
1. Use the survey package in `WOULD_YOU_PAY_FOR_THIS_STUDY.md`
2. Recruit 10-20 security practitioners
3. For each of 20 assumptions (10 Tier A + 10 Tier C), ask the 5 utility questions
4. Score willingness to pay, investigation intent, and incident history

---

## Evidence File Format

ASF v1 verifies assumptions against structured evidence. Supported formats:

### Evidence Schema Auto-Mapping

| Your Column | Maps To | Used For |
|-------------|---------|----------|
| `user`, `username`, `employee`, `email` | user | Identity & access checks |
| `group`, `department`, `team`, `unit` | group | Access control verification |
| `resource`, `application`, `system`, `service` | resource | Resource-level access checks |
| `permission`, `access`, `role`, `right` | permission | Permission level checks |
| `mfa`, `mfa_enabled`, `multi_factor`, `2fa`, `totp` | mfa | MFA compliance checks |
| `public`, `exposed`, `internet_facing`, `is_public` | public | Network exposure checks |
| `asset`, `host`, `server`, `system_name` | asset | Asset inventory checks |

### Example evidence CSV

```csv
user,group,resource,permission,department
alice.jones,Finance,payroll-system,read,Finance
bob.smith,Finance,payroll-system,write,Finance
```

### Verification Results

| Result | Meaning |
|--------|---------|
| `VERIFIED` | Evidence supports the assumption |
| `CONTRADICTED` | Evidence contradicts the assumption |
| `PARTIALLY_VERIFIED` | Some evidence supports, some contradicts |
| `UNKNOWN` | No matching evidence available |

---

## Core Research Papers (contained in this repo)

| Document | Location |
|----------|----------|
| **20 Architecture Simulation Results** | `benchmark/experiments/AGGREGATE_RESULTS.md` |
| **Multi-LLM Campaign Grand Aggregate** | `benchmark/experiments/MULTI_LLM_CAMPAIGN_GRAND_AGGREGATE.md` |
| **ASF v1 Diagnostic Tests** | `benchmark/experiments/asf_diagnostic_tests.md` |
| **ASF v1 Analyst Guide** | `benchmark/report/ASF_v1_Analyst_Guide.md` |
| **LLM-as-a-Judge Framework** | `benchmark/report/LLM_as_Judge_Framework.md` |
| **ASF Assumption Ontology** | `benchmark/report/ontology/ASF_Assumption_Ontology.md` |
| **Trust Graph Model** | `benchmark/report/ontology/trust_graph_model.json` |
| **20 Architecture Patterns** | `benchmark/assumption_knowledge_base/architecture_patterns.md` |
| **Assumption Gold Standard (1,229 predictions)** | `benchmark/assumption_knowledge_base/assumption_gold_standard.csv` |
| **Experiment Kit (printable)** | `benchmark/assumption_knowledge_base/experiment_kit.md` |
| **Would You Pay For This? Study** | `benchmark/experiments/WOULD_YOU_PAY_FOR_THIS_STUDY.md` |
| **Independent Derivation Test** | `benchmark/experiments/INDEPENDENT_DERIVATION_TEST.md` |

---

## Project Status

ASF is in **Phase 6 (Human Validation)**. The framework has been validated through:

- [x] Phase 1-5: ASF v1 tooling, ontology, methodology
- [x] Phase 6a: 20 architecture simulations (human architect vs ASF methodology)
- [x] Phase 6b: Multi-LLM campaign (5 AI personas × multi-judge AUS scoring)
- [ ] Phase 6c: Independent derivation test (blind AI reproduction)
- [ ] Phase 6d: "Would You Pay For This?" study (real human practitioners)
- [ ] Phase 7: ASF v2 engineering

**ASF v2 will not be built until Phase 6 proves the discovery methodology has real-world value with real security practitioners.**

---

## License

Research and educational use. See LICENSE file for details.

---

## Citation

If you reference this work in research, cite the assumption ontology and experimental methodology described in `benchmark/report/ontology/ASF_Assumption_Ontology.md` and `benchmark/experiments/AGGREGATE_RESULTS.md`.

```
@misc{asf2025,
  title={ASF: An Assumption Security Framework for Discovering Hidden Security Assumptions in Software Architecture},
  author={ASF Project},
  year={2025},
  howpublished={\url{https://github.com/moksh5936-2/asfassumption}}
}
```
