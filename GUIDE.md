# ASF Validator v0.1 User Guide

> **LEGACY DOCUMENT** — This guide describes the Python ASF Validator (v0.1),
> which is archived and not part of the ASF v2.x production runtime.
> The production system is the Go-native single binary (`asf`).
> See [LEGACY_PYTHON_REFERENCE.md](docs/LEGACY_PYTHON_REFERENCE.md) for details.

ASF Validator tests security **assumptions** against **evidence** to detect gaps between what an organization *believes* and what *actually exists*.

---

## Quick Start

```bash
# Install
pip install -e .

# Create default config
asf init

# Analyze a policy with evidence
asf analyze policy.txt -e iam.csv -e network.csv

# Analyze a directory of files
asf analyze ./policies/ -e ./evidence/

# JSON output
asf analyze policy.txt -e iam.csv --json

# Persist to database
asf analyze policy.txt -e iam.csv --persist
```

---

## CLI Reference

### `asf analyze`

```
asf analyze [PATHS...] [-e EVIDENCE] [--json] [--graph] [--persist]
```

| Argument | Description |
|----------|-------------|
| `PATHS` | Document files or directories (PDF, DOCX, TXT). Directories are scanned recursively for supported files. |
| `-e, --evidence` | Evidence files (CSV, JSON). Repeatable. |
| `--json` | Output results as JSON. |
| `--graph` | Export the relationship graph as JSON. |
| `--persist` | Save results to SQLite database. |

### `asf init`

Creates a default `asf.config.yaml` in the current directory.

---

## Evidence File Format

Evidence files are CSV or JSON files containing **structured records** that the verification engine uses to check assumptions.

### Column Name Auto-Mapping

The system automatically maps common column names to standard fields:

| Your Column | Maps To | Used For |
|-------------|---------|----------|
| `user`, `username`, `employee`, `name`, `email` | **user** | Identity & access checks |
| `group`, `department`, `team`, `unit`, `division`, `org` | **group** | Access control verification |
| `resource`, `application`, `system`, `service`, `target` | **resource** | Resource-level access checks |
| `permission`, `access`, `role`, `right`, `privilege`, `action` | **permission** | Permission level checks |
| `mfa`, `mfa_enabled`, `multi_factor`, `2fa`, `totp` | **mfa** | MFA compliance checks |
| `enabled`, `status`, `state`, `active` | **enabled** | Configuration checks |
| `public`, `exposed`, `internet_facing`, `is_public` | **public** | Network exposure checks |
| `asset`, `host`, `server`, `system`, `service_name` | **asset** | Asset inventory checks |

### Example: Access Control Evidence (payroll_acl.csv)

```csv
user,group,resource,permission,department
alice.jones,Finance,payroll-system,read,Finance
bob.smith,Finance,payroll-system,write,Finance
dave.wilson,Engineering,payroll-system,read,Engineering
```

### Example: Network Exposure Evidence (network_exposure.csv)

```csv
asset,environment,public,internet_facing,exposure,segment
payroll-app,production,false,false,internal,finance
customer-portal,production,true,true,public,dmz
```

### Example: IAM Export (iam_export.json)

```json
{
  "users": [
    {"username": "alice.jones", "department": "Finance", "mfa": true, "groups": ["Finance", "Payroll-Admin"]},
    {"username": "dave.wilson", "department": "Engineering", "mfa": true, "groups": ["Engineering", "Payroll-Read"]}
  ]
}
```

**Note**: JSON records are expected as an array or a single object. For nested structures, flatten them to records with key-value fields.

---

## Configuration (asf.config.yaml)

```yaml
# Database path for persistence
db_path: "asf_validator.db"

# Map your column names to standard fields
evidence_schema:
  user: "employee_id"         # if your CSV has "employee_id" instead of "user"
  group: "department_name"

# Additional regex patterns for claim extraction
claim_patterns:
  - "\\b(?:assert|guarantee|promise)\\s+that\\s+.+"

# Assumption type scoring weights (higher = more likely)
assumption_weights:
  ACCESS: 1.5     # Boost ACCESS detection
  NETWORK: 0.5    # Reduce NETWORK detection

# Custom severity overrides
severity_rules:
  - assumption_type: "ACCESS"
    min_confidence: 0.8
    severity: "CRITICAL"

# LLM provider (optional — system works without it)
llm:
  provider: "openai"           # or "ollama"
  model: "gpt-4o"
  # Or use env vars: ASF_LLM_API_KEY, ASF_LLM_BASE_URL, ASF_LLM_MODEL

# Default evidence directories
evidence_dirs:
  - "./evidence/"
```

---

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Web UI |
| GET | `/health` | Health check |
| POST | `/api/v1/documents` | Upload document, extract claims |
| POST | `/api/v1/evidence` | Upload evidence file |
| POST | `/api/v1/analyze` | Run full analysis |
| GET | `/api/v1/claims` | List persisted claims |
| GET | `/api/v1/assumptions` | List persisted assumptions |
| GET | `/api/v1/gaps` | List persisted gaps |
| GET | `/api/v1/graph` | Export graph as JSON |
| GET | `/api/v1/summary` | Summary statistics |

### Start the API server

```bash
uvicorn asf.api.app:app --reload
# Open http://localhost:8000 for the web UI
```

---

## Assumption Types

| Type | Description | Example |
|------|-------------|---------|
| `IDENTITY` | User identity & authentication | "MFA is required for all admins" |
| `ACCESS` | Who can access what | "Only Finance can access payroll" |
| `NETWORK` | Network posture & isolation | "DBs are not internet accessible" |
| `CONFIGURATION` | System configuration | "All data is encrypted" |
| `PROCESS` | Operational processes | "Changes must be approved" |
| `DOCUMENTATION` | Doc accuracy | "Architecture reflects reality" |
| `DEPENDENCY` | Service relationships | "App depends on auth service" |
| `GOVERNANCE` | Compliance & reviews | "Reviews are conducted quarterly" |

---

## Verification Results

| Result | Meaning |
|--------|---------|
| `VERIFIED` | Evidence supports the assumption |
| `CONTRADICTED` | Evidence contradicts the assumption |
| `PARTIALLY_VERIFIED` | Some evidence supports, some contradicts |
| `UNKNOWN` | No matching evidence available |

---

## Confidence Scoring

Confidence (0.0 – 1.0) is computed from four factors:

| Factor | Weight | Description |
|--------|--------|-------------|
| Verification confidence | 40% | How certain the matching logic is |
| Evidence freshness | 20% | How recent the evidence is |
| Evidence coverage | 20% | How much evidence was used vs available |
| Evidence completeness | 20% | How conclusive the verification result is |

---

## Understanding Findings

Each finding includes:

- **Assumption** — what the system assumed
- **Status** — VERIFIED, CONTRADICTED, PARTIALLY_VERIFIED, or UNKNOWN
- **Evidence Used** — number of evidence sources matched
- **Confidence** — 0-100%
- **Gap Type** — ACCESS_GAP, NETWORK_GAP, etc.
- **Severity** — CRITICAL, HIGH, MEDIUM, LOW, INFO
- **Explanation** — human-readable reason for the finding

---

## Sample Workflow

```bash
# 1. Create config
asf init

# 2. Place your policies in ./policies/ and evidence in ./evidence/

# 3. Run analysis
asf analyze ./policies/ -e ./evidence/ --persist

# 4. Start API to explore results in the browser
uvicorn asf.api.app:app --reload

# 5. Export graph
asf analyze ./policies/ -e ./evidence/ --graph > graph.json
```

---

## LLM Support (Optional)

The system works fully without an LLM. If you want AI-enhanced extraction:

```bash
# OpenAI-compatible API
export ASF_LLM_API_KEY="sk-..."
export ASF_LLM_MODEL="gpt-4o"
asf analyze policy.txt -e iam.csv

# Local Ollama
export ASF_LLM_PROVIDER="ollama"
export ASF_LLM_BASE_URL="http://localhost:11434"
asf analyze policy.txt -e iam.csv
```
