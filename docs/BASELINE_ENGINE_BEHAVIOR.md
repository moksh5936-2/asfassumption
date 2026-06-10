# ASF Engine Baseline Behavior

Generated from Python engine `v1.1.0` (commit `6baef2a`) and verified against Go native engine.

## Input: Sample Data

| File | Type | Records |
|------|------|---------|
| `sample_data/finance_policy.txt` | Text policy | 17 extracted claims |
| `sample_data/mfa_status.csv` | MFA evidence | 10 records |
| `sample_data/payroll_acl.csv` | ACL evidence | 10 records |
| `sample_data/network_exposure.csv` | Network evidence | 10 records |
| `sample_data/backup_config.csv` | Configuration evidence | 10 records |

## Summary Statistics

| Metric | Value |
|--------|-------|
| Claims extracted | 17 |
| Assumptions generated | 17 |
| Verified | 0 |
| Contradicted | 7 |
| Unknown | 5 |
| Partially verified (implied) | 5 |
| Critical gaps | 3 |
| Evidence files used | 4 |
| Graph nodes | 72 |
| Graph edges | 136 |

## Assumption Types

| Type | Count |
|------|-------|
| ACCESS | 5 |
| CONFIGURATION | 5 |
| NETWORK | 2 |
| GOVERNANCE | 2 |
| PROCESS | 2 |
| IDENTITY | 1 |

## Gap Types & Severities

| Gap Type | Severity | Count |
|----------|----------|-------|
| ACCESS_GAP | CRITICAL | 1 |
| IDENTITY_GAP | CRITICAL | 1 |
| NETWORK_GAP | CRITICAL | 1 |
| PROCESS_GAP | HIGH | 2 |
| GOVERNANCE_GAP | HIGH | 2 |
| CONFIGURATION_GAP | MEDIUM | 5 |
| EVIDENCE_GAP | LOW | 5 |

## Verification Results

CONTRADICTED appears for:
- ACCESS claims where evidence shows users outside the restricted group have access
- IDENTITY claims where MFA is not enforced for all users
- NETWORK claims where assets claimed isolated are publicly exposed
- PROCESS/GOVERNANCE claims where evidence shows pending reviews

UNKNOWN appears for:
- Claims where evidence doesn't match the assumption type (e.g., ACCESS claim with no user/permission columns in evidence)
- Claims where no compatible evidence columns can be found

PARTIALLY_VERIFIED appears for:
- CONFIGURATION claims where some resources are compliant and some aren't
- IDENTITY claims where MFA is partially implemented

## Confidence Scoring

Confidence is computed from:
1. **Freshness** (30%): How recent the evidence timestamps are (newer = higher)
2. **Coverage** (40%): Ratio of matched evidence to total available evidence
3. **Completeness** (30%): How detailed the verification reasoning is

Base confidence comes from the verification check function, then multiplied by weighted score.

## Graph Model

The graph is a directed multi-graph with 72 nodes and 136 edges:

- **Claim** nodes (17): Text claims extracted from documents
- **Assumption** nodes (17): Typed assumptions derived from claims
- **Evidence** nodes (4): Evidence sources loaded
- **Verification** nodes (17): Verification results per assumption
- **Gap** nodes (17): Gaps identified per non-verified assumption

Edge relationships: GENERATES, VERIFIES, SUPPORTS, CONTRADICTS, RELATES_TO, IDENTIFIES

## Fixture Data

Baseline outputs saved in `asf-tui/testdata/baseline/`:
- `baseline_json_output.json` — Python engine output (--json)
- `baseline_graph_output.json` — Python engine output (--graph)
- `native_json_output.json` — Go native engine output
- `native_graph_output.json` — Go native engine output

## Known Differences (Go vs Python)

1. **Confidence values differ slightly** — Same algorithm but floating-point arithmetic and randomized IDs cause minor variance. Functional meaning (high/medium/low) is consistent.
2. **IDs differ** — Random hex IDs generated per invocation; never compared across runs.
3. **Timestamps differ** — Created at analysis time; never compared across runs.
4. **No evidence in --json output** — Both engines omit the evidence list from the JSON summary output (it's only in the graph export).
