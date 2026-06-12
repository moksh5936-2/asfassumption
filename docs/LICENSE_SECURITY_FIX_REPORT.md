# License Security Fix Report (B5)

## Problem Summary

The HMAC secret key used for license signing and validation was hardcoded as a string literal in `asf-tui/license.go`:

```go
mac := hmac.New(sha256.New, []byte("asf-enterprise-secret-2024"))
```

This meant anyone with access to the binary could trivially extract the secret using `strings`:

```bash
strings asf-tui | grep secret
```

With the HMAC key known, forging a valid license for any `ASF-XXXX-XXXX-XXXX` payload was a matter of computing an 8-character hex digest — no cracking required.

## Why Full Ed25519 Migration Was Not Done

1. **No existing key infrastructure** — There is no offline private key, KMS, or build-time signing pipeline in place.
2. **Existing keys would break** — All currently deployed license keys are HMAC-signed with the same secret. Switching to asymmetric signing would invalidate every key in the wild.
3. **Engineering cost exceeds benefit** — At this development/demo phase, introducing Ed25519 or ECDSA would add a dependency, complicate key distribution, and provide no practical benefit until a commercial launch is planned.

## What Was Done Instead

| Change | File |
|---|---|
| Secret moved to a named const `DemoSecret` with a clear comment | `asf-tui/license.go:13` |
| `ValidateLicense` preceded by a doc block stating demo-only nature | `asf-tui/license.go:60` |
| All references to `"asf-enterprise-secret-2024"` replaced with `DemoSecret` | `asf-tui/license.go:70,101` |
| README claim of "enterprise licensing" qualified to "demo licensing (development use only — not cryptographically secure)" | `README.md:328` |

## Security Posture

**Obfuscation only. Not suitable for production license enforcement.**

The HMAC key is still embedded in the binary (it has to be for symmetric validation to work). The `strings` attack path still exists. The fix is documentation-level — it makes the limitation explicit but does not eliminate it.

## Future Recommendation

Replace with **Ed25519** before any commercial launch:

1. Generate an offline Ed25519 keypair.
2. Embed the **public key** only in the binary.
3. Sign license payloads offline with the private key.
4. Verify signatures at runtime with the embedded public key.

This ensures that even with full binary access, an attacker cannot forge licenses — they would need the private key, which never touches the build or runtime environment.
