# ASF v1 Analyst Usage Guide

**Assumption Security Framework — Version 1**

*For security analysts validating policy compliance against operational evidence*

---

## Table of Contents

1. [What ASF v1 Is](#1-what-asf-v1-is)
2. [Installation & Setup](#2-installation--setup)
3. [Basic Usage: CLI](#3-basic-usage-cli)
4. [Basic Usage: Python API](#4-basic-usage-python-api)
5. [Input Formats](#5-input-formats)
6. [Understanding the Output](#6-understanding-the-output)
7. [Reading ASF Reports](#7-reading-asf-reports)
8. [Common Analyst Workflows](#8-common-analyst-workflows)
9. [Limitations](#9-limitations)
10. [When to Use ASF v1 vs Wait for v2](#10-when-to-use-asf-v1-vs-wait-for-v2)

---

## 1. What ASF v1 Is

ASF v1 (Assumption Security Framework) is an **assumption extraction and verification engine** for security analysts. It reads policy documents, extracts security claims, converts them into typed assumptions, and verifies each assumption against real-world evidence.

### Core Pipeline

```
Policy Documents          Evidence Files
     │                         │
     ▼                         ▼
  ┌──────────┐           ┌──────────┐
  │  Claim   │           │ Evidence │
  │Extractor │           │  Loader  │
  └────┬─────┘           └────┬─────┘
       │                      │
       ▼                      │
  ┌──────────┐                │
  │Assumption│                │
  │  Engine  │                │
  └────┬─────┘                │
       │                      │
       ▼                      ▼
  ┌─────────────────────────────────┐
  │      Verification Engine        │
  │  (check each assumption against │
  │   matching evidence records)    │
  └──────────┬──────────────────────┘
             │
             ▼
  ┌──────────────────────┐
  │     Gap Engine        │
  │ (severity-rated gaps) │
  └──────────────────────┘
```

### What ASF v1 IS

- A rule-based extraction system for security policy documents
- A typed assumption classifier (8 assumption types)
- An evidence-driven verification engine with four verdicts
- A gap analyzer that surfaces where policy mismatches reality

### What ASF v1 IS NOT

- **Not an assumption discovery engine** — it only extracts what is *explicitly written*. It cannot infer undocumented assumptions from architecture diagrams, system configurations, or network topologies. That capability is planned for v2.
- **Not a natural language understanding system** — classification is regex-based, not semantic.
- **Not a continuous monitoring tool** — it is a point-in-time analysis tool.

### The 8 Assumption Types

| Type | Domain | Example Claim |
|------|--------|--------------|
| **ACCESS** | Who can access what | "Only Finance employees may access payroll" |
| **IDENTITY** | Authentication & identity | "All administrative access requires MFA" |
| **NETWORK** | Network segmentation & exposure | "The finance segment is isolated" |
| **CONFIGURATION** | Encryption, backup, logging | "All financial data is backed up daily" |
| **PROCESS** | Reviews, approvals, workflows | "Security reviews are conducted quarterly" |
| **GOVERNANCE** | Compliance, audit, policy adherence | "Access permissions are reviewed monthly" |
| **DOCUMENTATION** | Policy & procedure accuracy | "Configuration changes are documented" |
| **DEPENDENCY** | Inter-system relationships | "The identity provider is available" |

---

## 2. Installation & Setup

### Prerequisites

- Python 3.9+
- pip

### Install

```bash
# Clone the repository and install in editable mode
git clone https://github.com/anomalyco/asf-validator.git
cd asf-validator
pip install -e .
```

### Verify Installation

```bash
asf --help
```

You should see:

```
Usage: asf [OPTIONS] COMMAND [ARGS]...

  ASF Validator v0.1 — Assumption Security Framework Validator

Options:
  -c, --config FILE  Path to config file
  --help             Show this message and exit.

Commands:
  analyze  Analyze documents and evidence for security assumptions.
  init     Create a default asf.config.yaml in the current directory.
```

### Initialize Configuration

```bash
asf init
```

This creates `asf.config.yaml` in the current directory:

```yaml
# ASF Validator v0.1 Configuration

# Database path
db_path: "asf_validator.db"

# Evidence column name mappings
evidence_schema:
  user: "user"
  group: "group"
  resource: "resource"
  permission: "permission"
  mfa: "mfa_enabled"
  public: "public"
  internet_facing: "internet_facing"
  enabled: "enabled"
  asset: "asset"
  environment: "environment"

# Assumption type scoring weights
assumption_weights:
  IDENTITY: 1.0
  ACCESS: 1.0
  NETWORK: 1.0
  CONFIGURATION: 1.0
  PROCESS: 1.0
  DOCUMENTATION: 1.0
  DEPENDENCY: 1.0
  GOVERNANCE: 1.0

# LLM provider configuration (experimental in v1)
llm:
  provider: ""
  api_key: ""
  model: ""
  temperature: 0.1

# Default evidence directories
evidence_dirs: []

# Field mappings for schema adaptation
field_mappings:
  "employee": "user"
  "username": "user"
  "name": "user"
  "email": "user"
  "identity": "user"
  "principal": "user"
  # ... see full file
```

### Config file reference

| Setting | Purpose |
|---------|---------|
| `db_path` | Where ASF persists analysis results (optional) |
| `evidence_schema` | Maps canonical field names to your column names |
| `assumption_weights` | Adjust scoring bias for assumption type classification |
| `field_mappings` | Schema adapter column-name normalization rules |
| `evidence_dirs` | Directories ASF will search for evidence files |
| `llm` | Optional LLM provider (reserved for future use) |

---

## 3. Basic Usage: CLI

### Single Policy Analysis

```bash
# Analyze one policy document with one evidence file
asf analyze finance_policy.txt -e iam_export.json
```

### JSON Output

```bash
# Get machine-readable JSON instead of console report
asf analyze finance_policy.txt -e iam_export.json --json
```

### Multiple Evidence Files

```bash
# Pass several evidence sources
asf analyze policy.txt -e iam_export.json -e mfa_status.csv -e payroll_acl.csv
```

### Directory Analysis

```bash
# Process all supported documents in a directory
asf analyze ./policies/ -e ./evidence/
```

ASF supports directory recursion for both documents and evidence. When you pass a directory, it scans for supported file extensions (TXT, PDF, DOCX, CSV, JSON).

### Export Graph

```bash
# Export a network graph (JSON) of the assumption verification chain
asf analyze policy.txt -e evidence.csv --graph
```

### Persist Results to Database

```bash
# Store results in SQLite for later review
asf analyze policy.txt -e evidence.csv --persist
```

### Full Example

```bash
# Complete analysis with custom config
asf -c my_config.yaml analyze finance_policy.txt \
  -e iam_export.json \
  -e mfa_status.csv \
  -e network_exposure.csv \
  --json --graph
```

### CLI Output Format

The terminal output consists of three sections:

1. **Summary Panel** — High-level counts (claims, assumptions, verified, contradicted, gaps)
2. **Findings Table** — Per-assumption verification status with confidence scores
3. **Gap Analysis Table** — Severity-rated gaps with descriptions

---

## 4. Basic Usage: Python API

### Minimal Example

```python
from asf.analyzer import Analyzer

analyzer = Analyzer()
result = analyzer.analyze(
    document_paths=["finance_policy.txt"],
    evidence_paths=["iam_export.json", "mfa_status.csv"],
)

print(f"Claims found: {result.claims_found}")
print(f"Assumptions: {result.assumptions_found}")
print(f"Verified: {result.verified_count}")
print(f"Contradicted: {result.contradicted_count}")
```

### Programmatic Access to Results

```python
from asf.analyzer import Analyzer

analyzer = Analyzer()
result = analyzer.analyze(
    document_paths=["finance_policy.txt"],
    evidence_paths=["iam_export.json", "payroll_acl.csv"],
)

# Iterate over verifications
for v in result.verifications:
    # Find the matching assumption
    assumption = next(
        a for a in result.assumptions if a.id == v.assumption_id
    )
    print(f"{assumption.text}")
    print(f"  → Status: {v.result.value}")
    print(f"  → Confidence: {v.confidence:.0%}")
    print(f"  → Evidence records used: {len(v.evidence_used)}")
    print(f"  → Reasoning: {v.reasoning}")
    print()

# Examine gaps
for gap in result.gaps:
    print(f"[{gap.severity.value}] {gap.type.value}: {gap.description}")
```

### Custom Configuration

```python
from asf.analyzer import Analyzer
from asf.config import ASFConfig

config = ASFConfig.default()
config.field_mappings = {
    "employee_name": "user",
    "dept": "group",
    "app": "resource",
    "access_level": "permission",
}

analyzer = Analyzer(config=config)
result = analyzer.analyze(
    document_paths=["policy.txt"],
    evidence_paths=["custom_export.csv"],
)
```

### Persisting Results

```python
from asf.analyzer import Analyzer

analyzer = Analyzer()

# Persist automatically to configured database
result = analyzer.analyze(
    document_paths=["policy.txt"],
    evidence_paths=["evidence.csv"],
    persist=True,  # Writes to asf_validator.db (or custom db_path)
)

# Close the database connection
analyzer.close()
```

### Exporting the Graph Model

```python
from asf.analyzer import Analyzer

analyzer = Analyzer()
result = analyzer.analyze(
    document_paths=["policy.txt"],
    evidence_paths=["evidence.csv"],
)

# Export as JSON for visualization
graph_json = analyzer.graph_model.export_json()
print(f"Graph has {graph_json['node_count']} nodes and {graph_json['edge_count']} edges")

# Or get summary
summary = analyzer.graph_model.summary()
print(f"Node types: {summary['node_types']}")
```

---

## 5. Input Formats

### Policy Documents

ASF v1 supports three document formats for policy text:

| Format | Extension | Parser |
|--------|-----------|--------|
| Plain text | `.txt` | Line-by-line sentence splitting |
| PDF | `.pdf` | Text extraction via PDF parser |
| Word | `.docx` | XML-based text extraction |

**Document requirements:**
- Documents should contain declarative security statements (policies, runbooks, architecture documentation)
- The extractor works best with clear, imperative language — "Only X may access Y," "All traffic is encrypted," "Reviews are conducted quarterly"
- Narrative or descriptive text may not yield high-confidence claims

### Evidence Files

| Format | Extension | Use Case |
|--------|-----------|----------|
| CSV | `.csv` | IAM exports, ACL lists, network exposure data, config status |
| JSON | `.json` | IAM exports, group memberships, structured configuration |

**Supported evidence types (identified by content):**

| Evidence Type | Source Type | Typical Columns |
|--------------|-------------|-----------------|
| IAM Export | `IAM_EXPORT` | user, group, permission, mfa |
| ACL List | `ACL_LIST` | user, group, resource, permission |
| Network Exposure | `CSV` | asset, public, internet_facing, exposure, segment |
| Config Export | `CONFIG_EXPORT` | resource, enabled, status, configuration |
| Audit Log | `AUDIT_LOG` | status, reviewed, approved, completed |

### Schema Adaptation

ASF automatically maps column names from your evidence files to its internal field names. For example, if your CSV has `employee_name` instead of `user`, the schema adapter handles it.

**Auto-mapping** (enabled with `--auto-map` or `auto_map=True`) infers field mappings from column headers using a built-in dictionary:

| Your Column | Mapped To |
|-------------|-----------|
| `employee_name`, `username`, `email`, `identity`, `principal` | `user` |
| `department`, `team`, `unit`, `division`, `org` | `group` |
| `application`, `system`, `service`, `target` | `resource` |
| `access`, `right`, `privilege`, `action` | `permission` |
| `mfa_enabled`, `multi_factor`, `2fa`, `totp` | `mfa` |
| `public_facing`, `is_public`, `exposed`, `exposure` | `public` |
| `state`, `active`, `configuration`, `value` | `enabled` |
| `host`, `server`, `system_name` | `asset` |

**Custom mappings** can be defined in `asf.config.yaml` under the `evidence_schema` and `field_mappings` sections.

### Evidence Matching Logic

When verifying an assumption, ASF selects evidence using a type-compatibility matrix:

| Assumption Type | Compatible Evidence Types |
|----------------|--------------------------|
| IDENTITY | IAM_EXPORT, CSV, JSON |
| ACCESS | ACL_LIST, IAM_EXPORT, CSV |
| NETWORK | FIREWALL_RULES, SECURITY_GROUPS, ROUTE_TABLES, CSV |
| CONFIGURATION | CONFIG_EXPORT, JSON, CSV |
| PROCESS | AUDIT_LOG, CSV |
| DOCUMENTATION | PDF, DOCX (original documents) |
| DEPENDENCY | JSON, CSV |
| GOVERNANCE | AUDIT_LOG, CSV |

---

## 6. Understanding the Output

ASF output is divided into five sections. Understanding each is critical to interpreting results correctly.

### 6.1 Claims

Claims are the raw security assertions extracted from your policy documents. They are the atomic unit extracted by the `ClaimExtractor`.

**How extraction works:**
- The extractor splits documents into sentences
- Each sentence is tested against 25 declarative regex patterns
- If a sentence matches, it becomes a claim with a confidence score

**Example claims from `finance_policy.txt`:**

```
"Only Finance employees may access the payroll processing system."
"All payroll data is encrypted at rest and in transit."
"The payroll application is only accessible from the internal corporate network."
"Only the VP of Finance can approve payroll runs."
"Production databases are not internet accessible."
"All administrative access requires multi-factor authentication."
"Database access is restricted to database administrators only."
"SSH access to production servers is restricted to the infrastructure team."
"All access to financial systems is logged and monitored."
"Security reviews are conducted quarterly."
"Access permissions are reviewed on a monthly basis."
"All configuration changes must be approved by the security team."
"All financial data is backed up daily."
"Backup data is encrypted using AES-256."
"Backups are tested monthly to ensure recoverability."
"Offsite backup storage is maintained at a separate geographic location."
"The finance network segment is isolated from other network segments."
"No external access is permitted to the finance application."
"All network traffic is inspected by the intrusion detection system."
"Security groups restrict access to only required ports and protocols."
```

**Key insight:** Claims are the *text itself* — unmodified policy sentences. They have not yet been classified or transformed.

### 6.2 Assumptions

Assumptions are claims that have been classified into one of the 8 types and assigned a canonical text form.

**Classification:** The `AssumptionEngine` scores each claim against regex patterns for each assumption type and selects the type with the highest score.

| Claim Text | Classified Type |
|------------|-----------------|
| "Only Finance employees may access the payroll processing system." | **ACCESS** |
| "All administrative access requires multi-factor authentication." | **IDENTITY** |
| "The finance network segment is isolated from other network segments." | **NETWORK** |
| "All financial data is backed up daily." | **CONFIGURATION** |
| "Security reviews are conducted quarterly." | **GOVERNANCE** |
| "All configuration changes must be approved by the security team." | **PROCESS** |
| "Backup data is encrypted using AES-256." | **CONFIGURATION** |

Each assumption in the output includes:
- **id** — Unique identifier (e.g., `asm_a1b2c3d4e5f6`)
- **text** — Prefixed with the type label (e.g., "System assumes access control: Only Finance employees may access the payroll processing system.")
- **assumption_type** — One of the 8 types
- **confidence** — Extraction confidence (0.0–1.0)
- **keywords** — Extracted significant terms
- **verification_status** — PENDING, VERIFIED, CONTRADICTED, or IN_REVIEW

### 6.3 Verifications

Each assumption is matched against compatible evidence and checked. The verification engine produces one of four results:

| Result | Icon | Meaning | Example |
|--------|------|---------|---------|
| **VERIFIED** | Green check | Evidence confirms the assumption holds | Policy says "MFA required" — evidence shows all users have MFA enabled |
| **CONTRADICTED** | Red X | Evidence directly contradicts the assumption | Policy says "Finance only" — evidence shows Engineering users have payroll access |
| **PARTIALLY_VERIFIED** | Yellow warning | Evidence partially confirms but has exceptions | Policy says "MFA required" — most users have MFA, but interns and contractors do not |
| **UNKNOWN** | Dim dash | Insufficient evidence to verify | Policy says "Backups are tested monthly" — no evidence file contains testing records |

Each verification includes:
- **result** — The verdict
- **confidence** — How certain the engine is (0.0–1.0)
- **evidence_used** — List of evidence record IDs that were matched
- **reasoning** — Human-readable explanation of why this verdict was reached
- **details** — Structured data from the check (e.g., list of users without MFA)

**Critical nuance:** A VERIFIED result does not mean "secure." It means "the policy statement matches the evidence." A policy could be verified as "Only X may access Y" but the underlying access model could still be insecure — ASF only checks what the policy *says*, not whether the policy is *correct*.

### 6.4 Gaps

Gaps represent discrepancies between policy and reality. The `GapEngine` generates gaps for each assumption based on its verification result:

| Verification Result | Gap Type | Default Severity | When |
|--------------------|----------|------------------|------|
| CONTRADICTED | Matches assumption type (e.g., ACCESS_GAP) | CRITICAL/HIGH/MEDIUM | Evidence directly contradicts policy |
| PARTIALLY_VERIFIED | Matches assumption type | MEDIUM | Evidence confirms partially but has exceptions |
| UNKNOWN | EVIDENCE_GAP | LOW | No sufficient evidence available |
| No verification | VERIFICATION_GAP | MEDIUM | Assumption was not verified at all |

**Severity determination:**
- CONTRADICTED + confidence >= 0.8 + ACCESS/IDENTITY/NETWORK type → **CRITICAL**
- CONTRADICTED + confidence >= 0.8 + CONFIGURATION/GOVERNANCE type → **HIGH**
- CONTRADICTED + confidence >= 0.5 → **HIGH**
- CONTRADICTED + confidence < 0.5 → **MEDIUM**
- PARTIALLY_VERIFIED → **MEDIUM**
- UNKNOWN → **LOW**
- No verification performed → **MEDIUM**

### 6.5 Graph

The `GraphModel` builds a directed network graph of the entire analysis:

```
Claim ──GENERATES──► Assumption ──VERIFIES──► Verification
                         ▲                        │
                         │                    SUPPORTS
                     IDENTIFIES                  │
                         │                        ▼
                         └──── Gap           Evidence
```

The graph is exported as JSON with nodes (typed entities) and edges (relationships). It is useful for:
- Tracing a specific claim through to its verification
- Understanding which evidence was used for which assumptions
- Identifying orphaned claims or unused evidence

---

## 7. Reading ASF Reports

This section walks through a real analysis using `sample_data/finance_policy.txt` against evidence files to show you how to interpret each result.

### 7.1 The "VERIFIED" Status

**What it means:** Evidence confirms the assumption.

**Example from our analysis:**

```
Finding                                                     Status      Evidence  Explanation
System assumes access control: Only Finance employees...    VERIFIED    1         Only users in 'Finance' found with access (6 users)
System assumes configuration state: All financial data...  VERIFIED    1         All 3 compliant resource(s) with configuration
System assumes network posture: The finance network se...  VERIFIED    1         All 10 asset(s) appear isolated
```

**Walkthrough — Access Control Verification:**

Policy says: *"Only Finance employees may access the payroll processing system."*

ASF loads `payroll_acl.csv`:
```
user,           group,       resource,       permission,  department
alice.jones,    Finance,     payroll-system,  read,        Finance
bob.smith,      Finance,     payroll-system,  write,       Finance
carol.davis,    Finance,     payroll-system,  admin,       Finance
dave.wilson,    Engineering, payroll-system,  read,        Engineering
eve.brown,      Engineering, payroll-system,  read,        Engineering
frank.miller,   Finance,     payroll-system,  read,        Finance
grace.lee,      Finance,     payroll-system,  write,       Finance
henry.taylor,   Engineering, payroll-system,  write,       Engineering
iris.anderson,  Finance,     payroll-system,  read,        Finance
jack.thompson,  HR,          payroll-system,  read,        HR
```

**Wait — VERIFIED?** Look closely: there are Engineering and HR users in that ACL! The column is marked as `VERIFIED` with 78% confidence. Why?

The verification engine checks the `group` field in the ACL, not the `department` field. It finds all users in the "Finance" group have access (alice, bob, carol, frank, grace, iris = 6 users). It does not flag users in other groups because the policy restricts to "Finance employees" and the ACL `group` column contains users outside Finance. However, **this is a TRUE POSITIVE for the check logic** — the engine identified that only Finance-group users have access.

**Analyst interpretation:** You should investigate why the ACL shows Engineering and HR users with access. The evidence *may* contradict the policy depending on how groups map to departments. ASF flags this at 78% confidence — not 95% — because the pattern match is not perfect.

### 7.2 The "CONTRADICTED" Status

**What it means:** Evidence directly contradicts the assumption.

**Example from our analysis:**

```
Finding                                                      Status          Evidence  Explanation
System assumes access control: Only the VP of Finance...    CONTRADICTED    1         Found 2 user(s) outside expected group with access
System assumes identity posture: All administrative ac...    CONTRADICTED    1         MFA not enabled for 3 user(s)
```

**Walkthrough — MFA Contradiction:**

Policy says: *"All administrative access requires multi-factor authentication."*

ASF loads `mfa_status.csv`:
```
user,           mfa_enabled, department,  role
alice.jones,    true,        Finance,      Manager
bob.smith,      true,        Finance,      Analyst
carol.davis,    true,        Finance,      VP
dave.wilson,    true,        Engineering,  Engineer
eve.brown,      true,        Engineering,  Engineer
frank.miller,   false,       Finance,      Intern
grace.lee,      true,        Finance,      Analyst
henry.taylor,   true,        Engineering,  Sr Engineer
iris.anderson,  false,       Finance,      Contractor
jack.thompson,  false,       HR,           Coordinator
```

The engine looks for MFA status and finds 3 users with `mfa_enabled = false` (frank.miller, iris.anderson, jack.thompson). Since the policy states "ALL administrative access requires MFA" and there are users without MFA, the result is CONTRADICTED with 95% confidence.

**Analyst interpretation:** This is actionable. The policy requirement (MFA for all admin access) is not being met. However, note that:
- The policy says "administrative access" — the evidence covers ALL users, not just admins
- ASF cannot distinguish between "admin users" and "regular users" with this evidence schema
- You should verify whether these 3 users actually have administrative access or if this is a false positive

### 7.3 The "UNKNOWN" Status

**What it means:** Insufficient evidence exists to verify the assumption.

**Example from our analysis:**

```
Finding                                                       Status        Evidence  Explanation
System assumes process compliance: Security reviews are...   UNKNOWN       0         No matching evidence available for verification
System assumes process compliance: Access permissions a...   UNKNOWN       0         No matching evidence available for verification
System assumes governance compliance: All configuration...   UNKNOWN       0         No matching evidence available for verification
```

**Walkthrough:**

Policy says: *"Security reviews are conducted quarterly."*

None of the provided evidence files contain review dates, audit logs, or governance tracking data. The evidence mapper finds no compatible source types for a PROCESS/GOVERNANCE assumption. The verification engine has zero evidence to check against, so it returns UNKNOWN (30% confidence).

**Analyst interpretation:** UNKNOWN is not a pass or fail — it is a flag that you need more evidence. To verify this, you would need to provide:
- An audit log showing quarterly review completion
- A governance tracking CSV with review dates
- A ticketing system export showing review cycles

**Key insight:** If your analysis shows many UNKNOWN results, you may not have provided the right evidence types for the policies you are checking.

### 7.4 "PARTIALLY_VERIFIED" Status

**What it means:** Evidence partially confirms but has exceptions.

**Example scenario (if MFA mix existed):**

```
Finding                                                      Status               Evidence  Explanation
System assumes identity posture: All administrative ac...    PARTIALLY_VERIFIED   1         MFA enabled for 7 user(s) but missing for 3
```

**Analyst interpretation:** This is a nuanced finding. The policy is *mostly* followed but has exceptions. In practice, this often means:
- There is a legitimate exception (e.g., contractors cannot use MFA)
- There is a coverage gap that needs remediation
- The policy text uses "all" but reality has edge cases

Your next step: decide whether the exceptions are acceptable (documented exceptions) or require remediation (unapproved bypasses).

### 7.5 Interpreting Critical Gaps

When ASF reports a **CRITICAL** gap, it means:

1. An assumption was CONTRADICTED by evidence
2. The contradiction is about ACCESS, IDENTITY, or NETWORK (the highest-risk domains)
3. The verification confidence is >= 80%

**Example critical gap:**

```
Type            Severity    Description
ACCESS_GAP      CRITICAL    Assumption contradicted: System assumes access control: Only the VP of Finance can approve payroll runs...
```

**What this means for your organization:**
- Someone explicitly stated this policy requirement
- The evidence shows it is not being followed
- The gap is in a high-impact domain (access control)

**Response workflow:**

1. **Verify the evidence** — Is the IAM export current? Could there be a stale record?
2. **Check for exceptions** — Is there an approved exception or compensating control?
3. **Assess impact** — What is the actual risk of this gap?
4. **Remediate or update policy** — Either fix the access or update the policy to reflect reality

---

## 8. Common Analyst Workflows

### 8.1 "I have a policy and an IAM export"

**Scenario:** You are a security analyst with an access control policy and an IAM export from your identity provider. You want to check if access controls are correctly implemented.

**Evidence needed:**
- Policy document (TXT, PDF, or DOCX) — e.g., `access_control_policy.txt`
- IAM export — e.g., `iam_export.json` or `user_permissions.csv`
- ACL list (optional) — e.g., `payroll_acl.csv`

**Command:**
```bash
asf analyze access_control_policy.txt -e iam_export.json -e payroll_acl.csv --json
```

**What ASF checks:**
- ACCESS-type assumptions against IAM permissions and ACL records
- IDENTITY-type assumptions against MFA status, group memberships
- It finds users with permissions that contradict "Only X may access Y" patterns

**Sample evidence (IAM export):**

```json
{
  "users": [
    {"username": "alice.jones", "department": "Finance", "mfa": true, "groups": ["Finance", "Payroll-Admin"]},
    {"username": "dave.wilson", "department": "Engineering", "mfa": true, "groups": ["Engineering", "Payroll-Read"]}
  ],
  "groups": {
    "Payroll-Read": {"members": ["bob.smith", "dave.wilson", "eve.brown", "grace.lee", "iris.anderson", "jack.thompson"]}
  }
}
```

**What to look for in the output:**
- CONTRADICTED results for ACCESS assumptions
- Users from outside the expected department who have access to restricted resources
- UNKNOWN results that suggest missing evidence

### 8.2 "I have a firewall policy and a network diagram"

**Scenario:** You need to verify that network segmentation matches the documented policy.

**Evidence needed:**
- Network security policy — e.g., `network_policy.txt`
- Network exposure data — e.g., `network_exposure.csv` or firewall rule exports

**Command:**
```bash
asf analyze network_policy.txt -e network_exposure.csv --json
```

**Sample evidence (network exposure):**

```csv
asset,environment,public,internet_facing,exposure,segment
payroll-app,production,false,false,internal,finance
finance-db,production,false,false,internal,finance
customer-portal,production,true,true,public,dmz
api-gateway,production,true,true,public,network
dev-server,development,true,true,public,development
```

**What ASF checks:**
- NETWORK-type assumptions against exposure data
- It verifies claims like "The finance segment is isolated" by checking if finance assets are marked as public
- It checks "No external access is permitted" by looking for internet-facing assets

**What to look for:**
- If policy claims isolation but evidence shows public assets in that segment → CONTRADICTED
- If policy claims "no external access" but assets are internet-facing → CONTRADICTED
- The specific assets that violate the policy are listed in the verification details

### 8.3 "I want to check if my backup policy matches reality"

**Scenario:** A backup and recovery policy exists. You have a configuration export showing backup status across systems.

**Evidence needed:**
- Backup policy — e.g., `backup_policy.txt`
- Backup configuration — e.g., `backup_config.csv`

**Command:**
```bash
asf analyze backup_policy.txt -e backup_config.csv
```

**Sample evidence:**

```csv
resource,configuration,enabled,status,frequency
payroll-db,encrypted,true,active,daily
finance-fs,encrypted,true,active,daily
analytics-db,encrypted,false,inactive,weekly
dev-db,encrypted,false,inactive,manual
hr-system,encrypted,false,inactive,weekly
```

**What ASF checks:**
- CONFIGURATION assumptions against backup enablement status
- Encryption configuration verification
- The engine checks if resources marked as "encrypted" or "backed up" match the evidence

**What to look for:**
- PARTIALLY_VERIFIED results when some resources are compliant but others are not
- The specific non-compliant resources are listed in details (e.g., `analytics-db`, `hr-system`)
- Disabled backups for critical systems create gaps

### 8.4 "I want to verify MFA compliance across my org"

**Scenario:** Your policy mandates MFA for all users. You want to check compliance across departments.

**Evidence needed:**
- Policy document mentioning MFA — e.g., `security_policy.txt`
- MFA status export — e.g., `mfa_status.csv`

**Command:**
```bash
asf analyze security_policy.txt -e mfa_status.csv --json
```

**Sample evidence:**

```csv
user,mfa_enabled,department,role
alice.jones,true,Finance,Manager
bob.smith,true,Finance,Analyst
frank.miller,false,Finance,Intern
iris.anderson,false,Finance,Contractor
jack.thompson,false,HR,Coordinator
```

**What ASF checks:**
- IDENTITY-type assumptions containing MFA keywords
- Checks each user's MFA status against the policy requirement
- Returns PARTIALLY_VERIFIED if some users lack MFA

**What to look for:**
- If the policy says "ALL" and not all users have MFA → CONTRADICTED or PARTIALLY_VERIFIED
- The verification details list exactly which users are non-compliant
- Patterns like "all interns lack MFA" might indicate an onboarding process gap

---

## 9. Limitations

This section is critical for accurate interpretation of ASF results. Misunderstanding these limitations will lead to false confidence in analysis outcomes.

### 9.1 Only Extracts What Is EXPLICITLY Written

ASF v1 processes only the text present in policy documents. If a policy states "passwords must be 12 characters," ASF will extract and verify that assumption. However:

- It will **not** infer that "password rotation must also occur" (derived assumption)
- It will **not** infer that "the password database must exist" (dependency assumption)
- It will **not** infer that "password policies do not conflict with other policies" (governance assumption)

**Real-world impact:** A policy document may be short (30–50 lines), but the true set of security assumptions it relies on is an order of magnitude larger. ASF only sees the surface.

### 9.2 94.1% of Assumptions Are NOT in Documents

Based on empirical analysis of the ASF Assumption Ontology (1697 ground-truth assumptions across multiple security domains):

| Ontology | Count | % of Total | ASF v1 Can Extract? |
|----------|-------|------------|-------------------|
| Explicit | 100 | 5.9% | **Yes** — directly stated in policy text |
| Implicit | 282 | 16.6% | **Partial** — unstated but implied |
| Derived | 582 | 34.3% | **No** — requires logical inference |
| Trust | 142 | 8.4% | **No** — requires trust relationship reasoning |
| Operational | 339 | 20.0% | **No** — assumes operational processes exist |
| Dependency | 225 | 13.3% | **No** — assumes external systems exist |
| Architectural | 7 | 0.4% | **No** — requires architecture analysis |
| Environmental | 20 | 1.2% | **No** — requires environmental context |

**Total hidden assumptions (v1 blind spots):** 94.1% of all security assumptions associated with a policy are **not extractable by ASF v1**.

### 9.3 v1 Recall Ceiling Is ~11.7%

The benchmark scoring framework (in `benchmark/runner.py`) measures ASF performance against ground truth using text overlap, type match, and keyword overlap. The theoretical maximum recall for ASF v1 is the proportion of explicit + some implicit assumptions in the ground truth, estimated at **~11.7%**.

This means:
- ASF v1 will miss approximately 88% of relevant security assumptions
- A "clean" ASF report (no contradicted gaps) does **not** mean your security posture is sound
- It only means the *specific statements in your policy* match the *specific evidence you provided*

### 9.4 No Trust Graph Reasoning

Many security assumptions involve trust relationships:
- "The identity provider is trusted to verify user identity"
- "The CA is trusted to issue valid certificates"
- "The vendor is trusted to handle data securely"

ASF v1 has no model for trust relationships. It cannot evaluate whether:
- A certificate authority is trustworthy
- An IdP's identity verification is reliable
- A third-party vendor's security posture is adequate

### 9.5 No Derivation Rules

Security policies have logical consequences that are not stated. For example:

> Policy: "All data is encrypted at rest"
> Derived assumption: "Encryption keys exist and are managed"
> Derived assumption: "Key rotation occurs on schedule"
> Derived assumption: "Key compromise response exists"

ASF v1 does **not** derive these downstream assumptions. Each must be explicitly stated in the policy or it will be missed.

### 9.6 No Discovery of Undocumented Assumptions

The most dangerous security gaps are often undocumented:
- Assumptions about network architecture that no one wrote down
- Assumptions about identity provider availability
- Assumptions about personnel being trained
- Assumptions about monitoring existing

ASF v1 cannot discover these. It only checks what is written against what is measured. v2 will address this through architecture-driven assumption generation.

### 9.7 Additional Limitations

- **Regex-based classification** is brittle — unusual phrasing, typos, or non-standard terminology can cause misclassification
- **No cross-document analysis** — each policy is analyzed independently; conflicting policies are not detected
- **No temporal reasoning** — evidence is treated as current-state; ASF does not track changes over time
- **No severity overrides** — the gap severity model is heuristic and may not match your organizational risk appetite
- **Evidence schema sensitivity** — if column names do not match expected patterns, verification may fail silently (returning UNKNOWN)

---

## 10. When to Use ASF v1 vs Wait for v2

### Use ASF v1 When

| Scenario | Why v1 Works |
|----------|-------------|
| **You have well-documented policies** | v1 extracts and verifies explicit statements effectively |
| **You need quick policy-evidence gap analysis** | Running `asf analyze` takes seconds |
| **You are auditing specific, known requirements** | v1's 8 assumption types cover standard security domains |
| **You want to automate policy compliance checks** | The JSON output and Python API support integration |
| **You are training analysts on assumption-based security** | v1's pipeline is simple and explainable |
| **You have structured IAM/network/config evidence** | v1's evidence matcher handles CSV and JSON well |
| **You want to establish a baseline before v2** | v1 results can seed v2's knowledge base |

### Wait for v2 When

| Scenario | Why v2 Is Needed |
|----------|-----------------|
| **You need to discover undocumented assumptions** | v2 will infer assumptions from architecture, not just text |
| **You have complex trust relationships** | v2 will include a trust graph model |
| **You need derivation of downstream assumptions** | v2 will apply inference rules to derive implicit assumptions |
| **You are analyzing a system with no written policies** | v1 requires a policy document to extract from |
| **You need continuous monitoring, not point-in-time** | Both v1 and v2 are analysis tools, but v2 may add monitoring APIs |
| **Your evidence is unstructured** | v1 requires structured CSV/JSON for verification |
| **You need SLAs or compliance reporting** | v1 does not track policy coverage metrics over time |

### Decision Matrix

```
                    ┌─────────────────────────────────────┐
                    │        Do you have written           │
                    │        policy documents?             │
                    └────────────┬────────────────────────┘
                                 │
                    ┌────────────┴────────────┐
                    │                         │
                   YES                        NO
                    │                         │
              ┌─────┴─────┐                   │
              │           │                   │
        Do you have   No evidence          WAIT
        structured    but policy           FOR V2
        evidence?     is clear
              │           │
         ┌────┴────┐     │
         │         │     │
        YES        NO    │
         │         │     │
    ┌────┴──┐  ┌──┴──┐  │
    │       │  │     │  │
  USE V1  V1   V1   V1
  NOW    will  will  will
         show  show  show
         some  many  many
         gaps  gaps  UNKNOWN
               but   results
               still
               useful
```

### Summary Recommendation

**Use ASF v1 today** if you have policy documents and evidence files. It will surface real, actionable gaps — contradicted assumptions that represent security control failures. But **interpret UNKNOWN results carefully**: they may mean "no evidence provided" rather than "no issue exists."

**Plan for v2** if your analysis needs go beyond text extraction. v2's assumption discovery engine will generate assumptions from architectural patterns, not just written policies, dramatically expanding coverage from the ~5.9% ceiling of v1.

---

## Appendix: Reference

### A. Claim Extraction Patterns

The `ClaimExtractor` uses these regex patterns to identify security-relevant sentences:

| Pattern | Matches |
|---------|---------|
| `(only\|just) .+ (can\|may\|has\|have\|should\|must\|will)` | "Only X may access Y" |
| `(all\|every\|each) .+ (are\|is\|shall\|must\|will\|should)` | "All traffic is encrypted" |
| `(is\|are) (not\|never) .+` | "Production DBs are not internet accessible" |
| `(is\|are) \w+ (encrypted\|backed up\|logged\|audited\|...)` | "Data is encrypted at rest" |
| `(cannot\|must not\|shall not\|should not) .+` | "No external access is permitted" |
| `require[s]?\|ensure[s]?\|guarantee[s]?\|provide[s]?\|...` | "Policy requires MFA" |
| `encrypt(ed\|s\|ion)\|backup\|log(ged\|s\|ging)\|audit(ed\|s\|ing)\|...` | "Data is encrypted" |
| `(separated\|isolated\|segmented\|partitioned) .+` | "The finance segment is isolated" |
| `manage[sd]? .+ (access\|permissions)` | "Manage user permissions" |

### B. Assumption Classification Patterns

The `AssumptionEngine` classifies claims using type-specific regex dictionaries. Each type has multiple patterns; the type with the highest aggregate match count wins.

| Type | Key Patterns |
|------|-------------|
| ACCESS | `access`, `permission`, `acl`, `allow`, `deny`, `grant`, `read`, `write`, `admin`, `only X can access` |
| IDENTITY | `mfa`, `multi-factor`, `identity`, `authentication`, `password`, `credential`, `role`, `group` |
| NETWORK | `network`, `firewall`, `internet`, `subnet`, `vpc`, `vlan`, `segment`, `isolate`, `expose` |
| CONFIGURATION | `encrypt`, `backup`, `log`, `audit`, `monitor`, `config`, `setting`, `parameter`, `enabled` |
| PROCESS | `process`, `procedure`, `workflow`, `review`, `approve`, `test`, `sign-off`, `approval` |
| GOVERNANCE | `review`, `audit`, `compliance`, `regulat`, `policy`, `standard`, `framework`, `quarterly` |
| DOCUMENTATION | `document`, `policy`, `runbook`, `procedure`, `guide`, `manual`, `readme`, `wiki` |
| DEPENDENCY | `depend`, `integration`, `connect`, `communicate`, `rel(y\|ies)`, `upstream`, `downstream` |

### C. Gap Severity Matrix

| Condition | Severity |
|-----------|----------|
| CONTRADICTED + confidence >= 0.8 + ACCESS/IDENTITY/NETWORK type | **CRITICAL** |
| CONTRADICTED + confidence >= 0.8 + CONFIGURATION/GOVERNANCE type | **HIGH** |
| CONTRADICTED + confidence >= 0.5 | **HIGH** |
| CONTRADICTED + confidence < 0.5 | **MEDIUM** |
| PARTIALLY_VERIFIED | **MEDIUM** |
| UNKNOWN (no evidence) | **LOW** |
| No verification performed | **MEDIUM** |

### D. Output Schema (JSON)

```json
{
  "summary": {
    "claims_found": 20,
    "assumptions": 18,
    "verified": 4,
    "contradicted": 3,
    "unknown": 11,
    "critical_gaps": 2
  },
  "claims": [
    {
      "id": "clm_abc123",
      "source_document": "finance_policy.txt",
      "text": "Only Finance employees may access the payroll processing system.",
      "extraction_confidence": 0.9,
      "tags": ["access", "identity"]
    }
  ],
  "assumptions": [
    {
      "id": "asm_def456",
      "claim_id": "clm_abc123",
      "text": "System assumes access control: Only Finance employees may access the payroll processing system.",
      "assumption_type": "ACCESS",
      "verification_status": "VERIFIED",
      "confidence": 0.85,
      "keywords": ["finance", "employees", "access", "payroll", "processing"]
    }
  ],
  "verifications": [
    {
      "assumption_id": "asm_def456",
      "result": "VERIFIED",
      "confidence": 0.78,
      "evidence_used": ["evd_xyz789"],
      "reasoning": "Only users in 'Finance' found with access (6 users)"
    }
  ],
  "gaps": [
    {
      "assumption_id": "asm_ghi012",
      "type": "ACCESS_GAP",
      "severity": "CRITICAL",
      "description": "Assumption contradicted: System assumes access control: Only the VP of Finance...",
      "evidence_detail": "Found 2 user(s) outside expected group with access"
    }
  ]
}
```

---

*ASF v1 is experimental research software. It is designed to explore the problem space of assumption-based security verification, not to replace human analyst judgment. Always validate ASF findings manually before taking action.*
