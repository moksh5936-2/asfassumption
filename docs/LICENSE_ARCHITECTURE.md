# ASF License Architecture

> Version: 1.0.0 | June 2026

## Overview

ASF uses a deterministic HMAC-SHA256 license key system for enterprise licensing. The system validates license keys entirely offline with no phone-home mechanism.

## License Key Format

```
ASF-XXXX-XXXX-XXXX-XXXX-SSSSSSSS
│    │    │    │    │    │
│    │    │    │    │    └── 8-char HMAC-SHA256 signature (hex)
│    │    │    │    └─────── 4-char hex segment
│    │    │    └──────────── 4-char hex segment
│    │    └───────────────── 4-char hex segment
│    └────────────────────── 4-char hex segment
└─────────────────────────── "ASF" prefix
```

**Total:** 27 characters (ASF-XXXX-XXXX-XXXX-XXXX-SSSSSSSS)

## Validation Flow

```
User enters license key
        │
        ▼
Parse key format (regex validation)
        │
        ▼
Extract payload (ASF-XXXX-XXXX-XXXX-XXXX)
        │
        ▼
Extract signature (SSSSSSSS)
        │
        ▼
Compute HMAC-SHA256(payload, secret_key)
        │
        ▼
Compare computed signature with provided signature
        │
        ▼
┌───────┴───────┐
│ Match?         │
├───┬───────────┤
│   Yes          │   No
│   ▼            │   ▼
│ Valid License  │   Invalid License
│   │            │   │
│   ▼            │   ▼
│ ~/.asf/       │   Error message
│ license.key   │   displayed
└───────────────┘
```

## Code Location

- **License validation:** `asf-tui/license.go` — `validateLicenseKey()` function
- **CLI check:** `main.go` — `--license` flag
- **License file:** `~/.asf/license.key`

## Security Properties

| Property | Implementation |
|----------|---------------|
| Algorithm | HMAC-SHA256 |
| Secret key | Hardcoded in binary (deterministic) |
| Verification | Local only — no network calls |
| Side-channel | None — single comparison |
| Replay attack | Not applicable (no online validation) |
| Brute force | 2^32 space (8 hex chars) for signature |

## Limitations

### Security Weaknesses

1. **Hardcoded secret key** — The HMAC secret is compiled into the binary. Anyone with the binary can extract it via `strings` or reverse engineering.
2. **No expiration** — Licenses never expire. No time-based validation.
3. **No revocation** — Once a license key is generated, it cannot be revoked without a software update.
4. **No rate limiting** — The TUI allows unlimited validation attempts.
5. **No machine binding** — A license key works on any machine. No hardware fingerprinting.

### Feature Gating

The current license system validates keys but does not gate any features. The `--license` flag only shows the license status — it does not enable/disable functionality.

## Enterprise Activation

```bash
# Activate a license
echo 'ASF-XXXX-XXXX-XXXX-XXXX-SSSSSSSS' > ~/.asf/license.key

# Check license
asf --license

# Deactivate
rm ~/.asf/license.key
```

## Attack Surface

| Attack Vector | Risk | Mitigation |
|---------------|------|------------|
| Secret key extraction from binary | 🔴 High | Obfuscation or external key server needed |
| License key sharing | 🟡 Medium | Machine binding required |
| Static analysis bypass | 🟡 Medium | Runtime decryption could help |
| Timing attack | 🟢 Low | Constant-time HMAC comparison |

## Recommendations

1. **Implement feature gating** — Use license validation to unlock enterprise features
2. **Add machine binding** — Hash machine identifiers into the license payload
3. **Add expiration dates** — Embed validity period in the license payload
4. **Obfuscate the secret key** — Split across binary sections, XOR-encode, or use runtime derivation
5. **Add rate limiting** — Limit validation attempts in the TUI
6. **Consider online validation** — Optional phone-home for enterprise customers
