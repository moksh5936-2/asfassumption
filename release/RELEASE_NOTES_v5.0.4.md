# ASF0 v5.0.4 — Semantic Contradiction Engine

## Overview

ASF0 v5.0.4 introduces the Semantic Contradiction Engine — a knowledge-based reasoning layer that detects security contradictions a senior security architect would immediately recognize but ASF0 previously missed.

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
- **Availability** — high availability vs single instance, redundant vs single path
- **Trust** — zero trust vs implicit trust
- **Compliance** — compliant vs audit failure, verified vs unverified

### Detection Improvements

- Word-boundary matching prevents "http" matching inside "HTTPS" (false positive fix)
- Negation detection for "no X" and "without X" patterns
- "Un" prefix detection for "unencrypted", "unrestricted", etc.
- Security controls section excluded from semantic matching (structured labels are not semantic claims)
- Cross-domain encryption exclusion (transport vs storage pairs)
- Self-contradiction guard (single claim matching both concepts)
- Text-based dedup (duplicate statement pairs collapsed)

## Fixed Issues

- Previously obvious semantic contradictions could return 0 contradictions
- Encrypted vs plaintext conflicts now detected
- MFA vs password-only conflicts now detected
- Staff-only vs vendor access conflicts now detected
- Private vs public exposure conflicts now detected
- Logged vs unlogged access conflicts now detected
- False positives from HTTP matching inside HTTPS eliminated
- Clean architecture regression protected

## Benchmark Coverage

| Benchmark | Contradiction Types | Status |
|---|---|---|
| Healthcare | ENCRYPTION, AUTHENTICATION, AUTHORIZATION, NETWORK, COMPLIANCE | ✅ |
| Payroll | ENCRYPTION, AUTHENTICATION, AUTHORIZATION | ✅ |
| Cloud | ENCRYPTION, NETWORK, AUTHORIZATION, LOGGING | ✅ |
| Clean Architecture | 0 false positives | ✅ |

## Installation

```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash
```

Or download a binary directly from the GitHub release.

## Upgrade Notes

Upgrade from any previous version:

```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash -s -- --upgrade
```

## Known Limitations

- Contradiction severity is based on pair configuration, not dynamic context
- Cross-domain exclusion only applies to ENCRYPTION (transport vs storage)
- The knowledge base currently has 26 pairs — community contributions welcome
