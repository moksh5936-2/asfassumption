# Release Certification — End-to-End Product Test

## Commands Run

```bash
asf analyze testdata/asftest.yaml --json  # YAML benchmark
asf analyze testdata/asftest.md --json      # Markdown benchmark
```

## Results

### YAML (asftest.yaml)

| Metric | Value |
|--------|-------|
| Assumptions | 48 |
| Critical | 3 |
| High | 4 |
| Medium | 9 |
| Low | 32 |
| Partially Verified | 5 |
| Controls (gaps) | 43 |

**STRIDE Coverage:** All 6 categories present
- Spoofing, Tampering, Repudiation, Information Disclosure, Denial of Service, Elevation of Privilege

**Compliance:** HIPAA, SOC2, ISO27001 referenced in engine output

### Markdown (asftest.md)

| Metric | Value |
|--------|-------|
| Assumptions | 38 |
| Critical | 0 |
| High | 2 |
| Medium | 4 |
| Low | 32 |
| Partially Verified | 4 |
| Controls (gaps) | 34 |

**STRIDE Coverage:** All 6 categories present

**Compliance:** HIPAA, SOC2, ISO27001 referenced

## Partially Verified Assumptions (YAML)

1. `- What if MFA is not enforced on VPN?` — matched MFA control
2. `- AdminConsole requires MFA for all administrative access.` — matched MFA control
3. `MFA is enforced for all Auth0 user authentication` — matched MFA control
4. `Network segmentation isolates the PHI database in a private subnet` — matched Network Segmentation control
5. `Vendor risk assessments are conducted for ThirdPartyAnalytics` — matched Vendor Risk Assessment control

## High-Risk Themes Detected

| Theme | Status |
|-------|--------|
| Key_Management | ✅ (KMS references detected) |
| PHI_Access_Control | ✅ (PHI + access references detected) |
| Authentication | ✅ (MFA + auth references detected) |
| Third_Party_Dependencies | ✅ (Third-party references detected) |

## Verification Notes

- Binary rebuilt at 2026-06-12 18:12 after code changes
- `asf-tui analyze testdata/asftest.yaml --json` produces 48 assumptions (verified live)
- `asf-tui analyze testdata/asftest.md --json` produces 38 assumptions (verified live)

## Verdict

✅ **PASS** — E2E product analysis works correctly. Both YAML and Markdown produce structured output with assumptions, risk distribution, STRIDE coverage, compliance references, and partial verification from security controls.
