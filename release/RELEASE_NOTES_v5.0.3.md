# ASF0 v5.0.3 — Semantic Contradiction Engine

## Overview

ASF0 v5.0.3 introduces the Semantic Contradiction Engine — a knowledge-based reasoning layer that detects security contradictions a senior security architect would immediately recognize but ASF0 previously missed.

## Major Feature

### Semantic Contradiction Engine

Previously, ASF0 could only detect structural contradictions (e.g., same ID, duplicate statements). Obvious semantic contradictions like "All data is encrypted at rest" paired with "Backups are stored in plaintext" produced 0 findings.

The new engine detects 26 contradiction patterns across 9 categories:

- **Encryption** — encrypted vs plaintext, at-rest vs backup plaintext
- **Authentication** — MFA vs password-only, verified identity vs anonymous
- **Authorization** — restricted vs open, least privilege vs full access, staff vs vendor, admin vs everyone
- **Logging** — all access logged vs not logged, audit present vs disabled, monitoring enabled vs absent
- **Network** — private vs public, isolated vs direct, segmented vs flat, internal vs internet
- **Backup** — encrypted vs plaintext, redundant vs single, replicated vs single-region
- **Availability** — HA vs single instance, redundant path vs single path
- **Trust** — zero trust vs implicit trust, verified identity vs trusted by default
- **Compliance** — HIPAA vs unencrypted PHI, PCI vs card data exposed, SOC2 vs logging absent

Matching uses direct substring and token-based word matching with:
- Negation detection ("no two-factor" does not match "mfa")
- "Un" prefix detection ("unencrypted" does not match "encrypted")
- Self-contradiction guard (claims matching both sides are skipped)
- Cross-domain encryption context exclusion (transport vs storage)
- Text-based deduplication

## Benchmark Coverage

| Architecture | Contradictions | Types |
|---|---|---|
| Healthcare | 5 | ENCRYPTION, AUTHENTICATION, AUTHORIZATION, NETWORK, COMPLIANCE |
| Payroll | 3 | ENCRYPTION, AUTHENTICATION, AUTHORIZATION |
| Cloud | 4 | ENCRYPTION, NETWORK, AUTHORIZATION, LOGGING |
| Clean Architecture | 0 | Zero false-positive regression |

## Fixed Issues

- Previously obvious semantic contradictions returned 0 contradictions
- Encrypted vs plaintext conflicts now detected across different wording
- MFA vs password-only conflicts now detected with synonym support
- Staff-only vs vendor access conflicts now detected
- Private vs public exposure conflicts now detected
- Logged vs unlogged access conflicts now detected
- "without encryption" now correctly matches "unencrypted" concept
- "unencrypted" no longer falsely matches "encrypted" concept
- "unrestricted" no longer falsely matches "restricted" concept
- Security control labels no longer produce false contradictions
- Self-comparison and duplicate contradictions eliminated

## Installation

```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash
```

Or download a binary from the release assets.

## Upgrade

```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash -s -- --upgrade
```

## Known Limitations

- Token-based matching may produce edge-case false positives for distant word matches
- "un" prefix detection checks only the first occurrence of a word in text
- Security control labels are excluded from semantic matching (they are structured identifiers, not natural language statements)
- Windows installer (install.ps1) is not provided; use manual binary download
