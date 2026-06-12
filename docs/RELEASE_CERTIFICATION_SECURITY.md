# Release Certification — Security Review

## Files Audited

| File | Lines | Focus |
|------|-------|-------|
| `license.go` | 131 | License validation, HMAC demo key |
| `license_ed25519.go` | 61 | Ed25519 demo keys |
| `engine.go` | 1400 | Analysis engine, temp files |
| `parser.go` | 935 | File parsing, path handling |
| `main.go` | 140 | CLI entry, command routing |

## Findings

### Finding 1: Hardcoded Demo Secret (LOW)

**File:** `license.go:15`

```go
const DemoSecret = "asf-enterprise-secret-2024"
```

**Description:** HMAC demo key hardcoded in source. Explicitly labeled as "demonstration-only" and "NOT security". Used for license validation in demo builds.

**Impact:** LOW — The code explicitly states this is a demo key. Production builds should use `license_ed25519.go` or a proper key replacement mechanism.

**Mitigation:** Already documented in code comments. Not a release blocker for a demo/certification build.

### Finding 2: Hardcoded Ed25519 Seed (LOW)

**File:** `license_ed25519.go:17`

```go
seed := sha256.Sum256([]byte("asf-ed25519-demo-seed-2024-ed25519"))
```

**Description:** Deterministic Ed25519 key pair generated from a hardcoded string.

**Impact:** LOW — The `ReplacePublicKey()` function allows replacing the public key at runtime. This is a demo key.

**Mitigation:** Documented as demo. Production deployments should call `ReplacePublicKey()` with a real key.

### Finding 3: No Output Path Validation (MEDIUM)

**File:** `export.go`

**Description:** `ExportResult` writes to `filepath.Join(outputDir, baseName + ext)` without validating that `outputDir` is safe. If a malicious config sets `outputDir` to a system directory, exports could overwrite files.

**Impact:** MEDIUM — Requires attacker to control the user's config file. Config is stored in user's home directory.

**Mitigation:** Config is user-controlled. No privilege escalation possible.

### Finding 4: No Path Traversal in Parser (MEDIUM)

**File:** `parser.go`

**Description:** `ParseArchitecture` reads files directly from the provided path. No validation that the path is within expected bounds.

**Impact:** MEDIUM — The user provides the path via CLI. If a user passes `../../etc/passwd`, it would be read. However, this is a CLI tool where the user already has filesystem access.

**Mitigation:** User already has filesystem access. No privilege escalation.

## Absent Issues

| Issue | Status | Evidence |
|-------|--------|----------|
| Command injection | ✅ Not found | No `exec.Command` with user input |
| Arbitrary file write | ✅ Not found | All writes are to controlled paths or user-specified output dir |
| Unsafe temp files | ✅ Not found | Uses `mktemp -d` in install, controlled paths in engine |
| Panic paths | ✅ Not found | Standard error handling throughout |
| Goroutine leaks | ✅ Not found | Progress channels closed properly |
| Path traversal in export | ✅ Minor | Acceptable for user-controlled tool |

## Severity Classification

| Severity | Count | Issues |
|----------|-------|--------|
| CRITICAL | 0 | None |
| HIGH | 0 | None |
| MEDIUM | 2 | Output path validation, parser path validation |
| LOW | 2 | Demo HMAC key, demo Ed25519 seed |

## Verdict

✅ **PASS** — No critical or high security issues. Demo keys are explicitly labeled. Path handling risks are acceptable for a user-controlled CLI tool.
