# ASF Security Review

> Version: 1.0.0 | June 2026 | Classification: Internal

## Executive Summary

This document presents a security review of the ASF v1.0.0 codebase. ASF is a security tool that processes architecture documents locally. The review covers the application's own security posture, not the security assumptions it discovers.

**Overall Rating: MODERATE RISK** — The application has a small attack surface and no external-facing network services, but contains several areas of concern including an unverified Python dependency, lack of input sanitization in file parsers, and a hardcoded cryptographic secret.

---

## 1. Threat Model

### Assets

| Asset | Sensitivity | Location |
|-------|-------------|----------|
| Architecture documents | Confidential (customer designs) | Local filesystem |
| Assumption analysis results | Confidential | Local filesystem, TUI display |
| License key | Credential | `~/.asf/license.key` |
| User config | Low | `~/.asf/config.yaml` |
| ASF binary | Code integrity | `/usr/local/bin/asf` |

### Trust Boundaries

```
[Untrusted]                 [Trusted]                [Local Only]
──────────────    ┌─────────────────────┐    ┌──────────────────┐
│ User Input    │──▶│   ASF Process      │──▶│   Ollama API      │
│ (files, TUI)  │   │   (Go Binary)      │   │   localhost:11434  │
└──────────────┘   │   │                  │   └──────────────────┘
                   │   │   Python ASF     │   ┌──────────────────┐
                   │   │   (subprocess)   │──▶│   Filesystem I/O  │
                   │   └──────────────────┘   └──────────────────┘
                   └─────────────────────┘
```

- **Boundary 1:** Untrusted user files enter via file parsers
- **Boundary 2:** Untrusted TUI input (keyboard events)
- **Boundary 3:** Subprocess execution of Python ASF engine
- **Boundary 4:** Local network call to Ollama

---

## 2. Attack Surface

### File Parsing (HIGHEST RISK)

| Parser | Input Type | Risk | Vector |
|--------|-----------|------|--------|
| Draw.io XML | `.drawio` | 🟡 Medium | Malformed XML, entity expansion |
| Mermaid | `.mmd` | 🟢 Low | Regex-based parsing, no code execution |
| YAML | `.yaml` | 🟢 Low | Standard library parser |
| JSON | `.json` | 🟢 Low | Standard library parser |
| SVG | `.svg` | 🟡 Medium | XML parsing of untrusted SVG |
| OCR (Tesseract) | `.png/.jpg` | 🟢 Low | Tesseract subprocess with file path |
| Text/PDF/DOCX | `.txt/.md/.pdf/.docx` | 🟢 Low | Text extraction only |

**Analysis:** Go's `encoding/xml` is generally safe against XXE (no DTD processing by default). However, no explicit security configuration is set. The Mermaid parser uses simple regex patterns — injection risk is minimal.

### Python Subprocess

| Risk | Detail |
|------|--------|
| Command injection | Low — command is fixed string `python3 -m asf.cli.main analyze --json` |
| Argument injection | Low — input is piped via stdin, not command-line arguments |
| Python dependency trust | Medium — the Python ASF engine runs arbitrary Python code in a subprocess |
| Timeout failure | Low — 10-second timeout prevents hanging |

### TUI Input

| Risk | Detail |
|------|--------|
| Key injection | Low — Bubble Tea handles raw terminal input safely |
| Buffer overflow | Low — Go is memory-safe |
| Terminal escape sequences | Low — Bubble Tea parses terminal responses safely |

### Network

| Risk | Detail |
|------|--------|
| Ollama API call | Low — HTTP POST to localhost only, fixed URL |
| Data exfiltration | Low — no external network calls exist |
| Man-in-the-middle | Low — loopback interface only |

---

## 3. Dependency Vulnerabilities

| Dependency | Known CVEs | Risk |
|------------|-----------|------|
| bubbletea v1.3.10 | None known | 🟢 Low |
| lipgloss v1.1.0 | None known | 🟢 Low |
| go-pdf/fpdf v0.9.0 | None known | 🟢 Low |
| yaml.v3 v3.0.1 | None known (CVE-2022-3064 fixed in v3.0.1) | 🟢 Low |
| Go stdlib (1.24) | None known | 🟢 Low |

**Note:** No automated vulnerability scanning has been performed. This is a manual assessment based on public CVE databases.

---

## 4. Supply Chain Security

| Practice | Status | Recommendation |
|----------|--------|----------------|
| `go.sum` verification | ✅ In place | Ensure `go mod verify` runs in CI |
| Dependency pinning | ✅ In place | `go.mod` pins all direct dependencies |
| Vendor directory | ❌ Not used | Consider vendoring for air-gapped builds |
| SBOM generation | ❌ Not done | Use `go version -m` or `syft` |
| Dependency scanning | ❌ Not configured | Add Dependabot or Renovate |
| Binary signing | ❌ Not done | Add macOS codesign + notarization |

---

## 5. Security Controls

### Current Controls

| Control | Status | Detail |
|---------|--------|--------|
| No external network | ✅ | All analysis local, no telemetry |
| Deterministic analysis | ✅ | No AI randomness in core engine |
| Local-only AI | ✅ | Ollama runs on localhost |
| Config file permissions | ✅ | Stored in `~/.asf/` with user-only access |
| License validation | ✅ | Local HMAC, no phone-home |
| Gitignore for secrets | ✅ | `license.key` patterns excluded |

### Missing Controls

| Control | Priority | Detail |
|---------|----------|--------|
| Input size limits | 🔴 High | No maximum file size for parsers |
| Input validation | 🟡 Medium | No schema validation for YAML/JSON inputs |
| Python subprocess sandboxing | 🟡 Medium | Python runs with same permissions as ASF |
| Error message sanitization | 🟡 Medium | File paths may appear in error messages |
| Audit logging | 🟡 Medium | No log of which files were analyzed |
| Rate limiting | 🟢 Low | Not a network service |
| Crash reporting | 🟢 Low | No telemetry by design |

---

## 6. Vulnerability Summary

| ID | Vulnerability | Severity | Status |
|----|--------------|----------|--------|
| SEC-001 | Python subprocess not sandboxed | Medium | Accepted (design constraint) |
| SEC-002 | No input size limits on file parsers | Medium | Unaddressed |
| SEC-003 | Hardcoded HMAC secret in binary | Medium | Accepted (offline design) |
| SEC-004 | No macOS code signing | Medium | Unaddressed |
| SEC-005 | curl-pipe-bash installer | Low | Accepted (common pattern) |
| SEC-006 | No binary signature verification | Low | Unaddressed |
| SEC-007 | Error messages may leak file paths | Low | Unaddressed |
| SEC-008 | No audit trail of analysis runs | Low | Unaddressed |

---

## 7. Recommendations

### High Priority
1. Add maximum file size limits (e.g., 10MB) before parsing
2. Validate YAML/JSON schema before processing
3. Add macOS code signing and notarization

### Medium Priority
4. Run Python subprocess in restricted mode (e.g., `-I` flag for isolated mode)
5. Add basic audit logging (file analyzed, timestamp, result count)
6. Sanitize error messages to avoid leaking file paths
7. Add dependency vulnerability scanning (Dependabot)

### Low Priority
8. Consider binary vendoring for offline builds
9. Add SBOM generation to release workflow
10. Add `ed25519` signatures to release artifacts
