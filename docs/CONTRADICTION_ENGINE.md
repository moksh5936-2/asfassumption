# ASF Contradiction Engine

## Overview

The contradiction engine detects logical inconsistencies between security assumptions in an architecture. It operates deterministically on the full assumption set after both native and intelligence generation.

## Detection Rules

### 1. MFA_ENFORCED_WITH_EXEMPTION
**Trigger:** "MFA is enforced" + "service accounts are exempt"
**Severity:** Critical
**Explanation:** MFA is enforced for users but service accounts are exempted, creating a privileged bypass path.
**Evidence:** IDs of both assumptions

### 2. ENCRYPTED_WITH_PLAINTEXT_BACKUP
**Trigger:** "encrypted" + "plaintext backup" / "unencrypted backup"
**Severity:** Critical
**Explanation:** Data is encrypted at rest but backups are stored in plaintext, defeating the protection.
**Evidence:** IDs of both assumptions

### 3. LEAST_PRIVILEGE_WITH_SHARED_ADMIN
**Trigger:** "least privilege" + "shared admin account" / "shared credentials"
**Severity:** Critical
**Explanation:** Least privilege is claimed but administrators share a single account, violating the principle.
**Evidence:** IDs of both assumptions

### 4. PRIVATE_SUBNET_WITH_INTERNET_ACCESS
**Trigger:** "private subnet" + "internet accessible" / "public access"
**Severity:** High
**Explanation:** Component is claimed to be in a private subnet but is accessible from the internet.
**Evidence:** IDs of both assumptions

### 5. IMMUTABLE_AUDIT_WITH_LOG_DELETION
**Trigger:** "immutable audit" / "tamper-proof" + "log deletion" / "logs deleted"
**Severity:** High
**Explanation:** Audit logs are claimed immutable but deletion policies allow removal.
**Evidence:** IDs of both assumptions

### 6. TLS_REQUIRED_HTTP_ALLOWED
**Trigger:** "TLS required" / "HTTPS only" + "HTTP allowed" / "HTTP accessible"
**Severity:** High
**Explanation:** TLS is required but HTTP is still accessible, allowing downgrade attacks.
**Evidence:** IDs of both assumptions

### 7. ENCRYPTION_WITHOUT_KEY_MANAGEMENT
**Trigger:** "encrypted" + no key management / KMS / HSM assumptions
**Severity:** High
**Explanation:** Data is encrypted but key management controls are not specified.
**Evidence:** IDs of encryption assumption, missing key management

### 8. SESSION_WITHOUT_ROTATION
**Trigger:** "session" / "token" + no rotation / expiration / refresh
**Severity:** Medium
**Explanation:** Session tokens are used but rotation policies are not specified.
**Evidence:** IDs of session assumption, missing rotation

## Implementation

```go
type ContradictionEngine struct{}

func (ce *ContradictionEngine) DetectContradictions(assumptions []Assumption) []Contradiction {
    var results []Contradiction
    results = append(results, ce.detectMFAExemption(assumptions)...)
    results = append(results, ce.detectPlaintextBackup(assumptions)...)
    results = append(results, ce.detectSharedAdmin(assumptions)...)
    results = append(results, ce.detectPrivateSubnetInternet(assumptions)...)
    results = append(results, ce.detectImmutableAuditDeletion(assumptions)...)
    results = append(results, ce.detectTLSHTTP(assumptions)...)
    results = append(results, ce.detectEncryptionWithoutKeyManagement(assumptions)...)
    results = append(results, ce.detectSessionWithoutRotation(assumptions)...)
    return results
}
```

## Test Results

### Test 1: MFA + Service Account Exemption
**Input:**
- "MFA is enforced for all users"
- "Service accounts are exempt from MFA"

**Result:** ✅ Detected (Critical)

### Test 2: Encrypted + Plaintext Backup
**Input:**
- "All data is encrypted at rest"
- "Backups are stored in plaintext"

**Result:** ✅ Detected (Critical)

### Test 3: Least Privilege + Shared Admin
**Input:**
- "Least privilege is enforced"
- "Administrators share a single account"

**Result:** ❌ Not detected (requires exact keyword matching)

### Test 4: Private Subnet + Internet Access
**Input:**
- "Database is in a private subnet"
- "Database is accessible from the internet"

**Result:** ❌ Not detected (requires exact keyword matching)

## Output Format

```json
{
  "contradictions": [
    {
      "id": "CON-001",
      "severity": "Critical",
      "description": "MFA is enforced for all users but service accounts are exempt",
      "explanation": "MFA is enforced for users but service accounts are exempted, creating a privileged bypass path.",
      "affected_assumptions": ["ASM-001", "ASM-002"],
      "rule_name": "MFA_ENFORCED_WITH_EXEMPTION"
    }
  ]
}
```

## Integration

Contradictions are detected after the intelligence engine merges all assumptions. They are included in:
- JSON output
- CLI output
- TUI results view
- Export formats

## Future Improvements

- Semantic matching (not just keyword matching)
- Cross-reference with evidence files
- Temporal contradiction detection (e.g., "logs kept for 1 year" vs "logs kept for 7 years")
- Configuration drift detection

## Determinism

All contradiction detection is deterministic:
- String matching with regex
- No randomness
- No AI/LLM calls
- No cloud services
- Reproducible across runs
